# Story 19.5: Equipment Filter on Item Groups View

Status: done

## Story

As a user,
I want to filter item groups by equipment on the area view,
So that I can quickly find rooms or areas that have the equipment I need.

## Acceptance Criteria

1. **Given** I am on the item-groups view (e.g. `/areas/{areaId}/item-groups`)
   **When** I enter an equipment filter keyword
   **Then** item groups whose items do not match the filter are blurred and disabled
   **And** item groups with at least one matching item are shown normally

2. **Given** I clear the filter
   **When** the filter is removed
   **Then** all item groups are shown normally without blur

3. **Given** I use the advanced filter syntax (AND with `+`, exact match with quotes)
   **When** the filter is applied
   **Then** the same parsing rules from the existing equipment filter apply

## Tasks / Subtasks

- [x] Add equipment filter to ItemGroupsView (AC: 1, 2, 3)
  - [x] In `ItemGroupsView.vue`: added equipment filter input with the same styling
    as the ItemsView filter
  - [x] Fetches items per group to access equipment data
  - [x] Aggregates equipment across all items in each group
  - [x] Applies existing `useEquipmentFilter` parsing and matching logic
- [x] Apply blur/disable to non-matching groups (AC: 1)
  - [x] Non-matching item groups are blurred with the same CSS as ItemsView
  - [x] Blurred groups are disabled (not clickable)
- [x] Verify existing filter syntax works (AC: 3)
  - [x] Reuses `matchesFilter` from `useEquipmentFilter.ts` composable
  - [x] AND (`+`), OR (space), and exact match (quotes) all work
- [x] Verify E2E tests still pass

## Dev Notes

### Equipment Aggregation

Since item groups don't directly have equipment, the view fetches all items for each group
and aggregates their equipment arrays. An item group matches the filter if any of its items
has equipment matching the filter criteria.

### Reuse of Existing Filter Logic

The `useEquipmentFilter` composable from Story 17.1 provides `matchesFilter()` and
`parseFilter()`. This story reuses those functions, ensuring consistent filter behavior
across both views.

### References

- Epic 19 Story 19.5: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR68: `_bmad-output/planning-artifacts/prd.md`
- `web/src/views/ItemGroupsView.vue`
- `web/src/composables/useEquipmentFilter.ts` (from Story 17.1)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added equipment filter to ItemGroupsView with same UI pattern as ItemsView
- Fetches items per group to aggregate equipment for filtering
- Reuses `matchesFilter` from existing `useEquipmentFilter` composable
- Blur/disable applied to non-matching item groups
- All existing tests continue to pass

### File List

- `web/src/views/ItemGroupsView.vue` — Equipment filter input, item fetching, aggregation,
  blur/disable logic

## Change Log

- 2026-03-21: Story implemented and verified.
