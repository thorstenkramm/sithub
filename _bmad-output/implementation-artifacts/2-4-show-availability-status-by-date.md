# Story 2.4: Show Availability Status by Date

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an employee,
I want to see which desks are available for a selected date,
so that I can choose a free desk.

## Acceptance Criteria

1. **Given** I have selected a room and date  
   **When** desks are displayed  
   **Then** each desk shows available or occupied status for that date  
   **And** status updates when the date changes

## Tasks / Subtasks

- [x] Add availability support to desks API (AC: 1)
  - [x] Accept date query parameter and return availability status per desk
  - [x] Default to today when date is not provided
  - [x] Return 404 JSON:API error when the room does not exist
  - [x] Enforce auth middleware and JSON:API error handling
- [x] Add frontend date selection and status display (AC: 1)
  - [x] Add date selector for room desks
  - [x] Fetch desks for selected date and render availability status
  - [x] Update status when the date changes
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for availability query
  - [x] Frontend: unit test for date changes and status rendering
  - [x] Cypress E2E: date selection updates desk availability
- [x] Update API documentation (AC: 1)
  - [x] Add availability attribute and date query param to desks endpoint

## Dev Notes

- Bookings are full-day; use a single booking date for availability lookups.
- Use JSON:API envelopes with `application/vnd.api+json` content type.
- Keep handler free of storage concerns; use a small query helper for bookings.
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend: `internal/desks`, `internal/bookings`, `internal/api`, `internal/startup`.
- Frontend: `web/src/views`, `web/src/router`, `web/src/api`.

### References

- PRD FR8: `_bmad-output/planning-artifacts/prd.md` (Areas, Rooms, and Desks Discovery)
- Epic Story 2.4: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (JSON:API, naming)
- API docs rules: `_bmad-output/planning-artifacts/architecture.md` (OpenAPI 3.1 in `api-doc/`)

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Added availability lookup against bookings with date query support and bad request handling.
- Added date selector to desks view with status labels and refresh on date change.
- Tests added for backend availability, frontend date changes, and Cypress verification.
- Documented availability attribute and date query param in OpenAPI.

### File List

- _bmad-output/implementation-artifacts/2-4-show-availability-status-by-date.md
- _bmad-output/implementation-artifacts/sprint-status.yaml
- api-doc/endpoints/desks.yaml
- api-doc/openapi.yaml
- internal/api/errors.go
- internal/api/errors_test.go
- internal/bookings/store.go
- internal/bookings/store_test.go
- internal/desks/handler.go
- internal/desks/handler_test.go
- internal/spaces/config.go
- internal/startup/server.go
- web/src/api/desks.ts
- web/src/views/DesksView.vue
- web/src/views/DesksView.test.ts
- web/cypress/e2e/desks.cy.ts

### Change Log

- 2026-01-18: Story created and set to in-progress.
- 2026-01-18: Story completed.
