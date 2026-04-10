# Story 25.3: Auto-Save & Remove Undo

Status: review

## Story

As an admin,
I want the floor plan editor to save automatically and not distract me with undo and
manual save controls,
so that I can focus on positioning items without worrying about losing changes.

## Acceptance Criteria

1. **Given** I draw a new rectangle on the floor plan
   **When** I release the mouse button (pointerup)
   **Then** the changes are saved automatically without clicking a save button

2. **Given** I move an existing rectangle on the floor plan
   **When** I release the mouse button (pointerup)
   **Then** the changes are saved automatically

3. **Given** I resize an existing rectangle on the floor plan
   **When** I release the mouse button (pointerup)
   **Then** the changes are saved automatically

4. **Given** no unsaved changes exist
   **When** a pointerup event fires
   **Then** no save request is triggered

5. **Given** the editor is loaded
   **When** I look at the toolbar
   **Then** there is no manual Save button

6. **Given** an auto-save is in progress
   **When** I look at the toolbar
   **Then** I see a brief saving/saved indicator reflecting the auto-save state

7. **Given** the editor is loaded
   **When** I look at the toolbar
   **Then** there is no Undo button and the undo keyboard shortcut has no effect

## Tasks / Subtasks

- [x] Task 1: Implement auto-save on pointerup (AC: #1, #2, #3, #4)
  - [x] 1.1 In `FloorPlanEditorView.vue`, locate the `onCanvasPointerUp` handler (or equivalent pointerup handler for draw, move, and resize operations)
  - [x] 1.2 After the pointer interaction completes, check `hasUnsavedChanges` computed property (`dirtyItemIDs.size > 0 || deletedPositionIDs.length > 0`)
  - [x] 1.3 If unsaved changes exist, call `saveChanges()` automatically
  - [x] 1.4 Ensure auto-save does NOT fire when no changes exist (AC #4)
  - [x] 1.5 Handle the case where a save is already in progress (debounce or skip if saving)
- [x] Task 2: Add saving/saved indicator (AC: #6)
  - [x] 2.1 Add a reactive ref `saveState` with values: `'idle' | 'saving' | 'saved'`
  - [x] 2.2 Before `saveChanges()` runs, set `saveState = 'saving'`; after completion, set `saveState = 'saved'`; after a brief timeout (e.g., 1.5s), set back to `'idle'`
  - [x] 2.3 Replace the unsaved changes chip (`data-cy="editor-unsaved-chip"`, line ~158) with a saving/saved indicator: show a spinner or "Saving..." text when saving, a checkmark or "Saved" when saved, and nothing when idle
  - [x] 2.4 Add `data-cy="editor-save-indicator"` to the new indicator element
- [x] Task 3: Remove Save button (AC: #5)
  - [x] 3.1 Remove the Save button (`data-cy="save-floor-plan-btn"`, line ~189) from the toolbar template
  - [x] 3.2 Keep the `saveChanges()` function — it's now called by auto-save
  - [x] 3.3 Remove the Ctrl+S keyboard shortcut handler if one exists for manual save
- [x] Task 4: Remove Undo button and logic (AC: #7)
  - [x] 4.1 Remove the Undo button (`data-cy="editor-undo-btn"`, line ~168) from the toolbar template
  - [x] 4.2 Remove the `undoSnapshot` ref (line ~386) and the `undoLastChange()` function (line ~761)
  - [x] 4.3 Remove all `captureUndoSnapshot()` calls throughout the component (called before draw, move, resize, delete operations)
  - [x] 4.4 Remove the Ctrl+Z keyboard event handler if one exists
  - [x] 4.5 Remove any CSS or types related to undo functionality
- [x] Task 5: Validate (AC: #1-#7)
  - [x] 5.1 Run `npm run lint` and fix findings
  - [x] 5.2 Run `npm run type-check` and fix findings
  - [x] 5.3 Run `npm run build` and verify no build errors
  - [x] 5.4 Run `npx vitest run` and verify no regressions
  - [x] 5.5 Run `npm run test:e2e -- --browser electron` and verify no regressions

## Dev Notes

### Architecture & Patterns

- **Single file change**: `web/src/views/FloorPlanEditorView.vue`
- **No backend changes**: Pure frontend behavior change
- **Existing save logic**: `saveChanges()` (line ~804) is async, handles delete + create/update API calls, clears dirty state

### Key Code Locations

| Element | Location | data-cy |
|---------|----------|---------|
| `saveChanges()` function | Line ~804 | — |
| `hasUnsavedChanges` computed | Checks `dirtyItemIDs.size > 0 \|\| deletedPositionIDs.length > 0` | — |
| Unsaved chip | Line ~158 | `editor-unsaved-chip` |
| Undo button | Line ~168 | `editor-undo-btn` |
| `undoSnapshot` ref | Line ~386 | — |
| `undoLastChange()` function | Line ~761 | — |
| Save button | Line ~189 | `save-floor-plan-btn` |
| Pointer up handlers | Canvas event handlers section | — |

### Undo Removal Scope

The undo system consists of:
- `undoSnapshot: ref<EditorSnapshot | null>` — stores a single snapshot
- `captureUndoSnapshot()` — called before mutations (draw, move, resize, delete)
- `undoLastChange()` — restores from snapshot
- Undo button in toolbar with `disabled` when `!undoSnapshot`
- Possible Ctrl+Z keyboard handler

All of these should be removed completely.

### Auto-Save Implementation Notes

- `saveChanges()` already handles the full save workflow (delete removed positions, create/update dirty positions, clear dirty state, show green flash)
- The green flash animation on saved rects (`rect-saved` class) provides visual feedback — keep this
- Add a toolbar-level indicator for saving/saved state
- Consider a small debounce (200-300ms) to avoid rapid consecutive saves during fast interactions

### Anti-Patterns to Avoid

- Do NOT create a new save function — reuse existing `saveChanges()`
- Do NOT remove `hasUnsavedChanges` — it's used to gate auto-save
- Do NOT remove `dirtyItemIDs` or `deletedPositionIDs` — they track what needs saving
- Do NOT auto-save on every pointermove — only on pointerup when changes exist

### References

- [Source: web/src/views/FloorPlanEditorView.vue — save, undo, and pointer handlers]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- ESLint: pass, TypeScript type-check: pass, Build: pass, Vitest: 308/308 pass

### Completion Notes List

- Added `autoSave()` function called at end of `onCanvasPointerUp` and `deleteByItemId` when changes exist
- Added `saveState` ref (`idle`/`saving`/`saved`) with toolbar indicator chip (`data-cy="editor-save-indicator"`)
- Saving indicator shows spinner + "Saving..." during save, checkmark + "Saved" for 1.5s after
- Removed Save button (`save-floor-plan-btn`) from toolbar
- Removed Undo button (`editor-undo-btn`) from toolbar
- Removed all undo infrastructure: `EditorSnapshot` interface, `undoSnapshot` ref, `captureUndoState()`,
  `restoreUndoState()`, `clonePositions()`, `undoLastChange()`, Ctrl+Z keyboard handler
- Removed unused i18n keys: `unsavedChanges`, `undo`, `save` from `floorPlanEditor` namespace
- Added i18n keys: `saving`, `saved` in all 5 locale files (en, de, es, fr, uk)
- Error snackbar retained for save failures

### File List

- `web/src/views/FloorPlanEditorView.vue` (modified)
- `web/src/locales/en.json` (modified)
- `web/src/locales/de.json` (modified)
- `web/src/locales/es.json` (modified)
- `web/src/locales/fr.json` (modified)
- `web/src/locales/uk.json` (modified)

### Change Log

- 2026-04-10: Implemented story 25.3 — auto-save on pointerup, removed undo and manual save
