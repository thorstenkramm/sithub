# Story 11.8: Demo Users SQL File

Status: done

## Story

As a developer,
I want a demo users SQL file,
So that I can quickly set up a development environment with test data.

## Acceptance Criteria

1. **Given** the SQL file at `tools/database/demo-users.sql` exists
   **When** it is executed against the database
   **Then** 15 users are created: 2 admins and 13 regular users with local credentials
   **And** all passwords are bcrypt-hashed

## Tasks / Subtasks

- [x] Create `tools/database/demo-users.sql` (AC: 1)
  - [x] 2 admin users with is_admin = true
  - [x] 13 regular users with is_admin = false
  - [x] All passwords set to `SitHubDemo2026!!` (bcrypt hashed)
  - [x] All users have source "internal"
- [x] Update CI to seed demo users for E2E tests (AC: 1)
- [x] Update `cy.login()` custom command to use demo user credentials (AC: 1)

## Dev Notes

### Design Decisions

- All demo users share the same password for simplicity in development
- Password `SitHubDemo2026!!` meets the 14-character minimum requirement
- CI seeds demo users via `sqlite3` command before running Cypress tests
- `cy.login()` does `POST /api/v1/auth/login` programmatically

### References

- PRD FR35: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.8: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- `tools/database/demo-users.sql` creates 15 users with bcrypt-hashed passwords
- CI workflow updated to seed demo users before E2E tests
- Cypress `cy.login()` custom command uses programmatic API login

### Key Files

- `tools/database/demo-users.sql`
- `.github/workflows/ci.yml` (demo user seeding)
- `web/cypress/support/commands.ts` (cy.login custom command)

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
