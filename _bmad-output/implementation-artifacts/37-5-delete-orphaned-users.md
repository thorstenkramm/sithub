# Story 37.5: Delete Orphaned Users

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an administrator,
I want to delete users who left the organization,
so that they disappear from colleague selection and cannot sign in, while booking history stays intact.

## Acceptance Criteria

1. Migration `000004` adds `users.deleted` (default false). When the server starts, the migration applies
   cleanly and existing users are not deleted.
2. Each row in the Users table (story 37.4) has a delete icon. Clicking it opens a confirmation dialog;
   deleting my own account is not possible.
3. A user without any booking (as `user_id` or as `booked_by_user_id`, past or future): confirming the
   deletion removes the user from the database (hard delete).
4. A user with bookings: the dialog states that the user will be deactivated and how many future bookings
   (today onward) will be cancelled.
5. A user with bookings: confirming sets `deleted=true`, cancels bookings from today onward, keeps past
   bookings, and the user stays in the Users table marked as "deleted".
6. A soft-deleted user is not included in the colleague list (`GET /api/v1/colleagues`).
7. A soft-deleted local user who logs in with correct credentials is rejected with the same response as
   invalid credentials.
8. A soft-deleted user with an existing session cookie gets 401 on any request (the session check consults
   the database for local and Entra ID users).
9. A soft-deleted Entra ID user who logs in via Entra ID again gets `deleted` reset to false and is fully
   operational, including the colleague list.

FRs covered: FR185, FR186, FR187 (`_bmad-output/planning-artifacts/epics.md:685-695`).

## Tasks / Subtasks

- [x] Task 1: Migration `000004` for `users.deleted` (AC: #1)
  - [x] Add `internal/db/migrations/000004_add_users_deleted.up.sql`:
        `ALTER TABLE users ADD COLUMN deleted INTEGER NOT NULL DEFAULT 0;`
  - [x] Add `000004_add_users_deleted.down.sql`: `ALTER TABLE users DROP COLUMN deleted;` (SQLite
        DROP COLUMN is already used in `000003_add_border_width.down.sql`).
  - [x] Extend `internal/db/migrate_test.go`: after `RunMigrations`, a pre-existing user row (insert it
        after running `000003` only, or insert and re-run) has `deleted = 0`; running migrations twice is a
        no-op.
- [x] Task 2: Users store — `Deleted` field, colleague filter, Entra reactivation (AC: #1, #6, #9)
  - [x] Add `Deleted bool` to `users.Record` (`internal/users/record.go`); append `deleted` to
        `userColumns` and scan it as int in `scanRow` (same pattern as `is_admin`).
  - [x] `ListColleagues`: `SELECT id, display_name FROM users WHERE deleted = 0 ORDER BY display_name`.
  - [x] `UpsertEntraIDUser`: add `deleted = 0` to the `ON CONFLICT(email) DO UPDATE SET` clause. This is
        the reactivation path (AC #9); no other code is needed for it.
  - [x] Do NOT filter `deleted` in `FindByID`, `FindByEmail`, `FindDisplayNames`, `ListAll` (admin table
        must show deleted users; history and future export must still resolve names).
  - [x] Add `IsActive(ctx, db, id) (bool, error)`: `SELECT deleted FROM users WHERE id = ?`; a missing row
        returns `false, nil` (hard-deleted user). Single primary-key lookup, used by the session check.
  - [x] Update every hand-written `CREATE TABLE users` in tests that reads `userColumns` (see Dev Notes,
        "Test schemas") — add `deleted INTEGER NOT NULL DEFAULT 0` or switch to `db.RunMigrations`.
- [x] Task 3: Users store — deletion with booking handling in one transaction (AC: #3, #4, #5)
  - [x] Add sentinel `ErrCannotDeleteSelf` only if the check moves into the store; otherwise keep the
        handler check (see Task 4).
  - [x] Add `DeletionImpact(ctx, db, id) (Impact, error)` returning `HasBookings bool` and
        `FutureBookingCount int`:
        - `HasBookings`: `SELECT EXISTS(SELECT 1 FROM bookings WHERE user_id = ? OR booked_by_user_id = ?)`
        - `FutureBookingCount`: count with the SAME predicate as the cancellation (see next subtask).
  - [x] Add `DeleteOrDeactivate(ctx, db, id) (Result, error)` running in a single `db.BeginTx`:
        1. Ensure the user exists (else `ErrUserNotFound`).
        2. If no bookings reference the user → `DELETE FROM users WHERE id = ?`; `Result.Mode = "deleted"`.
        3. Otherwise → `UPDATE users SET deleted = 1, updated_at = ? WHERE id = ?`, then `SELECT` the
           future bookings to cancel (id, item_id, user_id, booking_date, is_guest, guest_name,
           guest_email), then `DELETE` them by id; `Result.Mode = "deactivated"`,
           `Result.CancelledBookings` holds the selected rows for notifications.
        4. Commit; on any error roll back and wrap with `%w`.
  - [x] Cancellation predicate (today = `time.Now().UTC().Format(time.DateOnly)`, matching
        `CountUserFutureBookings` in `internal/bookings/store.go:162-165`):
        `booking_date >= ? AND (user_id = ? OR (is_guest = 1 AND booked_by_user_id = ?))`.
        Bookings the deleted user made FOR another real colleague are kept (see Open questions).
  - [x] Implement the SQL directly in `internal/users` — `internal/bookings` already imports `users`, so
        `users` must not import `bookings` (import cycle).
  - [x] Remove or reimplement the old `DeleteUser` (hard delete without checks) so no caller can bypass
        the booking check; update `TestDeleteUser*` accordingly.
- [x] Task 4: Users handlers — preview endpoint, new delete semantics, `deleted` attribute (AC: #2-#5)
  - [x] Add `Deleted bool \`json:"deleted"\`` to `UserAttributes` and set it in `recordToAttributes`.
        The 37.4 list endpoint then exposes it automatically.
  - [x] New admin route `GET /api/v1/users/:id/deletion-impact` → resource type `user-deletion-impacts`,
        id = user id, attributes `has_bookings` (bool), `future_booking_count` (int), `mode`
        (`"delete"` | `"deactivate"`). 404 for unknown user.
  - [x] Rewrite `DeleteHandler` (`internal/users/handler.go:272-307`):
        - keep the self-delete guard → 400 "Cannot delete your own account" (check BEFORE any DB write);
        - REMOVE the "Only local users can be deleted" restriction (Entra ID users must be deletable,
          otherwise AC #9 is unreachable);
        - call `DeleteOrDeactivate`; hard delete → `204 No Content`; deactivation → `200` with the user
          resource (`deleted: true`) via `writeUser`;
        - already soft-deleted user → idempotent, same 200 path (cancels nothing new).
  - [x] Change signature to `DeleteHandler(store *sql.DB, notifier notifications.Notifier)` and emit one
        `notifications.EventBookingCanceled` per cancelled booking with `CanceledByUserID` = admin id
        (mirror `internal/bookings/handler.go:283-295`), so the live feed refreshes open floor plans and
        tables. Log one `slog.Info("user deactivated", ...)` with user id, admin id, cancelled count.
  - [x] Wire in `internal/startup/server.go:301` (`notifier` is already a `registerRoutes` parameter) and
        register the new `deletion-impact` route with `requireAuth, requireAdmin`.
  - [x] Use `api.WriteSingle` / error writers; terminal helper returns in non-top-level funcs need
        `//nolint:wrapcheck // Terminal response` (pattern from commit 74529fb).
- [x] Task 5: Block login and reject sessions of soft-deleted users (AC: #7, #8)
  - [x] `LocalLoginHandler` (`internal/auth/login_local.go:72-77`): AFTER `VerifyPassword` succeeds, if
        `rec.Deleted` return exactly the invalid-credentials response (401, "Invalid email or password",
        code `invalid_credentials`). Checking after bcrypt keeps timing identical (no enumeration).
        Do not update `last_login` for deleted users.
  - [x] Add `(*Service).EnsureActive(ctx, user *User) error` in `internal/auth/service.go` calling
        `users.IsActive`; return a new sentinel `ErrUserInactive` when deleted or missing. Document it as
        security-critical (golang rules). With a nil store (unit tests) return nil.
  - [x] Call it in `middleware.RequireAuth` (`internal/middleware/auth.go:19`) BEFORE
        `RefreshPermissions` → `api.WriteUnauthorized(c)` on error (applies to local and Entra ID users).
  - [x] Call it in `middleware.RedirectForbidden` (`internal/middleware/access.go:31`) → redirect to
        `/login` on error, so SPA page loads of a deleted user do not render the app shell.
  - [x] Product decision (2026-10-08): the same per-request lookup also refreshes `is_admin`. Let `IsActive`
        also return the stored `is_admin` (or load the record once), and copy it onto the context `*auth.User`
        for local and Entra ID users, so demoting an admin takes effect on the next request instead of after
        re-login. No extra query; Story 37.6 adds `is_data_exporter` to the same lookup.
  - [x] Optional hardening: on `ErrUserInactive` expire the `sithub_user` cookie (MaxAge -1) in
        `RequireAuth`. Not required by the AC.
- [x] Task 6: Server-side colleague guard (AC: #6)
  - [x] `resolveBookingParticipants` (`internal/bookings/handler.go:668-673`): treat a soft-deleted
        `for_user_id` like an unknown user (`400 for_user_id: user not found`). Hiding the user in the
        dropdown alone is not enough — the API must not accept new bookings for a deactivated user.
- [x] Task 7: API documentation (AC: #2-#5)
  - [x] `api-doc/endpoints/user.yaml` `delete`: update summary (no longer local-only), document 204 (hard
        delete) vs 200 (`UserSingleResponse`, deactivated), 400 only for self-deletion.
  - [x] New `api-doc/endpoints/user-deletion-impact.yaml`, referenced from `api-doc/openapi.yaml` paths
        (`/users/{user_id}/deletion-impact`); add a `UserDeletionImpact` schema.
  - [x] Add `deleted` (boolean) to `UserAttributes` in `api-doc/openapi.yaml` (~line 680).
  - [x] Note in the colleagues docs (if documented) that deactivated users are excluded.
  - [x] Lint: `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.
- [x] Task 8: Frontend — delete action, dialog, deleted marker (AC: #2, #4, #5)
  - [x] `web/src/api/users.ts`: add `deleted: boolean` to `UserListAttributes`; add
        `fetchUserDeletionImpact(id)` and `deleteUser(id)` (returns `null` on 204, user resource on 200).
        Add matching types to `web/src/api/types.ts` if shared types live there.
  - [x] In the Users view from story 37.4 (expected `web/src/views/UsersView.vue`), add an actions column
        with an icon button `mdi-delete` (`data-cy="user-delete-btn"`). Disable it for the own row
        (`item.id === authStore.userId`) and for rows already `deleted`; show a tooltip for the own row.
  - [x] Clicking loads `deletion-impact` and opens the existing `ConfirmDialog`
        (`web/src/components/ConfirmDialog.vue`, `confirmColor="error"`, `:loading` while deleting).
        Message: hard delete → "Delete {name} permanently?"; deactivate → "{name} has bookings and will be
        deactivated. {count} future booking(s) from today onward will be cancelled. Past bookings are kept."
        Use vue-i18n pluralization for the count.
  - [x] On confirm: call `deleteUser`, close dialog, reload the current table page (keep page/offset; if a
        hard delete empties the last page, step back one page), show a success snackbar per
        `.claude/rules/feedback.md` (`location="bottom"`, success 3000 ms, error 6000 ms + `closable`,
        `data-cy="users-feedback"`). Errors show the backend detail (story 37.1 pattern).
  - [x] Deleted rows: a small `v-chip` (`data-cy="user-deleted-chip"`, color `grey`) "deleted" next to the
        name; optionally reduced opacity for the row.
  - [x] i18n keys in ALL five locales (`web/src/locales/{en,de,es,fr,uk}.json`), e.g. `users.delete`,
        `users.deleteTitle`, `users.deleteHardMessage`, `users.deactivateMessage`, `users.deleted`,
        `users.deleteSuccess`, `users.deactivateSuccess`, `users.cannotDeleteSelf`. German: "Benutzer
        löschen", "deaktiviert", "gelöscht". Keep words short enough for mobile (FR195 / story 37.3).
- [x] Task 9: Go tests (AC: all)
  - [x] Store table tests (`internal/users/store_test.go`, real migrations): no bookings → hard delete;
        booking as `user_id` only (past) → deactivate, 0 cancelled; booking only as `booked_by_user_id`
        for a colleague → deactivate, colleague booking kept; future own + past own → future gone, past
        kept; guest booking by user in future → cancelled; today's booking → cancelled; unknown id →
        `ErrUserNotFound`; `DeletionImpact` count equals rows actually cancelled.
  - [x] `ListColleagues` excludes deleted; `UpsertEntraIDUser` on a deleted row sets `deleted=0`;
        `IsActive` true / false (deleted) / false (missing row).
  - [x] Middleware test: a local admin's cookie still says `is_admin=true`, the DB row is demoted → an admin
        route returns 403 on the next request (security regression test).
  - [x] Handler tests: self-delete 400 with no DB change; Entra user now deletable (replace
        `TestDeleteHandlerEntraIDUserRejected`); 204 vs 200 + `deleted:true`; notifier receives one
        cancel event per cancelled booking (use a recording notifier); non-admin 403 via route test;
        `deletion-impact` 200/404.
  - [x] Security regression tests (`internal/middleware`, `internal/auth`):
        - soft-deleted local user with a valid cookie → `LoadUser`+`RequireAuth` returns 401;
        - soft-deleted Entra ID user with a valid cookie → 401 (no Graph call needed to reach the 401);
        - HARD-deleted user with a still-valid cookie → 401 (today this passes — real bug fixed here);
        - active user → 200 (no regression);
        - local login of deleted user with correct password → body and status byte-identical to the
          wrong-password response;
        - `RedirectForbidden` with deleted user on `/` → 302 `/login`.
  - [x] Booking creation with `for_user_id` of a deleted user → 400.
  - [x] Run `go fmt ./...`, `go vet ./...`, `go test ./...`, `golangci-lint run ./...`, and
        `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exit-code 1` (repo gate is
        3 %, currently 1.51 %; do not add new clones).
- [x] Task 10: Frontend unit and E2E tests (AC: all)
  - [x] Vitest for the Users view: delete button disabled on own row and deleted rows; dialog shows the
        hard-delete vs deactivate text with the count; success snackbar; deleted chip rendered; API
        functions in `users.ts`. Keep coverage ≥ 75 % (`npx vitest run --coverage`).
  - [x] Cypress E2E `web/cypress/e2e/user-delete.cy.ts` against the real dev server (no response mocking;
        `cy.intercept` only for waits): seed users via admin API (`POST /api/v1/users` as Anna) with unique
        emails per run, seed bookings by logging in as the seeded user through `cy.request` and capturing
        its cookie.
        - should hard-delete a user without bookings and remove the row;
        - should show the deactivation text with the future booking count, mark the row "deleted", and
          remove the user's future bookings while past ones stay (assert via API);
        - should exclude the deactivated user from `GET /api/v1/colleagues`;
        - should reject local login of the deactivated user with the invalid-credentials response;
        - should return 401 for a request carrying the deactivated user's previously captured cookie;
        - should not offer a delete action on the admin's own row.
        AC #9 (Entra ID) cannot run in E2E (no Entra in CI) — cover it with the Go test in Task 9.
  - [x] Seeding past bookings is not possible via the API (past dates are rejected); assert "past kept" in
        Go tests and in E2E only assert that future bookings disappear.
  - [x] Run `npm run lint`, `npm run type-check`, `npx vitest run`, `npm run build`,
        `npm run test:e2e -- --browser electron`, TS jscpd, and `npm audit`.

### Review Findings

- [x] [Review][Patch] RequireAuth expired the session cookie on any EnsureActive error, including
      transient DB failures [internal/middleware/auth.go:22] — fixed: cookie is only expired on
      `errors.Is(err, auth.ErrUserInactive)`; transient errors keep the session and return 401.
- [x] [Review][Defer] Cookie name `sithub_user` hardcoded in middleware (third copy; `auth.userCookieName`
      is unexported) [internal/middleware/auth.go:24] — deferred, pre-existing duplication also in
      `internal/middleware/session.go`; export a shared constant from `internal/auth` in a follow-up.

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md` — `### Story 37.5` in "## Epic 37 Stories" (end of file) and
FR185-FR187 (`epics.md:685-695`). Product decisions are final (Epic 37 briefing): zero bookings → hard delete;
otherwise soft delete with future-booking cancellation, login block, session rejection (401), hidden from
colleagues, visible in admin table marked "deleted"; Entra ID re-login clears the flag; admins cannot delete
themselves.

### Dependency on story 37.4

37.4 builds the Users page (menu entry gated by `authStore.isAdmin`, route with `meta.requiresAdmin`,
`v-data-table-server`, `GET /api/v1/users?page[limit]&page[offset]`, `PaginationMeta`/`PaginationLinks` in
`internal/api`). This story only adds a column, a dialog and a chip to that view and does not touch
pagination. If 37.4 named things differently (view file, `data-cy` of the table), follow 37.4.

### Current state of the code (verified 2026-10-08)

- `users` table (`internal/db/migrations/000001_init.up.sql`): id, email UNIQUE, display_name, password_hash,
  user_source (`internal`|`entraid`), entra_id, is_admin, last_login, access_token, created_at, updated_at.
  Latest migration is `000003_add_border_width`; `000004` is reserved for this story (37.6 → `000005`,
  37.9 → `000006`). Migrations are embedded golang-migrate files run at startup
  (`internal/startup/server.go:57`).
- `users.DeleteUser` (`internal/users/store.go:211-225`) is a plain hard `DELETE`.
- `users.DeleteHandler` (`internal/users/handler.go:272-307`) rejects self-deletion via `currentUserID(c)`
  (uses the `GetID()` interface to avoid an `auth` import cycle) and rejects all Entra ID users. The Entra
  restriction must go.
- `users.ListColleagues` (`store.go:132-155`) returns all users → `ColleaguesHandler` → `/api/v1/colleagues`
  → `web/src/composables/useColleagues.ts` → `ColleagueSelect.vue`. Filtering in SQL is sufficient; the
  composable caches per component instance only, so a page reload picks up the change.
- `users.UpsertEntraIDUser` (`store.go:41-66`) uses `ON CONFLICT(email) DO UPDATE`; called from
  `auth.Service.FetchUser` (`internal/auth/service.go:204`) on every Entra login.
- `users.FindDisplayNames` (`store.go:306-359`) resolves booker names for history/matrix/presence; it must keep
  returning names of soft-deleted users. Hard delete only happens when no booking references the user, so no
  booking can lose its name.

### Session check — where and what it costs

- `middleware.LoadUser` (`internal/middleware/session.go`) only decodes the securecookie `sithub_user`; the
  cookie carries id, name, email, is_admin, is_permitted, auth_source and has no server-side session row.
- `middleware.RequireAuth` (`internal/middleware/auth.go:11-28`) calls `svc.RefreshPermissions`, which is a
  no-op for internal users and for Entra ID when no `users_group_id` is configured
  (`internal/auth/service.go:221-253`). Therefore a deleted (even hard-deleted) local user's cookie stays
  valid today. That is the gap AC #8 closes.
- Put the check in `RequireAuth`, not `LoadUser`: `LoadUser` runs on every request including static assets;
  `RequireAuth` guards every `/api/v1/*` route that needs a user (`server.go:247-311`), and admin routes chain
  `requireAuth, requireAdmin`. `RedirectForbidden` (`access.go`) also calls `RefreshPermissions` for SPA GETs.
- Cost: one primary-key lookup on a small SQLite table per authenticated API request (single-digit µs in WAL
  mode). When `users_group_id` is set, Entra users already trigger a Microsoft Graph round-trip per request,
  which dwarfs this. Do NOT add a cache: any TTL would delay revocation, and AC #8 requires immediate 401.
- Existing WebSocket connections (`/api/v1/live`) are authenticated at upgrade time only; an open socket of a
  deleted user is not closed. It only receives anonymized booking events (`internal/livefeed/event.go`), and
  the next reconnect gets 401. Accepted (see Open questions).

### Booking cancellation

- Cancelling a booking means `DELETE FROM bookings` (`internal/bookings/store.go:151-157`); there is no
  status column. Past bookings are untouched.
- Guest bookings store `user_id = "guest-xxxxxxxx"` and `booked_by_user_id` = host
  (`internal/bookings/handler.go:652-664`); on-behalf bookings store `user_id` = colleague and
  `booked_by_user_id` = booker.
- "Today" is the UTC date everywhere in the bookings package (`store.go:165`, `handler.go:315`); use the same.
- Emit `notifications.EventBookingCanceled` per cancelled booking via the `MultiNotifier` (webhook + live
  hub, `server.go:76-79`); the live hub anonymizes the payload (`livefeed/event.go:37-49`).
- No transactions exist in the codebase yet; `store.BeginTx(ctx, nil)` with `defer tx.Rollback()` (ignore
  `sql.ErrTxDone`) is fine with mattn/go-sqlite3. Doing the existence check, the booking check and the
  write in one transaction prevents a booking inserted between preview and delete from being orphaned.

### Login block

`LocalLoginHandler` (`internal/auth/login_local.go:35-104`) returns `invalid_credentials` for unknown email
(after a dummy bcrypt) and wrong password. Insert the `rec.Deleted` check after `VerifyPassword` so the deleted
case is timing- and byte-identical to a wrong password. Entra ID login needs no block — re-login reactivates
by design (AC #9).

### Test schemas

These tests create the `users` table by hand and will break once `userColumns` contains `deleted`:
`internal/middleware/auth_test.go`, `internal/auth/handlers_callback_success_test.go`,
`internal/auth/fetch_user_test.go`, `internal/areas/presence_handler_test.go`,
`internal/itemgroups/bookings_handler_test.go`, `internal/users/store_test.go`,
`internal/users/handler_test.go`, `internal/items/handler_test.go`. Prefer `db.RunMigrations` (as in
`internal/bookings/testhelpers_test.go:23` and `internal/floorplanpos/store_test.go:20`); otherwise add the
column. Many middleware/auth tests use `auth.NewService(cfg, nil)` — `EnsureActive` must treat a nil store as
"active" so they keep passing, and the new security tests must use a real migrated DB.

### Constraints and project rules

- Go: sentinel errors (`ErrUserInactive`), `%w` at the end, doc comments on security-critical functions
  (`EnsureActive`, `IsActive`, `DeleteOrDeactivate`), table-driven tests with testify `require`/`assert`,
  `t.Parallel()` where safe (`.claude/rules/golang.md`).
- JSON:API: attributes snake_case (`deleted`, `has_bookings`, `future_booking_count`), resource type
  kebab-case (`user-deletion-impacts`), `application/vnd.api+json` (`.claude/rules/json-api.md`).
- OpenAPI 3.1 split files, redocly lint (`.claude/rules/apidoc.md`).
- Vue: Composition API with `<script setup lang="ts">`, no `any`, Vuetify components, snackbar feedback
  (`.claude/rules/vue.md`, `.claude/rules/feedback.md`).
- Cypress: `data-cy` selectors, `should <outcome>` test names, no API mocking in E2E, no fixed waits
  (`.claude/rules/cypress.md`). E2E seed user: `anna@sithub.local` (admin) from
  `tools/database/demo-users.sql`; `cy.login()` uses `cy.session` — capture other users' cookies with
  `cy.request` instead of switching sessions.

### Project Structure Notes

- New: `internal/db/migrations/000004_add_users_deleted.{up,down}.sql`,
  `api-doc/endpoints/user-deletion-impact.yaml`, `web/cypress/e2e/user-delete.cy.ts`.
- Modified (backend): `internal/users/{record.go,store.go,handler.go}`, `internal/auth/{service.go,
  login_local.go}`, `internal/middleware/{auth.go,access.go}`, `internal/bookings/handler.go`
  (`resolveBookingParticipants`), `internal/startup/server.go` (route wiring), plus tests listed above.
- Modified (frontend): `web/src/api/users.ts`, the 37.4 Users view and its test, the five locale files.
- Modified (docs): `api-doc/openapi.yaml`, `api-doc/endpoints/user.yaml`.

### Open questions / assumptions

- Assumption: on-behalf bookings the deleted user made FOR another real colleague
  (`booked_by_user_id = X`, `user_id ≠ X`, not guest) are kept — the colleague still expects the seat.
  Guest bookings hosted by the deleted user from today onward ARE cancelled (host has left). The dialog
  count uses the same predicate. Recommended default; change only the predicate if the PO decides otherwise.
- Assumption: deletion preview is a separate admin endpoint (`GET /api/v1/users/:id/deletion-impact`)
  rather than extra attributes on every list row, to keep the paginated list query cheap.
- Assumption: hard delete returns 204, deactivation returns 200 with the user resource; deleting an already
  deactivated user is idempotent (200). The UI disables the delete icon for deleted rows.
- Assumption: Entra ID users become deletable (the old local-only restriction is removed). A hard-deleted
  Entra user is simply re-created on the next Entra login, which is acceptable.
- Assumption: no "reactivate" button for local users in this story; a soft-deleted local user stays
  deactivated (re-creation with the same email fails with 409 because of the UNIQUE index). Record as
  follow-up if needed.
- Assumption: open WebSocket connections of a deactivated user are not force-closed (see Session check).
- Assumption: no session-check cache (immediate revocation preferred over micro-optimization).

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 37.5; epics.md:681-695 (FR184-FR187)]
- [Source: internal/db/migrations/000001_init.up.sql; 000003_add_border_width.down.sql]
- [Source: internal/users/store.go:16-21,41-66,132-155,211-225,306-359,390-403]
- [Source: internal/users/handler.go:24-34,263-335; internal/users/record.go]
- [Source: internal/auth/login_local.go:35-104; internal/auth/service.go:40-51,155-253]
- [Source: internal/middleware/session.go; auth.go:11-28; access.go:12-41]
- [Source: internal/startup/server.go:57,76-82,224-311]
- [Source: internal/bookings/store.go:150-195; internal/bookings/handler.go:233-299,642-677]
- [Source: internal/notifications/notifier.go:14-45; internal/livefeed/event.go:30-55]
- [Source: web/src/api/users.ts; web/src/components/ConfirmDialog.vue; web/src/composables/useColleagues.ts]
- [Source: web/src/stores/useAuthStore.ts:5-36; web/cypress/support/commands.ts:38-63]
- [Source: api-doc/endpoints/user.yaml:75-104; api-doc/openapi.yaml UserAttributes (~680)]
- [Source: .claude/rules/golang.md; json-api.md; apidoc.md; vue.md; feedback.md; cypress.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 5.5 (claude-opus-5-5) and Claude Sonnet 5.5 (claude-sonnet-5-5); the session was interrupted by a
usage limit and resumed.

### Debug Log References

- Switching `userColumns` to include `deleted` broke tests that hand-build the `users` table, including three
  using `CREATE TABLE IF NOT EXISTS` (`areas`, `itemgroups`, `items`). The failures were silent: handlers fell back
  to "not reserved" because `FindByID` errored. All eight hand-written schemas now use the new
  `internal/testutil/dbtest.New` helper (real migrations); seed inserts needed an explicit `user_source`.
- E2E: the deactivation dialog read "...deactivated. one future booking from today onward will be cancelled"
  (lowercase sentence start from the interpolated count). Reworded in all five locales so the count follows a
  colon. `ConfirmDialog` gained `data-cy="confirm-dialog"`.
- E2E spec booked "tomorrow" and failed on a 409 when earlier failed runs had left bookings on the same
  database; it now books a random weekday 30-330 days ahead.
- Cypress occasionally failed to spawn Electron ("spawn Unknown system error -86") on an immediate rerun; a
  retry worked. Environment issue, not code.
- Final targeted Cypress failures were caused by dates that could land outside the configured test window and
  by an existing admin-users header-count assertion. The delete spec now derives future weekdays from the
  server date, and the admin-users assertion accepts the extra action column.
- `golangci-lint run --timeout=5m ./...` could not load packages locally (`context loading failed: no go
  files to analyze`) even with `--no-config`, from subdirectories, and after moving `web/node_modules` aside.
  `go vet ./...` and `go test ./... -count=1` both pass; this appears to be a local golangci-lint v2.14.0
  loader issue.
- Go jscpd at the story's strict `--threshold 0` reports the repo's existing 1.37% duplication baseline.
  The configured repo gate (`--threshold 3`) passes and this story did not raise the clone percentage.

### Completion Notes List

- Migration `000004_add_users_deleted` (up/down) with tests: existing users stay active, rerun is a no-op,
  down removes the column. `newMigrator` extracted in `internal/db/migrate.go` so tests can stop at version 3.
- `internal/users/deletion.go`: `DeleteOrDeactivate` (one transaction: hard delete if nothing references the
  user, otherwise `deleted = 1` and cancel bookings from today where the user is `user_id` or hosts a guest;
  bookings the user made for a real colleague are kept), `DeletionImpact` (same predicate as the cancellation,
  so the dialog count equals what is cancelled), `LookupAccess` (active + is_admin per request).
  `DeleteUser`/`ListAll` are gone, so nothing can bypass the booking check.
- `ListColleagues` excludes deleted users; `UpsertEntraIDUser` resets `deleted = 0` on re-login (AC 9).
- Handlers: `GET /api/v1/users/:id/deletion-impact`; `DeleteHandler(store, notifier)` returns 204 for hard
  delete and 200 with the user for deactivation, no longer rejects Entra ID users, keeps the self-delete guard
  before any write, emits one `booking.canceled` event per cancelled booking and logs `user deactivated`.
- Sessions: `Service.EnsureActive` (`ErrUserInactive`) runs in `RequireAuth` and `RedirectForbidden` for local
  and Entra ID users; hard-deleted users with a valid cookie now get 401 too (previously accepted). It also
  copies the stored `is_admin` onto the request user, so demoting an admin applies on the next request.
  Local login of a deactivated user returns the same 401 `invalid_credentials` after bcrypt verification.
  `resolveBookingParticipants` rejects a deleted `for_user_id` with 400.
- OpenAPI: `deleted` attribute, `deletion-impact` endpoint and schema, delete responses 200/204/400.
- Frontend: actions column with delete button (disabled for own and deleted rows, tooltip on own row),
  `ConfirmDialog` with hard-delete vs deactivate text (pluralized count), deleted chip, success/error
  snackbars, step back one page when a hard delete empties the last page, i18n in all five locales.
- Not implemented (per story assumptions): reactivation button for local users, force-closing open WebSockets.
- Gates passed: targeted backend suite; `go test ./... -count=1`; `go vet ./...`; Redocly lint valid
  (one pre-existing warning in `auth-providers.yaml`); `npm run lint`; `npm run type-check`;
  `npm run test:unit:coverage` (682 tests, coverage above 75%); `npm run build`; TS jscpd 0 clones;
  `npm audit --omit=dev --audit-level=high` (0 vulnerabilities); targeted delete Cypress specs; full isolated
  Cypress harness (22 specs, 102 tests).
- Gate caveats: `golangci-lint run --timeout=5m ./...` fails before analysis because the local tool cannot
  load any package; Go jscpd passes at the repo gate of 3% and fails only at the story's stricter 0% because
  of existing duplicated production Go blocks.

### File List

- internal/db/migrations/000004_add_users_deleted.up.sql (new)
- internal/db/migrations/000004_add_users_deleted.down.sql (new)
- internal/db/migrate.go
- internal/db/migrate_test.go
- internal/testutil/dbtest/dbtest.go (new)
- internal/users/record.go
- internal/users/store.go
- internal/users/deletion.go (new)
- internal/users/deletion_test.go (new)
- internal/users/handler.go
- internal/users/handler_test.go
- internal/users/store_test.go
- internal/auth/service.go
- internal/auth/login_local.go
- internal/auth/login_local_test.go
- internal/auth/fetch_user_test.go
- internal/auth/handlers_callback_success_test.go
- internal/middleware/auth.go
- internal/middleware/access.go
- internal/middleware/auth_test.go
- internal/middleware/access_test.go
- internal/bookings/handler.go
- internal/bookings/handler_test.go
- internal/startup/server.go
- internal/areas/presence_handler_test.go
- internal/itemgroups/bookings_handler_test.go
- internal/items/handler_test.go
- api-doc/openapi.yaml
- api-doc/endpoints/user.yaml
- api-doc/endpoints/users.yaml
- api-doc/endpoints/user-deletion-impact.yaml (new)
- web/src/api/types.ts
- web/src/api/users.ts
- web/src/api/users.test.ts
- web/src/views/AdminUsersView.vue
- web/src/views/AdminUsersView.test.ts
- web/src/components/ConfirmDialog.vue
- web/src/locales/{en,de,es,fr,uk}.json
- web/src/utils/date.ts
- web/src/utils/date.test.ts
- web/cypress/e2e/admin-users.cy.ts
- web/cypress/e2e/booking-errors.cy.ts
- web/cypress/e2e/user-delete.cy.ts (new)
- _bmad-output/implementation-artifacts/sprint-status.yaml

## Change Log

- 2026-10-08: Implemented story 37.5 (soft/hard delete of users, session and login blocking, Entra ID
  reactivation, admin UI with confirmation dialog). Status -> review.
- 2026-10-08: Code review (stories 37.4-37.6): fixed session cookie expiry to trigger only on
  `ErrUserInactive` (not on transient DB errors); deferred exporting a shared cookie-name constant.
  Status -> done.
- 2026-10-08: Re-review after patch fixes — clean, no new findings; gates green (go test, vet, vitest).
