# Story 18.7: Connection Lost Error Messaging

Status: done

## Story

As a user,
I want to see a clear "Connection to server lost" error when the backend is unavailable,
So that I understand the real problem instead of seeing misleading messages like
"no areas available."

## Acceptance Criteria

1. **Given** the backend server is unavailable
   **When** the frontend attempts to load data
   **Then** a clear error message "Connection to server lost" is displayed

2. **Given** the backend was available and then goes down
   **When** subsequent API calls fail
   **Then** the error message is shown instead of empty or misleading content

## Tasks / Subtasks

- [x] Update API client (`web/src/api/client.ts`)
  - [x] Wrap `fetch()` call in try/catch to detect network failures (`TypeError`)
  - [x] Throw `ApiError` with status `0` and message "Connection to server lost"
  - [x] Export `CONNECTION_LOST_MESSAGE` constant
  - [x] Export `isConnectionError()` helper function
- [x] Update views to show connection error message
  - [x] `AreasView.vue` — show `{{ areasError }}` instead of hardcoded text; catch connection
    error in `fetchMe`
  - [x] `ItemGroupsView.vue` — catch connection error in `fetchMe` and `fetchAreas`
  - [x] `ItemsView.vue` — catch connection error in `fetchMe`, `fetchAreas`, `loadItems`,
    `loadWeekData`
  - [x] `AreaPresenceView.vue` — catch connection error in `loadPresence` and `fetchAreas`
  - [x] `MyBookingsView.vue` — show `{{ bookingsError }}` instead of hardcoded text; catch
    connection error in `fetchMe`
  - [x] `BookingHistoryView.vue` — show `{{ historyError }}` instead of hardcoded text; catch
    connection error in `fetchMe`
  - [x] `ItemGroupBookingsView.vue` — catch connection error in `loadBookings` and `fetchAreas`
  - [x] `LoginView.vue` — show connection error instead of "Invalid email or password"
- [x] Add unit tests to `client.test.ts`
  - [x] `fetch` rejection throws `ApiError` with status 0
  - [x] `isConnectionError` returns true for status 0 ApiError
  - [x] `isConnectionError` returns false for other ApiErrors
  - [x] `isConnectionError` returns false for non-ApiError values
- [x] Run type-check, build, linting, and unit tests

## Dev Notes

### Error Detection Strategy

Network failures (server down, DNS failure, CORS block) cause `fetch()` to throw a `TypeError`
instead of returning a response. The API client now catches any exception from `fetch()` and
wraps it in an `ApiError` with status `0` and the connection lost message. This provides a
consistent error type that views can check with `isConnectionError()`.

### View Update Pattern

Each view that loads data on mount has two potential failure points:

1. `fetchMe()` — called before any data loading to get user info
2. Data loading call — `fetchAreas()`, `loadItems()`, etc.

Both were updated to detect connection errors and set the appropriate error state, preventing
the empty state ("No areas available", etc.) from showing when the real problem is a lost
connection.

### References

- Epic 18 Story 18.7: `_bmad-output/planning-artifacts/epics.md`
- FR66: `_bmad-output/planning-artifacts/prd.md`
- `web/src/api/client.ts`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Updated API client to wrap network failures in `ApiError(CONNECTION_LOST_MESSAGE, 0)`
- Updated 8 views to detect and display connection errors
- Views that showed hardcoded error text now display the actual error message
- Added 4 unit tests for connection error detection
- Added view-level tests covering connection lost rendering across the affected screens
- All 197 frontend tests, type-check, and build pass

### File List

- `web/src/api/client.ts` — Connection error wrapping, `isConnectionError()`, constant
- `web/src/api/client.test.ts` — 4 new tests for connection error handling
- `web/src/views/AreasView.vue` — Dynamic error text, connection error catch in fetchMe
- `web/src/views/ItemGroupsView.vue` — Connection error catch in fetchMe, fetchAreas,
  fetchItemGroups
- `web/src/views/ItemsView.vue` — Connection error catch in fetchMe, fetchAreas, loadItems,
  loadWeekData
- `web/src/views/AreaPresenceView.vue` — Connection error catch in loadPresence, fetchAreas
- `web/src/views/MyBookingsView.vue` — Dynamic error text, connection error catch in fetchMe
- `web/src/views/BookingHistoryView.vue` — Dynamic error text, connection error catch in fetchMe
- `web/src/views/ItemGroupBookingsView.vue` — Connection error catch in loadBookings, fetchAreas
- `web/src/views/LoginView.vue` — Connection error shown instead of misleading auth error
- `web/src/views/AreasView.test.ts` — Verifies connection lost rendering
- `web/src/views/ItemGroupsView.test.ts` — Verifies connection lost rendering
- `web/src/views/ItemsView.test.ts` — Verifies connection lost rendering
- `web/src/views/AreaPresenceView.test.ts` — Verifies connection lost rendering
- `web/src/views/MyBookingsView.test.ts` — Verifies connection lost rendering
- `web/src/views/BookingHistoryView.test.ts` — Verifies connection lost rendering
- `web/src/views/ItemGroupBookingsView.test.ts` — Verifies connection lost rendering
- `web/src/views/LoginView.test.ts` — Verifies connection lost rendering on login failure

## Change Log

- 2026-03-13: Story implemented and verified.
- 2026-03-13: Code review fixes added view-level coverage for all affected connection lost states.
