#!/usr/bin/env bash
# Build MyClip for macOS: preflight -> vet + tests -> mygo build -> codesign report.
#
#   scripts/build-mac.sh                    .app for this Mac's arch (fast, no DMG)
#   scripts/build-mac.sh --dmg              .app + .dmg
#   scripts/build-mac.sh --universal --dmg  arm64 + x86_64 in one .app/.dmg
#   scripts/build-mac.sh --sign "MyClip Dev"   identity (default: mygo.json macos.signingIdentity, else ad hoc)
#   scripts/build-mac.sh --open             launch the built app
#   scripts/build-mac.sh --skip-tests       skip vet + tests
#   scripts/build-mac.sh --cross            allow off-Mac build (unsigned, no DMG; compile check only)
#
# Output: the "out" dir of mygo.json.

source "$(dirname "$0")/lib.sh"

arch="" dmg=0 sign="" open_app=0 skip_tests=0 cross=0
while [ "$#" -gt 0 ]; do
	case "$1" in
	--dmg) dmg=1 ;;
	--universal) arch=universal ;;
	--arch) arch="${2:?--arch needs arm64|amd64|universal}"; shift ;;
	--sign) sign="${2:?--sign needs an identity}"; shift ;;
	--open) open_app=1 ;;
	--skip-tests) skip_tests=1 ;;
	--cross) cross=1 ;;
	-h|--help) sed -n '2,13p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
	*) die "unknown flag $1 (see --help)" ;;
	esac
	shift
done

on_mac=0
[ "$(uname -s)" = Darwin ] && on_mac=1

# --- preflight ---------------------------------------------------------------
step "preflight"
need go "install Go from https://go.dev/dl (go.mod pins $(sed -nE 's/^go //p' go.mod))"
need python3
python3 scripts/licenses.py --check
ok "$(go version)"
if [ "$on_mac" = 1 ]; then
	xcode-select -p >/dev/null 2>&1 || die "Xcode Command Line Tools missing — run: xcode-select --install"
	need codesign "comes with Xcode Command Line Tools"
	ok "macOS $(sw_vers -productVersion) ($(uname -m)), CLT at $(xcode-select -p)"
	go tool mygo doctor
elif [ "$cross" = 1 ]; then
	warn "not on macOS: building unsigned, no DMG — compile check only, not a runnable gate artifact"
	[ "$dmg" = 1 ] && die "--dmg needs macOS"
	[ -n "$sign" ] && die "--sign needs macOS"
else
	die "run this on a Mac (or pass --cross for an unsigned compile check)"
fi

if [ -z "$arch" ]; then
	case "$(uname -m)" in
	arm64|aarch64) arch=arm64 ;;
	x86_64|amd64) arch=amd64 ;;
	*) die "unknown CPU $(uname -m); pass --arch" ;;
	esac
fi
case "$arch" in arm64|amd64|universal) ;; *) die "--arch must be arm64, amd64 or universal" ;; esac

# --- vet + tests -------------------------------------------------------------
if [ "$skip_tests" = 0 ]; then
	"$ROOT/scripts/test.sh"
else
	warn "skipping vet + tests"
fi

# --- build -------------------------------------------------------------------
out="$(out_dir)"
args=(-platform "darwin/$arch")
[ "$dmg" = 0 ] && args+=(-skip-dmg)
[ -n "$sign" ] && args+=(-sign "$sign")

step "go tool mygo build ${args[*]}"
go tool mygo build "${args[@]}"

# mygo build writes <out>/darwin-<arch>/<name>.app.
app="$(find "$out/darwin-$arch" -maxdepth 1 -type d -name '*.app' -print 2>/dev/null | head -1)"
[ -n "$app" ] || die "no .app found under $out/darwin-$arch/"

# --- report ------------------------------------------------------------------
step "artifacts"
ok "app: $app ($(du -sh "$app" | cut -f1))"
exe="$app/Contents/MacOS/$(sed -nE 's/^[[:space:]]*"name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/p' mygo.json | head -1)"
[ -f "$exe" ] && ok "exe: $(du -h "$exe" | cut -f1) ($(file -b "$exe" | cut -c1-70))"
if [ "$dmg" = 1 ]; then
	while IFS= read -r d; do ok "dmg: $d ($(du -h "$d" | cut -f1))"; done < <(find "$out" -maxdepth 2 -name '*.dmg' -newer "$app")
fi

if [ "$on_mac" = 1 ]; then
	step "codesign"
	codesign -dv "$app" 2>&1 | grep -E '^(Identifier|Format|CodeDirectory|Signature|Authority|TeamIdentifier)' || true
	if codesign --verify --deep --strict "$app" 2>/dev/null; then
		ok "signature verifies"
	else
		warn "codesign --verify failed (expected for ad-hoc + Gatekeeper; fine for local runs)"
	fi
	if [ "$open_app" = 1 ]; then
		step "open $app"
		open "$app"
	fi
fi
