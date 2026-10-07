package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"omnirelay/internal/config"
	"omnirelay/internal/database"
	"omnirelay/internal/handlers"
	"omnirelay/internal/hub"
	"omnirelay/internal/middleware"
	"omnirelay/internal/proxy"
	"omnirelay/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const maxBodySize = 32 << 20 // 32MB

func main() {
	cfg := config.Load()

	db, err := database.Init(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	defer db.Close()

	authService := service.NewAuthService(db)
	authService.SetJWTSecret(cfg.JWTSecret)
	authService.SetEncryptKey(cfg.EncryptKey)
	providerService := service.NewProviderService(db, cfg)
	modelService := service.NewModelService(db)
	modelService.SetPricingCatalog(service.NewModelsDevCatalog())
	apiKeyService := service.NewAPIKeyService(db)
	usageService := service.NewUsageService(db)
	usageService.SetAPIKeyService(apiKeyService)
	performanceService := service.NewPerformanceService(db)

	h := hub.New()
	handlers.SetHub(h)
	proxyEngine := proxy.NewEngine(providerService, modelService, usageService, authService, h)

	r := gin.Default()

	// CORS: explicit origins from config, or same-origin-only in release
	// mode (the SPA is served by Caddy from the same host). In dev the
	// Vite dev server runs on a different origin, so default to it.
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	if cfg.CORSOrigins != "" {
		allowedOrigins = splitAndTrim(cfg.CORSOrigins, ",")
	} else if os.Getenv("GIN_MODE") == "release" {
		// Same-origin-only in the default production install.
		allowedOrigins = nil
	}
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key"},
		AllowCredentials: true,
	}
	if len(allowedOrigins) == 0 {
		// No origins configured: accept only same-origin requests so a
		// default install does not keep trusting localhost origins.
		corsConfig.AllowOriginWithContextFunc = func(c *gin.Context, origin string) bool {
			return isSameOrigin(origin, c.Request.Host)
		}
	} else {
		corsConfig.AllowOrigins = allowedOrigins
	}
	r.Use(cors.New(corsConfig))

	// Health check endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	admin := r.Group("/admin")
	{
		admin.POST("/auth/register", handlers.LoginRateLimit(), handlers.Register(authService))
		admin.POST("/auth/login", handlers.LoginRateLimit(), handlers.Login(authService))
		admin.POST("/auth/reset-password", handlers.LoginRateLimit(), handlers.ResetPassword(authService))
		admin.POST("/auth/2fa/verify", handlers.LoginRateLimit(), handlers.TwoFactorVerify(authService))
		admin.GET("/ws", handlers.WebSocketUpgrader(cfg.JWTSecret, allowedOrigins))

		adminAuth := admin.Group("")
		adminAuth.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			adminAuth.GET("/auth/2fa/status", handlers.TwoFactorStatus(authService))
			adminAuth.POST("/auth/2fa/setup", handlers.LoginRateLimit(), handlers.TwoFactorSetup(authService))
			adminAuth.POST("/auth/2fa/enable", handlers.LoginRateLimit(), handlers.TwoFactorEnable(authService))
			adminAuth.POST("/auth/2fa/disable", handlers.LoginRateLimit(), handlers.TwoFactorDisable(authService))
			adminAuth.POST("/auth/ws-token", handlers.WSToken(authService))

			adminAuth.GET("/providers", handlers.ListProviders(providerService, authService))
			adminAuth.GET("/models", handlers.ListModels(modelService, authService))
			adminAuth.GET("/models/source-list", handlers.ListSourceModels(modelService))

			adminAuth.GET("/api-keys", handlers.ListAPIKeys(apiKeyService))
			adminAuth.POST("/api-keys", handlers.CreateAPIKey(apiKeyService))
			adminAuth.DELETE("/api-keys/:id", handlers.DeleteAPIKey(apiKeyService))

			adminAuth.GET("/usage", handlers.ListUsage(usageService))
			adminAuth.GET("/stats", handlers.GetStats(usageService, apiKeyService, modelService))
			adminAuth.GET("/performance", handlers.GetPerformance(performanceService))

			adminOnly := adminAuth.Group("")
			adminOnly.Use(middleware.RequireAdmin())
			{
				adminOnly.POST("/providers", handlers.CreateProvider(providerService))
				adminOnly.PUT("/providers/:id", handlers.UpdateProvider(providerService))
				adminOnly.DELETE("/providers/:id", handlers.DeleteProvider(providerService))
				adminOnly.POST("/providers/:id/sync", handlers.SyncProviderModels(providerService, modelService))
				adminOnly.POST("/providers/:id/test", handlers.TestProvider(providerService, proxyEngine))
				adminOnly.POST("/providers/:id/keys", handlers.AddProviderKey(providerService))
				adminOnly.PATCH("/providers/:id/keys/:kid", handlers.SetProviderKeyActive(providerService))
				adminOnly.DELETE("/providers/:id/keys/:kid", handlers.DeleteProviderKey(providerService))

				adminOnly.POST("/models", handlers.CreateModel(modelService))
				adminOnly.PUT("/models/:id", handlers.UpdateModel(modelService))
				adminOnly.DELETE("/models/:id", handlers.DeleteModel(modelService))

				adminOnly.GET("/users", handlers.ListUsers(authService))
				adminOnly.DELETE("/users/:id", handlers.DeleteUser(authService))
				adminOnly.PUT("/users/:id/role", handlers.SetUserRole(authService))
				adminOnly.POST("/users/:id/reset-password", handlers.GenerateResetCode(authService))
				adminOnly.GET("/users/:id/providers", handlers.GetUserProviders(authService))
				adminOnly.PUT("/users/:id/providers", handlers.SetUserProviders(authService))
			}
		}
	}

	v1 := r.Group("/v1")
	v1.Use(middleware.APIKeyAuth(apiKeyService), bodySizeLimit())
	{
		v1.POST("/chat/completions", proxyEngine.HandleChatCompletions)
		v1.POST("/responses", proxyEngine.HandleResponses)
		v1.GET("/models", proxyEngine.HandleListModels)
		v1.GET("/models/*model", proxyEngine.HandleGetModel)
		v1.POST("/messages", proxyEngine.HandleMessages)
	}

	// Path-based routing: /:provider_key/v1/*endpoint
	pbr := r.Group("/")
	pbr.Use(middleware.APIKeyAuth(apiKeyService), bodySizeLimit())
	pbr.Any("/:provider_key/v1/*endpoint", proxyEngine.HandlePathRouted)
	pbr.Any("/:provider_key/v1beta/*endpoint", proxyEngine.HandlePathRouted)
	pbr.Any("/:provider_key/api/*endpoint", proxyEngine.HandlePathRouted)

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: r,
	}

	go func() {
		log.Printf("OmniRelay starting on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("received signal %v, shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}

// bodySizeLimit returns a middleware that enforces a maximum request body size.
func bodySizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		c.Next()
	}
}

// splitAndTrim splits a string by a separator and trims whitespace from each part.
func splitAndTrim(s, sep string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// isSameOrigin reports whether the Origin header matches the request Host
// (scheme aside, since browsers only send http/https origins here).
func isSameOrigin(origin, host string) bool {
	if origin == "" || host == "" {
		return false
	}
	for _, prefix := range []string{"http://", "https://"} {
		if strings.HasPrefix(origin, prefix) && strings.TrimPrefix(origin, prefix) == host {
			return true
		}
	}
	return false
}
