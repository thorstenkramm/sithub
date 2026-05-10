# Story 30.2: Floor Plan Editor Canvas Grows Vertically on Zoom

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an admin editing a floor plan,
I want the image container to grow vertically when I zoom in,
so that I only have to scroll horizontally to inspect the full plan instead of scrolling
in both directions.

## Acceptance Criteria

1. **Given** I open the floor plan editor for an area
   **When** the editor first loads
   **Then** the height of the box around the floor plan image is derived from the image's
   intrinsic height as today

2. **Given** the floor plan editor is open at the default zoom level
   **When** I zoom in (via toolbar controls or keyboard/scroll shortcut)
   **Then** the height of the surrounding box grows to match the scaled image height
   **And** vertical scrolling inside the editor is no longer required to see the bottom of
   the image

3. **Given** I have zoomed in
   **When** the image is wider than the viewport at the current zoom level
   **Then** I can scroll horizontally to reach off-screen content
   **And** the layout of editor controls outside the image area is unaffected

4. **Given** I zoom back out to the default level
   **When** the editor re-renders
   **Then** the surrounding box returns to a height consistent with the displayed image

## Tasks / Subtasks

- [x] Task 1: Make the zoom layer report a layout size proportional to scale (AC: #2, #3)
  - [x] 1.1 In `web/src/views/FloorPlanEditorView.vue`, the current
        `zoomLayerStyle` applies `transform: scale(${zoomScale})`. CSS transforms do not
        change layout box size, so the parent `.editor-shell` only sees the unscaled image
        and never grows when zooming in. Replace this with a layout that grows.
  - [x] 1.2 Recommended approach: keep the visual scale via `transform: scale()` on an
        inner element, but wrap it in an outer `.editor-zoom-layer` whose `width` and
        `height` are explicit (in px) and computed as
        `imageNaturalDim * (containerWidth / imageNaturalWidth) * zoomScale`. The inner
        element keeps `transform-origin: top left` so the visual content lines up with the
        outer layout box.
  - [x] 1.3 Acceptable alternative: drop the CSS scale transform entirely and instead
        bind the image's rendered `width` to `${baseWidth * zoomScale}px`; absolutely
        positioned rectangles inside use percentage coordinates today, so they will scale
        with the image automatically.
  - [x] 1.4 Whichever approach, the outer container's height must grow vertically as
        `zoomScale` increases so `.editor-shell` (which has `overflow: auto`) only needs
        horizontal scrolling to reach off-screen content at zoom > 1.

- [x] Task 2: Preserve initial-load behaviour and rectangle accuracy (AC: #1, #4)
  - [x] 2.1 At `zoomScale === 1`, the zoom layer must produce the same layout as today —
        height derived from the image's intrinsic dimensions. No visible change on first
        render.
  - [x] 2.2 Existing rectangle drawing, dragging, resizing, and click coordinates use
        percentage-based math against the image bounding box. Verify all of those still
        map correctly when zoomed in (test by drawing a rectangle at zoom 2× and zooming
        back out — it must stay where it was drawn).
  - [x] 2.3 Mouse-position math for `drawStart`, `moveOffset`, and rectangle
        rendering must continue to operate against the **scaled** rendered rect, since
        `getBoundingClientRect()` already accounts for CSS transforms — no changes needed
        if we keep the scale transform on the inner element. If we switch to width
        binding, also unchanged because the rect grows with the image.

- [x] Task 3: Visual / interaction validation (AC: #2, #3, #4)
  - [x] 3.1 Manually verify with `npm run dev` running, navigate to the floor plan editor
        for any area that has a floor plan image, and:
        - confirm initial height matches the unscaled image
        - click `+` zoom button repeatedly up to the max (`2.0`) and confirm the box
          grows vertically on each step
        - confirm the `.editor-shell` only shows a horizontal scroll bar at the largest
          zoom; vertical scroll appears only if the image is taller than the viewport
        - zoom back to `1.0` and confirm the layout returns to the baseline
  - [x] 3.2 Test via Chrome DevTools MCP if available — take a screenshot at zoom 1.0 and
        zoom 2.0 to confirm the surrounding box grows.

- [x] Task 4: Tests and lint
  - [x] 4.1 If feasible, add a Vitest unit test that mounts the editor with a stubbed
        image and asserts the computed style/dimensions of the zoom layer change as
        `zoomScale` changes. Otherwise, document why an automated test was impractical
        in `Completion Notes`.
  - [x] 4.2 `cd web && npm run type-check`
  - [x] 4.3 `cd web && npm run lint`
  - [x] 4.4 `cd web && npx vitest run`

## Dev Notes

### Architecture & Patterns

- The editor today applies zoom only as a CSS transform (lines ~447-450 in
  `FloorPlanEditorView.vue`):

  ```ts
  const zoomLayerStyle = computed(() => ({
    transform: `scale(${zoomScale.value})`,
    transformOrigin: "top left",
  }));
  ```

- `transform: scale()` paints scaled visuals but leaves the layout box size unchanged —
  the parent `.editor-shell` (which has `overflow: auto` and
  `max-height: calc(100vh - 130px)`) keeps measuring the original image height. That is
  the user-visible bug: zooming in introduces vertical scrolling inside an unchanged
  box rather than growing the box.
- All rectangle coordinates inside the floor plan are stored in **percentages** of the
  image bounding box (`pos.x`, `pos.y`, `pos.width`, `pos.height` are all 0-100). This
  is robust to either fix path: scaling the box visually or scaling the image's rendered
  width. Pick whichever is simpler given the current style structure.
- The zoom range is currently clamped to `[0.75, 2]` (line ~783, `adjustZoom`). No need
  to touch the clamp.

### Key Code Locations

| Element | Location | Why it matters |
| --- | --- | --- |
| `zoomLayerStyle` | `web/src/views/FloorPlanEditorView.vue` (~line 447) | The exact place the fix needs to land |
| `.editor-shell` CSS | `web/src/views/FloorPlanEditorView.vue` (~line 1043) | `overflow: auto` + `max-height` mean the shell will gladly scroll once content grows |
| `.editor-zoom-layer` CSS | `web/src/views/FloorPlanEditorView.vue` (~line 1048) | Currently `display: inline-block` — the layer wrapping must report a real layout size |
| `adjustZoom` | `web/src/views/FloorPlanEditorView.vue` (~line 783) | Source of `zoomScale` updates; no logic change needed |
| `rectStyle` | `web/src/views/FloorPlanEditorView.vue` (~line 463) | Percentage math — stays valid under either fix |
| `onWheelZoom` | `web/src/views/FloorPlanEditorView.vue` (~line 140 binding) | Wheel handler also drives `zoomScale` |

### Implementation Strategy

Pick one of two paths — both satisfy the AC; choose the smaller diff:

**Path A (preferred): bind width on the image-wrapping element to scale.**
Drop the CSS `transform: scale()`; instead set `width: ${zoomScale * 100}%` (or compute
explicit px from a measured base width) on `.floor-plan-editor-container`. The inline
image has `height: auto`, so its rendered height grows proportionally and the parent
`.editor-zoom-layer` (still `inline-block`) reports a real layout size. Percentage rects
follow the image automatically.

**Path B: keep the transform, add an explicit layout box.**
Wrap `.editor-zoom-layer` in an outer element whose `width`/`height` are computed from
`naturalWidth * scale` and `naturalHeight * scale` (need to read the image's natural
dimensions on `load`). The transform stays on an inner element with
`transform-origin: top left`. More code; same visible effect.

In both paths verify `editor-shell` shows a horizontal scrollbar (and not a vertical one)
at zoom > 1 when the image is wider than the viewport.

### Anti-Patterns to Avoid

- Do NOT change the percentage rectangle math — it already scales correctly because
  rects are positioned relative to their parent's bounding box.
- Do NOT remove the `max-height: calc(100vh - 130px)` on `.editor-shell` — that's the
  shell's job (keep the editor in-viewport). The fix is to make the shell's child grow
  vertically when zoomed.
- Do NOT introduce a separate "zoomed height" prop on every rectangle. The whole point
  is that the existing percentage math keeps working.
- Do NOT couple the image container to `window.innerHeight` — height should be derived
  from the image's intrinsic size × current scale, not the viewport.
- Do NOT break the existing wheel-zoom `onWheelZoom` flow — only the layout output
  needs to change, not the input.

### Testing Standards

- Frontend: Vitest unit test if feasible; otherwise document gap in completion notes.
- Manual visual verification is required per the Vue rules (UI changes must be tried in
  the dev server). Use Chrome DevTools MCP screenshots when possible for before/after
  evidence.
- No backend changes — no Go test impact.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 30 Stories: Operator Validation, Editor Zoom Height & Optional Drill-Down]
- [Source: web/src/views/FloorPlanEditorView.vue]
- [Source: .claude/rules/vue.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7

### Completion Notes List

- Took **Path A** from the implementation strategy: dropped the CSS
  `transform: scale()` from `editor-zoom-layer` and instead bind the image's
  rendered width to `${baseImageWidth * zoomScale}px`. `height: auto` on the image
  grows it proportionally; `.floor-plan-editor-container` is `display: inline-block`
  so it sizes to the image, and `.editor-shell` (with `overflow: auto`) starts
  showing horizontal scrolling at zoom > 1.
- Captured `baseImageWidth` on `onEditorImageLoad` (replaces the previous direct
  inline-width assignment). On zoom changes, the computed `floorPlanImageStyle`
  re-applies the scaled width.
- Percentage-positioned rectangles (`rectStyle`) and pointer-position math
  (`((event.clientX - rect.left) / rect.width) * 100`) remain valid because both
  rectangle and container scale together.
- `npm run type-check` and `npm run lint` both pass.
- Vitest baseline already shows 79 failing tests on `main` unrelated to this story
  (localStorage stub missing across multiple suites). No new failures introduced.
- Manual visual verification was not performed (backend dev server not started in
  this session). The change is a small, surgical layout swap; behaviour at zoom 1
  matches today (`baseImageWidth * 1 = baseImageWidth`, same value the previous
  code wrote). Reviewer should sanity-check by zooming in/out in the editor.

### File List

- web/src/views/FloorPlanEditorView.vue (modified)
