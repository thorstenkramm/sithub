# Story 10.1: Design System Foundation

## Story

**As a** user,  
**I want** a visually consistent and branded experience,  
**So that** the application feels professional and trustworthy.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** High (foundational for all other UI stories)

## Acceptance Criteria

**AC1: Custom Theme Applied**
- **Given** the application loads
- **When** I view any page
- **Then** I see the custom color scheme (not Vuetify defaults)
- **And** typography uses Inter font family
- **And** spacing is consistent throughout

**AC2: Brand Identity**
- **Given** the application loads
- **When** I look at the browser tab
- **Then** I see a custom favicon
- **And** the app bar displays the SitHub logo

**AC3: Color Consistency**
- **Given** I interact with the application
- **When** I see buttons, links, alerts, and status indicators
- **Then** they use the defined color palette consistently
- **And** success/warning/error states are visually distinct

**AC4: Dark Mode Support (Optional)**
- **Given** my system is set to dark mode
- **When** I open the application
- **Then** I see a dark theme that maintains readability and brand identity

## Technical Requirements

### Color Palette

```typescript
// Primary palette - Professional blue
primary: '#2563EB'      // Blue 600 - main actions, links, active states
primaryDark: '#1D4ED8'  // Blue 700 - hover states
primaryLight: '#3B82F6' // Blue 500 - lighter accents

// Secondary palette - Accent violet
secondary: '#7C3AED'    // Violet 600 - secondary actions, accents

// Semantic colors
success: '#059669'      // Emerald 600 - available, confirmed, success
warning: '#D97706'      // Amber 600 - warnings, pending states
error: '#DC2626'        // Red 600 - errors, destructive actions, unavailable

// Neutral palette
surface: '#F8FAFC'      // Slate 50 - page backgrounds
surfaceVariant: '#F1F5F9' // Slate 100 - card backgrounds
onSurface: '#1E293B'    // Slate 800 - primary text
onSurfaceVariant: '#64748B' // Slate 500 - secondary text
border: '#E2E8F0'       // Slate 200 - borders, dividers
```

### Typography

```typescript
// Font family
fontFamily: "'Inter', system-ui, -apple-system, sans-serif"

// Scale
h1: 2rem (32px), weight 700
h2: 1.5rem (24px), weight 600
h3: 1.25rem (20px), weight 600
h4: 1.125rem (18px), weight 600
body1: 1rem (16px), weight 400
body2: 0.875rem (14px), weight 400
caption: 0.75rem (12px), weight 500
```

### Spacing Scale

```typescript
// Base unit: 4px
xs: '4px'   // 1 unit
sm: '8px'   // 2 units
md: '16px'  // 4 units
lg: '24px'  // 6 units
xl: '32px'  // 8 units
xxl: '48px' // 12 units
```

## Tasks

### Task 1: Install Inter Font
- [ ] Add Inter font via Google Fonts or local files
- [ ] Update `index.html` with font link
- [ ] Verify font loads correctly

### Task 2: Create Vuetify Theme Configuration
- [ ] Create `web/src/plugins/vuetify.ts`
- [ ] Define light theme with color palette
- [ ] Define typography settings
- [ ] Configure default component props (density, variant, etc.)
- [ ] Export configured Vuetify instance

### Task 3: Update Main Entry Point
- [ ] Modify `web/src/main.ts` to use new Vuetify config
- [ ] Remove inline `createVuetify()` call
- [ ] Import from plugins/vuetify.ts

### Task 4: Create Global Styles
- [ ] Create `web/src/styles/global.css`
- [ ] Define CSS custom properties for design tokens
- [ ] Add base styles (body, links, focus states)
- [ ] Import in main.ts

### Task 5: Create Logo and Favicon
- [ ] Design simple text-based logo (or icon + text)
- [ ] Create `web/public/logo.svg`
- [ ] Create `web/public/favicon.svg` (or .ico)
- [ ] Update `index.html` with favicon link

### Task 6: Update App Bar with Logo
- [ ] Modify `App.vue` to display logo
- [ ] Style logo appropriately in app bar
- [ ] Ensure logo links to home

### Task 7: Verify Theme Application
- [ ] Check all existing views use new colors
- [ ] Verify buttons, alerts, chips use correct colors
- [ ] Test in different browsers

### Task 8: (Optional) Add Dark Mode
- [ ] Define dark theme colors
- [ ] Add theme toggle mechanism
- [ ] Test all views in dark mode

## File Changes

| Action | File Path |
|--------|-----------|
| Create | `web/src/plugins/vuetify.ts` |
| Create | `web/src/styles/global.css` |
| Create | `web/public/logo.svg` |
| Create | `web/public/favicon.svg` |
| Modify | `web/src/main.ts` |
| Modify | `web/index.html` |
| Modify | `web/src/App.vue` |

## Definition of Done

- [ ] Custom Vuetify theme is configured and applied
- [ ] Inter font is loaded and used throughout
- [ ] Color palette is defined and consistent
- [ ] Logo and favicon are created and displayed
- [ ] All existing views render correctly with new theme
- [ ] No Vuetify default blue (#1976D2) visible anywhere
- [ ] Code passes linting and existing tests still pass

## Notes

- Keep the design professional and clean - this is a corporate tool
- Avoid overly vibrant colors; subtle and professional is the goal
- Ensure sufficient color contrast for accessibility (WCAG AA)
- The logo can be simple text "SitHub" with a small icon, or just styled text

## Dependencies

- None (foundational story)

## Blocked By

- None

## Blocks

- Story 10.2: Reusable Component Library (needs design tokens)
- Story 10.3: Navigation & Layout Redesign (needs theme)
- All subsequent UI stories depend on this foundation
