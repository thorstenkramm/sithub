# Story 36.2: Desktop Table View for My Bookings

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a desktop user,
I want a comprehensive table view of My Bookings with a toggle to tiles,
so that I can scan my bookings efficiently on a large screen.

## Acceptance Criteria

1. On a desktop viewport with no prior choice, "My Bookings" shows the TABLE view by default.
2. On a mobile/narrow viewport with no prior choice, "My Bookings" shows the TILE view by default.
3. A tile/table toggle persists the user choice in localStorage and restores it on the next visit,
   overriding the viewport default.
4. Table view presents a scannable layout with relevant columns (date, area/item, status, and any
   on-behalf information).

## Tasks / Subtasks

- [ ] Task 1: Add a per-view persistence composable for My Bookings (AC: #1, #2, #3)
  - [ ] Create `web/src/composables/useMyBookingsViewPreference.ts` modeled on
        `useAreaViewPreference.ts` (`web/src/composables/useAreaViewPreference.ts:1-60`) but with a
        SINGLE global preference (no per-area keying). Use a flat storage key
        `sithub_my_bookings_view` via `getSafeLocalStorage()`
        (`web/src/composables/storage.ts:1-10`).
  - [ ] `load(isDesktop: boolean)`: when `!isDesktop` return `'cards'` (AC #2); on desktop read the
        stored value and return `'table'` or `'cards'`, defaulting to `'table'` when NOTHING is
        stored (AC #1) — this is the key difference from `useAreaViewPreference`, whose desktop
        default is `'cards'`. A stored value wins over the viewport default (AC #3).
  - [ ] `save(view: AreaView)`: persist `'table'`/`'cards'` to the flat key so it is restored on next
        visit (AC #3). Reuse the `AreaView` type (`'cards' | 'table'`) from
        `useAreaViewPreference.ts:6` (export/import it — do not redefine).
  - [ ] Add a Vitest spec `useMyBookingsViewPreference.test.ts` mirroring
        `web/src/composables/useAreaViewPreference.test.ts` (SSR-safe no-storage path, corrupted-JSON
        path, mobile forces cards, desktop-empty defaults to table, stored value overrides).
- [ ] Task 2: Add viewport detection + view-switch toggle to `MyBookingsView.vue` (AC: #1, #2, #3)
  - [ ] Add `isCompactViewport` ref and an `updateViewport()` using `window.matchMedia('(max-width:
        768px)')`, mirroring `web/src/views/ItemGroupsView.vue:465-476`; register/cleanup a `resize`
        listener (`onMounted`/`onUnmounted`) as in `ItemGroupsView.vue:482-485`.
  - [ ] Wire `useMyBookingsViewPreference` — destructure `{ activeView, load, save }` and add a
        `toggleView(val: boolean | null)` that calls `save(val ? 'table' : 'cards')`, mirroring
        `ItemGroupsView.vue:318-324` (minus the areaId param).
  - [ ] In `onMounted`, call `updateViewport()` then `load(!isCompactViewport.value)` BEFORE (or after)
        the existing `fetchMe`/`loadBookings` flow (`MyBookingsView.vue:148-165`), matching
        `ItemGroupsView.vue:482-484`.
  - [ ] Render the Tiles/Table `v-switch` cluster from `ItemGroupsView.vue:30-67` verbatim in
        pattern: label spans bound to `activeView`, a disabled+tooltip switch when `isCompactViewport`
        (tooltip text `bookings.viewTableDesktopOnly`) and an enabled switch otherwise bound to
        `toggleView`. Place it in the page header area, above the bookings list
        (`MyBookingsView.vue:2-52`). Keep `data-cy="view-switch"` and
        `data-cy="view-switch-container"`.
- [ ] Task 3: Render the table view (AC: #4)
  - [ ] Show the existing tile grid (`MyBookingsView.vue:41-52`) only when `activeView === 'cards'`;
        render a table when `activeView === 'table'`. Use a Vuetify `v-data-table` (self-contained,
        no new component needed) over a hand-rolled table.
  - [ ] Columns grounded in `MyBookingAttributes` (`web/src/api/bookings.ts:12-28`): Date
        (`booking_date`, formatted like `BookingHistoryView.vue:81` / `BookingCard` `formattedDate`),
        Item (`item_name`), Area/Group (`area_name` + `item_group_name`), Status, and an
        On-behalf/Guest column. Keep it scannable — one row per booking.
  - [ ] Status column: reuse `StatusChip.vue` with the same status derivation used by `BookingCard`
        (`web/src/components/BookingCard.vue:15-32`): `is_guest` -> `guest`; else `booked_for_me` ->
        `booked-for-me`; else `booked_by_user_id && !booked_for_me` -> `on-behalf`. Valid statuses are
        in `StatusChip.vue:16`.
  - [ ] On-behalf/Guest column: surface `guest_name` (`bookings.guest`) and, for on-behalf, the
        `booked_by_user_name` (`bookings.bookedBy`). NOTE: the "On behalf of <first> <last>" wording
        (FR168) is delivered by Story 36.3 — here just wire the existing data/labels so the column
        exists; 36.3 refines the text in both tile and table views.
  - [ ] Provide a cancel affordance per row (error-colored button, `data-cy="cancel-btn"`) that emits
        into the SAME `handleCancelBooking` flow already used by the tiles
        (`MyBookingsView.vue:115-146`) so the confirm dialog + snackbar behavior is unchanged.
  - [ ] Add i18n keys for column headers under the `bookings` block in all five locales
        (`web/src/locales/{en,de,fr,es,uk}.json`, `bookings` block at en.json:171-192). Reuse the
        existing shared `itemGroups.viewTiles`/`viewTable`/`viewTableDesktopOnly` keys
        (`en.json:95-97`) for the toggle labels — do not duplicate them.
- [ ] Task 4: Tests (AC: #1-#4)
  - [ ] Extend `web/src/views/MyBookingsView.test.ts` (structure at lines 16-116): default view is
        table on desktop with empty storage (stub `matchMedia` to non-matching); default is tiles on
        narrow viewport (matching `max-width: 768px`); toggling persists to `localStorage` and is
        restored on remount overriding the viewport default; table renders the expected columns and a
        StatusChip; cancel from a table row runs the existing confirm+snackbar path
        (`MyBookingsView.test.ts:188-204`).
  - [ ] Add/adjust a Cypress E2E under `web/cypress/e2e` covering the AC user journey (login -> My
        Bookings shows table on desktop, toggle to tiles persists) using `cy.login()` and network
        waits per `.claude/rules/cypress.md`.
  - [ ] Run `npm run type-check`, `npm run lint`, `npx vitest run`, `npm run build`; keep
        `npx jscpd --pattern "**/*.ts"` at threshold 0 (factor the composable to avoid duplicating
        `useAreaViewPreference`).

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md` — Epic 36 / Story 36.2 (lines 5388-5414),
FR167 (`epics.md:617-619, 912`). The view-switch UI, persistence composable, and viewport detection
already exist for the Areas/Item-Groups matrix and should be reused as the pattern.

### Current My Bookings view (tiles only)

`web/src/views/MyBookingsView.vue` renders a single tile grid today: `BookingCard` in a
`.card-grid`, driven by `bookings` (`MyBookingsView.vue:41-52`). Data comes from
`fetchMyBookings()` (`web/src/api/bookings.ts:158-160`) returning
`CollectionResponse<MyBookingAttributes>`. Cancel flow: `handleCancelBooking` opens a
`ConfirmDialog`, `confirmCancelBooking` calls `cancelBooking` then reloads, with success snackbar and
error `v-alert` (`MyBookingsView.vue:54-66, 115-146`). This cancel plumbing must be reused unchanged
by the table rows.

### Persistence + toggle precedent (Item Groups / Areas matrix)

`useAreaViewPreference` is the closest precedent (`web/src/composables/useAreaViewPreference.ts`):

- Type `AreaView = 'cards' | 'table'` (line 6).
- `load(areaId, isDesktop)` forces `'cards'` when `!isDesktop`, reads a JSON map from
  `sithub_area_view`, and defaults to `'cards'` (lines 15-39).
- `save(areaId, view)` writes the JSON map, deleting the entry for the `'cards'` default (lines
  41-57).

For My Bookings we need a GLOBAL (not per-area) preference and a `'table'` desktop default. Create a
sibling composable with a flat string key `sithub_my_bookings_view` and the inverted desktop default.
Everything else (SSR guard via `getSafeLocalStorage`, corrupted-data fallthrough) copies the pattern.

The toggle UI and wiring live in `web/src/views/ItemGroupsView.vue`:

- Template: Tiles label + disabled/tooltip switch when compact + enabled switch bound to `toggleView`
  + Table label (lines 30-67); `data-cy="view-switch-container"` and `data-cy="view-switch"`.
- Script: `const { activeView, load, save } = useAreaViewPreference()` and
  `toggleView` (lines 318-324).
- Conditional render: matrix shown `v-else-if="activeView === 'table'"` (line 141), cards otherwise.

> [!NOTE]
> The matrix table itself is a hand-rolled `<table class="matrix-table">`
> (`AreaWeeklyMatrixView.vue:13-16`) because it is a 2D availability grid. My Bookings is a flat list,
> so a Vuetify `v-data-table` is the appropriate, simpler choice — do not reuse the matrix table.

### Viewport detection precedent

Use `window.matchMedia('(max-width: 768px)')` for narrow detection, matching
`ItemGroupsView.vue:465-476` (`updateViewport`) and its `resize` listener registration
(`ItemGroupsView.vue:482-485`). Guard `typeof window.matchMedia !== 'function'` for SSR/JSDOM as that
code does. The 768px breakpoint is the project's "narrow" convention (also used in
`InteractiveFloorPlan.vue:1689` and `ItemGroupsView.vue:472`). This uses `matchMedia` rather than
Vuetify `useDisplay` for consistency with existing views; do not introduce `useDisplay`.

### Columns and data grounding (AC #4)

`MyBookingAttributes` (`web/src/api/bookings.ts:12-28`) provides every field needed with no API
change: `booking_date`, `item_name`, `item_group_name`, `area_name`, `booked_by_user_name`,
`booked_for_me`, `is_guest`, `guest_name`, `note`. Status derivation must match `BookingCard`
(`BookingCard.vue:15-32`) so tile and table agree; `StatusChip` handles `guest`, `booked-for-me`,
`on-behalf` (`StatusChip.vue:16, 35-36`), labels `status.onBehalf` etc. Date formatting mirrors
`BookingHistoryView.vue:81` / `BookingCard`'s `formattedDate`.

### Relationship to sibling stories

Story 36.3 (Named On-Behalf Bookings, `epics.md:5416`) refines the on-behalf TEXT ("On behalf of
<first> <last>") in BOTH tile and table views. Keep 36.2's on-behalf column thin (existing
`booked_by_user_name` / `guest_name`) so 36.3 can layer wording on top without rework. Do not
implement FR168 wording here.

### Project Structure Notes

- New: `web/src/composables/useMyBookingsViewPreference.ts` (+ `.test.ts`).
- Modified: `web/src/views/MyBookingsView.vue` (+ `MyBookingsView.test.ts`), the five
  `web/src/locales/*.json` files (new column-header keys), and a Cypress E2E spec.
- Reuses: `useAreaViewPreference` `AreaView` type, `getSafeLocalStorage`, `StatusChip`,
  `ConfirmDialog`, and existing shared `itemGroups.view*` i18n keys.
- No backend/API change.

### Testing standards summary

Vitest component tests already exist for this view (`MyBookingsView.test.ts`); extend them and add the
composable spec (table-driven, `require`/`assert` equivalents via Vitest `expect`). Stub
`window.matchMedia` per case to drive the viewport default. Coverage must stay >= 75%
(`.claude/rules/vue.md`). Run type-check, ESLint, jscpd (TS threshold 0), build. Add a Cypress E2E for
the AC journey using `cy.login()` and network-based waits, no fixed `cy.wait(ms)`
(`.claude/rules/cypress.md`). Follow `.claude/rules/feedback.md` — cancel feedback stays a snackbar.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 36.2 / FR167 :5388-5414, :617-619]
- [Source: web/src/views/MyBookingsView.vue:2-52,54-66,115-146,148-165]
- [Source: web/src/composables/useAreaViewPreference.ts:1-60]
- [Source: web/src/composables/storage.ts:1-10]
- [Source: web/src/views/ItemGroupsView.vue:30-67,141,318-324,465-476,482-485]
- [Source: web/src/api/bookings.ts:12-28,158-160,191-203]
- [Source: web/src/components/BookingCard.vue:15-32]
- [Source: web/src/components/StatusChip.vue:16,35-36]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue:13-16]
- [Source: web/src/locales/en.json:95-97,171-192]
- [Source: .claude/rules/vue.md, .claude/rules/cypress.md, .claude/rules/feedback.md]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
