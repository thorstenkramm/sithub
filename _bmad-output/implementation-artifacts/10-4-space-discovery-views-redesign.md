# Story 10.4: Space Discovery Views Redesign

## Story

**As a** user,  
**I want** visually appealing space discovery,  
**So that** browsing areas, rooms, and desks is enjoyable and efficient.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** Medium

## Acceptance Criteria

**AC1: Areas View Redesign**
- **Given** I am on the areas page
- **When** areas exist
- **Then** I see areas displayed as visually appealing cards
- **And** each card shows area name and description (if available)
- **And** cards have hover effects and are clickable
- **When** no areas exist
- **Then** I see the EmptyState component with illustration

**AC2: Rooms View Redesign**
- **Given** I am viewing rooms in an area
- **When** rooms exist
- **Then** I see room cards with name and desk count
- **And** each card shows availability summary (e.g., "3/5 desks available")
- **And** cards have visual indicators for availability status
- **When** no rooms exist
- **Then** I see the EmptyState component

**AC3: Desks View Redesign**
- **Given** I am viewing desks in a room
- **When** desks exist
- **Then** I see desks as cards in a grid layout
- **And** each card shows desk name, equipment, and status
- **And** available desks are visually distinct from booked desks
- **And** my bookings are highlighted differently
- **And** warnings are displayed prominently
- **When** no desks exist
- **Then** I see the EmptyState component

**AC4: Date Selection**
- **Given** I am on the desks view
- **When** I want to change the date
- **Then** I see the DatePicker component
- **And** selecting a new date refreshes the desk availability

**AC5: Quick Actions**
- **Given** I am on the areas view
- **When** I look at an area card
- **Then** I see quick action buttons: "View Rooms", "Today's Presence"
- **And** clicking them navigates to the appropriate view

**AC6: Loading States**
- **Given** data is being fetched
- **When** the page is loading
- **Then** I see skeleton loaders matching the card layout
- **And** there is no layout shift when content loads

## Technical Requirements

### Area Card Design
```
+----------------------------------+
| [Icon]                           |
| Area Name                        |
| Description text here...         |
|                                  |
| [View Rooms]  [Today's Presence] |
+----------------------------------+
```

### Room Card Design
```
+----------------------------------+
| Room Name                        |
| [===----] 3/5 available          |
|                                  |
| [View Desks]  [View Bookings]    |
+----------------------------------+
```

### Desk Card Design
```
+----------------------------------+
| Desk Name            [Available] |
| Equipment: Monitor, Keyboard     |
| ⚠️ Near window - may be cold    |
|                                  |
| [Book This Desk]                 |
+----------------------------------+
```

### Color Coding for Desk Status
| Status | Card Style |
|--------|------------|
| Available | Default card, success chip, primary "Book" button |
| Booked by others | Muted card, warning chip, disabled button |
| Booked by me | Primary border/highlight, "My Booking" chip, "Cancel" button |
| Unavailable | Grayed out, error chip |

## Tasks

### Task 1: Create AreaCard Component
- [ ] Create `web/src/components/AreaCard.vue`
- [ ] Display area name, description
- [ ] Add action buttons (View Rooms, Today's Presence)
- [ ] Add hover effects
- [ ] Write unit test

### Task 2: Create RoomCard Component
- [ ] Create `web/src/components/RoomCard.vue`
- [ ] Display room name, availability bar
- [ ] Add action buttons (View Desks, View Bookings)
- [ ] Write unit test

### Task 3: Create DeskCard Component
- [ ] Create `web/src/components/DeskCard.vue`
- [ ] Display desk name, status chip, equipment, warning
- [ ] Add book/cancel button based on status
- [ ] Style based on availability status
- [ ] Write unit test

### Task 4: Redesign AreasView
- [ ] Replace list with AreaCard grid
- [ ] Add LoadingState with cards skeleton
- [ ] Add EmptyState for no areas
- [ ] Ensure PageHeader with breadcrumbs

### Task 5: Redesign RoomsView
- [ ] Replace list with RoomCard grid
- [ ] Fetch and display availability summary per room
- [ ] Add LoadingState with cards skeleton
- [ ] Add EmptyState for no rooms
- [ ] Ensure PageHeader with breadcrumbs

### Task 6: Redesign DesksView
- [ ] Replace list with DeskCard grid
- [ ] Integrate DatePicker component
- [ ] Style cards based on booking status
- [ ] Add LoadingState with cards skeleton
- [ ] Add EmptyState for no desks
- [ ] Ensure PageHeader with breadcrumbs

### Task 7: Add Availability Summary to Rooms
- [ ] Create API endpoint or compute from existing data
- [ ] Display availability bar (visual indicator)
- [ ] Update on date change

### Task 8: Polish and Test
- [ ] Test all three views with real data
- [ ] Test empty states
- [ ] Test loading states
- [ ] Verify responsive layout (grid collapses on mobile)

## File Changes

| Action | File Path |
|--------|-----------|
| Create | `web/src/components/AreaCard.vue` |
| Create | `web/src/components/RoomCard.vue` |
| Create | `web/src/components/DeskCard.vue` |
| Modify | `web/src/views/AreasView.vue` |
| Modify | `web/src/views/RoomsView.vue` |
| Modify | `web/src/views/DesksView.vue` |
| Create | `web/src/components/__tests__/AreaCard.test.ts` |
| Create | `web/src/components/__tests__/RoomCard.test.ts` |
| Create | `web/src/components/__tests__/DeskCard.test.ts` |

## Definition of Done

- [ ] AreasView displays area cards in a responsive grid
- [ ] RoomsView displays room cards with availability summary
- [ ] DesksView displays desk cards with clear status indication
- [ ] All views have proper loading states (skeletons)
- [ ] All views have proper empty states (illustrations)
- [ ] DatePicker is integrated in DesksView
- [ ] Cards have hover effects and are clickable
- [ ] Responsive layout works on all screen sizes
- [ ] All existing tests still pass
- [ ] Code passes linting

## Notes

- Use CSS Grid for card layouts (2-3 columns on desktop, 1 on mobile)
- Consider adding subtle animations when cards appear
- Equipment list should truncate if too long, with tooltip for full list
- Availability bar can be a simple progress-bar style indicator

## Dependencies

- Story 10.1: Design System Foundation
- Story 10.2: Reusable Component Library
- Story 10.3: Navigation & Layout Redesign

## Blocked By

- Story 10.1
- Story 10.2
- Story 10.3

## Blocks

- Story 10.5: Booking Flow Redesign (uses DeskCard)
