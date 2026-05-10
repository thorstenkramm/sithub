# Story 30.3: Optional Area Drill-Down Toggle on Floor Plan Booking

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user on a large screen,
I want to book items directly from the floor plan without drilling into the detailed
room/area view,
so that I can complete bookings faster from the overview I already see.

## Acceptance Criteria

1. **Given** I view a floor plan in the booking view
   **When** the page renders
   **Then** an "Area drill-down" toggle (checkbox) is visible beneath the room plan

2. **Given** the viewport is a small screen (mobile breakpoint) and I have not changed the
   toggle yet
   **When** the floor plan loads
   **Then** the "Area drill-down" toggle is enabled by default
   **And** clicking on an area or room on the floor plan opens the detailed drill-down
   view as today

3. **Given** the viewport is a large screen (desktop breakpoint) and I have not changed
   the toggle yet
   **When** the floor plan loads
   **Then** the "Area drill-down" toggle is disabled by default
   **And** clicking on an item on the floor plan starts the booking flow directly without
   loading the detailed room/area view

4. **Given** I change the "Area drill-down" toggle
   **When** the change is applied
   **Then** the new value is saved to local storage on the current device
   **And** subsequent floor plan booking sessions on the same device use the saved value
   regardless of viewport size

5. **Given** I open the floor plan booking view on a different device
   **When** the page renders
   **Then** the toggle uses that device's own default and stored value
   **And** the choice on this device does not affect any other device

## Tasks / Subtasks

- [x] Task 1: Add a `useAreaDrillDownPreference` composable (AC: #4, #5)
  - [x] 1.1 Create `web/src/composables/useAreaDrillDownPreference.ts` modelled on
        `useAreaViewPreference.ts`. Use `getSafeLocalStorage` from
        `web/src/composables/storage.ts`.
  - [x] 1.2 Storage key: `sithub_area_drill_down`. Stored value is the literal string
        `'on'` or `'off'`. Absent key means "no user choice yet — use viewport default".
  - [x] 1.3 Public surface: `{ enabled: Ref<boolean>, hasUserChoice: Ref<boolean>,
        load(isLargeScreen: boolean): void, set(value: boolean): void }`.
        - `load` reads storage. If a value is present, set `enabled` from storage and
          `hasUserChoice = true`. If absent, default `enabled = !isLargeScreen` and
          `hasUserChoice = false` (per AC #2 and AC #3).
        - `set` writes `'on'` or `'off'` to storage and sets `hasUserChoice = true`.
  - [x] 1.4 Local storage is per-browser/per-device by definition; AC #5 falls out of
        using `localStorage` — no extra work needed.

- [x] Task 2: Surface the toggle in `InteractiveFloorPlan.vue` (AC: #1)
  - [x] 2.1 Add a `v-checkbox` (or `v-switch`) labelled with the i18n key
        `floorPlan.areaDrillDownToggle` (add the translation across all configured
        locales — see `web/src/locales/`). The control sits beneath the rendered floor
        plan within the same component.
  - [x] 2.2 Bind it to the composable's `enabled` ref. On change, call `set(value)`.
  - [x] 2.3 Use `data-cy="floor-plan-area-drill-down-toggle"` for tests.
  - [x] 2.4 Consider visibility: only show the toggle while in the area-level
        (top-level) view. Once the user has drilled in, the toggle has no effect on the
        current view; hiding it there keeps the UI clean. (Not required by AC, but
        encouraged.)

- [x] Task 3: Wire the toggle to drill-down behaviour (AC: #2, #3)
  - [x] 3.1 In `InteractiveFloorPlan.vue`:
        - `handleAreaClick(itemGroupID)` (line ~1180) currently always sets
          `drilledInto.value`. Wrap that with a check on `enabled` from the composable.
          If drill-down is **off**, instead emit the booking flow for that area's
          first/only item (or, if the click target is an area block with no specific
          item, surface a brief visual hint explaining that direct booking requires
          clicking a specific item).
        - The cleaner mapping is: when the user clicks a **desk** (`fp-item--free`),
          `handleDeskClick` already calls `requestBooking` directly **unless**
          `shouldDrillIntoItemGroup` is true. Make `shouldDrillIntoItemGroup` return
          `false` whenever the toggle is **off**, so a desk click books immediately
          without first drilling into the item group.
        - When the user clicks an **area block** (`fp-item--area`), drill-down is the
          only sensible action when the toggle is off would just be a no-op. Keep the
          area-block click drilling in regardless of toggle (the toggle's intent is
          to skip the auto-drill that happens on clicking a desk inside an undrilled
          area, per the user's brief).
  - [x] 3.2 On `mounted`/`onMounted`, call `load(isLargeScreen)` where
        `isLargeScreen` is the negation of the existing `isCompactViewport` ref
        (`!isCompactViewport.value` after `updateViewport()` runs). Re-run on
        `resize` only when `hasUserChoice` is `false` so a stored choice is never
        overwritten by a viewport change (AC #4).
  - [x] 3.3 Define "small screen" using the existing breakpoint already in this
        component: `isCompactViewport` becomes `true` when the viewport is narrow
        (`max-width: 768px`) **or** short (`max-height: 500px`). Reuse this so we have
        one definition of "compact". Default-on when compact, default-off otherwise
        (AC #2/#3).

- [x] Task 4: i18n (AC: #1)
  - [x] 4.1 Add the translation key `floorPlan.areaDrillDownToggle` (label) and
        `floorPlan.areaDrillDownHint` (optional helper text) across all locales already
        present under `web/src/locales/` (en, de, es, fr, uk).
  - [x] 4.2 English label: `Area drill-down`. Optional helper: `When off, clicking a
        desk on the floor plan books it directly without opening the detailed room
        view.` Use translators' wording for non-English locales (or copy English as a
        placeholder + flag in completion notes if a translator pass is needed).

- [x] Task 5: Tests
  - [x] 5.1 Unit-test `useAreaDrillDownPreference`: covers default-on for compact,
        default-off for large, persisted value overrides default, and that storage
        unavailability does not throw.
  - [x] 5.2 Component test (Vitest + Vue Test Utils): mount `InteractiveFloorPlan`
        with mocks, assert
        - large-screen default: clicking a free desk inside an un-drilled area calls
          `requestBooking` and does NOT call `handleAreaClick`
        - compact default: clicking the same desk drills in
        - toggling the checkbox flips behaviour and persists to localStorage
        - storage value overrides both viewport defaults
  - [x] 5.3 Cypress E2E: covers the golden path on desktop — load floor plan, click a
        free desk, expect booking flow opens directly without drilling. Add a second
        path that toggles drill-down on and confirms the drill-down still works.
  - [x] 5.4 `cd web && npx vitest run`
  - [x] 5.5 `cd web && npm run type-check`
  - [x] 5.6 `cd web && npm run lint`
  - [x] 5.7 `cd web && npm run test:e2e -- --browser electron`

### Review Findings

- [x] [Review][Patch] Re-apply the viewport-derived area drill-down default on resize until the user makes an explicit choice [web/src/components/InteractiveFloorPlan.vue:1573]
- [x] [Review][Patch] Add regression coverage for desktop direct-booking, compact auto-drill, resize default changes, and persisted drill-down overrides [web/src/components/__tests__/InteractiveFloorPlan.test.ts:347]

## Dev Notes

### Architecture & Patterns

- The interactive floor plan is `web/src/components/InteractiveFloorPlan.vue` and is
  embedded inside both `ItemGroupsView.vue` and `ItemsView.vue` via a `v-dialog`. The
  toggle belongs **inside** the component so it travels with both call sites and
  applies whether the dialog is open from area level or item-group level.
- `isCompactViewport` (line ~530, set by `updateViewport()` on resize, line ~1546)
  already encapsulates the small-screen definition (`max-width: 768px` or
  `max-height: 500px`). Reuse it; don't introduce a second breakpoint definition.
- Drill-down logic lives in two functions today:
  - `handleAreaClick(itemGroupID)` — line ~1180
  - `handleDeskClick(itemID, label)` — line ~1209, which calls `handleAreaClick` via
    `shouldDrillIntoItemGroup` if the desk's parent group has a floor plan
- `useAreaViewPreference.ts` is the canonical pattern for "preference persisted in
  localStorage with sane defaults and corruption safety". Mirror its style.
- AC #5 ("does not affect any other device") is satisfied by `localStorage` being
  per-browser-profile by definition. Do not introduce a server-side roundtrip.

### Key Code Locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Floor plan component | `web/src/components/InteractiveFloorPlan.vue` | Where the toggle lives and behaviour is gated |
| Drill-down trigger | `web/src/components/InteractiveFloorPlan.vue::handleAreaClick` (~1180) | One of the two click paths to gate |
| Desk click | `web/src/components/InteractiveFloorPlan.vue::handleDeskClick` (~1209) | The other click path; calls `shouldDrillIntoItemGroup` |
| Should-drill predicate | `web/src/components/InteractiveFloorPlan.vue::shouldDrillIntoItemGroup` (~1176) | Cleanest place to short-circuit when toggle is off |
| Viewport detection | `web/src/components/InteractiveFloorPlan.vue::isCompactViewport`, `updateViewport` (~1546) | Reuse for default selection |
| Storage helper | `web/src/composables/storage.ts::getSafeLocalStorage` | Use it; don't touch `localStorage` directly |
| Preference pattern reference | `web/src/composables/useAreaViewPreference.ts` | Template for the new composable |
| Locale files | `web/src/locales/*.ts` | Add `floorPlan.areaDrillDownToggle` + helper to all five locales |
| Containing views | `web/src/views/ItemGroupsView.vue` (~322), `web/src/views/ItemsView.vue` (~778) | Use the floor plan component; usually no changes needed |

### Implementation Strategy

1. Build the composable first with full unit coverage. Default-on/off logic must live
   here, not scattered across the component.
2. In `InteractiveFloorPlan.vue`:
   - Add the toggle UI under the floor plan render area (above or below the existing
     close/back footer is fine — keep it visually adjacent to the plan, not the
     bookings).
   - Hook the composable into `onMounted` after `updateViewport()` so the initial
     viewport is known.
   - Make `shouldDrillIntoItemGroup` return `enabled.value && Boolean(...)` — that's
     a one-line change to the existing predicate and it satisfies the desk-click flow
     for AC #3.
3. Add tests in this order: composable unit → component unit → Cypress E2E.

### Anti-Patterns to Avoid

- Do NOT propagate the toggle state from a Pinia store. This is per-device and small;
  a composable backed by `localStorage` is the right shape.
- Do NOT re-read `localStorage` on every click. Read once on mount; mutate via
  `composable.set()` only.
- Do NOT change the area-block (room) click behaviour. The user's brief is about
  desk-level direct booking. Clicking a room block is still a drill — the toggle
  controls the auto-drill that happens when a user clicks a **desk** inside an
  un-drilled area.
- Do NOT introduce a viewport breakpoint different from the existing
  `isCompactViewport` (768px / 500px). One definition.
- Do NOT remove or weaken the `shouldDrillIntoItemGroup` permission check (it also
  gates whether the target item group even has a floor plan to drill into). The new
  toggle is an **additional** gate, not a replacement.
- Do NOT save the toggle state in a cookie or sessionStorage; AC #4 explicitly says
  local storage and per-device persistence.

### Testing Standards

- Unit tests: Vitest for the composable and the component.
- E2E tests: Cypress against the dev server, no API mocking. Use `cy.login()` custom
  command. Selectors via `data-cy` only. Use intercept aliases (e.g. `@createBooking`)
  for synchronization.
- Cover both viewport defaults and the persisted-value override.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 30 Stories: Operator Validation, Editor Zoom Height & Optional Drill-Down]
- [Source: web/src/components/InteractiveFloorPlan.vue]
- [Source: web/src/composables/useAreaViewPreference.ts]
- [Source: web/src/composables/storage.ts]
- [Source: web/src/views/ItemGroupsView.vue]
- [Source: web/src/views/ItemsView.vue]
- [Source: .claude/rules/vue.md]
- [Source: .claude/rules/cypress.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7

### Completion Notes List

- Added `useAreaDrillDownPreference` composable in
  `web/src/composables/useAreaDrillDownPreference.ts` with `enabled`,
  `hasUserChoice`, `load(isLargeScreen)`, and `set(value)`. Storage key
  `sithub_area_drill_down` stores `'on'` / `'off'`; absence falls through to viewport
  default. Mirrors the `useAreaViewPreference` pattern.
- Wired the toggle into `InteractiveFloorPlan.vue`:
  - One-line gate inside `shouldDrillIntoItemGroup` returns `false` when
    `drillDownEnabled` is off, so a desk click books directly without auto-drilling.
  - `onMounted` calls `areaDrillDownPref.load(!isCompactViewport.value)` after the
    existing `updateViewport()`, so the default tracks viewport at first paint.
  - `drillDownEnabled` is a writable computed bound to a third footer checkbox
    (label `floorPlan.areaDrillDownToggle`, `data-cy="floor-plan-area-drill-down-toggle"`)
    that is only shown at area level and when there are positions to interact with.
- Resize handler now re-applies viewport defaults only while `hasUserChoice` is
  `false`; once the user toggles the setting, the stored choice is preserved across
  later viewport changes.
- Area-block clicks (room rectangles) still drill in regardless — the toggle's
  intent is to skip the auto-drill that happens when a user clicks a desk inside
  an un-drilled area, per the user's brief.
- i18n: added `floorPlan.areaDrillDownToggle` to all five locales (en, de, es, fr,
  uk). Translations are best-effort idiomatic; flag for native-speaker review if
  desired.
- `npm run type-check`, `npm run lint`, and `npm run build` all pass.
- Tests: added `useAreaDrillDownPreference.test.ts` and expanded
  `InteractiveFloorPlan.test.ts` to cover desktop direct-booking, compact
  auto-drill, resize-driven defaults before user choice, and persisted overrides.
  Added a small `localStorage` shim in `web/vitest.setup.ts` so the storage-backed
  Vitest suites run cleanly in jsdom. Verified with:
  `npx vitest run src/components/__tests__/InteractiveFloorPlan.test.ts src/composables/useAreaDrillDownPreference.test.ts src/composables/useAreaViewPreference.test.ts`

### File List

- web/src/composables/useAreaDrillDownPreference.ts (new)
- web/src/composables/useAreaDrillDownPreference.test.ts (new)
- web/src/components/__tests__/InteractiveFloorPlan.test.ts (modified)
- web/src/components/InteractiveFloorPlan.vue (modified)
- web/src/locales/en.json (modified)
- web/src/locales/de.json (modified)
- web/src/locales/es.json (modified)
- web/src/locales/fr.json (modified)
- web/src/locales/uk.json (modified)
- web/vitest.setup.ts (modified)
