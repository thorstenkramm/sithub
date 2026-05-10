# Story 30.1: Reject Duplicate Items in Areas Config at Startup

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an operator,
I want the server to refuse to start when the areas configuration contains duplicate item
identifiers,
so that misconfigurations are caught immediately instead of producing inconsistent booking
state at runtime.

## Acceptance Criteria

1. **Given** the areas configuration contains the same item identifier (e.g. `desk29`) more
   than once across any subareas or rooms
   **When** the server starts
   **Then** the server logs an error that names the duplicated identifier and the
   locations where it appears
   **And** the server exits with a non-zero status before opening its listening socket

2. **Given** the areas configuration contains no duplicate item identifiers
   **When** the server starts
   **Then** the server starts successfully and serves requests as before

3. **Given** duplicate detection is implemented
   **When** the unit tests run
   **Then** a test fixture covers at least one duplicate scenario and asserts the startup
   failure path
   **And** a separate fixture covers a clean configuration and asserts successful startup

## Tasks / Subtasks

- [x] Task 1: Add duplicate-ID detection to areas config validation (AC: #1, #2)
  - [x] 1.1 In `internal/areas/config.go`, extend `validateConfig` (or add a sibling
        validator called from `Load`) so that it walks the entire `Config` and detects:
        - duplicate area IDs across all areas
        - duplicate item-group IDs across all item groups in all areas
        - duplicate item IDs across all items in all item groups in all areas
  - [x] 1.2 The error message must name the duplicated ID and the locations where it
        appears, e.g. `duplicate item id "desk29" found in item groups "room-a" and
        "room-b"` (single error per duplicate is enough; do not need to enumerate every
        pair beyond the first collision).
  - [x] 1.3 Define a sentinel error (e.g. `ErrDuplicateID`) at package level and wrap it
        with `%w` per the project's error wrapping conventions.
  - [x] 1.4 Keep detection deterministic: iterate in slice order so the same fixture
        always produces the same error message.

- [x] Task 2: Wire the failure into startup so the server exits before listening (AC: #1)
  - [x] 2.1 Confirm `areas.Load` already returns the validation error to
        `loadAndValidateAreas` in `internal/startup/server.go`. No additional plumbing
        should be needed — the existing `fmt.Errorf("load areas config: %w", err)` path
        already aborts startup before the HTTP listener is bound.
  - [x] 2.2 If the existing plumbing exits without a clear log line, add a `slog.Error`
        in `loadAndValidateAreas` or its caller that surfaces the duplicate-ID message
        (do NOT downgrade the existing error return; the error itself must still cause a
        non-zero exit).

- [x] Task 3: Tests (AC: #3)
  - [x] 3.1 Add a table-driven test in `internal/areas/config_test.go` (or a new
        `validate_test.go` in the same package) that covers at minimum:
        - duplicate item IDs in two different item groups within one area
        - duplicate item IDs in two different areas
        - duplicate item-group IDs across two different areas
        - duplicate area IDs
        - clean config (no duplicates) — passes
  - [x] 3.2 Each failing case must assert `errors.Is(err, ErrDuplicateID)` and verify
        the offending ID appears in the error string.
  - [x] 3.3 Use `require` for setup, `assert` for behavioral checks, per project rules.

- [x] Task 4: Validation
  - [x] 4.1 `go fmt ./...`
  - [x] 4.2 `go vet ./...`
  - [x] 4.3 `golangci-lint run ./...`
  - [x] 4.4 `go test ./internal/areas/...`
  - [x] 4.5 Manually verify: craft a small areas YAML with `desk29` listed twice, run the
        server with `go run cmd/sithub/main.go run --config ./sithub.toml`, confirm it
        exits non-zero with the duplicate-ID message and never binds the listener.

## Dev Notes

### Architecture & Patterns

- The areas config is a single YAML loaded once at startup via
  `internal/areas/config.go::Load`. After parsing, the file-level `validateConfig`
  enforces that every area, item group, and item has both `id` and `name`. This is the
  natural place to extend with duplicate-ID detection — the validator already runs before
  the listener starts.
- `internal/startup/server.go::loadAndValidateAreas` invokes `areas.Load` first and then
  `areas.ValidateFloorPlans` and `areas.ValidateReservations`. All of these return
  errors that propagate up and abort startup before `e.Start(...)` is called. The new
  duplicate check must follow the same pattern so it cannot be bypassed.
- Existing helpers like `Config.FindItem` return the **first** match they encounter, so a
  duplicate today silently makes the second occurrence unreachable — exactly the kind of
  inconsistency this story prevents.

### Key Code Locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Config types | `internal/areas/config.go` (`Area`, `ItemGroup`, `Item`) | Source structures to walk for duplicates |
| Existing validator | `internal/areas/config.go::validateConfig` | Extend or call alongside it |
| Sentinel error pattern | `internal/areas/config.go::ErrReservationConflict` | Mirror this for `ErrDuplicateID` |
| Load entrypoint | `internal/areas/config.go::Load` | Already calls `validateConfig` and wraps errors |
| Startup invocation | `internal/startup/server.go::loadAndValidateAreas` | Confirms failure path aborts before listening |
| Existing tests | `internal/areas/config_test.go` (e.g. `TestLoadConfigMissingAreaID`) | Pattern for new table-driven cases |

### Implementation Strategy

1. Add `ErrDuplicateID` as a package-level sentinel: `var ErrDuplicateID = errors.New("duplicate id")`.
2. Add an unexported helper, e.g. `findDuplicateIDs(cfg *Config) error`, that walks the
   config once and tracks IDs in three maps (areas, item groups, items). On the first
   collision it returns
   `fmt.Errorf("%w: <kind> %q in <context-A> and <context-B>", ErrDuplicateID, id, ctxA, ctxB)`.
3. Call this helper from `validateConfig` (or from `Load` immediately after
   `validateConfig` succeeds). The existing wrapping in `Load`
   (`fmt.Errorf("validate areas config: %w", err)`) is sufficient — `errors.Is` still
   works through the wrap.
4. No changes to handler or HTTP code are required. The startup contract is already
   "any error from `loadAndValidateAreas` aborts startup before listening".

### Anti-Patterns to Avoid

- Do NOT add a runtime check inside `FindItem`/`FindItemGroup`/`FindArea` — that's a
  post-startup guard and the bug is configuration-time. Catch it once, at startup.
- Do NOT silently log a warning and continue. The acceptance criteria require a non-zero
  exit before the listener binds.
- Do NOT introduce a new top-level "validator" entrypoint in `internal/startup` if the
  check belongs in the `areas` package next to the other validators.
- Do NOT only check item IDs. Item-group and area IDs are also referenced by route
  parameters and would produce the same class of inconsistency if duplicated.
- Do NOT short-circuit detection on the first error in a way that hides the offending
  context. Returning a single duplicate is fine; returning only "duplicates exist" is
  not.

### Testing Standards

- Table-driven tests in the `areas` package, using `require` for setup and `assert` for
  behavior (per `.claude/rules/golang.md`).
- Keep fixtures inline (string literals + `os.WriteFile` to `t.TempDir()`), matching the
  existing `TestLoadConfig` pattern in `config_test.go`.
- Each failure case asserts `errors.Is(err, ErrDuplicateID)` and that the duplicated
  identifier appears in the error string.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 30 Stories: Operator Validation, Editor Zoom Height & Optional Drill-Down]
- [Source: internal/areas/config.go]
- [Source: internal/areas/config_test.go]
- [Source: internal/startup/server.go#loadAndValidateAreas]
- [Source: .claude/rules/golang.md#Error Handling]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7

### Completion Notes List

- Added `ErrDuplicateID` sentinel and `findDuplicateIDs` walker in
  `internal/areas/config.go`; called from existing `validateConfig`. Detects duplicate
  area, item-group, and item identifiers.
- Error messages name the duplicated id and include both occurrences (e.g.
  `duplicate id: item id "desk29" in item group "r1" is also defined in item group "r2"`).
- No startup plumbing changes needed — `loadAndValidateAreas` in
  `internal/startup/server.go` already aborts before binding the listener when
  `Load` returns an error.
- Tests added: clean config passes, duplicate item id within an area, duplicate item id
  across areas, duplicate item-group id across areas, duplicate area id; plus a
  `Load`-level fixture proving the YAML→error path.
- `go fmt`, `go vet`, `go test ./...` and `golangci-lint run ./internal/areas/...` all
  clean (only pre-existing `goconst` noise on test-fixture strings remains, matching
  baseline).

### File List

- internal/areas/config.go (modified)
- internal/areas/config_test.go (modified)
