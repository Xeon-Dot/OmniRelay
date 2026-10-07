package service

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"

	"omnirelay/internal/models"
)

const testEncryptKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func newTestTOTPService(t *testing.T) (*AuthService, int64) {
	t.Helper()
	svc := newTestAuthService(t)
	svc.SetJWTSecret("test-secret-for-jwt-signing")
	svc.SetEncryptKey(testEncryptKey)
	u, err := svc.Register(models.RegisterRequest{Username: "alice", Email: "alice@example.com", Password: "password1"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	return svc, u.ID
}

func TestTOTPEnableAndVerifyLoginFlow(t *testing.T) {
	svc, uid := newTestTOTPService(t)

	secret, _, err := svc.SetupTOTP(uid)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatalf("generate code: %v", err)
	}
	recovery, err := svc.EnableTOTP(uid, code)
	if err != nil {
		t.Fatalf("enable: %v", err)
	}
	if len(recovery) != 10 {
		t.Fatalf("expected 10 recovery codes, got %d", len(recovery))
	}

	login, err := svc.Login(models.LoginRequest{Email: "alice@example.com", Password: "password1"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if !login.Requires2FA || login.TwoFactorToken == "" {
		t.Fatal("expected pending 2FA token")
	}

	code2, _ := totp.GenerateCode(secret, time.Now())
	resp, err := svc.VerifyTwoFactor(login.TwoFactorToken, code2)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected session token")
	}
}

func TestTOTPRecoveryCodeSingleUse(t *testing.T) {
	svc, uid := newTestTOTPService(t)
	secret, _, _ := svc.SetupTOTP(uid)
	code, _ := totp.GenerateCode(secret, time.Now())
	recovery, _ := svc.EnableTOTP(uid, code)

	login, _ := svc.Login(models.LoginRequest{Email: "alice@example.com", Password: "password1"})
	if _, err := svc.VerifyTwoFactor(login.TwoFactorToken, recovery[0]); err != nil {
		t.Fatalf("recovery code should work once: %v", err)
	}
	login2, _ := svc.Login(models.LoginRequest{Email: "alice@example.com", Password: "password1"})
	if _, err := svc.VerifyTwoFactor(login2.TwoFactorToken, recovery[0]); err == nil {
		t.Fatal("recovery code should not be reusable")
	}
}

func TestTOTPDisableRequiresPassword(t *testing.T) {
	svc, uid := newTestTOTPService(t)
	secret, _, _ := svc.SetupTOTP(uid)
	code, _ := totp.GenerateCode(secret, time.Now())
	svc.EnableTOTP(uid, code)

	if err := svc.DisableTOTP(uid, "wrong", code); err == nil {
		t.Fatal("expected invalid password error")
	}
	code2, _ := totp.GenerateCode(secret, time.Now())
	if err := svc.DisableTOTP(uid, "password1", code2); err != nil {
		t.Fatalf("disable: %v", err)
	}
	enabled, _ := svc.TOTPStatus(uid)
	if enabled {
		t.Fatal("2FA should be disabled")
	}
}
