# Story 25.4: Enlarged Subarea Images on Drill-Down

Status: ready-for-dev

## Story

As a user,
I want subarea floor plan images to be displayed enlarged when I drill into them,
so that I can clearly see the layout and available items without zooming.

## Acceptance Criteria

1. **Given** I am viewing the floor plan booking view for an area
   **When** I click on a subarea rectangle to drill into it
   **Then** the subarea floor plan image renders at an enlarged size that fills the available
   viewport width

2. **Given** I have drilled into a subarea
   **When** the subarea floor plan is displayed at default zoom level
   **Then** no horizontal scrollbars appear; the image fits within the viewport width

3. **Given** I manually zoom in beyond the default level
   **When** the image exceeds the viewport width
   **Then** scrollbars appear as expected to allow navigation

## Tasks / Subtasks

- [ ] Task 1: Ensure drill-down image fills viewport width (AC: #1, #2)
  - [ ] 1.1 In `InteractiveFloorPlan.vue`, locate the drill-down handler: `handleAreaClick()` (line ~1116) sets `drilledInto` state, and the watcher (line ~968) resets zoom to 1 and loads the new floor plan
  - [ ] 1.2 Locate the image sizing logic (lines ~527-535): on image load, `img.style.width` is set to `shell.clientWidth` to fit width
  - [ ] 1.3 Verify that after drill-down, the image `onload` handler fires and sets the image width to fill the `fp-scroll-shell` container width
  - [ ] 1.4 If the image is not filling viewport width after drill-in, adjust the sizing logic to ensure `img.style.width = shell.clientWidth + 'px'` is called after the drilled floor plan image loads
  - [ ] 1.5 Ensure zoom resets to 1.0 on drill-in (already done at line ~970) so the fit-to-width calculation is correct
- [ ] Task 2: Prevent horizontal scrollbars at default zoom (AC: #2)
  - [ ] 2.1 Inspect `.fp-scroll-shell` CSS (lines ~1536-1545): `overflow: auto`, `max-height: calc(100vh - 260px)`
  - [ ] 2.2 Ensure the image width calculation accounts for any padding/margins inside the scroll shell so the image does not exceed the container width
  - [ ] 2.3 If the image has border, padding, or the zoom layer has margins, subtract those from `shell.clientWidth` when setting image width
  - [ ] 2.4 Test that at zoom = 1.0 after drill-in, no horizontal scrollbar appears
- [ ] Task 3: Verify zoom scrollbar behavior (AC: #3)
  - [ ] 3.1 Confirm that when user zooms in (via buttons or pinch), the image scales beyond viewport width and scrollbars appear naturally via `overflow: auto` on `.fp-scroll-shell`
  - [ ] 3.2 No code changes expected for this — just verify the existing behavior
- [ ] Task 4: Validate (AC: #1-#3)
  - [ ] 4.1 Run `npm run lint` and fix findings
  - [ ] 4.2 Run `npm run type-check` and fix findings
  - [ ] 4.3 Run `npm run build` and verify no build errors
  - [ ] 4.4 Run `npx vitest run` and verify no regressions
  - [ ] 4.5 Run `npm run test:e2e -- --browser electron` and verify no regressions

## Dev Notes

### Architecture & Patterns

- **Single file change**: `web/src/components/InteractiveFloorPlan.vue`
- **No backend changes**: Pure frontend sizing/CSS adjustment
- **Parent views**: Used by `ItemsView.vue` (single item group) and `ItemGroupsView.vue` (area-level with drill-down)

### Key Code Locations

| Element | Location | data-cy |
|---------|----------|---------|
| `InteractiveFloorPlan.vue` | `web/src/components/InteractiveFloorPlan.vue` | — |
| `handleAreaClick()` | Line ~1116 | — |
| Drill-down watcher | Line ~968 | — |
| Image sizing on load | Lines ~527-535 | — |
| `.fp-scroll-shell` CSS | Lines ~1536-1545 | `fp-scroll-shell` |
| `zoomScale` ref | Line ~514 | — |
| Zoom range | Lines ~1428 (clampZoom: 0.75-2.5) | — |
| `drilledInto` ref | Reactive ref tracking drill state | — |

### Drill-Down Flow

1. User clicks a subarea rectangle → `handleAreaClick(itemId)` fires
2. Sets `drilledInto.value` with itemGroupId, name, and floorPlan object
3. Watcher detects change, resets `zoomScale` to 1, loads new floor plan image
4. Image `onload` sets `img.style.width = shell.clientWidth` for fit-to-width

### Implementation Strategy

The fit-to-width logic likely already works for the initial load. The issue may be:
- The image `onload` not firing after drill-down (timing issue with reactive state)
- Padding/margin not accounted for in width calculation
- The scroll shell having different dimensions after drill-down

Test by drilling into a subarea and checking if the image fills width without scrollbars. Fix whichever of the above applies.

### Anti-Patterns to Avoid

- Do NOT change the zoom range (0.75-2.5) — keep existing behavior
- Do NOT modify `FloorPlanEditorView.vue` — this story is booking view only
- Do NOT use fixed pixel widths — must be responsive to container width
- Do NOT disable scrolling entirely — only prevent it at default zoom level

### References

- [Source: web/src/components/InteractiveFloorPlan.vue — drill-down and image sizing]
- [Source: web/src/views/ItemGroupsView.vue — parent view with area-level floor plan]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List

### Change Log
