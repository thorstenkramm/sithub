# Story 26.2: Hide Items Dropdown on Areas Tab

Status: done

## Story

As an admin,
I want the items dropdown to be hidden when I am on the Areas tab,
so that I am not confused by irrelevant controls while positioning subareas.

## Acceptance Criteria

1. **Given** I am on the Areas tab
   **When** I look at the toolbar
   **Then** the "Objekte" (Items) dropdown is not visible

2. **Given** I switch to the Items tab
   **When** I look at the toolbar
   **Then** the "Objekte" (Items) dropdown appears

## Tasks / Subtasks

- [x] Task 1: Add visibility condition to items dropdown (AC: #1, #2)
  - [x] 1.1 Change `v-if="selectedFloorPlan"` to `v-if="selectedFloorPlan && !(isAreaLevel && activeTab === 'areas')"` on the items v-select

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `!(isAreaLevel && activeTab === 'areas')` guard to items dropdown v-if

### File List

- `web/src/views/FloorPlanEditorView.vue` (modified)

### Change Log

- 2026-04-13: Implemented story 26.2
