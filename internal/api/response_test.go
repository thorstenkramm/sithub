//revive:disable-next-line var-naming
package api

import "testing"

func TestNewError(t *testing.T) {
	err := NewError(401, "Unauthorized", "Login required", "auth_required")
	if len(err.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(err.Errors))
	}

	item := err.Errors[0]
	if item.Status != "401" {
		t.Fatalf("expected status 401, got %s", item.Status)
	}
	if item.Title != "Unauthorized" || item.Detail != "Login required" || item.Code != "auth_required" {
		t.Fatalf("unexpected error payload: %#v", item)
	}
}
