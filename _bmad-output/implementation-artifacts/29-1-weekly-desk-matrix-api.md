# Story 29.1: Weekly Desk Matrix API

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want the frontend to load one weekly desk matrix for the selected area,
so that the table can render quickly and consistently without dozens of follow-up requests.

## Acceptance Criteria

1. **Given** the frontend requests weekly table data for an area and week
   **When** the backend responds
   **Then** the payload contains all subareas and desks for that area in the configured SitHub
   order
   **And** each desk contains one cell per requested visible day

2. **Given** the frontend requests 5 visible days because the current user has weekends disabled
   **When** the backend responds
   **Then** the payload contains Monday through Friday only

3. **Given** a matrix cell is occupied
   **When** the backend responds
   **Then** the cell includes the booker's display name and user ID
   **And** enough booking metadata is present for permitted cancellation actions

4. **Given** a matrix cell is free but reserved for other users
   **When** the backend responds
   **Then** the payload clearly distinguishes between `bookable` free cells and `locked`
   reserved cells

## Tasks / Subtasks

- [x] Task 1: Define the endpoint contract and wire the route (AC: #1, #2, #3, #4)
  - [x] 1.1 Add a new authenticated route in `internal/startup/server.go` for
        `GET /api/v1/areas/:area_id/item-groups/matrix`
  - [x] 1.2 Create a new handler file in `internal/itemgroups/` for the weekly matrix endpoint
        instead of expanding `availability_handler.go`
  - [x] 1.3 Reuse the existing ISO-week parsing and weekday generation helpers already present in
        `internal/itemgroups/availability_handler.go` rather than reimplementing week math
  - [x] 1.4 Define explicit response structs for:
        - group-level matrix resource attributes
        - visible day metadata
        - item row metadata
        - per-cell booking state
  - [x] 1.5 Add a frontend API wrapper in `web/src/api/` with matching TypeScript interfaces so
        later UI stories can consume the endpoint without raw fetches
  - [x] 1.6 Update OpenAPI docs in `api-doc/openapi.yaml` and a new endpoint file under
        `api-doc/endpoints/` for the new route and schemas

- [x] Task 2: Implement one-query-per-request booking hydration for the full visible week
      (AC: #1, #2, #3, #4)
  - [x] 2.1 Collect all item-group and item IDs for the selected area by iterating the config
        slices in their existing order; do not sort alphabetically
  - [x] 2.2 Add a shared bookings-store helper in `internal/bookings/store.go` to fetch booking
        records for a set of item IDs and visible dates in one query
  - [x] 2.3 Use that helper from the new matrix handler to avoid N x item-group or N x day query
        patterns
  - [x] 2.4 Resolve booker display names using `users.FindDisplayNames`; for guest bookings, use
        the stored guest name as the visible name
  - [x] 2.5 Resolve the current user's email in the same way the existing items endpoint does so
        reservation checks stay aligned with current behavior
  - [x] 2.6 Mark reservation state using `areas.IsReserved(...)` semantics so the matrix matches
        the existing `reserved` behavior from `GET /item-groups/:item_group_id/items`

- [x] Task 3: Finalize the response shape for later table stories (AC: #1, #3, #4)
  - [x] 3.1 Return one JSON:API resource per item group, preserving configured order
  - [x] 3.2 Include group-level `days[]` metadata so the frontend does not have to recompute the
        visible week labels from scratch
  - [x] 3.3 Include item-row metadata needed by later stories:
        - `item_id`
        - `item_name`
        - `equipment`
        - `warning`
        - `reserved`
  - [x] 3.4 Include per-cell metadata needed by later stories:
        - `date`
        - `availability`
        - `booker_name` when occupied
        - `booker_user_id` when occupied and not guest
        - `booked_by_me` using the same semantics as the current items endpoint
        - `booking_id` only when the current user is the booking owner or an admin
  - [x] 3.5 Keep the contract minimal: do not add UI-only formatting fields such as badge text,
        localized labels, or prebuilt initials

- [x] Task 4: Cover the endpoint with focused backend and client tests (AC: #1, #2, #3, #4)
  - [x] 4.1 Add backend handler tests for:
        - area not found
        - invalid week
        - 5-day vs 7-day output
        - configured item-group and item ordering
        - occupied cells with display-name resolution
        - guest bookings showing guest names
        - reserved free cells for a disallowed user
        - `booking_id` visibility rules for cancel-capable viewers only
  - [x] 4.2 Add tests for any new bookings-store helper added for matrix hydration
  - [x] 4.3 Add a small frontend API test for the new `web/src/api/` wrapper if that file is
        introduced in this story

- [x] Task 5: Validate the story implementation surface
  - [x] 5.1 Run `go test ./internal/itemgroups/... ./internal/bookings/...`
  - [x] 5.2 If shared auth or user lookup helpers are changed, also run
        `go test ./internal/items/... ./internal/auth/... ./internal/users/...`
  - [x] 5.3 If a new frontend API wrapper is added, run
        `cd web && npx vitest run src/api/<new-matrix-api-test>.ts`
  - [x] 5.4 If a new frontend API wrapper is added, run `cd web && npm run type-check`

## Dev Notes

### Architecture & Patterns

- Backend stack is Go + Echo with JSON:API responses and authenticated `/api/v1/*` routes.
- Frontend stack is Vue 3 + TypeScript + Vuetify, with API access constrained to `web/src/api`.
- Existing area-level weekly summary data already lives in `internal/itemgroups/availability_handler.go`.
  This story should mirror its request shape (`week`, optional `days`) but return desk-level matrix data.
- Existing item-level day data already lives in `internal/items/handler.go`; the new endpoint must
  stay semantically aligned with its `reserved`, `booker_name`, and `booked_by_me` behavior.
- Existing booking cancellation rules are broader than the table UX. Do not blindly expose every
  cancel-capable booking ID to every caller just because `DELETE /api/v1/bookings/:id` can allow it.

### Recommended API Contract

Use a collection response where each resource represents one item group for the selected area.

Recommended path:
- `GET /api/v1/areas/:area_id/item-groups/matrix?week=YYYY-Www&days=5|7`

Recommended resource type:
- `item-group-weekly-matrix`

Recommended attributes shape:

```json
{
  "item_group_id": "ig-1",
  "item_group_name": "Room 101",
  "days": [
    { "date": "2026-04-13", "weekday": "MO" }
  ],
  "items": [
    {
      "item_id": "desk-1",
      "item_name": "Desk 1",
      "equipment": ["Dock", "Monitor"],
      "warning": "Near window",
      "reserved": true,
      "cells": [
        {
          "date": "2026-04-13",
          "availability": "occupied",
          "booker_name": "Ada Lovelace",
          "booker_user_id": "user-1",
          "booked_by_me": false,
          "booking_id": "booking-1"
        }
      ]
    }
  ]
}
```

Contract guidance:
- `days` should reflect the requested visible week and use the same `MO`/`TU` style abbreviations
  as the existing availability endpoint.
- `reserved` belongs at the item-row level because reservation is static for the desk, while
  occupancy varies per day.
- `booking_id` should only be present when the current user is the booking owner or an admin.
  Do not expose it for read-only occupied cells.
- Do not include pre-rendered initials; the frontend can derive those from `booker_name`.

### Key Code Locations

| Element | Location | Why it matters |
|---------|----------|----------------|
| Existing weekly area summary endpoint | `internal/itemgroups/availability_handler.go` | Reuse `week`/`days` parsing and JSON:API response pattern |
| Existing item day endpoint | `internal/items/handler.go` | Canonical `reserved`, `booker_name`, `booked_by_me` semantics |
| Booking lookup helper | `internal/bookings/store.go` | Best place for shared item/date query helper |
| Cancellation authorization | `internal/bookings/handler.go` (`DeleteHandler`) | Defines what metadata is sensitive vs. useful |
| Route registration | `internal/startup/server.go` | Wire the new authenticated endpoint |
| Current frontend weekly availability client | `web/src/api/itemGroupAvailability.ts` | Pattern for the new `web/src/api` wrapper |
| Current area item-groups page week handling | `web/src/views/ItemGroupsView.vue` | Existing consumer already has selected week + weekend preference |
| OpenAPI summary endpoint docs | `api-doc/endpoints/item-group-availability.yaml` | Best template for the new endpoint docs |

### Implementation Strategy

1. Keep the endpoint area-scoped and week-scoped.
2. Preserve configured order by iterating `area.ItemGroups` and each group's `Items` slice exactly
   as loaded from YAML.
3. Use one booking query for the full requested matrix window, then map results by `item_id|date`.
4. Reuse the same reservation semantics as the current items endpoint by resolving the caller's
   email and passing item location to `areas.IsReserved(...)`.
5. Enrich occupied cells with display names using `users.FindDisplayNames`, falling back to guest
   names for guest bookings.
6. Keep the output contract intentionally UI-ready for later stories by including item-level
   `equipment` and `warning`, so future matrix stories do not need secondary per-desk fetches.

### Testing Requirements

- Backend tests must be added alongside the new handler in `internal/itemgroups/`.
- If you extract a shared helper into `internal/bookings/store.go`, add or extend store tests there.
- Keep JSON:API response assertions explicit: status code, content type, resource type, and
  attribute fields.
- Test both a regular user and an admin to verify `booking_id` exposure rules.
- Include at least one guest booking scenario so the display-name fallback is pinned.

### Anti-Patterns to Avoid

- Do NOT implement the matrix by calling the existing single-day items endpoint in a loop.
- Do NOT issue one SQL query per item group or per day when one request-scoped query can hydrate
  the full visible week.
- Do NOT sort item groups or desks alphabetically; configured order is part of the product
  contract.
- Do NOT add the actual table UI in this story; this story is the backend/data contract foundation.
- Do NOT expose cancellation metadata for occupied cells that the later table UX must keep
  read-only; for this story that means omitting `booking_id` unless the viewer is the booking
  owner or an admin.
- Do NOT introduce a separate mobile-specific variant of this endpoint; mobile does not use this
  view.

### Project Structure Notes

- Architecture guidance says data access should live behind domain repositories/helpers and not as
  raw SQL in handlers. Follow that guidance for the new endpoint even if some older handlers are
  looser.
- No `project-context.md` was found in this repo, so there are no extra project-local AI rules
  beyond the artifacts cited below.

### References

- [Source: private/epic-29.md#Context]
- [Source: _bmad-output/planning-artifacts/epics.md#Epic 29 Stories: Desktop Weekly Table View]
- [Source: _bmad-output/planning-artifacts/prd.md#Executive Summary]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Executive Summary]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Core User Experience]
- [Source: _bmad-output/planning-artifacts/architecture.md#Naming Patterns]
- [Source: _bmad-output/planning-artifacts/architecture.md#Structure Patterns]
- [Source: _bmad-output/planning-artifacts/architecture.md#Format Patterns]
- [Source: _bmad-output/planning-artifacts/architecture.md#Integration Boundaries]
- [Source: internal/startup/server.go]
- [Source: internal/itemgroups/availability_handler.go]
- [Source: internal/itemgroups/availability_handler_test.go]
- [Source: internal/items/handler.go]
- [Source: internal/bookings/store.go]
- [Source: internal/bookings/handler.go#DeleteHandler]
- [Source: web/src/api/itemGroupAvailability.ts]
- [Source: web/src/views/ItemGroupsView.vue]
- [Source: api-doc/endpoints/item-group-availability.yaml]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- Story creation only; no implementation logs yet
- All tests pass on first run; no regressions introduced

### Completion Notes List

- Epic 29 reordered so the backend matrix API is the first execution story
- Story context includes a recommended route and JSON contract to prevent frontend/backend drift
- Cancellation metadata guidance is intentionally narrower than raw backend delete permissions to
  match the planned table UX
- Implemented `GET /api/v1/areas/:area_id/item-groups/matrix?week=YYYY-Www&days=5|7`
- Created `MatrixHandlerDynamic` in a separate handler file, reusing `parseISOWeek` and
  `weekdayDates` helpers from the availability handler
- Added `FindMatrixBookings` to `bookings/store.go` for single-query multi-item multi-date
  booking hydration
- Response shape: one `item-group-weekly-matrix` resource per item group, with `days[]` metadata,
  `items[]` with per-row equipment/warning/reserved, and per-cell booking state
- `booking_id` exposed only to the booking owner or admins (AC #3)
- Reservation state uses `areas.IsReserved(...)` at the item-row level (AC #4)
- Guest bookings show `guest_name` as `booker_name`, no `booker_user_id` (AC #3)
- Frontend TypeScript API wrapper and test added
- OpenAPI 3.1 endpoint doc and schemas added and linted clean

### Change Log

- 2026-04-15: Implemented story 29-1 — weekly desk matrix API endpoint, store helper, frontend
  wrapper, OpenAPI docs, and comprehensive backend + frontend tests

### File List

- _bmad-output/implementation-artifacts/29-1-weekly-desk-matrix-api.md (updated)
- internal/itemgroups/matrix_handler.go (new)
- internal/itemgroups/matrix_handler_test.go (new)
- internal/bookings/store.go (modified — added FindMatrixBookings, MatrixBookingInfo)
- internal/bookings/store_test.go (modified — added FindMatrixBookings tests)
- internal/startup/server.go (modified — wired matrix route)
- web/src/api/itemGroupMatrix.ts (new)
- web/src/api/itemGroupMatrix.test.ts (new)
- api-doc/openapi.yaml (modified — added matrix path and schemas)
- api-doc/endpoints/item-group-matrix.yaml (new)
