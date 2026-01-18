#!/usr/bin/env bash
set -euo pipefail

rm -rf assets/web
mkdir -p assets/web
cp -R web/dist/. assets/web/
