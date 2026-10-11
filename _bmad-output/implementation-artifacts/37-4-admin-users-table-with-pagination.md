# Story 37.4: Admin Users Table with Pagination

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an administrator,
I want a paginated table of all users,
so that I can review accounts and their activity.

## Acceptance Criteria

1. Given I am an administrator, when I open the user menu, then I see a "Users" entry (gated by `isAdmin`,
   route with `meta.requiresAdmin`); non-admins do not see it and are redirected to Access Denied when
   opening the route directly.
2. Given the Users page, when it loads, then a server-paginated table (`v-data-table-server`) shows the
   columns email, name (display name), role(s), and last login (localized date-time, "never" if empty),
   20 rows per page, pagination beneath.
3. Given `GET /api/v1/users?page[limit]=20&page[offset]=40`, when called by an administrator, then the
   response contains the requested slice, `meta.total`, and JSON:API `links` (self, first, prev, next,
   last); `page[limit]` is capped at 100 and invalid values return 400.
4. Given the pagination support, when implemented, then the pagination parameters, `PaginationMeta`, and
   `PaginationLinks` live in `internal/api` so later stories reuse them, and `api-doc` documents the
   parameters.
5. Given more than 20 users, when I page through the table, then every user appears exactly once, sorted
   by name.

## Tasks / Subtasks

- [x] Task 1: Reusable JSON:API pagination helpers in `internal/api` (AC: #3, #4)
  - [x] Create `internal/api/pagination.go` with the shape described in Dev Notes
        "Pagination helper contract": `Page`, `ParsePage`, `ErrInvalidPagination`, `PaginationMeta`,
        `PaginationLinks`, `BuildPaginationLinks`, and constants `DefaultPageLimit = 20`,
        `MaxPageLimit = 100`.
  - [x] Extend `CollectionResponse` (`internal/api/response.go:28-31`) with optional
        `Meta *PaginationMeta` and `Links *PaginationLinks` (`omitempty`) so unpaginated endpoints
        keep their exact JSON output.
  - [x] Add `WritePaginatedCollection(c, resources, page, total, errLabel)` to `internal/api/write.go`;
        let `WriteCollection` and the new writer share one private write helper (jscpd gate).
  - [x] Add `internal/api/pagination_test.go`: table-driven tests for `ParsePage` (absent params,
        valid values, limit 101 -> 100, limit 0, negative limit, negative offset, non-numeric, empty
        string, float, huge number overflow) and for `BuildPaginationLinks` (first page has
        `prev: null`, last page has `next: null`, `total=0`, offset not aligned to limit, offset past
        the end, extra query params preserved, `self` echoes the effective capped limit).
- [x] Task 2: Paginated users store query (AC: #3, #5)
  - [x] Replace `ListAll` (`internal/users/store.go:101-128`, only caller is `ListHandler`) with
        `ListPage(ctx, db, limit, offset int) ([]Record, error)` and `Count(ctx, db) (int, error)`.
  - [x] Order deterministically: `ORDER BY display_name COLLATE NOCASE, id` so pages never overlap or
        skip rows when display names are equal (AC #5). Use `LIMIT ? OFFSET ?` placeholders.
  - [x] Return an empty, non-nil slice when no rows match (avoids `"data": null`).
  - [x] Update `store_test.go` (`TestListAll`, `TestListAll_Empty` at 208-237) to table-driven
        `TestListPage` cases: first page, middle page, last partial page, offset beyond total, tie on
        display name resolved by id, case-insensitive ordering; plus `TestCount`.
- [x] Task 3: Paginated `ListHandler` (AC: #3, #5)
  - [x] In `internal/users/handler.go:36-56` call `api.ParsePage(c)`; on `ErrInvalidPagination` return
        `api.WriteBadRequest(c, err.Error())` (detail names the offending parameter).
  - [x] Fetch `Count` and `ListPage`, map with `api.MapResources` + `recordToAttributes` (unchanged) and
        respond via `api.WritePaginatedCollection`.
  - [x] Update `handler_test.go` `TestListHandler` (76-93) into a table test: default page (no params ->
        limit 20), `page[limit]=20&page[offset]=40` slice with 45 seeded users, limit 500 capped to 100,
        invalid values -> 400 JSON:API error, `meta.total` and all five `links` asserted.
- [x] Task 4: API documentation (AC: #4)
  - [x] In `api-doc/openapi.yaml` add `components.parameters.PageLimit` / `PageOffset` (query names
        `page[limit]` / `page[offset]`, integer, `minimum` 1 / 0, limit `maximum` 100, default 20 / 0) and
        `components.schemas.PaginationMeta` / `PaginationLinks`; reference them from
        `UserCollectionResponse` (`openapi.yaml:705-715`) as `meta` and `links`.
  - [x] In `api-doc/endpoints/users.yaml` `get`: reference both parameters, document ordering
        (by display name), and add `400` (invalid pagination) and `403` (admin only) responses.
  - [x] Lint: `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.
- [x] Task 5: Frontend API layer (AC: #2, #3)
  - [x] In `web/src/api/types.ts` add `PaginationMeta`, `PaginationLinks`, and
        `PaginatedCollectionResponse<T>` (extends `CollectionResponse<T>` with required `meta` and
        `links`). 37.7 reuses these.
  - [x] In `web/src/api/users.ts` change `fetchUsers()` (currently unused, line 12-14) to
        `fetchUsers({ limit, offset })` building the query with `URLSearchParams`
        (`page[limit]`, `page[offset]`); extend `UserListAttributes` with `last_login`, `created_at`,
        `updated_at`.
  - [x] Add `web/src/api/users.test.ts` (or extend an existing api test) asserting the request URL.
- [x] Task 6: Users page, route and menu entries (AC: #1, #2, #5)
  - [x] Create `web/src/views/AdminUsersView.vue` (`<script setup lang="ts">`, Composition API) with a
        `v-data-table-server` (`data-cy="users-table"`), `:items-length="total"`,
        `:items-per-page="20"`, `:loading`, `hide-default-footer`, and a `v-pagination`
        (`data-cy="users-pagination"`) directly beneath; changing the page re-fetches with
        `offset = (page - 1) * 20`.
  - [x] Headers (all `sortable: false`, server order is fixed): email, name, role(s), last login. Row
        cells carry `data-cy="user-email"`, `user-name`, `user-roles`, `user-last-login`.
  - [x] Role(s) column: `v-chip`s produced by a small pure helper `userRoleLabels(attrs)` (Admin / User
        today; 37.6 appends "Data exporter").
  - [x] Last login: localized date-time via a new `formatDateTime(iso, locale)` in
        `web/src/utils/date.ts` (`Intl.DateTimeFormat`, `dateStyle: 'medium'`, `timeStyle: 'short'`),
        "never" i18n text when the value is empty or unparseable.
  - [x] Load failure is a persistent page state: show a `v-alert` (`data-cy="users-load-error"`) with
        the backend detail (`ApiError.detail`) per `.claude/rules/feedback.md`.
  - [x] Router (`web/src/router/index.ts`): add `{ path: '/admin/users', name: 'admin-users',
        component: () => import('../views/AdminUsersView.vue'), meta: { requiresAdmin: true } }` next to
        `floor-plan-editor` (lines 61-66). The existing guard (96-98) handles the redirect.
  - [x] `web/src/App.vue`: add a "Users" `v-list-item` gated by `authStore.isAdmin` in the user menu
        (after `floor-plan-editor-btn`, 67-76, `data-cy="users-menu-btn"`, icon `$users`) and in the mobile
        drawer (after 191-201, `data-cy="mobile-users-menu-btn"`, closes the drawer on click).
  - [x] i18n in all five locales (`web/src/locales/{en,de,es,fr,uk}.json`): `app.userMenu.users`, and a
        new `adminUsers` block (`title`, `colEmail`, `colName`, `colRoles`, `colLastLogin`, `never`,
        `roleAdmin`, `roleUser`, `loadError`). German: "Benutzer", "Zuletzt angemeldet", "nie".
- [x] Task 7: Frontend unit tests (AC: #1, #2, #5)
  - [x] `web/src/views/AdminUsersView.test.ts`: renders rows from a mocked `fetchUsers`; first call
        uses limit 20 / offset 0; selecting page 3 calls offset 40; "never" for empty `last_login`;
        localized date for a set value; admin chip for `is_admin`; error alert shows backend detail.
  - [x] Extend `web/src/router/index.test.ts`: non-admin navigating to `/admin/users` lands on
        `access-denied`; admin reaches `admin-users`.
  - [x] Extend `web/src/__tests__/App.test.ts`: Users entry visible for admins, hidden for non-admins.
  - [x] Unit-test `formatDateTime` and `userRoleLabels` (table-driven).
- [x] Task 8: Cypress E2E `web/cypress/e2e/admin-users.cy.ts` (AC: #1, #2, #3, #5)
  - [x] `before`: as admin (`cy.login()` -> anna), seed users via `cy.request('POST', '/api/v1/users')`
        with fixed emails `e2e-pager-01..10@sithub.local` so the total exceeds 20 together with the 15
        demo users; accept 201 or 409 so reruns against the same DB stay idempotent.
  - [x] "should show the Users menu entry for admins and open the table": menu -> `users-menu-btn` ->
        `cy.wait('@listUsers')` (intercept `GET /api/v1/users*`), assert 20 rows, 4 headers, pagination.
  - [x] "should page through all users exactly once sorted by name": walk every page via
        `users-pagination`, waiting on `@listUsers` each time; collect emails, assert no duplicates,
        count equals `meta.total` from the first response, names sorted case-insensitively.
  - [x] "should reject invalid pagination parameters": `cy.request({ url: '/api/v1/users?page[limit]=abc',
        failOnStatusCode: false })` -> 400 with `errors[0].detail`; `page[limit]=500` -> 100 or fewer rows
        and `links.self` contains the capped limit.
  - [x] "should hide the Users entry and deny the route for non-admins": `cy.login('alex@sithub.local',
        'SitHubDemo2026!!')`, entry absent, `cy.visit('/admin/users')` ends on `/access-denied`.
- [x] Task 9: Quality gates (all ACs)
  - [x] Go: `go fmt ./...`, `go vet ./...`, `go test ./...`, `golangci-lint run ./...`,
        `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exit-code 1`.
  - [x] Web: `npm run type-check`, `npm run lint`, `npx vitest run --coverage` (>= 75%), `npm run build`,
        `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0 --exit-code 1`,
        `npm run test:e2e -- --browser electron` against the dev server (backend on :9900, Vite on :5173).

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md` "Epic 37 Stories" / Story 37.4 (line 5802 ff.), FR184
(`epics.md:681-683`). Product decisions are final (Epic 37 briefing): one `display_name` column only, no
first/last name split; `is_admin` stays as is (the `is_data_exporter` flag comes in 37.6).

### Current state (verified 2026-10-08)

- Backend routes already exist and are admin-gated: `internal/startup/server.go:296-301`
  (`GET /api/v1/users` with `requireAuth, requireAdmin`). No route change is needed.
- `ListHandler` (`internal/users/handler.go:36-56`) returns ALL users via `ListAll`
  (`internal/users/store.go:101-128`, `ORDER BY display_name`, case-sensitive, no tie-breaker) and
  `api.WriteCollection`. `ListAll` has no other caller (`grep ListAll` -> handler only).
- `UserAttributes` (`handler.go:24-34`) already contains `email`, `display_name`, `is_admin`, `role`,
  `last_login`, `created_at`, `updated_at`. `last_login` is written as RFC3339 UTC by
  `UpdateLastLogin` (`store.go:247-265`) and is `''` for users who never logged in (demo seed data uses
  `''`). Keep the attribute name `last_login` (existing public field; renaming to `last_login_at` would
  break the documented API and is out of scope).
- No pagination exists anywhere yet: no `page[...]` parsing, no `PaginationMeta`, no
  `v-data-table-server`. `.claude/rules/json-api.md` "Pagination" mandates offset-based
  `page[limit]` / `page[offset]`; `.claude/rules/golang.md` "Shared Types" mandates
  `PaginationMeta` / `PaginationLinks` in `internal/api`.
- Frontend: no Users page, route or menu entry. `fetchUsers()` in `web/src/api/users.ts:12-14` is
  unused; `UserListAttributes` lacks `last_login`.
- `internal/api` error writers (`internal/api/errors.go`) write the response and return `nil`. In
  non-handler helper funcs that return `(x, error)`, the call needs `//nolint:wrapcheck // Terminal
  response` (pattern: `handler.go:249-254`, introduced by commit 74529fb).

### Pagination helper contract (reused by 37.7 data export)

Put everything in a new `internal/api/pagination.go`; keep it framework-light so 37.7 can call it with
additional filter query params.

```go
const (
    DefaultPageLimit = 20
    MaxPageLimit     = 100
)

// ErrInvalidPagination reports a malformed page[limit] or page[offset] value.
var ErrInvalidPagination = errors.New("invalid pagination parameter")

// Page is a validated offset-based page request.
type Page struct {
    Limit  int
    Offset int
}

// ParsePage reads page[limit] and page[offset] from the request query.
// Absent values default to DefaultPageLimit / 0. A limit above MaxPageLimit is capped.
// Non-integer values, limit < 1 and offset < 0 return ErrInvalidPagination
// (wrapped with the parameter name, e.g. "invalid pagination parameter: page[limit] must be ...").
func ParsePage(c echo.Context) (Page, error)

// PaginationMeta is the JSON:API top-level meta for paginated collections.
type PaginationMeta struct {
    Total  int `json:"total"`
    Limit  int `json:"limit"`
    Offset int `json:"offset"`
}

// PaginationLinks holds JSON:API pagination links. Prev/Next are null when unavailable.
type PaginationLinks struct {
    Self  string  `json:"self"`
    First string  `json:"first"`
    Prev  *string `json:"prev"`
    Next  *string `json:"next"`
    Last  string  `json:"last"`
}

// BuildPaginationLinks builds links from the request URL, preserving all non-page query params.
func BuildPaginationLinks(u *url.URL, page Page, total int) PaginationLinks
```

Rules for `BuildPaginationLinks`:

- Base is `u.Path` (e.g. `/api/v1/users`), relative links like the examples in `json-api.md`.
- Copy `u.Query()`, overwrite `page[limit]` / `page[offset]` per link, keep every other param (37.7 adds
  date range and area filters).
- `last` offset = `((total-1)/limit)*limit`, or `0` when `total == 0`.
- `prev` = `max(offset-limit, 0)` when `offset > 0`, else `nil`; `next` = `offset+limit` when
  `offset+limit < total`, else `nil`.
- `self` uses the effective (capped) limit, not the raw input.
- `url.Values.Encode()` percent-encodes brackets (`page%5Blimit%5D=20`); that is valid and Echo decodes
  it on the way in. Do not hand-build query strings.

`WritePaginatedCollection(c echo.Context, resources []Resource, page Page, total int, errLabel string)`
builds `CollectionResponse{Data, Meta: &PaginationMeta{...}, Links: &links}` from `c.Request().URL`.

### Response example (`GET /api/v1/users?page[limit]=20&page[offset]=40`, 45 users)

```json
{
  "data": [{ "id": "...", "type": "users", "attributes": { "email": "...", "display_name": "..." } }],
  "meta": { "total": 45, "limit": 20, "offset": 40 },
  "links": {
    "self": "/api/v1/users?page%5Blimit%5D=20&page%5Boffset%5D=40",
    "first": "/api/v1/users?page%5Blimit%5D=20&page%5Boffset%5D=0",
    "prev": "/api/v1/users?page%5Blimit%5D=20&page%5Boffset%5D=20",
    "next": null,
    "last": "/api/v1/users?page%5Blimit%5D=20&page%5Boffset%5D=40"
  }
}
```

### Store query

```sql
SELECT <userColumns> FROM users ORDER BY display_name COLLATE NOCASE, id LIMIT ? OFFSET ?;
SELECT COUNT(*) FROM users;
```

The two statements may run without a transaction; a user created between them only shifts `total` by
one, which is acceptable for an admin list. 37.5 adds `users.deleted`; the admin table must still list
soft-deleted users, so do NOT add a `deleted` filter to these queries later.

### Frontend implementation notes

- `v-data-table-server` emits `update:options`; with `hide-default-footer` drive paging solely from a
  `page` ref bound to `v-pagination` (`:length="Math.ceil(total / 20)"`) and a `watch(page, load)`.
  Load once on mount. Using the custom `v-pagination` gives a stable `data-cy` and avoids a per-page
  selector the AC does not ask for.
- Keep API calls in `web/src/api/users.ts`; the view can use `useApi()` (`web/src/composables/useApi.ts`)
  for `loading` but must read `ApiError.detail` itself for the alert (see 37.1, which standardizes
  error display; reuse its helper if 37.1 has landed).
- Vuetify locale messages are wired through the vue-i18n adapter (`web/src/plugins/i18n.ts:13-19`), so
  built-in table texts (loading, no data) are already translated.
- Locale-aware formatting: follow `MyBookingsView.vue:175, 199` (`const { t, locale } = useI18n()` and
  `locale.value || undefined`). Place `formatDateTime` beside `formatShortDate` in
  `web/src/utils/date.ts:10-19`; parse the ISO string as-is (it carries `Z`), do not append a time.
- The `$users` icon alias already exists (`web/src/plugins/vuetify.ts:128`, `mdiAccountGroup`).
- Mirror the existing admin menu item markup exactly (`App.vue:67-76` desktop, `App.vue:191-201`
  mobile) to keep jscpd happy and the menus consistent.

### Relationship to sibling stories

- 37.5 (Delete Orphaned Users) adds a delete icon per row and a "deleted" marker in this table: leave an
  actions slot easy to add (e.g. headers built in one `computed`), but do not add delete UI here.
- 37.6 (Data Exporter Role) adds `is_data_exporter` to attributes and an extra role chip: keep role
  rendering in the single `userRoleLabels` helper.
- 37.7 (Data Export Table) reuses `ParsePage`, `PaginationMeta`, `PaginationLinks`,
  `WritePaginatedCollection`, the TS `PaginatedCollectionResponse`, and the `v-data-table-server` +
  `v-pagination` pattern. Design them generically now; nothing in `internal/api` may reference users.

### Open questions / assumptions

- Requests without `page[...]` now return the first 20 users (default limit) instead of all users.
  Assumed acceptable: the only frontend caller is new, and the colleague picker uses
  `/api/v1/colleagues` (unchanged). Recommended default: keep.
- `page[limit]` above 100 is silently capped (per AC), not rejected; `page[offset]` beyond `total`
  returns `200` with empty `data` and valid links (not 400). Recommended default: as stated.
- Unknown `page[...]` keys (e.g. `page[number]`) are ignored rather than rejected. Recommended default:
  ignore.
- `meta` carries `limit` and `offset` in addition to the required `total` for client convenience.
- Name sorting is case-insensitive ASCII (`COLLATE NOCASE`); locale-aware collation for umlauts is not
  available in SQLite without extensions and is out of scope.
- The route path is `/admin/users` (consistent with `/admin/floor-plan-editor`).

### Project Structure Notes

- New: `internal/api/pagination.go` (+ `_test.go`), `web/src/views/AdminUsersView.vue` (+ `.test.ts`),
  `web/cypress/e2e/admin-users.cy.ts`.
- Modified: `internal/api/response.go`, `internal/api/write.go`, `internal/users/store.go`,
  `internal/users/handler.go` (+ both tests), `api-doc/openapi.yaml`, `api-doc/endpoints/users.yaml`,
  `web/src/api/types.ts`, `web/src/api/users.ts`, `web/src/utils/date.ts`, `web/src/router/index.ts`,
  `web/src/App.vue`, `web/src/locales/{en,de,es,fr,uk}.json`, and the related unit tests.
- No migration, no config change, no new dependency.

### Testing standards summary

- Go: table-driven tests with `t.Run`, `require` for setup, `assert` for behavior, `t.Parallel()` where
  the in-memory DB is per test (`setupHandlerDB`, `handler_test.go:27-53`). Seed users with `seedUser`
  (`handler_test.go:55-74`). Sentinel error `ErrInvalidPagination` checked with `errors.Is`. Wrap with
  `%w` at the end of format strings.
- Vitest: mock `../api/users` like `MatrixBookingPopover.test.ts:16`; coverage >= 75%.
- Cypress: data-cy selectors only, intercepts for waiting (`@listUsers`) not mocking, no `cy.wait(ms)`,
  test names "should ...". CI seeds 15 demo users from `tools/database/demo-users.sql` (admins
  `anna@`, `sam@`; non-admin `alex@`, password `SitHubDemo2026!!`) and logs in as anna
  (`.github/workflows/ci.yml:183-185`), so the spec must create the extra users itself.

### Project rules (from `.claude/rules/`)

- `json-api.md`: `application/vnd.api+json`, snake_case attributes, top-level `links` incl. pagination,
  errors as `errors[]` with `status`/`title`/`detail`, offset pagination via `page[limit]`/`page[offset]`.
- `golang.md`: shared types in `internal/api`, sentinel errors, doc comments on exported identifiers,
  document security-relevant input parsing (`ParsePage`), golangci-lint v2.5.0, jscpd threshold 0.
- `apidoc.md`: OpenAPI 3.1, one file per endpoint, redocly lint.
- `vue.md`: `<script setup lang="ts">`, no `any`, Vuetify components, all five locales.
- `feedback.md`: page load failure -> `v-alert`; transient actions -> `v-snackbar` (none in this story).
- `cypress.md`: E2E for every AC, run against the dev server.
- `markdown.md`: 120-char lines, lint with `npx --yes markdownlint-cli`.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.4 (5802-5835), FR184 (681-683)]
- [Source: internal/users/handler.go:24-56, 249-262, 309-325]
- [Source: internal/users/store.go:101-128, 247-265]
- [Source: internal/users/handler_test.go:27-93; internal/users/store_test.go:208-237]
- [Source: internal/api/response.go:16-40; internal/api/write.go:13-31; internal/api/errors.go]
- [Source: internal/startup/server.go:295-301]
- [Source: api-doc/openapi.yaml:45-46, 60-100, 652-715; api-doc/endpoints/users.yaml]
- [Source: web/src/api/users.ts; web/src/api/types.ts; web/src/api/client.ts:3-12, 62]
- [Source: web/src/router/index.ts:61-66, 96-98; web/src/stores/useAuthStore.ts:8, 22]
- [Source: web/src/App.vue:59-76, 162-201; web/src/plugins/vuetify.ts:128, 280-284]
- [Source: web/src/plugins/i18n.ts:9-20; web/src/utils/date.ts:10-19]
- [Source: tools/database/demo-users.sql; .github/workflows/ci.yml:168-185]
- [Source: .claude/rules/json-api.md#Pagination, .claude/rules/golang.md#Shared Types]

## Dev Agent Record

### Agent Model Used

Claude Opus 5.5 (claude-opus-5-5)

### Debug Log References

- First E2E run: table and pagination rendered empty although `GET /api/v1/users` returned 200. Cause: the
  Vuetify plugin registers components explicitly (`web/src/plugins/vuetify.ts`) and `VDataTableServer` /
  `VPagination` were not in the list. Registered both; re-run green. A unit guard for the registration was
  attempted but importing the real plugin in Vitest fails on Vuetify CSS imports, so the Cypress spec is the
  regression guard.
- `golangci-lint` `unparam` flagged the test helper `pageLink(path ...)`; replaced by a constant.
- Go jscpd at threshold 0 reports 12 clones; identical 12 clones exist at HEAD (none introduced here). CI
  gates Go duplication at 3 %.

### Completion Notes List

- `internal/api/pagination.go`: `Page`, `ParsePage`, `ErrInvalidPagination`, `PaginationMeta`,
  `PaginationLinks`, `BuildPaginationLinks`, `DefaultPageLimit = 20`, `MaxPageLimit = 100`. Nothing references
  users, so 37.7 can reuse it. Error text names the parameter, e.g.
  `page[limit] must be a positive integer: invalid pagination parameter` (`%w` at the end per golang rules).
  The next-link check uses `total-offset > limit` so a huge offset cannot overflow.
- `CollectionResponse` gained optional `Meta` / `Links` (`omitempty`); unpaginated endpoints keep their exact
  JSON (`{"data":[]}` asserted in a test). `WriteCollection` and `WritePaginatedCollection` share
  `writeCollection`.
- Users store: `ListAll` replaced by `ListPage` (`ORDER BY display_name COLLATE NOCASE, id LIMIT ? OFFSET ?`,
  empty non-nil slice) and `Count`. `ListHandler` returns 400 JSON:API errors for invalid parameters.
- OpenAPI: `components.parameters.PageLimit` / `PageOffset`, `components.schemas.PaginationMeta` /
  `PaginationLinks`, referenced from `UserCollectionResponse` and `endpoints/users.yaml` (ordering, 400, 403).
  Redocly lint valid (one pre-existing warning on another endpoint).
- Frontend: `PaginationMeta`, `PaginationLinks`, `PaginatedCollectionResponse<T>` in `api/types.ts`;
  `fetchUsers({ limit, offset })` via `URLSearchParams`; `formatDateTime` in `utils/date.ts` returns `null` for
  empty/unparseable input and the view renders the i18n "never"; `userRoleLabels` returns i18n keys.
- `AdminUsersView.vue`: `v-data-table-server` + `v-pagination` beneath, non-sortable headers built in one
  computed, rows derived in a computed (a language switch re-formats dates without refetching), a request
  counter discards out-of-order responses, load failures show a `v-alert` with the translated prefix plus the
  backend detail (`formatApiErrorMessage` from 37.1), 401/403 handled via `useAuthErrorHandler`.
- Route `/admin/users` (`admin-users`, `requiresAdmin`), desktop and mobile menu entries with `$users`,
  i18n in en/de/es/fr/uk.
- Tests: Go table tests for `ParsePage` (15 cases), `BuildPaginationLinks` (7), `WritePaginatedCollection`,
  `ListPage` (5 + empty), `Count`, `ListHandler` (3 + 3 invalid). Vitest: `users.test.ts`, `date.test.ts`,
  `userRoles.test.ts`, `AdminUsersView.test.ts` (10), router guard tests, App menu visibility. The story
  referenced `web/src/__tests__/App.test.ts`; the existing file is `web/src/App.test.ts`.
- Cypress `admin-users.cy.ts` (5 tests) covers AC 1, 2, 3, 5 and is idempotent on reruns.
- Gates: `go fmt`, `go vet`, `go test -race ./...`, `golangci-lint` (0 issues), redocly, `npm run type-check`,
  `npm run lint`, Vitest 630/630 (coverage 81.8 % statements), TS jscpd 0 clones, `npm run build`, Cypress
  full suite 86/86 against a freshly built binary on :8080 with demo users.

### File List

- internal/api/pagination.go (new)
- internal/api/pagination_test.go (new)
- internal/api/response.go
- internal/api/write.go
- internal/users/store.go
- internal/users/store_test.go
- internal/users/handler.go
- internal/users/handler_test.go
- api-doc/openapi.yaml
- api-doc/endpoints/users.yaml
- web/src/api/types.ts
- web/src/api/users.ts
- web/src/api/users.test.ts (new)
- web/src/utils/date.ts
- web/src/utils/date.test.ts (new)
- web/src/utils/userRoles.ts (new)
- web/src/utils/userRoles.test.ts (new)
- web/src/views/AdminUsersView.vue (new)
- web/src/views/AdminUsersView.test.ts (new)
- web/src/router/index.ts
- web/src/router/index.test.ts
- web/src/App.vue
- web/src/App.test.ts
- web/src/plugins/vuetify.ts
- web/src/locales/en.json
- web/src/locales/de.json
- web/src/locales/es.json
- web/src/locales/fr.json
- web/src/locales/uk.json
- web/cypress/e2e/admin-users.cy.ts (new)
- _bmad-output/implementation-artifacts/sprint-status.yaml

## Change Log

- 2026-10-08: Implemented story 37.4 (paginated admin users table, reusable JSON:API pagination helpers,
  API docs, Users menu entry and route, unit and E2E tests). Status -> review.
- 2026-10-08: Code review (stories 37.4-37.6) — no findings for this story; all ACs verified covered.
  Status -> done.
- 2026-10-08: Re-review after patch fixes — clean, no new findings; gates green (go test, vet, vitest).
