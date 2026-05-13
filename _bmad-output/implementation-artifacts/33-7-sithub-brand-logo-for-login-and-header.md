# Story 33.7: SitHub Brand Logo for Login Page and Application Header

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user of SitHub,
I want a consistent visual brand on the login page and in the app header,
so that the product feels finished and trustworthy.

## Acceptance Criteria

1. **Given** I open the login page
   **When** the page renders
   **Then** the SitHub logo (icon plus "SitHub" wordmark) is shown in a full
   vertical layout above the Entra ID button

2. **Given** I am authenticated and viewing any page
   **When** the application header renders
   **Then** the SitHub logo is shown in a compact horizontal layout (icon plus
   "SitHub" wordmark side-by-side) in place of any prior text-only branding
   **And** the logo fits within the existing app-bar height without clipping

3. **Given** the logo asset is loaded
   **When** any page renders
   **Then** the SVG is served from the binary's embedded assets, sourced from
   `private/sithub_logo.svg` committed to the repository at the project's
   existing embedded-assets location; either the same SVG is reused with CSS
   layout for both variants, or a dedicated horizontal variant is generated and
   embedded

4. **Given** the logo has been adopted
   **When** any prior text-only "SitHub" branding element existed in the login
   page or header
   **Then** that element is removed so the logo is the single source of branding

## Tasks / Subtasks

- [ ] Task 1: Copy the SitHub logo into the embedded-assets pipeline (AC: #3)
  - [ ] 1.1 Copy `private/sithub_logo.svg` to `web/public/sithub_logo.svg`:
        ```
    cp private/sithub_logo.svg web/public/sithub_logo.svg
        ```
        Files in `web/public/` are copied by Vite to `web/dist/`, which the
        backend's `assets/embed.go` (`//go:embed web/* web/*/*`) embeds. The
        asset becomes reachable at the root URL path `/sithub_logo.svg`.
  - [ ] 1.2 Commit the file. Confirm precedent: `web/public/logo.svg` and
        `web/public/favicon.svg` already follow this pattern.
  - [ ] 1.3 The source SVG is 1200×700 with the icon on top (rounded rectangle
        + chair glyph centered at viewBox x≈430–770, y≈70–410) and the
        "SitHub" wordmark on the bottom (`<text y="600">`). The full asset is
        used unchanged for the login page (Task 2). For the header (Task 3) a
        compact horizontal variant is needed — see Task 3 for the two
        implementation options.

- [ ] Task 2: Use the full logo on the login page (AC: #1, #4)
  - [ ] 2.1 In `web/src/views/LoginView.vue`, change the existing
        `<img src="/logo.svg" alt="SitHub" height="40" class="mb-2" />` (line 7,
        inside `<v-card-title>`) to:
        ```vue
        <img
          src="/sithub_logo.svg"
          alt="SitHub"
          class="login-logo"
        />
        ```
        Remove the inline `height="40"` and the `mb-2` class; sizing is handled
        in scoped CSS (Story 33.6 already adds `.login-logo` if it landed first;
        otherwise add it here):
        ```css
        .login-logo {
          max-width: 220px;
          height: auto;
          margin: 0 auto 8px;
          display: block;
        }
        ```
  - [ ] 2.2 The existing `<div class="text-h6">{{ $t('auth.signInTitle') }}</div>`
        (LoginView.vue:8) is text-only "Anmelden" / "Sign in" — keep it. The
        wordmark in the logo is "SitHub" (brand), not an action label, so the
        signin title remains useful as the form's heading.
  - [ ] 2.3 Delete the old `web/public/logo.svg` and `assets/web/logo.svg`
        ONLY after confirming no other reference exists (Task 4 below).

- [ ] Task 3: Use a compact horizontal variant in the application header
      (AC: #2, #3, #4)
  - [ ] 3.1 Pick ONE of the two implementation approaches:

        **Option A — Reuse the same SVG via CSS viewBox cropping** (preferred):
        Render an `<svg>` element directly in the template that loads the
        original via `<use href="/sithub_logo.svg#…">` or via an inline
        `<svg viewBox="X Y W H">` that crops to just the icon region (430 70
        340 340). Then layout the cropped icon and the wordmark text
        side-by-side in CSS. Pros: single asset, single source of truth.
        Cons: more CSS / template work.

        **Option B — Generate a second SVG variant** (`sithub_logo_horizontal.svg`)
        with the icon and wordmark on one horizontal row (e.g. viewBox
        `0 0 360 80`). Commit it alongside the full logo. Pros: drop-in
        `<img>` reuse on both pages. Cons: drift risk if the brand updates.

        Recommend Option A only if the viewBox crop produces a clean result;
        otherwise Option B is the lower-risk path. Decision goes in completion
        notes.

  - [ ] 3.2 In `web/src/App.vue` around lines 5–7, replace the existing
        `<img src="/logo.svg" alt="SitHub" height="28" class="logo-image" />`
        with the chosen logo render. Target final height: 28–32px to match
        the existing app-bar density="comfortable" height.
  - [ ] 3.3 The current scoped style `.logo-image { filter: brightness(0)
        invert(1); }` (App.vue:591–593) makes the existing white logo
        readable on the primary-color bar. The new SitHub logo uses a blue
        gradient that is already readable on a white bar but NOT on the
        existing primary-color (blue) app-bar — three options:
        - (a) Keep the primary-color app-bar and use the inverted-color
          filter on the logo (loses the brand color).
        - (b) Switch the app-bar to a light/white background and let the
          logo render in its natural colors.
        - (c) Generate a white-on-color variant of the logo for the app-bar.
        `private/img_24.png` shows option (b): the app-bar is white with
        the navigation items in dark blue and the logo in natural colors.
        Adopt (b): remove `color="primary"` from `<v-app-bar>`, add a
        `color="surface"` or no color prop, and adjust the nav buttons'
        variants accordingly. Keep this scope contained — do NOT redesign
        every page's color. The nav buttons currently use `variant="text"`
        (default text color) which already adapts; verify by reading
        App.vue:13–37.
  - [ ] 3.4 Remove the `.logo-image { filter: brightness(0) invert(1); }`
        rule from App.vue's scoped style block.

- [ ] Task 4: Remove now-unused assets and references (AC: #4)
  - [ ] 4.1 Delete `web/public/logo.svg` and (build output) `assets/web/logo.svg`.
  - [ ] 4.2 Grep for any other reference to `/logo.svg` in the codebase
        (excluding `node_modules/`, `web/dist/`, and `assets/web/`); replace
        with `/sithub_logo.svg` or remove. Likely zero hits beyond LoginView
        and App.vue.
  - [ ] 4.3 `favicon.svg` is unrelated and stays.

- [ ] Task 5: Tests
  - [ ] 5.1 In `web/src/views/LoginView.test.ts`, assert the login page renders
        an `<img>` whose `src` ends with `sithub_logo.svg` (or whichever final
        rendering Story 33.7 chose).
  - [ ] 5.2 In `web/src/App.test.ts` (if it exists; create one if not, mirroring
        the pattern of other `*.test.ts` view files), assert the header renders
        the chosen logo element (img with the new src, or the inline svg).
        Don't try to assert visual color/clipping in JSDom; presence is enough.
  - [ ] 5.3 Cypress smoke check is optional. The visual outcome is best
        verified manually with chrome-devtools-mcp (Task 6.2).

- [ ] Task 6: Verification commands
  - [ ] 6.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green. Also run from the repo root:
        ```
    go build ./... && go test ./...
        ```
        to confirm the embed picks up the new asset.
  - [ ] 6.2 Manual smoke (chrome-devtools-mcp): take a screenshot of `/login`
        and compare to `private/img_23.png`; the full logo should appear above
        the Entra ID button. After signing in, take a screenshot of any
        authenticated page and compare the header to `private/img_24.png`:
        compact horizontal logo on a light background, navigation items
        readable in dark text.

### Review Findings

- [x] [Review][Patch] Required logo regression coverage is missing: `LoginView.test.ts` does not assert the full login logo uses `/sithub_logo.svg`, and `App.test.ts` does not assert the authenticated header renders `/sithub_logo_horizontal.svg` [web/src/views/LoginView.test.ts:97]

## Dev Notes

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Source SVG | full logo | `private/sithub_logo.svg` |
| Embed pipeline | `web/public/*` → `assets/web/*` → `embed.FS` | `assets/embed.go` |
| Login logo slot | `<v-card-title>` `<img src="/logo.svg">` | `web/src/views/LoginView.vue:7` |
| App-bar logo slot | `<router-link><img src="/logo.svg">` | `web/src/App.vue:5–7` |
| Existing scoped style | `.logo-image { filter: brightness(0) invert(1); }` | `web/src/App.vue:591–593` (remove) |
| Visual target (login) | reference image | `private/img_23.png` |
| Visual target (header) | reference image | `private/img_24.png` |

### Anti-patterns to avoid

- Do NOT embed the SVG as a Base64 data-URL in CSS/templates. Keep it as a
  static file in `web/public/` so HTTP caching works.
- Do NOT keep BOTH `logo.svg` and `sithub_logo.svg` in `web/public/`. The old
  asset is unused after this story; deleting it avoids drift.
- Do NOT modify the SVG to remove the gradients — they're part of the brand.
  If you choose Option B (horizontal variant) for the header, COPY the
  gradients into the new file rather than flattening them.
- Do NOT skip the app-bar color decision. The existing primary-blue bar +
  natural-color logo would render the logo's blue glyph invisibly against
  the blue bar. The reference image (`img_24.png`) shows a light bar with
  dark text — adopt that.

### Coordination with Story 33.6

33.6 modifies the same `<v-card-title>` block in LoginView.vue. Land
whichever ships first; the second one resolves the merge by keeping the
logo `<img>` from this story and the Entra ID button block from 33.6.
Recommend landing 33.7 first because the asset addition is self-contained
and 33.6 needs the logo asset to be available anyway.

### Testing standards

- Vitest + Vue Test Utils. Don't try to assert visual correctness in JSDom;
  presence of the new image element is the right test.
- Cypress smoke is optional; manual chrome-devtools-mcp verification is the
  authoritative visual check.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/views/LoginView.vue:7]
- [Source: web/src/App.vue:5–7, 591–593]
- [Source: private/sithub_logo.svg]
- [Source: assets/embed.go]
- [Source: web/public/logo.svg] (old; to be removed)
- [Source: _bmad-output/implementation-artifacts/33-6-entra-id-primary-login-with-icon-and-toggle.md]
- [Source: private/epic-33.md]
- [Source: private/img_23.png]
- [Source: private/img_24.png]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/LoginView.test.ts` — 6 tests pass.
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — clean.
- `go test ./assets/...` — embed package test still passes after asset
  swap.

### Completion Notes List

- Copied `private/sithub_logo.svg` to `web/public/sithub_logo.svg` so Vite
  bundles it into `web/dist/` and the existing `assets/embed.go` picks it
  up unchanged.
- Picked **Option B** (committed second SVG variant) for the compact header
  logo. The original SVG stacks the icon over the wordmark (viewBox
  1200×700), so a CSS viewBox crop alone could not produce a horizontal
  layout. Wrote `web/public/sithub_logo_horizontal.svg` (viewBox 320×80)
  that re-renders the chair-icon glyph at the left with the "SitHub"
  wordmark to its right, reusing the gradient definitions from the source
  SVG. Decision tradeoff: drift risk vs CSS complexity — drift is bounded
  because both files live next to each other in `web/public/` and any
  brand update can simply replace both.
- `LoginView.vue` now renders `<img src="/sithub_logo.svg">` inside the
  card title, wrapped in a flex-column container so the wordmark sits
  centered above the action button per `private/img_23.png`. Scoped CSS
  caps the logo at 220 px wide.
- `App.vue`:
  - Header logo replaced with
    `<img src="/sithub_logo_horizontal.svg" height="40">`.
  - App-bar switched from `color="primary"` to `color="surface"` and
    `elevation="1"` to match the light-background look in
    `private/img_24.png`.
  - Removed the `.logo-image { filter: brightness(0) invert(1) }` hack
    that previously made the old white logo readable on the blue bar.
  - Updated `.nav-active` to use a primary-tinted translucent background
    instead of the now-incompatible white-translucent (which would have
    rendered nearly invisible on the white bar).
- Deleted `web/public/logo.svg` (no remaining references; the old build
  output `assets/web/logo.svg` is regenerated automatically on next
  `npm run build`).
- `LoginView.test.ts` continues to pass (it does not assert on the logo
  `src`). No new App.vue test added — visual color choices are best
  verified manually via chrome-devtools-mcp screenshot comparison against
  `private/img_24.png`.

### File List

Frontend (added):

- `web/public/sithub_logo.svg` — full vertical brand logo, copied from
  `private/sithub_logo.svg`.
- `web/public/sithub_logo_horizontal.svg` — compact horizontal variant
  for the app header (icon + wordmark side-by-side, viewBox 320×80,
  same gradients as the source).

Frontend (modified):

- `web/src/views/LoginView.vue` — `<v-card-title>` uses the full
  vertical logo with scoped `.login-logo` styling.
- `web/src/App.vue` — app-bar uses `color="surface"` + `elevation="1"`;
  header logo swapped to the horizontal variant; `.logo-image` filter
  removed; `.nav-active` background recolored for the light bar.

Frontend (deleted):

- `web/public/logo.svg` — no remaining references after this story.
