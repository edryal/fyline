#!/usr/bin/env bash

# runs the server + 2 clients (Alice, Bob) with debug logging
set -e
cd "$(dirname "$0")/.."

# Ctrl-C to kill everything
# kill all background jobs when this script exits
trap 'kill 0' EXIT

FYLINE_LOG=debug go run ./cmd/server &
sleep 1 # give the server a moment to start

FYLINE_LOG=debug FYLINE_USER=Alice go run ./cmd/client 2>&1 | sed 's/^/[Alice] /' &
FYLINE_LOG=debug FYLINE_USER=Bob go run ./cmd/client 2>&1 | sed 's/^/[Bob] /' &

wait
