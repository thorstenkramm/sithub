# Story 20.5: Floor Plan Editor (Admin)

Status: backlog

## Story

As an admin,
I want to draw rectangles on floor plan images to mark where items are located,
So that users can see and click items on the interactive floor plan.

## Acceptance Criteria

1. **Given** I am an admin and open the floor plan editor from settings
   **When** I select a floor plan
   **Then** the floor plan image is displayed with a list of unpositioned items on the left

2. **Given** I select an item from the list
   **When** I draw a rectangle on the floor plan image
   **Then** the rectangle is created with a label showing the item name
   **And** the cursor changes to crosshair during draw mode and reverts after

3. **Given** I have positioned items on the floor plan
   **When** I save
   **Then** all positions are persisted via the API
   **And** saved rectangles show a brief visual confirmation (green flash)

4. **Given** I want to reposition an item
   **When** I drag or resize its rectangle
   **Then** the position updates visually and can be saved

5. **Given** I want to remove a positioned item
   **When** I select the rectangle and press Delete or click the trash icon in the sidebar
   **Then** the rectangle is removed from the floor plan

6. **Given** I make a mistake (wrong position, accidental resize)
   **When** I press Ctrl+Z
   **Then** the last action is undone

7. **Given** a floor plan image is large and detailed
   **When** I need to position items precisely
   **Then** I can zoom in/out using Ctrl+scroll or a zoom slider

## Tasks / Subtasks

- [ ] Create floor plan editor view (AC: 1)
  - [ ] Create `FloorPlanEditorView.vue` in `web/src/views/`
  - [ ] Display floor plan image inside a `position: relative` container
  - [ ] Add floor plan selector (dropdown of available floor plans from areas config)
- [ ] Add route for floor plan editor (AC: 1)
  - [ ] Add admin-only route in Vue Router
  - [ ] Add navigation entry in settings menu
- [ ] Implement item list sidebar (AC: 1, 2, 5)
  - [ ] Show list of items for the selected floor plan's area/item group
  - [ ] Distinguish unpositioned items (available) from already-positioned ones (greyed out)
  - [ ] Clicking an unpositioned item enters draw mode for that item
  - [ ] Show trash icon next to selected/positioned items for deletion
- [ ] Implement draw mode with visual indicator (AC: 2)
  - [ ] On entering draw mode: change cursor to `crosshair`, highlight the sidebar item
  - [ ] On `pointerdown` + `pointermove` + `pointerup`, create rectangle
  - [ ] On completing draw or pressing Escape: exit draw mode, restore default cursor
- [ ] Implement rectangle drawing with DOM overlays (AC: 2, 4)
  - [ ] Absolutely positioned `<div>` rectangles inside the container
  - [ ] Store coordinates as percentages of the container size (resolution-independent)
  - [ ] Show item name as a label inside or above the rectangle
  - [ ] Drag to move: `pointerdown` on existing rectangle, update `top`/`left`
  - [ ] Resize: four small corner handles; on drag, update `width`/`height`
- [ ] Implement undo (AC: 6)
  - [ ] Store previous state before each mutation (draw, move, resize, delete)
  - [ ] Ctrl+Z reverts the last action (one level of undo is sufficient)
- [ ] Implement zoom (AC: 7)
  - [ ] Ctrl+scroll to zoom the floor plan container
  - [ ] Optional: zoom slider control
  - [ ] Ensure rectangle positions remain correct at all zoom levels (percentage-based)
- [ ] Implement save/load positions via API (AC: 3, 5)
  - [ ] Load existing positions from GET `/api/v1/floor-plan-positions`
  - [ ] Save new/updated positions via POST/PUT endpoints
  - [ ] Delete positions via DELETE endpoint
  - [ ] Show brief green flash on saved rectangles as visual confirmation
- [ ] Add to settings menu (AC: 1)
  - [ ] Add "Floor Plan Editor" entry visible to admin users only
- [ ] Add unit tests for editor components
- [ ] Verify E2E tests still pass

## Dev Notes

### Architecture Decision: HTML/CSS Overlays (no Konva.js)

Konva.js was evaluated and rejected. The requirements are simple rectangles on an image —
no rotation, freehand drawing, or complex layering. Using plain DOM keeps one rendering
paradigm across the app, allows Vuetify tooltips and components inside overlays, supports
Cypress `data-cy` selectors, and avoids ~140 KB of additional bundle weight.

### Technical Implementation Guide

#### Container structure

```vue
<div class="floor-plan-container" style="position: relative; display: inline-block;">
  <img :src="floorPlanUrl" @load="onImageLoad" draggable="false" />
  <div
    v-for="pos in positions"
    :key="pos.itemId"
    class="floor-plan-rect"
    :style="{
      position: 'absolute',
      left: pos.x + '%',
      top: pos.y + '%',
      width: pos.width + '%',
      height: pos.height + '%'
    }"
  >
    <span class="rect-label">{{ pos.label || pos.itemName }}</span>
    <div class="resize-handle top-left" />
    <div class="resize-handle top-right" />
    <div class="resize-handle bottom-left" />
    <div class="resize-handle bottom-right" />
  </div>
</div>
```

#### Coordinate system

All positions are stored as **percentages** (0-100) relative to the image's natural
dimensions. This ensures positions render correctly regardless of viewport size or dialog
width. Convert pointer events to percentages using
`(event.offsetX / containerWidth) * 100`.

#### Pointer event flow

1. **Draw mode**: Admin selects an unpositioned item, then `pointerdown` on the container
   records the start point. `pointermove` shows a preview rectangle. `pointerup` finalizes
   the rectangle and exits draw mode. Cursor is `crosshair` throughout.
2. **Move mode**: `pointerdown` on an existing rectangle (not a handle) starts drag.
   `pointermove` updates `left`/`top`. `pointerup` ends drag.
3. **Resize mode**: `pointerdown` on a corner handle. `pointermove` updates the opposite
   corner's coordinate. `pointerup` ends resize.

Use `pointer-events: none` on labels to prevent them from intercepting interactions.

#### No external dependencies needed

The editor needs ~80-120 lines of pointer event handling. Do not add Konva.js,
`@vueuse/gesture`, or other dragging libraries. Standard `pointerdown`/`pointermove`/
`pointerup` with `setPointerCapture` is sufficient and keeps the bundle lean.

### UX Recommendations (Sally)

#### Draw mode indicator

When the admin selects an item and enters draw mode, provide clear visual feedback:
cursor changes to `crosshair`, the selected sidebar item is highlighted with primary
color. After drawing or pressing Escape, the cursor and highlight revert.

#### Undo

Admins will make mistakes. Implement one level of undo (Ctrl+Z) by storing the previous
state before each mutation. This covers: accidental draws, wrong positions, unintended
resizes, and accidental deletes.

#### Visual save feedback

After clicking save, briefly flash saved rectangles green (200ms transition) to confirm
spatially which items were persisted.

#### Delete interaction

Two paths to delete: (1) select a rectangle by clicking it (blue selection border appears),
then press Delete key or click trash icon in sidebar. (2) Keyboard shortcut for power
users. Both paths are discoverable.

#### Zoom

Dense floor plans need zoom for precise positioning. Ctrl+scroll is the standard desktop
pattern. Ensure percentage-based coordinates remain correct at all zoom levels — the zoom
should be a CSS `transform: scale()` on the container.

### Dependencies

- Depends on Story 20.4 (Floor Plan Positions Database Schema and API)

### References

- Epic 20 Story 20.5: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR81: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

## Change Log

- 2026-03-22: Architecture decision — HTML/CSS overlays chosen over Konva.js.
- 2026-03-22: UX review — added AC 6 (undo), AC 7 (zoom), updated AC 2 (draw mode
  cursor), AC 3 (save feedback), AC 5 (delete interaction). Added tasks for undo, zoom,
  draw mode indicator, and delete UX.
