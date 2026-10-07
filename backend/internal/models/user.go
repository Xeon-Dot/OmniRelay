package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsAdmin      bool      `json:"is_admin"`
	TOTPEnabled  bool      `json:"totp_enabled"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token          string `json:"token"`
	User           User   `json:"user"`
	Requires2FA    bool   `json:"requires_2fa"`
	TwoFactorToken string `json:"two_factor_token,omitempty"`
}

type TwoFactorVerifyRequest struct {
	TwoFactorToken string `json:"two_factor_token" binding:"required"`
	Code           string `json:"code" binding:"required"`
}

type TwoFactorCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

type TwoFactorDisableRequest struct {
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type SetRoleRequest struct {
	IsAdmin bool `json:"is_admin"`
}

type ResetPasswordRequest struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

type SetUserProvidersRequest struct {
	ProviderIDs []int64 `json:"provider_ids"`
}
