# Story 37.9: API Key Management

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an administrator,
I want to create, list, and revoke API keys,
so that external systems can access SitHub without personal accounts.

**FRs covered:** FR189, FR190

## Acceptance Criteria

1. Migration `000006` adds an `api_keys` table (id, name, key_prefix, key_hash, scope, expires_at,
   created_at, created_by). It applies cleanly when the server starts, and the down migration drops it.
2. An administrator sees an "API Keys" entry in the user menu (desktop) and the mobile drawer. Non-admins
   do not see it and are redirected to Access Denied when they open the route directly
   (`meta.requiresAdmin`).
3. The API Keys page shows a table with name, the first 10 characters of each key (prefix), expiry date,
   and scope.
4. "Add key" opens a dialog with a name (required), an expiry dropdown (3 months, 12 months, 24 months),
   and a scope selector that offers only "Read bookings".
5. Submitting the dialog generates a random UUID v4 key server-side. The full key is shown exactly once
   with a copy button and a warning that it cannot be shown again. Only its SHA-256 hash and its
   10-character prefix are stored.
6. Revoking a key asks for confirmation; on confirm, the key is deleted and no longer listed.
7. All `/api/v1/api-keys` endpoints return 403 for authenticated non-admins (401 without a session).
   Responses never contain the hash; the full key appears only in the create (201) response.

## Tasks / Subtasks

- [x] Task 1: Migration `000006_api_keys` (AC: #1)
  - [x] Create `internal/db/migrations/000006_api_keys.up.sql` and `.down.sql` (schema in Dev Notes).
        Migrations are embedded via `//go:embed migrations/*.sql` in `internal/db/migrate.go`; no
        registration code is needed.
  - [x] Confirm `000004` (37.5) and `000005` (37.6) already exist. If either is missing, stop and
        flag it instead of shifting numbers; never reuse a migration number.
  - [x] Extend `internal/db/migrate_test.go`: after `RunMigrations`, insert and read back an
        `api_keys` row; assert the `UNIQUE` index on `key_hash` rejects a duplicate.
- [x] Task 2: Store layer `internal/apikeys/store.go` (AC: #1, #5, #6)
  - [x] New package `apikeys` (doc comment `// Package apikeys ...`), modelled on
        `internal/floorplanpos/store.go`: package-level funcs taking `ctx` and `*sql.DB`, sentinel
        `ErrNotFound`, `%w` wrapping at the end of format strings.
  - [x] Types: `Key` (ID, Name, KeyPrefix, Scope, ExpiresAt, CreatedAt, CreatedBy; no hash field),
        `CreateInput` (Name, Scope, ExpiresInMonths, CreatedBy).
  - [x] Constants: `ScopeReadBookings = "read-bookings"`, `PrefixLength = 10`, and the allowed
        expiry set `{3, 12, 24}`.
  - [x] Sentinel errors: `ErrNotFound`, `ErrInvalidName`, `ErrInvalidScope`, `ErrInvalidExpiry`.
  - [x] `HashKey(key string) string`: lowercase hex of `sha256.Sum256([]byte(key))`. Exported so
        Story 37.10 reuses the exact same hashing for bearer lookups.
  - [x] `Create(ctx, db, now time.Time, in CreateInput) (Key, string, error)`: validate, generate the
        key with `uuid.NewString()` (v4, crypto/rand), compute prefix, hash, and `expires_at`, insert,
        and return the record plus the plaintext key. Pass `now` in for deterministic tests.
  - [x] `List(ctx, db) ([]Key, error)`: `ORDER BY created_at DESC, name`; never selects `key_hash`.
  - [x] `Delete(ctx, db, id) error`: hard delete; `ErrNotFound` when no row was affected.
- [x] Task 3: HTTP handlers `internal/apikeys/handler.go` (AC: #5, #6, #7)
  - [x] `toResource(k Key) api.Resource` with type `api-keys` and snake_case attributes (see Dev
        Notes). Use `api.MapResources`, `api.WriteCollection`, `api.WriteSingle`.
  - [x] `ListHandler(db)`: `GET /api/v1/api-keys`, 200 collection.
  - [x] `CreateHandler(db)`: `POST /api/v1/api-keys`. Decode the JSON:API body, map validation
        sentinels to `api.WriteBadRequest` with a precise detail (shown by the Story 37.1 error UI),
        read the creator from `auth.GetUserFromContext(c).ID`, and respond 201 with the resource plus
        the one-time `key` attribute.
  - [x] `DeleteHandler(db)`: `DELETE /api/v1/api-keys/:id`, 204; 404 via `api.WriteNotFound`.
  - [x] Never log the plaintext key or the hash; do not pass them into error labels.
  - [x] Terminal helper calls in top-level funcs need `//nolint:wrapcheck // Terminal response`
        where golangci-lint asks for it (pattern from commit 74529fb).
- [x] Task 4: Route registration (AC: #7)
  - [x] In `registerRoutes` (`internal/startup/server.go`, admin block around lines 295-311), add
        `GET`, `POST /api/v1/api-keys` and `DELETE /api/v1/api-keys/:id`, each with
        `requireAuth, requireAdmin`. No trailing-slash variants.
  - [x] Add route tests to `internal/startup/server_test.go` next to
        `TestFloorPlanPositionsWriteRouteRequiresAdmin`: 401 without cookie, 403 with a non-admin
        cookie for all three routes (table-driven).
- [x] Task 5: OpenAPI documentation (AC: #7)
  - [x] Add `api-doc/endpoints/api-keys.yaml` (GET, POST) and `api-doc/endpoints/api-key.yaml`
        (DELETE), reference them in `api-doc/openapi.yaml` under `/api-keys` and `/api-keys/{id}`.
  - [x] Add component schemas `ApiKeyAttributes`, `ApiKeyCreateRequest`, `ApiKeyCollectionResponse`,
        `ApiKeyCreatedResponse` (with `key`, documented as returned only once). Document 400, 401,
        403, 404 with `ErrorResponse`, and the `scope` enum and `expires_in_months` enum.
  - [x] Run `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml` until clean. The bearer
        security scheme belongs to Story 37.10; do not add it here.
- [x] Task 6: Frontend API module `web/src/api/apiKeys.ts` (AC: #3, #5, #6)
  - [x] Follow `web/src/api/floorPlanPositions.ts`: `ApiKeyAttributes`, `CreatedApiKeyAttributes`
        (adds `key: string`), `ApiKeyScope = 'read-bookings'`, `ApiKeyExpiryMonths = 3 | 12 | 24`.
  - [x] `fetchApiKeys()`, `createApiKey({ name, scope, expires_in_months })`, `deleteApiKey(id)` using
        `apiRequest` and the `CollectionResponse`/`SingleResponse` types from `web/src/api/types.ts`.
  - [x] `deleteApiKey` must not call `response.json()` on 204. Reuse whatever
        `deleteFloorPlanPosition` relies on (`apiRequest<void>`) and verify it handles empty bodies;
        fix in `client.ts` only if it does not.
- [x] Task 7: API Keys view `web/src/views/ApiKeysView.vue` (AC: #3, #4, #5, #6)
  - [x] `<script setup lang="ts">`, Composition API, `PageHeader` with title "API Keys" and an
        "Add key" button (`data-cy="api-keys-add"`).
  - [x] `v-data-table` (client-side, not server-paginated; see assumptions) with columns name,
        key prefix (monospace, rendered as `<prefix>…`), expiry (localized date, red "Expired" chip
        when `expires_at` is in the past), scope (translated label), and a revoke icon button.
        Selectors: `api-keys-table`, `api-key-row`, `api-key-name`, `api-key-prefix`,
        `api-key-expiry`, `api-key-scope`, `api-key-revoke`.
  - [x] Empty state via the existing `EmptyState` component (`data-cy="api-keys-empty"`); load failure
        as persistent `v-alert` (`data-cy="api-keys-load-error"`) per `.claude/rules/feedback.md`.
  - [x] Create dialog (`data-cy="api-key-dialog"`), new component
        `web/src/components/ApiKeyCreateDialog.vue`. Form state: name text field (required, trimmed,
        max 100 chars, `data-cy="api-key-name-input"`), expiry `v-select` with 3/12/24 months, default
        12 (`data-cy="api-key-expiry-select"`), scope `v-select` with the single "Read bookings" option,
        preselected (`data-cy="api-key-scope-select"`). Submit (`data-cy="api-key-submit"`) is
        disabled while the name is blank or a request is in flight.
  - [x] After a 201, switch the same dialog to its "created" state: readonly monospace field with the
        full key (`data-cy="api-key-value"`), copy button (`data-cy="api-key-copy"`), and a
        `v-alert type="warning"` (`data-cy="api-key-once-warning"`) saying the key cannot be shown
        again. Only a "Done" button closes it (`persistent` dialog, no outside-click close). On close,
        clear the key ref and reload the list.
  - [x] Copy uses `navigator.clipboard.writeText`. Success snackbar "Key copied" (3000 ms). If the
        Clipboard API is missing or rejects (plain-http deployments), show a warning snackbar asking
        the user to copy manually and select the field text.
  - [x] Revoke uses the existing `ConfirmDialog.vue` (`confirmColor="error"`, message names the key).
        On success, snackbar "API key revoked" (3000 ms) and reload; on failure, error snackbar
        (6000 ms, `closable`) with the backend detail via the Story 37.1 helper.
  - [x] One page-level `v-snackbar` (`data-cy="api-keys-snackbar"`, `location="bottom"`) driven by
        color/timeout per `.claude/rules/feedback.md`.
- [x] Task 8: Route, menu, icons (AC: #2)
  - [x] `web/src/router/index.ts`: add `{ path: '/admin/api-keys', name: 'api-keys',
        component: () => import('../views/ApiKeysView.vue'), meta: { requiresAdmin: true } }`. The
        existing `beforeEach` guard redirects non-admins to `access-denied`.
  - [x] `web/src/App.vue`: add a `v-list-item` gated by `v-if="authStore.isAdmin"` in the user menu
        (next to the floor-plan editor entry and the 37.4 "Users" entry, `data-cy="api-keys-btn"`) and
        in the mobile drawer (`data-cy="mobile-api-keys-btn"`, `@click="mobileDrawer = false"`).
  - [x] `web/src/plugins/vuetify.ts`: add aliases `key: mdiKeyVariant` and `copy: mdiContentCopy`
        (import from `@mdi/js` like the others). Use `$key` for the menu entry and `$copy` for copy.
- [x] Task 9: i18n (AC: #2-#6)
  - [x] Add `app.userMenu.apiKeys` and a new `apiKeys` block (title, add, columns, expiry options,
        scope label, dialog texts, once-warning, copy/copied/copyManual, revoke title/message/success,
        expired chip, empty state) to all five locales: `en`, `de`, `uk`, `fr`, `es`
        (`web/src/locales/*.json`). German terms: "API-Schlüssel" for the menu and page title and
        "Buchungen lesen" for the scope.
- [x] Task 10: Tests (AC: #1-#7)
  - [x] Go store tests `internal/apikeys/store_test.go` (table-driven, `t.Parallel()`, testify):
        valid create returns a v4 UUID (`uuid.Parse` + `Version() == 4`); stored `key_hash` equals
        `HashKey(key)`; stored `key_prefix == key[:10]`; no column of the row contains the plaintext
        key; expiry for 3/12/24 months from a fixed `now`; invalid expiry (0, 1, 6, 36, -3) and
        invalid scope and blank/whitespace/101-char names return the matching sentinel; `List`
        ordering; `Delete` happy path and `ErrNotFound`.
  - [x] Go handler tests `internal/apikeys/handler_test.go`: 201 response has `key` and no
        `key_hash`; list response has neither `key` nor `key_hash` (assert on the raw JSON body);
        400 on malformed body and each validation error; 204/404 on delete; `created_by` equals the
        context user ID. Set the user with `c.Set("user", &auth.User{...})`.
  - [x] Route tests in `internal/startup/server_test.go` (Task 4).
  - [x] Vitest: `web/src/api/apiKeys.test.ts` (request URLs, methods, JSON:API body);
        `web/src/views/ApiKeysView.test.ts` (table rendering, expired chip, dialog validation, created
        state shows key and warning, copy success and fallback, close clears key, revoke confirm calls
        `deleteApiKey` and reloads, error snackbar shows detail); extend `web/src/App.test.ts` for the
        admin-only menu entries (desktop and mobile). Coverage stays at 75% or higher.
  - [x] Cypress E2E `web/cypress/e2e/api-keys.cy.ts` against the dev server, no mocked responses,
        intercept aliases `@listApiKeys`, `@createApiKey`, `@deleteApiKey`:
        should show the API Keys entry to admins and open the page; should create a key and show it
        once (UUID v4 regex, warning visible, copy button present); should list the new key with only
        its 10-character prefix (prefix equals the first 10 chars of the shown key, full key not in
        the page after closing); should revoke a key after confirmation; should hide the entry and
        show Access Denied for a non-admin (`alex@sithub.local`); `cy.request` with
        `failOnStatusCode: false` as non-admin to GET/POST/DELETE returns 403. Use unique names
        (`e2e-key-${Date.now()}`) and delete created keys in `afterEach` via `cy.request`.
- [x] Task 11: Quality gates
  - [x] Backend: `go fmt ./...`, `go vet ./...`, `golangci-lint run ./...`, `go test ./...`,
        `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exit-code 1`.
  - [x] Frontend (in `web/`): `npm run type-check`, `npm run lint`, `npx vitest run --coverage`,
        `npm run build`, `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0
        --exit-code 1`, `npm run test:e2e -- --browser electron`.
  - [x] Redocly lint for `api-doc`, markdownlint for any Markdown touched.

## Dev Notes

### Scope boundary

This story delivers key storage, admin CRUD, and the admin UI only. Bearer authentication
(`Authorization: Bearer <key>`), scope enforcement, constant-time comparison, the OpenAPI bearer
security scheme, and `docs/public-api.md` are Story 37.10. Do not touch `internal/middleware/session.go`
or `auth.go` here. Design the store so 37.10 can add `FindValidByHash(ctx, db, hash, now)` without
schema changes (unique index on `key_hash`, `expires_at` comparable as text).

### Migration `000006_api_keys`

```sql
-- up
CREATE TABLE api_keys (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  key_prefix TEXT NOT NULL,
  key_hash TEXT NOT NULL,
  scope TEXT NOT NULL CHECK (scope IN ('read-bookings')),
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL,
  created_by TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX idx_api_keys_key_hash ON api_keys(key_hash);

-- down
DROP TABLE api_keys;
```

`created_by` holds the creating admin's `users.id` but deliberately has **no** foreign key.
`PRAGMA foreign_keys=ON` is set in `internal/db/db.go:24`, and Story 37.5 hard-deletes users without
bookings; a FK would make deleting such an admin fail or cascade-delete keys. Timestamps follow the
existing TEXT/RFC 3339 convention (`internal/db/migrations/000001_init.up.sql`,
`floorplanpos.Create`). The `CHECK` on `scope` is cheap to extend later via a table rebuild; if that
is considered too rigid, drop the CHECK and rely on Go validation (keep the Go validation either way).

### Key generation, prefix, hash

- Key: `uuid.NewString()` from `github.com/google/uuid` (already in `go.mod`, used by
  `floorplanpos.Create`). It produces a random RFC 4122 version 4 UUID from `crypto/rand`, canonical
  lowercase form with dashes, e.g. `3f2b8c1e-9a4d-4c7e-b1f0-2d6e8a9c0b13`.
- Prefix: `key[:10]`, i.e. 8 hex chars, the dash, and 1 hex char (`3f2b8c1e-9`). This is what
  "first 10 characters" literally means; keep it, do not strip dashes.
- Hash: `hex.EncodeToString(sha256.Sum256([]byte(key)))` over the exact string returned to the client.
  No salt and no bcrypt: the key has 122 bits of randomness, and 37.10 needs a deterministic lookup by
  hash. Store lowercase hex (64 chars).
- The plaintext key exists only in the `Create` return value and the 201 response body. It must never
  be stored, logged, put into `slog` attributes, or kept in Pinia/localStorage on the frontend.

### Expiry computation

Request attribute `expires_in_months` (integer, one of 3, 12, 24; anything else is 400). Server
computes `expires_at = now.UTC().AddDate(0, months, 0).Format(time.RFC3339)`, e.g. created
`2026-10-08T09:15:00Z` with 12 months gives `2027-10-08T09:15:00Z`. Go's `AddDate` normalizes overflow
(created on 2026-11-30 plus 3 months gives 2027-03-02 because February has no 30th); accept that and
cover it in a table test so the behavior is explicit. The client never sends an absolute date.

### API contract

`GET /api/v1/api-keys` returns 200:

```json
{
  "data": [
    {
      "type": "api-keys",
      "id": "8c0e…",
      "attributes": {
        "name": "Finance export",
        "key_prefix": "3f2b8c1e-9",
        "scope": "read-bookings",
        "expires_at": "2027-10-08T09:15:00Z",
        "created_at": "2026-10-08T09:15:00Z",
        "created_by": "user-uuid"
      }
    }
  ]
}
```

`POST /api/v1/api-keys` request:

```json
{
  "data": {
    "type": "api-keys",
    "attributes": { "name": "Finance export", "scope": "read-bookings", "expires_in_months": 12 }
  }
}
```

Response 201: the same resource shape plus `"key": "3f2b8c1e-9a4d-4c7e-b1f0-2d6e8a9c0b13"` in
`attributes`. Use a dedicated attribute map for create so `toResource` (used by list) can never leak
`key`. `DELETE /api/v1/api-keys/:id` returns 204 with no body, 404 JSON:API error if unknown.

Validation details (English, passed through untranslated by the 37.1 error UI):
"name is required", "name must be at most 100 characters", "scope must be read-bookings",
"expires_in_months must be one of 3, 12, 24". Reject a `data.type` other than `api-keys` with 400 as
well. Duplicate names are allowed (the prefix disambiguates).

### Backend patterns to follow

- Package layout and style: `internal/floorplanpos/store.go` and `handler.go` (handler factories
  taking `*sql.DB`, `toResource`, sentinel `ErrNotFound`, `RowsAffected` check on delete).
- Response helpers: `internal/api/write.go` (`WriteCollection`, `WriteSingle`, `WriteInternalError`) and
  `internal/api/errors.go` (`WriteBadRequest`, `WriteNotFound`). These return `nil` after writing; do
  not continue after calling them.
- Admin gate: `middleware.RequireAdmin()` (`internal/middleware/require_admin.go`) returns 401 when no
  user is in context and 403 for non-admins; wire it exactly like the `/api/v1/users` routes in
  `internal/startup/server.go:296-301`.
- Creator: `auth.GetUserFromContext(c)` (`internal/auth/me.go:20`); after `RequireAuth` it is non-nil,
  but guard anyway and return `api.WriteUnauthorized` if nil.
- Document thread safety on exported types per `.claude/rules/golang.md` (package funcs are stateless
  and safe for concurrent use).
- The list is small; no pagination. Do not add `page[...]` parsing here.

### Frontend patterns to follow

- API module style: `web/src/api/floorPlanPositions.ts`; client `web/src/api/client.ts` (`ApiError`
  with `status` and `detail`).
- Admin page precedent: `web/src/views/FloorPlanEditorView.vue` route `/admin/floor-plan-editor` with
  `meta: { requiresAdmin: true }` (`web/src/router/index.ts:62-67`, guard at ~96). Story 37.4 adds
  the "Users" admin page and menu entry; place "API Keys" directly after it and mirror its markup.
- Menu: user `v-menu` in `web/src/App.vue` (~40-110, floor-plan editor entry ~67-76) and mobile drawer
  (~163-200). Both gate on `authStore.isAdmin`.
- Confirmations: `web/src/components/ConfirmDialog.vue` (`confirm-dialog-confirm` /
  `confirm-dialog-cancel` selectors, `loading` prop).
- Feedback: `.claude/rules/feedback.md` (snackbar bottom; success 3000 ms, error 6000 ms + closable,
  warning 4000 ms; `v-alert` only for persistent page states).
- Error text: use the Story 37.1 helper that shows a translated prefix plus `ApiError.detail`; do not
  add a new generic "please try again" message.
- Dates: format `expires_at` with the existing locale-aware date formatting used in
  `MyBookingsView.vue` / `BookingHistoryView.vue`; never parse it as local time without the `Z`.

### Security notes

- AC #7 is enforced by middleware, not by hiding UI; the route tests are mandatory.
- Never echo the hash. Assert this in handler tests on the raw response body (`assert.NotContains`).
- The created-state dialog is `persistent` so an accidental outside click does not lose the key
  before it was copied.
- Revocation is a hard delete; 37.10 then rejects the key because no hash matches.

### Open questions / assumptions

- Scope string: `read-bookings` (kebab-case, matching resource-type style). UI label "Read bookings" /
  "Buchungen lesen". Recommended default; 37.10 must use the same constant
  (`apikeys.ScopeReadBookings`).
- Expired keys stay listed with an "Expired" chip until an admin revokes them. No automatic cleanup.
- No pagination for the key list (expected to stay well below 100 rows).
- `created_by` is stored but not shown in the table (AC lists only name, prefix, expiry, scope);
  `created_at` is returned by the API for sorting and may be shown as a tooltip.
- Default expiry in the dialog is 12 months.
- Name: required, trimmed, max 100 characters, duplicates allowed.
- No FK on `created_by` (see Migration section).

### Project Structure Notes

- New backend: `internal/apikeys/{store.go,handler.go,store_test.go,handler_test.go}`,
  `internal/db/migrations/000006_api_keys.{up,down}.sql`.
- Modified backend: `internal/startup/server.go`, `internal/startup/server_test.go`,
  `internal/db/migrate_test.go`.
- New frontend: `web/src/api/apiKeys.ts` (+ test), `web/src/views/ApiKeysView.vue` (+ test),
  `web/src/components/ApiKeyCreateDialog.vue`, `web/cypress/e2e/api-keys.cy.ts`.
- Modified frontend: `web/src/router/index.ts`, `web/src/App.vue`, `web/src/App.test.ts`,
  `web/src/plugins/vuetify.ts`, `web/src/locales/{en,de,uk,fr,es}.json`.
- API docs: `api-doc/openapi.yaml`, `api-doc/endpoints/api-keys.yaml`, `api-doc/endpoints/api-key.yaml`.

### Testing standards summary

Go: table-driven tests, `t.Parallel()`, `require` for setup and `assert` for behavior, inject `now` for
expiry. Vitest coverage at 75% or higher. Cypress E2E covers every AC against the dev server with
`data-cy` selectors and intercepts only for waiting and assertions
[Source: .claude/rules/golang.md, .claude/rules/vue.md, .claude/rules/cypress.md].

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.9 (FR189, FR190), lines 5984-6023]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.10, lines 6025-6070 (out of scope here)]
- [Source: _bmad-output/planning-artifacts/epics.md, FR189-FR191, lines 700-711]
- [Source: internal/floorplanpos/store.go, internal/floorplanpos/handler.go]
- [Source: internal/api/write.go, internal/api/errors.go]
- [Source: internal/middleware/require_admin.go]
- [Source: internal/startup/server.go:224-228, 295-311]
- [Source: internal/startup/server_test.go:383-445]
- [Source: internal/db/db.go:21-25, internal/db/migrate.go:15-19, internal/db/migrate_test.go]
- [Source: internal/auth/me.go:20, internal/auth/service.go:41-48]
- [Source: web/src/App.vue:40-110, 163-200]
- [Source: web/src/router/index.ts:62-67, 77-100]
- [Source: web/src/api/floorPlanPositions.ts, web/src/api/client.ts, web/src/api/types.ts]
- [Source: web/src/components/ConfirmDialog.vue]
- [Source: web/src/plugins/vuetify.ts:120-170]
- [Source: .claude/rules/json-api.md, .claude/rules/apidoc.md, .claude/rules/feedback.md]

## Tasks / Subtasks (continued)

### Review Findings

- [x] [Review][Patch] German revoke-success snackbar is the infinitive, not a confirmation
  `[web/src/locales/de.json]` — fixed: "API-Schlüssel wurde widerrufen"
- [x] [Review][Patch] Expiry column renders time-of-day instead of the date-only convention referenced
  in Dev Notes `[web/src/views/ApiKeysView.vue]` — fixed: formatMediumDate on the RFC 3339 date part
- [x] [Review][Patch] Expired chip goes stale for up to 60 s past the expiry instant (minute tick +
  strict `<`) `[web/src/views/ApiKeysView.vue]` — fixed: isExpired uses `<=` and a documented
  minute-tick comparison
- [x] [Review][Patch] Hardcoded English "Connection to server lost" rendered in every locale
  `[web/src/api/client.ts]` — resolved as patch: common.connectionLost in all 5 locales, translator
  registered by the i18n plugin (English constant kept as fallback)
- [x] [Review][Patch] Key created but plaintext silently lost when the admin navigates away mid-create
  `[web/src/components/ApiKeyCreateDialog.vue]` — resolved as patch: the list always reloads on dialog
  close, so any committed mid-flight create is visible and revocable
- [x] [Review][Defer] Revoke erases all evidence a key existed (no tombstone/audit trail)
  `[internal/apikeys/store.go]` — deferred, deferred to story 37.10 design

## Dev Agent Record

### Agent Model Used

### Debug Log References

- Migration, store, handlers, route authorization and frontend tests were written first and failed
  before their implementation. Existing migrations 000004 and 000005 were confirmed before adding 000006.
- Implementation uses stateless storage functions, a separate public metadata map and a create-only key attribute.
  The creation response uses Cache-Control: no-store. UUID record IDs are separate from plaintext keys.
- Existing apiRequest always parsed JSON, including empty 204 responses. Added explicit 204 handling,
  verified through the real client with a Response whose json method must not be called.
- Route registration is in a small helper to keep registerRoutes within the existing function-length limit.
- Mobile E2E originally opened the drawer before initial navigation completed; Vuetify then closed it.
  The test now waits for the real areas response before opening the drawer. Targeted rerun: 4/4 pass.
- Literal Task 11 duplication gates fail on existing code: Go 12 clones, 1.31%; TypeScript 72 clones, 5.7%.
  No API-key clone is reported. The existing project TypeScript gate (excluding unit tests) passes at 0%.
  User approved applying the same project gates as Story 37.8: Go <=3% excluding tests/node_modules,
  TypeScript 0% excluding unit tests/node_modules. This supersedes Task 11's literal duplication commands.
- Approved gates pass: Go 1.31%, TypeScript 0%. Go uses --threshold 3 without --exit-code 1,
  matching the project gate; the latter option fails on any clone regardless of the threshold.
- Validation logs are in the session temporary directory, named api-keys-coverage.log,
  api-keys-e2e-full.log, api-keys-e2e-final.log, api-keys-e2e.log, api-keys-go-dup.log and api-keys-ts-dup.log.

### Completion Notes List

- All tasks complete: hashed storage, admin-only JSON:API CRUD, OpenAPI schemas, admin menu and route,
  localized key table, create-once dialog with clipboard fallback, confirmation-based revocation, five locales.
- Full Go test suite, go fmt, go vet and golangci-lint pass. Frontend type-check, ESLint and build pass.
- All 715 Vitest tests pass. Coverage: statements 82.64%, branches 75.04%, functions 77.30%, lines 84.24%.
- Final full Electron regression run: 108/108 passed across 24 specs against the isolated Vite dev server
  and real backend, including all four API-key tests. The earlier mobile test failure is resolved.
- Redocly exits successfully with its existing auth-providers missing-4XX warning; Markdown lint passes.
- Definition of Done passes with the user-approved duplication gates. All acceptance criteria and tasks
  are satisfied; story and sprint status are review.

### File List

- `_bmad-output/implementation-artifacts/37-9-api-key-management.md`
- `_bmad-output/implementation-artifacts/sprint-status.yaml`
- `internal/db/migrations/000006_api_keys.up.sql`
- `internal/db/migrations/000006_api_keys.down.sql`
- `internal/db/migrate_test.go`
- `internal/apikeys/store.go`
- `internal/apikeys/store_test.go`
- `internal/apikeys/handler.go`
- `internal/apikeys/handler_test.go`
- `internal/startup/server.go`
- `internal/startup/server_test.go`
- `api-doc/openapi.yaml`
- `api-doc/endpoints/api-keys.yaml`
- `api-doc/endpoints/api-key.yaml`
- `web/src/api/apiKeys.ts`
- `web/src/api/apiKeys.test.ts`
- `web/src/api/client.ts`
- `web/src/views/ApiKeysView.vue`
- `web/src/views/ApiKeysView.test.ts`
- `web/src/components/ApiKeyCreateDialog.vue`
- `web/src/App.vue`
- `web/src/App.test.ts`
- `web/src/router/index.ts`
- `web/src/router/index.test.ts`
- `web/src/plugins/vuetify.ts`
- `web/src/locales/en.json`
- `web/src/locales/de.json`
- `web/src/locales/fr.json`
- `web/src/locales/es.json`
- `web/src/locales/uk.json`
- `web/cypress/e2e/api-keys.cy.ts`

## Change Log

- 2026-10-08: Applied user-approved project duplication gates; both pass. Final full browser suite
  passes 108/108. Completed Task 11 and marked story and sprint status review.
- 2026-10-08: Code review (Blind Hunter, Edge Case Hunter, Acceptance Auditor) — 5 patches applied
  (German revoke message, date-only expiry, expired-chip boundary, translated connection-lost message,
  always-reload on dialog close), 1 finding deferred to story 37.10 (audit trail), 2 dismissed as noise.
  Story and sprint status set to done.
- 2026-10-08: Implemented API key management and regression coverage. Tasks 1-10 complete; Task 11
  blocked by existing repository duplication under the story's literal zero-duplication commands.
