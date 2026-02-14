# Story 12.2: Normalize Bookings Table

Status: done

## Story

As a developer,
I want the bookings table to reference users by user_id only,
So that display names are always current and not duplicated across tables.

## Acceptance Criteria

1. **Given** the bookings table contains `user_name`, `booked_by_user_id`, and
   `booked_by_user_name` columns
   **When** a new database migration is applied
   **Then** the redundant columns are removed from the bookings table
   **And** existing data is preserved (user_id remains as the foreign key reference)

2. **Given** a booking query needs to display a user's name
   **When** the query is executed
   **Then** the display name is resolved via JOIN with the users table
   **And** the name reflects the current value in the users table

3. **Given** the bookings API returns user information
   **When** a booking list or detail is requested
   **Then** user display names are included in the response via JOIN
   **And** the JSON:API response structure remains consistent

4. **Given** the migration has been applied
   **When** `go test -race ./...` is executed
   **Then** all existing tests pass with the normalized schema
   **And** booking creation and listing continue to work correctly

## Tasks / Subtasks

- [x] Create migration 000010 to normalize bookings table (AC: 1)
  - [x] Copy booked_by_user_id data to preserve on-behalf relationships
  - [x] Recreate table without user_name and booked_by_user_name columns
  - [x] Restore data from backup
- [x] Update `internal/bookings/store.go` queries to JOIN users table (AC: 2, 3)
  - [x] ListUserBookings: JOIN for user names
  - [x] FindBookingByID: JOIN for user names
  - [x] All booking list queries resolve names at query time
- [x] Add `internal/users/store.go` FindDisplayName function (AC: 2)
- [x] Create migration 000011 for performance indexes (AC: 4)
  - [x] Index on bookings(user_id)
  - [x] Index on bookings(booked_by_user_id)
  - [x] Index on bookings(item_id, booking_date)
- [x] Update all booking tests for normalized schema (AC: 4)

## Dev Notes

### Migration Strategy

Migration 000010 uses SQLite's table recreation pattern since ALTER TABLE DROP COLUMN
has limitations. Steps: create backup, drop original, create new table, restore data.

### Design Decisions

- User names are never stored in bookings; always resolved via JOIN
- `users.FindDisplayName` added for cases where a single name lookup is needed
- Performance indexes added to support JOIN queries efficiently

### Code Review Findings

- Code review identified that `for_user_id` was not validated against the users table.
  Fixed by adding `users.FindByID` validation in `resolveBookingParticipants`.

### References

- PRD FR41: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 12.2: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- Migration 000010 drops user_name and booked_by_user_name columns
- Migration 000011 adds performance indexes for JOIN queries
- All booking queries now resolve names via LEFT JOIN on users table
- `for_user_id` validated against users table before booking creation

### Key Files

- `migrations/000010_normalize_bookings.up.sql`
- `migrations/000010_normalize_bookings.down.sql`
- `migrations/000011_add_bookings_indexes.up.sql`
- `migrations/000011_add_bookings_indexes.down.sql`
- `internal/bookings/store.go` (JOIN queries)
- `internal/bookings/handler.go` (for_user_id validation)
- `internal/users/store.go` (FindDisplayName)

### Change Log

- 2026-02-09: Story created retroactively. Implementation was part of Epic 12 commit.
