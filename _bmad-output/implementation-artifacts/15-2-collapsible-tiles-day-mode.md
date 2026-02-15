# Story 15.2: Collapsible Tiles in Day Booking Mode

Status: done

## Story

As a user,
I want booked item tiles in day mode to be collapsed by default,
So that I can focus on available items and expand booked ones only when needed.

## Acceptance Criteria

1. **Given** I am in day booking mode
   **When** an item is booked
   **Then** the item tile hides equipment and warning details by default
   **And** a chevron-left icon appears in the tile header

2. **Given** I click the chevron on a folded booked item tile
   **When** the tile unfolds
   **Then** equipment chips and warning alerts become visible
   **And** the chevron rotates to chevron-down

3. **Given** I am in day booking mode
   **When** an item is available
   **Then** the item tile shows all details (equipment, warnings) without a chevron
   **And** the tile is not collapsible

## Tasks / Subtasks

- [x] Add fold/unfold state for day mode (AC: 1, 2)
  - [x] Added `expandedDayTiles` reactive Set to track which item IDs are expanded
  - [x] Added `toggleDayTileExpansion(itemId: string)` function
  - [x] Reset expanded state in `loadItems()` when items reload
- [x] Add chevron icon to booked item tiles (AC: 1, 2)
  - [x] Added chevron in `#append` slot, only for occupied items
  - [x] Use `mdi-chevron-left` when folded, `mdi-chevron-down` when expanded
  - [x] Added `data-cy="day-tile-chevron"` for testing
- [x] Hide equipment and warning on folded booked tiles (AC: 1)
  - [x] Equipment and warning v-if conditions include availability or expansion check
  - [x] Booker name and booking note remain visible when folded
- [x] Keep available tiles always expanded (AC: 3)
  - [x] Available tiles always show equipment and warning (no chevron)
- [x] Add unit tests
  - [x] Test that booked tiles start folded (no equipment visible)
  - [x] Test that expanding a booked tile shows equipment
  - [x] Test available tiles always show equipment
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

All changes in `web/src/views/ItemsView.vue`. No backend changes required.

### Day Mode Tile Structure (Current)

Each day mode tile (lines 202-323) contains:
- `v-card-item`: avatar, item name, status chip
- `v-card-text`: equipment chips, warning alert, booker name, booking note
- `v-card-actions`: book button (available) / cancel button (admin) / "not available" text

The equipment and warning sections are always shown regardless of booking status. This story
changes that: for booked items, equipment and warning are hidden by default behind a chevron.

### Relationship to Story 14.3

Story 14.3 removes the "Not available for \<date\>" text and increases booker/note font sizes.
Story 15.2 adds collapsible behavior to booked tiles. These stories are independent — 14.3
changes what's always visible, 15.2 changes what's hidden/shown on demand. They can be
implemented in either order.

### References

- Epic 15 Story 15.2: `_bmad-output/planning-artifacts/epics.md` (Epic 15 Stories section)
- FR51: `_bmad-output/planning-artifacts/prd.md`
- ItemsView day section: `web/src/views/ItemsView.vue` lines 200-323

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `expandedDayTiles` reactive Set and `toggleDayTileExpansion()` function
- Added chevron icon in day tile #append slot (occupied items only)
- Equipment and warning v-if conditions include expansion/availability check
- Booker name and booking note remain visible when folded
- Available tiles always show all details, no chevron
- Reset expanded state on `loadItems()` reload
- Added 3 unit tests: folded hides equipment, expanded shows equipment, available always shows
- Updated existing test to use available item for equipment/warning rendering
- All 138 unit tests pass, all 51 E2E tests pass
- Code review fix: reset expanded day tiles on item reload to prevent stale expansion
- Code review fix: add aria labels/expanded state to day tile chevron for accessibility

### Change Log

- 2026-02-14: Implemented Story 15.2 - collapsible tiles in day booking mode
- 2026-02-14: Code review fixes for day tile reset and accessibility

### File List

- web/src/views/ItemsView.vue (modified - day tile expansion, chevron, conditional equipment)
- web/src/views/ItemsView.test.ts (modified - updated and added tests for collapsible behavior)
- web/src/views/ItemsView.vue (modified - reset expanded day tiles, aria labels)