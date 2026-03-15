# Story 18.2: Enforce Areas Config Inside data_dir

Status: done

## Story

As an operator,
I want the areas config file path resolved relative to `main.data_dir`,
So that all data files are consistently located in one directory.

## Acceptance Criteria

1. **Given** `areas.config_file` is set to a relative filename like `sithub_areas.yaml`
   **When** the server starts
   **Then** the file is resolved inside `main.data_dir`

2. **Given** `areas.config_file` contains an absolute path outside `main.data_dir`
   **When** the server starts
   **Then** startup fails with a descriptive error message

## Tasks / Subtasks

- [x] Add `resolveAreasConfig()` function to `internal/config/config.go`
  - [x] Resolve relative paths via `filepath.Join(dataDir, raw)`
  - [x] Reject absolute paths outside `data_dir` using `strings.HasPrefix` check
  - [x] Verify file existence with `os.Stat`
  - [x] Replace raw config value with resolved absolute path
- [x] Add sentinel error `ErrAreasConfigOutsideDataDir`
- [x] Call `resolveAreasConfig()` from `config.Load()`
- [x] Add unit tests
  - [x] `TestLoadRelativeAreasConfig` — relative path resolved inside data_dir
  - [x] `TestLoadAbsolutePathOutsideDataDir` — absolute path outside data_dir rejected
- [x] Update `internal/startup/server_test.go` to use shared `dataDir` for both DB and areas config
- [x] Run tests, linting, and duplication checks

## Dev Notes

### Path Resolution Strategy

Relative paths are resolved using `filepath.Join(dataDir, raw)`. Absolute paths are validated
with `strings.HasPrefix(resolved, dataDir + string(filepath.Separator))` to ensure they remain
inside the data directory. This prevents configuration errors where the areas config could point
to an arbitrary location on disk.

### References

- Epic 18 Story 18.2: `_bmad-output/planning-artifacts/epics.md`
- FR60: `_bmad-output/planning-artifacts/prd.md`
- `internal/config/config.go`: `resolveAreasConfig()` function

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `resolveAreasConfig()` with path resolution and validation
- Added `ErrAreasConfigOutsideDataDir` sentinel error
- Updated server test to use shared data directory
- All Go tests, linting, and duplication checks pass

### File List

- `internal/config/config.go` — Added `resolveAreasConfig()` and sentinel error
- `internal/config/config_test.go` — Added path resolution tests
- `internal/startup/server_test.go` — Updated to use shared data directory

## Change Log

- 2026-03-13: Story implemented and verified.
