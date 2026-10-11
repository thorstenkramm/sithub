# Deferred Work

## Deferred from: code review of 37-9-api-key-management (2026-10-08)

- `internal/apikeys/store.go`: revoking a key hard-deletes the row, erasing its name, prefix, scope,
  creator and revocation time. A credential-management feature has no audit trail for
  "who had booking-read access, and who revoked which key when". Consider a `revoked_at` tombstone
  or audit log entry in the Story 37.10 design.

## Deferred from: code review of 37-8-export-bookings-as-csv (2026-10-08)

- `web/src/views/DataExportView.vue:218`: failed filter requests retain previous rows and totals beneath
  new filters. Pre-existing in 37.7; clear or explicitly mark stale table results.
- `internal/bookings/export_store.go:65`: COUNT and SELECT can observe different database states.
  Pre-existing in 37.7, related to the pagination concern below; CSV does not use the count.
- `internal/middleware/auth.go:22`: transient database errors from EnsureActive return 401 even though
  the session cookie is retained. Pre-existing in 37.5; distinguish infrastructure errors from inactive users.

## Deferred from: code review of 37-5-delete-orphaned-users (2026-10-08)

- Cookie name `sithub_user` is hardcoded in `internal/middleware/auth.go` (and already in
  `internal/middleware/session.go`); `internal/auth` keeps it as the unexported `userCookieName`.
  Export a shared constant (or a `ClearUserCookie` helper) from `internal/auth` so the three places
  cannot drift. Pre-existing duplication, cosmetic.

## Deferred from: code review of 37-7-data-export-table-with-date-and-area-filters (2026-10-08)

- E2E timezone flake: `web/cypress/e2e/data-export.cy.ts` seeds bookings on the UTC date while the view
  defaults to local-time dates, so the seeded rows fall outside the default range between local midnight
  and UTC midnight in non-UTC timezones. CI runs in UTC, so this never fails there; a real fix needs a
  product decision on local-vs-UTC "today" for booking creation validation.
- `internal/bookings/export_store.go` runs the COUNT and page SELECT without a transaction/snapshot, so
  bookings created or canceled while a data exporter pages can shift rows across page boundaries (same
  trade-off the users list accepted in story 37.4). Consider keyset pagination if finance reports rows
  skipped or duplicated.

## Deferred from: code review of 37-10-bearer-authentication-with-api-keys (2026-10-08)

- No rate limiting or log throttling for failed bearer attempts: every well-formed random UUID on any
  `/api/` path costs a SQLite lookup and a `slog.Warn` (`internal/middleware/api_key.go`). Story open
  question 5 decided "none" because the login limiter would throttle paginated exports; revisit with a
  per-IP limiter on rejections only if logs show abuse.
- `c.RealIP()` trusts client-supplied `X-Forwarded-For`/`X-Real-IP` because `e.IPExtractor` is never
  set, so the `remote_ip` in auth rejection logs is spoofable. Configure `echo.ExtractIPDirect()` or
  `ExtractIPFromXFFHeader` with trusted proxy ranges.
- No `last_used_at` (or success audit beyond a Debug log) for API keys, so admins cannot spot stale or
  abused keys before revoking (story open question 7).
- A corrupt `expires_at` value in `api_keys` is returned as `ErrInvalidAPIKey` and logged as
  `reason=unknown`, hiding a data-integrity problem behind what looks like a guessing attempt
  (`internal/apikeys/authenticate.go`).
- `apikeys.Create` formats timestamps as RFC 3339 without truncating `now` first, so a key expires up to
  one second earlier than the exact calendar-month instant (`internal/apikeys/store.go`, story 37.9).
- OpenAPI `components.parameters.PageLimit` declares `maximum: 100` while the server caps larger values
  to 100, so the schema rejects requests the server accepts (story 37.4).
- 401 responses from `RequireAuth` carry no `WWW-Authenticate` challenge (RFC 9110); only rejected bearer
  credentials do. Consider `WWW-Authenticate: Bearer` on export routes so API clients get a hint.
