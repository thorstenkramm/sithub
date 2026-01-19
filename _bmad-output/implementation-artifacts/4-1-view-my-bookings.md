# Story 4.1: View My Bookings

Status: complete

## Story

As an employee,
I want to see my upcoming bookings,
so that I can confirm my reservations.

## Acceptance Criteria

1. **Given** I am authenticated  
   **When** I open "My Bookings"  
   **Then** I see a list of my future bookings  
   **And** each entry includes desk, room, area, and date

## Tasks / Subtasks

- [x] Add backend endpoint to list user's bookings (AC: 1)
  - [x] Implement `GET /api/v1/bookings` returning current user's future bookings
  - [x] Include desk, room, area names via spaces config lookup
  - [x] Filter to only future bookings (booking_date >= today)
  - [x] Order by booking_date ascending
- [x] Add frontend "My Bookings" view (AC: 1)
  - [x] Create MyBookingsView component
  - [x] Add route `/my-bookings`
  - [x] Display list with desk, room, area, and date
  - [x] Add navigation link to "My Bookings"
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for listing bookings
  - [x] Frontend: unit test for MyBookingsView
  - [x] Note: Cypress E2E deferred - requires real booking data setup
- [x] Update API documentation (AC: 1)
  - [x] Document GET /api/v1/bookings endpoint

## Dev Notes

- Return only future bookings (booking_date >= today's date)
- Include related data: desk name, room name, area name
- Use JSON:API format with `application/vnd.api+json` content type
- Reuse existing auth middleware for user context
- Added `FindDeskLocation` method to spaces config for location lookup

### Project Structure Notes

- Backend: `internal/bookings` (extend existing package)
- Frontend: `web/src/views/MyBookingsView.vue`, `web/src/api/bookings.ts`

### References

- PRD FR12: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 4.1: `_bmad-output/planning-artifacts/epics.md`
- Architecture: `_bmad-output/planning-artifacts/architecture.md`
- Existing bookings code: `internal/bookings/handler.go`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Implemented `GET /api/v1/bookings` endpoint returning user's future bookings
- Added `ListUserBookings` store function with date filtering and ordering
- Added `FindDeskLocation` method to spaces config for desk/room/area lookup
- Created `MyBookingAttributes` struct with full location details
- Created MyBookingsView.vue with formatted date display
- Added navigation bar to App.vue with "My Bookings" link
- Added `/my-bookings` route
- Added `fetchMyBookings` API client function
- Added backend handler tests (unauthorized, returns future bookings, empty list)
- Added frontend unit tests (user name, empty state, bookings list, auth redirects)
- Added OpenAPI documentation for GET /api/v1/bookings

### File List

**New Files:**
- `web/src/views/MyBookingsView.vue` - My Bookings view component
- `web/src/views/MyBookingsView.test.ts` - Unit tests for MyBookingsView

**Modified Files:**
- `internal/bookings/handler.go` - Added ListHandler and MyBookingAttributes
- `internal/bookings/handler_test.go` - Added tests for ListHandler
- `internal/bookings/store.go` - Added ListUserBookings function and BookingRecord type
- `internal/bookings/testhelpers_test.go` - Updated seedTestDeskData signature
- `internal/bookings/store_test.go` - Updated seedTestDeskData calls
- `internal/spaces/config.go` - Added DeskLocation type and FindDeskLocation method
- `internal/spaces/config_test.go` - Added TestFindDeskLocation
- `internal/startup/server.go` - Registered GET /api/v1/bookings route
- `web/src/api/bookings.ts` - Added MyBookingAttributes interface and fetchMyBookings
- `web/src/api/bookings.test.ts` - Added tests for fetchMyBookings
- `web/src/router/index.ts` - Added /my-bookings route
- `web/src/router/index.test.ts` - Updated route test
- `web/src/App.vue` - Added navigation bar with "My Bookings" link
- `web/src/views/testHelpers.ts` - Added createFetchMeMocker helper
- `web/src/views/AreasView.test.ts` - Refactored to use createFetchMeMocker
- `web/src/views/RoomsView.test.ts` - Refactored to use createFetchMeMocker
- `web/src/views/DesksView.test.ts` - Minor cleanup
- `api-doc/endpoints/bookings.yaml` - Added GET endpoint documentation
- `api-doc/openapi.yaml` - Added MyBookingAttributes, MyBookingResource, MyBookingsCollectionResponse schemas

### Change Log

- 2026-01-19: Story created and set to in-progress.
- 2026-01-19: Implementation completed - all tasks done, tests passing.
