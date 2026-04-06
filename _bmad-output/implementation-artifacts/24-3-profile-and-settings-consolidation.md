# Story 24.3: Profile and Settings Consolidation

Status: done

## Story

As a user,
I want all settings accessible from a single Profile menu,
so that I don't have to choose between two overlapping menus to find what I need.

## Acceptance Criteria

1. **Given** I am logged in and viewing the app,
   **When** I look at the navigation,
   **Then** there is no separate "Settings" menu option; only the Profile menu
   (avatar/initials) exists.

2. **Given** I open the Profile menu,
   **When** I view the menu contents,
   **Then** all settings are present: theme selector, language selector, show weekends
   toggle, and change password option.

3. **Given** I open the Profile menu on mobile,
   **When** I view the menu contents,
   **Then** the same settings are available with the current profile layout styling.

4. **Given** I previously accessed a setting via the old Settings menu,
   **When** I look for it after the consolidation,
   **Then** the setting is accessible from the Profile menu with no functionality lost.

## Tasks / Subtasks

- [x] Task 1: Consolidate desktop menu sections in App.vue (AC: #1, #2, #4)
  - [x] 1.1 Remove the divider separating the Settings section from the Profile section in the desktop v-menu
  - [x] 1.2 Reorder: move profile action items (floor plan editor, avatar, change password) ABOVE the settings controls (theme, language, weekends)
  - [x] 1.3 Keep divider before logout to visually separate the destructive action
  - [x] 1.4 Verify all data-cy selectors are preserved

- [x] Task 2: Consolidate mobile menu sections in App.vue (AC: #3, #4)
  - [x] 2.1 Remove the divider separating the Settings section from the Profile section in the mobile v-navigation-drawer
  - [x] 2.2 Apply same reorder as desktop: profile actions above settings controls
  - [x] 2.3 Keep divider before logout
  - [x] 2.4 Verify all mobile data-cy selectors are preserved

- [x] Task 3: Write unit tests for consolidated menu (AC: #1-#4)
  - [x] 3.1 Test: desktop menu contains theme selector, language selector, weekends toggle, change password, and logout in a single section
  - [x] 3.2 Test: mobile menu contains the same items in a single section

- [x] Task 4: Run full validation suite
  - [x] 4.1 `npm run type-check` passes
  - [x] 4.2 `npm run lint` passes
  - [x] 4.3 `npx vitest run` — 308 tests, all pass
  - [x] 4.4 `npm run build` succeeds

## Dev Notes

### Architecture and Patterns

- **All menu code is in**: `web/src/App.vue`
- **Desktop menu**: `v-menu` at lines 40-146 with `v-list` at line 59
- **Mobile menu**: `v-navigation-drawer` at lines 157-265 with `v-list` at line 158
- **No separate Settings component** — settings and profile are sections in the same menu, separated by dividers

### Current Desktop Menu Structure (lines 59-145)

```
v-list
  ├── Header: user name + admin badge (lines 60-65)
  ├── v-divider (line 66)
  ├── SETTINGS: theme (67-85), language (86-100), weekends (101-108)
  ├── v-divider (line 109) ← REMOVE THIS
  └── PROFILE: floor plan editor (110-119), avatar (120-128),
      change password (129-138), logout (139-144)
```

### Current Mobile Menu Structure (lines 158-264)

```
v-list
  ├── Header: user name + admin badge (159-164)
  ├── v-divider (line 165)
  ├── NAVIGATION: Areas, My Bookings, History (166-183)
  ├── v-divider (line 184)
  ├── SETTINGS: theme (185-203), language (204-218), weekends (219-226)
  ├── v-divider (line 227) ← REMOVE THIS
  └── PROFILE: floor plan editor (228-238), avatar (239-247),
      change password (248-257), logout (258-263)
```

### Target Structure (Desktop)

```
v-list
  ├── Header: user name + admin badge
  ├── v-divider
  ├── Floor Plan Editor (admin only)
  ├── Avatar
  ├── Change Password (internal auth only)
  ├── Theme selector
  ├── Language selector
  ├── Weekends toggle
  ├── v-divider
  └── Logout
```

### Target Structure (Mobile)

```
v-list
  ├── Header: user name + admin badge
  ├── v-divider
  ├── NAVIGATION: Areas, My Bookings, History
  ├── v-divider
  ├── Floor Plan Editor (admin only)
  ├── Avatar
  ├── Change Password (internal auth only)
  ├── Theme selector
  ├── Language selector
  ├── Weekends toggle
  ├── v-divider
  └── Logout
```

### Key data-cy Selectors to Preserve

Desktop: `theme-selector`, `language-selector`, `show-weekends-toggle`, `floor-plan-editor-btn`, `avatar-btn`, `change-password-btn`, `change-password-icon`, `logout-btn`

Mobile: `mobile-theme-selector`, `mobile-language-selector`, `mobile-show-weekends-toggle`, `mobile-floor-plan-editor-btn`, `mobile-avatar-btn`, `mobile-change-password-btn`, `mobile-change-password-icon`, `mobile-logout-btn`

### Composables Used

- `useThemePreference()` → `themePreference`, `setThemePreference`
- `useLocalePreference()` → `localePreference`, `setLocalePreference`
- `useWeekendPreference()` → `showWeekends`

### Testing

- `web/src/App.test.ts` exists with 2 tests (router mounting, avatar rendering)
- No existing tests cover menu structure — new tests needed
- 306 tests currently passing

### Project Structure Notes

- All changes scoped to `web/src/App.vue` and `web/src/App.test.ts`
- No backend changes
- No new i18n keys needed — existing keys unchanged
- No new composables or components

### References

- [Source: web/src/App.vue — desktop menu lines 40-146]
- [Source: web/src/App.vue — mobile menu lines 157-265]
- [Source: web/src/App.test.ts — existing tests]
- [Source: _bmad-output/planning-artifacts/epics.md — Epic 24 Story 24.3]
- [Source: private/epic-24.md — original requirement]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

None — straightforward template reorder with no compilation or runtime issues.

### Completion Notes List

- Removed divider between Settings and Profile sections in desktop menu (was line 109)
- Removed divider between Settings and Profile sections in mobile menu (was line 227)
- Reordered both menus: profile actions (floor plan editor, avatar, change password) now appear above settings controls (theme, language, weekends)
- Logout remains last with its own divider for visual separation
- All data-cy selectors preserved (8 desktop + 8 mobile)
- Added 2 new tests verifying all menu items present in both desktop and mobile menus
- All validation gates pass: type-check, lint, 308 tests, build

### Change Log

- 2026-04-06: Story 24.3 implementation complete — profile and settings consolidation

### File List

- web/src/App.vue (modified)
- web/src/App.test.ts (modified)
