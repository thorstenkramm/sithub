# Story 36.1: Persistent Sessions Across Server Restarts

Status: done

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

- [x] Task 1: Add a keys package that loads-or-creates the persistent securecookie keys (AC: #1, #2, #4)
  - [x] Create `internal/auth/keys.go` with an exported `LoadOrCreateKeys(dataDir string)
        ([]byte, []byte, error)` that returns a 32-byte hash key and 32-byte block key.
  - [x] Persist the two keys in a single file `filepath.Join(dataDir, "cookie.key")`. On first
        start (file missing) generate both keys via `crypto/rand`, write them, reuse on later
        starts. Stored as two base64 lines.
  - [x] Write the file with `0o600` perms; create the dir with `os.MkdirAll(dir, 0o750)`.
  - [x] Validate a loaded file contains exactly two 32-byte keys; return the sentinel
        `ErrInvalidKeyFile` on malformed/short content (corrupt file fails loudly, no silent regen).
  - [x] Doc comment covers security-critical nature + key-rotation semantics.
- [x] Task 2: Wire the persistent keys into `auth.NewService` (AC: #1, #2, #3)
  - [x] Replaced the inline random generation in `NewService` with `LoadOrCreateKeys(cfg.Main.DataDir)`.
  - [x] Load failures wrapped: `return nil, fmt.Errorf("load cookie keys: %w", err)`. `io`/`crypto/rand`
        imports kept (still used by `NewState`).
  - [x] `cookieCodec` field and the Encode/Decode methods unchanged — only the key source changed.
- [x] Task 3: Handle the empty data_dir and test callers (AC: #2)
  - [x] `LoadOrCreateKeys("")` returns ephemeral in-memory keys (no file). DECISION (deviation from
        the "normalize empty to '.'" hint): an empty `dataDir` is treated as ephemeral rather than
        `"."` — this keeps the ~40 existing `NewService(&config.Config{}, ...)` test callers writing
        NO stray `cookie.key` into the repo, while production (viper default `data_dir = "."`) always
        persists. New tests use `t.TempDir()`-backed configs for the persistent path.
  - [x] Confirmed the keys file lives alongside `sithub.db` under `data_dir` (`internal/db/db.go:14-15`);
        no new directory contract for the production path.
- [x] Task 4: Optional config override for an externally-supplied key (AC: #1) — DATA_DIR-ONLY CHOSEN
  - [ ] (config field `cookie_key_file`) — intentionally NOT added.
  - [x] Chose the simpler data_dir-only approach (auto `cookie.key`); no new config field. This
        satisfies AC #1 ("a config value OR a file in data_dir"). Decision noted in Dev Agent Record.
- [x] Task 5: Documentation (AC: #4)
  - [x] Documented key-rotation/invalidation behavior in `docs/deployment-guide.md` (new
        "Session cookie keys" subsection). No `sithub.example.toml` change (no config field added).
  - [x] Added `/cookie.key` to `.gitignore`.
- [x] Task 6: Tests (AC: #1, #2, #3, #4)
  - [x] Unit tests in `internal/auth/keys_test.go`: first call generates + writes `0o600`;
        second call returns identical keys; malformed and wrong-length files return `ErrInvalidKeyFile`;
        empty dir is ephemeral (no file).
  - [x] Restart-continuity test: two Services over the same temp `data_dir`; `DecodeUser` of the first
        service's cookie succeeds on the second with unchanged identity (AC #3).
  - [x] Rotation test: after removing `cookie.key`, a fresh Service fails to `DecodeUser` the old
        cookie (AC #4).
  - [x] `go test ./...`, `go vet ./...`, `gofmt`, `golangci-lint run ./...` clean. Go jscpd: `keys.go`
        adds 0 clones; the repo's authoritative gate is `--threshold 3` (`run-all-tests.sh`), currently
        2.7% (pre-existing), which passes.

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

claude-opus-4-8

### Debug Log References

- `go test ./internal/auth/` → ok (new keys tests + auth regression)
- `go test ./...` → exit 0 (full suite, no regressions)
- `gofmt -l`, `go vet ./internal/auth/` → clean
- `golangci-lint run ./internal/auth/...` → 0 issues (G304 addressed with a justified `#nosec`
  directive on the preceding line, matching `internal/areas/config.go` convention)
- `npx jscpd` (Go, prod): `keys.go` in 0 clones; repo baseline 2.7% is pre-existing and under the
  authoritative `--threshold 3` gate
- `npx markdownlint docs/deployment-guide.md` → clean

### Completion Notes List

- Root cause fixed: `auth.NewService` generated fresh random securecookie keys on every start, so
  every `sithub_user`/`sithub_oauth_state` cookie failed to decode after a restart. Keys are now
  persisted.
- Added `internal/auth/keys.go` — `LoadOrCreateKeys(dataDir)`: generates a 32-byte hash+block key
  pair on first start, persists them base64 to `{dataDir}/cookie.key` (`0o600`, dir `0o750`), and
  reuses them thereafter. `ErrInvalidKeyFile` sentinel makes a corrupt/short file a hard startup
  error (no silent regenerate). Security-critical doc comment covers rotation/invalidation (AC #4).
- Wired into `NewService`; codec + Encode/Decode methods untouched (only the key source changed).
- DECISION (Task 3): empty `dataDir` → ephemeral in-memory keys (not persisted). This keeps the ~40
  existing `NewService(&config.Config{}, …)` test callers from writing a stray `cookie.key` into the
  repo, while production always sets `data_dir` (viper default `"."`) and therefore always persists.
- DECISION (Task 4): took the data_dir-only path; no `cookie_key_file` config field added (the story
  marks it optional and data_dir-only satisfies AC #1). No `sithub.example.toml` change needed.
- Documented rotation/backup behavior in `docs/deployment-guide.md`; added `/cookie.key` to
  `.gitignore`.
- Verified no stray `cookie.key` is created anywhere in the repo by the full test run.

### File List

- internal/auth/keys.go (new)
- internal/auth/keys_test.go (new)
- internal/auth/service.go (modified — NewService uses LoadOrCreateKeys)
- .gitignore (modified — /cookie.key)
- docs/deployment-guide.md (modified — Session cookie keys section)

### Change Log

- 2026-07-08: Implemented FR166 — persistent securecookie signing keys (`cookie.key` in `data_dir`)
  so sessions survive server restarts; corrupt-file hard error; rotation invalidates sessions
  (documented). Backend-only.
