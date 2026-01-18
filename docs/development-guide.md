# Development Guide

## Scan Method
- Quick scan (pattern-based). Source files were read for this project.

## Prerequisites (from README)
- Go (backend) using Echo framework
- Vue 3 + Vuetify (frontend)
- SQLite3 (embedded database)
- Entra ID (SSO integration)

## Local Development
- The backend lives in the repo root (`go.mod`, `cmd/`, `internal/`).
- The frontend lives in `web/` and is bundled into the binary via `assets/`.
- Run backend and frontend in separate terminals during development.

## Configuration
- `sithub.example.toml` (runtime configuration)
- `sithub_areas.example.yaml` (area/room/desk definitions)
- `sithub_areas.schema.json` (schema for YAML config)

## Build / Run / Test Commands
### Backend
- Build: `go build ./cmd/sithub`
- Run: `go run ./cmd/sithub --config ./sithub.toml`
- Format: `go fmt ./...`
- Vet: `go vet ./...`
- Lint: `golangci-lint run ./...` (v2.5.0)
- Duplication check: `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exitCode 1`
- Tests: `go test -race ./...`
- Coverage: `go test -covermode=atomic -coverprofile=coverage.out ./...`
- Coverage report: `go tool cover -func=coverage.out`

### Frontend
- Dev server: `npm run dev` (in `web/`)
- Build: `npm run build` (in `web/`)
- Lint: `npm run lint` (in `web/`)
- Type check: `npm run type-check` (in `web/`)
- Unit tests: `npm run test:unit` (in `web/`)
- Unit coverage: `npm run test:unit:coverage` (in `web/`)
- E2E (headless): `npm run test:e2e` (in `web/`)
- E2E (open): `npm run test:e2e:open` (in `web/`)

### Test Notes
- Cypress E2E tests require the dev server running (`npm run dev`).
- E2E tests must use real backend responses; use intercepts for waiting/assertions only.
- CI enforces 80% coverage for Go and frontend unit tests.

## Developer Notes
- Features include desk booking, real-time availability, notifications, and admin tools.
- Multi-language support is stated (EN/DE/ES/FR), but no locale files detected here.
