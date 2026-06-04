package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashed == password {
		t.Fatalf("Hashed password should not be equal to plain text password")
	}

	if !CheckPasswordHash(password, hashed) {
		t.Errorf("Password verification failed for correct password")
	}

	if CheckPasswordHash("wrongPassword", hashed) {
		t.Errorf("Password verification should have failed for incorrect password")
	}
}
