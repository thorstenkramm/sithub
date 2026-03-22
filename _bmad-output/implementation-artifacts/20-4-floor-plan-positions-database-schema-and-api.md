# Story 20.4: Floor Plan Positions Database Schema and API

Status: backlog

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

- [ ] Create database migration for `floor_plan_positions` table (AC: 1)
  - [ ] Add migration file with columns: id, floor_plan, item_id, label, x, y, width,
    height, created_at, updated_at
  - [ ] `label` is optional — when set, the floor plan viewer shows this text instead of
    deriving it from the item name (allows admins to use short labels like "T1")
  - [ ] Include appropriate indexes (floor_plan, item_id)
- [ ] Create store functions for floor plan positions (AC: 1, 2, 3, 4)
  - [ ] `Create(ctx, db, position)` to insert a new position
  - [ ] `FindByFloorPlan(ctx, db, floorPlan)` to list all positions for a floor plan
  - [ ] `Update(ctx, db, id, position)` to update a position
  - [ ] `Delete(ctx, db, id)` to remove a position
- [ ] Create API handlers (AC: 2, 3, 4)
  - [ ] GET `/api/v1/floor-plan-positions?floor_plan=<filename>` returns JSON:API collection
  - [ ] POST `/api/v1/floor-plan-positions` creates a new position (admin only)
  - [ ] PUT `/api/v1/floor-plan-positions/:id` updates a position (admin only)
  - [ ] DELETE `/api/v1/floor-plan-positions/:id` removes a position (admin only)
- [ ] Add routes to Echo router
- [ ] Write Go tests for store functions and handlers
  - [ ] Table-driven tests for CRUD operations
  - [ ] Test admin-only access on write endpoints
- [ ] Update API documentation
  - [ ] Add `floor-plan-positions.yaml` to `api-doc/`
  - [ ] Reference from `openapi.yaml`
  - [ ] Lint with `npx @redocly/cli lint`
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

### Completion Notes List

### File List

## Change Log

- 2026-03-22: UX review — added optional `label` column to schema for admin-controlled
  display text on floor plan rectangles.
