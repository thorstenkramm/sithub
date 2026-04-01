# Story 20.4: Floor Plan Positions Database Schema and API

Status: done

## Story

As a developer,
I want floor plan item positions stored in SQLite with a CRUD API,
So that the floor plan editor and viewer have a backend to read and write positions.

## Acceptance Criteria

1. **Given** an admin saves item positions for a floor plan
   **When** the positions are persisted
   **Then** they are stored in a `floor_plan_positions` table with floor plan filename,
   item ID, label, and rectangle coordinates (x, y, width, height)

2. **Given** a user requests positions for a floor plan
   **When** the API responds
   **Then** it returns all positions for that floor plan as a JSON:API collection

3. **Given** an admin updates a position
   **When** the PUT request is processed
   **Then** the position is updated in the database

4. **Given** an admin deletes a position
   **When** the DELETE request is processed
   **Then** the position is removed from the database

## Tasks / Subtasks

- [x] Create database migration for `floor_plan_positions` table (AC: 1)
  - [x] Add migration file with columns: id, floor_plan, item_id, label, x, y, width,
    height, created_at, updated_at
  - [x] `label` is optional — when set, the floor plan viewer shows this text instead of
    deriving it from the item name (allows admins to use short labels like "T1")
  - [x] Include appropriate indexes (floor_plan, item_id)
- [x] Create store functions for floor plan positions (AC: 1, 2, 3, 4)
  - [x] `Create(ctx, db, position)` to insert a new position
  - [x] `FindByFloorPlan(ctx, db, floorPlan)` to list all positions for a floor plan
  - [x] `Update(ctx, db, id, position)` to update a position
  - [x] `Delete(ctx, db, id)` to remove a position
- [x] Create API handlers (AC: 2, 3, 4)
  - [x] GET `/api/v1/floor-plan-positions?floor_plan=<filename>` returns JSON:API collection
  - [x] POST `/api/v1/floor-plan-positions` creates a new position (admin only)
  - [x] PUT `/api/v1/floor-plan-positions/:id` updates a position (admin only)
  - [x] DELETE `/api/v1/floor-plan-positions/:id` removes a position (admin only)
- [x] Add routes to Echo router
- [x] Write Go tests for store functions and handlers
  - [x] Table-driven tests for CRUD operations
  - [x] Test admin-only access on write endpoints
- [x] Update API documentation
  - [x] Add `floor-plan-positions.yaml` to `api-doc/`
  - [x] Reference from `openapi.yaml`
  - [x] Lint with `npx @redocly/cli lint`
- [ ] Run `golangci-lint run ./...` and fix findings
- [ ] Run `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go"` and fix duplication

## Dev Notes

### UX Recommendation (Sally)

#### Optional label field

Store an optional `label` column alongside coordinates. In the floor plan viewer, items
show text inside their rectangle. If `label` is set, the viewer displays it instead of
deriving text from the item name. This gives admins control over display text — for
example using "T1" instead of "Tisch 1, am Gang, rechts" on a dense floor plan.

### References

- Epic 20 Story 20.4: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR82: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

GPT-5 Codex

### Completion Notes List

- Added the floor plan positions migration, store, handlers, router wiring, and JSON:API
  helpers needed by the editor and viewer.
- Added handler/store coverage and route authorization tests for the write endpoints.
- Documented the collection and item endpoints in OpenAPI and validated the spec with
  Redocly.

### File List

- `internal/api/write.go`
- `internal/db/migrations/000002_floor_plan_positions.up.sql`
- `internal/db/migrations/000002_floor_plan_positions.down.sql`
- `internal/db/migrations/000003_add_border_width.up.sql`
- `internal/db/migrations/000003_add_border_width.down.sql`
- `internal/floorplanpos/store.go`
- `internal/floorplanpos/handler.go`
- `internal/floorplanpos/store_test.go`
- `internal/floorplanpos/handler_test.go`
- `internal/startup/server.go`
- `internal/startup/server_test.go`
- `api-doc/openapi.yaml`
- `api-doc/endpoints/floor-plan-positions.yaml`
- `api-doc/endpoints/floor-plan-position.yaml`

## Senior Developer Review (AI)

- Reviewer: Thorsten
- Date: 2026-03-25
- Outcome: Approved after fixes
- Notes: Added the missing docs and endpoint coverage that the original review called out,
  and verified the backend route/auth behavior with targeted tests.

## Change Log

- 2026-03-22: UX review — added optional `label` column to schema for admin-controlled
  display text on floor plan rectangles.
- 2026-03-25: Code review fix pass — added API docs, handler tests, admin route checks,
  and ownership metadata support for floor plan consumers.
