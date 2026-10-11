# Story 37.8: Export Bookings as CSV

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a data exporter,
I want to download the filtered bookings as a CSV file,
so that finance can process them in third-party software.

## Acceptance Criteria

1. Given the Data export page with a date range and area selected, when I click "Export as CSV" beneath
   the table, then a CSV file downloads containing all bookings matching the filters, not just the
   visible page.
2. Given the CSV, when opened, then it has a header row (name, email, booked item, booking date), UTF-8
   encoding with BOM so spreadsheet tools show umlauts correctly, RFC 4180 quoting, and ISO dates
   (`YYYY-MM-DD`).
3. Given the download, when the file is saved, then its name contains start and end date, e.g.
   `sithub-bookings-2026-10-03_2026-10-08.csv`.
4. Given the endpoint `GET /api/v1/bookings/export.csv` with the same filters as
   `GET /api/v1/bookings/export`, when called without the data exporter role, then it returns 403.
5. Given cell values starting with `=`, `+`, `-`, or `@`, when written to the CSV, then they are escaped
   to prevent spreadsheet formula injection.

FRs covered: FR194. [Source: _bmad-output/planning-artifacts/epics.md#Story 37.8, FR194]

## Tasks / Subtasks

- [x] Task 1: Make the Story 37.7 export query reusable without pagination (AC: #1, #4)
  - [x] Reuse the 37.7 building blocks in `internal/bookings/` (names as specified in the 37.7 story
        file `37-7-data-export-table-with-date-and-area-filters.md`): the query-parameter parser
        `parseExportParams(c, cfg)` in `export.go` (from/to/area validation,
        defaults "last 5 days including today", unknown area -> 400), the store query (all users,
        inclusive range, area filter via item IDs, sorted by `booking_date` then name), and the
        row mapper that resolves name/email (booked-for person, guest name/email for guests,
        soft-deleted users included) and the item label `"<item group> / <item>"`.
  - [x] Ensure the store query can run unpaginated (e.g. a zero-value page / `limit <= 0` means
        "no LIMIT/OFFSET" in 37.7's `ListExportBookings(ctx, store, filter, limit, offset)` from
        `export_store.go`, with `ExportFilter` unchanged; the CSV path ignores the returned total).
        Do NOT copy the SQL or the filter parser into a new function: jscpd runs with threshold 0 for Go.
  - [x] Reuse 37.7's `ExportRow` struct (name, email, item label, booking date); the JSON handler maps it to
        attributes and the CSV handler maps it to a `[]string`. Sorting must be identical in both outputs.
- [x] Task 2: CSV encoding helpers (AC: #2, #5)
  - [x] New file `internal/bookings/export_csv.go` with a writer that: writes the UTF-8 BOM
        (`"\xEF\xBB\xBF"`) first, then uses `encoding/csv` with `UseCRLF = true` (RFC 4180 line
        endings) and the default comma delimiter; `csv.Writer` already quotes fields containing `,`,
        `"`, CR/LF or a leading space and doubles embedded quotes.
  - [x] Header row: `Name`, `Email`, `Booked item`, `Booking date` (fixed English, independent of UI
        language; see "Open questions / assumptions").
  - [x] `escapeCSVCell(s string) string` (documented as security-critical per golang rules): if the
        value starts with `=`, `+`, `-`, `@`, TAB (`\t`) or CR (`\r`), prefix a single quote `'`
        (OWASP CSV injection guidance). Apply to every free-text cell (name, email, item label).
        Booking dates are server-generated `YYYY-MM-DD` and need no escaping.
  - [x] Check `csv.Writer.Error()` after `Flush()` and wrap errors (`fmt.Errorf("write csv: %w", err)`).
- [x] Task 3: Handler and route (AC: #1, #2, #3, #4)
  - [x] Add `ExportCSVHandler(getConfig areas.ConfigGetter, store *sql.DB) echo.HandlerFunc` next to the
        37.7 export handler. Flow: parse filters with the shared 37.7 parser (400 JSON:API on invalid
        input) -> run the unpaginated query -> only then set headers and stream the CSV. All error
        responses are written before the first CSV byte, so they stay JSON:API.
  - [x] `page[limit]` / `page[offset]` are ignored by the CSV endpoint (AC #1: all matching rows).
  - [x] Response headers: `Content-Type: text/csv; charset=utf-8`,
        `Content-Disposition: attachment; filename="sithub-bookings-<from>_<to>.csv"` (use the
        resolved/defaulted dates; they are validated `YYYY-MM-DD`, so no header injection is possible),
        `Cache-Control: no-store` (personal data).
  - [x] Register in `internal/startup/server.go` directly beside the 37.7 route with the same middleware
        slice `exportAccess` that 37.7 defines (`requireAuth` + `middleware.RequireDataExporter`):
        `e.GET("/api/v1/bookings/export.csv", bookings.ExportCSVHandler(getConfig, store), requireAuth,
        requireDataExporter)`. Echo matches the static path before `/bookings/:id` (only PATCH/DELETE
        exist there anyway). Administrators without `is_data_exporter` get 403.
  - [x] Terminal response writes in top-level funcs need `//nolint:wrapcheck // Terminal response`
        (pattern from commit 74529fb). Log only counts/filters with `slog`, never names or emails.
- [x] Task 4: Go tests (AC: #1-#5)
  - [x] Table test for `escapeCSVCell`: `=SUM(A1)`, `+1`, `-1`, `@cmd`, `\tx`, `\rx`, plain text, empty
        string, umlauts (`Jürgen Müßig` unchanged), value with `=` in the middle (unchanged).
  - [x] Handler tests (testify `require`/`assert`, `t.Parallel()`, helpers from
        `internal/bookings/testhelpers_test.go`): body starts with the BOM bytes; parse the rest with
        `csv.NewReader` and assert the header row; lines end in `\r\n`; a name containing `,` and `"`
        round-trips; dates are ISO; seeding 25 bookings with `page[limit]=20` still returns 25 rows;
        sort order equals the JSON endpoint; area filter; guest and on-behalf rows carry the
        booked-for person; soft-deleted user's booking included.
  - [x] Header tests: `Content-Type`, `Content-Disposition` filename for explicit and defaulted dates,
        `Cache-Control`.
  - [x] Authorization tests via the router setup used by the 37.7 tests: unauthenticated -> 401,
        regular user -> 403, admin without flag -> 403, data exporter -> 200. Invalid `from`/`to`,
        `from > to`, unknown area -> 400 with `application/vnd.api+json` body.
- [x] Task 5: OpenAPI documentation (AC: #4)
  - [x] New `api-doc/endpoints/bookings-export-csv.yaml`, referenced from `api-doc/openapi.yaml` as
        `/bookings/export.csv` next to `/bookings/export`. Reuse the 37.7 `from`/`to`/`area` parameter
        definitions via `$ref` (move them to `components/parameters` if 37.7 inlined them).
  - [x] Document `200` with `text/csv` (`schema: {type: string}`, example with BOM note, header row and
        one escaped row), response headers `Content-Disposition` and `Cache-Control`, and
        `400`/`401`/`403` with the shared `ErrorResponse` schema. State that pagination parameters are
        ignored.
  - [x] Lint: `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.
- [x] Task 6: Frontend download plumbing (AC: #1, #3)
  - [x] In `web/src/api/client.ts` extract the shared part of `apiRequest` (fetch, connection error ->
        `ApiError(CONNECTION_LOST_MESSAGE, 0)`, non-ok -> `ApiError` with `parseErrorDetail`) into a
        private `sendRequest(input, init, accept)` and add
        `apiDownload(input: RequestInfo): Promise<{ blob: Blob; filename: string | null }>` that sends
        `Accept: text/csv`, returns `response.blob()` and parses the filename from
        `Content-Disposition`. `apiRequest` keeps its exact behavior (existing tests must stay green).
  - [x] In the 37.7 API module `web/src/api/bookingExport.ts`, add `downloadBookingsCsv(filters)` that reuses the same
        query-string builder as the 37.7 list call (extract `buildExportQuery(filters)` if it is
        inlined; drop `page[...]`). Fallback filename when the header is missing:
        `sithub-bookings-${from}_${to}.csv`.
  - [x] New `web/src/utils/download.ts` with `saveBlob(blob: Blob, filename: string): void`:
        `URL.createObjectURL`, temporary `<a download>` element, `click()`, remove, and
        `URL.revokeObjectURL` afterwards. The session cookie is sent automatically by same-origin
        `fetch`; no token handling is needed.
- [x] Task 7: "Export as CSV" button on the Data export page (AC: #1, #3)
  - [x] In the 37.7 view (expected `web/src/views/DataExportView.vue`) add a `v-btn` beneath the
        `v-data-table-server`: `data-cy="data-export-csv-button"`, `prepend-icon="mdi-download"`,
        `:loading="isExporting"`, disabled while exporting or while the filter form is invalid.
  - [x] It exports the currently applied filters (the same reactive state that feeds the table), not
        the current page. Use `useApi()` or a local try/finally for loading state.
  - [x] Feedback per `.claude/rules/feedback.md`: success snackbar (`color="success"`, 3000 ms,
        `data-cy="data-export-csv-success"`); error snackbar (`color="error"`, 6000 ms, `closable`,
        `data-cy="data-export-csv-error"`) built with the Story 37.1 helper: translated prefix plus
        untranslated backend `detail`, or `HTTP <status>` when there is none. Reuse the page's
        existing snackbar if 37.7 added one.
  - [x] i18n keys in all five locales `web/src/locales/{en,de,es,fr,uk}.json` under the 37.7
        `dataExport` block: `exportCsv` (en "Export as CSV", de "Als CSV exportieren"),
        `exportCsvSuccess` (en "CSV file downloaded", de "CSV-Datei heruntergeladen"),
        `exportCsvFailed` (en "CSV export failed. API error:", de "CSV-Export fehlgeschlagen.
        API Fehler:").
- [x] Task 8: Vitest unit and component tests (AC: #1, #3)
  - [x] `web/src/api/client.test.ts`: `apiDownload` returns blob + filename; quoted and unquoted
        `filename=` parsing; missing header -> `null`; non-ok -> `ApiError` with `detail`; network
        failure -> status 0. Existing `apiRequest` tests unchanged.
  - [x] Test for `downloadBookingsCsv`: URL contains `from`, `to`, `area` and no `page[...]`; fallback
        filename.
  - [x] `web/src/utils/download.test.ts`: anchor gets `download` + blob URL, `click` called,
        `revokeObjectURL` called.
  - [x] View test: button renders beneath the table, calls the download with the applied filters,
        shows the loading state, success snackbar on success, error snackbar with backend detail on
        failure. Coverage stays >= 75%.
- [x] Task 9: Cypress E2E `web/cypress/e2e/data-export-csv.cy.ts` (AC: #1-#5)
  - [x] Setup via API (no UI, no mocks): log in as `anna@sithub.local` (admin), grant herself
        `is_data_exporter` through `PATCH /api/v1/users/:id` (37.6), create more than 20 bookings in
        the future within the booking limit (e.g. several items x days, including a guest booking with
        an umlaut name and a guest booking named `=HYPERLINK("x")`). Clean up created bookings and
        revoke the flag in `after`.
  - [x] Visit the Data export page, set the date range to the seeded dates, intercept
        `GET /api/v1/bookings/export.csv*` as `@exportCsv`, click `data-export-csv-button`, wait for
        `@exportCsv`, assert status 200 and the `Content-Disposition` filename.
  - [x] Verify the downloaded file with
        ``cy.readFile(`${Cypress.config('downloadsFolder')}/sithub-bookings-<from>_<to>.csv`, null)``
        (binary buffer): first three bytes are the BOM; decoded text has the header row, row count
        equals the table's total (more than one page), umlauts intact, the formula guest name starts
        with `'=`, dates match `^\d{4}-\d{2}-\d{2}$`.
  - [x] 403 check: revoke the flag, then `cy.request({ url: '/api/v1/bookings/export.csv', failOnStatusCode:
        false })` returns 403 (also covers "admin without flag").
  - [x] Add `web/cypress/downloads/` to `.gitignore` next to the existing Cypress artifact entries.
- [x] Task 10: Quality gates
  - [x] Go: `go fmt ./...`, `go vet ./...`, `golangci-lint run ./...`, `go test ./...`,
        `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exit-code 1`.
  - [x] Web (in `web/`): `npm run type-check`, `npm run lint`, `npx vitest run --coverage`,
        `npm run build`, `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0
        --exit-code 1`, `npm run test:e2e -- --browser electron` against a running dev stack.
  - [x] `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.

### Review Findings

- [x] [Review][Patch] Preserve embedded CR/LF in CSV fields while retaining CRLF record separators [internal/bookings/export_csv.go:32].
- [x] [Review][Patch] Map interrupted response-body downloads to the connection-lost ApiError [web/src/api/client.ts:84].
- [x] [Review][Patch] Add a working accessible close action to CSV error feedback [web/src/views/DataExportView.vue:99].
- [x] [Review][Patch] Exercise selected-area downloads and all formula prefixes in E2E [web/cypress/e2e/data-export-csv.cy.ts:58].
- [x] [Review][Defer] Failed table filter requests retain stale rows [web/src/views/DataExportView.vue:218] — deferred, pre-existing in 37.7.
- [x] [Review][Defer] Count and page queries lack a shared snapshot [internal/bookings/export_store.go:65] — deferred, pre-existing in 37.7; CSV ignores count.
- [x] [Review][Defer] Transient database errors return 401 [internal/middleware/auth.go:22] — deferred, pre-existing authentication behavior.

## Dev Notes

### Dependency on Stories 37.6 and 37.7

This story assumes 37.6 (`users.is_data_exporter`, a data-exporter guard middleware, `is_data_exporter`
in `GET /api/v1/me`) and 37.7 (Data export page, `GET /api/v1/bookings/export` with
`from`/`to`/`area` and `page[limit]`/`page[offset]`, pagination helpers from 37.4 in `internal/api`) are
done. Symbol names follow the 37.7 story file (`parseExportParams`, `ExportFilter`, `ExportRow`,
`ListExportBookings`, `exportAccess`, `web/src/api/bookingExport.ts`); if 37.7 shipped different names,
follow the code. Do not introduce a second filter parser, query, or name/email resolution.
[Source: _bmad-output/planning-artifacts/epics.md#Story 37.6, 37.7]

### Backend: where things live today

- Booking routes are registered in `registerRoutes` in `internal/startup/server.go` (bookings block
  ~268-276, admin users block with `requireAdmin` ~297-301). The CSV route goes next to the 37.7 export
  route with the same guard chain.
- Existing date-range handler pattern to mirror for validation messages:
  `HistoryHandlerDynamic` [Source: internal/bookings/handler.go:334-371] (`time.Parse(time.DateOnly, ...)`,
  `api.WriteBadRequest`).
- Item location and label: items are YAML-configured, not in SQLite. `cfg.FindItemLocation(itemID)`
  returns area, item group and item [Source: internal/areas/config.go:62-85];
  `collectAreaItemIDs(area)` gives the item IDs used to filter by area
  [Source: internal/bookings/handler.go:889-899]. Bookings whose item is no longer in the config are
  skipped with a `slog.Warn` in `writeBookingsCollection` (~409-413); keep the export consistent with
  whatever 37.7 decided for that case.
- Booking columns: id, item_id, user_id, booked_by_user_id, booking_date, is_guest, guest_name,
  guest_email, note, created_at, updated_at; `booking_date` is already stored as `YYYY-MM-DD`
  [Source: internal/bookings/store.go:64-115].
- JSON:API error writers live in `internal/api/errors.go` (`WriteBadRequest`, `WriteForbidden`, ...);
  they return nil after writing. `internal/api/write.go` only knows JSON; the CSV writer is
  bookings-specific and stays in `internal/bookings` (no new shared type needed).
- Security middleware (`secureConfig`, `contentSecurityPolicy` in `internal/startup/server.go`
  ~162-203) sets `X-Content-Type-Options: nosniff`; the explicit `text/csv` content type is therefore
  required. The CSP does not block `blob:` anchor downloads.

### CSV format details

- BOM: Excel only detects UTF-8 when the file starts with `EF BB BF`. Write it once before the header,
  not through `csv.Writer`.
- RFC 4180: comma delimiter, CRLF line endings (`csv.Writer.UseCRLF = true`), fields with comma, quote
  or line break are quoted and embedded quotes doubled. `encoding/csv` does all of this; do not hand-roll
  quoting.
- Formula injection: the leading `'` makes spreadsheets treat the cell as text. Because `csv.Writer`
  quotes on demand, a value like `=1+2,3` becomes `"'=1+2,3"`. Do not strip characters; finance must
  still see the original value after the apostrophe.
- Row source: one row per booking, same order as the table (date, then name). Name and email are those
  of the person the booking is for (guest name/email for guests; on-behalf -> the colleague; soft-deleted
  users included), item column `"<item group> / <item>"` (e.g. `Cube 1 / Tisch 3`).

### Frontend: download path

- `apiRequest` always sets `Accept`/`Content-Type: application/vnd.api+json` and calls
  `response.json()` [Source: web/src/api/client.ts:43-65], so it cannot download a file. Add
  `apiDownload` by extracting the shared fetch/error part; the duplicated error block in
  `cancelBooking` [Source: web/src/api/bookings.ts:191-203] shows what not to repeat.
- Use `fetch` + `Blob` + object URL rather than a plain `<a href>` to the endpoint: a navigation would
  replace the SPA with a raw JSON error page on 400/403, while `fetch` lets the page show the 37.1 error
  snackbar. Authentication is the `sithub_user` session cookie, sent automatically by same-origin
  `fetch` [Source: internal/middleware/session.go].
- The Data export menu entry and route are gated by 37.6/37.7 (`is_data_exporter`); nothing to add here.

### Feedback and i18n

Snackbar rules: success 3000 ms, error 6000 ms + `closable`, `location="bottom"`, always a `data-cy`
[Source: .claude/rules/feedback.md]. Error text follows 37.1: translated prefix, newline, untranslated
backend `detail` or `HTTP <status>` [Source: epics.md#Story 37.1]. Locales: `en`, `de`, `es`, `fr`, `uk`
in `web/src/locales/`.

### Testing standards summary

Go: table-driven tests with testify, `t.Parallel()`, `require.NoError` for setup
[Source: .claude/rules/golang.md]. Vitest coverage >= 75% [Source: .claude/rules/vue.md]. Cypress E2E
covers every AC against the real dev server; intercepts only for waiting/asserting, `data-cy`
selectors, setup and cleanup via API [Source: .claude/rules/cypress.md]. The E2E runner seeds demo
users from `tools/database/demo-users.sql` (`anna@sithub.local` is admin, not data exporter) and passes
`testUserEmail`/`testUserPassword` [Source: run-all-tests.sh:136-143]; `cy.login()` caches sessions
[Source: web/cypress/support/commands.ts:38-63]. Cypress downloads land in `cypress/downloads` by default
(`downloadsFolder`), which is trashed before each run; Electron supports blob downloads.

### Forward compatibility with Story 37.10

37.10 lets API keys with scope "read bookings" call both `/bookings/export` and `/bookings/export.csv`.
Keep the CSV route inside the same guard/group as the JSON export route so 37.10 only has to extend one
place. Do not add bearer handling here.

### Open questions / assumptions

#### Header labels

The header row uses fixed English labels `Name,Email,Booked item,Booking date`, independent of the UI
language, so third-party imports (and API-key clients in 37.10) get a stable format. Recommended
default; a localized header would require passing the locale and is not requested.

#### Delimiter

Comma per RFC 4180 (AC #2). German Excel opens comma files via the import wizard; no `sep=,` line is
added because it breaks the BOM detection in Excel and is not RFC 4180.

#### Pagination parameters

`page[limit]`/`page[offset]` are ignored by the CSV endpoint (AC #1 "all bookings"). There is no row
cap; the date range bounds the volume. If 37.7 introduced a maximum range, it applies to CSV as well.

#### Escape character set

Besides the four characters required by AC #5, leading TAB and CR are escaped too (OWASP
recommendation). The escape is a leading apostrophe.

#### Empty result

The button stays enabled when the table is empty; the export then contains only the header row.

#### File name

Only start and end date are in the name (`sithub-bookings-<from>_<to>.csv`); the area is not included.

### Project Structure Notes

- New: `internal/bookings/export_csv.go` (+ `export_csv_test.go`),
  `api-doc/endpoints/bookings-export-csv.yaml`, `web/src/utils/download.ts` (+ test),
  `web/cypress/e2e/data-export-csv.cy.ts`.
- Modified: 37.7 export handler/store files in `internal/bookings/` (unpaginated reuse only),
  `internal/startup/server.go`, `api-doc/openapi.yaml`, `web/src/api/client.ts` (+ test), the 37.7
  export API module and view, `web/src/locales/*.json` (5 files), `.gitignore`.
- No migration, no new Go dependency (`encoding/csv` is stdlib), no new npm dependency.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.8 (FR194), #Story 37.7 (FR193),
  #Story 37.6 (FR188), #Story 37.10]
- [Source: internal/startup/server.go:162-203 (security headers, CSP), 247-301 (routes)]
- [Source: internal/bookings/handler.go:334-371 (history date filters), 374-425, 889-899]
- [Source: internal/bookings/store.go:64-115]
- [Source: internal/areas/config.go:62-85]
- [Source: internal/api/write.go, internal/api/errors.go, internal/api/response.go]
- [Source: web/src/api/client.ts:1-65, web/src/api/bookings.ts:158-203]
- [Source: web/cypress/support/commands.ts:38-63, run-all-tests.sh:136-143,
  tools/database/demo-users.sql]
- [Source: .claude/rules/golang.md, json-api.md, apidoc.md, vue.md, cypress.md, feedback.md]
- RFC 4180 (CSV), OWASP "CSV Injection"

## Dev Agent Record

### Agent Model Used

GPT-6 Astra Fast (OpenAI)

### Debug Log References

- User decision: CSV requires explicit `from` and `to`, matching JSON export; no default dates.
- Continuation decision: use the existing project duplication gates (Go <= 3%, excluding test files
  and node_modules; TypeScript 0%, excluding node_modules and unit tests), superseding Task 10's
  literal duplication commands. No project-wide threshold changes.
- Shared filter validation is in `parseExportRequest`; `parseExportParams` preserves JSON pagination.
- Followed the actual 37.7 API location, `bookings.ts`, rather than creating a second export module.
- Used the registered SVG `$download` alias; omitted unused lint suppressions (targeted lint passes).
- New client-download tests live in `api/download.test.ts`; existing client tests remain unchanged.
- Full Cypress run against isolated Vite/backend: 103/104 passed. Breadcrumb request wait failed once;
  the isolated breadcrumb rerun passed 4/4 without code changes. CSV tests passed both runs.
- Initial gate failures were outside this story: gosec G124 at `internal/middleware/auth.go:27`;
  strict Go duplication 1.35% (12 existing clones); TS duplication 0.12% between existing
  `data-export.cy.ts` and `week-navigation.cy.ts`. No CSV clone reported. Resolved as recorded below.
- Inactive-user cookie expiry now uses `auth.Service.NewCookie`, preserving HTTPS and
  `force_secure_cookies` behavior. Regression tests first failed for HTTPS/proxy, then passed after
  the fix. G124 still flagged the dynamic factory result; a scoped, documented `#nosec G124`
  explains the tested TLS/config-driven Secure flag (plain HTTP remains supported).
- Both E2E specs now reuse `toLocalISODate`; frontend duplication is 0%. Go duplication is 1.35%,
  below the user-approved existing project threshold of 3%.
- Final full Cypress run on isolated Vite :5188 and fresh backend :9918: 104/104 passed across
  23 specs, including CSV download and breadcrumbs. No retry or test changes needed for breadcrumbs.

### Completion Notes List

- Three independent review layers completed; four patch findings fixed, three pre-existing findings
  recorded in deferred-work.md, one dismissed: buffering all rows is the explicitly specified query
  design, with no demonstrated capacity failure. No decision-needed findings.
- CSV review correction: stdlib UseCRLF removes bare CR inside fields. Encode one record with
  encoding/csv's default field handling, then replace only its final LF with CRLF; raw-byte regression
  tests cover CR, LF, CRLF and leading CR formula protection. This supersedes the Task 2 UseCRLF detail.
- Download body failures now produce status-0 ApiError. CSV errors have a translated close action;
  tested with actual VSnackbar/VBtn. Corrected Vitest's Vuetify dependency-inlining configuration
  location to enable the real-component test.
- Expanded real-backend E2E checks selected-area exclusion and all four formula prefixes. Exact
  pathname intercepts prevent the JSON alias from accidentally matching CSV requests.
- Post-fix verification: all Go tests, vet and golangci-lint pass; 694 Vitest tests pass (82.4%
  statements); type-check, ESLint, build, duplication gates and diff whitespace checks pass;
  CSV Electron tests 2/2 pass. Earlier full E2E run was 104/104; only the affected CSV spec was rerun
  after review fixes. Story and sprint status are done.

- Implemented shared unpaginated export, CSV BOM/CRLF/quoting, formula escaping, dated filename,
  no-store response and shared exporter authorization. Explicit dates required; page params ignored.
- Added filtered blob download with fallback filename, object-URL cleanup, loading/invalid-filter
  guards, translated success/error feedback in five languages, and real browser download tests.
- Verified: all Go tests and vet pass; targeted bookings/startup golangci-lint has 0 issues; web
  type-check, ESLint and build pass; 692 Vitest tests pass (82.39% statements); OpenAPI valid with
  one existing auth-providers warning; CSV Cypress 2/2 passes with 26 seeded rows and byte checks.
- Final verification: Go vet and repository-wide golangci-lint pass (0 issues); race-enabled Go
  tests pass with 87.1% total coverage; frontend type-check, lint and build pass; 692 Vitest tests
  pass with 82.39% statement coverage; Electron E2E 104/104 passes; approved duplication gates pass;
  OpenAPI valid with the existing auth-providers warning. All tasks complete; ready for review.

### File List

- internal/bookings/export.go
- internal/bookings/export_store.go
- internal/bookings/export_csv.go (new)
- internal/bookings/export_csv_test.go (new)
- internal/startup/server.go
- internal/startup/server_test.go
- internal/middleware/auth.go
- internal/middleware/auth_test.go
- api-doc/openapi.yaml
- api-doc/endpoints/bookings-export.yaml
- api-doc/endpoints/bookings-export-csv.yaml (new)
- web/src/api/client.ts
- web/src/api/bookings.ts
- web/src/api/download.test.ts (new)
- web/src/api/bookingDownload.test.ts (new)
- web/src/utils/download.ts (new)
- web/src/utils/download.test.ts (new)
- web/src/views/DataExportView.vue
- web/src/views/DataExportView.test.ts
- web/vitest.config.ts
- web/src/components/DatePickerField.vue
- web/src/plugins/vuetify.ts
- web/src/locales/en.json
- web/src/locales/de.json
- web/src/locales/es.json
- web/src/locales/fr.json
- web/src/locales/uk.json
- web/cypress/e2e/data-export-csv.cy.ts (new)
- web/cypress/e2e/data-export.cy.ts
- web/cypress/e2e/week-navigation.cy.ts
- .gitignore
- _bmad-output/implementation-artifacts/37-8-export-bookings-as-csv.md
- _bmad-output/implementation-artifacts/sprint-status.yaml
- _bmad-output/implementation-artifacts/deferred-work.md

## Change Log

- 2026-10-08: Completed three-layer code review and automatically fixed four findings; recorded three
  pre-existing deferred issues and verified regression checks. Status moved to done.
- 2026-10-08: Completed Task 10 using user-approved project duplication gates; fixed cookie-security
  policy on inactive-user expiry and reused the E2E date formatter. Full regression checks pass.
  Story and sprint status moved to review.
- 2026-10-08: Implemented CSV export and verified download end to end. Tasks 1-9 complete;
  Task 10 blocked by existing repository-wide lint/duplication findings. Status remains in-progress.
