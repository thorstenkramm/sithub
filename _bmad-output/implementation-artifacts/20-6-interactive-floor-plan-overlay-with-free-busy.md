# Story 20.6: Interactive Floor Plan Overlay with Free/Busy

Status: backlog

## Story

As a user,
I want to see free/busy status on the floor plan and book items by clicking them,
So that I can visually find and book available items.

## Acceptance Criteria

1. **Given** I open the floor plan overlay for an item group
   **When** the overlay renders
   **Then** the floor plan image is shown with positioned items overlaid as rectangles

2. **Given** a weekday selector appears at the top of the overlay
   **When** I select a day
   **Then** free items show a green outline, busy items show a red semi-transparent overlay,
   and items I have booked show a primary/blue highlight

3. **Given** the floor plan opens for the current week
   **When** today is within the week
   **Then** today is pre-selected and past days are disabled

4. **Given** the floor plan opens for a future week
   **When** the overlay renders
   **Then** Monday is pre-selected

5. **Given** I hover over a free item
   **When** the tooltip appears
   **Then** it shows the item name, equipment list, and any warning

6. **Given** I click on a free item
   **When** the click is processed
   **Then** a booking is created for the selected day, the item status updates to busy,
   and a snackbar with an "Undo" action is shown for 5 seconds

7. **Given** I click "Undo" on the booking snackbar within 5 seconds
   **When** the undo is processed
   **Then** the booking is cancelled and the item reverts to free

8. **Given** weekend visibility is off in settings
   **When** the weekday selector renders
   **Then** Saturday and Sunday are not shown

9. **Given** I switch the selected weekday
   **When** availability data is being fetched
   **Then** the rectangles show a subtle loading state (reduced opacity) until new data
   arrives

10. **Given** I view the floor plan on a mobile device
    **When** the floor plan is too large for the screen
    **Then** I can pinch-to-zoom and scroll to navigate

## Tasks / Subtasks

- [ ] Create interactive floor plan component (AC: 1)
  - [ ] Create `InteractiveFloorPlan.vue` component
  - [ ] Display floor plan image inside a `position: relative` container
  - [ ] Load positions from floor plan positions API
  - [ ] Render positioned items as absolutely positioned `<div>` overlays using percentage
    coordinates (same system as the editor in Story 20.5)
- [ ] Fetch positions and availability data (AC: 1, 2, 9)
  - [ ] Fetch floor plan positions for the current item group
  - [ ] Fetch availability/booking data for the selected day
  - [ ] Map free/busy/mine state to each positioned item
  - [ ] Show loading state (reduced opacity) during data fetch
- [ ] Render free/busy/mine overlays (AC: 2)
  - [ ] Free items: green border, `cursor: pointer`
  - [ ] Busy items (other users): red semi-transparent background, `pointer-events: none`
  - [ ] My bookings: primary/blue highlight to distinguish from other users' bookings
  - [ ] Item label inside rectangle (use `label` from API if set, else item name)
  - [ ] Update overlays reactively when selected day changes
- [ ] Implement hover tooltip on free items (AC: 5)
  - [ ] Vuetify `v-tooltip` on free item rectangles
  - [ ] Show: item name, equipment list, warning (if any)
- [ ] Implement weekday selector (AC: 2, 3, 4, 8)
  - [ ] `v-btn-toggle` with weekday labels at the top of the overlay
  - [ ] Pre-select today if within the current week, otherwise Monday
  - [ ] Disable past days with `disabled` prop
  - [ ] Respect `showWeekends` preference to hide Saturday and Sunday
- [ ] Implement booking on click with undo (AC: 6, 7)
  - [ ] Click handler on free item rectangles
  - [ ] Create booking via `createBooking` API for the selected day
  - [ ] Show snackbar with "Undo" action button, 5-second timeout
  - [ ] If undo clicked within timeout: cancel the booking, revert item to free
  - [ ] If timeout expires: snackbar auto-dismisses, booking stands
  - [ ] Update item status to busy/mine after successful booking
- [ ] Implement mobile zoom/scroll (AC: 10)
  - [ ] Wrap floor plan container in a scrollable/zoomable wrapper
  - [ ] Support pinch-to-zoom on touch devices
  - [ ] Ensure rectangle overlays scale correctly with zoom
- [ ] Add unit tests for interactive floor plan component
- [ ] Verify E2E tests still pass

## Dev Notes

### Architecture Decision: HTML/CSS Overlays (shared with Story 20.5)

The viewer uses the same DOM overlay approach as the editor. Positioned `<div>` elements
with percentage-based coordinates render on top of the floor plan `<img>`. This allows
Vuetify tooltips on hover, `data-cy` selectors for Cypress, and CSS transitions for
free/busy state changes.

### Technical Implementation Guide

#### Shared coordinate system with editor

Positions from the API are stored as percentages (0-100) of the image dimensions. The
viewer renders them identically to the editor — `left`, `top`, `width`, `height` in `%`
inside a `position: relative` container. No coordinate conversion needed.

#### Free/busy/mine styling

```css
.floor-plan-item {
  position: absolute;
  border: 2px solid transparent;
  transition: background-color 0.2s, border-color 0.2s;
}

.floor-plan-item--free {
  border-color: rgb(var(--v-theme-success));
  cursor: pointer;
}

.floor-plan-item--busy {
  border-color: rgb(var(--v-theme-error));
  background-color: rgba(var(--v-theme-error), 0.3);
  pointer-events: none;
}

.floor-plan-item--mine {
  border-color: rgb(var(--v-theme-primary));
  background-color: rgba(var(--v-theme-primary), 0.2);
}

.floor-plan-item--loading {
  opacity: 0.4;
  transition: opacity 0.2s;
}
```

#### Reusable component structure

Extract a shared `FloorPlanBase.vue` component used by both the editor and the viewer,
handling image loading and rectangle rendering. The editor adds drag/resize handlers. The
viewer adds click-to-book, free/busy styling, and tooltips.

### UX Recommendations (Sally)

#### Hover tooltip on free items

Before clicking, users need to know what they're booking. Show a Vuetify tooltip on hover
with: item name, equipment list, and any warning text. This prevents blind bookings.

#### One-click booking with undo

The whole point of the interactive floor plan is speed. Skip the confirmation dialog for
bookings — instead show a snackbar with an "Undo" action for 5 seconds (Gmail pattern).
Fast action, safety net.

#### "My booking" highlight

After booking, the item shouldn't look the same as items booked by others. Use the
primary/blue theme color for "my" bookings to provide a sense of ownership.

#### Loading state during day switch

When switching weekdays, availability is re-fetched. During the fetch, reduce rectangle
opacity to 0.4 to signal "updating" — prevents the user from acting on stale colors.

#### Mobile pinch-to-zoom

On a phone screen, floor plans are too dense to tap accurately. The container must
support pinch-to-zoom. Consider `touch-action: none` on the container and implementing
a CSS `transform: scale()` wrapper, or using `overflow: auto` with a scaled inner element.

### Dependencies

- Depends on Story 20.4 (Floor Plan Positions Database Schema and API)
- Shares coordinate system and overlay approach with Story 20.5 (Floor Plan Editor)

### References

- Epic 20 Story 20.6: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR79: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

## Change Log

- 2026-03-22: Architecture decision — HTML/CSS overlays with shared coordinate system.
- 2026-03-22: UX review — added AC 5 (hover tooltip), AC 6 updated (undo snackbar),
  AC 7 (undo action), AC 9 (loading state), AC 10 (mobile zoom). Updated AC 2 to include
  "mine" highlight. Added tasks for tooltip, undo, mobile zoom, and loading state.
