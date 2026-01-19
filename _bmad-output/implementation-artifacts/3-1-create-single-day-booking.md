# Story 3.1: Create Single-Day Booking

Status: complete

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an employee,
I want to book a desk for a specific date,
so that I can reserve my workspace.

## Acceptance Criteria

1. **Given** I have selected a desk and date  
   **When** I confirm the booking  
   **Then** the booking is created for that date  
   **And** it appears in "My Bookings"

## Tasks / Subtasks

- [x] Add bookings API endpoint (AC: 1)
  - [x] Implement `POST /api/v1/bookings` to create a single-day booking
  - [x] Enforce auth middleware and JSON:API error handling
- [x] Persist booking in database (AC: 1)
  - [x] Validate desk exists and store booking for the selected date
- [x] Update frontend booking flow (AC: 1)
  - [x] Add booking action to desks list
  - [x] Confirm booking and handle success/error states
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for booking creation
  - [x] Frontend: unit test for booking action
  - [x] Cypress E2E: booking a desk shows success message
- [x] Update API documentation (AC: 1)
  - [x] Add booking create endpoint and schema

## Dev Notes

- Bookings are full-day; store a single booking_date per reservation.
- Use JSON:API envelopes with `application/vnd.api+json` content type.
- Use optimistic conflict handling with a unique constraint on (desk_id, booking_date).
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend: `internal/bookings`, `internal/api`, `internal/startup`.
- Frontend: `web/src/views`, `web/src/api`.

### References

- PRD FR9: `_bmad-output/planning-artifacts/prd.md` (Booking)
- Epic Story 3.1: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (JSON:API, naming)
- API docs rules: `_bmad-output/planning-artifacts/architecture.md` (OpenAPI 3.1 in `api-doc/`)

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Implemented `POST /api/v1/bookings` endpoint with JSON:API request/response handling
- Added `CreateBooking` function with SQLite unique constraint conflict detection
- Added `FindDesk` method to spaces config for desk validation
- Added `WriteConflict` error helper to internal/api package
- Updated DesksView.vue with Book button and success/error handling
- Created bookings.ts API client for frontend
- Added comprehensive backend handler tests (12 test cases after code review)
- Added frontend unit tests for createBooking function
- Added Cypress E2E test for booking flow
- Updated OpenAPI documentation with CreateBookingRequest and BookingAttributes schemas

**Code Review Fixes (2026-01-19):**
- Issue 1 (HIGH): Added past date validation - cannot book dates in the past
- Issue 2 (MEDIUM): Added self-duplicate booking check with user-friendly message
- Issue 3 (MEDIUM): Added Content-Type header validation per JSON:API spec
- Issue 4 (MEDIUM): Added FindDesk unit test coverage
- Issue 5 (LOW): Improved frontend error test to verify status code propagation
- Issue 6 (LOW): Added structured logging with log/slog for bookings
- Issue 7 (LOW): Added data-cy-availability attribute for better E2E selectors
- Issue 8 (LOW): Added invalid JSON body test case

### File List

**New Files:**
- `internal/bookings/handler.go` - Booking create handler and CreateBooking function
- `internal/bookings/handler_test.go` - Handler tests
- `internal/bookings/testhelpers_test.go` - Shared test helpers
- `web/src/api/bookings.ts` - Frontend booking API client
- `web/src/api/bookings.test.ts` - Frontend booking API tests

**Modified Files:**
- `internal/api/errors.go` - Added WriteConflict and WriteUnsupportedMediaType helpers
- `internal/spaces/config.go` - Added FindDesk method
- `internal/spaces/config_test.go` - Added FindDesk unit test
- `internal/startup/server.go` - Registered POST /api/v1/bookings route
- `internal/bookings/store_test.go` - Updated to use shared test helpers
- `web/src/views/DesksView.vue` - Added Book button, booking UI, and data-cy attributes
- `web/cypress/e2e/desks.cy.ts` - Added booking E2E test with improved selectors
- `api-doc/endpoints/bookings.yaml` - Updated with full endpoint documentation
- `api-doc/openapi.yaml` - Added booking schemas

### Change Log

- 2026-01-18: Story created and set to in-progress.
- 2026-01-19: Implementation completed - all tasks done, tests passing.
- 2026-01-19: Code review completed - 8 issues fixed (1 HIGH, 3 MEDIUM, 4 LOW).
