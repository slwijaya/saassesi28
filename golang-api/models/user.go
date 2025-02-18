package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	FullName       string    `json:"full_name"`
	Package        string    `json:"package"`
	TrialExpiresAt time.Time `json:"trial_expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// 🔹 Mock Model (digunakan dalam unit test)
type MockUser struct {
	ID             int
	Email          string
	Password       string
	FullName       string
	Package        string
	TrialExpiresAt time.Time
}
