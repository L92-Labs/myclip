#!/usr/bin/env bash
# Logic gate: go vet + go test, cgo off (MyGo needs no cgo). Runs on macOS and Linux.
#
#   scripts/test.sh            vet + tests
#   scripts/test.sh -race      extra args go to go test

source "$(dirname "$0")/lib.sh"
need go "install Go from https://go.dev/dl"

step "go vet ./..."
CGO_ENABLED=0 go vet ./...
ok "vet"

step "go test ./...${*:+ $*}"
if [ "$#" -gt 0 ] && [[ " $* " == *" -race "* ]]; then
	# The race detector needs cgo.
	go test "$@" ./...
else
	CGO_ENABLED=0 go test "$@" ./...
fi
ok "tests"
