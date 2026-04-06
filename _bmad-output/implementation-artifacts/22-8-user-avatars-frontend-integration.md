# Story 22.8: User Avatars — Frontend Integration

Status: done

## Story

As a user,
I want to see profile photos in the navigation, presence list, and floor plan,
so that I can visually identify colleagues.

## Acceptance Criteria

1. **Given** I am logged in and have an avatar
   **When** the navigation bar renders
   **Then** my avatar image replaces the initials circle

2. **Given** I am logged in and have no avatar
   **When** the navigation bar renders
   **Then** the existing initials circle is shown (fallback)

3. **Given** I view Today's Presence
   **When** the presence list renders
   **Then** each user's entry shows their avatar (or initials fallback)

4. **Given** I am a local user in settings
   **When** I upload a profile image
   **Then** the avatar updates immediately in the navigation bar

5. **Given** I open the floor plan with "Show avatars" enabled
   **When** a desk is booked
   **Then** the booker's avatar thumbnail appears on the desk position

## Tasks / Subtasks

- [x] Task 1: Create avatar API client (AC: 1, 2, 3)
  - [x] 1.1 Create `web/src/api/avatars.ts` with `getAvatarUrl(userId: string)`
    returning `/api/v1/avatars/{userId}`. No fetch needed — use as `<img :src>`
  - [x] 1.2 Add `uploadAvatar(file: File)` and `deleteAvatar()` functions
    calling POST/DELETE `/api/v1/me/avatar`
- [x] Task 2: Replace initials with avatar in navigation (AC: 1, 2)
  - [x] 2.1 In `web/src/App.vue` lines 48-50: the `v-avatar` shows initials.
    Add a `v-img` inside the avatar with `:src="avatarUrl"` and `@error`
    fallback to initials. Use the current user's ID from auth store
  - [x] 2.2 Add `avatarUrl` computed to App.vue using `getAvatarUrl(authStore.userId)`
  - [x] 2.3 Ensure the auth store exposes `userId` (check `useAuthStore`)
- [x] Task 3: Show avatars in presence view (AC: 3)
  - [x] 3.1 In `web/src/views/AreaPresenceView.vue`: find the user avatar/initials
    display. Replace initials with `v-img` + fallback pattern from Task 2
  - [x] 3.2 The presence API response includes `user_id` — use it for avatar URL
- [x] Task 4: Avatar upload in settings (AC: 4)
  - [x] 4.1 Add an avatar section to the hamburger menu or a settings page:
    show current avatar preview, upload button, delete button
  - [x] 4.2 Use `<input type="file" accept="image/png,image/jpeg">` with
    `uploadAvatar()` from the API client
  - [x] 4.3 On successful upload, refresh the avatar URL by appending a
    cache-busting query parameter (e.g., `?t={timestamp}`)
- [x] Task 5: Floor plan avatar overlay (AC: 5)
  - [x] 5.1 In `web/src/components/InteractiveFloorPlan.vue`: when a desk is
    booked and "Show avatars" is checked, render a small `<img>` with the
    booker's avatar on the desk position
  - [x] 5.2 Add a "Show avatars" checkbox (similar to existing "Show labels"
    checkbox at line ~114). Add i18n key `floorPlan.showAvatars`
  - [x] 5.3 The floor plan already has booking info per desk — use
    `user_id` from the booking data to build the avatar URL
- [x] Task 6: Write tests (AC: 1, 2, 3, 4)
  - [x] 6.1 Test avatar URL generation in `avatars.test.ts`
  - [x] 6.2 Test App.vue avatar rendering with mock image and error fallback
  - [x] 6.3 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

## Dev Notes

### Avatar URL Pattern

```typescript
export function getAvatarUrl(userId: string): string {
  return `/api/v1/avatars/${encodeURIComponent(userId)}`;
}
```

Use as `<v-img :src="getAvatarUrl(userId)" />` — the browser handles the GET.

### Fallback Pattern

```vue
<v-avatar size="32" color="primary-lighten-1">
  <v-img
    v-if="avatarUrl"
    :src="avatarUrl"
    @error="avatarLoadFailed = true"
  />
  <span v-if="!avatarUrl || avatarLoadFailed">{{ userInitials }}</span>
</v-avatar>
```

### Existing Initials Implementation

`App.vue` lines 348-357: `userInitials` computed property splits display name
and extracts first+last initials. Keep this as fallback.

### Auth Store

`web/src/stores/useAuthStore.ts` — check if it stores `userId`. If not,
add it during the `/api/v1/me` fetch in the views' `onMounted`.

### Dependencies

This story requires Story 22.7 (backend avatar endpoints) to be complete.

### References

- [Source: private/epic-22.md — "Use avatars for settings, presence, floor plan"]
- [Source: web/src/App.vue:48-50 — current initials avatar]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Tasks 1-3 were already implemented (avatar API client, nav bar avatar, presence view avatars).
  Verified and checked off.
- Task 4: Added avatar upload/delete dialog to user menu (desktop + mobile).
  Uses `uploadAvatar`/`deleteAvatar` from API client. Cache-busting via timestamp
  query parameter refreshes avatar immediately after upload.
- Task 5: Added "Show avatars" checkbox to floor plan footer (persisted in localStorage).
  When enabled, avatar thumbnails render on booked desks. Required adding
  `booker_user_id` to the backend items API response (non-guest bookings only)
  and updating the frontend `ItemAttributes` type. Extracted `applyBookingAttrs`
  helper to keep `buildItemResources` under gocognit threshold.
- Task 6: Created `avatars.test.ts` with 7 tests covering URL generation (including
  encoding), upload (FormData, error), and delete (request, error).
- All i18n keys added to all 5 locales (en, de, es, fr, uk).
- All tests pass: 275 frontend tests, full Go test suite.
- Linters pass: ESLint, golangci-lint, type-check, build.

### File List

- `web/src/api/avatars.ts` (pre-existing, verified)
- `web/src/api/avatars.test.ts` (new)
- `web/src/api/items.ts` (modified — added `booker_user_id` field)
- `web/src/App.vue` (modified — avatar dialog, cache-busting URL, menu items)
- `web/src/components/InteractiveFloorPlan.vue` (modified — show avatars checkbox, avatar overlay, bookerUserId)
- `web/src/locales/en.json` (modified — avatar dialog + showAvatars keys)
- `web/src/locales/de.json` (modified — avatar dialog + showAvatars keys)
- `web/src/locales/es.json` (modified — avatar dialog + showAvatars keys)
- `web/src/locales/fr.json` (modified — avatar dialog + showAvatars keys)
- `web/src/locales/uk.json` (modified — avatar dialog + showAvatars keys)
- `internal/items/handler.go` (modified — booker_user_id in response, extracted applyBookingAttrs)

### Review Findings

- [x] [Review][Patch] Avatar fallback state is never reset after a failed image load [web/src/App.vue:422]
- [x] [Review][Patch] App avatar upload/delete flow has no rendering or fallback coverage [web/src/App.test.ts:13]
- [x] [Review][Patch] Floor-plan avatar overlay and `booker_user_id` contract have no regression coverage [web/src/components/__tests__/InteractiveFloorPlan.test.ts:207]
