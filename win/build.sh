#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p assets dist
if [[ ! -f assets/AppIcon.png ]]; then
  sips -z 128 128 ../macapp/AppIcon.png --out assets/AppIcon.png >/dev/null
fi
export GOOS=windows GOARCH=amd64 CGO_ENABLED=0
go build -ldflags "-H windowsgui -s -w" -o "dist/搜魂.exe" .
echo "built dist/搜魂.exe"
