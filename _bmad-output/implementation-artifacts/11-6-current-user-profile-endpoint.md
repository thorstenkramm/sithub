# Story 11.6: Current User Profile Endpoint

Status: done

## Story

As an authenticated user,
I want to retrieve my profile information,
So that the UI can display my identity and role.

## Acceptance Criteria

1. **Given** I am authenticated
   **When** I request `/me`
   **Then** I receive my id, email, name, role, and authentication source

2. **Given** I am not authenticated
   **When** I request `/me`
   **Then** I receive a 401 JSON:API error

## Tasks / Subtasks

- [x] Implement GET /api/v1/me endpoint (AC: 1, 2)
  - [x] Return user profile from session data
  - [x] Include id, email, display_name, role, auth_source
  - [x] Return 401 for unauthenticated requests
- [x] Update frontend auth store to use /me endpoint (AC: 1)
- [x] Add tests for /me endpoint (AC: 1, 2)

## Dev Notes

### Design Decisions

- `/me` returns the session user's profile; no database lookup needed
- Response uses JSON:API format with type "users"
- Frontend auth store fetches /me on app initialization

### References

- PRD FR33: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.6: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- `internal/auth/me.go` handles GET /api/v1/me
- Frontend `useAuthStore` calls /me on app mount
- Returns user profile with role and auth source

### Key Files

- `internal/auth/me.go`
- `web/src/api/me.ts`
- `web/src/stores/useAuthStore.ts`
- `api-doc/endpoints/me.yaml`

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
