# Story 22.3: Mobile Text Truncation Fixes

Status: done

## Story

As a mobile user,
I want to see full item names and dates without truncation,
so that I can distinguish between similar items like "Tisch 2, Fenster, rechts"
and "Tisch 2, Fenster, links".

## Acceptance Criteria

1. **Given** an item name longer than the card width
   **When** it renders in day mode item cards or week mode tile headers
   **Then** the text stays on a single line and truncates from the middle
   (e.g., "Tisch 1, am...rechts") so both prefix and suffix remain visible.
   Full name is shown on tap via tooltip.

2. **REMOVED** — replaced by middle-truncation in AC 1

3. **Given** a booking card shows item name, group, and area in the subtitle
   **When** the subtitle is too long for the card width
   **Then** the text wraps to a second line

4. **Given** the booking history page on a 390px screen
   **When** the date filter fields render
   **Then** "Von" and "Bis" fields stack vertically with readable date values

## Tasks / Subtasks

- [ ] Task 1: Revert white-space:normal on item card titles (AC: 1)
  - [ ] 1.1 In `web/src/views/ItemsView.vue` scoped CSS: remove the
    `.item-card :deep(.v-card-title) { white-space: normal }` rule added
    in the previous implementation — item names must stay single-line
- [ ] Task 2: Implement middle-truncation for item names (AC: 1)
  - [ ] 2.1 Create a utility function `middleTruncate(text, maxLen)` that
    keeps the first and last portions visible with "..." in the middle
    (e.g., "Tisch 1, am...rechts" for maxLen=20)
  - [ ] 2.2 Apply in day mode card title and week mode tile header
  - [ ] 2.3 Wrap the truncated name in a `v-tooltip` showing the full name
    on tap/hover
- [ ] Task 3: Fix BookingCard subtitle (AC: 3)
  - [ ] 3.1 In `web/src/components/BookingCard.vue` lines 34-36: the
    `v-card-subtitle` truncates by default. Override with
    `white-space: normal` to allow wrapping
- [ ] Task 4: Fix history date filter layout on mobile (AC: 4)
  - [ ] 4.1 In `web/src/views/BookingHistoryView.vue` lines 12-35: the date
    pickers use `d-flex flex-wrap` with `max-width: 200px`. On narrow screens
    200px is too small for the label + date value. Remove `max-width` or
    increase it, and ensure the flex container wraps to a column on mobile
    using `flex-column` below a breakpoint
  - [ ] 4.2 The labels "Von" and "Bis" (from `$t('history.fromDate')` and
    `$t('history.toDate')`) should be fully visible — verify the `v-text-field`
    label is not clipped
- [ ] Task 5: Run tests and lint (AC: 1, 2, 3, 4)
  - [ ] 5.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

## Dev Notes

### Root Cause

Vuetify's `v-card-title` and `v-card-subtitle` default to `white-space: nowrap`
and `text-overflow: ellipsis`. This works on desktop but fails on 390px mobile
screens where item names like "Tisch 2, Fenster, rechts" are 25+ characters.

### Recommended Approach

Use scoped CSS to override `white-space: normal` on the specific card
components. Avoid `!important` if possible — use component-scoped `:deep()`
selector. Example:

```css
:deep(.v-card-title) {
  white-space: normal;
  line-height: 1.4;
}
```

For the history date filter, switch to `d-flex flex-column` on mobile using
Vuetify's responsive display classes (`d-flex flex-column flex-sm-row`).

### Files to Change

| File | Lines | Change |
| --- | --- | --- |
| `web/src/views/ItemsView.vue` | ~202, ~384 | Card title white-space |
| `web/src/components/BookingCard.vue` | 34-36 | Subtitle white-space |
| `web/src/views/BookingHistoryView.vue` | 12-35 | Date filter stacking |

### References

- [Source: private/ux-observations.md — "Item names cut off everywhere"]
- [Source: private/ux-observations.md — "History date filter labels truncated"]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

### Review Findings

- [x] [Review][Patch] Item names are middle-truncated only when they exceed the rendered card width; widths are re-measured on load and resize, and regression coverage was added for fitting vs overflowing names [web/src/views/ItemsView.vue:206]
