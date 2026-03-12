# Story 17.1: Equipment Filter with Advanced Search

Status: done

## Story

As a user,
I want to filter items by equipment keywords,
So that I can quickly find a workspace with the tools I need.

## Acceptance Criteria

1. **Given** I am on the booking page (day or week mode)
   **When** I see the booking options card
   **Then** a text input labeled "Filter equipment" appears below the colleague option
   **And** an info icon appears next to the input

2. **Given** I click the info icon
   **When** the explanation popup appears
   **Then** it describes the search syntax:
   show only items having the filter keyword(s) in any of the equipment items;
   multiple keywords are combined with OR;
   use plus sign to combine with AND;
   use single or double quotation marks for exact matching;
   filters are case-insensitive;
   example: "27 inch display" + webcam

3. **Given** I type "webcam" into the filter input
   **When** the filter is applied
   **Then** items that have "webcam" in any of their equipment are shown normally
   **And** items without "webcam" in their equipment are blurred with an "equipment not
   available" overlay hint
   **And** blurred items are not removed from the list

4. **Given** I type `"27 inch display" + webcam` into the filter input
   **When** the filter is applied
   **Then** only items having both "27 inch display" (exact) AND "webcam" in their equipment
   are shown normally
   **And** all other items are blurred

5. **Given** I clear the filter input
   **When** the filter is removed
   **Then** all items are shown normally without blur

6. **Given** I am in week booking mode
   **When** I type a filter
   **Then** the same filtering logic applies to the week mode item tiles

## Tasks / Subtasks

- [x] Add filter input to booking options card (AC: 1)
  - [x] In `ItemsView.vue`: add a `v-text-field` labeled "Filter equipment" after the
    colleague/booking type section
  - [x] Add `data-cy="equipment-filter-input"` for testing
  - [x] Add an info icon (`$info` alias) next to the input
  - [x] Add `data-cy="equipment-filter-info"` for testing
- [x] Add info popup (AC: 2)
  - [x] On info icon click, show a `v-dialog` with the search syntax explanation
  - [x] Use the exact text from acceptance criteria
  - [x] Add `data-cy="equipment-filter-help"` for testing
- [x] Create equipment filter parser (AC: 3, 4, 5)
  - [x] Create `web/src/composables/useEquipmentFilter.ts` composable
  - [x] Parse filter string into search terms:
    - Split by `+` for AND groups
    - Within each group, quoted strings are exact matches
    - Unquoted terms are OR-combined keywords
    - All matching is case-insensitive
  - [x] Export a `matchesFilter(equipment: string[], filterText: string): boolean` function
  - [x] Export a `parseFilter(filterText: string): AndGroup[]` function for testing
  - [x] Return `true` if filter is empty (no filtering)
- [x] Apply filter to day mode items (AC: 3, 4, 5)
  - [x] Add reactive `equipmentFilter` ref
  - [x] `isItemFilteredOut()` computes whether each item matches the filter
  - [x] Non-matching items: add CSS class for blur effect and overlay
  - [x] Items are NOT removed from the DOM, only blurred
- [x] Apply filter to week mode items (AC: 6)
  - [x] `getWeekItemEquipment()` extracts equipment from first available day's data
  - [x] Same blur/overlay behavior as day mode
- [x] Add blur CSS and overlay (AC: 3)
  - [x] `.item-filtered-out` class with `filter: blur(3px)` and reduced opacity
  - [x] `.item-filtered-overlay` overlay div with "equipment not available" text
  - [x] `.item-filter-wrapper` relative container for positioning
- [x] Add unit tests
  - [x] `useEquipmentFilter.test.ts`: 19 tests for parser and matchesFilter
  - [x] `ItemsView.test.ts`: 4 new tests for filter rendering and blur behavior
  - [x] Test parser: OR keywords, AND with `+`, exact match with quotes
  - [x] Test `matchesFilter()` with various equipment arrays and filter strings
  - [x] Test empty filter returns all items as matching
  - [x] Test case-insensitivity
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

All filtering is client-side. Equipment data is already available in item attributes.
No backend changes required.

### Equipment Data Structure

In `web/src/api/items.ts`, the `ItemAttributes` interface (lines 4-12) has:

```typescript
equipment: string[];
```

Equipment is an array of strings like `["24 inch display", "webcam", "USB-C dock"]`.

### Filter Parsing Algorithm

Example: `"27 inch display" + webcam`

1. Split by `+` (with trimming): `['"27 inch display"', 'webcam']`
2. Each segment is an AND condition
3. For each segment:
   - If quoted (single or double): exact match against each equipment string
   - If unquoted: split by spaces into keywords, OR-combined
4. Item matches if ALL AND conditions are satisfied

Simple keyword example: `webcam monitor`

1. No `+` split: `['webcam monitor']`
2. Split into keywords: `['webcam', 'monitor']`
3. OR-combined: item matches if any equipment contains "webcam" OR "monitor"

### Blur Effect

Use CSS filter for blur and an absolute-positioned overlay:

```css
.item-filtered-out {
  filter: blur(3px);
  opacity: 0.5;
  pointer-events: none;
}
.item-filtered-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}
```

### Week Mode Equipment Access

In week mode, equipment data is available in `weekData[date]` per-day item attributes.
Since equipment is an item property (not date-dependent), use the first available day's data
to determine equipment for filtering:

```typescript
const getWeekItemEquipment = (itemId: string): string[] => {
  for (const dayItems of Object.values(weekData.value)) {
    const item = dayItems.find(i => i.id === itemId);
    if (item?.attributes.equipment?.length) return item.attributes.equipment;
  }
  return [];
};
```

### References

- Epic 17 Story 17.1: `_bmad-output/planning-artifacts/epics.md` (Epic 17 Stories section)
- FR58: `_bmad-output/planning-artifacts/prd.md`
- ItemsView: `web/src/views/ItemsView.vue`
- ItemAttributes: `web/src/api/items.ts` lines 4-12

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Pure frontend story: no backend changes required
- Created `useEquipmentFilter.ts` composable with `parseFilter()` and `matchesFilter()` exports
- Filter parser supports OR (space-separated), AND (`+`), and exact match (quotes)
- Exact matches require full-string equality; mixed quoted/unquoted terms are OR within a group
- Day mode items wrapped in `.item-filter-wrapper` div with conditional blur and overlay
- Week mode uses `getWeekItemEquipment()` to extract equipment from first available day's data
- Filter help dialog explains syntax with examples and icon button has an accessible label
- Added 2 more unit tests for edge cases in the parser/matcher
- Senior review fixes applied: exact AC text restored in help dialog, filter parsing reused across card renders, and week-mode/unit interaction coverage strengthened
- Added Cypress coverage for day-mode and week-mode equipment filtering and verified the full Cypress suite passes
- Noted unrelated working tree changes during review (not part of this story): internal/startup/server.go, internal/users/handler.go, internal/users/store.go, web/src/api/users.ts, web/src/components/BookingCard.vue, web/src/composables/useWeekendPreference.ts, web/src/composables/useWeekendPreference.test.ts, web/src/plugins/vuetify.ts, web/src/styles/global.css, web/src/views/AreaPresenceView.vue, floor_plans/

### File List

- `web/src/composables/useEquipmentFilter.ts` — Filter parser and matcher composable
- `web/src/composables/useEquipmentFilter.test.ts` — 20 unit tests including parsed-filter reuse
- `web/src/views/ItemsView.vue` — Filter input, help dialog, blur/overlay for day and week
- `web/src/views/ItemsView.test.ts` — 5 equipment filter tests including dialog click and week-mode blur
- `web/cypress/e2e/items.cy.ts` — Added day-mode equipment filter E2E coverage
- `web/cypress/e2e/week-booking.cy.ts` — Added week-mode equipment filter E2E coverage

## Senior Developer Review (AI)

### Reviewer

Thorsten

### Date

2026-03-12

### Findings Resolved

- Restored the filter help dialog wording to match the acceptance criteria text exactly.
- Reworked the view to reuse parsed filter groups instead of reparsing the filter for every card render.
- Strengthened unit tests so the help dialog test now verifies click-driven visibility instead of relying on an always-open stub.
- Added explicit week-mode filter assertions in unit tests.
- Added Cypress coverage for day-mode and week-mode filter behavior.

### Remaining Gap

- None.

## Change Log

- 2026-03-12: Senior review fixes applied for exact help text, parsed-filter reuse, stronger unit coverage, added Cypress coverage, and verified full test suite pass.
