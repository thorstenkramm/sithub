// Package users provides user management and persistence.
package users

import "errors"

// Record represents a user row from the database.
type Record struct {
	ID           string
	Email        string
	DisplayName  string
	PasswordHash string
	UserSource   string
	EntraID      string
	IsAdmin      bool
	LastLogin    string
	CreatedAt    string
	UpdatedAt    string
}

// Sentinel errors for user operations.
var (
	ErrEmailConflict = errors.New("email already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrNotLocalUser  = errors.New("operation only allowed for local users")
)
