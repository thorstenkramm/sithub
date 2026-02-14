# Story 12.5: Update Documentation and E2E Tests

Status: done

## Story

As a developer,
I want API documentation and E2E tests updated to use the new terminology,
So that the entire codebase is consistent and all tests pass.

## Acceptance Criteria

1. **Given** the OpenAPI documentation references `/rooms` and `/desks` endpoints
   **When** the documentation is updated
   **Then** all endpoint paths, schemas, and descriptions use the new terminology
   **And** `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml` passes

2. **Given** the Cypress E2E tests reference rooms/desks in selectors and assertions
   **When** the tests are updated
   **Then** all E2E tests use the new terminology in routes, selectors, and assertions
   **And** `npm run test:e2e -- --browser electron` passes all tests

3. **Given** the Go code duplication check runs
   **When** `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 3` is executed
   **Then** the duplication threshold is not exceeded

4. **Given** the TypeScript code duplication check runs
   **When** `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0` is
   executed
   **Then** the duplication threshold is not exceeded

## Tasks / Subtasks

- [x] Rewrite OpenAPI documentation (AC: 1)
  - [x] Rename schemas: Room* -> ItemGroup*, Desk* -> Item*, RoomBooking* -> ItemGroupBooking*
  - [x] Update paths: /rooms -> /item-groups, /desks -> /items
  - [x] Create new endpoint files: item-groups.yaml, items.yaml, item-group-bookings.yaml
  - [x] Add bookings-history.yaml (previously undocumented)
  - [x] Delete old endpoint files: rooms.yaml, desks.yaml, room-bookings.yaml, presence.yaml
  - [x] Update all attribute field names (desk_id -> item_id, room_name -> item_group_name, etc.)
- [x] Update Cypress E2E tests (AC: 2)
  - [x] Rename desks.cy.ts -> items.cy.ts
  - [x] Rename rooms.cy.ts -> item-groups.cy.ts
  - [x] Update routes, selectors, and assertions
  - [x] Update support/flows.ts helper
- [x] Run Redocly lint (AC: 1)
- [x] Run code duplication checks (AC: 3, 4)
- [x] Run full test suite to validate (AC: 1, 2, 3, 4)

## Dev Notes

### Code Review Findings

The code review identified that OpenAPI documentation had not been updated for the
Epic 12 rename. This was addressed as a post-review fix, generating a complete rewrite
of the API documentation.

### Design Decisions

- OpenAPI docs fully rewritten (not just find-replace) to ensure consistency
- `bookings-history.yaml` added as it was previously undocumented
- `presence.yaml` removed as a standalone endpoint (covered by area-presence.yaml)

### Test Results

- 32 Cypress E2E tests passing
- 90 Vitest unit tests passing
- Go 80.5% coverage
- All lint and duplication checks clean

### References

- PRD FR42: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 12.5: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- OpenAPI docs completely rewritten with new terminology
- 4 new endpoint files created, 4 old files deleted
- Cypress E2E tests renamed and updated
- Full test suite green: 32 E2E + 90 Vitest + Go 80.5% coverage

### Key Files

- `api-doc/openapi.yaml` (complete rewrite)
- `api-doc/endpoints/item-groups.yaml` (new)
- `api-doc/endpoints/items.yaml` (new)
- `api-doc/endpoints/item-group-bookings.yaml` (new)
- `api-doc/endpoints/bookings-history.yaml` (new)
- `api-doc/endpoints/bookings.yaml` (updated)
- `web/cypress/e2e/items.cy.ts` (was desks.cy.ts)
- `web/cypress/e2e/item-groups.cy.ts` (was rooms.cy.ts)
- `web/cypress/e2e/ui-framework.cy.ts` (updated)
- `web/cypress/support/flows.ts` (updated)

### Change Log

- 2026-02-09: Story created retroactively. Implementation was part of Epic 12 commit.
