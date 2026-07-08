# Story 36.1: Persistent Sessions Across Server Restarts

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a returning user,
I want to stay logged in after the server restarts,
so that routine restarts and deployments do not force me to sign in again.

## Acceptance Criteria

1. securecookie signing keys are loaded from a persistent source (a config value or a file in
   `data_dir`) instead of being generated randomly at startup.
2. On first start with no persistent key, the server generates one, persists it (e.g. in
   `data_dir`), and reuses it on later starts.
3. A user with an active session: after a server restart, the next request succeeds without
   re-authentication; identity unchanged.
4. If the persistent key is removed/rotated, existing sessions are invalidated (documented
   behavior).

## Tasks / Subtasks

- [ ] Task 1: Add a keys package that loads-or-creates the persistent securecookie keys (AC: #1, #2, #4)
  - [ ] Create `internal/auth/keys.go` with an exported `LoadOrCreateKeys(dataDir string)
        ([]byte, []byte, error)` that returns a 32-byte hash key and 32-byte block key.
  - [ ] Persist the two keys in a single file `filepath.Join(dataDir, "cookie.key")`. On first
        start (file missing) generate both keys via `crypto/rand` (mirror the existing
        `io.ReadFull(rand.Reader, key)` pattern at `internal/auth/service.go:81-88`), write them,
        and reuse on later starts. Store as two base64 lines (or a 64-byte binary blob) so the
        file is small and parseable.
  - [ ] Write the file with `0o600` perms; if the directory does not yet exist create it with
        `os.MkdirAll(dir, 0o750)` (mirror `ensureAvatarsDir` at
        `internal/startup/server.go:335-341`). The keys file MUST NOT be world-readable.
  - [ ] Validate that a loaded file contains exactly two 32-byte keys; return a wrapped sentinel
        error (e.g. `ErrInvalidKeyFile`) on malformed/short content so a corrupt file fails
        startup loudly instead of silently regenerating (which would invalidate all sessions).
  - [ ] Add a doc comment noting the security-critical nature and the key-rotation semantics
        (per `.claude/rules/golang.md` "Document Security-Critical Functions"): deleting or
        replacing `cookie.key` invalidates all existing session cookies (AC #4).
- [ ] Task 2: Wire the persistent keys into `auth.NewService` (AC: #1, #2, #3)
  - [ ] In `internal/auth/service.go:60-98`, replace the inline random key generation
        (`internal/auth/service.go:81-92`) with a call to `LoadOrCreateKeys(cfg.Main.DataDir)`
        and pass the returned keys to `securecookie.New(hashKey, blockKey)`.
  - [ ] Wrap load failures: `return nil, fmt.Errorf("load cookie keys: %w", err)`. Remove the
        now-unused `io`/`crypto/rand` imports from `service.go` only if `NewState` (which still
        uses them at `internal/auth/service.go:351-358`) is unaffected — `NewState` keeps them,
        so leave the imports.
  - [ ] Do NOT change the `cookieCodec` field type or the `EncodeUser`/`DecodeUser`/`EncodeState`/
        `DecodeState` methods (`internal/auth/service.go:122-156`) — only the key source changes.
- [ ] Task 3: Handle the empty/`.` data_dir and test callers (AC: #2)
  - [ ] `main.data_dir` defaults to `"."` (`internal/config/config.go:98`), and existing tests
        call `auth.NewService(cfg, nil)` with an empty `Config{}` (e.g.
        `internal/auth/service_test.go:18`, `internal/startup/server_test.go:448`). Ensure
        `LoadOrCreateKeys` treats an empty `dataDir` as `"."` so those callers still succeed, and
        confirm they write `cookie.key` into a temp/`.` dir without polluting the repo (prefer
        updating those tests to pass a `t.TempDir()`-backed config if they currently rely on `.`).
  - [ ] Verify `db.Open` already assumes `dataDir` exists and is writable
        (`internal/db/db.go:14-15`) — the keys file lives alongside `sithub.db`, so no new
        directory contract is introduced for the production path.
- [ ] Task 4: Optional config override for an externally-supplied key (AC: #1)
  - [ ] If (and only if) an explicit config value is desired in addition to the auto-file, add
        `CookieKeyFile string \`mapstructure:"cookie_key_file"\`` to `MainConfig`
        (`internal/config/config.go:46-51`) with a viper default of `""`
        (`internal/config/config.go:96-108`), resolved relative to `data_dir` like the areas
        paths (`resolveAreasConfig` at `internal/config/config.go:156-190`). When set, load keys
        from that path instead of the default `cookie.key`; when empty, use the default. Keep the
        auto-generate-on-first-start behavior for the default path.
  - [ ] If the reviewer prefers the simpler data_dir-only approach, skip the config key and note
        the decision in the Dev Agent Record. The default (no new config field) is acceptable and
        satisfies AC #1 ("a config value OR a file in data_dir").
- [ ] Task 5: Documentation (AC: #4)
  - [ ] Document key-rotation/invalidation behavior in `docs/deployment-guide.md` (and, if the
        config field from Task 4 is added, in `sithub.example.toml` under `[main]` following
        `.claude/rules/toml.md`): the `cookie.key` file must be preserved across upgrades to keep
        sessions valid; deleting or replacing it logs everyone out on next request.
  - [ ] Add `/cookie.key` to `.gitignore` (alongside `/sithub.db` at `.gitignore:64-66`) so a
        dev-generated key is never committed.
- [ ] Task 6: Tests (AC: #1, #2, #3, #4)
  - [ ] Unit tests in `internal/auth/keys_test.go`: first call on an empty temp dir generates and
        writes `cookie.key` with `0o600` perms; a second call returns the SAME keys (persistence);
        a malformed/truncated file returns `ErrInvalidKeyFile`.
  - [ ] Restart-continuity test: build a `Service` from a temp `data_dir`, `EncodeUser` a user,
        construct a SECOND `Service` from the same `data_dir`, and assert `DecodeUser` on the
        first-service cookie succeeds with an unchanged identity (AC #3). Mirror the round-trip
        style at `internal/auth/service_test.go:43-` (`TestServiceUserRoundTrip`).
  - [ ] Rotation test: after removing/overwriting `cookie.key`, a freshly-built `Service` fails to
        `DecodeUser` the old cookie (AC #4).
  - [ ] Run `go test ./...`, `go vet ./...`, `gofmt`, and `golangci-lint run ./...` clean; run
        `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exit-code 1`.

## Dev Notes

Source: `_bmad-output/planning-artifacts/epics.md` — Epic 36 (`### Epic 36`, line ~1160) and
`### Story 36.1` (line ~5363). FR166 detail at `_bmad-output/planning-artifacts/epics.md:613-617`:
after a restart, previously authenticated users remain logged in; the securecookie keys "must be
persistent across restarts (loaded from config/data, not regenerated randomly per start)."

### Root cause — keys are random per start

`auth.NewService` generates fresh 32-byte hash and block keys on every startup and hands them to
`securecookie.New`:

```go
hashKey := make([]byte, 32)
blockKey := make([]byte, 32)
if _, err := io.ReadFull(rand.Reader, hashKey); err != nil { ... }
if _, err := io.ReadFull(rand.Reader, blockKey); err != nil { ... }
return &Service{ cookieCodec: securecookie.New(hashKey, blockKey), ... }
```

`internal/auth/service.go:81-92`. Because the keys change on each restart, every previously-issued
`sithub_user` cookie fails to decode after a restart, forcing re-login. This is the entire bug.

### Where the cookie codec is used

The single `*securecookie.SecureCookie` (`internal/auth/service.go:33`) signs/encrypts both the
`sithub_user` session cookie and the `sithub_oauth_state` cookie (`internal/auth/service.go:24-27`):

- `EncodeState` / `DecodeState` — `internal/auth/service.go:122-138`
- `EncodeUser` / `DecodeUser` — `internal/auth/service.go:140-156`

Cookie read/verify happens in middleware: the session cookie is looked up by name in
`internal/middleware/session.go:13` and validated via the auth service in
`internal/middleware/load_user.go` (registered at `internal/startup/server.go:80-81`). None of
these need changes — making the keys stable makes decode-after-restart work automatically.

Login writes the cookie (`internal/auth/login_local.go`, `internal/auth/handlers.go`); logout
clears it (`internal/auth/logout.go`). No signature changes required.

### Persistence location and precedent

`data_dir` (`MainConfig.DataDir`, `internal/config/config.go:49`; default `"."` at
`internal/config/config.go:98`) is the established home for persistent state:

- SQLite DB: `filepath.Join(dataDir, "sithub.db")` — `internal/db/db.go:14-15`.
- Avatars dir created with `os.MkdirAll(dir, 0o750)` — `ensureAvatarsDir` at
  `internal/startup/server.go:335-341`.

Put the key file at `filepath.Join(dataDir, "cookie.key")` and follow the avatars precedent for
directory creation. Use `0o600` for the file itself (stricter than the dir) — it holds signing
secrets and must not be world-readable. This aligns with `.claude/rules/golang.md`
("Document Security-Critical Functions" — path/secret handling).

### NewService call site and construction order

`auth.NewService(cfg, store)` is called once at `internal/startup/server.go:70`, AFTER the DB is
opened (`internal/startup/server.go:46`) and avatars dir ensured (`:65`). `cfg.Main.DataDir` is
therefore already valid and writable at that point, so `LoadOrCreateKeys(cfg.Main.DataDir)` inside
`NewService` is safe. Keep the error wrapped: `init auth service: ...` is already applied at the
call site (`internal/startup/server.go:71-72`).

### Test-caller caution (empty data_dir)

`auth.NewService` is invoked in tests with an empty/partial `Config` and `nil` store:
`internal/auth/service_test.go:18`, `internal/startup/server_test.go:448`. With no `DataDir` set,
naive code would write `cookie.key` into the process CWD (the repo). Normalize an empty `dataDir`
to `"."` in `LoadOrCreateKeys`, and prefer updating those tests to point at `t.TempDir()` so the
suite never writes a stray key file into the working tree. Add `/cookie.key` to `.gitignore`
(`.gitignore:64-66`) as a backstop.

### Key rotation / invalidation semantics (AC #4)

securecookie verifies the HMAC and decrypts using the same keys that signed the cookie. If
`cookie.key` is deleted, a fresh pair is generated on next start; if it is edited/replaced, the new
keys differ. Either way, all outstanding `sithub_user` (and `sithub_oauth_state`) cookies fail
`Decode` and users are transparently redirected to login on their next request. This is the
intended, documented behavior — call it out in `docs/deployment-guide.md`. A truncated/corrupt file
must be treated as a hard error (not a silent regenerate) so operators do not unknowingly log out
their whole user base.

### Constraints

- Go 1.25, `gofmt` defaults, wrap errors with `%w` at end of string (`.claude/rules/golang.md`).
- Define `ErrInvalidKeyFile` as a package-level sentinel (`errors.New`) per the sentinel-error
  guidance; callers may branch on it. Document `LoadOrCreateKeys` as security-critical.
- No new third-party module needed — `crypto/rand`, `encoding/base64`, `os`, `path/filepath`, and
  the already-imported `github.com/gorilla/securecookie` cover everything.

### Project Structure Notes

- New: `internal/auth/keys.go` (+ `internal/auth/keys_test.go`).
- Modified: `internal/auth/service.go` (key source only, `NewService` at `:60-98`).
- Modified (Task 4, optional): `internal/config/config.go` (`MainConfig` + viper default +
  resolution), `sithub.example.toml` (`[main]` section, lines ~3-31).
- Modified: `.gitignore` (add `/cookie.key`), `docs/deployment-guide.md` (rotation behavior).
- Naming follows package conventions: package-level store/loader funcs, no package-name repetition
  (`.claude/rules/golang.md`). Runtime artifact `cookie.key` sits next to `sithub.db` and
  `avatars/` under `data_dir` — consistent with existing layout.

### Testing standards summary

Go table-driven tests with `testify` `require`/`assert`; use `t.TempDir()` for the key file; use
`require.NoError` for setup. Cover: generate-on-first-start, reuse-on-second-load, `0o600` perms,
malformed-file error, decode-after-restart (two Services / same dir), and rotation-invalidates.
Run `go test ./...`, `go vet ./...`, `gofmt`, `golangci-lint run ./...`, and the JSCPD duplication
check. No frontend or E2E work is required for this story (backend-only). [Source:
.claude/rules/golang.md]

### References

- [Source: internal/auth/service.go:24-27,31-38,60-98,81-92,122-156,351-358]
- [Source: internal/startup/server.go:46,65,70-73,80-81,335-341]
- [Source: internal/config/config.go:46-51,96-108,156-190]
- [Source: internal/db/db.go:14-15]
- [Source: internal/middleware/session.go:13]
- [Source: internal/auth/service_test.go:18,43]
- [Source: internal/startup/server_test.go:448]
- [Source: .gitignore:64-66]
- [Source: sithub.example.toml:3-31]
- [Source: _bmad-output/planning-artifacts/epics.md:613-617,5363]
- [Source: .claude/rules/golang.md#Error Handling; #Function Documentation]
- [Source: .claude/rules/toml.md]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List

### Change Log
