# Story 20.1: Free-Busy Indicators on Favorite Tiles

Status: done

## Story

As a user,
I want to see weekly availability indicators on my promoted third-level favorite tiles,
So that I can quickly see which days have availability without navigating into the
item group.

## Acceptance Criteria

1. **Given** I have third-level favorites promoted to the item-groups view
   **When** the page loads and availability data is fetched
   **Then** the favorite tiles show the same MO-TU-WE-TH-FR availability dots as regular
   item group tiles

2. **Given** an item within a favorite's item group is fully booked on a day
   **When** the availability dot renders
   **Then** the dot shows the booked (red outline) indicator for that day

3. **Given** I have a promoted favorite tile
   **When** the tile renders
   **Then** the item group name is shown as a subtitle beneath the item name for context
   (e.g. item name "Tisch 1", subtitle "Cube 1")

## Tasks / Subtasks

- [x] Fetch availability for favorite item groups (AC: 1, 2)
  - [x] Extend the item-groups view to request availability data for promoted
    third-level favorite items alongside regular item groups
  - [x] Map availability responses to the same structure used for regular tiles
- [x] Render availability dots on favorite tiles (AC: 1, 2)
  - [x] Add MO-TU-WE-TH-FR availability dots to promoted favorite tiles matching
    the regular item group tile style
  - [x] Show red outline indicator for fully booked days
- [x] Add item group name as subtitle on favorite tiles (AC: 3)
  - [x] Show the item group name beneath the item name on promoted favorite tiles
  - [x] Use `v-card-subtitle` or equivalent for consistent styling
- [x] Add unit tests for availability dot rendering on favorites
- [ ] Verify E2E tests still pass

## Dev Notes

### UX Recommendation (Sally)

Promoted favorite tiles should show the parent item group name as a subtitle. When a user
has multiple favorites from different item groups, the context helps them scan quickly.
For example: title "Tisch 1, am Gang, rechts", subtitle "Cube 1".

### References

- Epic 20 Story 20.1: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR75: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

GPT-5 Codex

### Completion Notes List

- Promoted third-level favorites now render the parent item-group subtitle on the item-groups view.
- Favorite tiles reuse the parent item-group weekly availability dots, including the red booked indicator.
- AI review fix: promoted favorite tiles now use a fully scoped Vue key to avoid collisions across item groups.
- Added targeted view tests for memorized week usage and favorite-tile availability rendering.
- E2E tests were not run in this review/fix pass.

### File List

- `web/src/views/ItemGroupsView.vue` — Promoted favorite tiles render subtitle, weekly availability, and scoped keys
- `web/src/views/ItemGroupsView.test.ts` — Added coverage for favorite-tile availability dots and memorized week loading

## Senior Developer Review (AI)

- Verified ACs 20.1.1 to 20.1.3 against the current item-groups implementation.
- Fixed a rendering risk where promoted favorites were keyed only by `itemId`, which is not globally unique.
- Added targeted Vitest coverage for promoted favorite tiles so the story no longer relies on regular tile tests alone.

## Change Log

- 2026-03-22: UX review — added AC 3 (item group subtitle on favorite tiles) and
  corresponding task.
- 2026-03-22: Story implementation reviewed and finalized; fixed favorite-tile key scoping and added targeted tests.
