# Story 16.3: Fix Change Password Icon

Status: done

## Story

As a user,
I want the Change Password menu item to show its icon,
So that the menu has a consistent visual appearance.

## Acceptance Criteria

1. **Given** I am a local user viewing the desktop user menu
   **When** I see the "Change Password" item
   **Then** an icon is displayed next to the text (consistent with other menu items)

2. **Given** I am a local user viewing the mobile navigation drawer
   **When** I see the "Change Password" item
   **Then** an icon is displayed next to the text

## Tasks / Subtasks

- [x] Investigate desktop icon rendering issue (AC: 1)
  - [x] Root cause: app uses SVG icon set (`mdi-svg` from `@mdi/js`), not webfont. The
    `mdi-lock-reset` syntax is for webfont, not SVG icons.
- [x] Fix desktop icon (AC: 1)
  - [x] Imported `mdiLockReset` from `@mdi/js` and added `lockReset: mdiLockReset` alias
    in `vuetify.ts`
  - [x] Changed `mdi-lock-reset` to `$lockReset` in App.vue desktop menu
- [x] Investigate and fix mobile icon (AC: 2)
  - [x] Same root cause. Changed `mdi-lock-reset` to `$lockReset` in App.vue mobile drawer
- [x] Verify E2E tests still pass

## Dev Notes

### Architecture: Frontend-Only Story

All changes in `web/src/App.vue`. No backend changes required.

### Current Icon State

The `mdi-lock-reset` icon is already in the template at both locations:
- Desktop: line 68 (`<v-icon size="small">mdi-lock-reset</v-icon>`)
- Mobile: line 124 (`<v-icon>mdi-lock-reset</v-icon>`)

Both locations use the `#prepend` template slot on `v-list-item`, same as other menu items.
The icon IS present in code, so this is likely a rendering issue rather than missing code.

### Possible Causes

- The `mdi-lock-reset` icon might not exist in the installed MDI icon set version. Check
  the `@mdi/font` package version in `package.json`.
- If the icon is missing, use an alternative like `mdi-lock-outline`, `mdi-key-change`,
  or `mdi-form-textbox-password`.
- The icon might render but be invisible due to color matching the background.

### Comparison with Working Icons

The "Sign out" item uses `$logout` (a custom alias defined in `vuetify.ts`). Other mobile
drawer items use custom aliases (`$area`, `$calendar`, `$history`). The Change Password item
uses a raw MDI icon name instead — this difference might cause a styling inconsistency.

### References

- Epic 16 Story 16.3: `_bmad-output/planning-artifacts/epics.md` (Epic 16 Stories section)
- FR57: `_bmad-output/planning-artifacts/prd.md`
- App.vue desktop menu: `web/src/App.vue` lines 62-71
- App.vue mobile drawer: `web/src/App.vue` lines 118-127
- Vuetify icon aliases: `web/src/plugins/vuetify.ts`

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Root cause: the project uses `mdi-svg` icon set from `@mdi/js` (not `@mdi/font` webfont).
  The syntax `mdi-lock-reset` is webfont class syntax; SVG icons require importing the JS
  path constant and registering an alias.
- Fix: imported `mdiLockReset` from `@mdi/js`, added `lockReset: mdiLockReset` alias in
  vuetify.ts, and changed both desktop and mobile references from `mdi-lock-reset` to
  `$lockReset`.
- Added stable `data-cy` hooks for change password icons and E2E coverage for desktop
  and mobile icon visibility.

### File List

- `web/src/plugins/vuetify.ts` — Added `mdiLockReset` import and `lockReset` alias
- `web/src/App.vue` — Changed icon references from `mdi-lock-reset` to `$lockReset`
- `web/cypress/e2e/password-change.cy.ts` — Added icon visibility assertions
