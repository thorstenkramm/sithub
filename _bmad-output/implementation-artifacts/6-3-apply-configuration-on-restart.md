# Story 6.3: Apply Configuration on Restart

Status: complete

## Story

As an operator,
I want configuration changes to apply on restart,
So that I can update spaces without manual migration steps.

## Acceptance Criteria

1. **Given** the config file has changed  
   **When** the server restarts  
   **Then** the new configuration is applied  
   **And** no manual data migration steps are required

## Tasks / Subtasks

- [x] Server reloads config files on startup (AC: 1)
- [x] Space changes reflected in API after restart (AC: 1)
- [x] No manual migration required for space changes (AC: 1)
- [x] Bookings reference desk_id which remains stable (AC: 1)

## Dev Notes

This functionality is inherent to the design - configuration is loaded fresh on each server start.

### Implementation Details

- Server loads configuration in `startup.Run()` on each start
- Space configuration is read from YAML file each time
- Bookings stored in SQLite reference desk_id (stable identifier)
- Operators can add/remove/modify spaces by editing config and restarting

### Design Decisions

- Configuration is file-based, not database-stored
- Desk IDs should remain stable to preserve booking references
- Removing a desk with existing bookings doesn't break the database
  (bookings just reference a desk that no longer exists in config)

### File List

**Existing files (implemented previously):**
- `internal/startup/server.go` - Loads config and spaces on startup
- `internal/config/config.go` - Server config loading
- `internal/spaces/config.go` - Space config loading

### Change Log

- Pre-sprint: Implemented as foundational work
- 2026-01-19: Story documented as complete (functionality already exists)
