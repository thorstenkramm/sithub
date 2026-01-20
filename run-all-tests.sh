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
run_step "Code duplication (frontend)" bash -c \
  "cd \"${WEB_DIR}\" && npx jscpd --pattern \"**/*.ts\" --ignore \"**/node_modules/**\" --threshold 0 --exitCode 1"
run_step "Frontend unit tests (coverage)" bash -c "cd \"${WEB_DIR}\" && npm run test:unit:coverage"

printf "\n${GREEN}🎉 All workflow tests completed successfully.${NC}\n"
