# Story 21.6: Booking Limits — Backend Configuration and Enforcement

Status: done

## Story

As an administrator,
I want to configure booking limits in the TOML config and areas YAML,
so that I can control how far in advance users can book and how many bookings
each person is allowed.

## Acceptance Criteria

1. **Given** the TOML config has `[bookings] weeks_in_advanced = 7`,
   **when** a user tries to book a date more than 7 weeks ahead,
   **then** the API returns a 400 error with a message indicating the date is too far
   in the future.

2. **Given** the TOML config has `[bookings] max_bookings_per_person = 10`,
   **when** a user already has 10 active (future) bookings and tries to create another,
   **then** the API returns a 409 error with a message like "you have reached the
   maximum of 10 active bookings".

3. **Given** an area in the YAML has `max_bookings_per_person: 5`,
   **when** a user has 5 active bookings within that area and tries to book another item
   in the same area,
   **then** the API returns a 409 error naming the area in the message.

4. **Given** an item group in the YAML has `max_bookings_per_person: 3`,
   **when** a user has 3 active bookings across items in that group,
   **then** the API returns a 409 error naming the item group.

5. **Given** an item in the YAML has `max_bookings_per_person: 2`,
   **when** a user has 2 active bookings for that specific item,
   **then** the API returns a 409 error naming both the item group and item
   (e.g. "Room 1, Desk 1").

6. **Given** the application starts,
   **when** an authenticated user calls `GET /api/v1/settings`,
   **then** the response includes `weeks_in_advanced` so the frontend can limit the
   week selector.

## Tasks / Subtasks

- [x] Task 1: Add BookingsConfig to TOML config (AC: 1, 2)
  - [x] 1.1 Add `BookingsConfig` struct with `WeeksInAdvanced` and `MaxBookingsPerPerson`
    to `internal/config/config.go`
  - [x] 1.2 Add viper defaults (`weeks_in_advanced: 5`, `max_bookings_per_person: 0`)
  - [x] 1.3 Add `[bookings]` section to `sithub.example.toml` with documentation
- [x] Task 2: Add `max_bookings_per_person` to areas YAML structs (AC: 3, 4, 5)
  - [x] 2.1 Add `MaxBookingsPerPerson int` field to `Area`, `ItemGroup`, and `Item`
    structs in `internal/areas/config.go`
  - [x] 2.2 Update `sithub_areas.schema.json` with `max_bookings_per_person` at all
    three levels
- [x] Task 3: Add store function to count future bookings (AC: 2, 3, 4, 5)
  - [x] 3.1 Add `CountUserFutureBookings(ctx, db, userID, itemIDs)` to
    `internal/bookings/store.go`
  - [x] 3.2 When `itemIDs` is nil, count all future bookings; otherwise filter by item IDs
- [x] Task 4: Add booking limit enforcement (AC: 1, 2, 3, 4, 5)
  - [x] 4.1 Add `BookingLimits` struct and `ErrBookingLimitExceeded` sentinel error
    to `internal/bookings/handler.go`
  - [x] 4.2 Add `enforceBookingLimits()` function checking item, item group, area, and
    global limits in order
  - [x] 4.3 Add `weeks_in_advanced` validation in `validateRequestFieldsMultiDay()`
  - [x] 4.4 Update `CreateHandlerDynamic()` signature to accept `*BookingLimits`
  - [x] 4.5 Wire limits enforcement into the create handler before booking creation
  - [x] 4.6 Update `CreateHandler()` to pass `nil` limits (backward-compatible)
- [x] Task 5: Add settings endpoint (AC: 6)
  - [x] 5.1 Create `internal/system/settings.go` with `SettingsHandler`
  - [x] 5.2 Register `GET /api/v1/settings` route in `internal/startup/server.go`
- [x] Task 6: Wire config through startup (AC: 1, 2, 3, 4, 5, 6)
  - [x] 6.1 Create `BookingLimits` from `cfg.Bookings` in `server.go`
  - [x] 6.2 Pass limits to `registerRoutes` and `CreateHandlerDynamic`
  - [x] 6.3 Pass `weeksInAdvanced` to `SettingsHandler`
- [x] Task 7: Add API documentation (AC: 6)
  - [x] 7.1 Create `api-doc/endpoints/settings.yaml`
  - [x] 7.2 Add `/settings` path to `api-doc/openapi.yaml`
  - [x] 7.3 Lint with redocly
- [x] Task 8: Write tests (AC: 1, 2, 3, 4, 5, 6)
  - [x] 8.1 Add `TestCountUserFutureBookingsAll` and `TestCountUserFutureBookingsFiltered`
    to `store_test.go`
  - [x] 8.2 Add `TestCreateHandlerItemLimitExceeded` to `handler_test.go`
  - [x] 8.3 Add `TestCreateHandlerItemGroupLimitExceeded` to `handler_test.go`
  - [x] 8.4 Add `TestCreateHandlerGlobalLimitExceeded` to `handler_test.go`
  - [x] 8.5 Add `TestCreateHandlerWeeksInAdvancedLimit` to `handler_test.go`
  - [x] 8.6 Add `TestCreateHandlerWithinLimitsSuccess` to `handler_test.go`
  - [x] 8.7 Add `TestSettingsHandlerReturnsConfig` to `settings_test.go`
  - [x] 8.8 Fix `server_test.go` calls to pass `nil` booking limits
  - [x] 8.9 Run full Go test suite — all pass
  - [x] 8.10 Run golangci-lint — 0 issues

## Dev Notes

### Limit Hierarchy

Limits are checked from most specific to least specific:
1. **Item level** — counts bookings for that specific item ID
2. **Item group level** — counts bookings for all item IDs in the group
3. **Area level** — counts bookings for all item IDs across all groups in the area
4. **Global level** — counts all future bookings regardless of item

All applicable limits are checked. The first one that is exceeded returns a 409.
A limit value of 0 means unlimited (not checked).

### Weeks in Advanced Calculation

The booking horizon is calculated as: current Monday + `weeks_in_advanced` weeks.
Dates on or after this cutoff are rejected with a 400 error.

### Guest Bookings

Guest bookings are excluded from limit checks because each guest gets a unique user ID
(`guest-{uuid}`), so their bookings never accumulate against a real user.

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `BookingsConfig` to TOML config with viper defaults
- Added `max_bookings_per_person` to Area, ItemGroup, Item YAML structs
- Added `CountUserFutureBookings` store function with optional item ID filtering
- Added `enforceBookingLimits` checking 4 scope levels with descriptive error messages
- Added `weeks_in_advanced` validation in date validation
- Created `GET /api/v1/settings` endpoint exposing `weeks_in_advanced`
- Updated `sithub_areas.schema.json` with `max_bookings_per_person` at all 3 levels
- 7 new Go tests, all passing; golangci-lint 0 issues
- OpenAPI spec lints clean

### File List

- `internal/config/config.go` — Added `BookingsConfig` struct and viper defaults
- `internal/areas/config.go` — Added `MaxBookingsPerPerson` to Area, ItemGroup, Item
- `internal/bookings/store.go` — Added `CountUserFutureBookings()`
- `internal/bookings/handler.go` — Added `BookingLimits`, `enforceBookingLimits()`,
  weeks validation, updated `CreateHandlerDynamic` signature
- `internal/bookings/store_test.go` — 2 new tests for count function
- `internal/bookings/handler_test.go` — 5 new tests for limits and weeks enforcement
- `internal/system/settings.go` — New settings endpoint handler
- `internal/system/settings_test.go` — Settings handler test
- `internal/startup/server.go` — Wired booking limits and settings route
- `internal/startup/server_test.go` — Fixed `registerRoutes` calls
- `sithub.example.toml` — Added `[bookings]` section
- `sithub_areas.schema.json` — Added `max_bookings_per_person` at area/group/item levels
- `api-doc/openapi.yaml` — Added `/settings` path
- `api-doc/endpoints/settings.yaml` — New settings endpoint spec
