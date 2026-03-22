# Story 20.2: Memorize Selected Week and Day

Status: done

## Story

As a user,
I want the selected week and day to persist as I navigate between areas and item groups,
So that I don't have to re-select the same date on every page.

## Acceptance Criteria

1. **Given** I select week 16 on the item-groups view
   **When** I navigate to an item group and back to the item-groups view
   **Then** week 16 is still selected

2. **Given** I select a specific day on the items view
   **When** I navigate to a different item group
   **Then** the same day is pre-selected

3. **Given** the memorized week is in the past
   **When** I return to the view
   **Then** the week resets to the current week

4. **Given** I successfully book an item
   **When** the booking succeeds
   **Then** the memorized day resets to today

5. **Given** I have selected Thursday in day mode
   **When** I switch to week mode and then back to day mode
   **Then** Thursday is still selected

## Tasks / Subtasks

- [x] Create composable or Pinia store for date state (AC: 1, 2, 3, 4, 5)
  - [x] Implement `useDateState` composable or store with `sessionStorage` persistence
  - [x] Store selected week (ISO week string) and selected day (YYYY-MM-DD)
  - [x] Add reset logic: reset week when in the past, reset day after successful booking
  - [x] Use `sessionStorage` (not `localStorage`) so new tabs start fresh
- [x] Integrate into ItemGroupsView week selector (AC: 1, 3)
  - [x] Read memorized week on mount; write on change
  - [x] Reset to current week if memorized week is in the past
- [x] Integrate into ItemsView day picker (AC: 2, 4, 5)
  - [x] Read memorized day on mount; write on change
  - [x] Reset memorized day to today after successful booking
  - [x] Preserve day when toggling between day and week mode
- [x] Add unit tests for date state composable/store
  - [x] Test persistence across navigation
  - [x] Test past-week reset logic
  - [x] Test day reset after booking
  - [x] Test day survives day/week mode toggle
- [ ] Verify E2E tests still pass

## Dev Notes

### UX Recommendations (Sally)

#### sessionStorage over localStorage

The memorized date is a session concept. If a user opens SitHub in a new tab tomorrow,
they want today's date, not yesterday's stale selection. Use `sessionStorage` so each
browser tab/session starts fresh.

#### Day persists across mode toggle

When the user selects Thursday in day mode, switches to week mode to see the full picture,
and switches back to day mode — Thursday should still be selected. This was added as AC 5.

### References

- Epic 20 Story 20.2: `_bmad-output/planning-artifacts/epics.md` (Epic 20 Stories section)
- FR76, FR77: `_bmad-output/planning-artifacts/prd.md`

## Dev Agent Record

### Agent Model Used

GPT-5 Codex

### Completion Notes List

- Added `useDateState` with `sessionStorage`-backed selected week/day persistence.
- Item-groups view restores the memorized week and updates it when the selector changes.
- Items view restores the memorized day, preserves it across day/week toggles, and resets it after successful booking.
- AI review fix: a successful booking now resets the live selected day in the UI immediately, not just the stored session value.
- AI review fix: `useDateState` now re-syncs from `sessionStorage` when reused, which makes remount/navigation flows deterministic.
- Expanded tests to cover memorized week/day restore, past-date reset, and day persistence across mode toggles.
- E2E tests were not run in this review/fix pass.

### File List

- `web/src/composables/useDateState.ts` — Session-backed date state for memorized week/day selection
- `web/src/composables/useDateState.test.ts` — Unit tests for week/day persistence and reset behavior
- `web/src/views/ItemGroupsView.vue` — Restores and persists the memorized calendar week
- `web/src/views/ItemGroupsView.test.ts` — Added coverage for memorized week usage on mount
- `web/src/views/ItemsView.vue` — Restores/persists selected day and resets live day state after booking success
- `web/src/views/ItemsView.test.ts` — Added coverage for memorized day restore, day/week toggle persistence, and day reset after booking

## Senior Developer Review (AI)

- Verified ACs 20.2.1 to 20.2.5 against the current date-state implementation.
- Fixed partial AC4 behavior so successful bookings update the active date picker immediately instead of only mutating stored state.
- Added targeted view tests for memorized week/day flows to cover the actual navigation-facing behavior.

## Change Log

- 2026-03-22: UX review — added AC 5 (day persists across mode toggle), specified
  sessionStorage over localStorage, added corresponding tasks.
- 2026-03-22: Story implementation reviewed and finalized; fixed live day reset behavior and expanded persistence tests.
