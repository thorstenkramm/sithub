---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/prd.md
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/product-brief-sithub-2026-01-17.md
  - /Users/thorsten/projects/thorsten/sithub/docs/index.md
  - /Users/thorsten/projects/thorsten/sithub/docs/project-overview.md
workflowType: 'architecture'
lastStep: 8
status: 'complete'
completedAt: '2026-01-17'
project_name: 'sithub'
user_name: 'Thorsten'
date: '2026-01-17'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each
architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
The system must support Entra ID SSO, role-based permissions (regular vs admin), discovery of areas/rooms/desks with
equipment and status, single-day booking with conflict handling, booking management (view/cancel), room and presence
overviews, and configuration-driven setup. Post-MVP capabilities include booking on behalf, guests, recurring bookings,
history, notifications, admin management UI, and later floor-map booking and analytics.

**Non-Functional Requirements:**

- Performance: for <=50 concurrent users, list navigation p95 <= 2s; booking/cancel p95 <= 3s.
- Reliability: restart without data loss; no partial records on conflicts.
- Security: all booking data behind Entra ID auth; unauthenticated access denied.
- Scalability: single-node deployment sufficient for MVP.
- Accessibility: WCAG A (labels, keyboard focus, labeled inputs).

**Scale & Complexity:**

- Primary domain: web application (SPA + API).
- Complexity level: low to medium (real-time availability + SSO + config-driven space model).
- Estimated architectural components: 5-7 (auth, booking, availability, config, UI, presence/overview, admin).

### Technical Constraints & Dependencies

- Entra ID SSO integration is mandatory.
- Space definition is file-based configuration.
- Single-binary distribution is a technical success criterion.
- Existing docs suggest a monolith with SPA + REST; treat as a candidate constraint to validate.

### Cross-Cutting Concerns Identified

- Real-time availability consistency and conflict handling
- Authorization boundaries for admin vs regular users
- Configuration validation and safe startup behavior
- Accessibility for mobile and desktop

## Starter Template Evaluation

### Technical Preferences Confirmed

- Backend: Go + Echo
- Frontend: Vue 3 + Vuetify
- Data store: SQLite (embedded)
- Packaging: single-binary distribution required
- Real-time: required (final mechanism TBD)
- Release targets: macOS (arm) and Linux (amd64) only; Windows out of scope

### Starter Options Considered (Based on Preferences)

#### Option A: Minimal Go + Echo API skeleton + separate Vue 3 + Vuetify SPA

- Establishes clear API/UI separation while preserving single-repo monolith.
- Keeps configuration and deployment flexible for single-binary packaging.

#### Option B: Go embedded assets (SPA bundled into binary) + Echo API

- Aligns tightly with the single-binary distribution requirement.
- Encourages early decisions on asset embedding, routing, and build pipeline.

#### Option C: Full-stack template with integrated build (Go + SPA + embed)

- Highest alignment with packaging goal, lowest setup effort.
- Risks coupling build tooling early; would need validation of maintenance quality.

### Implications of Each Option

- Option A optimizes developer workflow but may require extra packaging work to meet single-binary constraints.
- Option B bakes in the deployment requirement from day one, minimizing later refactors.
- Option C is fastest to start but depends on template quality and current maintenance; verify before adoption.

### Deployment/Packaging Constraints

- Build pipeline must produce native binaries for macOS (arm64) and Linux (amd64).
- Windows build tooling and packaging are intentionally excluded.

### Recommendation (Pending Version Verification)

Given the single-binary constraint, Option B is the safest architectural foundation. If a high-quality template
exists for Option C, it could be adopted after verifying active maintenance and toolchain compatibility.

**Note:** Template names and versions must be verified externally due to lack of web access in this environment.

## Core Architectural Decisions

### Category 1: Data Architecture

- Decision: Simple relational schema (SQLite with WAL).
- Decision: Versioned migrations using golang-migrate.
- Decision: No caching for MVP (direct reads from SQLite).
- Rationale: Keeps data model and operations simple, aligns with single-node deployment, and avoids premature
  optimization.

### Category 2: Real-Time Updates

- Decision: WebSockets for live availability updates with polling fallback.
- Rationale: Ensures responsive UX while providing resilience in environments that block WebSockets.

## Implementation Patterns & Consistency Rules

### Naming Patterns

- Database tables: snake_case, plural (e.g., `bookings`, `desks`).
- Database columns: snake_case; foreign keys use `{entity}_id`.
- Index naming: `idx_{table}_{column}`.
- API endpoints: plural nouns under `/api/v1`, no trailing slashes.
- Route params: `:id` format (Echo conventions).
- JSON fields: snake_case (JSON:API rule).
- Dates: RFC 3339 UTC; fields end with `_at`; ranges use `_from` / `_until`.

### Structure Patterns

- Backend organization: `internal/{domain}` for handlers/services/repos; shared JSON:API response types in
  `internal/api`.
- Frontend organization: feature-based folders; Composition API with `<script setup>`; Pinia stores in `src/stores`.
- Tests: backend tests co-located; frontend uses Vitest (unit) and Cypress (E2E).
- API docs: OpenAPI 3.1 in `api-doc/` with `openapi.yaml` as entry point and per-endpoint files.

### Format Patterns

- API responses follow JSON:API envelopes with `application/vnd.api+json` content type.
- Error responses use JSON:API `errors[]` with `status`, `title`, `detail`.
- Resource `type` values use kebab-case; attributes use snake_case.
- Pagination uses offset-based params: `page[limit]` and `page[offset]`.
- Trailing slashes are treated as distinct URLs; collection routes use `/resource`, not `/resource/`.

### Process Patterns

- Auth: Entra ID required for all API endpoints; unauthenticated requests return JSON:API errors (401/403).
- Config: TOML keys in snake_case with documented defaults per `toml.md`.
- Logging: structured logging with `log/slog` and error wrapping using `%w`.
- Tooling: Go 1.25 and golangci-lint 2.5.0; Node.js 24 with create-vue and Vite.
- Dev workflow: Vite proxies `/api` to `http://localhost:9900`; run Go and Vite in separate terminals.
- Containers: no Docker/K8s; use native `go` toolchain and local Node.js installs.
- Testing: acceptance criteria require Cypress E2E; use `data-cy` selectors and avoid mocked API in E2E.
- MCPs: Vuetify MCP and Chrome DevTools MCP are available for UI work and debugging.
- Documentation: keep Markdown lines <= 120 chars and lint with `markdownlint --fix`.

### Go Code Rules

- Use cobra for CLI entrypoints and viper for configuration loading.
- Use sentinel errors for domain conditions and wrap errors with `%w`.
- Document thread safety on types and security-sensitive helpers.
- Keep shared JSON:API response types in `internal/api` to avoid duplication.

### Frontend Rules

- Use Vue 3 with TypeScript, Pinia, Vue Router, and Vuetify.
- Use Composition API with `<script setup>`; no Options API or mixins.
- Use Vitest for unit tests and Cypress for E2E.

### API Documentation Rules

- OpenAPI 3.1 lives in `api-doc/openapi.yaml` with per-endpoint files.
- Lint API docs with `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.

## Project Structure & Boundaries

### Requirements to Modules Mapping

- Identity & Access -> `internal/auth`, `internal/middleware`
- Discovery (areas/rooms/desks) -> `internal/areas`, `internal/rooms`, `internal/desks`
- Booking creation/management -> `internal/bookings`
- Room/Presence overviews -> `internal/rooms`, `internal/presence`
- Config/setup -> `internal/config`, `internal/startup`
- JSON:API types -> `internal/api`

### Proposed Repository Structure

```text
sithub/
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── .gitignore
├── .golangci.yml
├── .github/
│   └── workflows/
│       └── ci.yml
├── api-doc/
│   ├── openapi.yaml
│   └── endpoints/
│       └── bookings.yaml
├── cmd/
│   └── sithub/
│       └── main.go            # cobra entrypoint
├── internal/
│   ├── api/                   # JSON:API response types + error builders
│   ├── auth/                  # Entra ID auth flow + role extraction
│   ├── middleware/            # auth/role guards, request logging
│   ├── config/                # viper config loading/validation
│   ├── db/                    # SQLite init + migrations runner
│   ├── areas/                 # area handlers/services/repos
│   ├── rooms/                 # room handlers/services/repos
│   ├── desks/                 # desk handlers/services/repos
│   ├── bookings/              # booking handlers/services/repos
│   ├── presence/              # today's presence logic
│   ├── system/                # health, version, uptime
│   └── startup/               # app wiring + router registration
├── migrations/                # golang-migrate SQL files
├── assets/
│   └── web/                   # embedded SPA build output (dist)
├── web/                       # Vue app root
│   ├── package.json
│   ├── package-lock.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── cypress.config.ts
│   ├── cypress/
│   │   ├── e2e/
│   │   ├── component/
│   │   ├── fixtures/
│   │   └── support/
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── router/
│   │   ├── stores/            # Pinia
│   │   ├── views/
│   │   ├── components/
│   │   ├── composables/
│   │   └── api/               # JSON:API client layer
│   └── tests/                 # Vitest unit tests
└── tools/
    └── embed/                 # build helpers to copy web/dist -> assets/web
```

### Integration Boundaries

- API boundary: `/api/v1/*` JSON:API only, no trailing slashes.
- Auth boundary: middleware enforces Entra ID on all API routes; role checks inside handlers.
- Data boundary: repositories in each domain; no direct SQL in handlers.
- Frontend boundary: API access only through `web/src/api` and composables; no raw fetches in views.

## Architecture Validation Results

### Coherence Validation

**Decision Compatibility:**
Go + Echo + SQLite (WAL) + Vue 3/Vuetify align with the single-binary, single-node goals. JSON:API rules and OpenAPI
3.1 documentation are compatible. Tooling rules (Go 1.25, Node 24, no Docker) fit the structure.

**Pattern Consistency:**
Naming, date/time, and pagination standards align with JSON:API rules. Go code rules align with shared types in
`internal/api`. Vue rules align with Composition API, Pinia, and Vuetify.

**Structure Alignment:**
The repo tree supports domain modules, OpenAPI docs, embedded SPA assets, and test layout.

### Requirements Coverage Validation

**Functional Coverage:**

- Identity & Access -> `internal/auth`, `internal/middleware`
- Discovery -> `internal/areas`, `internal/rooms`, `internal/desks`
- Booking -> `internal/bookings`
- Room/Presence -> `internal/rooms`, `internal/presence`
- Config/Setup -> `internal/config`, `internal/startup`

**Non-Functional Coverage:**

- Performance targets supported by simple relational model and single-node scope.
- Reliability supported by migrations and SQLite WAL.
- Security covered by Entra ID enforcement and JSON:API error handling.
- Accessibility supported by frontend rules and patterns.

### Implementation Readiness Validation

Critical decisions, patterns, and structure are explicit enough for multi-agent implementation.

### Gaps & Notes

**Important:**

- Technology versions in `.claude/rules/techstack.md` are not web-verified here.

**Nice-to-Have:**

- Define booking conflict resolution semantics in the data layer (optimistic vs pessimistic locking).

## Booking Conflict Resolution

- Decision: Use optimistic conflict handling with a unique constraint on (desk_id, booking_date).
- Behavior: If a second booking races, the write fails and returns a JSON:API conflict error.
- Rationale: Prevents double-booking with minimal complexity while keeping the system responsive.

## Booking Time Granularity

- Decision: Bookings are full-day only (no half-day or hourly bookings).
- Behavior: Booking records store a single booking_date per reservation.

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED
**Total Steps Completed:** 8
**Date Completed:** 2026-01-17
**Document Location:** `_bmad-output/planning-artifacts/architecture.md`

### Final Architecture Deliverables

#### Complete Architecture Document

- Architectural decisions documented and aligned with `.claude/rules`
- Implementation patterns ensuring multi-agent consistency
- Project structure with explicit boundaries and ownership
- Requirements to modules mapping
- Validation confirming coherence and coverage

#### Implementation Ready Foundation

- 4 key architectural decisions made
- 7 pattern categories defined
- 8 major architectural components specified
- 6 functional requirement categories supported

#### AI Agent Implementation Guide

- Technology stack with specified versions (pending web verification)
- Consistency rules to prevent implementation conflicts
- Project structure with clear boundaries
- Integration patterns and communication standards

### Implementation Handoff

**For AI Agents:**
Follow this architecture document before implementing any story. Keep all rules and boundaries intact.

**First Implementation Priority:**
Initialize the Go CLI entrypoint, Vue app scaffold, and asset embedding pipeline per the structure and rules.

**Development Sequence:**

1. Initialize project structure and config loading.
2. Set up database initialization and migrations.
3. Implement auth middleware and base API responses.
4. Build core domain modules (areas, rooms, desks, bookings, presence).
5. Implement frontend views and API client layer.
6. Add tests per Vitest and Cypress rules.
