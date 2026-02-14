# Story 11.1: Users Table and Entra ID User Sync

Status: done

## Story

As an operator,
I want all users stored in a central users table,
So that the system has a unified user directory regardless of authentication source.

## Acceptance Criteria

1. **Given** an Entra ID user logs in for the first time
   **When** the login completes
   **Then** the user is inserted into the users table with source "entraid", name, and email

2. **Given** an Entra ID user logs in again
   **When** the login completes
   **Then** the user's name and admin status are updated from Entra ID

3. **Given** a local user is created via the API
   **When** the creation succeeds
   **Then** the user exists in the users table with source "internal"

## Tasks / Subtasks

- [x] Create migration 000008 for users table (AC: 1, 2, 3)
  - [x] Columns: id, email, display_name, password_hash, user_source, is_admin, last_login, created_at, updated_at
  - [x] UNIQUE constraint on email
  - [x] CHECK constraint on user_source IN ('internal', 'entraid')
- [x] Create `internal/users/store.go` with CRUD functions (AC: 1, 2, 3)
  - [x] FindByID, FindByEmail, Create, Update, Delete, List
  - [x] Package-level functions (not methods): `users.FindByID(ctx, db, id)`
- [x] Update Entra ID callback to upsert user on login (AC: 1, 2)
  - [x] Insert on first login, update name/admin on subsequent logins
  - [x] Update last_login timestamp
- [x] Add tests for user store functions (AC: 1, 2, 3)

## Dev Notes

### Database Schema

The users table uses `user_source` enum ('internal', 'entraid') to distinguish authentication
sources. Email is unique across all sources to prevent identity conflicts.

### Design Decisions

- Store functions are package-level (not methods) following project convention
- `auth.User` struct extended with `GetID()` method for cross-package access
- `EncodeUser` takes `*User` pointer (not value) to avoid gocritic hugeParam

### References

- PRD FR28: `_bmad-output/planning-artifacts/prd.md`
- Epic Story 11.1: `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Completion Notes

- Migration 000008 creates users table with all required columns and constraints
- `internal/users/store.go` provides FindByID, FindByEmail, Create, Update, Delete, List
- Entra ID callback handler upserts user record on each login
- All store functions tested with in-memory SQLite

### Key Files

- `migrations/000008_create_users_table.up.sql`
- `migrations/000008_create_users_table.down.sql`
- `internal/users/store.go`
- `internal/auth/service.go` (updated for user upsert)
- `internal/auth/fetch_user_test.go`

### Change Log

- 2026-02-08: Story created retroactively. Implementation was part of Epic 11 commit.
