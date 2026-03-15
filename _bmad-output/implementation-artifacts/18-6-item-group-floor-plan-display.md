# Story 18.6: Item Group Floor Plan Display

Status: done

## Story

As a user,
I want to see a "Floor plan" button when viewing an item group that has a floor plan,
So that I can see the layout of individual items within the group.

## Acceptance Criteria

1. **Given** I am viewing an item group that has a floor plan configured
   **When** the page loads
   **Then** a "Floor plan" button with an appropriate icon appears beneath the day/week
   selector

2. **Given** I click the "Floor plan" button
   **When** the overlay opens
   **Then** the floor plan image is displayed with the item group name as heading
   **And** I can close the overlay

3. **Given** I am viewing an item group without a floor plan
   **When** the page loads
   **Then** no "Floor plan" button appears

## Tasks / Subtasks

- [x] Add reactive state to `ItemsView.vue`
  - [x] `itemGroupFloorPlan` ref for storing the floor plan path
  - [x] `showItemGroupFloorPlanDialog` ref for dialog visibility
  - [x] `itemGroupFloorPlanUrl` computed property to construct API URL
- [x] Capture `ig.attributes.floor_plan` in `onMounted` breadcrumb fetch
- [x] Add "Floor plan" button beneath the equipment filter area
  - [x] Conditionally shown with `v-if="itemGroupFloorPlan"`
  - [x] Uses `$map` icon and `variant="outlined"` styling
  - [x] `data-cy="item-group-floor-plan-btn"`
- [x] Add floor plan dialog with `v-img`
  - [x] Item group name as title
  - [x] `max-width="900"`, `max-height="600"`, `contain` mode
  - [x] Close button
  - [x] `data-cy="item-group-floor-plan-dialog"`, `data-cy="item-group-floor-plan-image"`
- [x] Run type-check, build, and unit tests

## Dev Notes

### Same Pattern as Story 18.5

Uses the same URL construction pattern as the area floor plan in `ItemGroupsView.vue`. The
`$map` icon alias was already added in Story 18.5, and the frontend now treats `floor_plan` as a
filename inside `areas.floor_plans`.

### References

- Epic 18 Story 18.6: `_bmad-output/planning-artifacts/epics.md`
- FR65: `_bmad-output/planning-artifacts/prd.md`
- `web/src/views/ItemsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added floor plan button and dialog to ItemsView following the same filename-based pattern as the area floor plan
- URL is constructed directly from the validated filename-only `floor_plan` attribute
- Button placed above equipment filter, conditionally shown
- Captured `ig.attributes.floor_plan` during breadcrumb area/item-group fetch
- Added frontend coverage for the floor plan button and dialog
- All 197 frontend tests, type-check, and build pass

### File List

- `web/src/views/ItemsView.vue` — Floor plan button, dialog, reactive state
- `web/src/views/ItemsView.test.ts` — Verifies floor plan dialog wiring

## Change Log

- 2026-03-13: Story implemented and verified.
- 2026-03-13: Code review fixes aligned floor plan URL handling with filename-only validation and added view tests.
