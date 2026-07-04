# Story 35.5: Warning-Text-Keyed Dismissals

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user who dismissed a warning,
I want the warning to reappear if its text changes,
so that I never miss new or updated warning information.

## Acceptance Criteria

1. After dismissing an item's warning via "Don't show again", changing that item's warning text and
   booking the item again shows the confirmation dialog despite the earlier dismissal.
2. A stored dismissal key is derived from the item ID combined with the current warning text (a hash),
   not from the item ID alone.
3. Any pre-existing dismissals stored under an item-only scheme no longer suppress warnings (a changed
   scheme invalidates them); the dialog is shown once more, after which the new-format dismissal
   applies.
4. When the warning text is unchanged, booking a previously dismissed item keeps the confirmation
   suppressed exactly as before.

## Tasks / Subtasks

- [x] Task 1: Verify the current implementation already satisfies text-keyed dismissal (AC: #1, #2, #4)
  - [x] Confirm `useWarningSuppression` keys on `itemId::hash(warning)` and that changing the warning
        text produces a different key (so the dialog reappears)
- [x] Task 2: Add regression tests (AC: #1, #2, #3, #4)
  - [x] Unit tests for `useWarningSuppression`: dismiss then same-text → suppressed; dismiss then
        changed-text → not suppressed; key format is `itemId::<hash>`; item-only legacy keys do not
        suppress
- [ ] Task 3: Confirm behavior end-to-end through the shared confirmation (AC: #1, #4) — DEFERRED
      to 35.4 (the shared confirmation does not exist yet); verified at the composable/unit level
  - [ ] With the shared confirmation from 35.4, verify a dismissed item stays suppressed and a
        text change re-shows the dialog (component/E2E level)
- [x] Task 4 (only if a change is required): make the key format explicit
  - [x] If the team wants the md5 form named in the brief, swap `hashWarning` for md5 — but this is
        OPTIONAL; the existing hash already meets FR165 (see Dev Notes). Do NOT add a crypto
        dependency unless the team explicitly chooses md5.

## Dev Notes

Source: `private/epic-35.md` — "If the warning text changes, the warning must be shown again ...
store the md5sum of the item-id and the warning text ...".
[Source: _bmad-output/planning-artifacts/epics.md#Story 35.5 / FR165]

### 🚨 This requirement is ALREADY implemented — scope is verify + test

`web/src/composables/useWarningSuppression.ts` already keys dismissals on the item ID **and** a hash
of the warning text, and its own doc-comment states "Suppression auto-resets when the warning text
changes":

```ts
const STORAGE_KEY = 'sithub_warning_suppressed';
function hashWarning(warning: string): string { /* base-36 string hash */ }
function makeKey(itemId: string, warning: string): string {
  return `${itemId}::${hashWarning(warning)}`;
}
// isWarningSuppressed(itemId, warning) -> loadSuppressed().has(makeKey(itemId, warning))
// suppressWarning(itemId, warning)     -> add makeKey(...) to the stored Set
```

So AC #1, #2, and #4 are satisfied by the current code. The brief suggested md5 as "one option" —
the existing non-crypto string hash achieves the same behavioral goal (a text change yields a
different key). There is **no md5/crypto library** in the frontend, so switching to md5 would add a
dependency for no functional gain. Recommendation: keep the existing hash; treat this story as
locking in the behavior with tests. Only do Task 4 if the team explicitly wants md5.

### AC #3 (legacy item-only dismissals)

The historical FR104 description implied per-item-only storage, but the shipped code already uses the
item+hash key — so there is unlikely to be any real item-only data in `sithub_warning_suppressed`.
The Set is keyed uniformly by `itemId::hash`; a value stored under the old per-item form (if any ever
existed) simply won't match `makeKey(itemId, currentWarning)`, so the dialog shows again — satisfying
AC #3 without a migration. Add a test asserting an item-only key does not suppress.

### Project Structure Notes

- Primary file: `web/src/composables/useWarningSuppression.ts` (likely no code change).
- Tests: add/extend `useWarningSuppression`'s unit test (create the spec if none exists).
- No API, no store, no UI change expected. If Task 4 (md5) is chosen, that is the only code change,
  plus its dependency.

### Testing standards summary

Vitest unit tests for the composable (mock/stub localStorage via the existing `getSafeLocalStorage`
seam). Cover: same-text suppressed, changed-text not suppressed, key shape, legacy-key no-suppress,
persistence round-trip. Run type-check, lint, vitest, build. [Source: .claude/rules/vue.md]

### References

- [Source: web/src/composables/useWarningSuppression.ts (full file)]
- [Source: web/src/composables/storage.ts (getSafeLocalStorage)]
- [Source: web/src/views/ItemsView.vue:1861-1877,1927-1953 (suppression call sites)]

## Dev Agent Record

### Agent Model Used

claude-fable-5

### Debug Log References

- `npx vitest run src/composables/__tests__/useWarningSuppression.test.ts` → 5 passed
- Full suite `npx vitest run` → 451 passed; type-check, lint, build clean

### Completion Notes List

- Confirmed `useWarningSuppression` already satisfies FR165: dismissals are keyed on
  `itemId::hash(warning)`, so a warning-text change yields a new key and the dialog reappears. No
  production code change was required.
- Added `useWarningSuppression.test.ts` (5 cases): same-text suppressed; changed-text NOT suppressed;
  key format is `itemId::<hash>` (not bare item id); a legacy item-only key does not suppress;
  persistence across composable instances.
- Task 4 (md5): intentionally NOT done — the existing non-crypto hash meets the requirement and the
  frontend has no crypto/md5 dependency; adding one would be gratuitous. Left as documented decision.
- Task 3 (end-to-end through the shared confirmation dialog): DEFERRED to Story 35.4, which
  introduces that shared dialog. Verified here only at the composable/unit level.

### File List

- web/src/composables/__tests__/useWarningSuppression.test.ts (new)

### Change Log

- 2026-07-04: Verified FR165 (text-keyed dismissal already implemented) and added regression tests.
