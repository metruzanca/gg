#!/usr/bin/env bash
# Render the demo GIFs with the patched vhs (mise-managed fork build).
#
# The patched vhs fixes a startup race: without it, the very first
# navigation to the freshly started ttyd server can fail with
# net::ERR_CONNECTION_REFUSED on some machines.
#
# Rendering happens in a throwaway git repo seeded with the tag v0.1.0,
# so the app starts from v0.1.0 and the tag created during the recording
# never touches this repository.
#
# Runtime deps: ttyd, ffmpeg, a chromium-based browser on PATH, and mise.
# On NixOS chromium may only be reachable through the system vhs wrapper's
# PATH, so add the nix chromium wrapper if none is found below.
set -euo pipefail

REPO="$(cd "$(dirname "$0")/.." && pwd)"
TAPE="$REPO/demos/tag-bump.tape"
VER=v0.11.2

if ! command -v chromium >/dev/null 2>&1; then
  for dir in /nix/store/*-chromium-*/bin; do
    if [ -x "$dir/chromium" ]; then
      export PATH="$dir:$PATH"
      break
    fi
  done
fi

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
cd "$WORK"

git init -q
git config user.email vhs@example.com
git config user.name vhs
echo demo > README.md
git add .
git commit -qm init
git tag v0.1.0

mise exec "go:github.com/metruzanca/vhs@$VER" -- vhs "$TAPE"
cp "$WORK/demos/tag-bump.gif" "$REPO/demos/tag-bump.gif"
echo "rendered $REPO/demos/tag-bump.gif"