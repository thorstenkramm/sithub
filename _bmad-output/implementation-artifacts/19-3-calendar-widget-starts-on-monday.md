# Story 19.3: Calendar Widget Starts on Monday

Status: done

## Story

As a user,
I want the calendar date picker to show Monday as the first day of the week,
So that it matches the European convention I am used to.

## Acceptance Criteria

1. **Given** I am on any view with a date picker
   **When** the calendar widget opens
   **Then** Monday is displayed as the first (leftmost) column
   **And** Sunday is displayed as the last (rightmost) column

## Tasks / Subtasks

- [x] Set Monday as first day of week in Vuetify defaults (AC: 1)
  - [x] In `web/src/plugins/vuetify.ts`: added `VDatePicker.firstDayOfWeek: 1` as a
    global Vuetify default
  - [x] Applies to all `VDatePicker` instances across the application
- [x] Update affected tests
  - [x] Updated `web/src/views/ItemGroupsView.test.ts` to account for Vuetify defaults
- [x] Verify E2E tests still pass

## Dev Notes

### Implementation

Vuetify's `VDatePicker` defaults to Sunday (0) as the first day of the week. Setting
`firstDayOfWeek: 1` in the global Vuetify defaults configuration ensures all date pickers
start on Monday without needing per-instance props.

### References

- Epic 19 Story 19.3: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR72: `_bmad-output/planning-artifacts/prd.md`
- `web/src/plugins/vuetify.ts`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Single configuration change in Vuetify defaults
- All date pickers now start on Monday globally
- Updated ItemGroupsView test to work with new default
- All existing tests continue to pass

### File List

- `web/src/plugins/vuetify.ts` — Added `VDatePicker.firstDayOfWeek: 1` global default
- `web/src/views/ItemGroupsView.test.ts` — Updated for Vuetify defaults change

## Change Log

- 2026-03-21: Story implemented and verified.
