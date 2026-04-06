# Story 22.4: Week Mode Mobile Readability

Status: done

## Story

As a mobile user,
I want to see who booked a desk in week mode without overlapping text,
so that I can quickly scan the week grid.

## Acceptance Criteria

1. **Given** a desk is booked in week mode on a 390px screen
   **When** the booker's name renders under the day column (folded/compact view)
   **Then** it shows initials (e.g., "AE" for Alex Employee) that fit the column
   without overflow

2. **Given** I tap on a booked day cell showing initials
   **When** the tooltip activates
   **Then** I see the full user name

3. **Given** the week tile is expanded
   **When** the booker's name renders in the expanded row
   **Then** the full name is shown (expanded view has sufficient width)

## Tasks / Subtasks

- [ ] Task 1: Show initials in compact week mode columns (AC: 1)
  - [ ] 1.1 In `web/src/views/ItemsView.vue` lines ~503-512: the compact view
    displays `getWeekDayBooker(item.id, date)` which returns the full name.
    Create a helper `getWeekDayBookerInitials(itemId, date)` that extracts
    initials from the full name (same logic as `userInitials` in App.vue
    lines 348-357)
  - [ ] 1.2 Use initials in the compact span (`.week-day-status-truncated`),
    keep full name in the existing `v-tooltip`
  - [ ] 1.3 Update CSS `.week-day-status` (lines 1710-1718): the `max-width: 60px`
    and `text-overflow: ellipsis` can be simplified since initials are always
    2-3 characters. Keep the constraint but remove `text-overflow: ellipsis`
- [ ] Task 2: Ensure expanded view shows full name (AC: 3)
  - [ ] 2.1 In `web/src/views/ItemsView.vue` lines ~578-581: the expanded view
    uses `getWeekDayBooker()` directly — keep as-is (full name)
  - [ ] 2.2 Verify the expanded grid `grid-template-columns: 40px 180px 1fr`
    (line 1735) provides enough space
- [ ] Task 3: Run tests and lint (AC: 1, 2, 3)
  - [ ] 3.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

## Dev Notes

### Initials Extraction Pattern

Reuse the proven pattern from `App.vue` lines 348-357:

```typescript
function getInitials(name: string): string {
  const parts = name.split(' ');
  if (parts.length >= 2) {
    return (parts[0]!.charAt(0) + parts[parts.length - 1]!.charAt(0)).toUpperCase();
  }
  return name.substring(0, 2).toUpperCase();
}
```

### Current CSS Constraints

- `.week-day-slot` min-width: 44px (line 1698)
- `.week-day-status` max-width: 60px, font-size: 0.7rem (lines 1710-1718)
- `.week-days-compact` grid gap: 4px (line 1691)

Initials ("AE", "TK") at 0.7rem will fit comfortably in 44-60px.

### Files to Change

- `web/src/views/ItemsView.vue` — compact week booker display, helper function

### References

- [Source: private/ux-observations.md — "User names overlap/truncate"]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
