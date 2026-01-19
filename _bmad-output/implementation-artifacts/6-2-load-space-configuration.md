# Story 6.2: Load Space Configuration

Status: complete

## Story

As an operator,
I want areas, rooms, desks, and equipment loaded from a config file,
So that space definitions are centrally managed.

## Acceptance Criteria

1. **Given** a valid space configuration file  
   **When** the server starts  
   **Then** the UI reflects the configured areas, rooms, desks, and equipment

## Tasks / Subtasks

- [x] Server loads YAML space configuration file (AC: 1)
- [x] Areas with id, name, description, floor_plan (AC: 1)
- [x] Rooms within areas with id, name, description (AC: 1)
- [x] Desks within rooms with id, name, equipment, warning (AC: 1)
- [x] Validation ensures required fields present (AC: 1)
- [x] API endpoints return configured spaces (AC: 1)

## Dev Notes

This functionality was implemented as foundational infrastructure work.

### Implementation Details

- Space configuration loaded via `spaces.Load()` function
- YAML format for nested structure readability
- Validation ensures all areas, rooms, desks have id and name
- Path configured via `spaces.config_file` in main config

### File List

**Existing files (implemented previously):**
- `internal/spaces/config.go` - Space configuration loading and validation
- `internal/spaces/config_test.go` - Tests
- `spaces.yaml.example` - Example space configuration

### Change Log

- Pre-sprint: Implemented as foundational work
- 2026-01-19: Story documented as complete (functionality already exists)
