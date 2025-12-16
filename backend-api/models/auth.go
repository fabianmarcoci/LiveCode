package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	IsOAuth      bool      `json:"is_oauth" db:"is_oauth"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=17"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	Success      bool         `json:"success"`
	FieldErrors  []FieldError `json:"field_errors,omitempty"`
	Message      string       `json:"message"`
	AccessToken  *string      `json:"access_token,omitempty"`
	RefreshToken *string      `json:"refresh_token,omitempty"`
	User         *UserData    `json:"user,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type UserData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	User         *UserData `json:"user,omitempty"`
}

type RefreshTokenResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token,omitempty"`
}
