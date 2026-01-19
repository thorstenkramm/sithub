# Story 5.1: Room Booking Overview

Status: complete

## Story

As an employee,
I want to see a room-level booking overview for a date,
So that I can understand room utilization.

## Acceptance Criteria

1. **Given** I select a room and date  
   **When** I view the overview  
   **Then** I see all booked desks and associated users for that date

## Tasks / Subtasks

- [x] Add backend endpoint for room bookings (AC: 1)
  - [x] Implement `GET /api/v1/rooms/:id/bookings?date=YYYY-MM-DD`
  - [x] Return list of bookings with desk name and user info
  - [x] Handle room not found (404)
  - [x] Default to today if date not provided
- [x] Add frontend room bookings view (AC: 1)
  - [x] Create RoomBookingsView component
  - [x] Show list of booked desks with user names
  - [x] Add date picker to change date
  - [x] Add route `/rooms/:roomId/bookings`
- [x] Add navigation to room bookings (AC: 1)
  - [x] Add "View Bookings" link from DesksView
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for room bookings
  - [x] Frontend: unit tests for RoomBookingsView
- [x] Update API documentation (AC: 1)
  - [x] Document GET /api/v1/rooms/:id/bookings endpoint

## Dev Notes

- Need to join bookings with user display names - but we don't have a users table
- Option 1: Return user_id only (like we did for admin cancel)
- Option 2: Store user display_name in bookings table when created
- For MVP, we'll store user_name in bookings table at creation time
- This requires a migration to add user_name column to bookings

### Project Structure Notes

- Backend: `internal/rooms/bookings_handler.go` (new)
- Frontend: `web/src/views/RoomBookingsView.vue` (new)
- Migration: Add user_name to bookings table

### References

- PRD FR15: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 5.1: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Added `user_name` column to bookings table via migration 000003
- Backend endpoint returns desk_id, desk_name, user_id, user_name, booking_date per booking
- Refactored common date parsing to `api.ParseBookingDate` and `api.ParseRoomRequest` to eliminate code duplication
- Frontend shows table with desk name and user name, with date picker navigation

### File List

**Backend:**
- `migrations/000003_add_user_name_to_bookings.up.sql` (new)
- `migrations/000003_add_user_name_to_bookings.down.sql` (new)
- `internal/api/response.go` (modified - added ParseBookingDate, ParseRoomRequest)
- `internal/api/response_test.go` (modified - added tests)
- `internal/bookings/handler.go` (modified - accept/store userName)
- `internal/bookings/testhelpers_test.go` (modified - helper for seeding with name)
- `internal/desks/handler.go` (modified - use api.ParseBookingDate)
- `internal/rooms/bookings_handler.go` (new)
- `internal/rooms/bookings_handler_test.go` (new)
- `internal/startup/server.go` (modified - register route)

**Frontend:**
- `web/src/api/roomBookings.ts` (new)
- `web/src/views/RoomBookingsView.vue` (new)
- `web/src/views/RoomBookingsView.test.ts` (new)
- `web/src/router/index.ts` (modified - add route)
- `web/src/router/index.test.ts` (modified - update test)
- `web/src/views/DesksView.vue` (modified - add link)

**API Docs:**
- `api-doc/endpoints/room-bookings.yaml` (new)
- `api-doc/openapi.yaml` (modified)

### Change Log

- 2026-01-19: Story created and set to in-progress.
- 2026-01-19: Story completed.
