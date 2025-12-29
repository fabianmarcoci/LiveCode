package utils

import (
	"errors"
	"livecode-api/config"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	if config.JWTSecret == "" {
		return nil, nil, errors.New("JWT_SECRET not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("invalid token claims")
	}

	return token, claims, nil
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func GenerateTokenPair(userID string, username string, email string) (*TokenPair, error) {
	if config.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}

	accessTokenExpiry, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"))
	if accessTokenExpiry == 0 {
		accessTokenExpiry = 15
	}

	accessTokenClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Duration(accessTokenExpiry) * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.JWTSecret))

	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshTokenExpiry, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS"))
	if refreshTokenExpiry == 0 {
		refreshTokenExpiry = 30
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(refreshTokenExpiry) * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.JWTSecret))

	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
