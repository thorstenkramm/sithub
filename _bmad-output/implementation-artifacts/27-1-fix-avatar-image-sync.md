# Story 27.1: Fix Avatar Image Sync from Entra ID

Status: done

## Story

As an Entra ID user,
I want my profile photo to sync correctly regardless of image format,
so that my avatar displays properly in SitHub.

## Acceptance Criteria

1. **Given** Microsoft Graph returns a JPEG, PNG, or other common image format
   **When** the avatar sync runs
   **Then** the image is decoded and re-encoded as PNG successfully

2. **Given** the avatar download or decoding fails
   **When** the error is logged
   **Then** the log message includes the user ID, HTTP status, content-type, and body size
   for diagnosis

3. **Given** a user has no profile photo or the sync fails
   **When** their avatar is displayed
   **Then** the fallback initials avatar is shown (no broken image)

4. **Given** an admin reads the FAQ
   **When** they look for avatar troubleshooting
   **Then** the README.md explains common causes and fixes for broken avatar sync

## Tasks / Subtasks

- [x] Task 1: Fix image decoding to handle all common formats (AC: #1)
  - [x] 1.1 In `internal/auth/avatar_handler.go`, the `SyncAvatar()` function (line ~127)
        already imports `image/jpeg` and uses `image.Decode()` which should handle both PNG
        and JPEG. The error "png: invalid format: not enough pixel data" suggests the
        `LimitReader` (line ~166) is truncating the image before full decode
  - [x] 1.2 Read the full response body into a byte buffer first (up to `maxAvatarSize`),
        then decode from the buffer — this avoids the LimitReader truncation issue
  - [x] 1.3 Add `_ "golang.org/x/image/webp"` import if WebP support is desired (optional)
  - [x] 1.4 Write a test with a real JPEG image to verify decoding works
- [x] Task 2: Improve error logging with diagnostics (AC: #2)
  - [x] 2.1 Before `image.Decode`, log the response `Content-Type` header and body size
  - [x] 2.2 On decode failure, include content-type and byte count in the error log
  - [x] 2.3 On non-200 status, include the status code (already done at line ~161)
- [x] Task 3: Ensure fallback avatar on failure (AC: #3)
  - [x] 3.1 Verify that when `SyncAvatar` fails, no corrupt `.png` file is left on disk
        (currently the file is only created after successful decode — verify this)
  - [x] 3.2 Verify `ServeAvatarHandler` returns 404 when no avatar file exists (line ~40-42)
  - [x] 3.3 Verify the frontend falls back to initials when avatar 404s
- [x] Task 4: Add FAQ section to README.md (AC: #4)
  - [x] 4.1 Add a troubleshooting section explaining: image format issues, Graph API
        permissions needed (`User.Read` with photo access), and how to re-trigger sync
- [x] Task 5: Validate (AC: #1-#4)
  - [x] 5.1 Run `go fmt ./...` and `go vet ./...`
  - [x] 5.2 Run `golangci-lint run ./...` and fix findings
  - [x] 5.3 Run `go test ./internal/auth/...` and verify no regressions
  - [x] 5.4 Run `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0`

### Review Findings

- [x] [Review][Patch] Avatar failure logs still omit required diagnostics fields [internal/auth/avatar_handler.go:161]

## Dev Notes

### Architecture & Patterns

- **Primary file**: `internal/auth/avatar_handler.go` — `SyncAvatar()` function
- **No frontend changes** for the sync fix itself (fallback already works)
- **README.md** for FAQ addition

### Root Cause Analysis

The error `png: invalid format: not enough pixel data` happens because:
1. Microsoft Graph returns the photo (possibly JPEG despite the endpoint name)
2. `io.LimitReader` wraps the body at 512KB + 1
3. `image.Decode` reads from the LimitReader
4. If the image is close to 512KB, the LimitReader may return EOF mid-decode

Fix: read the full body into `[]byte` first, check size, then decode from `bytes.Reader`.

### Key Code Locations

| Element | Location |
|---------|----------|
| `SyncAvatar()` | `internal/auth/avatar_handler.go:127` |
| `image.Decode` call | `internal/auth/avatar_handler.go:167` |
| `LimitReader` | `internal/auth/avatar_handler.go:166` |
| `ServeAvatarHandler` | `internal/auth/avatar_handler.go:26` |
| JPEG import | `internal/auth/avatar_handler.go:15` |

### Anti-Patterns to Avoid

- Do NOT remove the size limit — keep the 512KB max to prevent abuse
- Do NOT change `SyncAvatar` signature — it's called from a goroutine
- Do NOT fail the login on avatar errors — keep best-effort pattern

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List

### Change Log
