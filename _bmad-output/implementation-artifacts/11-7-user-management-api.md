# Story 11.7: User Management API

Status: done

## Story

As an admin,
I want to manage local users via the API,
So that I can create, update, and remove user accounts.

## Acceptance Criteria

1. **Given** I am an admin
   **When** I list users via `GET /users`
   **Then** I see all users (Entra ID and local) with their source and role

2. **Given** I am an admin
   **When** I create a local user via `POST /users`
   **Then** the user is created with source "internal" and a hashed password

3. **Given** I am an admin
   **When** I attempt to create or delete an Entra ID user via `/users`
   **Then** the request is rejected with a JSON:API error

4. **Given** I am a non-admin user
   **When** I attempt to create, update, or delete a user via `/users`
   **Then** the request is rejected with a 403 JSON:API error

5. **Given** I am a non-admin user
   **When** I list or read users via `/users`
   **Then** I receive the data (read access is allowed for all authenticated users)

## Tasks / Subtasks

- [x] Implement GET /api/v1/users (list all users) (AC: 1, 5)
- [x] Implement POST /api/v1/users (create local user) (AC: 2, 3, 4)
  - [x] Validate required fields (email, display_name, password)
  - [x] Hash password with bcrypt
  - [x] Enforce 14-character minimum
- [x] Implement GET /api/v1/users/{id} (read single user) (AC: 5)
- [x] Implement PATCH /api/v1/users/{id} (update user) (AC: 4)
- [x] Implement DELETE /api/v1/users/{id} (delete local user) (AC: 3, 4)
  - [x] Prevent deletion of Entra ID users
- [x] Add admin-only middleware for write operations (AC: 4)
- [x] Add OpenAPI documentation for all endpoints (AC: 1-5)
- [x] Add tests for all CRUD operations (AC: 1-5)

## Dev Notes

### Design Decisions

- Read access (GET) is available to all authenticated users
- Write access (POST, PATCH, DELETE) requires admin role
- Entra ID users can be read but not created/deleted via this API
- User list returns all users regardless of source

### References

- PRD FR34: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.7: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- Full CRUD API for users with admin-only write access
- Entra ID users are read-only via this endpoint
- All endpoints return JSON:API formatted responses

### Key Files

- `internal/users/store.go` (CRUD functions)
- `internal/startup/server.go` (route registration)
- `api-doc/endpoints/users.yaml`
- `api-doc/endpoints/user.yaml`
- `api-doc/openapi.yaml` (schema definitions)

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
