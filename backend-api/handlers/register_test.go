package handlers

import (
	"os"
	"testing"

	"livecode-api/database"
	"livecode-api/models"

	"github.com/joho/godotenv"
)

func TestRegisterUserSuccess(t *testing.T) {
	godotenv.Load("../.env")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("DATABASE_URL not set")
	}

	err := database.Connect(databaseURL)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	payload := models.RegisterRequest{
		Email:    "test@example6.com",
		Username: "@testuser",
		Password: "TestPass123!",
	}

	response, err := RegisterUserInternal(payload, database.DB)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
		t.Errorf("Message: %s", response.Message)

		if len(response.FieldErrors) > 0 {
			t.Errorf("Field errors:")
			for _, fieldErr := range response.FieldErrors {
				t.Errorf("  - %s: %s", fieldErr.Field, fieldErr.Message)
			}
		}

		t.FailNow()
	}

	if response.Message != "your account has been created successfully" {
		t.Errorf("Unexpected message: %s", response.Message)
	}

	if response.User == nil {
		t.Errorf("Expected user data, got nil")
	} else {
		if response.User.Email != payload.Email {
			t.Errorf("Expected email %s, got %s", payload.Email, response.User.Email)
		}
		if response.User.Username != payload.Username {
			t.Errorf("Expected username %s, got %s", payload.Username, response.User.Username)
		}
	}

	database.DB.Exec("DELETE FROM users WHERE email = $1", payload.Email)
}
