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
cleanup() { for p in "${pids[@]}"; do kill "$p" 2>/dev/null; done; rm -rf "$work"; }
trap cleanup EXIT

step "build"
go build -o "$work/myclip" . && go build -o "$work/cliphelper" ./_tools/cliphelper
ok "built"

export DISPLAY=:$((90 + RANDOM % 9))
Xvfb "$DISPLAY" -screen 0 1280x800x24 >"$work/xvfb.log" 2>&1 & pids+=($!)
sleep 1
wm=0
if command -v openbox >/dev/null; then openbox >"$work/wm.log" 2>&1 & pids+=($!); wm=1; sleep 1; fi
export HOME="$work/home" XDG_CONFIG_HOME="$work/home/.config"
cfg="$XDG_CONFIG_HOME/com.l92labs.myclip"
mkdir -p "$HOME"

launch() { "$work/myclip" >"$work/myclip.log" 2>&1 & app=$!; pids+=($app); sleep 3; }
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
copy "hello from linux"; copy "hunter2-secret" hint; copy "second clip"; copy "hello from linux"; settle
h="$(history)"
check "text captured" '[[ "$h" == *"\"second clip\""* ]]'
check "duplicate merged (2 copies)" 'grep -q $'"'"'^1\t2\t"hello from linux"'"'"' <<<"$h"'
check "password-manager copy not stored" '[[ "$h" != *hunter2* ]]'
check "data dir 0700, db 0600" '[ "$(stat -c %a "$cfg")" = 700 ] && [ "$(stat -c %a "$cfg/history.db")" = 600 ]'

if [ "$wm" = 1 ]; then
	step "popup"
	keys super+shift+c; sleep 2.5 # the first frame draws in software here
	keys down; sleep 0.5; keys return; sleep 1
	# Whole outputs, compared afterwards: pipefail + head would SIGPIPE the reader.
	check "Enter copies the second row back" '[ "$("$work/cliphelper" read | sed -n 1p)" = "second clip" ]'
	check "copy-back bumps it to the top" '[[ "$(history | sed -n 1p)" == *"second clip"* ]]'
	"$work/myclip" --toggle >/dev/null 2>&1; sleep 2 # a second launch: the popup opens
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
