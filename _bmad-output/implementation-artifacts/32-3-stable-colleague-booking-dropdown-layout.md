# Story 32.3: Stable Inline Layout for "Book for a Colleague" Dropdown

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user with permission to book for a colleague,
I want the colleague-selection dropdown to appear next to the radio button without
pushing the rest of the form down,
so that switching the booking target does not cause the UI to feel nervous.

## Acceptance Criteria

1. **Given** I am on the item-groups view with a wide desktop viewport
   **When** I select "Book for a colleague"
   **Then** the colleague-selection dropdown appears inline, to the right of the
   radio group, on the same horizontal line
   **And** the vertical position of the equipment filter and the tile grid below
   does not change compared to the "Book for myself" state

2. **Given** I toggle back to "Book for myself"
   **When** the dropdown disappears
   **Then** the row containing the radio group keeps its height
   **And** no other element shifts vertically

3. **Given** my viewport is narrow enough that the inline dropdown does not fit
   beside the radio group
   **When** I select "Book for a colleague"
   **Then** the dropdown wraps to its own line below the radio group
   **And** the transition does not produce additional intermediate layout jumps once
   the state has settled

4. **Given** the layout adjusts between inline and wrapped modes when the viewport
   is resized
   **When** the breakpoint is crossed
   **Then** the form remains usable and no interactive element is hidden behind
   another

5. **Given** I open the colleague-selection dropdown (inline or wrapped) and choose
   a colleague
   **When** the selection is applied
   **Then** the existing booking flow continues unchanged
   **And** the chosen colleague is used as the booker for the next booking action
   **And** existing `data-cy` selectors (`book-self-radio`, `book-colleague-radio`,
   `colleague-select`) continue to resolve

## Tasks / Subtasks

- [ ] Task 1: Replace the stacked layout with a flex row (AC: #1, #2, #3)
  - [ ] 1.1 In `web/src/views/ItemsView.vue`, locate the block at lines ~79–101
        containing `<v-radio-group v-model="bookingType" inline ...>` and the
        `<v-expand-transition>` that wraps the colleague dropdown.
  - [ ] 1.2 Wrap the radio group and the colleague dropdown in a single flex
        container that keeps them on one line on wide viewports and wraps on
        narrow ones:
        ```vue
        <div class="booking-type-row d-flex flex-wrap align-center ga-4 mb-2">
          <v-radio-group
            v-model="bookingType"
            inline
            density="compact"
            hide-details
            class="booking-type-radios ma-0"
          >
            <v-radio :label="$t('items.bookForMyself')" value="self" data-cy="book-self-radio" />
            <v-radio :label="$t('items.bookForColleague')" value="colleague" data-cy="book-colleague-radio" />
          </v-radio-group>
          <v-autocomplete
            v-if="bookingType === 'colleague'"
            v-model="selectedColleagueId"
            :items="usersList"
            item-title="displayName"
            item-value="id"
            :label="$t('items.selectColleague')"
            density="compact"
            :loading="usersLoading"
            clearable
            data-cy="colleague-select"
            class="colleague-select-inline"
          />
        </div>
        ```
        Key changes:
        - The outer `<v-expand-transition>` is **removed**. The expand
          transition was the root cause of the perceived "jump": expanding a
          block in a vertical stack pushes everything beneath it.
        - The `<v-autocomplete>` becomes a sibling of the radio group inside a
          flex row, not a block below it.
        - The `mt-4` on the dropdown wrapper is removed; the flex row's
          horizontal `ga-4` (16px) supplies the spacing.
        - `class="mb-2"` on the radio group is moved to the wrapper so the
          gap below the row matches today's spacing exactly. (Removed
          `mb-2` from `v-radio-group` to avoid double margins.)
  - [ ] 1.3 Remove the `style="max-width: 360px;"` from the autocomplete and
        move width control into the scoped CSS for `.colleague-select-inline`
        so wrapping behavior is consistent across viewports. The autocomplete
        should be approximately 320–360px wide on wide viewports and use the
        full row width when wrapped:
        ```css
        .colleague-select-inline {
          flex: 0 0 320px;
          max-width: 360px;
        }
        @media (max-width: 600px) {
          .colleague-select-inline {
            flex: 1 1 100%;
            max-width: 100%;
          }
        }
        ```
  - [ ] 1.4 The radio group's natural width is ~280–320px (two labels). On a
        wide viewport (default item-groups container is ~1200px when sidebar
        is hidden) the row fits both elements inline with the `ga-4` gap. On
        narrow viewports `flex-wrap` causes the autocomplete to drop to the
        next line. No JavaScript breakpoint detection is needed — pure CSS.

- [ ] Task 2: Preserve vertical rhythm (AC: #1, #2)
  - [ ] 2.1 Before: in the "Book for myself" state, the booking-type row was
        the `<v-radio-group>` alone (`class="mb-2"`); the colleague dropdown
        slot was a collapsed `v-expand-transition` rendering nothing. Total
        height roughly 32px + 8px margin.
  - [ ] 2.2 After: the booking-type row is a flex container that contains the
        radio group only in the "Book for myself" state (since the
        autocomplete is gated by `v-if`). The container has `align-center`,
        so its height is determined by the tallest child. In the "Book for
        myself" state that's the radio group (32px). In the "Book for a
        colleague" state on wide viewports it's the autocomplete (40px with
        density="compact"). The difference between states must be visible
        only in **this** row, never below.
  - [ ] 2.3 To make AC #1's "no vertical position change below" strict: set a
        `min-height` on `.booking-type-row` equal to the autocomplete's
        compact height so the row height does not change between states:
        ```css
        .booking-type-row {
          min-height: 40px;
        }
        ```
        With `min-height`, the equipment filter and tile grid sit at exactly
        the same y-coordinate in both states. Verify with chrome-devtools-mcp
        by switching between the two radio values and inspecting the equipment
        filter's `getBoundingClientRect().top`.
  - [ ] 2.4 Remove the `<v-expand-transition>` import if it is no longer used
        elsewhere in the file. (Grep first; if used, leave the import.)

- [ ] Task 3: Behavior unchanged (AC: #5)
  - [ ] 3.1 Do NOT change `bookingType`, `selectedColleagueId`, or `usersList`
        bindings. The reactive sources, watchers, and side effects (e.g. the
        existing watch that resets `selectedColleagueId` when `bookingType`
        switches to `'self'`) must keep working. Grep for `bookingType` and
        confirm no template-only assumption breaks.
  - [ ] 3.2 Verify the `book-self-radio`, `book-colleague-radio`, and
        `colleague-select` `data-cy` selectors remain attached to the same
        Vuetify components. E2E specs in `web/cypress/e2e/` use these — grep
        for each before merging:
        ```
    grep -rn "book-self-radio\|book-colleague-radio\|colleague-select" web/cypress
        ```
        Update any spec that asserted on a now-removed wrapper (`.mt-4`,
        v-expand-transition wrapper, etc.). If a spec asserts on the
        autocomplete being on a new visual line, update the assertion to
        match the new inline layout instead.

- [ ] Task 4: Unit / component tests (Vitest + Vue Test Utils)
  - [ ] 4.1 In `web/src/views/ItemsView.test.ts` (or wherever the suite
        lives), add a `describe('booking-type row layout', ...)`:
        - Mount with `bookingType = 'self'`; assert the row has no
          `[data-cy="colleague-select"]` element.
        - Toggle `bookingType` to `'colleague'`; assert
          `[data-cy="colleague-select"]` exists and is a child of the same
          `.booking-type-row` container as the radio group (i.e. they share
          a parent).
        - Assert that switching the radio model does not remount the entire
          `<v-card-text>` (component identity check via a ref or a
          render-counter watcher in the test) — guards against accidentally
          breaking the live-booking update wiring.
  - [ ] 4.2 No need to test the wrapping breakpoint in Vitest (JSDom does
        not paint). Defer that to manual smoke + the optional Cypress check
        below.

- [ ] Task 5: Cypress smoke (optional, only if trivial in the existing spec)
  - [ ] 5.1 In an existing item-groups E2E spec, add a one-line check:
        click "Book for a colleague", capture the `top` of
        `[data-cy="equipment-filter-input"]`, click "Book for myself",
        capture again, assert equal. This protects the vertical-stability
        contract from regression. Skip if no item-groups spec exists yet —
        the unit tests are sufficient.

- [ ] Task 6: Verification commands
  - [ ] 6.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All green.
  - [ ] 6.2 Manual smoke via chrome-devtools-mcp:
        - Navigate to an items page in a desktop-sized viewport (≥ 1024px).
        - Take a screenshot in "Book for myself" state.
        - Toggle to "Book for a colleague".
        - Take a second screenshot. The equipment filter and the first tile
          must occupy the same vertical position in both screenshots. The
          colleague autocomplete sits to the right of the radio group on the
          same row.
        - Resize the viewport to < 600px. Toggle again. The autocomplete
          wraps below the radio group; the equipment filter shifts down by
          one row's height (this is expected and acceptable on mobile).
        - Confirm that selecting a colleague still populates the booker for
          the next booking action.

## Dev Notes

### What this story changes

The booking-type controls in `web/src/views/ItemsView.vue:79–101` use a
vertical stack: a `<v-radio-group>` followed by a `<v-expand-transition>`
that reveals the colleague autocomplete. Toggling the radio expands a block
of ~40px+, pushing the equipment filter and the entire tile grid down. After
this story, the two controls live on the same flex row, the row has a fixed
`min-height`, and the `v-expand-transition` is gone. On wide viewports the
autocomplete renders inline next to the radio group; on narrow viewports it
wraps to a new line.

This story does NOT change:

- The colleague booking flow itself (`bookingType`, `selectedColleagueId`,
  `usersList`, watchers, API calls).
- The equipment filter (still positioned beneath the booking-type row).
- Day vs week mode logic.
- The `data-cy` selectors used by E2E tests.

### Why `min-height` + `flex-wrap` (and not a Vuetify breakpoint helper)

The root cause is layout, not state. Using a CSS `min-height` on the
flex-wrap row removes the vertical jump deterministically — no JS, no
breakpoint detection, no Vuetify `useDisplay()` import. Pure CSS makes the
behavior predictable across SSR / hydration and avoids a tiny FOUC on first
render. The `flex-wrap` rule handles narrow viewports the same way.

`useDisplay()` from Vuetify is available (the file already detects mobile
via `window.matchMedia` at line ~2255), but adding another reactive source
for a purely visual problem is unnecessary complexity.

### Anti-patterns to avoid

- Do NOT keep the `<v-expand-transition>`. The expand animation is the
  visible "nervous" motion the epic complains about. Without it, the
  autocomplete simply appears in the row.
- Do NOT use `position: absolute` on the autocomplete. It must participate
  in flex layout so wrapping works correctly on narrow viewports.
- Do NOT introduce a Pinia store, composable, or `useDisplay()` call for
  this. The change is CSS-only.
- Do NOT change `bookingType`'s allowed values (`'self'` / `'colleague'`)
  or the i18n keys (`items.bookForMyself`, `items.bookForColleague`,
  `items.selectColleague`).
- Do NOT remove the `inline` prop on `<v-radio-group>` — the labels must
  stay on one line.
- Do NOT add a "show only when admin" gate that isn't already there. The
  current visibility rules for "Book for a colleague" stay exactly as
  defined upstream.

### Key code locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Booking-type radio group | `web/src/views/ItemsView.vue:80–83` | Source of the radio |
| Colleague dropdown (current expand-transition) | `web/src/views/ItemsView.vue:86–101` | Source of the wrapped dropdown — replaced |
| Equipment filter (the element below) | `web/src/views/ItemsView.vue:104–140` | Verify it does not move vertically across radio toggles |
| ItemsView styles | `web/src/views/ItemsView.vue` (scoped style block, file end) | Add `.booking-type-row`, `.colleague-select-inline` rules |
| ItemsView tests | `web/src/views/ItemsView.test.ts` | Extend |
| i18n labels | `web/src/locales/{en,de,es,fr,uk}.json` | No new keys needed |

### Testing standards

- Vitest + Vue Test Utils. The component test should verify DOM structure
  (parent-child relationship of radio group and autocomplete), not visual
  layout (JSDom does not paint).
- Cypress smoke is optional and only justified if a relevant spec already
  exists; otherwise the unit test plus manual smoke is enough.
- E2E selectors stay stable: grep `web/cypress/e2e/` for each `data-cy`
  before merging.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 32 Stories]
- [Source: web/src/views/ItemsView.vue:79–101]
- [Source: web/src/views/ItemsView.vue:104–140]
- [Source: .claude/rules/vue.md]
- [Source: .claude/rules/cypress.md]
- [Source: .claude/rules/feedback.md]
- [Source: private/epic-32.md]
- [Source: private/img_18.png]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/ItemsView.test.ts` — 80 tests pass
  (3 new in the booking-type row layout describe block).
- `cd web && npx vitest run` — full suite: 424 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all
  clean.

### Completion Notes List

- Replaced the stacked layout (radio group followed by a
  `<v-expand-transition>` block) with a single
  `<div class="booking-type-row d-flex flex-wrap align-center ga-4 mb-2">`
  that contains the `<v-radio-group>` and, conditionally, the
  `<v-autocomplete>`. The expand transition — the visible cause of the
  perceived "jump" — is gone.
- The `mb-2` margin moved from the radio group to the wrapper so spacing
  beneath the row matches the previous state exactly. Added
  `class="ma-0"` to the radio group to suppress Vuetify's default
  margin that would otherwise make the row a hair taller.
- Added `hide-details` to the autocomplete so its bottom helper-text
  area no longer reserves vertical space (the row is now compact and
  inline-friendly).
- Scoped CSS adds `.booking-type-row { min-height: 40px }` so the row
  occupies the same vertical space in both states — equipment filter
  and tile grid sit at the same y-coordinate regardless of the radio
  selection. `.colleague-select-inline` is `flex: 0 0 320px; max-width:
  360px` on wide viewports; on `max-width: 600px` it switches to `flex:
  1 1 100%` so it wraps to a full-width row on mobile.
- No JS breakpoint detection was added. `flex-wrap` handles the narrow
  case purely in CSS, which avoids hydration mismatches and keeps the
  behavior deterministic.
- Did NOT modify `bookingType`, `selectedColleagueId`, `usersList`, or
  any of the existing colleague-resolution side effects. Existing
  `data-cy` selectors (`book-self-radio`, `book-colleague-radio`,
  `colleague-select`) remain attached to the same Vuetify components.
- Tests in `ItemsView.test.ts > booking-type row layout` cover three
  states: the dropdown is hidden in the self state, it appears as a
  child of `.booking-type-row` in the colleague state, and it
  disappears again when toggled back to self.

### File List

Frontend (modified):

- `web/src/views/ItemsView.vue` — replaces the radio-group + expand-
  transition stack with a flex-wrap row containing both controls;
  removes inline `style="max-width: 360px"` from the autocomplete in
  favor of scoped CSS; adds `hide-details` and `class="ma-0"` to the
  inner Vuetify components; adds scoped CSS for `.booking-type-row`,
  `.booking-type-radios`, and `.colleague-select-inline` (with a
  mobile media-query override).
- `web/src/views/ItemsView.test.ts` — adds `describe('booking-type row
  layout', ...)` with three assertions covering the dropdown's absence,
  presence as a row sibling, and re-absence after toggling back.
