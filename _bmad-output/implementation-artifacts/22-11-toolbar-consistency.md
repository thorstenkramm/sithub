# Story 22.11: Toolbar Consistency Across Views

Status: done

## Story

As a mobile user,
I want the booking toolbar (date/week selector, floor plan button, equipment filter)
to look and behave consistently across all views,
so that I always know where to find controls.

## Acceptance Criteria

1. **Given** I am on ItemGroupsView
   **When** the toolbar renders
   **Then** the equipment filter is inside the same card as the week selector
   (not in a separate card), and includes the info (i) button

2. **Given** I am on ItemsView in day mode
   **When** the date picker renders
   **Then** it uses full width (matching the week selector), has no calendar icon,
   and uses `density="compact"` matching the week selector height

3. **Given** a floor plan exists for the current area/item group
   **When** any view renders
   **Then** the RAUMPLAN button is always positioned next to the week/day
   selector on the same row

4. **Given** I switch between day and week mode
   **When** the selector outline renders
   **Then** both have the same visual height (both use `density="compact"`)

## Tasks / Subtasks

- [ ] Task 1: Move equipment filter into week selector card (ItemGroupsView) (AC: 1)
  - [ ] 1.1 In `ItemGroupsView.vue`: remove the separate `v-card` wrapping the
    equipment filter. Move the filter combobox + save button + info button
    into the existing week selector card
  - [ ] 1.2 Add the info (i) button (currently missing in ItemGroupsView)
- [ ] Task 2: Fix day mode date picker (AC: 2, 4)
  - [ ] 2.1 In `ItemsView.vue`: remove the calendar icon from DatePickerField
    (remove `prepend-inner-icon` or equivalent)
  - [ ] 2.2 Remove `max-width: 280px` constraint — use `max-width: 320px`
    matching the week selector
  - [ ] 2.3 Add `density="compact"` to the DatePickerField to match the week
    selector height
- [ ] Task 3: Consistent floor plan button position (AC: 3)
  - [ ] 3.1 In all three views, ensure the RAUMPLAN button is in the same
    flex row as the week/day selector with `align-end` alignment
- [ ] Task 4: Run tests and lint

## Dev Notes

### Files to Change

- `web/src/views/ItemGroupsView.vue` — merge filter card, add info button
- `web/src/views/ItemsView.vue` — date picker sizing/icon, floor plan position
- `web/src/components/DatePickerField.vue` — may need density prop support

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
