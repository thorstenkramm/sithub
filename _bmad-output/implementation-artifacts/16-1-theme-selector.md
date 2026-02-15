# Story 16.1: Theme Selector

Status: done

## Story

As a user,
I want to choose between light, dark, and auto themes,
So that the app matches my visual preference or adapts to my system setting.

## Acceptance Criteria

1. **Given** I click my username in the top right corner
   **When** I see the user menu
   **Then** I see a theme option with choices: auto (default), dark, light

2. **Given** I select "dark" theme
   **When** the selection is applied
   **Then** the Vuetify dark theme is activated immediately
   **And** my choice is persisted in localStorage
   **And** on my next visit, the dark theme is applied without manual selection

3. **Given** I select "auto" theme
   **When** my OS preference is dark mode
   **Then** the app uses dark theme
   **And** when my OS switches to light mode, the app follows

## Tasks / Subtasks

- [x] Add theme selector to desktop user menu (AC: 1)
  - [x] Added `v-list-item` with `v-btn-toggle` (auto/light/dark) after divider
  - [x] Added `data-cy="theme-selector"`
- [x] Add theme selector to mobile navigation drawer (AC: 1)
  - [x] Added matching `v-btn-toggle` in mobile drawer
  - [x] Added `data-cy="mobile-theme-selector"`
- [x] Implement theme switching logic (AC: 2, 3)
  - [x] Created `web/src/composables/useThemePreference.ts` composable
  - [x] Stores preference in localStorage key `sithub_theme`
  - [x] Defaults to `'auto'`
  - [x] Auto mode uses `matchMedia` with change listener and cleanup via `onScopeDispose`
  - [x] Applies theme via Vuetify's `useTheme()` composable
  - [x] Added guard for environments without `matchMedia` (jsdom/SSR)
- [x] Initialize theme on app mount (AC: 2, 3)
  - [x] Composable called in App.vue `<script setup>` — applies immediately
- [x] Add unit tests
  - [x] 6 tests: stored preferences, auto fallback, persistence, matchMedia usage
  - [x] Table-driven test for preference reading
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

New composable + App.vue changes. No backend changes required.

### Vuetify Theme Configuration

Both light and dark themes are already defined in `web/src/plugins/vuetify.ts` (line 198):

```typescript
theme: {
  defaultTheme: 'light',
  themes: {
    light: lightTheme,
    dark: darkTheme
  }
}
```

The `useTheme()` composable from Vuetify provides `theme.global.name` (a ref) that can be
set to `'light'` or `'dark'` to switch themes at runtime.

### Auto Mode Implementation

For auto mode, use the `prefers-color-scheme` media query:

```typescript
const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
const applyAuto = () => {
  theme.global.name.value = mediaQuery.matches ? 'dark' : 'light';
};
mediaQuery.addEventListener('change', applyAuto);
```

Remember to clean up the listener when switching away from auto mode.

### localStorage Key

Use `sithub_theme` as the key, consistent with the existing `sithub_booking_mode` pattern.

### References

- Epic 16 Story 16.1: `_bmad-output/planning-artifacts/epics.md` (Epic 16 Stories section)
- FR55: `_bmad-output/planning-artifacts/prd.md`
- Vuetify config: `web/src/plugins/vuetify.ts` lines 197-203
- App.vue user menu: `web/src/App.vue` lines 40-79
- App.vue mobile drawer: `web/src/App.vue` lines 90-135

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Created `useThemePreference` composable with auto/light/dark support
- Auto mode uses `matchMedia('(prefers-color-scheme: dark)')` with change listener
- Added `onScopeDispose` cleanup for media listener
- Added guard for `matchMedia` availability (jsdom test environment fallback)
- Added safe localStorage access and guarded `window` usage to avoid SSR/storage-blocked errors
- Both desktop and mobile menus use `v-btn-toggle` with mandatory selection
- App.test.ts updated to include Vuetify plugin (required for `useTheme()`)

### File List

- `web/src/composables/storage.ts` — Safe localStorage helper
- `web/src/composables/useThemePreference.ts` — New composable
- `web/src/composables/useThemePreference.test.ts` — 6 unit tests
- `web/src/App.vue` — Theme selector in desktop menu and mobile drawer
- `web/src/App.test.ts` — Added Vuetify plugin to test setup
