# Story 34.2: Force-Secure Cookies Behind a Reverse Proxy

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an operator running SitHub behind a TLS-terminating reverse proxy,
I want a configuration option that forces the `Secure` flag on all cookies,
so that session and OAuth-state cookies are never transmitted over plain HTTP regardless of the
scheme the backend observes.

## Acceptance Criteria

1. A new `[main] force_secure_cookies` boolean option exists in config (TOML, env var, and flag),
   documented in the example file, default `false`.
2. With `force_secure_cookies = true` and a request received over plain HTTP (from a TLS-terminating
   proxy), every `Set-Cookie` (session and OAuth-state) includes the `Secure` attribute.
3. With `force_secure_cookies = false` (default), the existing scheme-based behavior
   (`Secure` when `c.Scheme()` is HTTPS) is unchanged.
4. Documentation states the HTTPS-deployment requirement so operators do not enable it for a
   plain-HTTP deployment.
5. Unit tests cover both enabled (Secure present over HTTP) and disabled (existing behavior) cases
   for session and state cookies.

## Tasks / Subtasks

- [x] Task 1: Add the config field and default (AC: #1)
  - [x] `internal/config/config.go`: add `ForceSecureCookies bool \`mapstructure:"force_secure_cookies"\``
        to `MainConfig` (lines 45–50)
  - [x] Add `v.SetDefault("main.force_secure_cookies", false)` in the defaults block (~line 96)
- [x] Task 2: Add the CLI flag + override mapping (AC: #1)
  - [x] `cmd/sithub/main.go`: add field to `runOptions`, bind `--force-secure-cookies` in `bindFlags`,
        and map it to `main.force_secure_cookies` in `overrides` (guarded by `cmd.Flags().Changed`)
- [x] Task 3: Document the option (AC: #1, #4)
  - [x] `sithub.example.toml`: add a documented `force_secure_cookies` entry in `[main]` after
        `data_dir`, following the existing comment style (optional / override flag + env var /
        example / default), and state the HTTPS requirement
- [x] Task 4: Thread the flag to cookie creation (AC: #2, #3)
  - [x] Add `forceSecureCookies bool` to the `auth.Service` struct and set it from
        `cfg.Main.ForceSecureCookies` in `NewService`
  - [x] Add a `Service.NewCookie(c echo.Context, name, value string) *http.Cookie` method that sets
        `Secure: s.forceSecureCookies || c.Scheme() == schemeHTTPS`
  - [x] Replace the four cookie-creation sites to use it (see Dev Notes), including refactoring
        `LogoutHandler` to take `*Service`
  - [x] Update the `LogoutHandler()` call site in `internal/startup/server.go`
- [x] Task 5: Tests (AC: #5)
  - [x] Add tests for forced-secure (HTTP request → Secure set) and default (HTTP request → Secure
        not set) for both session and state cookies

### Review Findings

- [x] [Review][Patch] Default HTTP behavior is not tested for OAuth state cookies [internal/auth/handlers_test.go:41] — fixed by adding default-over-HTTP state-cookie coverage.

## Dev Notes

Source: AI-assisted security review, Finding 2 (Cookie Secure Flag May Be Absent Behind Reverse
Proxy, severity Medium). [Source: private/security-report-claude.md#Finding 2]
[Source: _bmad-output/planning-artifacts/epics.md#Story 34.2 / FR154]

### Cookie creation sites (all must respect the flag)

There are **four** places cookies are created, plus the logout clear:

1. `internal/auth/handlers.go:29` — `LoginHandler`: state cookie `sithub_oauth_state` via
   `newCookie(stateCookieName, encoded, c.Scheme() == schemeHTTPS)`.
2. `internal/auth/handlers.go:114` — `setUserCookieAndRedirect` (used by `CallbackHandler`): session
   cookie `sithub_user`.
3. `internal/auth/login_local.go:99` — `LocalLoginHandler`: session cookie `sithub_user`.
4. `internal/auth/logout.go:14-22` — `LogoutHandler`: clears `sithub_user` with an **inline**
   `http.Cookie` literal (does NOT use `newCookie`), `MaxAge: -1`.

The `newCookie` helper (handlers.go:98–107) takes `(name, value string, secure bool)` and sets
`Path:"/"`, `HttpOnly:true`, `Secure:secure`, `SameSite:http.SameSiteLaxMode`. Constants:
`stateCookieName = "sithub_oauth_state"`, `userCookieName = "sithub_user"` (service.go:24–26);
`schemeHTTPS = "https"` (logout.go:9). [Source: internal/auth/handlers.go:98-107]
[Source: internal/auth/service.go:24-26]

### Recommended approach — centralize on a Service method

The cleanest design (and easiest to test) is a method on `Service` so the flag and the request scheme
are both in scope:

```go
// internal/auth/service.go
func (s *Service) NewCookie(c echo.Context, name, value string) *http.Cookie {
    return &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        HttpOnly: true,
        Secure:   s.forceSecureCookies || c.Scheme() == schemeHTTPS,
        SameSite: http.SameSiteLaxMode,
    }
}
```

Then:

- `LoginHandler`/`CallbackHandler`/`LocalLoginHandler` (which already receive `svc *Service`) call
  `svc.NewCookie(c, name, value)` instead of `newCookie(...)`.
- `LogoutHandler` changes signature from `LogoutHandler()` to `LogoutHandler(svc *Service)` and builds
  the clear-cookie via `svc.NewCookie(c, userCookieName, "")` then sets `MaxAge = -1` on it.
- Update the route in `internal/startup/server.go` (~line 130) from `auth.LogoutHandler()` to
  `auth.LogoutHandler(authService)`. The `authService` is already constructed at server.go:64–67.

> [!NOTE]
> Keep the existing `newCookie` free function if other code/tests depend on it (TestNewCookie at
> handlers_test.go:26 calls it directly). Either keep `newCookie` and have `NewCookie` delegate, or
> migrate the test. Do not silently break `TestNewCookie`.

### Config plumbing (mirror an existing [main] field)

`MainConfig` (config.go:45–50) currently has `Listen`, `Port`, `DataDir` with `mapstructure` tags.
Viper uses env prefix `SITHUB` with `.`→`_` replacement, so `SITHUB_MAIN_FORCE_SECURE_COOKIES=true`
works automatically once the default is set. The flag override pattern in `cmd/sithub/main.go`
(`bindFlags` + `overrides` with `cmd.Flags().Changed(...)`) should mirror how `--listen`/`--data-dir`
are wired. Config reaches auth via `NewService(cfg, store)` (server.go:64–67), so read
`cfg.Main.ForceSecureCookies` there. [Source: internal/config/config.go:45-132]
[Source: cmd/sithub/main.go:47-89] [Source: internal/startup/server.go:64-67]

### Example TOML entry (follow existing style)

Add to `sithub.example.toml` `[main]` after `data_dir` (~line 21), matching the `.claude/rules/toml.md`
conventions:

```toml
  ## Force secure cookies, boolean, optional
  ## When true, session and OAuth state cookies always get the Secure flag, even when the
  ## backend is reached over plain HTTP. Enable this ONLY when a reverse proxy terminates TLS
  ## and the site is actually served over HTTPS to clients.
  ## Can be overridden with --force-secure-cookies flag or SITHUB_MAIN_FORCE_SECURE_COOKIES env var
  ## Example: true
  ## Default: false
  #force_secure_cookies = false
```

[Source: sithub.example.toml:3-21] [Source: .claude/rules/toml.md]

### Testing standards summary

Mirror the existing cookie tests. `logout_test.go:41-60` (`TestLogoutHandlerWithHTTPS`) shows the
pattern: build an Echo context with `httptest.NewRequest(method, "http://..."/"https://...", ...)`,
invoke the handler, then inspect `rec.Result().Cookies()` and assert `cookie.Secure`. Add:

- forced-on over HTTP → `assert.True(cookie.Secure)` (config `ForceSecureCookies: true`)
- default over HTTP → `assert.False(cookie.Secure)`

for both `sithub_user` (login_local / callback / logout) and `sithub_oauth_state` (login). Build the
`Service` with the flag set; the agent report sketches `newTestService(t, cfg)` — reuse whatever the
existing auth tests use to construct a `*Service`. testify `require`/`assert`, run
`golangci-lint run ./...`, `go vet ./...`, `go fmt ./...`. [Source: internal/auth/logout_test.go:41-60]
[Source: .claude/rules/golang.md#Testing]

### Project Structure Notes

Files to modify: `internal/config/config.go`, `cmd/sithub/main.go`, `sithub.example.toml`,
`internal/auth/service.go`, `internal/auth/handlers.go`, `internal/auth/login_local.go`,
`internal/auth/logout.go`, `internal/startup/server.go`. Independent of Story 34.1 — either order.

### References

- [Source: private/security-report-claude.md#Finding 2 — Cookie Secure Flag Behind Reverse Proxy]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 34.2 / FR154]
- [Source: internal/auth/handlers.go:29,98-107,114] [Source: internal/auth/login_local.go:99]
- [Source: internal/auth/logout.go:9-22] [Source: internal/auth/service.go:24-37,59-96]
- [Source: internal/config/config.go:45-132] [Source: cmd/sithub/main.go:47-89]
- [Source: internal/startup/server.go:64-67,130] [Source: sithub.example.toml:3-21]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `go test ./internal/auth/ ./internal/config/ ./cmd/...` → pass
- `golangci-lint run ./...` → 0 issues

### Completion Notes List

- Added `ForceSecureCookies bool` to `config.MainConfig` (`mapstructure:"force_secure_cookies"`) and
  `v.SetDefault("main.force_secure_cookies", false)`.
- Added `--force-secure-cookies` cobra flag + `main.force_secure_cookies` override mapping in
  `cmd/sithub/main.go`. Env override `SITHUB_MAIN_FORCE_SECURE_COOKIES` works automatically via viper.
- Added `forceSecureCookies` field to `auth.Service`, set from `cfg.Main.ForceSecureCookies` in
  `NewService`.
- Centralized cookie creation on a new `Service.NewCookie(c, name, value)` method (with
  `cookieSecure()` helper computing `forceSecureCookies || c.Scheme()==https`). Kept the free
  `newCookie` function so `TestNewCookie` is unaffected; `NewCookie` delegates to it.
- Updated all cookie sites: `LoginHandler` (state cookie), `setUserCookieAndRedirect` (callback
  session), `LocalLoginHandler` (local session), and refactored `LogoutHandler()` →
  `LogoutHandler(svc *Service)` (clears via `NewCookie` + `MaxAge=-1`). Updated the route
  registration in `server.go`.
- Documented `force_secure_cookies` in `sithub.example.toml` `[main]` with the HTTPS-deployment
  caveat, matching the existing comment style.
- Tests: updated existing logout tests for the new signature; added
  `TestLogoutHandlerForceSecureCookiesOverHTTP`, `TestLogoutHandlerDefaultOverHTTP`, and
  `TestServiceNewCookieForceSecureOverHTTP` (covers both state and session cookie names).

### File List

- internal/config/config.go (modified)
- cmd/sithub/main.go (modified)
- internal/auth/service.go (modified)
- internal/auth/handlers.go (modified)
- internal/auth/login_local.go (modified)
- internal/auth/logout.go (modified)
- internal/startup/server.go (modified)
- sithub.example.toml (modified)
- internal/auth/logout_test.go (modified)
- internal/auth/handlers_test.go (modified)

### Change Log

- 2026-06-29: Implemented FR154 force_secure_cookies config and threaded it through all cookie sites.
