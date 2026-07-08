# Story 36.5: Centered, Filled Presence Avatars

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user viewing Today's presence,
I want avatars centered and filling their circles,
so that the presence view looks clean and consistent.

## Acceptance Criteria

1. Presence avatar photos of varying aspect ratios render centered and cover the entire circular
   frame (`object-fit: cover`, `object-position: center`), with no off-center offset or
   clipped-corner gaps.
2. The initials fallback (no photo) is centered and fills the circle consistently with photo
   avatars.

## Tasks / Subtasks

- [ ] Task 1: Fix the presence avatar markup so photo and initials never stack (AC: #1, #2)
  - [ ] In `web/src/views/AreaPresenceView.vue` (~48-55), the `<v-avatar>` currently renders BOTH
        the `<v-img>` and the initials `<span>` unconditionally, so the fallback text sits under
        every photo. Gate them: render `<v-img>` only when a photo can load, and the initials
        `<span>` as the `v-else` — mirror the working precedent in `ItemsView.vue:303-319` and
        `App.vue:48-55`.
  - [ ] Add an `@error` handler + a `failedAvatars` set (reactive) so a missing/failed photo falls
        back to initials, exactly as `ItemsView.vue:314` (`@error="failedAvatars.add(...)"`) and
        `App.vue:52` (`@error="avatarLoadFailed = true"`) do. Key the set on `user_id`.
- [ ] Task 2: Center + cover the photo and fill the circle with initials (AC: #1, #2)
  - [ ] Ensure the `<v-img>` covers the circle. Vuetify's `v-img` defaults to `cover`, but make it
        explicit and robust for odd aspect ratios: `object-fit: cover; object-position: center;
        width: 100%; height: 100%` on the image, matching the `.fp-item-avatar` precedent
        (`InteractiveFloorPlan.vue:1863-1870`).
  - [ ] Give the initials `<span>` a full-fill, flex-centered style so it fills the circle like a
        photo: `width: 100%; height: 100%; display: (inline-)flex; align-items: center;
        justify-content: center; line-height: 1`. Copy the `.tile-booker-initials` precedent
        (`ItemsView.vue:2243-2255`) or `.fp-item-initials`
        (`InteractiveFloorPlan.vue:1872-1885`). Keep the existing `v-avatar` `color`/`variant` so
        the initials background/tone is unchanged.
  - [ ] Keep `border-radius: 50%` implicit via `<v-avatar>` (the avatar clips to a circle already);
        do not add a square radius. The `<v-img>`/`<span>` just need to fill it.
- [ ] Task 3: Tests (AC: #1, #2)
  - [ ] Vitest (`web/src/views/AreaPresenceView.test.ts`): add a case that an entry WITH a photo
        renders the `<v-img>` and NOT the initials span; an entry WITHOUT a photo (or after an
        image error) renders the initials span and NOT the `<v-img>`. Assert the initials span
        carries the fill/centering class.
  - [ ] Cypress (nice-to-have): visual check on `/areas/:areaId/presence` that avatars are circular
        and centered for both photo and initials rows.

## Dev Notes

The fix is CSS + a small markup gate. No API, store, or route changes. The presence view is the
"Today's presence" page: route `presence` → `/areas/:areaId/presence`
(`web/src/router/index.ts`), component `AreaPresenceView.vue`. Data comes from
`fetchAreaPresence` (`web/src/api/areaPresence.ts`); avatar URLs from
`getAvatarUrl(userId)` → `/api/v1/avatars/:userId` (`web/src/api/avatars.ts:2-4`).

### The broken markup (cite)

`web/src/views/AreaPresenceView.vue:48-55`:

```vue
<template #prepend>
  <v-avatar color="primary" variant="tonal" size="40">
    <v-img :src="getAvatarUrl(entry.attributes.user_id)" />
    <span class="text-body-2 font-weight-medium">
      {{ getInitials(entry.attributes.user_name) }}
    </span>
  </v-avatar>
</template>
```

Two problems versus every working avatar in the app:

1. The `<v-img>` and the initials `<span>` both render unconditionally, so the initials text is
   always painted (it sits under or beside the photo). There is no `v-if`/`v-else`.
2. There is no `@error` fallback, so a user with no photo shows a broken/empty image rather than
   clean initials, and the initials `<span>` has no fill/centering CSS — it is just a
   `text-body-2` span, so it is not guaranteed centered or circle-filling.

`getInitials` already exists (`AreaPresenceView.vue:169-177`) and is fine — keep it.

### Correct precedents to copy

The rest of the app already does this right. Follow one of these patterns:

Markup — `web/src/views/ItemsView.vue:303-319` (`<v-img v-if=... @error=...>` + `<span v-else
class="tile-booker-initials">`), and `web/src/App.vue:48-55` (`v-if` on img with `@error`, initials
`v-if` on the failure/no-user case). Both use a `failedAvatars` set / `avatarLoadFailed` flag.

Photo cover CSS — `.fp-item-avatar` in
`web/src/components/InteractiveFloorPlan.vue:1863-1870`:

```css
.fp-item-avatar {
  position: absolute;
  inset: 1px;
  width: calc(100% - 2px);
  height: calc(100% - 2px);
  object-fit: cover;
  border-radius: 2px;
}
```

For the presence view (inside a circular `<v-avatar>`), drop the `position: absolute`/`inset` and
just use `width: 100%; height: 100%; object-fit: cover; object-position: center` — the avatar
handles the circular clip.

Initials fill/center CSS — `.tile-booker-initials` in `web/src/views/ItemsView.vue:2243-2255`:

```css
.tile-booker-initials {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  font-size: 0.85em;
  font-weight: 600;
  color: white;
  background: rgba(var(--v-theme-error), 0.85);
  user-select: none;
  line-height: 1;
}
```

Reuse this (or `.fp-item-initials`, `InteractiveFloorPlan.vue:1872-1885`). Note the presence
avatar uses `color="primary" variant="tonal"`; keep that tone rather than the error-red background
above if you want to match the current look — the load-bearing parts are the
`width/height: 100%` + flex centering, not the color.

> [!NOTE]
> The scoped `<style>` block in `AreaPresenceView.vue:293-300` currently only holds `.note-text`.
> Add the new `.presence-avatar-img` / `.presence-avatar-initials` rules there.

### Project Structure Notes

- Modified: `web/src/views/AreaPresenceView.vue` (template ~48-55 + scoped `<style>` ~293-300) and
  its test `web/src/views/AreaPresenceView.test.ts`.
- No shared component exists for booker avatars yet (the pattern is duplicated across `App.vue`,
  `ItemsView.vue`, `InteractiveFloorPlan.vue`, `AreaWeeklyMatrixCell.vue`). Extracting one is out of
  scope for this small story; match the existing precedent inline.

### Testing standards summary

Vitest is the primary gate here (`AreaPresenceView.test.ts` already stubs `v-avatar`,
line 27). Test user-visible behavior: photo row shows the image and no initials; no-photo/error row
shows initials and no image. Run `npm run type-check`, `npm run lint`, `npx vitest run`, and
`npm run build`. A Cypress visual E2E on `/areas/:areaId/presence` is a nice-to-have (not required).
[Source: .claude/rules/vue.md, .claude/rules/cypress.md]

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 36.5 / FR173 (lines 5465-5482)]
- [Source: web/src/views/AreaPresenceView.vue:48-55,169-177,293-300]
- [Source: web/src/views/ItemsView.vue:303-319,2239-2255]
- [Source: web/src/App.vue:48-55]
- [Source: web/src/components/InteractiveFloorPlan.vue:1863-1885]
- [Source: web/src/api/avatars.ts:2-4]
- [Source: web/src/router/index.ts (route name "presence")]
- [Source: web/src/views/AreaPresenceView.test.ts:27]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
