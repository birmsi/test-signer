package helpers

import (
	"testing"
)

func TestGetUserIDFromJWT(t *testing.T) {
	validJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJqb3NlIn0.ZkhIJqEf0Uy_tIEcap9TtonWLXkNs96LW_7z21fKqhM"
	invalidJWT := "invalid.jwt.token"
	malformedJWT := "malformedjwt"

	// Test case: valid JWT
	userID, err := GetUserIDFromJWT(validJWT)
	if err != nil {
		t.Errorf("Unexpected error for valid JWT: %v", err)
	}
	if userID != "jose" {
		t.Errorf("Unexpected userID for valid JWT: got %s, want %s", userID, "expected_user_id")
	}

	// Test case: invalid JWT
	_, err = GetUserIDFromJWT(invalidJWT)
	if err == nil {
		t.Error("Expected error for invalid JWT, but got nil")
	}

	// Test case: malformed JWT
	_, err = GetUserIDFromJWT(malformedJWT)
	if err == nil {
		t.Error("Expected error for malformed JWT, but got nil")
	}
}
