# Story 19.4: Cancel Booking from Week View

Status: done

## Story

As a user,
I want to cancel my bookings directly from the week view,
So that I don't have to navigate to My Bookings to undo a booking.

## Acceptance Criteria

1. **Given** I am on the week view and a day/item has my booking (shown with a checkmark)
   **When** the page renders
   **Then** a small red cancel icon appears next to the checkmark

2. **Given** I click the red cancel icon
   **When** the cancellation is processed
   **Then** the booking is cancelled and the checkmark and cancel icon are removed
   **And** the item becomes bookable again for that day

3. **Given** the booking belongs to another user
   **When** the page renders
   **Then** no cancel icon is shown for that booking

## Tasks / Subtasks

- [x] Add cancel icon to "booked-by-me" week cells (AC: 1, 3)
  - [x] In `ItemsView.vue`: added red `$cancelCircle` icon on week cells where the
    current user has a booking
  - [x] Icon only appears for the current user's bookings, not other users'
  - [x] Added `mdiCloseCircle` icon to Vuetify aliases
- [x] Implement cancel booking action (AC: 2)
  - [x] Created `cancelWeekBooking()` function in `ItemsView.vue`
  - [x] On click, cancels the booking via API and refreshes week data
  - [x] Checkmark and cancel icon removed after successful cancellation
- [x] Change `myWeekBookings` from `Set` to `Map` (AC: 2)
  - [x] Changed from `Set<string>` to `Map<string, string>` to store booking IDs
  - [x] Booking IDs needed for the cancel API call
- [x] Register cancel icon in Vuetify (AC: 1)
  - [x] Added `mdiCloseCircle` icon and `$cancelCircle` alias in `vuetify.ts`
- [x] Verify E2E tests still pass

## Dev Notes

### Data Structure Change

The `myWeekBookings` data structure was changed from a `Set<string>` (tracking only
"itemId-date" keys) to a `Map<string, string>` (mapping "itemId-date" keys to booking IDs).
This allows the cancel function to look up the booking ID needed for the DELETE API call.

### Cancel Icon

Uses `mdiCloseCircle` from `@mdi/js`, registered as `$cancelCircle` alias in Vuetify. The
icon is rendered in red and only appears on cells where the current user owns the booking.

### References

- Epic 19 Story 19.4: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR70: `_bmad-output/planning-artifacts/prd.md`
- `web/src/views/ItemsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added cancel icon next to checkmark on user's own week bookings
- Changed `myWeekBookings` from `Set` to `Map` for booking ID tracking
- `cancelWeekBooking()` calls DELETE API and refreshes week data
- Registered `mdiCloseCircle` as `$cancelCircle` in Vuetify aliases
- All existing tests continue to pass

### File List

- `web/src/views/ItemsView.vue` — Cancel icon, `cancelWeekBooking()`, `Map` refactor
- `web/src/plugins/vuetify.ts` — Added `mdiCloseCircle` and `$cancelCircle` alias

## Change Log

- 2026-03-21: Story implemented and verified.
