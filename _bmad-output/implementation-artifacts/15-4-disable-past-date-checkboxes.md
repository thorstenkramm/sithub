# Story 15.4: Disable Past Date Checkboxes in Week Mode

Status: done

## Story

As a user,
I want past date checkboxes disabled in week booking mode,
So that I don't waste time selecting dates the backend would reject anyway.

## Acceptance Criteria

1. **Given** I am in week booking mode and the selected week includes past dates
   **When** I see the per-day checkboxes
   **Then** checkboxes for dates before today are disabled and visually grayed out
   **And** I cannot check or uncheck past date checkboxes

2. **Given** I am in week booking mode and the selected week is entirely in the future
   **When** I see the per-day checkboxes
   **Then** all free day checkboxes are enabled and interactive

## Tasks / Subtasks

- [x] Add past-date detection helper (AC: 1, 2)
  - [x] Added `isDateInPast(date: string): boolean` — simple string comparison `date < todayDate`
  - [x] Dates strictly before today are past; today itself is NOT past
- [x] Disable free-day checkboxes for past dates (AC: 1)
  - [x] Added `:disabled="isDateInPast(date)"` to free-day `v-checkbox` in both folded and
    expanded week views
  - [x] Disabled checkbox appears grayed out via `.week-day-past` class
  - [x] Added conditional `data-cy="week-day-checkbox-past"` when date is past
- [x] Prevent past dates from being toggled (AC: 1)
  - [x] Added early return `if (isDateInPast(date)) return` in `toggleWeekDay()`
  - [x] Safety check — checkbox is already disabled in the template
- [x] Apply visual graying (AC: 1)
  - [x] Added `.week-day-past` CSS class with `opacity: 0.5`
  - [x] Applied class to week-day-slot div when date is in the past
- [x] Updated E2E tests for past-date awareness
  - [x] Added `cy.clock()` to freeze time in week-booking tests that click checkboxes
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

All changes in `web/src/views/ItemsView.vue`. No backend changes required.

### Existing Date Reference

The `todayDate` constant is already defined at line 533:

```typescript
const todayDate = formatDate(new Date());
```

This returns a `YYYY-MM-DD` string. The `isDateInPast()` function can simply do a string
comparison: `date < todayDate`.

### Current Week Scenario

The current week (default selection) typically includes past dates (Mon-today) and future
dates (tomorrow-Fri). This is the primary use case for this feature. Future-only weeks
require no changes since all checkboxes are already interactive.

### No Impact on Other Statuses

Only "free" day checkboxes need the past-date disable. Days that are "booked-by-me",
"booked-by-other", or "unavailable" are already disabled by their respective checkbox
configurations (lines 366-384).

### References

- Epic 15 Story 15.4: `_bmad-output/planning-artifacts/epics.md` (Epic 15 Stories section)
- FR53: `_bmad-output/planning-artifacts/prd.md`
- ItemsView week checkboxes: `web/src/views/ItemsView.vue` lines 346-401
- todayDate: `web/src/views/ItemsView.vue` line 533

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `isDateInPast(date: string): boolean` helper using string comparison against `todayDate`
- Free-day checkboxes in both folded and expanded week views disabled for past dates
- Safety guard in `toggleWeekDay()` prevents toggling past dates programmatically
- `.week-day-past` CSS class applies `opacity: 0.5` to past date slots
- Conditional `data-cy="week-day-checkbox-past"` attribute for past-date checkboxes
- Updated 2 E2E tests with `cy.clock()` to freeze time to Monday for reliable checkbox interaction
- All 138 unit tests pass, all 51 E2E tests pass
- Code review fix: gray out past-date status text to match disabled state

### Change Log

- 2026-02-14: Implemented Story 15.4 - disable past date checkboxes in week mode
- 2026-02-14: Code review fix for past-date visual state

### File List

- web/src/views/ItemsView.vue (modified - isDateInPast helper, disabled checkboxes, CSS)
- web/cypress/e2e/week-booking.cy.ts (modified - added cy.clock for past-date awareness)
- web/src/views/ItemsView.vue (modified - past-date status styling)