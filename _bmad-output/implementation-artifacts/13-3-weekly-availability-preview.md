# Story 13.3: Weekly Availability Preview

Status: review

## Story

As a user,
I want to see a weekly availability preview on item group tiles,
So that I can quickly identify which days have open items without clicking into each group.

## Acceptance Criteria

1. **Given** I am viewing item groups within an area
   **When** the page loads
   **Then** I see a calendar week selector above the item group tiles
   **And** the current week is pre-selected
   **And** only the next 8 weeks are available for selection
   **And** each week option displays the Monday date and week number in locale-aware format
   (e.g., "2026-03-16 - Week 12")

2. **Given** I have selected a calendar week
   **When** the item group tiles are displayed
   **Then** each tile shows weekday indicators (MO, TU, WE, TH, FR)
   **And** green indicates at least one item is available on that day
   **And** red indicates all items in the group are fully booked on that day
   **And** each indicator uses a secondary visual cue in addition to color
   (e.g., filled circle for available, empty circle for booked) to meet WCAG A

3. **Given** the backend receives a request for weekly availability
   **When** `GET /api/v1/areas/:area_id/item-groups/availability?week=YYYY-Www` is called
   **Then** it returns per-day availability counts for each item group in the area
   **And** the response includes item group ID, day, total items, and available count
   **And** the response completes within the NFR1 performance target (2s at p95)

4. **Given** I change the selected calendar week
   **When** the new week is applied
   **Then** the weekday indicators update to reflect availability for the new week

## Tasks / Subtasks

- [x] Create backend availability endpoint (AC: 3)
  - [x] Add route `GET /api/v1/areas/:area_id/item-groups/availability`
  - [x] Parse `week` query parameter in ISO 8601 week format (`YYYY-Www`)
  - [x] Query bookings count per item group per day for the 5 weekdays (Mon-Fri)
  - [x] Query total items per item group from spaces config
  - [x] Return JSON:API collection with per-day availability data
  - [x] Add handler tests with table-driven cases
- [x] Define availability API response type (AC: 3)
  - [x] JSON:API resource type: `item-group-availability`
  - [x] Attributes: `item_group_id`, `item_group_name`, `days` (array of day objects)
  - [x] Each day object: `date`, `weekday`, `total`, `available`
- [x] Add frontend API service for availability (AC: 3)
  - [x] Create `fetchWeeklyAvailability(areaId, week)` in a new API file
  - [x] Define TypeScript types for the response
- [x] Add calendar week selector to ItemGroupsView (AC: 1, 4)
  - [x] Vuetify `v-select` component above the item group tiles
  - [x] Generate next 8 weeks starting from current week
  - [x] Format: "YYYY-MM-DD - Week NN"
  - [x] Default to current week on page load
  - [x] Fetch availability data when week selection changes
- [x] Add weekday indicators to item group tiles (AC: 2)
  - [x] Display MO/TU/WE/TH/FR indicators on each tile
  - [x] Green filled circle = available, Red empty circle = fully booked
  - [x] Use `aria-label` on each indicator for screen reader accessibility
  - [x] Ensure indicators are visible and understandable without color alone (WCAG A)
- [x] Add Vitest unit tests (AC: 1, 2, 3, 4)
  - [x] Test week selector card renders
  - [x] Test availability fetch on mount
  - [x] Test availability indicators display with weekday labels
  - [x] Test dot classes (available vs booked)
  - [x] Test graceful handling of availability fetch failure
- [x] Add Cypress E2E test (AC: 1, 2, 4)
  - [x] Verify week selector visible on item groups page
  - [x] Verify availability indicators with weekday labels (MO-FR)
  - [x] Select a different week, verify availability reloads
- [x] Update OpenAPI documentation (AC: 3)
  - [x] Add `item-group-availability.yaml` endpoint documentation
  - [x] Add schemas: DayAvailability, ItemGroupAvailabilityAttributes, etc.
  - [x] Lint with Redocly

## Dev Notes

### New Backend Endpoint

This is the first "aggregate" endpoint in the API. It queries across item groups
and bookings to return availability counts. The query pattern:

```sql
SELECT ig.id, ig.name, b.booking_date, COUNT(b.id) as booked_count
FROM item_groups ig
CROSS JOIN (SELECT date FROM weekdays) d
LEFT JOIN bookings b ON b.item_group_id = ig.id AND b.booking_date = d.date
WHERE ig.area_id = ?
GROUP BY ig.id, d.date
```

Note: Item groups and items come from YAML config (spaces), not from a database
table. The query needs to work with the spaces config for item counts and the
bookings table for booking counts.

### ISO 8601 Week Format

The `week` parameter uses ISO 8601 format: `2026-W12`. Go's `time` package does
not natively parse ISO weeks. Use a helper function to convert `YYYY-Www` to the
Monday date, then generate Mon-Fri dates.

### Week Selector Component

Use a Vuetify `v-select` with computed options. Each option value is the ISO week
string, and the display text is the formatted Monday date + week number.

### Performance Consideration

This endpoint queries 5 days x N item groups. For typical installations (< 50 item
groups), a single query with GROUP BY should complete well within the 2s target.
No caching needed for v1.

### References

- PRD FR36: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 13.3: `_bmad-output/planning-artifacts/epics.md`
- ItemGroupsView: `web/src/views/ItemGroupsView.vue`
- Bookings store: `internal/bookings/store.go`
- Spaces config: `internal/spaces/config.go`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None

### Completion Notes List

- Created `internal/itemgroups/availability_handler.go` with ISO week parsing,
  availability calculation, and JSON:API response
- Created `internal/itemgroups/availability_handler_test.go` with 8 tests covering
  parsing, weekday generation, and handler behavior
- Registered route in `internal/startup/server.go`
- Created `web/src/api/itemGroupAvailability.ts` with types and fetch function
- Updated `web/src/views/ItemGroupsView.vue` with week selector and availability
  indicators (filled/empty dots with WCAG-compliant dual visual cues)
- Updated `web/src/views/ItemGroupsView.test.ts` with 11 tests (7 new)
- Created `web/cypress/e2e/availability.cy.ts` with 3 E2E tests
- Updated `web/cypress/support/flows.ts` with availability intercept
- Created `api-doc/endpoints/item-group-availability.yaml` OpenAPI doc
- Updated `api-doc/openapi.yaml` with availability schemas and booker_name fix

### File List

- `internal/itemgroups/availability_handler.go` (new)
- `internal/itemgroups/availability_handler_test.go` (new)
- `internal/startup/server.go` (modified)
- `web/src/api/itemGroupAvailability.ts` (new)
- `web/src/views/ItemGroupsView.vue` (modified)
- `web/src/views/ItemGroupsView.test.ts` (modified)
- `web/cypress/e2e/availability.cy.ts` (new)
- `web/cypress/support/flows.ts` (modified)
- `api-doc/endpoints/item-group-availability.yaml` (new)
- `api-doc/openapi.yaml` (modified)
