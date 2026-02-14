# Story 11.3: Local User Login

Status: done

## Story

As a local user,
I want to log in with my email and password,
So that I can use SitHub without Entra ID.

## Acceptance Criteria

1. **Given** I am a local user with valid credentials
   **When** I enter my email and password in the login form and submit
   **Then** I am authenticated and see my name displayed

2. **Given** I enter an incorrect password
   **When** I submit the login form
   **Then** I see a descriptive error message
   **And** no session is created

## Tasks / Subtasks

- [x] Create POST /api/v1/auth/login endpoint (AC: 1, 2)
  - [x] Accept email and password in JSON body
  - [x] Verify password with bcrypt
  - [x] Use timing-safe comparison to prevent user enumeration
  - [x] Create securecookie session on success
  - [x] Return JSON:API error on failure
- [x] Add per-IP rate limiting (60 requests/min) (AC: 2)
- [x] Make Entra ID config optional (local-only mode) (AC: 1)
- [x] Add tests for login endpoint (AC: 1, 2)

## Dev Notes

### Security Considerations

- Timing-safe password verification prevents user enumeration attacks
- Per-IP rate limiting (60/min) protects against brute force
- bcrypt with default cost for password hashing
- Generic error message ("invalid credentials") regardless of whether email exists

### Design Decisions

- Entra ID configuration is optional; when omitted, only local auth is available
- Login endpoint returns the same securecookie session as Entra ID flow

### References

- PRD FR30: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.3: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- `internal/auth/login_local.go` handles POST /api/v1/auth/login
- Rate limiting via middleware (60 req/min per IP)
- Timing-safe verification prevents enumeration
- Entra ID config is now optional in `internal/config/config.go`

### Key Files

- `internal/auth/login_local.go`
- `internal/auth/login_local_test.go`
- `internal/config/config.go` (optional Entra ID)
- `api-doc/endpoints/auth-login.yaml`

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
