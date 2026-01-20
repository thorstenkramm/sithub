# Story 10.6: Booking Management Views Redesign

## Story

**As a** user,  
**I want** my bookings displayed beautifully,  
**So that** managing my reservations is pleasant.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** Medium

## Acceptance Criteria

**AC1: My Bookings Redesign**
- **Given** I open My Bookings
- **When** I have upcoming bookings
- **Then** I see bookings as attractive cards
- **And** each card shows desk, room, area, and date
- **And** cards are sorted by date (nearest first)
- **And** "Booked by" indicator shows if someone booked for me
- **And** "Guest" indicator shows for guest bookings

**AC2: Cancel Booking Flow**
- **Given** I want to cancel a booking
- **When** I click the cancel button
- **Then** I see a confirmation dialog (ConfirmDialog component)
- **And** the dialog clearly states what I'm cancelling
- **When** I confirm
- **Then** the booking is removed with a success message
- **And** the list updates immediately

**AC3: Booking History Redesign**
- **Given** I open Booking History
- **When** the page loads
- **Then** I see a DateRangePicker for filtering
- **And** I see past bookings in a clean list or table
- **And** each entry shows desk, room, area, and date
- **When** I select a date range
- **Then** the list filters to show only that range

**AC4: Empty States**
- **Given** I have no upcoming bookings
- **When** I view My Bookings
- **Then** I see the EmptyState component with a helpful message
- **And** I see a call-to-action to "Book a Desk"

- **Given** I have no booking history (or filtered results are empty)
- **When** I view Booking History
- **Then** I see the EmptyState component with appropriate message

**AC5: Loading States**
- **Given** bookings are being fetched
- **When** the page is loading
- **Then** I see skeleton loaders matching the content layout

**AC6: Presence and Room Bookings Views**
- **Given** I view Area Presence or Room Bookings
- **When** the data loads
- **Then** I see a clean, readable list of who is booked where
- **And** the design is consistent with other views

## Technical Requirements

### Booking Card Design (My Bookings)
```
+--------------------------------------------------+
| [Date Badge]                        [Cancel] [X] |
| Monday, January 20, 2026                         |
+--------------------------------------------------+
| Desk: Standing Desk A                            |
| Room: Conference Room 1                          |
| Area: Berlin Office                              |
|                                                  |
| [Booked by: Jane Smith]  (if applicable)        |
| [Guest: John Visitor]    (if guest booking)     |
+--------------------------------------------------+
```

### Booking History Table/List
```
+------------------------------------------------------------------+
| Date         | Desk            | Room              | Area        |
+------------------------------------------------------------------+
| Jan 15, 2026 | Standing Desk A | Conference Room 1 | Berlin      |
| Jan 14, 2026 | Hot Desk 3      | Open Space        | Berlin      |
| Jan 10, 2026 | Desk B2         | Quiet Room        | Munich      |
+------------------------------------------------------------------+
```

### Date Range Filter
```
+------------------------------------------+
| From: [Jan 1, 2026]  To: [Jan 31, 2026]  |
| [Apply Filter]  [Clear]                   |
+------------------------------------------+
```

## Tasks

### Task 1: Create BookingCard Component
- [ ] Create `web/src/components/BookingCard.vue`
- [ ] Display date prominently
- [ ] Show desk, room, area info
- [ ] Show "booked by" and "guest" indicators
- [ ] Add cancel button with emit
- [ ] Write unit test

### Task 2: Redesign MyBookingsView
- [ ] Replace list with BookingCard grid/list
- [ ] Add LoadingState with card skeletons
- [ ] Add EmptyState with "Book a Desk" CTA
- [ ] Integrate ConfirmDialog for cancel
- [ ] Handle cancel success/error

### Task 3: Create DateRangeFilter Component
- [ ] Create `web/src/components/DateRangeFilter.vue`
- [ ] Two date pickers: From and To
- [ ] Apply and Clear buttons
- [ ] Emit filter values
- [ ] Write unit test

### Task 4: Redesign BookingHistoryView
- [ ] Add DateRangeFilter component
- [ ] Display history in table or card format
- [ ] Add LoadingState
- [ ] Add EmptyState for no results
- [ ] Handle filter changes

### Task 5: Redesign AreaPresenceView
- [ ] Apply design system styling
- [ ] Use consistent card/list layout
- [ ] Add LoadingState and EmptyState
- [ ] Show user, desk, room info clearly

### Task 6: Redesign RoomBookingsView
- [ ] Apply design system styling
- [ ] Use consistent card/list layout
- [ ] Add LoadingState and EmptyState
- [ ] Show user, desk info clearly
- [ ] Add date picker for viewing different dates

### Task 7: Cancel Booking Integration
- [ ] Use ConfirmDialog for all cancel actions
- [ ] Show clear confirmation message
- [ ] Handle loading state during cancel
- [ ] Show success/error feedback

### Task 8: Polish and Test
- [ ] Test My Bookings with various states
- [ ] Test Booking History with date filtering
- [ ] Test Presence and Room Bookings views
- [ ] Test empty states
- [ ] Test cancel flow
- [ ] Test on mobile

## File Changes

| Action | File Path |
|--------|-----------|
| Create | `web/src/components/BookingCard.vue` |
| Create | `web/src/components/DateRangeFilter.vue` |
| Modify | `web/src/views/MyBookingsView.vue` |
| Modify | `web/src/views/BookingHistoryView.vue` |
| Modify | `web/src/views/AreaPresenceView.vue` |
| Modify | `web/src/views/RoomBookingsView.vue` |
| Create | `web/src/components/__tests__/BookingCard.test.ts` |
| Create | `web/src/components/__tests__/DateRangeFilter.test.ts` |

## Definition of Done

- [ ] MyBookingsView uses BookingCard components
- [ ] BookingHistoryView has DateRangeFilter
- [ ] AreaPresenceView is redesigned with consistent styling
- [ ] RoomBookingsView is redesigned with consistent styling
- [ ] All views have proper loading states
- [ ] All views have proper empty states
- [ ] Cancel booking uses ConfirmDialog
- [ ] All designs match the design system
- [ ] All existing tests still pass
- [ ] Code passes linting

## Notes

- Booking cards should show dates in a human-readable format
- History view could use a table on desktop, cards on mobile
- Consider adding sort options (date, area, room)
- Cancel button should have appropriate styling (warning/error color)

## Dependencies

- Story 10.1: Design System Foundation
- Story 10.2: Reusable Component Library
- Story 10.3: Navigation & Layout Redesign

## Blocked By

- Story 10.1
- Story 10.2
- Story 10.3

## Blocks

- None directly
