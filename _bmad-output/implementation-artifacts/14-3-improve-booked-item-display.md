# Story 14.3: Improve Booked Item Display

Status: done

## Story

As a user,
I want booked item details to be clearly readable and booking results to use icons,
So that I can quickly understand who booked what and whether my bookings succeeded.

## Acceptance Criteria

1. **Given** I am viewing items in day booking mode
   **When** an item is booked by another user
   **Then** the "Not available for \<date\>" message is not shown
   **And** the booker name is displayed at body-2 size or larger (not caption)
   **And** any booking note is displayed at body-2 size or larger

2. **Given** I submit bookings (day or week mode)
   **When** results are displayed
   **Then** each successful booking shows a green checkmark icon with item name and day
   **And** each failed booking shows a red warning icon with item name and error detail
   **And** raw text labels like "Booked" are replaced by the icons

## Tasks / Subtasks

- [x] Remove "Not available" text (AC: 1)
  - [x] In `ItemsView.vue`: replaced the "Not available for \<date\>" div with an empty spacer
  - [x] Removed unused `formattedDate` computed property
- [x] Increase booker name font size (AC: 1)
  - [x] In `ItemsView.vue`: changed `text-caption` to `text-body-2` on booker name element
- [x] Increase booking note font size (AC: 1)
  - [x] In `ItemsView.vue`: changed `text-caption` to `text-body-2` on booking note element
- [x] Day booking results already use icons (AC: 2)
  - [x] Verified: `v-alert` with `type="success"` and `type="error"` already displays icons
    automatically via Vuetify
- [x] Update week booking results to use icons consistently (AC: 2)
  - [x] Week results already show `mdi-check-circle` (green) and `mdi-close-circle` (red) icons
  - [x] Removed redundant "Booked" text label — the green icon alone now indicates success
- [x] Update E2E tests
  - [x] Updated `week-booking.cy.ts`: changed assertion from "Booked" to "Booking Results"
- [x] Verify all tests pass

## Dev Notes

### Architecture: Frontend-Only Story

All changes are in `web/src/views/ItemsView.vue`. No backend changes required.

### Day Mode Result Alerts

The day mode uses `v-alert` components for success/error feedback. These already have
`type="success"` and `type="error"` which provide colored backgrounds and icons automatically
through Vuetify. No additional icon changes were needed.

### Week Mode Results Already Have Icons

The week booking results section already renders `mdi-check-circle` (green) and
`mdi-close-circle` (red) icons per result line. The change removes the redundant "Booked" text
since the green icon already communicates success.

### Font Size Reference

Vuetify typography classes:
- `text-caption`: 12px / 0.75rem (previous — too small)
- `text-body-2`: 14px / 0.875rem (current — readable)

### References

- Epic 14 Story 14.3: `_bmad-output/planning-artifacts/epics.md` (Epic 14 Stories section)
- FR44, FR45, FR46: `_bmad-output/planning-artifacts/prd.md`
- ItemsView: `web/src/views/ItemsView.vue`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Replaced "Not available for \<date\>" text with empty spacer div
- Removed unused `formattedDate` computed property (dead code cleanup)
- Changed booker name from `text-caption` to `text-body-2` for readability
- Changed booking note from `text-caption` to `text-body-2` for readability
- Removed redundant "Booked" text from week booking results (icon suffices)
- Updated week-booking.cy.ts assertion for the changed result text
- All 132 unit tests pass, all 51 E2E tests pass
- Type check, ESLint, build, and code duplication checks all pass
- Code review fix: show item name + date in day-mode booking success and error alerts
- Code review fix: preserve booking success details after saving a note
- Code review fix: add unit test for on-behalf booking payload without display name
- Code review fix: update Cypress booking success/error assertions to match new alert text

### Change Log

- 2026-02-14: Implemented Story 14.3 - improved booked item display with larger fonts,
- 2026-02-14: Code review fixes for booking result details and tests
  removed redundant text, icon-only success indicators

### File List

- web/src/views/ItemsView.vue (modified - font sizes, removed "Not available" text, removed
- web/src/api/bookings.ts (modified - omit empty for_user_name on behalf)
- web/src/api/bookings.test.ts (modified - added on-behalf payload test)
- web/cypress/e2e/items.cy.ts (modified - updated booking alert assertions)
- web/cypress/e2e/booking-notes.cy.ts (modified - updated booking alert assertions)
  "Booked" label, removed unused formattedDate)
- web/cypress/e2e/week-booking.cy.ts (modified - updated result assertion)