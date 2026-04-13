# Story 26.1: Subarea Selection Respects Active Tab

Status: done

## Story

As an admin,
I want selecting a subarea from the dropdown to stay on the current tab,
so that I can position area rectangles without being forced into Items mode.

## Acceptance Criteria

1. **Given** I am on the Areas tab with an area-level floor plan loaded
   **When** I select "Open Space" from the subarea dropdown
   **Then** the toggle stays on "Areas" and does not switch to "Items"

2. **Given** I am on the Items tab
   **When** I select a subarea from the dropdown
   **Then** the toggle stays on "Items" (existing behavior preserved)

## Tasks / Subtasks

- [x] Task 1: Remove forced tab switch in onToolbarSubAreaSelect (AC: #1, #2)
  - [x] 1.1 Remove `activeTab.value = "items"` from `onToolbarSubAreaSelect()`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Removed the `if (subAreaId && isAreaLevel.value) { activeTab.value = "items"; }` block
  from `onToolbarSubAreaSelect()` — subarea selection now preserves whatever tab is active

### File List

- `web/src/views/FloorPlanEditorView.vue` (modified)

### Change Log

- 2026-04-13: Implemented story 26.1
