# Story 34.1: HTTP Security Response Headers

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an operator deploying SitHub,
I want every HTTP response to carry standard security headers,
so that the application is protected against clickjacking, MIME sniffing, and protocol downgrade
without any per-handler work.

## Acceptance Criteria

1. Every HTTP response (API, SPA, static asset) includes `X-Frame-Options: DENY`,
   `X-Content-Type-Options: nosniff`, `Strict-Transport-Security` with `max-age=31536000`,
   `Referrer-Policy: strict-origin-when-cross-origin`, and a `Content-Security-Policy`.
2. The Echo `Secure` middleware is registered in `internal/startup/server.go` before route
   registration so it applies to every route.
3. SitHub cannot be rendered in a cross-origin `<iframe>` (clickjacking blocked).
4. Under the CSP, the Vue SPA, Vuetify styles, Material Design Icons, user avatars, and floor-plan
   images all render correctly with no CSP violations in the browser console.
5. Existing unit and E2E suites pass with no regression to authenticated flows.

## Tasks / Subtasks

- [x] Task 1: Register Echo `Secure` middleware (AC: #1, #2)
  - [x] In `internal/startup/server.go`, add `e.Use(middleware.SecureWithConfig(...))` immediately
        after `e := echo.New()` / `e.HideBanner = true` (around line 38) — BEFORE `LoadUser`,
        `RedirectForbidden`, `registerRoutes`, and `registerSPAHandlers`
  - [x] Configure: `XFrameOptions: "DENY"`, `ContentTypeNosniff: "nosniff"`,
        `HSTSMaxAge: 31536000`, `ReferrerPolicy: "strict-origin-when-cross-origin"`,
        and `ContentSecurityPolicy` (see Dev Notes for the exact value)
- [x] Task 2: Resolve the Content-Security-Policy vs. Google Fonts conflict (AC: #4)
  - [x] Verify in a running browser whether `assets/web/index.html` still loads Google Fonts
        (`https://fonts.googleapis.com` / `https://fonts.gstatic.com`)
  - [x] EITHER include those hosts in the CSP `style-src`/`font-src`, OR self-host the font and
        remove the external `<link>` (preferred for the self-contained-binary philosophy)
- [x] Task 3: Add/extend server tests (AC: #1, #3, #5)
  - [x] In `internal/startup/server_test.go`, assert the security headers are present on a response
        via the existing `setupTestRouter` + `httptest` + `e.ServeHTTP` pattern
- [x] Task 4: Manual CSP verification (AC: #4)
  - [x] Build and run; load the SPA, open an area with floor plans and avatars; confirm zero CSP
        violations in the browser console (Chrome DevTools MCP is available)

### Review Findings

- [x] [Review][Patch] Security headers can be skipped on oversized-request 413 responses [internal/startup/server.go:40] — fixed by applying security headers before body limiting and asserting headers on oversized 413 responses.
- [x] [Review][Patch] HSTS is not present on every response as specified [internal/startup/server.go:146] — fixed by adding unconditional `Strict-Transport-Security: max-age=31536000` middleware and asserting it on plain HTTP responses.
- [x] [Review][Patch] CSP can block the live WebSocket in some browsers [internal/startup/server.go:148] — fixed by explicitly allowing same-request-host `ws://<host>` and `wss://<host>` in `connect-src` and asserting the CSP value.
- [x] [Review][Patch] Browser CSP acceptance was not verified [_bmad-output/implementation-artifacts/34-1-http-security-response-headers.md:176] — manual browser CSP smoke testing passed for the SPA, Vuetify styles, MDI, avatars, floor-plan images, and browser console CSP checks.
- [x] [Review][Patch] CSP WebSocket allowance is broader than same-origin live feed needs [internal/startup/server.go:150] — fixed by generating explicit `ws://<host>` and `wss://<host>` sources from the sanitized request host instead of allowing scheme-wide `ws:` / `wss:` origins.

## Dev Notes

### Where to register (exact location)

`internal/startup/server.go`:

- Echo instance is created at lines 37–38 (`e := echo.New()` / `e.HideBanner = true`).
- The ONLY middleware currently registered is at lines 74–75: `e.Use(middleware.LoadUser(authService))`
  then `e.Use(middleware.RedirectForbidden(authService))`. There is **no** CORS/Logger/Recover/GZip.
- Routes are registered at lines 88–90 (`registerRoutes(...)` then `registerSPAHandlers(...)`).

Register `Secure` right after line 38 so it sets response headers before anything else runs.
`Secure` only sets response headers and has no dependency on request context, so ordering before
`LoadUser` is correct and safe. [Source: internal/startup/server.go:37-90]

> [!IMPORTANT]
> The middleware import is already `github.com/labstack/echo/v4/middleware` (used for the rate
> limiter). Echo is **v4.15.0** ([Source: go.mod:10]), which has both `middleware.SecureWithConfig`
> and `middleware.BodyLimit`. Do not add a new dependency.

### Recommended middleware config

```go
e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
    XSSProtection:      "1; mode=block",
    ContentTypeNosniff: "nosniff",
    XFrameOptions:      "DENY",
    HSTSMaxAge:         31536000,
    ReferrerPolicy:     "strict-origin-when-cross-origin",
    ContentSecurityPolicy: "default-src 'self'; script-src 'self'; " +
        "img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; font-src 'self' data:",
}))
```

`'unsafe-inline'` in `style-src` is needed because Vuetify injects inline `<style>` at runtime —
without it, components will lose styling. This matches the project's own Artifact CSP guidance and
the security report's recommendation. [Source: private/security-report-claude.md#Finding 1]

### 🚨 CRITICAL GOTCHA — Google Fonts breaks the CSP

`assets/web/index.html` currently has **no inline scripts/styles** (good — `script-src 'self'`
works), and Vuetify CSS + MDI icons are served locally from `/assets/` (same-origin, fine).
**However, `index.html` loads Google Fonts from `https://fonts.googleapis.com`** (and the font files
come from `https://fonts.gstatic.com`). The CSP above does NOT permit those hosts, so the font
request will be blocked and the console will show a CSP violation. [Source: assets/web/index.html:15]

Pick one (Task 2):

1. **Allow the hosts** (smaller change): append to the CSP —
   `style-src 'self' 'unsafe-inline' https://fonts.googleapis.com` and
   `font-src 'self' data: https://fonts.gstatic.com`.
2. **Self-host the font** (preferred, fully self-contained — aligns with SitHub's embedded-binary
   philosophy): remove the Google Fonts `<link>` from the frontend, bundle the font via the build
   (or rely on the system font stack), and keep the strict `font-src 'self' data:`.

Confirm the actual behavior in a browser before choosing — the font link may already have been
removed in a prior epic. Do not ship a CSP that produces console violations (AC #4).

### HSTS behind a reverse proxy (interaction with Story 34.2)

Echo's `Secure` middleware emits the HSTS header. In the common deployment, TLS is terminated by a
reverse proxy and SitHub sees plain HTTP — that's expected; the header still instructs the browser.
Do not gate HSTS on `c.Scheme()`. This story and Story 34.2 (`force_secure_cookies`) are independent
and can land in either order. [Source: _bmad-output/planning-artifacts/epics.md#Story 34.2]

### Custom SPA error handler interaction

`registerSPAHandlers` (server.go:234–261) installs a custom `HTTPErrorHandler` that serves
`index.html` on GET 404s for non-API paths. The `Secure` middleware runs earlier in the chain and
only sets response headers, so it composes cleanly with the fallback — no special handling needed.
[Source: internal/startup/server.go:234-261]

### Project Structure Notes

- Single backend file touched for the middleware: `internal/startup/server.go`.
- Possible frontend touch only if choosing the self-host option in Task 2 (`assets/web/index.html`
  is generated from the Vite build under `web/`; change the source, then rebuild — do not hand-edit
  the embedded `assets/web/` output).
- Test file: `internal/startup/server_test.go`.

### Testing standards summary

- Go: table-driven where useful, `testify` `require`/`assert`, reuse `setupTestRouter(t)` (server_test.go:73–84)
  and the `httptest.NewRequest` + `httptest.NewRecorder` + `e.ServeHTTP(rec, req)` pattern
  (server_test.go:86–106). Assert headers via `rec.Result().Header.Get(...)`.
- Run `golangci-lint run ./...` (v2.5.0, 120-char lines), `go vet ./...`, `go fmt ./...`.
- Browser check via Chrome DevTools MCP for CSP violations.
  [Source: .claude/rules/golang.md] [Source: .claude/rules/vue.md]

### References

- [Source: private/security-report-claude.md#Finding 1 — Missing HTTP Security Headers]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 34.1 / FR153]
- [Source: internal/startup/server.go:37-90] (Echo init + middleware order + route registration)
- [Source: internal/startup/server.go:234-261] (SPA fallback error handler)
- [Source: go.mod:10] (echo v4.15.0)
- [Source: assets/web/index.html:15] (Google Fonts external link — CSP risk)
- [Source: internal/startup/server_test.go:73-106] (test setup + request pattern)

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `go test ./internal/startup/` → pass (TestSecurityHeadersPresent, TestBodyLimitRejectsOversizedBody)
- `golangci-lint run ./...` → 0 issues; `go vet ./...` → clean

### Completion Notes List

- Added `echomw "github.com/labstack/echo/v4/middleware"` import and registered
  `echomw.SecureWithConfig(secureConfig())` immediately after `e.HideBanner = true` (first
  response-header middleware, before routes).
- Added `secureConfig()` helper emitting `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`,
  `Strict-Transport-Security` (max-age 31536000), `Referrer-Policy: strict-origin-when-cross-origin`,
  `X-XSS-Protection`, and a CSP.
- CSP/Google Fonts decision: `assets/web/index.html` DOES load Inter from `fonts.googleapis.com`
  (stylesheet) and `fonts.gstatic.com` (fonts). Chose to allow those hosts in the CSP
  (`style-src ... https://fonts.googleapis.com`, `font-src ... https://fonts.gstatic.com`) plus
  explicit `ws:`/`wss:` in `connect-src` for the live-feed WebSocket — a backend-only change that
  doesn't touch the generated build output.
- Dependency: Echo's middleware subpackage pulls in `golang.org/x/time/rate`, which was absent from
  go.sum (project used a custom limiter). Ran `go mod tidy`; `go.mod`/`go.sum` now include
  `golang.org/x/time`.
- HSTS nuance: Echo's Secure middleware only emits `Strict-Transport-Security` for HTTPS requests
  (TLS or `X-Forwarded-Proto: https`), so review added a small unconditional HSTS middleware to meet
  AC #1's "every response" wording. Tests assert HSTS on plain HTTP and on oversized 413 responses.
- CSP WebSocket scope tightened after review: `connect-src` now keeps `'self'` and adds only
  same-request-host `ws://<host>` / `wss://<host>` sources when the Host value can be safely
  represented in a CSP source. The test parses `connect-src` and rejects broad `ws:` / `wss:`
  sources.
- Task 4 (live browser CSP smoke test) was completed manually. The SPA, Vuetify styles, MDI,
  avatars, and floor-plan images rendered correctly with no CSP browser-console violations observed.

### File List

- internal/startup/server.go (modified)
- internal/startup/server_test.go (modified)
- go.mod (modified)
- go.sum (modified)

### Change Log

- 2026-06-29: Implemented FR153 HTTP security response headers via Echo Secure middleware.
