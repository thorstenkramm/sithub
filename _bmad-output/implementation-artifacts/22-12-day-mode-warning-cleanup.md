# Story 22.12: Day Mode Warning Cleanup

Status: done

## Story

As a mobile user,
I want item warnings to not consume vertical space in the booking view,
so that I can see more items without scrolling.

## Acceptance Criteria

1. **Given** an item has a warning configured
   **When** the day mode tile renders in folded state
   **Then** the warning text block is NOT shown inline — only the warning
   icon (orange !) is visible in the #append slot

2. **Given** I tap the warning icon on a day mode tile
   **When** the tooltip/popup activates
   **Then** the full warning text is displayed

3. **Given** an item tile is expanded in day mode
   **When** the expanded content renders
   **Then** the warning text is shown (expanded view has room)

## Tasks / Subtasks

- [ ] Task 1: Remove inline warning from folded day mode tiles (AC: 1, 2)
  - [ ] 1.1 In `ItemsView.vue`: find the warning `v-alert` / warning block
    that renders inside the item card body in day mode (the orange block
    with warning text). Remove it from the folded state — it should only
    show when the tile is expanded
  - [ ] 1.2 Ensure the warning icon in the #append slot remains visible
    with its tooltip for the folded state
- [ ] Task 2: Keep warning in expanded view (AC: 3)
  - [ ] 2.1 Verify the expanded day tile still shows the warning text
- [ ] Task 3: Run tests and lint

## Dev Notes

### Current Behavior

Day mode folded tiles show both:
- Warning icon (orange !) in the #append slot with tooltip
- Warning text block (orange background) in the card body

This is redundant on mobile. Remove the body warning block from folded state.

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
