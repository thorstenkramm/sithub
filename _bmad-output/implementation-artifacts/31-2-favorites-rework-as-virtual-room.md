# Story 31.2: Favorites Rework as Virtual Room

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user with favorite desks across multiple areas,
I want my favorites grouped into a dedicated "Favorites" room with clear visual markers
in every view,
so that I can find and book my preferred desks quickly without scanning unrelated areas.

## Acceptance Criteria

1. **Given** I have at least one bookable item marked as a favorite
   **When** I open the area/room overview
   **Then** a tile labeled "Favorites" appears as the first tile
   **And** the tile behaves like any other room tile (drill-down to its items, free/busy
   indicators, identical interaction model)

2. **Given** I have no items marked as favorites
   **When** I open the area/room overview
   **Then** the "Favorites" tile is not shown

3. **Given** I am on a screen that previously allowed adding an area or room to favorites
   **When** the page renders
   **Then** no control to favorite an area or room is available
   **And** only bookable items (desks) can be added to favorites

4. **Given** I am viewing the weekly table view and one or more items are favorites
   **When** the table renders
   **Then** each favorite item row displays a heart icon
   **And** clicking the heart icon removes that item from my favorites
   **And** no sorting or filtering by favorites is offered in the table view

5. **Given** I am viewing a floor plan and one or more items shown on the plan are
   favorites
   **When** an item is in the free state
   **Then** a heart icon is rendered with its center positioned exactly at the
   bottom-right corner of the item marker (matching the reference image
   `epic-31-favorite-heart.png`)
   **And** clicking the heart icon removes the item from my favorites

6. **Given** an item shown on the floor plan is busy
   **When** the floor plan renders
   **Then** no heart icon is shown for that item, regardless of favorite status

7. **Given** I remove an item from favorites via any of the heart icons (tile, table,
   floor plan)
   **When** the change is applied
   **Then** the item is removed across all views consistently and the "Favorites"
   tile disappears once no favorites remain

## Tasks / Subtasks

- [ ] Task 1: Trim `useFavorites` to items-only and migrate stored data (AC: #3)
  - [ ] 1.1 Remove `favoriteItemGroups`, `isItemGroupFavorite`,
        `toggleItemGroupFavorite` from `web/src/composables/useFavorites.ts` and its
        return signature. Keep `favoriteItems`, `isItemFavorite`,
        `toggleItemFavorite`, and `favoriteItemsForArea`.
  - [ ] 1.2 On first load, also `localStorage.removeItem('sithub_favorite_item_groups')`
        once so stale data is purged. Wrap in `getSafeLocalStorage()` and tolerate
        absence (covered by existing helper).
  - [ ] 1.3 Update `useFavorites.test.ts` to drop item-group cases and add coverage
        for the cleanup of `sithub_favorite_item_groups` on first load.

- [ ] Task 2: Remove favorite-on-area/room UI affordances (AC: #3)
  - [ ] 2.1 In `web/src/views/ItemGroupsView.vue`:
        - Delete the `ig-favorite-heart` button (around line ~298) and its handler
          `handleToggleItemGroupFavorite` (around line ~513). Remove the related
          computed in `sortedItemGroups.igFavs`; render item groups in YAML order
          only.
        - Delete the third-level "favorites promoted" block (around lines ~145–224).
          Favorites are no longer rendered here — they live in the new virtual
          area, see Task 3.
  - [ ] 2.2 Search for any other "ig-favorite" or "favorite-item-tile" references
        in the codebase and remove or update them. Common spots: tests under
        `web/src/views/__tests__/`, e2e specs in `web/cypress/e2e/`.
  - [ ] 2.3 In `web/src/views/AreasView.vue`: confirm there is no per-area favorite
        button today (there is none as of `main`, but verify after the refactor).

- [ ] Task 3: Add the virtual "Favorites" area (AC: #1, #2)
  - [ ] 3.1 Define a constant `FAVORITES_AREA_ID = '__favorites__'` in
        `web/src/composables/useFavorites.ts` and export it.
  - [ ] 3.2 In `AreasView.vue`, after `fetchAreas()` succeeds, prepend a synthetic
        area object to `areas.value` **only when** `favoriteItems.value.length > 0`:
        ```ts
        {
          id: FAVORITES_AREA_ID,
          attributes: {
            name: t('favorites.areaName'),
            description: t('favorites.areaSubtitle'),
            icon: '$heart'
          }
        }
        ```
        It must be the first tile (AC #1). When the count drops to 0 (after a
        toggle), the tile must disappear without page reload — bind via a
        `computed` over `favoriteItems`.
  - [ ] 3.3 Make the Favorites tile clickable like any other; it routes to a new
        view `FavoritesView` (Task 4). Hide the "Today's presence" button only on
        the Favorites tile (it has no real area to compute presence for).
  - [ ] 3.4 Add i18n keys `favorites.areaName` ("Favorites") and
        `favorites.areaSubtitle` (e.g. "Your saved desks") to all locales under
        `web/src/locales/` (en, de, es, fr, uk). Translations are best-effort; flag
        for native-speaker review if needed.

- [ ] Task 4: Wire the Favorites route and view (AC: #1)
  - [ ] 4.1 Add a router entry in `web/src/router/index.ts`:
        ```ts
        { path: '/favorites', name: 'favorites',
          component: () => import('../views/FavoritesView.vue') }
        ```
  - [ ] 4.2 In `AreasView.vue`, route the Favorites tile to `{ name: 'favorites' }`;
        all other tiles continue to route to `item-groups`.
  - [ ] 4.3 Create `web/src/views/FavoritesView.vue`. It reuses the same UI
        components as `ItemGroupsView` would but treats the favorites set as a
        single virtual item group containing all favorite items. Concretely:
        - Read `favoriteItems` from `useFavorites()`.
        - For each unique `(areaId, itemGroupId)` represented in the favorites,
          fetch availability via the existing `fetchItemGroupAvailability` API and
          render the same weekly availability indicators ItemGroupsView already
          uses for AC #1's "free/busy indicators".
        - Render each favorite as a card identical in structure to the
          `favorite-item-tile` cards that previously lived in `ItemGroupsView`
          (reuse the markup, just relocated). Keep the `data-cy="favorite-item-tile"`
          selector for E2E continuity.
        - Clicking a card navigates to the matching `items` route
          `{ name: 'items', params: { itemGroupId }, query: { areaId } }` — same
          as before.
  - [ ] 4.4 The breadcrumb on `FavoritesView` shows
        `Home > Favorites` (use existing `PageHeader` component pattern from
        `AreasView` and `ItemGroupsView`).

- [ ] Task 5: Heart icon on the weekly matrix table (AC: #4)
  - [ ] 5.1 In `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue`,
        add a heart icon to the sticky desk-name cell, rendered **only** when the
        item is a favorite. Visual: small filled heart (`$heart`,
        `color="error"`, size 14), placed inline after the desk name and any
        equipment/warning icons.
  - [ ] 5.2 Click handler removes the item from favorites via `toggleItemFavorite`.
        Use `data-cy="matrix-favorite-heart-${item.item_id}"`.
  - [ ] 5.3 The heart is **only** an outlined-then-filled heart for items that are
        already favorites. Non-favorites in the table do not show any heart icon
        at all (the user's brief is "favorites must be marked"; adding/removing
        from the table is one-way: removal). Adding to favorites still happens via
        the day-mode item view (`item-favorite-heart` data-cy already exists in
        `ItemsView.vue`).
  - [ ] 5.4 No sort or filter control by favorites in the table view (AC #4).
        Confirm none exists today and add no new control.
  - [ ] 5.5 The matrix row needs an `isFavorite` boolean. Inject it via
        `AreaWeeklyMatrixView.vue` by mapping the `useFavorites()` set onto each
        row before passing to `AreaWeeklyMatrixRow`. Pass as a new prop.

- [ ] Task 6: Heart icon on the floor plan (AC: #5, #6, #7)
  - [ ] 6.1 In `web/src/components/InteractiveFloorPlan.vue`, locate every
        `.fp-item--free` rectangle (lines ~193–208 for the area-level block, and
        lines ~213–232 for the item-level block).
  - [ ] 6.2 Add a child element only when the item is a favorite:
        ```html
        <v-icon
          v-if="isFloorPlanItemFavorite(pos.itemId)"
          class="fp-favorite-heart"
          size="14"
          color="error"
          :data-cy="`fp-favorite-heart-${pos.itemId}`"
          @click.stop="removeFloorPlanFavorite(pos)"
        >$heart</v-icon>
        ```
  - [ ] 6.3 CSS: position the heart absolutely so its **center** sits exactly on
        the bottom-right corner of the rectangle (AC #5):
        ```css
        .fp-favorite-heart {
          position: absolute;
          right: 0;
          bottom: 0;
          transform: translate(50%, 50%);
          background: rgba(var(--v-theme-surface), 0.9);
          border-radius: 50%;
          padding: 2px;
          pointer-events: auto;
          cursor: pointer;
          z-index: 2;
        }
        ```
        The translate(50%, 50%) pulls the icon outward by half its own width/height
        so the visual center aligns with the corner. Ensure `.fp-item--free` has
        `position: relative` (it already does in the existing styles; verify).
  - [ ] 6.4 The reference image is at
        `_bmad-output/planning-artifacts/epic-31-favorite-heart.png`. Compare the
        rendered output side-by-side: the heart's center must coincide with the
        corner — not inset, not floating outside.
  - [ ] 6.5 Add the heart **only** to the free-state rectangles (AC #6). Do NOT
        add it to `fp-item--busy` or `fp-item--reserved`. The current code has two
        rendering paths (clickable area-level desks at ~193–208 and item-level
        free items at ~220–232). Both need the icon; busy and reserved paths
        must remain unchanged.
  - [ ] 6.6 Click on the heart removes the favorite (`@click.stop` is essential
        so the desk's own click — which starts a booking — does not fire). The
        Vuetify `v-icon` accepts native click events.
  - [ ] 6.7 `removeFloorPlanFavorite(pos)` resolves `(areaId, itemGroupId, itemId,
        itemName)` from the floor plan's positional data already available in the
        component (`enrichedPositions` maps each `pos` back to its item group).
  - [ ] 6.8 `isFloorPlanItemFavorite(itemId)` reads from `useFavorites().favoriteItems`
        — match by `itemId` only (the floor plan is scoped to one area, so an
        `itemId` collision with a different area's favorite is not possible in
        practice; if it ever became possible, gate by the floor plan's own
        `areaId`).

- [ ] Task 7: i18n & visual polish (AC: #1, #4, #5)
  - [ ] 7.1 Add to all locales under `web/src/locales/`:
        - `favorites.areaName`
        - `favorites.areaSubtitle`
        - `favorites.empty` (e.g. "No favorites yet — open a desk and tap the
          heart to add one.")
        - `favorites.removed` (snackbar text after removal — already exists for
          the existing favorites flow as `items.removedFromFavorites`; reuse if
          wording fits, otherwise add).
  - [ ] 7.2 Tooltip on every heart-removal icon: localized "Remove from favorites"
        (`favorites.removeTooltip`).

- [ ] Task 8: Tests
  - [ ] 8.1 Unit (Vitest):
        - `useFavorites.test.ts`: drop item-group cases; assert that
          `sithub_favorite_item_groups` is cleared on first load; cover
          add/remove/check for items.
        - `AreaWeeklyMatrixRow.test.ts` (new or existing): asserts the
          `matrix-favorite-heart-*` element is rendered for a favorite row only,
          and clicking it calls the toggle.
  - [ ] 8.2 Component (Vitest + Vue Test Utils):
        - `InteractiveFloorPlan.test.ts`: asserts the heart icon is rendered
          inside `fp-item--free` for a favorite item, is absent on `fp-item--busy`,
          and clicking it removes the favorite (via a stubbed `useFavorites`).
        - `AreasView.test.ts`: asserts the Favorites tile is the first tile when
          favorites exist and is absent otherwise.
        - `FavoritesView.test.ts` (new): asserts the view renders one card per
          favorite and routes to `items` on click.
  - [ ] 8.3 Cypress E2E:
        - `cypress/e2e/favorites.cy.ts` (new or extend existing): cover the
          golden path — log in, add a favorite via day-mode item view, return
          home, see Favorites tile as first, click into it, remove via heart, see
          tile disappear.
        - Also cover the floor plan and table-view heart removal paths if a
          floor-plan-enabled area is available in the dev fixture.
  - [ ] 8.4 Run:
        ```
    cd web
        npx vitest run
        npm run type-check
        npm run lint
        npm run build
        npm run test:e2e -- --browser electron
        ```

## Dev Notes

### Architecture & Patterns

- Favorites today are a per-device localStorage construct
  (`web/src/composables/useFavorites.ts`). Two separate sets are stored:
  `sithub_favorite_item_groups` (for areas/rooms) and `sithub_favorite_items`
  (for desks). This story removes the first set entirely and keeps the second
  as the only source of truth.
- The frontend is the only authoritative location for favorites — no backend
  table, no API. Do not introduce one.
- The new "virtual area" is **purely a frontend concept**: the backend never
  hears about it. Routing must therefore detect `__favorites__` (or whatever
  constant is chosen) before issuing area-scoped API calls. The simplest route
  is to give it its own component (`FavoritesView.vue`) so we never accidentally
  hit `/api/v1/areas/__favorites__/...`.
- The weekly matrix and floor plan currently treat favorites as decoration only
  in the day-mode item view. Adding hearts to the matrix and floor plan is the
  net new UI surface in this story.
- AC #5's heart placement is precise: image
  `_bmad-output/planning-artifacts/epic-31-favorite-heart.png` shows the heart's
  center sitting exactly on the rectangle corner (half inside, half outside the
  rectangle). The CSS recipe is `position: absolute; right: 0; bottom: 0;
  transform: translate(50%, 50%);` — confirm by overlaying the screenshot on the
  reference image during implementation.

### Key Code Locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Favorites composable | `web/src/composables/useFavorites.ts` | Trim to items-only, add storage cleanup |
| Item-group favorites UI | `web/src/views/ItemGroupsView.vue` (lines ~145–224, ~298–308, ~513–528) | Remove area/room favoriting and the third-level promoted favorites block |
| Day-mode item heart | `web/src/views/ItemsView.vue` (line ~239) | Keep — this is still where users **add** to favorites |
| Areas overview | `web/src/views/AreasView.vue` | Prepend the virtual Favorites tile when count > 0 |
| Router | `web/src/router/index.ts` | Add `/favorites` route |
| New view | `web/src/views/FavoritesView.vue` | New file; reuses card markup from ItemGroupsView's removed block |
| Weekly matrix row | `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue` | Add heart in sticky desk-name cell |
| Weekly matrix view | `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue` | Pass `isFavorite` per row |
| Floor plan free rectangles | `web/src/components/InteractiveFloorPlan.vue` (lines ~193–208 and ~213–232) | Add heart child for free favorites |
| Locales | `web/src/locales/{en,de,es,fr,uk}.json` | Add new `favorites.*` keys |

### Implementation Strategy

1. Trim the composable first and run the existing tests — they will fail, but
   the failure surface tells you exactly what UI references must be removed.
2. Remove the area/room favoriting UI (Task 2). The third-level promoted block
   in `ItemGroupsView.vue` is the biggest deletion.
3. Build `FavoritesView.vue` by lifting the `favorite-item-tile` card markup
   verbatim from the removed ItemGroupsView block; this preserves the visual
   pattern and the `data-cy="favorite-item-tile"` selector that existing E2E
   tests rely on.
4. Wire the virtual area in AreasView and the router. Verify navigation works
   end-to-end before touching the matrix or floor plan.
5. Add the matrix heart (small change — sticky cell only).
6. Add the floor plan heart (the most visually fiddly part — keep the CSS
   isolated to `.fp-favorite-heart`).
7. i18n + tests last.

### Anti-patterns to Avoid

- Do NOT introduce a backend favorites endpoint. Local storage is the right
  shape per Story 19.8 and AC #3.
- Do NOT keep the `favoriteItemGroups` set "for backwards compatibility". The
  user's brief explicitly says area/room favoriting is no longer supported;
  leaving the set in place would just collect drift.
- Do NOT render the heart inside the floor plan rectangle (inset). AC #5 calls
  for the heart's center to land on the corner — half inside, half outside.
- Do NOT render the heart on `fp-item--busy` or `fp-item--reserved` (AC #6).
- Do NOT add a "favorites" filter or sort to the weekly matrix (AC #4). Heart
  is decorative + one-way removal only in the table.
- Do NOT use the matrix heart to **add** favorites. The matrix heart only
  appears for items that are already favorites and is purely a removal
  affordance. Adding still happens via the day-mode item view.
- Do NOT call `/api/v1/areas/__favorites__/...` or any other backend route for
  the synthetic area. Detect the synthetic ID before any fetch.

### Testing Standards

- Unit tests: Vitest, mock `getSafeLocalStorage` where needed.
- Component tests: Vue Test Utils + Vitest, mock `useFavorites` per the
  established pattern in `InteractiveFloorPlan.test.ts`.
- E2E: Cypress against the dev server with `cy.login()` custom command and
  `data-cy` selectors only. Use intercept aliases for API synchronization
  (`@listAreas`, `@itemGroupAvailability`).
- Maintain the existing `data-cy="favorite-item-tile"` selector when relocating
  the card markup so existing E2E specs do not break unnecessarily.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 31 Stories: Live Updates, Favorites Rework & Areas Config Hint]
- [Source: _bmad-output/planning-artifacts/epic-31-favorite-heart.png]
- [Source: web/src/composables/useFavorites.ts]
- [Source: web/src/views/AreasView.vue]
- [Source: web/src/views/ItemGroupsView.vue]
- [Source: web/src/views/ItemsView.vue]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue]
- [Source: web/src/components/InteractiveFloorPlan.vue]
- [Source: web/src/router/index.ts]
- [Source: _bmad-output/implementation-artifacts/19-8-favorites.md]
- [Source: _bmad-output/implementation-artifacts/22-6-favorites-heart-icon-visibility-fix.md]
- [Source: _bmad-output/implementation-artifacts/23-1-booking-tile-heart-icon-position.md]
- [Source: .claude/rules/vue.md]
- [Source: .claude/rules/cypress.md]
- [Source: .claude/rules/feedback.md]

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
