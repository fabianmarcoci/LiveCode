package routes

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"livecode-api/database"
	"livecode-api/handlers"
	"livecode-api/models"
	"livecode-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Register(c *gin.Context) {
	validatedPayload, exists := c.Get("validated_payload")
	if !exists {
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
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Success: false,
			Message: "An unexpected error occurred. Please try again.",
		})
		return
	}

	statusCode := http.StatusCreated
	if !response.Success {
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, response)
}

func CheckFieldAvailable(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")

	if field == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"available": nil})
		return
	}

	available, err := handlers.CheckFieldAvailableInternal(field, value, database.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"available": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}

func Login(c *gin.Context) {
	validatedPayload, exists := c.Get("validated_payload")
	if !exists {
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
		c.JSON(http.StatusInternalServerError, models.LoginResponse{
			Success: false,
			Message: "An unexpected error occurred. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Refresh token is required.",
		})
		return
	}

	_, claims, err := utils.VerifyJWT(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid or expired refresh token.",
		})
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID in token.",
		})
		return
	}

	var user models.UserData
	query := `SELECT id, username, email FROM users WHERE id = $1 LIMIT 1`
	err = database.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error.",
		})
		return
	}

	accessTokenExpiry, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"))
	if accessTokenExpiry == 0 {
		accessTokenExpiry = 15
	}

	accessTokenClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Duration(accessTokenExpiry) * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate access token.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"access_token": accessTokenString,
		"message":      "Access token refreshed successfully.",
	})
}
