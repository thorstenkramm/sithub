# Story 11.4: Self-Service Password Change

Status: done

## Story

As a local user,
I want to change my own password,
So that I can maintain my account security.

## Acceptance Criteria

1. **Given** I am authenticated as a local user
   **When** I submit a password change via `/me` with a new password of 14+ characters
   **Then** the password is updated
   **And** the old password no longer works

2. **Given** I submit a new password shorter than 14 characters
   **When** the request is processed
   **Then** it is rejected with a validation error

3. **Given** I am an Entra ID user
   **When** I attempt to change my password via `/me`
   **Then** the request is rejected (Entra ID passwords are managed externally)

## Tasks / Subtasks

- [x] Add PATCH /api/v1/me endpoint for password change (AC: 1, 2, 3)
  - [x] Accept current_password and new_password
  - [x] Verify current password before allowing change
  - [x] Enforce 14-character minimum on new password
  - [x] Reject password change for Entra ID users
- [x] Add Cypress E2E tests for password change flow (AC: 1, 2)
- [x] Add unit tests for password change handler (AC: 1, 2, 3)

## Dev Notes

### Design Decisions

- Password change requires current password verification (not just authentication)
- Entra ID users get a clear error explaining passwords are managed externally
- Minimum 14 characters enforced at both API and UI level

### References

- PRD FR31: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.4: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- `internal/auth/me.go` handles PATCH /api/v1/me for password changes
- Validates current password, enforces 14-char minimum, rejects Entra ID users
- Cypress `password-change.cy.ts` covers 5 E2E test cases

### Key Files

- `internal/auth/me.go`
- `internal/auth/update_me_test.go`
- `web/cypress/e2e/password-change.cy.ts`
- `api-doc/endpoints/me.yaml`

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
