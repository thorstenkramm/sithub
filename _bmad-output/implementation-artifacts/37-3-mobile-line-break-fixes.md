# Story 37.3: Mobile Line-Break Fixes

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a mobile user,
I want labels to wrap at word boundaries,
so that the interface is readable and looks polished.

## Acceptance Criteria

1. Given a 360px wide viewport and the item group bookings page, when the breadcrumb shows "Büro 2.EG" and
   "Cube 1", then no breadcrumb entry is split inside a word; long entries truncate with an ellipsis or wrap
   between words.
2. Given a 360px wide viewport and the floor plan dialog, when the footer shows the checkboxes
   "Beschriftungen anzeigen" and "Avatare anzeigen" and the close button, then labels are not split inside a
   word; the footer wraps controls onto separate lines if needed.
3. Given the fix, when reviewed, then other mobile views with the same breadcrumb and checkbox components are
   checked for the same problem, and a Cypress test at 360px asserts the affected labels render on word
   boundaries (no element narrower than its longest word).

## Tasks / Subtasks

- [x] Task 1: Keep breadcrumb entries intact and let the trail wrap between entries (AC: #1)
  - [x] In `web/src/components/PageHeader.vue`, the `<ol>` (line 5) is `d-flex align-center ga-1` with no
        wrapping, so every `<li>` (lines 6-11) is squeezed to its min-content width and "Büro 2.EG" /
        "Cube 1" break onto two lines each (see `private/img_42.png`). Add `flex-wrap` to the `<ol>` (Vuetify
        utility class `flex-wrap`) and a small row gap so whole entries move to the next line instead.
  - [x] Give the crumb label (both the `router-link.breadcrumb-link`, lines 12-19, and the current-page
        `<span>`, lines 20-22) a shared class, e.g. `breadcrumb-label`, styled with
        `white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 100%;` so an entry never
        breaks internally and an entry wider than the viewport truncates with an ellipsis.
  - [x] Add `min-width: 0; max-width: 100%;` to the `<li>` so the ellipsis can actually engage inside the
        flex row; keep the chevron (`v-icon`, lines 23-29) with `flex-shrink: 0` so it never collapses.
  - [x] Add `:title="crumb.text"` to the label so a truncated entry stays readable on hover/long-press.
  - [x] Do not change `data-cy` attributes (`breadcrumbs`, `breadcrumb-item-N`, `breadcrumb-link-N`,
        `breadcrumb-text-N`); `cypress/e2e/breadcrumbs.cy.ts` and `PageHeader.test.ts` depend on them.
- [x] Task 2: Stop the floor-plan footer checkbox labels from breaking mid-word (AC: #2)
  - [x] In `web/src/components/InteractiveFloorPlan.vue`, the footer (`.fp-footer`, lines 335-370) holds a
        `div.d-flex.align-center.ga-4` (line 336) with three `v-checkbox`es (`fp-show-labels` 337-344,
        `fp-show-avatars` 345-352, `floor-plan-area-drill-down-toggle` 353-360) plus the close button
        `fp-close-btn` (362-369). Neither the footer nor the inner div wraps, so at 360px the checkboxes shrink
        and Vuetify's label rule `white-space: normal; word-break: break-word`
        (`node_modules/vuetify/lib/components/VSelectionControl/VSelectionControl.css:11-12`, Vuetify 3.13.6)
        splits "Beschriftungen" and "Avatare" mid-word (see `private/img_43.png`).
  - [x] Add `flex-wrap: wrap` to `.fp-footer` (CSS lines 2112-2118) and add the `flex-wrap` class plus a
        row gap to the inner checkbox container (line 336).
  - [x] In the scoped style add a rule scoped to the footer only, e.g.
        `.fp-footer :deep(.v-selection-control .v-label) { white-space: nowrap; word-break: normal; }`
        and `.fp-footer :deep(.v-checkbox) { flex: 0 0 auto; }` so each checkbox keeps its natural width and
        moves to its own line when it does not fit.
  - [x] Keep the close button at the end: give `.fp-close-btn` `margin-left: auto` so it right-aligns on
        whatever line it lands on. Do not hide or shorten any label.
  - [x] Check the compact (`max-width: 768px`) sticky footer override (CSS lines 2164-2182): with wrapping the
        footer becomes taller; confirm the floor-plan scroll area still shows and nothing is hidden behind
        the sticky footer.
- [x] Task 3: Review other mobile places using the same components (AC: #3)
  - [x] Breadcrumbs: `PageHeader` is the only breadcrumb implementation (no `v-breadcrumbs` in `web/src`).
        It is used by `AreasView`, `ItemGroupsView`, `ItemsView`, `ItemGroupBookingsView`,
        `AreaPresenceView`, `MyBookingsView`, `BookingHistoryView`, `FloorPlanEditorView`; the fix in Task 1
        covers all of them. Spot-check `ItemsView` and `AreaPresenceView` at 360px.
  - [x] Floor-plan dialog: `InteractiveFloorPlan` is opened from `ItemGroupsView.vue` (area floor plan,
        dialog `floor-plan-dialog`, lines 225-246) and `ItemsView.vue` (`item-group-floor-plan-dialog`,
        lines 761-781). Both share the footer, so Task 2 covers both; verify both at 360px, including the
        area-level case where the third checkbox "Bereichs-Drilldown" is visible.
  - [x] Other selection controls with labels at 360px: `WarningConfirmDialog.vue:49-57`
        ("dont show again" checkbox), `v-switch` in `ItemGroupsView.vue:36,52` and `MyBookingsView.vue:31,47`
        use the same Vuetify label rule. Check them in German; fix only if a word actually breaks, using the
        same targeted `:deep(.v-label)` rule rather than a global override. Record the result in the
        Completion Notes.
  - [x] Checkboxes without labels (week-day checkboxes in `ItemsView.vue` ~493-590, booking-day checkboxes in
        `InteractiveFloorPlan.vue` 406-433) are not affected; leave them alone.
- [x] Task 4: Tests (AC: #1, #2, #3)
  - [x] Add a reusable Cypress helper (e.g. `assertNoMidWordBreak(selector)` in `web/cypress/support/flows.ts`
        or `commands.ts`) that, for each matched element, walks its text nodes, creates a `Range` per word
        and asserts `range.getClientRects().length === 1` (a word that spans two lines yields two rects).
        This is the concrete form of "no element narrower than its longest word"; additionally assert
        `el.getBoundingClientRect().width >= widest word rect width` for elements that are not ellipsised.
  - [x] New spec `web/cypress/e2e/mobile-line-breaks.cy.ts` with `cy.viewport(360, 740)` and German locale.
        Note `cypress/support/e2e.ts` forces `sithub_locale = 'en'` on `window:before:load`; set
        `sithub_locale = 'de'` in `cy.visit(..., { onBeforeLoad })` and assert a German string is visible
        (e.g. "Beschriftungen anzeigen") to prove the override took effect.
  - [x] Test `should keep breadcrumb entries on word boundaries at 360px`: open the item group bookings page
        (`/item-groups/:itemGroupId/bookings`, route `item-group-bookings`) for an area named "Büro 2.EG"
        and an item group named "Cube 1"; run the helper on `[data-cy^="breadcrumb-link-"]` and
        `[data-cy^="breadcrumb-text-"]`, and assert each crumb label is on a single line.
  - [x] Test `should keep floor-plan footer labels on word boundaries at 360px`: open the floor plan from
        `ItemsView` (`item-group-floor-plan-btn`), wait for `@positions`, run the helper on the label of
        `[data-cy="fp-show-labels"]`, `[data-cy="fp-show-avatars"]` and `[data-cy="fp-close-btn"]`, and assert
        all three are visible inside the viewport (no horizontal overflow of `[data-cy="fp-footer"]`).
  - [x] Extend `web/src/components/__tests__/PageHeader.test.ts` with a case asserting the `<ol>` carries the
        wrap class and each crumb label carries `breadcrumb-label` and a `title` with the full text (JSDOM
        has no layout, so pixel checks stay in Cypress).
  - [x] Run `npm run type-check`, `npm run lint`, `npx vitest run` (coverage >= 75%), `npm run build`,
        jscpd for TS, and `npm run test:e2e -- --browser electron` against `npm run dev`.

### Review Findings

No actionable findings. Shared breadcrumbs and both floor-plan footers pass the 360px German browser checks.

## Dev Notes

Source: [Source: `_bmad-output/planning-artifacts/epics.md` Story 37.3 / FR195]. FR195: the mobile view has no
words broken mid-word; breadcrumb entries (e.g. "Büro 2.EG", "Cube 1") and the floor-plan dialog checkbox
labels ("Beschriftungen anzeigen", "Avatare anzeigen") wrap only at word boundaries or truncate with an
ellipsis; at 360px width no label is split inside a word
[Source: `_bmad-output/planning-artifacts/epics.md:723-726,5777-5800`].

This is a CSS/template-only story. No backend, API, locale or store change.

### Breadcrumb: why it breaks today

`PageHeader.vue` renders the trail as a non-wrapping flex row:

```vue
<ol class="d-flex align-center ga-1 text-body-2">
  <li v-for="(crumb, index) in breadcrumbs" class="d-flex align-center" ...>
```

[Source: web/src/components/PageHeader.vue:5-11]. When the four entries ("Startseite", "Büro 2.EG", "Cube 1",
"Buchungen") do not fit in 360px, flexbox shrinks every `<li>` to its min-content width and the label text
wraps inside the entry ("Büro" / "2.EG", "Cube" / "1"), see `private/img_42.png`. The scoped style only
covers `list-style`, padding and link underline [Source: web/src/components/PageHeader.vue:62-76]. The fix
is `flex-wrap` on the `<ol>` plus `nowrap` + ellipsis on each label, so entries stay whole and move to the
next line as a unit.

### Floor-plan footer: why it breaks today

Footer template [Source: web/src/components/InteractiveFloorPlan.vue:335-370]:

```vue
<div class="fp-footer" data-cy="fp-footer">
  <div class="d-flex align-center ga-4">
    <v-checkbox ... :label="$t('floorPlan.showLabels')" data-cy="fp-show-labels" />
    <v-checkbox ... :label="$t('floorPlan.showAvatars')" data-cy="fp-show-avatars" />
    <v-checkbox v-if="areaLevel && !drilledInto && ..." data-cy="floor-plan-area-drill-down-toggle" />
  </div>
  <v-btn variant="text" data-cy="fp-close-btn" class="fp-close-btn">{{ $t('common.close') }}</v-btn>
</div>
```

`.fp-footer` is `display: flex; justify-content: space-between` without wrapping
[Source: web/src/components/InteractiveFloorPlan.vue:2112-2118]; in compact mode it becomes a sticky bar
[Source: web/src/components/InteractiveFloorPlan.vue:2164-2182]. Vuetify's selection-control label uses
`white-space: normal` and `word-break: break-word`
[Source: web/node_modules/vuetify/lib/components/VSelectionControl/VSelectionControl.css:11-12], so a squeezed
checkbox breaks "Beschriftungen" into "Beschr/iftunge/n" (`private/img_43.png`). German strings:
`floorPlan.showLabels` "Beschriftungen anzeigen", `floorPlan.showAvatars` "Avatare anzeigen",
`floorPlan.areaDrillDownToggle` "Bereichs-Drilldown", `common.close` "Schließen"
[Source: web/src/locales/de.json:262-264,305].

Keep the override scoped to `.fp-footer` with `:deep(...)`. A global `.v-label` override would change every
form in the app and is out of scope.

### Compact viewport

`InteractiveFloorPlan` sets `isCompactViewport` from `(max-width: 768px)` or `(max-height: 500px)`
[Source: web/src/components/InteractiveFloorPlan.vue:619,1854-1863]; both host dialogs are `fullscreen` in that
mode (`ItemGroupsView.vue:225-230`, `ItemsView.vue:761-766`). At 360px the footer is the sticky variant, so
test in that state.

### Cypress notes

- E2E locale is forced to English in `web/cypress/support/e2e.ts:6-8`. English labels ("Show labels",
  "Show avatars") are short and would not reproduce the bug, so the new spec must switch to German.
- Existing floor-plan and breadcrumb specs show how to open the views: `breadcrumbs.cy.ts` (flows
  `resetAndLogin`, `setupItemsPageIntercepts`, `openFirstAreaAndWaitItemGroups`) and
  `floor-plan-booking.cy.ts:43-57` (`item-group-floor-plan-btn`, `@positions`).
- Use `data-cy` selectors and intercepts with aliases as synchronization points, no fixed waits
  [Source: .claude/rules/cypress.md].

### Open questions / assumptions

- Test data (product decision 2026-10-08, no stubs): extend the E2E areas seed with an area "Büro 2.EG"
  containing an item group "Cube 1" (with at least one item). Update BOTH seeds identically: the
  `TEST_SPACES` heredoc in `run-all-tests.sh` and `.sithub-test-areas.yaml` in `.github/workflows/ci.yml`
  (~137-166). Story 37.2 extends the same seed with a floor plan; keep the additions compatible. Existing specs
  that count areas or item groups must still pass. Intercepts are used only as sync points (Cypress rule).
- AC #1 allows "truncate with an ellipsis or wrap between words". Assumption: an individual breadcrumb entry
  never wraps internally (even at a space); the trail wraps between entries, and a single entry wider than
  the viewport is ellipsised with a `title` tooltip. This matches the product intent from `img_42.png`.
- If "Bereichs-Drilldown" (hyphenated) breaks at the hyphen, that counts as a word boundary and is accepted;
  with `white-space: nowrap` it will not break at all.

### Project Structure Notes

- Modified: `web/src/components/PageHeader.vue` (template classes, `title`, scoped CSS).
- Modified: `web/src/components/InteractiveFloorPlan.vue` (footer container class, scoped CSS for
  `.fp-footer`, `.fp-close-btn`).
- Modified: `web/src/components/__tests__/PageHeader.test.ts`.
- New: `web/cypress/e2e/mobile-line-breaks.cy.ts`; helper added to `web/cypress/support/flows.ts` or
  `commands.ts` (keep it thin).
- Possibly modified after Task 3 review: `WarningConfirmDialog.vue`, `ItemGroupsView.vue`, `MyBookingsView.vue`
  (only if a mid-word break is observed).

### Testing standards summary

Vue rules: Composition API, Vuetify components, no `any`, run type-check, lint, Vitest (coverage >= 75%),
build, jscpd (`npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0 --exit-code 1`), and
Cypress with Electron against the dev server [Source: .claude/rules/vue.md]. Every AC needs E2E coverage;
test names use "should ..." [Source: .claude/rules/cypress.md]. Visual check with Chrome DevTools MCP at
360px in German for the item group bookings page, the items page and both floor-plan dialogs; take
before/after screenshots [Source: .claude/rules/mcps.md].

### References

- [Source: web/src/components/PageHeader.vue:4-32,62-76]
- [Source: web/src/components/InteractiveFloorPlan.vue:335-370,619,1854-1863,2112-2118,2164-2182]
- [Source: web/src/views/ItemGroupsView.vue:225-246]
- [Source: web/src/views/ItemsView.vue:761-781]
- [Source: web/src/views/ItemGroupBookingsView.vue:3-6,105-120]
- [Source: web/src/router/index.ts:47-49]
- [Source: web/src/locales/de.json:262-264,305]
- [Source: web/cypress/support/e2e.ts:4-8]
- [Source: web/cypress/e2e/breadcrumbs.cy.ts]
- [Source: web/cypress/e2e/floor-plan-booking.cy.ts:14-57]
- [Source: private/img_42.png, private/img_43.png]
- [Source: `_bmad-output/planning-artifacts/epics.md:723-726,5777-5800`]

## Dev Agent Record

### Agent Model Used

Claude Opus 5.5 (claude-opus-5-5)

### Debug Log References

- RED: the new spec failed before the fixes ("Büro 2.EG" on two lines, "Beschriftungen" split). A first
  breadcrumb check only looked at single words and passed against the old `PageHeader` (the screenshot shows a
  break between "Büro" and "2.EG", not inside a word); `assertSingleLine` was added so the test catches it.
- Visual check at 360px in German via temporary Cypress screenshots (bookings page, area floor plan incl.
  drill-down checkbox): crumbs wrap between whole entries, footer controls each on their own line, close button
  right-aligned, floor plan scroll area stays visible above the sticky footer. Temporary spec removed.

### Completion Notes List

- `PageHeader.vue`: `<ol>` gets `flex-wrap`; each `<li>` `breadcrumb-item` (`min-width: 0; max-width: 100%`);
  link and current-page span share `breadcrumb-label` (`nowrap`, `overflow: hidden`, ellipsis) plus
  `:title="crumb.text"`; chevron `breadcrumb-separator` with `flex-shrink: 0`. `data-cy` attributes unchanged.
- `InteractiveFloorPlan.vue` footer: `.fp-footer` wraps, inner `fp-footer-controls` container wraps; scoped
  `:deep` rules keep `.v-checkbox` at natural width and the label `nowrap` (footer only, no global `.v-label`
  override); `.fp-close-btn` `margin-left: auto`. No label hidden or shortened.
- Task 3 review: `PageHeader` is the only breadcrumb implementation, so all eight views are covered; ItemsView and
  AreaPresenceView breadcrumbs verified by E2E. Both floor-plan dialogs (item view and area level incl.
  "Bereichs-Drilldown") verified by E2E. View-switch labels ("Kacheln"/"Tabelle", same markup in
  `ItemGroupsView` and `MyBookingsView`) are single words and verified on the area page by E2E; the `v-switch`es
  have no labels. `WarningConfirmDialog` checkbox ("Nicht mehr anzeigen.") sits in a full-width
  `flex-column` row and cannot be squeezed; not reproducible without warned items in the seed, no change needed.
  No other file changed.
- Test data: reuses the "Büro 2.EG" / "Cube 1" seed added in story 37.2 (both `run-all-tests.sh` and `ci.yml`).
  `seedFloorPlanPositions` moved to `cypress/support/flows.ts` and shared by both specs.
- Helpers in `cypress/support/flows.ts`: `assertNoMidWordBreak` (one line box per word Range; non-ellipsised
  elements at least as wide as their widest word) and `assertSingleLine`.
- Tests: Vitest 58 files / 588 tests pass (coverage statements 81.46 %, lines 82.88 %). Cypress full suite 81/81
  incl. new `mobile-line-breaks.cy.ts` (4 tests, 360px, German). type-check, ESLint, jscpd (0 clones), build
  pass.

### File List

- `web/src/components/PageHeader.vue`
- `web/src/components/__tests__/PageHeader.test.ts`
- `web/src/components/InteractiveFloorPlan.vue`
- `web/cypress/e2e/mobile-line-breaks.cy.ts` (new)
- `web/cypress/support/flows.ts`
- `web/cypress/e2e/week-navigation.cy.ts`
- `_bmad-output/implementation-artifacts/sprint-status.yaml`

## Change Log

- 2026-10-08: Implemented story 37.3. Breadcrumb entries stay whole and the trail wraps between entries;
  floor-plan footer wraps controls instead of breaking labels; 360px German E2E with word-boundary helpers.
  Status -> review.

## Code Review Record

- Date: 2026-10-08. Parallel Blind Hunter, Edge Case Hunter and Acceptance Auditor reviews; actionable findings auto-fixed as requested.
- Validation: 600 Vitest tests pass in 58 files; statements 81.56%, lines 83.01%. All 12 story Cypress tests pass against the Vite dev server and an isolated real backend; week/mobile specs rerun after the final loading fix (11/11).
- Type-check, ESLint, production build, TypeScript duplication check (0 clones), Go tests and golangci-lint (0 issues) pass. No unresolved review findings.
- 2026-10-08: Code review complete; status -> done.
