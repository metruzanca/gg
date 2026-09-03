# Race investigation

`vhs` launches a headless browser and immediately navigates it to a freshly
started `ttyd` server. On some machines the browser's very first navigation
fails with `net::ERR_CONNECTION_REFUSED`, so recording aborts. The purpose of
this directory is to reproduce and quantify that across operating systems.

## What we found so far

- **NixOS (this dev box):** stock vhs fails **0/5**, patched vhs passes **5/5**
  (see `racecheck.sh` output above).
- Root cause: after `launcher.Launch()` returns, vhs connects and navigates
  immediately; a short wait resolves it. Fix lives in the
  `metruzanca/vhs` fork (branches `fix-race-condition` and `main`).
- Separately, upstream `main` has an unrelated rendering regression: `Evaluate`
  cancels the recording context *before* `Render` runs, and `Render` now passes
  that cancelled context into `exec.CommandContext(ffmpeg)`, so **no video/GIF is
  produced at all on `main`** (releases/tags are unaffected). This is worth
  filing upstream too.

## Files

| File | Purpose |
| --- | --- |
| `race-min.tape` | Minimal tape used for the probe. |
| `racecheck.sh` | Runs stock + patched vhs N times, prints pass/fail counts. |
| `Dockerfile` | Debian container that builds both vhs binaries and runs the probe. |

## Run it

```sh
# local (any OS, needs go, ttyd, ffmpeg, a chromium browser, mise)
RUNS=10 ./demos/investigate/racecheck.sh

# container (no docker/podman on the NixOS box at the moment; run elsewhere)
docker build -t vhs-racecheck vhs-racecheck -f demos/investigate/Dockerfile demos/investigate
docker run --rm -e RUNS=10 vhs-racecheck
```

Default binaries:

- **stock:** the `vhs` on `PATH` (e.g. `go install github.com/charmbracelet/vhs@v0.11.0`)
- **patched:** `mise which vhs` (the `go:github.com/metruzanca/vhs@v0.11.2` tool)

Override with `VHS_STOCK=...` and/or `VHS_FIXED=...`.

## macOS

```sh
brew install ttyd ffmpeg          # or install Chrome.app from google.com/chrome
go install github.com/charmbracelet/vhs@v0.11.0   # stock, for the probe
# then, in this repo with mise active:
go install github.com/metruzanca/vhs@v0.11.2
RUNS=10 VHS_STOCK="$(go env GOPATH)/bin/vhs" ./demos/investigate/racecheck.sh
```

Report the pass/fail numbers alongside the NixOS numbers above.