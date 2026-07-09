# Story 36.3: Named On-Behalf Bookings in My Bookings

Status: done

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

- [x] Task 1: Expose the colleague's full name on the My Bookings API (AC: #1, #3)
  - [x] In `writeBookingsCollection`, added `rec.UserID` (when it differs from the current user) to
        the `userIDSet` so `users.FindDisplayNames` resolves the colleague's display name.
  - [x] Added a new `for_user_name` (`json:"for_user_name,omitempty"`) attribute to
        `MyBookingAttributes`. Populated only for the on-behalf-by-me case
        (`rec.BookedByUserID == currentUserID && rec.UserID != currentUserID`) via a new helper
        `buildMyBookingAttributes` (extracted to keep gocognit within threshold).
  - [x] Left the existing `booked_by_user_name` / `booked_for_me` behaviour unchanged.
- [x] Task 2: Add the field to the frontend booking type (AC: #1, #3)
  - [x] Added `for_user_name?: string;` to `MyBookingAttributes` in `web/src/api/bookings.ts`.
- [x] Task 3: Render the named on-behalf hint in the tile view (AC: #1, #2, #3)
  - [x] In `BookingCard.vue`, added a caption line rendering
        `t('bookings.onBehalfOf', { name: for_user_name })` guarded by
        `booked_by_user_id && !booked_for_me && for_user_name`, with `data-cy="on-behalf-of"`.
  - [x] Added the `bookings.onBehalfOf` key ("On behalf of {name}") to all five locale files;
        the bare `status.onBehalf` chip label is unchanged.
  - [x] Self-bookings show neither the chip (existing `v-else-if` guard) nor the caption.
- [x] Task 4: Render the named on-behalf hint in the table view (AC: #1, #3)
  - [x] Story 36.2 landed in the same change; the table's For/Guest cell renders the same
        "On behalf of {name}" text via `for_user_name` + `bookings.onBehalfOf`, so both views read
        identically.
- [x] Task 5: Tests (AC: #1, #2, #3)
  - [x] Go: added `TestListHandlerIncludesForUserNameForOnBehalfByMe` asserting `for_user_name` is
        present with the colleague's display name for on-behalf-by-me, and absent for both a
        self-booking and a `booked_for_me` booking.
  - [x] Frontend: added `web/src/components/BookingCard.test.ts` (on-behalf shows the name,
        self-booking shows no hint, booked-for-me shows no on-behalf caption) plus a table-view
        assertion in `MyBookingsView.test.ts`.

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

claude-opus-4-8

### Debug Log References

- Go gate: `go test ./...` (incl. new `TestListHandlerIncludesForUserNameForOnBehalfByMe`),
  `go vet ./...`, `gofmt -l`, `golangci-lint run ./...` — all clean.
- golangci-lint initially flagged gocognit 28 on `writeBookingsCollection` after the new branch;
  resolved by extracting `buildMyBookingAttributes`.
- Frontend gate: `npm run type-check`, `npm run lint`, `npx vitest run` (485 tests pass),
  `npm run build` — all clean.

### Completion Notes List

- The API change is the crux: the colleague (`rec.UserID`) was never resolved before, so
  `for_user_name` is a genuinely new attribute, emitted only for on-behalf-by-me bookings via
  `omitempty`. Self-bookings and booked-for-me bookings emit no field.
- Attribute name `for_user_name` matches the existing create-request attribute and JSON:API
  snake_case rules; the app stores a single `display_name` (full name), which satisfies FR168.
- Tile and table views share the `bookings.onBehalfOf` key so wording is identical across views.
- Implemented sequentially after Story 36.2; the table cell wires the same field, so no rework was
  needed.

### File List

- `internal/bookings/handler.go` (modified — `ForUserName` attribute, colleague lookup, extracted
  `buildMyBookingAttributes` helper)
- `internal/bookings/handler_test.go` (modified — `TestListHandlerIncludesForUserNameForOnBehalfByMe`)
- `web/src/api/bookings.ts` (modified — `for_user_name?` on `MyBookingAttributes`)
- `web/src/components/BookingCard.vue` (modified — on-behalf-of caption)
- `web/src/components/BookingCard.test.ts` (new)
- `web/src/views/MyBookingsView.vue` (table For/Guest cell renders `onBehalfOf`; shared with 36.2)
- `web/src/views/MyBookingsView.test.ts` (table on-behalf-name assertion)
- `web/src/locales/{en,de,es,fr,uk}.json` (modified — `bookings.onBehalfOf`)

### Change Log

- 2026-07-09: Added the `for_user_name` API attribute for on-behalf-by-me bookings (FR168) and
  rendered "On behalf of {name}" in both the tile (`BookingCard`) and table (36.2) views, with Go
  handler tests and a `BookingCard` component test.
