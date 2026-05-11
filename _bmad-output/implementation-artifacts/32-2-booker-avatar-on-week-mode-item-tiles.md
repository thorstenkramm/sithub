# Story 32.2: Booker Avatar on Week-Mode Item Tiles

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user browsing the item-groups view in week mode,
I want each booked weekday cell on a tile to show the booker's avatar with their name
on hover or tap,
so that I can scan a whole week and immediately see who has booked which day.

## Acceptance Criteria

1. **Given** I switch the item-groups view to week mode
   **When** a weekday cell on a tile represents a booking by another user
   **Then** the cell displays a circular avatar of that day's booker
   **And** the avatar treatment is visually consistent with day mode (Story 32.1) and
   the floor plan

2. **Given** I hover over an avatar in a weekday cell on a desktop viewport
   **When** the tooltip appears
   **Then** it shows the booker's full display name

3. **Given** I tap an avatar in a weekday cell on a touch viewport
   **When** the action is recognized
   **Then** the booker's full display name is shown without navigating away from the
   view

4. **Given** the booker has no synced or uploaded avatar image (or the image fails to
   load)
   **When** the weekday cell renders
   **Then** an initials-based circular avatar is shown using `getInitials(booker_name)`
   from `web/src/utils/text.ts` (replacing the current ad-hoc
   `getBookerInitials` two-char span)

5. **Given** a weekday cell is free
   **When** the tile renders
   **Then** the existing free-state visuals (checkbox, "frei" label) are shown
   unchanged
   **And** no avatar is rendered for that cell

6. **Given** a weekday cell is booked-by-me (a booking owned by me, regardless of
   whether the booker is me or a colleague I booked on behalf of)
   **When** the folded tile renders
   **Then** the cell shows a circular avatar of the booker (image when available,
   initials fallback) with no visible name text beneath
   **And** the booker's full display name is shown on hover (desktop) and tap
   (touch) via the tooltip
   **And** the avatar `src`, `alt`, and initials are all sourced from the API's
   `booker_user_id` / `booker_name` (via `getWeekDayBookerUserId` and
   `getWeekDayBooker`), NOT from `authStore` — so a booking made on behalf of a
   colleague displays the colleague's identity, not the current user's

7. **Given** a weekday cell is booked-by-me
   **When** the expanded tile renders
   **Then** the row shows the avatar AND the booker's display name (sourced from
   `getWeekDayBooker`, NOT from `authStore.userName`)

8. **Given** a weekday cell is booked-by-me and the date is not in the past
   **When** I click the checked checkbox in that cell
   **Then** the existing week-cancel confirmation dialog opens (no separate red
   cancel-X icon is rendered)
   **And** confirming the dialog cancels the booking via the existing
   `requestWeekCancel` / `confirmWeekCancel` flow

9. **Given** a weekday cell is booked-by-me and the date is in the past
   **When** the tile renders
   **Then** the checkbox is disabled and clicking it has no effect

## Tasks / Subtasks

- [ ] Task 1: Compact (folded) week-mode cell — replace initials span with avatar
      (AC: #1, #4, #5, #6)
  - [ ] 1.1 In `web/src/views/ItemsView.vue`, locate the folded week-day cell block
        at lines ~569–579 (the `<template v-else-if="getWeekDayStatus(...) ===
        'booked-by-other'">` containing the `v-tooltip` and the
        `getBookerInitials(...)` span).
  - [ ] 1.2 Replace the current `<span class="week-day-status text-caption
        text-error">{{ getBookerInitials(...) }}</span>` with a small circular
        avatar inside the same tooltip activator:
        ```vue
        <v-avatar
          v-bind="tooltipProps"
          size="24"
          class="week-day-avatar"
          :data-cy="`week-day-avatar-${item.id}-${date}`"
        >
          <v-img
            v-if="getWeekDayBookerUserId(item.id, date) &&
                  !weekAvatarFailed.has(getWeekDayBookerUserId(item.id, date)!)"
            :src="getAvatarUrl(getWeekDayBookerUserId(item.id, date)!)"
            :alt="getWeekDayBooker(item.id, date)"
            @error="weekAvatarFailed.add(getWeekDayBookerUserId(item.id, date)!)"
          />
          <span v-else class="week-day-initials">
            {{ getInitials(getWeekDayBooker(item.id, date)) }}
          </span>
        </v-avatar>
        ```
        The avatar replaces the bare initials span. The tooltip content
        (`{{ getWeekDayBooker(...) }}`) and `v-tooltip location="top"` wrapper
        stay unchanged.
  - [ ] 1.3 Add a new helper near `getWeekDayBooker` (line ~1315):
        ```ts
        const getWeekDayBookerUserId = (
          itemId: string,
          date: string,
        ): string | undefined => {
          const dayItems = weekData.value[date];
          return dayItems?.find(i => i.id === itemId)?.attributes.booker_user_id;
        };
        ```
        The `Item` API type already includes `booker_user_id?: string` (see
        `web/src/api/items.ts:11`), so this is a typed lookup. No backend change.
  - [ ] 1.4 Define `const weekAvatarFailed = reactive(new Set<string>())` near
        the other refs at the top of `<script setup>`. (Use `reactive` not
        `ref(new Set())` so `.add()` is reactive without reassigning. The
        floor plan uses an immutable replacement pattern; the matrix cell uses a
        per-cell local ref — for a list rendered hundreds of times, a single
        shared reactive Set is the right shape.)
  - [ ] 1.5 Replace the existing ad-hoc `getBookerInitials` (lines ~1322–1329)
        with a thin wrapper over the shared `getInitials` utility:
        ```ts
        const getBookerInitials = (itemId: string, date: string): string =>
          getInitials(getWeekDayBooker(itemId, date));
        ```
        Or, since the only remaining caller is the new initials fallback span
        inside `v-avatar`, inline the call to `getInitials(getWeekDayBooker(...))`
        and delete `getBookerInitials` entirely. The two-char-substring fallback
        for single-word names is acceptable to drop because `getInitials` handles
        single-word names with a single uppercase char, which is fine for the
        avatar.
  - [ ] 1.6 Import `getInitials` from `../utils/text` and `getAvatarUrl` from
        `../api/avatars` at the top of `<script setup>`. (32-1 already imports
        these — verify your branch order; if 32-1 lands first, the imports may
        already exist.)
  - [ ] 1.7 Add `import { reactive } from 'vue'` to the Vue imports (the file
        already imports `ref`, `computed`, etc.).

- [ ] Task 2: Expanded week-mode row — same change to the one-line-per-day variant
      (AC: #1, #4)
  - [ ] 2.1 Locate the expanded-week block at lines ~640–643 (the `<span
        v-else-if="getWeekDayStatus(item.id, date) === 'booked-by-other'">`
        currently rendering the full booker name).
  - [ ] 2.2 Replace the bare full-name span with the same `v-tooltip` +
        `v-avatar` pattern used in Task 1, but the avatar sits next to the
        existing full booker name (do not delete the name — the expanded view
        intentionally shows it). Use size 28 here (slightly larger because
        the row has more vertical room) — match the row gap with `class="ml-2"`
        or the existing flex layout the row already uses.
  - [ ] 2.3 Add `data-cy="week-day-avatar-expanded-${item.id}-${date}"` to the
        expanded-avatar element so E2E can distinguish folded vs expanded.

- [ ] Task 3: Styling
  - [ ] 3.1 Add CSS in `<style scoped>`:
        ```css
        .week-day-avatar {
          flex-shrink: 0;
        }
        .week-day-initials {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 100%;
          height: 100%;
          font-size: 0.7rem;
          font-weight: 600;
          color: white;
          background: rgba(var(--v-theme-error), 0.85);
          line-height: 1;
          user-select: none;
        }
        ```
        The `v-avatar` provides the circle and clipping; `.week-day-initials`
        fills it.
  - [ ] 3.2 Ensure the new avatar still fits inside `.week-day-slot` without
        increasing the row height — the existing slot already accommodates a
        checkbox (~24px) and the booker-initials span; size 24 keeps parity.
        Check the mobile (`week-days-compact`) and desktop (`week-days`)
        variants of the row via chrome-devtools-mcp.

- [ ] Task 4: Unit / component tests (Vitest + Vue Test Utils)
  - [ ] 4.1 In `web/src/views/__tests__/ItemsView.test.ts` (or the existing
        path used by the project), extend the suite with
        `describe('week-mode booker avatar', ...)`:
        - Asserts that a `booked-by-other` weekday cell with a `booker_user_id`
          renders `[data-cy^="week-day-avatar-"]` containing an `<img>` whose
          `src` matches `getAvatarUrl(...)`.
        - Asserts that on image load error, the same element falls back to a
          `<span class="week-day-initials">` containing the initials.
        - Asserts that a `booked-by-other` cell with `booker_user_id` absent
          renders the initials fallback directly.
        - Asserts that `free` and `booked-by-me` cells render no avatar
          element.
  - [ ] 4.2 Reuse the existing week-data fixtures in `ItemsView.test.ts`.
        Mock the avatar endpoint by inspecting `<img src=...>` only — no
        network call is needed in JSDom.

- [ ] Task 5: Verification commands
  - [ ] 5.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All green.
  - [ ] 5.2 Manual smoke (chrome-devtools-mcp acceptable): start backend and
        frontend (see `.claude/rules/vite.md`), log in as a non-admin user,
        navigate to a room, switch to week mode, confirm each
        booked-by-other day shows a circular avatar (image or initials), and
        the tooltip surfaces the full booker name on hover (desktop) and tap
        (mobile emulation).

### Review Findings

- [x] [Review][Patch] Expanded week-mode booked-by-other rows render avatar and booker name as separate grid children, so the name is auto-placed into the next grid row instead of staying in the status column with the avatar [web/src/views/ItemsView.vue:678]
- [x] [Review][Patch] Week-mode avatar regression coverage misses required image-error fallback and booked-by-me absence states; the current "free or booked-by-me" test only renders a free cell [web/src/views/ItemsView.test.ts:946]

## Dev Notes

### What this story changes

Week-mode tiles in `web/src/views/ItemsView.vue` currently show a 2-character
initials span (`getBookerInitials`) inside the booked-by-other weekday cell.
After this story they will show a 24px circular avatar instead — image when
available, initials inside the circle when not — with the tooltip already in
place.

The expanded row variant (one-line-per-day) gets the avatar **alongside** the
existing full name (the expanded view exists to show more information; do not
remove the name).

This story does NOT change:

- Day-mode tile rendering — Story 32.1 owns that surface.
- Free / booked-by-me / unavailable cells. Only `booked-by-other` cells gain
  the avatar.
- The weekly desktop matrix from Epic 29 (`AreaWeeklyMatrixCell.vue`) — it
  already has its own avatar implementation.
- The week-data API. `booker_user_id` is already returned in the existing
  `Item` payload (`web/src/api/items.ts:11`).

### Why a shared `reactive(new Set<string>())` for failures

The week view can render dozens of cells across 5–7 days. A per-cell `ref`
(matrix-cell pattern) would mean re-checking the same user's failed avatar
many times if the same colleague has booked multiple days. A single shared
`Set<string>` keyed by user id avoids redundant 404s by remembering failure
across all cells once a load has failed.

Use `reactive(new Set<string>())` — its `.add()` is reactive. Do **not** use
`ref(new Set())` with `.value.add(...)`, which mutates without triggering an
update (Vue tracks `.value` reassignment, not nested mutations on a non-proxy
Set). The floor plan currently re-creates the Set with `new Set([...prev, id])`
to work around this; the `reactive(Set)` form is cleaner.

### Replacement of `getBookerInitials`

The current implementation (`ItemsView.vue:1322–1329`) handles single-word
names by returning the first two characters uppercased. The shared
`getInitials` returns only the first character for single-word names. This
is a deliberate simplification that matches the floor plan's behavior; users
have been seeing `getInitials` output on the floor plan for some time
without complaint. If a future audit shows the single-char output is
ambiguous, that's a separate change to `text.ts` and applies everywhere.

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Build avatar URL | `getAvatarUrl(userId)` | `web/src/api/avatars.ts:1–4` |
| Derive initials | `getInitials(name)` | `web/src/utils/text.ts:33–40` |
| Closest precedent (Vuetify) | Weekly matrix cell | `web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:46–90` |
| Tooltip pattern | Existing folded-week `v-tooltip` | `web/src/views/ItemsView.vue:569–579` |
| Day-mode tile avatar (Story 32.1) | Use the same visual treatment | `_bmad-output/implementation-artifacts/32-1-booker-avatar-on-day-mode-item-tiles.md` |

### Coordination with Story 32.1

If 32-1 lands first, you can reuse:

- The `failedAvatars` reactive Set name (if 32-1 named it identically — verify
  in main, otherwise keep a separate `weekAvatarFailed` to avoid coupling the
  two surfaces).
- The shared CSS classes if 32-1 introduced reusable ones. As of writing,
  32-1 uses `.tile-booker-avatar` / `.tile-booker-initials` — week mode uses
  smaller (24px) variants and is named `.week-day-avatar` /
  `.week-day-initials`. Keep them distinct to avoid cross-impact.

If 32-2 lands first (Story 32.1 not yet merged), the helpers and patterns
introduced here should be reused by 32-1 verbatim; do not duplicate them.

### Anti-patterns to avoid

- Do NOT pass the booker user id through the URL or any new API. The data is
  already in the existing `Item` payload — the helper just reads it.
- Do NOT use `<v-avatar size="32">` here. The week-mode cells are denser; 24
  matches the surrounding checkbox heights.
- Do NOT remove the existing tooltip wrapper or the `data-cy="week-day-..."`
  selectors used by E2E tests in `web/cypress/e2e/`. Grep first if in doubt.
- Do NOT introduce a per-cell composable (`useWeekDayAvatar`) for one shared
  Set and two helper functions — that's premature abstraction.

### Key code locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Folded week-day cell `booked-by-other` | `web/src/views/ItemsView.vue:569–579` | Replace span with avatar |
| Expanded week-day cell `booked-by-other` | `web/src/views/ItemsView.vue:640–643` | Add avatar next to full name |
| Week-data status helpers | `web/src/views/ItemsView.vue:1302–1329` | Add `getWeekDayBookerUserId`; simplify `getBookerInitials` |
| `Item` API type | `web/src/api/items.ts:11` | `booker_user_id?: string` is already present |
| Matrix-cell avatar precedent | `web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:46–90` | Reference Vuetify pattern |
| Day-mode avatar (Story 32.1) | `web/src/views/ItemsView.vue:315–323` (post-32.1) | Keep visual consistency |
| Locales | `web/src/locales/{en,de,es,fr,uk}.json` | No new keys needed |

### Testing standards

- Vitest + Vue Test Utils, extending the existing `ItemsView.test.ts`. No new
  test file.
- Cover the four cell states from AC #1 / #4 / #5 / #6: image present,
  image errored, no `booker_user_id`, free, booked-by-me. (booked-by-me
  intentionally renders no avatar; assert absence.)
- No new Cypress E2E. Existing week-mode specs (`web/cypress/e2e/`) use the
  current `data-cy` selectors; do not remove or rename them.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 32 Stories]
- [Source: web/src/views/ItemsView.vue:569–579]
- [Source: web/src/views/ItemsView.vue:640–643]
- [Source: web/src/views/ItemsView.vue:1302–1329]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue]
- [Source: web/src/api/items.ts]
- [Source: web/src/api/avatars.ts]
- [Source: web/src/utils/text.ts]
- [Source: _bmad-output/implementation-artifacts/32-1-booker-avatar-on-day-mode-item-tiles.md]
- [Source: _bmad-output/implementation-artifacts/28-2-floor-plan-avatar-tooltip-booker-name.md]
- [Source: .claude/rules/vue.md]
- [Source: .claude/rules/cypress.md]
- [Source: private/epic-32.md]
- [Source: private/img_17.png]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/ItemsView.test.ts` — 80 tests pass
  (3 new for the week-mode avatar describe block).
- `cd web && npx vitest run` — full suite: 424 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all
  clean.

### Completion Notes List

- **UX revision (2026-05-11, round 2):** the round-1 implementation used
  `authStore.userName` for the displayed name on booked-by-me cells, which
  surfaced the wrong identity whenever the current user booked on behalf of
  a colleague (booking owned by me, booker is the colleague). The folded
  tile also showed a name text under the checkbox, which the user wanted
  removed in favor of tooltip-only. The user also asked for the red cancel-X
  icon to disappear, with cancellation triggered by clicking the booked
  checkbox itself. Round-2 changes:
    1. Booked-by-me and booked-by-other now share a single avatar template
       in both folded and expanded variants. All sources of identity
       (`getAvatarUrl(booker_user_id)`, `<v-img :alt>`, `getInitials`, and
       the displayed name in the expanded row) come from `getWeekDayBooker`
       and `getWeekDayBookerUserId` (i.e. the API booking record), never
       from `authStore`. This fixes the wrong-name bug.
    2. The folded variant no longer renders any name text beneath the
       avatar; the booker's name appears only in the tooltip on hover/tap,
       matching `private/img_17.png`.
    3. The expanded variant continues to show the booker name next to the
       avatar (now sourced correctly), with text color `text-primary` for
       me-bookings and `text-error` for others' bookings.
    4. The red `week-cancel-btn` v-icon was removed from both variants.
    5. The booked-by-me v-checkbox lost its hard-coded `disabled` attribute
       and is now interactive when the date is not in the past. Its
       `@update:model-value` handler calls `requestWeekCancel`, which opens
       the existing confirmation dialog. Past-date booked-by-me cells keep
       the checkbox `disabled` since cancellation of past bookings is not
       supported.
- **Round-1 history (2026-05-11):** the initial implementation followed
  the original AC #6 wording, which excluded the `booked-by-me` state from
  the avatar treatment entirely. That was wrong; round-1 added the avatar
  to booked-by-me but still used `authStore` for the name. Round-2 fixed
  both that name-source bug and the UX gap above.
- Replaced the folded week-day `booked-by-other` initials span (formerly
  rendered via the ad-hoc `getBookerInitials` two-character substring
  helper) with a `<v-avatar size="24" data-cy="week-day-avatar-${itemId}-
  ${date}">`. The avatar renders `<v-img :src="getAvatarUrl(booker_user_id)">`
  when the id is present and the image has not failed, otherwise it falls
  back to a `<span class="week-day-initials">` populated by
  `getInitials(getWeekDayBooker(...))`. The existing `<v-tooltip>` wrapper
  and the full-name body text remain unchanged.
- Replaced the expanded-row `booked-by-other` block analogously, using
  `size="28"` and `data-cy="week-day-avatar-expanded-..."`. The full
  booker name span remains (the expanded view exists to show more info);
  the avatar sits to its left.
- Added a helper `getWeekDayBookerUserId(itemId, date)` next to
  `getWeekDayBooker` that reads `booker_user_id` from `weekData.value`.
  The `Item` API type already exposes `booker_user_id` (see
  `web/src/api/items.ts:11`).
- Removed the now-unused `getBookerInitials` ad-hoc helper. The two-char
  substring fallback for single-word names is replaced by `getInitials`'s
  single-char behavior, which matches the floor plan and the matrix cell.
- The `failedAvatars` reactive Set introduced in Story 32.1 is shared
  here. Tests confirm both day-mode and week-mode failures coexist
  without cross-talk.
- Scoped CSS: `.week-day-avatar { flex-shrink: 0 }` and
  `.week-day-initials` mirror the day-mode recipe at a smaller font size
  (0.7rem) to fit the 24px avatar.
- Tests in `ItemsView.test.ts > week mode rendering > booker avatar on
  booked-by-other cells` cover three states: image present, no
  booker_user_id, and absence on free / booked-by-me cells. Uses
  `vi.setSystemTime` to pin the week so `weekData` materializes
  deterministically.

### File List

Frontend (modified):

- `web/src/views/ItemsView.vue` — folded and expanded week-day
  `booked-by-other` blocks now render a `<v-avatar>` with image +
  initials fallback inside the existing `<v-tooltip>`. Adds
  `getWeekDayBookerUserId` helper. Removes the ad-hoc
  `getBookerInitials` (replaced by direct `getInitials(...)`). Adds
  scoped CSS for `.week-day-avatar` and `.week-day-initials`.
- `web/src/views/ItemsView.test.ts` — assertions covering: image render,
  initials fallback, free-state absence, booked-by-me image render,
  booker-name-not-authStore-name (book-on-behalf scenario), folded
  no-name-text invariant, no-red-cancel-icon invariant, and the
  checkbox-click → cancel-confirmation-dialog wiring.
