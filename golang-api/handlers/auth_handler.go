package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang-api/config"
	"golang-api/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Secret Key untuk JWT
var jwtSecret = []byte("your_secret_key")

// 🔹 Fungsi untuk menghasilkan JWT Token
func GenerateJWT(email string, packageType string, trialExpiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"email":         email,
		"package":       packageType,
		"trial_expires": trialExpiresAt.Unix(),
		"exp":           time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 🔹 ParseJWT untuk memverifikasi token dan mengambil claims
func ParseJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 🔹 Fungsi untuk hashing password sebelum disimpan ke database
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// 🔹 Fungsi untuk memverifikasi password saat login
func checkPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// 🔹 RegisterHandler untuk pendaftaran pengguna baru dengan logging
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	config.Logger.Info("Received request for RegisterHandler") // Logging masuk ke handler

	if r.Method != http.MethodPost {
		config.Logger.Warn("Invalid request method for RegisterHandler")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		config.Logger.Warn("Invalid input in RegisterHandler: ", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi input
	if newUser.Email == "" || newUser.Password == "" || newUser.FullName == "" {
		config.Logger.Warn("Missing required fields in RegisterHandler")
		http.Error(w, "Email, full name, and password are required", http.StatusBadRequest)
		return
	}

	trialExpiresAt := time.Now().Add(5 * time.Minute)

	// 🔹 Logging sebelum hashing password
	config.Logger.Info("Hashing password for new user: ", newUser.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		config.Logger.Error("Failed to hash password for user: ", newUser.Email, " Error: ", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 🔹 Logging sebelum memasukkan data user ke database
	config.Logger.Info("Inserting new user into database: ", newUser.Email)
	query := `INSERT INTO users (email, password, full_name, package, trial_expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err = config.DB.QueryRow(query, newUser.Email, hashedPassword, newUser.FullName, newUser.Package, trialExpiresAt).Scan(&newUser.ID, &newUser.CreatedAt)
	if err != nil {
		config.Logger.Error("Error inserting user into database: ", err)
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	// 🔹 Logging sebelum pembuatan token JWT
	config.Logger.Info("Generating JWT token for new user: ", newUser.Email)
	token, err := GenerateJWT(newUser.Email, newUser.Package, trialExpiresAt)
	if err != nil {
		config.Logger.Error("Failed to generate token for user: ", newUser.Email, " Error: ", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 🔹 Logging sukses register user
	config.Logger.Info("User registered successfully: ", newUser.Email)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"token":   token,
	})
}

// 🔹 LoginHandler untuk autentikasi pengguna dengan logging
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	config.Logger.Info("Received request for LoginHandler") // Logging masuk ke handler

	if r.Method != http.MethodPost {
		config.Logger.Warn("Invalid request method for LoginHandler")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginUser models.User
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		config.Logger.Warn("Invalid input in LoginHandler: ", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if loginUser.Email == "" || loginUser.Password == "" {
		config.Logger.Warn("Missing email or password in LoginHandler")
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// 🔹 Logging sebelum mengambil user dari database
	config.Logger.Info("Fetching user from database: ", loginUser.Email)
	var storedUser models.User
	var storedHashedPassword string
	query := `SELECT id, email, password, full_name, package, trial_expires_at, created_at FROM users WHERE email=$1`
	err := config.DB.QueryRow(query, loginUser.Email).Scan(&storedUser.ID, &storedUser.Email, &storedHashedPassword, &storedUser.FullName, &storedUser.Package, &storedUser.TrialExpiresAt, &storedUser.CreatedAt)
	if err != nil {
		config.Logger.Warn("User not found in LoginHandler: ", loginUser.Email)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 🔹 Logging sebelum memverifikasi password
	config.Logger.Info("Verifying password for user: ", loginUser.Email)
	if err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(loginUser.Password)); err != nil {
		config.Logger.Warn("Invalid password attempt for user: ", loginUser.Email)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 🔹 Logging sebelum pembuatan token JWT
	config.Logger.Info("Generating JWT token for user: ", loginUser.Email)
	token, err := GenerateJWT(storedUser.Email, storedUser.Package, storedUser.TrialExpiresAt)
	if err != nil {
		config.Logger.Error("Failed to generate token for user: ", storedUser.Email, " Error: ", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 🔹 Logging sukses login user
	config.Logger.Info("User logged in successfully: ", loginUser.Email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}

// Tambahkan struct handler agar bisa menerima dependency injection
type AuthHandler struct {
	DB DatabaseInterface // Pastikan pakai interface, bukan langsung `config.DB`
}

// Interface untuk menggantikan sql.DB dalam unit test
type DatabaseInterface interface {
	QueryRow(query string, args ...interface{}) *models.User
}

// 🔹 LoginHandler untuk autentikasi pengguna dengan logging
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginUser models.User
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if loginUser.Email == "" || loginUser.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// 🔹 Logging sebelum mengambil user dari database
	var storedUser *models.User
	storedUser = h.DB.QueryRow("SELECT id, email, password, package FROM users WHERE email = $1", loginUser.Email)

	if storedUser == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 🔹 Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginUser.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 🔹 Generate JWT Token
	token, err := GenerateJWT(storedUser.Email, storedUser.Package, time.Now().Add(5*time.Minute))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
