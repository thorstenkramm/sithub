# Story 33.6: Entra ID Primary Login With Official Icon and "More Options" Toggle

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a corporate user landing on the login page,
I want the Entra ID button to be the obvious primary action with the official
Microsoft icon,
so that I can sign in with one click without being distracted by the local-
credentials form.

## Acceptance Criteria

1. **Given** Entra ID is configured on the server
   **When** I open the login page unauthenticated
   **Then** I see only the SitHub brand logo, the Entra ID sign-in button (showing
   the official Microsoft Entra ID color icon), and a clickable "more login
   options" link
   **And** the local username/password form is hidden by default

2. **Given** I click "more login options"
   **When** the page updates
   **Then** the local credentials form expands into view
   **And** the link text changes to "less login options"

3. **Given** the credentials form is expanded
   **When** I click "less login options"
   **Then** the form collapses again
   **And** the link text reverts to "more login options"

4. **Given** Entra ID is NOT configured on the server
   **When** I open the login page
   **Then** the local credentials form is shown by default
   **And** the Entra ID button and toggle link are not rendered, so users are not
   locked out

5. **Given** the Entra ID icon is rendered on the button
   **When** the page loads
   **Then** the icon SVG is served from the binary's embedded assets (no external
   network request to Wikimedia or any other host); the SVG source was downloaded
   from
   `https://upload.wikimedia.org/wikipedia/commons/8/8c/Microsoft_Entra_ID_color_icon.svg`
   and committed to the repository under the existing embedded-assets location

## Tasks / Subtasks

- [ ] Task 1: Download the official Entra ID icon and place it in the embedded
      assets directory (AC: #5)
  - [ ] 1.1 Fetch the SVG from
        `https://upload.wikimedia.org/wikipedia/commons/8/8c/Microsoft_Entra_ID_color_icon.svg`.
        Recommended command:
        ```
    curl -fSL -o web/public/entra-id-icon.svg \
          https://upload.wikimedia.org/wikipedia/commons/8/8c/Microsoft_Entra_ID_color_icon.svg
        ```
        Files placed under `web/public/` are copied by Vite to `web/dist/`, which
        the backend's `assets/embed.go` (`//go:embed web/* web/*/*`) embeds into
        the binary. Confirm precedent: `web/public/logo.svg` and
        `web/public/favicon.svg` are embedded the same way.
  - [ ] 1.2 Commit the file under `web/public/entra-id-icon.svg`. Do NOT
        commit a derivative under `assets/web/` — `assets/web/` is the build
        output directory (gitignored apart from the embed marker).

- [ ] Task 2: Add a backend `GET /api/v1/auth/providers` endpoint exposing whether
      Entra ID is configured (AC: #4)
  - [ ] 2.1 Add a handler that returns a JSON:API single-resource document:
        ```json
        {
          "data": {
            "type": "auth-providers",
            "id": "current",
            "attributes": {
              "entraid": true,
              "local": true
            }
          }
        }
        ```
        Source the `entraid` boolean from `Config.EntraIDConfigured()` already
        defined at `internal/config/config.go:146`. `local` is always `true`
        today (local auth is always available), but exposing it explicitly
        future-proofs the contract.
  - [ ] 2.2 Register the route in the existing public auth route group so the
        endpoint does not require authentication (the login page must reach it
        before the user is logged in). Mirror the existing
        `/api/v1/auth/login` registration pattern in `internal/auth/`.
  - [ ] 2.3 Add a Go test for the handler covering both states (configured /
        not configured).
  - [ ] 2.4 Update the OpenAPI spec under `./api-doc/` per
        `.claude/rules/apidoc.md`. Add `auth-providers.yaml` (or extend an
        existing auth file). Lint with `npx @redocly/cli lint --lint-config off
        ./api-doc/openapi.yaml`.

- [ ] Task 3: Frontend — add a typed client for the new endpoint
  - [ ] 3.1 In `web/src/api/auth.ts` (the same module that exports `loginLocal`,
        `logout`), add:
        ```ts
        export interface AuthProvidersAttributes {
          entraid: boolean;
          local: boolean;
        }
        export type AuthProvidersResponse =
          JsonApiResponse<'auth-providers', AuthProvidersAttributes>;
        export const fetchAuthProviders = () =>
          api.get<AuthProvidersResponse>('/auth/providers');
        ```
        Follow the existing API client style in the file.

- [ ] Task 4: Rebuild the login page layout (AC: #1, #2, #3, #4, #5)
  - [ ] 4.1 In `web/src/views/LoginView.vue`, refactor the template (lines 1–67)
        as follows:
        ```vue
        <template>
          <v-container class="fill-height" fluid>
            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" md="4">
                <v-card elevation="2">
                  <v-card-title class="text-center pt-6">
                    <img :src="logoSrc" alt="SitHub" class="login-logo" />
                  </v-card-title>
                  <v-card-text>
                    <!-- Entra ID primary action (only when configured) -->
                    <v-btn
                      v-if="entraIdAvailable"
                      block
                      size="large"
                      variant="outlined"
                      :loading="entraIdLoading"
                      :disabled="entraIdLoading"
                      data-cy="login-entraid"
                      class="login-entraid-btn"
                      @click="handleEntraIdLogin"
                    >
                      <template #prepend>
                        <img src="/entra-id-icon.svg" alt="Entra ID" class="entra-id-icon" />
                      </template>
                      {{ $t('auth.signInWithEntraId') }}
                    </v-btn>
                    <div
                      v-if="entraIdAvailable"
                      class="text-center mt-3"
                    >
                      <a
                        href="#"
                        class="text-caption text-medium-emphasis login-more-options"
                        data-cy="login-toggle-local"
                        @click.prevent="showLocalForm = !showLocalForm"
                      >
                        {{ showLocalForm ? $t('auth.lessLoginOptions') : $t('auth.moreLoginOptions') }}
                      </a>
                    </div>
                    <!-- Local credentials form -->
                    <v-expand-transition>
                      <div v-if="!entraIdAvailable || showLocalForm">
                        <v-divider v-if="entraIdAvailable" class="my-4" />
                        <v-form ...>
                          <!-- existing email + password + submit unchanged -->
                        </v-form>
                      </div>
                    </v-expand-transition>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </template>
        ```
        Notes:
        - `logoSrc` is supplied by Story 33.7 (defaults to `/logo.svg` until
          33.7 lands, then `/sithub_logo.svg`); for 33.6 in isolation use the
          existing `/logo.svg`.
        - `entraIdAvailable` is a ref populated on mount via `fetchAuthProviders`
          (Task 3); default `false` until resolved.
        - `showLocalForm` is a ref, default `false` when Entra ID is available,
          `true` when not.

- [ ] Task 5: Script-setup wiring (AC: #1, #2, #3, #4)
  - [ ] 5.1 In `<script setup>` of `LoginView.vue`, add:
        ```ts
        import { onMounted, ref } from 'vue';
        import { fetchAuthProviders } from '../api/auth';

        const entraIdAvailable = ref(false);
        const showLocalForm = ref(false);

        onMounted(async () => {
          try {
            const resp = await fetchAuthProviders();
            entraIdAvailable.value = resp.data.attributes.entraid;
            // When Entra ID is unavailable, default to showing the local form.
            showLocalForm.value = !entraIdAvailable.value;
          } catch {
            // If the providers endpoint fails (older server, network), fall back
            // to showing both options so the user can still sign in.
            entraIdAvailable.value = true;
            showLocalForm.value = false;
          }
        });
        ```

- [ ] Task 6: i18n (AC: #2, #3)
  - [ ] 6.1 Add to all five locales (`web/src/locales/{en,de,es,fr,uk}.json`):
        - `auth.moreLoginOptions` — EN: "more login options"; DE: "weitere
          Anmeldeoptionen"
        - `auth.lessLoginOptions` — EN: "less login options"; DE: "weniger
          Anmeldeoptionen"
        Spanish, French, Ukrainian translations: best-effort; flag for
        native-speaker review in completion notes if uncertain.

- [ ] Task 7: Scoped CSS for new elements (AC: #1, #5)
  - [ ] 7.1 Add scoped styles in `LoginView.vue`:
        ```css
        .login-logo {
          max-width: 220px;
          height: auto;
          margin: 0 auto 8px;
          display: block;
        }
        .login-entraid-btn {
          text-transform: none;
          font-weight: 500;
        }
        .entra-id-icon {
          width: 20px;
          height: 20px;
          display: inline-block;
        }
        .login-more-options {
          text-decoration: none;
          cursor: pointer;
        }
        .login-more-options:hover { text-decoration: underline; }
        ```

- [ ] Task 8: Tests (Vitest + Vue Test Utils)
  - [ ] 8.1 In `web/src/views/LoginView.test.ts`, add tests:
        - `fetchAuthProviders` mock returning `entraid: true` → Entra ID button
          visible, "more login options" link visible, local form hidden by
          default.
        - Click the link → local form visible, link text changes.
        - `fetchAuthProviders` mock returning `entraid: false` → Entra ID button
          NOT rendered, no toggle link, local form visible from the start.
        - `fetchAuthProviders` mock throws → fallback shows both (Entra ID button
          present, local form visible). Document the fallback explicitly so a
          future change does not regress it.
  - [ ] 8.2 Existing login submission tests (`handleLogin`) continue unchanged.
        Adjust any test that asserted "credentials form is always visible" — it
        is now gated.

- [ ] Task 9: Verification commands
  - [ ] 9.1 Backend:
        ```
    go test ./internal/auth/...
        golangci-lint run ./...
        ```
  - [ ] 9.2 Frontend (from `web/`):
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green.
  - [ ] 9.3 Manual smoke (chrome-devtools-mcp):
        - With Entra ID configured: open `/login`, screenshot, compare to
          `private/img_23.png`. Confirm the icon next to the button text, the
          "more login options" link beneath, and the absence of the credentials
          form. Click the link; the form expands, link text changes.
        - With Entra ID NOT configured: temporarily clear the relevant config
          values, restart the backend, open `/login`. The Entra ID button and
          toggle link are gone; the credentials form is shown.

## Dev Notes

### Backend contract

The endpoint MUST work without authentication (the login page calls it before
the user signs in). Mirror existing public routes such as `/api/v1/auth/login`
in `internal/auth/`. Response is a JSON:API single-resource document per
`.claude/rules/json-api.md`.

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Existing login page | LoginView.vue | `web/src/views/LoginView.vue` |
| Existing auth API client | `web/src/api/auth.ts` | reuse for `fetchAuthProviders` |
| Backend config introspection | `Config.EntraIDConfigured()` | `internal/config/config.go:146` |
| Public route registration pattern | `/api/v1/auth/login` handler | `internal/auth/` |
| Embedded asset pipeline | `web/public/*` → `assets/web/*` | `assets/embed.go` |
| Existing `<img src="/logo.svg">` precedent | `LoginView.vue:7`, `App.vue:6` | already embed-served |

### Anti-patterns to avoid

- Do NOT inline the Entra ID SVG inside the `.vue` template — keep it as a
  static asset under `web/public/` so the build pipeline embeds it once.
- Do NOT use the Wikimedia URL at runtime. The asset MUST be embedded; offline
  installs and corporate firewalls would otherwise break the button.
- Do NOT remove the existing local-login flow code. The form is only conditionally
  rendered; the underlying `handleLogin` handler stays intact.
- Do NOT change the OAuth click handler. `handleEntraIdLogin` still redirects to
  `/oauth/login` — only the visual presentation changes.
- Do NOT default `showLocalForm` to `true` on first paint when Entra ID is
  available — that would briefly flash the credentials form during the mount
  (FOUC).
- Do NOT make the providers endpoint authenticated. The login page is the
  unauthenticated context; making it auth-gated creates a chicken-and-egg
  problem.

### Coordination with Story 33.7

33.7 ships the SitHub brand logo. If 33.7 has already landed, use its
`logoSrc` directly (probably a `/sithub_logo.svg` import or a `<SitHubLogo
variant="full" />` component). If 33.6 ships first, the login page keeps the
existing `/logo.svg` placeholder and 33.7 swaps it out as part of its own
diff.

### Testing standards

- Vitest + Vue Test Utils, extending `LoginView.test.ts`. Mock
  `fetchAuthProviders` per test.
- Go: table-driven tests for the handler (configured / not configured /
  partially configured — `Config.EntraIDConfigured()` already gates on all 5
  required fields).

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/views/LoginView.vue]
- [Source: web/src/api/auth.ts]
- [Source: internal/config/config.go:60, 146]
- [Source: assets/embed.go]
- [Source: web/public/logo.svg] (precedent for embedded SVG)
- [Source: .claude/rules/json-api.md]
- [Source: .claude/rules/apidoc.md]
- [Source: .claude/rules/golang.md]
- [Source: .claude/rules/vue.md]
- [Source: private/epic-33.md]
- [Source: private/img_23.png]
- [Source: https://upload.wikimedia.org/wikipedia/commons/8/8c/Microsoft_Entra_ID_color_icon.svg]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `go test ./internal/auth/...` — providers handler tests pass (configured /
  not configured states). Full Go test suite green.
- `cd web && npx vitest run src/views/LoginView.test.ts` — 6 tests pass
  (4 new for providers gating + existing 2 updated for the toggle).
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — clean.
- `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml` —
  validation passes.
- `golangci-lint run ./...` — new `internal/auth/providers.go` clean
  (extracted `providerEntraID` / `providerLocal` constants to satisfy
  `goconst`). The 149 pre-existing repo-wide warnings are not Epic-33
  introduced.

### Completion Notes List

- Downloaded the official Microsoft Entra ID color icon SVG from the
  Wikimedia URL into `web/public/entra-id-icon.svg` (1301 bytes). Files in
  `web/public/` are copied to `web/dist/` by Vite and embedded by
  `assets/embed.go`. The login page references the icon via `/entra-id-icon.svg`
  — served from the binary, no external network call at runtime.
- Backend: added `GET /api/v1/auth/providers` (public, no auth required) in
  a new `internal/auth/providers.go`. Source of truth is the existing
  `Config.EntraIDConfigured()`; exposed via a new `Service.EntraIDConfigured()`
  method that reports `s.oauthConfig != nil` (the service-level invariant
  set by `NewService`). Returns JSON:API single-resource with attributes
  `{ entraid: bool, local: true }`. Registered alongside the other
  `/api/v1/auth/...` routes in `internal/startup/server.go`.
- Added Go tests covering both configured / not-configured states.
- OpenAPI: added `api-doc/endpoints/auth-providers.yaml` and a
  `/auth/providers` entry in `api-doc/openapi.yaml`. Redocly lint passes.
- Frontend: added `fetchAuthProviders()` to `web/src/api/auth.ts` with a
  typed `AuthProvidersAttributes` interface. The new client uses the existing
  `apiRequest` helper.
- `LoginView.vue` refactored:
  - Entra ID button is now the primary action, rendered first with the
    embedded Microsoft icon (`<template #prepend>`).
  - "more login options" / "less login options" link beneath the button
    toggles the local credentials form via `v-expand-transition`.
  - `onMounted` fetches the providers contract and gates the layout:
    Entra ID configured → only Entra ID + toggle visible; Entra ID not
    configured → local form shown by default, button + toggle hidden.
  - Fallback on fetch failure: render both (Entra ID button + local form)
    so users can still authenticate during transient backend issues.
- i18n: added `auth.moreLoginOptions` and `auth.lessLoginOptions` to all
  five locales (en, de, es, fr, uk). Translations are best-effort for es,
  fr, uk; flag for native-speaker review.
- Tests in `LoginView.test.ts`:
  - Updated mock setup to include `fetchAuthProviders` and a helper that
    returns either `{ entraid: true }` or `{ entraid: false }`.
  - Added `v-expand-transition` and `v-divider` to the stub set.
  - Updated the two existing tests so they expand the local form via the
    toggle before submitting / interacting with the Entra ID button.
  - Added a new `auth providers gating` describe block with 4 tests
    covering all four configurations: Entra ID only, toggle expand/collapse,
    no Entra ID, and the fetch-failure fallback.

### File List

Backend (added):

- `internal/auth/providers.go` — `EntraIDConfigured()` method on `*Service`,
  `ProvidersHandler` returning the JSON:API document.
- `internal/auth/providers_test.go` — tests for both configured / not
  configured cases.
- `api-doc/endpoints/auth-providers.yaml` — OpenAPI spec for the new
  endpoint.

Backend (modified):

- `internal/startup/server.go` — registered
  `GET /api/v1/auth/providers` next to the existing public auth routes.
- `api-doc/openapi.yaml` — added the `/auth/providers` path entry.

Frontend (added):

- `web/public/entra-id-icon.svg` — official Microsoft Entra ID color icon
  fetched from Wikimedia and committed for binary embedding.

Frontend (modified):

- `web/src/api/auth.ts` — new `AuthProvidersAttributes` type and
  `fetchAuthProviders()` client.
- `web/src/views/LoginView.vue` — refactored layout (Entra ID primary,
  collapsible local form, embedded icon, scoped styles).
- `web/src/views/LoginView.test.ts` — providers mock + 4 new tests +
  2 updated tests.
- `web/src/locales/{en,de,es,fr,uk}.json` — new `auth.moreLoginOptions` /
  `auth.lessLoginOptions` keys.
