# Cross-platform reproduction of the vhs startup race.
#
# Build:
#   docker build -t vhs-racecheck -f demos/investigate/Dockerfile demos/investigate
# Run:
#   docker run --rm -e RUNS=10 vhs-racecheck
#
# The script runs the STOCK upstream vhs (v0.11.0) and the PATCHED fork
# (v0.11.2) over a tiny tape and counts how often each one fails to open the
# ttyd page. On the NixOS dev box the stock binary fails ~always and the
# patched one ~never; running this container reports the numbers on a stock
# Debian setup so the two environments can be compared.
#
# NOTE: this uses Debian (same as upstream vhs's own Dockerfile) because
# Ubuntu ships chromium as a snap stub that cannot run in containers.

FROM tsl0922/ttyd:alpine AS ttyd

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git ffmpeg chromium \
  && rm -rf /var/lib/apt/lists/*

COPY --from=ttyd /usr/bin/ttyd /usr/bin/ttyd

# vhs v0.11.0 requires go >= 1.25.8; Debian's golang-go is older, so fetch a
# pinned toolchain.
ARG GO_VERSION=1.26.6
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
      | tar -C /usr/local -xzf - \
 && export PATH=/usr/local/go/bin:$PATH \
 && go install github.com/charmbracelet/vhs@v0.11.0 \
 && cp "$(go env GOPATH)/bin/vhs" /usr/local/bin/vhs-stock \
 && go install github.com/metruzanca/vhs@v0.11.2 \
 && cp "$(go env GOPATH)/bin/vhs" /usr/local/bin/vhs-patched

COPY race-min.tape /race/race-min.tape
COPY racecheck.sh /race/racecheck.sh

ENV VHS_STOCK=vhs-stock \
    VHS_FIXED=vhs-patched \
    RACE_TAPE=/race/race-min.tape \
    VHS_NO_SANDBOX=1

CMD ["/race/racecheck.sh"]