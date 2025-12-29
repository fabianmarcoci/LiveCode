package handlers

import (
	"database/sql"
	"errors"
	"livecode-api/config"
	"os"
	"strconv"
	"time"

	"livecode-api/models"
	"livecode-api/utils"

	"github.com/golang-jwt/jwt/v5"
)

func RefreshTokenInternal(refreshToken string, db *sql.DB) (models.RefreshTokenResponse, error) {
	_, claims, err := utils.VerifyJWT(refreshToken)
	if err != nil {
		return models.RefreshTokenResponse{
			Success: false,
			Message: "Invalid or expired refresh token.",
		}, nil
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return models.RefreshTokenResponse{
			Success: false,
			Message: "Invalid user ID in token.",
		}, nil
	}

	var user models.UserData
	query := `SELECT id, username, email FROM users WHERE id = $1 LIMIT 1`
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		return models.RefreshTokenResponse{
			Success: false,
			Message: "User not found.",
		}, nil
	}

	if err != nil {
		return models.RefreshTokenResponse{}, errors.New("database error during token refresh")
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
	accessTokenString, err := accessToken.SignedString([]byte(config.JWTSecret))

	if err != nil {
		return models.RefreshTokenResponse{}, errors.New("failed to generate access token")
	}

	return models.RefreshTokenResponse{
		Success:     true,
		Message:     "Access token refreshed successfully.",
		AccessToken: accessTokenString,
	}, nil
}
