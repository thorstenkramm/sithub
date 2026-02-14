# Story 12.1: Rename Backend Internals and Database Column

Status: done

## Story

As a developer,
I want the Go packages, structs, and database column to use domain-neutral terminology,
So that the codebase foundation is ready for the public API rename.

## Acceptance Criteria

1. **Given** the Go packages `internal/rooms/` and `internal/desks/` exist
   **When** the rename is applied
   **Then** they are consolidated or renamed to use "item group" and "item" terminology
   **And** all Go struct fields, function names, and variables use the new terminology

2. **Given** the database has `desk_id` columns in the bookings table
   **When** a new migration is applied
   **Then** the column is renamed to `item_id`
   **And** the unique constraint on (desk_id, booking_date) becomes (item_id, booking_date)
   **And** existing bookings are preserved with correct references

3. **Given** the API routes remain unchanged in this story
   **When** the internal rename is applied
   **Then** the existing API routes still work (backward compatible)
   **And** `go test -race ./...` succeeds
   **And** `golangci-lint run ./...` reports no errors

## Tasks / Subtasks

- [x] Rename `internal/rooms/` to `internal/itemgroups/` (AC: 1)
  - [x] Rename handler functions and types
  - [x] Update all references in `internal/startup/server.go`
- [x] Rename `internal/desks/` to `internal/items/` (AC: 1)
  - [x] Rename handler functions, types, and query strings
- [x] Create migration 000009 to rename desk_id to item_id (AC: 2)
  - [x] ALTER TABLE bookings RENAME COLUMN desk_id TO item_id
  - [x] Recreate unique constraint
- [x] Update `internal/bookings/` to use item_id terminology (AC: 1)
  - [x] Rename struct fields, function parameters, SQL queries
- [x] Update `internal/areas/presence_handler.go` (AC: 1)
- [x] Update `internal/notifications/notifier.go` (AC: 1)
- [x] Run `go test -race ./...` and `golangci-lint run ./...` (AC: 3)

## Dev Notes

### Design Decisions

- `internal/rooms/` package split into `internal/itemgroups/` (list handler + bookings handler)
- `internal/desks/` became `internal/items/` with updated query column references
- Migration uses ALTER TABLE RENAME COLUMN (SQLite 3.25+)

### References

- PRD FR4-FR16 (reworded): `_bmad-output/planning-artifacts/prd.md`
- Epic Story 12.1: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- Go packages renamed: rooms -> itemgroups, desks -> items
- Migration 000009 renames desk_id to item_id in bookings table
- All internal references updated across bookings, areas, notifications packages

### Key Files

- `internal/itemgroups/handler.go` (was rooms/handler.go)
- `internal/itemgroups/bookings_handler.go` (was rooms/bookings_handler.go)
- `internal/items/handler.go` (was desks/handler.go)
- `migrations/000009_rename_desk_id_to_item_id.up.sql`
- `migrations/000009_rename_desk_id_to_item_id.down.sql`
- `internal/bookings/handler.go` (desk_id -> item_id)
- `internal/bookings/store.go` (desk_id -> item_id)
- `internal/startup/server.go` (package imports)

### Change Log

- 2026-02-09: Story created retroactively. Implementation was part of Epic 12 commit.
