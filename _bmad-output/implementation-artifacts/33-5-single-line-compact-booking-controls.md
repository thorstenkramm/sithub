# Story 33.5: Single-Line Compact Booking Controls Layout

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user on the item-groups view,
I want all booking controls on a single line on a wide screen,
so that the controls take less vertical space and leave more room for tiles.

## Acceptance Criteria

1. **Given** I am on the item-groups view on a wide desktop viewport
   **When** the booking-controls card renders
   **Then** the day/week toggle, date or week selector, equipment filter input
   (with its info icon), "Book for a colleague" checkbox, and colleague-selection
   dropdown all appear on a single horizontal row
   **And** none of them is truncated or hidden behind another

2. **Given** the same wide viewport
   **When** I compare the new controls-card height to the previous (multi-row)
   layout
   **Then** the new card is visibly shorter

3. **Given** I narrow the viewport progressively
   **When** the row no longer fits horizontally
   **Then** the controls wrap onto additional rows naturally (Vuetify flex-wrap),
   with no element overlapping another

4. **Given** I switch between day mode and week mode
   **When** the controls re-render
   **Then** the single-line layout is preserved in both modes (the date picker
   swaps for the week selector but the rest stays in place)

## Tasks / Subtasks

- [ ] Task 1: Consolidate the controls into a single flex-wrap row (AC: #1, #3)
  - [ ] 1.1 In `web/src/views/ItemsView.vue` around lines 25–146, the
        booking-controls card currently has THREE separate horizontal blocks:
        - The mode toggle (`v-btn-toggle` with `mode-day-btn`/`mode-week-btn`)
          standalone, with `mb-4`.
        - A `d-flex flex-wrap align-end ga-4 mb-4` row containing the date
          picker / week selector and the floor-plan button.
        - The `.booking-type-row` containing the new checkbox + autocomplete
          (per Story 33.4).
        - The equipment filter row (`d-flex align-center ga-2 mt-4`).
  - [ ] 1.2 Wrap ALL of these in a single flex container under the existing
        `<v-card-text>`:
        ```vue
        <v-card-text>
          <div class="booking-controls-row d-flex flex-wrap align-center ga-3">
            <v-btn-toggle ... />            <!-- day/week toggle -->
            <DatePickerField v-if="bookingMode === 'day'" ... />
            <v-select v-if="bookingMode === 'week'" ... />  <!-- week selector -->
            <v-btn v-if="itemGroupFloorPlan" ... />          <!-- floor plan btn -->
            <div class="d-flex align-center ga-2 equipment-filter-cluster">
              <v-combobox ... data-cy="equipment-filter-input" />
              <v-tooltip ...><v-btn icon ... /></v-tooltip>   <!-- save/delete  -->
              <v-btn icon ... data-cy="equipment-filter-info" />  <!-- info     -->
            </div>
            <v-checkbox ... data-cy="book-colleague-checkbox" />
            <v-autocomplete ... data-cy="colleague-select"
              class="colleague-select-inline" />
          </div>
        </v-card-text>
        ```
        Use `align-center` so all controls line up vertically along their
        centers. Use a smaller gap (`ga-3` = 12px) than today's `ga-4` to fit
        all five on common desktop widths.
  - [ ] 1.3 Remove the now-vestigial classes from Story 32.3/33.4:
        `.booking-type-row` and `.booking-type-radios`. Keep
        `.colleague-select-inline` but adjust its `flex-basis` (Task 3).

- [ ] Task 2: Per-control width tuning so the row fits on common desktops (AC: #1)
  - [ ] 2.1 Day/week toggle: natural width (~120px).
  - [ ] 2.2 DatePickerField / week selector: cap width via the existing
        `style="max-width: 320px;"`; on the consolidated row reduce to 240px so
        we save horizontal space. `style="max-width: 240px; min-width: 200px;"`.
  - [ ] 2.3 Floor-plan button: natural width, prepend icon + label.
  - [ ] 2.4 Equipment-filter cluster: cap the combobox at 240px on the new row.
        The cluster (combobox + save/delete icon + info icon) sits as one flex
        unit and should `flex: 1 1 280px` so it absorbs remaining horizontal
        space gracefully.
  - [ ] 2.5 `book-colleague-checkbox`: natural width — Vuetify renders it as
        ~180–220px depending on locale.
  - [ ] 2.6 `colleague-select-inline`: redefine in scoped CSS as
        ```css
        .colleague-select-inline { flex: 0 0 260px; max-width: 320px; }
        @media (max-width: 600px) { .colleague-select-inline { flex: 1 1 100%; max-width: 100%; } }
        ```
        Slightly narrower than before so the row fits at common widths.

- [ ] Task 3: Card height shrinks (AC: #2)
  - [ ] 3.1 Remove the `mb-4` from the mode toggle, the date row, the
        booking-type row, and the `mt-4` from the equipment-filter row. Spacing
        is supplied by the row's `ga-3`.
  - [ ] 3.2 Keep the outer `<v-card class="mb-6">` so the card still has
        breathing room from the tile grid below.
  - [ ] 3.3 Verify the new card is visibly shorter via chrome-devtools-mcp:
        screenshot the controls-card before vs after, measure both heights, the
        new one should be ~50–60% of the old one (single row vs 3 rows).

- [ ] Task 4: Wrapping behavior (AC: #3)
  - [ ] 4.1 The single flex-wrap container handles narrow viewports natively.
        Verify by resizing the viewport in chrome-devtools-mcp through breakpoints
        (1440 → 1024 → 900 → 768 → 600 → 414); the row should wrap progressively,
        no element overlaps another, and the bottom-row tile grid never gets
        covered by the controls.
  - [ ] 4.2 On mobile (<600px), the autocomplete becomes full-width per the
        media-query in Task 2.6. The other controls wrap to subsequent lines
        naturally.

- [ ] Task 5: Same layout in both day and week modes (AC: #4)
  - [ ] 5.1 The `v-if`-gated date picker (day mode) and week selector (week
        mode) live in the same flex slot, so the rest of the row layout does not
        shift when the mode toggles.
  - [ ] 5.2 Confirm by switching modes via chrome-devtools-mcp and verifying
        the equipment filter, checkbox, and dropdown do not change their
        x-position.

- [ ] Task 6: Tests (Vitest + Vue Test Utils)
  - [ ] 6.1 In `web/src/views/ItemsView.test.ts`:
        - Update or remove the `booking-type row layout` describe block (it was
          asserting `.booking-type-row` which is being deleted).
        - Add a `describe('compact booking controls row', ...)` that asserts:
          - `.booking-controls-row` exists at mount in day mode.
          - All five control selectors live inside it:
            `mode-day-btn`, `items-date`, `equipment-filter-input`,
            `book-colleague-checkbox`, `colleague-select`.
          - Switching `bookingMode` to `'week'` swaps `items-date` for
            `week-selector` but keeps every other selector inside the same row.

- [ ] Task 7: Verification commands
  - [ ] 7.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green.
  - [ ] 7.2 Manual smoke (chrome-devtools-mcp): on a wide viewport (≥1440px)
        confirm all controls appear on one row; resize down and confirm graceful
        wrapping; compare card height visually with `private/img_21.png` and
        `private/img_22.png` (the target).

## Dev Notes

### Coordination with Story 33.4

33.4 ships the checkbox+autocomplete in place of the radio group. 33.5
collapses the entire controls area to a single row. If both land in the same
PR, do 33.4's logic change first (it's the safer one), then 33.5's layout
change second. If 33.5 lands first standalone, the radio group stays inside
the new row (visually awkward but functional); 33.4 then drops it cleanly.

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Existing card text container | `<v-card class="mb-6"><v-card-text>` | `web/src/views/ItemsView.vue:25` |
| Mode toggle | `<v-btn-toggle v-model="bookingMode">` | `web/src/views/ItemsView.vue:28–37` |
| Date picker | `<DatePickerField>` | `web/src/views/ItemsView.vue:41–51` |
| Week selector | `<v-select data-cy="week-selector">` | `web/src/views/ItemsView.vue:54–65` |
| Floor-plan button | `<v-btn data-cy="item-group-floor-plan-btn">` | `web/src/views/ItemsView.vue:67–76` |
| Equipment filter cluster | `<div class="d-flex align-center ga-2 mt-4">` | `web/src/views/ItemsView.vue:108–144` |
| Visual target (week) | `private/img_21.png` | — |
| Visual target (day) | `private/img_22.png` | — |

### Anti-patterns to avoid

- Do NOT use a CSS Grid for this — the controls have variable natural widths
  and grid would force them into rigid columns. Flex with `flex-wrap` and
  per-control widths is the right shape.
- Do NOT add custom JS breakpoint detection. The `flex-wrap` rule handles
  narrow viewports purely in CSS.
- Do NOT collapse the mode toggle into the same flex item as the date picker.
  They are separate concerns; keep them as siblings so the toggle does not
  jump when the date picker re-mounts on mode change.
- Do NOT remove the equipment filter info / save icons — they are part of
  the equipment-filter cluster and must stay attached to the input.

### Testing standards

- Vitest + Vue Test Utils, extending `ItemsView.test.ts`. No new test file.
- Visual height comparison is a manual smoke item; no JSDom assertion.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/views/ItemsView.vue:25–146]
- [Source: _bmad-output/implementation-artifacts/33-4-single-checkbox-and-always-rendered-colleague-dropdown.md]
- [Source: private/epic-33.md]
- [Source: private/img_21.png]
- [Source: private/img_22.png]
- [Source: .claude/rules/vue.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/ItemsView.test.ts` — 89 tests pass
  (2 new for the compact booking controls row).
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all clean.
- ItemsView bundle: 43.57 KB (down from 44.54 KB pre-Epic-32) thanks to the
  removed conditional rendering branch.

### Completion Notes List

- Collapsed the three separate horizontal blocks (mode toggle, date row,
  booking-type row, equipment filter row) into a single
  `<div class="booking-controls-row d-flex flex-wrap align-center ga-3">`
  containing every control in the order specified in `private/img_21.png`
  and `img_22.png`: mode toggle, date picker / week selector, floor plan
  button, equipment filter cluster (input + save/delete icon + info icon),
  book-for-colleague checkbox, colleague-select autocomplete.
- Per-control widths tuned via scoped CSS so the row fits at common desktop
  widths (≥1280px). On viewports `<= 600px`, the date input, equipment
  filter cluster, and colleague-select all flex to full width so the row
  wraps gracefully.
- Set `.booking-controls-row { min-height: 40px }` so the row height does
  not jitter when the autocomplete enables/disables.
- Removed the vestigial `.booking-type-row`, `.booking-type-radios`, and
  Story-32.3's `.colleague-select-inline` block (replaced with the new
  per-control CSS rules).
- Added a `compact booking controls row` describe block in
  `ItemsView.test.ts` with two tests verifying day-mode and week-mode both
  render all controls inside the single `.booking-controls-row` container.
  The day-mode `items-date` selector is asserted via the DatePickerField
  component name because the test stub for `v-menu` doesn't render the
  text-field inside the menu's activator slot.

### File List

Frontend (modified):

- `web/src/views/ItemsView.vue` — booking-controls card collapsed into a
  single flex-wrap row; per-control widths set via new scoped CSS; obsolete
  classes removed.
- `web/src/views/ItemsView.test.ts` — added `compact booking controls row`
  describe block.
