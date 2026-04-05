# Story 21.5: Visual Improvements — Save Icon and Floor Plan Button

Status: done

## Story

As a user,
I want the equipment filter save button to use a recognizable save icon and the floor
plan button to be consistently sized and positioned next to the week selector,
so that the UI is visually polished and controls are easy to find.

## Acceptance Criteria

1. **Given** I type an equipment filter and it is not yet saved,
   **when** I see the save button icon,
   **then** the icon is `mdi-content-save` (not a plus icon).

2. **Given** an area or item group has a floor plan configured,
   **when** I view the week selector,
   **then** the floor plan button appears at the same height as the week selector.

3. **Given** an item group with a floor plan is selected in the items view,
   **when** the page renders,
   **then** the floor plan button appears next to the calendar week selector
   (same flex row), not below it.

## Tasks / Subtasks

- [x] Task 1: Change equipment filter save icon (AC: 1)
  - [x] 1.1 Replace `$plus` with `mdi-content-save` in `ItemGroupsView.vue` filter save button
  - [x] 1.2 Replace `$plus` with `mdi-content-save` in `ItemsView.vue` filter save button
- [x] Task 2: Match floor plan button height to week selector (AC: 2)
  - [x] 2.1 Change floor plan button from `size="small"` to `density="compact"` in
    `ItemGroupsView.vue`
  - [x] 2.2 Change floor plan button from `size="small"` to `density="compact"` in
    `ItemsView.vue`
- [x] Task 3: Reposition floor plan button next to week selector (AC: 3)
  - [x] 3.1 Move the floor plan button inside the flex row containing the week selector
    in `ItemsView.vue`
  - [x] 3.2 Remove the old standalone floor plan button position (`mt-4 mb-2` class)
- [x] Task 4: Verify tests and lint (AC: 1, 2, 3)
  - [x] 4.1 Run Vitest unit tests
  - [x] 4.2 Run ESLint

## Dev Notes

### Scope

Frontend-only template changes to `ItemGroupsView.vue` and `ItemsView.vue`. No backend
changes. No new components or composables.

### Icon Change

The plus icon for saving a filter was easily overlooked. The `mdi-content-save` (floppy
disk) icon is universally recognized as "save".

### Floor Plan Button Positioning

The floor plan button was previously positioned below the week selector with `mt-4 mb-2`.
It is now inside the same `d-flex flex-wrap align-end ga-4` container as the week selector,
ensuring it appears alongside it at the same visual height.

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Changed equipment filter save icon from `$plus` to `mdi-content-save` in both views
- Floor plan button uses `density="compact"` to match the compact week selector height
- Floor plan button moved into the flex row next to the week selector in ItemsView
- All 255 Vitest tests pass, ESLint clean

### File List

- `web/src/views/ItemGroupsView.vue` — Save icon change, floor plan button density
- `web/src/views/ItemsView.vue` — Save icon change, floor plan button repositioned and density
