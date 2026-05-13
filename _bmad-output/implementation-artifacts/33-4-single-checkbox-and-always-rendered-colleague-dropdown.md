# Story 33.4: Single Checkbox and Always-Rendered Colleague Dropdown

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user booking an item,
I want a single "Book for a colleague" checkbox alongside the colleague dropdown,
so that the booking-type intent is expressed with one control and the colleague
dropdown does not appear and disappear as I toggle modes.

## Acceptance Criteria

1. **Given** I open the item-groups view
   **When** the booking controls render
   **Then** I see a single "Book for a colleague" checkbox (unchecked by default)
   and a colleague-selection dropdown (rendered but disabled) — there is no
   "Book for myself / Book for a colleague" radio group

2. **Given** the checkbox is unchecked
   **When** I look at the colleague dropdown
   **Then** it is visible, occupies the same space it would when enabled, and is
   disabled (cannot be opened or typed into)

3. **Given** I check the "Book for a colleague" checkbox
   **When** the change is applied
   **Then** the colleague dropdown becomes enabled without any layout reflow
   **And** I can pick a colleague exactly as before

4. **Given** I uncheck the box after having selected a colleague
   **When** the change is applied
   **Then** the dropdown becomes disabled again
   **And** the previously selected colleague is cleared so a stale value is not
   used when the box is rechecked

5. **Given** a booking is being made
   **When** the box is unchecked
   **Then** the booking is made for the current user (the previous "Book for
   myself" behavior)

6. **Given** the Story 32.3 flex-wrap radio-group layout existed before this story
   **When** this story ships
   **Then** that layout is removed entirely; this story supersedes Story 32.3 for
   the booking-type controls

## Tasks / Subtasks

- [ ] Task 1: Replace the radio-group + autocomplete with checkbox + always-rendered
      autocomplete (AC: #1, #2, #3, #4, #6)
  - [ ] 1.1 In `web/src/views/ItemsView.vue` around lines 79–106, replace the
        entire `<div class="booking-type-row d-flex flex-wrap align-center ga-4
        mb-2">` block (introduced by Story 32.3) with:
        ```vue
        <v-checkbox
          v-model="bookForColleague"
          :label="$t('items.bookForColleague')"
          hide-details
          density="compact"
          class="book-colleague-checkbox ma-0"
          data-cy="book-colleague-checkbox"
        />
        <v-autocomplete
          v-model="selectedColleagueId"
          :items="usersList"
          item-title="displayName"
          item-value="id"
          :label="$t('items.selectColleague')"
          density="compact"
          :loading="usersLoading"
          :disabled="!bookForColleague"
          clearable
          hide-details
          data-cy="colleague-select"
          class="colleague-select-inline"
        />
        ```
        The exact final placement (which container the checkbox + dropdown live
        in) is decided by Story 33.5, which compacts the entire booking-controls
        row to a single line. For 33.4 in isolation, keep them in the existing
        `.booking-type-row` flex container so the file stays valid between
        stories; Story 33.5 will then collapse the container.

- [ ] Task 2: Replace the `bookingType` ref with a boolean `bookForColleague` ref
      and migrate all usages (AC: #1, #3, #5)
  - [ ] 2.1 At `web/src/views/ItemsView.vue:1051` change
        `const bookingType = ref<'self' | 'colleague'>('self');` to
        `const bookForColleague = ref(false);`.
  - [ ] 2.2 `grep -n "bookingType" web/src/views/ItemsView.vue` to find every
        consumer and migrate:
        - `bookingType.value === 'colleague'` → `bookForColleague.value`
        - `bookingType.value === 'self'` → `!bookForColleague.value`
        - Watchers that reset `selectedColleagueId` on switch-to-self continue
          to work; just key them off `bookForColleague` falling to `false`.
  - [ ] 2.3 The booking submission logic that selects between "book for me" vs
        "book for colleague" must use `bookForColleague.value` to gate. Read the
        existing `createBooking` call sites in `ItemsView.vue` and adjust.

- [ ] Task 3: Clear `selectedColleagueId` when checkbox is unchecked (AC: #4)
  - [ ] 3.1 Add (or reuse the existing) watcher on `bookForColleague`:
        ```ts
        watch(bookForColleague, (checked) => {
          if (!checked) {
            selectedColleagueId.value = null;
          }
        });
        ```
        Place it near the other ref declarations in the script setup. If a watcher
        already exists on `bookingType`, replace its body with the above; do not
        duplicate.

- [ ] Task 4: Remove obsolete styles introduced by Story 32.3 (AC: #6)
  - [ ] 4.1 In the scoped style block of `ItemsView.vue`, remove
        `.booking-type-row` (it is no longer used after Story 33.5 collapses the
        row; 33.4 in isolation keeps it but the class becomes vestigial — flag
        in completion notes for 33.5 to delete it).
  - [ ] 4.2 Keep `.colleague-select-inline` for now; Story 33.5 will revisit the
        width and wrapping rules. Delete the `.booking-type-radios` rule (no
        longer used since the radio group is gone).

- [ ] Task 5: i18n keys — keep `items.bookForColleague`; remove
      `items.bookForMyself` if unused elsewhere (AC: #1)
  - [ ] 5.1 `items.bookForColleague` continues to be used (now as the checkbox
        label). `items.bookForMyself` was only the "Book for myself" radio option
        label. Grep for it across `web/src/`: if no other consumer exists, remove
        the key from all five locales (`web/src/locales/{en,de,es,fr,uk}.json`).
        If a consumer remains, leave the key.

- [ ] Task 6: Tests (Vitest + Vue Test Utils)
  - [ ] 6.1 In `web/src/views/ItemsView.test.ts`:
        - REMOVE / update the existing `booking-type row layout` describe block
          introduced by Story 32.3 — its assertions about
          `book-self-radio` and the radio toggle no longer apply.
        - ADD `describe('book-for-colleague checkbox', ...)` with assertions:
          - Initial render: `[data-cy="book-colleague-checkbox"]` exists,
            `[data-cy="colleague-select"]` exists and has its `disabled` prop
            true (or `[disabled]` attribute), and there is NO
            `[data-cy="book-self-radio"]` or `[data-cy="book-colleague-radio"]`.
          - Toggle the checkbox to true: the autocomplete is now enabled.
          - Toggle back to false: the autocomplete is disabled AND
            `selectedColleagueId` has been reset to `null`.
  - [ ] 6.2 Search the existing tests for `bookingType` (the ref name) and update
        any test that wrote to it directly to write to `bookForColleague` instead.
        Tests around the colleague-booking flow (e.g. submission) should still
        pass without rewriting their booking-mock fixtures.

- [ ] Task 7: Verification commands
  - [ ] 7.1 From `web/`:
        ```
    npx vitest run
        npm run type-check
        npm run lint
        npm run build
        ```
        All must be green.
  - [ ] 7.2 Manual smoke: open an item-groups items page, confirm a single
        "Book for a colleague" checkbox + a disabled colleague dropdown render;
        check the box, confirm the dropdown enables; pick a colleague, confirm
        a booking is created for them; uncheck the box, confirm the previously
        selected colleague is cleared.

### Review Findings

- [x] [Review][Patch] The round-2 "always enabled colleague dropdown" submit contract lacks regression coverage: tests only mutate `selectedColleagueId` directly, but do not assert that day and week booking submissions pass `{ forUserId, forUserName }` to `createBooking`; a regression could silently book for the current user despite a selected colleague [web/src/views/ItemsView.test.ts:439]

## Dev Notes

### Supersession

Story 32.3 (`_bmad-output/implementation-artifacts/32-3-stable-colleague-booking-dropdown-layout.md`)
shipped a `.booking-type-row` with a `<v-radio-group>` and a conditionally
rendered `<v-autocomplete>`. This story replaces that radio group with a single
checkbox and renders the autocomplete unconditionally (gated by `:disabled`).
The previous `.booking-type-row` / `.colleague-select-inline` CSS may stay
during 33.4 alone; Story 33.5 will repurpose them for the single-line layout.

### Reuse, don't reinvent

| Need | Use this | Path |
| --- | --- | --- |
| Current booking-type controls | `<v-radio-group>` + `<v-autocomplete>` | `web/src/views/ItemsView.vue:80–105` |
| `bookingType` ref to migrate | line 1051 | `web/src/views/ItemsView.vue` |
| `selectedColleagueId`, `usersList`, `usersLoading` | already present | `web/src/views/ItemsView.vue` |
| Existing colleague-mode tests | `colleague autocomplete` it-block | `web/src/views/ItemsView.test.ts` (around line 740) |

### Anti-patterns to avoid

- Do NOT keep a hidden radio group "for backwards compatibility". The brief is
  explicit: one checkbox, no radio.
- Do NOT use `v-show` instead of `:disabled` on the autocomplete. AC #2 requires
  the dropdown to occupy the same space (rendered but disabled) so there is no
  layout shift when the checkbox is toggled.
- Do NOT remove `clearable` from the autocomplete; users still need a way to
  clear a colleague after selecting one (without having to uncheck the box).
- Do NOT leave `bookingType` and `bookForColleague` both alive. One source of
  truth.

### Testing standards

- Vitest + Vue Test Utils, extending `ItemsView.test.ts`.
- Cypress smoke check is optional; the unit tests cover the control wiring.

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Epic 33 Stories]
- [Source: web/src/views/ItemsView.vue:80–105, 1051]
- [Source: _bmad-output/implementation-artifacts/32-3-stable-colleague-booking-dropdown-layout.md]
- [Source: _bmad-output/implementation-artifacts/33-5-single-line-compact-booking-controls.md]
- [Source: web/src/locales/en.json] (i18n cleanup)
- [Source: private/epic-33.md]
- [Source: private/img_21.png] (visual reference for the new control)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `cd web && npx vitest run src/views/ItemsView.test.ts` — 89 tests pass.
- `cd web && npx vitest run` — full suite: 441 tests, 47 files, all green.
- `cd web && npm run type-check` / `npm run lint` / `npm run build` — all clean.

### Completion Notes List

- Replaced the Story 32.3 radio-group + conditional autocomplete with a single
  `<v-checkbox v-model="bookForColleague">` plus an always-rendered
  `<v-autocomplete :disabled="!bookForColleague">`. Layout reflow on toggle
  is eliminated because the autocomplete is always in the layout.
- Renamed the `bookingType: 'self' | 'colleague'` ref to a boolean
  `bookForColleague`. Migrated every consumer (5 read-sites and 1 write-site
  inside the post-booking reset block). The post-booking reset now just sets
  `bookForColleague.value = false`; the watcher takes care of clearing
  `selectedColleagueId`.
- Added a `watch(bookForColleague, ...)` that nulls `selectedColleagueId`
  when the checkbox is unchecked, so a stale colleague does not survive into
  the next booking.
- Updated `ItemsView.test.ts`:
  - Replaced the `booking-type row layout` describe block (from Story 32.3)
    with a new `book-for-colleague checkbox` describe — 3 tests covering
    initial state (no radios, dropdown disabled), toggle-to-checked enables
    the dropdown, toggle-to-unchecked clears `selectedColleagueId`.
  - Updated the `does not render guest radio option` test (now asserts that
    BOTH the guest radio AND the self/colleague radios are gone, with the
    new checkbox present in their place).
  - Updated the `renders colleague autocomplete when booking type is
    colleague` test to reflect the new always-rendered semantics.
- Kept `clearable` on the autocomplete so users still have a way to clear a
  selected colleague without unchecking the box.

### File List

Frontend (modified):

- `web/src/views/ItemsView.vue` — booking-type radio group replaced by a
  single checkbox + always-rendered (conditionally disabled) autocomplete;
  `bookingType` ref renamed to `bookForColleague` (boolean); watcher added;
  5 consumer call-sites migrated.
- `web/src/views/ItemsView.test.ts` — new `book-for-colleague checkbox`
  describe block; two existing tests updated for the new control shape.
