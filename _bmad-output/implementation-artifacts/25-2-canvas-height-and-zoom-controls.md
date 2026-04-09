# Story 25.2: Canvas Height & Zoom Controls

Status: ready-for-dev

## Story

As an admin,
I want a taller canvas and compact zoom controls,
so that I can see and edit the floor plan image with less scrolling and a cleaner toolbar.

## Acceptance Criteria

1. **Given** I open the floor plan editor
   **When** the editor loads
   **Then** the canvas area uses approximately double the vertical space compared to the
   previous layout

2. **Given** I look at the zoom controls in the editor toolbar
   **When** I inspect their layout
   **Then** the zoom percentage label appears between the minus and plus buttons, not next
   to them

3. **Given** I click the plus or minus zoom buttons
   **When** the zoom level changes
   **Then** the percentage label between the buttons updates to reflect the current zoom
   factor

## Tasks / Subtasks

- [ ] Task 1: Double the canvas height (AC: #1)
  - [ ] 1.1 In `FloorPlanEditorView.vue`, locate `.editor-shell` CSS (line ~1031): change `max-height: calc(100vh - 230px)` to approximately `calc(100vh - 130px)` to roughly double vertical space
  - [ ] 1.2 Verify the editor canvas is visibly taller and the floor plan image has more room
  - [ ] 1.3 Test on a smaller viewport to ensure the editor remains usable (no overflow issues)
- [ ] Task 2: Reposition zoom percentage label (AC: #2, #3)
  - [ ] 2.1 In the toolbar section (lines ~142-156), locate the zoom controls: currently ordered as minus button, slider, plus button, then percentage label
  - [ ] 2.2 Rearrange to: minus button, percentage label (e.g., "100%"), plus button — remove or reposition the slider
  - [ ] 2.3 Ensure the percentage label reactively displays `Math.round(zoomScale * 100) + '%'`
  - [ ] 2.4 Style the label to be visually centered between the buttons with appropriate spacing
  - [ ] 2.5 Verify clicking plus/minus updates the displayed percentage
- [ ] Task 3: Validate (AC: #1-#3)
  - [ ] 3.1 Run `npm run lint` and fix findings
  - [ ] 3.2 Run `npm run type-check` and fix findings
  - [ ] 3.3 Run `npm run build` and verify no build errors
  - [ ] 3.4 Run `npx vitest run` and verify no regressions
  - [ ] 3.5 Run `npm run test:e2e -- --browser electron` and verify no regressions

## Dev Notes

### Architecture & Patterns

- **Single file change**: `web/src/views/FloorPlanEditorView.vue`
- **No backend changes**: Pure frontend CSS/template refactor
- **No store involved**: Editor uses local Composition API state

### Key Code Locations

| Element | Location | data-cy |
|---------|----------|---------|
| `.editor-shell` max-height | Line ~1031 | — |
| Zoom slider | Line ~142 | `editor-zoom-slider` |
| Zoom minus button | Line ~143 | — |
| Zoom plus button | Line ~155 | — |
| Zoom percentage display | Line ~156 | — |
| `zoomScale` ref | Line ~367 | — |

### Implementation Notes

- The zoom controls currently use a `v-slider` between minus/plus buttons with the percentage shown after the plus button
- Replace the slider with a static percentage label between the buttons for a more compact layout
- The `zoomScale` ref holds the current zoom (e.g., 1.0 = 100%)
- Zoom step size is defined in the existing plus/minus click handlers — keep the same increment
- The `InteractiveFloorPlan.vue` component (booking view) already uses a compact minus/percentage/plus layout — match that pattern

### Anti-Patterns to Avoid

- Do NOT change the zoom step size or range — keep existing behavior
- Do NOT modify `InteractiveFloorPlan.vue` — this story is editor-only
- Do NOT add new reactive variables — reuse `zoomScale`

### References

- [Source: web/src/views/FloorPlanEditorView.vue — zoom controls and CSS]
- [Source: web/src/components/InteractiveFloorPlan.vue — compact zoom layout reference]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List

### Change Log
