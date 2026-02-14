# Story 12.4: Rename Rooms and Desks to Items in Frontend

Status: done

## Story

As a user,
I want the UI to use domain-neutral terminology,
So that I see consistent labels regardless of what I am booking.

## Acceptance Criteria

1. **Given** the frontend uses components and routes named `RoomsView`, `DesksView`, etc.
   **When** the rename is applied
   **Then** components and routes use item-group and item terminology
   **And** the Vue Router paths match the new API routes

2. **Given** the areas list view shows a "VIEW ROOM" button on each area tile
   **When** the rename is applied
   **Then** the button reads "BOOK"
   **And** no UI element references "room" or "desk" in user-facing text

3. **Given** the item detail view shows a "BOOK THIS DESK" button
   **When** the rename is applied
   **Then** the button reads "BOOK THIS ITEM"

4. **Given** a booking is successfully created
   **When** the confirmation message is displayed
   **Then** the message references the item name from the configuration
   (e.g., "Parking Lot 1 booked successfully") rather than the generic term "desk"

5. **Given** the Pinia stores and API service files reference rooms/desks
   **When** the rename is applied
   **Then** all store names, API paths, and TypeScript types use the new terminology
   **And** `npm run type-check` and `npm run build` succeed without errors
   **And** `npx vitest run` passes all unit tests

## Tasks / Subtasks

- [x] Rename Vue components (AC: 1)
  - [x] RoomsView.vue -> ItemGroupsView.vue
  - [x] DesksView.vue -> ItemsView.vue
  - [x] RoomBookingsView.vue -> ItemGroupBookingsView.vue
- [x] Rename API service files (AC: 5)
  - [x] rooms.ts -> itemGroups.ts
  - [x] desks.ts -> items.ts
  - [x] roomBookings.ts -> itemGroupBookings.ts
- [x] Update Vue Router paths and names (AC: 1)
  - [x] /areas/:areaId/rooms -> /areas/:areaId/item-groups
  - [x] /rooms/:roomId/desks -> /item-groups/:itemGroupId/items
  - [x] /rooms/:roomId/bookings -> /item-groups/:itemGroupId/bookings
- [x] Update UI labels (AC: 2, 3, 4)
  - [x] "VIEW ROOM" -> "BOOK"
  - [x] "BOOK THIS DESK" -> "BOOK THIS ITEM"
  - [x] Confirmation messages use item name from config
- [x] Update all Vitest unit tests (AC: 5)
- [x] Run type-check and build (AC: 5)

## Dev Notes

### Design Decisions

- Button labels simplified: "BOOK" instead of "VIEW ROOM" (cleaner UX)
- Confirmation messages use the actual item name from configuration
- All TypeScript types updated to match new API field names

### References

- PRD FR4-FR16, FR42: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 12.4: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- All Vue components, routes, API services renamed
- UI labels use domain-neutral terminology throughout
- 90 Vitest unit tests passing with renamed types
- Type-check and build succeed

### Key Files

- `web/src/views/ItemGroupsView.vue` (was RoomsView.vue)
- `web/src/views/ItemsView.vue` (was DesksView.vue)
- `web/src/views/ItemGroupBookingsView.vue` (was RoomBookingsView.vue)
- `web/src/api/itemGroups.ts` (was rooms.ts)
- `web/src/api/items.ts` (was desks.ts)
- `web/src/api/itemGroupBookings.ts` (was roomBookings.ts)
- `web/src/router/index.ts`
- `web/src/App.vue`

### Change Log

- 2026-02-09: Story created retroactively. Implementation was part of Epic 12 commit.
