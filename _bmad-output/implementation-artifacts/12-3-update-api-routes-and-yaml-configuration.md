# Story 12.3: Update API Routes and YAML Configuration

Status: done

## Story

As an operator,
I want the API routes and YAML configuration to use domain-neutral terminology,
So that the public interface reflects the flexible item model.

## Acceptance Criteria

1. **Given** the API routes use `/rooms/:room_id/desks` and `/areas/:area_id/rooms`
   **When** the route rename is applied
   **Then** the routes use `/item-groups/:item_group_id/items` and
   `/areas/:area_id/item-groups`
   **And** JSON:API resource types use the new terminology (e.g., `item-groups`, `items`)

2. **Given** the YAML configuration uses `rooms` and `desks` keys
   **When** the spaces config loader is updated
   **Then** it reads `items` keys at both hierarchy levels (item groups and items)
   **And** the `sithub_areas.example.yaml` is updated with the new keys
   **And** the example includes diverse item types (office desks, parking lots)

3. **Given** the route rename is applied
   **When** `go test -race ./...` is executed
   **Then** all tests pass with the new route paths
   **And** `golangci-lint run ./...` reports no errors

## Tasks / Subtasks

- [x] Update API routes in `internal/startup/server.go` (AC: 1)
  - [x] `/areas/:area_id/rooms` -> `/areas/:area_id/item-groups`
  - [x] `/rooms/:room_id/desks` -> `/item-groups/:item_group_id/items`
  - [x] `/rooms/:room_id/bookings` -> `/item-groups/:item_group_id/bookings`
- [x] Update JSON:API resource types (AC: 1)
  - [x] "rooms" -> "item-groups"
  - [x] "desks" -> "items"
- [x] Update `internal/spaces/config.go` YAML parsing (AC: 2)
  - [x] Parse `items` key at both levels instead of `rooms`/`desks`
  - [x] Update struct field names and tags
- [x] Update `sithub_areas.example.yaml` (AC: 2)
  - [x] Use `items` keys throughout
  - [x] Add diverse examples (parking lots, lab benches)
- [x] Update all handler and store tests (AC: 3)

## Dev Notes

### Design Decisions

- YAML config uses `items` at both hierarchy levels (flat key name)
- Example config demonstrates booking diverse resource types
- API resource types use kebab-case per JSON:API convention

### References

- PRD FR4-FR16, FR18 (reworded): `_bmad-output/planning-artifacts/prd.md`
- Epic Story 12.3: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- All API routes updated to use item-groups/items terminology
- JSON:API resource types updated to "item-groups" and "items"
- YAML config parser reads "items" key at both levels
- Example config shows office desks, parking lots, and meeting rooms

### Key Files

- `internal/startup/server.go` (route definitions)
- `internal/spaces/config.go` (YAML parsing)
- `internal/spaces/config_test.go`
- `internal/api/response.go` (resource type constants)
- `sithub_areas.example.yaml`

### Change Log

- 2026-02-09: Story created retroactively. Implementation was part of Epic 12 commit.
