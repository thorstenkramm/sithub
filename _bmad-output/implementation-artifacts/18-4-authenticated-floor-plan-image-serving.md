# Story 18.4: Authenticated Floor Plan Image Serving

Status: done

## Story

As a user,
I want floor plan images served through an authenticated endpoint,
So that floor plans are protected from unauthorized access.

## Acceptance Criteria

1. **Given** I am authenticated
   **When** I request a floor plan image via the API
   **Then** the image is returned with the correct content type

2. **Given** I am not authenticated
   **When** I request a floor plan image
   **Then** the request is denied with 401

3. **Given** I request an image that does not exist
   **When** the request is processed
   **Then** I receive a 404 response

## Tasks / Subtasks

- [x] Create `FloorPlanHandler()` in `internal/areas/floor_plan_handler.go`
  - [x] Accept `floorPlansDir` parameter
  - [x] Return 404 if floor plans not configured (empty dir)
  - [x] Validate filename: reject slashes and serve only filename-only references
  - [x] Check file extension against MIME type map
  - [x] Read file and serve with correct content type via `c.Blob()`
  - [x] Return 404 for missing files or unsupported formats
- [x] Define MIME type map: `.jpg`/`.jpeg` → `image/jpeg`, `.png` → `image/png`,
  `.svg` → `image/svg+xml`
- [x] Register route in `internal/startup/server.go`
  - [x] `GET /api/v1/floor-plans/:filename` behind `requireAuth` middleware
  - [x] Pass `floorPlansDir` from config
- [x] Create OpenAPI documentation `api-doc/endpoints/floor-plans.yaml`
- [x] Add path reference to `api-doc/openapi.yaml`
- [x] Add unit tests in `internal/areas/floor_plan_handler_test.go`
  - [x] ServesImage — PNG served with correct content type
  - [x] ServesSVG — SVG served with correct content type
  - [x] NotFound — 404 for missing file
  - [x] NotConfigured — 404 when floor plans dir is empty
  - [x] UnsupportedFormat — 404 for `.txt` extension
  - [x] PathTraversal — 404 for `../` in filename
- [x] Add route-level auth coverage in `internal/startup/server_test.go`
  - [x] Unauthenticated request to `/api/v1/floor-plans/:filename` returns 401
- [x] Run tests, linting, and duplication checks

## Dev Notes

### Security

The handler rejects filenames containing `/` or `\` characters and only serves filename-only
references. This prevents directory traversal attacks. The route is also behind the
`requireAuth` middleware so only authenticated users can access floor plans.

### References

- Epic 18 Story 18.4: `_bmad-output/planning-artifacts/epics.md`
- FR63: `_bmad-output/planning-artifacts/prd.md`
- `internal/areas/floor_plan_handler.go`
- `api-doc/endpoints/floor-plans.yaml`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Created floor plan handler with security-focused filename validation
- Handler uses `c.Blob()` for efficient binary response
- Added 6 unit tests covering all acceptance criteria and security edge cases
- Added a startup-level auth test to prove the route returns 401 without authentication
- Created OpenAPI spec for the endpoint
- All Go tests, API doc linting, and duplication checks pass

### File List

- `internal/areas/floor_plan_handler.go` — Floor plan image serving handler
- `internal/areas/floor_plan_handler_test.go` — 6 unit tests
- `internal/startup/server.go` — Route registration with `requireAuth`
- `internal/startup/server_test.go` — Verifies unauthenticated floor plan requests return 401
- `api-doc/endpoints/floor-plans.yaml` — OpenAPI endpoint spec
- `api-doc/openapi.yaml` — Added path reference

## Change Log

- 2026-03-13: Story implemented and verified.
- 2026-03-13: Code review fixes added route-level auth coverage and aligned filename validation notes.
