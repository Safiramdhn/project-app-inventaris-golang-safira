package models

import "time"

type LoginRequest struct {
	Username string
	Password string
}

type Session struct {
	UserID       int
	SessionToken string
	ExpiresAt    time.Time
	IsActive     bool
}
