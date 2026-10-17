#!/usr/bin/env bash
# Development run on macOS: live reload, ad-hoc signed dev app.
#
#   scripts/dev-mac.sh                 run
#   scripts/dev-mac.sh -lang fr        extra args go to mygo dev

source "$(dirname "$0")/lib.sh"
[ "$(uname -s)" = Darwin ] || die "mygo dev needs macOS (use scripts/test.sh here)"
need go "install Go from https://go.dev/dl"

exec go tool mygo dev "$@"
