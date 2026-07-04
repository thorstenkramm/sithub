# Backlog — Candidate Epics

Un-scheduled epic candidates. These are not yet decomposed into stories and carry no FR numbers.
Promote a candidate to a real epic (in `epics.md`) via `create-epics-and-stories` when it is
scheduled.

Origin: Epic 34 retrospective and the Dependabot triage of 2026-07-04. Each candidate is a major
dependency migration deliberately deferred; a matching `ignore` rule for its `semver-major` updates
exists in `.github/dependabot.yml` and must be removed when the migration is done so future majors
flow again.

## Candidate: Vuetify 3 -> 4 Migration

Major, breaking upgrade of the core UI framework (Dependabot PR #4, closed; currently on 3.11.6).
Must be a deliberate migration, not an automated bump.

- Scope: work through the official Vuetify 4 upgrade guide and breaking-changes list (Vuetify MCP
  `get_upgrade_guide` / `get_v4_breaking_changes` are available); update component usage, prop/slot
  changes, and theme/design-token APIs across all views; verify tree-shaking/bundle config and MDI
  icon integration; re-run unit + Cypress E2E and check for visual regressions.
- Highest-risk areas: floor plan editor, weekly desk matrix, item tiles, dialogs.
- Risks: broad UI regression surface; subtle theming/token changes; interaction with the recently
  merged Vite 8, Vue 3.5, and plugin-vue updates.
- Done when: all views render correctly, unit + E2E green, no visual regressions, and the `vuetify`
  major-ignore rule is removed from `.github/dependabot.yml`.

## Candidate: vitest 3 -> 4 Migration

`@vitest/coverage-v8@4.x` requires `vitest@4.x` as a peer (Dependabot PR #5, closed on ERESOLVE);
the two must move together. Vite 8 is already on `main`, so this migration should also confirm
vitest 4 works with Vite 8.

- Scope: bump `vitest` and `@vitest/coverage-v8` to 4.x together; address test-runner config/API
  changes; keep coverage at or above the enforced threshold (80% Go / 75% frontend); update the
  documented pin in `.claude/rules/vue.md` (currently `@vitest/coverage-v8@^3.2.6`); confirm the full
  unit suite passes on vitest 4.
- Risks: test-runner config and API changes; coverage-reporter differences; jsdom interplay.
- Done when: all unit tests pass on vitest 4, coverage maintained, docs updated, and the `vitest` +
  `@vitest/coverage-v8` major-ignore rules are removed from `.github/dependabot.yml`.

## Candidate: TypeScript 5.9 -> 6 Migration

Dependabot's TypeScript 6 bump (PR #10, closed) failed `type-check`; the toolchain (`vue-tsc`,
`@vue/eslint-config-typescript`) does not yet support TS 6. Partly gated on upstream toolchain
support.

- Scope: wait for / confirm `vue-tsc` and `@vue/eslint-config-typescript` releases that support
  TypeScript 6; bump `typescript`; fix new or stricter type errors surfaced across the codebase;
  verify `type-check`, `lint`, and `build` are green.
- Risks: new/stricter type errors project-wide; timing depends on external toolchain readiness.
- Done when: `type-check` + `lint` + `build` are green on TypeScript 6, and the `typescript`
  major-ignore rule is removed from `.github/dependabot.yml`.
