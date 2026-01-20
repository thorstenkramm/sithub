#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

pushd "$script_dir/web" >/dev/null
npm ci
npm run build
popd >/dev/null

# Note: vite outputs directly to ../assets/web, no copy needed

go build -o "$script_dir/sithub" ./cmd/sithub
