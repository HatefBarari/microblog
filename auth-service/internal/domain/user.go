package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
	Verified     bool
	CreatedAt    time.Time
}
