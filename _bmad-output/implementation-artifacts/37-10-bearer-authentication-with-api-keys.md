# Story 37.10: Bearer Authentication with API Keys

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an integrator,
I want to authenticate with an API key,
so that scripts and finance tools can read bookings automatically.

## Acceptance Criteria

1. A request with `Authorization: Bearer <key>` whose SHA-256 hash matches a stored, unexpired key is
   authenticated as that key (not as a user). The bearer header is checked before the session cookie.
2. An unknown, revoked, or expired key gets 401 with a JSON:API error.
3. A key with scope "read bookings" can call `GET /api/v1/bookings/export` and
   `GET /api/v1/bookings/export.csv` (Stories 37.7/37.8). The request succeeds with the same filters and
   pagination as in the UI.
4. A key with scope "read bookings" gets 403 on any other endpoint (admin, booking creation, users,
   api-keys).
5. Key comparison uses the stored hash with a constant-time comparison. The key value is never logged.
6. `docs/public-api.md` explains how an administrator creates a key, the bearer header, the "read
   bookings" scope and its endpoints with parameters, expiry and revocation, and a working curl example.
   It passes markdownlint.
7. The OpenAPI docs in `api-doc` document the bearer security scheme and the export endpoints and pass
   the redocly lint.
8. A key with scope "read bookings" can also call `GET /api/v1/avatars/{user_id}` (read-only; uploading and
   deleting avatars stays denied), so integrations such as door signage can show the people listed by the
   export. Sessions keep their existing access. Added 2026-10-11 at the product owner's request.

FRs covered: FR191, FR192.
[Source: _bmad-output/planning-artifacts/epics.md#Story 37.10 (line ~6025), FR191/FR192 (lines ~708-714)]

## Tasks / Subtasks

- [x] Task 1: API-key principal type and context helpers (AC: #1, #4)
  - [x] Add `internal/auth/api_key.go` with `type APIKeyPrincipal struct { ID, Name, Scope string }`
        and `const apiKeyContextKey = "api_key"`.
  - [x] Add `SetAPIKeyInContext(c, *APIKeyPrincipal)` and `GetAPIKeyFromContext(c) *APIKeyPrincipal`
        mirroring `GetUserFromContext` (`internal/auth/me.go:18-26`). The principal is stored under its
        own context key and is NEVER stored under `"user"`.
  - [x] Add `func (p *APIKeyPrincipal) HasScope(scope string) bool` (exact string match).
  - [x] Doc comments explain why a key is not a user (see Dev Notes "Why a key is not a user").
- [x] Task 2: Key verification in the api-keys package from 37.9 (AC: #1, #2, #5)
  - [x] Add `Authenticate(ctx, store, rawKey string, now time.Time) (*Key, error)` to the package that
        37.9 created (`internal/apikeys`).
  - [x] Reject malformed input before any DB access: the raw key must parse as a UUID
        (`uuid.Parse` or a strict 36-char check). Return the sentinel `ErrInvalidAPIKey`.
  - [x] Hash the raw key with `apikeys.HashKey` from 37.9 (lowercase hex SHA-256; do not re-implement).
        Look the row up via the unique index on `key_hash` (37.9 suggests `FindValidByHash`), then confirm
        with `subtle.ConstantTimeCompare(stored, computed) == 1`. No row or no match: `ErrInvalidAPIKey`.
  - [x] Expiry: the key is expired when `!now.Before(expiresAt)` (UTC). Return the sentinel
        `ErrAPIKeyExpired`.
  - [x] Revoked keys are deleted rows (37.9), so they fall into `ErrInvalidAPIKey` automatically.
  - [x] Doc comment marks `Authenticate` as security-critical (`.claude/rules/golang.md`
        "Document Security-Critical Functions").
- [x] Task 3: Global `LoadAPIKey` middleware with bearer-before-cookie precedence (AC: #1, #2, #5)
  - [x] New `internal/middleware/api_key.go`: `LoadAPIKey(store *sql.DB) echo.MiddlewareFunc`
        (inject a `now func() time.Time` via an unexported constructor for tests).
  - [x] Only act on `/api/` paths and only when the `Authorization` header uses the `Bearer` scheme
        (case-insensitive scheme, single space, non-empty token). Any other scheme (e.g. `Basic` from a
        reverse proxy) is ignored and the request continues to cookie auth.
  - [x] Valid key: `auth.SetAPIKeyInContext`, then `next`. Invalid/expired/empty token: respond 401
        immediately. Never fall back to the session cookie when a bearer token was presented.
  - [x] 401 responses: JSON:API error via a new `api.WriteUnauthorizedDetail(c, detail)` (refactor
        `WriteUnauthorized` to call it, same pattern as `WriteForbiddenDetail` in `internal/api/errors.go`,
        to keep jscpd clean). Details: "Invalid API key" / "API key has expired". Set
        `WWW-Authenticate: Bearer error="invalid_token"` (RFC 6750).
  - [x] Update `LoadUser` (`internal/middleware/session.go:10-23`): if `auth.GetAPIKeyFromContext(c) != nil`,
        skip cookie decoding so `"user"` stays unset.
  - [x] Register `e.Use(middleware.LoadAPIKey(store))` BEFORE `e.Use(middleware.LoadUser(authService))`
        in `internal/startup/server.go:81`. Extract both lines plus `RedirectForbidden` into a helper
        (e.g. `registerAuthMiddleware(e, authService, store)`) so tests use the real order.
- [x] Task 4: Default-deny for API keys on every user route (AC: #4)
  - [x] `RequireAuth` (`internal/middleware/auth.go:12-28`): if an API-key principal is present, return
        `api.WriteForbiddenDetail(c, "API key not permitted for this endpoint")` before the user check.
  - [x] `RequireAdmin` (`internal/middleware/require_admin.go:11-24`): same 403 for defense in depth.
  - [x] Apply the same 403 rule to `RequireDataExporter()` from 37.6
        (`internal/middleware/require_data_exporter.go`) for defense in depth.
- [x] Task 5: Scope-based access to the export endpoints (AC: #3, #4)
  - [x] Add a composite middleware, e.g. `RequireExportAccess(authSvc, store, scope)`: with an API-key
        principal, allow only if `HasScope(scope)` (else 403); without one, delegate to
        `RequireAuth(authSvc)` + the 37.6 data-exporter check (unchanged session behavior).
  - [x] Replace the `exportAccess` slice (`requireAuth, requireDataExporter`, defined once by 37.7 in
        `registerRoutes` and reused by 37.8) with this composite, so both export routes change in one place.
        Use `apikeys.ScopeReadBookings` (`"read-bookings"`, 37.9); do not hard-code a new string.
  - [x] Verify the 37.7/37.8 export handlers never dereference `auth.GetUserFromContext(c)` without a
        nil check. With a key, the user is nil. Fix any such access.
  - [x] Make no other route API-key-aware. The api-keys routes from 37.9 stay admin-only and return 403
        for keys via Task 4.
- [x] Task 6: Logging hygiene (AC: #5)
  - [x] Never log the raw header, the raw key, the prefix of a presented key, or the computed hash.
  - [x] On rejection log at `Warn` with `reason` (`malformed|unknown|expired`), `path` and
        `remote_ip` only. On success optionally log at `Debug` with `api_key_id` and `api_key_name`.
  - [x] Confirm no request logger middleware is registered that dumps headers (none exists today).
- [x] Task 7: Go unit and security regression tests (AC: #1-#5)
  - [x] `internal/apikeys/*_test.go` (table-driven): valid; malformed; unknown UUID; deleted (revoked);
        expired (`expires_at = now-1s`); boundary `expires_at == now` is expired; uppercase variant of a
        valid key is rejected (hash is over the exact string).
  - [x] `internal/middleware/api_key_test.go` (table-driven): no header passes through; `Bearer <valid>`
        and `bearer <valid>` set the principal; `Bearer` with empty token, `Bearer not-a-uuid`, unknown,
        revoked and expired give 401 with JSON:API body, `application/vnd.api+json`, and
        `WWW-Authenticate`; `Basic dXNlcjpwYXNz` is ignored; non-`/api/` path is ignored.
  - [x] Precedence regression: valid admin cookie plus invalid bearer gives 401 (no cookie fallback).
        Valid admin cookie plus valid key on `GET /api/v1/users` gives 403 (key wins, admin cookie ignored).
  - [x] `RequireAuth`/`RequireAdmin` with only a key principal give 403.
  - [x] Composite middleware table: key with scope 200; key with another scope 403; exporter user 200;
        admin without flag 403; nobody 401.
  - [x] Route-sweep regression in `internal/startup/server_test.go`: build the router with the real
        middleware helper from Task 3 and a migrated store, seed a valid key, iterate `e.Routes()` and
        assert that every `/api/` route except the allowlist (two export routes plus public
        `/ping`, `/auth/login`, `/auth/logout`, `/auth/providers`) returns 403 for a bearer-only request.
        This catches future routes that accidentally accept keys.
  - [x] Logging regression: capture `slog` into a buffer (`slog.SetDefault` with a `TextHandler`,
        restore in `t.Cleanup`, no `t.Parallel()`), send valid, unknown and expired keys, and assert the
        buffer contains neither the key string nor its hash.
- [x] Task 8: Cypress E2E against the dev server (AC: #1-#4)
  - [x] New `web/cypress/e2e/api-key-auth.cy.ts`, only `cy.request` (no UI, no mocking). Admin
        `anna@sithub.local` via `cy.login()` creates a key through `POST /api/v1/api-keys`. Read the full
        key from `data.attributes.key` of the 201 response (37.9).
  - [x] Call `cy.clearCookies()`, then with `Authorization: Bearer <key>`: export JSON gives 200 with
        `data` and pagination `meta`/`links`; export CSV gives 200 with `text/csv`;
        `GET /api/v1/users`, `GET /api/v1/api-keys`, `POST /api/v1/bookings` and `GET /api/v1/me` give 403;
        a random UUID gives 401 (`failOnStatusCode: false`).
  - [x] Precedence: logged in as admin (cookie present) plus valid key on `/api/v1/users` gives 403.
  - [x] Revocation: as admin, delete the key, clear cookies, and verify the same key now gets 401.
  - [x] Expired keys cannot be created through the API (min 3 months); expiry is covered by Go tests.
        Note this in the spec file.
  - [x] Use unique key names (timestamp suffix) and revoke created keys in `afterEach` so reruns stay clean.
- [x] Task 9: `docs/public-api.md` (AC: #6)
  - [x] Write the document following the outline in Dev Notes. Follow `.claude/rules/markdown.md`: real
        headings, 120-char lines, GFM alerts, no single-item lists.
  - [x] Add a link to it in `docs/index.md`.
  - [x] Run `npx --yes markdownlint-cli docs/public-api.md` until clean.
  - [x] Manually run both curl examples against a local server with a real key. Record the commands
        (key redacted) in the Debug Log.
- [x] Task 10: OpenAPI (AC: #7)
  - [x] `api-doc/openapi.yaml`: add `components.securitySchemes.bearerAuth`
        (`type: http`, `scheme: bearer`, `bearerFormat: UUID`, with a description). Keep the global
        `security: [cookieAuth]`.
  - [x] Ensure `api-doc/endpoints/bookings-export.yaml` and `bookings-export-csv.yaml` exist (37.7/37.8
        may have created them) and are referenced under `/bookings/export` and `/bookings/export.csv`.
        Each operation sets `security: [{cookieAuth: []}, {bearerAuth: []}]` (OR semantics), documents
        `from`, `to`, `area`, `page[limit]`, `page[offset]` (JSON only), and 400/401/403 responses using
        `ErrorResponse`.
  - [x] Mention in the api-keys endpoint docs (37.9) that bearer keys get 403 there.
  - [x] `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml` passes.
- [x] Task 11: Quality gates
  - [x] `go fmt ./...`, `go vet ./...`, `golangci-lint run ./...`, `go test ./...`.
  - [x] Go jscpd (`npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 3`); baseline is 1.51%.
  - [x] `npm run lint`, `npm run type-check` and `npm run test:e2e -- --browser electron` (Cypress spec
        is TypeScript), or run `./run-all-tests.sh`.

### Review Findings

Code review 2026-10-08 (Blind Hunter, Edge Case Hunter, Acceptance Auditor; auto-fix): 0 decision-needed,
12 patch (all fixed), 7 deferred, 11 dismissed.

- [x] [Review][Patch] `Bearer<key>` without separator fell back to the session cookie [internal/middleware/api_key.go:47]
- [x] [Review][Patch] Credential-store failures were not logged; cancellations looked like outages [internal/middleware/api_key.go:66]
- [x] [Review][Patch] No Debug audit log with `api_key_id`/`api_key_name` on successful key use [internal/middleware/api_key.go:38]
- [x] [Review][Patch] Wrong-scope 403 lacked RFC 6750 `insufficient_scope` challenge [internal/middleware/export_access.go:16]
- [x] [Review][Patch] Route sweep only checked status code, never checked invalid bearer plus admin cookie, and did not verify allowlist entries exist [internal/startup/api_key_test.go:35]
- [x] [Review][Patch] Precedence/export router tests looped without `t.Run`; one message was wrong for the Basic case [internal/startup/api_key_test.go:60]
- [x] [Review][Patch] "no header" middleware case sent an empty `Authorization` header instead of none [internal/middleware/api_key_test.go:37]
- [x] [Review][Patch] `TestAPIKeyRoutesRequireAdmin` and `TestBookingExportRouteAccess` bypassed `registerAuthMiddleware` [internal/startup/server_test.go:446]
- [x] [Review][Patch] Cypress precedence test never proved the admin cookie was present [web/cypress/e2e/api-key-auth.cy.ts:70]
- [x] [Review][Patch] Docs claimed public endpoints stay public under an invalid bearer; no warning about proxy-injected Bearer headers [docs/public-api.md:30]
- [x] [Review][Patch] Errors table cell over 75 chars; curl example lacked the `export SITHUB_API_KEY` line [docs/public-api.md:101]
- [x] [Review][Patch] api-keys 403 descriptions did not mention bearer denial [api-doc/endpoints/api-keys.yaml:20]
- [x] [Review][Defer] No rate limiting or log throttling for failed bearer attempts [internal/middleware/api_key.go:84] — deferred, spec open question 5
- [x] [Review][Defer] `c.RealIP()` trusts `X-Forwarded-For` because `e.IPExtractor` is unset [internal/startup/server.go] — deferred, pre-existing
- [x] [Review][Defer] No `last_used_at` usage tracking for keys [internal/apikeys] — deferred, spec open question 7
- [x] [Review][Defer] Corrupt `expires_at` is logged as `reason=unknown` instead of an integrity error [internal/apikeys/authenticate.go:45] — deferred, deliberate fail-closed choice
- [x] [Review][Defer] RFC 3339 formatting truncates expiry to whole seconds (up to 1s early) [internal/apikeys/store.go:86] — deferred, pre-existing (37.9)
- [x] [Review][Defer] OpenAPI `PageLimit` declares `maximum: 100` while the server caps larger values [api-doc/openapi.yaml] — deferred, pre-existing (37.4)
- [x] [Review][Defer] Generic 401 from `RequireAuth` carries no `WWW-Authenticate` challenge [internal/middleware/auth.go:21] — deferred, pre-existing

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md`, Epic 37 (`### Epic 37`, line ~1264) and
`### Story 37.10` (line ~6025). Product decisions are final (briefing): keys are server-generated UUID v4,
stored as SHA-256 hash plus 10-char prefix, admins revoke by deleting, and scope "read bookings" covers ONLY
the export endpoints of 37.7/37.8.

### Dependencies (assumed implemented as described in epics.md)

- 37.4: shared pagination helpers (`page[limit]`, `page[offset]`, `PaginationMeta`/`PaginationLinks` in
  `internal/api`).
- 37.6: `users.is_data_exporter`, `auth.User.IsDataExporter` (refreshed from the DB in `RequireAuth`) and
  `middleware.RequireDataExporter()` (`internal/middleware/require_data_exporter.go`, 401 without user,
  403 without flag, no admin fallback). [Source: 37-6-data-exporter-role.md Task 5]
- 37.7/37.8: `GET /api/v1/bookings/export` (JSON:API, paginated) and `GET /api/v1/bookings/export.csv`,
  both registered with the shared `exportAccess := []echo.MiddlewareFunc{requireAuth, requireDataExporter}`
  slice, explicitly left for 37.10 to replace. [Source: 37-7-...md Task 4; 37-8-...md "Forward
  compatibility with Story 37.10"]
- 37.9: package `internal/apikeys`, migration `000006_api_keys` (`key_hash TEXT` with a UNIQUE index,
  `scope` CHECK `'read-bookings'`, `expires_at` RFC 3339 TEXT, `created_by` without FK), constants
  `ScopeReadBookings = "read-bookings"` and `PrefixLength = 10`, `HashKey(key) string` (lowercase hex
  SHA-256, exported for this story), create response attribute `key`. Keys are `uuid.NewString()`.
  [Source: 37-9-api-key-management.md Tasks 1-3, Dev Notes]

Before coding, open the merged code of these stories and use their actual names. If a name differs from
the names above, follow the code, not this story.

### Current auth chain (verified 2026-10-08)

Global middleware in `internal/startup/server.go:40-82`, in this order:

1. `echomw.SecureWithConfig`, `strictTransportSecurity`, `contentSecurityPolicy` (Epic 34 security headers;
   they only set response headers and do not affect bearer requests).
2. `echomw.BodyLimitWithConfig` (2 MB; GET export requests are unaffected).
3. `middleware.LoadUser(authService)` (`internal/middleware/session.go:10-23`): decodes the securecookie
   `sithub_user` and sets `c.Set("user", *auth.User)`. It never fails the request.
4. `middleware.RedirectForbidden(authService)` (`internal/middleware/access.go:13-41`): applies to non-API
   GET paths only (`isBypassPath` skips `/api/`), so it is irrelevant for bearer requests.

Per route: `requireAuth := middleware.RequireAuth(authService)` (`internal/middleware/auth.go:12-28`, 401
without user, calls `RefreshPermissions`, 403 if not permitted), and
`requireAdmin := middleware.RequireAdmin()` (`internal/middleware/require_admin.go:11-24`).
Routes are listed in `registerRoutes` (`internal/startup/server.go:224-312`).

### CSRF, origin and CORS: nothing blocks bearer requests

- There is NO CSRF or origin-check middleware. CSRF mitigation is `SameSite=Lax` on the session cookie
  (`internal/auth/handlers.go:105`). The only `Origin` check is the WebSocket upgrade in
  `internal/livefeed/handler.go:49-70`, and that route stays key-forbidden anyway.
- Bearer auth is CSRF-safe by construction because browsers never attach `Authorization` automatically.
  That is also why a presented bearer token must win over the cookie: a request that carries both is
  coming from a script that set the header explicitly.
- No CORS middleware is registered, so browser JavaScript from other origins cannot read responses.
  Server-to-server clients (curl, finance tools) are unaffected. Document "server-side clients only" in
  `docs/public-api.md`; do NOT add CORS in this story.
- Rate limiting exists only on `POST /api/v1/auth/login` (`internal/middleware/rate_limit.go`,
  `server.go:237-239`). See the open questions for bearer failures.

### Why a key is not a user

Handlers authorize by reading `auth.GetUserFromContext(c)` (`internal/auth/me.go:18-26`) and its fields
`ID`, `IsAdmin`, `IsPermitted`. They also use `auth.ResolveUserEmail`, booking ownership checks
(`internal/bookings/handler.go:110,256`) and `users` lookups by ID. A synthetic `*auth.User` for a key would:

- pass `RequireAuth` on every route (default-allow instead of default-deny),
- let booking creation, `/me`, colleagues, avatars and the live feed act on a fake user ID,
- become admin-capable through a single wrong boolean.

Therefore the key principal lives under its own context key (`"api_key"`), `"user"` stays nil, and
`RequireAuth` rejects keys with 403. Access is granted only where a route explicitly opts in through the
scope-aware composite middleware (Task 5). The route-sweep test (Task 7) enforces this for future routes.

### Status codes

- No or unknown/revoked/expired/malformed bearer token: 401 (`WWW-Authenticate: Bearer error="invalid_token"`).
- Valid key on a route outside its scope (including all `requireAuth` and admin routes): 403.
- Valid key without the "read bookings" scope on the export routes: 403.
- Validation errors on export (from > to, unknown area, bad dates): 400, unchanged from 37.7.

### Constant-time comparison and hashing

- Hash with `apikeys.HashKey` (lowercase hex of `sha256.Sum256`), exactly as 37.9 stores it.
- Look up by `key_hash` through the unique index, then compare the stored and computed hex strings with
  `crypto/subtle.ConstantTimeCompare`. The SQL lookup operates on a hash, not the secret, so index timing
  reveals nothing usable; the explicit constant-time check satisfies AC #5 and guards against a future
  change to prefix-based lookup. Do not look up by `key_prefix` (about 36 bits, not unique).
- A UUID v4 has 122 random bits; brute force is infeasible. Constant-time comparison is still required by
  the AC and costs nothing.
- Store timestamps and compare in UTC (`.claude/rules/json-api.md`, Date and Time). Parse `expires_at` the
  same way 37.9 writes it.

### Logging redaction

No request-logging middleware is registered (only `slog` calls in startup). Echo's default `HTTPErrorHandler`
does not log headers. The risk is new code: never pass `c.Request().Header`, the token, its prefix or its
hash to `slog`. Prefixes of a presented invalid key are attacker-controlled. Logging them adds noise and,
for a valid key, leaks 10 of 36 chars.

### `docs/public-api.md` outline

Write in English; product name "SitHub". Suggested structure:

1. `# SitHub Public API`: purpose (machine access to booking data for finance and scripts), server-side
   clients only.
2. `## Create an API key`: admin opens the user menu, then "API Keys", then "Add key"; enters a name and
   picks an expiry (3, 12 or 24 months) and the scope "Read bookings". Use a `> [!WARNING]` alert: the key
   is shown once; copy it into a secret store. Only the first 10 characters remain visible in the table.
3. `## Authenticate`: header `Authorization: Bearer <key>`; HTTPS only; bearer takes precedence over a
   session cookie.
4. `## Scopes`: small table "Scope / Grants" with one row: Read bookings, the two export endpoints. All
   other endpoints return 403.
5. `## Endpoints`
   - `### GET /api/v1/bookings/export`: parameters `from`, `to` (YYYY-MM-DD, inclusive, required),
     `area` (top-level area ID, optional, omit for all areas), `page[limit]` (default 20, max per 37.4),
     `page[offset]`. Short JSON:API response sample (snake_case attributes as implemented by 37.7) and
     sorting (date, then name).
   - `### GET /api/v1/bookings/export.csv`: same filters, no pagination, UTF-8 with BOM, RFC 4180, columns,
     file name pattern `sithub-bookings-<from>_<to>.csv`.
6. `## Errors`: 400/401/403 with one JSON:API error example.
7. `## Expiry and revocation`: expired keys return 401; revoke by deleting the key; revocation is
   immediate; to rotate, create a new key and then delete the old one.
8. `## Examples`: working curl commands. Use `-g` (`--globoff`) because curl otherwise treats `[ ]` in
   `page[limit]` as a glob:

   ```bash
   export SITHUB_API_KEY="<your-key>"
   curl -g -H "Authorization: Bearer ${SITHUB_API_KEY}" \
     "https://sithub.example.com/api/v1/bookings/export?from=2026-10-01&to=2026-10-08&page[limit]=50&page[offset]=0"
   curl -H "Authorization: Bearer ${SITHUB_API_KEY}" -o bookings.csv \
     "https://sithub.example.com/api/v1/bookings/export.csv?from=2026-10-01&to=2026-10-08&area=<area-id>"
   ```

   The repo's `.markdownlint.json` sets `MD013.code_blocks: false`, so long URLs inside code blocks are
   fine. Mention reading the key from an environment variable instead of typing it into the shell history.

### OpenAPI snippet

```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: UUID
      description: API key created by an administrator (scope "read bookings"); see docs/public-api.md.
```

Per export operation: `security: [{cookieAuth: []}, {bearerAuth: []}]`. The global default remains
`cookieAuth` (`api-doc/openapi.yaml:9-10,56-60`). Existing endpoint files live in `api-doc/endpoints/`.

### Test conventions (reuse, do not reinvent)

- testify `require` for setup and `assert` for behavior, table-driven with `t.Run`
  (`.claude/rules/golang.md#Testing`).
- Middleware helper `runMiddleware` (`internal/middleware/helpers_test.go:13-36`). Model the precedence and
  invalid-credential tests on `TestInvalidSessionCookieRejected`
  (`internal/middleware/security_session_test.go`, Story 34.5 style: one security property per test,
  with an assertion message that states the property).
- Router-level helpers in `internal/startup/server_test.go`: `setupStartupTestStore` (migrated temp DB,
  :324), `setupTestRouter` (:383, currently registers only `LoadUser`; switch it to the new middleware
  helper), `newTestAuthService` (:445), `testUserCookie` (:461).
- Seed keys with `apikeys.Create(ctx, db, now, in)` using a `now` far in the past for expired keys (or a
  direct INSERT with `HashKey`); `Create` returns the plaintext key for the request header.
- Cypress: `cy.login()` (`web/cypress/support/commands.ts:38-63`) caches sessions with `cy.session`.
  `cy.request` sends the cookies, so call `cy.clearCookies()` for bearer-only checks. The E2E runner seeds
  `tools/database/demo-users.sql`; `anna@sithub.local` and `sam@sithub.local` are admins, and the password
  is `SitHubDemo2026!!`. Run against the dev server (`npm run dev`, backend on :9900) per
  `.claude/rules/cypress.md`; the Vite proxy forwards the `Authorization` header unchanged.

### Project rules checklist

- Go: sentinel errors (`ErrInvalidAPIKey`, `ErrAPIKeyExpired`) checked with `errors.Is`; `%w` at the end;
  security-critical functions documented; no package-name stutter (`apikeys.Authenticate`, not
  `apikeys.AuthenticateAPIKey`). Terminal `api.Write*` calls in top-level funcs need
  `//nolint:wrapcheck // Terminal response` (commit 74529fb pattern).
- JSON:API errors with `application/vnd.api+json` via `internal/api` helpers; do not create new error
  structs (`.claude/rules/golang.md#Shared Types`).
- OpenAPI 3.1, split files, redocly lint (`.claude/rules/apidoc.md`).
- Markdown rules for `docs/public-api.md` (`.claude/rules/markdown.md`); lint with
  `npx --yes markdownlint-cli`.
- No frontend UI change is required in this story. Key management UI belongs to 37.9.

### Open questions / assumptions

#### 1. Scope identifier

The stored scope string is defined by 37.9 as `apikeys.ScopeReadBookings = "read-bookings"`. Reuse the
constant; do not introduce a second literal. The docs show the label "Read bookings".

#### 2. Invalid bearer with a valid cookie

Decision taken here: respond 401 and do not fall back to the cookie. Silent fallback would hide broken
integrations and blur which principal acted.

#### 3. Non-Bearer Authorization schemes

Assumption: ignore them (for example `Basic` added by a reverse proxy) and continue with cookie auth.
Rejecting them would break deployments behind basic-auth proxies.

#### 4. Bearer evaluation scope

Assumption: only `/api/` paths evaluate the header. SPA and static routes ignore it.

#### 5. Rate limiting failed bearer attempts

Recommended default: none in this story. UUID v4 keys have 122 random bits, and the existing limiter
counts all requests, which would throttle legitimate paginated exports. Revisit if logs show abuse.

#### 6. Key lifetime versus creator

Assumption: a key stays valid if the admin who created it is later demoted or deleted. `created_by` is
informational only. Revocation is the explicit control.

#### 7. Usage tracking

There is no `last_used_at` column in the 37.9 schema. Out of scope; a later story can add it.

#### 8. Names from earlier stories

Names above are taken from the 37.6/37.7/37.9 story files (not yet merged). If the merged code differs,
follow the code.

### Project Structure Notes

- New: `internal/auth/api_key.go`, `internal/middleware/api_key.go` (+ tests),
  `web/cypress/e2e/api-key-auth.cy.ts`, `docs/public-api.md`.
- Modified: `internal/middleware/session.go`, `internal/middleware/auth.go`,
  `internal/middleware/require_admin.go`, `internal/middleware/require_data_exporter.go` (37.6), `internal/api/errors.go`
  (`WriteUnauthorizedDetail`), the 37.9 api-keys package (`Authenticate`), `internal/startup/server.go`
  (middleware helper, export routes), `internal/startup/server_test.go`, `api-doc/openapi.yaml`,
  `api-doc/endpoints/bookings-export*.yaml`, `docs/index.md`.
- No new migration. The `api_keys` table comes from 37.9 (`000006`).
- No new Go module dependency. Use `crypto/sha256`, `crypto/subtle`, and the UUID package that 37.9
  already uses for generation.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.10, #Story 37.6-37.9, FR189-FR194]
- [Source: internal/startup/server.go:40-82,224-312]
- [Source: internal/middleware/session.go:10-23], [internal/middleware/auth.go:12-28],
  [internal/middleware/require_admin.go:11-24], [internal/middleware/access.go:13-49],
  [internal/middleware/rate_limit.go:12-65]
- [Source: internal/auth/me.go:18-40], [internal/auth/service.go:30-48,222-251],
  [internal/auth/handlers.go:105]
- [Source: internal/api/errors.go] (`WriteUnauthorized`, `WriteForbiddenDetail`)
- [Source: internal/livefeed/handler.go:49-70]
- [Source: internal/middleware/helpers_test.go], [internal/middleware/security_session_test.go],
  [internal/startup/server_test.go:324-470]
- [Source: _bmad-output/implementation-artifacts/34-5-security-regression-test-coverage.md]
- [Source: api-doc/openapi.yaml], [web/cypress/support/commands.ts:38-63], [run-all-tests.sh:70-145],
  [tools/database/demo-users.sql]
- [Source: .claude/rules/golang.md, json-api.md, apidoc.md, markdown.md, cypress.md]
- RFC 6750 (Bearer Token Usage), section 3 `WWW-Authenticate`

## Dev Agent Record

### Agent Model Used

### Debug Log References

- Red/green tests verified principal helpers, Authenticate, LoadAPIKey and composite middleware before implementation.
- RequireAdmin and RequireDataExporter share requireRole in the actual code; its key rejection supplies both gates.
  RequireExportAccess takes the auth service and scope; it does not need an unused database parameter.
- Router-level regressions live in internal/startup/api_key_test.go to keep the large server_test.go manageable.
  setupTestRouter now uses the same registerAuthMiddleware helper as production with a migrated database.
- Authentication failures log fixed reasons, route templates and remote IP only; no headers or query strings.
  Storage failures fail closed with a generic 500 and do not expose driver messages. No header logger is registered.
- Manual curl verification used a real key from POST /api/v1/api-keys on the isolated local backend.
  The Authorization header was passed via curl's stdin configuration to avoid exposing the key in commands/logs.
  Equivalent commands with the key redacted:

  ```bash
  curl -g --fail-with-body -H 'Authorization: Bearer <redacted>' \
    'http://127.0.0.1:9900/api/v1/bookings/export?from=2026-10-01&to=2026-10-08&page[limit]=50&page[offset]=0'
  curl --fail-with-body -H 'Authorization: Bearer <redacted>' \
    'http://127.0.0.1:9900/api/v1/bookings/export.csv?from=2026-10-01&to=2026-10-08&area=test_area'
  ```

  Both returned 200; JSON data/meta/links and limit=50 verified, CSV BOM and header verified. Key revoked afterward.
- Session temporary logs: bearer-e2e.log, bearer-e2e-full.log, bearer-vitest.log and bearer-server.log.

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created.

- Separate machine principals, exact-string hashing, constant-time verification, inclusive expiry rejection,
  bearer-before-cookie middleware and explicit export scope access implemented. Other user routes deny keys.
- Malformed/unknown/revoked/expired keys, exact expiry instant, uppercase variants, storage failure, duplicate
  Authorization headers, invalid bearer plus admin cookie, route sweep and log redaction have regression tests.
- Both export handlers are already identity-agnostic; they need no user dereference changes.
- Expired credentials cannot be created via the API (minimum 3 months); Go tests cover expiry and its boundary.
- Public integration guide, index link and OpenAPI bearer OR cookie schemes document real behavior and examples.
- Full Go tests, vet and golangci-lint pass. Race tests pass for apikeys, auth, middleware and startup.
  Frontend type-check and lint pass; all 715 Vitest tests pass. Full Electron suite: 112/112 across 25 specs.
- Go duplication is 1.31%, below Task 11's 3% gate. Redocly passes with the pre-existing auth-providers 4XX warning.
- Markdown lint and diff whitespace checks pass. Definition of Done passes; story and sprint are marked review.

- Follow-up (2026-10-11): `middleware.RequireAuthOrScope` (shared `keyOrSession` helper with
  `RequireExportAccess`) guards `GET /api/v1/avatars/:user_id`. Other keys get 403 with the
  `insufficient_scope` challenge; `TestBearerRouteSweep` allow-lists only this route, and `TestBearerAvatar`
  covers 200, 404 for a missing photo, 401 and the still-denied avatar upload/delete. `docs/public-api.md`
  and `api-doc/endpoints/avatars.yaml` document it.

### File List

- `_bmad-output/implementation-artifacts/37-10-bearer-authentication-with-api-keys.md`
- `_bmad-output/implementation-artifacts/sprint-status.yaml`
- `internal/auth/api_key.go`
- `internal/auth/api_key_test.go`
- `internal/apikeys/authenticate.go`
- `internal/apikeys/authenticate_test.go`
- `internal/api/errors.go`
- `internal/middleware/api_key.go`
- `internal/middleware/api_key_test.go`
- `internal/middleware/api_key_logging_test.go`
- `internal/middleware/export_access.go`
- `internal/middleware/export_access_test.go`
- `internal/middleware/session.go`
- `internal/middleware/auth.go`
- `internal/middleware/require_role.go`
- `internal/startup/server.go`
- `internal/startup/server_test.go`
- `internal/startup/api_key_test.go`
- `web/cypress/e2e/api-key-auth.cy.ts`
- `docs/public-api.md`
- `docs/index.md`
- `api-doc/openapi.yaml`
- `api-doc/endpoints/bookings-export.yaml`
- `api-doc/endpoints/bookings-export-csv.yaml`
- `api-doc/endpoints/api-keys.yaml`
- `api-doc/endpoints/api-key.yaml`
- `_bmad-output/implementation-artifacts/deferred-work.md`

## Change Log

- 2026-10-08: Implemented scoped bearer authentication, default-deny authorization, security tests,
  real-backend E2E coverage and integration/OpenAPI documentation.
- 2026-10-08: Completed all quality gates and acceptance criteria; marked ready for code review.
- 2026-10-08: Code review: fixed 12 findings (bearer scheme detection, store-failure logging,
  insufficient_scope challenge, stronger route sweep and precedence tests, docs/OpenAPI accuracy);
  7 deferred to deferred-work.md. Status set to done.
