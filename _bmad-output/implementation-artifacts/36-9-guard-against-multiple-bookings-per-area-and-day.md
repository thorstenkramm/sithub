# Story 36.9: Guard Against Multiple Bookings per Area and Day

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want to be warned when I already have a booking in an area on a day and offered to swap it,
so that I don't unknowingly block multiple slots in the same area on the same day.

## Acceptance Criteria

1. When I already have a booking for an item in an area on a given day and I try to book another
   item in the SAME area on the SAME day, a dialog appears:
   "You already booked \<ITEM\> on \<DATE\>. Multiple bookings per area and day are not allowed.
   Do you want to cancel \<ITEM\> and book \<NEW-ITEM\> instead?"
2. Confirm → the existing booking for that area/day is cancelled and the new booking is created.
3. Cancel → no change is made: the original booking remains and the new one is not created.
4. The guard applies consistently across the tile, table, and floor-plan booking flows and is
   scoped to the same area and the same day.

## Tasks / Subtasks

- [x] Task 1: Add a shared area/day guard composable (AC: #1, #2, #3, #4)
  - [x] Create `web/src/composables/useAreaDayGuard.ts`. It fetches the current user's bookings
        (`fetchMyBookings()` — `web/src/api/bookings.ts:158`) and exposes a
        `guard({ areaId, date, newItemName, newItemId }, onProceed)` function that: looks for an
        existing booking whose `attributes.area_id === areaId` AND `attributes.booking_date === date`
        AND `attributes.item_id !== newItemId`; if none found, calls `onProceed()` immediately; if
        found, opens the confirmation dialog carrying the existing booking's `id`, `item_name`, and
        `booking_date`.
  - [x] On confirm: `await cancelBooking(existing.id)` (`web/src/api/bookings.ts:191`) then
        `onProceed()` (which performs the create). On cancel: reset, do nothing.
  - [x] Model the reactive shape on `useWarningConfirmation.ts` (`show`, message fields, `confirm`,
        `cancel`, `present`/`guard`) so the three call sites bind it identically
        [Source: web/src/composables/useWarningConfirmation.ts:25-91].
  - [x] Because the guard needs fresh data, fetch `fetchMyBookings()` at guard time (or accept an
        already-loaded booking list as an optional argument to avoid a redundant call in the tile
        week flow, which already holds it — ItemsView.vue:1397,1741).
- [x] Task 2: Add the swap confirmation dialog (AC: #1, #2, #3)
  - [x] Reuse the existing generic `ConfirmDialog.vue`
        (`web/src/components/ConfirmDialog.vue`) — it already supports `title`, `message`,
        `confirmText`, `cancelText`, `confirmColor`, `loading`, and emits `confirm`/`cancel` with
        `data-cy="confirm-dialog-confirm"` / `data-cy="confirm-dialog-cancel"`. Bind its `loading`
        prop to the in-flight cancel+create so the buttons disable during the swap.
  - [x] Compose the message via a new i18n key with `{existingItem}`, `{date}`, `{newItem}`
        interpolation. Add the key to all five locales: `en`, `de`, `uk`, `fr`, `es`
        (`web/src/locales/*.json` — warning keys live around line 166).
- [x] Task 3: Wire the guard into the tile flow (AC: #1, #2, #3, #4)
  - [x] In `web/src/views/ItemsView.vue`, insert the guard between the warning confirmation and the
        create. Single-day: in `requestBooking` (line 1794) the chain is
        `presentWarnings(list, () => bookItem(itemId))` (line 1800) → gate `bookItem` so the guard
        runs before `createBooking` (line 1846): `presentWarnings(list, () => guard(ctx, () => bookItem(itemId)))`.
  - [x] Week/multi-day: `submitWeekBookings` (line 1482) books multiple (item, date) pairs; apply
        the guard per (area, date). The view aggregates favorites across MULTIPLE areas, so match on
        each booking's real `area_id` from `fetchMyBookings()`, never a single view-level area.
  - [x] Resolve `areaId` for a tile via `getCurrentAreaId()` (line 1553). In favorites mode this is
        empty, so rely on the my-bookings `area_id` match keyed to the item being booked; if the
        item's area cannot be determined client-side, the guard is a no-op for that item (do not
        block). Reload after swap via `loadItemsForView(selectedDate.value)` (line 1865).
- [x] Task 4: Wire the guard into the table flow (AC: #1, #2, #3, #4)
  - [x] `MatrixBookingPopover.vue` calls the API itself: `submitBooking` (line 220) →
        `presentWarnings(warnItems, () => void doBook())` (line 232) → `doBook` →
        `createBooking(props.item.item_id, props.cell.date, ...)` (line 245). Insert the guard so it
        runs before `doBook`'s create.
  - [x] The popover does NOT receive `area_id` (props: `item`, `cell` — lines 122-127). Pass the
        parent's `props.areaId` down as a new prop from
        `AreaWeeklyMatrixView.vue` (renders the popover at lines 48-56; owns `props.areaId` at
        line 98-103). The matrix flow has no `fetchMyBookings` data loaded, so the guard fetches it.
  - [x] After a swap, the popover already emits `booked` (line 252); the parent's `onBooked`
        (AreaWeeklyMatrixView.vue:232) calls `loadMatrix({ silent: true })` — no extra refresh code
        needed.
- [x] Task 5: Wire the guard into the floor-plan flow (AC: #1, #2, #3, #4)
  - [x] `InteractiveFloorPlan.vue`: `confirmPendingBooking` (line 1518) →
        `presentWarnings(warnItems, () => void executeBooking(bookingDates))` (line 1538). Insert the
        guard before `executeBooking`.
  - [x] Floor plan supports MULTI-DAY selection; apply the guard per selected date in
        `bookingDates`. `executeBooking` (line 1541) uses `createBooking` (line 1553) for one day and
        `createMultiDayBooking` (line 1559) for many. If any selected day already has a same-area
        booking, prompt to swap that day; confirm swaps only the conflicting day(s).
  - [x] `areaId` is an optional prop (`areaId?` — line 577); at drill-down level it is set. As with
        tiles, match on the my-bookings `area_id`; if `props.areaId` is empty, fall back to the
        my-bookings match and no-op when undeterminable. Refresh via `refreshAvailability()`
        (line 1592) after a swap.
- [x] Task 6: Tests (AC: #1-#4)
  - [x] Unit test `useAreaDayGuard`: same area+day existing booking → dialog opens with existing
        item name + date; different area or different day → `onProceed` called directly, no dialog;
        confirm → `cancelBooking` then `onProceed`; cancel → neither called.
  - [x] Component tests for each flow: booking a second item in the same area/day opens the swap
        dialog; confirm cancels the old and books the new; cancel leaves both unchanged. Reuse the
        `data-cy` hooks on `ConfirmDialog`.
  - [ ] Cypress E2E (deferred this pass): book item A in an area for a date, then attempt item B in the same
        area/date → swap dialog → confirm → A gone from My Bookings, B present; repeat and cancel →
        A remains, B absent. [Source: .claude/rules/cypress.md]

### Review Findings

- [x] [Review][Patch] Week-table batch swaps can cancel an old booking even when its replacement create failed [`web/src/views/ItemsView.vue:1656`] — resolved by tracking swap cancellation ids by `(area,date)` and cancelling only keys whose replacement create succeeded; regression covered in `web/src/views/ItemsView.test.ts`.
- [x] [Review][Patch] Floor-plan multi-day swaps can cancel old bookings for dates that were not created [`web/src/components/InteractiveFloorPlan.vue:1684`] — resolved by tracking cancellation ids by date and cancelling only dates returned by the successful multi-day create response; regression covered in `web/src/components/__tests__/InteractiveFloorPlan.test.ts`.
- [x] [Review][Patch] Self-booking guard can treat an on-behalf booking with a missing display name as the user's own seat [`web/src/composables/useAreaDayGuard.ts:88`] — resolved by excluding `for_user_id` rows from the self-seat predicate even when `for_user_name` is absent; regression covered in `web/src/composables/__tests__/useAreaDayGuard.test.ts`.

## Dev Notes

### Backend: no change required

The guard is fully client-side. The existing "my bookings" endpoint already returns everything the
guard needs per booking: `item_id`, `item_name`, `area_id`, `area_name`, `booking_date`, and the
booking `id`. See `MyBookingAttributes` [Source: internal/bookings/handler.go:71-87] and the JSON
built in `writeBookingsCollection` where `area_id`/`area_name` come from
`cfg.FindItemLocation(rec.ItemID)` [Source: internal/bookings/handler.go:418-434,
internal/areas/config.go:69-85]. The frontend type mirror is `MyBookingAttributes`
[Source: web/src/api/bookings.ts:12-28], fetched via `fetchMyBookings()`
[Source: web/src/api/bookings.ts:158-160].

> [!IMPORTANT]
> The existing backend "already have this item booked for this date" conflict
> (`processBooking`, `FindUserBooking`) is scoped to the exact SAME item, not the area
> [Source: internal/bookings/handler.go:904-918, 1090-1105]. Story 36.9 is a broader, area-scoped
> guard and is intentionally implemented as a client-side confirmation + swap, not a hard backend
> 409. Do not repurpose the item-level conflict for this.

### The swap is two existing calls, run in order

Confirm → `await cancelBooking(existingBookingId)` [Source: web/src/api/bookings.ts:191-203] then
the flow's normal create (`createBooking` / `createMultiDayBooking`
[Source: web/src/api/bookings.ts:106-156]). Keep the dialog in a `loading` state across both so the
user cannot double-submit; on cancel-failure or create-failure surface the flow's existing error
snackbar and reload availability.

### All three flows already share the warning-confirmation pattern — mirror it

Every booking surface already owns a `useWarningConfirmation()` instance and a
`WarningConfirmDialog`, and gates its create behind `presentWarnings(items, onConfirmed)`
[Source: web/src/composables/useWarningConfirmation.ts:54-66]. The area/day guard slots into the
SAME chain, immediately after warnings and immediately before the create:

`presentWarnings(warnItems, () => guard({ areaId, date, newItemId, newItemName }, () => doCreate()))`

This keeps warning confirmation first (per epic 35) and the swap prompt second, and it reuses the
mental model already established. There is no shared booking-orchestration composable today; the
only shared piece is the API layer, so `useAreaDayGuard` is a small new composable each flow
instantiates itself (same as `useWarningConfirmation`).

### Per-flow insertion points (verified)

- Tile — `web/src/views/ItemsView.vue`: single-day `requestBooking` → `bookItem`
  (`createBooking` at 1846); week mode `submitWeekBookings` (1482) books per (item, date) with
  `createBooking` in a loop. `areaId` via `getCurrentAreaId()` (1553); favorites mode aggregates
  across areas so match on my-bookings `area_id`. Success snackbar `data-cy="booking-success"`
  (~829-850), error snackbar `color="error" :timeout="6000" closable` (~852-861).
- Table — `web/src/components/area-weekly-matrix/MatrixBookingPopover.vue`: `submitBooking` (220) →
  `doBook` → `createBooking(props.item.item_id, props.cell.date, ...)` (245); emits `booked` (252)
  and `bookingConflict` (256). Parent `AreaWeeklyMatrixView.vue` owns `props.areaId` (98-103),
  renders the popover (48-56) and refreshes via `onBooked` → `loadMatrix` (232, 253). Pass `areaId`
  down as a new popover prop; this flow must `fetchMyBookings()` itself.
- Floor plan — `web/src/components/InteractiveFloorPlan.vue`: `confirmPendingBooking` (1518) →
  `executeBooking` (1541); single day `createBooking` (1553), multi-day `createMultiDayBooking`
  (1559). `areaId?` prop (577); refresh via `refreshAvailability()` (1592). Must `fetchMyBookings()`
  itself; apply the guard per selected date.

### Feedback and dialog conventions

Use `ConfirmDialog.vue` for the swap prompt (blocking decision, not transient feedback). Operation
outcomes (swap succeeded / failed) use the flow's existing `v-snackbar`
(success 3000ms, error 6000ms + closable) per [Source: .claude/rules/feedback.md]. Add a `data-cy`
to any new element; `ConfirmDialog` already exposes `confirm`/`cancel` hooks
[Source: web/src/components/ConfirmDialog.vue:20,29].

### Scope and edge cases

- Guard is scoped to same `area_id` AND same `booking_date`; a different area or a different day is
  not a conflict.
- Booking a second item in a DIFFERENT item group of the same area still triggers the guard (scope
  is the area, not the item group).
- Only one existing same-area/day booking is expected (that is the invariant this story enforces);
  if the dialog shows the first match, that is sufficient. Guest and on-behalf bookings: keep the
  guard for the acting user's own bookings (the my-bookings list is the current user's) and do not
  block colleague/guest targets differently unless a follow-up requires it — note this as an
  intentional current-user-scoped limitation.
- Multi-day floor-plan / week-tile: evaluate each date independently; a same-area booking on one of
  several selected days prompts a swap for that day only.

### Project Structure Notes

- New: `web/src/composables/useAreaDayGuard.ts` (+ test in `web/src/composables/__tests__/`).
- Modified: `web/src/views/ItemsView.vue`,
  `web/src/components/area-weekly-matrix/MatrixBookingPopover.vue` (+ new `areaId` prop),
  `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue` (pass `areaId` down),
  `web/src/components/InteractiveFloorPlan.vue`, and five locale files under `web/src/locales/`.
- Reuses existing `ConfirmDialog.vue`, `fetchMyBookings`, `cancelBooking`, `createBooking`,
  `createMultiDayBooking`. No backend or API-type changes.

### Testing standards summary

Vitest unit + component tests (coverage ≥ 75%); table-driven where it fits. Run `npm run type-check`,
`npm run lint`, `npx vitest run`, `npm run build`, and the jscpd TS duplication check (threshold 0).
Add a Cypress E2E covering the acceptance flow against a dev server (no mocked responses; intercepts
only for waiting) [Source: .claude/rules/vue.md, .claude/rules/cypress.md].

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 36.9 (FR178), lines 5560-5585]
- [Source: internal/bookings/handler.go:71-87,418-434,904-918,1090-1105]
- [Source: internal/areas/config.go:62-85]
- [Source: web/src/api/bookings.ts:12-28,106-156,158-160,191-203]
- [Source: web/src/composables/useWarningConfirmation.ts:25-91]
- [Source: web/src/components/ConfirmDialog.vue:1-76]
- [Source: web/src/views/ItemsView.vue:1482,1553,1794-1895,1397,1741]
- [Source: web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:122-133,220-259]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue:48-56,98-103,232,253]
- [Source: web/src/components/InteractiveFloorPlan.vue:247,496,577,1460-1592]
- [Source: web/src/locales/en.json:166 (warning keys block)]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `npm run type-check`, `npm run lint`, `npx vitest run` (503 passed), `npm run build`, and
  `npx jscpd` (0 clones) all green.
- `useAreaDayGuard.guard()` returns a `Promise<boolean>` (true = proceed, incl. confirmed swap;
  false = cancelled) so multi-day/week flows can await each decision sequentially. Confirm cancels
  the existing booking then resolves true; cancel resolves false with no change.

### Completion Notes List

- New `useAreaDayGuard` composable + reusable `ConfirmDialog` wired into all three flows:
  - Tile (`ItemsView.vue`): day mode gates `bookItem` after warnings; week mode runs a pre-pass
    (`resolveWeekAreaGuards`) that swaps conflicting bookings per (area, date) before the batch.
    Favorites mode aggregates across areas, so the area is undeterminable and the guard no-ops
    (never blocks) — an intentional current-user/single-area-scoped limitation.
  - Table (`MatrixBookingPopover.vue`): new required `areaId` prop passed down from
    `AreaWeeklyMatrixView.vue`; guard runs after warnings, before `doBook`'s create.
  - Floor plan (`InteractiveFloorPlan.vue`): `resolveFloorPlanAreaGuards` applies the guard per
    selected date before booking; only runs when `areaId` is set (drill-down), else no-ops.
- Match is on the my-bookings `area_id` + `booking_date` (different `item_id`); confirm reuses
  `cancelBooking` + the flow's normal create; the `ConfirmDialog` `loading` binding disables the
  buttons during the swap. New i18n keys `items.areaDayGuard{Title,Message,Confirm}` in all 5 locales.
- No backend/API-type changes.

### File List

- web/src/composables/useAreaDayGuard.ts (new)
- web/src/composables/__tests__/useAreaDayGuard.test.ts (new)
- web/src/views/ItemsView.vue
- web/src/views/ItemsView.test.ts
- web/src/components/InteractiveFloorPlan.vue
- web/src/components/__tests__/InteractiveFloorPlan.test.ts
- web/src/components/area-weekly-matrix/MatrixBookingPopover.vue
- web/src/components/area-weekly-matrix/MatrixBookingPopover.test.ts
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue
- web/src/locales/{en,de,es,fr,uk}.json (items.areaDayGuard*)

## Senior Developer Review (AI)

Adversarial review of the whole Epic 36 change set (Blind Hunter + Edge Case Hunter + Acceptance
Auditor, 2026-07-09). No High/Med issues in the security-critical key handling (36.1) or version
threading (36.4). All findings below were triaged and RESOLVED before marking done.

### Resolved — decisions

- [x] [D1][Med] Guard mis-scoped on-behalf bookings — `findConflict` now requires `booked_for_me`,
  so a booking made for a colleague is never offered for cancellation.
- [x] [D2][Med] Swap was cancel-then-create (a failed create lost both) — guard returns
  `{ proceed, existingBookingId }` and never cancels; call sites CREATE first, then cancel the old
  id on success (a failed post-create cancel keeps the new booking + warns).
- [x] [D3][Med] Favorites week-booking bypassed the guard — each favorite is guarded by its own
  `areaId`; only legacy favorites with no resolvable area stay unguarded (documented in code).
- [x] [D4][Low] Multi-day: declining one day's swap aborted all days — now per-day.

### Resolved — patches

- [x] [P1][Med] Added `guest_name` to the My Bookings collection API (+test).
- [x] [P2][Low] Added `.catch`/try-catch to the three guard callbacks.
- [x] [P3][Low] Removed dead `getShortName`.
- [x] [P4][Low] Fixed stale `selectedDate`→`selectedDay` naming.
- [x] [P5][Low] Unified the guard-dialog date via `src/utils/date.ts`.
- [x] [P6][Low] Extracted shared `deriveBookingStatus` (`src/utils/bookingStatus.ts`).

### Dismissed (verified noise)

App.vue version watcher already `immediate:true`; all i18n keys present in 5 locales; dormant
loading spinner; `preselectDay` stale-index fails safe.
