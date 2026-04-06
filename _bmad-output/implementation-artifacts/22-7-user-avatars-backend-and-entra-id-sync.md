# Story 22.7: User Avatars — Backend and Entra ID Sync

Status: ready-for-dev

## Story

As a user,
I want my profile photo stored and served by SitHub,
so that colleagues can identify me visually across the app.

## Acceptance Criteria

1. **Given** a user logs in via Entra ID
   **When** the login completes
   **Then** their Microsoft Graph profile photo is downloaded and stored at
   `{data_dir}/avatars/{user_id}.png`

2. **Given** the Entra ID user has no profile photo
   **When** the login completes
   **Then** no avatar file is created and no error occurs

3. **Given** a local user opens settings
   **When** they upload a profile image
   **Then** it is stored at `{data_dir}/avatars/{user_id}.png` with max 512 KB

4. **Given** a local user has an avatar
   **When** they delete it from settings
   **Then** the file is removed from disk

5. **Given** an avatar exists for a user
   **When** any authenticated user requests `GET /api/v1/avatars/{user_id}`
   **Then** the PNG image is served with appropriate cache headers

6. **Given** no avatar exists for a user
   **When** the avatar endpoint is called
   **Then** a 404 is returned

## Tasks / Subtasks

- [ ] Task 1: Create avatars directory on startup (AC: 1, 3)
  - [ ] 1.1 In `internal/startup/server.go`: after resolving `DataDir`, create
    `{data_dir}/avatars/` directory if it doesn't exist (`os.MkdirAll`)
- [ ] Task 2: Download avatar on Entra ID login (AC: 1, 2)
  - [ ] 2.1 In `internal/auth/service.go` `FetchUser()` (line ~158): after
    fetching user info from Graph API, also call
    `https://graph.microsoft.com/v1.0/me/photo/$value` with the access token
  - [ ] 2.2 If the photo endpoint returns 200, save the image to
    `{data_dir}/avatars/{user_id}.png`. If 404 or error, skip silently
  - [ ] 2.3 Create a helper function `syncAvatar(ctx, token, userID, avatarsDir)`
    in a new file `internal/auth/avatar.go`
  - [ ] 2.4 Call `syncAvatar` from `CallbackHandler` after `FetchUser` succeeds.
    Pass the avatars directory path from config
- [ ] Task 3: Avatar upload endpoint for local users (AC: 3, 4)
  - [ ] 3.1 Create `internal/auth/avatar_handler.go` with:
    - `UploadAvatarHandler(avatarsDir string)` — `POST /api/v1/me/avatar`
      accepts multipart form upload, validates PNG/JPEG, max 512 KB, converts
      to PNG, saves to `{avatarsDir}/{user_id}.png`
    - `DeleteAvatarHandler(avatarsDir string)` — `DELETE /api/v1/me/avatar`
      removes the file
  - [ ] 3.2 Register routes in `server.go` under `requireAuth`
- [ ] Task 4: Avatar serving endpoint (AC: 5, 6)
  - [ ] 4.1 Create handler `ServeAvatarHandler(avatarsDir string)` —
    `GET /api/v1/avatars/{user_id}` serves the PNG file with
    `Content-Type: image/png` and `Cache-Control: max-age=300`
  - [ ] 4.2 Return 404 if file doesn't exist
  - [ ] 4.3 Register route in `server.go` under `requireAuth`
- [ ] Task 5: Write tests (AC: 1, 2, 3, 4, 5, 6)
  - [ ] 5.1 Test avatar sync: mock HTTP for Graph photo endpoint, verify file
    created/not-created
  - [ ] 5.2 Test upload handler: valid PNG, oversized file, invalid format
  - [ ] 5.3 Test serve handler: avatar exists (200), avatar missing (404)
  - [ ] 5.4 Test delete handler: file removed
  - [ ] 5.5 Run `go test ./...`, `golangci-lint run ./...`
- [ ] Task 6: Update API documentation (AC: 5)
  - [ ] 6.1 Add `api-doc/endpoints/avatar.yaml` for GET/POST/DELETE
  - [ ] 6.2 Add paths to `api-doc/openapi.yaml`
  - [ ] 6.3 Lint with redocly

## Dev Notes

### Microsoft Graph Photo API

```
GET https://graph.microsoft.com/v1.0/me/photo/$value
Authorization: Bearer {access_token}
```

Returns binary image data (typically JPEG). Convert to PNG for consistent
storage. Returns 404 if no photo is set.

### File Storage

```
{data_dir}/avatars/
  7c937bdb-a9ec-4a07-b6fe-346be95b8e95.png   (Anna Admin)
  830232de-a6a0-41af-b33a-c459322dab43.png   (Alex Employee)
```

Use `image/png` codec from Go stdlib. For JPEG→PNG conversion use
`image/jpeg` decode + `image/png` encode.

### Config Wiring

The `avatarsDir` path = `filepath.Join(cfg.Main.DataDir, "avatars")`.
Pass it to `registerRoutes` like `floorPlansDir` is passed today.

### Existing Patterns to Follow

- File serving: `areas.FloorPlanHandler(floorPlansDir)` in `server.go` line 157
- Auth context: `auth.GetUserFromContext(c)` for user ID
- Multipart upload: Go stdlib `c.FormFile("avatar")` with Echo

### Anti-Patterns

- Do NOT store avatars in the database (binary blobs waste SQLite space)
- Do NOT serve avatars without authentication
- Do NOT block login if avatar sync fails — it's optional
- Do NOT resize images server-side (the 512 KB limit is sufficient)

### References

- [Source: private/epic-22.md — "Use avatars"]
- [Source: internal/auth/service.go:158 — FetchUser flow]
- [Source: internal/auth/handlers.go:44 — CallbackHandler]

## Dev Agent Record

### Agent Model Used

### Completion Notes List

### File List
