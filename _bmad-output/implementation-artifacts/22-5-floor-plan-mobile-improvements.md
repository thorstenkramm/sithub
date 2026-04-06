# Story 22.5: Floor Plan Mobile Improvements

Status: done

## Story

As a mobile user,
I want the floor plan to be readable on my phone and adapt to dark mode,
so that I can use it effectively without manual zoom adjustments.

## Acceptance Criteria

1. **Given** I open the floor plan on a 390px-wide screen
   **When** it renders
   **Then** the initial zoom level fits the floor plan width to the viewport

2. **Given** dark mode is active
   **When** the floor plan image renders
   **Then** a CSS filter reduces brightness contrast with the dark UI

3. **Given** I open the floor plan editor on a viewport narrower than 768px
   **When** the editor page renders
   **Then** a banner recommends using a desktop screen

## Tasks / Subtasks

- [ ] Task 1: Auto-zoom floor plan to fit mobile viewport (AC: 1)
  - [ ] 1.1 In `web/src/components/InteractiveFloorPlan.vue`: the `zoomScale`
    ref starts at 1.0 (line 440). On mount, calculate the ratio of viewport
    width to floor plan image natural width, and set `zoomScale` to that ratio
    (clamped to the existing 0.75-2.0 range)
  - [ ] 1.2 Add an `onMounted` or `onLoad` handler for the `<img>` element
    (line 113, class `fp-image-fit`) that reads `img.naturalWidth` and
    computes the initial zoom. Use `nextTick` to ensure the container is
    rendered
  - [ ] 1.3 Only apply auto-zoom on mobile viewports (check
    `window.innerWidth < 768`). On desktop, keep the default 1.0 zoom
- [ ] Task 2: Apply dark mode filter to floor plan image (AC: 2)
  - [ ] 2.1 In `web/src/components/InteractiveFloorPlan.vue`: detect dark mode
    via Vuetify's `useTheme()` composable (`theme.global.current.value.dark`)
  - [ ] 2.2 Apply a CSS class to the `.fp-image-fit` image when dark mode is
    active: `filter: brightness(0.85) contrast(1.1)` — reduces glare without
    fully inverting the image
  - [ ] 2.3 Also apply the same filter in the floor plan editor
    (`web/src/views/FloorPlanEditorView.vue` line ~200)
- [ ] Task 3: Add desktop recommendation banner in floor plan editor (AC: 3)
  - [ ] 3.1 In `web/src/views/FloorPlanEditorView.vue`: add a `v-alert` with
    `type="info"` that shows only on narrow viewports (use
    `useDisplay()` composable from Vuetify or `window.matchMedia`)
  - [ ] 3.2 Message: `$t('floorPlanEditor.desktopRecommended')` — add i18n key
    to all locale files: "Der Raumplan-Editor funktioniert am besten auf einem
    Desktop-Bildschirm." (de) / "The floor plan editor works best on a desktop
    screen." (en) / equivalents for es, fr, uk
- [ ] Task 4: Run tests and lint (AC: 1, 2, 3)
  - [ ] 4.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

## Dev Notes

### Auto-Zoom Calculation

```typescript
const img = document.querySelector('.fp-image-fit') as HTMLImageElement;
if (img && img.naturalWidth > 0 && window.innerWidth < 768) {
  const containerWidth = img.parentElement?.clientWidth ?? window.innerWidth;
  zoomScale.value = Math.min(
    Math.max(containerWidth / img.naturalWidth, 0.75), 2.0
  );
}
```

Wait for the image `load` event before reading `naturalWidth`.

### Dark Mode Detection

Vuetify provides `useTheme()`:

```typescript
import { useTheme } from 'vuetify';
const theme = useTheme();
const isDark = computed(() => theme.global.current.value.dark);
```

### Zoom Range

Current zoom: 0.75 to 2.0 (InteractiveFloorPlan.vue lines 769-771).
The `.fp-scroll-shell` handles overflow with `overflow: auto` (line 1419).
Mobile compact mode uses `.fp-scroll-shell--compact` with `flex: 1; min-height: 0`.

### Files to Change

| File | Change |
| --- | --- |
| `web/src/components/InteractiveFloorPlan.vue` | Auto-zoom, dark filter |
| `web/src/views/FloorPlanEditorView.vue` | Dark filter, desktop banner |
| `web/src/locales/*.json` | Add `floorPlanEditor.desktopRecommended` key |

### References

- [Source: private/ux-observations.md — "Floor plan too small", "dark mode"]
- [Source: private/ux-observations.md — "Floor Plan Editor Not Practical"]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

### Review Findings

- [x] [Review][Patch] Mobile floor-plan auto-fit is lost on week changes because `initialLoad()` resets `zoomScale` to `1` even when the image source is unchanged and the load handler does not rerun [web/src/components/InteractiveFloorPlan.vue:837]
