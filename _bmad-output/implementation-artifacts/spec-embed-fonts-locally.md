---
title: 'Embed all frontend fonts locally'
type: 'bugfix'
created: '2026-11-25T15:20:05Z'
status: 'done'
baseline_commit: '9b7fd545c0b49134ae2f67c29e73b098b7d26918'
context:
  - '../../README.md'
  - '../../docs/development-guide.md'
  - '../planning-artifacts/architecture.md'
---

<frozen-after-approval reason="human-owned intent — do not modify unless human renegotiates">

## Intent

**Problem:** The application downloads Inter stylesheets and font files from Google Fonts, disclosing user network metadata to a third party and making typography dependent on external connectivity. This also contradicts SitHub's self-contained, single-binary distribution model.

**Approach:** Bundle every webfont used by the application into the frontend build and resulting Go binary, remove all Google Fonts references, and restrict Content Security Policy to same-origin or data-font sources. Preserve the existing Inter visual design, weights, languages, and already-local Material Design Icons.

## Boundaries & Constraints

**Always:**
- Serve every application webfont same-origin from assets embedded in `assets.Web`; preserve normal Inter 400/500/600/700 and MDI.
- Retain full language subsets, including Cyrillic for Ukrainian, plus `font-display: swap` and existing system fallbacks.
- Make no requests to `fonts.googleapis.com` or `fonts.gstatic.com` in development, production, or the binary; add no alternative remote provider.
- Rebuild generated `assets/web` output from source; never hand-edit or commit it.

**Ask First:** Changing font families, weights, languages, styles, or visual design.

**Never:** Retain Google URLs as fallbacks, permit external fonts in CSP, or modify unrelated BMAD/prototype templates with illustrative Google URLs.

## I/O & Edge-Case Matrix

| Scenario | Input / State | Expected Output / Behavior | Error Handling |
|----------|--------------|---------------------------|----------------|
| HAPPY_PATH | Production SPA, cache disabled | Inter and MDI load from SitHub; zero Google requests | N/A |
| OFFLINE_CLIENT | External network unavailable, SitHub reachable | All weights/icons render without external requests or CSP violations | Existing system fallbacks cover local asset failure |
| CYRILLIC_TEXT | Ukrainian text at supported weights | Local Cyrillic subset loads same-origin | No remote fallback |
| BINARY_DISTRIBUTION | Rebuild frontend before Go compile | `assets.Web` contains Inter WOFF2 and MDI assets | Test fails when assets are absent |
| CSP_RESPONSE | Response with application CSP | `font-src` is `'self' data:`; no external font style host | Tests reject both Google hosts |

</frozen-after-approval>

## Code Map

- `web/index.html` -- Google preconnects/stylesheet; sole runtime CDN entry point.
- `web/src/main.ts` -- CSS entry point already bundling Vuetify and MDI styles.
- `web/package.json` and `web/package-lock.json` -- Dependency declaration and reproducible resolution.
- `web/src/styles/global.css` -- Existing Inter stack and weight tokens; preserve unchanged.
- `internal/startup/server.go` and `server_test.go` -- CSP generation and its security regression coverage.
- `assets/embed.go` and `embed_test.go` -- Existing frontend `embed.FS` and asset test location.
- `web/cypress/e2e/fonts.cy.ts` -- New focused browser privacy and font-loading coverage.
- `web/vite.config.ts` -- Existing hashed same-origin output into `assets/web`; no changes needed.

## Tasks & Acceptance

**Execution:**
- [x] `web/package.json`, `web/package-lock.json`, and `web/src/main.ts` -- Add static Inter and import normal 400/500/600/700 full-subset CSS -- Bundle all existing weights, including Cyrillic and swap behavior.
- [x] `web/index.html` -- Remove Google preconnects and stylesheet -- Eliminate third-party font loading.
- [x] `internal/startup/server.go` -- Remove Google CSP sources and update comments -- Enforce local fonts/styles.
- [x] `internal/startup/server_test.go` -- Assert exact local-only directives and forbidden Google hosts -- Prevent CSP regressions.
- [x] `assets/embed_test.go` -- Assert embedded Inter/MDI assets and scan embedded HTML/CSS for Google URLs -- Validate production binary input.
- [x] `web/cypress/e2e/fonts.cy.ts` -- Load the real app, verify all four Inter faces, and fail on Google requests -- Test actual browser behavior.

**Acceptance Criteria:**
- Given a fresh production binary with external networking blocked, when the SPA loads, then all webfonts come from embedded same-origin assets with zero Google requests.
- Given English or Ukrainian text at weights 400/500/600/700, when rendered, then normal Inter is locally available with swap and existing fallbacks.
- Given an application response, when CSP is inspected, then `font-src` is exactly `'self' data:` and neither Google host is allowed.
- Given a frontend rebuild, when icons render and changes are inspected, then MDI still works and no generated `assets/web` output is committed.

### Review Findings

- [x] Build the frontend before local Go checks so tests read fresh embedded assets.
- [x] Import the Cypress `Interception` type and type-check the font spec directly.
- [x] Clear the HTTP cache and block cross-origin browser requests during the font check.
- [x] Load a font-backed MDI icon and verify its versioned WOFF2 response by pathname.
- [x] Assert that visible application navigation selects Inter as its first family.
- [x] Resolve CSS font URLs against their stylesheet and verify the exact embedded path.
- [x] Check the font and style CSP directives independently of unrelated directive order.
- [x] Document the seeded test credentials required by the focused Cypress command.

## Spec Change Log

## Design Notes

Use `@fontsource/inter` static entries, matching npm-managed MDI without renaming the family or manually copying binaries. Vite emits hashed same-origin WOFF2 files; full weight CSS includes every subset, while browsers fetch only applicable unicode ranges. Rebuild before Go asset tests to avoid validating stale output.

## Verification

**Commands:**
- `cd web && npm ci && npm run type-check && npm run lint && npm run test:unit:coverage && npm run build` -- expected: all checks and asset generation pass.
- `go test -race ./assets ./internal/startup && go test -race ./...` -- expected: privacy regressions and complete backend suite pass against fresh output.
- With a running local test binary seeded from `tools/database/demo-users.sql`, run from `web/`:
  ```sh
  CYPRESS_testUserEmail=anna@sithub.local CYPRESS_testUserPassword='SitHubDemo2026!!' \
    npx cypress run --browser electron --config baseUrl=http://localhost:8080 --spec cypress/e2e/fonts.cy.ts
  ```
  The cold-cache check passed against a fresh binary: all four Latin/Cyrillic weights and MDI loaded locally.
- `rg -n 'fonts\.(googleapis|gstatic)\.com' web/index.html assets/web internal/startup` -- expected: no runtime matches.

## Suggested Review Order

**Local font packaging**

- Bundle four static Inter weights and every language subset through Vite.
  [`main.ts:4`](../../web/src/main.ts#L4)

- Keep the HTML entry point free of third-party font loaders.
  [`index.html:9`](../../web/index.html#L9)

**Privacy enforcement**

- Restrict font and stylesheet sources at the HTTP response boundary.
  [`server.go:197`](../../internal/startup/server.go#L197)

**Regression coverage**

- Verify embedded licenses, all font faces, MDI, and remote-host absence.
  [`embed_test.go:22`](../../assets/embed_test.go#L22)

- Observe real browser requests and each Latin/Cyrillic weight response.
  [`fonts.cy.ts:14`](../../web/cypress/e2e/fonts.cy.ts#L14)

- Check the local font and style directives and reject both Google hosts.
  [`server_test.go:108`](../../internal/startup/server_test.go#L108)

**Distribution metadata**

- Declare the npm-managed Inter dependency for reproducible local builds.
  [`package.json:19`](../../web/package.json#L19)

- Preserve the upstream OFL notice with redistributed Inter assets.
  [`inter-OFL-1.1.txt:1`](../../web/public/licenses/inter-OFL-1.1.txt#L1)
