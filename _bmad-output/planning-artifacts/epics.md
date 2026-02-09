---
stepsCompleted: [step-01-validate-prerequisites, step-02-design-epics, step-03-create-stories, step-04-final-validation]
inputDocuments:
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/prd.md
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/architecture.md
lastEdited: '2026-02-08'
editHistory:
  - date: '2026-02-07'
    changes: "Updated Epic 1 for dual-source auth (Entra ID + local). Added FR28-FR35. Added Epic 11: User Management & Local Authentication with 8 stories. Updated NFR3, additional requirements, and coverage map."
  - date: '2026-02-08'
    changes: "Domain rename: reworded FR4-FR19 (rooms/desks to items). Added FR36-FR42 (weekly availability, booking notes, week booking, booker display, breadcrumbs, schema normalization, UI labels). Added Epic 12: Domain Rename & Data Normalization and Epic 13: Enhanced Booking Experience."
---

# sithub - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for sithub, decomposing the requirements from the PRD and
Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: Users can authenticate via Entra ID SSO or local credentials (email and password).
Acceptance: the login page presents a username/password form and a "Login via Entra ID"
button; both methods result in an authenticated session with the user's name displayed.
FR2: The system determines user roles (regular vs admin) based on authentication source.
For Entra ID users, admin status is synced from group membership on every login and cannot
be changed locally. For local users, admin status is managed locally by administrators.
Acceptance: admins see admin-only controls; regular users do not; Entra ID admin status
reflects current group membership after each login.
FR3: Users can access the application only if they are authenticated. Acceptance:
unauthenticated users see only the login page and cannot view any booking data.
FR4: Users can view a list of available areas. Acceptance: after login, the UI lists all
configured areas.
FR5: Users can view a list of item groups within a selected area. Acceptance: selecting an
area shows only its item groups.
FR6: Users can view a list of items within a selected item group. Acceptance: selecting an
item group lists its items.
FR7: Users can view equipment details for each item. Acceptance: each item entry shows its
equipment list if configured.
FR8: Users can view current booking status for items. Acceptance: item entries show
available/occupied status for the selected date.
FR9: Users can book an item for a single day. Acceptance: selecting an item and date creates
a booking that appears in "My Bookings."
FR10: The system prevents double-booking of the same item for the same day. Acceptance: the
second attempt is rejected and no duplicate booking is created.
FR11: Users receive a message when a selected item becomes unavailable during booking.
Acceptance: the message states the item is no longer available for that date and prompts
the user to choose another item.
FR12: Users can view their upcoming bookings ("My Bookings"). Acceptance: the list includes
item, item group, area, and date for each future booking; booking notes are displayed if
present.
FR13: Users can cancel their own bookings. Acceptance: cancelling removes the booking from
all relevant lists and frees the item.
FR14: Admin users can cancel any booking. Acceptance: admins can cancel another user's
booking and the affected user sees the cancellation reflected in their list.
FR15: Users can view an item-group-level booking overview. Acceptance: for a selected item
group and date, the overview lists all booked items and associated users; booking notes are
displayed if present.
FR16: Users can view "Today's presence" for an area (who is in the office today).
Acceptance: the view lists all users with bookings in that area for today; booking notes
are displayed if present.
FR17: Operators can configure server settings via a configuration file. Acceptance: invalid
settings prevent startup with a descriptive error; valid settings allow startup.
FR18: Operators can configure areas, item groups, items, and equipment via a configuration
file. Acceptance: after restart, the UI reflects the updated space definitions.
FR19: The system can load and apply configuration changes on startup. Acceptance:
configuration changes take effect after restart without manual data migration steps.
FR20: Users can book on behalf of another user. Acceptance: the booking appears in both
users' booking lists and either can cancel.
FR21: Users can book items for guests outside the organization. Acceptance: a guest booking
stores guest name and contact and is visible as a guest booking in overviews.
FR22: Users can book multi-day or recurring reservations. Acceptance: the system creates
individual daily bookings and reports conflicts per day.
FR23: Users can view booking history. Acceptance: users can see past bookings with date
range filtering.
FR24: Users can receive notifications related to bookings. Acceptance: booking
creation/cancellation triggers a notification via the configured channel within 5 minutes.
FR25: Admins can manage item groups and items via an admin UI. Acceptance: admins can
add/edit/remove item groups and items; changes appear in discovery lists after save.
FR26: Users can book items using a graphical floor-map view. Acceptance: an item selected
on the map can be booked for a chosen date.
FR27: Admins can access advanced reporting and analytics. Acceptance: admins can view usage
summaries by area/item group and date range.
FR28: All users (Entra ID and local) are stored in a users table. Entra ID users are
inserted on first login and updated on subsequent logins. Acceptance: after login, the user
exists in the users table with correct source, name, and email.
FR29: Email addresses are unique across all authentication sources. Creating a local user
with an email that exists for an Entra ID user is blocked, and vice versa. Acceptance:
attempting to create a duplicate email returns an error regardless of authentication source.
FR30: Local users can log in with email and password. Acceptance: entering valid credentials
in the login form creates an authenticated session; invalid credentials show a descriptive
error message.
FR31: Local users can change their own password via the `/me` endpoint. Acceptance: after
changing the password, the old password no longer works and the new password (minimum 14
characters) is accepted.
FR32: Admin users can reset the password of any local user via the `/users/{id}` endpoint.
Acceptance: the affected user can log in with the new password; Entra ID user passwords
cannot be reset this way.
FR33: The system provides a `/me` endpoint returning the current user's profile information.
Acceptance: authenticated requests to `/me` return the user's id, email, name, role, and
authentication source.
FR34: The system provides a `/users` endpoint for user management. Acceptance: admins can
list, create, read, update, and delete local users; non-admin users can only read. Entra ID
users cannot be created or deleted via this endpoint.
FR35: A demo users SQL file is provided for development setup. Acceptance: running the SQL
file creates 15 users (2 admins, 13 regular users) with local credentials in the database.
FR36: Users can view a weekly availability preview for item groups. Acceptance: the item
group list view includes a calendar week selector (next 8 weeks, current week pre-selected);
each item group tile displays weekday indicators (MO-FR) colored green (at least one item
available) or red (fully booked) for the selected week; the week selector displays dates in
locale-aware format with week number (e.g., "2026-03-16 - Week 12").
FR37: Users can add, view, and edit free-text notes on their bookings. Acceptance: after
completing a booking, a confirmation message includes an "add note" action; notes are
visible in My Bookings, Today's Presence, and item detail views; notes longer than the
display width are truncated with an expand indicator that opens the full text; notes are
editable from the My Bookings view.
FR38: Users can toggle between day booking mode and week booking mode. Acceptance: the
selected mode is persisted in browser local storage and restored on next visit; in week
mode, the date selector becomes a calendar week selector (next 8 weeks); item tiles show
per-day checkboxes with booker names; existing bookings by other users cannot be unchecked;
a single "Confirm My Booking" button submits all selected days at once.
FR39: Users can see the display name of the person who booked an item. Acceptance: in the
item detail view, booked items show the booker's display name alongside the booking status.
FR40: Breadcrumbs in the navigation hierarchy are clickable and navigate to the
corresponding view. Acceptance: clicking any breadcrumb segment navigates to that level of
the hierarchy; the current page breadcrumb is not clickable.
FR41: The bookings table references users by `user_id` only; display names are resolved
via JOIN with the users table. Acceptance: the bookings table does not store denormalized
user name columns; all booking queries that require user names perform a JOIN; existing
bookings are migrated to remove redundant columns.
FR42: UI action labels use domain-neutral terminology. Acceptance: labels read "BOOK"
instead of "VIEW DESK" or "VIEW ROOM"; "BOOK THIS ITEM" instead of "BOOK THIS DESK";
booking confirmation messages reference the item name from the configuration (e.g.,
"Parking Lot 1 booked successfully") rather than the generic term "desk."

### NonFunctional Requirements

NFR1: For expected usage (<=50 concurrent users), list navigation actions complete within
2 seconds at p95; booking and cancellation complete within 3 seconds at p95.
NFR2: The system can be restarted without data loss; bookings remain intact after restart
and conflicts do not create partial records.
NFR3: All booking data requires authenticated access (Entra ID or local credentials);
unauthenticated requests are denied. Local user passwords are stored as bcrypt hashes;
plaintext passwords are never persisted or logged. Minimum password length is 14 characters.
Data at rest is stored without application-layer encryption; in-transit encryption is managed
outside the application.
NFR4: Single-node deployment is sufficient; no clustering or horizontal scaling is required
for MVP usage levels.
NFR5: Meets WCAG A: all interactive elements have accessible names, keyboard focus is
visible, and form inputs are labeled.

### Additional Requirements

- Go 1.25 with Echo, SQLite (WAL), and JSON:API responses using `application/vnd.api+json`.
- CLI uses cobra; configuration uses viper with TOML and documented keys.
- Migrations handled via golang-migrate.
- Single-binary distribution with embedded frontend assets.
- Real-time availability via WebSockets with polling fallback.
- Booking conflicts handled optimistically with unique constraint on (item_id, booking_date).
- Bookings are full-day only; store a single booking_date per booking.
- Target builds: macOS (arm64) and Linux (amd64) only; Windows out of scope.
- No Docker or Kubernetes workflows.
- Dual-source authentication: Entra ID SSO (optional) and local credentials (always available).
- Central users table storing both Entra ID and local users with `user_source` enum.
- Unified session mechanism: both auth paths produce gorilla/securecookie signed cookies.
- bcrypt password hashing for local users; minimum 14 characters.
- Admin role sync: Entra ID users sync `is_admin` from group membership on every login;
  local admin managed via DB.
- Email uniqueness enforced at DB level across authentication sources.
- `test_auth` mechanism removed; replaced by real local users.
- Demo users SQL file (`tools/database/demo-users.sql`) with 15 users for development and
  testing.
- OpenAPI 3.1 docs in `api-doc/` with per-endpoint files; lint with Redocly.
- Vue 3 + Vuetify + Pinia + Vue Router; Composition API with `<script setup>`.
- Vitest for unit tests, Cypress for E2E with `data-cy` selectors and real API responses.
- Vite dev server proxies `/api` to `http://localhost:9900`.
- Domain-neutral terminology: "items" instead of "desks", "item groups" instead of "rooms"
  throughout codebase, API, and UI.
- Bookings table normalized: user names resolved via JOIN, no denormalized name columns.

### FR Coverage Map

FR1: Epic 1 - Dual-source auth
FR2: Epic 1 - Role determination
FR3: Epic 1 - Access control
FR4: Epic 2 - List areas (reworded in Epic 12)
FR5: Epic 2 - List item groups (reworded in Epic 12)
FR6: Epic 2 - List items (reworded in Epic 12)
FR7: Epic 2 - Item equipment (reworded in Epic 12)
FR8: Epic 2 - Booking status (reworded in Epic 12)
FR9: Epic 3 - Single-day booking (reworded in Epic 12)
FR10: Epic 3 - Prevent double-booking (reworded in Epic 12)
FR11: Epic 3 - Item-taken feedback (reworded in Epic 12)
FR12: Epic 4 - My Bookings (reworded in Epic 12)
FR13: Epic 4 - Cancel booking (reworded in Epic 12)
FR14: Epic 4 - Admin cancel (reworded in Epic 12)
FR15: Epic 5 - Item group booking overview (reworded in Epic 12)
FR16: Epic 5 - Today's presence (reworded in Epic 12)
FR17: Epic 6 - Server configuration
FR18: Epic 6 - Space configuration (reworded in Epic 12)
FR19: Epic 6 - Apply on restart
FR20: Epic 7 - Book on behalf
FR21: Epic 7 - Guest bookings
FR22: Epic 7 - Multi-day bookings
FR23: Epic 7 - Booking history
FR24: Epic 7 - Notifications
FR25: Epic 8 - Admin management UI
FR26: Epic 9 - Floor maps
FR27: Epic 9 - Analytics
FR28: Epic 11 - Users table
FR29: Epic 11 - Email uniqueness
FR30: Epic 11 - Local login
FR31: Epic 11 - Password change
FR32: Epic 11 - Admin password reset
FR33: Epic 11 - /me endpoint
FR34: Epic 11 - /users endpoint
FR35: Epic 11 - Demo users
FR36: Epic 13 - Weekly availability preview
FR37: Epic 13 - Booking notes
FR38: Epic 13 - Week booking mode
FR39: Epic 13 - Booker display name
FR40: Epic 13 - Clickable breadcrumbs
FR41: Epic 12 - Schema normalization
FR42: Epic 12 - UI label consistency

## Epic List

### Epic 1: Dual-Source Authentication & Access Control

Users can authenticate via Entra ID SSO or local credentials, and only authenticated users
can access SitHub.
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

### Epic 11: User Management & Local Authentication

User management API, local login, password management, and demo users for development.
**FRs covered:** FR28, FR29, FR30, FR31, FR32, FR33, FR34, FR35

### Epic 12: Domain Rename & Data Normalization

Users interact with consistent, domain-neutral terminology that supports booking any kind
of resource (desks, parking lots, lab benches), not just desks and rooms. The data model is
normalized to eliminate redundant columns.
**FRs covered:** FR41, FR42 (plus codebase migration of reworded FR4-FR19)

### Epic 13: Enhanced Booking Experience

Users get powerful new booking capabilities: weekly availability previews, booking notes,
week-at-a-time booking mode, visibility of who booked what, and clickable navigation
breadcrumbs.
**FRs covered:** FR36, FR37, FR38, FR39, FR40

<!-- Repeat for each epic in epics_list (N = 1, 2, 3...) -->

## Epic 1 Stories: Dual-Source Authentication & Access Control

Users can authenticate via Entra ID SSO or local credentials, and only authenticated users
can access SitHub.
**FRs covered:** FR1, FR2, FR3

### Story 1.1: Dual-Source Login

**FRs covered:** FR1

As an employee,
I want to sign in via Entra ID or local credentials,
So that I can access SitHub regardless of my company's identity provider.

**Acceptance Criteria:**

**Given** I am not authenticated
**When** I open SitHub
**Then** I see a login page with a username/password form and a "Login via Entra ID" button

**Given** I click "Login via Entra ID"
**When** I complete the Entra ID flow
**Then** I return to SitHub with my name displayed

**Given** I enter valid local credentials in the login form
**When** I submit the form
**Then** I am authenticated and see my name displayed

**Given** I enter invalid local credentials
**When** I submit the form
**Then** I see a descriptive error message

### Story 1.2: Source-Dependent Role Determination

**FRs covered:** FR2

As an admin,
I want my role determined from my authentication source,
So that I see admin-only controls.

**Acceptance Criteria:**

**Given** my Entra ID account is in the admin group
**When** I log in via Entra ID
**Then** the system syncs my admin status from group membership
**And** admin-only controls are visible

**Given** I am a local user with admin role in the database
**When** I log in with local credentials
**Then** admin-only controls are visible

**Given** I am removed from the Entra ID admin group
**When** I log in again
**Then** admin-only controls are no longer visible

### Story 1.3: Access Denied for Unauthenticated Users

**FRs covered:** FR3

As a company operator,
I want unauthenticated users blocked,
So that booking data is protected.

**Acceptance Criteria:**

**Given** I am not authenticated
**When** I attempt to access any protected page
**Then** I am redirected to the login page

**Given** I am not authenticated
**When** I make an API request to a protected endpoint
**Then** the API returns a JSON:API error with 401 status

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
**And** it appears in "My Bookings"

**Given** I attempt to book a desk for a past date
**When** I submit the booking
**Then** the system rejects the booking
**And** I see a clear error message that past dates cannot be booked

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
**When** I open "My Bookings"
**Then** I see a list of my future bookings
**And** each entry includes desk, room, area, and date

**Given** I have no upcoming bookings
**When** I open "My Bookings"
**Then** I see an empty state with a helpful message and action to book a desk

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

**Given** I am an admin viewing a room booking overview (Epic 5)
**When** I see another user's booking
**Then** I see an admin cancel action on that booking

**Given** I am an admin
**When** I cancel another user's booking
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

Transform the application from basic Vuetify defaults into a polished, modern desk booking
experience with consistent design language, reusable components, and excellent mobile support.
**FRs covered:** NFR5 (Accessibility), enhances all existing FRs

### Story 10.1: Design System Foundation

**FRs covered:** NFR5

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

**FRs covered:** NFR5

As a user,
I want a consistent UI experience across all pages,
So that the application feels polished and predictable.

**Acceptance Criteria:**

**Given** I navigate between different views
**When** I interact with common UI patterns (empty states, loading, confirmations)
**Then** they look and behave consistently across the application
**And** all components follow the design system from Story 10.1

**Components to create:**

- `PageHeader.vue` - Page title, breadcrumbs, actions
- `EmptyState.vue` - Illustrated empty states with action
- `LoadingState.vue` - Skeleton loaders matching content layout
- `ConfirmDialog.vue` - Confirmation modal with customizable actions
- `DatePicker.vue` - Vuetify date picker with consistent styling
- `DateRangePicker.vue` - Date range selection for filters
- `StatusChip.vue` - Consistent status indicators (available, booked, etc.)

### Story 10.3: Navigation & Layout Redesign

**FRs covered:** NFR5

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

**FRs covered:** FR4, FR5, FR6, FR7, FR8, NFR5

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

**FRs covered:** FR9, NFR5

As a user,
I want an intuitive and delightful booking experience,
So that reserving a desk feels effortless.

**Acceptance Criteria:**

**Given** I want to book a desk
**When** I click the book button
**Then** I see a clear booking dialog/flow
**And** date selection uses a proper calendar picker

**Given** I complete a booking
**When** the booking succeeds
**Then** I see a success confirmation with booking details
**And** I have clear next actions (view bookings, book another)

### Story 10.6: Booking Management Views Redesign

**FRs covered:** FR12, FR13, NFR5

As a user,
I want my bookings displayed beautifully,
So that managing my reservations is pleasant.

**Acceptance Criteria:**

**Given** I open My Bookings
**When** the page loads
**Then** I see bookings as cards with all relevant info
**And** cancel action has a confirmation dialog

**Given** I have no upcoming bookings
**When** I open My Bookings
**Then** I see a helpful empty state with an action to book a desk

### Story 10.7: Mobile Responsiveness

**FRs covered:** NFR5

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

## Epic 11 Stories: User Management & Local Authentication

User management API, local login, password management, and demo users for development.
**FRs covered:** FR28, FR29, FR30, FR31, FR32, FR33, FR34, FR35

### Story 11.1: Users Table and Entra ID User Sync

**FRs covered:** FR28

As an operator,
I want all users stored in a central users table,
So that the system has a unified user directory regardless of authentication source.

**Acceptance Criteria:**

**Given** an Entra ID user logs in for the first time
**When** the login completes
**Then** the user is inserted into the users table with source "entraid", name, and email

**Given** an Entra ID user logs in again
**When** the login completes
**Then** the user's name and admin status are updated from Entra ID

**Given** a local user is created via the API
**When** the creation succeeds
**Then** the user exists in the users table with source "internal"

### Story 11.2: Email Uniqueness Across Sources

**FRs covered:** FR29

As an operator,
I want email addresses unique across all authentication sources,
So that identity conflicts are prevented.

**Acceptance Criteria:**

**Given** an Entra ID user exists with email `alex@example.com`
**When** an admin attempts to create a local user with the same email
**Then** the request is rejected with a JSON:API error

**Given** a local user exists with email `dana@example.com`
**When** an Entra ID user with the same email logs in for the first time
**Then** the login fails with a descriptive error

### Story 11.3: Local User Login

**FRs covered:** FR30

As a local user,
I want to log in with my email and password,
So that I can use SitHub without Entra ID.

**Acceptance Criteria:**

**Given** I am a local user with valid credentials
**When** I enter my email and password in the login form and submit
**Then** I am authenticated and see my name displayed

**Given** I enter an incorrect password
**When** I submit the login form
**Then** I see a descriptive error message
**And** no session is created

### Story 11.4: Self-Service Password Change

**FRs covered:** FR31

As a local user,
I want to change my own password,
So that I can maintain my account security.

**Acceptance Criteria:**

**Given** I am authenticated as a local user
**When** I submit a password change via `/me` with a new password of 14+ characters
**Then** the password is updated
**And** the old password no longer works

**Given** I submit a new password shorter than 14 characters
**When** the request is processed
**Then** it is rejected with a validation error

**Given** I am an Entra ID user
**When** I attempt to change my password via `/me`
**Then** the request is rejected (Entra ID passwords are managed externally)

### Story 11.5: Admin Password Reset

**FRs covered:** FR32

As an admin,
I want to reset any local user's password,
So that I can help users who are locked out.

**Acceptance Criteria:**

**Given** I am an admin
**When** I reset a local user's password via `/users/{id}`
**Then** the user can log in with the new password

**Given** I attempt to reset an Entra ID user's password
**When** the request is processed
**Then** it is rejected with a JSON:API error

### Story 11.6: Current User Profile Endpoint

**FRs covered:** FR33

As an authenticated user,
I want to retrieve my profile information,
So that the UI can display my identity and role.

**Acceptance Criteria:**

**Given** I am authenticated
**When** I request `/me`
**Then** I receive my id, email, name, role, and authentication source

**Given** I am not authenticated
**When** I request `/me`
**Then** I receive a 401 JSON:API error

### Story 11.7: User Management API

**FRs covered:** FR34

As an admin,
I want to manage local users via the API,
So that I can create, update, and remove user accounts.

**Acceptance Criteria:**

**Given** I am an admin
**When** I list users via `GET /users`
**Then** I see all users (Entra ID and local) with their source and role

**Given** I am an admin
**When** I create a local user via `POST /users`
**Then** the user is created with source "internal" and a hashed password

**Given** I am an admin
**When** I attempt to create or delete an Entra ID user via `/users`
**Then** the request is rejected with a JSON:API error

**Given** I am a non-admin user
**When** I attempt to create, update, or delete a user via `/users`
**Then** the request is rejected with a 403 JSON:API error

**Given** I am a non-admin user
**When** I list or read users via `/users`
**Then** I receive the data (read access is allowed for all authenticated users)

### Story 11.8: Demo Users SQL File

**FRs covered:** FR35

As a developer,
I want a demo users SQL file,
So that I can quickly set up a development environment with test data.

**Acceptance Criteria:**

**Given** the SQL file at `tools/database/demo-users.sql` exists
**When** it is executed against the database
**Then** 15 users are created: 2 admins and 13 regular users with local credentials
**And** all passwords are bcrypt-hashed

## Epic 12 Stories: Domain Rename & Data Normalization

Users interact with consistent, domain-neutral terminology that supports booking any kind
of resource (desks, parking lots, lab benches), not just desks and rooms. The data model is
normalized to eliminate redundant columns.
**FRs covered:** FR41, FR42 (plus codebase migration of reworded FR4-FR19)

### Story 12.1: Rename Backend Internals and Database Column

**FRs covered:** FR4-FR16 (reworded, partial)

As a developer,
I want the Go packages, structs, and database column to use domain-neutral terminology,
So that the codebase foundation is ready for the public API rename.

**Acceptance Criteria:**

**Given** the Go packages `internal/rooms/` and `internal/desks/` exist
**When** the rename is applied
**Then** they are consolidated or renamed to use "item group" and "item" terminology
**And** all Go struct fields, function names, and variables use the new terminology

**Given** the database has `desk_id` columns in the bookings table
**When** a new migration is applied
**Then** the column is renamed to `item_id`
**And** the unique constraint on (desk_id, booking_date) becomes (item_id, booking_date)
**And** existing bookings are preserved with correct references

**Given** the API routes remain unchanged in this story
**When** the internal rename is applied
**Then** the existing API routes still work (backward compatible)
**And** `go test -race ./...` succeeds
**And** `golangci-lint run ./...` reports no errors

### Story 12.2: Normalize Bookings Table

**FRs covered:** FR41

As a developer,
I want the bookings table to reference users by user_id only,
So that display names are always current and not duplicated across tables.

**Acceptance Criteria:**

**Given** the bookings table contains `user_name`, `booked_by_user_id`, and
`booked_by_user_name` columns
**When** a new database migration is applied
**Then** the redundant columns are removed from the bookings table
**And** existing data is preserved (user_id remains as the foreign key reference)

**Given** a booking query needs to display a user's name
**When** the query is executed
**Then** the display name is resolved via JOIN with the users table
**And** the name reflects the current value in the users table

**Given** the bookings API returns user information
**When** a booking list or detail is requested
**Then** user display names are included in the response via JOIN
**And** the JSON:API response structure remains consistent

**Given** the migration has been applied
**When** `go test -race ./...` is executed
**Then** all existing tests pass with the normalized schema
**And** booking creation and listing continue to work correctly

### Story 12.3: Update API Routes and YAML Configuration

**FRs covered:** FR4-FR16 (reworded, partial), FR18 (reworded)

As an operator,
I want the API routes and YAML configuration to use domain-neutral terminology,
So that the public interface reflects the flexible item model.

**Acceptance Criteria:**

**Given** the API routes use `/rooms/:room_id/desks` and `/areas/:area_id/rooms`
**When** the route rename is applied
**Then** the routes use `/item-groups/:item_group_id/items` and
`/areas/:area_id/item-groups`
**And** JSON:API resource types use the new terminology (e.g., `item-groups`, `items`)

**Given** the YAML configuration uses `rooms` and `desks` keys
**When** the spaces config loader is updated
**Then** it reads `items` keys at both hierarchy levels (item groups and items)
**And** the `sithub_areas.example.yaml` is updated with the new keys
**And** the example includes diverse item types (office desks, parking lots)

**Given** the route rename is applied
**When** `go test -race ./...` is executed
**Then** all tests pass with the new route paths
**And** `golangci-lint run ./...` reports no errors

### Story 12.4: Rename Rooms and Desks to Items in Frontend

**FRs covered:** FR4-FR16 (reworded), FR42

As a user,
I want the UI to use domain-neutral terminology,
So that I see consistent labels regardless of what I am booking.

**Acceptance Criteria:**

**Given** the frontend uses components and routes named `RoomsView`, `DesksView`, etc.
**When** the rename is applied
**Then** components and routes use item-group and item terminology
**And** the Vue Router paths match the new API routes

**Given** the areas list view shows a "VIEW ROOM" button on each area tile
**When** the rename is applied
**Then** the button reads "BOOK"
**And** no UI element references "room" or "desk" in user-facing text

**Given** the item detail view shows a "BOOK THIS DESK" button
**When** the rename is applied
**Then** the button reads "BOOK THIS ITEM"

**Given** a booking is successfully created
**When** the confirmation message is displayed
**Then** the message references the item name from the configuration
(e.g., "Parking Lot 1 booked successfully") rather than the generic term "desk"

**Given** the Pinia stores and API service files reference rooms/desks
**When** the rename is applied
**Then** all store names, API paths, and TypeScript types use the new terminology
**And** `npm run type-check` and `npm run build` succeed without errors
**And** `npx vitest run` passes all unit tests

### Story 12.5: Update Documentation and E2E Tests

**FRs covered:** FR42 (partial)

As a developer,
I want API documentation and E2E tests updated to use the new terminology,
So that the entire codebase is consistent and all tests pass.

**Acceptance Criteria:**

**Given** the OpenAPI documentation references `/rooms` and `/desks` endpoints
**When** the documentation is updated
**Then** all endpoint paths, schemas, and descriptions use the new terminology
**And** `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml` passes

**Given** the Cypress E2E tests reference rooms/desks in selectors and assertions
**When** the tests are updated
**Then** all E2E tests use the new terminology in routes, selectors, and assertions
**And** `npm run test:e2e -- --browser electron` passes all tests

**Given** the Go code duplication check runs
**When** `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 3` is executed
**Then** the duplication threshold is not exceeded

**Given** the TypeScript code duplication check runs
**When** `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0` is
executed
**Then** the duplication threshold is not exceeded

## Epic 13 Stories: Enhanced Booking Experience

Users get powerful new booking capabilities: weekly availability previews, booking notes,
week-at-a-time booking mode, visibility of who booked what, and clickable navigation
breadcrumbs.
**FRs covered:** FR36, FR37, FR38, FR39, FR40

### Story 13.1: Clickable Breadcrumbs

**FRs covered:** FR40

As a user,
I want breadcrumbs to be clickable links,
So that I can navigate back to any level of the hierarchy quickly.

**Acceptance Criteria:**

**Given** I am viewing items within an item group
**When** I see the breadcrumb showing "Home > Office 1st Floor > Room 101"
**Then** "Home" and "Office 1st Floor" are clickable links
**And** clicking "Home" navigates to the areas list
**And** clicking "Office 1st Floor" navigates to the item groups for that area
**And** "Room 101" (the current page) is not clickable

**Given** I am viewing item groups within an area
**When** I see the breadcrumb showing "Home > Office 1st Floor"
**Then** "Home" is a clickable link that navigates to the areas list
**And** "Office 1st Floor" (the current page) is not clickable

### Story 13.2: Booker Display Name

**FRs covered:** FR39

As a user,
I want to see who has booked an item,
So that I know which colleagues are in the office or using a resource.

**Acceptance Criteria:**

**Given** I am viewing items in an item group for a specific date
**When** an item is booked
**Then** I see the booker's display name alongside the booking status
**And** the display name is resolved from the users table (not stored in the booking)

**Given** an item is available (not booked)
**When** the item is displayed
**Then** no booker name is shown
**And** the item is clearly marked as available

**Given** I am viewing Today's Presence for an area
**When** the presence list is displayed
**Then** each entry shows the user's display name

### Story 13.3: Weekly Availability Preview

**FRs covered:** FR36

As a user,
I want to see a weekly availability preview on item group tiles,
So that I can quickly identify which days have open items without clicking into each group.

**Acceptance Criteria:**

**Given** I am viewing item groups within an area
**When** the page loads
**Then** I see a calendar week selector above the item group tiles
**And** the current week is pre-selected
**And** only the next 8 weeks are available for selection
**And** each week option displays the Monday date and week number in locale-aware format
(e.g., "2026-03-16 - Week 12")

**Given** I have selected a calendar week
**When** the item group tiles are displayed
**Then** each tile shows weekday indicators (MO, TU, WE, TH, FR)
**And** green indicates at least one item is available on that day
**And** red indicates all items in the group are fully booked on that day
**And** each indicator uses a secondary visual cue in addition to color
(e.g., filled circle for available, empty circle for booked) to meet WCAG A

**Given** the backend receives a request for weekly availability
**When** `GET /api/v1/areas/:area_id/item-groups/availability?week=YYYY-Www` is called
**Then** it returns per-day availability counts for each item group in the area
**And** the response includes item group ID, day, total items, and available count
**And** the response completes within the NFR1 performance target (2s at p95)

**Given** I change the selected calendar week
**When** the new week is applied
**Then** the weekday indicators update to reflect availability for the new week

### Story 13.4: Booking Notes

**FRs covered:** FR37

As a user,
I want to add, view, and edit notes on my bookings,
So that I can communicate useful information to colleagues (e.g., "arriving after noon").

**Acceptance Criteria:**

**Given** I have just completed a booking
**When** the success confirmation message is displayed
**Then** I see an "add note" action within the confirmation
**And** clicking it opens a text input where I can type a free-text note
**And** the note is saved to the booking

**Given** I am viewing My Bookings
**When** a booking has a note
**Then** the note is displayed as a single line on the booking card
**And** if the note is longer than the available width, it is truncated
**And** a truncation indicator (icon) signals there is more text
**And** clicking the indicator opens a dialog (desktop) or bottom sheet (mobile)
showing the full note text

**Given** I am viewing My Bookings
**When** I want to edit a note on one of my bookings
**Then** I can modify the note text and save the changes
**And** the updated note is reflected immediately

**Given** I am viewing Today's Presence for an area
**When** a booking has a note
**Then** the note is displayed with the same truncation behavior as My Bookings

**Given** I am viewing items in an item group
**When** a booked item has a note
**Then** the note is displayed alongside the booker's name with truncation behavior

**Given** the backend receives a note update request
**When** the note is saved via the API
**Then** the booking record is updated with the new note text
**And** a `note` field is added to the bookings table via migration
**And** the JSON:API response includes the note in booking attributes

### Story 13.5: Week Booking Mode

**FRs covered:** FR38

As a user,
I want to switch to week booking mode and book multiple days at once,
So that I can reserve my workspace for an entire week efficiently.

**Acceptance Criteria:**

**Given** I am viewing items in an item group
**When** I see the booking mode toggle
**Then** I can switch between "book by day" and "book by week" modes
**And** the selected mode is persisted in browser local storage
**And** on my next visit, the previously selected mode is restored

**Given** I have selected week booking mode
**When** the date selector is displayed
**Then** it becomes a calendar week selector showing the next 8 weeks
**And** each week option displays the Monday date and week number

**Given** I have selected a week in week booking mode
**When** the item tiles are displayed
**Then** each tile shows a per-day breakdown (MO through FR) with checkboxes
**And** days booked by other users show the booker's name in red text
**And** days booked by other users have their checkboxes disabled (cannot be unchecked)
**And** days I have already booked show my name with a checked checkbox
**And** free days show "free" in green text with an unchecked checkbox

**Given** I am viewing week booking mode on a screen narrower than 600px
**When** the item tiles are displayed
**Then** the per-day breakdown uses a compact layout suitable for mobile
**And** touch targets meet the minimum 44px size requirement

**Given** I have checked one or more free days across one or more items
**When** I look below the item tiles
**Then** I see a single "Confirm My Booking" button
**And** the individual "BOOK THIS ITEM" buttons are not shown in week mode

**Given** I click "Confirm My Booking" with multiple days selected
**When** the bookings are submitted
**Then** each day/item combination is submitted as an individual API request
to `POST /api/v1/bookings`
**And** results are collected and displayed per-day
**And** each successful booking appears in My Bookings
**And** if any day fails (e.g., concurrent booking), the error is reported for that day
**And** successful bookings are not rolled back due to a single day's failure

**Given** I switch back to day booking mode
**When** the view updates
**Then** the standard single-day booking interface is restored
**And** the "BOOK THIS ITEM" button reappears on each item tile
