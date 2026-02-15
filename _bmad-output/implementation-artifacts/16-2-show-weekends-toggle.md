# Story 16.2: Show Weekends Toggle

Status: done

## Story

As a user,
I want to optionally show weekends on booking pages,
So that I can book Saturday and Sunday if my workplace supports it.

## Acceptance Criteria

1. **Given** I click my username in the top right corner
   **When** I see the user menu
   **Then** I see a "Show weekends" checkbox (unchecked by default)

2. **Given** I enable "Show weekends"
   **When** I view the booking page in week mode
   **Then** Saturday and Sunday columns appear alongside Monday through Friday
   **And** my preference is persisted in localStorage

3. **Given** I enable "Show weekends"
   **When** I view the weekly availability indicators on item group tiles
   **Then** the indicators include Saturday and Sunday

4. **Given** I disable "Show weekends"
   **When** I view any booking page
   **Then** only Monday through Friday are shown

## Tasks / Subtasks

- [x] Add "Show weekends" checkbox to desktop user menu (AC: 1)
  - [x] Added `v-list-item` with `v-checkbox` and `data-cy="show-weekends-toggle"`
- [x] Add "Show weekends" checkbox to mobile drawer (AC: 1)
  - [x] Added matching checkbox with `data-cy="mobile-show-weekends-toggle"`
- [x] Create weekends preference composable (AC: 2, 4)
  - [x] Created `web/src/composables/useWeekendPreference.ts`
  - [x] localStorage key `sithub_show_weekends`, default `false`
  - [x] Exports reactive `showWeekends` ref with watch for persistence
- [x] Update `useWeekSelector` to support 7-day mode (AC: 2, 3)
  - [x] `getWeekdayDates()` accepts `includeWeekends` param (5 or 7 days)
  - [x] Extended `WEEKDAY_LABELS` to 7: `['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU']`
  - [x] Extended `WEEKDAY_LABELS_SHORT` to 7: `['M', 'T', 'W', 'T', 'F', 'S', 'S']`
  - [x] `useWeekSelector()` accepts optional `showWeekends` Ref
  - [x] `selectedWeekDates` computed uses `showWeekends` to decide 5 or 7 days
- [x] Update ItemsView week mode grid (AC: 2)
  - [x] CSS grid uses inline style for dynamic column count
  - [x] Extended `getFullDayLabel` fallback to include Saturday/Sunday
  - [x] Watch on `[selectedWeek, showWeekends]`
- [x] Update ItemGroupsView availability indicators (AC: 3)
  - [x] Backend availability handler updated: `?days=7` query parameter for 7-day mode
  - [x] Backend `weekdayDates()` accepts `count` parameter
  - [x] Backend `weekdayAbbreviation()` extended with explicit SA/SU cases
  - [x] Frontend `fetchWeeklyAvailability` passes `days=7` when weekends enabled
  - [x] Watch on `[selectedWeek, showWeekends]`
- [x] Add unit tests
  - [x] Frontend: `getWeekdayDates` returns 5 vs 7 dates, weekend labels, composable tests
  - [x] Backend: `TestWeekdayDatesFullWeek`, `TestAvailabilityHandlerReturnsSevenDays`,
    `weekdayAbbreviation` SA/SU assertions
  - [x] `useWeekendPreference.test.ts`: localStorage persistence, default false
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend + Possible Backend Changes

This story primarily changes frontend composables and views. The backend availability API
may need updating to return 7-day data when weekends are requested.

### useWeekSelector Composable Changes

Key changes in `web/src/composables/useWeekSelector.ts`:

- `getWeekdayDates(monday, includeWeekends?)`: change loop from `i < 5` to
  `i < (includeWeekends ? 7 : 5)`
- Labels arrays: extend from 5 to 7 entries
- `getWeekdayLabel(index)`: already handles arbitrary indices via `labels[index] ?? ''`

### Backend Availability API

The availability endpoint handler in `internal/areas/presence_handler.go` or the bookings
store may filter to Mon-Fri. If so, it needs a parameter to include weekends. Check the
`getWeekdayDates` usage in the backend to understand the current behavior.

Alternatively, the frontend can always request 7 dates and the backend returns whatever
data exists — if no bookings exist for Saturday/Sunday, availability = total (all available).

### CSS Grid Consideration

The `week-days-compact` class uses `repeat(5, 1fr)`. This should become dynamic:

```css
.week-days-compact-7 {
  grid-template-columns: repeat(7, 1fr);
}
```

Or use inline style binding: `:style="{ gridTemplateColumns: \`repeat(${dayCount}, 1fr)\` }"`

### References

- Epic 16 Story 16.2: `_bmad-output/planning-artifacts/epics.md` (Epic 16 Stories section)
- FR56: `_bmad-output/planning-artifacts/prd.md`
- useWeekSelector: `web/src/composables/useWeekSelector.ts`
- ItemsView week grid: `web/src/views/ItemsView.vue` lines 346-404
- ItemGroupsView indicators: `web/src/views/ItemGroupsView.vue` lines 70-88
- App.vue user menu: `web/src/App.vue` lines 40-79

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Full-stack story: frontend composables, views, and backend handler/tests updated
- Backend supports `?days=7` query param on availability endpoint
- `weekdayDates()` now takes explicit `count` parameter (breaking change, all callers updated)
- Fixed pre-existing `FULL_DAY_LABELS` undefined bug in ItemsView (from Epic 15)
- Fixed pre-existing test state leak in ItemsView.test.ts (`sithub_booking_mode` localStorage)
- Fixed `v-tooltip` and `v-card-item` stubs to render named slots (activator, append)
- Updated `ItemGroupsView.test.ts` to expect third `undefined` arg in availability call
- Added safe localStorage access for weekend and booking mode preferences

### File List

- `web/src/composables/storage.ts` — Safe localStorage helper
- `web/src/composables/useWeekendPreference.ts` — New composable
- `web/src/composables/useWeekendPreference.test.ts` — 4 unit tests
- `web/src/composables/useWeekSelector.ts` — Extended for 7-day support
- `web/src/composables/useWeekSelector.test.ts` — Added 7-day and weekend label tests
- `web/src/views/ItemsView.vue` — Weekend support in week mode grid
- `web/src/views/ItemsView.test.ts` — Fixed stubs and test state leak
- `web/src/views/ItemGroupsView.vue` — Weekend support in availability indicators
- `web/src/views/ItemGroupsView.test.ts` — Updated availability call expectation
- `web/src/api/itemGroupAvailability.ts` — Added `days` parameter
- `web/src/App.vue` — Show weekends checkbox in desktop and mobile menus
- `internal/itemgroups/availability_handler.go` — `?days=7` support, SA/SU abbreviations
- `internal/itemgroups/availability_handler_test.go` — 7-day and SA/SU tests
