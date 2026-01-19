# Story 7.1: Book on Behalf of Another User

Status: in-progress

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

- [ ] Update booking creation to accept optional `for_user_id` (AC: 1)
  - [ ] Modify POST /api/v1/bookings to accept `for_user_id` and `for_user_name`
  - [ ] Store `booked_by_user_id` to track who made the booking
  - [ ] Validate that current user is authenticated
- [ ] Update booking response to include booking creator info (AC: 1)
  - [ ] Add `booked_by_user_id`, `booked_by_user_name` to response
  - [ ] Indicate if booking was made on behalf
- [ ] Update My Bookings to show bookings made for me (AC: 1)
  - [ ] List bookings where `user_id` matches OR `booked_by_user_id` matches
  - [ ] Show indication if booking was made by someone else
- [ ] Allow either party to cancel (AC: 1)
  - [ ] Original booker can cancel
  - [ ] User for whom booking was made can cancel
- [ ] Add frontend "Book for colleague" option (AC: 1)
  - [ ] Add user search/input field when booking
  - [ ] Show whose booking it is in My Bookings
- [ ] Add tests (AC: 1)
- [ ] Update API documentation (AC: 1)

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

(To be filled after completion)

### File List

(To be filled after completion)

### Change Log

- 2026-01-19: Story created and set to in-progress.
