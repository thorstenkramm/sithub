# Story 36.4: Application Version Reporting

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an operator and user,
I want to know which SitHub version is running,
so that I can verify deployments and report issues accurately.

## Acceptance Criteria

1. `sithub version` CLI command prints the version string and exits with code 0.
2. The GitHub release workflow injects the release tag as the build version (Go ldflags);
   non-release builds report a sensible fallback such as `dev`.
3. The API exposes the running version in a JSON:API-consistent response.
4. The UI shows the version in the user settings view (the user menu).

## Tasks / Subtasks

- [x] Task 1: Add a package-level `version` variable and a `version` CLI command (AC: #1, #2)
  - [x] In `cmd/sithub/main.go`, declare a package-level `var version = "dev"` (default fallback).
        GoReleaser's default ldflags target `main.version`, so this exact variable name and package
        (`package main`) matter — see Dev Notes.
  - [x] Register a `versionCmd` (`cobra.Command`, `Use: "version"`) alongside the existing `runCmd`
        via `rootCmd.AddCommand(...)`. Its `Run` prints the version to stdout and returns nil
        (cobra exits 0). Extracted into `newVersionCmd()` for testability.
  - [x] Print exactly the version string via `cmd.Println(version)` — a bare, parseable string.
- [x] Task 2: Confirm/adjust release-tag ldflags injection (AC: #2)
  - [x] Verified `.goreleaser.yml` defines NO `ldflags:` under `builds:` (grep found none), so
        GoReleaser's default `-X main.version={{.Version}} ...` applies and injects the tag into
        `main.version`. No edit made or required.
  - [x] Confirmed no change to `.github/workflows/release.yml` — triggers on `[0-9]+.[0-9]+.[0-9]+`
        tag push with `fetch-depth: 0` and runs `goreleaser release --clean`.
  - [x] Did NOT touch `build.sh` — the `"dev"` default already satisfies the fallback (decision:
        avoid over-engineering, per task guidance).
- [x] Task 3: Expose the version via the API (AC: #3)
  - [x] Added `Version(version string) echo.HandlerFunc` in `internal/system/version.go`, mirroring
        `SettingsHandler`: `VersionAttributes{ Version string \`json:"version"\` }`,`Type: "version"`,
        `ID: "version"`,`api.SingleResponse`, JSON:API content type header,`c.JSON`.
  - [x] Registered `e.GET("/api/v1/version", system.Version(version), requireAuth)` next to settings.
  - [x] Threaded `version` from `main` -> `startup.Run(ctx, cfg, version)` ->
        `registerRoutes(..., version)` (least-invasive: mirrors how `bookingLimits` is threaded;
        no global state).
- [x] Task 4: Show the version in the user settings menu (AC: #4)
  - [x] Added `web/src/api/version.ts` mirroring `api/settings.ts`.
  - [x] In `App.vue`, `loadVersion()` fetches once on auth into `appVersion` ref; rendered as a
        low-emphasis caption line in BOTH the desktop `v-menu` (`data-cy="app-version"`) AND the mobile
        drawer (`data-cy="mobile-app-version"`), guarded by `v-if="appVersion"`.
  - [x] Added `app.userMenu.version` to all five locale files
        (en/de/fr: "Version", es: "Versión", uk: "Версія").
- [x] Task 5: Tests (AC: #1-#4)
  - [x] Go: `TestVersionCommandPrintsVersion` executes `newVersionCmd()` and asserts stdout.
        `TestVersion` handler test mirrors `ping_test.go` (status 200, JSON:API content type,
        `data.type == "version"`, attribute equals injected value).
  - [x] Frontend: `web/src/api/version.test.ts` for `fetchVersion()` (mirrors `settings.test.ts`);
        `App.test.ts` mocks `./api/version` and asserts `[data-cy="app-version"]` renders `1.2.3`.
  - [ ] Optional Cypress: NOT added (decision: covered by Vitest App render test + fetch unit test;
        the optional E2E was deemed unnecessary for this low-risk display line).

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md#Story 36.4` (epics.md:5438-5463) and FR169-FR172
(epics.md:624-632). Epic 36 is "User Feedback — Persistent Sessions, Versioning & Booking UX"
(epics.md:1160-1173).

### CLI: single package-level version variable + `version` command (AC #1, #2)

`cmd/sithub/main.go` is a small cobra program: a `rootCmd` (`Use: "sithub"`) with one subcommand
`runCmd` (`Use: "run"`) added via `rootCmd.AddCommand(runCmd)` (main.go:16-44). There is currently
NO version variable anywhere in `cmd/` or `internal/` (grep found none). Add:

```go
// version is set at build time via ldflags (-X main.version=...); "dev" for non-release builds.
var version = "dev"
```

Then a sibling command:

```go
versionCmd := &cobra.Command{
    Use:   "version",
    Short: "Print the SitHub version",
    Run: func(cmd *cobra.Command, _ []string) {
        cmd.Println(version)
    },
}
rootCmd.AddCommand(versionCmd)
```

cobra returns exit code 0 when a command's `Run`/`RunE` succeeds; the existing `os.Exit(1)` only
fires on `rootCmd.Execute()` error (main.go:42-44), so `version` naturally exits 0.

### Build injection: GoReleaser defaults already target `main.version` (AC #2)

`.goreleaser.yml` (repo root) has `builds:` with `main: ./cmd/sithub`, `binary: sithub`, and NO
custom `ldflags:` block. GoReleaser's DEFAULT build ldflags are:

```text
-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser
```

Because the CLI's `var version` lives in `package main` under `./cmd/sithub`, GoReleaser's default
`-X main.version={{.Version}}` populates it with the git tag automatically — no `.goreleaser.yml`
edit is strictly required. `.github/workflows/release.yml` already: triggers on tags matching
`[0-9]+.[0-9]+.[0-9]+` (release.yml:5-6), checks out with `fetch-depth: 0` so tags are available
(release.yml:18-19), and runs `goreleaser release --clean` (release.yml:41-45).

> [!IMPORTANT]
> If anyone later adds a custom `ldflags:` block to `.goreleaser.yml`, it overrides the defaults and
> MUST re-include `-X main.version={{.Version}}`, or the release tag will stop being injected.

The `"dev"` default in code covers the non-release / `go build` / `build.sh` path (build.sh builds
with a plain `go build -o ... ./cmd/sithub` and no ldflags, so `version` stays `"dev"`).

### API: new `/api/v1/version` endpoint, JSON:API-consistent (AC #3)

Follow the established `system` package pattern. `SettingsHandler` (internal/system/settings.go:16-35)
is the closest template — it is a closure that captures a runtime value (`weeksInAdvanced`) and
returns an `echo.HandlerFunc` producing an `api.SingleResponse` with the JSON:API content type. Do
the same for the version:

- `VersionAttributes struct { Version string \`json:"version"\` }` (snake_case per json-api.md; single
  word `version` is already snake_case).
- Resource `Type: "version"`, `ID: "version"`, `Attributes: VersionAttributes{...}`.
- Set `c.Response().Header().Set(echo.HeaderContentType, api.JSONAPIContentType)` then
  `c.JSON(http.StatusOK, resp)`, wrapping any error with `%w` (settings.go:29-33).

Shared response types are in `internal/api/response.go`: reuse `api.SingleResponse`,
`api.Resource`, `api.JSONAPIContentType` (response.go:14,16-26) — do NOT define new envelope types
(per golang.md "Shared Types").

Route registration goes in `internal/startup/server.go`. Public `ping` is at server.go:243; the
authenticated `settings` route is at server.go:251 (`e.GET("/api/v1/settings",
system.SettingsHandler(weeksInAdvanced), requireAuth)`). Register version the same way with
`requireAuth` since the UI consumes it while logged in:

```go
e.GET("/api/v1/version", system.Version(version), requireAuth)
```

Wiring the value: `weeksInAdvanced` reaches the handler as a plain int derived inside the server
builder (server.go:247-251). The `version` string originates in `package main`. Read `startup.Run`
and the server builder signatures and pass `version` through the same path the config/limits use;
prefer adding it to the existing config/args flow rather than introducing global state (golang.md:
"avoid global state beyond config wiring"). `main.go` calls `startup.Run(cmd.Context(), cfg)`
(main.go:32) — the cleanest option is to add a `version` field carried alongside `cfg`, or extend
`startup.Run`'s signature; choose whichever touches the fewest call sites while staying explicit.

### Frontend: user settings view = the App.vue user menu (AC #4)

There is NO standalone "SettingsView.vue". User settings live in `App.vue` as the user menu, which
appears twice:

- Desktop `v-menu` (App.vue:40-147): contains theme selector (`data-cy="theme-selector"`, ~97),
  language selector (~116), show-weekends toggle (~131), a `v-divider`, and the logout item
  (`data-cy="logout-btn"`, ~140-145).
- Mobile `v-navigation-drawer` (App.vue:158-267): the same controls repeated with `mobile-*`
  `data-cy` values, ending in `data-cy="mobile-logout-btn"` (~260).

Add the version line to BOTH so desktop and mobile match. Use a low-emphasis caption style
consistent with existing labels (e.g. `class="text-caption text-medium-emphasis"`, as used at
App.vue:98). Give it `data-cy="app-version"` (and a `mobile-app-version` variant if you prefer
symmetry with the other mobile ids).

API client: add `web/src/api/version.ts` mirroring `web/src/api/settings.ts` (settings.ts:1-10),
using `apiRequest` from `./client` and `SingleResponse` from `./types` (types.ts:1-9). Attribute
key is `version` (snake_case matches the Go struct tag). Fetch once in `App.vue`'s `<script setup>`
(the file already imports preference composables at App.vue:403-405 and uses `computed`/`ref`);
store the result in a `ref<string>` and guard the render so nothing shows before it loads.

i18n: `App.vue` labels use `$t('app.userMenu.*')` (App.vue:98,117,134,144). The `userMenu` block
exists in all five locale files under `web/src/locales/` (`de,en,es,fr,uk.json`; loaded via
`web/src/plugins/i18n.ts`). Add a `version` key to each. English value: `"Version"`.

### Project Structure Notes

Files to add/modify:

- Modify: `cmd/sithub/main.go` (version var + version command).
- Verify only: `.goreleaser.yml`, `.github/workflows/release.yml` (defaults already inject the tag).
- Add: `internal/system/version.go` (+ `internal/system/version_test.go`).
- Modify: `internal/startup/server.go` (route + thread `version` through) and whatever `startup.Run`
  / server-builder signature the wiring requires.
- Add: `web/src/api/version.ts` (+ `web/src/api/version.test.ts`).
- Modify: `web/src/App.vue` (fetch + render in both menus) and its test.
- Modify: `web/src/locales/{en,de,es,fr,uk}.json` (add `app.userMenu.version`).

No database migration and no schema change — version is compile-time metadata, not persisted.

### Testing standards summary

Go: table-driven where useful, `require`/`assert`, mirror `internal/system/ping_test.go` (ping_test.go:14-39)
for the handler test (status 200, `api.JSONAPIContentType`, decode into `api.SingleResponse`, assert
`Data.Type` and attribute). Run `golangci-lint run ./...`, `go fmt ./...`, `go vet ./...`,
and the jscpd Go duplication check. [Source: .claude/rules/golang.md]

Frontend: Vitest for `fetchVersion()` (mirror `api/settings.test.ts`) and an `App.vue` render test.
Run `npm run type-check`, `npm run lint`, `npx vitest run`, `npm run build`; coverage stays >= 75%.
Optional Cypress E2E using the existing `cy.login()` command. [Source: .claude/rules/vue.md,
.claude/rules/cypress.md]

### References

- [Source: cmd/sithub/main.go:14-44]
- [Source: .goreleaser.yml (root) builds/archives sections]
- [Source: .github/workflows/release.yml:1-6,18-19,41-45]
- [Source: internal/system/settings.go:12-35]
- [Source: internal/system/ping.go:14-33]
- [Source: internal/system/ping_test.go:14-39]
- [Source: internal/api/response.go:14,16-31]
- [Source: internal/startup/server.go:243,247-251]
- [Source: web/src/api/settings.ts:1-10]
- [Source: web/src/api/client.ts:1-40; web/src/api/types.ts:1-9]
- [Source: web/src/App.vue:40-147,158-267,403-405]
- [Source: web/src/locales/en.json (app.userMenu), web/src/plugins/i18n.ts]
- [Source: _bmad-output/planning-artifacts/epics.md:5438-5463,624-632]
- [Source: .claude/rules/json-api.md; .claude/rules/golang.md; .claude/rules/vue.md]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- Go gate: `go build ./...`, `go vet ./...`, `gofmt -l` (clean), `go test ./...` (all pass),
  `golangci-lint run ./...` -> `0 issues.`
- Frontend gate: `npm run type-check` (clean), `npm run lint` (0 warnings),
  `npx vitest run` -> 51 files / 469 tests pass, `npm run build` -> built OK.
- jscpd: TS check shows no clones involving the new `version.ts`/`version.test.ts` (pre-existing
  4.7% elsewhere). Go check reports one clone between `version.go` and `ping.go` — the shared
  JSON:API response-writing tail — which is the intentional, story-mandated system-package pattern
  (same idiom already duplicated across `ping.go`/`settings.go`).

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created.
- CLI: `var version = "dev"` in `package main`; `version` subcommand via `newVersionCmd()` (extracted
  for unit testability) prints the bare version and exits 0.
- Build injection: verified `.goreleaser.yml` has no custom `ldflags:`, so GoReleaser's default
  `-X main.version={{.Version}}` populates `main.version` from the git tag. No release-config edits.
- API: `GET /api/v1/version` (authenticated, `requireAuth`) returns JSON:API single resource
  `{ data: { type: "version", id: "version", attributes: { version } } }`.
- Wiring: `version` threaded via `startup.Run(ctx, cfg, version)` -> `registerRoutes(..., version)`;
  updated all 6 test call sites accordingly (5 `registerRoutes`, 1 `Run`).
- UI: version shown as a low-emphasis caption in both desktop menu and mobile drawer, fetched once
  on authentication; hidden until loaded.
- Locale note: es="Versión", uk="Версія", en/de/fr="Version".

### File List

- Modified: `cmd/sithub/main.go`
- Modified: `cmd/sithub/main_test.go`
- Added: `internal/system/version.go`
- Added: `internal/system/version_test.go`
- Modified: `internal/startup/server.go`
- Modified: `internal/startup/server_test.go`
- Added: `web/src/api/version.ts`
- Added: `web/src/api/version.test.ts`
- Modified: `web/src/App.vue`
- Modified: `web/src/App.test.ts`
- Modified: `web/src/locales/en.json`, `de.json`, `es.json`, `fr.json`, `uk.json`

### Change Log

- 2026-07-08: Story drafted (ready-for-dev) — FR169-FR172 application version reporting across CLI,
  release build injection, API endpoint, and settings-menu UI.
- 2026-07-09: Implemented all tasks — CLI `version` command + `var version`, `GET /api/v1/version`
  authenticated endpoint, App.vue user-menu display (desktop + mobile), i18n keys in all five
  locales, and Go + Vitest tests. All gates green. Status -> review.
