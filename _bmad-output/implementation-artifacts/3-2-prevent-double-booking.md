# Story 3.2: Prevent Double-Booking

Status: complete

<!-- Note: This story was implemented as part of Story 3-1 -->

## Story

As an employee,
I want the system to prevent duplicate bookings for the same desk and day,
so that I don't book a desk that's already taken.

## Acceptance Criteria

1. **Given** a desk is already booked for a date  
   **When** another booking is attempted for the same desk and date  
   **Then** the request is rejected  
   **And** no duplicate booking is created

## Tasks / Subtasks

- [x] Database constraint for uniqueness (AC: 1)
  - [x] SQLite unique constraint on (desk_id, booking_date) prevents duplicates
- [x] Backend conflict detection (AC: 1)
  - [x] Detect SQLite constraint violation and return 409 Conflict
  - [x] Return JSON:API error with descriptive message
- [x] Self-duplicate detection (AC: 1)
  - [x] Check if same user already has booking for same desk+date
  - [x] Return user-friendly message "You already have a booking for this desk on this date"
- [x] Add tests (AC: 1)
  - [x] Backend: test double-booking returns 409 Conflict
  - [x] Backend: test self-duplicate returns specific message

## Dev Notes

- This story was fully implemented as part of Story 3-1 (Create Single-Day Booking).
- The unique constraint on (desk_id, booking_date) in the bookings table prevents duplicates at the database level.
- The backend handler detects the constraint violation via `ErrConflict` sentinel error and returns HTTP 409.
- Additional self-duplicate check provides better UX when user tries to book same desk twice.

### Implementation Details

**Database (from migration 001):**
```sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_bookings_desk_date ON bookings(desk_id, booking_date);
```

**Backend conflict handling (internal/bookings/store.go):**
- `ErrConflict` sentinel error returned when SQLite UNIQUE constraint violated
- Detected via `sqlite3.ErrConstraintUnique` error code

**Backend handler (internal/bookings/handler.go):**
- `FindUserBooking` checks for self-duplicate before insert
- Returns 409 with appropriate message for both cases

### References

- PRD FR10: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 3.2: `_bmad-output/planning-artifacts/epics.md`
- Story 3-1 implementation: `_bmad-output/implementation-artifacts/3-1-create-single-day-booking.md`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Story requirements were implemented as part of Story 3-1
- Unique constraint on (desk_id, booking_date) prevents database-level duplicates
- Handler returns 409 Conflict for both conflict scenarios:
  - Self-duplicate: "You already have a booking for this desk on this date"
  - Other-user conflict: "This desk is already booked for this date"
- All tests passing in `./run-all-tests.sh`

### File List

**Relevant Files (implemented in Story 3-1):**
- `internal/bookings/store.go` - ErrConflict sentinel, constraint detection
- `internal/bookings/handler.go` - Conflict handling, self-duplicate check
- `internal/bookings/handler_test.go` - Conflict test cases
- `internal/api/errors.go` - WriteConflict helper

### Change Log

- 2026-01-19: Story created and marked complete (implemented as part of 3-1).
