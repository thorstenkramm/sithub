# Story 19.6: Equipment Filter Saving

Status: done

## Story

As a user,
I want to save my equipment filters for reuse,
So that I don't have to retype the same filter keywords every time I book.

## Acceptance Criteria

1. **Given** I have typed a filter into the equipment filter input
   **When** I click the save icon next to the input
   **Then** the filter is saved to browser local storage
   **And** a confirmation is shown

2. **Given** I have saved filters
   **When** I focus the equipment filter input
   **Then** a combobox dropdown shows my saved filters alongside free-text input

3. **Given** I select a saved filter from the combobox
   **When** the filter loads
   **Then** the save icon becomes a delete icon

4. **Given** I click the delete icon on a loaded saved filter
   **When** the deletion is confirmed
   **Then** the filter is removed from local storage
   **And** the input is cleared

5. **Given** I have no saved filters
   **When** the page loads
   **Then** the input behaves as a regular text field with no dropdown entries

## Tasks / Subtasks

- [x] Create `useSavedFilters` composable (AC: 1, 2, 4, 5)
  - [x] Created `web/src/composables/useSavedFilters.ts`
  - [x] Manages saved filters in browser local storage
  - [x] Exports `savedFilters`, `saveFilter()`, `deleteFilter()`, `isSavedFilter()`
- [x] Replace `v-text-field` with `v-combobox` on ItemsView (AC: 2, 3)
  - [x] Updated equipment filter input from `v-text-field` to `v-combobox`
  - [x] Dropdown shows saved filters when they exist
  - [x] Free-text input still works alongside saved items
- [x] Replace `v-text-field` with `v-combobox` on ItemGroupsView (AC: 2, 3)
  - [x] Same combobox upgrade as ItemsView
- [x] Add save/delete toggle icon with tooltips (AC: 1, 3, 4)
  - [x] When filter is new (unsaved): save icon shown with tooltip
  - [x] When filter is a saved filter: delete icon shown with tooltip
  - [x] Save click persists to local storage and shows confirmation
  - [x] Delete click removes from local storage and clears input
- [x] Add unit tests
  - [x] Created `web/src/composables/useSavedFilters.test.ts`
  - [x] Tests for save, delete, list, and persistence in localStorage
- [x] Verify E2E tests still pass

## Dev Notes

### Composable Design

The `useSavedFilters` composable wraps localStorage access with a reactive interface.
Saved filters are stored as a JSON array under a fixed key. The composable provides:

- `savedFilters` — reactive list of saved filter strings
- `saveFilter(filter)` — adds a filter to the list
- `deleteFilter(filter)` — removes a filter from the list
- `isSavedFilter(filter)` — checks if a filter is already saved

### Combobox Behavior

The `v-combobox` component allows both free-text input and selection from a dropdown list.
When no saved filters exist, it behaves identically to a text field. The save/delete toggle
uses tooltips to indicate the available action.

### References

- Epic 19 Story 19.6: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR69: `_bmad-output/planning-artifacts/prd.md`
- `web/src/composables/useSavedFilters.ts`
- `web/src/views/ItemsView.vue`
- `web/src/views/ItemGroupsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Created `useSavedFilters` composable with localStorage persistence
- Replaced `v-text-field` with `v-combobox` on both ItemsView and ItemGroupsView
- Save/delete toggle icon with tooltips on both views
- AI review fix: save/delete actions now show confirmation feedback in both views
- AI review fix: deleting a saved filter now clears the current combobox input
- Unit tests cover save, delete, list, and localStorage persistence
- All existing tests continue to pass

### File List

- `web/src/composables/useSavedFilters.ts` — Saved filters composable with localStorage
- `web/src/composables/useSavedFilters.test.ts` — Unit tests for saved filters
- `web/src/views/ItemsView.vue` — Combobox upgrade, save/delete toggle, confirmation feedback
- `web/src/views/ItemsView.test.ts` — View tests for saved-filter confirmations and input clearing
- `web/src/views/ItemGroupsView.vue` — Combobox upgrade, save/delete toggle, confirmation feedback
- `web/src/views/ItemGroupsView.test.ts` — View tests for saved-filter confirmations and input clearing

## Senior Developer Review (AI)

- Fixed missing AC behavior: save/delete actions now show confirmation messages in both booking views.
- Fixed missing AC behavior: deleting a saved filter now clears the current input in both booking views.
- Added targeted view tests for the saved-filter flows so the ACs are covered by Vitest.

## Change Log

- 2026-03-21: Story implemented and verified.
- 2026-03-21: Applied AI review fixes for saved-filter confirmation and delete-input clearing.
