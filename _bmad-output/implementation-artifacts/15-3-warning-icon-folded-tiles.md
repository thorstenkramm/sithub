# Story 15.3: Warning Icon on Folded Tiles

Status: done

## Story

As a user,
I want to know about item warnings even when a tile is folded,
So that I can make informed decisions without expanding every tile.

## Acceptance Criteria

1. **Given** a tile is folded and the item has a warning
   **When** I see the tile header
   **Then** a warning icon is visible

2. **Given** I click the warning icon on a folded tile
   **When** the popup or tooltip appears
   **Then** the full warning message is displayed
   **And** I do not need to unfold the tile to read the warning

3. **Given** a tile is folded and the item has no warning
   **When** I see the tile header
   **Then** no warning icon appears

## Tasks / Subtasks

- [x] Add warning icon to day mode folded tiles (AC: 1, 2, 3)
  - [x] Added `mdi-alert` icon with warning color in day tile #append slot
  - [x] Wrapped in `v-tooltip` showing full warning text
  - [x] Only visible when tile is folded AND has a warning attribute
  - [x] Added `data-cy="folded-warning-icon"` for testing
- [x] Add warning icon to week mode folded tiles (AC: 1, 2, 3)
  - [x] Added `mdi-alert` icon in week tile #append slot
  - [x] Uses `getWeekItemAttributes()` to get warning from any day's item data
  - [x] Wrapped in `v-tooltip` with full warning text
  - [x] Added `data-cy="week-folded-warning-icon"` for testing
- [x] Hide warning icon when tile is expanded (AC: 1)
  - [x] Icon hidden via `!expandedDayTiles.has()` / `!expandedWeekTiles.has()` condition
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

All changes in `web/src/views/ItemsView.vue`. No backend changes required.

### Dependency on Stories 15.1 and 15.2

This story depends on the fold/unfold mechanism introduced in Stories 15.1 (week mode) and
15.2 (day mode). The `expandedWeekTiles` and `expandedDayTiles` Sets must exist before this
story can add conditional warning icons.

### Warning Data Availability

In day mode, `entry.attributes.warning` is directly available on each item.

In week mode, the warning is available per-day in `weekData[date]` for each item. Since
warnings are item properties (not date-dependent), use the first available day's data:

```typescript
const getWeekItemWarning = (itemId: string): string | undefined => {
  for (const dayItems of Object.values(weekData.value)) {
    const item = dayItems.find(i => i.id === itemId);
    if (item?.attributes.warning) return item.attributes.warning;
  }
  return undefined;
};
```

### References

- Epic 15 Story 15.3: `_bmad-output/planning-artifacts/epics.md` (Epic 15 Stories section)
- FR52: `_bmad-output/planning-artifacts/prd.md`
- ItemsView: `web/src/views/ItemsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `mdi-alert` warning icon with `v-tooltip` in day tile #append slot (folded only)
- Added `mdi-alert` warning icon with `v-tooltip` in week tile #append slot (folded only)
- `getWeekItemAttributes()` helper used to extract warning from week data
- Icons automatically hidden when tile is expanded (full warning alert visible instead)
- Added `v-tooltip` to test stubs in unit tests
- All 138 unit tests pass, all 51 E2E tests pass
- Code review fix: make folded warning icons focusable with accessible labels
- Code review fix: add unit tests for folded warning icon visibility

### Change Log

- 2026-02-14: Implemented Story 15.3 - warning icon on folded tiles
- 2026-02-14: Code review fixes for folded warning icon accessibility and tests

### File List

- web/src/views/ItemsView.vue (modified - warning icons in day and week tile headers)
- web/src/views/ItemsView.test.ts (modified - added v-tooltip to stubs)
- web/src/views/ItemsView.vue (modified - warning icon buttons with aria labels)
- web/src/views/ItemsView.test.ts (modified - warning icon tests)