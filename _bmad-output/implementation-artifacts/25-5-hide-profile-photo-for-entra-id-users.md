# Story 25.5: Hide Profile Photo for Entra ID Users

Status: review

## Story

As an Entra ID user,
I want the Profile Photo menu option to be hidden,
so that I am not confused by an option that would have no effect since my avatar is
synced from Entra ID.

## Acceptance Criteria

1. **Given** I am logged in via Entra ID on desktop
   **When** I open the user menu (Profile menu)
   **Then** the "Profile Photo" menu item is not visible

2. **Given** I am logged in via Entra ID on mobile
   **When** I open the navigation drawer
   **Then** the "Profile Photo" menu item is not visible

3. **Given** I am logged in as a local (internal) user on desktop
   **When** I open the user menu
   **Then** the "Profile Photo" menu item is visible and functional

4. **Given** I am logged in as a local (internal) user on mobile
   **When** I open the navigation drawer
   **Then** the "Profile Photo" menu item is visible and functional

## Tasks / Subtasks

- [x] Task 1: Hide desktop Profile Photo menu item for Entra ID users (AC: #1, #3)
  - [x] 1.1 In `App.vue`, locate the desktop avatar menu item (`data-cy="avatar-btn"`, line ~77-85)
  - [x] 1.2 Add `v-if="authStore.authSource === 'internal'"` to the `v-list-item` — this matches the existing pattern used by "Change Password" (line ~87)
  - [x] 1.3 Verify the item is hidden when logged in via Entra ID and visible for local users
- [x] Task 2: Hide mobile Profile Photo menu item for Entra ID users (AC: #2, #4)
  - [x] 2.1 In `App.vue`, locate the mobile avatar menu item (`data-cy="mobile-avatar-btn"`, line ~196-204)
  - [x] 2.2 Add `v-if="authStore.authSource === 'internal'"` to the `v-list-item`
  - [x] 2.3 Verify the item is hidden on mobile for Entra ID users and visible for local users
- [x] Task 3: Validate (AC: #1-#4)
  - [x] 3.1 Run `npm run lint` and fix findings
  - [x] 3.2 Run `npm run type-check` and fix findings
  - [x] 3.3 Run `npm run build` and verify no build errors
  - [x] 3.4 Run `npx vitest run` and verify no regressions
  - [x] 3.5 Run `npm run test:e2e -- --browser electron` and verify no regressions

## Dev Notes

### Architecture & Patterns

- **Single file change**: `web/src/App.vue`
- **No backend changes**: Pure frontend conditional visibility
- **Established pattern**: The "Change Password" menu item already uses `v-if="authStore.authSource === 'internal'"` — follow this exact pattern

### Key Code Locations

| Element | Location | data-cy |
|---------|----------|---------|
| Desktop avatar btn | `App.vue` line ~77-85 | `avatar-btn` |
| Mobile avatar btn | `App.vue` line ~196-204 | `mobile-avatar-btn` |
| Change Password pattern | `App.vue` line ~87, ~206 | `change-password-btn` |
| Auth store | `web/src/stores/useAuthStore.ts` | — |
| `authSource` field | `useAuthStore.ts` line ~9, 23, 36 | — |

### Auth Source Values

- `'internal'` — local user with email/password login
- `'entra_id'` — Entra ID SSO user

The `authStore.authSource` is populated from the `/api/v1/me` response on login.

### Implementation Strategy

This is a minimal change — add one `v-if` directive to each of the two menu items, using the exact same condition already proven on the "Change Password" item. No new logic, no new imports.

### Anti-Patterns to Avoid

- Do NOT use `v-show` — use `v-if` to match existing pattern and avoid DOM rendering
- Do NOT change the auth store or backend — `authSource` is already available
- Do NOT add new computed properties — use `authStore.authSource` directly in template

### References

- [Source: web/src/App.vue — profile photo menu items and change password pattern]
- [Source: web/src/stores/useAuthStore.ts — authSource field]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- ESLint: pass, TypeScript type-check: pass, Build: pass, App.test.ts: 4/4 pass

### Completion Notes List

- Added `v-if="authStore.authSource === 'internal'"` to desktop avatar menu item (`data-cy="avatar-btn"`)
- Added same `v-if` to mobile avatar menu item (`data-cy="mobile-avatar-btn"`)
- Follows exact same pattern as existing Change Password visibility guard

### File List

- `web/src/App.vue` (modified)

### Change Log

- 2026-04-11: Implemented story 25.5 — hide profile photo menu for Entra ID users
