# Story 22.13: Week Mode Tile Alignment with Day Mode

Status: done

## Story

As a user,
I want week mode item tiles to show the same information as day mode tiles,
so that I can make informed booking decisions without switching modes.

## Acceptance Criteria

1. **Given** I am in week mode
   **When** an item tile renders
   **Then** it shows an availability chip (e.g., "Verfügbar 7/7" or "3/5 frei")
   respecting the weekend toggle for the total count

2. **Given** an item is 100% free for the week
   **When** the avatar icon renders
   **Then** the avatar background color is green (success)

3. **Given** an item is 0% free (fully booked all days)
   **When** the avatar icon renders
   **Then** the avatar background color is red (error)

4. **Given** an item is partially booked
   **When** the avatar icon renders
   **Then** the avatar background color is blue (primary)

5. **Given** I am in week mode
   **When** an item tile renders in folded state
   **Then** the equipment list is visible (same as day mode)

6. **Given** I have selected days for booking in week mode
   **When** the book button renders
   **Then** it shows "Buchen (N Tage)" / "Book (N Days)" — short label

7. **Given** I have selected days for booking
   **When** I scroll through the item list
   **Then** the book button remains visible at the bottom of the viewport
   (sticky/fixed position)

## Tasks / Subtasks

- [ ] Task 1: Add availability chip to week mode tiles (AC: 1)
  - [ ] 1.1 Calculate availability ratio from weekData: count free days
    vs total days (5 or 7 based on weekend toggle)
  - [ ] 1.2 Render a StatusChip or custom chip showing "Verfügbar N/M"
    using i18n key
- [ ] Task 2: Color-code item avatar by availability (AC: 2, 3, 4)
  - [ ] 2.1 Compute avatar color: `success` (100% free), `error` (0% free),
    `primary` (partial)
  - [ ] 2.2 Apply to the `v-avatar` color prop on week mode tiles
- [ ] Task 3: Show equipment on folded week tiles (AC: 5)
  - [ ] 3.1 In week mode folded view: add the equipment chip list below
    the item title (same template as day mode equipment section)
- [ ] Task 4: Shorten book button text (AC: 6)
  - [ ] 4.1 Add i18n key: `"bookDays": "Buchen ({count} Tage)"`
    (en: "Book ({count} Days)")
  - [ ] 4.2 Replace the existing confirm button label
- [ ] Task 5: Sticky book button (AC: 7)
  - [ ] 5.1 Make the week booking submit button `position: sticky` at
    the bottom of the viewport. Use `bottom: 0` with a background
    matching the page surface color
  - [ ] 5.2 Only show when selections exist (current behavior)
- [ ] Task 6: Run tests and lint

## Dev Notes

### Availability Calculation

```typescript
const freeCount = selectedWeekDates.value.filter(
  date => getWeekDayStatus(item.id, date) === 'free'
).length;
const totalDays = selectedWeekDates.value.length; // 5 or 7
const ratio = freeCount / totalDays;
// color: ratio === 1 ? 'success' : ratio === 0 ? 'error' : 'primary'
```

### Sticky Button Pattern

```css
.week-book-footer {
  position: sticky;
  bottom: 0;
  z-index: 2;
  background: rgb(var(--v-theme-surface));
  padding: 12px 16px;
  border-top: 1px solid rgb(var(--v-theme-outline));
}
```

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
