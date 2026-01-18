# Story 1.2: Role Determination from Entra ID Groups

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an admin,
I want my role determined from Entra ID group membership,
So that I see admin-only controls.

## Acceptance Criteria

1. **Given** my Entra ID account is in the admin group  
   **When** I log in  
   **Then** the system marks me as admin  
   **And** admin-only cancellation controls are visible

## Tasks / Subtasks

- [x] Determine admin membership from Entra ID groups (AC: 1)
  - [x] Fetch group IDs for the authenticated user
  - [x] Match against configured `admins_group_id` (and `users_group_id` if set)
  - [x] Persist admin flag in auth session cookie
- [x] Expose admin flag via API (AC: 1)
  - [x] Include `is_admin` in `GET /api/v1/me` JSON:API response
- [x] Show admin-only controls in UI (AC: 1)
  - [x] Gate admin-only cancellation controls on `is_admin`
- [x] Add tests (AC: 1)
  - [x] Backend: admin flag set for matching group membership
  - [x] Frontend: admin-only controls hidden for non-admin users

## Dev Notes

- Use Entra ID settings from `sithub.example.toml` (`users_group_id`, `admins_group_id`).
- If `users_group_id` is configured, admins must belong to both groups.
- Enforce JSON:API responses with `application/vnd.api+json` content type for API errors.
- Use `internal/auth` for OAuth flow logic and `internal/middleware` for auth enforcement.
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend handlers: `internal/auth`, `internal/middleware`, `internal/startup` for router wiring.
- Shared JSON:API responses in `internal/api`.
- Frontend user state via Pinia store in `web/src/stores`.

### References

- PRD FR2: `_bmad-output/planning-artifacts/prd.md` (Identity & Access)
- Epic Story 1.2: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (Auth patterns, JSON:API)
- Entra ID config fields: `sithub.example.toml`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

None.
### Completion Notes List
- Added Entra ID group membership fetch and admin determination with users/admin groups handling.
- Added `is_admin` to `/api/v1/me` and persisted admin flag in the auth cookie.
- Added admin-only UI marker in Areas view, gated by `is_admin`.
- Tests: `./run-all-tests.sh`; `go test -race ./...`; `npm run test:unit:coverage`.

### File List
- internal/auth/fetch_user_test.go
- internal/auth/handlers_callback_success_test.go
- internal/auth/me.go
- internal/auth/me_test.go
- internal/auth/service.go
- web/src/api/me.ts
- web/src/views/AreasView.test.ts
- web/src/views/AreasView.vue

### Change Log
- 2026-01-18: Implemented admin role detection from Entra ID groups and UI gating.
