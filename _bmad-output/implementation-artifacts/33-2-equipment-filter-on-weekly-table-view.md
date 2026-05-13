# Story 33.2: Equipment Filter on Weekly Table View

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user on the weekly desktop table view,
I want the equipment filter at the top of the table to actually filter rows,
so that I can find desks by equipment without leaving the table.

## Acceptance Criteria

1. **Given** I am on the weekly desktop table view and I type a keyword in the
   equipment filter input
   **When** the filter is applied
   **Then** every row whose item does not match the parsed filter is visually dimmed
   and its cells become non-interactive
   **And** rows whose items match remain at full opacity and remain interactive

2. **Given** the filter is applied
   **When** I clear the input (via the clear icon or by backspacing)
   **Then** every row returns to its normal state

3. **Given** the filter matches every item in the table
   **When** the filter is applied
   **Then** no row is dimmed

4. **Given** the filter applies to the table view
   **When** I switch back and forth between the table and card views
   **Then** the filter input value is preserved within the same session per the
   existing saved-filters behavior

## Tasks / Subtasks

- [ ] Task 1: Inherit Story 33.1's `?? ''` null-guard fix here (AC: #2)
  - [ ] 1.1 `web/src/views/ItemGroupsView.vue` has the SAME `equipmentFilter` /
        `parsedEquipmentFilter` pattern as `ItemsView.vue`, and therefore the SAME
        v-combobox clearable null bug. At line 380 the computed is:
        ```ts
        const parsedEquipmentFilter = computed(() => parseFilter(equipmentFilter.value));
        ```
        Change to `parseFilter(equipmentFilter.value ?? '')` exactly as in Story 33.1.
        Also add the same defensive watcher just below line 355
        (`const equipmentFilter = ref('');`).
  - [ ] 1.2 The line 394 guard `if (!equipmentFilter.value) return false;` already
        short-circuits for empty/null — keep it; together with the `?? ''` fix this
        is now fully robust.

- [ ] Task 2: Pipe the parsed filter down to `AreaWeeklyMatrixView` (AC: #1, #3)
  - [ ] 2.1 In `web/src/views/ItemGroupsView.vue` at line 136, add a new prop to
        the `<AreaWeeklyMatrixView>` element:
        ```vue
        <AreaWeeklyMatrixView
          v-else-if="activeView === 'table'"
          :area-id="route.params.areaId as string"
          :week="selectedWeek"
          :show-weekends="showWeekends"
          :parsed-equipment-filter="parsedEquipmentFilter"
        />
        ```
  - [ ] 2.2 In `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue`,
        accept the new prop:
        ```ts
        import type { AndGroup } from '../../composables/useEquipmentFilter';
        const props = defineProps<{
          areaId: string;
          week: string;
          showWeekends: boolean;
          parsedEquipmentFilter?: AndGroup[];
        }>();
        ```
        Default to `[]` if not provided so the matrix continues to render correctly
        when mounted from a context that does not pass the prop.
  - [ ] 2.3 Forward `parsedEquipmentFilter` to each child section / row component
        the matrix already renders (see `AreaWeeklyMatrixRoomSection.vue` /
        `AreaWeeklyMatrixRow.vue`). The exact pass-through depends on the existing
        prop chain — read the current template to confirm, but the pattern is the
        same as the `showWeekends` prop already threads through.

- [ ] Task 3: Compute "row is filtered out" inside the row component (AC: #1, #3)
  - [ ] 3.1 In `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue`,
        accept the parsed filter as a prop and compute a `isFilteredOut` boolean:
        ```ts
        import { matchesParsedFilter, type AndGroup } from '../../composables/useEquipmentFilter';
        const props = defineProps<{
          item: MatrixItem;
          // ...existing props...
          parsedEquipmentFilter?: AndGroup[];
        }>();
        const isFilteredOut = computed(() =>
          (props.parsedEquipmentFilter?.length ?? 0) > 0 &&
          !matchesParsedFilter(props.item.equipment ?? [], props.parsedEquipmentFilter ?? [])
        );
        ```
        `MatrixItem.equipment` is `string[]` per
        `web/src/api/itemGroupMatrix.ts` (confirm exact field name when reading the
        type).
  - [ ] 3.2 Bind a `matrix-row--filtered-out` class to the row's root `<tr>` based
        on `isFilteredOut`. Add scoped CSS:
        ```css
        .matrix-row--filtered-out {
          opacity: 0.35;
          pointer-events: none;
          filter: grayscale(0.3);
        }
        ```
        Match the visual treatment of the existing card-view
        `.item-filtered-out` rule in `ItemsView.vue` / `ItemGroupsView.vue` for
        consistency (grep for the existing class to match its visuals exactly).
  - [ ] 3.3 Add `data-cy="matrix-row-filtered-out"` to the root `<tr>` only when
        `isFilteredOut` is true so E2E specs can target it.

- [ ] Task 4: Saved-filter persistence works exactly as today (AC: #4)
  - [ ] 4.1 No new state is introduced. The existing `equipmentFilter` ref in
        `ItemGroupsView.vue` persists across the card/table view toggle because
        both views are children of the same component. Verify by reading the
        `activeView` toggle logic — both branches render under the same parent.
  - [ ] 4.2 Saved filters (`useSavedFilters`) require no change — they continue
        to populate the combobox `:items` list, which the table view now respects
        through the shared `parsedEquipmentFilter` computed.

- [ ] Task 5: Tests (Vitest + Vue Test Utils)
  - [ ] 5.1 In `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.test.ts`
        (or create alongside, matching the pattern used by other matrix tests),
        add tests:
        - With `parsedEquipmentFilter` = `[]`, the row root does not carry
          `matrix-row--filtered-out`.
        - With a matching parsed group (e.g. `{ exact: [], keywords: ['monitor'] }`)
          and `item.equipment = ['Monitor', 'Keyboard']`, the row does NOT carry
          the class.
        - With a non-matching parsed group (e.g. `{ exact: [], keywords: ['foo'] }`)
          and `item.equipment = ['Monitor']`, the row DOES carry the class.
  - [ ] 5.2 In `web/src/views/ItemGroupsView.test.ts`, add one test that mounts
        `ItemGroupsView` in table mode with a non-matching filter set and asserts
        the matrix renders with at least one `matrix-row--filtered-out` row.
        Re-use the existing matrix view stub if present; otherwise mount the real
        matrix with a small fixture.

- [ ] Task 6: Verification commands
  - [ ] 6.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green.
  - [ ] 6.2 Manual smoke: open `/item-groups?areaId=...` on a desktop viewport,
        click Table view, type a keyword in the equipment filter input at the top,
        confirm non-matching rows dim; clear the input, confirm all rows return.

### Review Findings

- [x] [Review][Patch] Weekly table filter is not covered through the actual `ItemGroupsView` input/view path; current tests only pass `parsedEquipmentFilter` directly to `AreaWeeklyMatrixView`, so they do not prove `ig-equipment-filter` state is preserved and propagated when switching to table view [web/src/views/ItemGroupsView.test.ts:549]
- [x] [Review][Patch] Filtered matrix rows replace the stable `matrix-row-<id>` selector with `matrix-row-filtered-out-<id>`, so the same row cannot be addressed by its normal selector once filtered and this deviates from adding a filtered marker to the existing row [web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue:5]

## Dev Notes

### Where the filter input already lives

The equipment filter at the top of the page IS the same `<v-combobox
data-cy="ig-equipment-filter">` element used by the card view
(`web/src/views/ItemGroupsView.vue:78–93`). The bug is that today the table view
is rendered as a sibling of the card grid but never receives the filter state.
This story wires the missing link; it does NOT introduce a second filter input.

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Filter parser | `parseFilter`, `matchesParsedFilter` | `web/src/composables/useEquipmentFilter.ts` |
| Type for AND-group array | `AndGroup` | `web/src/composables/useEquipmentFilter.ts:14` |
| Current filter ref + computed | `equipmentFilter`, `parsedEquipmentFilter` | `web/src/views/ItemGroupsView.vue:355, 380` |
| Matrix view mount point | `<AreaWeeklyMatrixView>` | `web/src/views/ItemGroupsView.vue:136` |
| Row component | `AreaWeeklyMatrixRow.vue` | `web/src/components/area-weekly-matrix/` |
| Matrix item type with `equipment` | `MatrixItem` | `web/src/api/itemGroupMatrix.ts` |
| Card-view dim CSS recipe | `.item-filtered-out` | `web/src/views/ItemGroupsView.vue` scoped style |

### Anti-patterns to avoid

- Do NOT add a second equipment filter input above the matrix. Reuse the existing
  one in `ItemGroupsView.vue`.
- Do NOT mutate `MatrixItem` shape to add a `_filteredOut` flag — compute it in
  the row component from the prop. Mutating API DTOs leads to stale state across
  refreshes.
- Do NOT short-circuit the matrix fetch based on the filter (i.e. don't pass the
  filter as a query parameter to the matrix API). The filter is a UI affordance;
  filtering happens in the browser so toggling does not cause network round-trips.
- Do NOT forget the null-guard fix from Task 1 — without it this story's wiring
  inherits the same clear-X bug as Story 33.1.

### Coordination with Story 33.1

If 33.1 lands first, copy the `?? ''` guard pattern; if this story lands first,
apply the same pattern to both `ItemsView.vue:1195` AND `ItemGroupsView.vue:380`.
Either way, the watcher snippet from 33.1 Task 1.3 is a 4-line addition and
should land in both files.

### Testing standards

- Vitest + Vue Test Utils. New tests extend existing files; do not introduce new
  test files unless none exist for the component you touch.
- Reuse the `MatrixItem` fixture pattern from existing
  `AreaWeeklyMatrixView.test.ts` (`web/src/components/area-weekly-matrix/`).

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/views/ItemGroupsView.vue:78–105, 136, 355, 380–396]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue]
- [Source: web/src/api/itemGroupMatrix.ts]
- [Source: web/src/composables/useEquipmentFilter.ts]
- [Source: _bmad-output/implementation-artifacts/33-1-equipment-filter-reset-on-item-groups-view.md]
- [Source: private/epic-33.md]
- [Source: .claude/rules/vue.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/components/area-weekly-matrix/ src/views/ItemGroupsView.test.ts` — 68 tests pass
  (3 new matrix-equipment-filter cases).
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all clean.

### Completion Notes List

- Threaded the existing `parsedEquipmentFilter` computed in `ItemGroupsView.vue`
  down through `AreaWeeklyMatrixView` → `AreaWeeklyMatrixRoomSection` →
  `AreaWeeklyMatrixRow` as a new optional `parsedEquipmentFilter?: AndGroup[]`
  prop. Default `[]` so the matrix continues to render correctly when mounted
  from a context that doesn't pass the prop.
- `AreaWeeklyMatrixRow.vue` computes `isFilteredOut` from
  `matchesParsedFilter(item.equipment, props.parsedEquipmentFilter)`; rows that
  don't match get a `matrix-row--filtered-out` CSS class
  (`opacity: 0.35; pointer-events: none; filter: grayscale(0.3)`) and a
  `data-cy="matrix-row-filtered-out-{itemId}"` selector.
- Saved-filter persistence requires no new state — the existing
  `equipmentFilter` ref in `ItemGroupsView.vue` already persists across the
  card/table view toggle. No backend or API changes.
- The null-guard fix from Story 33.1 inoculates this story too (the same
  `equipmentFilter` ref in `ItemGroupsView.vue` is fixed there).
- Added 3 unit tests in `AreaWeeklyMatrixView.test.ts > equipment filter`:
  no rows dimmed for empty filter, partial dim for partial-match, all-rows
  dim for no-match.

### File List

Frontend (modified):

- `web/src/views/ItemGroupsView.vue` — passes `parsedEquipmentFilter` prop
  to `<AreaWeeklyMatrixView>`.
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue` — accepts
  and forwards the prop; imports `AndGroup` type.
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRoomSection.vue` —
  forwards the prop to row.
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue` — accepts
  the prop, computes `isFilteredOut`, applies CSS class + data-cy on the
  `<tr>` root; adds scoped `.matrix-row--filtered-out` styles.
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts` —
  extended `mountMatrix` helper to accept `parsedEquipmentFilter`; added
  `equipment filter` describe block with 3 tests.
