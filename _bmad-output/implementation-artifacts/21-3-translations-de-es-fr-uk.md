# Story 21.3: German, Spanish, French, and Ukrainian Translations

Status: done

## Story

As a user,
I want all UI text translated into German, Spanish, French, and Ukrainian,
so that I can use SitHub fully in my preferred language.

## Acceptance Criteria

1. **Given** the language is set to German (or Spanish, French, Ukrainian),
   **when** I navigate through the app,
   **then** all labels, buttons, messages, headings, placeholders, and error messages
   appear in the selected language.

2. **Given** translation files exist for all four languages,
   **when** a developer inspects them,
   **then** every key present in the English file has a corresponding entry in each
   translation file with no missing keys.

3. **Given** the backend returns error messages (e.g., booking conflicts),
   **when** the frontend displays them,
   **then** the messages are localized using frontend translation keys, not raw backend
   strings.

## Tasks / Subtasks

- [x] Task 1: Create German translation file (AC: 1, 2)
  - [x] 1.1 Create `web/src/locales/de.json` with all keys from en.json translated
- [x] Task 2: Create Spanish translation file (AC: 1, 2)
  - [x] 2.1 Create `web/src/locales/es.json` with all keys from en.json translated
- [x] Task 3: Create French translation file (AC: 1, 2)
  - [x] 3.1 Create `web/src/locales/fr.json` with all keys from en.json translated
- [x] Task 4: Create Ukrainian translation file (AC: 1, 2)
  - [x] 4.1 Create `web/src/locales/uk.json` with all keys from en.json translated
- [x] Task 5: Verify key completeness (AC: 2)
  - [x] 5.1 Write or run a script to verify all en.json keys exist in each translation
- [x] Task 6: Run tests (AC: 1)
  - [x] 6.1 Run unit tests
  - [x] 6.2 Run E2E tests
  - [x] 6.3 Run lint

## Dev Notes

### Architecture and Patterns

Each translation file must mirror the exact key structure of `web/src/locales/en.json`.
All keys must be present — no missing keys allowed. vue-i18n will fall back to English
for any missing keys, but completeness is required per AC2.

### Translation Quality

Provide natural, contextual translations — not literal word-for-word. For example:
- "My Bookings" in German should be "Meine Buchungen" (not "Mein Buchungen")
- Buttons like "Book" should use the imperative form appropriate to the language
- Error messages should sound natural in the target language

### Plural Forms

vue-i18n pipe `|` syntax is used for plurals. Each language file must maintain the
same pipe-separated plural structure as en.json.

### Anti-Pattern Prevention

- DO NOT change en.json structure — translation files must mirror it exactly
- DO NOT add new keys that don't exist in en.json
- DO NOT translate operator-defined content (area names, item names from YAML)
- DO NOT translate the `CONNECTION_LOST_MESSAGE` constant in api/client.ts

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None

### Completion Notes List

- Created all 4 translation files with 204 keys each matching en.json exactly
- German (de.json): natural German — Bereiche, Meine Buchungen, Buchen, Buchung stornieren
- Spanish (es.json): Castilian Spanish — Zonas, Mis Reservas, Reservar, Cancelar reserva
- French (fr.json): natural French — Espaces, Mes Reservations, Reserver, Annuler la reservation
- Ukrainian (uk.json): natural Ukrainian with 3-form plurals where applicable
- Key completeness verified programmatically: all 204 keys present in every file
- E2E tests required adding locale pinning in cypress/support/e2e.ts to prevent auto-detection from OS locale (German machine) breaking English text assertions
- All 252 unit tests pass, all 54 E2E tests pass, lint/type-check/build clean

### Change Log

- 2026-04-04: Story 21.3 implementation — complete translations for de/es/fr/uk

### File List

Modified files:
- web/src/locales/de.json (complete German translation, 204 keys)
- web/src/locales/es.json (complete Spanish translation, 204 keys)
- web/src/locales/fr.json (complete French translation, 204 keys)
- web/src/locales/uk.json (complete Ukrainian translation, 204 keys)
- web/cypress/support/e2e.ts (force English locale for E2E tests)
