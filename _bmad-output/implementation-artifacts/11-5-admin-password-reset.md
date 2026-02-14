# Story 11.5: Admin Password Reset

Status: done

## Story

As an admin,
I want to reset any local user's password,
So that I can help users who are locked out.

## Acceptance Criteria

1. **Given** I am an admin
   **When** I reset a local user's password via `/users/{id}`
   **Then** the user can log in with the new password

2. **Given** I attempt to reset an Entra ID user's password
   **When** the request is processed
   **Then** it is rejected with a JSON:API error

## Tasks / Subtasks

- [x] Add password field to PATCH /api/v1/users/{id} (AC: 1, 2)
  - [x] Only admins can reset passwords
  - [x] Only local users' passwords can be reset
  - [x] Enforce 14-character minimum
  - [x] Hash with bcrypt before storing
- [x] Add tests for admin password reset (AC: 1, 2)

## Dev Notes

### Design Decisions

- Admin password reset uses the existing PATCH /users/{id} endpoint
- Password field is optional in the update request; when present, triggers reset
- Entra ID user password reset returns descriptive error

### References

- PRD FR32: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.5: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- PATCH /api/v1/users/{id} accepts optional `password` field for admin resets
- Validates user is local (not Entra ID) before allowing password change
- bcrypt hashing applied before storage

### Key Files

- `internal/users/store.go` (Update function handles password hash)
- `api-doc/endpoints/user.yaml`

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
