#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

pushd "$script_dir/web" >/dev/null
npm ci
npm run build
popd >/dev/null

"$script_dir/tools/embed/copy.sh"

go build -o "$script_dir/sithub" ./cmd/sithub
