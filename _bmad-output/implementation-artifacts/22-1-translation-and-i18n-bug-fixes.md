# Story 22.1: Translation and i18n Bug Fixes

Status: done

## Story

As a user with a non-English language selected,
I want all UI labels, error messages, and abbreviations displayed in my chosen language,
so that the app feels fully localized without English fragments leaking through.

## Acceptance Criteria

1. **Given** the UI language is German
   **When** weekday abbreviation dots render on item group tiles (availability indicators)
   **Then** they show MO, DI, MI, DO, FR, SA, SO (not English MO, TU, WE, TH, FR, SA, SU)

2. **Given** the UI language is German
   **When** weekday headers render in week mode (ItemsView folded tiles)
   **Then** they show M, D, M, D, F, S, S (short) or MO, DI, MI, DO, FR, SA, SO (long)

3. **Given** the UI language is German and a day is free in week mode
   **When** the availability label renders below "frei"
   **Then** "n/a" is no longer shown (remove the label entirely — "frei" is sufficient)

4. **Given** a booking limit error occurs and the UI language is German
   **When** the error snackbar displays
   **Then** the message is in German, e.g. "Sie haben das Maximum von 2 aktiven
   Buchungen für 'Tiefgaragenstellplätze' erreicht"

5. **Given** translations are updated for all languages
   **When** a developer inspects locale files
   **Then** new weekday keys exist in en.json, de.json, es.json, fr.json, uk.json

## Tasks / Subtasks

- [ ] Task 1: Add weekday i18n keys to all locale files (AC: 1, 2, 5)
  - [ ] 1.1 Add keys to `web/src/locales/en.json` under a new `weekdays` section:
    `"weekdays": { "mo": "MO", "tu": "TU", "we": "WE", "th": "TH", "fr": "FR",
    "sa": "SA", "su": "SU", "moShort": "M", "tuShort": "T", "weShort": "W",
    "thShort": "T", "frShort": "F", "saShort": "S", "suShort": "S" }`
  - [ ] 1.2 Add German keys to `de.json`:
    `"mo": "MO", "tu": "DI", "we": "MI", "th": "DO", "fr": "FR", "sa": "SA",
    "su": "SO"` and short forms `"moShort": "M", "tuShort": "D", ...`
  - [ ] 1.3 Add Spanish, French, Ukrainian equivalents to `es.json`, `fr.json`, `uk.json`
- [ ] Task 2: Localize weekday labels in useWeekSelector (AC: 2)
  - [ ] 2.1 In `web/src/composables/useWeekSelector.ts`: remove hardcoded
    `WEEKDAY_LABELS` and `WEEKDAY_LABELS_SHORT` arrays (lines 50-51)
  - [ ] 2.2 Change `getWeekdayLabel` to accept a `t` function parameter or create
    a new composable `useLocalizedWeekdays()` that returns localized labels
    from i18n keys. The function must remain usable from non-setup contexts
    (it's called from `ItemsView.vue` and `ItemGroupsView.vue`)
  - [ ] 2.3 Update `useWeekSelector.test.ts` — tests for `getWeekdayLabel`
    (lines 76-98) need to account for the i18n dependency
- [ ] Task 3: Localize availability indicator weekday labels (AC: 1)
  - [ ] 3.1 In `web/src/views/ItemGroupsView.vue` lines 132 and 217:
    `{{ day.weekday }}` renders the backend's English abbreviation directly.
    Map it through a lookup: `day.weekday` (backend: "MO","TU"...) →
    i18n key (e.g., `weekdays.tu` → "DI" in German)
  - [ ] 3.2 Create a helper function `localizeWeekday(backendAbbrev: string): string`
    that maps "MO"→`t('weekdays.mo')`, "TU"→`t('weekdays.tu')`, etc.
    Place in a shared composable or utility
- [ ] Task 4: Remove "n/a" label from week mode (AC: 3)
  - [ ] 4.1 In `web/src/views/ItemsView.vue`: find all locations where
    `$t('items.notAvailable')` renders "n/a" in week mode (folded + expanded tiles)
    and remove or hide the element. Keep `$t('items.free')` ("frei")
  - [ ] 4.2 Alternatively, change the translation value in all locale files from
    "n/a" to empty string, but only if other views don't depend on it
- [ ] Task 5: Localize booking limit error messages (AC: 4)
  - [ ] 5.1 Add i18n key to all locale files:
    `"bookingLimitExceeded": "You have reached the maximum of {count} active
    bookings for {scope}"` (en.json) /
    `"bookingLimitExceeded": "Sie haben das Maximum von {count} aktiven Buchungen
    für {scope} erreicht"` (de.json) / etc.
  - [ ] 5.2 Add `"bookingLimitExceededGlobal"` key for the case without scope:
    `"You have reached the maximum of {count} active bookings"` / German equivalent
  - [ ] 5.3 In `web/src/views/ItemsView.vue`, update `localizeItemsBookingConflict()`
    (line ~1227): instead of passing through the raw backend English string,
    parse the count (regex: `maximum of (\d+)`) and scope (text after `for `)
    from the backend detail, then use `t('items.bookingLimitExceeded', { count, scope })`
- [ ] Task 6: Fix fallback day names (AC: 2)
  - [ ] 6.1 In `web/src/views/ItemsView.vue` line ~1030: the `getFullDayLabel`
    fallback array `['Monday', 'Tuesday', ...]` is hardcoded English.
    Replace with i18n-based fallback or remove the fallback entirely
    (the `Intl.DateTimeFormat` primary path should always work)
- [ ] Task 7: Run tests and lint (AC: 1, 2, 3, 4, 5)
  - [ ] 7.1 Run `npx vitest run` — all tests pass
  - [ ] 7.2 Run `npm run lint` — clean
  - [ ] 7.3 Run `npm run type-check` — clean
  - [ ] 7.4 Run `npm run build` — clean

## Dev Notes

### Architecture: Backend vs Frontend Weekday Labels

The backend availability API (`/api/v1/areas/:id/item-groups/availability`)
returns English weekday abbreviations in the `weekday` field of `DayAvailability`.
See `internal/itemgroups/availability_handler.go` line 204:
`Weekday: weekdayAbbreviation(day.Weekday())`.

Do NOT change the backend — the API should remain locale-agnostic.
The frontend must map backend abbreviations to localized labels.

### Key Files and Line References

| File | What to change |
| --- | --- |
| `web/src/composables/useWeekSelector.ts:50-56` | Hardcoded WEEKDAY_LABELS |
| `web/src/views/ItemGroupsView.vue:132,217` | `day.weekday` rendered directly |
| `web/src/views/ItemsView.vue:454` | `getWeekdayLabel(dayIdx, isMobile)` |
| `web/src/views/ItemsView.vue:517` | `$t('items.notAvailable')` = "n/a" |
| `web/src/views/ItemsView.vue:1030` | Hardcoded English fallback days |
| `web/src/views/ItemsView.vue:1227-1239` | `localizeItemsBookingConflict()` |
| `web/src/locales/*.json` | Add weekday keys, update notAvailable |

### Previous Story Learnings (from 21-7)

- Booking limit errors come as 409 with detail prefix `"booking limit exceeded:"`
- The `localizeItemsBookingConflict` already detects this prefix (line ~1227)
- Error display now uses `v-snackbar` per `.claude/rules/feedback.md`
- The `errorSnackbarMessage` ref is the single point for error display

### Anti-Pattern Prevention

- Do NOT add a locale/language parameter to the backend API
- Do NOT change the backend `weekdayAbbreviation()` function
- Do NOT hardcode translations outside of locale JSON files
- The `getWeekdayLabel` function is used in non-component contexts (tests,
  data-cy attributes) — ensure the refactored version handles this

### Testing Notes

- `useWeekSelector.test.ts` lines 76-98 test `getWeekdayLabel` with hardcoded
  expectations ("MO", "TU", etc.) — these need updating for i18n
- `ItemGroupsView.test.ts` line 102 has a hardcoded weekdays array — update
- Vitest tests use `createTestI18n()` helper — ensure new keys are included

### References

- [Source: _bmad-output/planning-artifacts/epics.md — Story 22.1]
- [Source: private/epic-22.md — bugs section]
- [Source: private/ux-observations.md — weekday labels, "n/a"]
- [Source: .claude/rules/feedback.md — snackbar convention]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

### Review Findings

- [x] [Review][Patch] Main ItemGroupsView availability indicators still render raw backend weekday abbreviations, so non-favorite tiles remain partly untranslated [web/src/views/ItemGroupsView.vue:217]
