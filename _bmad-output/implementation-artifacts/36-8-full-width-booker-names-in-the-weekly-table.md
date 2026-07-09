# Story 36.8: Full-Width Booker Names in the Weekly Table

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user reading the weekly table,
I want colleagues' names shown as fully as the cell allows,
so that I can recognize people by their first names.

## Acceptance Criteria

1. Weekly-table booker-name cells use the maximum available cell width with decent padding; the
   previous 60px max-width cap is removed.
2. A name that fits shows the full first and last name.
3. A name too long for the cell is truncated from the END (ellipsis) so the first name stays
   visible.

## Tasks / Subtasks

- [x] Task 1: Show the full booker name (first + last) in the occupied cell (AC: #2)
  - [x] In `AreaWeeklyMatrixCell.vue`, the occupied cell currently renders `shortName`
        (`{{ shortName }}` at template ~57), which is `getShortName(cell.booker_name)` (~87). That
        helper abbreviates the first name to an initial ("Ada Lovelace" -> "A. Lovelace",
        `utils/text.ts:20-26`) and hard-caps length at 14. That conflicts with AC #2 ("full first
        and last name"). Replace the displayed value with the full `cell.booker_name` so the first
        name is shown in full; let CSS width + ellipsis (Task 2) handle overflow instead of the
        helper's length cap.
  - [x] Keep `initials` (used only in the past-day branch, template ~10 / ~86) unchanged — AC #2/#3
        are about the interactive occupied cell, not the muted past-day pill.
  - [x] `booker_name` is the full display name already present on the cell (`itemGroupMatrix.ts:13`,
        `booker_name?: string`) — no API or type change and no new prop threading is required.
- [x] Task 2: Remove the 60px cap, use available width + padding, truncate from the END (AC: #1, #3)
  - [x] In `AreaWeeklyMatrixCell.vue` `<style scoped>`, edit the `.cell-short-name` selector
        (~199-207): delete `max-width: 60px;` (~206). The `<td>.matrix-cell` already sets
        `min-width: 80px` and `padding: 4px` (~118-119) and `.cell-content` adds `padding: 2px 4px`
        (~132), so the name now uses the full cell width with decent padding.
  - [x] Keep `white-space: nowrap; overflow: hidden; text-overflow: ellipsis;` (~203-205). With no
        `max-width` and the default LTR text direction, `text-overflow: ellipsis` truncates from the
        END, so the last name is dropped first and the first name stays visible (AC #3). Do NOT add
        `direction: rtl` — that would truncate from the front.
  - [x] Add `min-width: 0;` to `.cell-short-name` if needed so the span can shrink inside the
        flex `.cell-content` (~125-133) and actually ellipsis instead of overflowing; verify in the
        browser (see visual notes).
  - [x] Leave the avatar branch (`.cell-avatar`, template ~46-56) as-is; when an avatar is shown the
        name sits beside it and must still ellipsis within the remaining width.
- [x] Task 3: Tests (AC: #1, #2, #3)
  - [x] Update the existing Vitest expectation in `AreaWeeklyMatrixView.test.ts:471` — the
        occupied-cell test asserts `matrix-cell-initials` text equals `'A. Lovelace'`; after Task 1
        it must equal the full name `'Ada Lovelace'`. The tooltip assertion (~472,
        `matrix-cell-tooltip` == `'Ada Lovelace'`) stays unchanged.
  - [x] Add a case: a long booker name (e.g. `'Alexander Seidemann-Klamant'`) renders in full in the
        `matrix-cell-initials` span text (JSDOM does not compute pixel truncation, so assert the raw
        text is the full name and rely on the CSS class for visual truncation).
  - [x] Run `npm run type-check`, `npm run lint`, `npx vitest run`, and `npm run build`.

## Dev Notes

Source: [Source: _bmad-output/planning-artifacts/epics.md#Story 36.8 / FR177]. FR177 text:
booker-name cells use the maximum available cell width (with decent padding) and show the full first
and last name; when the name does not fit it is truncated from the end so the first name stays
visible; the previous 60px max-width cap is removed
[Source: _bmad-output/planning-artifacts/epics.md:647-651].

### Where the cell renders and what it shows today

The weekly table renders one `AreaWeeklyMatrixCell.vue` per day/item. The occupied branch shows an
optional avatar plus a name span:

- Template: `<span class="cell-short-name text-caption" data-cy="matrix-cell-initials">{{ shortName }}</span>`
  [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:57].
- `shortName = computed(() => getShortName(props.cell.booker_name))`
  [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:87].

Note the `data-cy` is `matrix-cell-initials` for historical reasons even though it holds a name — do
not rename it (tests and any E2E depend on it).

### The blocker for AC #2: getShortName abbreviates the first name

`getShortName` returns first-initial + last name and caps at 14 chars:
`"Thorsten Kramm" -> "T. Kramm"`, `"Ada Lovelace" -> "A. Lovelace"`
[Source: web/src/utils/text.ts:14-26]. AC #2 requires the FULL first and last name, so the cell must
display `cell.booker_name` directly rather than the abbreviated `shortName`. Overflow handling then
becomes purely a CSS concern (Task 2). `getInitials` is a separate helper used by the past-day pill
[Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:86,10] and stays as-is.

### The 60px cap and truncation direction

The cap lives in one selector:

```css
.cell-short-name {
  font-weight: 600;
  font-size: 0.7rem;
  color: rgb(var(--v-theme-primary));
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 60px;   /* <- remove this line */
}
```

[Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:199-207]. Removing
`max-width: 60px` lets the name expand to the cell width. The cell (`<td>`) has `min-width: 80px`,
`padding: 4px` [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:118-119] and
the inner `.cell-content` flex container adds `padding: 2px 4px`
[Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:125-133] — that is the
"decent padding" from AC #1. `text-overflow: ellipsis` on LTR text truncates from the END by default,
satisfying AC #3 with no extra properties.

### Data / API

No backend or type change. `booker_name` is already the full name on the cell resource
[Source: web/src/api/itemGroupMatrix.ts:10-17] and is what the tooltip already renders in full
[Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:60].

### Related but out of scope

- The desk-label column has its own `max-width: 180px` ellipsis in
  `AreaWeeklyMatrixRow.vue:127-135` — that is the item name, not a booker name; do not touch it.
- Tile/floor-plan views use `getShortName`/`middleTruncate` elsewhere; this story only changes the
  weekly-table cell.

### Project Structure Notes

- Modified: `web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue` (template line ~57 +
  CSS `.cell-short-name`).
- Modified: `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts` (update the
  `'A. Lovelace'` expectation to `'Ada Lovelace'`, add a long-name case).
- No changes to `getShortName`/`getInitials` (`web/src/utils/text.ts`) — they remain in use by other
  views and by the past-day pill.

### Testing standards summary

Vitest component tests for the matrix live in
`web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts`. JSDOM does not lay out pixels,
so assert on the rendered text (full name) and trust the CSS class for the visual ellipsis. Run
type-check, lint, `vitest run`, build. Keep coverage >= 75%. [Source: .claude/rules/vue.md]

Visual verification (Chrome DevTools MCP, [Source: .claude/rules/vue.md]): open the weekly table for
an area with occupied cells, confirm (a) short names like "Ada Lovelace" show in full, (b) a long
name fills the cell width with left/right padding and (c) an overly long name shows a trailing
ellipsis with the first name still readable. Screenshot before/after for the change log.

### References

- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:57,86-87,118-119,125-133,199-207]
- [Source: web/src/utils/text.ts:14-40]
- [Source: web/src/api/itemGroupMatrix.ts:10-17]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts:447-473]
- [Source: _bmad-output/planning-artifacts/epics.md:647-651,5537-5558]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `cd web && npm run type-check` — clean (vue-tsc --noEmit, no errors).
- `cd web && npm run lint` — clean (eslint --max-warnings 0). Removed the now-unused `getShortName`
  import and `shortName` computed to keep lint green.
- `cd web && npx vitest run src/components/area-weekly-matrix` — 3 files, 45 tests passed.
- `cd web && npm run build` — built successfully.

### Completion Notes List

- Occupied-cell name span now renders the full `cell.booker_name` instead of the abbreviated
  `getShortName(...)` output, satisfying AC #2 (full first + last name).
- Removed `max-width: 60px` from `.cell-short-name` and added `min-width: 0` so the span shrinks
  inside the flex `.cell-content` and truncates with the existing LTR `text-overflow: ellipsis`
  (drops the last name first, keeps the first name visible — AC #1, #3). No `direction: rtl` added.
- `getShortName` import and the `shortName` computed became unused after the template change and were
  removed; `getInitials`/`initials` remain for the past-day pill (unchanged, per Task 1).
- No API/type change: `booker_name?: string` already carries the full name on `MatrixCell`.
- Updated the existing occupied-cell test expectation from `'A. Lovelace'` to `'Ada Lovelace'`
  (tooltip assertion unchanged) and added a long-name case (`'Alexander Seidemann-Klamant'`) that
  asserts the full text renders and the span carries the `cell-short-name` class for visual
  truncation (JSDOM does not compute pixel layout).

### File List

- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue` (modified)
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts` (modified)

## Change Log

- 2026-07-09: Story 36.8 (FR177) implemented. Weekly-table booker-name cell now shows the full
  `booker_name` at full available cell width; removed the 60px cap and added `min-width: 0` so the
  name truncates from the end (last name first) via the existing ellipsis. Dropped the unused
  `getShortName` usage in this cell. Updated Vitest expectation to the full name and added a
  long-name rendering case. Gates (type-check, lint, matrix vitest, build) all pass.
