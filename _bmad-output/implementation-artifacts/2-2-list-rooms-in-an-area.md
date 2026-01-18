# Story 2.2: List Rooms in an Area

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an employee,
I want to see rooms for a selected area,
so that I can choose a room.

## Acceptance Criteria

1. **Given** I am viewing an area  
   **When** I select the area  
   **Then** I see only rooms belonging to that area  
   **And** rooms outside the area are not shown

## Tasks / Subtasks

- [x] Add rooms API endpoint (AC: 1)
  - [x] Implement `GET /api/v1/areas/{area_id}/rooms` returning JSON:API collection
  - [x] Return 404 JSON:API error when the area does not exist
  - [x] Enforce auth middleware and JSON:API error handling
- [x] Add frontend room list view (AC: 1)
  - [x] Add route for area rooms
  - [x] Fetch rooms for selected area and render list
  - [x] Link from areas list to the selected area
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for area filtering and not-found
  - [x] Frontend: unit test for list rendering and navigation
  - [x] Cypress E2E: selecting an area shows only its rooms
- [x] Update API documentation (AC: 1)
  - [x] Add room resource schema and update endpoint response

## Dev Notes

- Room data comes from YAML config via `spaces.config_file`.
- Use JSON:API envelopes with `application/vnd.api+json` content type.
- Use `internal/rooms` for handler; keep handler free of storage concerns.
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend: `internal/rooms`, `internal/api`, `internal/startup`, `internal/spaces`.
- Frontend: `web/src/views`, `web/src/router`, `web/src/api`.

### References

- PRD FR5: `_bmad-output/planning-artifacts/prd.md` (Areas, Rooms, and Desks Discovery)
- Epic Story 2.2: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (JSON:API, naming)
- API docs rules: `_bmad-output/planning-artifacts/architecture.md` (OpenAPI 3.1 in `api-doc/`)

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Added rooms listing endpoint from YAML config with 404 JSON:API errors on missing areas.
- Added RoomsView, router path, and navigation from areas list.
- Tests: `./run-all-tests.sh`; Cypress `rooms.cy.ts` with test auth.

### File List

- api-doc/endpoints/rooms.yaml
- api-doc/openapi.yaml
- internal/api/errors.go
- internal/api/errors_test.go
- internal/api/response.go
- internal/api/response_test.go
- internal/api/write.go
- internal/areas/handler.go
- internal/rooms/handler.go
- internal/rooms/handler_test.go
- internal/spaces/config.go
- internal/spaces/config_test.go
- internal/startup/server.go
- web/cypress/e2e/rooms.cy.ts
- web/src/api/rooms.ts
- web/src/router/index.test.ts
- web/src/router/index.ts
- web/src/views/AreasView.test.ts
- web/src/views/AreasView.vue
- web/src/views/RoomsView.test.ts
- web/src/views/RoomsView.vue
- web/src/views/testHelpers.ts

### Change Log

- 2026-01-18: Implemented room list API, UI navigation, tests, and docs.
