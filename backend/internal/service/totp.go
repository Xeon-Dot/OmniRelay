package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"omnirelay/internal/crypto"
	"omnirelay/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// issueTwoFactorToken returns a short-lived JWT that only proves the user
// passed password auth; it must not be accepted as a session token (scope:"2fa").
func (s *AuthService) issueTwoFactorToken(userID int64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"scope":    "2fa",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// validTwoFactorToken checks a pending-token from /auth/login.
func (s *AuthService) validTwoFactorToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired two-factor token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["scope"] != "2fa" {
		return 0, errors.New("invalid two-factor token")
	}
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid two-factor token claims")
	}
	return int64(uid), nil
}

// IssueSessionToken mints the regular 24h session JWT for a user.
func (s *AuthService) issueSessionToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// VerifyTwoFactor completes login: accepts a TOTP code or a single-use recovery code.
func (s *AuthService) VerifyTwoFactor(twoFactorToken, code string) (*models.LoginResponse, error) {
	userID, err := s.validTwoFactorToken(twoFactorToken)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = s.db.QueryRow(
		"SELECT id, username, email, is_admin, COALESCE(totp_enabled, 0), created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.TOTPEnabled, &user.CreatedAt)
	if err != nil || !user.TOTPEnabled {
		return nil, errors.New("two-factor authentication not enabled")
	}

	if ok := s.checkUserTOTP(userID, code); ok {
		return s.finishLogin(user)
	}
	if s.consumeRecoveryCode(userID, code) {
		return s.finishLogin(user)
	}
	return nil, errors.New("invalid two-factor code")
}

func (s *AuthService) finishLogin(user models.User) (*models.LoginResponse, error) {
	tokenString, err := s.issueSessionToken(user)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{Token: tokenString, User: user}, nil
}

func (s *AuthService) checkUserTOTP(userID int64, code string) bool {
	var enc string
	if err := s.db.QueryRow("SELECT totp_secret_encrypted FROM users WHERE id = ?", userID).Scan(&enc); err != nil || enc == "" {
		return false
	}
	secret, err := crypto.Decrypt(enc, s.encryptKey)
	if err != nil {
		return false
	}
	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{Period: 30, Skew: 1, Digits: otp.DigitsSix, Algorithm: otp.AlgorithmSHA1})
	return err == nil && valid
}

// SetupTOTP generates a new TOTP secret for the user and stores it (encrypted)
// in a disabled state; it becomes active on EnableTOTP.
func (s *AuthService) SetupTOTP(userID int64) (secret, otpauthURL string, err error) {
	var email string
	if err := s.db.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email); err != nil {
		return "", "", errors.New("user not found")
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OmniRelay",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}
	enc, err := crypto.Encrypt(key.Secret(), s.encryptKey)
	if err != nil {
		return "", "", fmt.Errorf("encrypt secret: %w", err)
	}
	if _, err := s.db.Exec("UPDATE users SET totp_secret_encrypted = ?, totp_enabled = 0 WHERE id = ?", enc, userID); err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// EnableTOTP activates 2FA after verifying the first code, and issues recovery codes.
func (s *AuthService) EnableTOTP(userID int64, code string) ([]string, error) {
	if !s.checkUserTOTP(userID, code) {
		return nil, errors.New("invalid code")
	}
	if _, err := s.db.Exec("UPDATE users SET totp_enabled = 1 WHERE id = ?", userID); err != nil {
		return nil, err
	}
	return s.generateRecoveryCodes(userID)
}

// DisableTOTP requires the account password and a valid TOTP code (or recovery code).
func (s *AuthService) DisableTOTP(userID int64, password, code string) error {
	var hash string
	if err := s.db.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&hash); err != nil {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("invalid password")
	}
	if !s.checkUserTOTP(userID, code) && !s.consumeRecoveryCode(userID, code) {
		return errors.New("invalid two-factor code")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("UPDATE users SET totp_enabled = 0, totp_secret_encrypted = '' WHERE id = ?", userID); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM recovery_codes WHERE user_id = ?", userID); err != nil {
		return err
	}
	return tx.Commit()
}

// TOTPStatus reports whether the user currently has 2FA enabled.
func (s *AuthService) TOTPStatus(userID int64) (bool, error) {
	var enabled bool
	err := s.db.QueryRow("SELECT COALESCE(totp_enabled, 0) FROM users WHERE id = ?", userID).Scan(&enabled)
	return enabled, err
}

func (s *AuthService) generateRecoveryCodes(userID int64) ([]string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("DELETE FROM recovery_codes WHERE user_id = ?", userID); err != nil {
		return nil, err
	}
	codes := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		plain := hex.EncodeToString(b)
		plain = plain[:4] + "-" + plain[4:]
		sum := sha256.Sum256([]byte(plain))
		if _, err := tx.Exec("INSERT INTO recovery_codes (user_id, code_hash) VALUES (?, ?)", userID, hex.EncodeToString(sum[:])); err != nil {
			return nil, err
		}
		codes = append(codes, plain)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return codes, nil
}

func (s *AuthService) consumeRecoveryCode(userID int64, code string) bool {
	sum := sha256.Sum256([]byte(code))
	res, err := s.db.Exec("UPDATE recovery_codes SET used = 1 WHERE user_id = ? AND code_hash = ? AND used = 0", userID, hex.EncodeToString(sum[:]))
	if err != nil {
		return false
	}
	n, err := res.RowsAffected()
	return err == nil && n > 0
}
