# Story 14.2: Streamline Booking Form

Status: done

## Story

As a user,
I want a simplified booking form with fewer options,
So that the interface is less overwhelming and common tasks are faster.

## Acceptance Criteria

1. **Given** I am on the booking page
   **When** I see the booking type options
   **Then** "Book for guest" is not available as an option
   **And** only "Book for myself" and "Book for colleague" are shown

2. **Given** I am on the booking page in day booking mode
   **When** I see the booking options
   **Then** the "Book multiple days" checkbox is not shown
   **And** no additional dates field appears

3. **Given** I select "Book for colleague"
   **When** the colleague fields appear
   **Then** I see a dropdown listing existing users by display name (fetched from
   `GET /api/v1/users`)
   **And** selecting a user from the dropdown sets the booking to use that user's ID
   **And** the free-text colleague email and name fields are removed

## Tasks / Subtasks

- [x] Remove guest radio option (AC: 1)
  - [x] In `ItemsView.vue`: remove the `v-radio` for "Book for guest"
  - [x] In `ItemsView.vue`: remove the guest fields `v-expand-transition` block
  - [x] Remove `guestName`, `guestEmail` refs and all guest-related validation/logic in script
  - [x] Remove `GuestBookingOptions` import and guest param from `bookItem()` and
    `submitWeekBookings()`
- [x] Remove multi-day checkbox and fields (AC: 2)
  - [x] In `ItemsView.vue`: remove multi-day checkbox, additional dates field, and expand
    transition
  - [x] Remove `multiDayBooking`, `additionalDates` refs from script
  - [x] Remove multi-day branch in `bookItem()` that calls `createMultiDayBooking`
  - [x] Remove `createMultiDayBooking` import
- [x] Replace colleague text fields with user dropdown (AC: 3)
  - [x] Create `web/src/api/users.ts` with a `fetchUsers()` function calling
    `GET /api/v1/users`
  - [x] In `ItemsView.vue`: replace the two text fields with a `v-autocomplete` that loads
    users on mount and filters by display name
  - [x] `v-autocomplete` item-title displays `displayName`, item-value is user `id`
  - [x] Replace `colleagueId` / `colleagueName` refs with `selectedColleagueId` ref
  - [x] Update `bookItem()` and `submitWeekBookings()` to use `selectedColleagueId` for the
    `forUserId` param and resolve the name from the users list
- [x] Update unit tests
  - [x] In `ItemsView.test.ts`: added users mock, tests for removed elements and new dropdown
  - [x] Add test for user dropdown rendering
- [x] Verify E2E tests still pass
  - [x] Updated `ui-framework.cy.ts` radio count assertion (3 -> 2) and removed guest label check

## Dev Notes

### Architecture: Frontend + API Client

This story requires changes to `ItemsView.vue` (removing and replacing form elements) and
creation of a new API client file `web/src/api/users.ts`.

### Backend API Already Exists

The `GET /api/v1/users` endpoint was created in Epic 11 (Story 11.7). It returns a JSON:API
collection of users with attributes including `display_name`, `email`, `source`, and `is_admin`.
All authenticated users have read access (not admin-only).

### Colleague Booking API Contract

The `POST /api/v1/bookings` endpoint accepts `for_user_id` in the request body for booking on
behalf. Currently the frontend sends `for_user_id` (email) and `for_user_name`. After this
story, the frontend will send the user's database `id` as `for_user_id` and resolve the
display name from the fetched users list.

### Guest Booking Removal Scope

Removing the guest radio and fields from the UI only. The backend `POST /api/v1/bookings`
endpoint still accepts guest booking parameters — no backend changes needed. The guest booking
feature can be re-added in a future story if needed.

### Multi-Day Removal Scope

The "Book multiple days" checkbox in day mode is removed because week booking mode (Story 13.5)
now provides a better UX for multi-day bookings. The `createMultiDayBooking` API function in
`web/src/api/bookings.ts` is kept for potential future use but no longer referenced by any view.

### References

- Epic 14 Story 14.2: `_bmad-output/planning-artifacts/epics.md` (Epic 14 Stories section)
- FR47, FR48, FR49: `_bmad-output/planning-artifacts/prd.md`
- ItemsView: `web/src/views/ItemsView.vue`
- Bookings API: `web/src/api/bookings.ts`
- Users API (backend): `internal/users/handler.go`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Removed guest radio option and all guest-related template, refs, validation, and logic
- Removed multi-day checkbox, additional dates field, and createMultiDayBooking branch
- Created `web/src/api/users.ts` with `fetchUsers()` calling `GET /api/v1/users`
- Replaced colleague text fields with `v-autocomplete` dropdown using user database IDs
- Users list loaded on mount (non-blocking) with loading state
- `resolveColleagueName()` helper resolves display name from users list for API calls
- Updated ItemsView.test.ts: added users mock, tests for removed guest/multi-day elements
- Updated ui-framework.cy.ts: adjusted radio count and removed guest label assertion
- All 132 unit tests pass, all 51 E2E tests pass
- Type check, ESLint, build, and code duplication checks all pass
- Code review fix: allow on-behalf bookings without empty for_user_name when user list is unavailable
- Code review fix: update colleague dropdown test to actually switch booking type

### Change Log

- 2026-02-14: Implemented Story 14.2 - streamlined booking form by removing guest and
- 2026-02-14: Code review fixes for on-behalf booking payload and colleague dropdown test
  multi-day options, replaced colleague text fields with user dropdown

### File List

- web/src/views/ItemsView.vue (modified - removed guest/multi-day, added user dropdown)
- web/src/api/users.ts (new - fetchUsers API client)
- web/src/views/ItemsView.test.ts (modified - updated mocks and assertions)
- web/cypress/e2e/ui-framework.cy.ts (modified - updated radio button assertions)
- web/src/api/bookings.ts (modified - on-behalf payload omits empty names)