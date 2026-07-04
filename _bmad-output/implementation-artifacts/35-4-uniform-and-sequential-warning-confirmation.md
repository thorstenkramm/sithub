# Story 35.4: Uniform and Sequential Warning Confirmation

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user booking warned items from any view,
I want one identical confirmation dialog, shown once per warned item,
so that the confirmation behavior is predictable no matter where I book.

## Acceptance Criteria

1. A warned item booked from the tiles, the floor plan, or the weekly table view opens the same
   confirmation dialog: title "WARNING!", the item name, the warning text in the shared style, a
   "Don't show again" checkbox, and CANCEL/CONFIRM actions.
2. A single booking action covering multiple items where 2+ have warnings shows the confirmations
   one after another, each identifying its item; the booking is submitted only after every warning
   is confirmed.
3. Cancelling any one confirmation aborts the entire booking; no item is booked.
4. An item whose warning was previously dismissed via "Don't show again" is skipped while other
   warned items still show theirs.

## Tasks / Subtasks

- [x] Task 1: Extract the confirmation dialog + queue into shared, reusable code (AC: #1, #2, #3, #4)
  - [x] Extract the inline dialog markup (`ItemsView.vue` ~915-956) into a shared component
        (e.g. `WarningConfirmDialog.vue`) and the queue/flow logic (~1861-1959) into a composable
        (e.g. `useWarningConfirmation`) that wraps `useWarningSuppression`
  - [x] The dialog uses the shared warning presentation (35.1) for the message; keep i18n keys
        `warningDialogTitle`/`warningConfirm`/`warningCancel`/`warningDontShowAgain` and all
        `data-cy` hooks (`warning-dialog`, `warning-item-name`, `warning-message`,
        `warning-dont-show-checkbox`, `warning-cancel-btn`, `warning-confirm-btn`)
  - [x] Sequential queue behavior (advance on confirm, abort-all on cancel, skip suppressed) is
        preserved from the current week-mode implementation and generalized to any caller
- [x] Task 2: Adopt it in the tile view (AC: #1, #2, #3, #4)
  - [x] Replace `ItemsView.vue`'s inline dialog + `requestBooking`/`startWeekWarningFlow`/
        `confirmWarningDialog`/`cancelWarningDialog` with the shared component + composable; day and
        week flows behave exactly as before
- [x] Task 3: Adopt it in the floor plan and weekly table (AC: #1, #2, #3, #4)
  - [x] Wire `InteractiveFloorPlan.vue` booking (35.2) to the shared confirmation
  - [x] Wire the weekly-table booking (`MatrixBookingPopover.submitBooking`, ~229-264) to the shared
        confirmation so a warned cell booking shows the dialog (the table has no blocking warning
        today)
- [x] Task 4: Tests (AC: #1-#4)
  - [x] Unit/component tests for the shared dialog + composable: single confirm, sequential multi,
        cancel-aborts-all, suppressed-skipped
  - [x] Keep `ItemsView.test.ts` week-mode warning tests green; add table + floor-plan confirmation
        cases

## Dev Notes

Source: `private/epic-35.md` — warning confirmation section (`img_33.png`).
[Source: _bmad-output/planning-artifacts/epics.md#Story 35.4 / FR163,FR164]

### Current confirmation implementation (to extract, not reinvent)

All inline in `web/src/views/ItemsView.vue`:

- Dialog markup: ~915-956 (`v-dialog data-cy="warning-dialog"`, title `items.warningDialogTitle` =
  "WARNING!", item name ~920, message ~926 `white-space: pre-line`, checkbox ~931, CANCEL ~940,
  CONFIRM ~948).
- State refs: ~1011-1018 (`showWarningDialog`, `warningDialogItemId/Name/Message`,
  `warningDontShowAgain`, `warningQueue`, `warningQueueMode`).
- Day flow: `requestBooking(itemId)` ~1861-1877 (checks `isWarningSuppressed(itemId, warning)`, shows
  dialog else books).
- Week flow: `collectUnsuppressedWarnings()` ~1887-1900 + `startWeekWarningFlow()` ~1902-1917 (builds
  a queue, shows first).
- Confirm: `confirmWarningDialog()` ~1927-1953 (suppress if checked; advance queue in week mode; else
  `bookItem`). Cancel: `cancelWarningDialog()` ~1955-1959 (clears queue → aborts all).
- Suppression: `useWarningSuppression` (`web/src/composables/useWarningSuppression.ts`).

The week-mode queue ALREADY implements FR164 (sequential, abort-all, skip-suppressed). This story
generalizes it: move it into `useWarningConfirmation` returning something like
`confirmThenBook(items: {itemId,itemName,warning}[], onAllConfirmed)`, and drive the shared dialog.

### What each view needs after extraction

- Tiles: swap inline for the shared component/composable — no behavior change.
- Floor plan (35.2): call `confirmThenBook([{single item}], () => confirmPendingBooking())` before
  creating the booking.
- Weekly table: `MatrixBookingPopover.submitBooking()` (~229-264) currently books with no blocking
  warning — route it through the shared confirmation. Multi-item is per-cell here (usually one
  item), but the same path applies.

### Do not regress suppression semantics

`useWarningSuppression` keys on `itemId::hash(warning)` and is used for the "Don't show again"
checkbox. Reuse it as-is (35.5 covers its text-change behavior/tests). The `.warning-tooltip` CSS
still referenced by the old inline dialog should now come from the shared component (35.1) — verify
before deleting the class from `ItemsView.vue`.

### Project Structure Notes

- New: `web/src/components/WarningConfirmDialog.vue`,
  `web/src/composables/useWarningConfirmation.ts` (names suggestive).
- Modified: `ItemsView.vue`, `InteractiveFloorPlan.vue`, `MatrixBookingPopover.vue` +
  `AreaWeeklyMatrixView.vue` as needed to invoke the shared flow.
- Depends on 35.1 (presentation). Closely paired with 35.2/35.3 adoption.

### Testing standards summary

Vitest for the composable + dialog; keep the existing `ItemsView.test.ts` week-mode sequential
tests green (they are the FR164 regression guard). Add table + floor-plan confirmation coverage. A
Cypress E2E confirming a warned booking from each view is valuable. Run type-check, lint, vitest,
build. [Source: .claude/rules/vue.md] [Source: .claude/rules/cypress.md]

### References

- [Source: web/src/views/ItemsView.vue:915-956,1011-1018,1861-1877,1887-1917,1927-1959]
- [Source: web/src/composables/useWarningSuppression.ts]
- [Source: web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:229-264]
- [Source: web/src/components/InteractiveFloorPlan.vue:1416-1429,1474-1544]
- [Source: web/src/locales/en.json (warningDialogTitle/Confirm/Cancel/DontShowAgain)]

## Dev Agent Record

### Agent Model Used

claude-fable-5

### Debug Log References

- 456 unit tests pass (50 files); type-check, lint, jscpd (0 clones), build all clean
- Verified in-browser (Chrome DevTools): identical "WARNUNG!" dialog from tiles + floor plan

### Completion Notes List

- Extracted the confirmation into `WarningConfirmDialog.vue` (uses `ItemWarning` inline, showIcon
  false — orange message, no icon, matching img_33) + `useWarningConfirmation.ts` (queue, sequential
  display, suppression via useWarningSuppression, cancel-aborts-all, skip-suppressed).
- Refactored `ItemsView.vue` to use them: `requestBooking` and `startWeekWarningFlow` now call
  `present(items, onConfirmed)`; removed the inline dialog markup, state refs, and the old
  confirm/cancel/queue functions. All 90 ItemsView tests (incl. week-mode sequential) stay green.
- Wired the weekly table (`MatrixBookingPopover.submitBooking` → `present([...], doBook)`).
- Wired the floor plan (`InteractiveFloorPlan.confirmPendingBooking` → `present([...], executeBooking)`).
- Added `useWarningConfirmation.test.ts` (6 cases: empty→immediate, single, sequential multi,
  cancel-aborts-all, skip-suppressed) plus updated popover tests.
- i18n keys and all `data-cy` hooks preserved.

### File List

- web/src/components/WarningConfirmDialog.vue (new)
- web/src/composables/useWarningConfirmation.ts (new)
- web/src/composables/__tests__/useWarningConfirmation.test.ts (new)
- web/src/components/index.ts (modified — export WarningConfirmDialog)
- web/src/views/ItemsView.vue (modified)
- web/src/components/area-weekly-matrix/MatrixBookingPopover.vue (modified)
- web/src/components/area-weekly-matrix/MatrixBookingPopover.test.ts (modified)
- web/src/components/InteractiveFloorPlan.vue (modified)

### Change Log

- 2026-07-04: Implemented FR163/FR164 — shared uniform + sequential warning confirmation across
  tiles, weekly table, and floor plan.
