# Story 4.3: Admin Cancel Any Booking

Status: complete

## Story

As an admin,
I want to cancel any booking,
So that I can resolve conflicts.

## Acceptance Criteria

1. **Given** I am an admin  
   **When** I cancel another user's booking  
   **Then** the booking is removed from all relevant lists  
   **And** the affected user sees the cancellation

## Tasks / Subtasks

- [x] Add backend admin cancel capability (AC: 1)
  - [x] Modify DELETE handler to allow admins to delete any booking
  - [x] Keep existing user-can-delete-own logic
  - [x] Add tests for admin cancel scenarios
- [x] Add frontend admin cancellation controls (AC: 1)
  - [x] Update useAuthStore to include isAdmin flag
  - [x] Show admin cancel controls where bookings are displayed (DesksView)
  - [x] Call DELETE endpoint and refresh list
- [x] Update API documentation (AC: 1)
  - [x] Document admin access to DELETE endpoint
  - [x] Document booking_id and booker_name fields in desk response (admin-only)

## Dev Notes

- Story 4-2 added the basic DELETE endpoint - now extend it for admin use
- Admin check: `user.IsAdmin` is already available in auth context
- Frontend: `is_admin` is returned by `/api/v1/me` endpoint
- Admin should be able to cancel any booking, not just their own
- Return 404 for non-admin trying to delete others' bookings (existing behavior)

### Project Structure Notes

- Backend: `internal/bookings/handler.go` (modify DeleteHandler)
- Frontend: `web/src/stores/useAuthStore.ts`, various views

### References

- PRD FR14: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 4.3: `_bmad-output/planning-artifacts/epics.md`
- Story 4-2: `_bmad-output/implementation-artifacts/4-2-cancel-my-booking.md`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Backend DELETE handler now checks `user.IsAdmin` - admins can delete any booking
- Non-admins still see 404 when trying to delete others' bookings (existing behavior)
- Desks endpoint now returns `booking_id` for admins when desk is occupied
- Admin cancel button shown in DesksView for occupied desks
- All tests passing (47 frontend unit tests, backend handler tests)

### File List

**Backend:**
- `internal/bookings/handler.go` - Modified DeleteHandler for admin support
- `internal/bookings/handler_test.go` - Added admin cancel tests (table-driven)
- `internal/bookings/store.go` - Added `FindDeskBookings` function
- `internal/bookings/store_test.go` - Added tests for FindDeskBookings
- `internal/desks/handler.go` - Refactored to include booking info for admins

**Frontend:**
- `web/src/stores/useAuthStore.ts` - Added `isAdmin` state
- `web/src/stores/useAuthStore.test.ts` - Added isAdmin tests
- `web/src/api/desks.ts` - Added `booking_id`, `booker_name` to DeskAttributes
- `web/src/views/DesksView.vue` - Admin cancel button, auth store integration
- `web/src/views/DesksView.test.ts` - Added Pinia setup

**API Docs:**
- `api-doc/endpoints/booking.yaml` - Updated DELETE description for admin
- `api-doc/endpoints/desks.yaml` - Added description for admin fields
- `api-doc/openapi.yaml` - Added booking_id, booker_name to DeskAttributes

### Change Log

- 2026-01-19: Story created and set to in-progress.
- 2026-01-19: Story completed - all acceptance criteria met.
