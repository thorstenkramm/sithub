# Story 19.7: Custom Icons in Areas YAML

Status: done

## Story

As an operator,
I want to specify custom MDI icons for areas, item groups, and items in the areas YAML,
So that the UI reflects the actual purpose of each space with meaningful icons.

## Acceptance Criteria

1. **Given** an area, item group, or item in the areas YAML has an `icon` field
   (e.g. `icon: mdi-office-building`)
   **When** the frontend renders that entity
   **Then** the specified MDI icon is displayed instead of the default icon

2. **Given** an entity does not have an `icon` field
   **When** the frontend renders that entity
   **Then** the current default icon is used

3. **Given** the `icon` field contains an invalid or unknown MDI icon name
   **When** the server starts
   **Then** a warning is logged but the server does not fail to start
   **And** the frontend falls back to the default icon

4. **Given** the areas API returns the `icon` attribute
   **When** the frontend receives the response
   **Then** the icon value is available for rendering at all three levels
   (area, item group, item)

## Tasks / Subtasks

- [x] Add `Icon` field to Go structs (AC: 1, 4)
  - [x] Added `Icon` field to `Area`, `ItemGroup`, and `Item` structs in
    `internal/areas/config.go`
  - [x] Field is optional and parsed from YAML configuration
- [x] Expose icon via API responses (AC: 4)
  - [x] Updated `internal/areas/handler.go` to include `icon` in area attributes
  - [x] Updated `internal/itemgroups/handler.go` to include `icon` in item group attributes
  - [x] Updated `internal/items/handler.go` to include `icon` in item attributes
- [x] Add backend tests (AC: 4)
  - [x] Updated `internal/areas/handler_test.go` to verify icon in API responses
- [x] Update frontend API types (AC: 4)
  - [x] Added `icon` field to `web/src/api/areas.ts` type
  - [x] Added `icon` field to `web/src/api/itemGroups.ts` type
  - [x] Added `icon` field to `web/src/api/items.ts` type
- [x] Install `@mdi/font` and register iconset (AC: 1)
  - [x] Installed `@mdi/font` package
  - [x] Registered `mdiFont` iconset in `web/src/plugins/vuetify.ts`
  - [x] Imported `@mdi/font/css/materialdesignicons.css` in `web/src/main.ts`
- [x] Implement icon inheritance in frontend views (AC: 1, 2)
  - [x] Updated `AreasView.vue` to render custom icon from API or fall back to default
  - [x] Updated `ItemGroupsView.vue` with icon inheritance: item group -> area -> default
  - [x] Updated `ItemsView.vue` with icon inheritance: item -> item group -> area -> default
- [x] Update schema and example files (AC: 1)
  - [x] Updated `sithub_areas.schema.json` with `icon` field definition
  - [x] Updated `sithub_areas.example.yaml` with icon examples
- [x] Verify E2E tests still pass

## Dev Notes

### Icon Inheritance Chain

Icons inherit through the hierarchy: item -> item group -> area -> default icon. If an item
has no icon set, it checks its parent item group, then the area, and finally uses the
hardcoded default.

### @mdi/font Integration

The `@mdi/font` package provides the full Material Design Icons font. The `mdiFont` iconset
is registered in Vuetify so that icon names like `mdi-office-building` resolve to the correct
glyph. This allows operators to use any MDI icon name in their YAML configuration.

### References

- Epic 19 Story 19.7: `_bmad-output/planning-artifacts/epics.md` (Epic 19 Stories section)
- FR73: `_bmad-output/planning-artifacts/prd.md`
- `internal/areas/config.go`
- `sithub_areas.schema.json`
- `sithub_areas.example.yaml`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Full-stack story: backend config, API handlers, frontend types, views, and Vuetify config
- Added `Icon` field to Area, ItemGroup, and Item Go structs
- All three API handlers expose icon in JSON:API attributes
- Installed `@mdi/font` and registered `mdiFont` iconset in Vuetify
- Icon inheritance: item -> item group -> area -> default across all three views
- AI review fix: startup now warns on invalid configured icon names instead of silently accepting them
- AI review fix: frontend icon rendering now falls back to default icons when configured values are invalid
- Updated JSON schema and example YAML with icon field
- All existing tests continue to pass

### File List

- `internal/areas/config.go` — `Icon` field on Area, ItemGroup, Item structs
- `internal/areas/config_test.go` — Tests for invalid configured icon detection
- `internal/areas/handler.go` — Icon in area API response attributes
- `internal/areas/handler_test.go` — Test for icon in API response
- `internal/itemgroups/handler.go` — Icon in item group API response attributes
- `internal/items/handler.go` — Icon in item API response attributes
- `internal/startup/server.go` — Startup warnings for invalid configured icon names
- `web/package.json` — `@mdi/font` dependency for runtime icon rendering
- `web/package-lock.json` — Locked `@mdi/font` dependency tree
- `web/src/api/areas.ts` — `icon` field in AreaAttributes type
- `web/src/api/itemGroups.ts` — `icon` field in ItemGroupAttributes type
- `web/src/api/items.ts` — `icon` field in ItemAttributes type
- `web/src/plugins/vuetify.ts` — Registered `mdiFont` iconset
- `web/src/main.ts` — Imported `@mdi/font/css/materialdesignicons.css`
- `web/src/utils/icons.ts` — Shared configured-icon validation and fallback helper
- `web/src/utils/icons.test.ts` — Unit tests for configured-icon validation and fallback
- `web/src/views/AreasView.vue` — Custom icon rendering with validated default fallback
- `web/src/views/ItemGroupsView.vue` — Icon inheritance from area with validated default fallback
- `web/src/views/ItemsView.vue` — Icon inheritance from item group, area, with validated default fallback
- `sithub_areas.schema.json` — `icon` field definition
- `sithub_areas.example.yaml` — Icon usage examples

## Senior Developer Review (AI)

- Added startup warnings for invalid configured icon names so operators get feedback without blocking server startup.
- Added a shared frontend icon resolver so invalid configured values fall back to the established default icons.
- Updated schema validation and added tests for icon validation/fallback behavior.

## Change Log

- 2026-03-21: Story implemented and verified.
- 2026-03-21: Applied AI review fixes for invalid-icon warnings and frontend fallback handling.
