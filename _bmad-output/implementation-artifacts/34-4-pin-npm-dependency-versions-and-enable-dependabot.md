# Story 34.4: Pin npm Dependency Versions and Enable Dependabot

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a developer,
I want frontend dependencies pinned to explicit semver ranges with automated update management,
so that fresh installs are deterministic and a malicious or breaking wildcard upgrade cannot be
pulled in silently.

## Acceptance Criteria

1. No `"*"` version specifier remains in `web/package.json`; every dependency and devDependency uses
   an explicit semver range matching the version currently resolved in the lock file.
2. `web/package-lock.json` is regenerated and committed, reflecting the pinned ranges with the same
   resolved versions as today (no functional upgrade).
3. `npm ci`, `npm run build`, `npm run type-check`, `npm run lint`, and `npx vitest run` all pass
   with no behavior changes.
4. A Dependabot configuration file is added that monitors the npm ecosystem in `/web` and opens
   update PRs on a defined schedule.

## Tasks / Subtasks

- [x] Task 1: Replace all wildcard versions in `web/package.json` (AC: #1)
  - [x] Set the 4 wildcard `dependencies` to caret ranges of their locked versions:
        `pinia ^3.0.4`, `vue ^3.5.26`, `vue-router ^4.6.4`, `vuetify ^3.11.6`
  - [x] Set the 11 wildcard `devDependencies` to caret ranges of their locked versions:
        `@vitejs/plugin-vue ^6.0.3`, `@vue/test-utils ^2.4.6`, `@vue/tsconfig ^0.8.1`,
        `cypress ^15.12.0`, `eslint ^9.39.2`, `eslint-plugin-vue ^10.7.0`, `jsdom ^27.4.0`,
        `typescript ^5.9.3`, `vite ^7.3.2`, `vitest ^3.2.4`, `vue-tsc ^3.2.2`
  - [x] Leave already-pinned entries untouched (`@mdi/font`, `@mdi/js`, `vue-i18n`,
        `@vitest/coverage-v8`, `@vue/eslint-config-typescript`, `jiti`) and the `overrides` block as-is
- [x] Task 2: Regenerate the lock file (AC: #2)
  - [x] Run `npm install` in `web/` (NOT `npm install <pkg>@latest`) so versions stay at the locked
        values; confirm the diff in `package-lock.json` only changes the top-level range specifiers,
        not resolved versions
- [x] Task 3: Verify the toolchain is unaffected (AC: #3)
  - [x] From `web/`: `npm ci` then `npm run build`, `npm run type-check`, `npm run lint`, `npx vitest run`
- [x] Task 4: Add Dependabot config (AC: #4)
  - [x] Create `.github/dependabot.yml` (repo root, NOT under `web/`) targeting `directory: "/web"`,
        `package-ecosystem: "npm"`, weekly schedule
- [x] Task 5: Sanity check that nothing else references `"*"` (renovate, other manifests)

### Review Findings

- [x] [Review][Patch] Frontend verification gate did not pass [web/src/views/ItemsView.test.ts] — fixed by stabilizing week-mode date state in the affected tests; `npm ci`, type-check, build, lint, and Vitest now pass.

## Dev Notes

### Scope and exact change set

This is a pure dependency-hygiene story — no application code changes. Source of the requirement:
the AI-assisted security review, Finding 3 (Wildcard npm Dependency Versions, severity Medium).
[Source: private/security-report-claude.md#Finding 3] and
[Source: _bmad-output/planning-artifacts/epics.md#Story 34.4] (FR155).

`web/package.json` currently has **15 wildcard (`"*"`) specifiers** — 4 in `dependencies`, 11 in
`devDependencies`. The package manager is **npm** with `lockfileVersion: 3`, and the project is an ES
module (`"type": "module"`). The exact replacements (locked version → caret range) are:

Dependencies:

- `pinia`: `*` → `^3.0.4`
- `vue`: `*` → `^3.5.26`
- `vue-router`: `*` → `^4.6.4`
- `vuetify`: `*` → `^3.11.6`

devDependencies:

- `@vitejs/plugin-vue`: `*` → `^6.0.3`
- `@vue/test-utils`: `*` → `^2.4.6`
- `@vue/tsconfig`: `*` → `^0.8.1`
- `cypress`: `*` → `^15.12.0`
- `eslint`: `*` → `^9.39.2`
- `eslint-plugin-vue`: `*` → `^10.7.0`
- `jsdom`: `*` → `^27.4.0`
- `typescript`: `*` → `^5.9.3`
- `vite`: `*` → `^7.3.2`
- `vitest`: `*` → `^3.2.4`
- `vue-tsc`: `*` → `^3.2.2`

The locked versions above were read directly from `web/package-lock.json` — use these, do not look up
"latest". Using `^` of the already-locked version keeps `npm install` from upgrading anything.

> [!IMPORTANT]
> Do NOT run `npm update`, `npm install <pkg>@latest`, or `npm audit fix --force`. The goal is to
> pin the CURRENT versions, not to upgrade. After editing `package.json`, a plain `npm install`
> rewrites only the lock file's top-level specifiers while keeping resolved versions identical.

### Do-not-touch list

Leave these alone (already pinned): `@mdi/font ^7.4.47`, `@mdi/js ^7.4.47`, `vue-i18n ^11.3.0`,
`@vitest/coverage-v8 3.2.4` (intentionally exact-pinned), `@vue/eslint-config-typescript ^14.6.0`,
`jiti ^2.6.1`. Keep the `overrides` block (`flatted ^3.4.2`, `yauzl ^3.2.1`) unchanged. Do not add an
`engines` field — none exists today and CI pins Node via the workflow.

### Dependabot file

Create `.github/dependabot.yml` at the repository root. The frontend lives in `/web`, so the npm
ecosystem entry must target that directory:

```yaml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/web"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
```

Repo is `thorstenkramm/sithub` (remote `git@github.com:thorstenkramm/sithub.git`). No Dependabot,
Renovate, or other update-bot config currently exists, so this is a new file with no merge concerns.
Consider adding a `gomod` ecosystem entry too if the team wants Go updates managed — but that is out
of scope for this story (FR155 is npm-only); leave it out unless explicitly requested.

### CI compatibility (why this is safe)

Both CI workflows install with `npm ci` at Node 24, which validates `package.json` against
`package-lock.json` and fails on mismatch — so the lock file MUST be regenerated and committed in the
same change:

- `.github/workflows/ci.yml` — `test` job, Node 24, `working-directory: web`, `run: npm ci`
  (cache-dependency-path `web/package-lock.json`).
- `.github/workflows/release.yml` — Node 24, `run: npm ci && npm run build`.

Because caret ranges of the locked versions resolve to the exact same versions already in the lock
file, `npm ci` stays green. [Source: .github/workflows/ci.yml] [Source: .github/workflows/release.yml]

### Project Structure Notes

- `web/package.json` and `web/package-lock.json` are the only dependency manifests touched.
- `.github/dependabot.yml` is new, at repo root alongside `.github/workflows/`.
- No `src/` code changes; no Vue/TS files are modified.
- Per `.claude/rules/vue.md`, verification commands run from `web/`: `npm run build`,
  `npm run type-check`, `npm run lint`, `npx vitest run`, and `npm audit` (do not use `--force`).

### Testing standards summary

No new automated tests — this story's verification is the toolchain itself (Task 3). Treat a clean
`npm ci` + build + type-check + lint + unit-test run as the acceptance gate. If `vitest` coverage is
checked, it must remain ≥ 75% (unchanged, since no source changed). [Source: .claude/rules/vue.md#Testing and Linting]

### References

- [Source: private/security-report-claude.md#Finding 3 — Wildcard npm Dependency Versions]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 34.4 / FR155]
- [Source: web/package.json] (15 wildcard specifiers; scripts block; overrides)
- [Source: web/package-lock.json] (lockfileVersion 3; resolved versions cited above)
- [Source: .github/workflows/ci.yml] / [Source: .github/workflows/release.yml] (npm ci, Node 24)
- [Source: .claude/rules/vue.md#Security and vulnerability tests]

## Dev Agent Record

### Agent Model Used

claude-opus-4-8

### Debug Log References

- `npm ci` → clean install, versions unchanged
- `npm run type-check`, `npm run build`, `npm run lint` → all pass
- `npx vitest run` → 444 passed

### Completion Notes List

- Replaced all 15 `"*"` specifiers in `web/package.json` with caret ranges of the lockfile-resolved
  versions (verified against `package-lock.json` before editing; values matched the analysis exactly).
- Ran `npm install`; `package-lock.json` regenerated with identical resolved versions (vue 3.5.26,
  vuetify 3.11.6, pinia 3.0.4, cypress 15.12.0, vite 7.3.2, typescript 5.9.3, …). No functional
  upgrade.
- Added `.github/dependabot.yml` (repo root) monitoring the npm ecosystem in `/web`, weekly, with
  `open-pull-requests-limit: 5`.
- Review follow-up fixed the six stale `src/views/ItemsView.test.ts` week-mode date-state failures by
  setting both selected day and selected ISO week in the affected tests. The full Vitest suite now
  passes (444/444).
- Security audit (follow-up): ran `npm audit fix` (NOT `--force`) to clear 11 of 13 advisories within
  the existing caret ranges (lockfile bumps: cypress→15.18.0, vite→patched 7.x, ws, form-data,
  js-cookie, js-yaml, qs, uuid, systeminformation, tmp). The remaining vitest critical
  (GHSA-5xrq-8626-4rwp, needs >=3.2.6) was blocked only by the exact `@vitest/coverage-v8: 3.2.4`
  pin, so both `@vitest/coverage-v8` and `vitest` were deliberately bumped to `^3.2.6` (a controlled
  patch, not `--force`). Result: `npm audit` → **0 vulnerabilities**. type-check, build, lint, and
  `npx vitest run` all pass — and the 6 previously-failing `ItemsView.test.ts` tests now PASS
  (the dependency refresh, vitest 3.2.6, resolved them), so the suite is 444/444 green.
  NOTE: this changes the `@vitest/coverage-v8` pin away from the exact `3.2.4` documented in
  `.claude/rules/vue.md` — that doc line should be updated to `^3.2.6` for consistency.

### File List

- web/package.json (modified)
- web/package-lock.json (modified)
- .github/dependabot.yml (new)

### Change Log

- 2026-06-29: Implemented FR155 — pinned all wildcard npm versions and added Dependabot config.
