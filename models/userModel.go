package models

import "time"

type User struct {
	ID           int    `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
	Email        string `json:"email,omitempty"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserDTO struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
