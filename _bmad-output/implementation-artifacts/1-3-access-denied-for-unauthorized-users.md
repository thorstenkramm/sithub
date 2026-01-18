# Story 1.3: Access Denied for Unauthorized Users

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a company operator,
I want unauthorized users blocked,
so that booking data is protected.

## Acceptance Criteria

1. **Given** my account is not permitted by Entra ID settings  
   **When** I attempt to access SitHub  
   **Then** I see an access-denied screen  
   **And** I cannot view booking data
2. **Given** I access protected API routes  
   **When** I am unauthenticated or forbidden  
   **Then** API requests return a JSON:API error with 401/403 status

## Tasks / Subtasks

- [x] Enforce permitted users based on Entra ID group settings (AC: 1, 2)
  - [x] Define "permitted" using `users_group_id` (and `admins_group_id` for admin) from `sithub.example.toml`
  - [x] If `users_group_id` is configured, deny non-members even if authenticated
  - [x] Decide and implement test-auth behavior for permission checks (local E2E support)
- [x] Return JSON:API errors for unauthorized/forbidden API access (AC: 2)
  - [x] Add a JSON:API 403 helper in `internal/api` (mirroring `WriteUnauthorized`)
  - [x] Update auth/authorization middleware to return 401 when missing auth and 403 when forbidden
- [x] Add access-denied UI route and screen (AC: 1)
  - [x] Create an access-denied view with clear copy and next steps
  - [x] Add a router path (e.g., `/access-denied`) and ensure SPA fallback allows direct navigation
  - [x] Redirect forbidden users to the access-denied screen after login
- [x] Add tests for access denial behavior (AC: 1, 2)
  - [x] Backend: forbidden membership yields 403 JSON:API error
  - [x] Backend: unauthenticated access yields 401 JSON:API error
  - [x] Frontend: access-denied screen renders and forbidden responses redirect there
  - [x] Cypress E2E: forbidden user sees access-denied screen (use test auth or group config)

## Dev Notes

- Enforce Entra ID access using `entraid.users_group_id`; if empty, allow all users who can authenticate.
- Admins still require `entraid.admins_group_id` (and `users_group_id` when configured).
- Use JSON:API error envelopes with `application/vnd.api+json` content type; add a 403 error helper alongside 401.
- Prefer returning 403 for authenticated-but-forbidden users; use 401 for unauthenticated.
- Use `internal/auth` for OAuth flow and group checks, `internal/middleware` for enforcement, and `internal/api` for error helpers.
- Frontend should handle 403 by routing to the access-denied screen (do not show booking data).
- Add `data-cy` selectors for new UI elements to support Cypress.
- Use `log/slog` and wrap errors with `%w`.

### Project Structure Notes

- Backend: `internal/auth`, `internal/middleware`, `internal/api`, `internal/startup`.
- Frontend: `web/src/views`, `web/src/router`, `web/src/api`, `web/src/stores`.

### Previous Story Intelligence

- Story 1.2 stores admin info in the auth cookie and exposes `is_admin` via `GET /api/v1/me`.
- The SPA redirects to `/oauth/login` on 401 in `web/src/views/AreasView.vue`.
- SPA assets are served from `assets/web` with a fallback in `internal/startup/server.go`.

### Git Intelligence Summary

- Recent work touched auth + group checks and SPA routing assumptions (see commits `33fd91a`, `569a779`).
- Avoid changing the OAuth callback flow without preserving JSON:API error responses.

### References

- PRD FR3: `_bmad-output/planning-artifacts/prd.md` (Identity & Access)
- Epic Story 1.3: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (Auth, JSON:API, WCAG A)
- Entra ID config fields: `sithub.example.toml`
- Current auth flow: `internal/auth/handlers.go`, `internal/middleware/auth.go`

## Dev Agent Record

### Agent Model Used

dev - Amelia

### Debug Log References

- `go test -race` emits a macOS sqlite linker warning in `internal/db`.
- Cypress warns about `term-size` "Bad CPU type" on macOS.

### Completion Notes List

- Enforced user permissions via `users_group_id` and added access-token revalidation per request.
- Added test-auth `permitted` override and 403 JSON:API helper.
- Added access-denied route/view plus server-side redirect for forbidden SPA routes.
- Tests: `./run-all-tests.sh`; `npm run test:unit:coverage`;
  `npm run test:e2e:auth:unauth`; `CYPRESS_testAuthEnabled=true CYPRESS_testAuthPermitted=false npm run test:e2e:auth:test`.

### File List
- internal/middleware/access.go
- internal/middleware/access_test.go
- internal/middleware/helpers_test.go
- internal/startup/server.go
- cmd/sithub/main.go
- internal/api/errors.go
- internal/api/errors_test.go
- internal/auth/fetch_user_test.go
- internal/auth/handlers.go
- internal/auth/handlers_callback_success_test.go
- internal/auth/service.go
- internal/auth/test_auth.go
- internal/auth/test_auth_test.go
- internal/config/config.go
- internal/config/config_test.go
- internal/middleware/auth.go
- internal/middleware/auth_test.go
- sithub.example.toml
- web/cypress/e2e/auth.cy.ts
- web/src/router/index.test.ts
- web/src/router/index.ts
- web/src/views/AccessDeniedView.test.ts
- web/src/views/AccessDeniedView.vue
- web/src/views/AreasView.test.ts
- web/src/views/AreasView.vue

### Change Log
- 2026-01-18: Enforced Entra ID access control with forbidden handling and access-denied UI.
