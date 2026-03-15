#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEB_DIR="${ROOT_DIR}/web"

RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
NC="\033[0m"

log_step() {
	printf "\n${YELLOW} 🔎 %s${NC}\n" "$1"
}

log_ok() {
	printf "${GREEN} ✅ %s${NC}\n" "$1"
}

log_fail() {
	printf "${RED} ❌ %s${NC}\n" "$1"
}

run_step() {
	local label="$1"
	shift
	log_step "${label}"
	if "$@"; then
		log_ok "${label} passed"
	else
		log_fail "${label} failed"
		return 1
	fi
}

cd "${ROOT_DIR}"

# Remove web/node_modules before Go steps to prevent npm packages
# that ship .go files (e.g. flatted) from being picked up by Go tooling.
# npm ci will recreate node_modules later.
rm -rf "${WEB_DIR}/node_modules"

run_step "golangci-lint config verify" golangci-lint config verify
run_step "golangci-lint run" golangci-lint run --timeout=5m ./...
run_step "Go tests (race + coverage)" bash -c \
	'go test -race -covermode=atomic -coverprofile=coverage.out ./... &&
   total=$(go tool cover -func=coverage.out | awk "/^total:/ {print \$3}") &&
   pct=${total%\%} &&
   echo "Go coverage: ${pct}%" &&
   awk -v pct="$pct" '\''BEGIN {exit !(pct >= 80)}'\'''

run_step "Markdown lint" npx markdownlint README.md
run_step "API doc lint" npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml
run_step "Code duplication (Go)" npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 3

run_step "Install frontend deps" bash -c "cd \"${WEB_DIR}\" && npm ci"
run_step "Frontend type-check" bash -c "cd \"${WEB_DIR}\" && npm run type-check"
run_step "Frontend lint" bash -c "cd \"${WEB_DIR}\" && npm run lint"
run_step "Code duplication (frontend)" bash -c \
	"cd \"${WEB_DIR}\" && npx jscpd --pattern \"**/*.ts\" --ignore \"**/node_modules/**,**/*.test.ts\" --threshold 0 --exitCode 1"
run_step "Frontend unit tests (coverage)" bash -c "cd \"${WEB_DIR}\" && npm run test:unit:coverage"
run_step "Frontend build" bash -c "cd \"${WEB_DIR}\" && npm run build"

# Security tests
run_step "Frontend NPM Audit" bash -c "cd \"${WEB_DIR}\" && npm audit"
run_step "Frontend Trivy Scan" bash -c "trivy fs --include-dev-deps --disable-telemetry ./web"
run_step "Backend Trivy Scan" bash -c "trivy fs --skip-dirs web --skip-dirs .codex --include-dev-deps ."

# Cypress E2E tests require a running server
log_step "Cypress E2E tests"

# Build the server
go build -o "${ROOT_DIR}/sithub" ./cmd/sithub

# Create temporary data dir and config files
TEST_DATA_DIR="$(mktemp -d -t sithub-e2e.XXXXXX)"
TEST_CONFIG="${ROOT_DIR}/.sithub-test.toml"
TEST_SPACES="${TEST_DATA_DIR}/areas.yaml"

cat >"${TEST_SPACES}" <<EOF
areas:
  - id: test_area
    name: "Test Area"
    items:
      - id: test_room
        name: "Test Room"
        items:
          - id: desk_1
            name: "Desk 1"
            equipment: [Monitor, Keyboard]
          - id: desk_2
            name: "Desk 2"
            equipment: [Monitor]
EOF

cat >"${TEST_CONFIG}" <<EOF
[main]
port = 8080
listen = "127.0.0.1"
data_dir = "${TEST_DATA_DIR}"

[log]
file = "-"
level = "info"

[areas]
config_file = "areas.yaml"
EOF

# Start the server in the background
"${ROOT_DIR}/sithub" run --config "${TEST_CONFIG}" &
SERVER_PID=$!

# Cleanup function to stop the server
cleanup() {
	if [ -n "${SERVER_PID:-}" ] && kill -0 "${SERVER_PID}" 2>/dev/null; then
		kill "${SERVER_PID}" 2>/dev/null || true
		wait "${SERVER_PID}" 2>/dev/null || true
	fi
	rm -f "${TEST_CONFIG}"
	rm -rf "${TEST_DATA_DIR}"
}
trap cleanup EXIT

# Wait for the server to be ready
for i in {1..30}; do
	if curl -s http://localhost:8080/health >/dev/null 2>&1; then
		break
	fi
	sleep 1
done

# Check if server started successfully
if ! curl -s http://localhost:8080/health >/dev/null 2>&1; then
	log_fail "Server failed to start"
	exit 1
fi

# Seed demo users for local authentication
sqlite3 "${TEST_DATA_DIR}/sithub.db" <"${ROOT_DIR}/tools/database/demo-users.sql"

# Run Cypress E2E tests (headless, Electron only)
if cd "${WEB_DIR}" && npx cypress run --browser electron --config baseUrl=http://localhost:8080 --env testUserEmail=anna@sithub.local,testUserPassword=SitHubDemo2026!!; then
	log_ok "Cypress E2E tests passed"
else
	log_fail "Cypress E2E tests failed"
	exit 1
fi

# Cleanup is handled by trap
cd "${ROOT_DIR}"

printf "\n${GREEN}🎉 All workflow tests completed successfully.${NC}\n"
