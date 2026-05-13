# Story 33.3: Widen Weekly-Table Booking-Cancel Popover

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user cancelling my own booking from the weekly table view,
I want the cancel-confirmation popover to be tall and wide enough to show the cancel
button fully,
so that I can complete the action without scrolling inside the popover.

## Acceptance Criteria

1. **Given** I am on the weekly table view and I click one of my own booked cells
   **When** the booking-cancel popover opens
   **Then** the popover container is sized so that the Person, Platz, Datum lines
   and both the "Schliessen" and "Buchung stornieren" buttons are fully visible
   without any internal scrolling
   **And** the popover remains anchored to the cell that opened it

2. **Given** the popover is open
   **When** the buttons render
   **Then** neither button is clipped at the bottom edge of the container

## Tasks / Subtasks

- [ ] Task 1: Widen and lift container height constraints (AC: #1, #2)
  - [ ] 1.1 Edit `web/src/components/area-weekly-matrix/MatrixCancelPopover.vue:7`
        — increase the `<v-menu max-width="300">` value. Recommended:
        `max-width="380"` (slightly wider than today; matches the longest expected
        booker name "Alexander Seidemann-Klamant" plus padding without wrapping).
  - [ ] 1.2 Add an explicit `min-width="320"` to the same `<v-menu>` so the
        popover never collapses too narrow on short content.
  - [ ] 1.3 The root cause of the button-clipping in `private/img_20.png` is the
        v-menu's default `max-height` (Vuetify defaults to viewport-height
        constrained, with internal `overflow-y: auto`). When the activator cell
        is near the bottom edge of the table the popover's max-height shrinks and
        the action row falls below the fold. Add an explicit, generous
        `max-height="none"` on the v-menu so the popover always renders at its
        natural content height, and rely on the location preference flipping
        (Vuetify auto-flips `bottom` → `top` when there's no room) to keep it on
        screen.
  - [ ] 1.4 If `max-height="none"` causes layout issues on very small viewports
        (which is unlikely since the table view is desktop-only per Story 29.x),
        switch to `max-height="400"` as a fallback. Verify with chrome-devtools-mcp
        by reducing the viewport height and opening the popover from a
        bottom-row cell.

- [ ] Task 2: Verify the inner card already lays out cleanly (no change expected)
  - [ ] 2.1 Inner `<v-card class="pa-4">` (line 11) — keep as is. The
        `<v-card-actions class="pa-0">` (line 33) already lets the actions span
        the natural row height; the clipping was external (the v-menu container),
        not internal.
  - [ ] 2.2 Confirm the existing `<v-spacer />` + `<v-btn>` layout produces a
        right-aligned action row at the bottom of the card. No change needed.

- [ ] Task 3: Tests (Vitest + Vue Test Utils)
  - [ ] 3.1 In `web/src/components/area-weekly-matrix/MatrixCancelPopover.test.ts`,
        add a test that mounts the popover open and asserts the rendered
        `<v-menu>` exposes the new `max-width` (380) and `min-width` (320). The
        existing test file already mocks `cancelBooking` and renders the popover;
        use its mount helper.
  - [ ] 3.2 Add a behavioural test: assert that both
        `[data-cy="matrix-cancel-close"]` AND `[data-cy="matrix-cancel-confirm"]`
        are present in the rendered DOM when the popover is open. (This guards
        against any future regression that drops one of the buttons or
        accidentally adds a clipping overflow.) DOM presence is not the same as
        "visually not clipped", but it is the strongest assertion JSDom permits;
        the visual-clipping check is a manual smoke item below.

- [ ] Task 4: Verification commands
  - [ ] 4.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green.
  - [ ] 4.2 Manual smoke (chrome-devtools-mcp acceptable): start the dev server,
        navigate to the weekly table view on a desktop viewport, scroll so that
        a booked-by-me cell is near the bottom of the viewport, click it, confirm
        the popover renders above the cell (Vuetify auto-flip) and that both the
        "Schliessen" and "Buchung stornieren" buttons are fully visible.

## Dev Notes

### Root cause

The clipping in `private/img_20.png` happens because Vuetify's `<v-menu>`
applies a viewport-height max with internal scroll. When the activator is low in
the viewport, the popover's max-height collapses to the gap between the
activator and the bottom edge, and the action row (rendered last in the
`<v-card>`) falls into the overflow region. The popover gets a scrollbar (in
Vuetify it's often invisible) and the user sees only the top half of the
content.

The fix is to lift the height constraint and trust Vuetify's auto-flip behavior
to position the menu above the activator when there's no room below. Width is
bumped slightly because the current 300px wraps the booker name and date in
some locales.

### Key code locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Popover root | `web/src/components/area-weekly-matrix/MatrixCancelPopover.vue:2–55` | The only file to edit |
| Existing test | `web/src/components/area-weekly-matrix/MatrixCancelPopover.test.ts` | Extend |
| Activator wiring | `<AreaWeeklyMatrixCell>` (sibling file) | Not changed; verify by inspection |

### Anti-patterns to avoid

- Do NOT replace `<v-menu>` with `<v-dialog>`. The popover's anchored positioning
  (next to the clicked cell) is intentional UX — a centered dialog would be a
  regression.
- Do NOT scroll the inner card. The card should always render its content at
  natural height; height management belongs to the v-menu wrapper.
- Do NOT widen the cell to give the popover more room — the cells must stay
  compact for the matrix to be dense.

### Testing standards

- Vitest + Vue Test Utils. Extend the existing
  `MatrixCancelPopover.test.ts`. No new test file.
- The visual clipping behavior is hard to assert in JSDom; the DOM-presence
  assertion plus the explicit prop values are the right testable proxies.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/components/area-weekly-matrix/MatrixCancelPopover.vue:7, 11, 33]
- [Source: web/src/components/area-weekly-matrix/MatrixCancelPopover.test.ts]
- [Source: private/epic-33.md]
- [Source: private/img_20.png]
- [Source: .claude/rules/vue.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/components/area-weekly-matrix/MatrixCancelPopover.test.ts` — 6 tests pass.
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all clean.

### Completion Notes List

- On `<v-menu>` in `MatrixCancelPopover.vue`: bumped `max-width="300"` to
  `380` and added `min-width="320"` + `max-height="none"`. The
  `max-height="none"` lifts Vuetify's default viewport-height max that was
  collapsing the action row out of view when the activator cell sat near the
  bottom of the visible matrix; Vuetify's auto-flip places the menu above the
  cell when there's no room below.
- Inner `<v-card>` and `<v-card-actions>` left as-is. The clipping was an
  external (v-menu container) problem, not internal.
- Test note: the test stub for `v-menu` consumes `maxWidth` as a typed prop
  so the attribute is not reflected to the DOM. The new test verifies
  `min-width` and `max-height` (both pass through `$attrs`) plus DOM presence
  of both action buttons — the strongest assertion JSDom permits for the
  no-clipping contract.

### File List

Frontend (modified):

- `web/src/components/area-weekly-matrix/MatrixCancelPopover.vue` — `<v-menu>`
  gains `min-width="320"`, `max-width="380"`, `max-height="none"`.
- `web/src/components/area-weekly-matrix/MatrixCancelPopover.test.ts` — added
  a regression test asserting both cancel-action buttons are present and the
  popover carries the new size constraints.
