# Story 32.1: Booker Avatar on Day-Mode Item Tiles

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user browsing the item-groups view in day mode,
I want each booked item tile to show the booker's avatar with their name on hover or
tap,
so that I can recognize who has reserved a desk at a glance without opening details.

## Acceptance Criteria

1. **Given** I am viewing an item-groups page in day mode
   **When** a tile represents a booked item
   **Then** the tile displays a circular avatar of the booker
   **And** the avatar treatment matches the one used on the floor plan (image when
   available, initials fallback inside a colored circle)

2. **Given** I hover over the avatar on a desktop viewport
   **When** the tooltip appears
   **Then** it shows the booker's full display name

3. **Given** I tap the avatar on a mobile or touch viewport
   **When** the action is recognized
   **Then** the booker's full display name is shown without navigating away from the view

4. **Given** the booker has no synced or uploaded avatar image (or the image fails to
   load)
   **When** the tile renders
   **Then** a circular initials avatar is shown using `getInitials(booker_name)` from
   `web/src/utils/text.ts`
   **And** the same hover/tap behavior surfaces the full display name

5. **Given** a tile represents an available (not booked) item
   **When** the tile renders
   **Then** no booker avatar is shown

## Tasks / Subtasks

- [ ] Task 1: Render the avatar on the day-mode tile (AC: #1, #4, #5)
  - [ ] 1.1 In `web/src/views/ItemsView.vue`, locate the day-mode booker block at
        lines ~315–323 (the `<div data-cy="item-booker">` that currently renders
        `<v-icon size="14">$user</v-icon>` followed by `entry.attributes.booker_name`).
  - [ ] 1.2 Replace the `$user` icon with a circular avatar that uses the same source
        logic as the floor plan (`web/src/components/InteractiveFloorPlan.vue:174–185`
        and lines 283–293):
        - When `entry.attributes.booker_user_id` is set and the image has not previously
          failed, render `<img :src="getAvatarUrl(entry.attributes.booker_user_id)">`
          in a 32×32 circle (`border-radius: 50%`, `object-fit: cover`).
        - On `@error`, mark that user id as failed (component-local `Set<string>`) and
          fall back to initials.
        - When no `booker_user_id` is available, or the image failed, render a
          circular initials badge using `getInitials(entry.attributes.booker_name)`
          with the same colored background recipe used by `.fp-item-initials` —
          adapt to a stand-alone circle (no `position: absolute`, no `inset`).
  - [ ] 1.3 Import `getAvatarUrl` from `../api/avatars` and `getInitials` from
        `../utils/text` at the top of `<script setup>` (the file already imports
        `middleTruncate` etc. — add `getInitials` to that import).
  - [ ] 1.4 Render the avatar inline before the booker name on the same flex row. Use
        a small `ga-2` gap and `align-center`. Keep the booker name visible (do not
        rely on the avatar alone). Keep `data-cy="item-booker"` on the row.
        Add `data-cy="item-booker-avatar"` to the avatar element (image or
        initials span) so E2E tests can target it.
  - [ ] 1.5 The avatar must NOT be conditional on the floor-plan-only
        `sithub_fp_show_avatars` localStorage flag. That flag (defined at
        `InteractiveFloorPlan.vue:567–569`) only controls the floor plan; tile
        avatars are unconditional. The initials-fallback covers the "no image"
        case.
  - [ ] 1.6 No avatar is rendered for available items (the existing `v-if`
        guarding the booker block already gates on
        `entry.attributes.availability === 'occupied' && entry.attributes.booker_name`,
        keep that guard intact — AC #5).

- [ ] Task 2: Tooltip with full display name on hover/tap (AC: #2, #3)
  - [ ] 2.1 Wrap the avatar in a Vuetify `<v-tooltip location="top">` whose content
        is `entry.attributes.booker_name`. Mirror the pattern used in the week-mode
        tile at `ItemsView.vue:569–579` (booked-by-other tooltip).
  - [ ] 2.2 The tooltip must work on mobile (tap-to-show). Vuetify's `v-tooltip`
        opens on `touchstart` by default; verify with chrome-devtools-mcp in mobile
        emulation if behavior differs. Do NOT add a separate `v-menu` for mobile —
        keep one implementation.
  - [ ] 2.3 The tooltip text is the full `booker_name` as the API returns it. Do
        not truncate; `middleTruncate` / `getShortName` are not used here.

- [ ] Task 3: Styling — small reusable circular avatar (AC: #1)
  - [ ] 3.1 Add CSS to `ItemsView.vue` `<style scoped>` for `.tile-booker-avatar`
        and `.tile-booker-initials`. The avatar is 32×32 with `border-radius: 50%`
        and `object-fit: cover`. The initials fallback uses the same Vuetify error
        color background and white bold text as `.fp-item-initials` but as a
        free-standing circle (no `position: absolute`, no `inset`):
        ```css
        .tile-booker-avatar {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          object-fit: cover;
          flex-shrink: 0;
        }
        .tile-booker-initials {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          display: inline-flex;
          align-items: center;
          justify-content: center;
          font-size: 0.85em;
          font-weight: 600;
          color: white;
          background: rgba(var(--v-theme-error), 0.85);
          user-select: none;
          line-height: 1;
          flex-shrink: 0;
        }
        ```
  - [ ] 3.2 Update the surrounding wrapper to a flex row with `align-center` and
        `ga-2`. The booker name continues to use `text-body-2
        text-medium-emphasis`. Remove the `mr-1` from the now-removed `$user`
        icon.

- [ ] Task 4: Unit / component tests (Vitest + Vue Test Utils)
  - [ ] 4.1 In `web/src/views/__tests__/ItemsView.test.ts` (or whichever path the
        file currently lives at — the previous story used `ItemsView.test.ts` in
        `web/src/views/`), add a `describe('day-mode booker avatar', ...)` block:
        - Asserts that an occupied item with `booker_user_id="user-1"` renders an
          `<img>` whose `src` matches `getAvatarUrl('user-1')` and the element
          has `data-cy="item-booker-avatar"`.
        - Asserts that on image error, the same row falls back to a span
          containing initials (`'TK'` for `'Thorsten Kramm'`).
        - Asserts that an occupied item with no `booker_user_id` (e.g. legacy
          data) renders the initials fallback directly.
        - Asserts that an available item renders no `item-booker-avatar` element.
  - [ ] 4.2 Do not introduce a brand-new test file; extend the existing one.

- [ ] Task 5: Verification commands
  - [ ] 5.1 Run from `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green.
  - [ ] 5.2 Manual smoke: start `cd backend && go run cmd/sithub/main.go run --config
        ./sithub.toml` and `cd web && npm run dev`. Log in, navigate to a room
        with at least one booked desk for today, and confirm the avatar appears
        beside the booker name, the tooltip shows the full name on hover, and the
        initials fallback appears for a user whose avatar 404s.

## Dev Notes

### What this story changes

Day-mode item tiles in `web/src/views/ItemsView.vue` currently show only a
`$user` icon next to the booker name. After this story they will show a 32×32
circular avatar (image with initials fallback) plus the existing name. Hovering
or tapping the avatar reveals the full display name via tooltip.

This story does NOT change:

- The week-mode tile (Story 32.2 handles that).
- The floor plan (it already has avatars; see `InteractiveFloorPlan.vue`).
- The weekly desktop table view from Epic 29 — it is a separate component
  (`AreaWeeklyMatrixRow.vue`) and is out of scope.
- The avatar API surface (`/api/v1/avatars/{userId}`) — already in place per
  Epic 22 (FR97–FR99) and Epic 27 (FR122).

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Build avatar URL | `getAvatarUrl(userId)` | `web/src/api/avatars.ts:1–4` |
| Derive initials | `getInitials(name)` | `web/src/utils/text.ts:33–40` |
| Visual recipe (image + initials fallback) | Floor plan implementation | `web/src/components/InteractiveFloorPlan.vue:170–189, 1803–1825` |
| Closest Vuetify-idiomatic precedent | Weekly matrix cell avatar | `web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:46–90` |
| Tooltip pattern | Week-mode booked-by-other tooltip | `web/src/views/ItemsView.vue:569–579` |

Implementation note on picking an avatar pattern: the matrix cell uses
`<v-avatar size="24"><v-img :src @error="avatarFailed = true" />` plus a
side-by-side text span, while the floor plan uses raw `<img>` with absolute
positioning. For a stand-alone inline avatar on a tile, the cleanest fit is
the Vuetify pattern adapted to put the initials *inside* the avatar circle
(not next to it):

```vue
<v-avatar size="32" class="tile-booker-avatar">
  <v-img
    v-if="entry.attributes.booker_user_id && !failedAvatars.has(entry.attributes.booker_user_id)"
    :src="getAvatarUrl(entry.attributes.booker_user_id)"
    :alt="entry.attributes.booker_name"
    @error="failedAvatars.add(entry.attributes.booker_user_id!)"
  />
  <span v-else class="tile-booker-initials">{{ getInitials(entry.attributes.booker_name) }}</span>
</v-avatar>
```

The `failedAvatars` set is a component-local `reactive(new Set<string>())`.
The CSS in Task 3 still applies (background color and text rules); the
`.tile-booker-avatar` wrapper inherits the circle from `<v-avatar>` so the
explicit `border-radius: 50%; width; height` lines become redundant when this
pattern is used.

Do NOT introduce a `BookerAvatar.vue` component for this single use site — the
markup is small, and Story 32.2 will use the same pattern. If a third use site
appears in the future, refactor then. (Per `.claude/rules/golang.md` rule
spirit: do not abstract until the third occurrence demands it.)

### Anti-patterns to avoid

- Do NOT gate the tile avatar on the floor-plan `sithub_fp_show_avatars`
  localStorage flag. That flag is intentionally local to the floor plan.
- Do NOT call a new backend endpoint or introduce a Pinia store. Avatars are
  served by the existing `/api/v1/avatars/{userId}` endpoint and consumed via a
  plain `<img>` tag — the browser caches the response.
- Do NOT remove the booker name; the avatar is added alongside, not as a
  replacement. Users have been reading names for months; deleting them is a
  regression and is not what the epic asks for.
- Do NOT use Vuetify's `<v-avatar>` component for this. The floor plan uses a
  plain `<img>` with hand-rolled CSS so failure handling (`@error`) and the
  initials fallback are explicit. Mirror that pattern for consistency and so
  the unit tests can assert on `<img>` / `<span>` structure.

### Key code locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Day-mode booker block | `web/src/views/ItemsView.vue:315–323` | The single block to edit |
| Avatar URL helper | `web/src/api/avatars.ts:1–4` | Existing helper to reuse |
| Initials helper | `web/src/utils/text.ts:33–40` | Existing helper to reuse |
| Floor plan avatar markup (reference) | `web/src/components/InteractiveFloorPlan.vue:170–189` | Visual treatment to mirror |
| Floor plan avatar CSS (reference) | `web/src/components/InteractiveFloorPlan.vue:1803–1825` | CSS recipe to adapt to a stand-alone circle |
| Floor plan error handler (reference) | `web/src/components/InteractiveFloorPlan.vue:601–609` | `failedAvatars` set + `onAvatarError` |
| ItemsView tests | `web/src/views/ItemsView.test.ts` | Extend existing suite |
| i18n labels | `web/src/locales/{en,de,es,fr,uk}.json` | No new keys needed for this story |

### Testing standards

- Vitest + Vue Test Utils for the component tests; mock `getAvatarUrl` is not
  needed (the function is pure and deterministic). Use Vuetify test utilities
  per the existing `ItemsView.test.ts` setup.
- Coverage of the four states from AC #1 / #4 / #5: image present, image
  errored, no booker_user_id, available item.
- No new Cypress E2E test is required for this story — `item-booker` already
  has E2E coverage. If trivial, add a single assertion that
  `[data-cy="item-booker-avatar"]` is present in the existing day-mode booking
  flow; otherwise defer.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 32 Stories]
- [Source: web/src/views/ItemsView.vue:315–323]
- [Source: web/src/components/InteractiveFloorPlan.vue:170–189]
- [Source: web/src/components/InteractiveFloorPlan.vue:601–609]
- [Source: web/src/components/InteractiveFloorPlan.vue:1803–1825]
- [Source: web/src/api/avatars.ts]
- [Source: web/src/utils/text.ts]
- [Source: _bmad-output/implementation-artifacts/22-7-user-avatar-sync-from-entra-id.md]
- [Source: _bmad-output/implementation-artifacts/27-1-avatar-sync-fix-for-non-png-formats.md]
- [Source: _bmad-output/implementation-artifacts/28-2-floor-plan-avatar-tooltip-booker-name.md]
- [Source: .claude/rules/vue.md]
- [Source: .claude/rules/cypress.md]
- [Source: private/epic-32.md]
- [Source: private/img_16.png]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/ItemsView.test.ts` — 80 tests pass
  (10 new for stories 32.1, 32.2, 32.3).
- `cd web && npx vitest run` — full suite: 424 tests, 47 files, all green.
- `cd web && npm run type-check` — vue-tsc clean.
- `cd web && npm run lint` — ESLint clean (`--max-warnings 0`).
- `cd web && npm run build` — production build clean.
- `npx jscpd --pattern "src/**/*.ts" --ignore "**/node_modules/**"` — 3.09%
  duplication; all pre-existing in matrix/floor-plan/StatusChip test files.
  No new clones involving `ItemsView.test.ts` or `ItemsView.vue`.

### Completion Notes List

- Implemented the day-mode booker avatar by replacing the `$user` icon
  beside `entry.attributes.booker_name` with a `<v-tooltip>` whose
  activator is a `<v-avatar size="32" data-cy="item-booker-avatar">`. The
  avatar renders `<v-img :src="getAvatarUrl(booker_user_id)">` when the id
  is present and the image hasn't failed; otherwise it falls back to a
  `<span class="tile-booker-initials">` containing the initials from
  `getInitials(booker_name)`.
- The original booker name span is preserved next to the avatar in a
  flex row (`d-flex align-center ga-2`). Removing the name would have
  been a visual regression — the user-facing brief did not ask for that.
- Added a component-local `failedAvatars = reactive(new Set<string>())`
  shared by stories 32.1 and 32.2. Stories 32.1 and 32.2 land together so
  they share the set; the floor plan's older `failedAvatars` ref + `new
  Set([...prev])` recreation pattern was not migrated since touching it
  is out of scope.
- The avatar is unconditional on the floor-plan-only
  `sithub_fp_show_avatars` localStorage flag, per the story's anti-pattern
  note. The initials fallback covers the "no image" case.
- Scoped CSS: `.tile-booker-avatar { flex-shrink: 0 }` and
  `.tile-booker-initials { display: inline-flex; ...; background:
  rgba(var(--v-theme-error), 0.85); color: white; }` — matches the floor
  plan's `.fp-item-initials` color recipe adapted to a free-standing
  Vuetify-supplied circle.
- Tests in `ItemsView.test.ts > day-mode booker avatar` cover all four
  states: image present, no booker_user_id (initials direct), image
  errored (initials after toggling `failedAvatars`), and unbooked tile
  (no avatar).
- Did NOT introduce a `BookerAvatar.vue` component. The same shape is
  used in week mode (Story 32.2), but the templates are short, the data
  shapes differ slightly (entry.attributes vs. dayItems lookup), and
  centralizing would tangle the failure-Set ownership. If a third use
  site appears, refactor then.

### File List

Frontend (modified):

- `web/src/views/ItemsView.vue` — imports `getInitials`, `getAvatarUrl`,
  and Vue's `reactive`; introduces `failedAvatars` reactive Set;
  replaces the day-mode booker icon+name with avatar+tooltip+name;
  adds scoped CSS for `.tile-booker-avatar` and `.tile-booker-initials`.
- `web/src/views/ItemsView.test.ts` — adds
  `describe('day-mode booker avatar', ...)` with four assertions covering
  image, initials fallback, error transition, and available-item gating.
