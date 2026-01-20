package spaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Store provides database operations for spaces.
type Store struct {
	db *sql.DB
}

// NewStore creates a new spaces store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// AreaRecord represents an area row in the database.
type AreaRecord struct {
	ID          string
	Name        string
	Description string
	FloorPlan   string
	CreatedAt   string
	UpdatedAt   string
}

// RoomRecord represents a room row in the database.
type RoomRecord struct {
	ID          string
	AreaID      string
	Name        string
	Description string
	FloorPlan   string
	CreatedAt   string
	UpdatedAt   string
}

// DeskRecord represents a desk row in the database.
type DeskRecord struct {
	ID        string
	RoomID    string
	Name      string
	Equipment []string
	Warning   string
	CreatedAt string
	UpdatedAt string
}

// ListAreas returns all areas.
func (s *Store) ListAreas(ctx context.Context) ([]AreaRecord, error) {
	query := `SELECT id, name, description, floor_plan, created_at, updated_at 
	          FROM areas ORDER BY name`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query areas: %w", err)
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck // Error not actionable in defer

	var areas []AreaRecord
	for rows.Next() {
		var a AreaRecord
		if err := rows.Scan(
			&a.ID, &a.Name, &a.Description, &a.FloorPlan, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan area: %w", err)
		}
		areas = append(areas, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate areas: %w", err)
	}
	return areas, nil
}

// GetArea returns an area by ID.
func (s *Store) GetArea(ctx context.Context, id string) (*AreaRecord, error) {
	var a AreaRecord
	query := `SELECT id, name, description, floor_plan, created_at, updated_at 
	          FROM areas WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.Name, &a.Description, &a.FloorPlan, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get area: %w", err)
	}
	return &a, nil
}

// CreateArea inserts a new area.
func (s *Store) CreateArea(ctx context.Context, a *AreaRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	a.CreatedAt = now
	a.UpdatedAt = now
	query := `INSERT INTO areas (id, name, description, floor_plan, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query,
		a.ID, a.Name, a.Description, a.FloorPlan, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert area: %w", err)
	}
	return nil
}

// UpdateArea updates an existing area.
//
//nolint:dupl // Update methods share structure but operate on different entities
func (s *Store) UpdateArea(ctx context.Context, a *AreaRecord) error {
	a.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE areas SET name = ?, description = ?, floor_plan = ?, updated_at = ? 
	          WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query,
		a.Name, a.Description, a.FloorPlan, a.UpdatedAt, a.ID)
	if err != nil {
		return fmt.Errorf("update area: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteArea removes an area by ID.
func (s *Store) DeleteArea(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM areas WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete area: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ListRooms returns all rooms for an area.
func (s *Store) ListRooms(ctx context.Context, areaID string) ([]RoomRecord, error) {
	query := `SELECT id, area_id, name, description, floor_plan, created_at, updated_at 
	          FROM rooms WHERE area_id = ? ORDER BY name`
	rows, err := s.db.QueryContext(ctx, query, areaID)
	if err != nil {
		return nil, fmt.Errorf("query rooms: %w", err)
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck // Error not actionable in defer

	var rooms []RoomRecord
	for rows.Next() {
		var r RoomRecord
		if err := rows.Scan(
			&r.ID, &r.AreaID, &r.Name, &r.Description, &r.FloorPlan, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan room: %w", err)
		}
		rooms = append(rooms, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rooms: %w", err)
	}
	return rooms, nil
}

// GetRoom returns a room by ID.
func (s *Store) GetRoom(ctx context.Context, id string) (*RoomRecord, error) {
	var r RoomRecord
	query := `SELECT id, area_id, name, description, floor_plan, created_at, updated_at 
	          FROM rooms WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&r.ID, &r.AreaID, &r.Name, &r.Description, &r.FloorPlan, &r.CreatedAt, &r.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get room: %w", err)
	}
	return &r, nil
}

// CreateRoom inserts a new room.
func (s *Store) CreateRoom(ctx context.Context, r *RoomRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	r.CreatedAt = now
	r.UpdatedAt = now
	query := `INSERT INTO rooms (id, area_id, name, description, floor_plan, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query,
		r.ID, r.AreaID, r.Name, r.Description, r.FloorPlan, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert room: %w", err)
	}
	return nil
}

// UpdateRoom updates an existing room.
//
//nolint:dupl // Update methods share structure but operate on different entities
func (s *Store) UpdateRoom(ctx context.Context, r *RoomRecord) error {
	r.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE rooms SET name = ?, description = ?, floor_plan = ?, updated_at = ? 
	          WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query,
		r.Name, r.Description, r.FloorPlan, r.UpdatedAt, r.ID)
	if err != nil {
		return fmt.Errorf("update room: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteRoom removes a room by ID.
func (s *Store) DeleteRoom(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM rooms WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete room: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ListDesks returns all desks for a room.
func (s *Store) ListDesks(ctx context.Context, roomID string) ([]DeskRecord, error) {
	query := `SELECT id, room_id, name, equipment, warning, created_at, updated_at 
	          FROM desks WHERE room_id = ? ORDER BY name`
	rows, err := s.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("query desks: %w", err)
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck // Error not actionable in defer

	var desks []DeskRecord
	for rows.Next() {
		var d DeskRecord
		var equipmentJSON string
		if err := rows.Scan(
			&d.ID, &d.RoomID, &d.Name, &equipmentJSON, &d.Warning, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan desk: %w", err)
		}
		if equipmentJSON != "" {
			if err := json.Unmarshal([]byte(equipmentJSON), &d.Equipment); err != nil {
				// Fallback: treat raw string as single equipment item
				d.Equipment = []string{equipmentJSON}
			}
		}
		desks = append(desks, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate desks: %w", err)
	}
	return desks, nil
}

// GetDesk returns a desk by ID.
func (s *Store) GetDesk(ctx context.Context, id string) (*DeskRecord, error) {
	var d DeskRecord
	var equipmentJSON string
	query := `SELECT id, room_id, name, equipment, warning, created_at, updated_at 
	          FROM desks WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.RoomID, &d.Name, &equipmentJSON, &d.Warning, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get desk: %w", err)
	}
	if equipmentJSON != "" {
		if err := json.Unmarshal([]byte(equipmentJSON), &d.Equipment); err != nil {
			d.Equipment = []string{equipmentJSON}
		}
	}
	return &d, nil
}

// CreateDesk inserts a new desk.
func (s *Store) CreateDesk(ctx context.Context, d *DeskRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	d.CreatedAt = now
	d.UpdatedAt = now
	equipmentJSON, err := json.Marshal(d.Equipment)
	if err != nil {
		return fmt.Errorf("marshal equipment: %w", err)
	}
	query := `INSERT INTO desks (id, room_id, name, equipment, warning, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query,
		d.ID, d.RoomID, d.Name, string(equipmentJSON), d.Warning, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert desk: %w", err)
	}
	return nil
}

// UpdateDesk updates an existing desk.
func (s *Store) UpdateDesk(ctx context.Context, d *DeskRecord) error {
	d.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	equipmentJSON, err := json.Marshal(d.Equipment)
	if err != nil {
		return fmt.Errorf("marshal equipment: %w", err)
	}
	query := `UPDATE desks SET name = ?, equipment = ?, warning = ?, updated_at = ? 
	          WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query,
		d.Name, string(equipmentJSON), d.Warning, d.UpdatedAt, d.ID)
	if err != nil {
		return fmt.Errorf("update desk: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteDesk removes a desk by ID.
func (s *Store) DeleteDesk(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM desks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete desk: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// SyncFromConfig synchronizes the database with a YAML config.
// This is called on startup to seed the database.
func (s *Store) SyncFromConfig(ctx context.Context, cfg *Config) error {
	for _, area := range cfg.Areas {
		if err := s.syncArea(ctx, &area); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) syncArea(ctx context.Context, area *Area) error {
	existing, err := s.GetArea(ctx, area.ID)
	if err != nil {
		return fmt.Errorf("check area %s: %w", area.ID, err)
	}
	if existing == nil {
		if err := s.CreateArea(ctx, &AreaRecord{
			ID:          area.ID,
			Name:        area.Name,
			Description: area.Description,
			FloorPlan:   area.FloorPlan,
		}); err != nil {
			return fmt.Errorf("create area %s: %w", area.ID, err)
		}
	}

	for i := range area.Rooms {
		if err := s.syncRoom(ctx, area.ID, &area.Rooms[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) syncRoom(ctx context.Context, areaID string, room *Room) error {
	existing, err := s.GetRoom(ctx, room.ID)
	if err != nil {
		return fmt.Errorf("check room %s: %w", room.ID, err)
	}
	if existing == nil {
		if err := s.CreateRoom(ctx, &RoomRecord{
			ID:          room.ID,
			AreaID:      areaID,
			Name:        room.Name,
			Description: room.Description,
			FloorPlan:   room.FloorPlan,
		}); err != nil {
			return fmt.Errorf("create room %s: %w", room.ID, err)
		}
	}

	for i := range room.Desks {
		if err := s.syncDesk(ctx, room.ID, &room.Desks[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) syncDesk(ctx context.Context, roomID string, desk *Desk) error {
	existing, err := s.GetDesk(ctx, desk.ID)
	if err != nil {
		return fmt.Errorf("check desk %s: %w", desk.ID, err)
	}
	if existing == nil {
		if err := s.CreateDesk(ctx, &DeskRecord{
			ID:        desk.ID,
			RoomID:    roomID,
			Name:      desk.Name,
			Equipment: desk.Equipment,
			Warning:   desk.Warning,
		}); err != nil {
			return fmt.Errorf("create desk %s: %w", desk.ID, err)
		}
	}
	return nil
}

// LoadConfig builds a Config from the database.
func (s *Store) LoadConfig(ctx context.Context) (*Config, error) {
	areas, err := s.ListAreas(ctx)
	if err != nil {
		return nil, fmt.Errorf("list areas: %w", err)
	}

	cfg := &Config{}
	for _, a := range areas {
		area, err := s.loadArea(ctx, &a)
		if err != nil {
			return nil, err
		}
		cfg.Areas = append(cfg.Areas, *area)
	}

	return cfg, nil
}

func (s *Store) loadArea(ctx context.Context, a *AreaRecord) (*Area, error) {
	area := &Area{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		FloorPlan:   a.FloorPlan,
	}

	rooms, err := s.ListRooms(ctx, a.ID)
	if err != nil {
		return nil, fmt.Errorf("list rooms for area %s: %w", a.ID, err)
	}

	for _, r := range rooms {
		room, err := s.loadRoom(ctx, &r)
		if err != nil {
			return nil, err
		}
		area.Rooms = append(area.Rooms, *room)
	}

	return area, nil
}

func (s *Store) loadRoom(ctx context.Context, r *RoomRecord) (*Room, error) {
	room := &Room{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		FloorPlan:   r.FloorPlan,
	}

	desks, err := s.ListDesks(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("list desks for room %s: %w", r.ID, err)
	}

	for _, d := range desks {
		room.Desks = append(room.Desks, Desk{
			ID:        d.ID,
			Name:      d.Name,
			Equipment: d.Equipment,
			Warning:   d.Warning,
		})
	}

	return room, nil
}
