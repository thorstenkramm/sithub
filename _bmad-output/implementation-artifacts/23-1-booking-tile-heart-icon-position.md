# Story 23.1: Booking Tile Heart Icon Position

Status: done

## Story

As a user,
I want the favorite heart icon correctly positioned on booking tiles,
so that the tile layout is clean and consistent in both day and week modes.

## Acceptance Criteria

1. **Given** I view items in day booking mode
   **When** a tile renders with the heart/favorite icon
   **Then** the heart icon is on the second line (status row) after the
   availability chip, before the info/warning icon and chevron

2. **Given** I view items in week booking mode
   **When** a tile renders with the heart/favorite icon
   **Then** the heart icon is in the same position as day mode (second line,
   after availability chip)

3. **Given** the tile is rendered in either mode
   **When** I inspect the layout
   **Then** a unit test verifies the heart icon is inside the status row,
   not in v-card-actions

## Tasks / Subtasks

- [x] Task 1: Write failing test for heart icon position (AC: 3)
  - [x] 1.1 In `web/src/views/ItemsView.test.ts`: add test that asserts the
    heart icon (`[data-cy="item-favorite-heart"]`) is inside the subtitle/status
    row, NOT inside `[data-cy="day-item-actions"]`. Test must fail first.
  - [x] 1.2 Add same test for week mode: assert `[data-cy="week-item-favorite-heart"]`
    is inside the week status row, NOT inside `[data-cy="week-item-actions"]`.
- [x] Task 2: Move heart icon in day mode (AC: 1)
  - [x] 2.1 In `web/src/views/ItemsView.vue` day mode template (~lines 219-258):
    move the heart `v-btn` from `v-card-actions` (~line 352) into the subtitle
    row, between the warning icon and the `v-spacer` that precedes the chevron.
  - [x] 2.2 Remove the now-empty `v-card-actions` section for available items
    (or keep only for book/cancel button if needed). If the book button remains,
    keep `v-card-actions` but without the heart.
  - [x] 2.3 Verify tests from Task 1 now pass for day mode.
- [x] Task 3: Move heart icon in week mode (AC: 2)
  - [x] 3.1 In `web/src/views/ItemsView.vue` week mode template (~lines 408-451):
    move the heart `v-btn` from `v-card-actions` (~line 646) into the availability
    row, same position as day mode.
  - [x] 3.2 Remove the now-empty `v-card-actions` section in week mode (it only
    contained the heart).
  - [x] 3.3 Verify tests from Task 1 now pass for week mode.
- [x] Task 4: Run full test suite and linters (AC: 1, 2, 3)
  - [x] 4.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`,
    `npm run build`
  - [x] 4.2 Fix any failures

## Dev Notes

### Current Problem

The heart icon is in `v-card-actions` (bottom of card) in both day and week modes.
The desired layout puts it on the **second line** (status row) alongside the
availability chip, warning icon, and chevron toggle.

### Desired Tile Layout (from epic-23.md)

1. **Line 1:** Item name (truncated from center if needed)
2. **Line 2:** Availability icon + heart icon + info/warning icon + chevron
3. **Line 3:** Equipment chips
4. **Line 4:** Day checkboxes (week mode) or book button (day mode)

### Current Structure

**Day mode** (`ItemsView.vue` ~lines 197-365):

```text
v-card-item (prepend avatar)
  v-card-title → item name
  subtitle row → status chip + warning icon + v-spacer + chevron
v-card-text → equipment, warning alert, booker name, note
v-card-actions → book/cancel btn + v-spacer + HEART ← WRONG
```

**Week mode** (`ItemsView.vue` ~lines 390-657):

```text
v-card-item (prepend avatar)
  v-card-title → item name
  availability row → chip + warning icon + v-spacer + chevron
v-card-text → equipment, week days grid, warning
v-card-actions → v-spacer + HEART ← WRONG
```

### Target Structure

Move heart button into the subtitle/availability row in both modes:

```text
subtitle row → status chip + HEART + warning icon + v-spacer + chevron
```

The heart `v-btn` stays identical (same icon, click handler, data-cy). Only
its location in the template changes.

### Key data-cy Selectors

- Day mode heart: `data-cy="item-favorite-heart"`
- Week mode heart: `data-cy="week-item-favorite-heart"`
- Day actions container: `data-cy="day-item-actions"`
- Week actions container: `data-cy="week-item-actions"`

### Existing Tests to Update

`web/src/views/ItemsView.test.ts`:
- Lines ~496-513: day mode heart icon test (currently checks inside `day-item-actions`)
- Lines ~540-553: week mode heart icon test (currently checks inside `week-item-actions`)

These tests need to assert the heart is in the **status row**, not in
`v-card-actions`.

### Files to Modify

- `web/src/views/ItemsView.vue` — move heart button in template (both modes)
- `web/src/views/ItemsView.test.ts` — update heart icon position tests

### Do NOT Change

- Heart icon behavior (toggle, color, click handler, data-cy attribute)
- `useFavorites` composable — it works correctly
- Any other tile content or ordering beyond the heart position

### References

- [Source: private/epic-23.md — "Booking tile" section with desired layout]
- [Source: web/src/views/ItemsView.vue — current implementation]
- [Source: web/src/views/ItemsView.test.ts — existing heart icon tests]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Task 1: Wrote failing tests asserting heart icon is inside `[data-cy="day-status-row"]`
  and `[data-cy="week-status-row"]` respectively, and NOT in `v-card-actions`. Both tests
  failed as expected (RED phase).
- Task 2: Moved heart `v-btn` from `v-card-actions` into the day-mode status row div
  (after StatusChip, before warning icon). Added `data-cy="day-status-row"`. Removed heart
  from `v-card-actions` (kept book/cancel buttons). Day test passes.
- Task 3: Moved heart `v-btn` from `v-card-actions` into the week-mode status row div
  (after availability chip, before warning icon). Added `data-cy="week-status-row"`.
  Removed the now-empty week `v-card-actions` entirely. Week test passes.
- Task 4: All 280 tests pass. Type-check, ESLint, build all clean.

### File List

- `web/src/views/ItemsView.vue` (modified — moved heart icon to status row in both modes)
- `web/src/views/ItemsView.test.ts` (modified — updated heart position assertions)

### Review Findings

- [x] [Review][Patch] Empty v-card-actions gap for non-admin occupied items in day mode — conditionally render v-card-actions only when a button is visible [web/src/views/ItemsView.vue:~337]
- [x] [Review][Patch] Tautological week-mode negative test asserts heart absent from removed container — assert container itself does not exist [web/src/views/ItemsView.test.ts:~554]
- [x] [Review][Defer] Week mode cards lose bottom-alignment CSS hook after v-card-actions removal — deferred, pre-existing grid alignment pattern
