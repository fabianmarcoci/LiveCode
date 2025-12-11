package handlers

import (
	"database/sql"
	"errors"

	"livecode-api/models"
	"livecode-api/utils"

	"github.com/google/uuid"
)

func CheckFieldAvailableInternal(field string, value string, db *sql.DB) (*bool, error) {
	if field != "email" && field != "username" {
		return nil, nil
	}

	query := "SELECT id FROM users WHERE " + field + " = $1 LIMIT 1"

	var id string
	err := db.QueryRow(query, value).Scan(&id)

	if err == sql.ErrNoRows {
		available := true
		return &available, nil
	}

	if err != nil {
		return nil, err
	}

	available := false
	return &available, nil
}

func RegisterUserInternal(payload models.RegisterRequest, db *sql.DB) (models.RegisterResponse, error) {
	var fieldErrors []models.FieldError

	var emailID string
	emailErr := db.QueryRow(
		"SELECT id FROM users WHERE email = $1 LIMIT 1",
		payload.Email,
	).Scan(&emailID)

	if emailErr != nil && emailErr != sql.ErrNoRows {
		return models.RegisterResponse{}, errors.New("database error during email check")
	}

	if emailErr == nil {
		fieldErrors = append(fieldErrors, models.FieldError{
			Field:   "email",
			Message: "This email is already taken.",
		})
	}

	var usernameID string
	usernameErr := db.QueryRow(
		"SELECT id FROM users WHERE username = $1 LIMIT 1",
		payload.Username,
	).Scan(&usernameID)

	if usernameErr != nil && usernameErr != sql.ErrNoRows {
		return models.RegisterResponse{}, errors.New("database error during username check")
	}

	if usernameErr == nil {
		fieldErrors = append(fieldErrors, models.FieldError{
			Field:   "username",
			Message: "This username is already taken.",
		})
	}

	if len(fieldErrors) > 0 {
		return models.RegisterResponse{
			Success:     false,
			FieldErrors: fieldErrors,
			Message:     "Account could not be created.",
		}, nil
	}

	passwordHash, err := utils.HashPassword(payload.Password)
	if err != nil {
		return models.RegisterResponse{}, errors.New("an unexpected error occurred. Please try again")
	}

	userID := uuid.New().String()

	_, err = db.Exec(
		"INSERT INTO users (id, username, email, password_hash, is_oauth) VALUES ($1, $2, $3, $4, $5)",
		userID, payload.Username, payload.Email, passwordHash, false,
	)

	if err != nil {
		return models.RegisterResponse{}, errors.New("an unexpected error occurred. Please try again")
	}

	tokens, err := utils.GenerateTokenPair(userID, payload.Username, payload.Email)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	return models.RegisterResponse{
		Success:      true,
		Message:      "Your account has been created successfully.",
		AccessToken:  &tokens.AccessToken,
		RefreshToken: &tokens.RefreshToken,
		User: &models.UserData{
			ID:       userID,
			Username: payload.Username,
			Email:    payload.Email,
		},
	}, nil
}
