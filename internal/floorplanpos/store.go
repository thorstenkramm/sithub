// Package floorplanpos provides storage and handlers for floor plan item positions.
package floorplanpos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound indicates the requested position does not exist.
var ErrNotFound = errors.New("position not found")

// Position represents an item's rectangle on a floor plan.
type Position struct {
	ID          string
	FloorPlan   string
	ItemID      string
	Label       string
	X           float64
	Y           float64
	Width       float64
	Height      float64
	BorderWidth int
	CreatedAt   string
	UpdatedAt   string
}

// CreateInput holds fields for creating a position.
type CreateInput struct {
	FloorPlan   string
	ItemID      string
	Label       string
	X           float64
	Y           float64
	Width       float64
	Height      float64
	BorderWidth int
}

// UpdateInput holds fields for updating a position.
type UpdateInput struct {
	Label       *string
	X           *float64
	Y           *float64
	Width       *float64
	Height      *float64
	BorderWidth *int
}

// Create inserts a new floor plan position.
func Create(ctx context.Context, db *sql.DB, input *CreateInput) (*Position, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.NewString()

	bw := input.BorderWidth
	if bw < 1 || bw > 5 {
		bw = 2
	}

	_, err := db.ExecContext(ctx,
		`INSERT INTO floor_plan_positions
		(id, floor_plan, item_id, label, x, y, width, height, border_width, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, input.FloorPlan, input.ItemID, input.Label,
		input.X, input.Y, input.Width, input.Height, bw,
		now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("create position: %w", err)
	}

	return &Position{
		ID: id, FloorPlan: input.FloorPlan, ItemID: input.ItemID,
		Label: input.Label, X: input.X, Y: input.Y,
		Width: input.Width, Height: input.Height, BorderWidth: bw,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

// FindByFloorPlan returns all positions for a given floor plan filename.
func FindByFloorPlan(ctx context.Context, db *sql.DB, floorPlan string) ([]Position, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, floor_plan, item_id, label, x, y, width, height, border_width, created_at, updated_at
		FROM floor_plan_positions WHERE floor_plan = ? ORDER BY item_id`,
		floorPlan,
	)
	if err != nil {
		return nil, fmt.Errorf("query positions: %w", err)
	}
	defer func() {
		_ = rows.Close() //nolint:errcheck // Defer close, error not critical
	}()

	var result []Position
	for rows.Next() {
		var p Position
		if err := rows.Scan(&p.ID, &p.FloorPlan, &p.ItemID, &p.Label,
			&p.X, &p.Y, &p.Width, &p.Height, &p.BorderWidth,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan position: %w", err)
		}
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate positions: %w", err)
	}
	return result, nil
}

// Update modifies an existing position.
func Update(ctx context.Context, db *sql.DB, id string, input UpdateInput) (*Position, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	setClauses := []string{"updated_at = ?"}
	args := []interface{}{now}

	if input.Label != nil {
		setClauses = append(setClauses, "label = ?")
		args = append(args, *input.Label)
	}
	if input.X != nil {
		setClauses = append(setClauses, "x = ?")
		args = append(args, *input.X)
	}
	if input.Y != nil {
		setClauses = append(setClauses, "y = ?")
		args = append(args, *input.Y)
	}
	if input.Width != nil {
		setClauses = append(setClauses, "width = ?")
		args = append(args, *input.Width)
	}
	if input.Height != nil {
		setClauses = append(setClauses, "height = ?")
		args = append(args, *input.Height)
	}
	if input.BorderWidth != nil {
		setClauses = append(setClauses, "border_width = ?")
		args = append(args, *input.BorderWidth)
	}

	args = append(args, id)

	// setClauses contains only hardcoded column names, not user input.
	query := "UPDATE floor_plan_positions SET " + //nolint:gosec // G202 false positive
		strings.Join(setClauses, ", ") + " WHERE id = ?"

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("update position: %w", err)
	}

	n, errRows := result.RowsAffected()
	if errRows != nil {
		n = 0
	}
	if n == 0 {
		return nil, ErrNotFound
	}

	return FindByID(ctx, db, id)
}

// Delete removes a position by ID.
func Delete(ctx context.Context, db *sql.DB, id string) error {
	result, err := db.ExecContext(ctx,
		`DELETE FROM floor_plan_positions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}
	n, errRows := result.RowsAffected()
	if errRows != nil {
		n = 0
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// FindByID returns a single position by ID.
func FindByID(ctx context.Context, db *sql.DB, id string) (*Position, error) {
	var p Position
	err := db.QueryRowContext(ctx,
		`SELECT id, floor_plan, item_id, label, x, y, width, height, border_width, created_at, updated_at
		FROM floor_plan_positions WHERE id = ?`, id,
	).Scan(&p.ID, &p.FloorPlan, &p.ItemID, &p.Label,
		&p.X, &p.Y, &p.Width, &p.Height, &p.BorderWidth,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("find position: %w", err)
	}
	return &p, nil
}
