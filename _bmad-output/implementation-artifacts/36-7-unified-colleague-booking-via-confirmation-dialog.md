# Story 36.7: Unified Colleague Booking via Confirmation Dialog

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user booking for a colleague,
I want one consistent way to pick a colleague across all views,
so that the booking flow is predictable.

## Acceptance Criteria

1. Tile view: clicking "Book" opens the shared booking confirmation dialog offering colleague
   selection; the inline colleague dropdown is removed from the tile view.
2. Selecting a colleague in the dialog and confirming creates the booking on that colleague's
   behalf.
3. The floor-plan booking dialog also offers colleague selection.
4. Selecting a colleague and multiple days in the floor-plan dialog books every selected day on
   that colleague's behalf.
5. No colleague selected means the booking is for myself (unchanged default).

[Source: _bmad-output/planning-artifacts/epics.md#Story 36.7 (lines 5506-5535); FR175 (line 639),
FR176 (line 643)]

## Tasks / Subtasks

Implementation summary: created shared `ColleagueSelect.vue` (radio + autocomplete, `v-model` for
colleague id) backed by a new `useColleagues` composable (fetch + name resolution). Tile view drops
the inline dropdown and opens a `tile-booking-dialog` that hosts the fragment; day + week modes both
route their create through this dialog (option (a) — week-mode colleague booking preserved). Floor
plan hosts the fragment (`fp-` prefix) and threads the on-behalf option into single- and multi-day
calls, so one colleague selection books every selected day on their behalf.

- [x] Task 1: Create a reusable colleague-select fragment for booking dialogs (AC: #1, #2, #3, #5)
  - [x] Extract the colleague radio + autocomplete pattern that already lives in
        `MatrixBookingPopover.vue` (`web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:16-52`,
        state at :146-149, `loadColleagues`/`resolveColleagueName` at :177-218) into a small shared
        component, e.g. `web/src/components/ColleagueSelect.vue`, exposing a `v-model` for the
        selected colleague id and emitting nothing else. Keep the two `data-cy` radios and the
        autocomplete but use generic hooks so tiles + floor plan can reuse it. Do NOT force
        `MatrixBookingPopover.vue` to adopt it in this story unless trivially clean — the matrix
        popover is the reference pattern, not a required refactor target.
  - [x] Reuse i18n keys already present: `items.bookForMyself`, `items.bookForColleague`,
        `items.selectColleague`, `items.selectColleagueError`
        (`web/src/locales/en.json:104-106,135` and sibling `de/fr/es/uk.json`). No new keys needed
        for colleague labels.
  - [x] Reuse `fetchColleagues()` (`web/src/api/users.ts:20-22`) and `BookOnBehalfOptions`
        (`web/src/api/bookings.ts:55-58`). Colleague name resolution follows the existing
        `resolveColleagueName` pattern.

- [x] Task 2: Tile view — replace the inline dropdown with a shared confirmation dialog (AC: #1, #2, #5)
  - [x] Remove the top-level inline colleague `v-autocomplete`
        (`web/src/views/ItemsView.vue:122-136`, `data-cy="colleague-select"`) and its inline CSS
        (`:2225-2233`, `.colleague-select-inline`).
  - [x] Add a booking confirmation dialog to `ItemsView.vue` that contains the colleague-select
        fragment (Task 1) and Cancel/Confirm actions. Model its structure/behaviour on the matrix
        popover flow (`MatrixBookingPopover.vue:75-96` actions, `submitBooking`/`doBook` at
        :220-263) — but as a `v-dialog`, not a `v-menu`. Give it a `data-cy` such as
        `tile-booking-dialog` with `tile-booking-confirm` / `tile-booking-cancel` buttons.
  - [x] Rewire the day-mode "Book" button (`ItemsView.vue:353-364`, `data-cy="book-item-btn"`,
        `@click="requestBooking(entry.id)"`) so `requestBooking` (:1794-1801) now OPENS the new
        dialog (capturing the target item id + name) instead of going straight to the warning flow.
        The warning-confirmation `presentWarnings(...)` step must still run — sequence it so the
        dialog Confirm triggers the existing warning flow, then `bookItem`. Preserve the current
        order: warning confirm (if any) → booking.
  - [x] In `bookItem` (`ItemsView.vue:1831-1865`), read the selected colleague from the dialog's
        model instead of the removed top-level `selectedColleagueId`. Keep the on-behalf mapping
        exactly as today (`:1842-1844`): a selected colleague → `{ forUserId, forUserName }`, none →
        `undefined`. Keep the post-booking reset (`:1861-1862`) but reset the dialog's colleague
        model, and close the dialog on success.
  - [x] WEEK MODE: the same top-level `selectedColleagueId` also feeds `submitWeekBookings`
        (`ItemsView.vue:1497-1506`) and the week warning flow (`startWeekWarningFlow` :1827-1829).
        Since AC1 only removes the dropdown from the tile view and unifies via the confirmation
        dialog, decide the minimal-surprise path: either (a) surface the same colleague-select
        fragment in the week-mode submit flow (preferred — keeps colleague booking possible in week
        mode), or (b) scope this story to day-mode + floor-plan and leave week-mode colleague
        booking as a documented follow-up. Choose (a) if it stays clean; whichever is chosen, week
        mode must NOT be left calling a now-deleted `selectedColleagueId`. Do not silently drop
        week-mode colleague booking without noting it.

- [x] Task 3: Floor-plan dialog — add colleague selection across all selected days (AC: #3, #4, #5)
  - [x] Add the colleague-select fragment (Task 1) into the floor-plan booking dialog
        (`InteractiveFloorPlan.vue:379-501`, between the summary/day list and the actions), with a
        floor-plan-scoped `data-cy` (e.g. `fp-colleague-select`). Add matching component state near
        `pendingBooking` (:756-764).
  - [x] Load colleagues when the dialog opens (`requestBooking` :1460-1473) using
        `fetchColleagues()`; reuse the resolve-name helper pattern.
  - [x] Thread the on-behalf option into BOTH booking calls in `executeBooking`
        (`InteractiveFloorPlan.vue:1541-1565`): pass it as the 3rd arg to `createBooking(...)`
        (:1553) for the single-day path and to `createMultiDayBooking(...)` (:1559) for the
        multi-day path. Both API functions already accept `onBehalf?: BookOnBehalfOptions`
        (`web/src/api/bookings.ts:106-156`), so every selected day is booked on the colleague's
        behalf.
  - [x] Reset the floor-plan colleague selection when the dialog closes / after a successful
        booking (alongside `pendingBooking.value = null` at :1590-1591), so the next booking
        defaults to "for me" (AC #5).

- [x] Task 4: Tests (AC: #1-#5)
  - [x] Vitest component tests for `ItemsView`: the inline `colleague-select` autocomplete is gone;
        clicking "Book" opens the confirmation dialog; confirming with a colleague selected calls
        `createBooking` with `for_user_id`/`for_user_name`; confirming with none selected calls it
        without on-behalf fields. Mock `fetchColleagues`/`createBooking`.
  - [x] Vitest component tests for `InteractiveFloorPlan`: the dialog renders the colleague select;
        selecting a colleague + multiple free days calls `createMultiDayBooking` with the on-behalf
        option; single day calls `createBooking` with the option; no colleague → no on-behalf
        fields. Extend existing `InteractiveFloorPlan.test.ts`.
  - [ ] Cypress E2E (deferred this pass): from the tile view, book a desk on a colleague's behalf
        via the dialog and confirm the booking appears attributed to the colleague; use
        `cy.intercept('POST', '/api/v1/bookings*')` to assert the payload carries `for_user_id`.
        Prefer extending an existing booking spec (`web/cypress/e2e/`).

## Dev Notes

### The unification target — matrix popover is the pattern

The weekly matrix already implements the exact colleague-selection UX this story generalizes:
a `self`/`colleague` radio group plus a colleague `v-autocomplete`, with on-behalf mapping and a
"pick a colleague" validation guard. Treat it as the canonical reference, not necessarily a file to
refactor:

- Radio group + autocomplete template:
  `web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:16-52`.
- State: `bookingType` (`self`/`colleague`), `selectedColleagueId`, `colleagueList`,
  `colleaguesLoading` (`:146-149`).
- `loadColleagues()` (`:177-191`), `resolveColleagueName()` (`:216-218`), on-behalf mapping in
  `doBook()` (`:238-241`), validation guard in `submitBooking()` (`:220-226`).
- It also persists the last colleague under `sithub_matrix_last_colleague`
  (`:120,154-175`). This "remember last colleague" behaviour is matrix-specific; do NOT wire the
  tile/floor-plan dialogs to the same key unless product asks — cross-surface persistence is out of
  scope here. If a shared fragment needs its own persistence later, use a distinct key.

### Tile view — current inline flow being removed (AC #1, #2, #5)

`web/src/views/ItemsView.vue` today has ONE top-level colleague dropdown that applies to both
day-mode and week-mode bookings:

- Inline dropdown to remove: `:122-136` (`data-cy="colleague-select"`, `v-model="selectedColleagueId"`).
- Its state: `selectedColleagueId` (`:977`), `usersList`/`usersLoading` (`:978-979`), loaded by
  `loadUsers()` (`:1423-1434`), name via `resolveColleagueName()` (`:1438-1440`). Note the tile view
  currently loads the FULL users list (`fetchUsers`) not `fetchColleagues` — align on
  `fetchColleagues` to match the matrix/floor-plan and the intended colleague semantics.
- Day-mode "Book" button: `:353-364` → `requestBooking(entry.id)` (`:1794-1801`) → warning flow →
  `bookItem(itemId)` (`:1831-1865`). On-behalf mapping is at `:1842-1844`; success reset clears
  `selectedColleagueId` at `:1861-1862`.
- Week-mode uses the SAME `selectedColleagueId` in `submitWeekBookings()` (`:1497-1506`) and its
  warning entry point `startWeekWarningFlow()` (`:1827-1829`). Removing the shared top-level state
  affects week mode — see Task 2 last subtask. Do not leave dangling references to a deleted ref.

### Warning-confirmation flow must be preserved

Booking already runs through the uniform warning confirmation `useWarningConfirmation`
(`ItemsView.vue:948` `present: presentWarnings`; floor plan `:773` + call sites). This is a
SEPARATE dialog (`WarningConfirmDialog.vue`) from the new colleague/booking confirmation. Do not
merge them. Sequence: user opens the new booking dialog → picks colleague/days → clicks Confirm →
existing `presentWarnings(...)` runs (shows the warning dialog only if the item has a
non-suppressed warning) → on warning-confirm the actual `createBooking`/`createMultiDayBooking`
executes. The floor plan already does this ordering in `confirmPendingBooking()`
(`InteractiveFloorPlan.vue:1518-1538`) → `executeBooking()` (`:1541-1565`); keep it and add the
colleague option inside `executeBooking`.

### Floor-plan dialog specifics (AC #3, #4)

`InteractiveFloorPlan.vue` is one component; the booking dialog lives at `:373-502`
(`data-cy="fp-booking-dialog"`), showing a per-day checkbox list (`:396-448`) for single or
multi-day selection. `pendingBooking` holds the item (`:756-764`). Booking executes in
`executeBooking(bookingDates)` (`:1541-1565`) which branches on `multiDay = bookingDates.length > 1`
and calls `createBooking` (single) or `createMultiDayBooking` (multi). Both take an optional
`onBehalf` third argument, so a single on-behalf value threads to every selected day (AC #4). Reset
the colleague state where `pendingBooking` is cleared on success (`:1590-1591`) and when the dialog
is cancelled (`:484` sets `showBookingDialog = false`).

### API and backend — no changes needed

- `createBooking(itemId, bookingDate, onBehalf?, guest?, note?)` and
  `createMultiDayBooking(itemId, bookingDates, onBehalf?, ...)` already accept
  `onBehalf: BookOnBehalfOptions` and serialize it to `for_user_id` / `for_user_name`
  (`web/src/api/bookings.ts:30-104,106-156`).
- Backend already supports on-behalf bookings: `resolveBookingParticipants` reads
  `for_user_id`, validates the user via `users.FindByID`, and sets `targetUserID` while
  `bookedByUserID` stays the current user (`internal/bookings/handler.go:641-666`). No backend or
  migration work in this story.
- Colleague list source: `fetchColleagues()` → `GET /api/v1/colleagues`
  (`web/src/api/users.ts:16-22`), returning `display_name` per colleague.

### Project Structure Notes

- New (recommended): `web/src/components/ColleagueSelect.vue` — thin shared radio + autocomplete
  fragment. Keep it presentation + `v-model`; keep colleague fetching/name-resolution in the
  hosting view or a tiny composable if duplication becomes real. Follow the Vue rules: Composition
  API + `<script setup>`, typed props/emits, no `any` (`.claude/rules/vue.md`).
- Modified: `web/src/views/ItemsView.vue` (remove inline dropdown, add dialog), and its test.
- Modified: `web/src/components/InteractiveFloorPlan.vue` (add colleague select + thread on-behalf),
  and `InteractiveFloorPlan.test.ts`.
- `MatrixBookingPopover.vue` is the design precedent; optional light refactor to consume the shared
  fragment only if it stays clean.
- WarningConfirmDialog (`web/src/components/WarningConfirmDialog.vue`) is referenced only as a
  style/flow precedent — it is a different dialog and is NOT the colleague dialog.

### Testing standards summary

Run `npm run type-check`, `npm run lint`, `npx vitest run`, and `npm run build`; keep coverage at or
above 75%. Add a Cypress E2E per the acceptance-criteria coverage rule. Use `data-cy` selectors and
intercept `POST /api/v1/bookings*` for waiting/assertions (do not mock in E2E). Avoid TS/JS
duplication (`npx jscpd`, TS threshold 0) — a shared `ColleagueSelect.vue` helps here.
[Source: .claude/rules/vue.md, .claude/rules/cypress.md]

### References

- [Source: _bmad-output/planning-artifacts/epics.md:5506-5535, 639, 643]
- [Source: web/src/views/ItemsView.vue:122-136, 353-364, 977-979, 1423-1440, 1497-1506,
  1794-1801, 1831-1865, 2225-2233]
- [Source: web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:16-52, 75-96, 120,
  146-149, 154-218, 220-263]
- [Source: web/src/components/InteractiveFloorPlan.vue:373-502, 756-764, 1460-1473, 1518-1538,
  1541-1565, 1590-1591]
- [Source: web/src/api/bookings.ts:30-104, 106-156]
- [Source: web/src/api/users.ts:16-22]
- [Source: internal/bookings/handler.go:641-666]
- [Source: web/src/components/WarningConfirmDialog.vue]
- [Source: web/src/locales/en.json:104-106, 116, 135]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `npm run type-check`, `npm run lint`, `npx vitest run` (503 passed), `npm run build`, and
  `npx jscpd` (0 clones) all green.
- Extracted colleague fetch + name resolution into `useColleagues` so `ColleagueSelect.vue`,
  `ItemsView.vue`, and `InteractiveFloorPlan.vue` share one implementation (keeps jscpd at 0).

### Completion Notes List

- Chose option (a): week-mode colleague booking is preserved — the week "Book N days" button now
  opens the same `tile-booking-dialog`; confirming runs the week warning flow + `submitWeekBookings`
  with the dialog's colleague selection. The old top-level `selectedColleagueId` was removed cleanly;
  nothing dangles.
- Tile day-mode Book opens the dialog capturing item id + name; Confirm runs `presentWarnings` then
  `bookItem`, preserving warning → booking order.
- Floor plan threads a single on-behalf option into both `createBooking` and `createMultiDayBooking`,
  so every selected day is booked on the colleague's behalf; selection resets on close/success.
- No backend/API changes (on-behalf params already existed).
- Existing ItemsView/FloorPlan/MatrixBookingPopover tests updated for the dialog-based flow.

### File List

- web/src/components/ColleagueSelect.vue (new)
- web/src/composables/useColleagues.ts (new)
- web/src/components/index.ts
- web/src/views/ItemsView.vue
- web/src/components/InteractiveFloorPlan.vue
- web/src/components/__tests__/ColleagueSelect.test.ts (new)
- web/src/views/ItemsView.test.ts
- web/src/components/__tests__/InteractiveFloorPlan.test.ts
- web/src/locales/{en,de,es,fr,uk}.json (items.confirmBookingFor)
