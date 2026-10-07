package handlers

import (
	"net/http"
	"omnirelay/internal/apiresponse"
	"omnirelay/internal/models"
	"omnirelay/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}

		user, err := svc.Register(req)
		if err != nil {
			apiresponse.AbortAdminConflict(c, err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": user})
	}
}

func Login(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}

		resp, err := svc.Login(req)
		if err != nil {
			apiresponse.AbortAdminError(c, http.StatusUnauthorized, err.Error(), "unauthorized")
			return
		}

		resetLoginRateLimit(c)
		c.JSON(http.StatusOK, resp)
	}
}

func TwoFactorStatus(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		enabled, err := svc.TOTPStatus(c.GetInt64("user_id"))
		if err != nil {
			apiresponse.AbortAdminInternal(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"enabled": enabled})
	}
}

func TwoFactorSetup(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, url, err := svc.SetupTOTP(c.GetInt64("user_id"))
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"secret": secret, "otpauth_url": url})
	}
}

func TwoFactorEnable(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.TwoFactorCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		codes, err := svc.EnableTOTP(c.GetInt64("user_id"), req.Code)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"recovery_codes": codes})
	}
}

func TwoFactorDisable(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.TwoFactorDisableRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		if err := svc.DisableTOTP(c.GetInt64("user_id"), req.Password, req.Code); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "two-factor authentication disabled"})
	}
}

func TwoFactorVerify(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.TwoFactorVerifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		resp, err := svc.VerifyTwoFactor(req.TwoFactorToken, req.Code)
		if err != nil {
			apiresponse.AbortAdminError(c, http.StatusUnauthorized, err.Error(), "unauthorized")
			return
		}
		resetLoginRateLimit(c)
		c.JSON(http.StatusOK, resp)
	}
}

// WSToken mints a short-lived token for the /admin/ws upgrade. Browsers
// cannot set headers on a WebSocket handshake, so the token must travel in
// the URL; a scoped 60s token limits the damage when it leaks into access
// logs.
func WSToken(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := svc.IssueWSToken(c.GetInt64("user_id"), c.GetString("username"))
		if err != nil {
			apiresponse.AbortAdminInternal(c, "failed to issue websocket token")
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func ListUsers(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := svc.ListUsers()
		if err != nil {
			apiresponse.AbortAdminInternal(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func DeleteUser(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, "invalid user id")
			return
		}
		requesterID := c.GetInt64("user_id")
		if err := svc.DeleteUser(id, requesterID); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}

func SetUserRole(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, "invalid user id")
			return
		}
		var req models.SetRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		if err := svc.SetRole(id, req.IsAdmin); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "role updated"})
	}
}

func GenerateResetCode(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, "invalid user id")
			return
		}
		code, err := svc.GenerateResetCode(id)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": code})
	}
}

func ResetPassword(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		if err := svc.ResetPasswordWithCode(req.Code, req.NewPassword); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
	}
}

func GetUserProviders(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, "invalid user id")
			return
		}
		ids, err := svc.GetUserProviders(id)
		if err != nil {
			apiresponse.AbortAdminInternal(c, err.Error())
			return
		}
		if ids == nil {
			ids = []int64{}
		}
		c.JSON(http.StatusOK, gin.H{"provider_ids": ids})
	}
}

func SetUserProviders(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			apiresponse.AbortAdminBadRequest(c, "invalid user id")
			return
		}
		var req models.SetUserProvidersRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.AbortAdminBadRequest(c, err.Error())
			return
		}
		if err := svc.SetUserProviders(id, req.ProviderIDs); err != nil {
			apiresponse.AbortAdminInternal(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "providers updated"})
	}
}
