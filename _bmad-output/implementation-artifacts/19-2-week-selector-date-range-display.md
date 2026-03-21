# Story 19.2: Week Selector Date Range Display

Status: done

## Story

As a user,
I want the calendar week selector to show both the first and last day of each week,
So that I can immediately see which date range a calendar week covers.

## Acceptance Criteria

1. **Given** I am on a view with the week selector
   **When** I open the week selector dropdown
   **Then** each option shows the format "DD.MM.-DD.MM.YYYY - Week N"
   (e.g. "23.03.-29.03.2026 - Week 13")

2. **Given** the show weekends toggle is off
   **When** I view the week selector
   **Then** the date range still shows Monday through Sunday (full week),
   regardless of the weekends setting

## Tasks / Subtasks

- [x] Update week selector label format (AC: 1, 2)
  - [x] In `useWeekSelector.ts`: replace `Intl.DateTimeFormat` formatting with
    `formatDayMonth` and `formatDayMonthYear` helper functions
  - [x] Label format: "DD.MM.-DD.MM.YYYY - Week N" (e.g. "23.03.-29.03.2026 - Week 13")
  - [x] Always show Monday through Sunday range regardless of weekends toggle
- [x] Add unit tests
  - [x] Updated `useWeekSelector.test.ts` with assertions for new date range format
  - [x] Verified format matches "DD.MM.-DD.MM.YYYY - Week N" pattern
- [x] Verify E2E tests still pass

## Dev Notes

### Formatting Approach

Replaced `Intl.DateTimeFormat` usage with custom `formatDayMonth` and `formatDayMonthYear`
helper functions to produce the exact "DD.MM." and "DD.MM.YYYY" format. The week label
always uses the full Monday-to-Sunday range so users see the complete week span even when
weekends are hidden.

### References

- Epic 19 Story 19.2: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR71: `_bmad-output/planning-artifacts/prd.md`
- `web/src/composables/useWeekSelector.ts`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Replaced `Intl.DateTimeFormat` with `formatDayMonth`/`formatDayMonthYear` helpers
- Week label always shows full Monday-Sunday range
- Updated unit tests to verify new format
- All existing tests continue to pass

### File List

- `web/src/composables/useWeekSelector.ts` — New format helpers, updated label computation
- `web/src/composables/useWeekSelector.test.ts` — Updated tests for date range format

## Change Log

- 2026-03-21: Story implemented and verified.
