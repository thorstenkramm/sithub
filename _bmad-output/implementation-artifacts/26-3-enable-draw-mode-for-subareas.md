# Story 26.3: Enable Draw Mode for Subareas on Areas Tab

Status: done

## Story

As an admin,
I want to draw a rectangle for an unpositioned subarea when I select it on the Areas tab,
so that I can position subareas on the floor plan.

## Acceptance Criteria

1. **Given** I am on the Areas tab and select an unpositioned subarea from the dropdown
   **When** the selection is made
   **Then** the editor enters draw mode (crosshair cursor) for that subarea

2. **Given** I am on the Areas tab and select a positioned subarea from the dropdown
   **When** the selection is made
   **Then** the editor selects that subarea's rectangle on the canvas

## Tasks / Subtasks

- [x] Task 1: Wire subarea selection to selectSidebarItem on Areas tab (AC: #1, #2)
  - [x] 1.1 In `onToolbarSubAreaSelect()`, when on Areas tab, find the subarea's EditableItem
        (scope === "area") and call `selectSidebarItem(item)` to enter draw mode or select rect

### Review Findings

- [x] [Review][Patch] Subarea selection clears its own draw/select state [web/src/views/FloorPlanEditorView.vue:506]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added Areas-tab logic to `onToolbarSubAreaSelect()`: looks up the subarea in
  `allEditableItems` with `scope === "area"` and calls `selectSidebarItem(item)`
- Reuses existing draw/select logic — unpositioned subareas get crosshair, positioned get selected

### File List

- `web/src/views/FloorPlanEditorView.vue` (modified)

### Change Log

- 2026-04-13: Implemented story 26.3
