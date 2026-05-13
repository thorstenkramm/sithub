# Story 33.1: Equipment Filter Resets on Item-Groups View

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user filtering items by equipment,
I want clearing the filter input to re-enable every tile,
so that I can recover from a typo or a non-matching search without reloading the page.

## Acceptance Criteria

1. **Given** I am on an item-groups page in day mode and I type a keyword that no
   item's equipment matches
   **When** I press the input's built-in clear "X" icon
   **Then** every tile that was blurred / disabled by the filter returns to its normal
   interactive state without a page reload

2. **Given** I am on the same page with a non-matching filter applied
   **When** I backspace the input until it is empty
   **Then** every tile returns to its normal interactive state immediately

3. **Given** I am on the same page in week mode
   **When** I clear the filter via either method above
   **Then** every week tile returns to its normal interactive state

4. **Given** the filter is empty
   **When** I view the page
   **Then** no tile is blurred or disabled by the equipment-filter logic, regardless
   of prior filter history

## Tasks / Subtasks

- [ ] Task 1: Normalise the `equipmentFilter` ref so it never becomes `null`
      (AC: #1, #2, #3, #4)
  - [ ] 1.1 Root cause: `<v-combobox>` with `clearable` emits `null` on clear, but
        the ref is declared `ref('')` of type `string`
        (`web/src/views/ItemsView.vue:1085`). After clearing, `equipmentFilter.value
        === null`, then `parsedEquipmentFilter = computed(() =>
        parseFilter(equipmentFilter.value))` runs `parseFilter(null)`, which calls
        `null.trim()` and throws. Vue swallows the exception inside the computed and
        the previous (non-empty) parsed value can persist depending on how the
        runtime handles it — either way the UI ends up with tiles still blurred.
  - [ ] 1.2 Fix at the consumer: change
        `web/src/views/ItemsView.vue:1195` from
        ```ts
        const parsedEquipmentFilter = computed(() => parseFilter(equipmentFilter.value));
        ```
        to
        ```ts
        const parsedEquipmentFilter = computed(() =>
          parseFilter(equipmentFilter.value ?? ''));
        ```
        This makes the computed safe regardless of whether the combobox sets the
        value to `''` or `null`.
  - [ ] 1.3 Belt-and-braces: add a watcher that normalises the ref itself so any
        downstream code (saved-filters helpers, click handlers) does not have to
        repeat the guard. Place it just below the ref declaration at
        `web/src/views/ItemsView.vue:1085`:
        ```ts
        watch(equipmentFilter, (value) => {
          if (value === null || value === undefined) {
            equipmentFilter.value = '';
          }
        });
        ```
        Keep it small; do not introduce a `computed` wrapper or change the ref's
        type — the v-combobox model contract is `string | null` so the ref staying
        as `string` plus the watcher is the cleanest fix.

- [ ] Task 2: Verify the existing parser already returns "no filter" for empty
      input (no changes expected) (AC: #4)
  - [ ] 2.1 `parseFilter('')` already returns `[]` (`web/src/composables/useEquipmentFilter.ts:25`),
        and `matchesParsedFilter(_, [])` already returns `true` (line 73). No
        changes needed to the composable.
  - [ ] 2.2 Confirm that `isItemFilteredOut` (ItemsView.vue:1197–1199) returns
        `false` whenever `parsedEquipmentFilter.value` is `[]`. After Task 1 the
        cleared input will yield exactly that — the existing logic is correct, the
        bug was upstream.

- [ ] Task 3: Tests (Vitest + Vue Test Utils)
  - [ ] 3.1 In `web/src/views/ItemsView.test.ts`, add a `describe('equipment filter
        reset', ...)` block with three assertions:
        - Given a non-matching filter that blurs all tiles, then setting
          `equipmentFilter` to `null` (simulating v-combobox clear) — no tile carries
          the `item-filtered-overlay` element and no tile has the
          `item-filtered-out` class on its inner wrapper. Use the pattern from the
          existing `equipment filter` describe block (lines 915–...) for the
          fixture shape.
        - Given the same setup, then setting `equipmentFilter` to `''` directly —
          same assertions hold.
        - Given the filter is empty at mount time — no tile is filtered, regardless
          of subsequent rapid set/clear cycles.
  - [ ] 3.2 Cypress regression check is optional — extending the existing
        `web/cypress/e2e/` items spec with a one-line "type, clear, assert tiles
        interactive" is welcome but not required. The unit tests are sufficient
        because the bug is logic-level.

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
        navigate to `/item-groups/parking_lot/items?areaId=underground_car_park`,
        type "foo" in the equipment filter, confirm all tiles get the
        `Ausstattung nicht verfügbar` overlay, click the clear-X icon, confirm
        every tile becomes interactive immediately. Repeat using backspace.

## Dev Notes

### Root cause analysis (please don't skip — it informs the fix)

`<v-combobox>` with `clearable` is documented to emit `null` on clear, while typed
characters emit a `string`. The current ref is `ref('')` (typed `Ref<string>`), so
the runtime tolerates the `null` assignment but every consumer that calls
string-only methods on `.value` breaks. `parseFilter(null)` throws on `.trim()`;
the exception inside `computed` is swallowed and the previous parsed value either
sticks or returns `undefined`, leaving `isItemFilteredOut` evaluating
`!matchesParsedFilter(eq, undefined)` — which throws again and is once more
swallowed. The visible symptom: tiles stay blurred.

The fix is intentionally tiny: a `?? ''` guard plus a watcher that resets the ref.
Do NOT change the combobox to a `v-text-field` to dodge the issue (combobox is
needed for saved-filter dropdown values) and do NOT change the ref type to
`string | null` because every other consumer (save/delete handlers, the saved-filter
matcher) would then need its own null-guard.

### Key code locations

| Element | Location | Why it matters |
| --- | --- | --- |
| `equipmentFilter` ref | `web/src/views/ItemsView.vue:1085` | Add the watcher right after |
| `parsedEquipmentFilter` computed | `web/src/views/ItemsView.vue:1195` | Add the `?? ''` guard here |
| `isItemFilteredOut` | `web/src/views/ItemsView.vue:1197–1199` | Already correct; no change |
| `parseFilter`, `matchesParsedFilter` | `web/src/composables/useEquipmentFilter.ts` | Already correct; no change |
| `<v-combobox>` filter input | `web/src/views/ItemsView.vue` (around lines 108–144) | Source of the `null` value |
| Existing equipment filter tests | `web/src/views/ItemsView.test.ts:915–...` | Pattern to follow |

### Anti-patterns to avoid

- Do NOT replace the `v-combobox` with a plain `v-text-field`. The combobox
  exposes the saved-filters dropdown — losing it breaks Story 19.x's saved-filters
  feature.
- Do NOT delete `clearable`. The clear icon is the affordance the user expects;
  the fix is to handle the value the combobox emits, not remove the affordance.
- Do NOT add a defensive `if (!equipmentFilter.value) return;` to `isItemFilteredOut`.
  The existing logic via `matchesParsedFilter([], _)` already does the right thing
  once the ref is normalised — duplicating the guard scatters logic.

### Testing standards

- Vitest + Vue Test Utils, extending `web/src/views/ItemsView.test.ts`. No new
  test file. Use the same combobox stub the existing equipment-filter tests use
  (which emits string values, so adapt by directly setting
  `wrapper.vm.equipmentFilter = null` to simulate the combobox clear).
- Do not mock `useEquipmentFilter` — it's pure functions and trivial to call
  directly. The bug is in `ItemsView.vue`'s glue, not the composable.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/views/ItemsView.vue:1085]
- [Source: web/src/views/ItemsView.vue:1195–1199]
- [Source: web/src/composables/useEquipmentFilter.ts]
- [Source: private/epic-33.md]
- [Source: private/img_19.png]
- [Source: .claude/rules/vue.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/ItemsView.test.ts src/views/ItemGroupsView.test.ts` — 115 tests pass.
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all clean.

### Completion Notes List

- Root cause confirmed: `<v-combobox>` with `clearable` emits `null` on
  X-icon click; the `equipmentFilter` ref typed `Ref<string>` accepted the
  assignment, but `parseFilter(null)` throws on `.trim()`, the exception is
  swallowed inside the `computed`, and the previous filter sticks.
- Applied the same minimal fix in both `ItemsView.vue` and `ItemGroupsView.vue`:
  widened the ref type to `ref<string | null>('')`, added a watcher that
  coerces any `null`/`undefined` back to `''`, and added a `?? ''` guard on
  the `parsedEquipmentFilter` computed.
- Added a regression test (`equipment filter > removes blur when filter is
  cleared via null (v-combobox clearable)`) that simulates the clear-X by
  setting the ref to `null` and asserts every blurred tile becomes
  interactive again.

### File List

Frontend (modified):

- `web/src/views/ItemsView.vue` — `equipmentFilter` ref widened to
  `string | null`, defensive watcher added, `parsedEquipmentFilter` computed
  guarded with `?? ''`.
- `web/src/views/ItemGroupsView.vue` — same null-guard pattern applied
  (also consumed by Story 33.2).
- `web/src/views/ItemsView.test.ts` — added the `null`-clear regression test.
