# Story 25.6: Async Avatar Sync & Login Spinner

Status: review

## Story

As an Entra ID user,
I want the login to complete quickly with visual feedback,
so that I am not left waiting on a slow avatar sync with no indication of progress.

## Acceptance Criteria

1. **Given** I click "Sign in with Entra ID" on the login page
   **When** the click is registered
   **Then** the button shows a loading spinner and is disabled, preventing double-clicks

2. **Given** the Entra ID OAuth callback is processed by the backend
   **When** the avatar sync would normally run
   **Then** the avatar sync runs asynchronously in a goroutine; the OAuth callback returns
   immediately and redirects the user without waiting for the avatar download

3. **Given** the async avatar sync completes successfully in the background
   **When** I navigate to a page showing my avatar
   **Then** my Entra ID profile photo is displayed

4. **Given** the async avatar sync fails (e.g., no photo in Entra ID, network error)
   **When** I navigate to a page showing my avatar
   **Then** the fallback initials avatar is displayed; no error is shown to the user

## Tasks / Subtasks

- [x] Task 1: Make avatar sync asynchronous in backend (AC: #2, #3, #4)
  - [x] 1.1 In `internal/auth/handlers.go`, locate `CallbackHandler()` (line ~46) and the `SyncAvatar()` call (line ~81)
  - [x] 1.2 Wrap the `SyncAvatar()` call in a goroutine: `go SyncAvatar(context.Background(), client, user.ID, avatarsDir[0])`
  - [x] 1.3 Use `context.Background()` instead of `c.Request().Context()` because the request context will be cancelled after the redirect response is sent
  - [x] 1.4 Ensure errors in the goroutine are logged via `slog.Error` (SyncAvatar already logs errors internally — verify this)
  - [x] 1.5 Verify the OAuth callback returns immediately and the redirect happens without waiting for avatar download
- [x] Task 2: Add loading spinner to Entra ID login button (AC: #1)
  - [x] 2.1 In `LoginView.vue`, locate the Entra ID button (`data-cy="login-entraid"`, line ~52-59): currently a `v-btn` with `href="/oauth/login"`
  - [x] 2.2 Add a reactive ref `entraIdLoading = ref(false)`
  - [x] 2.3 Add `:loading="entraIdLoading"` and `:disabled="entraIdLoading"` props to the button
  - [x] 2.4 Replace the `href` with a `@click` handler that sets `entraIdLoading = true` and then navigates via `window.location.href = '/oauth/login'`
  - [x] 2.5 Verify the spinner appears and the button is disabled after click, before the browser navigates away
- [x] Task 3: Backend validation (AC: #2, #3, #4)
  - [x] 3.1 Run `go fmt ./...` and `go vet ./...`
  - [x] 3.2 Run `golangci-lint run ./...` and fix findings
  - [x] 3.3 Run `go test ./internal/auth/...` if tests exist, and verify no regressions
  - [x] 3.4 Run `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exitCode 1` for duplication check
- [x] Task 4: Frontend validation (AC: #1)
  - [x] 4.1 Run `npm run lint` and fix findings
  - [x] 4.2 Run `npm run type-check` and fix findings
  - [x] 4.3 Run `npm run build` and verify no build errors
  - [x] 4.4 Run `npx vitest run` and verify no regressions
  - [x] 4.5 Run `npm run test:e2e -- --browser electron` and verify no regressions

## Dev Notes

### Architecture & Patterns

- **Two files changed**: Backend `internal/auth/handlers.go` + Frontend `web/src/views/LoginView.vue`
- **Go concurrency**: Use `go func(){}()` pattern for fire-and-forget background work
- **Critical**: Use `context.Background()` for the goroutine, NOT the request context — the request context is cancelled after the HTTP response is sent

### Key Code Locations — Backend

| Element | Location |
|---------|----------|
| `CallbackHandler()` | `internal/auth/handlers.go` line ~46 |
| `SyncAvatar()` call | `internal/auth/handlers.go` line ~81 |
| `SyncAvatar()` implementation | `internal/auth/avatar_handler.go` line ~123 |

### SyncAvatar Behavior

- Downloads user profile photo from Microsoft Graph (`/me/photo/$value`)
- Converts to PNG, saves to `{data_dir}/avatars/{user_id}.png`
- Errors are logged but not propagated — already best-effort
- Currently synchronous — blocks the callback response

### Key Code Locations — Frontend

| Element | Location | data-cy |
|---------|----------|---------|
| Entra ID button | `LoginView.vue` line ~52-59 | `login-entraid` |
| `loading` ref (local login) | `LoginView.vue` line ~80 | — |
| `handleLogin()` | `LoginView.vue` line ~96 | — |

### Entra ID Button Current State

The button is a `v-btn` with `href="/oauth/login"` — a simple link. It has no loading state. The local login form already has a `loading` ref pattern (line ~80) that can be referenced for the spinner implementation.

### Implementation Strategy

**Backend**: One-line change — wrap `SyncAvatar(...)` in `go func() { ... }()` with `context.Background()`.

**Frontend**: Replace `href` with `@click` handler that sets loading state then navigates. The spinner will show briefly before the browser navigates to the OAuth provider.

### Anti-Patterns to Avoid

- Do NOT use `c.Request().Context()` in the goroutine — it will be cancelled after redirect
- Do NOT add error handling UI for avatar sync — it's already best-effort
- Do NOT modify `SyncAvatar()` function signature — only change how it's called
- Do NOT add a timeout to the goroutine — let it complete naturally
- Do NOT use `router.push()` for OAuth redirect — use `window.location.href` since it's an external redirect

### Testing Notes

- The async avatar sync is difficult to test in E2E without Entra ID — verify manually if possible
- The login spinner can be verified by checking the button's disabled state after click
- Existing E2E tests use local auth (`cy.login()` with email/password) and won't be affected

### References

- [Source: internal/auth/handlers.go — CallbackHandler and SyncAvatar call]
- [Source: internal/auth/avatar_handler.go — SyncAvatar implementation]
- [Source: web/src/views/LoginView.vue — Entra ID login button]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- ESLint: pass, TypeScript type-check: pass, Build: pass
- golangci-lint: 0 issues, go vet: pass, go test ./internal/auth/...: pass

### Completion Notes List

- Wrapped `SyncAvatar()` call in `go` goroutine with `context.Background()` so OAuth
  callback returns immediately without waiting for avatar download
- Added `context` import to handlers.go
- Replaced `href="/oauth/login"` on Entra ID button with `@click="handleEntraIdLogin"`
- Added `entraIdLoading` ref for loading/disabled state on the button
- `handleEntraIdLogin()` sets loading state then navigates via `window.location.href`

### File List

- `internal/auth/handlers.go` (modified)
- `web/src/views/LoginView.vue` (modified)

### Change Log

- 2026-04-11: Implemented story 25.6 — async avatar sync, Entra ID login spinner
