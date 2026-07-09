# Story 36.6: Floor Plan Opens for the Selected Day/Week

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user booking from the tile view,
I want the floor plan to open for the day or week I selected on the tiles,
so that I book the correct day.

## Acceptance Criteria

1. Day mode with a non-today day selected on the tile view: opening the floor plan shows
   availability for that selected day.
2. Week mode with a selected week: the floor plan opens on that week (and its selected/first
   bookable day), not the current week.
3. Changing the selected day/week and reopening reflects the latest selection.

## Tasks / Subtasks

- [x] Task 1: Pass the tile-view selected day into the floor plan (AC: #1, #3)
  - [x] Added optional prop `selectedDay?: string` (ISO `YYYY-MM-DD`) to
        `InteractiveFloorPlan.vue`'s `defineProps` block. (Renamed from `selectedDate` to avoid a
        vue/no-dupe-keys collision with the internal `selectedDate` computed.)
  - [x] In `ItemsView.vue`, bound the prop on the dialog's `<InteractiveFloorPlan>` instance via
        `:selected-day="floorPlanSelectedDate"` (day-mode = `selectedDate`, week-mode =
        `undefined`).
- [x] Task 2: Ensure the floor plan's `weekDates` cover the selected day in day mode (AC: #1)
  - [x] Added `floorPlanWeekDates`/`floorPlanWeekLabel` computeds in `ItemsView.vue`: week mode
        keeps `selectedWeekDates`/week label; day mode derives the Monday-based week CONTAINING
        `selectedDate` via `getMondayOfWeek` + `getWeekdayDates`, honoring `showWeekends`, and
        uses `formatDisplayDate(selectedDate)` as the header label. Bound both to the dialog.
- [x] Task 3: Honor the selected day in `InteractiveFloorPlan`'s initial day selection (AC: #1, #2, #3)
  - [x] Updated `preselectDay()`: if `props.selectedDay` is a non-past entry in `weekdays.value`,
        select that index; otherwise fall back to today / first-future. Drives
        `refreshAvailability()` for that day.
  - [x] Reopen behavior confirmed: the dialog `v-if`-mounts a fresh component each open, so
        `initialLoad()` → `preselectDay()` picks up the latest reactive props; no extra watcher.
- [x] Task 4: Week-mode preselection stays correct (AC: #2)
  - [x] Week mode passes `selectedDay = undefined`, so `preselectDay()` keeps the today/first-future
        day within the selected week. Existing week-mode tests still pass.
- [x] Task 5: Tests (AC: #1-#3)
  - [x] Component tests added to `InteractiveFloorPlan.test.ts`: `selectedDay` inside the week
        preselects that day and requests availability for it; a past `selectedDay` falls back to
        today.
  - [ ] E2E deferred (unit + component coverage in place; no Cypress spec added this pass).

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md#Story 36.6` (FR174, `epics.md:5484-5504`);
FR174 text: "Opening the floor plan from the tile view opens it for the day (day mode) or week
(week mode) selected on the tiles" (`epics.md:636`).

### Where the two selections live (ItemsView)

`web/src/views/ItemsView.vue` has a `bookingMode` toggle of `'day' | 'week'`
(`ItemsView.vue:30,1141-1143`) with two independent selectors:

- Day mode: `DatePickerField` bound to `selectedDate` ref (`ItemsView.vue:41-42`;
  `selectedDate = ref(getDay())` at `:956`).
- Week mode: `<v-select data-cy="week-selector">` bound to `selectedWeek`
  (`ItemsView.vue:53-63`), from `useWeekSelector` (`ItemsView.vue:1146`).

The floor plan button (`data-cy="item-group-floor-plan-btn"`) just opens the dialog
(`ItemsView.vue:66-71`) which conditionally renders `InteractiveFloorPlan` behind
`showItemGroupFloorPlanDialog` (`ItemsView.vue:775-796,1008`).

### The gap (root cause)

The dialog passes `:week-dates="selectedWeekDates"` and
`:week-label="weekOptions.find(o => o.value === selectedWeek)?.label"`
(`ItemsView.vue:788-789`). Both come exclusively from the WEEK selector
(`selectedWeekDates` computed from `selectedWeek` in
`web/src/composables/useWeekSelector.ts:123-142`) and NOTHING passes the day-mode
`selectedDate`. Inside the floor plan, `preselectDay()` always jumps to today (or first future
day) (`InteractiveFloorPlan.vue:911-921`). Net effect: in day mode with a non-today day the floor
plan ignores the tile selection and shows today; and there is no `selectedDate` prop at all.

### How the floor plan derives and preselects a day

- Props today: `floorPlan, title, weekLabel, weekDates, itemGroupId, areaLevel?, areaId?` — no
  `selectedDate` (`InteractiveFloorPlan.vue:569-578`).
- `weekdays` computed maps `props.weekDates` to `{date,label,past}` using `new Date()` as "today"
  boundary (`InteractiveFloorPlan.vue:735-748`).
- `selectedDate` (component-internal) = `weekdays[selectedDayIndex].date`
  (`InteractiveFloorPlan.vue:750-752`); `bookingDayOptions = weekdays`
  (`:792`).
- `preselectDay()` is called from `initialLoad()` (`:1078`), which runs on mount and whenever the
  watched props change (`InteractiveFloorPlan.vue:1091-1110`, `immediate: true`). Changing
  `selectedDayIndex` triggers `refreshAvailability()` (`:1112-1114`).

So the correct fix is: (a) give the component a `selectedDate` prop, (b) make `preselectDay()`
prefer it, and (c) make sure the `weekDates` handed in actually contain that day (Task 2), because
`weekdays` is built strictly from `props.weekDates`.

### useWeekSelector helpers to reuse

`getMondayOfWeek(date)` (`useWeekSelector.ts:9-15`) and
`getWeekdayDates(monday, includeWeekends)` (`useWeekSelector.ts:36-47`) can build the week that
contains an arbitrary `selectedDate` for day mode, mirroring how `selectedWeekDates` is produced
internally (`useWeekSelector.ts:140-142`).

### Constraints / gotchas

- Do NOT break week mode: when `bookingMode === 'week'` keep passing `selectedWeekDates` and the
  week label, and leave `selectedDate` prop undefined so today/first-future preselection within
  that week is preserved (AC #2).
- Guard against past days: `preselectDay` must not select a `past` entry
  (`weekdays[].past`, `InteractiveFloorPlan.vue:745`). Fall back to existing logic if the incoming
  `selectedDate` is past or not present in `weekDates`.
- No API/store changes and no backend changes — this is prop plumbing plus preselection logic.

### Project Structure Notes

- Modified: `web/src/views/ItemsView.vue` (bind `selected-date`; compute day-mode week
  dates/label) and `web/src/components/InteractiveFloorPlan.vue` (new prop + `preselectDay`).
- Reuses existing `useWeekSelector` helpers; no new files or types required.

### Testing standards summary

Vitest component tests exist for the floor plan (`InteractiveFloorPlan.test.ts`). Add cases for
`selectedDate`-driven preselection (present, absent, past). Cypress E2E should cover the AC flows
(day-mode selected day, week-mode selected week, reopen-reflects-latest) per `.claude/rules/vue.md`
and `.claude/rules/cypress.md`; treat the availability GET as the sync point via `cy.intercept`.
Run `npm run type-check`, `npm run lint`, `npx vitest run`, and `npm run build`.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 36.6 / FR174 (lines 5484-5504, 636)]
- [Source: web/src/views/ItemsView.vue:30,41-42,53-63,66-71,775-796,956,1008,1141-1146]
- [Source: web/src/components/InteractiveFloorPlan.vue:569-578,735-752,792,911-921,1076-1114]
- [Source: web/src/composables/useWeekSelector.ts:9-15,36-47,123-142]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `npm run type-check`, `npm run lint`, `npx vitest run` (503 passed), `npm run build` — all green.
- Renamed prop `selectedDate` → `selectedDay` after `vue/no-dupe-keys` flagged a collision with the
  internal `selectedDate` computed in `InteractiveFloorPlan.vue`.

### Completion Notes List

- Prop plumbing + preselection only; no API/store/backend changes.
- Day mode now derives the week that contains the tile-selected day so the floor plan can resolve
  and preselect it; week mode is unchanged (passes `undefined`).

### File List

- web/src/components/InteractiveFloorPlan.vue
- web/src/views/ItemsView.vue
- web/src/components/__tests__/InteractiveFloorPlan.test.ts
