# Story 37.2: Week Dropdown and Navigation Arrows in Floor Plan and Tile View

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want to switch and step through weeks directly in the floor plan and the tile view,
so that I can plan bookings without leaving the view I am in.

## Acceptance Criteria

1. Floor plan opened from any entry path (tile day mode, tile week mode, item view, drill-down) shows the same
   week dropdown as the tile view, with the same options (current week up to the advance-booking limit from
   `/api/v1/settings`).
2. Selecting another week in the floor plan dropdown reloads availability for that week and keeps the selected
   weekday (first workday, i.e. first non-past day, if that weekday is not available).
3. The tile view and the floor plan show a "previous week" and a "next week" arrow button next to the week
   dropdown; each click changes the week by exactly one week.
4. When the current week is selected, "previous week" is disabled.
5. When the last week allowed by `weeks_in_advanced` is selected, "next week" is disabled.
6. German UI: the floor plan header shows the calendar week as "KW 41" (translated label) instead of "Week 41";
   English keeps "Week 41".
7. Arrows and dropdown have accessible labels and `data-cy` attributes (`week-prev`, `week-next`); an E2E test
   covers stepping weeks in both views including both disabled bounds.

## Tasks / Subtasks

- [x] Task 1: Translate the week label in `useWeekSelector` (AC: #6)
  - [x] Add an optional translator to `useWeekSelector` (e.g. third parameter
        `t?: (key: string, named?: Record<string, unknown>) => string`, mirroring the optional-`t` pattern of
        `getWeekdayLabel`, `useWeekSelector.ts:60-68`). Replace the hardcoded `- Week ${weekNum}` in the option
        label (`useWeekSelector.ts:117`) with `t('weekSelector.weekNumber', { week: weekNum })`; fall back to
        `Week N` when no `t` is given so existing unit tests keep passing.
  - [x] Add locale keys in all five locales (`web/src/locales/{en,de,es,fr,uk}.json`), new section
        `weekSelector`: `weekNumber` (en "Week {week}", de "KW {week}", es "Semana {week}", fr "Semaine {week}",
        uk "Тиждень {week}"), `previousWeek` (en "Previous week", de "Vorherige Woche"), `nextWeek`
        (en "Next week", de "Nächste Woche"); translate es/fr/uk accordingly.
  - [x] Pass `t` from `useI18n()` in `ItemGroupsView.vue:411` and `ItemsView.vue:1218` so the computed labels
        react to locale changes.
- [x] Task 2: Create a shared `WeekSelector.vue` component (AC: #3, #4, #5, #7)
  - [x] New `web/src/components/WeekSelector.vue` (`<script setup lang="ts">`): props `options: WeekOption[]`,
        `label: string`, optional `selectDataCy` (default `week-selector`); `defineModel<string>()` for the
        selected ISO week (`YYYY-Www`).
  - [x] Render: `v-btn` icon `$chevronLeft` (`data-cy="week-prev"`, `:aria-label="t('weekSelector.previousWeek')"`),
        the existing `v-select` (`item-title="label"`, `item-value="value"`, `density="compact"`, `hide-details`),
        `v-btn` icon `$chevronRight` (`data-cy="week-next"`, `:aria-label="t('weekSelector.nextWeek')"`).
        Icons `$chevronLeft`/`$chevronRight` are already registered aliases (`web/src/plugins/vuetify.ts`).
  - [x] Step logic based on the option index: prev = `options[idx - 1]`, next = `options[idx + 1]`;
        `prevDisabled = idx <= 0`, `nextDisabled = idx < 0 || idx >= options.length - 1`. Never compute dates
        here; the options list is the single source of the allowed range.
  - [x] Export from `web/src/components/index.ts`.
- [x] Task 3: Use `WeekSelector` in the tile views (AC: #3, #4, #5)
  - [x] `ItemGroupsView.vue:13-23`: replace the bare `v-select` with `<WeekSelector v-model="selectedWeek"
        :options="weekOptions" :label="$t('itemGroups.calendarWeek')" />`. Keep `data-cy="week-selector"` on the
        inner select (existing tests rely on it: `ItemGroupsView.test.ts:246`, `availability.cy.ts`). Adapt the
        `.ig-week-selector` width rules (`ItemGroupsView.vue:609-613`, mobile block ~627) so select + arrows fit.
  - [x] `ItemsView.vue:53-64` (week mode only): same replacement; adapt `.booking-date-input`
        (`ItemsView.vue:2462-2466`, mobile ~2480). Day mode keeps the `DatePickerField` without arrows (see Open
        questions).
  - [x] Existing watchers already reload on week change (`ItemGroupsView.vue:584-590`, `ItemsView.vue:2382-2387`)
        and memorize via `setWeek`; no extra loading code needed.
- [x] Task 4: Week dropdown + arrows in `InteractiveFloorPlan` header (AC: #1, #3, #6)
  - [x] Add props `weekOptions: WeekOption[]` and `week: string` plus emit `update:week` (or `defineModel('week')`)
        to `InteractiveFloorPlan.vue:594-609`. Keep `weekDates` as the parent-computed input (it drives
        `weekdays`, `InteractiveFloorPlan.vue:766-779`).
  - [x] Render `<WeekSelector>` in `fp-header` (`InteractiveFloorPlan.vue:3-45`) with
        `select-data-cy="fp-week-selector"` and label `$t('itemGroups.calendarWeek')`. A distinct select
        `data-cy` avoids clashes with the page's `week-selector` behind the dialog.
  - [x] Header label (`InteractiveFloorPlan.vue:20`): show the selected option's label (translated, "KW 41") in
        every entry path. Either derive it inside the component from `weekOptions`/`week` and drop the
        `weekLabel` prop, or keep the prop and feed it the option label; in day mode it must no longer show
        `formatDisplayDate(selectedDate)` (`ItemsView.vue:1235-1240`).
  - [x] Check the header on compact viewports (`@media (max-width: 768px), (max-height: 500px)`,
        `InteractiveFloorPlan.vue:~2163`): selector must not overflow at 360 px; wrap it onto its own line if needed.
- [x] Task 5: Keep the weekday when the week changes (AC: #2)
  - [x] Today the `weekDates` watcher (`InteractiveFloorPlan.vue:1163-1182`) calls `initialLoad()` →
        `preselectDay()` (`:970-993`), which jumps to `selectedDay`/today/first-future and loses the weekday.
  - [x] Extend `preselectDay` with a `preserveIndex` mode: when only `weekDates` changed (not `floorPlan`,
        `itemGroupId`, `areaLevel` – the existing `resetView` flag), keep `selectedDayIndex` if
        `weekdays[idx]` exists and is not past; otherwise pick the first non-past day (fallback `0`).
  - [x] Drill-down must survive a week change: `drilledInto` is only reset when `resetView` is true – keep it so.
        Availability for the drilled group is reloaded through `refreshAvailability()` (`:1138-1146`).
  - [x] Note: if the preserved index equals the old index, the `selectedDayIndex` watcher (`:1184-1186`) does not
        fire; `initialLoad()` already calls `refreshAvailability()`, so no double load is needed.
- [x] Task 6: Wire every entry path (AC: #1, #2)
  - [x] Area tile view (`ItemGroupsView.vue:233-243`): bind `:week-options="weekOptions"` and
        `v-model:week="selectedWeek"`. The floor plan and the tile view share the week, so the tiles follow a
        change made in the floor plan (and `setWeek` memorizes it).
  - [x] Drill-down: handled inside the area-level floor plan (`handleAreaClick`, `InteractiveFloorPlan.vue:1373`;
        `drillBack`, `:1386`); the header selector stays visible in both levels – verify only.
  - [x] Item view, week mode (`ItemsView.vue:768-779`): `v-model:week="selectedWeek"`, `:week-options`,
        `:week-dates="selectedWeekDates"`.
  - [x] Item view, day mode: add a local `floorPlanWeek` ref, set to `getISOWeekString(getMondayOfWeek(selectedDate))`
        when the dialog opens (`showItemGroupFloorPlanDialog`, `ItemsView.vue:69`); derive
        `floorPlanWeekDates` from it (replacing `dayModeWeekDates`, `ItemsView.vue:1229-1234`). Keep passing
        `floorPlanSelectedDate` so the first open still lands on the tile day (story 36.6). Do not write the
        floor plan's week back into `selectedDate`.
  - [x] Reuse `selectedMonday` logic: extract the ISO-week-to-Monday parsing (`useWeekSelector.ts:126-137`) into an
        exported helper (e.g. `getMondayOfISOWeek(isoWeek)`) instead of duplicating it in `ItemsView` (jscpd gate).
- [x] Task 7: Unit / component tests (AC: #1-#7)
  - [x] `useWeekSelector.test.ts`: translated label via `t` ("KW 41"), fallback "Week N", new ISO-week helper.
  - [x] New `web/src/components/__tests__/WeekSelector.test.ts`: prev/next step by one option, prev disabled on
        first option, next disabled on last option, aria-labels present, emits `update:modelValue`.
  - [x] `InteractiveFloorPlan.test.ts`: renders `fp-week-selector`, `week-prev`, `week-next`; changing `weekDates`
        keeps the weekday index and requests availability for the new date; falls back to first non-past day;
        drill-down survives a week change. Update `mountComponent` defaults (`:221-235`) with the new props.
  - [x] `ItemGroupsView.test.ts` / `ItemsView.test.ts`: arrows present next to `week-selector`; ItemsView day mode
        passes the day's week to the floor plan.
- [x] Task 8: Cypress E2E (AC: #1-#7)
  - [x] New `web/cypress/e2e/week-navigation.cy.ts`. Freeze the clock to a Monday
        (`cy.clock(..., ['Date'])`) and observe `GET /api/v1/settings` (`cy.intercept(...).as('settings')`,
        no stub) to compute the expected option count `weeks_in_advanced + 1` from the response body.
  - [x] Tile view (area page and item page week mode): `week-prev` disabled initially; click `week-next` until
        disabled, assert the option count, the dropdown label changes by one week per click and
        `@availability` / `@listItems` fire; step back to the start and assert `week-prev` disabled again.
  - [x] Floor plan: open from item view day mode, item view week mode and area level (drill into a group);
        assert `fp-week-selector`, `week-prev`, `week-next` within the dialog (`cy.get('[data-cy="..."]')`
        scoped via `.within()` on `item-group-floor-plan-dialog` / `floor-plan-dialog`); step weeks, assert the
        items request date moves by 7 days and the selected weekday button stays active; both bounds disabled.
  - [x] German UI: switch locale (localStorage key `sithub_locale`, `useLocalePreference.ts:8`) and assert the floor
        plan header contains "KW"; English contains "Week".
- [x] Task 9: Quality gates
  - [x] `npm run type-check`, `npm run lint`, `npx vitest run --coverage` (>= 75 %), `npm run build`,
        `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0 --exit-code 1`,
        `npm run test:e2e -- --browser electron` against `npm run dev` + backend.

### Review Findings

- [x] [Review][Patch][P2] Prevent stale availability and loading completion from earlier week/drill-down requests overwriting the current view [web/src/components/InteractiveFloorPlan.vue:1050] — fixed with request IDs, captured-context checks and shared initialization generations; regression completes week responses in reverse order.
- [x] [Review][Patch][P2] Clamp restored and day-mode floor-plan weeks to the loaded booking horizon [web/src/composables/useWeekSelector.ts:157; web/src/views/ItemsView.vue:1235] — fixed centrally when options change and for the dialog model; tile day selection stays unchanged.
- [x] [Review][Patch][P3] Reject invalid numeric ISO weeks instead of silently rolling into a different week/year [web/src/composables/useWeekSelector.ts:36] — fixed by range validation and ISO round-trip; tests cover W00, W54 and nonexistent W53.

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md#Story 37.2` (FR181-FR183, `epics.md:5735-5775`, FR text
`epics.md:671-680`). Product decision (final): the floor plan always shows the week dropdown plus prev/next arrows on
every entry path; the "Week" label is translated ("KW" in German). Screenshot `private/img_41.png` shows the
target: header "Büro 2.EG 05.10.-11.10.2026 - Week 41" (arrow marks the untranslated "Week") with a
"Kalenderwoche" dropdown to the right; arrows go next to that dropdown.

This is a frontend-only story. No backend, API or OpenAPI changes.

### Current state (verified 2026-10-08)

- `useWeekSelector(showWeekends?, maxWeeks?)` (`web/src/composables/useWeekSelector.ts:102-150`) builds
  `weekOptions` = current week + `maxWeeks` weeks (default 7 until settings arrive), label
  `DD.MM.-DD.MM.YYYY - Week N` with hardcoded English "Week" (`:117`). `selectedWeek` is the ISO string
  `YYYY-Www`; `selectedMonday` parses it (`:126-137`); `selectedWeekDates` honors `showWeekends`.
- Advance limit: `bookings.weeks_in_advanced` (backend default 5, `internal/config/config.go:106`; local
  `sithub.toml` sets 2), served by `/api/v1/settings`. Both views load it non-blocking
  (`ItemGroupsView.vue:503-504`, `ItemsView.vue:2291-2292`) into `weeksInAdvanced`, which feeds
  `useWeekSelector`. Semantics match the backend and `maxBookingDate` (`ItemsView.vue:1002-1011`): the last
  allowed week is the last option. So "next disabled" = last option selected.
- Area tile view `ItemGroupsView.vue`: `v-select data-cy="week-selector"` (`:13-23`), floor plan dialog
  `data-cy="floor-plan-dialog"` (`:225-245`) passes `week-label` (option label), `week-dates`,
  `:area-level="true"`. Week restored from `useDateState` (`:412-418`), watcher reloads + memorizes (`:584-590`).
- Item view `ItemsView.vue`: `bookingMode` day/week toggle (`:28-37`), day `DatePickerField` (`:40-50`), week
  `v-select` (`:53-64`, only in week mode), floor plan dialog `item-group-floor-plan-dialog` (`:761-781`).
  Floor plan inputs from story 36.6: `dayModeWeekDates`, `floorPlanWeekDates`, `floorPlanWeekLabel` (day mode =
  formatted date, not a week), `floorPlanSelectedDate` (`:1226-1243`).
- `InteractiveFloorPlan.vue`: header `fp-header` (`:3-45`) shows only `title` + `weekLabel` text (`:20`) and
  zoom controls; weekday buttons `fp-day-*` (`:46-67`). Props (`:594-609`): `floorPlan, title, weekLabel,
  weekDates, itemGroupId, areaLevel?, areaId?, selectedDay?`. No week selector and no arrows exist anywhere.
- Drill-down: area-level floor plan sets `drilledInto` (`handleAreaClick`, `:1373`; `drillBack`, `:1386`); the
  `weekDates` watcher (`:1163-1182`) resets `drilledInto` only when floor plan / group / level change.
- Availability reload: `refreshAvailability()` (`:1138-1146`) → `loadAvailability()` (`:1015`,
  `fetchItems(groupId, selectedDate)`) or `loadAreaAvailability()` (`:1047`, per group). E2E sync point:
  `GET /api/v1/item-groups/*/items*` (alias `@listItems` in `cypress/support/flows.ts:36-44`).

### Entry paths into the floor plan (all must show the selector)

1. Area page tile/table toolbar → `area-floor-plan-btn` → area-level floor plan (`ItemGroupsView`).
2. Drill-down from (1) into an item group (inside `InteractiveFloorPlan`).
3. Item page in day mode → `item-group-floor-plan-btn` (`ItemsView`, week containing the selected day).
4. Item page in week mode → same button (`ItemsView`, selected week).

There is no other mount of `InteractiveFloorPlan` (`grep` over `web/src` confirms only these two views).

### Design guidance

- One component (`WeekSelector.vue`) for select + arrows, used in three places. It only moves between
  `options` entries; the parent stays the owner of the week state. This keeps the bounds identical to the
  dropdown (AC #1, #4, #5) and avoids duplicated date math (jscpd gate).
- Floor plan keeps receiving `weekDates` from the parent; it emits week changes upward. Do not call
  `useWeekSelector` or `fetchSettings` inside `InteractiveFloorPlan` – that would create a second, unsynced
  week state and a second settings request.
- Day-mode preselection from story 36.6 must keep working on first open: `selectedDay` wins on the initial load,
  index preservation applies only to subsequent week changes.
- Past days: on the current week, preserved weekdays that are already past fall back to the first non-past day
  (`weekdays[].past`, `:775`). The existing `preselectDay` uses `toISOString()` (UTC) for "today" (`:984`); do not
  copy that pattern into new code – compare local dates.
- Feedback: no snackbar needed (no user operation succeeds/fails here); loading uses the existing
  `initialLoading` spinner below the header so the selector stays visible (`.claude/rules/feedback.md`).
- Story 37.3 handles mobile line breaks; still verify the new header row at 360 px does not overflow.

### Testing standards summary

- Vitest + Vue Test Utils, Composition API with `<script setup>`, no `any` (`.claude/rules/vue.md`); coverage
  >= 75 %.
- Cypress (`.claude/rules/cypress.md`): `data-cy` selectors only, intercepts as sync points, no fixed waits,
  `should <user outcome>` test names, freeze dates with `cy.clock(..., ['Date'])`.
- CI E2E data (`.github/workflows/ci.yml:137-166`) has no floor plans. Product decision (2026-10-08): no
  stubs. Extend the E2E seed with a small floor plan for one item group (SVG under the floor-plan location the
  config requires, plus positions seeded the way the app stores them) in BOTH `run-all-tests.sh` (`TEST_SPACES`
  heredoc) and `.github/workflows/ci.yml`. Story 37.3 adds an area "Büro 2.EG" / item group "Cube 1" to the same
  seed; reuse it for the floor plan if convenient. Existing specs that count areas must still pass; intercepts
  only as sync points.

### Project Structure Notes

- New: `web/src/components/WeekSelector.vue`, `web/src/components/__tests__/WeekSelector.test.ts`,
  `web/cypress/e2e/week-navigation.cy.ts`.
- Modified: `web/src/composables/useWeekSelector.ts` (+ test), `web/src/components/InteractiveFloorPlan.vue`
  (+ test), `web/src/views/ItemGroupsView.vue` (+ test), `web/src/views/ItemsView.vue` (+ test),
  `web/src/components/index.ts`, `web/src/locales/{en,de,es,fr,uk}.json`.
- Recent commits (`74529fb`, `025e1e0`) are backend refactors and dependency updates; nothing overlaps this story.
  Story 37.1 (error messages) touches `InteractiveFloorPlan.vue` and `ItemsView.vue` error paths – expect merge
  proximity, no functional overlap.

### Open questions / assumptions

- Day mode in the item view has no week dropdown today. Assumption: the tile view keeps its date picker without
  arrows in day mode; arrows appear wherever a week dropdown exists (area page, item page week mode, floor plan).
- In item-view day mode, a week change inside the floor plan stays local to the dialog and does not change the
  tile view's `selectedDate`. In week mode and on the area page, the floor plan and tile view share
  `selectedWeek` (changes propagate and are memorized). Recommended default; revisit if users expect sync.
- "First workday if the weekday is not available" is interpreted as the first non-past day of the target week.
- If a memorized week lies beyond the limit once settings load, the selector has no matching option; recommend
  clamping `selectedWeek` to the last option after `weeksInAdvanced` is set (small guard in both views).
- The dropdown option labels also use the translated week token ("KW 41"), not only the header.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.2 (lines 5735-5775), FR181-FR183 (lines 671-680)]
- [Source: private/img_41.png]
- [Source: web/src/composables/useWeekSelector.ts:60-68,102-150]
- [Source: web/src/components/InteractiveFloorPlan.vue:3-67,594-613,766-783,970-993,1015,1047,1138-1188,1373-1388]
- [Source: web/src/views/ItemGroupsView.vue:13-23,225-245,411-418,503-504,584-590,609-613]
- [Source: web/src/views/ItemsView.vue:28-71,761-781,1002-1011,1211-1243,2291-2292,2355-2387,2462-2484]
- [Source: web/src/plugins/vuetify.ts (icon aliases chevronLeft/chevronRight)]
- [Source: web/cypress/e2e/floor-plan-booking.cy.ts, web/cypress/support/flows.ts:36-44]
- [Source: _bmad-output/implementation-artifacts/36-6-floor-plan-opens-for-the-selected-day-week.md]
- [Source: .claude/rules/vue.md, cypress.md, feedback.md, markdown.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 5.5 (claude-opus-5-5)

### Debug Log References

- RED confirmed for the translated week token and `getMondayOfISOWeek` (4 failing) before implementation;
  `InteractiveFloorPlan` "keeps the weekday" test fails against the old `preselectDay` (jumps to Monday).
- ItemsView floor-plan test first used `futureDay()` (today + 3 = Sunday, hidden without weekends); switched to
  next week's Wednesday for a weekday-independent test.
- E2E: shared `@listItems` / `@availability` aliases were satisfied by initial-load requests; replaced by
  per-request intercepts keyed on `week=` / `date=` query (`nextAvailabilityFor`, `nextItemsFor`).

### Completion Notes List

- `useWeekSelector(showWeekends?, maxWeeks?, t?)`: optional `WeekTranslator`; option label uses
  `weekSelector.weekNumber` ("Week 41" / "KW 41"), English fallback without `t`. ISO-week parsing extracted
  into exported `getMondayOfISOWeek` (used by `selectedMonday` and `ItemsView`). Both views pass `t`.
- New `weekSelector` locale section (`weekNumber`, `previousWeek`, `nextWeek`) in en/de/es/fr/uk.
- New `WeekSelector.vue` (`defineModel`, props `options`, `label`, `selectDataCy`): prev/next arrows
  (`week-prev`/`week-next`, translated `aria-label`) move strictly between neighbouring options; disabled on the
  first/last option or when the value is not an option. Exported from `components/index.ts`.
- Tile views: `ItemGroupsView` and `ItemsView` (week mode) use `WeekSelector`; `data-cy="week-selector"` kept
  on the inner select. Width rules widened to 320px (new `.booking-week-input` in ItemsView), 100% on mobile.
  Day mode keeps the date picker without arrows (open-question default).
- `InteractiveFloorPlan`: `weekLabel` prop replaced by `weekOptions` + `v-model:week`; header shows the selected
  option label (`fp-week-label`, translated) and a `WeekSelector` (`fp-week-selector`). A week-only change
  (no floor plan / group / level change) keeps the weekday via `preserveSelectedDay` (fallback: first non-past
  day) and keeps the drill-down. Compact layout puts the selector on its own full-width line.
- Wiring: area page shares `selectedWeek` with the floor plan; item page week mode shares `selectedWeek`; item
  page day mode uses a dialog-local week (initialised from the selected day on open), does not write back to
  `selectedDate`, and still passes `selectedDay` for the first open (story 36.6).
- E2E seed (product decision, no stubs): area "Büro 2.EG" (`buero_2eg`, `e2e-area.svg`) with item group "Cube 1"
  (`cube_1`, `e2e-room.svg`, desks `cube_1_desk_1/2`) added to `run-all-tests.sh` and `ci.yml`; SVGs live in
  `web/cypress/fixtures/floor-plans/` (CI points `floor_plans` at it, the script copies it into its data dir).
  Positions are seeded idempotently by the spec through the admin API (`POST /api/v1/floor-plan-positions`).
- Not implemented (open-question recommendation, not a task): clamping a memorized week beyond the limit to the
  last option. The arrows are both disabled in that case, the dropdown still allows picking a valid week.
- Tests: Vitest 58 files / 587 tests pass, coverage statements 81.46 %, lines 82.88 %. Cypress full suite 77/77
  (new `week-navigation.cy.ts` 7 tests incl. both bounds in tile and floor-plan views, drill-down, weekday
  preservation, 360px fit, German "KW"). type-check, ESLint, jscpd (0 clones), build pass.

### File List

- `web/src/composables/useWeekSelector.ts`
- `web/src/composables/useWeekSelector.test.ts`
- `web/src/components/WeekSelector.vue` (new)
- `web/src/components/__tests__/WeekSelector.test.ts` (new)
- `web/src/components/index.ts`
- `web/src/components/InteractiveFloorPlan.vue`
- `web/src/components/__tests__/InteractiveFloorPlan.test.ts`
- `web/src/views/ItemGroupsView.vue`
- `web/src/views/ItemGroupsView.test.ts`
- `web/src/views/ItemsView.vue`
- `web/src/views/ItemsView.test.ts`
- `web/src/locales/en.json`
- `web/src/locales/de.json`
- `web/src/locales/fr.json`
- `web/src/locales/es.json`
- `web/src/locales/uk.json`
- `web/cypress/e2e/week-navigation.cy.ts` (new)
- `web/cypress/fixtures/floor-plans/e2e-area.svg` (new)
- `web/cypress/fixtures/floor-plans/e2e-room.svg` (new)
- `run-all-tests.sh`
- `.github/workflows/ci.yml`
- `_bmad-output/implementation-artifacts/sprint-status.yaml`

## Change Log

- 2026-10-08: Implemented story 37.2. Shared `WeekSelector` with prev/next arrows in area/item tile views and
  the floor plan header, translated week label ("KW"), weekday/drill-down preservation on week change, E2E seed
  with floor plans, unit and E2E tests. Status -> review.

## Code Review Record

- Date: 2026-10-08. Parallel Blind Hunter, Edge Case Hunter and Acceptance Auditor reviews; actionable findings auto-fixed as requested.
- Validation: 600 Vitest tests pass in 58 files; statements 81.56%, lines 83.01%. All 12 story Cypress tests pass against the Vite dev server and an isolated real backend; week/mobile specs rerun after the final loading fix (11/11).
- Type-check, ESLint, production build, TypeScript duplication check (0 clones), Go tests and golangci-lint (0 issues) pass. No unresolved review findings.
- 2026-10-08: Code review complete; status -> done.
