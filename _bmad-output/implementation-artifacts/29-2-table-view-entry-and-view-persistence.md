# Story 29.2: Table View Entry & View Persistence

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a desktop user,
I want a Table view option on the area item-groups page and for SitHub to remember my last
desktop view,
so that I can return directly to the weekly matrix without extra setup.

## Acceptance Criteria

1. **Given** I open `/areas/:areaId/item-groups` on a desktop viewport
   **When** the page renders
   **Then** I see a `Table view` action alongside the existing area actions in the same control
   area as `Floor plan`

2. **Given** I open `/areas/:areaId/item-groups` on a mobile viewport
   **When** I look at the `Table view` action
   **Then** it is disabled
   **And** hovering or long-pressing it explains that the table view is available on desktop only

3. **Given** I switch from the default card view to `Table view` on desktop
   **When** I leave the page and come back later to the same area
   **Then** SitHub restores `Table view` as the active desktop view for that area context

## Tasks / Subtasks

- [x] Task 1: Add the desktop-only table-view action to the area page (AC: #1, #2)
  - [x] 1.1 Update `web/src/views/ItemGroupsView.vue` to render a `Table view` action in the same
        action row that currently contains the `Floor plan` button
  - [x] 1.2 Reuse the existing viewport detection approach already present in `ItemGroupsView.vue`
        so the button can be enabled on desktop and disabled on compact/mobile viewports
  - [x] 1.3 Implement the disabled-mobile tooltip using a wrapper element, since Vuetify tooltips
        do not attach reliably to disabled buttons directly
  - [x] 1.4 Add stable `data-cy` selectors for:
        - desktop table-view button
        - mobile-disabled state
        - explanatory tooltip text
  - [x] 1.5 Add i18n strings for the new button label and mobile explanation in all supported
        locale files

- [x] Task 2: Persist the selected area view per area in local storage (AC: #3)
  - [x] 2.1 Add a small persistence helper or composable that uses `getSafeLocalStorage()` rather
        than raw `window.localStorage`
  - [x] 2.2 Scope the persisted preference by area ID so one area can remember `Table view`
        without forcing the same choice onto every area
  - [x] 2.3 Default to the current card view when no preference is stored
  - [x] 2.4 Restore the memorized view on mount only for desktop-capable viewports; compact/mobile
        viewports must not activate the table view

- [x] Task 3: Introduce the view-switching shell for the weekly matrix (AC: #1, #3)
  - [x] 3.1 Add a dedicated component boundary for the table view at
        `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue`, even if the initial
        shell is lightweight until Story 29.3
  - [x] 3.2 Keep the existing week selector shared above both views so week context remains
        shared; do not expand the current equipment-filter behavior into the matrix in this story
  - [x] 3.3 Ensure switching between card view and table view happens in-place on the same route;
        do not introduce a new route or page

- [x] Task 4: Cover the behavior with focused tests
  - [x] 4.1 Extend `web/src/views/ItemGroupsView.test.ts` for:
        - desktop button visible
        - mobile button disabled
        - tooltip text rendered
        - memorized view restored for the same area
        - another area not inheriting the stored choice by accident
  - [x] 4.2 If a new persistence composable/helper is introduced, add its own unit tests

- [x] Task 5: Validation
  - [x] 5.1 Run `cd web && npx vitest run src/views/ItemGroupsView.test.ts`
  - [x] 5.2 Run `cd web && npm run type-check`
  - [x] 5.3 Run `cd web && npm run lint -- src/views/ItemGroupsView.vue`

### Review Findings

- [x] [Review][Patch] Mobile-disabled table switch has no explanatory tooltip [web/src/views/ItemGroupsView.vue:36]

## Dev Notes

### Architecture & Patterns

- This story is the table-view entry point and state-restoration shell; it does **not** implement
  the full matrix layout yet.
- The current page already owns:
  - selected week
  - weekend preference
  - area floor plan button
  - compact viewport detection
- The current equipment filter is tied to the card-grid browsing flow. Do not invent matrix-filter
  semantics in this story just because both views share the same route shell.
- Keep the table view on the existing `/areas/:areaId/item-groups` route. The UX spec explicitly
  places it beside `Floor plan`, not behind a new route or nested page.

### Key Code Locations

| Element | Location | Why it matters |
|---------|----------|----------------|
| Area item-groups page | `web/src/views/ItemGroupsView.vue` | Main integration point |
| Current viewport detection | `web/src/views/ItemGroupsView.vue` | Reuse for desktop-only gating |
| Week persistence | `web/src/composables/useDateState.ts` | Existing session-backed week memory |
| Weekend preference | `web/src/composables/useWeekendPreference.ts` | Existing localStorage preference |
| Safe localStorage helper | `web/src/composables/storage.ts` | Use for view persistence |
| Existing tests | `web/src/views/ItemGroupsView.test.ts` | Extend rather than creating a parallel test style |

### Implementation Strategy

1. Add a small `activeAreaView` state on `ItemGroupsView.vue`.
2. Persist it in localStorage keyed by area ID.
3. Keep the current card-grid view as the fallback/default.
4. Mount a dedicated table-view shell component when the table view is active so Story 29.3 can
   expand that same component instead of replacing ad hoc markup later.

### Anti-Patterns to Avoid

- Do NOT create a separate `/table` route.
- Do NOT hide or fork the existing week selector for table view.
- Do NOT silently apply the current equipment filter to the matrix without explicit UX rules.
- Do NOT access `window.localStorage` directly when `getSafeLocalStorage()` is already available.
- Do NOT rely on the disabled button itself as the tooltip activator.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 29 Stories: Desktop Weekly Table View]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Executive Summary]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Core User Experience]
- [Source: web/src/views/ItemGroupsView.vue]
- [Source: web/src/views/ItemGroupsView.test.ts]
- [Source: web/src/composables/useDateState.ts]
- [Source: web/src/composables/useWeekendPreference.ts]
- [Source: web/src/composables/storage.ts]

## Dev Agent Record

### Agent Model Used

GPT-5

### Debug Log References

- Story creation only; no implementation logs yet

### Completion Notes List

- Planning assumes the current card view remains the default fallback
- Story intentionally creates a reusable table-view shell so Story 29.3 can expand it cleanly

### File List

- _bmad-output/implementation-artifacts/29-2-table-view-entry-and-view-persistence.md (new story file)
