---
stepsCompleted: [step-01-validate-prerequisites, step-02-design-epics, step-03-create-stories, step-04-final-validation]
inputDocuments:
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/prd.md
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/architecture.md
---

# sithub - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for sithub, decomposing the requirements from the PRD and
Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: Users can authenticate via Entra ID SSO. Acceptance: unauthenticated users are redirected to Entra ID and, after
successful login, return to SitHub with their name displayed.
FR2: The system can determine user roles (regular vs admin) based on Entra ID group membership. Acceptance: admins see
admin-only cancellation controls; regular users do not.
FR3: Users can access the application only if they are permitted by the configured Entra ID settings. Acceptance:
unauthorized users receive an access-denied screen and cannot view any booking data.
FR4: Users can view a list of available areas. Acceptance: after login, the UI lists all configured areas.
FR5: Users can view a list of rooms within a selected area. Acceptance: selecting an area shows only its rooms.
FR6: Users can view a list of desks within a selected room. Acceptance: selecting a room lists its desks.
FR7: Users can view desk equipment details for each desk. Acceptance: each desk entry shows its equipment list.
FR8: Users can view current booking status for desks. Acceptance: desk entries show available/occupied status for the
selected date.
FR9: Users can book a desk for a single day. Acceptance: selecting a desk and date creates a booking that appears in
“My Bookings.”
FR10: The system prevents double-booking of the same desk for the same day. Acceptance: the second attempt is rejected
and no duplicate booking is created.
FR11: Users receive a message when a selected desk becomes unavailable during booking. Acceptance: the message states
the desk is no longer available for that date and prompts the user to choose another desk.
FR12: Users can view their upcoming bookings (“My Bookings”). Acceptance: the list includes desk, room, area, and date
for each future booking.
FR13: Users can cancel their own bookings. Acceptance: cancelling removes the booking from all relevant lists and
frees the desk.
FR14: Admin users can cancel any booking. Acceptance: admins can cancel another user’s booking and the affected user
sees the cancellation reflected in their list.
FR15: Users can view a room-level booking overview. Acceptance: for a selected room and date, the overview lists all
booked desks and associated users.
FR16: Users can view “Today’s presence” for an area (who is in the office today). Acceptance: the view lists all users
with bookings in that area for today.
FR17: Operators can configure server settings via a configuration file. Acceptance: invalid settings prevent startup
with a clear error; valid settings allow startup.
FR18: Operators can configure areas, rooms, desks, and equipment via a configuration file. Acceptance: after restart,
the UI reflects the updated space definitions.
FR19: The system can load and apply configuration changes on startup. Acceptance: configuration changes take effect
after restart without manual data migration steps.
FR20: Users can book on behalf of another user. Acceptance: the booking appears in both users’ booking lists and either
can cancel.
FR21: Users can book desks for guests outside the organization. Acceptance: a guest booking stores guest name and
contact and is visible as a guest booking in overviews.
FR22: Users can book multi-day or recurring reservations. Acceptance: the system creates individual daily bookings and
reports conflicts per day.
FR23: Users can view booking history. Acceptance: users can see past bookings with date range filtering.
FR24: Users can receive notifications related to bookings. Acceptance: booking creation/cancellation triggers a
notification via the configured channel within 5 minutes.
FR25: Admins can manage rooms/desks via an admin UI. Acceptance: admins can add/edit/remove rooms/desks and changes
appear in discovery lists after save.
FR26: Users can book desks using a graphical floor-map view. Acceptance: a desk selected on the map can be booked for a
chosen date.
FR27: Admins can access advanced reporting and analytics. Acceptance: admins can view usage summaries by area/room and
date range.

### NonFunctional Requirements

NFR1: For expected usage (<=50 concurrent users), list navigation actions complete within 2 seconds at p95; booking and
cancellation complete within 3 seconds at p95.
NFR2: The system can be restarted without data loss; bookings remain intact after restart and conflicts do not create
partial records.
NFR3: All booking data requires authenticated access via Entra ID; unauthenticated requests are denied.
NFR4: Data at rest is stored without application-layer encryption; in-transit encryption is managed outside the
application.
NFR5: Single-node deployment is sufficient; no clustering or horizontal scaling is required for MVP usage levels.
NFR6: Meets WCAG A: all interactive elements have accessible names, keyboard focus is visible, and form inputs are
labeled.

### Additional Requirements

- Go 1.25 with Echo, SQLite (WAL), and JSON:API responses using `application/vnd.api+json`.
- CLI uses cobra; configuration uses viper with TOML and documented keys.
- Migrations handled via golang-migrate.
- Single-binary distribution with embedded frontend assets.
- Real-time availability via WebSockets with polling fallback.
- Booking conflicts handled optimistically with unique constraint on (desk_id, booking_date).
- Bookings are full-day only; store a single booking_date per booking.
- Target builds: macOS (arm64) and Linux (amd64) only; Windows out of scope.
- No Docker or Kubernetes workflows.
- OpenAPI 3.1 docs in `api-doc/` with per-endpoint files; lint with Redocly.
- Vue 3 + Vuetify + Pinia + Vue Router; Composition API with `<script setup>`.
- Vitest for unit tests, Cypress for E2E with `data-cy` selectors and real API responses.
- Vite dev server proxies `/api` to `http://localhost:9900`.

### FR Coverage Map

FR1: Epic 1  
FR2: Epic 1  
FR3: Epic 1  
FR4: Epic 2  
FR5: Epic 2  
FR6: Epic 2  
FR7: Epic 2  
FR8: Epic 2  
FR9: Epic 3  
FR10: Epic 3  
FR11: Epic 3  
FR12: Epic 4  
FR13: Epic 4  
FR14: Epic 4  
FR15: Epic 5  
FR16: Epic 5  
FR17: Epic 6  
FR18: Epic 6  
FR19: Epic 6  
FR20: Epic 7  
FR21: Epic 7  
FR22: Epic 7  
FR23: Epic 7  
FR24: Epic 7  
FR25: Epic 8  
FR26: Epic 9  
FR27: Epic 9

## Epic List

### Epic 1: Secure Sign-In & Access Control

Users can authenticate via Entra ID and only permitted users can access SitHub.
**FRs covered:** FR1, FR2, FR3

### Epic 2: Space Discovery & Availability

Users can browse areas, rooms, desks, equipment, and availability.
**FRs covered:** FR4, FR5, FR6, FR7, FR8

### Epic 3: Single-Day Booking & Conflict Handling

Users can book a desk for a day and get clear messaging when a desk is taken.
**FRs covered:** FR9, FR10, FR11

### Epic 4: Booking Management & Admin Overrides

Users can view/cancel their bookings; admins can cancel any booking.
**FRs covered:** FR12, FR13, FR14

### Epic 5: Room & Presence Overviews

Users can view room booking summaries and today’s presence.
**FRs covered:** FR15, FR16

### Epic 6: Operator Configuration & Startup

Operators configure server and spaces via config files and changes apply on restart.
**FRs covered:** FR17, FR18, FR19

### Epic 7: Advanced Booking Options (Post-MVP)

Bookings on behalf, guests, recurring, history, notifications.
**FRs covered:** FR20, FR21, FR22, FR23, FR24

### Epic 8: Admin Management UI (Future)

Admins manage rooms/desks in a UI.
**FRs covered:** FR25

### Epic 9: Floor Maps & Analytics (Future)

Graphical floor map booking and analytics.
**FRs covered:** FR26, FR27

<!-- Repeat for each epic in epics_list (N = 1, 2, 3...) -->

## Epic 1 Stories: Secure Sign-In & Access Control

Users can authenticate via Entra ID and only permitted users can access SitHub.
**FRs covered:** FR1, FR2, FR3

### Story 1.1: Entra ID SSO Login

**FRs covered:** FR1

As an employee,
I want to sign in via Entra ID,
So that I can access SitHub without a separate account.

**Acceptance Criteria:**

**Given** I am not authenticated  
**When** I open SitHub  
**Then** I am redirected to Entra ID for login  
**And** after successful login I return to SitHub and see my name displayed

### Story 1.2: Role Determination from Entra ID Groups

**FRs covered:** FR2

As an admin,
I want my role determined from Entra ID group membership,
So that I see admin-only controls.

**Acceptance Criteria:**

**Given** my Entra ID account is in the admin group  
**When** I log in  
**Then** the system marks me as admin  
**And** admin-only cancellation controls are visible

### Story 1.3: Access Denied for Unauthorized Users

**FRs covered:** FR3

As a company operator,
I want unauthorized users blocked,
So that booking data is protected.

**Acceptance Criteria:**

**Given** my account is not permitted by Entra ID settings  
**When** I attempt to access SitHub  
**Then** I see an access-denied screen  
**And** API requests return a JSON:API error with 401/403 status

## Epic 2 Stories: Space Discovery & Availability

Users can browse areas, rooms, desks, equipment, and availability.
**FRs covered:** FR4, FR5, FR6, FR7, FR8

### Story 2.1: List Areas

**FRs covered:** FR4

As an employee,
I want to see the list of areas,
So that I can choose where to book.

**Acceptance Criteria:**

**Given** I am authenticated  
**When** I open the app  
**Then** I see all configured areas  
**And** the list is empty-safe (shows zero areas without error)

### Story 2.2: List Rooms in an Area

**FRs covered:** FR5

As an employee,
I want to see rooms for a selected area,
So that I can choose a room.

**Acceptance Criteria:**

**Given** I am viewing an area  
**When** I select the area  
**Then** I see only rooms belonging to that area  
**And** rooms outside the area are not shown

### Story 2.3: List Desks with Equipment

**FRs covered:** FR6, FR7

As an employee,
I want to see desks and their equipment for a room,
So that I can pick a suitable desk.

**Acceptance Criteria:**

**Given** I am viewing a room  
**When** I open the room  
**Then** I see the list of desks in that room  
**And** each desk shows its equipment list

### Story 2.4: Show Availability Status by Date

**FRs covered:** FR8

As an employee,
I want to see which desks are available for a selected date,
So that I can choose a free desk.

**Acceptance Criteria:**

**Given** I have selected a room and date  
**When** desks are displayed  
**Then** each desk shows available or occupied status for that date  
**And** status updates when the date changes

## Epic 3 Stories: Single-Day Booking & Conflict Handling

Users can book a desk for a day and get clear messaging when a desk is taken.
**FRs covered:** FR9, FR10, FR11

### Story 3.1: Create Single-Day Booking

**FRs covered:** FR9

As an employee,
I want to book a desk for a specific date,
So that I can reserve my workspace.

**Acceptance Criteria:**

**Given** I have selected a desk and date  
**When** I confirm the booking  
**Then** the booking is created for that date  
**And** it appears in “My Bookings”

### Story 3.2: Prevent Double-Booking

**FRs covered:** FR10

As an employee,
I want the system to prevent duplicate bookings for the same desk and day,
So that I don’t book a desk that’s already taken.

**Acceptance Criteria:**

**Given** a desk is already booked for a date  
**When** another booking is attempted for the same desk and date  
**Then** the request is rejected  
**And** no duplicate booking is created

### Story 3.3: Desk-Taken Feedback

**FRs covered:** FR11

As an employee,
I want a clear message when the desk becomes unavailable during booking,
So that I can choose another desk.

**Acceptance Criteria:**

**Given** I am booking a desk and it becomes unavailable  
**When** I submit the booking  
**Then** I see a message that the desk is no longer available for that date  
**And** I am prompted to choose another desk

## Epic 4 Stories: Booking Management & Admin Overrides

Users can view/cancel their bookings; admins can cancel any booking.
**FRs covered:** FR12, FR13, FR14

### Story 4.1: View My Bookings

**FRs covered:** FR12

As an employee,
I want to see my upcoming bookings,
So that I can confirm my reservations.

**Acceptance Criteria:**

**Given** I am authenticated  
**When** I open “My Bookings”  
**Then** I see a list of my future bookings  
**And** each entry includes desk, room, area, and date

### Story 4.2: Cancel My Booking

**FRs covered:** FR13

As an employee,
I want to cancel my booking,
So that I can free the desk if plans change.

**Acceptance Criteria:**

**Given** I have a future booking  
**When** I cancel it  
**Then** the booking is removed from my list  
**And** the desk becomes available for that date

### Story 4.3: Admin Cancel Any Booking

**FRs covered:** FR14

As an admin,
I want to cancel any booking,
So that I can resolve conflicts.

**Acceptance Criteria:**

**Given** I am an admin  
**When** I cancel another user’s booking  
**Then** the booking is removed from all relevant lists  
**And** the affected user sees the cancellation

## Epic 5 Stories: Room & Presence Overviews

Users can view room booking summaries and today’s presence.
**FRs covered:** FR15, FR16

### Story 5.1: Room Booking Overview

**FRs covered:** FR15

As an employee,
I want to see a room-level booking overview for a date,
So that I can understand room utilization.

**Acceptance Criteria:**

**Given** I select a room and date  
**When** I view the overview  
**Then** I see all booked desks and associated users for that date

### Story 5.2: Today’s Presence by Area

**FRs covered:** FR16

As an employee,
I want to see who is in the office today by area,
So that I can coordinate with colleagues.

**Acceptance Criteria:**

**Given** I view today’s presence for an area  
**When** the list is displayed  
**Then** I see all users with bookings in that area for today

## Epic 6 Stories: Operator Configuration & Startup

Operators configure server and spaces via config files and changes apply on restart.
**FRs covered:** FR17, FR18, FR19

### Story 6.1: Load Server Configuration

**FRs covered:** FR17

As an operator,
I want the server to load settings from a config file,
So that I can control listen address, port, and data directory.

**Acceptance Criteria:**

**Given** a valid configuration file  
**When** the server starts  
**Then** the server loads the settings  
**And** invalid settings prevent startup with a clear error

### Story 6.2: Load Space Configuration

**FRs covered:** FR18

As an operator,
I want areas, rooms, desks, and equipment loaded from a config file,
So that space definitions are centrally managed.

**Acceptance Criteria:**

**Given** a valid space configuration file  
**When** the server starts  
**Then** the UI reflects the configured areas, rooms, desks, and equipment

### Story 6.3: Apply Configuration on Restart

**FRs covered:** FR19

As an operator,
I want configuration changes to apply on restart,
So that I can update spaces without manual migration steps.

**Acceptance Criteria:**

**Given** the config file has changed  
**When** the server restarts  
**Then** the new configuration is applied  
**And** no manual data migration steps are required

## Epic 7 Stories: Advanced Booking Options (Post-MVP)

Bookings on behalf, guests, recurring, history, notifications.
**FRs covered:** FR20, FR21, FR22, FR23, FR24

### Story 7.1: Book on Behalf of Another User

**FRs covered:** FR20

As an employee,
I want to book a desk on behalf of another user,
So that we can sit together.

**Acceptance Criteria:**

**Given** I book a desk for another user  
**When** the booking is created  
**Then** it appears in both users’ booking lists  
**And** either user can cancel it

### Story 7.2: Guest Booking

**FRs covered:** FR21

As an employee,
I want to book a desk for a guest,
So that visitors can reserve a spot.

**Acceptance Criteria:**

**Given** I create a guest booking  
**When** the booking is saved  
**Then** the guest name and contact are stored  
**And** the booking is labeled as a guest booking in overviews

### Story 7.3: Multi-Day or Recurring Booking

**FRs covered:** FR22

As an employee,
I want to book multiple days or a recurring schedule,
So that I can plan ahead.

**Acceptance Criteria:**

**Given** I select multiple dates or a recurrence  
**When** I submit the booking  
**Then** the system creates individual daily bookings  
**And** conflicts are reported per day

### Story 7.4: Booking History

**FRs covered:** FR23

As an employee,
I want to view my booking history,
So that I can review past reservations.

**Acceptance Criteria:**

**Given** I open booking history  
**When** I filter by date range  
**Then** I see past bookings within that range

### Story 7.5: Booking Notifications

**FRs covered:** FR24

As an employee,
I want to receive booking notifications,
So that I stay informed about changes.

**Acceptance Criteria:**

**Given** I create or cancel a booking  
**When** the action completes  
**Then** a notification is sent within 5 minutes via the configured channel

## Epic 8 Stories: Admin Management UI (Future)

Admins manage rooms/desks in a UI.
**FRs covered:** FR25

### Story 8.1: Manage Rooms and Desks in Admin UI

**FRs covered:** FR25

As an admin,
I want to add, edit, and remove rooms and desks,
So that I can keep space definitions up to date.

**Acceptance Criteria:**

**Given** I am an admin  
**When** I add, edit, or remove a room or desk  
**Then** changes are saved  
**And** discovery lists reflect the updates

## Epic 9 Stories: Floor Maps & Analytics (Future)

Graphical floor map booking and analytics.
**FRs covered:** FR26, FR27

### Story 9.1: Book via Floor Map

**FRs covered:** FR26

As an employee,
I want to select a desk from a floor map,
So that I can book visually.

**Acceptance Criteria:**

**Given** I open the floor map  
**When** I select a desk and date  
**Then** the desk can be booked for that date

### Story 9.2: Usage Analytics

**FRs covered:** FR27

As an admin,
I want to view usage summaries by area, room, and date range,
So that I can understand utilization.

**Acceptance Criteria:**

**Given** I open the analytics view  
**When** I select area, room, and date range  
**Then** I see usage summaries for the selection

## Epic 10 Stories: UI/UX Redesign

Transform the application from basic Vuetify defaults into a polished, modern desk booking experience with consistent design language, reusable components, and excellent mobile support.
**FRs covered:** NFR6 (Accessibility), enhances all existing FRs

### Story 10.1: Design System Foundation

**FRs covered:** NFR6

As a user,
I want a visually consistent and branded experience,
So that the application feels professional and trustworthy.

**Acceptance Criteria:**

**Given** the application loads  
**When** I view any page  
**Then** I see consistent colors, typography, and spacing throughout  
**And** the color scheme reflects a professional brand identity  
**And** a logo and favicon are displayed

**Technical Requirements:**
- Custom Vuetify theme in `web/src/plugins/vuetify.ts`
- Color palette: primary, secondary, success, warning, error, surface colors
- Typography: Inter font family with defined scale
- Spacing tokens following 4px base unit
- Logo SVG and favicon in `web/public/`

### Story 10.2: Reusable Component Library

**FRs covered:** NFR6

As a developer,
I want reusable UI components,
So that the interface is consistent and maintainable.

**Acceptance Criteria:**

**Given** I am building a view  
**When** I need common UI patterns  
**Then** I can import pre-built components from `web/src/components/`  
**And** components follow the design system

**Components to create:**
- `PageHeader.vue` - Page title, breadcrumbs, actions
- `EmptyState.vue` - Illustrated empty states with action
- `LoadingState.vue` - Skeleton loaders matching content layout
- `ConfirmDialog.vue` - Confirmation modal with customizable actions
- `DatePicker.vue` - Vuetify date picker with consistent styling
- `DateRangePicker.vue` - Date range selection for filters
- `StatusChip.vue` - Consistent status indicators (available, booked, etc.)

### Story 10.3: Navigation & Layout Redesign

**FRs covered:** NFR6

As a user,
I want clear navigation and context awareness,
So that I always know where I am and can easily move around.

**Acceptance Criteria:**

**Given** I am on any page  
**When** I look at the navigation  
**Then** I see the current page highlighted in the nav  
**And** I see breadcrumbs showing my location in the hierarchy  
**And** I can access main sections (Areas, My Bookings, History) from any page

**Given** I am on a mobile device  
**When** I open the navigation  
**Then** I see a drawer menu that works well on small screens

**Technical Requirements:**
- Redesigned `App.vue` with improved header
- Breadcrumb component integrated into layout
- Mobile navigation drawer
- User menu with name and logout

### Story 10.4: Space Discovery Views Redesign

**FRs covered:** FR4, FR5, FR6, FR7, FR8, NFR6

As a user,
I want visually appealing space discovery,
So that browsing areas, rooms, and desks is enjoyable and efficient.

**Acceptance Criteria:**

**Given** I am viewing the areas list  
**When** the page loads  
**Then** I see areas displayed as cards with visual hierarchy  
**And** empty state shows an illustration and helpful message

**Given** I am viewing rooms in an area  
**When** the page loads  
**Then** I see room cards with desk availability summary  
**And** breadcrumbs show: Home > [Area Name]

**Given** I am viewing desks in a room  
**When** the page loads  
**Then** I see desks as visual cards with clear status indicators  
**And** equipment and warnings are displayed attractively  
**And** available vs booked desks are visually distinct

### Story 10.5: Booking Flow Redesign

**FRs covered:** FR9, FR20, FR21, FR22, NFR6

As a user,
I want an intuitive and delightful booking experience,
So that reserving a desk feels effortless.

**Acceptance Criteria:**

**Given** I want to book a desk  
**When** I click the book button  
**Then** I see a clear booking dialog/flow  
**And** date selection uses a proper calendar picker  
**And** multi-day selection is visual (calendar-based, not text input)

**Given** I want to book for a colleague  
**When** I toggle "Book for someone else"  
**Then** I can search/select a colleague (or enter details)  
**And** the flow is clearly differentiated from personal booking

**Given** I complete a booking  
**When** the booking succeeds  
**Then** I see a success confirmation with booking details  
**And** I have clear next actions (view bookings, book another)

### Story 10.6: Booking Management Views Redesign

**FRs covered:** FR12, FR13, FR23, NFR6

As a user,
I want my bookings displayed beautifully,
So that managing my reservations is pleasant.

**Acceptance Criteria:**

**Given** I open My Bookings  
**When** the page loads  
**Then** I see bookings as cards with all relevant info  
**And** "booked by" and "guest" bookings are visually distinguished  
**And** cancel action has a confirmation dialog

**Given** I open Booking History  
**When** the page loads  
**Then** I see a date range picker for filtering  
**And** past bookings are displayed in a clean list/table  
**And** empty state is handled gracefully

### Story 10.7: Mobile Responsiveness

**FRs covered:** NFR6

As a mobile user,
I want the app to work well on my phone,
So that I can book desks on the go.

**Acceptance Criteria:**

**Given** I access the app on a mobile device  
**When** I view any page  
**Then** the layout adapts to the screen size  
**And** touch targets are appropriately sized (min 44px)  
**And** navigation is accessible via drawer menu  
**And** forms and dialogs are usable on small screens

**Technical Requirements:**
- Responsive breakpoints for all views
- Touch-friendly interactions
- Viewport-appropriate font sizes
- No horizontal scrolling on mobile
