# Story 5.2: Today's Presence by Area

Status: complete

## Story

As an employee,
I want to see who is in the office today by area,
So that I can coordinate with colleagues.

## Acceptance Criteria

1. **Given** I view today's presence for an area  
   **When** the list is displayed  
   **Then** I see all users with bookings in that area for today

## Tasks / Subtasks

- [x] Add backend endpoint for area presence (AC: 1)
  - [x] Implement `GET /api/v1/areas/:id/presence?date=YYYY-MM-DD`
  - [x] Return list of users with bookings in that area for the date
  - [x] Default to today if date not provided
  - [x] Handle area not found (404)
- [x] Add frontend area presence view (AC: 1)
  - [x] Create AreaPresenceView component
  - [x] Show list of users present in the area
  - [x] Add date picker (default to today)
  - [x] Add route `/areas/:areaId/presence`
- [x] Add navigation to area presence (AC: 1)
  - [x] Add "View Presence" link from AreasView
- [x] Add tests (AC: 1)
  - [x] Backend: handler tests for area presence
  - [x] Frontend: unit tests for AreaPresenceView
- [x] Update API documentation (AC: 1)
  - [x] Document GET /api/v1/areas/:id/presence endpoint

## Dev Notes

- Similar to room bookings, but aggregates by area (all rooms in the area)
- Returns user info for all bookings in all rooms of the area
- Can reuse patterns from room bookings implementation

### Project Structure Notes

- Backend: `internal/areas/presence_handler.go` (new)
- Frontend: `web/src/views/AreaPresenceView.vue` (new)

### References

- PRD FR16: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 5.2: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.

### Completion Notes List

- Backend endpoint returns user_id, user_name, desk_id, desk_name, room_id, room_name for each booking in the area
- Extracted `api.BuildINClause` helper to share IN clause building logic between rooms and areas handlers
- Frontend shows list of present users with room and desk location, date picker for navigation
- Added "View Presence" link to AreasView

### File List

**Backend:**
- `internal/areas/presence_handler.go` (new)
- `internal/areas/presence_handler_test.go` (new)
- `internal/api/response.go` (modified - added BuildINClause)
- `internal/api/response_test.go` (modified - added tests)
- `internal/rooms/bookings_handler.go` (modified - use BuildINClause)
- `internal/startup/server.go` (modified - register route)

**Frontend:**
- `web/src/api/areaPresence.ts` (new)
- `web/src/views/AreaPresenceView.vue` (new)
- `web/src/views/AreaPresenceView.test.ts` (new)
- `web/src/router/index.ts` (modified - add route)
- `web/src/router/index.test.ts` (modified - update test)
- `web/src/views/AreasView.vue` (modified - add link)
- `web/src/views/RoomBookingsView.test.ts` (modified - use shared helpers)

**API Docs:**
- `api-doc/endpoints/area-presence.yaml` (new)
- `api-doc/openapi.yaml` (modified)

### Change Log

- 2026-01-19: Story created and set to in-progress.
- 2026-01-19: Story completed.
