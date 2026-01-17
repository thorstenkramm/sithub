# Cypress Testing Guidelines for a Vue 3 + TypeScript

This document adapts general Cypress best practices
for **Vue 3 + TypeScript** projects talking to a **REST backend**.

---

## Goals

- Tests are **stable** (low flakiness)
- Tests are **readable** (clear intent)
- Tests are **maintainable** (refactors don't cause mass failures)
- Tests provide **fast feedback** (component tests for logic, E2E for key flows)

All requirements outlined as "acceptance criteria" must be covered by an end-to-end (E2E) test when implementing new
features.

## Running tests

Run the dev server separately with `npm run dev` before executing Cypress E2E tests.
Before starting, stop any other dev servers that might be running, for example: `pkill -f 'npm run dev'`.

**Always** run Cypress E2E tests against a dev server. Do not mock API responses in E2E tests; use intercepts only for
waiting and assertions. Component tests may stub network responses when needed.

---

## Recommended test split (what to test where)

### Cypress Component Tests

Use for complex UI behavior without full app setup:

- booking form validation and error states
- date range selection logic
- floor plan selection interactions

### Cypress E2E Tests

Use for a small set of critical user journeys:

- Login -> book a desk -> confirm it appears in `My Bookings`
- Admin creates a room -> users can see and book desks

Rule of thumb: **Prefer component tests first.** Add E2E only when multiple parts must work together.

---

## Top 5 Dos and Don'ts (project-specific)
<!-- markdownlint-disable MD024 -->

## 1) Stable selectors

### ✅ Do

- Use a consistent attribute naming scheme (example):
  - `data-cy="foo-row"`
  - `data-cy="foo-name"`
  - `data-cy="toolbar-upload"`
  - `data-cy="dialog-rename"`
- Keep selectors **semantic**, not presentational.

### ❌ Don't

- Select by class names, DOM depth, or Tailwind classes.
- Couple selectors to how it looks instead of what it is.

**Convention recommendation:** Reserve `data-cy` only for Cypress. Avoid reusing it for styling or app logic.

---

## 2) Make network the synchronization point (avoid timing)

### ✅ Do

- Treat backend calls as the primary done signal:
  - `cy.intercept('GET', '/api/v1/foos*').as('listFoo')`
  - `cy.wait('@listFoo')`
- Assert on response bodies when useful (e.g., count, filenames) to diagnose failures faster.

### ❌ Don't

- Use fixed waits like `cy.wait(2000)` to stabilize.
- Assume rendering finishes within a timeframe.

---

## 3) Control app state (repeatable tests)

### ✅ Do

- Prefer **programmatic login** or token injection (faster, less brittle).
- Seed server data through an API endpoint if available (or a dedicated test mode).
- Reset between tests:
  - clear cookies, localStorage, and sessionStorage
  - reset backend test data if possible

### ❌ Don't

- Rely on previous tests to create data.
- Depend on execution order.

**If the backend is not in test mode:** Use a dedicated dev or staging dataset that is reset between runs. Keep E2E
tests against real responses, and use `cy.intercept()` only for observing and waiting.

---

## 4) Test user-visible behavior, not Vue internals

### ✅ Do

- Assert what a user can see:
  - the item appears in the list
  - the dialog shows a validation error
  - the output panel contains the expected result
- Assert accessibility-related states where practical:
  - focus moves correctly when using keyboard navigation
  - buttons are disabled/enabled correctly

### ❌ Don't

- Assert Pinia store contents in E2E tests.
- Assert component refs or internal computed values.

---

## 5) Keep tests small, expressive, and DRY (but not over-abstracted)

### ✅ Do

- Extract repeated sequences into:
  - `cy.login()` custom command
  - `cy.createTestFile()` (backend seed helper)
  - "page object"-style helpers sparingly (for stable screens like login)
- Keep helpers **thin** and readable.

### ❌ Don't

- Build large, generic abstractions that hide intent.
- Create helpers that perform many unrelated steps.

---

## Suggested folder structure

```text
cypress/
  e2e/
    auth.cy.ts
    booking-create.cy.ts
    booking-cancel.cy.ts
    room-browse.cy.ts
    admin-rooms.cy.ts
  component/
    DeskList.cy.ts
    FloorPlan.cy.ts
    BookingDialog.cy.ts
    DateRangePicker.cy.ts
  fixtures/
    areas.small.json
    rooms.mixed.json
    bookings.overlap.json
    bookings.empty.json
  support/
    commands.ts
    e2e.ts
```

---

## Naming and organization conventions

### Spec file naming

- E2E: `feature-action.cy.ts` (e.g., `foo-rename.cy.ts`)
- Component: component name (e.g., `RenameDialog.cy.ts`)

### Test naming

- Use: **`should <user outcome>`**
  - `should rename a file and show the new name in the list`
  - `should show an error when renaming to an existing filename`

---

## Recommended Cypress patterns for this app

## Programmatic login (concept)

- Prefer a login that avoids the UI where possible:
  - Call backend auth endpoint and set token in local storage / cookie
  - Or use a backend test token in local dev (only in test environment)

Benefits:

- Faster tests
- Less flakiness
- Avoids UI changes breaking auth tests

Keep **one** UI login E2E test (smoke test) to validate login page.

---

## Intercept patterns (REST backend)

Use intercepts with clear aliases:

- `@listAreas`
- `@listRooms`
- `@listDesks`
- `@createBooking`
- `@cancelBooking`

Assert:

- HTTP status
- request payload correctness for write operations (rename/move)
- minimal response correctness for diagnosis

---

## Fixtures strategy

- Use fixtures to cover edge cases:
  - empty areas or rooms
  - desks with missing equipment metadata
  - overlapping bookings and conflicts
  - multi-day bookings across week boundaries
  - timezone or locale-specific date formatting

Keep fixtures small and intentionally named:

- `areas.empty.json`
- `rooms.mixed.json`
- `bookings.overlap.json`

## Date and time handling

- Prefer fixed dates in fixtures and tests to avoid flakiness from "today" or timezones.
- If the UI depends on the current time, freeze it with `cy.clock()` and set a consistent timezone.


## Final guideline

> Prefer **clear intent** over clever abstractions.  
> Optimize for a future teammate understanding a failing test in 60 seconds.
