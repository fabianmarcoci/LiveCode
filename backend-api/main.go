package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"livecode-api/config"
	"livecode-api/database"
	"livecode-api/middleware"
	"livecode-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DatabaseURL string
	Port        string
	GinMode     string
	JWTSecret   string
}

func getEnvOrSecret(envKey, secretPath string) string {
	value := os.Getenv(envKey)
	if value != "" {
		return value
	}

	secretFile := os.Getenv(envKey + "_FILE")
	if secretFile != "" {
		data, err := os.ReadFile(secretFile)
		if err != nil {
			middleware.Logger.Warn("failed to read secret file",
				zap.String("file", secretFile),
				zap.Error(err),
			)
			return ""
		}
		return strings.TrimSpace(string(data))
	}

	return ""
}

func main() {
	middleware.InitLogger()

	cfg := loadConfig()
	config.Init(cfg.JWTSecret)

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		middleware.Logger.Fatal("database connection failed",
			zap.Error(err),
		)
	}
	defer database.Close()

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		middleware.Logger.Fatal("database migrations failed",
			zap.Error(err),
		)
	}

	router := setupRouter()

	runServer(router, cfg.Port)
}

func loadConfig() *Config {
	if os.Getenv("DOCKER_ENV") != "true" {
		_ = godotenv.Load("../.env.local")
		_ = godotenv.Load("../.env")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	ginMode := getEnvOrSecret("GIN_MODE", "/run/secrets/gin_mode")
	jwtSecret := getEnvOrSecret("JWT_SECRET", "/run/secrets/jwt_secret")

	if os.Getenv("DOCKER_ENV") == "true" && databaseURL == "" {
		pgUser := os.Getenv("POSTGRES_USER")
		pgDB := os.Getenv("POSTGRES_DB")
		pgPassword := getEnvOrSecret("POSTGRES_PASSWORD", "/run/secrets/postgres_password")

		if pgUser == "" || pgDB == "" || pgPassword == "" {
			middleware.Logger.Fatal("POSTGRES_USER, POSTGRES_DB, POSTGRES_PASSWORD required for Docker")
		}

		databaseURL = "postgresql://" + pgUser + ":" + pgPassword + "@postgres:5432/" + pgDB + "?sslmode=disable"
	}

	if jwtSecret == "" {
		middleware.Logger.Fatal("JWT_SECRET is required")
	}
	if port == "" {
		port = "3000"
	}
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	if databaseURL == "" {
		middleware.Logger.Fatal("DATABASE_URL is required")
	}

	middleware.Logger.Info("configuration loaded",
		zap.String("port", port),
		zap.String("gin_mode", ginMode),
		zap.Bool("jwt_from_secret_file", os.Getenv("JWT_SECRET_FILE") != ""),
		zap.Bool("database_from_secrets", os.Getenv("DOCKER_ENV") == "true"),
	)

	return &Config{
		DatabaseURL: databaseURL,
		Port:        port,
		GinMode:     ginMode,
		JWTSecret:   jwtSecret,
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RequestLogger())

	refreshTokenLimiter := middleware.NewRateLimiter(3, 3)
	authLimiter := middleware.NewRateLimiter(5, 5)
	checkFieldLimiter := middleware.NewRateLimiter(10, 10)
	clientMonitoringLimiter := middleware.NewRateLimiter(2, 2)

	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/refresh", refreshTokenLimiter.Limit(), middleware.ValidateRefreshToken(), routes.RefreshToken)
			authRoutes.POST("/register", authLimiter.Limit(), middleware.ValidateRegisterInput(), routes.Register)
			authRoutes.POST("/login", authLimiter.Limit(), middleware.ValidateLoginInput(), routes.Login)
			authRoutes.GET("/check-field", checkFieldLimiter.Limit(), middleware.ValidateCheckFieldAvailable(), routes.CheckFieldAvailable)
		}

		clientMonitoringRoutes := v1.Group("/monitoring")
		clientMonitoringRoutes.Use(clientMonitoringLimiter.Limit(), middleware.ValidateClientErrorLog())
		{
			clientMonitoringRoutes.POST("/client-errors", routes.LogClientError)
		}

		protectedRoutes := v1.Group("")
		protectedRoutes.Use(middleware.AuthMiddleware())
		{
			protectedRoutes.GET("/profile", routes.GetProfile)
		}
	}

	return router
}

func healthCheck(c *gin.Context) {
	if err := database.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "error",
			"database": "disconnected",
			"error":    err.Error(),
		})
		return
	}

	if err := database.CheckMigrationsApplied(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":     "error",
			"database":   "connected",
			"migrations": "not_applied",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"database":   "connected",
		"migrations": "applied",
	})
}

func runServer(router *gin.Engine, port string) {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		middleware.Logger.Info("server starting",
			zap.String("address", "http://localhost:"+port),
			zap.String("port", port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			middleware.Logger.Fatal("server failed to start",
				zap.Error(err),
			)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	middleware.Logger.Info("shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		middleware.Logger.Fatal("server forced to shutdown",
			zap.Error(err),
		)
	}

	middleware.Logger.Info("server exited gracefully")
}
