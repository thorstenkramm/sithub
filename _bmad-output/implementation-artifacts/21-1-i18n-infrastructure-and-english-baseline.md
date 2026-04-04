# Story 21.1: i18n Infrastructure and English Baseline

Status: done

## Story

As a developer,
I want vue-i18n configured with all existing UI strings extracted into an English message file,
so that the app is ready for translation without changing user-visible behavior.

## Acceptance Criteria

1. **Given** the app starts, **when** no language preference is stored,
   **then** the UI renders in English, identical to the current behavior.

2. **Given** any component renders text, **when** the text is user-visible
   (labels, buttons, messages, headings, placeholders),
   **then** the text comes from the i18n message file, not hardcoded in the template.

3. **Given** the English message file exists, **when** a developer inspects it,
   **then** all keys are organized by feature area (e.g., `auth.login`,
   `bookings.cancel`, `settings.theme`) and use dot-notation nesting.

## Tasks / Subtasks

- [x] Task 1: Install and configure vue-i18n (AC: 1)
  - [x] 1.1 Add `vue-i18n` dependency to `web/package.json`
  - [x] 1.2 Create `web/src/plugins/i18n.ts` — configure vue-i18n with English as default locale
  - [x] 1.3 Register i18n plugin in `web/src/main.ts` (before Vuetify)
  - [x] 1.4 Configure Vuetify's locale adapter to use vue-i18n in `web/src/plugins/vuetify.ts`
- [x] Task 2: Create English message file with organized keys (AC: 3)
  - [x] 2.1 Create `web/src/locales/en.json` with all extracted strings
  - [x] 2.2 Organize keys by feature area using dot-notation
- [x] Task 3: Extract hardcoded strings from all components (AC: 2)
  - [x] 3.1 Components: `BookingCard.vue`, `ConfirmDialog.vue`, `DatePickerField.vue`,
        `EmptyState.vue`, `InteractiveFloorPlan.vue`, `PageHeader.vue`, `StatusChip.vue`
  - [x] 3.2 Views: `LoginView.vue`, `AreasView.vue`, `ItemGroupsView.vue`, `ItemsView.vue`,
        `MyBookingsView.vue`, `BookingHistoryView.vue`, `ItemGroupBookingsView.vue`,
        `AreaPresenceView.vue`, `FloorPlanEditorView.vue`, `AccessDeniedView.vue`
  - [x] 3.3 Root: `App.vue` (navigation bar, user menu, theme labels, footer)
  - [x] 3.4 Composables with user-facing strings: `useAuthErrorHandler.ts`
- [x] Task 4: Update unit tests (AC: 1, 2)
  - [x] 4.1 Add i18n test helper/mock for Vitest
  - [x] 4.2 Update existing component tests to provide i18n plugin
- [x] Task 5: Verify E2E tests still pass (AC: 1)
  - [x] 5.1 Run Cypress E2E suite — selectors use `data-cy`, so text changes should not break tests
  - [x] 5.2 Fix any tests that assert on hardcoded English text

## Dev Notes

### Architecture and Patterns

#### vue-i18n Setup

Install `vue-i18n` (latest stable, v10.x for Vue 3). Register in `main.ts` **before** Vuetify
so the locale adapter can bind to it.

```typescript
// web/src/plugins/i18n.ts
import { createI18n } from 'vue-i18n';
import en from '@/locales/en.json';

export const i18n = createI18n({
  legacy: false,          // Composition API mode
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en }
});
```

```typescript
// web/src/main.ts — registration order
app.use(pinia);
app.use(i18n);   // BEFORE vuetify
app.use(router);
app.use(vuetify);
```

#### Vuetify Locale Adapter

Vuetify has built-in component strings (date picker labels, data table headers, etc.).
Use Vuetify's `createVueI18nAdapter` to bridge vue-i18n into Vuetify so both share one
locale. This is configured in `web/src/plugins/vuetify.ts`:

```typescript
import { createVueI18nAdapter } from 'vuetify/locale/adapters/vue-i18n';
import { useI18n } from 'vue-i18n';
import { i18n } from './i18n';

// Inside createVuetify():
locale: {
  adapter: createVueI18nAdapter({ i18n, useI18n })
}
```

Vuetify's own locale messages (`$vuetify.datePicker.*`, etc.) are handled automatically
by the adapter. No need to duplicate Vuetify's internal strings in `en.json`.

#### Message File Organization

`web/src/locales/en.json` — flat JSON with dot-nested keys:

```json
{
  "app": {
    "name": "SitHub",
    "navigation": { "areas": "Areas", "myBookings": "My Bookings", "history": "History" },
    "userMenu": { "theme": "Theme", ... }
  },
  "auth": {
    "signIn": "Sign in to SitHub",
    "email": "Email",
    "password": "Password",
    "signInButton": "Sign in",
    "signInWithEntraId": "Sign in with Entra ID",
    "invalidCredentials": "Invalid email or password",
    "genericError": "An error occurred. Please try again."
  },
  "bookings": { ... },
  "areas": { ... },
  "items": { ... },
  "status": {
    "available": "Available",
    "booked": "Booked",
    "mine": "My Booking",
    "unavailable": "Unavailable",
    "guest": "Guest",
    "pending": "Pending",
    "bookedForMe": "Booked for you",
    "onBehalf": "On behalf"
  },
  "common": {
    "confirm": "Confirm",
    "cancel": "Cancel",
    "save": "Save",
    "delete": "Delete",
    "close": "Close",
    "loading": "Loading...",
    "noData": "No data available"
  }
}
```

#### String Extraction in Templates

Replace hardcoded strings with `$t()` in templates and `t()` in `<script setup>`:

```vue
<!-- Before -->
<v-btn>Sign in</v-btn>

<!-- After -->
<v-btn>{{ $t('auth.signInButton') }}</v-btn>
```

In script setup:

```typescript
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
errorMessage.value = t('auth.invalidCredentials');
```

#### StatusChip.vue — Special Pattern

`StatusChip.vue` has a `statusConfig` record mapping status types to labels.
Replace the hardcoded `label` values with i18n keys looked up at render time:

```typescript
const { t } = useI18n();
const statusConfig: Record<StatusType, { color: string; icon: string; labelKey: string }> = {
  available: { color: 'success', icon: '...', labelKey: 'status.available' },
  booked: { color: 'error', icon: '...', labelKey: 'status.booked' },
  // ...
};
// Use computed or function to resolve: t(statusConfig[status].labelKey)
```

#### ConfirmDialog.vue — Default Props

The component has default props `confirmText: 'Confirm'` and `cancelText: 'Cancel'`.
Change defaults to use i18n:

```typescript
const { t } = useI18n();
const props = withDefaults(defineProps<{...}>(), {
  // Cannot call t() in withDefaults — use empty string sentinel
  confirmText: '',
  cancelText: ''
});
const resolvedConfirmText = computed(() => props.confirmText || t('common.confirm'));
const resolvedCancelText = computed(() => props.cancelText || t('common.cancel'));
```

### Source Files to Modify

| File | Changes |
|------|---------|
| `web/package.json` | Add `vue-i18n` dependency |
| `web/src/main.ts` | Register i18n plugin before Vuetify |
| `web/src/plugins/vuetify.ts` | Add locale adapter |
| `web/src/composables/useAuthErrorHandler.ts` | Replace hardcoded error strings |
| `web/src/components/StatusChip.vue` | Replace 8 status labels |
| `web/src/components/ConfirmDialog.vue` | Replace default button text |
| `web/src/components/BookingCard.vue` | Replace booking card labels |
| `web/src/components/EmptyState.vue` | Verify props pass-through (no hardcoded text) |
| `web/src/components/PageHeader.vue` | Verify props pass-through |
| `web/src/components/InteractiveFloorPlan.vue` | Replace floor plan UI strings |
| `web/src/views/LoginView.vue` | Replace ~8 strings |
| `web/src/views/AreasView.vue` | Replace area labels, empty state text |
| `web/src/views/ItemGroupsView.vue` | Replace item group labels |
| `web/src/views/ItemsView.vue` | Replace item/booking labels |
| `web/src/views/MyBookingsView.vue` | Replace booking management strings |
| `web/src/views/BookingHistoryView.vue` | Replace history labels |
| `web/src/views/ItemGroupBookingsView.vue` | Replace overview labels |
| `web/src/views/AreaPresenceView.vue` | Replace presence labels |
| `web/src/views/FloorPlanEditorView.vue` | Replace editor labels |
| `web/src/views/AccessDeniedView.vue` | Replace access denied message |
| `web/src/App.vue` | Replace nav labels, user menu items, footer |

### Files to Create

| File | Purpose |
|------|---------|
| `web/src/plugins/i18n.ts` | vue-i18n plugin configuration |
| `web/src/locales/en.json` | English message file |

### Testing Guidance

#### Vitest Unit Tests

Create a test helper that provides a mock i18n instance:

```typescript
// web/src/__tests__/helpers/i18n.ts
import { createI18n } from 'vue-i18n';
import en from '@/locales/en.json';

export function createTestI18n() {
  return createI18n({
    legacy: false,
    locale: 'en',
    messages: { en }
  });
}
```

Mount components with i18n in tests:

```typescript
import { createTestI18n } from '@/__tests__/helpers/i18n';
const wrapper = mount(MyComponent, {
  global: { plugins: [createTestI18n()] }
});
```

#### Cypress E2E Tests

E2E tests use `data-cy` selectors and should not break. However, any test that asserts
on visible text content (e.g., `cy.contains('Sign in')`) may need updating. Verify by
running the full E2E suite.

### Anti-Pattern Prevention

- **DO NOT** create a Pinia store for locale in this story. Locale preference persistence
  is Story 21.2's scope. This story only sets up the infrastructure with English as the
  hardcoded default.
- **DO NOT** add language switching UI. That is Story 21.2.
- **DO NOT** create translation files for other languages. That is Story 21.3.
- **DO NOT** use `vue-i18n` legacy mode. Use Composition API mode (`legacy: false`).
- **DO NOT** use `$t()` with interpolation for strings that don't need it — keep it simple.
- **DO NOT** extract area/item names from the YAML config — those are operator-defined
  and not translatable.
- **DO NOT** inline long translation keys. Keep keys short and semantic.
- **DO NOT** forget to handle dynamic strings in composables (e.g., `useAuthErrorHandler.ts`
  constructs error messages).

### Previous Epic Learnings

From Epic 20 and prior work:
- The codebase uses `<script setup lang="ts">` consistently — maintain this pattern.
- Composables follow the `useX` naming convention and use module-level shared refs
  for cross-component state (see `useWeekendPreference.ts`).
- localStorage access must go through `getSafeLocalStorage()` from `composables/storage.ts`.
- Snackbar confirmations use a bottom snackbar pattern (FR78) — ensure snackbar messages
  are also extracted to i18n.

### Project Structure Notes

- All new plugin files go in `web/src/plugins/`
- Locale message files go in `web/src/locales/` (new directory)
- Test helpers go in `web/src/__tests__/helpers/` (create if needed)
- Follows existing project conventions for file naming and module structure

### References

- [Source: epics.md — Story 21.1, FR85]
- [Source: web/src/main.ts — plugin registration order]
- [Source: web/src/plugins/vuetify.ts — Vuetify config, locale adapter needed]
- [Source: web/src/composables/useThemePreference.ts — localStorage pattern reference]
- [Source: web/src/composables/storage.ts — safe localStorage access]
- [Source: web/src/components/StatusChip.vue — hardcoded status labels]
- [Source: web/src/components/ConfirmDialog.vue — default button text]
- [vue-i18n docs: https://vue-i18n.intlify.dev/]
- [Vuetify i18n adapter: https://vuetifyjs.com/en/features/internationalization/]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None

### Completion Notes List

- Installed vue-i18n v10 and configured with Composition API mode (legacy: false)
- Created i18n plugin at `web/src/plugins/i18n.ts` with English as default/fallback locale
- Registered i18n before Vuetify in `web/src/main.ts` plugin chain
- Configured Vuetify locale adapter via `createVueI18nAdapter` in vuetify.ts
- Created comprehensive `web/src/locales/en.json` with ~200 keys organized by feature area
- Extracted all hardcoded strings from 10 views, 5 components, and App.vue
- Used vue-i18n pipe-separated plurals for count-dependent messages
- Added `@` resolve alias to both vite.config.ts and vitest.config.ts for `@/locales` imports
- Created test helper `web/src/__tests__/helpers/i18n.ts` for providing i18n plugin in tests
- Updated 14 test files to provide i18n plugin when mounting components
- Updated main.test.ts to verify i18n plugin registration
- useAuthErrorHandler.ts had no hardcoded UI strings (only navigation logic) — no changes needed
- EmptyState.vue and PageHeader.vue are pass-through components (all text via props) — no changes needed
- All 245 unit tests pass, all 54 E2E tests pass, ESLint clean, build succeeds

### Change Log

- 2026-04-04: Story 21.1 implementation — vue-i18n infrastructure and English baseline

### File List

New files:
- web/src/plugins/i18n.ts
- web/src/locales/en.json
- web/src/__tests__/helpers/i18n.ts

Modified files:
- web/package.json (added vue-i18n dependency)
- web/package-lock.json (updated lockfile)
- web/src/main.ts (registered i18n plugin)
- web/src/plugins/vuetify.ts (added locale adapter)
- web/vite.config.ts (added @ resolve alias)
- web/vitest.config.ts (added @ resolve alias)
- web/src/App.vue (extracted all strings to i18n)
- web/src/components/StatusChip.vue (i18n label keys)
- web/src/components/ConfirmDialog.vue (i18n default button text)
- web/src/components/DatePickerField.vue (i18n default label)
- web/src/components/BookingCard.vue (i18n strings)
- web/src/components/InteractiveFloorPlan.vue (i18n alert text)
- web/src/views/LoginView.vue (i18n strings)
- web/src/views/AreasView.vue (i18n strings)
- web/src/views/AccessDeniedView.vue (i18n strings)
- web/src/views/ItemGroupsView.vue (i18n strings)
- web/src/views/ItemsView.vue (i18n strings)
- web/src/views/MyBookingsView.vue (i18n strings)
- web/src/views/BookingHistoryView.vue (i18n strings)
- web/src/views/ItemGroupBookingsView.vue (i18n strings)
- web/src/views/AreaPresenceView.vue (i18n strings)
- web/src/views/FloorPlanEditorView.vue (i18n strings)
- web/src/views/testHelpers.ts (re-exported createTestI18n)
- web/src/App.test.ts (added i18n plugin)
- web/src/main.test.ts (added i18n mock and assertion)
- web/src/views/AreasView.test.ts (added i18n plugin)
- web/src/views/AreaPresenceView.test.ts (added i18n plugin)
- web/src/views/BookingHistoryView.test.ts (added i18n plugin)
- web/src/views/ItemGroupBookingsView.test.ts (added i18n plugin)
- web/src/views/ItemGroupsView.test.ts (added i18n plugin)
- web/src/views/ItemsView.test.ts (added i18n plugin)
- web/src/views/MyBookingsView.test.ts (added i18n plugin)
- web/src/views/LoginView.test.ts (added i18n plugin)
- web/src/views/AccessDeniedView.test.ts (added i18n plugin)
- web/src/components/__tests__/StatusChip.test.ts (added i18n plugin)
- web/src/components/__tests__/ConfirmDialog.test.ts (added i18n plugin)
- web/src/components/__tests__/InteractiveFloorPlan.test.ts (added i18n plugin)
