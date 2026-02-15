# Story 15.1: Collapsible Tiles in Week Booking Mode

Status: done

## Story

As a user,
I want to expand item tiles in week mode to see full details,
So that the default view is compact and I can drill into specifics on demand.

## Acceptance Criteria

1. **Given** I am in week booking mode viewing item tiles
   **When** I see a tile
   **Then** a chevron-left icon appears in the tile header

2. **Given** I click the chevron on a folded tile
   **When** the tile unfolds
   **Then** the chevron rotates to chevron-down
   **And** the compact M-F row is replaced by one line per day
   **And** each line shows the full day name (Monday, Tuesday, etc.)
   **And** each line shows the full booker display name (not truncated)
   **And** equipment chips and warning alerts are visible below the daily breakdown

3. **Given** I click the chevron on an unfolded tile
   **When** the tile folds
   **Then** the chevron rotates back to chevron-left
   **And** the compact M-F row is restored

4. **Given** I am viewing a folded tile with truncated booker names
   **When** I hover over a truncated name
   **Then** a tooltip shows the full display name

## Tasks / Subtasks

- [x] Add fold/unfold state management (AC: 1, 2, 3)
  - [x] Add a `expandedWeekTiles` reactive Set to track which item IDs are expanded
  - [x] Add `toggleWeekTileExpansion(itemId: string)` function
  - [x] Reset expanded state on week data reload
- [x] Add chevron icon to week tile header (AC: 1, 3)
  - [x] In `ItemsView.vue` week tile `v-card-item`: added chevron button in `#append` slot
  - [x] Use `mdi-chevron-left` when folded, `mdi-chevron-down` when expanded
  - [x] Add `data-cy="week-tile-chevron"` for testing
- [x] Implement expanded view layout (AC: 2)
  - [x] Conditional rendering: folded (compact M-F row) vs expanded (vertical one-line-per-day)
  - [x] Each expanded row: full day name, checkbox, full booker name or "free"
  - [x] Equipment chips and warning alert shown below daily rows in expanded view
  - [x] `getWeekItemAttributes()` helper extracts equipment/warning from first available day
  - [x] Extracted `FULL_DAY_LABELS` constant, replaced local `dayLabels` in `submitWeekBookings`
- [x] Add tooltip to truncated names in folded view (AC: 4)
  - [x] Wrapped `booked-by-other` status span in `v-tooltip` for full name on hover
- [x] Update CSS for expanded layout
  - [x] Added `.week-day-expanded` and `.week-day-expanded-label` classes
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

All changes in `web/src/views/ItemsView.vue`. No backend changes required.

### Week Tile Structure (Current)

The week mode tiles (lines 327-404) currently show:
- `v-card-item` with avatar and item name
- `v-card-text` with a flex/grid layout of 5 day slots, each containing:
  - Day label (MO, TU, etc.)
  - Checkbox (free/booked-by-me/booked-by-other/unavailable)
  - Status text (free/name/n-a)

Equipment and warnings are NOT currently displayed in week mode tiles. The data is available
in `weekData` (per-day item attributes include `equipment` and `warning`).

### Equipment/Warning Source

Each day's `fetchItems()` response includes `equipment` and `warning` in item attributes.
Since equipment and warnings are item properties (not booking properties), they are the same
for all days. Use the first available day's data to display equipment and warning in the
expanded view.

### Day Labels for Expanded View

The `WEEKDAY_LABELS` array in `useWeekSelector.ts` has short labels (MO, TU, etc.). For the
expanded view, use full names: `['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday']`.
Note: the `submitWeekBookings()` function at line 722 already has a `dayLabels` array with
full names — consider extracting this to a shared constant.

### References

- Epic 15 Story 15.1: `_bmad-output/planning-artifacts/epics.md` (Epic 15 Stories section)
- FR50, FR54: `_bmad-output/planning-artifacts/prd.md`
- ItemsView week section: `web/src/views/ItemsView.vue` lines 326-404
- useWeekSelector: `web/src/composables/useWeekSelector.ts`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `expandedWeekTiles` reactive Set and `toggleWeekTileExpansion()` function
- Added chevron icon (mdi-chevron-left/down) in week tile #append slot
- Implemented folded/expanded conditional view: compact M-F row vs vertical full-day list
- `getWeekItemAttributes()` extracts equipment/warning from any day's data
- Expanded view shows equipment chips and warning alerts below daily breakdown
- Added `v-tooltip` on booked-by-other names in folded view for full name display
- Extracted `FULL_DAY_LABELS` constant shared with `submitWeekBookings()`
- Added `.week-day-expanded` and `.week-day-expanded-label` CSS classes
- All 138 unit tests pass, all 51 E2E tests pass
- Type check, ESLint, build, and code duplication checks all pass
- Code review fix: compute full day labels from dates to support locale and weekends
- Code review fix: show tooltips only when names are actually truncated in folded week view
- Code review fix: memoize week item attributes to avoid repeated scans

### Change Log

- 2026-02-14: Implemented Story 15.1 - collapsible tiles in week booking mode
- 2026-02-14: Code review fixes for week tile label, tooltip, and attribute lookup

### File List

- web/src/views/ItemsView.vue (modified - week tile expansion, chevron, expanded view)
- web/src/views/ItemsView.vue (modified - week tile labels, tooltip gating, attribute memoization)
- web/src/views/ItemsView.test.ts (modified - warning icon behavior tests)