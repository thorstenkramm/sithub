# Story 13.1: Clickable Breadcrumbs

Status: done

## Story

As a user,
I want breadcrumbs to be clickable links,
So that I can navigate back to any level of the hierarchy quickly.

## Acceptance Criteria

1. **Given** I am viewing items within an item group
   **When** I see the breadcrumb showing "Home > Office 1st Floor > Room 101"
   **Then** "Home" and "Office 1st Floor" are clickable links
   **And** clicking "Home" navigates to the areas list
   **And** clicking "Office 1st Floor" navigates to the item groups for that area
   **And** "Room 101" (the current page) is not clickable

2. **Given** I am viewing item groups within an area
   **When** I see the breadcrumb showing "Home > Office 1st Floor"
   **Then** "Home" is a clickable link that navigates to the areas list
   **And** "Office 1st Floor" (the current page) is not clickable

## Tasks / Subtasks

- [x] Fix breadcrumb `to` properties in ItemGroupsView.vue (AC: 2)
  - [x] "Home" breadcrumb already links to `/` - verify it works
  - [x] Last breadcrumb (area name) remains non-clickable (current page) - already correct
- [x] Fix breadcrumb `to` properties in ItemsView.vue (AC: 1)
  - [x] "Home" links to `/`
  - [x] Area name links to `/areas/:areaId/item-groups` (requires knowing `areaId`)
  - [x] Last breadcrumb (item group name) remains non-clickable - already correct
- [x] Fix breadcrumb `to` properties in ItemGroupBookingsView.vue (AC: 1)
  - [x] "Home" links to `/`
  - [x] Area name links to `/areas/:areaId/item-groups`
  - [x] Item group name links to `/item-groups/:itemGroupId/items`
  - [x] Last breadcrumb ("Bookings") remains non-clickable - already correct
- [x] Propagate `areaId` via query parameter (Decision: Path A) (AC: 1)
  - [x] In ItemGroupsView: pass `query: { areaId }` when navigating to items/bookings
  - [x] In ItemsView: read `areaId` from `route.query.areaId` for breadcrumb link
  - [x] In ItemGroupBookingsView: read `areaId` from `route.query.areaId` for breadcrumb link
  - [x] Ensure `areaId` is also forwarded from ItemsView to ItemGroupBookingsView links
  - [x] Guard: if `areaId` is missing (direct URL), render area breadcrumb as non-clickable text
  - [x] Audit all `<router-link>` templates and navigation calls that target items/bookings views
- [x] Add/update Vitest unit tests for breadcrumb navigation (AC: 1, 2)
  - [x] Test breadcrumb `to` values when `areaId` is present in query
  - [x] Test breadcrumb fallback when `areaId` is missing (no broken links)
- [x] Add Cypress E2E test for breadcrumb click navigation (AC: 1, 2)
  - [x] Full 4-level chain: Areas → ItemGroups → Items → Bookings, verify area breadcrumb clickable at each level
  - [x] Click area breadcrumb on Items page, assert URL is `/areas/:areaId/item-groups`
  - [x] Click Home breadcrumb, assert URL is `/`
  - [x] Cross-area test: navigate Area A → items, go back, navigate Area B → items, verify breadcrumb links to Area B

## Dev Notes

### Current Breadcrumb Infrastructure

Breadcrumbs are rendered in `web/src/components/PageHeader.vue` using a custom implementation:

```vue
<nav v-if="breadcrumbs?.length" class="breadcrumbs mb-2" aria-label="Breadcrumb">
  <ol class="d-flex align-center ga-1 text-body-2">
    <li v-for="(crumb, index) in breadcrumbs" :key="index" class="d-flex align-center">
      <router-link
        v-if="crumb.to && index < breadcrumbs.length - 1"
        :to="crumb.to"
        class="breadcrumb-link text-primary"
      >
        {{ crumb.text }}
      </router-link>
      <span v-else class="text-medium-emphasis">{{ crumb.text }}</span>
      ...
    </li>
  </ol>
</nav>
```

The `BreadcrumbItem` type is:

```typescript
export interface BreadcrumbItem {
  text: string;
  to?: RouteLocationRaw;
}
```

The logic already supports clickable breadcrumbs via `router-link` when `to` is provided
and the item is not the last one. **No changes needed to PageHeader.vue.**

### Current Breadcrumb State Per View

| View | Breadcrumbs | Issue |
|------|-------------|-------|
| AreasView | `[{ text: 'Home' }]` | None (root) |
| ItemGroupsView | `[{ text: 'Home', to: '/' }, { text: areaName }]` | Correct |
| ItemsView | `[Home, Area(?), ItemGroup]` | Area link is conditional/broken |
| ItemGroupBookingsView | `[Home, Area->/, ItemGroup, Bookings]` | Area links to `/` not area |

### Key Problem: Missing `areaId` Context

`ItemsView` and `ItemGroupBookingsView` routes only receive `itemGroupId`:
- Route: `/item-groups/:itemGroupId/items`
- Route: `/item-groups/:itemGroupId/bookings`

They need `areaId` to build the breadcrumb link `/areas/:areaId/item-groups`.

### Decision: Query Parameter Approach (Tree of Thoughts Analysis)

Four approaches were evaluated. **Path A (Query Parameter)** was selected.

**Path A - Query Parameter (SELECTED)**
Pass `areaId` as a query param when navigating from ItemGroupsView. Frontend-only,
zero backend changes, no new state management. Graceful degradation: if `areaId` is
missing (direct URL access), the area breadcrumb falls back to non-clickable text.

```typescript
// ItemGroupsView navigates with areaId in query:
router.push({
  name: 'items',
  params: { itemGroupId },
  query: { areaId: route.params.areaId }
})

// ItemsView reads areaId from query:
const areaId = computed(() => route.query.areaId as string | undefined)
```

**Rejected alternatives:**
- Path B (Route Restructure): Nesting items under `/areas/:areaId/...` would be a
  breaking URL change affecting bookmarks, shared links, and all tests. High risk.
- Path C (Pinia Store): State lost on page refresh/direct URL navigation. Adds
  complexity without solving the bookmark/refresh case.
- Path D (API Derive): Requires backend changes to include `areaId` in item-group
  API responses. Over-engineered for a breadcrumb link.

### Route Structure Reference

```
/                                          → AreasView
/areas/:areaId/item-groups                 → ItemGroupsView
/item-groups/:itemGroupId/items            → ItemsView
/item-groups/:itemGroupId/bookings         → ItemGroupBookingsView
```

### Testing Requirements

- Vitest: Test breadcrumb computed properties return correct `to` values
- Vitest: Test missing `areaId` fallback (no broken `/areas/undefined/...` links)
- Cypress E2E: Full 4-level navigation chain with breadcrumb clicks at each level
- Cypress E2E: Click breadcrumbs and assert destination URL, not just visual presence
- Cypress E2E: Cross-area navigation to verify no stale `areaId` carry-over

### Project Structure Notes

- Breadcrumb type defined in `web/src/components/PageHeader.vue`
- Each view defines its own breadcrumbs as a `computed` property
- Follow existing pattern: no new components or composables needed

### References

- PRD FR40: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 13.1: `_bmad-output/planning-artifacts/epics.md`
- PageHeader component: `web/src/components/PageHeader.vue`
- Router: `web/src/router/index.ts`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None - clean implementation with no debugging required.

### Completion Notes List

- Implemented areaId propagation via query parameter (Path A from Tree of Thoughts analysis)
- Updated ItemGroupsView to pass `areaId` in query when navigating to items and bookings
- Updated ItemsView breadcrumbs to read `areaId` from query and link to correct area page
- Updated ItemGroupBookingsView breadcrumbs with area link, item group link to items view
- Updated "View Item Group Bookings" link in ItemsView to forward areaId
- Graceful degradation: missing areaId renders area breadcrumb as non-clickable text
- Added 6 new Vitest tests (breadcrumb links with/without areaId, navigation with areaId)
- Added 4 new Cypress E2E tests (area breadcrumb click, Home click, item groups Home click, cross-area)
- All 96 Vitest tests pass, all 36 Cypress E2E tests pass
- Zero code duplication (0%), type-check clean, ESLint clean, build clean

### File List

- `web/src/views/ItemGroupsView.vue` (modified: goToItems passes areaId, View Bookings link passes areaId)
- `web/src/views/ItemsView.vue` (modified: breadcrumbs read areaId from query, area link corrected, bookings link forwards areaId)
- `web/src/views/ItemGroupBookingsView.vue` (modified: breadcrumbs read areaId from query, area and item group links corrected)
- `web/src/views/ItemGroupsView.test.ts` (modified: added areaId navigation test)
- `web/src/views/ItemsView.test.ts` (modified: added query to route mock, added breadcrumb tests)
- `web/src/views/ItemGroupBookingsView.test.ts` (modified: added query to route mock, added breadcrumb tests)
- `web/cypress/e2e/breadcrumbs.cy.ts` (new: 4 E2E tests for breadcrumb navigation)

### Change Log

- 2026-02-09: Story implemented. All ACs satisfied. 96 Vitest + 36 Cypress E2E tests pass.