# Story 13.4: Booking Notes

Status: done

## Story

As a user,
I want to add, view, and edit notes on my bookings,
So that I can communicate useful information to colleagues (e.g., "arriving after noon").

## Acceptance Criteria

1. **Given** I have just completed a booking
   **When** the success confirmation message is displayed
   **Then** I see an "add note" action within the confirmation
   **And** clicking it opens a text input where I can type a free-text note
   **And** the note is saved to the booking

2. **Given** I am viewing My Bookings
   **When** a booking has a note
   **Then** the note is displayed as a single line on the booking card
   **And** if the note is longer than the available width, it is truncated
   **And** a truncation indicator (icon) signals there is more text
   **And** clicking the indicator opens a dialog (desktop) or bottom sheet (mobile)
   showing the full note text

3. **Given** I am viewing My Bookings
   **When** I want to edit a note on one of my bookings
   **Then** I can modify the note text and save the changes
   **And** the updated note is reflected immediately

4. **Given** I am viewing Today's Presence for an area
   **When** a booking has a note
   **Then** the note is displayed with the same truncation behavior as My Bookings

5. **Given** I am viewing items in an item group
   **When** a booked item has a note
   **Then** the note is displayed alongside the booker's name with truncation behavior

6. **Given** the backend receives a note update request
   **When** the note is saved via the API
   **Then** the booking record is updated with the new note text
   **And** a `note` field is added to the bookings table via migration
   **And** the JSON:API response includes the note in booking attributes

## Tasks / Subtasks

- [x] Add database migration for `note` column (AC: 6)
  - [x] Create migration `000012_add_booking_note.up.sql`
  - [x] Create migration `000012_add_booking_note.down.sql`
  - [x] Verify migration runs cleanly on existing data
- [x] Add note update API endpoint (AC: 6)
  - [x] `PATCH /api/v1/bookings/:id` with JSON:API payload
  - [x] Authorization: only booking owner, booker, or admin can update
  - [x] Validate note length (max 500 characters)
  - [x] Return updated booking in JSON:API format with `note` in attributes
  - [x] Add handler tests (6 test functions)
- [x] Include `note` in all booking API responses (AC: 6)
  - [x] Add `note` field to all booking attribute structs
  - [x] Update store queries to SELECT the note column
  - [x] Update all booking-related API responses
- [x] Add note to booking creation flow (AC: 1)
  - [x] After successful booking in ItemsView, show "Add note" link
  - [x] Clicking opens a v-dialog with v-textarea
  - [x] Save note via PATCH endpoint
  - [x] Booking ID captured from creation response
- [x] Display note on BookingCard (AC: 2, 3)
  - [x] Show single-line truncated note on the card
  - [x] Truncation indicator icon (mdi-arrow-expand)
  - [x] Click opens v-dialog with full text
  - [x] Edit capability via "Edit note" button
  - [x] Save edits via PATCH endpoint, update local state
- [x] Display note on Today's Presence view (AC: 4)
  - [x] Show truncated note with expand button
- [x] Display note on items view alongside booker name (AC: 5)
  - [x] Show truncated note below booker name on occupied items
- [x] Add frontend API service for note update (AC: 1, 3)
  - [x] Create `updateBookingNote(bookingId, note)` function
- [x] Add Vitest unit tests (AC: 1, 2, 3)
  - [x] Test updateBookingNote API function (2 tests)
- [x] Add Cypress E2E test (AC: 1, 2, 3, 5)
  - [x] Show add note action after booking (AC: 1)
  - [x] Add note via dialog after booking (AC: 1)
  - [x] Display and edit note on BookingCard (AC: 2, 3)
  - [x] Show add note button when no note exists (AC: 3)
  - [x] Show note on occupied items (AC: 5)
- [x] Update OpenAPI documentation (AC: 6)
  - [x] Add `note` to BookingAttributes, MyBookingAttributes, ItemGroupBookingAttributes, AreaPresenceAttributes, ItemAttributes
  - [x] Document PATCH endpoint with UpdateBookingRequest schema
  - [x] Lint with Redocly: valid, 0 errors

## Dev Notes

### Database Migration

SQLite does not support `ALTER TABLE DROP COLUMN` reliably across all versions.
The down migration should recreate the table without the `note` column using the
standard CREATE-INSERT-DROP-RENAME pattern.

### Note Length

Max 500 characters is a reasonable limit for short messages like "arriving after
noon" or "please don't touch the monitor setup". Enforce both in the backend
(validation) and frontend (character counter).

### PATCH Endpoint Design

Use `PATCH /api/v1/bookings/:id` per JSON:API spec. The request body should
follow JSON:API format:

```json
{
  "data": {
    "type": "bookings",
    "id": "abc123",
    "attributes": {
      "note": "Arriving after 2pm"
    }
  }
}
```

### Truncation Behavior

Use CSS `text-overflow: ellipsis` with `overflow: hidden` and `white-space: nowrap`
for single-line truncation. The indicator icon should only appear when the text is
actually truncated (compare scrollWidth vs clientWidth, or use a fixed character
threshold like 60 chars).

### Dependency on Story 13.2

AC 5 (note on items view alongside booker name) depends on Story 13.2 adding the
booker name to the items response. If 13.2 is not complete, AC 5 can be deferred
or the booker name part can be omitted.

### References

- PRD FR37: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 13.4: `_bmad-output/planning-artifacts/epics.md`
- BookingCard: `web/src/components/BookingCard.vue`
- Bookings handler: `internal/bookings/handler.go`
- Bookings store: `internal/bookings/store.go`
- MyBookingsView: `web/src/views/MyBookingsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- All 6 acceptance criteria implemented and verified
- Go tests: all 16 packages pass, 0 lint issues
- Frontend: 106 Vitest tests pass, ESLint clean, 0% TS duplication
- Cypress: 45 E2E tests pass (5 new booking-notes tests)
- OpenAPI: valid with Redocly, PATCH endpoint documented
- Migration number adjusted to 000012 (000011 already existed)

### File List

**Backend (Go)**

- `migrations/000012_add_booking_note.up.sql` (new)
- `migrations/000012_add_booking_note.down.sql` (new)
- `internal/bookings/store.go` (Note field, UpdateNote function, query updates)
- `internal/bookings/handler.go` (PatchHandler, Note fields, response builders)
- `internal/bookings/handler_test.go` (6 new test functions)
- `internal/startup/server.go` (PATCH route)
- `internal/items/handler.go` (note in occupied item attrs)
- `internal/items/handler_test.go` (schema update)
- `internal/areas/presence_handler.go` (note in query, struct, resource)
- `internal/areas/presence_handler_test.go` (schema update)
- `internal/itemgroups/bookings_handler.go` (note in struct, query, resource)
- `internal/itemgroups/bookings_handler_test.go` (schema update)

**Frontend (Vue/TypeScript)**

- `web/src/api/bookings.ts` (note fields, updateBookingNote function)
- `web/src/api/bookings.test.ts` (2 new updateBookingNote tests)
- `web/src/api/areaPresence.ts` (note field)
- `web/src/api/itemGroupBookings.ts` (note field)
- `web/src/api/items.ts` (note field)
- `web/src/components/BookingCard.vue` (note display/edit/truncation UI)
- `web/src/views/MyBookingsView.vue` (note-updated handler)
- `web/src/views/ItemsView.vue` (add note after booking, note on occupied items)
- `web/src/views/AreaPresenceView.vue` (note display with expand dialog)

**Cypress E2E**

- `web/cypress/e2e/booking-notes.cy.ts` (new, 5 tests)

**OpenAPI**

- `api-doc/openapi.yaml` (note in schemas, UpdateBookingRequest)
- `api-doc/endpoints/booking.yaml` (PATCH endpoint)
- `api-doc/endpoints/bookings.yaml` (note in examples)