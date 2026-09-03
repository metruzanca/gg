#!/usr/bin/env bash
# Race probe: how often does stock vs patched vhs fail to start.
#
# The patched fork (v0.11.2) fixes a startup race in vhs: immediately after
# the headless browser is launched, its very first navigation to the freshly
# started ttyd server can fail with net::ERR_CONNECTION_REFUSED. On this
# NixOS box that failure is essentially deterministic with the stock binary;
# elsewhere it is rare but still reproducible. This script quantifies it so
# different machines/OSes can be compared apples-to-apples.
#
# Usage:
#   ./demos/investigate/racecheck.sh [RUNS]     # RUNS = number of attempts (default 6)
#
# Env overrides:
#   VHS_STOCK  binary for the unpatched vhs   (default: the `vhs` on PATH)
#   VHS_FIXED  binary for the patched fork    (default: mise-managed vhs)
#
# deps: bash; the vhs binaries; for the patched vhs also `mise` (or set VHS_FIXED).
# Runtime deps for both: ttyd, ffmpeg and a chromium-based browser on PATH.
set -u

RUNS="${1:-${RUNS:-6}}"
REPO="$(cd "$(dirname "$0")/../.." && pwd)"
TAPE="${RACE_TAPE:-$REPO/demos/investigate/race-min.tape}"
STOCK="${VHS_STOCK:-vhs}"
FIXED="${VHS_FIXED:-}"
if [ -z "$FIXED" ] && command -v mise >/dev/null 2>&1; then
  FIXED="$(mise which vhs)"
fi

# vhs drives rod, which needs a chromium-based browser it can find via
# launcher.LookPath(). If none is on PATH it falls back to a rod-managed
# download. On NixOS the repo's chromium may only be reachable through the
# system vhs wrapper's PATH, so add it like demos/render.sh does.
if ! command -v chromium >/dev/null 2>&1; then
  for dir in /nix/store/*-chromium-*/bin; do
    if [ -x "$dir/chromium" ]; then
      export PATH="$dir:$PATH"
      break
    fi
  done
fi

# Recording writes demo-gif output; put it in a throwaway dir so repeated
# runs never clobber the repo.
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
cd "$WORK"

run_n() {
  local bin=$1
  [ -n "$bin" ] && [ -x "$(command -v "$bin" 2>/dev/null || echo "$bin")" ] || { echo "$bin: not executable, skipping"; return; }
  local i pass=0 fail=0
  for i in $(seq "$RUNS"); do
    if "$bin" "$TAPE" >"$WORK/out.log" 2>&1; then
      pass=$((pass + 1))
    else
      fail=$((fail + 1))
    fi
  done
  printf '%-60s %d/%d passes, %d failures\n' "$bin" "$pass" "$RUNS" "$fail"
}

echo "tape: $TAPE"
echo "attempts per binary: $RUNS"
echo
echo "-- stock (unpatched) vhs --"
run_n "$STOCK"
echo
echo "-- patched fork vhs --"
run_n "$FIXED"