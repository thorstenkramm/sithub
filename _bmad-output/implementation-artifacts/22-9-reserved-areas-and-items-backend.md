# Story 22.9: Reserved Areas and Items — Backend

Status: done

## Story

As an operator,
I want to restrict areas and items to specific users via YAML configuration,
so that shared resources can be reserved for designated teams or individuals.

## Acceptance Criteria

1. **Given** an area has `reserved_for: [anna@sithub.local, tk@system42.io]`
   **When** `alex@sithub.local` attempts to book any item in that area
   **Then** the booking is rejected with a 403 error naming the area

2. **Given** an item group has `reserved_for: [tk@system42.io]`
   **When** `anna@sithub.local` (who has area access) attempts to book in that group
   **Then** the booking is rejected with a 403 error naming the item group

3. **Given** a child item has `reserved_for: [user2@example.com]` but the parent
   area does not include `user2@example.com`
   **When** the server starts
   **Then** startup fails with a validation error explaining the conflict

4. **Given** `reserved_for` is missing or null at any level
   **When** a booking is attempted
   **Then** no reservation restriction applies at that level

5. **Given** a booking is rejected due to reservation
   **When** the error response is returned
   **Then** it includes a clear message: "This area/item is reserved. You do not
   have access."

6. **Given** the items list API is called by a user without access
   **When** the response renders
   **Then** reserved items include a `reserved: true` attribute so the frontend
   can disable them

## Tasks / Subtasks

- [ ] Task 1: Add `reserved_for` to YAML structs (AC: 1, 4)
  - [ ] 1.1 In `internal/areas/config.go`: add `ReservedFor []string` field
    with `yaml:"reserved_for,omitempty"` to `Area`, `ItemGroup`, and `Item`
    structs
  - [ ] 1.2 Update `sithub_areas.schema.json`: add `reserved_for` property
    (type: array of strings) at all three levels
- [ ] Task 2: Validate hierarchical reservation consistency (AC: 3)
  - [ ] 2.1 In `internal/areas/config.go`: add `ValidateReservations(cfg *Config)`
    function called from `Load()`
  - [ ] 2.2 For each item/item-group with `reserved_for`, check that ALL listed
    emails are also in the parent's `reserved_for` (if the parent has one).
    If not, return an error: "item '{itemID}' reserves for '{email}' but
    parent area '{areaID}' does not include this user"
  - [ ] 2.3 The check is: if parent has `reserved_for` set (non-empty), then
    child's `reserved_for` entries must be a subset of parent's
- [ ] Task 3: Enforce reservation in booking handler (AC: 1, 2, 5)
  - [ ] 3.1 In `internal/bookings/handler.go`: add a reservation check in
    `CreateHandlerDynamic` after `FindItemLocation` and before limits check.
    Look up the user's email from the database using their user ID
  - [ ] 3.2 Create `checkReservation(userEmail string, loc *areas.ItemLocation) error`
    that checks item → item group → area reservation lists. If ANY level
    has `reserved_for` set and the user's email is not in it, reject
  - [ ] 3.3 Return 403 Forbidden with message via `api.WriteForbidden` or a
    new `api.WriteReserved` helper
- [ ] Task 4: Add `reserved` flag to items API response (AC: 6)
  - [ ] 4.1 In `internal/items/handler.go`: the items list handler returns
    item attributes. Add a `reserved` boolean field that is `true` when
    the item (or its parents) has `reserved_for` set and the current user's
    email is NOT in the list
  - [ ] 4.2 This requires passing the authenticated user's email to the items
    handler — retrieve from `auth.GetUserFromContext(c)`
- [ ] Task 5: Write tests (AC: 1, 2, 3, 4, 5, 6)
  - [ ] 5.1 Test YAML validation: valid hierarchy, invalid hierarchy (child
    email not in parent), empty reserved_for
  - [ ] 5.2 Test booking rejection: user not in reserved list, user in list
    (allowed), no reservation (allowed)
  - [ ] 5.3 Test items API: `reserved` flag set correctly
  - [ ] 5.4 Run `go test ./...`, `golangci-lint run ./...`
- [ ] Task 6: Update API documentation (AC: 5, 6)
  - [ ] 6.1 Update `api-doc/endpoints/bookings.yaml` with 403 response
  - [ ] 6.2 Update item attributes schema with `reserved` boolean
  - [ ] 6.3 Lint with redocly

## Dev Notes

### Reservation Hierarchy Logic

```
Area reserved_for: [A, B, C]
  ItemGroup reserved_for: [A, B]      ← subset of area
    Item reserved_for: [A]            ← subset of item group
    Item reserved_for: null           ← inherits item group [A, B]
  ItemGroup reserved_for: null        ← inherits area [A, B, C]
```

Check order for booking: item → item group → area. The FIRST level that has
`reserved_for` set determines access. If none have it, booking is open.

### User Email Lookup

The booking handler has `user.ID` from the auth context. To check reservation,
you need the user's email. Options:

1. Add `Email` field to `auth.User` struct (requires updating cookie encoding)
2. Query `users.FindByID()` which already returns email

Option 2 is simpler and the query is already used in the booking flow.

### Existing Pattern: `handleBookingLimits`

Follow the same extraction pattern: create `handleReservation()` that returns
nil on success or writes the error response + returns `errResponseWritten`.

### Config Validation on Startup

Call `ValidateReservations()` in `server.go` after `areas.Load()` — same
location where `ValidateFloorPlans()` is called (line ~66).

### Files to Change

| File | Change |
| --- | --- |
| `internal/areas/config.go` | Add ReservedFor field, validation |
| `internal/bookings/handler.go` | Reservation check in create flow |
| `internal/items/handler.go` | Add `reserved` flag to response |
| `internal/startup/server.go` | Call ValidateReservations |
| `sithub_areas.schema.json` | Add reserved_for schema |
| `api-doc/` | Update API docs |

### References

- [Source: private/epic-22.md — "Reserved areas and items"]
- [Source: internal/areas/config.go — Area/ItemGroup/Item structs]
- [Source: internal/bookings/handler.go — CreateHandlerDynamic flow]

### Review Findings

- [x] [Review][Patch] Apply reservation checks to the actual booking target instead of only the acting user [internal/bookings/handler.go:493]
- [x] [Review][Patch] Return a reservation-specific 403 message and document the `reserved` item attribute in the API surface [internal/api/errors.go:30]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
