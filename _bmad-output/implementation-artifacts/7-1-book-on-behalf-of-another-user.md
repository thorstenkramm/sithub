# Story 7.1: Book on Behalf of Another User

Status: done

## Story

As an employee,
I want to book a desk on behalf of another user,
So that we can sit together.

## Acceptance Criteria

1. **Given** I book a desk for another user  
   **When** the booking is created  
   **Then** it appears in both users' booking lists  
   **And** either user can cancel it

## Tasks / Subtasks

- [x] Update booking creation to accept optional `for_user_id` (AC: 1)
  - [x] Modify POST /api/v1/bookings to accept `for_user_id` and `for_user_name`
  - [x] Store `booked_by_user_id` to track who made the booking
  - [x] Validate that current user is authenticated
- [x] Update booking response to include booking creator info (AC: 1)
  - [x] Add `booked_by_user_id`, `booked_by_user_name` to response
  - [x] Indicate if booking was made on behalf (`booked_for_me` flag)
- [x] Update My Bookings to show bookings made for me (AC: 1)
  - [x] List bookings where `user_id` matches OR `booked_by_user_id` matches
  - [x] Show indication if booking was made by someone else (chip display)
- [x] Allow either party to cancel (AC: 1)
  - [x] Original booker can cancel
  - [x] User for whom booking was made can cancel
- [x] Add frontend "Book for colleague" option (AC: 1)
  - [x] Add checkbox toggle and user input fields when booking
  - [x] Show whose booking it is in My Bookings via chip display
- [x] Add tests (AC: 1)
- [x] Update API documentation (AC: 1)

## Dev Notes

### Database Changes

Need migration to add `booked_by_user_id` and `booked_by_user_name` columns:
- `booked_by_user_id` - user who created the booking (may be different from user_id)
- `booked_by_user_name` - display name of user who created the booking

### Design Decisions

- Booking `user_id` = the person who will use the desk
- Booking `booked_by_user_id` = the person who made the booking
- For self-bookings, these are the same
- No user lookup/directory service - frontend must provide both user_id and user_name

### References

- PRD FR20: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 7.1: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Migration 000004 adds `booked_by_user_id` and `booked_by_user_name` columns to bookings table
- Backend handler accepts optional `for_user_id` and `for_user_name` in POST request
- My Bookings returns bookings where user is either the target or the booker
- Cancel allows target user OR booker to cancel (or admin)
- Frontend has checkbox toggle to expand "Book for colleague" fields
- MyBookingsView shows chip indicating who booked on behalf

### File List

Backend:
- `migrations/000004_add_booked_by_columns.up.sql` - Add columns
- `migrations/000004_add_booked_by_columns.down.sql` - Remove columns
- `internal/bookings/handler.go` - CreateHandler, ListHandler, DeleteHandler updated
- `internal/bookings/store.go` - BookingRecord struct, ListUserBookings, FindBookingByID updated
- `internal/bookings/testhelpers_test.go` - Added `seedTestBookingFull` helper
- `internal/bookings/handler_test.go` - Added 6 new tests (refactored to table-driven)

Frontend:
- `web/src/api/bookings.ts` - Updated types and createBooking function
- `web/src/views/DesksView.vue` - Added "Book for colleague" checkbox and fields
- `web/src/views/MyBookingsView.vue` - Added chip display for on-behalf bookings
- `web/src/views/MyBookingsView.test.ts` - Updated mocks and added test
- `web/src/api/bookings.test.ts` - Added test for on-behalf booking

API Docs:
- `api-doc/openapi.yaml` - Added for_user_id, for_user_name, booked_by_* fields to schemas
- `api-doc/endpoints/bookings.yaml` - Updated descriptions and examples

### Change Log

- 2026-01-19: Story created and set to in-progress.
- 2026-01-19: Backend implementation complete with tests.
- 2026-01-19: Frontend and API documentation complete. Story done.
