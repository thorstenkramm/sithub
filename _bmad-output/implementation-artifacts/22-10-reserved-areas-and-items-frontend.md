# Story 22.10: Reserved Areas and Items — Frontend

Status: ready-for-dev

## Story

As a user,
I want to see which items I cannot book because they are reserved for others,
so that I do not waste time trying to book restricted items.

## Acceptance Criteria

1. **Given** I am viewing items where some are reserved for other users
   **When** the item list renders in day mode
   **Then** reserved items are disabled and visually blurred (like the equipment
   filter blur pattern)

2. **Given** I am viewing items in week mode
   **When** reserved items render
   **Then** their checkboxes are disabled and the tile is visually dimmed

3. **Given** a floor plan shows items reserved for other users
   **When** the floor plan renders
   **Then** reserved desk positions are grayed out with a lock icon

4. **Given** I tap on a reserved/disabled item
   **When** the interaction is processed
   **Then** no booking action occurs and a tooltip or message explains
   the item is reserved

## Tasks / Subtasks

- [ ] Task 1: Handle `reserved` flag in items API response (AC: 1, 2)
  - [ ] 1.1 In `web/src/api/items.ts`: update `ItemAttributes` interface to
    include `reserved?: boolean`
  - [ ] 1.2 In `web/src/views/ItemsView.vue` day mode: check `entry.attributes.reserved`
    and apply the existing blur pattern (`.item-filtered-out` class and
    `.item-filtered-overlay` from the equipment filter, lines ~1754-1772)
  - [ ] 1.3 Show overlay message `$t('items.reserved')` — add i18n key:
    "Reserviert" (de) / "Reserved" (en) / equivalents
- [ ] Task 2: Disable reserved items in week mode (AC: 2)
  - [ ] 2.1 In `web/src/views/ItemsView.vue` week mode: check `reserved` flag
    on the item. Disable all day checkboxes and dim the tile card
  - [ ] 2.2 Reuse the `.item-filtered-out` CSS class for visual consistency
- [ ] Task 3: Gray out reserved items on floor plan (AC: 3)
  - [ ] 3.1 In `web/src/components/InteractiveFloorPlan.vue`: the desk positions
    already have booking status. Add a `reserved` state check — if the item
    is reserved for others, render the position with reduced opacity and a
    lock icon overlay (`mdi-lock`)
  - [ ] 3.2 The floor plan component receives item data. Ensure the `reserved`
    flag is passed through from the items API
- [ ] Task 4: Add reserved tooltip/message (AC: 4)
  - [ ] 4.1 Wrap reserved items in a `v-tooltip` showing
    `$t('items.reservedTooltip')`: "Dieses Objekt ist reserviert." (de) /
    "This item is reserved." (en)
  - [ ] 4.2 Prevent click events on reserved items (`@click.prevent` or
    `:disabled="true"` on the book button)
- [ ] Task 5: Add i18n keys (AC: 1, 4)
  - [ ] 5.1 Add to all locale files under `items`:
    `"reserved": "Reserviert"`, `"reservedTooltip": "Dieses Objekt ist
    reserviert. Sie haben keinen Zugriff."`
- [ ] Task 6: Write tests (AC: 1, 2, 3, 4)
  - [ ] 6.1 Test that reserved items render with blur/overlay in day mode
  - [ ] 6.2 Test that reserved items have disabled checkboxes in week mode
  - [ ] 6.3 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

## Dev Notes

### Existing Blur Pattern to Reuse

The equipment filter already implements a blur+overlay pattern:

```vue
<div v-if="isItemFilteredOut(...)" class="item-filtered-overlay">
  <span>{{ $t('items.equipmentNotAvailable') }}</span>
</div>
<v-card :class="['item-card', { 'item-filtered-out': isItemFilteredOut(...) }]">
```

CSS (ItemsView.vue lines 1754-1772):

```css
.item-filtered-out { filter: blur(3px); opacity: 0.5; pointer-events: none; }
.item-filtered-overlay { position: absolute; inset: 0; display: flex;
  align-items: center; justify-content: center; z-index: 1; }
```

Reuse this exact pattern for reserved items. Add a new check function
`isItemReserved(entry)` alongside the existing `isItemFilteredOut()`.

### Dependencies

This story requires Story 22.9 (backend reserved_for enforcement and
`reserved` flag in items API) to be complete.

### Floor Plan Positioning Data

The floor plan component loads positions from `/api/v1/floor-plan-positions`
and booking data from the items/bookings APIs. The `reserved` flag needs to
flow through to the position rendering logic.

### Files to Change

| File | Change |
| --- | --- |
| `web/src/api/items.ts` | Add `reserved` to ItemAttributes |
| `web/src/views/ItemsView.vue` | Day + week mode reserved state |
| `web/src/components/InteractiveFloorPlan.vue` | Floor plan reserved overlay |
| `web/src/locales/*.json` | Add reserved i18n keys |

### References

- [Source: private/epic-22.md — "disabled and blurred"]
- [Source: web/src/views/ItemsView.vue:1754-1772 — existing blur pattern]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
