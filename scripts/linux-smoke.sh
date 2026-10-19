#!/usr/bin/env bash
# End-to-end Linux check on a virtual display (Xvfb, + openbox for focus):
# capture, duplicate merge, password-manager hint skipped, popup hotkey,
# Enter copy-back, capture-safe on unreadable settings. Needs Xvfb,
# python3 with libXtst; openbox for the popup legs.
#
#   scripts/linux-smoke.sh

source "$(dirname "$0")/lib.sh"
need Xvfb "apt install xvfb"
need python3
[ "$(uname -s)" = Linux ] || die "Linux only"

work="$(mktemp -d)"
pids=()
cleanup() {
	for p in "${pids[@]}"; do kill "$p" 2>/dev/null || true; done
	for p in "${pids[@]}"; do wait "$p" 2>/dev/null || true; done
	rm -rf "$work"
}
trap cleanup EXIT

step "build"
app_bin="${MYCLIP_SMOKE_BINARY:-$work/myclip}"
if [ -z "${MYCLIP_SMOKE_BINARY:-}" ]; then go build -o "$app_bin" .; fi
[ -x "$app_bin" ] || die "not executable: $app_bin"
go build -o "$work/cliphelper" ./_tools/cliphelper
ok "built"

Xvfb -displayfd 3 -screen 0 1280x800x24 3>"$work/display" >"$work/xvfb.log" 2>&1 & pids+=($!)
for _ in {1..50}; do [ ! -s "$work/display" ] || break; sleep 0.1; done
[ -s "$work/display" ] || die "Xvfb did not start"
export DISPLAY=":$(cat "$work/display")"
wm=0
if command -v openbox >/dev/null; then openbox >"$work/wm.log" 2>&1 & pids+=($!); wm=1; sleep 1; fi
export XDG_CONFIG_HOME="$work/home/.config"
cfg="$XDG_CONFIG_HOME/com.l92labs.myclip"
mkdir -p "$XDG_CONFIG_HOME"

launch() { "$app_bin" >"$work/myclip.log" 2>&1 & app=$!; pids+=($app); sleep 3; kill -0 "$app" || die "app exited during startup"; }
copies=()
copy() { "$work/cliphelper" write "$1" "${2:-none}" 2 & copies+=($!); sleep 1.5; }
settle() { wait "${copies[@]}"; copies=(); }
history() { "$work/cliphelper" dump "$cfg/history.db"; }
keys() { python3 _tools/xkeys.py "$@"; }
fail=0
# check retries for up to 5 s: the app works on its own clock.
check() {
	for _ in 1 2 3 4 5 6 7 8 9 10; do
		if eval "$2"; then ok "$1"; return; fi
		sleep 0.5
	done
	warn "FAIL: $1"; fail=1
}

step "capture"
launch
copy "hello from linux"
check "first copy captured before continuing" '[[ "$(history)" == *"hello from linux"* ]]'
copy "hunter2-secret" hint; copy "second clip"; copy "hello from linux"; settle
h="$(history)"
check "text captured" '[[ "$h" == *"\"second clip\""* ]]'
check "duplicate merged (2 copies)" 'grep -Eq $'"'"'^[0-9]+\t2\t"hello from linux"'"'"' <<<"$h"'
check "password-manager copy not stored" '[[ "$h" != *hunter2* ]]'
check "data dir 0700, db 0600" '[ "$(stat -c %a "$cfg")" = 700 ] && [ "$(stat -c %a "$cfg/history.db")" = 600 ]'

if [ "$wm" = 1 ]; then
	step "popup"
	keys super+shift+c; sleep 2.5 # the first frame draws in software here
	if [ -n "${MYCLIP_SMOKE_EVIDENCE_DIR:-}" ]; then
		need scrot
		mkdir -p "$MYCLIP_SMOKE_EVIDENCE_DIR"
		scrot -o "$MYCLIP_SMOKE_EVIDENCE_DIR/linux-popup.png"
	fi
	keys down; sleep 0.5; keys return; sleep 1
	# Whole outputs, compared afterwards: pipefail + head would SIGPIPE the reader.
	check "Enter copies the second row back" '[ "$("$work/cliphelper" read | sed -n 1p)" = "second clip" ]'
	check "copy-back bumps it to the top" '[[ "$(history | sed -n 1p)" == *"second clip"* ]]'
	"$app_bin" --toggle >/dev/null 2>&1; sleep 2 # a second launch: the popup opens
	keys down; sleep 0.5; keys return; sleep 1 # the other clip, so the check is not vacuous
	check "myclip --toggle opens the popup" '[ "$("$work/cliphelper" read | sed -n 1p)" = "hello from linux" ]'
else
	warn "openbox missing: popup legs skipped"
fi

step "capture-safe"
kill "$app"; sleep 1
printf '{' >"$cfg/settings.json"
launch
copy "copied while settings unreadable"; settle
check "app still running (capture-safe)" 'kill -0 "$app"'
check "nothing captured" '[[ "$(history)" != *unreadable* ]]'
check "settings file kept" '[ "$(cat "$cfg/settings.json")" = "{" ]'

[ "$fail" = 0 ] || die "linux smoke failed (log: see above)"
ok "linux smoke passed"
