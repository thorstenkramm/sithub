# Architecture

## Executive Summary
SitHub is a desk booking web application with a Go (Echo) REST backend and a Vue 3 + Vuetify frontend. The system is distributed as a single binary with embedded frontend assets and uses an embedded SQLite3 database. This document is based on a quick scan (pattern-based) and README details; source code was not present in this repository snapshot.

## System Context
- **Users:** Office employees and admins managing shared desks and rooms.
- **Core Capabilities:** Desk booking, real-time availability, notifications, admin operations, and SSO via Entra ID.
- **Deployment:** Single-binary distribution, reverse proxy for SSL termination.

## Technology Stack
- **Backend:** Go, Echo framework
- **Frontend:** Vue 3, Vuetify
- **API:** REST (JSON:API)
- **Database:** SQLite3 (embedded)
- **Auth:** Entra ID SSO with group-based access
- **CI/CD:** GitHub Actions (mentioned in README)

## Architecture Pattern
- **Pattern:** SPA frontend + REST API backend.
- **Likely layering:** HTTP handlers/controllers -> services -> data access (inferred).
- **Frontend integration:** Static assets embedded into the backend binary.

## Data Architecture
- **Storage:** Single embedded SQLite database.
- **Expected entities (from product domain):** Areas, Rooms, Desks, Desk Equipment, Bookings, Users, Groups, Guests.
- **Schema details:** Not available in this repo snapshot.

## API Design
- **Style:** JSON:API compliant REST endpoints.
- **Auth:** Entra ID SSO; group-based authorization.
- **Endpoint inventory:** Not available (no route files detected).

## Component Overview
- **Frontend components:** Not detected in this repository snapshot.
- **Admin UI:** Likely part of the Vue app (from README).

## Source Tree Highlights
- **Observed:** Config files and examples (`sithub.example.toml`, `sithub_areas.example.yaml`, `sithub_areas.schema.json`).
- **Missing:** Source directories (`src/`, `app/`, `client/`, `server/`) and build manifests.

## Development Workflow
- **Build/Test:** Not detected (no `go.mod`, `package.json`, or build scripts).
- **Configuration:** Example TOML and YAML files are present.

## Deployment Architecture
- **Artifact:** Single binary with embedded frontend assets.
- **Runtime:** Reverse proxy for SSL termination (per README).
- **Database:** Local SQLite file.

## Testing Strategy
- **Evidence:** No test files detected in this repo snapshot.

## Gaps and Next Steps
- Provide source directories and build manifests to document:
  - API routes and request/response schemas
  - Database schema and migrations
  - Frontend component inventory and state management
  - Testing strategy and CI/CD workflows
