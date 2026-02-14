# Story 14.1: Simplify Action Labels Across Views

Status: done

## Story

As a user,
I want concise action labels that get me to my destination faster,
So that I spend less time reading and more time booking.

## Acceptance Criteria

1. **Given** I am viewing the areas list (Home)
   **When** I see an area tile
   **Then** the action button reads "Select" instead of "View Item Groups"

2. **Given** I am viewing item groups within an area
   **When** I see the page
   **Then** the page title "Item Groups" and subtitle "Select an item group to view available
   items" are removed
   **And** the action button on each tile reads "Select" instead of "View Items"

3. **Given** I am viewing items within an item group
   **When** I see the page
   **Then** the page title "Items" and subtitle "Select an item to book for your chosen date"
   are removed

4. **Given** I am viewing available items in day booking mode
   **When** I see an available item tile
   **Then** the booking button reads "Book" instead of "Book This Item"

## Tasks / Subtasks

- [x] Rename area tile button (AC: 1)
  - [x] In `AreasView.vue` line 57: change "View Item Groups" to "Select"
- [x] Remove ItemGroupsView page header text (AC: 2)
  - [x] In `ItemGroupsView.vue` line 4: remove `title="Item Groups"` (pass empty string)
  - [x] In `ItemGroupsView.vue` line 5: remove `subtitle` prop entirely
- [x] Rename item group tile button (AC: 2)
  - [x] In `ItemGroupsView.vue` line 97: change "View Items" to "Select"
- [x] Remove ItemsView page header text (AC: 3)
  - [x] In `ItemsView.vue` line 4: remove `title="Items"` (pass empty string)
  - [x] In `ItemsView.vue` line 5: remove `subtitle` prop entirely
- [x] Rename booking button (AC: 4)
  - [x] In `ItemsView.vue` line 305: change "Book This Item" to "Book"
- [x] Update unit tests (AC: 2)
  - [x] In `ItemGroupsView.test.ts` line 107: update assertion that checks for "Item Groups"
- [x] Verify E2E tests still pass (AC: 1, 2, 3, 4)
  - [x] Run `npm run test:e2e -- --browser electron` — no text assertions to update
    (tests use `data-cy` selectors)

## Dev Notes

### Architecture: Frontend-Only Story

This story requires NO backend changes. All changes are label text modifications in three
Vue view components. The scope is deliberately minimal — touch only the specified strings.

### PageHeader Component

The `PageHeader` component (`web/src/components/PageHeader.vue`) accepts `title` and
`subtitle` props. When `title` is empty string and `subtitle` is not provided, the
component renders only the breadcrumbs and the action slot. The `<h1>` tag renders
unconditionally on line 37 — to avoid an empty heading element, pass `title` as empty
string or consider removing the title/subtitle props entirely for these views. Verify the
rendered output does not leave empty heading tags for accessibility.

### Existing Test Coverage

- `ItemGroupsView.test.ts` line 107 asserts `expect(wrapper.text()).toContain('Item Groups')`
  — this must be updated or removed.
- `AreasView.test.ts` asserts `toContain('Areas')` — NOT affected (Areas title stays).
- `ItemsView.test.ts` has no assertions on page title text.
- All Cypress E2E tests use `data-cy` selectors, not button text — no E2E changes needed.

### Breadcrumbs Are Not Affected

Breadcrumbs display area and item group *names* from API data (e.g., "Office 2nd Floor"),
not the page titles being removed. Breadcrumb behavior is unchanged by this story.

### Project Structure Notes

- All changes in `web/src/views/` — standard Vue SFC files
- No new files, no new dependencies, no API changes
- Consistent with Epic 12 terminology (domain-neutral labels)

### References

- Epic 14 Story 14.1: `_bmad-output/planning-artifacts/epics.md` (Epic 14 Stories section)
- FR43: `_bmad-output/planning-artifacts/prd.md` (Navigation & UI Consistency)
- AreasView: `web/src/views/AreasView.vue`
- ItemGroupsView: `web/src/views/ItemGroupsView.vue`
- ItemsView: `web/src/views/ItemsView.vue`
- PageHeader: `web/src/components/PageHeader.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Changed "View Item Groups" to "Select" in AreasView.vue
- Removed title and subtitle from ItemGroupsView PageHeader (pass empty string)
- Changed "View Items" to "Select" in ItemGroupsView.vue
- Removed title and subtitle from ItemsView PageHeader (pass empty string)
- Changed "Book This Item" to "Book" in ItemsView.vue
- Added `v-if="title"` to PageHeader `<h1>` to prevent empty heading element (accessibility)
- Updated ItemGroupsView.test.ts: replaced title text assertion with breadcrumb existence check
- All 128 unit tests pass, all 51 E2E tests pass
- Type check, ESLint, build, and code duplication checks all pass
- Code review fix: update area and item group card aria-labels to match "Select" action text
- Code review fix: add unit tests to assert "Select" labels on area and item group tiles

### Change Log

- 2026-02-14: Implemented Story 14.1 - simplified action labels across 3 views, fixed
- 2026-02-14: Code review fixes for aria-label consistency and tests
  PageHeader accessibility for empty titles

### File List

- web/src/views/AreasView.vue (modified - button label)
- web/src/views/ItemGroupsView.vue (modified - title/subtitle removed, button label)
- web/src/views/ItemsView.vue (modified - title/subtitle removed, button label)
- web/src/components/PageHeader.vue (modified - added v-if on h1 for accessibility)
- web/src/views/ItemGroupsView.test.ts (modified - updated title assertion)
- web/src/views/AreasView.test.ts (modified - added select label assertion)
- web/src/views/AreasView.vue (modified - aria-label updated to Select)
- web/src/views/ItemGroupsView.vue (modified - aria-label updated to Select)
- web/src/views/ItemGroupsView.test.ts (modified - added select label assertion)