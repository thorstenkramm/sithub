# Story 35.2: Floor-Plan Warnings for Free Items

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user booking from the floor plan,
I want to see a warning indicator on free items and confirm the warning when booking,
so that I am as well informed on the floor plan as I am on the tile views.

## Acceptance Criteria

1. A free item with a warning shows the shared warning icon on the top-level floor plan.
2. A free item with a warning shows the same shared warning icon on the drill-down floor plan.
3. Hovering the icon (either level) shows the warning message in the shared style.
4. A booked item shows only booker information — no warning icon and no warning in its hover.
5. Initiating a booking of a free warned item from the floor plan opens the warning confirmation
   dialog; confirming completes the booking, cancelling aborts without creating a booking.

## Tasks / Subtasks

- [x] Task 1: Add the warning icon overlay to free floor-plan items (AC: #1, #2, #3)
  - [x] In `InteractiveFloorPlan.vue`, render the shared warning icon (from 35.1) on free items in
        BOTH the non-area loop (`enrichedPositions`, template ~223-250) and the area-view desk loop
        (`deskPositions`, template ~143-218)
  - [x] Position it as a corner overlay consistent with existing icons — mirror the `.fp-item-lock`
        pattern (absolute, top-right, `z-index: 1`, CSS ~1827-1834); pick a corner that does not
        collide with the favorite heart (bottom-right, ~1766-1778) or the lock (top-right for
        reserved) — e.g. top-left for the warning
  - [x] The icon's hover shows the warning message in the shared style
- [x] Task 2: Remove warnings from booked items (AC: #4)
  - [x] In `enrichedPositions` tooltip construction (~1162-1172), stop appending `item.warning` for
        occupied items so booked items show only booker + name; keep the warning available for the
        free-item icon/hover
- [x] Task 3: Add warning confirmation to the floor-plan booking flow (AC: #5)
  - [x] Hook the shared warning confirmation (from 35.4) into `requestBooking`/`confirmPendingBooking`
        so booking a free warned item shows the dialog before the booking is created
  - [x] Remove or supersede the inline `fp-booking-warning` alert (~442-451) in favor of the uniform
        confirmation, respecting suppression via `useWarningSuppression`
- [x] Task 4: Tests (AC: #1-#5)
  - [x] Component tests: free warned item shows icon + hover message (both modes); booked item shows
        no warning; booking a warned free item opens the confirmation

## Dev Notes

Source: `private/epic-35.md` — floor-plan section (`img_30.png` shows the icon top-right + orange
hover on a free item). [Source: _bmad-output/planning-artifacts/epics.md#Story 35.2 / FR160,FR161]

### Component and rendering paths

`web/src/components/InteractiveFloorPlan.vue` is ONE component with two modes via `isAreaView`
(~706):

- Non-area (drill-down / item-group) items: template ~221-302, driven by `enrichedPositions`
  (computed ~1150-1186). Free items ~223-250, reserved ~252-268, busy/booked ~270-300.
- Area view (top-level): area rectangles ~121-140 + desk items ~143-218, driven by `deskPositions`
  (~1120-1148). Desk items reuse the same per-item logic.

Both paths need the warning icon on FREE items. Free/booked status is derived at ~1174-1181
(`availability === "occupied"` → busy; `reserved` → reserved; else free).

### Icon overlay positioning (precedent)

- Lock (reserved): `.fp-item-lock` absolute top-right 4px, `z-index: 1` (CSS ~1827-1834).
- Favorite heart: `.fp-favorite-heart` bottom-right, `z-index: 2` (CSS ~1766-1778).
- Avatar/initials (booked): inset:1px fill (CSS ~1803-1825).

Place the warning icon where it won't overlap these — **top-left** is free. Reuse the shared
`ItemWarning.vue` (35.1) in icon mode; keep the icon small (the matrix uses size 14, tiles 18).

### 🚨 Booked items must not show warnings (AC #4)

Today the warning is appended to the combined `v-tooltip` text for ALL items including occupied ones
(`enrichedPositions` ~1170-1172, joined with `\n` at ~1180). Per FR160, booked items show booker
only. Remove the warning from the occupied-item tooltip; surface it exclusively through the
free-item warning icon/hover.

### Booking flow hook (AC #5)

`requestBooking(itemID, label)` (~1416-1429) opens the booking dialog; `confirmPendingBooking()`
(~1474-1544) creates the booking. The dialog currently shows an inline `fp-booking-warning`
`v-alert` (~442-451, populated via `bookingItemWarning` at ~1336-1368) — that is NOT a blocking
confirmation. Replace it with the shared confirmation flow from 35.4: before the booking is created
(single or multi-day), if the item has a warning that is not suppressed, show the uniform dialog;
confirm → proceed, cancel → abort. Depends on 35.4's shared confirmation being available; if 35.2 is
implemented first, reuse the existing `useWarningSuppression` + a temporary call into the
ItemsView-style dialog, but prefer sequencing 35.4 before/with this task.

### Data

Warning field: `item.attributes.warning` (`web/src/api/items.ts:8`), carried into the component as
`ItemData.warning` (~668). No API change.

### Project Structure Notes

- Modified: `web/src/components/InteractiveFloorPlan.vue` (+ its test).
- Reuses `ItemWarning.vue` (35.1) and the shared confirmation (35.4).

### Testing standards summary

Vitest component tests exist for the floor plan (`InteractiveFloorPlan.test.ts`). Add cases for the
warning icon on free items (both modes), booked-item has no warning, and the confirmation on
booking. Run type-check, lint, vitest, build. A Cypress E2E for booking a warned floor-plan item is
a nice-to-have. [Source: .claude/rules/vue.md]

### References

- [Source: web/src/components/InteractiveFloorPlan.vue:143-218,221-302,1103-1186,1416-1429,1474-1544,442-451,1336-1368,1728-1834]
- [Source: web/src/api/items.ts:8]

## Dev Agent Record

### Agent Model Used

claude-fable-5

### Debug Log References

- 18 InteractiveFloorPlan tests pass; full suite 456; type-check/lint/build clean
- Verified in-browser (Chrome DevTools) against a seeded floor plan — see notes

### Completion Notes List

- Added a warning icon (shared `ItemWarning`, icon-variant "plain") on FREE floor-plan items in both
  the non-area (`enrichedPositions`) and area (`deskPositions`) loops, positioned top-left
  (`.fp-warning-icon`, z-index 2) so it doesn't collide with the lock/heart. Exposed a `warning`
  field on both computed maps, set only when the item is free.
- Removed the warning from the occupied-item tooltip so booked items show booker info only (FR160).
- Removed the inline `fp-booking-warning` alert from the booking dialog; gated
  `confirmPendingBooking` through the shared confirmation (split into `confirmPendingBooking` →
  `present([...], () => executeBooking(dates))`).
- Browser verification: free desks 1 & 3 (warned) show the orange icon, desk 2 (no warning) doesn't;
  booking warned desk 1 → uniform "WARNUNG!" dialog (orange message, matches img_33) → confirm books
  it; the now-booked desk shows only the "AA" booker with NO warning icon (FR160/FR161 confirmed).

### File List

- web/src/components/InteractiveFloorPlan.vue (modified)

### Change Log

- 2026-07-04: Implemented FR160/FR161 — floor-plan warning icons on free items + uniform booking
  confirmation; booked items show booker only.
