# Story 24.2: Sequential Warning Dialogs (Week Mode)

Status: done

## Story

As a user,
I want warnings for multiple items shown one after another when booking in week mode,
so that I can review each item's restrictions before confirming the full week booking.

## Acceptance Criteria

1. **Given** I am in week booking mode and have selected days on multiple items that have
   warnings,
   **When** I click "Confirm My Booking",
   **Then** the warning dialogs are shown sequentially, one per item with a warning, each
   identifying the item by name.

2. **Given** a sequential warning dialog is displayed for item A,
   **When** I click CONFIRM,
   **Then** the next item's warning dialog is shown (or booking proceeds if no more warnings
   remain).

3. **Given** a sequential warning dialog is displayed for item B,
   **When** I click CANCEL,
   **Then** the entire week booking is aborted and no bookings are created for any item.

4. **Given** I have previously suppressed warnings for some items via "Don't show again",
   **When** I book a week that includes those items,
   **Then** the suppressed items' warning dialogs are skipped; only unsuppressed warnings
   are shown.

5. **Given** all items in my week booking have their warnings suppressed,
   **When** I click "Confirm My Booking",
   **Then** the booking proceeds immediately with no warning dialogs.

## Tasks / Subtasks

- [x] Task 1: Implement warning queue logic in ItemsView.vue (AC: #1, #2, #3, #4, #5)
  - [x] 1.1 Add state refs for the warning queue: `warningQueue` (array of `{itemId, itemName, warning}`), `warningQueueMode` (ref to distinguish day vs week dialog usage)
  - [x] 1.2 Create `collectUnsuppressedWarnings()` helper: extract unique item IDs from `weekSelections`, find items with unsuppressed warnings, return as queue array
  - [x] 1.3 Create `startWeekWarningFlow()` function: called by "Confirm My Booking" button instead of `submitWeekBookings()`; calls `collectUnsuppressedWarnings()`; if queue is empty proceed directly to `submitWeekBookings()`; otherwise set `warningQueueMode = 'week'` and show first dialog
  - [x] 1.4 Update `confirmWarningDialog()`: if `warningQueueMode === 'week'`, shift queue and show next dialog (or call `submitWeekBookings()` if queue exhausted); for day mode keep existing behavior
  - [x] 1.5 Update `cancelWarningDialog()`: if `warningQueueMode === 'week'`, clear queue and abort entirely (do not call `submitWeekBookings()`); reset all state
  - [x] 1.6 Update week-mode "Confirm My Booking" button `@click` from `submitWeekBookings` to `startWeekWarningFlow`

- [x] Task 2: Write unit tests for sequential warning flow (AC: #1-#5)
  - [x] 2.1 Test: week confirm with 2 warned items shows first warning dialog, CONFIRM shows second, CONFIRM proceeds with booking
  - [x] 2.2 Test: CANCEL on first warning dialog aborts entire booking (createBooking not called)
  - [x] 2.3 Test: CANCEL on second warning dialog aborts entire booking
  - [x] 2.4 Test: suppressed items are skipped in the queue; only unsuppressed shown
  - [x] 2.5 Test: all warnings suppressed — booking proceeds immediately without dialog
  - [x] 2.6 Test: single warned item in week mode shows one dialog, CONFIRM proceeds

- [x] Task 3: Run full validation suite
  - [x] 3.1 `npm run type-check` passes
  - [x] 3.2 `npm run lint` passes
  - [x] 3.3 `npx vitest run` — 306 tests, all pass
  - [x] 3.4 `npm run build` succeeds

## Dev Notes

### Architecture and Patterns

- **Week booking flow**: `submitWeekBookings()` at ~line 1426 in `ItemsView.vue`; collects entries from `weekSelections` Set (keys: `itemId::date`), executes parallel `createBooking()` calls via `Promise.allSettled()`
- **Week confirm button**: `@click="submitWeekBookings"` at ~line 669 with `data-cy="week-confirm-btn"`
- **Items data**: `items` ref contains all fetched items including `attributes.warning`
- **Warning dialog**: Already exists from Story 24.1 — `showWarningDialog`, `warningDialogItemId`, `warningDialogItemName`, `warningDialogMessage`, `warningDontShowAgain` refs + `v-dialog` template
- **Suppression composable**: `useWarningSuppression` with `isWarningSuppressed(itemId, warning)` and `suppressWarning(itemId, warning)` — keys by `itemId::warningHash`

### Implementation Strategy

The key insight: reuse the existing warning dialog from Story 24.1. Add a queue mechanism that:
1. Before week booking, collects unique items with unsuppressed warnings
2. Shows the same dialog sequentially for each
3. On CONFIRM: advance queue; on CANCEL: abort all
4. When queue exhausted: call `submitWeekBookings()`

No new dialog template needed — extend existing `confirmWarningDialog` and `cancelWarningDialog` handlers with queue awareness via `warningQueueMode`.

### Week Selection Parsing

```typescript
// weekSelections contains keys like "itemId::date"
const getWeekSelectionKey = (itemId: string, date: string) => `${itemId}::${date}`;

// To extract unique item IDs from selections:
const uniqueItemIds = new Set(
  [...weekSelections.value].map(key => key.split('::')[0])
);
```

### Previous Story Intelligence (Story 24.1)

- **Dialog guard**: `if (showWarningDialog.value) return;` prevents double-click race
- **State cleanup**: Both `confirmWarningDialog` and `cancelWarningDialog` reset all dialog state refs
- **Suppression hash**: `isWarningSuppressed(itemId, warning)` — keys by itemId + warning content hash
- **i18n keys**: Already added in Story 24.1 (`warningDialogTitle`, `warningConfirm`, `warningCancel`, `warningDontShowAgain`)
- **Test pattern**: Clear `createBookingMock` in beforeEach; use `data-cy` selectors; `mountView()` + `flushPromises()` + trigger click

### Review Findings from Story 24.1 to Carry Forward

- Always reset all dialog state refs on both confirm and cancel paths
- Guard entry points against dialog already being open
- `confirmWarningDialog` and `cancelWarningDialog` need to handle both day and week modes

### Project Structure Notes

- All changes scoped to `web/src/views/ItemsView.vue` and its test file
- No new files needed — extends existing composable and dialog infrastructure
- No backend changes — frontend-only story
- No new i18n keys needed — reuses keys from Story 24.1

### References

- [Source: web/src/views/ItemsView.vue — submitWeekBookings() ~line 1426]
- [Source: web/src/views/ItemsView.vue — weekSelections ref line 1124]
- [Source: web/src/views/ItemsView.vue — week-confirm-btn ~line 669]
- [Source: web/src/views/ItemsView.vue — requestBooking/confirmWarningDialog ~line 1560]
- [Source: web/src/composables/useWarningSuppression.ts — isWarningSuppressed/suppressWarning]
- [Source: _bmad-output/implementation-artifacts/24-1-warning-confirmation-dialog.md]
- [Source: _bmad-output/planning-artifacts/epics.md — Epic 24 Story 24.2]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- vue-tsc doesn't narrow types through `continue` or `||` in SFC scripts — used explicit `!` assertions and separate `filter` chains
- v-btn test stub causes double-fire of click handlers due to `v-bind="$attrs"` passing through event listeners + `$emit('click')` — fixed with custom `inheritAttrs: false` stub in week mode tests
- `ref` array `.shift()` mutation not reactive in vue-tsc — used immutable `slice(1)` instead

### Completion Notes List

- Added `warningQueue` and `warningQueueMode` state refs for queue-based sequential warning display
- Created `collectUnsuppressedWarnings()` — extracts unique items from `weekSelections`, filters to unsuppressed warnings using `findWeekItem()` to look up items from `weekData`
- Created `startWeekWarningFlow()` — entry point for week mode: builds queue, shows first dialog or proceeds directly if empty
- Extended `confirmWarningDialog()` with queue mode: advances queue in-place (updates dialog content without close/reopen), calls `submitWeekBookings()` when exhausted
- Extended `cancelWarningDialog()` with queue mode: clears queue and aborts
- Added `resetWarningDialogState()` helper to DRY up state cleanup
- Added `.stop` modifier on dialog buttons to prevent native event propagation
- Updated week confirm button from `submitWeekBookings` to `startWeekWarningFlow`
- 6 new integration tests covering all 5 acceptance criteria
- All validation gates pass: type-check, lint, 306 tests, build

### Change Log

- 2026-04-06: Story 24.2 implementation complete — sequential warning dialogs for week booking mode

### File List

- web/src/views/ItemsView.vue (modified)
- web/src/views/ItemsView.test.ts (modified)
