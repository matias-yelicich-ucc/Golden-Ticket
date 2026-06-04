package utils

import (
	"os"
	"testing"
)

func TestJWTTokenLifecycle(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET")

	userID := uint(42)
	rol := "cliente"

	token, err := GenerateToken(userID, rol)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatalf("Generated token is empty")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Rol != rol {
		t.Errorf("Expected Rol %s, got %s", rol, claims.Rol)
	}

	_, err = ValidateToken("invalid.token.string")
	if err == nil {
		t.Errorf("Validating an invalid token string should fail")
	}
}
