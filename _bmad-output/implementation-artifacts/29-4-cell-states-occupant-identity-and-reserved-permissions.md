# Story 29.4: Cell States, Occupant Identity & Reserved Permissions

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want every cell to communicate availability and permissions immediately,
so that I can understand the board without extra clicks.

## Acceptance Criteria

1. **Given** I look at a free cell I am allowed to book
   **When** the matrix renders
   **Then** the cell looks like a normal bookable free cell with minimal text

2. **Given** I look at a free cell in a reserved room or desk that I am not allowed to book
   **When** the matrix renders
   **Then** the cell shows a lock indicator
   **And** it is not clickable

3. **Given** I look at an occupied cell
   **When** the matrix renders
   **Then** it shows the occupant using avatar plus initials in the compact cell layout
   **And** hovering reveals the full person name

4. **Given** I hover a desk label
   **When** that desk has equipment configured
   **Then** I see the equipment hints on hover

5. **Given** I am not an admin and the occupied cell is not my own booking
   **When** I interact with it
   **Then** it is non-clickable and exposes no extra popup content

## Tasks / Subtasks

- [x] Task 1: Implement final free/locked/occupied cell rendering states (AC: #1, #2, #3, #5)
  - [x] 1.1 Expand the matrix cell component from Story 29.3 to render distinct visual states for:
        - bookable free
        - locked free
        - occupied
        - muted past-day cells
  - [x] 1.2 Render reserved free cells with a subtle lock indicator and disabled interaction
  - [x] 1.3 Render occupied cells with occupant identity taking visual priority over reservation
        status so booked reserved desks still show the person
  - [x] 1.4 Keep non-admin read-only occupied cells inert; do not open any popup or menu for them

- [x] Task 2: Add occupant identity details using existing avatar infrastructure (AC: #3)
  - [x] 2.1 Use `getAvatarUrl(booker_user_id)` when a real user ID is present
  - [x] 2.2 Always show initials in the compact cell layout; if the avatar image fails or no
        `booker_user_id` exists, the initials must still carry the state
  - [x] 2.3 Derive initials from `booker_name` on the frontend; do not ask the backend to send
        precomputed initials
  - [x] 2.4 Add a tooltip containing the full `booker_name` on hover

- [x] Task 3: Add desk-label hover enrichment (AC: #4)
  - [x] 3.1 Add a tooltip to the sticky desk-name cell when `equipment` is present
  - [x] 3.2 Keep the row label itself uncluttered; equipment belongs in hover enrichment, not the
        default dense label

- [x] Task 4: Preserve clean interaction boundaries for later stories
  - [x] 4.1 Keep free bookable cells visually alive so Story 29.5 can attach booking behavior to
        them cleanly
  - [x] 4.2 Keep occupied cancel-eligible cells identifiable in component state, but do not add
        the cancellation popover until Story 29.6

- [x] Task 5: Tests and validation
  - [x] 5.1 Add component tests covering:
        - free bookable cell rendering
        - locked reserved free cell rendering and non-clickable behavior
        - occupied cell initials/avatar rendering
        - tooltip full-name rendering
        - equipment tooltip on desk label
        - non-admin occupied cell inert behavior
  - [x] 5.2 Run `cd web && npx vitest run` for the new matrix cell tests
  - [x] 5.3 Run `cd web && npm run type-check`
  - [x] 5.4 Run `cd web && npm run lint`

### Review Findings

- [x] [Review][Patch] Occupied matrix cells render short names instead of initials [web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue:40]

## Dev Notes

### Architecture & Patterns

- Story 29.1 already defined the API contract to include `booker_name`, `booker_user_id`,
  `booked_by_me`, `reserved`, `equipment`, and `warning`.
- The current app already has avatar infrastructure and an established tooltip pattern:
  use those instead of introducing a second identity system.
- Occupied reserved desks should still show the occupant. The product already corrected a similar
  priority bug in the floor plan flow.

### Key Code Locations

| Element | Location | Why it matters |
|---------|----------|----------------|
| Avatar URL helper | `web/src/api/avatars.ts` | Reuse current avatar endpoint |
| Floor plan avatar behavior | `web/src/components/InteractiveFloorPlan.vue` | Existing identity rendering reference |
| App-level avatar fallback | `web/src/App.vue` | Existing image-fallback pattern |
| Matrix components from Story 29.3 | `web/src/components/area-weekly-matrix/*` | Primary implementation surface |

### Implementation Strategy

1. Finalize the visual semantics of each cell state.
2. Add avatar + initials rendering and hover name disclosure.
3. Add equipment hover to the desk label.
4. Keep interactions narrow: bookable free cells remain available for Story 29.5; read-only
   occupied cells stay inert.

### Anti-Patterns to Avoid

- Do NOT ask the backend for precomputed initials.
- Do NOT hide occupants on reserved occupied desks.
- Do NOT open any read-only popup for non-admin occupied cells.
- Do NOT clutter the sticky desk label with inline equipment chips.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 29 Stories: Desktop Weekly Table View]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Core User Experience]
- [Source: web/src/api/avatars.ts]
- [Source: web/src/App.vue]
- [Source: web/src/components/InteractiveFloorPlan.vue]

## Dev Agent Record

### Agent Model Used

GPT-5

### Debug Log References

- Story creation only; no implementation logs yet

### Completion Notes List

- Story assumes the matrix layout from 29.3 already exists
- Occupant identity is intentionally prioritized over reserved-state decoration for occupied cells

### File List

- _bmad-output/implementation-artifacts/29-4-cell-states-occupant-identity-and-reserved-permissions.md (new story file)
