# Story 29.5: Direct Booking from Free Cells

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want to book directly from a free table cell,
so that I can act from the weekly overview without leaving the matrix.

## Acceptance Criteria

1. **Given** I click a bookable free cell
   **When** the booking UI opens
   **Then** it appears as a lightweight desktop popover anchored to that cell

2. **Given** the booking popover is open
   **When** I inspect its controls
   **Then** `Book for myself` is selected by default
   **And** I can switch to `Book for colleague`
   **And** the colleague picker appears only when `Book for colleague` is selected

3. **Given** I previously booked for a colleague from the table
   **When** I switch the booking popover to `Book for colleague` again
   **Then** the last selected colleague is preselected

4. **Given** the booking popover is open
   **When** I inspect its contents
   **Then** the note field is visible immediately
   **And** any booking warning is shown inline in the same popover instead of a second dialog

5. **Given** I enter a note and confirm a booking successfully
   **When** the request completes
   **Then** the note is stored with the created booking

6. **Given** I confirm the booking successfully
   **When** the request completes
   **Then** the popover closes
   **And** the cell updates immediately in place without navigation away from the matrix

7. **Given** another user books the same desk before my confirmation succeeds
   **When** the booking request returns a conflict
   **Then** the popover stays open
   **And** it shows an inline error explaining that the desk is no longer available

## Tasks / Subtasks

- [x] Task 1: Add an anchored booking popover to bookable free cells (AC: #1, #2)
  - [x] 1.1 Add a matrix-level booking state that tracks the active free cell and anchors a
        lightweight Vuetify popover/menu to that cell
  - [x] 1.2 Reuse the current booking-type pattern from `ItemsView.vue`:
        - `Book for myself` selected by default
        - `Book for colleague` reveals the colleague picker
  - [x] 1.3 Load colleagues through the existing `/api/v1/colleagues` client rather than
        inventing a new colleague source
  - [x] 1.4 Keep the booking UI desktop-anchored; do not fall back to a page-level navigation or
        fullscreen dialog

- [x] Task 2: Support note entry in the same booking confirmation (AC: #4, #5, #6)
  - [x] 2.1 Extend the booking create contract to accept an optional note; the current create
        payload does not already support it
  - [x] 2.2 Update:
        - `web/src/api/bookings.ts`
        - `internal/bookings/handler.go`
        - relevant booking tests
        - OpenAPI booking request schema
  - [x] 2.3 Bind the visible note field directly to that optional create payload
  - [x] 2.4 Do **not** reuse the current post-booking note dialog from `ItemsView.vue` for this
        flow; the table view requires note entry inside the initial confirmation

- [x] Task 3: Keep warnings inline inside the same popover (AC: #4)
  - [x] 3.1 Surface item warnings inline inside the popover instead of opening the existing
        warning modal from `ItemsView.vue`
  - [x] 3.2 Reuse `useWarningSuppression()` if you preserve warning-suppression behavior for this
        flow; do not create a second warning-persistence system

- [x] Task 4: Remember the last colleague choice for future colleague bookings (AC: #3)
  - [x] 4.1 Persist the last selected colleague in localStorage using `getSafeLocalStorage()`
  - [x] 4.2 Restore that colleague only when the user switches the radio to `Book for colleague`
  - [x] 4.3 Keep `Book for myself` as the default selection every time the popover opens

- [x] Task 5: Handle success and conflict outcomes in place (AC: #6, #7)
  - [x] 5.1 On success, close the popover and refresh matrix availability in place
  - [x] 5.2 On `409 Conflict`, keep the popover open, show inline error feedback, and refresh the
        affected matrix data so the user sees the new occupied state
  - [x] 5.3 Do not navigate away from the matrix after any booking outcome

- [x] Task 6: Tests and validation
  - [x] 6.1 Add frontend tests covering:
        - anchored popover open/close
        - default self-booking state
        - colleague picker reveal
        - remembered colleague preselection
        - visible note field
        - note stored with the created booking payload
        - inline warning rendering
        - success closes + refreshes
        - conflict stays open + shows inline error
  - [x] 6.2 Add backend/API tests for the create-booking contract extension for `note`
  - [x] 6.3 Run `go test ./internal/bookings/...`
- [x] 6.4 Run `cd web && npx vitest run`
- [x] 6.5 Run `cd web && npm run type-check`
- [x] 6.6 Run `cd web && npm run lint`

### Review Findings

- [x] [Review][Patch] Conflict path does not refresh the matrix after a failed booking [web/src/components/area-weekly-matrix/MatrixBookingPopover.vue:253]

## Dev Notes

### Architecture & Patterns

- This story is not purely frontend. The current booking create payload does not carry a note
  field, so note-in-popover requires a backend contract extension.
- Reuse the existing colleague API and booking-on-behalf semantics from `ItemsView.vue`.
- Reuse the existing warning-suppression composable if warning suppression remains part of the
  booking flow.

### Key Code Locations

| Element | Location | Why it matters |
|---------|----------|----------------|
| Current create booking client | `web/src/api/bookings.ts` | Extend create payload with optional `note` |
| Booking create handler | `internal/bookings/handler.go` | Accept and persist note on create |
| Existing colleague client | `web/src/api/users.ts` | Reuse `fetchColleagues()` |
| Existing booking form logic | `web/src/views/ItemsView.vue` | Reuse booking-type and colleague patterns |
| Warning suppression | `web/src/composables/useWarningSuppression.ts` | Existing persistence pattern |
| Safe localStorage | `web/src/composables/storage.ts` | Remember last colleague |

### Implementation Strategy

1. Anchor a popover to the clicked free cell.
2. Reuse colleague-booking UI patterns from `ItemsView.vue`.
3. Carry note in the same confirmation flow.
4. Keep warnings inline inside the popover.
5. Refresh matrix data in place after success or conflict.

### Anti-Patterns to Avoid

- Do NOT open a page-level modal or navigate away from the matrix.
- Do NOT use the old post-booking note dialog for this flow.
- Do NOT introduce a second colleague source or warning-persistence mechanism.
- Do NOT close the popover on a conflict response.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 29 Stories: Desktop Weekly Table View]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Core User Experience]
- [Source: web/src/api/bookings.ts]
- [Source: web/src/api/users.ts]
- [Source: web/src/views/ItemsView.vue]
- [Source: web/src/composables/useWarningSuppression.ts]
- [Source: web/src/composables/storage.ts]
- [Source: internal/bookings/handler.go]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Completion Notes List

- Extended booking create contract with optional `note` field (backend + frontend + OpenAPI)
- Anchored `v-menu` popover to free cells via `provide/inject` pattern (no event prop drilling)
- Reused `fetchColleagues` API, warning suppression composable, and `getSafeLocalStorage`
- Last-selected colleague persisted in localStorage (`sithub_matrix_last_colleague`)
- 409 conflict keeps popover open with inline error; success closes and refreshes matrix
- Inline warning rendering inside popover with suppress button

### File List

- internal/bookings/handler.go (modified: note on create)
- internal/bookings/handler_test.go (modified: 2 new tests)
- web/src/api/bookings.ts (modified: note in create payload)
- web/src/api/itemGroupMatrix.ts (unchanged)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixCell.vue (modified: click handlers)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixRow.vue (modified: pass item prop)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue (modified: popover wiring)
- web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.test.ts (modified: new tests)
- web/src/components/area-weekly-matrix/MatrixBookingPopover.vue (new)
- web/src/components/area-weekly-matrix/MatrixBookingPopover.test.ts (new)
- web/src/components/area-weekly-matrix/matrixTypes.ts (new)
- web/src/components/area-weekly-matrix/testHelpers.ts (new)
- web/src/locales/en.json (modified: matrix keys)
- web/src/locales/de.json (modified: matrix keys)
- web/src/locales/es.json (modified: matrix keys)
- web/src/locales/fr.json (modified: matrix keys)
- web/src/locales/uk.json (modified: matrix keys)
- api-doc/openapi.yaml (modified: note in CreateBookingRequestAttributes)
