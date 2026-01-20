# Story 10.5: Booking Flow Redesign

## Story

**As a** user,  
**I want** an intuitive and delightful booking experience,  
**So that** reserving a desk feels effortless.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** Medium

## Acceptance Criteria

**AC1: Single Day Booking**
- **Given** I click "Book" on an available desk
- **When** the booking dialog opens
- **Then** I see the desk name and selected date clearly
- **And** I can confirm the booking with one click
- **When** I confirm
- **Then** I see a success message with booking details
- **And** the desk status updates immediately

**AC2: Multi-Day Booking**
- **Given** I want to book multiple days
- **When** I click "Book Multiple Days" or toggle multi-day mode
- **Then** I see a calendar where I can select multiple dates
- **And** I can click individual dates to toggle selection
- **And** I see the count of selected dates
- **When** I confirm
- **Then** I see results showing successful bookings and any conflicts

**AC3: Book for Colleague**
- **Given** I want to book for someone else
- **When** I toggle "Book for colleague"
- **Then** I see input fields for colleague name and email
- **And** the booking is clearly labeled as "on behalf of"
- **When** I confirm
- **Then** the booking shows who it's for

**AC4: Guest Booking**
- **Given** I want to book for a guest
- **When** I toggle "Book for guest"
- **Then** I see input fields for guest name and optional email
- **And** the booking is clearly labeled as a guest booking
- **When** I confirm
- **Then** the booking shows the guest indicator

**AC5: Booking Confirmation**
- **Given** I complete a booking
- **When** the booking succeeds
- **Then** I see a success dialog/notification
- **And** I see options: "View My Bookings" or "Book Another Desk"
- **And** the dialog auto-closes after a few seconds (or manual close)

**AC6: Error Handling**
- **Given** a booking fails (e.g., desk no longer available)
- **When** the error occurs
- **Then** I see a clear error message explaining what happened
- **And** I am prompted to try another desk or date
- **And** the error state is visually distinct (not just text)

## Technical Requirements

### Booking Dialog Structure
```
+------------------------------------------+
| Book Desk                            [X] |
+------------------------------------------+
| Desk: Standing Desk A                    |
| Room: Conference Room 1                  |
| Date: January 20, 2026                   |
|                                          |
| [ ] Book for multiple days               |
| [ ] Book for colleague                   |
| [ ] Book for guest                       |
|                                          |
| [Cancel]                    [Confirm]    |
+------------------------------------------+
```

### Multi-Day Calendar View
```
+------------------------------------------+
| Select Dates                         [X] |
+------------------------------------------+
|        January 2026                      |
|  Mo Tu We Th Fr Sa Su                   |
|           1  2  3  4  5                  |
|   6  7  8  9 10 11 12                   |
|  13 14 15 [16] 17 [18] 19               |
|  [20] 21 22 23 24 25 26                 |
|  27 28 29 30 31                         |
|                                          |
| Selected: 3 days                         |
| [Cancel]                    [Confirm]    |
+------------------------------------------+
```

### Success State
```
+------------------------------------------+
| ✓ Booking Confirmed!                     |
+------------------------------------------+
| You've booked:                           |
| Standing Desk A                          |
| Conference Room 1, Berlin Office         |
| January 20, 2026                         |
|                                          |
| [View My Bookings]  [Book Another Desk]  |
+------------------------------------------+
```

## Tasks

### Task 1: Create BookingDialog Component
- [ ] Create `web/src/components/BookingDialog.vue`
- [ ] Implement single-day booking view
- [ ] Show desk, room, area, date info
- [ ] Add confirm/cancel buttons
- [ ] Write unit test

### Task 2: Add Multi-Day Selection
- [ ] Add "Book for multiple days" toggle
- [ ] Integrate calendar component for date selection
- [ ] Allow selecting/deselecting individual dates
- [ ] Show selected date count
- [ ] Handle multi-day API response (success + conflicts)

### Task 3: Add Colleague Booking
- [ ] Add "Book for colleague" toggle
- [ ] Show name and email input fields
- [ ] Validate inputs before submission
- [ ] Pass colleague info to API

### Task 4: Add Guest Booking
- [ ] Add "Book for guest" toggle
- [ ] Show guest name (required) and email (optional) fields
- [ ] Validate inputs before submission
- [ ] Pass guest info to API

### Task 5: Create BookingSuccess Component
- [ ] Create `web/src/components/BookingSuccess.vue`
- [ ] Show booking confirmation details
- [ ] Add action buttons (View Bookings, Book Another)
- [ ] Optional: auto-close with countdown

### Task 6: Create BookingError Component
- [ ] Create `web/src/components/BookingError.vue`
- [ ] Show clear error message
- [ ] Suggest next actions
- [ ] Allow retry or close

### Task 7: Integrate into DesksView
- [ ] Replace current booking logic with BookingDialog
- [ ] Handle dialog open/close state
- [ ] Connect to booking API
- [ ] Update desk status after booking

### Task 8: Handle Multi-Day Results
- [ ] Display successful bookings
- [ ] Display conflicts separately
- [ ] Allow partial success (some days booked, some conflicted)

### Task 9: Polish and Test
- [ ] Test single-day booking flow
- [ ] Test multi-day booking flow
- [ ] Test colleague booking flow
- [ ] Test guest booking flow
- [ ] Test error scenarios
- [ ] Test on mobile

## File Changes

| Action | File Path |
|--------|-----------|
| Create | `web/src/components/BookingDialog.vue` |
| Create | `web/src/components/BookingSuccess.vue` |
| Create | `web/src/components/BookingError.vue` |
| Create | `web/src/components/MultiDayCalendar.vue` |
| Modify | `web/src/views/DesksView.vue` |
| Create | `web/src/components/__tests__/BookingDialog.test.ts` |
| Create | `web/src/components/__tests__/BookingSuccess.test.ts` |

## Definition of Done

- [ ] Single-day booking works via new dialog
- [ ] Multi-day booking works with calendar selection
- [ ] Colleague booking works with name/email inputs
- [ ] Guest booking works with guest indicator
- [ ] Success confirmation shows booking details
- [ ] Error states are handled gracefully
- [ ] All booking types work correctly
- [ ] UI is polished and consistent with design system
- [ ] All existing tests still pass
- [ ] Code passes linting

## Notes

- The backend API already supports all booking types (single, multi-day, colleague, guest)
- Multi-day calendar can use Vuetify's v-date-picker with multiple selection
- Consider keyboard navigation for calendar
- Success/error feedback should be obvious and not easily missed

## Dependencies

- Story 10.1: Design System Foundation
- Story 10.2: Reusable Component Library
- Story 10.4: Space Discovery Views Redesign (DeskCard integration)

## Blocked By

- Story 10.1
- Story 10.2

## Blocks

- None directly, but improves overall UX
