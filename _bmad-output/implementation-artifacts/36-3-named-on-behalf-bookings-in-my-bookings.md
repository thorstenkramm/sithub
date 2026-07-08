# Story 36.3: Named On-Behalf Bookings in My Bookings

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user who books for colleagues,
I want My Bookings to show for whom I booked,
so that I can tell my on-behalf bookings apart.

## Acceptance Criteria

1. A booking made on behalf of a colleague shows "On behalf of \<first name\> \<last name\>"
   (the colleague's full name) in My Bookings, in **both** the tile view (`BookingCard.vue`) and the
   table view (introduced by story 36.2).
2. A self-booking shows no "on behalf" hint.
3. The colleague full name comes from the booking record (API response) and renders consistently in
   both views.

## Tasks / Subtasks

- [ ] Task 1: Expose the colleague's full name on the My Bookings API (AC: #1, #3)
  - [ ] In `writeBookingsCollection` (`internal/bookings/handler.go:393-459`), the current logic only
        resolves a display name for `rec.BookedByUserID` (the booker). For an on-behalf booking made
        BY the current user, `rec.UserID` is the colleague and `rec.BookedByUserID` is the current
        user — so the colleague's name is never looked up today. Add the colleague's user ID
        (`rec.UserID`, when it differs from the current user) to the `userIDSet` collected at
        `handler.go:398-403` so `users.FindDisplayNames` (`internal/users/store.go:322`) resolves it.
  - [ ] Add a new attribute `for_user_name` (snake_case per JSON:API) to `MyBookingAttributes`
        (`internal/bookings/handler.go:72-87`), e.g. `ForUserName string
        \`json:"for_user_name,omitempty"\``. Populate it in the resource loop
        (`handler.go:436-446`) ONLY for the on-behalf-by-me case: when
        `rec.BookedByUserID == currentUserID && rec.UserID != currentUserID`, set
        `attrs.ForUserName = displayNames[rec.UserID]`. Do not set it for self-bookings or for
        bookings made FOR the current user (`booked_for_me`).
  - [ ] Leave the existing `booked_by_user_name` / `booked_for_me` behaviour unchanged — that path
        (`handler.go:437-446`) serves the reverse case (someone booked FOR me) and is out of scope.
- [ ] Task 2: Add the field to the frontend booking type (AC: #1, #3)
  - [ ] Add `for_user_name?: string;` to `MyBookingAttributes` in `web/src/api/bookings.ts:12-28`,
        keeping the snake_case name aligned with the API attribute.
- [ ] Task 3: Render the named on-behalf hint in the tile view (AC: #1, #2, #3)
  - [ ] In `BookingCard.vue`, the on-behalf chip at `web/src/components/BookingCard.vue:27-32`
        currently renders a bare "On behalf" `StatusChip` with no name. Show the colleague's full
        name. Prefer adding a caption line (mirroring the existing `booked_by`
        block at `BookingCard.vue:39-46`) that renders
        `t('bookings.onBehalfOf', { name: booking.attributes.for_user_name })` when
        `booking.attributes.booked_by_user_id && !booking.attributes.booked_for_me &&
        booking.attributes.for_user_name`. Keep it consistent with the `booked-for-me` caption
        pattern; add a `data-cy` such as `on-behalf-of` for testability.
  - [ ] Add the `bookings.onBehalfOf` key ("On behalf of {name}") to all five locale files
        (`web/src/locales/{en,de,fr,es,uk}.json` — `bookings` section, near `bookedBy` at
        `en.json:183`). Existing bare `status.onBehalf` (`en.json:283`) stays for the chip label.
  - [ ] A self-booking (no `booked_by_user_id`, `booked_for_me = false`) shows neither the chip nor
        the caption — verify the existing `v-else-if` at `BookingCard.vue:27-28` already guards this.
- [ ] Task 4: Render the named on-behalf hint in the table view (AC: #1, #3)
  - [ ] Story 36.2 (`36-2-desktop-table-view-for-my-bookings`, still backlog) introduces the My
        Bookings table with a status / on-behalf column. In that table's on-behalf cell, render the
        same "On behalf of {name}" text driven by `for_user_name`, reusing the
        `bookings.onBehalfOf` key so both views read identically (AC #3).
  - [ ] If 36.2 has not landed when this story is implemented, coordinate sequencing: the backend
        field (Task 1) and shared i18n key (Task 3) are prerequisites the table view will consume;
        the tile view (Task 3) is fully implementable independently.
- [ ] Task 5: Tests (AC: #1, #2, #3)
  - [ ] Go: extend the list/collection handler tests to assert `for_user_name` is present with the
        colleague's display name for an on-behalf-by-me booking, and absent for a self-booking and
        for a `booked_for_me` booking. Reuse `seedTestUser` / `seedTestBookingFull`
        (`internal/bookings/testhelpers_test.go:33-65`); see the existing on-behalf coverage
        `TestCreateHandlerBookOnBehalf` (`internal/bookings/handler_test.go:635`) and the
        `booked_by_user_name` assertions at `handler_test.go:900-908`.
  - [ ] Frontend: component test for `BookingCard.vue` — on-behalf booking (with `for_user_name`)
        shows "On behalf of <name>"; self-booking shows no hint. Follow the existing stub/assertion
        pattern in `web/src/views/MyBookingsView.test.ts:55-82`.

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md` — Story 36.3 (Epic 36 Stories) and FR168
(`epics.md:621-622`, `epics.md:913`). [Source: _bmad-output/planning-artifacts/epics.md#Story 36.3]

### The core data problem (why an API change is required)

A booking row has two user references (`internal/bookings/store.go:40-53`):

- `UserID` — who the booking is **for** (the colleague, for an on-behalf booking).
- `BookedByUserID` — who **made** the booking (the current user, for an on-behalf booking).

The My Bookings list (`ListUserBookingsRange`, `store.go:64-84`) returns rows where the current user
is either `user_id` OR `booked_by_user_id`, so a booking I made for a colleague appears in my list
via `booked_by_user_id`.

Today `writeBookingsCollection` (`internal/bookings/handler.go:393-459`) only looks up display names
for `BookedByUserID` (`handler.go:398-403`, `handler.go:409`) and only sets `booked_by_user_name`
for the reverse case — someone booked FOR me (`handler.go:437-446`). The **colleague's** name
(`rec.UserID`) is never resolved or emitted. Therefore FR168's "On behalf of <first> <last>" cannot
be rendered from the current response — the colleague name must be added to the API. `FindDisplayNames`
(`internal/users/store.go:322`) already returns the `display_name` column (full name), so no new query
type is needed — just include `rec.UserID` in the lookup set and map it to a new `for_user_name`
attribute.

### JSON:API conventions for the new attribute

Attributes use snake_case; the new field is `for_user_name` (mirrors the create-request attribute
`for_user_name` at `handler.go:45` and the FE `BookOnBehalfOptions.forUserName` at
`web/src/api/bookings.ts:55-58`). Use `omitempty` so self-bookings emit no field.
[Source: .claude/rules/json-api.md#Case style for API responses]

### Frontend rendering paths

- Tile view: `web/src/components/BookingCard.vue`. On-behalf chip guard at lines 27-32; the parallel
  "booked for me" caption to mirror is at lines 39-46. The chip itself (`StatusChip.vue:36`,
  `status.onBehalf`) stays as a compact badge; the **name** goes in a caption line so it does not
  crowd the chip. `MyBookingsView.vue` renders `BookingCard` (`web/src/views/MyBookingsView.vue:42-51`).
- Table view: introduced by story 36.2 (`36-2-desktop-table-view-for-my-bookings`, currently
  `backlog` per `_bmad-output/implementation-artifacts/sprint-status.yaml:315`). Its status /
  on-behalf column must use the same `for_user_name` + `bookings.onBehalfOf` key.
- `BookingHistoryView.vue` (`web/src/views/BookingHistoryView.vue`) is a separate list (past
  bookings) and is NOT the "table view" this story refers to — leave it unchanged unless product
  says otherwise.

> [!NOTE]
> Story 36.2 (the table view) is a soft dependency. The backend field and the shared i18n key are the
> integration contract; implement the tile view regardless, and wire the table cell when/if 36.2 is
> present.

### Existing on-behalf create flow (context, no change needed)

On-behalf bookings are created via `for_user_id` / `for_user_name` in the create payload
(`web/src/api/bookings.ts:55-58,84-89`; `internal/bookings/handler.go:44-45,660`). The colleague
selection UI lives in `web/src/views/ItemsView.vue:1497-1506,1840-1846`. This story only changes how
existing on-behalf bookings are **displayed** in My Bookings, not how they are created.

### Project Structure Notes

- Backend modified: `internal/bookings/handler.go` (`MyBookingAttributes` + `writeBookingsCollection`).
  No store or migration change — `UserID` is already selected (`store.go:71-84,99-102`) and
  `FindDisplayNames` already exists.
- Frontend modified: `web/src/api/bookings.ts` (type), `web/src/components/BookingCard.vue` (tile),
  the 36.2 table component (if present), and `web/src/locales/{en,de,fr,es,uk}.json` (new key).
- Aligns with shared-type guidance: no new duplicate response types; the attribute is added to the
  existing `MyBookingAttributes`. [Source: .claude/rules/golang.md#Shared Types]
- No conflicts detected. Variance: FR168 speaks of "first name + last name"; the app stores a single
  `display_name` full name (`internal/users/store.go:346`), so "full name" is satisfied by
  `display_name`. No separate first/last columns exist.

### Testing standards summary

- Go: table-driven tests with `require`/`assert`, in-memory SQLite; use existing helpers
  `seedTestUser` / `seedTestBookingFull` (`internal/bookings/testhelpers_test.go:33-65`).
  Run `golangci-lint run ./...`, `go vet ./...`, `go test ./...`.
  [Source: .claude/rules/golang.md#Testing]
- Frontend: Vitest component tests following `web/src/views/MyBookingsView.test.ts:55-82`. Run
  `npm run type-check`, `npm run lint`, `npx vitest run`, `npm run build`. A Cypress E2E confirming
  the named hint after an on-behalf booking is a nice-to-have. [Source: .claude/rules/vue.md]

### References

- [Source: internal/bookings/handler.go:72-87,393-459,437-446]
- [Source: internal/bookings/store.go:40-53,64-84,99-102]
- [Source: internal/users/store.go:320-346]
- [Source: internal/bookings/testhelpers_test.go:33-65]
- [Source: internal/bookings/handler_test.go:635,900-908]
- [Source: web/src/api/bookings.ts:12-28,55-58]
- [Source: web/src/components/BookingCard.vue:27-46]
- [Source: web/src/components/StatusChip.vue:16,36]
- [Source: web/src/views/MyBookingsView.vue:42-51]
- [Source: web/src/views/MyBookingsView.test.ts:55-82]
- [Source: web/src/locales/en.json:183,283]
- [Source: _bmad-output/planning-artifacts/epics.md:621-622,913]
- [Source: _bmad-output/implementation-artifacts/sprint-status.yaml:315]
- [Source: .claude/rules/json-api.md#Case style for API responses]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
