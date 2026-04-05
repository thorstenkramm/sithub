# Story 21.4: Swap Booking Date and Item in My Bookings

Status: done

## Story

As a user,
I want to see the booking date as the primary line on my booking cards,
so that I can scan my upcoming bookings chronologically at a glance.

## Acceptance Criteria

1. **Given** I am viewing "My Bookings",
   **when** a booking card renders,
   **then** the booking date appears on the first line (card title) and the booked item
   name appears on the second line (card subtitle).

2. **Given** a booking card displays,
   **when** I look at the subtitle,
   **then** I see the item name, item group name, and area name separated by bullets
   (e.g. "Desk 1 &bull; Room 1 &bull; Office").

3. **Given** a booking has a status chip (guest, booked-for-me, on-behalf),
   **when** the card renders,
   **then** the status chip appears next to the date on the first line.

## Tasks / Subtasks

- [x] Task 1: Reorder BookingCard layout (AC: 1, 2, 3)
  - [x] 1.1 Move formatted date into `v-card-title` (first line)
  - [x] 1.2 Move item name into `v-card-subtitle` with item group and area (second line)
  - [x] 1.3 Change avatar icon from `$desk` to `$calendar` to match the new date-first layout
  - [x] 1.4 Keep status chips on the title line next to the date
  - [x] 1.5 Remove the old date row from `v-card-text` (it is now in the title)
- [x] Task 2: Verify existing tests still pass (AC: 1)
  - [x] 2.1 Run Vitest unit tests
  - [x] 2.2 Run ESLint

## Dev Notes

### Scope

This is a frontend-only change in `BookingCard.vue`. The backend API response structure
is unchanged — only the presentation order in the card template changes.

### Before / After

**Before:** Item name (title) > Item group + Area (subtitle) > Date (card text)
**After:** Date (title) > Item name + Item group + Area (subtitle)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Reordered BookingCard.vue template: date is now v-card-title, item/group/area is v-card-subtitle
- Avatar icon changed from $desk to $calendar
- Status chips remain on the title line next to the formatted date
- Removed the standalone date row from v-card-text since date is now in the title
- All 255 Vitest tests pass, ESLint clean

### File List

- `web/src/components/BookingCard.vue` — Swapped date to title line, item to subtitle line
