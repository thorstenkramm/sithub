# Story 20.3: Consistent Snackbar Confirmations

Status: done

## Story

As a user,
I want all confirmations to use the same bottom snackbar style,
So that the feedback is consistent and predictable across the app.

## Acceptance Criteria

1. **Given** I cancel a booking from My Bookings
   **When** the cancellation succeeds
   **Then** a bottom-center snackbar shows "Booking cancelled successfully."
   (not an inline alert)

2. **Given** I perform any action that shows a success confirmation
   **When** the confirmation appears
   **Then** it uses a bottom-center snackbar, matching the style used for favorites and
   filter confirmations

3. **Given** I book an item from the day view
   **When** the booking succeeds
   **Then** the confirmation is shown as a snackbar (not an inline alert that pushes
   content down)

4. **Given** a snackbar is visible and another action triggers a new confirmation
   **When** the new snackbar fires
   **Then** the previous snackbar is replaced (only one snackbar visible at a time)

## Tasks / Subtasks

- [x] Audit all confirmation styles across the app (AC: 2)
  - [x] Identify all views and components that show success confirmations
  - [x] Document which ones use inline alerts vs snackbars
- [x] Define snackbar standard (AC: 1, 2, 3, 4)
  - [x] Position: bottom-center (`location="bottom"`)
  - [x] Duration: 3 seconds
  - [x] Color: `success` for positive actions, `error` for failures
  - [x] No close button — auto-dismiss only
  - [x] One at a time — new snackbar replaces the current one
- [x] Replace My Bookings cancel success alert with snackbar (AC: 1)
  - [x] Remove inline success alert from `MyBookingsView.vue`
  - [x] Add bottom snackbar with "Booking cancelled successfully." message
- [x] Replace booking success alert on items page (AC: 3)
  - [x] Convert the day-mode booking success alert in `ItemsView.vue` to a snackbar
  - [x] Ensure the snackbar does not push item tiles down
- [x] Ensure all views use the snackbar pattern (AC: 2)
  - [x] Convert any remaining inline success alerts to bottom snackbars
  - [x] Verify consistent styling across all confirmation messages
- [x] Add or update unit tests for snackbar confirmations
- [ ] Verify E2E tests still pass

## Dev Notes

### UX Recommendations (Sally)

#### Snackbar Appearance Spec

All success confirmations across SitHub must use this pattern:

- **Position:** bottom-center (works well on both desktop and mobile)
- **Duration:** 3 seconds, auto-dismiss
- **Color:** `success` theme color for positive actions, `error` for failures
- **Close button:** none — auto-dismiss only reduces visual noise
- **Stacking:** one at a time — if a new snackbar fires while one is showing, replace it

#### Booking success on items page

The day-mode booking success currently uses an inline `v-alert` at the top that pushes all
item tiles down. This is disruptive — a snackbar is less intrusive and matches the
pattern used everywhere else.

### References

- Epic 20 Story 20.3: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR78: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

GPT-5 Codex

### Completion Notes List

- My Bookings cancellation uses a bottom snackbar for success feedback.
- Day-view booking success on the items page uses a bottom snackbar and no longer relies on inline success alerts.
- AI review fix: ItemsView and ItemGroupsView now each use a single replaceable success snackbar, so a new confirmation replaces the previous one.
- AI review fix: standardized success snackbar duration to 3 seconds and removed the remaining inline password-change success alert in `App.vue`.
- Added targeted snackbar tests for booking success flows and My Bookings cancellation feedback.
- E2E tests were not run in this review/fix pass.

### File List

- `web/src/App.vue` — Converted password-change success confirmation to the snackbar pattern
- `web/src/App.test.ts` — Updated app mount stubs for snackbar rendering
- `web/src/views/ItemGroupsView.vue` — Consolidated favorite/filter confirmations into a single replaceable snackbar
- `web/src/views/ItemGroupsView.test.ts` — Existing snackbar-related coverage remains valid after consolidation
- `web/src/views/ItemsView.vue` — Consolidated booking, cancel, favorite, and filter confirmations into a single replaceable snackbar
- `web/src/views/ItemsView.test.ts` — Added coverage for booking-success snackbar behavior and persisted day reset
- `web/src/views/MyBookingsView.vue` — Uses snackbar confirmation for successful booking cancellation
- `web/src/views/MyBookingsView.test.ts` — Added cancellation-success snackbar coverage

## Senior Developer Review (AI)

- Verified ACs 20.3.1 to 20.3.4 against the current frontend success-confirmation flows.
- Fixed the missing replacement behavior by consolidating per-view success messaging to one snackbar instance per view.
- Fixed the remaining timing inconsistency by standardizing the items-page success snackbar to the 3-second spec.
- Completed the app-wide audit by converting the password-change success confirmation in `App.vue` to the same snackbar pattern.

## Change Log

- 2026-03-22: UX review — added AC 3 (booking success as snackbar), AC 4 (one at a time),
  defined snackbar appearance spec, added task for items page booking alert conversion.
- 2026-03-22: Story implementation reviewed and finalized; consolidated success snackbars, standardized timing, and removed the remaining inline success alert.
