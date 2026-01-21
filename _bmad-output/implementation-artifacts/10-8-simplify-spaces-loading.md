# Story 10.8: Simplify Spaces Loading

## Story

**As a** system administrator,  
**I want** spaces (areas, rooms, desks) to be loaded directly from the YAML config file,  
**So that** configuration is simple, predictable, and doesn't require database management.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** High (fixes confusing behavior)

## Context

The current implementation syncs spaces from YAML to SQLite database on startup. This causes confusion because:
- Old data persists in the database even when removed from YAML
- Users see spaces that don't exist in their config file
- The sync only adds new items, never removes or updates existing ones

The new approach: **YAML is the single source of truth**. Spaces are loaded into memory at startup. No database tables for spaces.

## Acceptance Criteria

**AC1: Spaces Loaded from YAML Only**
- **Given** the server starts with a valid `spaces.config_file` in config
- **When** I call `GET /api/v1/areas`
- **Then** I see exactly the areas defined in the YAML file
- **And** no more, no less

**AC2: No Database Persistence for Spaces**
- **Given** the server is running
- **When** I inspect the SQLite database
- **Then** there are no `areas`, `rooms`, or `desks` tables

**AC3: Bookings Reference Desk IDs as Strings**
- **Given** a booking exists for `desk_id = "desk_101_1"`
- **When** I query bookings
- **Then** the desk_id is stored as a plain string (no foreign key)

**AC4: Booking Validation Against In-Memory Config**
- **Given** I try to book a desk with `desk_id = "nonexistent"`
- **When** the booking request is processed
- **Then** I receive a 404 error "desk not found"

**AC5: Config Changes Require Restart**
- **Given** the server is running
- **When** I modify the YAML file
- **Then** changes are NOT reflected until server restart

## Technical Requirements

### Database Schema Changes

**Remove tables:**
- `areas`
- `rooms`
- `desks`

**Modify `bookings` table:**
- Remove foreign key constraint on `desk_id`
- Keep `desk_id` as `TEXT NOT NULL`

### Code Changes

**Remove:**
- `internal/spaces/store.go` - all CRUD and sync methods for areas/rooms/desks
- `internal/admin/handlers.go` - space management endpoints (CreateArea, UpdateArea, DeleteArea, CreateRoom, UpdateRoom, DeleteRoom, CreateDesk, UpdateDesk, DeleteDesk)
- `internal/admin/config_holder.go` - no longer needed
- Related test files

**Modify:**
- `internal/spaces/config.go` - keep `LoadFromFile()`, remove DB-related code
- `internal/startup/server.go` - load config into memory, remove sync calls
- `internal/api/handlers.go` - use in-memory config for listing areas/rooms/desks
- `internal/booking/` - validate desk_id against in-memory config, not DB

**Keep:**
- `internal/spaces/config.go` - Config, Area, Room, Desk structs and YAML loading
- Booking store and handlers (modified to use string desk_id)

### API Changes

**Remove endpoints:**
- `POST /api/v1/admin/areas`
- `PUT /api/v1/admin/areas/:id`
- `DELETE /api/v1/admin/areas/:id`
- `POST /api/v1/admin/areas/:areaId/rooms`
- `PUT /api/v1/admin/rooms/:id`
- `DELETE /api/v1/admin/rooms/:id`
- `POST /api/v1/admin/rooms/:roomId/desks`
- `PUT /api/v1/admin/desks/:id`
- `DELETE /api/v1/admin/desks/:id`

**Keep endpoints (unchanged behavior):**
- `GET /api/v1/areas`
- `GET /api/v1/areas/:id/rooms`
- `GET /api/v1/rooms/:id/desks`

## Tasks

### Task 1: Create New Database Migration
- [ ] Create migration to drop `areas`, `rooms`, `desks` tables
- [ ] Modify `bookings` table to remove FK constraint on `desk_id`
- [ ] Update `internal/db/schema.go`

### Task 2: Simplify Spaces Package
- [ ] Keep only `config.go` with YAML loading
- [ ] Remove `store.go` (or gut it completely)
- [ ] Update `Config` struct to have lookup methods: `GetArea(id)`, `GetRoom(id)`, `GetDesk(id)`
- [ ] Update tests

### Task 3: Update Startup
- [ ] Modify `internal/startup/server.go`
- [ ] Load spaces config into memory once
- [ ] Remove `SyncFromConfig` call
- [ ] Pass config (or getter) to handlers that need it

### Task 4: Update API Handlers
- [ ] Modify `ListAreasHandler` to use in-memory config
- [ ] Modify `ListRoomsHandler` to use in-memory config
- [ ] Modify `ListDesksHandler` to use in-memory config
- [ ] Remove ConfigHolder dependency

### Task 5: Remove Admin Space Endpoints
- [ ] Remove admin handlers for spaces CRUD
- [ ] Remove routes from `server.go`
- [ ] Delete `internal/admin/config_holder.go`
- [ ] Update admin handler tests

### Task 6: Update Booking Validation
- [ ] Modify booking creation to validate desk_id against in-memory config
- [ ] Return 404 if desk doesn't exist in config
- [ ] Update booking tests

### Task 7: Clean Up Tests
- [ ] Remove tests for deleted functionality
- [ ] Update integration tests
- [ ] Ensure all tests pass

### Task 8: Update Documentation
- [ ] Update API docs (remove admin space endpoints)
- [ ] Document that spaces require restart to reload

## File Changes

| Action | File Path |
|--------|-----------|
| Modify | `internal/db/schema.go` |
| Modify | `internal/spaces/config.go` |
| Delete | `internal/spaces/store.go` |
| Delete | `internal/spaces/store_test.go` |
| Delete | `internal/admin/config_holder.go` |
| Delete | `internal/admin/config_holder_test.go` |
| Modify | `internal/admin/handlers.go` |
| Modify | `internal/admin/handlers_test.go` |
| Modify | `internal/startup/server.go` |
| Modify | `internal/startup/server_test.go` |
| Modify | `internal/api/handlers.go` |
| Modify | `internal/api/handlers_test.go` |
| Modify | `internal/booking/handlers.go` |
| Modify | `internal/booking/handlers_test.go` |
| Modify | `docs/api.md` |

## Definition of Done

- [ ] `areas`, `rooms`, `desks` tables removed from schema
- [ ] Spaces loaded from YAML into memory at startup
- [ ] No admin CRUD endpoints for spaces
- [ ] Bookings store desk_id as plain string
- [ ] Booking creation validates desk exists in config
- [ ] All existing API endpoints work correctly
- [ ] All tests pass
- [ ] No references to removed code remain

## Notes

- Delete existing `sithub.db` file after this change (schema incompatible)
- This is a breaking change - document in release notes
- Future enhancement: hot reload could be added later if needed

## Dependencies

- None

## Blocked By

- None

## Blocks

- None (can be done in parallel with other Epic 10 stories)
