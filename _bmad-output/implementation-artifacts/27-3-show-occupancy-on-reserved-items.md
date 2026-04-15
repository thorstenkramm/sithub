# Story 27.3: Show Occupancy on Reserved Items in Regular Booking View

Status: done

## Story

As a user,
I want to see who is sitting where on reserved items in the list view,
so that I know room occupancy even though I cannot book there myself.

## Acceptance Criteria

1. **Given** I view items in a reserved area (e.g., Finance & People Area)
   **When** I look at the item list
   **Then** I can see free/busy status and booker names — the opaque veil is removed

2. **Given** I see a reserved item that is free
   **When** I look at it
   **Then** it shows a badge with a lock icon and "reserved" text

3. **Given** I try to book a reserved item
   **When** I interact with it
   **Then** the booking action is blocked (no book button or it is disabled)

## Tasks / Subtasks

- [x] Task 1: Remove opaque overlay from reserved items in day mode (AC: #1)
  - [x] 1.1 In `ItemsView.vue`, locate the overlay at line ~176-184: the `v-if` condition
        includes `entry.attributes.reserved` which renders a full opaque overlay hiding
        all item content
  - [x] 1.2 Remove `|| entry.attributes.reserved` from the overlay's `v-if` condition
        (line ~177) so reserved items are no longer fully hidden
  - [x] 1.3 Remove `|| entry.attributes.reserved` from the card's `item-filtered-out`
        class condition (line ~190) so the card is not dimmed
- [x] Task 2: Add lock badge on free reserved items in day mode (AC: #2)
  - [x] 2.1 Add a `v-chip` or badge inside the item card that shows when
        `entry.attributes.reserved === true` — display a lock icon (`mdi-lock`) with
        "Reserved" text
  - [x] 2.2 Style the badge to be visible but non-intrusive (e.g., small chip at top-right
        or below the item name)
  - [x] 2.3 Add `data-cy="item-reserved-badge"` for testability
- [x] Task 3: Block booking action on reserved items in day mode (AC: #3)
  - [x] 3.1 The "Book" button for day-mode items: add `:disabled` or `v-if` condition
        to hide/disable when `entry.attributes.reserved === true`
  - [x] 3.2 Verify clicking a reserved item does not open a booking dialog
- [x] Task 4: Remove opaque overlay from reserved items in week mode (AC: #1)
  - [x] 4.1 Locate the week-mode overlay (line ~381 area) — same pattern as day mode
  - [x] 4.2 Remove the reserved condition from the overlay `v-if` and the card class
  - [x] 4.3 Ensure week-mode checkboxes for reserved items remain disabled (already
        handled by `isWeekItemReserved()` at line ~509)
- [x] Task 5: Add lock badge on reserved items in week mode (AC: #2)
  - [x] 5.1 Add a similar lock badge/chip for week-mode reserved items
  - [x] 5.2 Ensure the badge is visible in both expanded and collapsed tile states
- [x] Task 6: Validate (AC: #1-#3)
  - [x] 6.1 Run `npm run lint` and fix findings
  - [x] 6.2 Run `npm run type-check` and fix findings
  - [x] 6.3 Run `npm run build` and verify no build errors
  - [x] 6.4 Run `npx vitest run` and verify no regressions
  - [x] 6.5 Run `npm run test:e2e -- --browser electron` and verify no regressions

## Dev Notes

### Architecture & Patterns

- **Primary file**: `web/src/views/ItemsView.vue`
- **No backend changes** — the `reserved` field and booking prevention already work
- The current implementation uses an opaque overlay (`item-filtered-overlay`) that
  completely hides the item card content. The fix replaces this with a visible card
  plus a "reserved" badge.

### Key Code Locations

| Element | Location | data-cy |
|---------|----------|---------|
| Day-mode overlay | `ItemsView.vue:176-184` | `item-reserved` |
| Day-mode card class | `ItemsView.vue:190` | `item-entry` |
| Week-mode overlay | `ItemsView.vue:~381` | `item-reserved` |
| Book button (day) | `ItemsView.vue` — look for `book-item-btn` | `book-item-btn` |
| `isWeekItemReserved()` | `ItemsView.vue:1187` | — |
| `getOverlayLabel()` | `ItemsView.vue:1164` | — |

### Current Overlay HTML Structure

```html
<div v-if="isItemFilteredOut(...) || entry.attributes.reserved"
     class="item-filtered-overlay">
  <span>Reserved</span>
</div>
<v-card :class="['item-card', { 'item-filtered-out': ... || entry.attributes.reserved }]">
  <!-- item content hidden behind overlay -->
</v-card>
```

After fix:
```html
<div v-if="isItemFilteredOut(...)"
     class="item-filtered-overlay">
  <span>Equipment not available</span>
</div>
<v-card :class="['item-card', ...]">
  <!-- item content visible -->
  <v-chip v-if="entry.attributes.reserved" ...>Reserved</v-chip>
  <!-- book button disabled for reserved items -->
</v-card>
```

### Anti-Patterns to Avoid

- Do NOT allow booking of reserved items — booking must remain blocked
- Do NOT change the `reserved` field from the backend
- Do NOT remove the equipment filter overlay — only decouple reserved from filtered-out

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List

### Change Log
