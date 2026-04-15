# Story 29.3: Weekly Matrix Layout, Sticky Orientation & Room State

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want a dense but readable weekly matrix with collapsible room sections,
so that I can scan the whole floor quickly without losing orientation.

## Acceptance Criteria

1. **Given** I open `Table view`
   **When** the matrix renders
   **Then** it shows all subareas of the selected area in one long table
   **And** the subareas and desks appear in the exact configured SitHub order

2. **Given** I view the matrix
   **When** I scroll vertically
   **Then** the weekday header remains sticky
   **And** the left desk-name column remains sticky

3. **Given** I open the matrix for the first time
   **When** the room sections render
   **Then** all rooms are expanded by default

4. **Given** I collapse or expand a room using its dedicated chevron
   **When** I reopen the table later or switch to another week
   **Then** the previous collapsed state is restored from local storage

5. **Given** a room is collapsed
   **When** I look at its header
   **Then** I see compact occupied counts for each visible day of the current week

6. **Given** the selected week contains past days
   **When** the matrix renders
   **Then** those past-day columns stay visible
   **And** their cells are visually muted and non-interactive

## Tasks / Subtasks

- [x] Task 1: Build the table-view matrix component on top of Story 29.1 API data (AC: #1)
  - [x] 1.1 Create or expand a purpose-built matrix component using Vue + Vuetify primitives;
        do **not** use `v-data-table`
  - [x] 1.2 Consume `fetchWeeklyMatrix()` from `web/src/api/itemGroupMatrix.ts`
  - [x] 1.3 Render one room section per item-group resource and one desk row per item while
        preserving backend order exactly
  - [x] 1.4 Keep the component scoped to desktop use inside the `ItemGroupsView` table shell
  - [x] 1.5 Add loading and error states for the matrix itself so the rest of the page remains
        stable during refreshes

- [x] Task 2: Implement sticky orientation and dense week columns (AC: #2, #6)
  - [x] 2.1 Render a dedicated weekday header row from the `days[]` metadata returned by Story
        29.1; use the page week/weekend state only to parameterize the matrix request
  - [x] 2.2 Keep the weekday header sticky while vertically scrolling
  - [x] 2.3 Keep the left desk-name column sticky so users never lose row context
  - [x] 2.4 Compare each column date against today and mark past-day columns/cells as muted and
        non-interactive
  - [x] 2.5 Keep columns compact enough to avoid horizontal scrolling where possible without
        collapsing the header into ambiguity

- [x] Task 3: Add room collapse behavior and persisted state (AC: #3, #4, #5)
  - [x] 3.1 Use a dedicated chevron/control in each room header; do not collapse on header click
  - [x] 3.2 Default all rooms to expanded when no prior preference exists
  - [x] 3.3 Persist collapsed room IDs in localStorage, scoped by area ID
  - [x] 3.4 Restore the same collapsed state after a week change and on later visits
  - [x] 3.5 Compute collapsed-room summaries as occupied counts per visible day for the current
        week

- [x] Task 4: Establish the component structure for later cell-interaction stories
  - [x] 4.1 Break the matrix into maintainable subcomponents or clearly isolated sections, e.g.:
        - matrix container
        - room section/header
        - desk row
        - matrix cell
  - [x] 4.2 Pass enough props/state through this structure so Stories 29.4-29.6 can add richer
        cell rendering and interactions without rewriting the layout

- [x] Task 5: Tests and validation
  - [x] 5.1 Add component tests covering:
        - room and desk ordering
        - sticky header/left-column marker classes
        - default expanded state
        - collapse persistence across remounts
        - collapsed summary counts
        - past-day muted rendering
  - [x] 5.2 Run `cd web && npx vitest run` for the new matrix component tests
  - [x] 5.3 Run `cd web && npm run type-check`
  - [x] 5.4 Run `cd web && npm run lint`

## Dev Notes

### Architecture & Patterns

- This story owns structure, orientation, and room-state persistence.
- It should not introduce booking or cancellation popovers yet; those belong to Stories 29.5 and
  29.6.
- The UX design explicitly chose a purpose-built custom matrix component using Vuetify primitives
  instead of a generic data-table abstraction.

### Recommended Component Boundary

Recommended structure:
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue`
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRoomSection.vue`
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue`
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue`

This is a recommendation, not a hard requirement. The important constraint is keeping layout and
interaction concerns separable for the later stories.

### Key Code Locations

| Element | Location | Why it matters |
|---------|----------|----------------|
| Area page shell | `web/src/views/ItemGroupsView.vue` | Hosts the table view |
| Story 29.1 API client | `web/src/api/itemGroupMatrix.ts` | Source of matrix data and authoritative `days[]` column metadata |
| Week selection context | `web/src/composables/useWeekSelector.ts` | Drives request params for week selection, not the rendered header truth |
| Weekend preference | `web/src/composables/useWeekendPreference.ts` | 5-day vs 7-day matrix |
| Safe localStorage | `web/src/composables/storage.ts` | Collapse persistence |
| UX decision | `_bmad-output/planning-artifacts/ux-design-specification.md#Design System Foundation` | No generic data-table |

### Implementation Strategy

1. Land the matrix container and fetch cycle.
2. Treat the matrix API `days[]` metadata as the source of truth for rendered columns and basic
   cell placeholders.
3. Add sticky positioning and the scroll shell.
4. Add collapse persistence and collapsed summaries.
5. Keep cell rendering simple enough that Story 29.4 can enrich it without tearing up the layout.

### Anti-Patterns to Avoid

- Do NOT use `v-data-table` or another generic table abstraction.
- Do NOT reorder rooms or desks on the frontend.
- Do NOT hide past-day columns; mute them instead.
- Do NOT make the entire room header clickable; use a dedicated chevron/control.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 29 Stories: Desktop Weekly Table View]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Core User Experience]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Design System Foundation]
- [Source: web/src/api/itemGroupMatrix.ts]
- [Source: web/src/views/ItemGroupsView.vue]
- [Source: web/src/composables/useWeekSelector.ts]
- [Source: web/src/composables/useWeekendPreference.ts]
- [Source: web/src/composables/storage.ts]

## Dev Agent Record

### Agent Model Used

GPT-5

### Debug Log References

- Story creation only; no implementation logs yet

### Completion Notes List

- Story assumes Story 29.1 delivers the weekly matrix endpoint and client wrapper first
- Layout story intentionally stops short of booking/cancellation interactions

### File List

- _bmad-output/implementation-artifacts/29-3-weekly-matrix-layout-sticky-orientation-and-room-state.md (new story file)
