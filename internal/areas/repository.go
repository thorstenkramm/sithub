// Package areas provides area data access.
package areas

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository provides access to area data.
// Repository is safe for concurrent use after construction.
type Repository struct {
	db *sql.DB
}

// Area represents a bookable area.
type Area struct {
	ID        string
	Name      string
	SortOrder int
	CreatedAt string
	UpdatedAt string
}

// NewRepository returns a new area repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List returns all areas ordered by sort order then name.
func (r *Repository) List(ctx context.Context) (_ []Area, err error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, sort_order, created_at, updated_at
		FROM areas
		ORDER BY sort_order ASC, name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list areas: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close areas rows: %w", closeErr)
		}
	}()

	areas := []Area{}
	for rows.Next() {
		var area Area
		if err := rows.Scan(&area.ID, &area.Name, &area.SortOrder, &area.CreatedAt, &area.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan area: %w", err)
		}
		areas = append(areas, area)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate areas: %w", err)
	}

	return areas, nil
}
