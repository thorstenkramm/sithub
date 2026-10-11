# Story 37.1: Show Backend Error Messages in the UI

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want to see the exact reason the server rejected my action,
so that I understand what went wrong without contacting support.

## Acceptance Criteria

1. Given a booking request fails with a JSON:API error (e.g. 400 "booking_date is too far in the future
   (maximum 3 weeks in advance): 2026-08-24"), when the error snackbar is shown, then it shows a translated
   prefix (de: "Objekt konnte nicht gebucht werden. API Fehler:") followed on a new line by the untranslated
   backend `errors[].detail`.
2. Given a failed API response carries no `errors[].detail` (e.g. network error, empty body), when the error
   is shown, then the translated prefix is followed by the HTTP status (e.g. "HTTP 502") instead of a generic
   "please try again" text.
3. Given the error paths in tile booking (ItemsView), floor-plan booking (InteractiveFloorPlan), weekly table
   booking (MatrixBookingPopover), cancellation, colleague and guest booking, and profile/settings, when any
   of them fails, then each shows the backend detail as above; specialised messages for known 409 cases may
   stay as the prefix but must not hide the backend detail.
4. Given a user books a reserved item they are not allowed to book (403 with a reservation message), when the
   request fails, then the user stays on the page and sees the backend message; there is no redirect to
   Access Denied.
5. Given an expired session (401), when any request fails with 401, then the existing redirect-to-login
   behavior is unchanged.
6. Given the shared error handling, when implemented, then one helper in the API layer formats prefix +
   detail (no per-view duplication), it is covered by unit tests, and the E2E test books beyond the advance
   limit and asserts the backend detail is visible.

[Source: _bmad-output/planning-artifacts/epics.md#Story 37.1 (lines 5695-5733); FR179 (line 660),
FR180 (line 666)]

## Tasks / Subtasks

- [x] Task 1: Shared error formatting helper in the API layer (AC: #1, #2, #6)
  - [x] Add `formatApiErrorMessage(prefix: string, err: unknown): string` to `web/src/api/client.ts`. It
        returns `` `${prefix}\n${reason}` `` where `reason` is, in order: `err.detail` (an `ApiError` with a
        non-empty detail, verbatim, never translated); `CONNECTION_LOST_MESSAGE` for `status === 0`
        (see Open questions); `` `HTTP ${err.status}` `` for any other `ApiError`; `err.message` for a plain
        `Error`; `String(err)` as last resort. The prefix is passed in already translated so the API layer
        stays free of `vue-i18n`.
  - [x] Add `throwApiError(response: Response): Promise<never>` (or similar) in `client.ts` that runs
        `parseErrorDetail` and throws `ApiError`, and use it in `apiRequest` (`client.ts:59-62`),
        `cancelBooking` (`web/src/api/bookings.ts:193-204`) and both avatar functions
        (`web/src/api/avatars.ts:11-29`, which today throw a plain `Error` without detail). This removes the
        duplicated parse-and-throw block and keeps jscpd green.
  - [x] Optionally let `parseErrorDetail` (`client.ts:23-34`) join multiple `errors[].detail` entries with
        `\n` (today only `errors[0]` is read). Keep the single-error behaviour identical.
  - [x] Vitest in `web/src/api/client.test.ts`: detail present, detail missing (`HTTP 502`), empty/invalid
        body, status 0, plain `Error`, non-Error value, multi-line output (exactly one `\n` between prefix
        and reason), `throwApiError` with and without detail. Extend `bookings.test.ts` / add
        `avatars.test.ts` cases for the new `ApiError` with detail.

- [x] Task 2: i18n prefixes (AC: #1, #2)
  - [x] Add prefix keys in a new `apiErrors` block in all five locales (`web/src/locales/en.json`,
        `de.json`, `fr.json`, `es.json`, `uk.json`). Suggested keys: `book`, `cancel`, `saveNote`,
        `undoBooking`, `checkAvailability`, `uploadAvatar`, `deleteAvatar`, `changePassword`, `load`.
        German `book` MUST be exactly "Objekt konnte nicht gebucht werden. API Fehler:"; English
        "Unable to book item. API error:". Other prefixes follow the pattern "[what failed]. API Fehler:" /
        "[what failed]. API error:".
  - [x] Remove or stop using the "Please try again" variants that the new prefixes replace
        (`items.unableToBook` `en.json:143`/`de.json:143`, `items.unableToCancelBooking` :153,
        `items.unableToSaveNote` :155, `bookings.unableToCancel` :189, `app.avatarDialog.uploadFailed` :48,
        `floorPlan.bookingCouldNotBeCompleted` :285, `matrix.cancelFailed` :320, ...). Delete keys that
        end up unused in all locales; do not leave orphans.

- [x] Task 3: Split 401 from 403 handling for user actions (AC: #4, #5)
  - [x] `useAuthErrorHandler` (`web/src/composables/useAuthErrorHandler.ts:11-21`) redirects 401 to
        `/login` AND 403 to `/access-denied`. Add an option (e.g. `handleAuthError(err, { forbidden:
        'report' })`) or a sibling `handleUnauthorized(err)` that only handles 401. Page-load callers keep
        the 403 redirect; booking/cancel/note action callers in `ItemsView.vue` (1430, 2128, 2170, 2200)
        use the 401-only variant so a 403 falls through to the snackbar.
  - [x] `ItemGroupsView.vue:438-447` has a private copy of the same handler; it is only used for page
        loads (491, 536), so leave its behaviour unchanged (or replace it with the composable if trivial).
  - [x] `InteractiveFloorPlan.vue` and `MatrixBookingPopover.vue` do not call any auth handler today, so a
        401 there currently just shows the generic text. Do not add new behaviour beyond the AC: keep 401
        handling as is (see Open questions).

- [x] Task 4: Tile view (ItemsView) error paths (AC: #1, #3, #4)
  - [x] `bookItem` catch (`web/src/views/ItemsView.vue:2127-2152`): replace `t('items.unableToBook')`
        fallback and the 404 `itemNotFound` text with `formatApiErrorMessage(t('apiErrors.book'), err)`.
        Keep the `${itemName} - ${date}:` context prefix. For 409 keep `localizeItemsBookingConflict` as the
        prefix and append the backend detail via the helper. Limit dialog (`:2135-2140`,
        `limitDialogMessage`, already `white-space: pre-line` at :903) gets the same treatment.
  - [x] Week mode `submitWeekBookings` per-day errors (`:1674-1684`): `localizeItemsBookingError`
        (`:1541-1553`) returns only localized text; make it return prefix + backend detail via the helper.
        Check how the per-day messages are rendered so that the `\n` is visible.
  - [x] Area-day guard `.catch` (`:2038-2040`) -> helper with `apiErrors.book`.
  - [x] Week cancel (`:1429-1432`), admin cancel (`:2169-2177`), note save (`:2199-2204`) -> helper with
        `apiErrors.cancel` / `apiErrors.saveNote`. The 404 "booking not found" texts may stay as prefix,
        but the detail must follow.
  - [x] Swap-cancel failures (`:1705-1710`, `:2105-2106`, `areaDayGuardCancelFailed`) currently swallow
        the error in a bare `catch {}`; capture the error and append its detail.
  - [x] Error snackbar (`:889-898`, `data-cy="booking-error-text"`): add `white-space: pre-line` (CSS
        class, not inline style if a scoped block exists) so the second line renders.

- [x] Task 5: Floor plan (InteractiveFloorPlan) error paths (AC: #1, #3, #4)
  - [x] `formatBookingError` (`web/src/components/InteractiveFloorPlan.vue:1582-1590`) -> 409 keeps
        `floorPlan.itemAlreadyBooked` / `selectedDaysAlreadyBooked` as prefix, all other cases use
        `apiErrors.book`; both append backend detail via the helper. Caller at `:1769-1771`.
  - [x] Availability check failure (`:1492-1504`), swap-cancel failure (`:1669-1674`), undo failure
        (`:1786-1793`) -> capture `err` and use the helper.
  - [x] Multi-day partial conflicts (`:1729-1733`, `formatConflictMessage`) are a 200 response with
        `conflicts[]`, not an error; leave unchanged.
  - [x] Error snackbar (`:530-539`, `data-cy="fp-booking-error"`): add `white-space: pre-line`, align with
        `feedback.md` (timeout `6000`, `closable`; today 4000 and not closable).

- [x] Task 6: Weekly table (matrix) error paths (AC: #1, #3, #4)
  - [x] `MatrixBookingPopover.vue`: guard `.catch` (`:288-290`) and `doBook` catch (`:327-334`) -> helper.
        Keep `matrix.deskNoLongerAvailable` as the 409 prefix and still emit `bookingConflict`. Swap-cancel
        catch (`:314-318`) -> append detail.
  - [x] `MatrixCancelPopover.vue:98-99` -> helper with `apiErrors.cancel`.
  - [x] Render the inline `v-alert` (`MatrixBookingPopover.vue:65-73`, `data-cy="matrix-booking-error"`)
        with `white-space: pre-line`.

- [x] Task 7: My Bookings, note editing, profile/settings (AC: #1, #3)
  - [x] `MyBookingsView.vue:276-288` cancel -> helper (keep `notFoundOrCancelled` as 404 prefix).
  - [x] `BookingCard.vue:265-273` `saveNote` has no `catch` at all (unhandled rejection, silent failure).
        Add a catch that surfaces `formatApiErrorMessage(t('apiErrors.saveNote'), err)` (emit to parent or
        show a snackbar inside the card, following `feedback.md`).
  - [x] `App.vue` avatar upload/delete (`:516-518`, `:532-534`) -> helper (requires Task 1 avatar change).
  - [x] `App.vue` `getPasswordErrorMessage` (`:569-587`): keep the localized specialised messages as
        prefix but append the backend detail; the generic fallback uses `apiErrors.changePassword`. A 401
        here means "wrong current password", NOT session expiry; do not redirect.
  - [x] Pre-line rendering for the dialog `v-alert`s (`App.vue:302-309`, `:373-380`).

- [x] Task 8: Page-level load errors (AC: #2, FR179 "every failed API call")
  - [x] Generic load fallbacks (`items.unableToLoad` `ItemsView.vue:1807,1877`, `unableToLoadWeekly`
        :1504,1951, `itemGroups.unableToLoad` `ItemGroupsView.vue:547`, `AreaWeeklyMatrixView.vue:264`,
        `presence.unableToLoad` `AreaPresenceView.vue:205`, `itemGroupBookings.unableToLoad`
        `ItemGroupBookingsView.vue:150`) -> helper with `apiErrors.load`. Keep 404 "not found" and
        connection-lost handling as they are; keep these as persistent `v-alert`s (`feedback.md`).
  - [x] `useApi.ts:13` stores `err.message` ("Request failed: 500"); store the detail instead
        (detail when present, otherwise "HTTP" plus status via the same fallback as the helper) so `MyBookingsView` /
        `BookingHistoryView` error alerts show it. Update `useApi.test.ts`.

- [x] Task 9: Tests (AC: #1-#6)
  - [x] Vitest: helper (Task 1). Component tests that mock the API to reject with
        `new ApiError('Request failed: 400', 400, '<detail>')` and assert the rendered text contains the
        prefix and the detail for: `ItemsView` day booking, week booking, cancel; `InteractiveFloorPlan`
        booking; `MatrixBookingPopover` booking (incl. 409 with detail); `MatrixCancelPopover`;
        `MyBookingsView` cancel; `BookingCard` note save; `App.vue` password + avatar.
  - [x] Vitest: 403 from `createBooking` in `ItemsView` does NOT call `router.push('/access-denied')` and
        shows the reservation message; 401 still sets `window.location.href = '/login'`; page-load 403
        still redirects to `/access-denied`.
  - [x] Update existing tests that assert old strings (`unableToBook`, `bookingCouldNotBeCompleted`,
        `deskNoLongerAvailable`, `cancelFailed`, ...) in `ItemsView.test.ts`,
        `components/__tests__/InteractiveFloorPlan.test.ts`, matrix tests, `MyBookingsView.test.ts`.
  - [x] Cypress E2E `web/cypress/e2e/booking-errors.cy.ts`: `cy.login()`; in `onBeforeLoad` set
        `sessionStorage.sithub_selected_day` to a date beyond the advance limit (read
        `weeks_in_advanced` via `cy.request('/api/v1/settings')` and compute a weekday after
        `nextMonday + weeks*7`; `useDateState.loadDay` only rejects past days, so the view opens on it).
        Visit an item group in day mode, `cy.intercept('POST', '/api/v1/bookings*').as('createBooking')`,
        click `book-item-btn`, confirm `tile-booking-confirm`, `cy.wait('@createBooking')` and assert
        status 400; assert `booking-error-text` contains the English prefix and
        "booking_date is too far in the future". No response mocking (`cypress.md`).
  - [x] If the dev area YAML used by E2E contains a `reserved_for` item the test user is excluded from,
        add a second E2E asserting the 403 message and `cy.location('pathname')` unchanged; otherwise cover
        AC #4 in Vitest only and note this in Completion Notes.
  - [x] Go: no backend change is expected. Run `go test ./...` and `golangci-lint run ./...` anyway; the
        detail texts are already asserted in `internal/bookings/handler_test.go:1585` (advance limit) and
        `:1640` (reservation).

- [x] Task 10: Quality gates
  - [x] `npm run type-check`, `npm run lint`, `npx vitest run --coverage` (>= 75 %), `npm run build`,
        `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0 --exit-code 1`,
        `npm run test:e2e -- --browser electron` against `npm run dev` + backend on :9900.

### Review Findings

- [x] [Review][Patch][P2] Preserve valid backend details in mixed malformed error arrays and ignore non-string details [web/src/api/client.ts:16] — fixed with runtime validation and regression coverage.
- [x] [Review][Patch][P3] Report swapped-booking cancellation failures even when the rejected value is falsy [web/src/views/ItemsView.vue:1710; web/src/components/InteractiveFloorPlan.vue:1720] — fixed with a separate failure flag; regression cases cover null, undefined and empty string.

## Dev Notes

### Current state of error handling

The data needed is already there; the views throw it away.

- `ApiError` (`web/src/api/client.ts:3-12`) carries `status` and `detail`. `parseErrorDetail`
  (`:23-34`) reads `errors[0].detail`. `apiRequest` throws `ApiError('Request failed: <status>', status,
  detail)` (`:59-62`); network failures throw `ApiError(CONNECTION_LOST_MESSAGE, 0)` (`:55-57`).
- `cancelBooking` (`web/src/api/bookings.ts:193-204`) duplicates the parse-and-throw block because it uses
  raw `fetch` (DELETE returns 204 without a body). `uploadAvatar`/`deleteAvatar`
  (`web/src/api/avatars.ts:11-29`) throw plain `Error` and lose the detail entirely.
- Views map errors to fixed localized strings. The worst offender for this story is
  `ItemsView.vue:2133` (`items.unableToBook` = "Objekt konnte nicht gebucht werden. Bitte versuchen Sie es
  erneut.", `de.json:143`), which hides e.g. the 400 advance-limit message produced at
  `internal/bookings/handler.go:781-785`.
- 409 mapping by substring: `localizeItemsBookingConflict` (`ItemsView.vue:1521-1539`),
  `isLimitError` (`:1555-1558`). These are allowed to stay as the prefix (AC #3). Do not remove the limit
  dialog flow.

### Reserved-item 403 (AC #4)

Backend returns `403` with detail `This item "<name>" is reserved. You do not have access.`
(`internal/bookings/handler.go:555-564`, message built in `reservationForbiddenMessage` :586-599). The
frontend `bookItem` catch calls `handleAuthError` first (`ItemsView.vue:2128`), and
`useAuthErrorHandler.ts:16-19` pushes `/access-denied` for any 403. That redirect is correct for page loads
(area not accessible) but wrong for booking actions. Split the handling (Task 3); do not change the
backend status code.

### Helper design

Keep it a pure function in `client.ts` with the translated prefix as a parameter, so it is trivially unit
testable and does not need an i18n instance. Each call site becomes one line:

```ts
errorSnackbarMessage.value = formatApiErrorMessage(t('apiErrors.book'), err);
```

For specialised 409 messages, pass the specialised text as the prefix:
`formatApiErrorMessage(t('matrix.deskNoLongerAvailable'), err)`. The backend detail is never translated,
never lower-cased and never truncated.

### Rendering the line break

Vuetify snackbars and alerts collapse `\n`. Add `white-space: pre-line` to the text element of every
surface that shows these messages (ItemsView snackbar, floor-plan snackbar, matrix popover alert, App.vue
dialog alerts, My Bookings alert). `ItemsView.vue:903` already uses this for the limit dialog.

### Feedback conventions

Per `.claude/rules/feedback.md` transient operation errors belong in `v-snackbar` (`color="error"`,
`timeout` 6000, `closable`, `location="bottom"`, `data-cy`). Existing in-dialog/in-popover `v-alert`s
(matrix popover, password/avatar dialogs, My Bookings cancel) are kept as they are in this story; only their
text changes (see Open questions). Fix the floor-plan snackbar timeout/closable while touching it.

### Out of scope

- No backend or API changes, no new endpoints, no OpenAPI changes.
- Guest booking: the backend supports `is_guest`/`guest_name`/`guest_email` (`createBooking` has a `guest`
  parameter), but no current UI calls it with a guest (`grep createBooking web/src` shows only on-behalf
  usage). Guest errors are covered automatically because they go through the same `createBooking` call
  sites; no separate UI work.
- Login errors (`LoginView.vue:128-160`) have their own specialised mapping; leave unchanged unless the
  generic fallback is trivially switched to the helper.

### Project Structure Notes

- Modified: `web/src/api/client.ts`, `web/src/api/bookings.ts`, `web/src/api/avatars.ts`,
  `web/src/composables/useAuthErrorHandler.ts`, `web/src/composables/useApi.ts`,
  `web/src/views/ItemsView.vue`, `web/src/components/InteractiveFloorPlan.vue`,
  `web/src/components/area-weekly-matrix/MatrixBookingPopover.vue`,
  `web/src/components/area-weekly-matrix/MatrixCancelPopover.vue`, `web/src/views/MyBookingsView.vue`,
  `web/src/components/BookingCard.vue`, `web/src/App.vue`, load-error views listed in Task 8, and all
  five locale files.
- New: `web/cypress/e2e/booking-errors.cy.ts` (naming per `cypress.md`: `feature-action.cy.ts`).
- Vue rules: Composition API, `<script setup>`, no `any` (use `unknown` + `instanceof ApiError` narrowing),
  i18n for every user-visible prefix (`.claude/rules/vue.md`).
- No Go code touched, so no Go jscpd/migration impact.

### Testing standards summary

Vitest + Vue Test Utils for helper and component behaviour (coverage >= 75 %); one Cypress E2E per AC where
the dev dataset allows it, using `data-cy` selectors and `cy.intercept` only to wait/assert, never to mock
(`.claude/rules/cypress.md`). Test names follow "should [user outcome]". Run the full quality gate in
Task 10.

### Open questions / assumptions

- Network error (status 0): AC #2 says "followed by the HTTP status", but a network failure has no status.
  Default: show `CONNECTION_LOST_MESSAGE` ("Connection to server lost") as the second line instead of
  "HTTP 0".
- Multiple `errors[]` entries: default is to join all `detail` values with a newline; single-error output
  is unchanged.
- Existing inline `v-alert`s inside dialogs/popovers stay (only text changes); migrating them to snackbars
  is not required by the ACs. MyBookings cancel uses a `v-alert` that violates `feedback.md`; converting it
  to a snackbar is recommended if it stays a small change.
- 401 inside floor plan and matrix popover: today no redirect happens there. AC #5 says "unchanged", so no
  redirect is added; the 401 message is simply shown with the prefix.
- Page-level load errors (Task 8) are included because FR179 says "every failed API call", although the
  AC list names only action paths. If time-boxed, Tasks 1-7 + 9 are the must-haves.
- Reserved-item E2E depends on the dev areas YAML; if it has no suitable `reserved_for` item, AC #4 is
  covered by Vitest only.

### References

- [Source: _bmad-output/planning-artifacts/epics.md:5695-5733, 660-670]
- [Source: web/src/api/client.ts:3-65]
- [Source: web/src/api/bookings.ts:193-204]
- [Source: web/src/api/avatars.ts:11-29]
- [Source: web/src/composables/useAuthErrorHandler.ts:11-21]
- [Source: web/src/composables/useApi.ts:7-18]
- [Source: web/src/composables/useDateState.ts:60-66]
- [Source: web/src/views/ItemsView.vue:889-903, 1002-1011, 1429-1432, 1521-1558, 1674-1710, 2038-2040,
  2105-2106, 2127-2204]
- [Source: web/src/components/InteractiveFloorPlan.vue:530-539, 1492-1504, 1582-1590, 1669-1674,
  1729-1733, 1769-1771, 1786-1793]
- [Source: web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:65-73, 288-290, 314-334]
- [Source: web/src/components/area-weekly-matrix/MatrixCancelPopover.vue:95-99]
- [Source: web/src/views/MyBookingsView.vue:276-288]
- [Source: web/src/components/BookingCard.vue:265-273]
- [Source: web/src/App.vue:302-309, 373-380, 516-534, 569-600]
- [Source: web/src/views/ItemGroupsView.vue:438-447]
- [Source: web/src/locales/en.json:48, 143-155, 189, 285-288, 317-320 and sibling de/fr/es/uk]
- [Source: internal/bookings/handler.go:555-564, 586-599, 781-785]
- [Source: internal/bookings/handler_test.go:1585, 1640]
- [Source: .claude/rules/feedback.md, vue.md, cypress.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 5.5 (claude-opus-5-5)

### Debug Log References

- RED confirmed for `formatApiErrorMessage`/`throwApiError`/multi-detail/avatar `ApiError` (13 failing),
  `handleUnauthorized` (2 failing) and `useApi` detail storage (2 failing) before implementation.
- `App.test.ts`: a route assertion was dropped because the router guard redirects to `/login` in the test
  setup regardless of the password error; the inline alert text proves no redirect path was taken.

### Completion Notes List

- Helper: `formatApiErrorMessage(prefix, err)` plus `apiErrorReason(err)` in `web/src/api/client.ts`. Reason
  order: backend detail (verbatim), `CONNECTION_LOST_MESSAGE` for status 0, `HTTP <status>`, `Error.message`,
  `String(err)`. `parseErrorDetail` joins multiple `errors[].detail` with `\n`. `throwApiError(response)`
  replaces the duplicated parse-and-throw in `apiRequest`, `cancelBooking` and both avatar calls (avatars now
  throw `ApiError` with detail).
- i18n: new `apiErrors` block (`book`, `cancel`, `saveNote`, `undoBooking`, `checkAvailability`,
  `uploadAvatar`, `deleteAvatar`, `changePassword`, `load`) in all five locales; de `book` is exactly
  "Objekt konnte nicht gebucht werden. API Fehler:". 18 orphaned "Please try again"/generic keys removed from
  all locales.
- 401/403: `useAuthErrorHandler` gained `handleUnauthorized` (401 only). Booking, week cancel, admin cancel,
  note save (ItemsView) and My Bookings cancel use it, so a 403 reservation message stays on the page.
  Page loads keep `handleAuthError`. `ItemGroupsView` private handler left unchanged (page loads only).
- Specialised 409/404 texts stay as prefix (`localizeItemsBookingConflict`, `itemNotFound`,
  `bookingNotFound`, `notFoundOrCancelled`, `floorPlan.itemAlreadyBooked`, `matrix.deskNoLongerAvailable`,
  password specialised texts); the backend detail always follows on a new line. Limit dialog included.
- Swap-cancel failures (ItemsView day/week, floor plan, matrix) now capture the error and append its detail.
- Line break rendering: global utility class `.api-error-text { white-space: pre-line }` in
  `styles/global.css`, applied to every snackbar/alert showing these messages and to the week results row.
- Floor-plan error snackbar aligned with `feedback.md` (6000 ms, closable). My Bookings cancel error moved from
  a page `v-alert` to an error snackbar (`data-cy="cancel-error"` kept). `BookingCard.saveNote` now catches
  errors and shows them in a card-local error snackbar (`data-cy="note-save-error"`).
- Page-load fallbacks (Task 8) use `apiErrors.load`; 404 and connection-lost branches unchanged. `useApi`
  stores `apiErrorReason(err)`.
- Open-question defaults applied: status 0 shows "Connection to server lost"; multiple details joined by
  newline; dialog/popover `v-alert`s kept; no 401 redirect added in floor plan/matrix; undo failure uses
  `apiErrors.undoBooking` (replaces `floorPlan.undoFailedStillActive`).
- AC #4 E2E: the E2E/CI areas config (`test_room`, desks 1-2) has no `reserved_for` item, so AC #4 is covered
  by Vitest only (`ItemsView.test.ts` 403 test, `useAuthErrorHandler.test.ts`, `MyBookingsView.test.ts`).
- Tests: Vitest 57 files / 565 tests pass; coverage statements 81.21 %, lines 82.67 % (branches 73.53 %, no
  threshold configured). Cypress full suite 70/70 passed against a CI-equivalent backend on :9900 incl. new
  `booking-errors.cy.ts` (real 400 from the advance limit, no mocking). type-check, ESLint, jscpd (0 clones),
  build, `go test ./...` and `golangci-lint` (0 issues) pass. No Go code changed.

### File List

- `web/src/api/client.ts`
- `web/src/api/client.test.ts`
- `web/src/api/bookings.ts`
- `web/src/api/avatars.ts`
- `web/src/api/avatars.test.ts`
- `web/src/composables/useAuthErrorHandler.ts`
- `web/src/composables/__tests__/useAuthErrorHandler.test.ts`  (new)
- `web/src/composables/useApi.ts`
- `web/src/composables/useApi.test.ts`
- `web/src/styles/global.css`
- `web/src/views/ItemsView.vue`
- `web/src/views/ItemsView.test.ts`
- `web/src/views/MyBookingsView.vue`
- `web/src/views/MyBookingsView.test.ts`
- `web/src/views/AreaPresenceView.vue`
- `web/src/views/ItemGroupBookingsView.vue`
- `web/src/views/ItemGroupsView.vue`
- `web/src/views/BookingHistoryView.vue`
- `web/src/components/InteractiveFloorPlan.vue`
- `web/src/components/__tests__/InteractiveFloorPlan.test.ts`
- `web/src/components/area-weekly-matrix/MatrixBookingPopover.vue`
- `web/src/components/area-weekly-matrix/MatrixBookingPopover.test.ts`
- `web/src/components/area-weekly-matrix/MatrixCancelPopover.vue`
- `web/src/components/area-weekly-matrix/MatrixCancelPopover.test.ts`
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue`
- `web/src/components/BookingCard.vue`
- `web/src/components/BookingCard.test.ts`
- `web/src/App.vue`
- `web/src/App.test.ts`
- `web/src/locales/en.json`
- `web/src/locales/de.json`
- `web/src/locales/fr.json`
- `web/src/locales/es.json`
- `web/src/locales/uk.json`
- `web/cypress/e2e/booking-errors.cy.ts`  (new)
- `_bmad-output/implementation-artifacts/sprint-status.yaml`

## Change Log

- 2026-10-08: Implemented story 37.1. Shared API error formatting helper, `apiErrors` i18n prefixes, 401/403
  split for user actions, backend detail shown on all booking/cancel/note/profile/load error paths, unit and
  E2E tests. Status -> review.

## Code Review Record

- Date: 2026-10-08. Parallel Blind Hunter, Edge Case Hunter and Acceptance Auditor reviews; actionable findings auto-fixed as requested.
- Validation: 600 Vitest tests pass in 58 files; statements 81.56%, lines 83.01%. All 12 story Cypress tests pass against the Vite dev server and an isolated real backend; week/mobile specs rerun after the final loading fix (11/11).
- Type-check, ESLint, production build, TypeScript duplication check (0 clones), Go tests and golangci-lint (0 issues) pass. No unresolved review findings.
- 2026-10-08: Code review complete; status -> done.
