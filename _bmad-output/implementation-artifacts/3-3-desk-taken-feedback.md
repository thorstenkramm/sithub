# Story 3.3: Desk-Taken Feedback

Status: complete

## Story

As an employee,
I want a clear message when the desk becomes unavailable during booking,
so that I can choose another desk.

## Acceptance Criteria

1. **Given** I am booking a desk and it becomes unavailable  
   **When** I submit the booking  
   **Then** I see a message that the desk is no longer available for that date  
   **And** I am prompted to choose another desk

## Tasks / Subtasks

- [x] Parse error detail from backend response (AC: 1)
  - [x] Update ApiError to include detail from JSON:API error response
  - [x] Extract detail message when available
- [x] Improve conflict error message (AC: 1)
  - [x] Show backend's specific message (self-duplicate vs other-user conflict)
  - [x] Add prompt to choose another desk
- [x] Refresh desk list after conflict (AC: 1)
  - [x] Auto-refresh availability so user sees updated status
- [x] Add tests (AC: 1)
  - [x] Frontend: unit test for error detail parsing
  - [x] Cypress E2E: test conflict feedback message

## Dev Notes

- Backend already returns detailed messages in JSON:API error format
- Self-duplicate: "You already have this desk booked for this date"
- Other-user conflict: "Desk is already booked for this date"
- Need to extract `errors[0].detail` from response body

### References

- PRD FR11: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 3.3: `_bmad-output/planning-artifacts/epics.md`
- Story 3-1 implementation: `_bmad-output/implementation-artifacts/3-1-create-single-day-booking.md`
- Story 3-2 implementation: `_bmad-output/implementation-artifacts/3-2-prevent-double-booking.md`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Added `detail` property to `ApiError` class to capture JSON:API error detail
- Added `parseErrorDetail` helper to extract `errors[0].detail` from response
- Updated DesksView.vue to show backend's detail message + "Please choose another desk" prompt
- Added auto-refresh of desk list after conflict so user sees updated availability
- Added 3 unit tests for error detail parsing (with detail, without detail, JSON parse error)
- Added 2 Cypress E2E tests (conflict message, self-duplicate message)
- Refactored tests to avoid code duplication using `expectApiError` helper

### File List

**Modified Files:**
- `web/src/api/client.ts` - Added detail property to ApiError, parseErrorDetail function
- `web/src/api/client.test.ts` - Added error detail parsing tests
- `web/src/views/DesksView.vue` - Updated conflict handling with detail message and prompt
- `web/cypress/e2e/desks.cy.ts` - Added conflict feedback E2E tests

### Change Log

- 2026-01-19: Story created and set to in-progress.
- 2026-01-19: Implementation completed - all tasks done, tests passing.
