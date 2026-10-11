# Story 37.6: Data Exporter Role

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an administrator,
I want to grant the "data exporter" role to selected users,
so that finance can access booking data without becoming administrators.

## Acceptance Criteria

1. Migration `000005` adds `users.is_data_exporter` (INTEGER NOT NULL DEFAULT 0). When the server starts, the
   migration applies cleanly and no existing user gains the role.
2. In the admin Users table (Story 37.4), the role column shows "Admin" and/or "Data exporter" (or "User" when
   neither flag is set), and the administrator can toggle "Data exporter" per user. The toggle calls
   `PATCH /api/v1/users/:id` with the attribute `is_data_exporter`.
3. A user without the flag (including administrators) who requests data-export endpoints or routes gets
   403 (API) or the Access Denied page (frontend route).
4. When the flag is granted or revoked, the change takes effect on the user's next request without a new login,
   because the flag is read from the database on every request, not from the session cookie.
5. `GET /api/v1/me` includes the boolean attribute `is_data_exporter`, so the frontend can gate the menu entry
   (the menu entry itself is added in Story 37.7).

## Tasks / Subtasks

- [x] Task 1: Add migration `000005_add_is_data_exporter` (AC: #1)
  - [x] Create `internal/db/migrations/000005_add_is_data_exporter.up.sql`:
        `ALTER TABLE users ADD COLUMN is_data_exporter INTEGER NOT NULL DEFAULT 0;`
  - [x] Create the matching `.down.sql`: `ALTER TABLE users DROP COLUMN is_data_exporter;`
        (same style as `000003_add_border_width`).
  - [x] Confirm `000004` (Story 37.5, `users.deleted`) exists before this one. If 37.5 is not merged yet,
        stop and sequence the work; never reuse or renumber `000004`.
- [x] Task 2: Persist the flag in the users store (AC: #1, #2, #4)
  - [x] Add `IsDataExporter bool` to `users.Record` (`internal/users/record.go`).
  - [x] Append `is_data_exporter` to `userColumns` and scan it in `scanRow` (int -> bool, like `is_admin`)
        in `internal/users/store.go`.
  - [x] Add `IsDataExporter *bool` to `UpdateFields` and a matching `is_data_exporter = ?` clause in
        `UpdateUser` (hardcoded column name, `boolToInt`).
  - [x] Leave `UpsertEntraIDUser`'s `ON CONFLICT ... DO UPDATE SET` list untouched: an Entra ID login must
        NEVER overwrite `is_data_exporter` (unlike `is_admin`, which is driven by the Entra admin group).
  - [x] `CreateLocalUser` keeps the column default (false); do not add a create-time flag.
- [x] Task 3: Expose and accept the flag in the admin users API (AC: #2)
  - [x] Add `` IsDataExporter bool `json:"is_data_exporter"` `` to `users.UserAttributes` and fill it in
        `recordToAttributes`. Keep `role` unchanged (`admin`/`user`) for backward compatibility.
  - [x] Add `` IsDataExporter *bool `json:"is_data_exporter,omitempty"` `` to the PATCH attributes and map it
        into `UpdateFields` in `applyFieldUpdates`.
  - [x] The PATCH attributes struct is currently declared twice (named `updateUserRequest` and an identical
        anonymous struct in the `applyFieldUpdates` signature). Extract a named `updateUserAttributes` type
        and use it in both places so the new field cannot drift.
- [x] Task 4: Read the flag from the database on every request and expose it via `/api/v1/me` (AC: #4, #5)
  - [x] Add `` IsDataExporter bool `json:"is_data_exporter"` `` to `auth.User` (`internal/auth/service.go`).
  - [x] In the per-request database user check introduced by Story 37.5 (the lookup that rejects
        soft-deleted users with 401), copy `rec.IsDataExporter` onto the context `*auth.User`. The cookie value
        of this field must never be trusted; it is always overwritten from the DB record.
  - [x] Set `IsDataExporter: rec.IsDataExporter` when building the user in `LocalLoginHandler`
        (`internal/auth/login_local.go`) and in `Service.FetchUser` (Entra ID) so the login response is
        already correct.
  - [x] Add `attrIsDataExporter = "is_data_exporter"` to `internal/auth/json_fields.go` and include it in
        `respondWithUserProfile` (`internal/auth/me.go`). This covers `GET /api/v1/me`, `PATCH /api/v1/me`,
        and the local login response.
- [x] Task 5: Add the `RequireDataExporter` middleware (AC: #3)
  - [x] Create `internal/middleware/require_data_exporter.go`, modeled on `RequireAdmin`: nil user -> 401
        via `api.WriteUnauthorized`, `!user.IsDataExporter` -> 403 via `api.WriteForbidden`. Do NOT fall back to
        `IsAdmin`.
  - [x] Document in the doc comment that it must run after `RequireAuth` (which performs the DB refresh).
  - [x] Do not register any data-export endpoint here; Stories 37.7/37.8 mount
        `requireAuth, requireDataExporter` on `/api/v1/bookings/export*`. Create the middleware instance next to
        `requireAdmin` in `registerRoutes` only when a route uses it (unused variables fail the linter).
- [x] Task 6: Update the OpenAPI documentation (AC: #2, #5)
  - [x] `api-doc/openapi.yaml`: add `is_data_exporter: boolean` to `components.schemas.UserAttributes`
        (and to its `required` list) and to `UpdateUserRequest.data.attributes`.
  - [x] Lint with `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.
- [x] Task 7: Frontend auth state and route guard (AC: #3, #5)
  - [x] `web/src/api/me.ts`: add `is_data_exporter: boolean` to `UserAttributes`.
  - [x] `web/src/stores/useAuthStore.ts`: add `isDataExporter` ref, accept `is_data_exporter` in `setUser`,
        reset it in `clearUser`, return it.
  - [x] Update both `setUser` callers: `web/src/router/index.ts` (~84-90) and `web/src/views/LoginView.vue`
        (~143-149).
  - [x] Add a `meta.requiresDataExporter` check to `router.beforeEach` (after the `requiresAdmin` check, ~96)
        that redirects to `access-denied` when `!authStore.isDataExporter`. Admin status must not bypass it.
  - [x] Add the route `/data-export` (name `data-export`, `meta: { requiresDataExporter: true }`) with a
        minimal placeholder `web/src/views/DataExportView.vue` (page title only, `data-cy="data-export-page"`).
        Story 37.7 replaces its content and adds the menu entry. See "Open questions / assumptions".
- [x] Task 8: Role column and toggle in the Users table (AC: #2)
  - [x] `web/src/api/users.ts`: add `is_data_exporter: boolean` to `UserListAttributes`. Reuse the
        `updateUser` (PATCH) client function if Story 37.4/37.5 added one; otherwise add
        `updateUser(id: string, attributes: Partial<...>)` sending a JSON:API body.
  - [x] In the Users view from Story 37.4, render the role column as chips derived from the two booleans:
        "Admin" when `is_admin`, "Data exporter" when `is_data_exporter`, "User" when neither. Do not derive it
        from the `role` string.
  - [x] Add a per-row toggle (e.g. `v-switch` density compact, `hide-details`) labelled "Data exporter" with
        `data-cy="user-data-exporter-toggle"`. Disable it while its request is pending; on success update the
        row from the response; on failure revert and show the backend detail (Story 37.1 error pattern).
  - [x] Feedback via `v-snackbar` per `.claude/rules/feedback.md` (success 3000 ms, error 6000 ms + closable,
        `location="bottom"`, `data-cy="users-snackbar"` or the snackbar 37.4/37.5 already use).
  - [x] Soft-deleted users (Story 37.5) may still be toggled; no special handling needed.
  - [x] When the admin toggles their own row, also update `authStore.isDataExporter` (or refetch `/me`) so
        their own menu reflects the change without reload.
- [x] Task 9: i18n (AC: #2)
  - [x] Add keys for "Admin", "Data exporter", "User", the toggle label and the success/error snackbar
        prefixes to all five locale files `web/src/locales/{en,de,es,fr,uk}.json`, under the namespace Story
        37.4 introduced for the Users page. German: "Datenexport" for the role, "Benutzer" for user.
- [x] Task 10: Tests (AC: #1-#5)
  - [x] Go `internal/db`: extend migration test to assert the column exists with default 0 for a row inserted
        before/without the column value.
  - [x] Go `internal/users/store_test.go`: table test for `UpdateUser` with `IsDataExporter` true/false/nil;
        `UpsertEntraIDUser_Update` keeps a previously granted `is_data_exporter`.
  - [x] Go `internal/users/handler_test.go`: PATCH `{"is_data_exporter": true}` returns the attribute;
        PATCH without the attribute leaves it unchanged; list/get include it.
  - [x] Go `internal/middleware/require_data_exporter_test.go`: table test (exporter -> 200, admin without flag
        -> 403 `forbidden`, regular user -> 403, nil user -> 401 `auth_required`) using `runMiddleware`.
  - [x] Go: test that the per-request DB refresh (37.5) picks up a flag change for an existing cookie user
        (grant in DB -> next request sees `IsDataExporter == true`; revoke -> false).
  - [x] Go `internal/auth/me_test.go`: `/me` response contains `is_data_exporter`.
  - [x] Vitest: `useAuthStore.test.ts` (set/clear flag), `router/index.test.ts` (route exists with
        `requiresDataExporter`; guard redirects admin without flag; allows exporter), Users view test (role
        chips for all combinations; toggle calls PATCH and updates the row; failure reverts and shows snackbar),
        `api/users` test for the PATCH body.
  - [x] Cypress E2E `web/cypress/e2e/admin-data-exporter-role.cy.ts` (real backend, intercepts only for
        waiting): see "Testing standards summary".
- [x] Task 11: Quality gates
  - [x] Go: `go fmt ./...`, `go vet ./...`, `go test ./...`, `golangci-lint run ./...`,
        `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exit-code 1` (watch for a clone
        between `require_admin.go` and `require_data_exporter.go`; if flagged, extract a shared
        `requireFlag(func(*auth.User) bool)` helper).
  - [x] Frontend: `npm run type-check`, `npm run lint`, `npx vitest run --coverage` (>= 75 %),
        `npm run build`, TS jscpd, `npm run test:e2e -- --browser electron`.
  - [x] Redocly lint for `api-doc`.

### Review Findings

- [x] [Review][Patch] Test named "should render the four column headers" asserted six headers
      [web/src/views/AdminUsersView.test.ts:185] — fixed: renamed to "should render the six column headers".

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md#Story 37.6` (epics.md:5884-5913) and FR188
(epics.md:696-699). Product decision (final): `is_admin` stays; `is_data_exporter` is an independent flag;
administrators do NOT get data export unless they also hold the flag.

Dependencies: Story 37.4 (admin Users page, `v-data-table-server`, pagination helpers in `internal/api`) and
Story 37.5 (migration `000004` `users.deleted`, per-request DB user check, soft-delete UI) are assumed
implemented as described in epics.md. Stories 37.7/37.8 consume `RequireDataExporter`, the `/data-export`
route and `authStore.isDataExporter`.

### Database (AC #1)

Migrations are embedded golang-migrate files in `internal/db/migrations/`; the latest in the repo today is
`000003_add_border_width` (single `ALTER TABLE ... ADD COLUMN ... NOT NULL DEFAULT`). Epic 37 numbering is
fixed: 37.5 -> `000004`, 37.6 -> `000005`, 37.9 -> `000006`. SQLite supports `DROP COLUMN` (already used by
`000003_add_border_width.down.sql`). The `users` table is defined in `000001_init.up.sql:4-16`; booleans are
`INTEGER NOT NULL DEFAULT 0`. The default 0 guarantees no existing user gains the role.

Demo seed `tools/database/demo-users.sql` lists columns explicitly and does not need changing (the default
applies). "Felix Finance" (`felix@sithub.local`, non-admin) is a natural E2E data-exporter candidate.

### Users store and admin API (AC #2)

- `internal/users/record.go:7-18` `Record`; `internal/users/store.go:18-20` `userColumns`,
  `store.go:390-403` `scanRow` (positional scan, so column order and scan order must match).
- `store.go:157-162` `UpdateFields`, `store.go:165-209` `UpdateUser` builds the SET clause from hardcoded
  column names (keep the `//nolint:gosec // G202` justification valid; never interpolate input).
- `store.go:43-66` `UpsertEntraIDUser`: its conflict update list sets `is_admin` from the Entra group. Do not
  add `is_data_exporter` there; the flag is managed only by administrators in SitHub.
- `internal/users/handler.go:24-34` `UserAttributes`, `handler.go:159-168` `updateUserRequest`,
  `handler.go:236-261` `applyFieldUpdates` (duplicated anonymous struct), `handler.go:309-324`
  `recordToAttributes`. Routes: `internal/startup/server.go:296-301` (GET/POST/PATCH/DELETE `/api/v1/users`,
  `requireAuth, requireAdmin`).
- Story 37.4 may have changed `ListHandler`/`ListAll` for pagination; add the attribute through
  `recordToAttributes` so list, get, create and update responses all include it.
- Use `api.WriteSingle`/`api.WriteCollection` and the `//nolint:wrapcheck // Terminal response` pattern for
  terminal helper returns (commit 74529fb).

### Session, `/me` and per-request DB read (AC #4, #5)

- `auth.User` (`internal/auth/service.go:41-48`) is serialized into the signed `sithub_user` cookie
  (`EncodeUser`/`DecodeUser`, service.go:137-153) and decoded by the global `middleware.LoadUser`
  (`internal/middleware/session.go:10-23`). Adding a field is backward compatible: old cookies decode with
  `IsDataExporter == false`, and the DB refresh overwrites it anyway.
- Today `RequireAuth` (`internal/middleware/auth.go:12-28`) calls `svc.RefreshPermissions`, which is a no-op
  for local users (service.go:223-253). Story 37.5 adds a DB lookup per request for local and Entra ID users
  to reject soft-deleted accounts. Hook the flag hydration into exactly that lookup (one query, no second DB
  round trip). If 37.5 placed the lookup elsewhere (e.g. a `Service.LoadCurrentUser` helper), use that.
- `MeHandler` (`internal/auth/me.go:42-51`) returns the context user via `respondWithUserProfile`
  (me.go:143-156) using constants from `internal/auth/json_fields.go:6-14`. Because the context user is
  refreshed before the handler runs, `/me` reflects the current DB value.
- `LocalLoginHandler` builds the `User` at `login_local.go:82-89`; `FetchUser` builds it at
  service.go:209-216 from the upsert result `rec`.
- `userRole()` (`login_local.go:106-111`) stays as is (`admin`/`user`).

### Middleware (AC #3)

`RequireAdmin` (`internal/middleware/require_admin.go:10-24`) is the template. Keep `RequireDataExporter`
equally small, add a doc comment stating that admins are not implicitly exporters (security-relevant per
`.claude/rules/golang.md`). Test helper `runMiddleware` lives in `internal/middleware/helpers_test.go:13`;
existing assertions check `resp.Errors[0].Code` for `forbidden` / `auth_required`
(`require_admin_test.go:15-54`).

### Frontend (AC #2, #3, #5)

- Auth store `web/src/stores/useAuthStore.ts:4-37` (Composition API Pinia store). `setUser` callers:
  `web/src/router/index.ts:84-90` and `web/src/views/LoginView.vue:143-149`.
- Router guard `web/src/router/index.ts:76-102`; `requiresAdmin` route example at index.ts:61-66; Access
  Denied route `access-denied` (public) at index.ts:67-72. Router unit tests mock the store and `fetchMe`
  (`web/src/router/index.test.ts`).
- `web/src/api/me.ts:4-10` `UserAttributes`; `web/src/api/users.ts:4-14` `UserListAttributes` and
  `fetchUsers` (37.4 extends it with pagination). Use `apiRequest` from `web/src/api/client.ts`.
- `App.vue` gates admin menu entries with `v-if="authStore.isAdmin"` (desktop ~62-76, mobile drawer
  ~167-200). Do not add the "Data export" menu entry here; that is Story 37.7. The `/me` flag is what 37.7 uses.
- Error display: follow Story 37.1 (translated prefix + untranslated backend `errors[].detail`).

### Project Structure Notes

Files to add:

- `internal/db/migrations/000005_add_is_data_exporter.up.sql` and `.down.sql`
- `internal/middleware/require_data_exporter.go` (+ `_test.go`)
- `web/src/views/DataExportView.vue` (placeholder, replaced by 37.7)
- `web/cypress/e2e/admin-data-exporter-role.cy.ts`

Files to modify:

- `internal/users/record.go`, `store.go`, `handler.go` (+ tests)
- `internal/auth/service.go`, `login_local.go`, `me.go`, `json_fields.go` (+ tests), plus the file where
  37.5 implemented the per-request DB check (probably `internal/middleware/auth.go` or `internal/auth/service.go`)
- `api-doc/openapi.yaml` (`UserAttributes`, `UpdateUserRequest`)
- `web/src/api/me.ts`, `web/src/api/users.ts`, `web/src/stores/useAuthStore.ts`, `web/src/router/index.ts`,
  `web/src/views/LoginView.vue`, the 37.4 Users view and its test, `web/src/locales/*.json`

No new shared response types: reuse `internal/api` (`WriteSingle`, `WriteCollection`, `WriteForbidden`,
`WriteUnauthorized`). JSON attribute is snake_case `is_data_exporter`.

### Testing standards summary

Go: table-driven tests with `testify` (`require` for setup, `assert` for behavior), `t.Parallel()` where the
existing helpers allow it. Reuse `setupHandlerDB`/`seedUser` (`internal/users/handler_test.go:27,56`) and
`setupTestDB` (`store_test.go:13`). [Source: .claude/rules/golang.md]

Vitest: store, router guard, Users view, API client. Coverage >= 75 %. [Source: .claude/rules/vue.md]

Cypress E2E (every AC, real backend, `data-cy` selectors, `cy.intercept` only to wait/assert, no stubbed
responses) in `admin-data-exporter-role.cy.ts`:

- `should show Admin, Data exporter or User in the role column` (login as `anna@sithub.local`).
- `should grant and revoke the data exporter role via the toggle` (wait on `@updateUser`, assert request body
  `is_data_exporter`, response 200, chip changes, success snackbar).
- `should expose is_data_exporter on /me without re-login`: log in as Felix in a `cy.request` session,
  grant via admin API, then `GET /api/v1/me` as Felix shows `true`; revoke shows `false`.
- `should show Access Denied on /data-export for an admin without the flag` and `... for a regular user`.
- `should open /data-export for a data exporter`.
- Reset: each test resets Felix's flag to false via `PATCH /api/v1/users/:id` as admin in `beforeEach`/
  `afterEach` so order does not matter. Use `cy.login(email, password)` (`web/cypress/support/commands.ts:38`);
  demo password `SitHubDemo2026!!`.

The API-level 403 on export endpoints cannot be exercised end-to-end until 37.7 adds them; cover it in the
middleware Go test now and in 37.7's E2E later.

### Open questions / assumptions

- Placeholder route: `/data-export` with a minimal `DataExportView.vue` is added now so the route guard (AC #3)
  is testable end-to-end; no menu entry until 37.7. Recommended default; 37.7 replaces the view content.
- `role` string attribute stays `admin`/`user` for backward compatibility; the UI derives role chips from the
  booleans. No `roles` array is introduced.
- Admins may toggle their own data exporter flag (no self-restriction, unlike self-delete in 37.5).
- Entra ID users can be flagged; the flag survives Entra logins because `UpsertEntraIDUser` does not touch it.
- `is_admin` is refreshed from the DB by the 37.5 per-request lookup (product decision 2026-10-08); add
  `is_data_exporter` to that same lookup, do not introduce a second query.
- If golangci-lint flags an unused `requireDataExporter` variable in `registerRoutes`, defer wiring it to 37.7.

### References

- [Source: _bmad-output/planning-artifacts/epics.md:696-699 (FR188), 5884-5913 (Story 37.6),
  5836-5882 (Story 37.5), 5802-5834 (Story 37.4), 5915-5952 (Story 37.7)]
- [Source: internal/db/migrations/000001_init.up.sql:4-16; 000003_add_border_width.*.sql]
- [Source: internal/users/record.go:7-18; internal/users/store.go:18-20,43-66,157-209,390-403]
- [Source: internal/users/handler.go:24-34,159-168,236-261,309-324]
- [Source: internal/auth/service.go:41-48,137-153,204-216,223-253]
- [Source: internal/auth/me.go:42-51,143-156; internal/auth/json_fields.go:6-14;
  internal/auth/login_local.go:82-89,106-111]
- [Source: internal/middleware/require_admin.go:10-24; auth.go:12-28; session.go:10-23;
  require_admin_test.go:15-54]
- [Source: internal/startup/server.go:247,254,296-301]
- [Source: api-doc/openapi.yaml UserAttributes (~652-687), UpdateUserRequest (~747-770)]
- [Source: web/src/stores/useAuthStore.ts:4-37; web/src/router/index.ts:61-102;
  web/src/views/LoginView.vue:143-149; web/src/api/me.ts:4-10; web/src/api/users.ts:4-14; web/src/App.vue:62-76]
- [Source: tools/database/demo-users.sql; web/cypress/support/commands.ts:38]
- [Source: .claude/rules/golang.md; json-api.md; apidoc.md; vue.md; cypress.md; feedback.md]

## Dev Agent Record

### Agent Model Used

Claude Sonnet 5.5 (claude-sonnet-5-5)

### Debug Log References

- `RequireDataExporter` was flagged by jscpd as a clone of `RequireAdmin` (13th clone, baseline is 12). Both now
  delegate to `requireRole(func(*auth.User) bool)` in `internal/middleware/require_role.go`, as the story
  suggested; back to 12 clones (1.42 %).
- The view test stubs `v-switch` as a checkbox. It re-syncs to `modelValue` after the emit, mirroring a
  controlled Vuetify switch, which is what makes the "failed update reverts" test meaningful.

### Completion Notes List

- Migration `000005_add_is_data_exporter` (up/down). Tests: the flag defaults to 0 for an existing admin and the
  down migration drops the column. Applied cleanly on the already populated E2E database (version 4 to 5).
- Store: `Record.IsDataExporter`, `userColumns`/`scanRow`, `UpdateFields.IsDataExporter` (a nil pointer leaves
  the value untouched). `UpsertEntraIDUser` does not touch it; a test proves an Entra ID login keeps a granted
  role.
- API: `is_data_exporter` on user attributes (list, get, create, update) and in the PATCH body; `role` stays
  `admin`/`user`. The duplicated anonymous PATCH struct became `updateUserAttributes`.
- Session: `users.LookupAccess` (the 37.5 per-request lookup) returns the flag; `Service.EnsureActive` copies it
  onto the context user, so grant/revoke applies on the next request even with a stale cookie (tested in both
  directions). The flag is also set at local login and in Entra `FetchUser`, and returned by `/api/v1/me`,
  `PATCH /me` and the login response.
- `RequireDataExporter` is implemented and unit-tested (exporter 200, admin without flag 403, user 403,
  missing user 401) but not yet mounted on a route; 37.7/37.8 do that. A placeholder `/data-export` route and
  `DataExportView` make the frontend guard testable end to end.
- Frontend: `authStore.isDataExporter`, `meta.requiresDataExporter` guard (administrators are not exempt),
  role chips from both flags via `userRoleLabels` (Admin and/or Data exporter, "User" if neither), a
  "Data exporter" column with a `v-switch` per row (disabled while the request runs, backend detail shown on
  failure, own-row toggle also updates the auth store), `updateUser` API function, i18n in five locales.
- Cypress `admin-data-exporter-role.cy.ts` (6 tests): role column, grant/revoke via toggle with request and
  response assertions, `/me` follows the flag without re-login, Access Denied for an admin without the flag and
  for a regular user, `/data-export` opens for an exporter. The spec locates Felix's table page through the API,
  so it does not break as the number of users grows, and resets the flag before and after.
- Gates: `go test -race`, golangci-lint 0 issues, Go jscpd 1.42 %, redocly valid, type-check, eslint, Vitest
  652/652 (coverage 82.1 %), TS jscpd 0 clones.
- Known limit: `authStore.isDataExporter` is loaded once per page session, like `isAdmin`. A role change made
  elsewhere shows in the UI after the next reload; the API enforces the current value on every request.

### File List

- internal/db/migrations/000005_add_is_data_exporter.up.sql (new)
- internal/db/migrations/000005_add_is_data_exporter.down.sql (new)
- internal/db/migrate_test.go
- internal/users/record.go
- internal/users/store.go
- internal/users/deletion.go
- internal/users/deletion_test.go
- internal/users/handler.go
- internal/users/handler_test.go
- internal/auth/service.go
- internal/auth/login_local.go
- internal/auth/json_fields.go
- internal/auth/me.go
- internal/auth/me_test.go
- internal/middleware/require_role.go (new)
- internal/middleware/require_admin.go
- internal/middleware/require_data_exporter.go (new)
- internal/middleware/require_data_exporter_test.go (new)
- api-doc/openapi.yaml
- web/src/api/me.ts
- web/src/api/users.ts
- web/src/api/users.test.ts
- web/src/stores/useAuthStore.ts
- web/src/stores/useAuthStore.test.ts
- web/src/router/index.ts
- web/src/router/index.test.ts
- web/src/views/LoginView.vue
- web/src/views/DataExportView.vue (new)
- web/src/views/AdminUsersView.vue
- web/src/views/AdminUsersView.test.ts
- web/src/utils/userRoles.ts
- web/src/utils/userRoles.test.ts
- web/src/locales/{en,de,es,fr,uk}.json
- web/cypress/e2e/admin-data-exporter-role.cy.ts (new)
- _bmad-output/implementation-artifacts/sprint-status.yaml

## Change Log

- 2026-10-08: Implemented story 37.6 (independent data exporter role: migration, admin API and toggle,
  per-request database refresh, `/me` attribute, `RequireDataExporter` middleware, route guard). Status -> review.
- 2026-10-08: Code review (stories 37.4-37.6): fixed a stale test name in `AdminUsersView.test.ts`.
  Status -> done.
- 2026-10-08: Re-review after patch fixes — clean, no new findings; gates green (go test, vet, vitest).
