# Story 19.1: Fix Cancel Booking Dialog Not Closing

Status: done

## Story

As a user,
I want the cancel booking confirmation dialog to close after I confirm the cancellation,
So that I am not left with a stale dialog on screen.

## Acceptance Criteria

1. **Given** I am on the My Bookings page and click cancel on a booking
   **When** the confirmation dialog appears and I click "Cancel Booking"
   **Then** the booking is removed from the list
   **And** the confirmation dialog closes automatically

## Tasks / Subtasks

- [x] Fix dialog not closing after cancellation (AC: 1)
  - [x] In `MyBookingsView.vue`: add `showCancelDialog.value = false` in the `finally`
    block of `confirmCancelBooking`
  - [x] Dialog now closes regardless of success or failure
- [x] Verify E2E tests still pass

## Dev Notes

### Root Cause

The `confirmCancelBooking` function in `MyBookingsView.vue` was missing the dialog close
statement. After the API call completed (or failed), the dialog remained open because
`showCancelDialog` was never set back to `false`. Adding the assignment in the `finally`
block ensures the dialog closes in all cases.

### References

- Epic 19 Story 19.1: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR67: `_bmad-output/planning-artifacts/prd.md`
- `web/src/views/MyBookingsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Single-line bug fix: added `showCancelDialog.value = false` in `finally` block
- No new tests required; existing E2E tests cover the cancel booking flow
- All existing tests continue to pass

### File List

- `web/src/views/MyBookingsView.vue` — Added dialog close in `finally` block of
  `confirmCancelBooking`

## Change Log

- 2026-03-21: Story implemented and verified.
