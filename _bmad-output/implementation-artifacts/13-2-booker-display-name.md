# Story 13.2: Booker Display Name

Status: done

## Story

As a user,
I want to see who has booked an item,
So that I know which colleagues are in the office or using a resource.

## Acceptance Criteria

1. **Given** I am viewing items in an item group for a specific date
   **When** an item is booked
   **Then** I see the booker's display name alongside the booking status
   **And** the display name is resolved from the users table (not stored in the booking)

2. **Given** an item is available (not booked)
   **When** the item is displayed
   **Then** no booker name is shown
   **And** the item is clearly marked as available

3. **Given** I am viewing Today's Presence for an area
   **When** the presence list is displayed
   **Then** each entry shows the user's display name

## Tasks / Subtasks

- [x] Update items API to include booker display name (AC: 1)
  - [x] Modify the items list handler to resolve booker names for ALL users (not just admins)
  - [x] Add `booker_name` field to the items JSON:API attributes (only when booked)
  - [x] `booking_date` query parameter already existed for availability
  - [x] Ensure guest bookings show the guest name instead of a user lookup
- [x] Update ItemsView.vue to display booker name (AC: 1, 2)
  - [x] Show booker name on booked item tiles (below equipment, using text-medium-emphasis)
  - [x] Ensure available items show no name and are clearly marked "Available"
  - [x] `data-cy="item-booker"` attribute already present (repurposed from admin-only)
- [x] Verify Today's Presence already shows display names (AC: 3)
  - [x] AreaPresenceView already shows `user_name` from the presence API
  - [x] Presence handler resolves names via `users.FindDisplayNames` - confirmed working
- [x] Update Vitest unit tests (AC: 1, 2)
  - [x] Test items API response with booker name present (occupied)
  - [x] Test items API response with available item (no booker name)
- [x] Add Cypress E2E test (AC: 1, 2, 3)
  - [x] Mocked items view with occupied item, verify booker name visible
  - [x] Verify available item shows no booker name

## Dev Notes

### Current Items API Response

The items endpoint returns item data per item group. Current `ItemAttributes`:

```typescript
interface ItemAttributes {
  name: string
  equipment: string[]
  is_available: boolean
}
```

The `is_available` field is date-dependent. The booker's name is NOT currently
included. This story adds it.

### Backend Change: Items Handler

The items list handler in `internal/items/` needs to:
1. Accept a `date` query parameter (may already exist for availability)
2. LEFT JOIN bookings table for the given date
3. LEFT JOIN users table to resolve the booker's display name
4. Return `booked_by_name` in the JSON:API attributes (null/omitted when available)

For guest bookings, use the `guest_name` from the bookings record directly instead
of looking up a user.

### Frontend Change: ItemsView.vue

Item tiles currently show name, equipment tags, and an availability indicator.
Add the booker's display name below the status when booked. Use subtle styling
(e.g., `text-medium-emphasis`) to keep the name secondary to the item name.

### Today's Presence (AC: 3)

The `AreaPresenceView` and its backend handler already return `user_name` per
booking. This AC is likely already satisfied. Verify and add a test if missing.

### API Type Update

```typescript
interface ItemAttributes {
  name: string
  equipment: string[]
  is_available: boolean
  booked_by_name?: string  // NEW: only present when booked
}
```

### References

- PRD FR39: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 13.2: `_bmad-output/planning-artifacts/epics.md`
- Items handler: `internal/items/`
- ItemsView: `web/src/views/ItemsView.vue`
- Presence handler: `internal/areas/presence_handler.go`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None needed - clean implementation.

### Completion Notes List

- Backend items handler (`internal/items/handler.go`) was already returning `booker_name` but
  only for admin users. Changed to resolve and return `booker_name` for ALL users.
- `booking_id` remains admin-only.
- `FindItemBookings` in `internal/bookings/store.go` now also selects `is_guest` and
  `guest_name` from bookings table. `ItemBookingInfo` struct updated with `IsGuest` and
  `GuestName` fields.
- `resolveBookerNames` now handles guest bookings by using `GuestName` directly instead of
  doing a users table lookup.
- Frontend `ItemsView.vue` updated to show booker name for all users (removed admin-only gate).
- AC 3 (Today's Presence) was already satisfied - `presence_handler.go` resolves display names
  via `users.FindDisplayNames`.
- All 16 Go packages pass, 98 Vitest tests pass, 37 Cypress E2E tests pass (1 new).
- Go lint clean, TypeScript type-check clean, ESLint clean, build clean.
- Go duplication 2.18% (under 3% threshold), TS duplication 0%.

### File List

- `internal/items/handler.go` - Removed admin-only gating for booker_name, simplified
  loadItemBookings, enhanced resolveBookerNames for guest bookings
- `internal/items/handler_test.go` - Updated non-admin test, added guest booking test
- `internal/bookings/store.go` - Added IsGuest/GuestName to ItemBookingInfo, updated query
- `internal/bookings/store_test.go` - Added guest booking info test
- `web/src/api/items.ts` - Updated comment on booker_name field
- `web/src/views/ItemsView.vue` - Show booker_name for all users (not just admin)
- `web/src/views/ItemsView.test.ts` - Added booker name visibility tests
- `web/cypress/e2e/items.cy.ts` - Added E2E test for booker name display
- `web/cypress/support/flows.ts` - Added bookerName param to createMockItem