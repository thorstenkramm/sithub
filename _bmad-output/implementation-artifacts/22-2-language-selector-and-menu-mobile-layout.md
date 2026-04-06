# Story 22.2: Language Selector and Menu Mobile Layout

Status: done

## Story

As a mobile user,
I want the language and theme buttons to fit the navigation drawer without clipping,
so that I can read and tap each option.

## Acceptance Criteria

1. **Given** I open the hamburger menu on a 390px-wide screen
   **When** the language buttons render
   **Then** all language names and flags are fully visible (no text clipping)

2. **Given** I open the hamburger menu on a 390px-wide screen
   **When** the theme toggle renders
   **Then** all three options (Automatisch, Hell, Dunkel) are fully readable

## Tasks / Subtasks

- [ ] Task 1: Fix language selector layout in mobile drawer (AC: 1)
  - [ ] 1.1 In `web/src/App.vue` lines 190-204: the `.locale-grid` CSS class
    (lines 429-433) uses `grid-template-columns: repeat(3, 1fr)` with `gap: 4px`.
    On a ~256px drawer this clips button text. Change to
    `grid-template-columns: repeat(2, 1fr)` for 2 columns, giving each button ~120px
  - [ ] 1.2 Alternatively, replace the grid of buttons with a `v-select` dropdown
    (saves space, fits any drawer width). Use the flag emoji + language name as
    item titles
- [ ] Task 2: Fix theme toggle in mobile drawer (AC: 2)
  - [ ] 2.1 In `web/src/App.vue` lines 171-189: the `v-btn-toggle` for theme uses
    `size="small"`. On narrow drawer, "AUTOMATISCH" clips. Either use shorter
    labels on mobile (e.g., "Auto", "Hell", "Dunkel") or switch to a `v-select`
- [ ] Task 3: Apply same fix to desktop menu if affected (AC: 1, 2)
  - [ ] 3.1 Check `web/src/App.vue` lines 62-95 (desktop user menu) for the same
    overflow. The desktop menu should have more width but verify
- [ ] Task 4: Run tests and lint (AC: 1, 2)
  - [ ] 4.1 Run `npx vitest run`, `npm run lint`, `npm run type-check`, `npm run build`

## Dev Notes

### Current Layout

Mobile drawer: `v-navigation-drawer` at `location="right"`, default width ~256px.
Language grid: `display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px`
(App.vue line 431). Six buttons in 3 columns = 2 rows. Each button ~80px wide —
not enough for "AUTOMATISCH" (12 chars) or flag emoji + "UKRAINSKA".

### Recommended Approach

A `v-select` dropdown is the cleanest mobile solution. It uses zero horizontal
space beyond one line, shows all options in a Vuetify menu overlay, and supports
the flag emoji + text as item titles. This approach also eliminates the theme
toggle overflow.

### Files to Change

- `web/src/App.vue` — mobile drawer section (lines 143-242), CSS (lines 429-433)

### References

- [Source: private/ux-observations.md — "Language selector buttons clipped"]
- [Source: private/epic-22.md — "Odd language selector on mobile"]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List

### Review Findings

- [x] [Review][Patch] Mobile theme selector hardcodes "Auto" instead of the localized `themeAuto` label, so the drawer regresses to English copy in non-English locales [web/src/App.vue:325]
