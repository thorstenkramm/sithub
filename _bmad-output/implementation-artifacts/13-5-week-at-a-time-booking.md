# Story 13.5: Week Booking Mode

Status: done

## Story

As a user,
I want to switch to week booking mode and book multiple days at once,
So that I can reserve my workspace for an entire week efficiently.

## Acceptance Criteria

1. **Given** I am viewing items in an item group
   **When** I see the booking mode toggle
   **Then** I can switch between "book by day" and "book by week" modes
   **And** the selected mode is persisted in browser local storage
   **And** on my next visit, the previously selected mode is restored

2. **Given** I have selected week booking mode
   **When** the date selector is displayed
   **Then** it becomes a calendar week selector showing the next 8 weeks
   **And** each week option displays the Monday date and week number

3. **Given** I have selected a week in week booking mode
   **When** the item tiles are displayed
   **Then** each tile shows a per-day breakdown (MO through FR) with checkboxes
   **And** days booked by other users show the booker's name in red text
   **And** days booked by other users have their checkboxes disabled (cannot be unchecked)
   **And** days I have already booked show my name with a checked checkbox
   **And** free days show "free" in green text with an unchecked checkbox

4. **Given** I am viewing week booking mode on a screen narrower than 600px
   **When** the item tiles are displayed
   **Then** the per-day breakdown uses a compact layout suitable for mobile
   **And** touch targets meet the minimum 44px size requirement

5. **Given** I have checked one or more free days across one or more items
   **When** I look below the item tiles
   **Then** I see a single "Confirm My Booking" button
   **And** the individual "BOOK THIS ITEM" buttons are not shown in week mode

6. **Given** I click "Confirm My Booking" with multiple days selected
   **When** the bookings are submitted
   **Then** each day/item combination is submitted as an individual API request
   to `POST /api/v1/bookings`
   **And** results are collected and displayed per-day
   **And** each successful booking appears in My Bookings
   **And** if any day fails (e.g., concurrent booking), the error is reported for that day
   **And** successful bookings are not rolled back due to a single day's failure

7. **Given** I switch back to day booking mode
   **When** the view updates
   **Then** the standard single-day booking interface is restored
   **And** the "BOOK THIS ITEM" button reappears on each item tile

## Tasks / Subtasks

- [x] Add booking mode toggle to ItemsView (AC: 1, 7)
  - [x] Add a `v-btn-toggle` for "Day" vs "Week" mode
  - [x] Persist selected mode in `localStorage` key `sithub_booking_mode`
  - [x] Restore mode from localStorage on component mount
  - [x] Switching modes re-renders the date selector and item tiles
- [x] Replace date picker with week selector in week mode (AC: 2)
  - [x] Extract shared `useWeekSelector` composable from ItemGroupsView
  - [x] Display Monday date + week number per option
  - [x] Default to current week
- [x] Fetch per-day booking data for selected week (AC: 3)
  - [x] Call items endpoint for each weekday (Mon-Fri) in parallel
  - [x] Collect booking status per item per day: free, booked-by-me, booked-by-other
  - [x] Resolve booker names from existing API response
- [x] Render per-day breakdown on item tiles (AC: 3, 4)
  - [x] Show MO/TU/WE/TH/FR columns with checkboxes
  - [x] Free days: green "free" text, unchecked checkbox, enabled
  - [x] Booked by other: red booker name, disabled checkbox
  - [x] Booked by me: user name, checked checkbox, disabled
  - [x] Mobile layout (< 600px): compact grid, 44px touch targets
- [x] Add "Confirm My Booking" button (AC: 5)
  - [x] Show below item tiles when in week mode and at least one day is checked
  - [x] Hide individual "BOOK THIS ITEM" buttons in week mode
  - [x] Display count: "Confirm My Booking (3 days)"
- [x] Implement batch booking submission (AC: 6)
  - [x] Collect all checked day/item combinations
  - [x] Submit each as individual `POST /api/v1/bookings` request
  - [x] Use `Promise.allSettled()` to handle partial failures
  - [x] Display results: green checkmarks for success, red X for failures
  - [x] Show specific error messages per failed day
  - [x] Refresh item tiles after submission to reflect new state
- [x] Add Vitest unit tests (AC: 1, 2, 3)
  - [x] Test useWeekSelector composable (15 tests)
  - [x] Test mode toggle default and localStorage persistence
  - [x] Test week mode rendering and data fetching
- [x] Add Cypress E2E test (AC: 1, 3, 5, 6, 7)
  - [x] Switch to week mode, verify week selector appears
  - [x] Per-day breakdown with checkboxes
  - [x] Select days, verify "Confirm" button appears
  - [x] Submit bookings, verify success feedback
  - [x] Switch back to day mode, verify standard UI restored
  - [x] Verify localStorage persistence across page reload
- [x] Refactor ItemGroupsView to use shared `useWeekSelector` composable

## Dev Notes

### Architecture: Frontend-Heavy Story

This story requires NO backend changes. The existing `POST /api/v1/bookings` endpoint
already supports single-day booking creation. The week mode simply batches multiple
calls from the frontend.

The existing multi-day booking endpoint (`booking_dates` array) could also be used
per item, but individual requests per day/item give better granularity for partial
failure handling as specified in AC 6.

### Booking Mode Persistence

Use `localStorage.getItem('sithub_booking_mode')` / `setItem()`. Values: `'day'`
(default) or `'week'`. Check on component mount:

```typescript
const bookingMode = ref<'day' | 'week'>(
  (localStorage.getItem('sithub_booking_mode') as 'day' | 'week') || 'day'
)

watch(bookingMode, (mode) => {
  localStorage.setItem('sithub_booking_mode', mode)
})
```

### Per-Day Data Fetching Strategy

Two options for fetching week data:

**Option A - 5 parallel item fetches:** Call the items endpoint 5 times (one per
weekday) with the date parameter. Simple, uses existing API.

**Option B - New batch endpoint:** Create a backend endpoint that returns items
with availability for multiple dates. More efficient but requires backend work.

Recommend **Option A** for v1. Five parallel requests complete quickly and avoid
backend changes. If performance becomes an issue, Option B can be added later.

### Dependencies

- **Story 13.2 (Booker Display Name):** Required for showing booker names in the
  per-day breakdown. If 13.2 is not complete, show "Booked" instead of the name.
- **Story 13.3 (Weekly Availability):** The week selector pattern can be shared.
  Extract a `useWeekSelector` composable if both stories need it.

### Mobile Layout (AC: 4)

The per-day breakdown needs a compact representation at < 600px. Consider:
- Abbreviated day labels (M/T/W/T/F instead of MO/TU/WE/TH/FR)
- Stacked layout instead of horizontal row
- Checkboxes sized to 44px minimum touch target

### Partial Failure Handling (AC: 6)

Use `Promise.allSettled()` to ensure all requests complete regardless of individual
failures:

```typescript
const results = await Promise.allSettled(
  selections.map(({ itemId, date }) =>
    createBooking({ item_id: itemId, booking_date: date })
  )
)
```

Display results in a summary dialog or inline feedback per day.

### References

- PRD FR38: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 13.5: `_bmad-output/planning-artifacts/epics.md`
- ItemsView: `web/src/views/ItemsView.vue`
- Bookings API: `web/src/api/bookings.ts`
- Bookings handler: `internal/bookings/handler.go`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Review fixes (2026-02-14): week mode uses item IDs, derives "booked by me" from My Bookings instead of display names, validates colleague/guest inputs before submit, surfaces weekly load errors, hides multi-day toggle in week mode.
- All 7 acceptance criteria implemented and verified
- Frontend-only story: no backend changes required
- Extracted `useWeekSelector` composable to share week logic with ItemGroupsView
- Refactored ItemGroupsView to use shared composable (removed duplication)
- Vitest: 128 tests pass (22 new: 15 useWeekSelector + 7 ItemsView week mode)
- Cypress: 51 E2E tests pass (6 new week-booking tests)
- ESLint: clean, TypeScript type-check: clean, build: passes
- jscpd: 0% TS duplication

### File List

**Frontend (Vue/TypeScript)**

- `.gitignore` (review fix: ignore local artifacts)
- `web/src/composables/useWeekSelector.ts` (new - shared week selector logic)
- `web/src/composables/useWeekSelector.test.ts` (new - 15 unit tests)
- `web/src/views/ItemsView.vue` (week mode toggle, week selector, per-day breakdown,
  confirm button, batch submission, booking results)
- `web/src/views/ItemsView.test.ts` (7 new week mode tests)
- `web/src/views/ItemGroupsView.vue` (refactored to use shared useWeekSelector)

**Cypress E2E**

- `web/cypress/e2e/week-booking.cy.ts` (new, 6 tests)