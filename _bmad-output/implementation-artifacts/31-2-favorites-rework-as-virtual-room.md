# Story 31.2: Favorites Rework as Virtual Room

Status: done

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

### Review Findings

- [x] [Review][Patch] Favorites state is not shared between `useFavorites()` callers, so removal in one mounted component does not update other mounted consumers and can leave the Favorites tile/heart markers stale, violating AC #7's "removed across all views consistently" requirement [web/src/composables/useFavorites.ts:77]
- [x] [Review][Patch] `FavoritesView` keys availability by `itemGroupId` only, so favorites from different areas with the same item-group id overwrite each other's indicators and render incorrect free/busy dots for cross-area favorites [web/src/views/FavoritesView.vue:141]
- [x] [Review][Patch] The Favorites virtual tile is hidden behind the real-area empty state, so a user with favorites but no loaded real areas sees `areas-empty` instead of the required first Favorites tile [web/src/views/AreasView.vue:19]
- [x] [Review][Patch] Required regression coverage for the new matrix and floor-plan heart removal paths is missing: no test asserts `matrix-favorite-heart-*`, no test asserts `fp-favorite-heart-*`, and the story's own completion notes claim these paths are covered when they are not [web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue:27]
- [x] [Review][Patch] Favorites-mode week bookings are disabled by the old single-room guard: `submitWeekBookings()` returns when `activeItemGroupId` is null, so the `/favorites` week view can show selected days and a confirm button but never creates bookings [web/src/views/ItemsView.vue:1502]
- [x] [Review][Patch] Favorites-mode week reload does not preserve or prune week state correctly: `loadFavoriteWeekData()` always clears `weekBookingResults` even when called with `keepResults=true`, and silent live refresh does not prune selections that became unavailable [web/src/views/ItemsView.vue:1728]
- [x] [Review][Patch] Favorites-mode item loading swallows initial fetch failures per item group, so connection loss or API errors render as an empty/partial favorites room instead of the normal ItemsView error state expected from a real room [web/src/views/ItemsView.vue:1668]
- [x] [Review][Patch] Removing a favorite from the favorites week view removes the row from `weekData` but leaves any selected `weekSelections` for that item intact, allowing stale selections to survive after the item disappears [web/src/views/ItemsView.vue:1094]

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

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run` — 404 tests pass (47 files)
- `cd web && npm run type-check && npm run lint && npm run build` — all green

### Completion Notes List

- **`useFavorites` trimmed to items-only.** Removed `favoriteItemGroups`,
  `isItemGroupFavorite`, `toggleItemGroupFavorite`, and the legacy storage
  read. Added a one-shot purge that removes the legacy
  `sithub_favorite_item_groups` key from localStorage on first use, gated by
  a module-level flag with a test-only reset. Exported a
  `FAVORITES_AREA_ID` constant for views that need to refer to the synthetic
  area.
- **ItemGroupsView cleanup.** Deleted the third-level "favorites promoted"
  block, the `ig-favorite-heart` button, the `sortedItemGroups` shaping
  (now just YAML order), and the two toggle handlers. `useFavorites()` is
  still called for its purge side effect — leaving it absent would mean the
  legacy key only got purged once the user actually opened the home view.
- **Virtual `Favorites` area** on `AreasView`: prepended a tile gated on
  `favoriteItems.length > 0` so it disappears when the last favorite is
  removed without a page reload. The tile routes to a new
  `/favorites` page rather than into the regular `item-groups` flow so the
  synthetic id never reaches the area-scoped API.
- **`FavoritesView.vue`** (new): renders one card per favorite item with the
  same availability indicator dots used in `ItemGroupsView`. Availability
  fetch groups favorites by `areaId` and issues one
  `fetchWeeklyAvailability` per area so a user with N favorites spread over
  K areas makes K calls, not N. The heart button removes the favorite; a
  `watch` on `favoriteItems` re-fetches indicators when the set changes;
  `useLiveBookingRefresh` keeps the indicators current when other users
  book/cancel.
- **Matrix heart icon.** Added an inline `$heart` icon in the sticky
  desk-name cell of `AreaWeeklyMatrixRow.vue`, rendered only when the row
  is a favorite (it does not appear for non-favorites — the table is
  removal-only per AC #4). Click / Enter / Space all remove the favorite.
  `data-cy="matrix-favorite-heart-{itemId}"`. The view passes `areaId`
  down, and the section forwards `itemGroupId` + `itemGroupName` so the
  row can build a complete `ItemFavorite` for `toggleItemFavorite`.
- **Floor plan heart icon.** Added a `$heart` icon inside every
  `.fp-item--free` rectangle (both area-level desk paths and item-level
  free paths), rendered only for favorited items via
  `isFloorPlanItemFavorite(itemId)`. Busy and reserved rectangles never
  show the heart (AC #6). The CSS recipe is exactly the one specified in
  the story: `position: absolute; right: 0; bottom: 0;
  transform: translate(50%, 50%)` so the icon's geometric center sits on
  the bottom-right corner of the rectangle. `.fp-item--free` got
  `overflow: visible` so the heart can extend outside (the rest of the
  `.fp-item` keeps `overflow: hidden` so labels still get clipped on busy
  rectangles). The heart's `@click.stop` is essential — without it the
  click would fall through to `requestBooking` / `handleDeskClick`.
- **Floor plan needs `areaId`.** Added an optional `areaId` prop to
  `InteractiveFloorPlan.vue` and wired both call sites
  (`ItemGroupsView.vue` passes `route.params.areaId`, `ItemsView.vue`
  passes `getCurrentAreaId()`) so favorite matching is correctly scoped.
- **i18n.** Added `favorites.areaName`, `favorites.areaSubtitle`,
  `favorites.emptyTitle`, `favorites.emptyMessage`, and
  `favorites.removeTooltip` to all five locales (en, de, es, fr, uk).
  Translations are idiomatic best-effort; flag for native-speaker review if
  desired.

### Round 2 — UX rework after first review

**Problem reported on 2026-05-11.** The first implementation built a
dedicated `FavoritesView.vue` whose cards looked like room/area tiles
(small icon + name + a row of green dots labelled MO–FR + a tiny `SELECT`
button + a heart). Three things were wrong:

1. The green dots referred to the current week without any selector, so
   users could not tell which week they referenced.
2. The card layout looked like a *room* card, not a *desk* card, so the
   page looked nothing like the room view a user expected.
3. The page lacked all the booking affordances a real room view has
   (day/week toggle, date picker, equipment filter, `BOOK` button,
   `Book for myself / Book for colleague` admin selector).

The story's original brief was *"It should behave like all other area
aka rooms"*. The first implementation was a half-measure.

**Rework approach.** Drop `FavoritesView.vue` and mount `ItemsView.vue`
on the `/favorites` route. ItemsView is the canonical day/week booking
surface; reusing it (rather than duplicating its template into
FavoritesView) keeps the two views visually and behaviourally identical
forever. ItemsView gains a small `favoritesMode` switch driven by
`route.meta.favoritesMode === true` which:

- Skips the single-item-group lookup that drives breadcrumbs, floor plan,
  and the "VIEW ITEM GROUP BOOKINGS" link (those gates were already there,
  keyed on `activeItemGroupId` / `itemGroupFloorPlan`).
- Aggregates `fetchItems` across every distinct `itemGroupId` in the user's
  favorites list, then filters to the favourited `itemId`s, preserving the
  user's favourite order across the merged list. Three new loaders:
  `loadFavoriteItems`, `silentReloadFavoriteItems`, `loadFavoriteWeekData`.
- Wraps every dispatch call site (watchers, post-book reload, post-cancel
  reload, live-feed refresh) in mode-aware helpers: `loadItemsForView`,
  `silentReloadItemsForView`, `loadWeekDataForView`. The single-group
  loaders stay unchanged; only the dispatchers branch.
- Renders an empty state with `favorites.emptyTitle` /
  `favorites.emptyMessage` (i18n keys already added in round 1).
- Makes `isItemFav` look up by `itemId` against `favoriteItems` (the
  desk's heart is always filled in favourites mode).
- `handleToggleItemFav` looks up the favourite entry by `itemId` for the
  original `areaId`/`itemGroupId`, calls `toggleItemFavorite`, then
  immediately filters `items.value` and `weekData[date]` so the card
  disappears without a refetch. The snackbar shows
  "removed from favourites".
- Renders the favourites breadcrumb (`Home > Favorites`) via a branch
  in the existing `breadcrumbs` computed.

The `FavoritesView.vue` stub and its test were deleted.

**Net effect:** the `/favorites` page now renders an exact ItemsView
layout — DAY/WEEK toggle, date picker, equipment filter, admin
"Book for myself / Book for colleague" selector, per-desk cards with
availability chip, equipment chips, warning icon, heart, BOOK button —
populated with the user's favourited desks aggregated across all the
rooms they came from. Booking a desk works identically; cancelling
works identically; live updates (from story 31.1) work identically.
The `FLOOR PLAN` button and `VIEW ITEM GROUP BOOKINGS` link are absent
because favourites span rooms (the existing gates take care of this
for free).

**Tests added.** Five new ItemsView favourites-mode tests cover:

1. Breadcrumb is `Home > Favorites` and the single-item-group lookup
   (`fetchAreas` / `fetchItemGroups`) is skipped entirely.
2. Day mode aggregates `fetchItems` per favourited item group and
   merges the results.
3. Items returned by the API that are *not* favourites are filtered
   out client-side.
4. Empty state shows the favourites copy and skips API calls when
   `favoriteItems` is empty.
5. Clicking the heart removes the desk from the visible list and the
   empty state takes over.

All 410 tests pass (`vitest run`). Type-check, lint, and production
build are clean.

### Deferred follow-ups (intentionally not in this story)

- **Cypress E2E for the full Favorites flow.** Unit + component tests
  cover the wiring (useFavorites, AreasView, ItemGroupsView,
  AreaWeeklyMatrixRow, InteractiveFloorPlan, ItemsView favourites
  mode). A cross-view E2E that exercises adding a favorite from
  ItemsView, navigating home, drilling into Favorites (now ItemsView
  in favourites mode), removing via floor-plan/matrix hearts, and
  seeing the tile disappear belongs in a follow-up alongside the
  live-updates E2E deferred from story 31.1.
- **Smoke test in the dev backend.** Blocked at smoke-test time by two
  pre-existing data issues in the local `private/sithub_areas.yaml`
  (one missing closing quote, plus a duplicate `desk28` id that the
  story 30.1 validation now refuses). One of these I fixed (the missing
  quote); the duplicate ID is owned by the user and not in this story's
  scope. The unit/component tests cover the contract.

### File List

Frontend (deleted):

- `web/src/views/FavoritesView.vue` — superseded by ItemsView in
  favourites mode.
- `web/src/views/FavoritesView.test.ts` — superseded by the new
  `ItemsView.test.ts > favorites mode` describe block.

Frontend (modified, round 2):

- `web/src/router/index.ts` — `/favorites` route now mounts
  `ItemsView.vue` with `meta: { favoritesMode: true }`.
- `web/src/views/ItemsView.vue` — new `favoritesMode` computed,
  favourites-mode breadcrumb, three favourites-aggregator loaders, three
  mode-aware dispatch wrappers, `isItemFav` and `handleToggleItemFav`
  updated, empty-state copy and icon switch on `favoritesMode`.
- `web/src/views/ItemsView.test.ts` — five new tests in a
  `favorites mode` describe block.

Frontend (modified, round 1 — unchanged from first pass):

- `web/src/composables/useFavorites.ts` (items-only API + legacy purge +
  `FAVORITES_AREA_ID` constant + `__resetLegacyPurgeForTests` helper)
- `web/src/composables/useFavorites.test.ts` (drop item-group cases, add
  legacy-purge coverage)
- `web/src/views/AreasView.vue` (prepend Favorites tile; new
  `goToFavorites`)
- `web/src/views/AreasView.test.ts` (favorites tile presence/absence
  coverage)
- `web/src/views/ItemGroupsView.vue` (remove favorites UI, simplify
  `sortedItemGroups`, pass `area-id` to floor plan)
- `web/src/views/ItemGroupsView.test.ts` (assert favorites UI is gone)
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue` (pass
  `area-id` through)
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRoomSection.vue`
  (forward `area-id`/`itemGroupId`/`itemGroupName` to row)
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue` (inline
  heart icon in sticky desk-name cell)
- `web/src/components/InteractiveFloorPlan.vue` (new `areaId` prop, heart
  icon inside every `.fp-item--free` rectangle, corner-anchored CSS)
- `web/src/locales/en.json`, `de.json`, `es.json`, `fr.json`, `uk.json`
  (new `favorites.*` keys)
