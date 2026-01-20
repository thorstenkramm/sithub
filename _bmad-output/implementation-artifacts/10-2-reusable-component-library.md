# Story 10.2: Reusable Component Library

## Story

**As a** developer,  
**I want** reusable UI components,  
**So that** the interface is consistent and maintainable.

## Status

- **Epic:** 10 - UI/UX Redesign
- **Status:** ready-for-dev
- **Priority:** High (required for view redesigns)

## Acceptance Criteria

**AC1: PageHeader Component**
- **Given** I am on any page
- **When** the page renders
- **Then** I see a consistent header with title, optional subtitle, and breadcrumbs
- **And** optional action buttons are displayed on the right

**AC2: EmptyState Component**
- **Given** a list has no items
- **When** the empty state is displayed
- **Then** I see an illustration, a message, and an optional action button
- **And** the design is visually appealing (not just plain text)

**AC3: LoadingState Component**
- **Given** data is being fetched
- **When** the loading state is displayed
- **Then** I see skeleton loaders that match the expected content layout
- **And** the loading state doesn't cause layout shift when content loads

**AC4: ConfirmDialog Component**
- **Given** I trigger a destructive action (e.g., cancel booking)
- **When** the confirmation is required
- **Then** I see a dialog with clear title, message, and confirm/cancel buttons
- **And** the confirm button is styled appropriately for the action (e.g., red for delete)

**AC5: DatePicker Component**
- **Given** I need to select a date
- **When** I click the date input
- **Then** I see a Vuetify date picker with consistent styling
- **And** the selected date is formatted consistently (YYYY-MM-DD for API, localized for display)

**AC6: StatusChip Component**
- **Given** I view a desk or booking status
- **When** the status is displayed
- **Then** I see a consistently styled chip
- **And** colors match: available (success), booked (warning), my booking (primary), unavailable (error)

## Technical Requirements

### Component Specifications

#### PageHeader.vue
```typescript
Props:
  - title: string (required)
  - subtitle?: string
  - breadcrumbs?: Array<{ text: string, to?: string }>
Slots:
  - actions: For action buttons on the right
```

#### EmptyState.vue
```typescript
Props:
  - title: string (required)
  - message?: string
  - icon?: string (mdi icon name)
  - actionText?: string
  - actionTo?: string | RouteLocationRaw
Emits:
  - action: When action button is clicked
```

#### LoadingState.vue
```typescript
Props:
  - type: 'list' | 'cards' | 'table' | 'detail' (default: 'list')
  - count?: number (number of skeleton items, default: 3)
```

#### ConfirmDialog.vue
```typescript
Props:
  - modelValue: boolean (v-model for open state)
  - title: string (required)
  - message: string (required)
  - confirmText?: string (default: 'Confirm')
  - cancelText?: string (default: 'Cancel')
  - confirmColor?: string (default: 'primary')
  - loading?: boolean
Emits:
  - update:modelValue
  - confirm
  - cancel
```

#### DatePicker.vue
```typescript
Props:
  - modelValue: string (YYYY-MM-DD format)
  - label?: string
  - min?: string (minimum selectable date)
  - max?: string (maximum selectable date)
  - disabled?: boolean
Emits:
  - update:modelValue
```

#### StatusChip.vue
```typescript
Props:
  - status: 'available' | 'booked' | 'mine' | 'unavailable' | 'guest' | 'pending'
  - label?: string (override default label)
  - size?: 'small' | 'default'
```

## Tasks

### Task 1: Create PageHeader Component
- [ ] Create `web/src/components/PageHeader.vue`
- [ ] Implement title, subtitle, breadcrumbs
- [ ] Add actions slot
- [ ] Style according to design system
- [ ] Write unit test

### Task 2: Create EmptyState Component
- [ ] Create `web/src/components/EmptyState.vue`
- [ ] Design with icon/illustration placeholder
- [ ] Implement action button
- [ ] Style with appropriate spacing and typography
- [ ] Write unit test

### Task 3: Create LoadingState Component
- [ ] Create `web/src/components/LoadingState.vue`
- [ ] Implement skeleton variants (list, cards, table)
- [ ] Use Vuetify's v-skeleton-loader
- [ ] Write unit test

### Task 4: Create ConfirmDialog Component
- [ ] Create `web/src/components/ConfirmDialog.vue`
- [ ] Implement v-dialog with v-model
- [ ] Style confirm button based on action type
- [ ] Add loading state for async confirms
- [ ] Write unit test

### Task 5: Create DatePicker Component
- [ ] Create `web/src/components/DatePicker.vue`
- [ ] Wrap Vuetify v-date-picker or v-date-input
- [ ] Handle date formatting (ISO for value, localized for display)
- [ ] Style consistently with design system
- [ ] Write unit test

### Task 6: Create StatusChip Component
- [ ] Create `web/src/components/StatusChip.vue`
- [ ] Define color mappings for each status
- [ ] Implement default labels for each status
- [ ] Write unit test

### Task 7: Create Component Index
- [ ] Create `web/src/components/index.ts`
- [ ] Export all components for easy importing
- [ ] Document usage in comments

### Task 8: Update Existing Views (Basic Integration)
- [ ] Replace one hardcoded empty state with EmptyState component
- [ ] Replace one loading state with LoadingState component
- [ ] Verify components work in real context

## File Changes

| Action | File Path |
|--------|-----------|
| Create | `web/src/components/PageHeader.vue` |
| Create | `web/src/components/EmptyState.vue` |
| Create | `web/src/components/LoadingState.vue` |
| Create | `web/src/components/ConfirmDialog.vue` |
| Create | `web/src/components/DatePicker.vue` |
| Create | `web/src/components/StatusChip.vue` |
| Create | `web/src/components/index.ts` |
| Create | `web/src/components/__tests__/PageHeader.test.ts` |
| Create | `web/src/components/__tests__/EmptyState.test.ts` |
| Create | `web/src/components/__tests__/LoadingState.test.ts` |
| Create | `web/src/components/__tests__/ConfirmDialog.test.ts` |
| Create | `web/src/components/__tests__/DatePicker.test.ts` |
| Create | `web/src/components/__tests__/StatusChip.test.ts` |

## Definition of Done

- [ ] All 6 components are created and functional
- [ ] Each component has a unit test
- [ ] Components use design system colors and typography
- [ ] Components are exported from index.ts
- [ ] At least one view uses the new components (proof of integration)
- [ ] All existing tests still pass
- [ ] Code passes linting

## Notes

- Components should be "dumb" - no API calls, just props and events
- Use Vuetify components internally where appropriate
- Ensure all components work with both light and dark themes
- Add appropriate ARIA attributes for accessibility

## Dependencies

- Story 10.1: Design System Foundation (needs theme/colors)

## Blocked By

- Story 10.1

## Blocks

- Story 10.3: Navigation & Layout Redesign (uses PageHeader)
- Story 10.4-10.6: View redesigns (use all components)
