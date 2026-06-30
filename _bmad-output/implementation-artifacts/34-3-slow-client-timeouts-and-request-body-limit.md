# Story 34.3: Slow-Client Timeouts and Request Body Limit

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an operator,
I want the HTTP server to drop slow or oversized connections,
so that a single malicious or misbehaving client cannot hold resources open or exhaust memory with a
large body.

## Acceptance Criteria

1. The `http.Server` in `internal/startup/server.go` configures `ReadTimeout`, `WriteTimeout`, and
   `IdleTimeout` in addition to the existing `ReadHeaderTimeout`.
2. A client that sends a request body slower than `ReadTimeout` has its connection closed rather than
   held open indefinitely.
3. Echo `BodyLimit` middleware is registered with a 2 MB limit; a request with a JSON body larger
   than 2 MB returns HTTP 413 (Request Entity Too Large).
4. The 2 MB limit does NOT break the existing 4 MB avatar upload (`POST /api/v1/me/avatar`); normal
   booking requests also still succeed.
5. The existing test suite passes with no regression to normal request handling.

## Tasks / Subtasks

- [x] Task 1: Add server timeouts (AC: #1, #2)
  - [x] In `internal/startup/server.go` (struct literal at lines 93–97), add `ReadTimeout`,
        `WriteTimeout`, and `IdleTimeout` alongside `ReadHeaderTimeout: 5 * time.Second`
        (suggested: Read 30s, Write 60s, Idle 120s — see Dev Notes)
- [x] Task 2: Register BodyLimit with an avatar-route skipper (AC: #3, #4)
  - [x] Register `e.Use(middleware.BodyLimitWithConfig(...))` with `Limit: "2M"` and a `Skipper` that
        returns true for `POST /api/v1/me/avatar`
  - [x] Place it as the FIRST middleware (before `Secure`/`LoadUser`) so oversized bodies are
        rejected early
- [x] Task 3: Tests (AC: #3, #4, #5)
  - [x] Add a server test: a request with a >2 MB body to a normal endpoint returns 413
  - [x] Add/verify a test that an avatar upload up to 4 MB is NOT rejected by the global limit
  - [x] Confirm a normal small request still passes

### Review Findings

- [x] [Review][Patch] Avatar body-limit test does not prove the 4 MB upload path still works [internal/startup/server_test.go:127] — fixed by exercising an authenticated multipart avatar upload above 2 MB and below 4 MB through the registered route.
- [x] [Review][Patch] Slow-client timeout behavior is not exercised [internal/startup/server.go:100] — fixed by adding a TCP-level test that sends headers, stalls past `ReadTimeout`, and asserts the request body read fails.
- [x] [Review][Patch] Normal booking requests are not covered under the 2 MB body limit [internal/startup/server_test.go:110] — fixed by driving a normal authenticated booking request through registered routes with the global body-limit middleware and asserting `201 Created`.
- [x] [Review][Patch] Avatar upload route bypasses the global limit without a route-specific 4 MB middleware cap [internal/startup/server.go:130] — fixed by adding route-specific `echomw.BodyLimit("4M")` on `POST /api/v1/me/avatar` and asserting a 5 MB multipart upload returns 413 before auth.

## Dev Notes

Source: AI-assisted security review, Roadmap items #5 (read/write timeouts — Slowloris) and #6 (body
limit). The epic explicitly opted IN to both despite the report's note that timeouts were previously
"excluded per DoS policy". [Source: private/security-report-claude.md#Priority 3 Hardening]
[Source: _bmad-output/planning-artifacts/epics.md#Story 34.3 / FR157,FR158]

### Server timeouts — exact location

The `http.Server` is built at `internal/startup/server.go:93-97`:

```go
server := &http.Server{
    Addr:              addr,
    Handler:           e,
    ReadHeaderTimeout: 5 * time.Second,
}
```

Add the three missing timeouts. Recommended values (from the report's example):

```go
server := &http.Server{
    Addr:              addr,
    Handler:           e,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       30 * time.Second,
    WriteTimeout:      60 * time.Second,
    IdleTimeout:       120 * time.Second,
}
```

`time` is already imported (used by `ReadHeaderTimeout` and the rate limiter). Keep `IdleTimeout >
ReadTimeout`. [Source: internal/startup/server.go:93-97]

> [!CAUTION]
> Check for any streaming/long-lived endpoints before finalizing `WriteTimeout`. SitHub has a
> **live-feed / WebSocket hub** (`livefeed.NewHub()`, passed into `registerRoutes`). A WebSocket
> connection that outlives `WriteTimeout`/`IdleTimeout` could be cut. Echo upgrades the connection
> via `http.Hijacker`, which detaches it from `http.Server` timeouts in most setups — but VERIFY the
> live-feed E2E tests still pass (the repo recently raised live-feed Eventually timeouts to avoid
> `-race` flakes, per recent git history). If WebSockets drop, exclude/relax the timeout for the
> live-feed route or rely on the hub's own keepalive.

### 🚨 CRITICAL GOTCHA — BodyLimit vs. the 4 MB avatar upload

A global `BodyLimit("2M")` will reject avatar uploads **before** the handler runs, returning a raw
413 instead of the handler's friendly 400. The avatar limit is enforced in the handler:

- `internal/auth/avatar_handler.go:23` — `const maxAvatarSize = 4 * 1024 * 1024 // 4 MB`
- `internal/auth/avatar_handler.go:64-67` — rejects `file.Size > maxAvatarSize` with a 400
- Route: `internal/startup/server.go:176-177` — `e.POST("/api/v1/me/avatar", auth.UploadAvatarHandler(avatarsDir), requireAuth)`

Register BodyLimit with a `Skipper` that bypasses the avatar route so its own 4 MB check remains
authoritative:

```go
e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
    Limit: "2M",
    Skipper: func(c echo.Context) bool {
        return c.Request().Method == http.MethodPost && c.Path() == "/api/v1/me/avatar"
    },
}))
```

Note `c.Path()` returns the registered route pattern, which for this fixed path equals
`/api/v1/me/avatar`. Verify the skipper matches by asserting an avatar-sized (e.g. 3 MB) upload is
accepted in a test. [Source: internal/auth/avatar_handler.go:23,64-67]
[Source: internal/startup/server.go:176-177]

### Middleware ordering

Current order: `LoadUser` (server.go:74) → `RedirectForbidden` (server.go:75) → routes (88–90).
Put `BodyLimit` FIRST (reject oversized requests before any other work), then `Secure` (Story 34.1
if landed), then the existing two. Echo is **v4.15.0**, which has `BodyLimit` /
`BodyLimitWithConfig`. [Source: go.mod:10] [Source: internal/startup/server.go:74-90]

### Project Structure Notes

- Only `internal/startup/server.go` is modified for both changes; tests in
  `internal/startup/server_test.go`.
- No new dependency; no config option (limits are fixed per FR157/FR158). If the team later wants
  configurable limits, that's a separate story.

### Testing standards summary

Reuse `setupTestRouter(t)` (server_test.go:73–84) and the `httptest.NewRequest` +
`httptest.NewRecorder` + `e.ServeHTTP(rec, req)` pattern. For the 413 test, send a request body
larger than 2 MB (e.g. `bytes.NewReader(make([]byte, 3<<20))`) to a normal endpoint and assert
`http.StatusRequestEntityTooLarge`. For the avatar test, post a ~3 MB multipart body to
`/api/v1/me/avatar` and assert it is NOT 413 (it will hit auth/handler logic instead). testify
`require`/`assert`; run `golangci-lint run ./...`, `go vet ./...`, `go fmt ./...`, and the live-feed
E2E/integration tests. [Source: internal/startup/server_test.go:73-106]
[Source: .claude/rules/golang.md#Testing]

### References

- [Source: private/security-report-claude.md#Priority 3 — items 5 & 6]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 34.3 / FR157,FR158]
- [Source: internal/startup/server.go:74-97,176-177] (middleware order, server struct, avatar route)
- [Source: internal/auth/avatar_handler.go:23,64-67] (4 MB avatar limit in handler)
- [Source: go.mod:10] (echo v4.15.0 — BodyLimit available)
- [Source: internal/startup/server_test.go:73-106] (test setup)

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `go test ./internal/startup/` → pass (TestBodyLimitRejectsOversizedBody, TestBodyLimitSkipsAvatarUpload)
- full `go test ./...` → pass

### Completion Notes List

- Added `ReadTimeout: 30s`, `WriteTimeout: 60s`, `IdleTimeout: 120s` to the `http.Server` struct
  alongside the existing `ReadHeaderTimeout: 5s`.
- Registered `echomw.BodyLimitWithConfig` with `Limit: "2M"` early in the middleware stack after the
  security-header middleware, with a `Skipper` that bypasses `POST /api/v1/me/avatar` (matched via
  `c.Path()`). Introduced an `avatarUploadPath` constant reused at the POST/DELETE avatar routes
  (avoids goconst on the now-3 occurrences).
- WebSocket safety verified: the live-feed `writePump` calls `SetWriteDeadline(now+10s)` before every
  write/ping (ping interval 30s) and the pong handler refreshes the read deadline to 60s, so the
  hijacked connection's deadlines are continuously refreshed below the 60s server WriteTimeout — the
  live feed is not severed. Full suite (incl. livefeed) passes.
- Tests: `TestBodyLimitRejectsOversizedBody` (3 MB body → 413, with security headers) and
  `TestBodyLimitSkipsAvatarUpload` (authenticated multipart avatar upload above 2 MB and below 4 MB
  succeeds through the registered route).
- Review follow-up tests added: `TestHTTPServerReadTimeoutClosesSlowRequestBody` exercises a stalled
  request body against a real `http.Server`; `TestBodyLimitAllowsNormalBookingRequest` proves a
  normal booking still succeeds under the 2 MB global limit; `TestAvatarUploadHasRouteSpecificBodyLimit`
  proves the avatar route has its own 4 MB middleware cap.

### File List

- internal/startup/server.go (modified)
- internal/startup/server_test.go (modified)

### Change Log

- 2026-06-29: Implemented FR157 (server timeouts) and FR158 (2 MB body limit with avatar skipper).
