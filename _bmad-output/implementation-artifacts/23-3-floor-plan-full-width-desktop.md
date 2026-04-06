# Story 23.3: Floor Plan Full-Width Desktop Layout

Status: review

## Story

As a user,
I want the floor plan to use the full available width on desktop,
so that I can see floor plan details without unnecessary whitespace.

## Acceptance Criteria

1. **Given** I am viewing a floor plan on a desktop viewport (>= 960px)
   **When** the floor plan dialog renders
   **Then** the dialog uses the full available width of the viewport

2. **Given** I am viewing a floor plan on a mobile viewport
   **When** the floor plan renders
   **Then** the existing fullscreen mobile behavior is unchanged

## Tasks / Subtasks

- [x] Task 1: Remove max-width constraint from floor plan dialogs (AC: 1, 2)
  - [x] 1.1 In `web/src/views/ItemGroupsView.vue` line 268: remove
    `max-width="1100"` from the `v-dialog` for the area floor plan.
  - [x] 1.2 In `web/src/views/ItemsView.vue` line 753: remove
    `max-width="1100"` from the `v-dialog` for the item group floor plan.
- [x] Task 2: Write tests (AC: 1)
  - [x] 2.1 In `web/src/views/ItemGroupsView.test.ts`: add test verifying
    the floor plan dialog has no max-width attribute on desktop.
  - [x] 2.2 In `web/src/views/ItemsView.test.ts`: add same test for the
    item group floor plan dialog.
- [x] Task 3: Run full test suite and linters (AC: 1, 2)
  - [x] 3.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`,
    `npm run build`
  - [x] 3.2 Fix any failures

## Dev Notes

### Root Cause

The floor plan dialog has `max-width="1100"` as a Vuetify prop on `v-dialog`
in both ItemGroupsView and ItemsView. On wide desktop screens this caps the
dialog at 1100px, leaving large margins on both sides.

### Width Constraint Chain

```text
App.vue → v-main (no constraint)
  → View (.page-container, max-width: 1400px)
    → v-dialog (max-width="1100") ← THE CONSTRAINT
      → v-card.floor-plan-dialog-card (height: 100%)
        → InteractiveFloorPlan (no width constraint)
```

### Fix

Remove `max-width="1100"` from both `v-dialog` elements. Vuetify dialogs
without max-width expand to fill the viewport. The existing `fullscreen`
prop on compact viewports (`isCompactFloorPlanViewport`) is unchanged.

### Affected Dialogs

| View | Line | data-cy | Prop to remove |
|------|------|---------|---------------|
| ItemGroupsView.vue | ~268 | `floor-plan-dialog` | `max-width="1100"` |
| ItemsView.vue | ~753 | `item-group-floor-plan-dialog` | `max-width="1100"` |

### Existing Tests

- `ItemGroupsView.test.ts:344` — tests fullscreen on compact viewports
- `ItemsView.test.ts:947` — tests fullscreen on compact viewports

These test the mobile behavior which is unchanged. New tests verify the
desktop behavior (no max-width constraint).

### Stub for v-dialog in Tests

Both test files stub `v-dialog` as:

```js
'v-dialog': {
  props: ['modelValue', 'fullscreen', 'persistent'],
  template: '<div v-if="modelValue" v-bind="$attrs" ...><slot /></div>'
}
```

The stub does NOT include `maxWidth` in props, so the Vuetify `max-width`
prop is passed through as an HTML attribute. After removal, the attribute
simply won't be present — tests can assert its absence.

### Files to Modify

- `web/src/views/ItemGroupsView.vue` — remove max-width from dialog
- `web/src/views/ItemsView.vue` — remove max-width from dialog
- `web/src/views/ItemGroupsView.test.ts` — add no-max-width test
- `web/src/views/ItemsView.test.ts` — add no-max-width test

### Do NOT Change

- Mobile fullscreen behavior (isCompactFloorPlanViewport)
- InteractiveFloorPlan component internals
- `.page-container` max-width (affects all views, not just floor plans)
- Floor plan dialog card/body CSS

### References

- [Source: private/epic-23.md — "Floor plan on desktop" section with screenshot]
- [Source: web/src/views/ItemGroupsView.vue:268 — max-width="1100"]
- [Source: web/src/views/ItemsView.vue:753 — max-width="1100"]
- [Source: web/src/styles/global.css:200 — .page-container max-width]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Task 1: Removed `max-width="1100"` from v-dialog in both ItemGroupsView.vue
  (area floor plan) and ItemsView.vue (item group floor plan). Mobile fullscreen
  behavior unchanged.
- Task 2: Added tests in both test files asserting the dialog has no max-width
  or maxwidth attribute on desktop.
- Task 3: All 285 tests pass. Type-check, lint, build clean.

### File List

- `web/src/views/ItemGroupsView.vue` (modified — removed max-width from dialog)
- `web/src/views/ItemsView.vue` (modified — removed max-width from dialog)
- `web/src/views/ItemGroupsView.test.ts` (modified — added no-max-width test)
- `web/src/views/ItemsView.test.ts` (modified — added no-max-width test)
