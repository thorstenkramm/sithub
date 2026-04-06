package areas

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateReservationsValid(t *testing.T) {
	t.Parallel()
	cfg := &Config{
		Areas: []Area{
			{
				ID:          "area-1",
				Name:        "Office",
				ReservedFor: []string{"a@test.com", "b@test.com"},
				ItemGroups: []ItemGroup{
					{
						ID:          "room-1",
						Name:        "Room 1",
						ReservedFor: []string{"a@test.com"},
						Items: []Item{
							{ID: "desk-1", Name: "Desk 1"},
						},
					},
				},
			},
		},
	}
	assert.NoError(t, ValidateReservations(cfg))
}

func TestValidateReservationsConflict(t *testing.T) {
	t.Parallel()
	cfg := &Config{
		Areas: []Area{
			{
				ID:          "area-1",
				Name:        "Office",
				ReservedFor: []string{"a@test.com"},
				ItemGroups: []ItemGroup{
					{
						ID:          "room-1",
						Name:        "Room 1",
						ReservedFor: []string{"b@test.com"}, // b not in area
						Items:       []Item{{ID: "desk-1", Name: "Desk 1"}},
					},
				},
			},
		},
	}
	err := ValidateReservations(cfg)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrReservationConflict))
	assert.Contains(t, err.Error(), "b@test.com")
}

func TestValidateReservationsNoRestriction(t *testing.T) {
	t.Parallel()
	cfg := &Config{
		Areas: []Area{
			{
				ID:   "area-1",
				Name: "Office",
				ItemGroups: []ItemGroup{
					{
						ID:    "room-1",
						Name:  "Room 1",
						Items: []Item{{ID: "desk-1", Name: "Desk 1"}},
					},
				},
			},
		},
	}
	assert.NoError(t, ValidateReservations(cfg))
}

func TestIsReservedItemLevel(t *testing.T) {
	t.Parallel()
	loc := &ItemLocation{
		Area:      &Area{ID: "a"},
		ItemGroup: &ItemGroup{ID: "ig"},
		Item:      &Item{ID: "i", ReservedFor: []string{"allowed@test.com"}},
	}
	assert.True(t, IsReserved(loc, "denied@test.com"))
	assert.False(t, IsReserved(loc, "allowed@test.com"))
}

func TestIsReservedAreaLevel(t *testing.T) {
	t.Parallel()
	loc := &ItemLocation{
		Area:      &Area{ID: "a", ReservedFor: []string{"allowed@test.com"}},
		ItemGroup: &ItemGroup{ID: "ig"},
		Item:      &Item{ID: "i"},
	}
	assert.True(t, IsReserved(loc, "denied@test.com"))
	assert.False(t, IsReserved(loc, "allowed@test.com"))
}

func TestIsReservedNoRestriction(t *testing.T) {
	t.Parallel()
	loc := &ItemLocation{
		Area:      &Area{ID: "a"},
		ItemGroup: &ItemGroup{ID: "ig"},
		Item:      &Item{ID: "i"},
	}
	assert.False(t, IsReserved(loc, "anyone@test.com"))
}
