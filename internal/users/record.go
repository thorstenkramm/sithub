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
	// IsDataExporter grants access to booking data export. It is independent of IsAdmin.
	IsDataExporter bool
	// Deleted marks a deactivated (soft-deleted) user who may no longer sign in.
	Deleted   bool
	CreatedAt string
	UpdatedAt string
}

// Sentinel errors for user operations.
var (
	ErrEmailConflict = errors.New("email already exists")
	ErrUserNotFound  = errors.New("user not found")
)
