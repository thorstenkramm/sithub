# Story 22.6: Favorites Heart Icon Visibility Fix

Status: done

## Story

As a user,
I want to see the favorite heart icon on all item tiles,
so that I can manage my favorites regardless of other badges shown.

## Acceptance Criteria

1. **Given** any item tile (day mode, week mode, or item group view)
   **When** the tile renders
   **Then** the favorite heart icon is in the same position across all views:
   in `v-card-actions` row, right-aligned (matching the ItemGroupsView pattern)

2. **Given** a tile has both warning badge and favorite heart
   **When** the tile renders on a narrow mobile screen
   **Then** both are visible — heart in card-actions, warning in #append slot

## Tasks / Subtasks

- [ ] Task 1: Investigate and fix heart visibility on day mode tiles (AC: 1)
  - [ ] 1.1 In `web/src/views/ItemsView.vue` lines ~209-220: the heart button
    is in `v-card-title` with `class="ml-1"`. The warning icon is in the
    `#append` template slot (lines 222-257). On narrow screens, the flex
    row may push the heart off-screen when the item name is long.
    Verify by testing with Chrome DevTools on iPhone 14
  - [ ] 1.2 If the heart is hidden due to flex overflow, move it to a
    dedicated row below the title, or ensure `flex-wrap` allows it to
    wrap to a new line
  - [ ] 1.3 Check z-index: the `.item-filtered-overlay` uses `z-index: 1`
    (line 1764). Ensure the heart button is not obscured by this overlay
    when equipment filter is active
- [ ] Task 2: Verify heart visibility on week mode tiles (AC: 2)
  - [ ] 2.1 In `web/src/views/ItemsView.vue` lines ~386-397: the week mode
    heart is also in `v-card-title` with `class="ml-1"`. The warning icon
    is in `#append` (lines 401-421). Same potential flex overflow issue
  - [ ] 2.2 Test with long item names and both icons present
- [ ] Task 3: Verify on ItemGroupsView (AC: 1)
  - [ ] 3.1 In `web/src/views/ItemGroupsView.vue` lines 159-167 (favorite
    items) and 244-254 (item groups): these have heart in `v-card-actions`
    with `v-spacer` — different layout, less likely to overflow. Verify
- [ ] Task 4: Run tests and lint (AC: 1, 2)
  - [ ] 4.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`
  - [ ] 4.2 Visual verification with Chrome DevTools MCP on iPhone 14

## Dev Notes

### Root Cause Hypothesis

The heart icon is in a `d-flex align-center` row with the item name.
When the item name is long and the warning icon is in the `#append` slot
(right side), the flex container may overflow on a 390px screen, pushing
the heart beyond the visible area. The heart `ml-1` margin doesn't help
if there's no room.

### Possible Fix

Option A: Move heart icon to `v-card-actions` (like ItemGroupsView does).
This gives it a dedicated row and guaranteed visibility.

Option B: Keep in title but ensure `flex-shrink: 0` on the heart button
so it never gets squeezed by the flex layout.

### Screenshot Reference

See `private/epic-22.md`: "The heart icon from the favorites feature is
not visible." with attached screenshot.

### Files to Change

- `web/src/views/ItemsView.vue` — day mode (lines ~209-220) and week mode
  (lines ~386-397) heart icon positioning

### References

- [Source: private/epic-22.md — "heart icon not visible"]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

### Review Findings

- [x] [Review][Patch] Day and week item tiles now place the favorite heart in the right-aligned `v-card-actions` row, matching `ItemGroupsView`, with regression coverage for both day and week layouts [web/src/views/ItemsView.vue:214]
