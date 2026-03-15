# Story 18.1: Rename Config Section from [spaces] to [areas]

Status: done

## Story

As an operator,
I want the configuration to use consistent `[areas]` terminology everywhere,
So that there is no confusion between "spaces" and "areas" in configuration, code,
and documentation.

## Acceptance Criteria

1. **Given** I have a `sithub.toml` with `[areas]` section
   **When** the server starts
   **Then** it reads configuration from the `[areas]` table

2. **Given** I use `--areas-config-file` flag or `SITHUB_AREAS_CONFIG_FILE` env var
   **When** the server starts
   **Then** it applies the override correctly

3. **Given** the codebase
   **When** searching for the term "space" or "spaces"
   **Then** no references to the old terminology exist (data models, CLI flags, env vars)

## Tasks

1. Rename Go package `internal/spaces` to `internal/areas`
2. Rename `SpacesConfig` struct to `AreasConfig`, field `Spaces` to `Areas`
3. Rename sentinel error `ErrMissingSpacesConfig` to `ErrMissingAreasConfig`
4. Update viper defaults from `spaces.*` to `areas.*`
5. Rename CLI flag `--spaces-config-file` to `--areas-config-file`
6. Update env var from `SITHUB_SPACES_CONFIG_FILE` to `SITHUB_AREAS_CONFIG_FILE`
7. Update TOML section `[spaces]` to `[areas]` in sithub.toml and sithub.example.toml
8. Update all variable names (`spacesConfig` → `areasConfig`, etc.)
9. Update all import paths and references across handler packages
10. Update CI workflow config references
11. Run tests, linting, and duplication checks

## Dev Agent Record

### Agent Model Used

GPT-5 Codex

### Completion Notes List

- Renamed the backend package, config struct field, sentinel error, flags, and env var usage from `spaces` to `areas`
- Updated config defaults, startup wiring, imports, and CI references to use `areas`
- Reconciled the story record so the implementation can be audited like the rest of Epic 18
- Verified with `go test ./...`

### File List

- `internal/areas/config.go` — Renamed package from `spaces` to `areas`
- `internal/config/config.go` — Renamed config field, defaults, and sentinel errors to `areas`
- `internal/startup/server.go` — Loads areas config via the renamed package and field
- `cmd/sithub/main.go` — Uses `--areas-config-file` override wiring
- `cmd/sithub/options_test.go` — Verifies CLI override mapping
- `.github/workflows/ci.yml` — Updated test config references to `[areas]`
- `sithub.example.toml` — Documents `[areas]` and `config_file`

## Change Log

- 2026-03-13: Story implemented and documentation synced after code review fixes.
