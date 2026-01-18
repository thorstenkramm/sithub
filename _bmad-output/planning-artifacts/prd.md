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
workflowType: 'prd'---

# Product Requirements Document - sithub

**Author:** Thorsten
**Date:** 2026-01-17


## Executive Summary

## Differentiators

- Single-executable distribution to minimize operational overhead
- Entra ID SSO with group-based access control
- Mobile-first, no-pinch booking experience
- File-based space configuration for simple administration


SitHub is a SPA desk-booking web app for shared offices that replaces manual Confluence tables with a fast, mobile-friendly booking flow. Employees can see real-time availability, book single-day desks, and manage their bookings; admins can cancel any booking; IT configures Entra ID and space definitions via TOML/YAML. The MVP prioritizes a low-friction UX and simple deployment (single executable), with success measured by user preference over the Confluence workflow after a 5-day trial.


## Success Criteria

### User Success
- Users can book a desk quickly and easily on mobile or desktop.
- Users describe SitHub as “worlds better than Confluence.”
- After a 5-day trial, users do not want to return to the Confluence workflow.

### Business Success
- Post-trial preference >= 80% of test users favor SitHub.
- Reversion rate: 0 teams request a return to Confluence.

### Technical Success
- Entra ID authentication works reliably end-to-end.
- Installation and setup documentation is clear and complete, including all Entra ID steps.
- Distribution is a single executable bundling frontend and backend.

### Measurable Outcomes
- Trial preference rate >= 80%.
- Reversion rate = 0.
- Successful Entra ID setup using documented steps in a test environment.

## Product Scope

### MVP - Minimum Viable Product
- Entra ID login
- List view of areas, rooms, and desks with equipment
- Book a desk for a single day
- Cancel a booking
- Room booking overview (per room)
- “Today’s presence” view (who is in the office by area)

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

### Journey 1: Employee (Happy Path) — “Fast and Easy Booking”
**Opening Scene:**  
Alex, an employee, opens SitHub on a phone or desktop to book a desk for tomorrow.

**Rising Action:**  
- Sees a clean list of areas → selects one  
- Chooses a room → sees desk list with equipment  
- Selects a desk → books for a single day  
- Visits “My Bookings” to confirm

**Climax (Aha Moment):**  
“No pinching or zooming needed—everything is readable and actionable on mobile.”

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

### Journey 4: IT / Setup — “Low-Friction Launch”
**Opening Scene:**  
Sam from IT needs to deploy SitHub quickly.

**Rising Action:**  
- Creates a new Entra ID application  
- Fills all required values in `sithub.example.toml`  
- Defines areas/rooms/desks using `sithub_areas.example.yaml` and validates against `sithub_areas.schema.json`  
- Starts the single executable

**Climax:**  
Authentication works and the UI loads with the defined spaces.

**Resolution:**  
Setup is complete and low-maintenance (target: under 30 minutes).

### Journey Requirements Summary

- Entra ID SSO login with clear error handling
- Responsive list-based navigation (area → room → desk)
- Single-day booking flow with immediate confirmation
- “My Bookings” view with cancel action
- Room booking overview and “Today’s presence” per area
- Conflict handling when a desk becomes unavailable
- Role-based permissions: users cancel own bookings; admins can cancel any
- File-based configuration via TOML + YAML schema
- Single-binary distribution (frontend + backend)


## Web App Specific Requirements

### Project-Type Overview
- Single Page Application (SPA)
- Desktop + mobile support across major browsers
- No SEO requirements
- Real-time availability updates
- Accessibility target: WCAG A

### Technical Architecture Considerations
- SPA client + REST backend (JSON:API)
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
- 1 backend engineer (Go/Echo, REST, SQLite)  
- 1 frontend engineer (Vue 3/Vuetify, responsive UX)  
- Light DevOps/IT support for Entra ID setup + packaging  
- Optional UX input for mobile-first flow and accessibility

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**
- Employee happy-path booking  
- Employee edge case (desk taken)  
- Admin cancellation  
- IT setup and configuration

**Must-Have Capabilities:**
- Entra ID login  
- Area -> room -> desk list with equipment  
- Single-day booking + cancel  
- “My Bookings”  
- Room booking overview + today’s presence  
- Real-time availability updates  
- Single executable (frontend + backend)  
- Clear setup docs (including Entra ID steps)

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

- FR1: Users can authenticate via Entra ID SSO. Acceptance: unauthenticated users are redirected to Entra ID and, after successful login, return to SitHub with their name displayed.
- FR2: The system can determine user roles (regular vs admin) based on Entra ID group membership. Acceptance: admins see admin-only cancellation controls; regular users do not.
- FR3: Users can access the application only if they are permitted by the configured Entra ID settings. Acceptance: unauthorized users receive an access-denied screen and cannot view any booking data.

### Areas, Rooms, and Desks Discovery

- FR4: Users can view a list of available areas. Acceptance: after login, the UI lists all configured areas.
- FR5: Users can view a list of rooms within a selected area. Acceptance: selecting an area shows only its rooms.
- FR6: Users can view a list of desks within a selected room. Acceptance: selecting a room lists its desks.
- FR7: Users can view desk equipment details for each desk. Acceptance: each desk entry shows its equipment list.
- FR8: Users can view current booking status for desks. Acceptance: desk entries show available/occupied status for the selected date.

### Booking Creation

- FR9: Users can book a desk for a single day. Acceptance: selecting a desk and date creates a booking that appears in “My Bookings.”
- FR10: The system prevents double-booking of the same desk for the same day. Acceptance: the second attempt is rejected and no duplicate booking is created.
- FR11: Users receive a message when a selected desk becomes unavailable during booking. Acceptance: the message states the desk is no longer available for that date and prompts the user to choose another desk.

### Booking Management

- FR12: Users can view their upcoming bookings (“My Bookings”). Acceptance: the list includes desk, room, area, and date for each future booking.
- FR13: Users can cancel their own bookings. Acceptance: cancelling removes the booking from all relevant lists and frees the desk.
- FR14: Admin users can cancel any booking. Acceptance: admins can cancel another user’s booking and the affected user sees the cancellation reflected in their list.

### Room and Presence Overviews

- FR15: Users can view a room-level booking overview. Acceptance: for a selected room and date, the overview lists all booked desks and associated users.
- FR16: Users can view “Today’s presence” for an area (who is in the office today). Acceptance: the view lists all users with bookings in that area for today.

### Configuration & Setup (Operator Capabilities)

- FR17: Operators can configure server settings via a configuration file. Acceptance: invalid settings prevent startup with a clear error; valid settings allow startup.
- FR18: Operators can configure areas, rooms, desks, and equipment via a configuration file. Acceptance: after restart, the UI reflects the updated space definitions.
- FR19: The system can load and apply configuration changes on startup. Acceptance: configuration changes take effect after restart without manual data migration steps.

### Post-MVP (Phase 2+)

- FR20: Users can book on behalf of another user. Acceptance: the booking appears in both users’ booking lists and either can cancel.
- FR21: Users can book desks for guests outside the organization. Acceptance: a guest booking stores guest name and contact and is visible as a guest booking in overviews.
- FR22: Users can book multi-day or recurring reservations. Acceptance: the system creates individual daily bookings and reports conflicts per day.
- FR23: Users can view booking history. Acceptance: users can see past bookings with date range filtering.
- FR24: Users can receive notifications related to bookings. Acceptance: booking creation/cancellation triggers a notification via the configured channel within 5 minutes.
- FR25: Admins can manage rooms/desks via an admin UI. Acceptance: admins can add/edit/remove rooms/desks and changes appear in discovery lists after save.

### Future (Phase 3+)

- FR26: Users can book desks using a graphical floor-map view. Acceptance: a desk selected on the map can be booked for a chosen date.
- FR27: Admins can access advanced reporting and analytics. Acceptance: admins can view usage summaries by area/room and date range.


## Non-Functional Requirements

### Performance
- For expected usage (<=50 concurrent users), list navigation actions complete within 2 seconds at p95; booking and cancellation complete within 3 seconds at p95.

### Reliability
- The system can be restarted without data loss; bookings remain intact after restart and conflicts do not create partial records.

### Security
- All booking data requires authenticated access via Entra ID; unauthenticated requests are denied.
- Data at rest is stored without application-layer encryption; in-transit encryption is managed outside the application.

### Scalability
- Single-node deployment is sufficient; no clustering or horizontal scaling is required for MVP usage levels.

### Accessibility
- Meets WCAG A: all interactive elements have accessible names, keyboard focus is visible, and form inputs are labeled.
