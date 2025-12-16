package handlers

import (
	"database/sql"
	"errors"

	"livecode-api/models"
	"livecode-api/utils"
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

	tokens, err := utils.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		return models.LoginResponse{}, errors.New("token generation failed")
	}

	return models.LoginResponse{
		Success:      true,
		Message:      "Login successful.",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         &user,
	}, nil
}
