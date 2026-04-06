# Story 22.8: User Avatars — Frontend Integration

Status: ready-for-dev

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

- [ ] Task 1: Create avatar API client (AC: 1, 2, 3)
  - [ ] 1.1 Create `web/src/api/avatars.ts` with `getAvatarUrl(userId: string)`
    returning `/api/v1/avatars/{userId}`. No fetch needed — use as `<img :src>`
  - [ ] 1.2 Add `uploadAvatar(file: File)` and `deleteAvatar()` functions
    calling POST/DELETE `/api/v1/me/avatar`
- [ ] Task 2: Replace initials with avatar in navigation (AC: 1, 2)
  - [ ] 2.1 In `web/src/App.vue` lines 48-50: the `v-avatar` shows initials.
    Add a `v-img` inside the avatar with `:src="avatarUrl"` and `@error`
    fallback to initials. Use the current user's ID from auth store
  - [ ] 2.2 Add `avatarUrl` computed to App.vue using `getAvatarUrl(authStore.userId)`
  - [ ] 2.3 Ensure the auth store exposes `userId` (check `useAuthStore`)
- [ ] Task 3: Show avatars in presence view (AC: 3)
  - [ ] 3.1 In `web/src/views/AreaPresenceView.vue`: find the user avatar/initials
    display. Replace initials with `v-img` + fallback pattern from Task 2
  - [ ] 3.2 The presence API response includes `user_id` — use it for avatar URL
- [ ] Task 4: Avatar upload in settings (AC: 4)
  - [ ] 4.1 Add an avatar section to the hamburger menu or a settings page:
    show current avatar preview, upload button, delete button
  - [ ] 4.2 Use `<input type="file" accept="image/png,image/jpeg">` with
    `uploadAvatar()` from the API client
  - [ ] 4.3 On successful upload, refresh the avatar URL by appending a
    cache-busting query parameter (e.g., `?t={timestamp}`)
- [ ] Task 5: Floor plan avatar overlay (AC: 5)
  - [ ] 5.1 In `web/src/components/InteractiveFloorPlan.vue`: when a desk is
    booked and "Show avatars" is checked, render a small `<img>` with the
    booker's avatar on the desk position
  - [ ] 5.2 Add a "Show avatars" checkbox (similar to existing "Show labels"
    checkbox at line ~114). Add i18n key `floorPlan.showAvatars`
  - [ ] 5.3 The floor plan already has booking info per desk — use
    `user_id` from the booking data to build the avatar URL
- [ ] Task 6: Write tests (AC: 1, 2, 3, 4)
  - [ ] 6.1 Test avatar URL generation in `avatars.test.ts`
  - [ ] 6.2 Test App.vue avatar rendering with mock image and error fallback
  - [ ] 6.3 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

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

### Completion Notes List

### File List
