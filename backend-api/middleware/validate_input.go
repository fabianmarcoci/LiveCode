package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	emailRegex    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	usernameRegex = regexp.MustCompile(`^@[a-z0-9]{3,16}$`)
	lowerRegex    = regexp.MustCompile(`[a-z]`)
	upperRegex    = regexp.MustCompile(`[A-Z]`)
	digitRegex    = regexp.MustCompile(`\d`)
	specialRegex  = regexp.MustCompile(`[^A-Za-z0-9]`)
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
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

		errors := []ValidationError{}

		if len(payload.Email) > 255 {
			errors = append(errors, ValidationError{
				Field:   "email",
				Message: "Email must not exceed 255 characters",
			})
		}

		if !emailRegex.MatchString(payload.Email) {
			errors = append(errors, ValidationError{
				Field:   "email",
				Message: "Invalid email format",
			})
		}

		if len(payload.Username) > 17 {
			errors = append(errors, ValidationError{
				Field:   "username",
				Message: "Username must not exceed 17 characters",
			})
		}

		if !usernameRegex.MatchString(payload.Username) {
			errors = append(errors, ValidationError{
				Field:   "username",
				Message: "Username must start with @ and contain 3-16 lowercase letters or digits",
			})
		}

		if len(payload.Password) < 8 {
			errors = append(errors, ValidationError{
				Field:   "password",
				Message: "Password must be at least 8 characters",
			})
		}

		if len(payload.Password) > 72 {
			errors = append(errors, ValidationError{
				Field:   "password",
				Message: "Password must not exceed 72 characters",
			})
		}

		hasLower := lowerRegex.MatchString(payload.Password)
		hasUpper := upperRegex.MatchString(payload.Password)
		hasDigit := digitRegex.MatchString(payload.Password)
		hasSpecial := specialRegex.MatchString(payload.Password)

		if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
			errors = append(errors, ValidationError{
				Field:   "password",
				Message: "Password must include lowercase, uppercase, number and special character",
			})
		}

		if containsNullBytes(payload.Email) || containsNullBytes(payload.Username) || containsNullBytes(payload.Password) {
			errors = append(errors, ValidationError{
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

		errors := []ValidationError{}

		if len(payload.Identifier) == 0 {
			errors = append(errors, ValidationError{
				Field:   "identifier",
				Message: "Email or username is required",
			})
		}

		if len(payload.Identifier) > 255 {
			errors = append(errors, ValidationError{
				Field:   "identifier",
				Message: "Email or username must not exceed 255 characters",
			})
		}

		if len(payload.Password) == 0 {
			errors = append(errors, ValidationError{
				Field:   "password",
				Message: "Password is required",
			})
		}

		if len(payload.Password) > 72 {
			errors = append(errors, ValidationError{
				Field:   "password",
				Message: "Password must not exceed 72 characters",
			})
		}

		if containsNullBytes(payload.Identifier) || containsNullBytes(payload.Password) {
			errors = append(errors, ValidationError{
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
