# Story 20.7: First-Level Floor Plan Drill-Down

Status: backlog

## Story

As a user,
I want to click on an area in the first-level floor plan to open its detail floor plan,
So that I can drill down from the building overview to individual items.

## Acceptance Criteria

1. **Given** I open the floor plan for an area that has sub-areas with their own floor plans
   **When** the overlay renders
   **Then** each sub-area is shown with its positioned rectangle, a fraction label
   (e.g. "2/8 free"), and a color indicating availability

2. **Given** all items within a sub-area are booked for the selected day
   **When** the sub-area renders
   **Then** it shows a red semi-transparent overlay

3. **Given** a sub-area has some availability
   **When** the sub-area renders
   **Then** it shows a color gradient: green when mostly free, orange when few left,
   red when full

4. **Given** I click on a sub-area rectangle
   **When** the click is processed
   **Then** the detail floor plan for that sub-area opens with item-level free/busy state
   **And** a smooth zoom/fade transition animates the drill-down

5. **Given** I have drilled into a sub-area floor plan
   **When** I want to go back
   **Then** a breadcrumb trail at the top of the overlay shows the navigation path
   (e.g. "Büro 2.EG > Cube 1") and I can click the area name to return

## Tasks / Subtasks

- [ ] Detect area-level vs item-group-level floor plan (AC: 1)
  - [ ] Determine whether the current floor plan represents an area with sub-areas
    or an item group with individual items
  - [ ] Load sub-area positions and their associated floor plans
- [ ] Render sub-area rectangles with aggregated free/busy (AC: 1, 2, 3)
  - [ ] Show positioned sub-area rectangles on the area floor plan
  - [ ] Aggregate availability across all items within each sub-area
  - [ ] Show fraction label inside rectangle (e.g. "2/8 free")
  - [ ] Color gradient: green (mostly free) → orange (few left) → red (full)
  - [ ] Full red semi-transparent overlay when all items are booked
- [ ] Implement drill-down navigation with transition (AC: 4)
  - [ ] Handle click on sub-area rectangle
  - [ ] Animate drill-down with a fade-and-scale transition (~200ms)
  - [ ] Open the detail floor plan for the clicked sub-area
  - [ ] Show item-level free/busy state on the detail floor plan (reuse
    `InteractiveFloorPlan.vue` from Story 20.6)
- [ ] Implement breadcrumb back navigation (AC: 5)
  - [ ] Show breadcrumb trail at the top of the overlay inside the dialog
  - [ ] Format: "Area Name > Sub-Area Name"
  - [ ] Area name is clickable to navigate back to the area-level floor plan
  - [ ] Do not rely on browser back button — the floor plan is an overlay, not a page
- [ ] Add unit tests for drill-down logic and aggregated availability
- [ ] Verify E2E tests still pass

## Dev Notes

### UX Recommendations (Sally)

#### Breadcrumb back navigation

When drilling from the area floor plan into a sub-area, users need a clear way back. A
breadcrumb trail at the top of the overlay (e.g. `Büro 2.EG > Cube 1`) with the area
name clickable provides spatial context and navigation. Do not rely on the browser back
button — the floor plan is a dialog overlay, not a route change.

#### Partial availability visualization

A binary "all booked = red, otherwise green" misses the middle ground. A sub-area with
8 desks and only 2 free should feel different from one with 8 free. Show:

- A fraction label: "2/8 free"
- Color gradient: green when >50% free, orange when 1-50% free, red when full

This gives the user instant visual triage without needing to drill in.

#### Drill-down transition

A smooth fade-and-scale animation (~200ms) when drilling from area to sub-area helps the
user maintain spatial context. An abrupt content swap feels disorienting — even a simple
opacity + scale transition makes it feel intentional.

### Dependencies

- Depends on Story 20.6 (Interactive Floor Plan Overlay with Free/Busy)

### References

- Epic 20 Story 20.7: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR80: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

## Change Log

- 2026-03-22: UX review — added AC 3 (partial availability gradient), AC 4 updated
  (drill-down transition), AC 5 (breadcrumb navigation). Added tasks for fraction labels,
  color gradient, breadcrumb trail, and transition animation.
