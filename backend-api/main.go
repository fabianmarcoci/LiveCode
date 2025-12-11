package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"livecode-api/database"
	"livecode-api/middleware"
	"livecode-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	GinMode     string
}

func main() {
	cfg := loadConfig()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer database.Close()

	router := setupRouter()

	runServer(router, cfg.Port)
}

func loadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	ginMode := os.Getenv("GIN_MODE")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if port == "" {
		port = "3000"
	}
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	return &Config{
		DatabaseURL: databaseURL,
		Port:        port,
		GinMode:     ginMode,
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	authLimiter := middleware.NewRateLimiter(5, 5)
	checkFieldLimiter := middleware.NewRateLimiter(10, 10)
	generalLimiter := middleware.NewRateLimiter(50, 10)

	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")
	v1.Use(generalLimiter.Limit())
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/refresh", routes.RefreshToken)
			authRoutes.POST("/register", authLimiter.Limit(), middleware.ValidateRegisterInput(), routes.Register)
			authRoutes.POST("/login", authLimiter.Limit(), middleware.ValidateLoginInput(), routes.Login)
			authRoutes.GET("/check-field", checkFieldLimiter.Limit(), routes.CheckFieldAvailable)
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

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "connected",
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
		log.Printf("Server starting on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
