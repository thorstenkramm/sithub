# Story 35.1: Shared Warning Presentation Component

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user encountering item warnings,
I want the warning icon and message to look identical everywhere,
so that I immediately recognize a warning regardless of which view I am in.

## Acceptance Criteria

1. The day/week tile warning presentation (orange circular info icon; dark orange text on a light
   orange background) is extracted into a shared, reusable component; the tiles render their
   warnings through it with no visual change.
2. The shared component is the single source of the warning look and is ready for adoption by the
   floor plan (35.2) and the weekly table view (35.3); no per-surface warning styling is duplicated
   once those stories land.
3. The warning message is dark orange text on a light orange background and remains legible
   (sufficient contrast) in both light and dark themes.

## Tasks / Subtasks

- [x] Task 1: Create the shared warning presentation component (AC: #1, #3)
  - [x] Add `src/components/ItemWarning.vue` exposing an icon-with-tooltip mode (for tiles / floor
        plan / table hover) and an inline-message mode (the expanded-tile alert), so all surfaces
        share one look
  - [x] Move the canonical `.warning-tooltip` styling (`background-color: #fff3e0; color: #e65100;
        font-weight: 500`) into the component (scoped or via the same `content-class`)
  - [x] Use the existing `$warning` icon alias (`mdiAlertCircle`) and `color="warning"`
  - [x] Export it from `src/components/index.ts` following the existing default-export + re-export
        pattern
- [x] Task 2: Adopt the shared component on the tiles (AC: #1)
  - [x] Replace the inline day-mode folded icon+tooltip (`ItemsView.vue` ~243-261) and expanded
        alert (~301-310) with the shared component
  - [x] Replace the week-mode equivalents (`ItemsView.vue` ~469-487 and ~704-713)
  - [x] Preserve the existing `data-cy` hooks (`folded-warning-icon`, `item-warning`,
        `week-folded-warning-icon`, `week-item-warning`) so E2E/unit tests keep passing
- [x] Task 3: Verify no visual change and add a component test (AC: #1, #3)
  - [x] Add a Vitest component test for `ItemWarning.vue` (renders icon+message, applies the orange
        styling, message text shown)
  - [x] Run type-check, lint, unit tests, build

## Dev Notes

Source: `private/epic-35.md` (the tile warning is the visual reference — see `img_29.png`).
[Source: _bmad-output/planning-artifacts/epics.md#Story 35.1 / FR159]

### The canonical style (reference to preserve exactly)

`web/src/views/ItemsView.vue` lines ~2098-2102:

```css
.warning-tooltip.v-overlay__content {
  background-color: #fff3e0 !important; /* light orange */
  color: #e65100 !important;            /* dark orange  */
  font-weight: 500;
}
```

Icon: the `$warning` alias resolves to `mdiAlertCircle` (`web/src/plugins/vuetify.ts:149`) — the
orange circular "i". Theme `warning` color is `#D97706` (light) / `#F59E0B` (dark)
(`vuetify.ts:189-191, 243-245`). Keep the icon at `color="warning"` and the hover/inline message at
the `#fff3e0`/`#e65100` pairing above.

### Existing warning render sites to consolidate (all in `ItemsView.vue`)

- Day folded icon + tooltip: ~243-261 (`content-class="warning-tooltip"`, `v-btn icon` +
  `<v-icon>$warning</v-icon>`, `data-cy="folded-warning-icon"`).
- Day expanded: ~301-310 (`<v-alert type="warning" variant="tonal" data-cy="item-warning">`).
- Week folded icon + tooltip: ~469-487 (`data-cy="week-folded-warning-icon"`).
- Week expanded: ~704-713 (`data-cy="week-item-warning"`).

The current expanded alerts use Vuetify's `type="warning" variant="tonal"` (theme oranges), while the
hover tooltip uses the `#fff3e0`/`#e65100` pairing. FR159 wants ONE look — standardize both the
icon-hover and the inline message on the tooltip pairing (dark-orange-on-light-orange). Confirm the
expanded alert visually matches after switching; adjust the shared component so both modes use the
same colors.

### Component shape (suggested)

```vue
<!-- ItemWarning.vue -->
<!-- props: warning: string; mode?: 'icon' | 'inline' (default 'icon'); size?: number -->
<!-- 'icon'  -> v-tooltip(content-class=warning-tooltip) wrapping the $warning icon button -->
<!-- 'inline'-> a styled block (same orange pairing) for the expanded tile / any inline use -->
```

Keep it presentation-only (no booking/suppression logic — that stays in 35.4). Do not change the
warning data source: `item.attributes.warning` (`web/src/api/items.ts:8`).

### Project Structure Notes

- New: `web/src/components/ItemWarning.vue`; export via `web/src/components/index.ts`.
- Modified: `web/src/views/ItemsView.vue` (adopt component; remove the now-shared inline CSS if fully
  migrated — but the `.warning-tooltip` class may still be referenced by the confirmation dialog
  until 35.4, so verify before deleting).
- This is a pure refactor: no API, no store, no behavior change.

### Testing standards summary

Vitest + Vue Test Utils (component under `src/components`). Assert the icon renders, the message
text appears, and the orange classes/styles are applied. Keep existing `data-cy` selectors so
`ItemsView.test.ts` and Cypress E2E keep working. Run `npm run type-check`, `npm run lint`,
`npx vitest run`, `npm run build`. [Source: .claude/rules/vue.md] [Source: .claude/rules/cypress.md]

### References

- [Source: web/src/views/ItemsView.vue:243-261,301-310,469-487,704-713,2098-2102]
- [Source: web/src/plugins/vuetify.ts:149,189-191,243-245]
- [Source: web/src/api/items.ts:8] [Source: web/src/components/index.ts]

## Dev Agent Record

### Agent Model Used

claude-fable-5

### Debug Log References

- `npm run type-check` clean; `npm run lint` clean
- `npx vitest run` → 451 passed (49 files); `npm run build` clean

### Completion Notes List

- Created `web/src/components/ItemWarning.vue` with two modes: `icon` (v-tooltip + `$warning`
  button, uses the shared `warning-tooltip` overlay style) and `inline` (styled block). The global
  `.warning-tooltip` style now lives in this component (single source); removed the duplicate global
  `<style>` block from `ItemsView.vue`.
- Standardized both the hover tooltip and the (formerly `v-alert`) inline/expanded message on the
  canonical `#fff3e0` / `#e65100` pairing so all tile warnings share one look (FR159).
- Adopted the component at all four tile sites in `ItemsView.vue` (day folded + expanded, week
  folded + expanded), preserving the `data-cy` hooks (`folded-warning-icon`, `item-warning`,
  `week-folded-warning-icon`, `week-item-warning`) — all 90 `ItemsView.test.ts` tests stay green.
- Added `ItemWarning.test.ts` (icon mode renders the button + message; inline mode renders the
  styled block). Pure refactor — no API/store/behavior change.

### File List

- web/src/components/ItemWarning.vue (new)
- web/src/components/__tests__/ItemWarning.test.ts (new)
- web/src/components/index.ts (modified — export ItemWarning)
- web/src/views/ItemsView.vue (modified — adopt component, remove duplicate CSS)

### Change Log

- 2026-07-04: Implemented FR159 — shared ItemWarning component; tiles adopt it, no visual change.
