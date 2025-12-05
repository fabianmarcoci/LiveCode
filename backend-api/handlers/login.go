package handlers

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	"livecode-api/models"
	"livecode-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoginUserInternal(payload models.LoginRequest, db *sql.DB) (models.LoginResponse, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1 OR username = $1 LIMIT 1`

	var user models.UserData
	var passwordHash string

	err := db.QueryRow(query, payload.Identifier).Scan(&user.ID, &user.Username, &user.Email, &passwordHash)

	if err == sql.ErrNoRows {
		return models.LoginResponse{
			Success: false,
			Message: "Invalid credentials.",
		}, nil
	}

	if err != nil {
		return models.LoginResponse{}, errors.New("database error during login")
	}

	if !utils.CheckPasswordHash(payload.Password, passwordHash) {
		return models.LoginResponse{
			Success: false,
			Message: "Invalid credentials.",
		}, nil
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return models.LoginResponse{}, errors.New("JWT_SECRET not configured")
	}

	accessTokenExpiry, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"))
	if accessTokenExpiry == 0 {
		accessTokenExpiry = 15
	}

	refreshTokenExpiry, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"))
	if refreshTokenExpiry == 0 {
		refreshTokenExpiry = 43200
	}

	accessTokenClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Duration(accessTokenExpiry) * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return models.LoginResponse{}, errors.New("failed to generate access token")
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Duration(refreshTokenExpiry) * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
		"jti":     uuid.New().String(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return models.LoginResponse{}, errors.New("failed to generate refresh token")
	}

	return models.LoginResponse{
		Success:      true,
		Message:      "Login successful.",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		User:         &user,
	}, nil
}
