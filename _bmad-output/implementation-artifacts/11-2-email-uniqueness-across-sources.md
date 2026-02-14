# Story 11.2: Email Uniqueness Across Sources

Status: done

## Story

As an operator,
I want email addresses unique across all authentication sources,
So that identity conflicts are prevented.

## Acceptance Criteria

1. **Given** an Entra ID user exists with email `alex@example.com`
   **When** an admin attempts to create a local user with the same email
   **Then** the request is rejected with a JSON:API error

2. **Given** a local user exists with email `dana@example.com`
   **When** an Entra ID user with the same email logs in for the first time
   **Then** the login fails with a descriptive error

## Tasks / Subtasks

- [x] Enforce UNIQUE constraint on email column in users table (AC: 1, 2)
- [x] Return JSON:API error on duplicate email during user creation (AC: 1)
- [x] Handle Entra ID login conflict with existing local user (AC: 2)
- [x] Add tests for duplicate email scenarios (AC: 1, 2)

## Dev Notes

### Design Decisions

- Email uniqueness enforced at DB level (UNIQUE constraint) for safety
- Application layer provides descriptive JSON:API errors on constraint violations
- Both creation paths (API and Entra ID login) check for conflicts

### References

- PRD FR29: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.2: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- UNIQUE constraint on email column catches duplicates at DB level
- User creation handler returns 409 Conflict on duplicate email
- Entra ID callback returns error when email matches existing local user

### Key Files

- `migrations/000008_create_users_table.up.sql` (UNIQUE on email)
- `internal/users/store.go` (Create checks for duplicates)
- `internal/auth/service.go` (Entra ID upsert handles conflicts)

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
