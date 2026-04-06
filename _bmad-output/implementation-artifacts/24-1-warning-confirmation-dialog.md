# Story 24.1: Warning Confirmation Dialog (Day Mode)

Status: done

## Story

As a user,
I want to see a confirmation dialog with the item's warning before booking,
so that I am aware of restrictions and can decide whether to proceed or choose a different item.

## Acceptance Criteria

1. **Given** I click BOOK on an item that has a warning in day booking mode,
   **When** the warning dialog appears,
   **Then** it displays the item name (truncated with ellipsis if longer than the dialog width),
   the warning text, a CONFIRM button, and a CANCEL button.

2. **Given** the warning dialog is displayed,
   **When** I click CONFIRM,
   **Then** the booking proceeds as normal.

3. **Given** the warning dialog is displayed,
   **When** I click CANCEL,
   **Then** the booking is aborted and I remain on the booking view with no booking created.

4. **Given** the warning dialog is displayed with a "Don't show again" checkbox,
   **When** I check "Don't show again" and click CONFIRM,
   **Then** the booking proceeds and the suppression is stored in localStorage keyed by item ID.

5. **Given** I have previously checked "Don't show again" for an item,
   **When** I book that same item again,
   **Then** the warning dialog is skipped and the booking proceeds immediately.

6. **Given** an item has no warning configured,
   **When** I click BOOK,
   **Then** no warning dialog is shown and the booking proceeds as before.

## Tasks / Subtasks

- [x] Task 1: Create warning suppression composable (AC: #4, #5)
  - [x] 1.1 Create `web/src/composables/useWarningSuppression.ts` with `isWarningSuppressed(itemId)` and `suppressWarning(itemId)` functions
  - [x] 1.2 Use localStorage key `sithub_warning_suppressed` storing a JSON array of suppressed item IDs
  - [x] 1.3 Use `getSafeLocalStorage()` from `composables/storage.ts` for safe access
  - [x] 1.4 Write unit tests for the composable

- [x] Task 2: Add i18n keys for warning dialog (AC: #1)
  - [x] 2.1 Add keys to all 5 locale files (en, de, es, fr, uk): `items.warningDialogTitle`, `items.warningConfirm`, `items.warningCancel`, `items.warningDontShowAgain`
  - [x] 2.2 English values: "WARNING!", "CONFIRM", "CANCEL", "Don't show again."

- [x] Task 3: Add warning confirmation dialog template to ItemsView.vue (AC: #1)
  - [x] 3.1 Add dialog state refs: `showWarningDialog`, `warningDialogItemId`, `warningDialogItemName`, `warningDialogMessage`, `warningDontShowAgain`
  - [x] 3.2 Add `v-dialog` with `data-cy="warning-dialog"`, max-width 400, persistent; display item name (truncated via CSS `text-overflow: ellipsis`), warning text, checkbox, CONFIRM and CANCEL buttons
  - [x] 3.3 Add `data-cy` selectors: `warning-dialog`, `warning-confirm-btn`, `warning-cancel-btn`, `warning-dont-show-checkbox`

- [x] Task 4: Intercept day-mode booking flow with warning check (AC: #1, #2, #3, #4, #5, #6)
  - [x] 4.1 Extract a `requestBooking(itemId)` function that checks for warning + suppression before calling `bookItem(itemId)`
  - [x] 4.2 If item has warning and not suppressed: populate dialog state and show dialog; CONFIRM handler calls `bookItem()` (and `suppressWarning()` if checkbox checked); CANCEL handler resets dialog state
  - [x] 4.3 If item has no warning or warning is suppressed: call `bookItem()` directly
  - [x] 4.4 Update day-mode BOOK button `@click` to call `requestBooking(entry.id)` instead of `bookItem(entry.id)`

- [x] Task 5: Write unit tests for warning dialog integration (AC: #1-#6)
  - [x] 5.1 Test: clicking BOOK on item with warning shows dialog with correct item name and warning text
  - [x] 5.2 Test: clicking CONFIRM proceeds with booking (bookItem called)
  - [x] 5.3 Test: clicking CANCEL hides dialog without booking
  - [x] 5.4 Test: checking "Don't show again" + CONFIRM stores suppression in localStorage
  - [x] 5.5 Test: suppressed item skips dialog and books directly
  - [x] 5.6 Test: item without warning books directly without dialog

- [x] Task 6: Run full validation suite
  - [x] 6.1 `npm run type-check` passes
  - [x] 6.2 `npm run lint` passes
  - [x] 6.3 `npx vitest run` — 299 tests, all pass
  - [x] 6.4 `npm run build` succeeds

### Review Findings

- [x] [Review][Decision] Warning suppression keyed by item ID only — changed warnings stay hidden [HIGH] — Fixed: suppression now keys by itemId::warningHash; changed warnings auto-reset suppression
- [x] [Review][Patch] Double-click race: BOOK on item-B while dialog shows for item-A overwrites dialog state [MED] — Fixed: added `if (showWarningDialog.value) return;` guard in requestBooking()
- [x] [Review][Patch] confirmWarningDialog does not clear dialog state refs unlike cancelWarningDialog [LOW] — Fixed: confirmWarningDialog now resets all dialog state refs
- [x] [Review][Defer] Unbounded localStorage growth — no eviction or size limit — deferred, acceptable for MVP item counts
- [x] [Review][Defer] No user-scoped suppression key — shared browser leaks preferences — deferred, pre-existing pattern (all localStorage keys are unscoped)

## Dev Notes

### Architecture and Patterns

- **Booking view**: All booking logic lives in `web/src/views/ItemsView.vue` (1700+ lines)
- **Day-mode booking**: `bookItem(itemId)` function at ~line 1510 calls `createBooking()` API
- **BOOK button click**: `@click="bookItem(entry.id)"` at ~line 350 in day-mode template
- **Warning field**: `ItemAttributes.warning?: string` defined in `web/src/api/items.ts` line 8
- **Items data**: `items` ref in ItemsView contains the fetched items with attributes including `warning`

### Existing Dialog Patterns to Follow

- **ConfirmDialog component**: `web/src/components/ConfirmDialog.vue` — reusable, but does NOT have a checkbox slot, so use inline `v-dialog` in ItemsView (matches the booking-limit-dialog pattern from Story 23.2)
- **Booking limit dialog pattern (Story 23.2)**: `showLimitDialog` ref + inline `v-dialog` in ItemsView.vue with `data-cy="booking-limit-dialog"`. Follow this exact pattern.
- **Dialog props**: `max-width="400"`, `persistent`, `v-card` wrapper with `v-card-title`, `v-card-text`, `v-card-actions`

### localStorage Pattern

- Use `getSafeLocalStorage()` from `web/src/composables/storage.ts`
- Storage key convention: `sithub_*` prefix (e.g., `sithub_theme`, `sithub_show_weekends`)
- New key: `sithub_warning_suppressed` — store JSON array of item IDs
- Pattern: read on init, write on change

### i18n Pattern

- Locale files: `web/src/locales/{en,de,es,fr,uk}.json`
- Keys grouped under feature namespace: `items.*` for booking-related strings
- All 5 locales must be updated simultaneously

### Testing Patterns

- Unit tests co-located: `web/src/views/ItemsView.test.ts`
- Composable tests: `web/src/composables/__tests__/` or co-located `.test.ts`
- Use `data-cy` selectors in tests for consistency
- Mock localStorage using vi.stubGlobal or direct property mock
- 282 tests currently passing (as of Story 23.2)

### Previous Story Intelligence (Story 23.2)

- Dialog template follows `v-dialog > v-card > v-card-title + v-card-text + v-card-actions` structure
- State refs: boolean for show/hide, string for message content
- `data-cy` selectors on dialog container and action buttons
- Error routing: different error types → different UI (dialog vs snackbar)
- Full suite: type-check, lint, vitest, build must all pass

### Git Intelligence

- Latest commit: `a2e3cf6` feat: epic 23 — UI bug fixes
- Recent pattern: single commit per epic with all story changes bundled

### Warning Data Flow

```
YAML config (warning field on item)
  → Go backend serves via /api/v1/items
    → ItemAttributes.warning?: string
      → ItemsView.vue reads entry.attributes.warning
        → Currently: tooltip icon (folded) + v-alert (expanded)
        → NEW: intercept bookItem() to show confirmation dialog
```

### Project Structure Notes

- All changes scoped to `web/src/` (frontend-only story)
- New file: `web/src/composables/useWarningSuppression.ts`
- Modified files: `web/src/views/ItemsView.vue`, locale JSON files
- No backend changes required — warning field already exists in API

### References

- [Source: web/src/views/ItemsView.vue — bookItem() function ~line 1510]
- [Source: web/src/api/items.ts — ItemAttributes interface line 8]
- [Source: web/src/composables/storage.ts — getSafeLocalStorage()]
- [Source: web/src/components/ConfirmDialog.vue — reusable dialog pattern]
- [Source: _bmad-output/planning-artifacts/epics.md — Epic 24 Story 24.1]
- [Source: _bmad-output/implementation-artifacts/23-2-booking-limit-error-modal.md — dialog pattern]
- [Source: sithub_areas.example.yaml — warning field in YAML config]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- Initial composable tests: 8/8 pass
- ItemsView warning dialog tests: first run had 2 failures due to uncleaned createBookingMock call history; fixed by adding mockClear() in test beforeEach
- Final full suite: 299/299 pass (37 files)

### Completion Notes List

- Created `useWarningSuppression` composable with `isWarningSuppressed()` and `suppressWarning()` functions; stores JSON array of item IDs in localStorage under `sithub_warning_suppressed`
- Added i18n keys (`warningDialogTitle`, `warningConfirm`, `warningCancel`, `warningDontShowAgain`) to all 5 locale files
- Added persistent warning dialog to ItemsView.vue following the booking-limit-dialog pattern: item name with ellipsis truncation, warning text, dont-show-again checkbox, CONFIRM/CANCEL buttons
- Created `requestBooking()` function that intercepts the day-mode booking flow: checks for warning + suppression before calling `bookItem()`
- Updated BOOK button `@click` from `bookItem(entry.id)` to `requestBooking(entry.id)`
- 8 composable unit tests + 6 integration tests for the warning dialog behavior
- All validation gates pass: type-check, lint, 299 tests, build

### Change Log

- 2026-04-06: Story 24.1 implementation complete — warning confirmation dialog for day booking mode

### File List

- web/src/composables/useWarningSuppression.ts (new)
- web/src/composables/useWarningSuppression.test.ts (new)
- web/src/views/ItemsView.vue (modified)
- web/src/views/ItemsView.test.ts (modified)
- web/src/locales/en.json (modified)
- web/src/locales/de.json (modified)
- web/src/locales/es.json (modified)
- web/src/locales/fr.json (modified)
- web/src/locales/uk.json (modified)
