# Story 28.1: Preserve Selected Date After Booking

Status: done

## Story

As a user,
I want the date selector to stay on my selected date after completing a booking,
so that I can continue browsing availability for the same date without being sent back to
today.

## Acceptance Criteria

1. **Given** I am on the items page with a future date selected (e.g. 30 April 2026)
   **When** I book an item and the booking confirmation completes
   **Then** the date picker still shows 30 April 2026
   **And** the displayed items reflect booking status for 30 April 2026 (including my new
   booking shown as occupied)

2. **Given** I am in week booking mode with a future week selected
   **When** I confirm bookings for that week
   **Then** the week selector stays on the same week
   **And** my new bookings are shown as booked in the week view

## Tasks / Subtasks

- [x] Task 1: Remove unconditional date reset after booking (AC: #1, #2)
  - [x] 1.1 Open `web/src/views/ItemsView.vue`, locate the `bookItem()` function
        (around line 1678-1764)
  - [x] 1.2 Find the date reset block (lines ~1715-1719):
        ```javascript
        resetDayToToday();
        const resetDate = getDay();
        const dayChanged = selectedDate.value !== resetDate;
        selectedDate.value = resetDate;
        ```
  - [x] 1.3 Remove the `resetDayToToday()` call — do NOT reset sessionStorage
  - [x] 1.4 Remove the `selectedDate.value = resetDate` assignment — keep the current
        `selectedDate.value` unchanged
  - [x] 1.5 Keep the items reload logic (`loadItems()`) but trigger it with the current
        `selectedDate.value` instead of `resetDate`
  - [x] 1.6 Verify the watch on `selectedDate` (line ~1950-1960) does not interfere —
        since the date is not changing, the watch should NOT fire, avoiding a double reload

- [x] Task 2: Verify week mode preserves selected week (AC: #2)
  - [x] 2.1 Check the week booking confirmation flow in `ItemsView.vue` for similar
        reset patterns
  - [x] 2.2 If the week selector is also reset after booking, apply the same fix:
        preserve `selectedWeek` value
  - [x] 2.3 Verify `useDateState.ts` composable (`web/src/composables/useDateState.ts`)
        — the `setDay()` function (lines ~106-109) should still persist the user's
        selected date to sessionStorage after booking

- [x] Task 3: Validation
  - [x] 3.1 Manual test: select a future date, book an item, verify date stays
  - [x] 3.2 Manual test: select a future week, book in week mode, verify week stays
  - [x] 3.3 Run `cd web && npx vitest run` — all unit tests pass
  - [x] 3.4 Run `cd web && npm run type-check` — no type errors
  - [x] 3.5 Run `cd web && npm run lint` — no lint errors
  - [x] 3.6 Run `cd web && npm run build` — builds cleanly
  - [x] 3.7 Run `cd web && npm run test:e2e -- --browser electron` — E2E tests pass

## Dev Notes

### Architecture & Patterns

This is a frontend-only bug fix. No backend changes needed.

**Primary file:** `web/src/views/ItemsView.vue`
**Supporting file:** `web/src/composables/useDateState.ts`

### Root Cause Analysis

After a successful booking, `bookItem()` unconditionally calls `resetDayToToday()` which
clears the sessionStorage memorized date and resets `selectedDate` to today. This was
likely added to ensure users see their "current day" bookings, but it breaks the UX when
booking for future dates.

The fix is straightforward: after booking, reload the items for the **currently selected
date** instead of resetting to today. The user intentionally selected that date and expects
to stay there.

### Key Code Locations

| Element | Location | Notes |
|---------|----------|-------|
| `bookItem()` | `ItemsView.vue` ~line 1678 | Booking action handler |
| Date reset block | `ItemsView.vue` ~line 1715-1719 | The bug — remove these lines |
| `selectedDate` ref | `ItemsView.vue` ~line 979 | Local reactive date state |
| `resetDayToToday()` | `useDateState.ts` ~line 80-83 | Resets sessionStorage to today |
| `getDay()` | `useDateState.ts` ~line 98-104 | Retrieves stored date |
| `setDay()` | `useDateState.ts` ~line 106-109 | Persists date to sessionStorage |
| Watch on selectedDate | `ItemsView.vue` ~line 1950-1960 | Triggers loadItems on change |

### Anti-Patterns to Avoid

- Do NOT add a new prop or emit for date preservation — the fix is removing the reset,
  not adding new state management
- Do NOT change `useDateState.ts` — the composable is correct; the bug is in the caller
- Do NOT add `cy.wait()` or fixed delays in tests — use intercept aliases
- Do NOT touch the week selector component itself — the fix is in the booking handler

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- Type-check: clean
- Lint: clean
- Unit tests: 318/318 pass
- Build: clean

### Completion Notes List

- Removed `resetDayToToday()` call and `selectedDate.value = resetDate` assignment from `bookItem()`
- Simplified reload logic: always reload with current `selectedDate.value` after booking
- Removed unused `resetDayToToday` from destructuring of `useDateState()`
- Updated test: renamed "resets the live selected day" to "preserves the selected day" and updated assertions to expect the stored future date instead of today
- Week mode already correct — `submitWeekBookings()` never reset the week selector

### Change Log

- 2026-04-15: Fixed date selector jumping to today after booking (Story 28.1)

### File List

- web/src/views/ItemsView.vue (modified — removed date reset after booking)
- web/src/views/ItemsView.test.ts (modified — updated test to match new behavior)
