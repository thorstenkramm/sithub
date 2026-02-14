---
stepsCompleted: [step-01-init, step-02-discovery, step-03-success, step-04-journeys, step-05-domain, step-06-innovation, step-07-project-type, step-08-scoping, step-09-functional, step-10-nonfunctional, step-11-polish]
inputDocuments:
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/product-brief-sithub-2026-01-17.md
  - /Users/thorsten/projects/thorsten/sithub/docs/index.md
  - /Users/thorsten/projects/thorsten/sithub/docs/project-overview.md
  - /Users/thorsten/projects/thorsten/sithub/docs/architecture.md
  - /Users/thorsten/projects/thorsten/sithub/docs/source-tree-analysis.md
  - /Users/thorsten/projects/thorsten/sithub/docs/component-inventory.md
  - /Users/thorsten/projects/thorsten/sithub/docs/development-guide.md
  - /Users/thorsten/projects/thorsten/sithub/docs/deployment-guide.md
  - /Users/thorsten/projects/thorsten/sithub/docs/api-contracts-root.md
  - /Users/thorsten/projects/thorsten/sithub/docs/data-models-root.md
  - /Users/thorsten/projects/thorsten/sithub/README.md
documentCounts:
  briefCount: 1
  researchCount: 0
  brainstormingCount: 0
  projectDocsCount: 9
classification:
  projectType: web_app
  domain: general
  complexity: low
  projectContext: brownfield
date: 2026-01-17
lastEdited: 2026-02-14
workflowType: 'prd'
editHistory:
  - date: 2026-02-07
    changes: "Added multi-source auth (EntraID + local), user management API, password management. Removed test_auth. Fixed validation findings (implementation leakage, post-MVP traceability)."
  - date: 2026-02-08
    changes: "Domain rename: rooms/desks to generic items (FR4-FR18 reworded). Added FR36-FR42: weekly availability, booking notes, week booking mode, booker visibility, clickable breadcrumbs, schema normalization, UI label consistency. Resolved all 6 validation observations. Updated all user journeys and UX requirements."
  - date: 2026-02-14
    changes: "Added FR43-FR53: UI label simplification, booking form streamlining, collapsible item tiles, past date protection, theme selector, show weekends toggle, equipment filter. Added Epics 14-17."
---

# Product Requirements Document - sithub

**Author:** Thorsten
**Date:** 2026-01-17

## Executive Summary

SitHub is a SPA resource-booking web app for shared offices that replaces manual Confluence
tables with a fast, mobile-friendly booking flow. Employees can see real-time availability, book
items (desks, parking lots, lab benches, or any reservable resource) for a single day or an
entire week, and manage their bookings; admins can cancel any booking; IT configures space
definitions via YAML. SitHub supports dual authentication: Entra ID SSO for enterprise
environments and local credentials for teams without Entra ID. The MVP prioritizes a
low-friction UX and simple deployment (single executable), with success measured by user
preference over the Confluence workflow after a 5-day trial.

## Differentiators

- Single-executable distribution to minimize operational overhead
- Dual authentication: Entra ID SSO or local credentials (works with or without Entra ID)
- Domain-agnostic item booking (desks, parking lots, lab equipment, meeting rooms, etc.)
- Mobile-first, no-pinch booking experience
- File-based space configuration for simple administration

## Success Criteria

### User Success

- Users can book an item quickly and easily on mobile or desktop.
- Users describe SitHub as "worlds better than Confluence."
- After a 5-day trial, users do not want to return to the Confluence workflow.

### Business Success

- Post-trial preference >= 80% of test users favor SitHub.
- Reversion rate: 0 teams request a return to Confluence.

### Technical Success

- Both Entra ID and local authentication work reliably end-to-end.
- Application is fully usable for companies without Entra ID using local credentials only.
- Installation and setup documentation is clear and complete, including Entra ID steps and
  local user setup.
- Distribution is a single executable bundling frontend and backend.

### Measurable Outcomes

- Trial preference rate >= 80%.
- Reversion rate = 0.
- Successful Entra ID setup using documented steps in a test environment.
- Successful local-only deployment using documented steps without Entra ID configured.

## Product Scope

### MVP - Minimum Viable Product

- Dual authentication: Entra ID SSO or local credentials
- User management API (`/users` CRUD, `/me` endpoint)
- Self-service and admin password management for local users
- List view of areas and items (hierarchical) with equipment
- Book an item for a single day or an entire week
- Weekly availability preview for item groups
- Booking notes (add, view, edit)
- Cancel a booking
- Item group booking overview
- "Today's presence" view (who is in the office by area)

### Out of Scope (MVP)

- A UI to manage internal users is not desired for MVP.
- User self-registration is not supported; local users are created by admins or via SQL.
- "Book for colleague" and "Book for guest" UI is deferred; backend API exists but is not
  exposed in the frontend.

### Growth Features (Post-MVP)

- Booking on behalf of others (UI; backend API exists from MVP)
- Guest bookings (UI; backend API exists from MVP)
- Notifications
- Booking history
- Multi-day and recurring bookings

### Vision (Future)

- Graphical floor maps
- Admin management UI and advanced controls
- Advanced reporting and analytics

## User Journeys

### Journey 1: Employee (Happy Path) -- "Fast and Easy Booking"

**Opening Scene:**
Alex, an employee, opens SitHub on a phone or desktop to book for tomorrow.

**Rising Action:**

- Sees the login page with two options: a login form and an "Login via Entra ID" button
- Logs in with their preferred method
- Sees a clean list of areas -> selects one
- Sees item groups (e.g., rooms, parking sections) -> selects one
- Sees items (e.g., desks, parking lots) with equipment and availability
- Optionally checks the weekly availability preview to find the best day
- Books an item for a single day (or toggles to week mode and books multiple days at once)
- Adds an optional note to the booking (e.g., "arriving after noon")
- Visits "My Bookings" to confirm

**Climax (Aha Moment):**
"No pinching or zooming needed -- everything is readable and actionable on mobile."

**Resolution:**
Alex feels it was fast and easy and trusts the booking is secured.

---

### Journey 2: Employee (Edge Case) -- "Item Taken"

**Opening Scene:**
Alex selects an item and is about to book.

**Rising Action:**
Another user books the same item moments earlier.

**Climax:**
The system responds with a clear message: "Sorry. Someone else already picked this item."

**Resolution:**
Alex quickly picks another item and completes booking.

---

### Journey 3: Admin/Operations -- "Resolve Conflicts"

**Opening Scene:**
Kim (admin) checks bookings and sees a conflict or needs to free capacity.

**Rising Action:**

- Reviews item group bookings and today's presence
- Cancels a booking that needs resolution (admin privilege)

**Climax:**
Conflict is resolved immediately.

**Resolution:**
Kim says, "This was fast and easy."

---

### Journey 4: IT / Setup -- "Low-Friction Launch"

**Opening Scene:**
Sam from IT needs to deploy SitHub quickly.

**Rising Action (with Entra ID):**

- Creates a new Entra ID application
- Configures server settings and Entra ID connection
- Defines areas, item groups, items, and equipment in the space configuration file
- Starts the single executable

**Rising Action (without Entra ID):**

- Configures server settings (no Entra ID section needed)
- Imports demo users via the provided SQL file or creates users via the API
- Defines areas, item groups, items, and equipment in the space configuration file
- Starts the single executable

**Climax:**
Authentication works and the UI loads with the defined spaces.

**Resolution:**
Setup is complete and low-maintenance (target: under 30 minutes).

---

### Journey 5: Employee -- "Local Login and Password Change"

**Opening Scene:**
Dana works at a company without Entra ID. She uses SitHub with local credentials.

**Rising Action:**

- Opens SitHub and sees the login page with a username/password form
- Enters email and password, logs in
- Books an item as usual
- Later, changes her password via her profile

**Climax:**
"This works just like any other web app -- no special enterprise setup needed."

**Resolution:**
Dana manages her own credentials and uses SitHub without Entra ID dependency.

### Journey Requirements Summary

- Dual authentication: Entra ID SSO or local credentials with email and password
- Login page with form fields and "Login via Entra ID" button
- Self-service password change for local users
- Admin password reset for local users
- Users table storing both Entra ID and local users
- Email uniqueness enforced across authentication sources
- Responsive list-based navigation (area -> item group -> item)
- Single-day and week booking flow with immediate confirmation
- Booking notes (add, view, edit)
- Weekly availability preview for item groups with color-coded weekday indicators
- Booker display name shown on booked items
- "My Bookings" view with cancel action and note editing
- Item group booking overview and "Today's presence" per area
- Conflict handling when an item becomes unavailable
- Role-based permissions: users cancel own bookings; admins can cancel any
- Clickable breadcrumbs for navigation hierarchy
- File-based space configuration
- Single-binary distribution (frontend + backend)
- Demo users SQL file for development and testing setup

## Web App Specific Requirements

### Project-Type Overview

- Single Page Application (SPA)
- Desktop + mobile support across major browsers
- No SEO requirements
- Real-time availability updates
- Accessibility target: WCAG A

### Technical Architecture Considerations

- SPA client with REST backend
- Real-time updates for availability (push or polling)
- Responsive layout for mobile + desktop

### Browser & Accessibility

- Supported browsers: modern Chrome, Edge, Firefox, Safari (desktop + mobile)
- Accessibility compliance: WCAG A

### Implementation Considerations

- Ensure real-time availability sync without conflicts
- Optimize UI for small screens (no pinch/zoom)

## UX/UI Requirements

- Primary flow is list-based and linear: area -> item group -> item -> confirm
- Booking status is visible at-a-glance with clear available/unavailable states
- Weekly availability is visible via color-coded weekday indicators (MO-FR) on item group
  tiles; green indicates at least one item available, red indicates fully booked
- Users can toggle between day and week booking modes; the selected mode persists across
  sessions via browser local storage
- Booking confirmation includes the item name from the configuration (e.g., "Parking Lot 1
  booked successfully") and an option to add a note
- Booking notes are visible in My Bookings, Today's Presence, and item detail views;
  notes longer than the display width are truncated with an expand indicator
- Booked items display the booker's display name
- Breadcrumbs are clickable and navigate to the corresponding hierarchy level
- Error states are plain-language and actionable (e.g., item taken, retry)
- "My Bookings" is reachable from the primary navigation on mobile and desktop
- Layout remains readable and fully operable on small screens without zoom
- UI uses accessible labels, focus states, and contrast consistent with WCAG A
- UI action labels use domain-neutral terminology ("BOOK" instead of "VIEW DESK/ROOM")

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach:** Experience MVP + Platform MVP
**Resource Requirements:**

- 1 backend engineer
- 1 frontend engineer
- Light DevOps/IT support for Entra ID setup + packaging
- Optional UX input for mobile-first flow and accessibility

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**

- Employee happy-path booking (Entra ID or local login)
- Employee edge case (item taken)
- Employee local login and password change
- Admin cancellation
- IT setup and configuration (with or without Entra ID)

**Must-Have Capabilities:**

- Dual authentication (Entra ID SSO or local credentials)
- User management API and password management
- Area -> item group -> item list with equipment
- Single-day booking + week booking + cancel
- Weekly availability preview
- Booking notes (add, view, edit)
- "My Bookings"
- Item group booking overview + today's presence
- Real-time availability updates
- Single executable (frontend + backend)
- Clear setup docs (including Entra ID steps and local auth setup)

### Post-MVP Features

**Phase 2 (Post-MVP):**

- Booking on behalf of others (UI)
- Guest bookings (UI)
- Notifications
- Booking history
- Multi-day and recurring bookings

**Phase 3 (Expansion):**

- Graphical floor maps
- Admin management UI + advanced controls
- Reporting and analytics

### Risk Mitigation Strategy

**Technical Risks:**

- Real-time availability sync + conflict resolution
- Entra ID setup reliability
- Single-binary packaging of frontend + backend
*Mitigation:* build early spikes for live updates and Entra ID, validate packaging early.

**Market Risks:**

- Users may tolerate the Confluence workflow and resist change
*Mitigation:* run the 5-day trial, collect explicit preference feedback.

**Resource Risks:**

- Limited capacity across backend/frontend/ops
*Mitigation:* keep strict MVP boundaries and defer non-essentials.

## Functional Requirements

### Identity & Access

- FR1: Users can authenticate via Entra ID SSO or local credentials (email and password).
  Acceptance: the login page presents a username/password form and a "Login via Entra ID"
  button; both methods result in an authenticated session with the user's name displayed.
- FR2: The system determines user roles (regular vs admin) based on authentication source.
  For Entra ID users, admin status is synced from group membership on every login and cannot
  be changed locally. For local users, admin status is managed locally by administrators.
  Acceptance: admins see admin-only controls; regular users do not; Entra ID admin status
  reflects current group membership after each login.
- FR3: Users can access the application only if they are authenticated. Acceptance:
  unauthenticated users see only the login page and cannot view any booking data.

### User Management

- FR28: All users (Entra ID and local) are stored in a users table. Entra ID users are
  inserted on first login and updated on subsequent logins. Acceptance: after login, the
  user exists in the users table with correct source, name, and email.
- FR29: Email addresses are unique across all authentication sources. Creating a local user
  with an email that exists for an Entra ID user is blocked, and vice versa. Acceptance:
  attempting to create a duplicate email returns an error regardless of authentication source.
- FR30: Local users can log in with email and password. Acceptance: entering valid credentials
  in the login form creates an authenticated session; invalid credentials show a descriptive
  error message.
- FR31: Local users can change their own password via the `/me` endpoint. Acceptance: after
  changing the password, the old password no longer works and the new password (minimum 14
  characters) is accepted.
- FR32: Admin users can reset the password of any local user via the `/users/{id}` endpoint.
  Acceptance: the affected user can log in with the new password; Entra ID user passwords
  cannot be reset this way.
- FR33: The system provides a `/me` endpoint returning the current user's profile information.
  Acceptance: authenticated requests to `/me` return the user's id, email, name, role, and
  authentication source.
- FR34: The system provides a `/users` endpoint for user management. Acceptance: admins can
  list, create, read, update, and delete local users; non-admin users can only read. Entra ID
  users cannot be created or deleted via this endpoint.
- FR35: A demo users SQL file is provided for development setup. Acceptance: running the SQL
  file creates 15 users (2 admins, 13 regular users) with local credentials in the database.

### Areas and Items Discovery

- FR4: Users can view a list of available areas. Acceptance: after login, the UI lists all
  configured areas.
- FR5: Users can view a list of item groups within a selected area. Acceptance: selecting an
  area shows only its item groups.
- FR6: Users can view a list of items within a selected item group. Acceptance: selecting an
  item group lists its items.
- FR7: Users can view equipment details for each item. Acceptance: each item entry shows its
  equipment list if configured.
- FR8: Users can view current booking status for items. Acceptance: item entries show
  available/occupied status for the selected date.
- FR36: Users can view a weekly availability preview for item groups. Acceptance: the item
  group list view includes a calendar week selector (next 8 weeks, current week pre-selected);
  each item group tile displays weekday indicators (MO-FR) colored green (at least one item
  available) or red (fully booked) for the selected week; the week selector displays dates in
  locale-aware format with week number (e.g., "2026-03-16 - Week 12").
- FR39: Users can see the display name of the person who booked an item. Acceptance: in the
  item detail view, booked items show the booker's display name alongside the booking status.

### Booking Creation

- FR9: Users can book an item for a single day. Acceptance: selecting an item and date creates
  a booking that appears in "My Bookings."
- FR10: The system prevents double-booking of the same item for the same day. Acceptance:
  the second attempt is rejected and no duplicate booking is created.
- FR11: Users receive a message when a selected item becomes unavailable during booking.
  Acceptance: the message states the item is no longer available for that date and prompts
  the user to choose another item.
- FR38: Users can toggle between day booking mode and week booking mode. Acceptance: the
  selected mode is persisted in browser local storage and restored on next visit; in week
  mode, the date selector becomes a calendar week selector (next 8 weeks); item tiles show
  per-day checkboxes with booker names; existing bookings by other users cannot be unchecked;
  a single "Confirm My Booking" button submits all selected days at once.

### Booking Management

- FR12: Users can view their upcoming bookings ("My Bookings"). Acceptance: the list includes
  item, item group, area, and date for each future booking; booking notes are displayed if
  present.
- FR13: Users can cancel their own bookings. Acceptance: cancelling removes the booking from
  all relevant lists and frees the item.
- FR14: Admin users can cancel any booking. Acceptance: admins can cancel another user's
  booking and the affected user sees the cancellation reflected in their list.
- FR37: Users can add, view, and edit free-text notes on their bookings. Acceptance: after
  completing a booking, a confirmation message includes an "add note" action; notes are
  visible in My Bookings, Today's Presence, and item detail views; notes longer than the
  display width are truncated with an expand indicator that opens the full text; notes are
  editable from the My Bookings view.

### Item Group and Presence Overviews

- FR15: Users can view an item-group-level booking overview. Acceptance: for a selected item
  group and date, the overview lists all booked items and associated users; booking notes are
  displayed if present.
- FR16: Users can view "Today's presence" for an area (who is in the office today).
  Acceptance: the view lists all users with bookings in that area for today; booking notes
  are displayed if present.

### Configuration & Setup (Operator Capabilities)

- FR17: Operators can configure server settings via a configuration file. Acceptance: invalid
  settings prevent startup with a descriptive error message stating the failure reason; valid
  settings allow startup.
- FR18: Operators can configure areas, item groups, items, and equipment via a configuration
  file. Acceptance: after restart, the UI reflects the updated space definitions.
- FR19: The system can load and apply configuration changes on startup. Acceptance:
  configuration changes take effect after restart without manual data migration steps.

### Navigation & UI Consistency

- FR40: Breadcrumbs in the navigation hierarchy are clickable and navigate to the
  corresponding view. Acceptance: clicking any breadcrumb segment navigates to that level of
  the hierarchy; the current page breadcrumb is not clickable.
- FR42: UI action labels use domain-neutral terminology. Acceptance: labels read "BOOK"
  instead of "VIEW DESK" or "VIEW ROOM"; "BOOK THIS ITEM" instead of "BOOK THIS DESK";
  booking confirmation messages reference the item name from the configuration (e.g.,
  "Parking Lot 1 booked successfully") rather than the generic term "desk."
- FR43: Navigation action labels are simplified for speed. Acceptance: "VIEW ITEM GROUPS"
  becomes "SELECT" on area tiles; "VIEW ITEMS" becomes "SELECT" on item group tiles;
  redundant page titles and subtitles are removed from item group and item views; "BOOK
  THIS ITEM" becomes "BOOK" on available item tiles in day mode.
- FR44: Redundant "Not available" text is removed from booked items in day mode. Acceptance:
  booked items show only the "booked" status chip, the booker's name, and any booking notes;
  the "Not available for \<date\>" message no longer appears.
- FR45: Booker name and booking notes on occupied items use a readable font size.
  Acceptance: both are displayed at body-2 size or larger (not caption size); information is
  legible without zooming on desktop and mobile.
- FR46: Booking result feedback uses icons. Acceptance: successful bookings show a green
  checkmark icon; failed bookings show a red warning icon; the item name and day label
  accompany each icon; raw error text is replaced by concise icon-based feedback.

### Booking Form Simplification

- FR47: Guest booking option is removed from the booking form. Acceptance: the "Book for
  guest" radio button and associated guest name/email fields no longer appear.
- FR48: Multi-day booking checkbox is removed from day booking mode. Acceptance: the "Book
  multiple days" checkbox and associated additional dates field no longer appear; week
  booking mode is the replacement for multi-day booking.
- FR49: Colleague booking uses a user dropdown instead of free-text input. Acceptance: when
  "Book for colleague" is selected, a dropdown lists existing users by display name from the
  database; booking is submitted with the selected user's ID; free-text email and name
  fields are removed.

### Collapsible Item Tiles

- FR50: Item tiles in week booking mode are collapsible. Acceptance: each tile has a chevron
  toggle; folded state shows the compact M-F row with checkboxes and truncated names;
  unfolded state shows one line per day with full day names, full booker names, equipment,
  and warnings; the chevron rotates between left (folded) and down (unfolded).
- FR51: Item tiles in day booking mode are collapsible for booked items. Acceptance: booked
  item tiles hide equipment and warnings by default; a chevron toggle reveals the full
  details; available item tiles remain fully expanded (no collapsing).
- FR52: Folded tiles with warnings show a warning icon. Acceptance: when a tile is folded and
  the item has a warning, a warning icon appears in the tile header; clicking the icon
  displays the warning message in a tooltip or popup; the icon is visible in both day and
  week modes.
- FR53: Past date checkboxes are disabled in week booking mode. Acceptance: checkboxes for
  dates before today are disabled and visually grayed out; users cannot select past dates for
  booking; only future and today's dates have active checkboxes.
- FR54: Full booker name appears on hover in week mode. Acceptance: in the folded tile view,
  hovering over a truncated booker name shows the full display name in a tooltip.

### User Preferences

- FR55: Users can select a visual theme. Acceptance: the user menu offers auto (default),
  dark, and light theme options; the selection is persisted in localStorage; the Vuetify
  theme updates immediately on selection; auto follows the OS preference.
- FR56: Users can toggle weekend visibility. Acceptance: the user menu includes a "Show
  weekends" checkbox; when enabled, all booking pages (day and week mode) show Saturday and
  Sunday; the preference is persisted in localStorage; the default is weekends hidden.
- FR57: The Change Password menu item displays its icon. Acceptance: the icon is visible next
  to "Change Password" in both desktop and mobile navigation menus.

### Equipment Filtering

- FR58: Users can filter items by equipment keywords. Acceptance: a text input labeled
  "Filter equipment" appears on the booking page; items not matching the filter are blurred
  with an "equipment not available" overlay; an info icon explains the search syntax;
  keywords are combined with OR by default; the plus sign combines with AND; single or
  double quotes trigger exact matching; filtering is case-insensitive; the filter works in
  both day and week booking modes; filtering is implemented on the frontend only.

### Data Integrity

- FR41: The bookings table references users by `user_id` only; display names are resolved
  via JOIN with the users table. Acceptance: the bookings table does not store denormalized
  user name columns; all booking queries that require user names perform a JOIN; existing
  bookings are migrated to remove redundant columns.

### Post-MVP (Phase 2+)

These requirements extend the MVP booking experience for power users and team coordinators.
They trace to the Growth Features scope and would require new user journeys when implemented.

- FR20: Users can book on behalf of another user. Acceptance: the booking appears in both
  users' booking lists and either can cancel.
- FR21: Users can book items for guests outside the organization. Acceptance: a guest booking
  stores guest name and contact and is visible as a guest booking in overviews.
- FR22: Users can book multi-day or recurring reservations. Acceptance: the system creates
  individual daily bookings and reports conflicts per day.
- FR23: Users can view booking history. Acceptance: users can see past bookings with date
  range filtering.
- FR24: Users can receive notifications related to bookings. Acceptance: booking
  creation/cancellation triggers a notification via the configured channel within 5 minutes.
- FR25: Admins can manage item groups and items via an admin UI. Acceptance: admins can
  add/edit/remove item groups and items; changes appear in discovery lists after save.

### Future (Phase 3+)

These requirements support the long-term vision of visual space management and data-driven
decisions. They trace to the Vision scope and would require new user journeys and potentially
new user personas when implemented.

- FR26: Users can book items using a graphical floor-map view. Acceptance: an item selected
  on the map can be booked for a chosen date.
- FR27: Admins can access advanced reporting and analytics. Acceptance: admins can view usage
  summaries by area/item group and date range.

## Non-Functional Requirements

### Performance

- For expected usage (<=50 concurrent users), list navigation actions complete within
  2 seconds at p95; booking and cancellation complete within 3 seconds at p95.

### Reliability

- The system can be restarted without data loss; bookings remain intact after restart and
  conflicts do not create partial records.

### Security

- All booking data requires authenticated access (Entra ID or local credentials);
  unauthenticated requests are denied.
- Local user passwords are stored as cryptographic hashes; plaintext passwords are never
  persisted or logged.
- Local user passwords require a minimum length of 14 characters.
- Data at rest is stored without application-layer encryption; in-transit encryption is
  managed outside the application.

### Scalability

- Single-node deployment is sufficient; no clustering or horizontal scaling is required for
  MVP usage levels.

### Accessibility

- Meets WCAG A: all interactive elements have accessible names, keyboard focus is visible,
  and form inputs are labeled.

### Availability

- For MVP, no formal uptime SLA is required; the system is recoverable by operator restart
  within minutes of failure.
