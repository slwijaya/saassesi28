package tests

import (
	"testing"
	"time"

	"golang-api/handlers"

	"github.com/stretchr/testify/assert"
)

// 🔹 Test JWT Token Generation dan Parsing
func TestJWTTokenTrial(t *testing.T) {
	email := "user@example.com"
	packageType := "Basic"
	trialExpiresAt := time.Now().Add(5 * time.Minute)

	token, err := handlers.GenerateJWT(email, packageType, trialExpiresAt)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verifikasi isi token
	claims, _ := handlers.ParseJWT(token)
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, packageType, claims["package"])
}
