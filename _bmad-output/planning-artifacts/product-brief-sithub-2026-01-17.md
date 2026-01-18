---
stepsCompleted: [1, 2, 3, 4, 5]
inputDocuments:
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
date: 2026-01-17
author: Thorsten
---

# Product Brief: sithub

<!-- Content will be appended sequentially through collaborative workflow steps -->


## Executive Summary

SitHub is a desk booking web application for shared offices that replaces manual, error-prone Confluence tables with a fast, mobile-friendly booking experience. Employees can see real-time availability, book desks in seconds, and receive notifications, while admins manage areas, rooms, desks, and access controls. SitHub prioritizes simplicity: single-binary deployment, embedded SQLite database, and Entra ID SSO deliver low operational overhead without sacrificing usability.

---

## Core Vision

### Problem Statement

Shared offices rely on manual desk booking workflows—often Confluence tables that must be recreated weekly. This is time-consuming, error-prone, and delivers a poor user experience, especially on mobile.

### Problem Impact

- Employees waste time finding available desks and updating manual tables.
- Admins spend effort creating and maintaining weekly booking sheets.
- Errors and conflicts are common, leading to frustration and inefficiency.
- The experience is clunky on mobile even when apps exist.

### Why Existing Solutions Fall Short

- Current approaches are manual and fragile.
- Enterprise desk-booking tools are often expensive and overbuilt for smaller teams.
- Many tools don’t prioritize low-maintenance deployment or frictionless setup.

### Proposed Solution

A user-friendly desk booking web app with real-time availability, quick booking (single day to multi-day), and notifications. Admins manage areas, rooms, desks, and access via a YAML configuration. SitHub supports Entra ID SSO, runs as a single binary, and uses an embedded SQLite database—making setup and maintenance minimal.

### Key Differentiators

- **Single-binary deployment + embedded database** → low ops and cost
- **Entra ID SSO with group-based access** → enterprise-ready authentication
- **Mobile-first, user-friendly UX** → designed for fast bookings
- **Simple configuration via YAML** → admin control without heavy tooling


## Target Users

### Primary Users

**Employees (Desk Bookers)**  
- **Context:** Company employees who need desks in shared offices and work across mobile and desktop.  
- **Goals:** Find available desks quickly, book in seconds, sit near colleagues when needed.  
- **Pain Today:** Manual Confluence tables are tedious and error-prone; mobile UX is frustrating.  
- **Success Moment:** “I can book a desk (or two) in seconds and I’m done—no tables, no confusion.”

Key needs:
- Mobile-optimized and desktop-friendly booking
- Real-time availability visibility
- Easy booking for themselves or on behalf of another user

### Secondary Users

**Admins (Office/Operations)**  
- Manage bookings and resolve conflicts  
- Can cancel any booking  
- Ensure fairness and smooth desk allocation

**IT / Platform Owners**  
- Configure system via `sithub.example.toml` and `sithub_areas.example.yaml`  
- Set up Entra ID SSO and group-based access  
- Prefer low-maintenance, file-based configuration (no admin UI required)

### User Journey

**Employee Journey**  
- **Discovery:** Access SitHub via company link/SSO.  
- **Onboarding:** Entra ID login; sees area/room/desk layout.  
- **Core Usage:**  
  - Browse availability (all bookings visible).  
  - Book a desk for self or another user (e.g., to sit together).  
- **Success Moment:** “Booked in seconds on my phone.”  
- **Ongoing:** Manage/cancel own bookings; can cancel bookings they made on behalf of others.

**Admin Journey**  
- **Discovery:** Granted admin rights through Entra ID group.  
- **Core Usage:**  
  - Oversee bookings  
  - Cancel any booking when needed  
- **Success Moment:** Fewer conflicts, quick resolution without manual cleanup.

**IT Journey**  
- **Setup:** Configure server + Entra ID credentials (`sithub.example.toml`).  
- **Space Definition:** Define areas/rooms/desks and equipment via YAML + schema.  
- **Success Moment:** System runs with minimal maintenance and no UI dependency.


## Success Metrics

**User Success Metrics**
- Users can complete a desk booking quickly and easily (qualitative feedback: "fast and easy").
- "Wow" reaction: users explicitly describe SitHub as "worlds better than Confluence."

### Business Objectives
- After a 5-day test period, users do not want to return to the Confluence process.

### Key Performance Indicators
- Post-trial preference: >= 80% of test users state they prefer SitHub over Confluence.
- Reversion rate: 0 teams request return to Confluence after the test period.


## MVP Scope

### Core Features

- Entra ID login
- List view of areas, rooms, and desks with desk equipment details
- Book a desk for a single day
- Cancel a booking
- Room booking overview (per room)
- “Today’s presence” view showing who is in the office (by area)

### Out of Scope for MVP

- Graphical floor map view
- Multi-day or recurring bookings (weeks/series)
- Booking history
- Admin features
- Booking on behalf of others
- Guest bookings (non-organization users)
- Notifications

### MVP Success Criteria

- Users can complete a single-day booking quickly and easily.
- After a 5-day test period, users prefer SitHub over the Confluence table process.
- “Wow” reaction that SitHub is “worlds better” than the current workflow.

### Future Vision

- Graphical floor maps
- Multi-day and recurring bookings
- Booking history and reporting
- Admin management UI and advanced controls
- Booking on behalf of others and for guests
- Notifications and reminders
