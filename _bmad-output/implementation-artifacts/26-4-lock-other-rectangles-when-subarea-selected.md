# Story 26.4: Lock Other Rectangles When Subarea Is Selected

Status: done

## Story

As an admin,
I want only the selected subarea to be editable on the canvas,
so that I cannot accidentally move or delete other subareas.

## Acceptance Criteria

1. **Given** I have selected "Open Space" for editing on the Areas tab
   **When** I try to click, move, or delete another subarea's rectangle (e.g., "Cube 1")
   **Then** the other rectangle does not respond to interaction

2. **Given** I have a subarea selected
   **When** I look at the other subarea rectangles on the canvas
   **Then** they appear visually distinct (e.g., dimmed or dashed) to indicate they are locked

## Tasks / Subtasks

- [x] Task 1: Filter activePositions to selected subarea only (AC: #1)
  - [x] 1.1 In `activePositions` computed, when on Areas tab with a `selectedSubAreaId`,
        return only the position matching that subarea ID
- [x] Task 2: Move all other positions to contextPositions (AC: #2)
  - [x] 2.1 In `contextPositions` computed, when on Areas tab with a `selectedSubAreaId`,
        return all positions except the selected subarea (they render as dashed, non-interactive)

### Review Findings

- [x] [Review][Patch] Auto-selected unpositioned first subarea can leave the Areas tab non-interactive on first load [web/src/views/FloorPlanEditorView.vue:1010]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- `activePositions`: added early return when `activeTab === "areas" && selectedSubAreaId`
  filtering to only `pos.itemId === selectedSubAreaId`
- `contextPositions`: added branch when `selectedSubAreaId` is set on Areas tab, returning
  all positions except the selected one — these render as dashed outlines with `pointer-events: none`

### File List

- `web/src/views/FloorPlanEditorView.vue` (modified)

### Change Log

- 2026-04-13: Implemented story 26.4
