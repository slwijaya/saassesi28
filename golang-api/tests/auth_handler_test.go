package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-api/handlers"
	"golang-api/models"
	"golang-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 🔹 Test Login dengan Mock Database
func TestLoginHandlerWithMock(t *testing.T) {
	mockDB := new(mocks.MockDatabase)
	mockUser := &models.User{
		ID:       1,
		Email:    "user1@example.com",
		Password: "$2a$10$5q1qUwJvZ8rqy8OGEX.tyeZQG6fu4AOJMCrLGJdP9eniJTLqIrI.C", // Hash bcrypt
		Package:  "Basic",
	}

	// Setting mock agar mengembalikan user jika email cocok
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockUser)

	// Inisialisasi handler dengan Mock Database
	authHandler := handlers.AuthHandler{DB: mockDB}

	// Buat request dummy
	requestPayload := map[string]string{
		"email":    "user1@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(requestPayload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	authHandler.LoginHandler(recorder, req)

	// Validasi response
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token") // Harus ada token
}
