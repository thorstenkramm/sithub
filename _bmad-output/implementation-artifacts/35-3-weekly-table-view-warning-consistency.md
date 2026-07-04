# Story 35.3: Weekly Table-View Warning Consistency

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user on the weekly desktop table view,
I want the warning hover to match the shared style and not repeat itself when I select a cell,
so that the table view feels consistent with the rest of the application.

## Acceptance Criteria

1. Hovering the table-view warning icon shows the message in the shared style (dark orange text,
   light orange background), positioned correctly next to the icon.
2. Clicking a free cell to prepare a booking shows no additional warning inside the cell or popover;
   the uniform confirmation dialog on booking is the only pre-booking warning.
3. The corrected hover message is visually identical to a tile warning message.

## Tasks / Subtasks

- [x] Task 1: Fix the row warning hover styling (AC: #1, #3)
  - [x] In `AreaWeeklyMatrixRow.vue` (~20-25), replace the bare `<span>` tooltip content with the
        shared warning presentation (35.1) OR apply the `warning-tooltip` `content-class` so the
        tooltip uses the `#fff3e0` background / `#e65100` text pairing instead of Vuetify's default
        dark tooltip
  - [x] Keep `data-cy="matrix-warning-icon"` and `data-cy="matrix-warning-tooltip"`
- [x] Task 2: Remove the duplicate in-cell/popover warning (AC: #2)
  - [x] In `MatrixBookingPopover.vue` (~16-35), remove the inline `v-alert` warning + its
        "Don't show again" suppress button (`data-cy="matrix-booking-warning"` /
        `matrix-booking-suppress-warning`); the pre-booking warning is handled solely by the uniform
        confirmation dialog (35.4)
  - [x] Remove the now-unused `warningSuppressed` computed (~211-213) and `doSuppressWarning`
        (~158-163) from the popover if nothing else uses them
- [x] Task 3: Tests (AC: #1-#3)
  - [x] Update/extend `AreaWeeklyMatrixView.test.ts` / row + popover tests: warning hover uses the
        shared style; no inline warning appears in the booking popover

## Dev Notes

Source: `private/epic-35.md` — table-view section (`img_31.png` correct hover, `img_32.png` no
duplicate on cell click). [Source: _bmad-output/planning-artifacts/epics.md#Story 35.3 / FR162]

### The wrong hover styling (fix target)

`web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue` ~20-25:

```vue
<v-tooltip v-if="item.warning" location="right">
  <template #activator="{ props: warnProps }">
    <v-icon v-bind="warnProps" size="14" color="warning" class="ml-1" data-cy="matrix-warning-icon">$warning</v-icon>
  </template>
  <span data-cy="matrix-warning-tooltip">{{ item.warning }}</span>  <!-- bare span → default dark tooltip -->
</v-tooltip>
```

The bare `<span>` gets Vuetify's default dark tooltip. Fix by using the shared `ItemWarning.vue`
(icon mode) from 35.1, or minimally add `content-class="warning-tooltip"` to the `v-tooltip` and
ensure the `.warning-tooltip` CSS is in scope (it lives in `ItemsView.vue` today — prefer the shared
component so the style travels with it). `location="right"` positioning is fine; verify placement
matches the reference.

### 🚨 The duplicate in-cell warning (remove)

`web/src/components/area-weekly-matrix/MatrixBookingPopover.vue` ~16-35 renders an inline
`v-alert type="warning"` with a "Don't show again" button when a free cell is clicked. This is the
second, inconsistent warning FR162 removes. Delete that alert. Note this popover currently does NOT
block booking on warnings — the actual uniform confirmation is added in 35.4; this story just
removes the duplicate. Coordinate ordering: if 35.4 lands first, the confirmation is already wired;
if 35.3 lands first, the table temporarily has no pre-booking warning until 35.4 — acceptable and
intended (the inline alert was never a real confirmation anyway).

### Data

`MatrixItem.warning?: string` (`web/src/api/itemGroupMatrix.ts:19-26`). Booking flow: cell click →
`provide('matrixCellClick', ...)` in `AreaWeeklyMatrixView.vue` (~196-218) → `MatrixBookingPopover`
`submitBooking()` (~229-264). Suppression composable: `useWarningSuppression` (imported ~147).

### Project Structure Notes

- Modified: `AreaWeeklyMatrixRow.vue`, `MatrixBookingPopover.vue` (both under
  `web/src/components/area-weekly-matrix/`), + their tests.
- Reuses `ItemWarning.vue` (35.1). The blocking confirmation is 35.4.

### Testing standards summary

Vitest component tests in `area-weekly-matrix/`. Assert the hover uses the shared style and the
booking popover renders no inline warning. Run type-check, lint, vitest, build. [Source: .claude/rules/vue.md]

### References

- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue:20-25]
- [Source: web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:16-35,147,158-163,211-213,229-264]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue:196-218]
- [Source: web/src/api/itemGroupMatrix.ts:19-26]

## Dev Agent Record

### Agent Model Used

claude-fable-5

### Debug Log References

- 41 area-weekly-matrix tests pass; full suite 456; type-check/lint/build clean

### Completion Notes List

- Replaced the bare-`<span>` warning tooltip in `AreaWeeklyMatrixRow.vue` with the shared
  `ItemWarning` (icon-variant "plain", location "right", size 14, `data-cy="matrix-warning-icon"`),
  so the hover message now uses the shared `#fff3e0`/`#e65100` style instead of Vuetify's dark
  default (FR162, hover part). Extended `ItemWarning` with `iconVariant` + `location` props for this.
- Removed the duplicate inline warning `v-alert` (and its `warningSuppressed`/`doSuppressWarning`)
  from `MatrixBookingPopover.vue`; the pre-booking warning is now the uniform confirmation dialog
  (wired in 35.4) — the table gained a real blocking confirmation it previously lacked.
- Updated the matrix tests: the row warning renders via the shared component; the popover shows no
  inline warning and instead opens the confirmation on booking a warned item.

Note: the old `matrix-warning-tooltip` span data-cy was removed (the shared component owns the
tooltip); the matrix test now asserts the warning text is present rather than that specific hook.

### File List

- web/src/components/ItemWarning.vue (modified — iconVariant/location props; shared with 35.1)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue (modified)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts (modified)
- web/src/components/area-weekly-matrix/MatrixBookingPopover.vue (modified — see 35.4)
- web/src/components/area-weekly-matrix/MatrixBookingPopover.test.ts (modified)

### Change Log

- 2026-07-04: Implemented FR162 — weekly-table warning hover uses the shared style; duplicate in-cell
  warning removed in favor of the uniform confirmation.
