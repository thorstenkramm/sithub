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
lastEdited: 2026-02-07
workflowType: 'prd'
editHistory:
  - date: 2026-02-07
    changes: "Added multi-source auth (EntraID + local), user management API, password management. Removed test_auth. Fixed validation findings (implementation leakage, post-MVP traceability)."
---

# Product Requirements Document - sithub

**Author:** Thorsten
**Date:** 2026-01-17

## Executive Summary

SitHub is a SPA desk-booking web app for shared offices that replaces manual Confluence tables
with a fast, mobile-friendly booking flow. Employees can see real-time availability, book
single-day desks, and manage their bookings; admins can cancel any booking; IT configures
space definitions via YAML. SitHub supports dual authentication: Entra ID SSO for enterprise
environments and local credentials for teams without Entra ID. The MVP prioritizes a
low-friction UX and simple deployment (single executable), with success measured by user
preference over the Confluence workflow after a 5-day trial.

## Differentiators

- Single-executable distribution to minimize operational overhead
- Dual authentication: Entra ID SSO or local credentials (works with or without Entra ID)
- Mobile-first, no-pinch booking experience
- File-based space configuration for simple administration

## Success Criteria

### User Success

- Users can book a desk quickly and easily on mobile or desktop.
- Users describe SitHub as “worlds better than Confluence.”
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
- List view of areas, rooms, and desks with equipment
- Book a desk for a single day
- Cancel a booking
- Room booking overview (per room)
- "Today's presence" view (who is in the office by area)

### Out of Scope (MVP)

- A UI to manage internal users is not desired for MVP.
- User self-registration is not supported; local users are created by admins or via SQL.

### Growth Features (Post-MVP)

- Booking on behalf of others
- Guest bookings
- Notifications
- Booking history
- Multi-day and recurring bookings

### Vision (Future)

- Graphical floor maps
- Admin management UI and advanced controls
- Advanced reporting and analytics

## User Journeys

### Journey 1: Employee (Happy Path) — "Fast and Easy Booking"

**Opening Scene:**
Alex, an employee, opens SitHub on a phone or desktop to book a desk for tomorrow.

**Rising Action:**

- Sees the login page with two options: a login form and an "Login via Entra ID" button
- Logs in with their preferred method
- Sees a clean list of areas -> selects one
- Chooses a room -> sees desk list with equipment
- Selects a desk -> books for a single day
- Visits "My Bookings" to confirm

**Climax (Aha Moment):**
"No pinching or zooming needed - everything is readable and actionable on mobile."

**Resolution:**
Alex feels it was fast and easy and trusts the booking is secured.

---

### Journey 2: Employee (Edge Case) — “Desk Taken”

**Opening Scene:**  
Alex selects a desk and is about to book.

**Rising Action:**  
Another user books the same desk moments earlier.

**Climax:**  
The system responds with a clear message: “Sorry. Someone else already picked this desk.”

**Resolution:**  
Alex quickly picks another desk and completes booking.

---

### Journey 3: Admin/Operations — “Resolve Conflicts”

**Opening Scene:**  
Kim (admin) checks bookings and sees a conflict or needs to free capacity.

**Rising Action:**  

- Reviews room bookings  
- Cancels a booking that needs resolution (admin privilege)

**Climax:**  
Conflict is resolved immediately.

**Resolution:**  
Kim says, “This was fast and easy.”

---

### Journey 4: IT / Setup — "Low-Friction Launch"

**Opening Scene:**
Sam from IT needs to deploy SitHub quickly.

**Rising Action (with Entra ID):**

- Creates a new Entra ID application
- Configures server settings and Entra ID connection
- Defines areas, rooms, desks, and equipment in the space configuration file
- Starts the single executable

**Rising Action (without Entra ID):**

- Configures server settings (no Entra ID section needed)
- Imports demo users via the provided SQL file or creates users via the API
- Defines areas, rooms, desks, and equipment in the space configuration file
- Starts the single executable

**Climax:**
Authentication works and the UI loads with the defined spaces.

**Resolution:**
Setup is complete and low-maintenance (target: under 30 minutes).

---

### Journey 5: Employee — "Local Login and Password Change"

**Opening Scene:**
Dana works at a company without Entra ID. She uses SitHub with local credentials.

**Rising Action:**

- Opens SitHub and sees the login page with a username/password form
- Enters email and password, logs in
- Books a desk as usual
- Later, changes her password via her profile

**Climax:**
"This works just like any other web app - no special enterprise setup needed."

**Resolution:**
Dana manages her own credentials and uses SitHub without Entra ID dependency.

### Journey Requirements Summary

- Dual authentication: Entra ID SSO or local credentials with email and password
- Login page with form fields and "Login via Entra ID" button
- Self-service password change for local users
- Admin password reset for local users
- Users table storing both Entra ID and local users
- Email uniqueness enforced across authentication sources
- Responsive list-based navigation (area -> room -> desk)
- Single-day booking flow with immediate confirmation
- "My Bookings" view with cancel action
- Room booking overview and "Today's presence" per area
- Conflict handling when a desk becomes unavailable
- Role-based permissions: users cancel own bookings; admins can cancel any
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

- Primary flow is list-based and linear: area → room → desk → confirm
- Booking status is visible at-a-glance with clear available/unavailable states
- Booking and cancellation confirmations are immediate and explicit
- Error states are plain-language and actionable (e.g., desk taken, retry)
- “My Bookings” is reachable from the primary navigation on mobile and desktop
- Layout remains readable and fully operable on small screens without zoom
- UI uses accessible labels, focus states, and contrast consistent with WCAG A

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
- Employee edge case (desk taken)
- Employee local login and password change
- Admin cancellation
- IT setup and configuration (with or without Entra ID)

**Must-Have Capabilities:**

- Dual authentication (Entra ID SSO or local credentials)
- User management API and password management
- Area -> room -> desk list with equipment
- Single-day booking + cancel
- "My Bookings"
- Room booking overview + today's presence
- Real-time availability updates
- Single executable (frontend + backend)
- Clear setup docs (including Entra ID steps and local auth setup)

### Post-MVP Features

**Phase 2 (Post-MVP):**

- Booking on behalf of others  
- Guest bookings  
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
  be changed locally. For local users, admin status is managed via the database.
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
  in the login form creates an authenticated session; invalid credentials show a clear error.
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

### Areas, Rooms, and Desks Discovery

- FR4: Users can view a list of available areas. Acceptance: after login, the UI lists all configured areas.
- FR5: Users can view a list of rooms within a selected area. Acceptance: selecting an area shows only its rooms.
- FR6: Users can view a list of desks within a selected room. Acceptance: selecting a room lists its desks.
- FR7: Users can view desk equipment details for each desk. Acceptance: each desk entry shows its equipment list.
- FR8: Users can view current booking status for desks. Acceptance: desk entries show
  available/occupied status for the selected date.

### Booking Creation

- FR9: Users can book a desk for a single day. Acceptance: selecting a desk and date creates
  a booking that appears in "My Bookings."
- FR10: The system prevents double-booking of the same desk for the same day. Acceptance:
  the second attempt is rejected and no duplicate booking is created.
- FR11: Users receive a message when a selected desk becomes unavailable during booking.
  Acceptance: the message states the desk is no longer available for that date and prompts
  the user to choose another desk.

### Booking Management

- FR12: Users can view their upcoming bookings ("My Bookings"). Acceptance: the list includes
  desk, room, area, and date for each future booking.
- FR13: Users can cancel their own bookings. Acceptance: cancelling removes the booking from
  all relevant lists and frees the desk.
- FR14: Admin users can cancel any booking. Acceptance: admins can cancel another user's
  booking and the affected user sees the cancellation reflected in their list.

### Room and Presence Overviews

- FR15: Users can view a room-level booking overview. Acceptance: for a selected room and
  date, the overview lists all booked desks and associated users.
- FR16: Users can view "Today's presence" for an area (who is in the office today).
  Acceptance: the view lists all users with bookings in that area for today.

### Configuration & Setup (Operator Capabilities)

- FR17: Operators can configure server settings via a configuration file. Acceptance: invalid
  settings prevent startup with a clear error; valid settings allow startup.
- FR18: Operators can configure areas, rooms, desks, and equipment via a configuration file.
  Acceptance: after restart, the UI reflects the updated space definitions.
- FR19: The system can load and apply configuration changes on startup. Acceptance:
  configuration changes take effect after restart without manual data migration steps.

### Post-MVP (Phase 2+)

These requirements extend the MVP booking experience for power users and team coordinators.
They trace to the Growth Features scope and would require new user journeys when implemented.

- FR20: Users can book on behalf of another user. Acceptance: the booking appears in both
  users' booking lists and either can cancel.
- FR21: Users can book desks for guests outside the organization. Acceptance: a guest booking
  stores guest name and contact and is visible as a guest booking in overviews.
- FR22: Users can book multi-day or recurring reservations. Acceptance: the system creates
  individual daily bookings and reports conflicts per day.
- FR23: Users can view booking history. Acceptance: users can see past bookings with date
  range filtering.
- FR24: Users can receive notifications related to bookings. Acceptance: booking
  creation/cancellation triggers a notification via the configured channel within 5 minutes.
- FR25: Admins can manage rooms/desks via an admin UI. Acceptance: admins can add/edit/remove
  rooms/desks and changes appear in discovery lists after save.

### Future (Phase 3+)

These requirements support the long-term vision of visual space management and data-driven
decisions. They trace to the Vision scope and would require new user journeys and potentially
new user personas when implemented.

- FR26: Users can book desks using a graphical floor-map view. Acceptance: a desk selected
  on the map can be booked for a chosen date.
- FR27: Admins can access advanced reporting and analytics. Acceptance: admins can view usage
  summaries by area/room and date range.

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
