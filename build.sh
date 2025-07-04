#!/bin/sh
set -e

CGO_ENABLED=0 go build \
  -trimpath \
  -buildvcs=false \
  -ldflags="-s -w" \
  -o goweb .