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

func TestMapResources(t *testing.T) {
	items := []string{"a", "b"}
	resources := MapResources(items, func(item string) Resource {
		return Resource{
			Type: "letters",
			ID:   item,
		}
	})

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}
	if resources[1].ID != "b" {
		t.Fatalf("unexpected resource id: %s", resources[1].ID)
	}
}

func TestParseBookingDate(t *testing.T) {
	t.Run("empty defaults to today", func(t *testing.T) {
		date, err := ParseBookingDate("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if date == "" {
			t.Fatal("expected non-empty date")
		}
	})

	t.Run("valid date", func(t *testing.T) {
		date, err := ParseBookingDate("2025-12-25")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if date != "2025-12-25" {
			t.Fatalf("expected 2025-12-25, got %s", date)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := ParseBookingDate("not-a-date")
		if err == nil {
			t.Fatal("expected error for invalid date")
		}
	})
}

func TestParseRoomRequest(t *testing.T) {
	t.Run("valid params", func(t *testing.T) {
		params, err := ParseRoomRequest("room-1", "2025-01-15")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if params.RoomID != "room-1" {
			t.Fatalf("expected room-1, got %s", params.RoomID)
		}
		if params.BookingDate != "2025-01-15" {
			t.Fatalf("expected 2025-01-15, got %s", params.BookingDate)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := ParseRoomRequest("room-1", "invalid")
		if err == nil {
			t.Fatal("expected error for invalid date")
		}
	})
}
