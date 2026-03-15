# Story 18.5: Area Floor Plan Display

Status: done

## Story

As a user,
I want to see a "Floor plan" button when viewing an area that has a floor plan,
So that I can see where items are located.

## Acceptance Criteria

1. **Given** I am viewing an area that has a floor plan configured
   **When** the page loads
   **Then** a "Floor plan" button with an appropriate icon appears next to the calendar week selector

2. **Given** I click the "Floor plan" button
   **When** the overlay opens
   **Then** the floor plan image is displayed with the area name as heading
   **And** I can close the overlay

3. **Given** I am viewing an area without a floor plan
   **When** the page loads
   **Then** no "Floor plan" button appears

## Tasks / Subtasks

- [x] Add `mdiMap` icon import and `map` alias to `web/src/plugins/vuetify.ts`
- [x] Add reactive state to `ItemGroupsView.vue`
  - [x] `areaFloorPlan` ref for storing the floor plan path
  - [x] `showFloorPlanDialog` ref for dialog visibility
  - [x] `floorPlanUrl` computed property to construct API URL
- [x] Capture `area.attributes.floor_plan` in `onMounted`
- [x] Add "Floor plan" button next to the calendar week selector
  - [x] Conditionally shown with `v-if="areaFloorPlan"`
  - [x] Uses `$map` icon and `variant="outlined"` styling
  - [x] `data-cy="area-floor-plan-btn"`
- [x] Add floor plan dialog with `v-img`
  - [x] Area name as title
  - [x] `max-width="900"`, `max-height="600"`, `contain` mode
  - [x] Close button
  - [x] `data-cy="floor-plan-dialog"`, `data-cy="floor-plan-image"`
- [x] Run type-check, build, and unit tests

## Dev Notes

### URL Construction

The `floor_plan` attribute is treated as a filename inside the configured `areas.floor_plans`
directory. The frontend constructs the authenticated API URL directly from that filename:

```typescript
const floorPlanUrl = computed(() => {
  if (!areaFloorPlan.value) return '';
  return `/api/v1/floor-plans/${encodeURIComponent(areaFloorPlan.value)}`;
});
```

### References

- Epic 18 Story 18.5: `_bmad-output/planning-artifacts/epics.md`
- FR64: `_bmad-output/planning-artifacts/prd.md`
- `web/src/views/ItemGroupsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `mdiMap` icon and `map` alias to vuetify plugin
- Added floor plan button and dialog to ItemGroupsView
- Button conditionally shown based on `area.attributes.floor_plan`
- URL is constructed directly from the validated filename-only `floor_plan` attribute
- Added frontend coverage for visible, hidden, and dialog-open states
- All 197 frontend tests, type-check, and build pass

### File List

- `web/src/plugins/vuetify.ts` — Added `mdiMap` import and `map` alias
- `web/src/views/ItemGroupsView.vue` — Floor plan button, dialog, reactive state
- `web/src/views/ItemGroupsView.test.ts` — Verifies visible, hidden, and dialog-open states

## Change Log

- 2026-03-13: Story implemented and verified.
- 2026-03-13: Code review fixes moved the area floor plan control onto the area booking screen and added view tests.
