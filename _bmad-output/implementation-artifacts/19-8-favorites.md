# Story 19.8: Favorites

Status: done

## Story

As a user,
I want to mark item groups and items as favorites,
So that my most-used spaces appear first and are quick to find.

## Acceptance Criteria

1. **Given** I am on the item-groups view (second level)
   **When** I see an item group tile
   **Then** a heart outline icon is visible on the tile

2. **Given** I click the heart outline on an item group
   **When** the favorite is saved
   **Then** a confirmation "{item group name} saved as favorite." is shown
   **And** the icon becomes a red-filled heart
   **And** the favorite is persisted in browser local storage

3. **Given** I click a red-filled heart on an item group
   **When** the favorite is removed
   **Then** a confirmation "{item group name} removed from favorites." is shown
   **And** the icon reverts to a heart outline

4. **Given** I am on the items view (third level)
   **When** I see an item tile
   **Then** a heart outline icon is visible on the tile
   **And** clicking it saves/removes the favorite with confirmation
   "{item group name} {item name} saved/removed as favorite."

5. **Given** I have third-level favorites
   **When** I view the item-groups page (second level)
   **Then** my third-level favorites appear as bookable tiles on that page

6. **Given** I am on the item-groups view with favorites
   **When** the page renders
   **Then** items are ordered: (1) third-level favorites A-Z,
   (2) second-level favorites A-Z, (3) remaining item groups in YAML order
   with second-level favorites subtracted

## Tasks / Subtasks

- [x] Create `useFavorites` composable (AC: 2, 3, 4)
  - [x] Created `web/src/composables/useFavorites.ts`
  - [x] Manages favorites in browser local storage
  - [x] Supports both item group and item favorites
  - [x] Exports `toggleFavorite()`, `isFavorite()`, `getFavorites()`
- [x] Add heart icons to item group tiles (AC: 1, 2, 3)
  - [x] Updated `ItemGroupsView.vue` with heart outline/filled toggle on each tile
  - [x] Clicking toggles favorite state with snackbar confirmation
- [x] Add heart icons to item tiles (AC: 4)
  - [x] Updated `ItemsView.vue` with heart outline/filled toggle on each tile
  - [x] Confirmation message includes item group name and item name
- [x] Promote third-level favorites to second-level view (AC: 5)
  - [x] In `ItemGroupsView.vue`: third-level (item) favorites appear as bookable tiles
    on the item-groups page
- [x] Implement sorting logic (AC: 6)
  - [x] Sort order: (1) third-level favorites A-Z, (2) second-level favorites A-Z,
    (3) remaining item groups in YAML order
  - [x] Second-level favorites removed from "remaining" group to avoid duplication
- [x] Register heart icons in Vuetify (AC: 1)
  - [x] Added heart outline and filled heart icons to Vuetify aliases in `vuetify.ts`
- [x] Add snackbar confirmations (AC: 2, 3, 4)
  - [x] Snackbar shows on favorite add/remove with entity name
- [x] Add unit tests
  - [x] Created `web/src/composables/useFavorites.test.ts`
  - [x] Tests for toggle, persistence, retrieval, and item vs item group favorites
- [x] Verify E2E tests still pass

## Dev Notes

### Local Storage Structure

Favorites are stored in localStorage as JSON. Item group favorites and item favorites are
tracked separately, allowing the composable to distinguish between second-level and
third-level favorites for the promotion and sorting logic.

### Sorting Algorithm

The item-groups view sorts tiles in three tiers:

1. Third-level (item) favorites promoted to second level, sorted alphabetically
2. Second-level (item group) favorites, sorted alphabetically
3. Remaining item groups in their original YAML order, excluding any that appear as
   second-level favorites

### Snackbar Confirmations

Both views show a snackbar message when a favorite is toggled. The message format differs:

- Item groups: "{name} saved as favorite." / "{name} removed from favorites."
- Items: "{group name} {item name} saved as favorite." / similar for removal

### References

- Epic 19 Story 19.8: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR74: `_bmad-output/planning-artifacts/prd.md`
- `web/src/composables/useFavorites.ts`
- `web/src/views/ItemGroupsView.vue`
- `web/src/views/ItemsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Created `useFavorites` composable with localStorage persistence
- Heart icons on both ItemGroupsView and ItemsView tiles
- Third-level favorites promoted to second-level view as bookable tiles
- Three-tier sorting: promoted item favorites, group favorites, remaining
- AI review fix: favorites are now scoped by area and item-group path to avoid ID collisions
- Snackbar confirmations on all favorite toggles
- Registered heart icons in Vuetify aliases
- All existing tests continue to pass

### File List

- `web/src/composables/useFavorites.ts` — Favorites composable with area-scoped and item-group-scoped persistence
- `web/src/composables/useFavorites.test.ts` — Unit tests for scoped favorites persistence and collision prevention
- `web/src/views/ItemGroupsView.vue` — Heart icons, promoted favorites, area-scoped sorting
- `web/src/views/ItemsView.vue` — Heart icons, favorite toggle with scoped persistence
- `web/src/plugins/vuetify.ts` — Heart icon aliases

## Senior Developer Review (AI)

- Fixed a storage bug where repeated item-group or item IDs from different branches of the hierarchy could overwrite each other.
- Scoped favorite persistence by area and item-group path, and added tests to prevent regressions.

## Change Log

- 2026-03-21: Story implemented and verified.
- 2026-03-21: Applied AI review fixes for scoped favorites persistence and collision handling.
