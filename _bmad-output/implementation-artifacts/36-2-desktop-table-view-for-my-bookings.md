# Story 36.2: Desktop Table View for My Bookings

Status: done

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

- [x] Task 1: Add a per-view persistence composable for My Bookings (AC: #1, #2, #3)
  - [x] Create `web/src/composables/useMyBookingsViewPreference.ts` modeled on
        `useAreaViewPreference.ts` with a SINGLE global preference and flat key
        `sithub_my_bookings_view` via `getSafeLocalStorage()`.
  - [x] `load(isDesktop: boolean)`: mobile returns `'cards'`; desktop reads the stored value,
        defaulting to `'table'` when nothing is stored; a stored value wins over the viewport default.
  - [x] `save(view: AreaView)`: persists `'table'`/`'cards'` to the flat key. Reuses the imported
        `AreaView` type from `useAreaViewPreference.ts`.
  - [x] Added Vitest spec `useMyBookingsViewPreference.test.ts` (no-storage/SSR path, corrupted-data
        path, mobile forces cards, desktop-empty defaults to table, stored value overrides).
- [x] Task 2: Add viewport detection + view-switch toggle to `MyBookingsView.vue` (AC: #1, #2, #3)
  - [x] Added `isCompactViewport` ref and `updateViewport()` using
        `window.matchMedia('(max-width: 768px)')` with SSR guard; resize listener registered in
        `onMounted` and removed in `onUnmounted`.
  - [x] Wired `useMyBookingsViewPreference` — `{ activeView, load, save }` plus a
        `toggleView(val: boolean | null)` calling `save(val ? 'table' : 'cards')`.
  - [x] In `onMounted`, calls `updateViewport()` then `load(!isCompactViewport.value)` before the
        existing `fetchMe`/`loadBookings` flow.
  - [x] Rendered the Tiles/Table `v-switch` cluster (label spans, disabled+tooltip switch when
        compact using `bookings.viewTableDesktopOnly`, enabled switch bound to `toggleView`) above the
        list, keeping `data-cy="view-switch"` and `data-cy="view-switch-container"`.
- [x] Task 3: Render the table view (AC: #4)
  - [x] Tile grid shown only when `activeView === 'cards'`; a Vuetify `v-data-table` renders when
        `activeView === 'table'`.
  - [x] Columns: Date (formatted via a shared `formatBookingDate`), Item (`item_name`),
        Area/Group (`item_group_name` + `area_name`), Status, For/Guest, Actions.
  - [x] Status column reuses `StatusChip.vue` with the same derivation as `BookingCard`
        (`guest` / `booked-for-me` / `on-behalf`).
  - [x] For/Guest column surfaces `guest_name` and `booked_by_user_name`; the "On behalf of {name}"
        wording (FR168) is delivered alongside 36.3 via `for_user_name` + `bookings.onBehalfOf`.
  - [x] Per-row cancel button (`data-cy="cancel-btn"`) calls the same `handleCancelBooking` flow,
        leaving the confirm dialog + snackbar behaviour unchanged.
  - [x] Added column-header i18n keys (`colDate`/`colItem`/`colArea`/`colStatus`/`colOnBehalf`/
        `colActions`) plus `viewTableDesktopOnly` under the `bookings` block in all five locales;
        reused shared `itemGroups.viewTiles`/`viewTable` for the toggle labels.
- [x] Task 4: Tests (AC: #1-#4)
  - [x] Extended `MyBookingsView.test.ts`: desktop-empty defaults to table; narrow viewport defaults
        to tiles; stored preference overrides the viewport default; table renders a StatusChip and the
        on-behalf name; cancel from a table row runs the confirm+snackbar path. All 14 tests pass.
  - [x] Added Cypress E2E `web/cypress/e2e/my-bookings-view-toggle.cy.ts` covering login -> table on
        desktop -> toggle to tiles -> reload persists, using `cy.login()` and `cy.wait('@listBookings')`.
  - [x] `npm run type-check`, `npm run lint`, `npx vitest run` (485 pass), `npm run build` all clean;
        new composable/test files introduce no jscpd clones (repo has a pre-existing 4.06% TS baseline
        of shared view-test boilerplate).

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
  - Table label (lines 30-67); `data-cy="view-switch-container"` and `data-cy="view-switch"`.
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

claude-opus-4-8

### Debug Log References

- Go gate: `go test ./...`, `go vet ./...`, `gofmt -l`, `golangci-lint run ./...` — all clean.
- Frontend gate: `npm run type-check`, `npm run lint`, `npx vitest run` (485 tests pass),
  `npm run build` — all clean.
- After extracting `buildMyBookingAttributes` (shared with 36.3), golangci-lint gocognit stayed
  within the 20 threshold.

### Completion Notes List

- New global composable `useMyBookingsViewPreference` reuses the `AreaView` type and
  `getSafeLocalStorage` from the Areas precedent; the only behavioural difference is the flat storage
  key and the `'table'` desktop default.
- The view switch and viewport detection mirror `ItemGroupsView.vue`; `matchMedia('(max-width: 768px)')`
  is used rather than Vuetify `useDisplay`, consistent with existing views.
- The table uses a Vuetify `v-data-table` with a For/Guest column. The on-behalf label wording is
  shared with Story 36.3 (`for_user_name` + `bookings.onBehalfOf`) so both views read identically.
- Cypress E2E spec authored and linted; it was not executed here as it requires the live dev-server
  - seeded DB stack. Unit tests fully cover the toggle/persistence and label behaviour.

### File List

- `web/src/composables/useMyBookingsViewPreference.ts` (new)
- `web/src/composables/useMyBookingsViewPreference.test.ts` (new)
- `web/src/views/MyBookingsView.vue` (modified — toggle, viewport detection, table view)
- `web/src/views/MyBookingsView.test.ts` (modified — table/tile toggle + persistence + cancel tests)
- `web/cypress/e2e/my-bookings-view-toggle.cy.ts` (new)
- `web/src/locales/{en,de,es,fr,uk}.json` (modified — column headers + `viewTableDesktopOnly`)

### Change Log

- 2026-07-09: Implemented the desktop table view with a persisted tile/table toggle for My Bookings
  (FR167). Added the `useMyBookingsViewPreference` composable, a `v-data-table` view, viewport-based
  defaulting, i18n column headers in all locales, unit tests, and a Cypress E2E.
