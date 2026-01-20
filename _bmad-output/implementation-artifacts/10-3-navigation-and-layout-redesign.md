# Story 10.3: Navigation & Layout Redesign

## Story

**As a** user,  
**I want** clear navigation and context awareness,  
**So that** I always know where I am and can easily move around.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** High (affects all views)

## Acceptance Criteria

**AC1: Improved App Bar**
- **Given** I am on any page
- **When** I look at the app bar
- **Then** I see the SitHub logo on the left
- **And** I see navigation links (Areas, My Bookings, History)
- **And** the current page/section is visually highlighted
- **And** I see my user name and a menu on the right

**AC2: Breadcrumb Navigation**
- **Given** I am viewing a room's desks
- **When** I look at the page header
- **Then** I see breadcrumbs: Home > [Area Name] > [Room Name]
- **And** each breadcrumb is clickable to navigate back

**AC3: User Menu**
- **Given** I click on my user name/avatar
- **When** the menu opens
- **Then** I see my full name and email
- **And** I see a logout option
- **And** (if admin) I see an "Admin" indicator

**AC4: Mobile Navigation**
- **Given** I am on a mobile device (< 768px)
- **When** I view the app bar
- **Then** navigation links are hidden
- **And** I see a hamburger menu icon
- **When** I tap the hamburger icon
- **Then** a drawer opens with all navigation options

**AC5: Active Route Indication**
- **Given** I am on the My Bookings page
- **When** I look at the navigation
- **Then** "My Bookings" is visually highlighted as active
- **And** other links are not highlighted

## Technical Requirements

### App Bar Structure (Desktop)
```
+------------------------------------------------------------------------+
| [Logo]     Areas    My Bookings    History          [User Name ▼]      |
+------------------------------------------------------------------------+
```

### App Bar Structure (Mobile)
```
+----------------------------------------+
| [☰]        [Logo]           [Avatar]   |
+----------------------------------------+
```

### Drawer Menu (Mobile)
```
+---------------------------+
| [User Name]               |
| [user@email.com]          |
| [Admin Badge] (if admin)  |
+---------------------------+
| Areas                     |
| My Bookings               |
| History                   |
+---------------------------+
| Logout                    |
+---------------------------+
```

### Breadcrumb Logic
| Route | Breadcrumbs |
|-------|-------------|
| `/` | Home |
| `/areas/:areaId/rooms` | Home > [Area Name] |
| `/rooms/:roomId/desks` | Home > [Area Name] > [Room Name] |
| `/rooms/:roomId/bookings` | Home > [Area Name] > [Room Name] > Bookings |
| `/areas/:areaId/presence` | Home > [Area Name] > Presence |
| `/my-bookings` | Home > My Bookings |
| `/bookings/history` | Home > Booking History |

## Tasks

### Task 1: Redesign App.vue App Bar
- [ ] Update app bar with logo, nav links, user menu
- [ ] Style nav links with active state (using router-link-active)
- [ ] Add user name display on right side
- [ ] Ensure proper spacing and alignment

### Task 2: Create User Menu Component
- [ ] Create `web/src/components/UserMenu.vue`
- [ ] Display user name, email, admin status
- [ ] Add logout functionality
- [ ] Style according to design system

### Task 3: Implement Mobile Navigation Drawer
- [ ] Add v-navigation-drawer to App.vue
- [ ] Show/hide based on screen size
- [ ] Add hamburger menu icon for mobile
- [ ] Include all nav links in drawer
- [ ] Include user info and logout in drawer

### Task 4: Create Breadcrumb Integration
- [ ] Use PageHeader component with breadcrumbs prop
- [ ] Create composable `useBreadcrumbs()` to generate crumbs from route
- [ ] Fetch area/room names for dynamic crumbs
- [ ] Integrate into each view

### Task 5: Update AreasView with New Layout
- [ ] Add PageHeader with breadcrumbs
- [ ] Verify navigation works correctly
- [ ] Test mobile view

### Task 6: Update RoomsView with New Layout
- [ ] Add PageHeader with breadcrumbs (Home > Area)
- [ ] Fetch area name for breadcrumb
- [ ] Verify navigation works correctly

### Task 7: Update DesksView with New Layout
- [ ] Add PageHeader with breadcrumbs (Home > Area > Room)
- [ ] Fetch area and room names for breadcrumbs
- [ ] Verify navigation works correctly

### Task 8: Update Remaining Views
- [ ] MyBookingsView - Add PageHeader
- [ ] BookingHistoryView - Add PageHeader
- [ ] AreaPresenceView - Add PageHeader with breadcrumbs
- [ ] RoomBookingsView - Add PageHeader with breadcrumbs

### Task 9: Responsive Testing
- [ ] Test on desktop (> 1024px)
- [ ] Test on tablet (768px - 1024px)
- [ ] Test on mobile (< 768px)
- [ ] Verify drawer works on mobile
- [ ] Verify no horizontal scroll

## File Changes

| Action | File Path |
|--------|-----------|
| Modify | `web/src/App.vue` |
| Create | `web/src/components/UserMenu.vue` |
| Create | `web/src/composables/useBreadcrumbs.ts` |
| Modify | `web/src/views/AreasView.vue` |
| Modify | `web/src/views/RoomsView.vue` |
| Modify | `web/src/views/DesksView.vue` |
| Modify | `web/src/views/MyBookingsView.vue` |
| Modify | `web/src/views/BookingHistoryView.vue` |
| Modify | `web/src/views/AreaPresenceView.vue` |
| Modify | `web/src/views/RoomBookingsView.vue` |
| Create | `web/src/components/__tests__/UserMenu.test.ts` |

## Definition of Done

- [ ] App bar displays logo, nav links, and user menu
- [ ] Current route is highlighted in navigation
- [ ] Breadcrumbs are displayed on all relevant pages
- [ ] Mobile drawer navigation works correctly
- [ ] All navigation links work correctly
- [ ] User can logout from user menu
- [ ] Responsive design works at all breakpoints
- [ ] All existing tests still pass
- [ ] Code passes linting

## Notes

- Use Vue Router's `router-link-active` class for active states
- Breadcrumb data may require additional API calls or route meta
- Consider caching area/room names to avoid repeated fetches
- Drawer should close when a link is clicked on mobile

## Dependencies

- Story 10.1: Design System Foundation
- Story 10.2: Reusable Component Library (PageHeader)

## Blocked By

- Story 10.1
- Story 10.2

## Blocks

- Story 10.7: Mobile Responsiveness (refines mobile nav)
