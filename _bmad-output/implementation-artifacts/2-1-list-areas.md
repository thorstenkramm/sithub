# Story 2.1: List Areas

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an employee,
I want to see the list of areas,
so that I can choose where to book.

## Acceptance Criteria

1. **Given** I am authenticated  
   **When** I open the app  
   **Then** I see all configured areas  
   **And** the list is empty-safe (shows zero areas without error)

## Tasks / Subtasks

- [x] Add areas config loader (AC: 1)
  - [x] Load areas from YAML configuration file
  - [x] Validate required fields (id, name)
- [x] Add areas API endpoint (AC: 1)
  - [x] Implement `GET /api/v1/areas` returning JSON:API collection
  - [x] Enforce auth middleware and JSON:API error handling
- [x] Add frontend list view (AC: 1)
  - [x] Fetch areas and render a list with empty-state copy
  - [x] Ensure accessible labels and visible focus states (WCAG A)
- [x] Add tests (AC: 1)
  - [x] Backend: config loader and handler tests for empty list
  - [x] Frontend: unit test for empty and non-empty states
  - [x] Cypress E2E: authenticated user sees areas list
- [x] Add API documentation (AC: 1)
  - [x] Add endpoint doc in `api-doc/endpoints/areas.yaml`
  - [x] Wire the path in `api-doc/openapi.yaml`

## Dev Notes

- Use JSON:API envelopes with `application/vnd.api+json` content type.
- Use `internal/areas` for handlers and `internal/spaces` for YAML loading.
- For now, data can be seeded in tests; list output comes from YAML config.
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend: `internal/areas`, `internal/api`, `internal/startup`.
- Frontend: `web/src/views`, `web/src/api`.

### References

- PRD FR4: `_bmad-output/planning-artifacts/prd.md` (Areas, Rooms, and Desks Discovery)
- Epic Story 2.1: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (JSON:API, naming)
- API docs rules: `_bmad-output/planning-artifacts/architecture.md` (OpenAPI 3.1 in `api-doc/`)

## Dev Agent Record

### Agent Model Used

sm - Bob

### Debug Log References

None.

### Completion Notes List

- Implemented YAML-backed spaces loader and JSON:API `GET /api/v1/areas`.
- Added spaces config path to `sithub.example.toml` and CLI flags.
- Added Vue list rendering with empty state and error/loading affordances.
- Tests: `./run-all-tests.sh` (includes `go test -race` and frontend unit coverage).
- Cypress E2E spec added for areas list (not run; requires dev server).

### File List

- api-doc/endpoints/areas.yaml
- api-doc/openapi.yaml
- internal/areas/handler.go
- internal/areas/handler_test.go
- internal/spaces/config.go
- internal/spaces/config_test.go
- internal/startup/server.go
- internal/startup/server_test.go
- sithub.example.toml
- sithub.toml
- web/cypress/e2e/areas.cy.ts
- web/src/api/areas.ts
- web/src/api/types.ts
- web/src/views/AreasView.test.ts
- web/src/views/AreasView.vue

### Change Log

- 2026-01-18: Implemented areas list API, UI, tests, and documentation.
