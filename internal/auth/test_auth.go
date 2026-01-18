package auth

import "strings"

// TestUser returns a user only when test auth is explicitly enabled.
// This is for local E2E runs and must not be enabled in production.
func (s *Service) TestUser() *User {
	if !s.testAuth.Enabled {
		return nil
	}

	id := strings.TrimSpace(s.testAuth.UserID)
	if id == "" {
		id = "test-user"
	}
	name := strings.TrimSpace(s.testAuth.UserName)
	if name == "" {
		name = "Test User"
	}

	return &User{ID: id, Name: name}
}
