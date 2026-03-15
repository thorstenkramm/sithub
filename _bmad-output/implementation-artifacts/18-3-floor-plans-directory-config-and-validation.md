# Story 18.3: Floor Plans Directory Configuration & Validation

Status: done

## Story

As an operator,
I want to configure a floor plans directory and have all image references validated
at startup,
So that runtime errors from missing or invalid images are caught early.

## Acceptance Criteria

1. **Given** `areas.floor_plans` is set to a directory name
   **When** the server starts
   **Then** the directory is resolved inside `main.data_dir` and its existence is validated

2. **Given** `areas.floor_plans` points to a non-existent directory
   **When** the server starts
   **Then** the server exits with a descriptive error

3. **Given** the areas config references floor plan images
   **When** the server starts
   **Then** each referenced image is checked for existence in the floor plans directory

4. **Given** a referenced image has an unsupported format (not jpg, png, svg)
   **When** the server starts
   **Then** the server exits with a descriptive error listing the invalid file

5. **Given** `areas.floor_plans` is not set
   **When** the server starts
   **Then** floor plan features are disabled and no validation occurs

## Tasks / Subtasks

- [x] Add `FloorPlansDir` field to `AreasConfig` struct in `internal/config/config.go`
- [x] Add `resolveFloorPlansDir()` function to resolve and validate directory path
  - [x] Resolve relative paths inside `data_dir`
  - [x] Reject absolute paths outside `data_dir`
  - [x] Verify directory existence
- [x] Add sentinel errors: `ErrFloorPlansDirOutsideDataDir`, `ErrFloorPlansDirNotFound`
- [x] Add `ValidateFloorPlans()` function to `internal/areas/config.go`
  - [x] Define `supportedFloorPlanExts` map (`.jpg`, `.jpeg`, `.png`, `.svg`)
  - [x] Validate each `floor_plan` reference in areas and item groups
  - [x] Check file existence in floor plans directory
  - [x] Check file extension against supported formats
- [x] Call `ValidateFloorPlans()` from `internal/startup/server.go` after loading areas config
- [x] Add `--areas-floor-plans` CLI flag and `SITHUB_AREAS_FLOOR_PLANS` env var
- [x] Keep legacy `floor_plans_dir` config/env/flag support as compatibility fallback
- [x] Update `sithub.example.toml` and `sithub.toml` with `floor_plans` setting
- [x] Add unit tests
  - [x] Config tests: relative dir, non-existent dir, dir outside data_dir, empty (skipped)
  - [x] Areas tests: valid images, unsupported format, missing file, no references
- [x] Run tests, linting, and duplication checks

## Dev Notes

### Architecture: Split Validation

Directory resolution and existence checking happen in `config.Load()` (config package).
Image reference validation happens in `startup/server.go` via `areas.ValidateFloorPlans()`.
This split avoids import cycles between `config` and `areas` packages. The implementation now
matches the epic-facing contract (`areas.floor_plans`, `--areas-floor-plans`,
`SITHUB_AREAS_FLOOR_PLANS`) while still accepting the earlier `floor_plans_dir` variant as a
compatibility fallback.

### Floor Plan References

`floor_plan` values are validated as filenames only. Nested paths such as
`floor_plans/office.svg` are rejected so the config, API, and frontend all agree on a single
identifier for each image.

### Supported Formats

`.jpg`, `.jpeg`, `.png`, `.svg` — matching the formats supported by the floor plan image
serving endpoint.

### References

- Epic 18 Story 18.3: `_bmad-output/planning-artifacts/epics.md`
- FR61, FR62: `_bmad-output/planning-artifacts/prd.md`
- `internal/config/config.go`: `resolveFloorPlansDir()`
- `internal/areas/config.go`: `ValidateFloorPlans()`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Added `FloorPlansDir` to `AreasConfig` struct with mapstructure tag
- Created `resolveFloorPlansDir()` with same path resolution pattern as areas config
- Created `ValidateFloorPlans()` and `validateFloorPlanFile()` in areas package
- Added CLI flag and environment variable override for `areas.floor_plans`
- Added compatibility fallback for the earlier `floor_plans_dir` naming
- Updated example TOML with documentation
- All Go tests, linting, and duplication checks pass

### File List

- `internal/config/config.go` — Added `FloorPlansDir` field, `resolveFloorPlansDir()`, sentinel errors
- `internal/config/config_test.go` — Added floor plans directory tests
- `internal/areas/config.go` — Added `ValidateFloorPlans()`, `validateFloorPlanFile()`, supported extensions
- `internal/areas/config_test.go` — Added floor plan validation tests
- `internal/startup/server.go` — Added floor plan validation call after areas config load
- `cmd/sithub/main.go` — Added `--areas-floor-plans` flag and legacy fallback
- `cmd/sithub/options_test.go` — Verifies floor plan CLI override mapping
- `sithub.example.toml` — Added `floor_plans` setting documentation
- `sithub_areas.example.yaml` — Documents filename-only `floor_plan` references

## Change Log

- 2026-03-13: Story implemented and verified.
- 2026-03-13: Code review fixes aligned the config contract with Epic 18 and tightened filename validation.
