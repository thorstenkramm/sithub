# Story 1.1: Entra ID SSO Login

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an employee,
I want to sign in via Entra ID,
so that I can access SitHub without a separate account.

## Acceptance Criteria

1. **Given** I am not authenticated  
   **When** I open SitHub  
   **Then** I am redirected to Entra ID for login  
   **And** after successful login I return to SitHub and see my name displayed

## Tasks / Subtasks

- [x] Implement Entra ID login initiation (AC: 1)
  - [x] Add `/oauth/login` route to build authorization URL using config values
  - [x] Redirect unauthenticated users to `/oauth/login`
- [x] Implement Entra ID callback handling (AC: 1)
  - [x] Add `/oauth/callback` route to exchange auth code for token
  - [x] Fetch user profile and set authenticated session context
  - [x] Ensure JSON:API error responses on failures
- [x] Show authenticated user name in UI (AC: 1)
  - [x] Add `GET /api/v1/me` (or similar) JSON:API endpoint for user info
  - [x] Surface user name in the main UI shell
- [x] Add tests for the login flow (AC: 1)
  - [x] Cypress E2E: unauthenticated visit redirects to Entra ID
  - [x] Cypress E2E: successful callback returns to app and shows user name
- [x] Add build script for local testing
  - [x] Create `build.sh` to build frontend, embed assets, and compile backend to `./sithub`
  - [x] Add tests for the build script

## Dev Notes

- Use Entra ID settings from `sithub.example.toml` (`authorize_url`, `token_url`, `redirect_uri`, `client_id`,
  `client_secret`, and optional group IDs).
- For E2E tests, enable test auth with `SITHUB_TEST_AUTH_ENABLED=true` and optional
  `SITHUB_TEST_USER_ID`/`SITHUB_TEST_USER_NAME`.
- Enforce JSON:API responses with `application/vnd.api+json` content type for API errors.
- Use `internal/auth` for OAuth flow logic and `internal/middleware` for auth enforcement.
- Ensure Entra ID auth is required before accessing booking data.
- Use `log/slog` and error wrapping with `%w`.

### Project Structure Notes

- Backend handlers: `internal/auth`, `internal/middleware`, `internal/startup` for router wiring.
- Shared JSON:API responses in `internal/api`.
- Frontend user state via Pinia store in `web/src/stores`.

### References

- PRD FR1: `_bmad-output/planning-artifacts/prd.md` (Identity & Access)
- Epic Story 1.1: `_bmad-output/planning-artifacts/epics.md`
- Architecture rules: `_bmad-output/planning-artifacts/architecture.md` (Auth patterns, JSON:API)
- Entra ID config fields: `sithub.example.toml`

## Dev Agent Record

### Agent Model Used

SM - Bob

### Debug Log References
- Cypress Electron emitted `term-size` "Bad CPU type" warnings on macOS; tests still completed.

### Completion Notes List
- Implemented test-auth callback path and TOML config support to validate login without external Entra ID.
- Added CLI flags for all documented config options, including test auth overrides.
- Cypress auth E2E covers unauthenticated redirect and test-auth callback flow; added cookie clearing per test.
- Tests: `golangci-lint run --timeout=5m ./...`; `go test -race ./...`; `npm run test:unit:coverage`;
  `npm run test:e2e -- --browser electron --spec cypress/e2e/auth.cy.ts`.
- Adjusted startup shutdown test to align with clean shutdown behavior.
- Cypress runs: unauth backend (`SITHUB_TEST_AUTH_ENABLED=0`) and test-auth backend (`SITHUB_TEST_AUTH_ENABLED=1`) with
  `CYPRESS_testAuthEnabled=true` for the callback test.
- Added `build.sh` to build frontend assets, embed them, and compile `./sithub` for local testing.

### File List
- assets/embed.go
- assets/embed_test.go
- build.sh
- cmd/sithub/build_script_test.go
- cmd/sithub/main.go
- cmd/sithub/main_test.go
- go.mod
- go.sum
- internal/api/errors.go
- internal/api/response.go
- internal/api/response_test.go
- internal/auth/fetch_user_test.go
- internal/auth/handlers.go
- internal/auth/handlers_callback_success_test.go
- internal/auth/me.go
- internal/auth/service.go
- internal/auth/test_auth.go
- internal/auth/test_auth_test.go
- internal/config/config.go
- internal/config/config_test.go
- internal/db/db.go
- internal/db/db_test.go
- internal/db/migrate.go
- internal/db/migrate_test.go
- internal/middleware/auth.go
- internal/middleware/load_user_test.go
- internal/middleware/session.go
- internal/startup/server.go
- internal/startup/server_test.go
- internal/system/ping.go
- README.md
- sithub.example.toml
- web/cypress/e2e/auth.cy.ts
- web/index.html
- web/package.json
- web/src/api/client.test.ts
- web/vite.config.ts

### Change Log
- 2026-01-17: Added Cypress auth E2E coverage and dev-server support for `/oauth`.
- 2026-01-17: Added test-auth callback path and updated auth E2E scripts.
- 2026-01-17: Added build script and Go test coverage for build output.
- 2026-01-17: Added TOML/CLI test auth config, expanded CLI flags, and lint-driven fixes.
- 2026-01-17: Updated startup shutdown test to expect clean shutdown.
