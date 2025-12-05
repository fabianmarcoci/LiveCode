package handlers

import (
	"os"
	"testing"

	"livecode-api/database"
	"livecode-api/models"
	"livecode-api/utils"
)

func TestLoginUserInternal_Success(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	if err := database.Connect(databaseURL); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	testEmail := "login_test@example.com"
	testUsername := "@logintest"
	testPassword := "TestPassword123"

	passwordHash, err := utils.HashPassword(testPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (id, username, email, password_hash, is_oauth) VALUES (gen_random_uuid(), $1, $2, $3, false) ON CONFLICT (email) DO NOTHING",
		testUsername, testEmail, passwordHash,
	)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	payload := models.LoginRequest{
		Identifier: testEmail,
		Password:   testPassword,
	}

	response, err := LoginUserInternal(payload, database.DB)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}

	if response.Message != "Login successful." {
		t.Errorf("Expected message 'Login successful.', got: %s", response.Message)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if response.RefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

	if response.User == nil {
		t.Fatal("Expected user data, got nil")
	}

	if response.User.Email != testEmail {
		t.Errorf("Expected email %s, got: %s", testEmail, response.User.Email)
	}

	database.DB.Exec("DELETE FROM users WHERE email = $1", testEmail)
}

func TestLoginUserInternal_InvalidCredentials(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	if err := database.Connect(databaseURL); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	payload := models.LoginRequest{
		Identifier: "nonexistent@example.com",
		Password:   "password123",
	}

	response, err := LoginUserInternal(payload, database.DB)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if response.Success {
		t.Error("Expected success=false, got true")
	}

	if response.Message != "Invalid credentials." {
		t.Errorf("Expected 'Invalid credentials.', got: %s", response.Message)
	}

	if response.AccessToken != "" {
		t.Error("Expected empty access token")
	}
}

func TestLoginUserInternal_WrongPassword(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	if err := database.Connect(databaseURL); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	testEmail := "wrong_password_test@example.com"
	testUsername := "@wrongpasstest"
	correctPassword := "CorrectPassword123"

	passwordHash, err := utils.HashPassword(correctPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (id, username, email, password_hash, is_oauth) VALUES (gen_random_uuid(), $1, $2, $3, false) ON CONFLICT (email) DO NOTHING",
		testUsername, testEmail, passwordHash,
	)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	payload := models.LoginRequest{
		Identifier: testEmail,
		Password:   "WrongPassword123",
	}

	response, err := LoginUserInternal(payload, database.DB)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if response.Success {
		t.Error("Expected success=false, got true")
	}

	if response.Message != "Invalid credentials." {
		t.Errorf("Expected 'Invalid credentials.', got: %s", response.Message)
	}

	database.DB.Exec("DELETE FROM users WHERE email = $1", testEmail)
}
