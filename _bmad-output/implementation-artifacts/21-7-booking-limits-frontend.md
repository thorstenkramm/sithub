# Story 21.7: Booking Limits — Frontend Integration and Error Display

Status: done

## Story

As a user,
I want to see only the weeks I am allowed to book and receive clear error messages
when I exceed booking limits,
so that I understand the constraints and can plan accordingly.

## Acceptance Criteria

1. **Given** the backend settings return `weeks_in_advanced: 5`,
   **when** the week selector renders,
   **then** it shows only the current week plus 5 additional weeks (6 total options).

2. **Given** the backend settings return `weeks_in_advanced: 5`,
   **when** I use the day-mode date picker,
   **then** dates beyond the allowed booking window are not selectable.

3. **Given** I try to book an item and the backend returns a booking limit error,
   **when** the error message displays,
   **then** the message appears as a **red snackbar** (same position and shape as the
   green success snackbar) — not as an inline v-alert banner. The message is
   human-readable and includes the limit number and scope
   (e.g. "you have reached the maximum of 2 active bookings for 'Room 1, Desk 1'").

4. **Given** the settings endpoint fails or is unavailable,
   **when** the page loads,
   **then** the week selector falls back to a default of 8 weeks and no error is shown.

5. **Given** any booking operation fails (limit exceeded, conflict, cancellation error,
   note save error, or validation error),
   **when** the error is shown in `ItemsView`,
   **then** it uses a **red error snackbar** (`color="error"`) at the bottom of the
   viewport — the same `v-snackbar` component and position as the green success snackbar.
   The old inline `v-alert` banner (`data-cy="booking-error"`) must be removed entirely.

6. **Given** success and error feedback are displayed,
   **when** a user compares them,
   **then** both use the same `v-snackbar` component at the same viewport position,
   differing only in color (`success` vs `error`) and the `data-cy` attribute.

## Tasks / Subtasks

- [x] Task 1: Create settings API function (AC: 1, 4)
  - [x] 1.1 Create `web/src/api/settings.ts` with `fetchSettings()` and
    `SettingsAttributes` interface
- [x] Task 2: Update week selector to accept dynamic max weeks (AC: 1, 4)
  - [x] 2.1 Add optional `maxWeeks` Ref parameter to `useWeekSelector()`
  - [x] 2.2 Use `(maxWeeks?.value ?? 7) + 1` for the loop count (current + N weeks)
- [x] Task 3: Fetch settings and wire to week selector (AC: 1, 4)
  - [x] 3.1 In `ItemGroupsView.vue`: add `weeksInAdvanced` ref, fetch settings on mount,
    pass to `useWeekSelector`
  - [x] 3.2 In `ItemsView.vue`: add `weeksInAdvanced` ref, fetch settings on mount,
    pass to `useWeekSelector`
- [x] Task 4: Limit day-mode date picker max date (AC: 2)
  - [x] 4.1 Add `maxBookingDate` computed property in `ItemsView.vue` based on
    `weeksInAdvanced`
  - [x] 4.2 Bind `:max="maxBookingDate"` to `DatePickerField`
- [x] Task 5: Display booking limit error messages (AC: 3)
  - [x] 5.1 Update `localizeItemsBookingConflict()` in `ItemsView.vue` to detect
    "booking limit exceeded" in the error detail and pass through the user-facing message
- [x] Task 6: Replace inline v-alert errors with error snackbar (AC: 5, 6)
  - [x] 6.1 In `ItemsView.vue`: remove the `v-alert` block at `data-cy="booking-error"`
  - [x] 6.2 Add a second `v-snackbar` for errors: `color="error"`, `location="bottom"`,
    `data-cy="booking-error"`, bound to a new `showErrorSnackbar` computed
  - [x] 6.3 Replaced `bookingErrorMessage` and `bookingErrorDetails` refs with a single
    `errorSnackbarMessage` ref. Detailed errors formatted as
    `${itemName} - ${date}: ${detail}`, simple errors as plain string
  - [x] 6.4 Error snackbar `timeout` set to 6000ms with `closable`
  - [x] 6.5 Removed `clearErrorMessage()` function — snackbar auto-dismisses
  - [x] 6.6 Kept `data-cy="booking-error"` and `data-cy="booking-error-text"` selectors
- [x] Task 7: Update tests for snackbar errors (AC: 5, 6)
  - [x] 7.1 Verified: no Vitest tests in `ItemsView.test.ts` reference the old
    `v-alert` / `booking-error` selector — no changes needed
  - [x] 7.2 Verified: Cypress E2E tests use `[data-cy="booking-error-text"]` which is
    preserved on the snackbar span — no changes needed
- [x] Task 8 (was Task 6): Write tests (AC: 1, 3, 4)
  - [x] 6.1 Add `settings.test.ts` for `fetchSettings()`
  - [x] 6.2 Add week selector tests for `maxWeeks` parameter in
    `useWeekSelector.test.ts`
  - [x] 6.3 Run full Vitest suite — 255 tests pass
  - [x] 6.4 Run ESLint — clean
  - [x] 6.5 Run type-check — clean
  - [x] 6.6 Run build — clean

## Dev Notes

### Settings Fetch Strategy

Settings are fetched non-blocking in `onMounted` after auth succeeds. On failure, the
week selector uses a default of 7 additional weeks (8 total). This ensures the UI
always works even if the settings endpoint is temporarily unavailable.

### Error Message Passthrough

The backend returns booking limit errors as 409 responses with detail strings like:
`"booking limit exceeded: you have reached the maximum of 2 active bookings for 'Room 1, Desk 1'"`

The frontend detects the "booking limit exceeded" prefix and extracts the user-facing
message after the colon. This avoids needing frontend translation keys for every possible
limit combination — the backend already generates clear, contextual messages.

### Feedback Consistency (UX Review Fix)

The original implementation used an inline `v-alert` banner for booking errors and a
`v-snackbar` for success confirmations. This created an inconsistent experience:
users had to look in different places for feedback depending on whether something
succeeded or failed.

**Convention (see `.claude/rules/feedback.md`):** All user-facing feedback in SitHub
must use `v-snackbar` at the bottom of the viewport. Success = green, error = red.
No inline `v-alert` banners for transient operation feedback.

The error snackbar uses a longer timeout (6000ms vs 3000ms for success) because error
messages tend to be longer and users need more time to read and understand them.

### Day-Mode Max Date

The max date for the day picker is calculated as: next Monday + `weeksInAdvanced * 7 - 1`
days. This matches the backend's horizon calculation.

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Created `api/settings.ts` with `fetchSettings()` and types
- Updated `useWeekSelector` to accept optional `maxWeeks` ref
- Both `ItemGroupsView` and `ItemsView` fetch settings on mount and wire to week selector
- Day-mode date picker has computed `maxBookingDate` bound to `:max`
- `localizeItemsBookingConflict` passes through booking limit error messages from backend
- 4 new frontend tests (settings API, week selector maxWeeks, reactivity)
- **Rework (UX review):** Replaced inline `v-alert` error banner with `v-snackbar`
  (`color="error"`, `timeout=6000`, `closable`, `location="bottom"`)
- Consolidated `bookingErrorMessage` + `bookingErrorDetails` into single
  `errorSnackbarMessage` ref with `showErrorSnackbar` computed
- Removed `clearErrorMessage()` — snackbar auto-dismisses
- No test changes needed — `data-cy` selectors preserved
- All 255 tests pass, ESLint/type-check/build clean

### File List

- `web/src/api/settings.ts` — New settings API function
- `web/src/api/settings.test.ts` — Settings API test
- `web/src/composables/useWeekSelector.ts` — Added `maxWeeks` parameter
- `web/src/composables/useWeekSelector.test.ts` — 2 new tests for maxWeeks
- `web/src/views/ItemGroupsView.vue` — Fetch settings, pass weeksInAdvanced to selector
- `web/src/views/ItemsView.vue` — Fetch settings, maxBookingDate computed, limit error
  passthrough
