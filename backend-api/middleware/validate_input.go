package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"livecode-api/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	emailRegex    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	usernameRegex = regexp.MustCompile(`^@[a-z0-9]{3,16}$`)
	lowerRegex    = regexp.MustCompile(`[a-z]`)
	upperRegex    = regexp.MustCompile(`[A-Z]`)
	digitRegex    = regexp.MustCompile(`\d`)
	specialRegex  = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func ValidateCheckFieldAvailable() gin.HandlerFunc {
	return func(c *gin.Context) {
		field := c.Query("field")
		value := c.Query("value")

		if field == "" || value == "" {
			GetLogger(c).Warn("check_field_missing_params",
				zap.String("field", field),
				zap.String("value_empty", value),
			)
			c.JSON(http.StatusBadRequest, gin.H{"available": nil})
			c.Abort()
			return
		}

		value = strings.TrimSpace(strings.ToLower(value))

		var validationErr *models.FieldError
		switch field {
		case "email":
			validationErr = ValidateEmail(value)
		case "username":
			validationErr = ValidateUsername(value)
		default:
			GetLogger(c).Warn("check_field_invalid_field_type",
				zap.String("field", field),
			)
			c.JSON(http.StatusBadRequest, gin.H{"available": nil})
			c.Abort()
			return
		}

		if validationErr != nil {
			GetLogger(c).Warn("check_field_validation_failed",
				zap.String("field", field),
				zap.String("error", validationErr.Message),
			)
			c.JSON(http.StatusBadRequest, gin.H{"available": nil})
			c.Abort()
			return
		}

		if containsNullBytes(value) {
			GetLogger(c).Warn("check_field_null_bytes",
				zap.String("field", field),
			)
			c.JSON(http.StatusBadRequest, gin.H{"available": nil})
			c.Abort()
			return
		}

		c.Set("validated_field", field)
		c.Set("validated_value", value)
		c.Next()
	}
}

func ValidateRegisterInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			c.Abort()
			return
		}

		payload.Email = strings.TrimSpace(strings.ToLower(payload.Email))
		payload.Username = strings.TrimSpace(strings.ToLower(payload.Username))

		errors := []models.FieldError{}

		if emailErr := ValidateEmail(payload.Email); emailErr != nil {
			errors = append(errors, *emailErr)
		}

		if usernameErr := ValidateUsername(payload.Username); usernameErr != nil {
			errors = append(errors, *usernameErr)
		}

		passwordErrors := ValidatePassword(payload.Password)
		errors = append(errors, passwordErrors...)

		if containsNullBytes(payload.Email) || containsNullBytes(payload.Username) || containsNullBytes(payload.Password) {
			errors = append(errors, models.FieldError{
				Field:   "general",
				Message: "Invalid characters detected",
			})
		}

		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validation failed",
				"errors":  errors,
			})
			c.Abort()
			return
		}

		c.Set("validated_payload", payload)
		c.Next()
	}
}

func ValidateLoginInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			c.Abort()
			return
		}

		payload.Identifier = strings.TrimSpace(strings.ToLower(payload.Identifier))

		errors := []models.FieldError{}

		if len(payload.Identifier) == 0 {
			errors = append(errors, models.FieldError{
				Field:   "identifier",
				Message: "Email or username is required",
			})
		}

		if len(payload.Identifier) > 255 {
			errors = append(errors, models.FieldError{
				Field:   "identifier",
				Message: "Email or username must not exceed 255 characters",
			})
		}

		if len(payload.Password) == 0 {
			errors = append(errors, models.FieldError{
				Field:   "password",
				Message: "Password is required",
			})
		}

		if len(payload.Password) > 72 {
			errors = append(errors, models.FieldError{
				Field:   "password",
				Message: "Password must not exceed 72 characters",
			})
		}

		if containsNullBytes(payload.Identifier) || containsNullBytes(payload.Password) {
			errors = append(errors, models.FieldError{
				Field:   "general",
				Message: "Invalid characters detected",
			})
		}

		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validation failed",
				"errors":  errors,
			})
			c.Abort()
			return
		}

		c.Set("validated_payload", payload)
		c.Next()
	}
}

func containsNullBytes(s string) bool {
	return strings.Contains(s, "\x00")
}

func ValidateEmail(email string) *models.FieldError {
	if len(email) > 255 {
		return &models.FieldError{
			Field:   "email",
			Message: "Email must not exceed 255 characters",
		}
	}

	if !emailRegex.MatchString(email) {
		return &models.FieldError{
			Field:   "email",
			Message: "Invalid email format",
		}
	}

	return nil
}

func ValidateUsername(username string) *models.FieldError {
	if len(username) > 17 {
		return &models.FieldError{
			Field:   "username",
			Message: "Username must not exceed 17 characters",
		}
	}

	if !usernameRegex.MatchString(username) {
		return &models.FieldError{
			Field:   "username",
			Message: "Username must start with @ and contain 3-16 lowercase letters or digits",
		}
	}

	return nil
}

func ValidatePassword(password string) []models.FieldError {
	errors := []models.FieldError{}

	if len(password) < 8 {
		errors = append(errors, models.FieldError{
			Field:   "password",
			Message: "Password must be at least 8 characters",
		})
	}

	if len(password) > 72 {
		errors = append(errors, models.FieldError{
			Field:   "password",
			Message: "Password must not exceed 72 characters",
		})
	}

	hasLower := lowerRegex.MatchString(password)
	hasUpper := upperRegex.MatchString(password)
	hasDigit := digitRegex.MatchString(password)
	hasSpecial := specialRegex.MatchString(password)

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		errors = append(errors, models.FieldError{
			Field:   "password",
			Message: "Password must include lowercase, uppercase, number and special character",
		})
	}

	return errors
}
