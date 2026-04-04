# Story 21.2: Language Selector with Flags and Auto-Detection

Status: done

## Story

As a user,
I want to choose my preferred UI language from the settings page,
so that I can use SitHub in my native language.

## Acceptance Criteria

1. **Given** I open the user menu, **when** I see the language selector,
   **then** it shows options: Auto, English, Deutsch, Espanol, Francais, Ukrainska —
   each with a country flag emoji (GB for English, DE for German, ES for Spanish,
   FR for French, UA for Ukrainian).

2. **Given** I select "Deutsch", **when** the selection is applied,
   **then** the entire UI switches to German immediately without page reload.

3. **Given** I select "Auto", **when** my browser's preferred language is German,
   **then** the UI renders in German.

4. **Given** I select "Auto", **when** my browser's preferred language is not one
   of the supported languages, **then** the UI falls back to English.

5. **Given** I select a language and close the browser, **when** I reopen SitHub,
   **then** the previously selected language is restored from local storage.

## Tasks / Subtasks

- [x] Task 1: Create useLocalePreference composable (AC: 2, 3, 4, 5)
  - [x] 1.1 Create `web/src/composables/useLocalePreference.ts` following
        useThemePreference pattern
  - [x] 1.2 Support locale values: auto, en, de, es, fr, uk
  - [x] 1.3 Implement auto-detection from navigator.language
  - [x] 1.4 Persist to localStorage via getSafeLocalStorage
  - [x] 1.5 Update i18n.global.locale.value on change
- [x] Task 2: Add language selector to App.vue (AC: 1)
  - [x] 2.1 Add language selector in desktop user menu (below theme selector)
  - [x] 2.2 Add language selector in mobile drawer (below theme selector)
  - [x] 2.3 Show flag emoji + label for each option
- [x] Task 3: Add language-related i18n keys (AC: 1)
  - [x] 3.1 Add language selector labels to en.json
- [x] Task 4: Update i18n plugin to support all locales (AC: 2)
  - [x] 4.1 Import all locale files in i18n.ts
  - [x] 4.2 Register all locales in messages config
- [x] Task 5: Add unit tests (AC: 2, 3, 4, 5)
  - [x] 5.1 Test useLocalePreference composable
  - [x] 5.2 Update App.test.ts for language selector
- [x] Task 6: Verify E2E tests still pass (AC: 1)
  - [x] 6.1 Run full E2E suite

## Dev Notes

### Architecture and Patterns

#### useLocalePreference composable

Follow the useThemePreference pattern:
- Storage key: `sithub_locale`
- Type: `'auto' | 'en' | 'de' | 'es' | 'fr' | 'uk'`
- Auto-detection: `navigator.language.split('-')[0]` mapped to supported locales
- Update i18n: `i18n.global.locale.value = resolvedLocale`
- Use `getSafeLocalStorage()` for persistence

#### Language Options

| Value | Label | Flag |
|-------|-------|------|
| auto | Auto | (globe or auto icon) |
| en | English | GB flag |
| de | Deutsch | DE flag |
| es | Espanol | ES flag |
| fr | Francais | FR flag |
| uk | Ukrainska | UA flag |

Use Unicode flag emojis for flags.

#### i18n Plugin Updates

`web/src/plugins/i18n.ts` needs to import all locale files and register them.
Translation files for de/es/fr/uk will be created in Story 21.3 — for now,
create empty stub files so the selector works (falling back to English via
vue-i18n's fallbackLocale mechanism).

### Anti-Pattern Prevention

- DO NOT create a Pinia store for locale — use a composable like theme preference
- DO NOT add locale persistence to the backend — localStorage only
- DO NOT translate area/item names from YAML config
- DO NOT use `$i18n.locale` directly in components — use the composable

### Source Files to Modify

| File | Changes |
|------|---------|
| web/src/composables/useLocalePreference.ts | New composable |
| web/src/plugins/i18n.ts | Import all locales |
| web/src/App.vue | Add language selector |
| web/src/locales/en.json | Add language labels |

### Files to Create

| File | Purpose |
|------|---------|
| web/src/composables/useLocalePreference.ts | Locale preference composable |
| web/src/locales/de.json | German stub (empty or minimal) |
| web/src/locales/es.json | Spanish stub |
| web/src/locales/fr.json | French stub |
| web/src/locales/uk.json | Ukrainian stub |

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None

### Completion Notes List

- Created useLocalePreference composable with module-level shared state, sync watcher, auto-detection
- Language selector uses v-select with flag emojis (Unicode) in both desktop and mobile menus
- i18n plugin updated to import all 5 locale files (en + 4 stubs)
- Stub locale files created for de/es/fr/uk (empty — will be filled in Story 21.3)
- 7 unit tests for useLocalePreference, all 252 unit tests pass, 54 E2E tests pass

### Change Log

- 2026-04-04: Story 21.2 implementation — language selector with auto-detection

### File List

New files:
- web/src/composables/useLocalePreference.ts
- web/src/composables/useLocalePreference.test.ts
- web/src/locales/de.json (stub)
- web/src/locales/es.json (stub)
- web/src/locales/fr.json (stub)
- web/src/locales/uk.json (stub)

Modified files:
- web/src/plugins/i18n.ts (import all locales)
- web/src/App.vue (language selector in both menus)
- web/src/locales/en.json (language selector labels)
