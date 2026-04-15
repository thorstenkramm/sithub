# Story 27.2: Fix Reserved Area Display on Floor Plan

Status: done

## Story

As a user,
I want reserved areas on the floor plan to show correct availability and allow drill-down,
so that I can see who is in a reserved room even though I cannot book there.

## Acceptance Criteria

1. **Given** I view the floor plan and "People & Finance" is a reserved area with 3 of 4
   desks free
   **When** I look at the area overlay
   **Then** it shows "3/4 free" with a green indicator (not "0/4" red)

2. **Given** a desk in a reserved area is booked by someone
   **When** I look at the floor plan
   **Then** their avatar is displayed on that desk

3. **Given** I click on a reserved area on the floor plan
   **When** I drill down into it
   **Then** I can see individual desk availability and who has booked

4. **Given** I see a free desk in a reserved area after drill-down
   **When** I look at that desk
   **Then** it is blurred/dimmed with a "reserved" message and I cannot book it

## Tasks / Subtasks

- [x] Task 1: Fix area free/busy count to include reserved items (AC: #1)
  - [x] 1.1 In `InteractiveFloorPlan.vue`, locate the availability count computation
        (line ~871-875): `response.data.filter(item => item.attributes.availability ===
        "available" && item.attributes.reserved !== true).length`
  - [x] 1.2 Remove the `&& item.attributes.reserved !== true` condition — reserved items
        that are available should count as free for display purposes
  - [x] 1.3 Verify the area overlay shows correct "X/Y free" counts
- [x] Task 2: Verify avatars display on booked desks in reserved areas (AC: #2)
  - [x] 2.1 Check that the avatar rendering in `deskPositions` computed (line ~993) does not
        skip reserved items — avatars should render for any booked desk regardless of
        reservation status
  - [x] 2.2 If avatars are missing, ensure the `booker_name` / `booker_user_id` fields are
        populated for reserved item bookings
- [x] Task 3: Verify drill-down works for reserved areas (AC: #3)
  - [x] 3.1 Check that `handleAreaClick()` (line ~1118) does not block drill-down into
        reserved areas
  - [x] 3.2 After drill-down, verify item positions and availability load correctly
- [x] Task 4: Show reserved overlay on free desks after drill-down (AC: #4)
  - [x] 4.1 In the drill-down desk view, reserved items with `availability === "available"`
        must show a "reserved" overlay (blurred/dimmed, non-bookable)
  - [x] 4.2 Check the `deskPositions` computed (line ~1015-1025): currently `reserved`
        status overrides `free` — a reserved available item shows as "reserved" with the
        existing `fp-item--reserved` CSS class, which is correct behavior
  - [x] 4.3 Verify the booking dialog does NOT open when clicking a reserved item
- [x] Task 5: Validate (AC: #1-#4)
  - [x] 5.1 Run `npm run lint` and fix findings
  - [x] 5.2 Run `npm run type-check` and fix findings
  - [x] 5.3 Run `npm run build` and verify no build errors
  - [x] 5.4 Run `npx vitest run` and verify no regressions

### Review Findings

- [x] [Review][Patch] Occupied reserved desks still render as locked reserved items instead of busy avatars [web/src/components/InteractiveFloorPlan.vue:1024]

## Dev Notes

### Architecture & Patterns

- **Primary file**: `web/src/components/InteractiveFloorPlan.vue`
- **No backend changes** — the `reserved` field is already served correctly by the API
- The core bug is line ~874: `item.attributes.reserved !== true` excludes reserved items
  from the free count, making reserved areas show "0/N free" (all red)

### Key Code Locations

| Element | Location |
|---------|----------|
| Area availability count | `InteractiveFloorPlan.vue:871-875` |
| `deskPositions` computed | `InteractiveFloorPlan.vue:993-1030` |
| `handleAreaClick()` | `InteractiveFloorPlan.vue:1118` |
| `fp-item--reserved` CSS | `InteractiveFloorPlan.vue:1629` |
| Area overlay rendering | Template lines ~120-160 |

### Implementation Strategy

The main fix is a one-line change: remove `&& item.attributes.reserved !== true` from the
availability count filter. The rest is verification that existing behavior is correct:
- Avatars already render for booked items regardless of reservation
- Drill-down already works for all areas with floor plans
- Reserved items after drill-down already show the `fp-item--reserved` overlay

### Anti-Patterns to Avoid

- Do NOT allow booking of reserved items — the overlay and booking block must remain
- Do NOT change the backend `reserved` field logic
- Do NOT remove the `fp-item--reserved` CSS styling

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List

### Change Log
