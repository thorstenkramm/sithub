# Story 10.7: Mobile Responsiveness

## Story

**As a** mobile user,  
**I want** the app to work well on my phone,  
**So that** I can book desks on the go.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** Medium

## Acceptance Criteria

**AC1: Responsive Layout**
- **Given** I access the app on a mobile device (< 768px)
- **When** I view any page
- **Then** the layout adapts to single-column
- **And** cards stack vertically
- **And** there is no horizontal scrolling

**AC2: Touch-Friendly Interactions**
- **Given** I am using touch input
- **When** I interact with buttons, links, and inputs
- **Then** touch targets are at least 44px in size
- **And** there is adequate spacing between interactive elements
- **And** hover states have touch equivalents

**AC3: Mobile Navigation**
- **Given** I am on a mobile device
- **When** I view the app bar
- **Then** I see a hamburger menu icon
- **When** I tap the hamburger icon
- **Then** a drawer slides in with all navigation options
- **When** I tap a navigation link
- **Then** the drawer closes and I navigate to the page

**AC4: Mobile Forms**
- **Given** I am booking a desk on mobile
- **When** I fill out the booking form
- **Then** inputs are appropriately sized for touch
- **And** date pickers work well on mobile
- **And** the keyboard doesn't obscure important content

**AC5: Mobile Tables/Lists**
- **Given** I view booking history on mobile
- **When** the table would be too wide
- **Then** I see a mobile-friendly card layout instead
- **And** all information is still accessible

**AC6: Breakpoint Consistency**
- **Given** I resize the browser window
- **When** I cross breakpoint thresholds
- **Then** the layout transitions smoothly
- **And** there are no jarring layout shifts

## Technical Requirements

### Breakpoints
```css
/* Mobile first approach */
/* Default: < 600px (mobile) */
/* sm: >= 600px (large mobile / small tablet) */
/* md: >= 960px (tablet) */
/* lg: >= 1280px (desktop) */
/* xl: >= 1920px (large desktop) */
```

### Responsive Grid
```vue
<!-- Cards should stack on mobile, multi-column on larger screens -->
<v-row>
  <v-col cols="12" sm="6" md="4" lg="3" v-for="item in items">
    <ItemCard :item="item" />
  </v-col>
</v-row>
```

### Touch Target Sizes
- Minimum button height: 44px
- Minimum clickable area: 44x44px
- Spacing between targets: 8px minimum

### Mobile-Specific Considerations
| Component | Desktop | Mobile |
|-----------|---------|--------|
| Navigation | Horizontal links | Drawer |
| Card grid | 2-4 columns | 1 column |
| Tables | Full table | Card list |
| Date picker | Calendar popup | Full-screen calendar |
| Dialogs | Centered modal | Bottom sheet or full-screen |

## Tasks

### Task 1: Audit Current Responsive Issues
- [ ] Test all views at mobile width (375px, 414px)
- [ ] Document all layout issues
- [ ] Document all touch target issues
- [ ] Document all usability issues

### Task 2: Fix App Shell Responsiveness
- [ ] Ensure mobile nav drawer works correctly
- [ ] Fix any app bar overflow issues
- [ ] Test user menu on mobile

### Task 3: Fix AreasView Responsiveness
- [ ] Cards should be single column on mobile
- [ ] Touch targets should be adequate
- [ ] Test empty state on mobile

### Task 4: Fix RoomsView Responsiveness
- [ ] Cards should be single column on mobile
- [ ] Availability bar should fit
- [ ] Action buttons should be touch-friendly

### Task 5: Fix DesksView Responsiveness
- [ ] Cards should be single column on mobile
- [ ] Date picker should work well on mobile
- [ ] Booking dialog should be mobile-friendly (possibly full-screen)
- [ ] Equipment list should truncate appropriately

### Task 6: Fix MyBookingsView Responsiveness
- [ ] Booking cards should be full-width on mobile
- [ ] Cancel button should be easily tappable
- [ ] Empty state should look good on mobile

### Task 7: Fix BookingHistoryView Responsiveness
- [ ] Date range filter should stack on mobile
- [ ] History list should use cards on mobile (not table)
- [ ] Filter controls should be touch-friendly

### Task 8: Fix Presence/RoomBookings Views
- [ ] Lists should be readable on mobile
- [ ] Date picker should work on mobile
- [ ] All touch targets should be adequate

### Task 9: Fix All Dialogs
- [ ] Booking dialog: consider bottom sheet on mobile
- [ ] Confirm dialog: ensure buttons are tappable
- [ ] Success/error dialogs: should be readable

### Task 10: Cross-Browser Mobile Testing
- [ ] Test on iOS Safari
- [ ] Test on Android Chrome
- [ ] Test on tablet sizes
- [ ] Fix any platform-specific issues

### Task 11: Add Responsive Utilities
- [ ] Create CSS utilities for hiding/showing at breakpoints
- [ ] Create composable for detecting mobile (`useIsMobile`)
- [ ] Document responsive patterns for future development

## File Changes

| Action | File Path |
|--------|-----------|
| Modify | `web/src/App.vue` |
| Modify | `web/src/views/AreasView.vue` |
| Modify | `web/src/views/RoomsView.vue` |
| Modify | `web/src/views/DesksView.vue` |
| Modify | `web/src/views/MyBookingsView.vue` |
| Modify | `web/src/views/BookingHistoryView.vue` |
| Modify | `web/src/views/AreaPresenceView.vue` |
| Modify | `web/src/views/RoomBookingsView.vue` |
| Modify | `web/src/components/BookingDialog.vue` |
| Modify | `web/src/components/DateRangeFilter.vue` |
| Create | `web/src/composables/useBreakpoints.ts` |
| Modify | `web/src/styles/global.css` |

## Definition of Done

- [ ] All views work correctly at 375px width (iPhone SE)
- [ ] All views work correctly at 414px width (iPhone Plus)
- [ ] All views work correctly at 768px width (iPad)
- [ ] No horizontal scrolling on any view
- [ ] All touch targets are at least 44px
- [ ] Navigation drawer works on mobile
- [ ] All dialogs are usable on mobile
- [ ] All forms are usable on mobile
- [ ] Tested on iOS Safari and Android Chrome
- [ ] No visual bugs at breakpoint transitions
- [ ] All existing tests still pass
- [ ] Code passes linting

## Notes

- Use Vuetify's built-in responsive classes where possible
- Test with actual mobile devices, not just browser dev tools
- Consider using `@media (hover: hover)` for hover-only styles
- Mobile performance is important - avoid heavy animations
- Consider reduced motion preferences for accessibility

## Testing Devices/Sizes

| Device | Width | Priority |
|--------|-------|----------|
| iPhone SE | 375px | High |
| iPhone 12/13/14 | 390px | High |
| iPhone Plus/Max | 414px | High |
| iPad Mini | 768px | Medium |
| iPad | 810px | Medium |
| iPad Pro | 1024px | Low |

## Dependencies

- Story 10.1: Design System Foundation
- Story 10.2: Reusable Component Library
- Story 10.3: Navigation & Layout Redesign
- Story 10.4: Space Discovery Views Redesign
- Story 10.5: Booking Flow Redesign
- Story 10.6: Booking Management Views Redesign

## Blocked By

- All previous UI stories (10.1 - 10.6)

## Blocks

- None (final polish story)
