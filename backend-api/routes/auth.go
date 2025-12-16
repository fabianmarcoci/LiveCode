package routes

import (
	"net/http"

	"livecode-api/database"
	"livecode-api/handlers"
	"livecode-api/middleware"
	"livecode-api/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckFieldAvailable(c *gin.Context) {
	field, _ := c.Get("validated_field")
	value, _ := c.Get("validated_value")

	available, err := handlers.CheckFieldAvailableInternal(field.(string), value.(string), database.DB)

	if err != nil {
		middleware.GetLogger(c).Error("check_field_database_error",
			zap.String("field", field.(string)),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"available": nil})
		return
	}

	middleware.GetLogger(c).Info("check_field_success",
		zap.String("field", field.(string)),
		zap.Bool("available", *available),
	)

	c.JSON(http.StatusOK, gin.H{"available": available})
}

func Register(c *gin.Context) {
	validatedPayload, exists := c.Get("validated_payload")
	if !exists {
		middleware.GetLogger(c).Error("register_validation_missing")
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Success: false,
			Message: "Validation error occurred.",
		})
		return
	}

	payload := validatedPayload.(struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	})

	req := models.RegisterRequest{
		Email:    payload.Email,
		Username: payload.Username,
		Password: payload.Password,
	}

	response, err := handlers.RegisterUserInternal(req, database.DB)

	if err != nil {
		middleware.GetLogger(c).Error("register_failed",
			zap.String("error", err.Error()),
			zap.String("email", req.Email),
			zap.String("username", req.Username),
		)
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Success: false,
			Message: "An unexpected error occurred. Please try again.",
		})
		return
	}

	if !response.Success {
		middleware.GetLogger(c).Warn("register_validation_failed",
			zap.String("email", req.Email),
			zap.String("username", req.Username),
			zap.Int("field_errors_count", len(response.FieldErrors)),
		)
		c.JSON(http.StatusOK, response)
		return
	}

	middleware.GetLogger(c).Info("register_success",
		zap.String("user_id", response.User.ID),
		zap.String("username", response.User.Username),
	)

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	validatedPayload, exists := c.Get("validated_payload")
	if !exists {
		middleware.GetLogger(c).Error("login_validation_missing")
		c.JSON(http.StatusInternalServerError, models.LoginResponse{
			Success: false,
			Message: "Validation error occurred.",
		})
		return
	}

	payload := validatedPayload.(struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	})

	req := models.LoginRequest{
		Identifier: payload.Identifier,
		Password:   payload.Password,
	}

	response, err := handlers.LoginUserInternal(req, database.DB)

	if err != nil {
		middleware.GetLogger(c).Error("login_failed",
			zap.String("error", err.Error()),
			zap.String("identifier", req.Identifier),
		)
		c.JSON(http.StatusInternalServerError, models.LoginResponse{
			Success: false,
			Message: "An unexpected error occurred. Please try again.",
		})
		return
	}

	if !response.Success {
		middleware.GetLogger(c).Warn("login_invalid_credentials",
			zap.String("identifier", req.Identifier),
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	middleware.GetLogger(c).Info("login_success",
		zap.String("user_id", response.User.ID),
		zap.String("username", response.User.Username),
	)

	c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {
	validatedPayload, exists := c.Get("validated_payload")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Validation error occurred.",
		})
		return
	}

	payload := validatedPayload.(struct {
		RefreshToken string `json:"refresh_token"`
	})

	response, err := handlers.RefreshTokenInternal(payload.RefreshToken, database.DB)

	if err != nil {
		middleware.GetLogger(c).Error("refresh_token_failed",
			zap.String("error", err.Error()),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again.",
		})
		return
	}

	if !response.Success {
		middleware.GetLogger(c).Warn("refresh_token_invalid",
			zap.String("message", response.Message),
		)
	}

	statusCode := http.StatusOK
	if !response.Success {
		statusCode = http.StatusUnauthorized
	}

	c.JSON(statusCode, response)
}
