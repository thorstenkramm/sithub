# Story 2.3: List Desks with Equipment

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an employee,
I want to see desks and their equipment for a room,
so that I can pick a suitable desk.

## Acceptance Criteria

1. **Given** I am viewing a room  
   **When** I open the room  
   **Then** I see the list of desks in that room  
   **And** each desk shows its equipment list

## Tasks / Subtasks

- [x] Add desks API endpoint (AC: 1)
  - [x] Implement `GET /api/v1/rooms/{room_id}/desks` returning JSON:API collection
  - [x] Return 404 JSON:API error when the room does not exist
  - [x] Enforce auth middleware and JSON:API error handling
- [x] Add frontend desks list view (AC: 1)
  - [x] Add route for room desks
  - [x] Fetch desks for selected room and render list with equipment
  - [x] Link from rooms list to the selected room
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for room filtering and not-found
  - [x] Frontend: unit test for list rendering and navigation
  - [x] Cypress E2E: selecting a room shows its desks and equipment
- [x] Update API documentation (AC: 1)
  - [x] Add desk resource schema and update endpoint response

## Dev Notes

- Desk data comes from YAML config via `spaces.config_file`.
- Use JSON:API envelopes with `application/vnd.api+json` content type.
- Use `internal/desks` for handler; keep handler free of storage concerns.
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend: `internal/desks`, `internal/api`, `internal/startup`, `internal/spaces`.
- Frontend: `web/src/views`, `web/src/router`, `web/src/api`.

### References

- PRD FR6/FR7: `_bmad-output/planning-artifacts/prd.md` (Areas, Rooms, and Desks Discovery)
- Epic Story 2.3: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (JSON:API, naming)
- API docs rules: `_bmad-output/planning-artifacts/architecture.md` (OpenAPI 3.1 in `api-doc/`)

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Added desks API endpoint with room filtering, JSON:API responses, and 404 handling.
- Added desks list view and routing, plus room link navigation and equipment display.
- Added backend, frontend, and Cypress coverage for desks listing.
- Updated OpenAPI schema/endpoint documentation for desks collection.

### File List

- `internal/desks/handler.go`
- `internal/desks/handler_test.go`
- `internal/spaces/config.go`
- `internal/spaces/config_test.go`
- `internal/startup/server.go`
- `internal/api/response.go`
- `internal/api/write.go`
- `api-doc/openapi.yaml`
- `api-doc/endpoints/desks.yaml`
- `web/src/api/desks.ts`
- `web/src/router/index.ts`
- `web/src/router/index.test.ts`
- `web/src/views/RoomsView.vue`
- `web/src/views/RoomsView.test.ts`
- `web/src/views/AreasView.test.ts`
- `web/src/views/DesksView.vue`
- `web/src/views/DesksView.test.ts`
- `web/src/views/testHelpers.ts`
- `web/cypress/e2e/rooms.cy.ts`
- `web/cypress/e2e/desks.cy.ts`
- `web/cypress/support/flows.ts`

### Change Log

- 2026-01-18: Story created and set to in-progress.
- 2026-01-18: Story completed.
