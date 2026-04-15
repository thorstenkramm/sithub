# Story 29.6: In-Place Cancellation for Own & Admin-Allowed Bookings

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want eligible occupied cells to support lightweight cancellation in place,
so that I can correct bookings without leaving the weekly overview.

## Acceptance Criteria

1. **Given** I click a cell containing my own booking
   **When** the cancellation UI opens
   **Then** it uses the same lightweight anchored popover pattern as booking
   **And** it shows only person, desk, and date before the `Cancel booking` action

2. **Given** I am an admin and click a cell containing someone else's booking
   **When** the cancellation UI opens
   **Then** I can cancel that booking from the same anchored popover pattern

3. **Given** I am not an admin and the occupied cell belongs to another user
   **When** I interact with that cell
   **Then** no popup opens
   **And** the cell remains read-only

4. **Given** I confirm a cancellation successfully
   **When** the request completes
   **Then** the popover closes
   **And** the cell updates immediately in place to its new free or locked state

## Tasks / Subtasks

- [x] Task 1: Enable cancellation popovers only for eligible occupied cells (AC: #1, #2, #3)
  - [x] 1.1 Use the `booking_id` exposure rules established in Story 29.1 to determine which
        occupied cells can open the cancellation popover
  - [x] 1.2 Keep non-eligible occupied cells inert and read-only
  - [x] 1.3 Reuse the same anchored popover infrastructure introduced in Story 29.5 rather than
        creating a second overlay system

- [x] Task 2: Implement the lightweight cancellation confirmation content (AC: #1, #2)
  - [x] 2.1 Show only:
        - person
        - desk
        - date
        - cancel / close actions
  - [x] 2.2 Do not include notes, extra profile details, or a read-only info view for
        non-admin users

- [x] Task 3: Cancel bookings in place and refresh the matrix (AC: #4)
  - [x] 3.1 Call the existing `cancelBooking()` client
  - [x] 3.2 On success, close the popover and refresh the matrix data in place
  - [x] 3.3 Ensure the cell transitions to the correct post-cancel state:
        - free bookable when the desk is available to the current user
        - locked free when the desk remains reserved for the current user
  - [x] 3.4 Surface cancellation failures inline or via local matrix-scoped feedback without
        navigating away from the table

- [x] Task 4: Tests and validation
  - [x] 4.1 Add frontend tests covering:
        - own booking opens cancel popover
        - admin can open cancel popover on another user's booking
        - non-admin other-user booking stays inert
        - popover content includes only person, desk, and date
        - success closes and refreshes to free/locked state
  - [x] 4.2 Run `cd web && npx vitest run`
  - [x] 4.3 Run `cd web && npm run type-check`
  - [x] 4.4 Run `cd web && npm run lint`

## Dev Notes

### Architecture & Patterns

- This story intentionally narrows the table UX to:
  - own bookings
  - admin cancellations
- The broader backend delete authorization also allows the person who booked on behalf of someone
  else to cancel, but that does **not** have to surface in this table UX unless the epic is later
  expanded.

### Key Code Locations

| Element | Location | Why it matters |
|---------|----------|----------------|
| Cancel client | `web/src/api/bookings.ts` | Existing cancellation API |
| Delete authorization rules | `internal/bookings/handler.go#DeleteHandler` | Backend permission boundary |
| Matrix popover infrastructure | Story 29.5 component work | Reuse same anchored interaction language |

### Implementation Strategy

1. Use `booking_id` presence as the eligibility gate for the UI.
2. Reuse the anchored popover infrastructure from booking.
3. Keep the confirmation terse.
4. Refresh matrix data immediately after success so the board becomes the confirmation.

### Anti-Patterns to Avoid

- Do NOT open a popup for non-admin occupied cells belonging to another user.
- Do NOT add a second, different overlay style for cancellation.
- Do NOT navigate away from the matrix after cancellation.
- Do NOT expose extra booking details beyond person, desk, and date.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 29 Stories: Desktop Weekly Table View]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Core User Experience]
- [Source: web/src/api/bookings.ts]
- [Source: internal/bookings/handler.go#DeleteHandler]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Reused the same anchored `v-menu` popover infrastructure from Story 29.5
- Own booking + admin cells are interactive (cursor pointer, hover); others stay inert
- Cancel popover shows only person, desk, date + cancel action (terse by design)
- Success closes popover and refreshes matrix; failure shows inline error
- Non-admin other-user cells remain read-only with no popup

### File List

- web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue (modified: occupied click)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue (modified: cancel popover)
- web/src/components/area-weekly-matrix/MatrixCancelPopover.vue (new)
- web/src/components/area-weekly-matrix/MatrixCancelPopover.test.ts (new)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts (modified: new tests)

