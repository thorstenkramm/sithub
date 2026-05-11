# Story 31.3: Areas Config Location Hint in Example TOML

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an operator setting up SitHub,
I want the example configuration to state that the areas config file must live inside
`data_dir`,
so that I do not waste time placing the file in an unsupported location.

## Acceptance Criteria

1. **Given** I read `sithub.example.toml`
   **When** I look at the setting that points to the areas configuration file
   **Then** its inline comment includes the sentence "Must be inside data_dir."
   **And** the wording matches the existing hint used for the floor plans directory
   setting

2. **Given** the example file is updated
   **When** the TOML linter or formatter runs
   **Then** the file remains valid TOML and follows the existing comment style
   described in `.claude/rules/toml.md`

## Tasks / Subtasks

- [ ] Task 1: Update `sithub.example.toml` — `[areas]` `config_file` comment (AC: #1)
  - [ ] 1.1 Edit `sithub.example.toml` line ~48 (the comment block above
        `#config_file = "./sithub_areas.yaml"`).
  - [ ] 1.2 Add the sentence `Must be inside data_dir.` to the description line so
        the comment reads (matching the floor-plans wording at line ~55):
        `## Path to the YAML file that defines areas, rooms, and desks. Must be inside data_dir.`
        Keep the rest of the comment block (`mandatory`, override flag, env var,
        Example, Default) unchanged and in the same order.
  - [ ] 1.3 Do not change the `config_file` value or its commented-out state. The
        change is documentation-only.
  - [ ] 1.4 Run `git diff sithub.example.toml` and confirm the diff is exactly the
        added sentence and nothing else.

- [ ] Task 2: Verify (AC: #2)
  - [ ] 2.1 Manual TOML parse check:
        `go run ./cmd/sithub run --config ./sithub.example.toml --help`
        does not print a TOML parse error. (The file is fully commented out, so
        config validation should not trip — only TOML syntax matters.)
  - [ ] 2.2 If a TOML linter is installed locally (e.g. `taplo`), run it; otherwise
        rely on the Go config loader as the de facto check. There is no project-level
        TOML lint script today.
  - [ ] 2.3 Confirm the wording matches `.claude/rules/toml.md` style: ASCII-only,
        consistent phrasing with the `floor_plans` line.

- [ ] Task 3: No-op verification of runtime behaviour
  - [ ] 3.1 Story 18.2 (`enforce-areas-config-inside-data-dir`) already enforces
        the constraint at startup. This story only documents it. Do **not** add or
        change any Go validation code — that is out of scope.
  - [ ] 3.2 Spot-check `internal/startup/server.go::loadAndValidateAreas` and
        `internal/areas/...` to confirm the runtime check still exists; if it has
        been removed for any reason, raise a question in completion notes rather
        than fixing it under this story.

## Dev Notes

### What changes and what does not

- Pure documentation change in one file: `sithub.example.toml`.
- No backend code change. No frontend change. No tests to add — story 18.2 already
  has runtime tests for the constraint.
- The change is intentionally tiny because the wording is the only correct fix:
  users read `sithub.example.toml` to learn the configuration surface, and the
  current comment for `config_file` does not mention the data_dir rule even though
  the rule is enforced at startup.

### Reference text

Look at the existing comment block for `floor_plans` (line ~55 in
`sithub.example.toml`):

```toml
  ## Floor plans directory, string, optional
  ## Can be overridden with --areas-floor-plans flag or SITHUB_AREAS_FLOOR_PLANS environment variable
  ## Directory containing floor plan images (jpg, png, svg). Must be inside data_dir.
  ## If not set, floor plan features are disabled.
  ## Example: "floor_plans"
  ## Default: none
  #floor_plans = "floor_plans"
```

The pattern is: the `Must be inside data_dir.` sentence is appended to the
description line, not added on its own line. Mirror that.

### Anti-patterns to avoid

- Do NOT introduce a new comment line such as
  `## Note: Must be inside data_dir.`. The convention is to append the constraint
  to the description sentence (see floor_plans).
- Do NOT touch any other section of `sithub.example.toml`. Scope is one comment.
- Do NOT change the runtime validation in `internal/startup/server.go` or the
  `areas` package. That is already covered by story 18.2.
- Do NOT add a test asserting the comment exists. Documentation comments are not
  worth testing.

### Testing Standards

- No automated test required.
- Manual verification: `git diff sithub.example.toml` shows the single-sentence
  addition; `go run ./cmd/sithub run --config ./sithub.example.toml --help` runs
  without TOML parse errors.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 31 Stories: Live Updates, Favorites Rework & Areas Config Hint]
- [Source: sithub.example.toml]
- [Source: .claude/rules/toml.md]
- [Source: _bmad-output/implementation-artifacts/18-2-enforce-areas-config-inside-data-dir.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `git diff sithub.example.toml` shows the single-sentence addition and
  nothing else.
- `go run ./cmd/sithub run --config ./sithub.example.toml --help` prints the
  usage banner with no TOML parse error.

### Completion Notes List

- Appended `Must be inside data_dir.` to the description line of
  `[areas].config_file` in `sithub.example.toml`, mirroring the wording
  already used for `[areas].floor_plans`. No new comment line introduced —
  the constraint is appended to the existing description per the style
  shown in `.claude/rules/toml.md` and the floor-plans precedent.
- Did not touch the runtime validation. Story 18.2 already enforces "areas
  config must live inside `data_dir`" at startup (spot-checked:
  `internal/startup/server.go::loadAndValidateAreas` still wires through
  `internal/areas.Load` which calls `EnsureInsideDataDir`).
- No tests added — the change is a documentation comment only.

### File List

Documentation (modified):

- `sithub.example.toml` (one line: appended `Must be inside data_dir.` to
  `[areas].config_file` description)
