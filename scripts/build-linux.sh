#!/usr/bin/env bash
# Build MyClip for Linux: vet + tests -> mygo build (binary, .deb, .tar.gz
# with install.sh). Runs on any OS (no cgo); run the app on Linux.
#
#   scripts/build-linux.sh               this machine's arch (amd64 or arm64)
#   scripts/build-linux.sh --arch arm64
#   scripts/build-linux.sh --skip-tests
#
# Runtime: GTK 3; libayatana-appindicator3 for the menu icon; xdotool (X11)
# or wtype / ydotool (Wayland) to paste; xprop to tell apps apart.

source "$(dirname "$0")/lib.sh"

arch="" skip_tests=0
while [ "$#" -gt 0 ]; do
	case "$1" in
	--arch) arch="${2:?--arch needs amd64|arm64}"; shift ;;
	--skip-tests) skip_tests=1 ;;
	-h|--help) sed -n '2,10p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
	*) die "unknown flag $1 (see --help)" ;;
	esac
	shift
done
if [ -z "$arch" ]; then
	case "$(uname -m)" in
	aarch64|arm64) arch=arm64 ;;
	*) arch=amd64 ;;
	esac
fi
need go "install Go from https://go.dev/dl"

[ "$skip_tests" = 1 ] || "$ROOT/scripts/test.sh"

out="$(out_dir)"
step "go tool mygo build -platform linux/$arch"
go tool mygo build -platform "linux/$arch"
step "artifacts"
find "$out" -maxdepth 2 \( -name 'myclip' -o -name '*.deb' -o -name '*.tar.gz' \) -newer go.mod -path "*linux*" -print |
	while IFS= read -r f; do ok "$f ($(du -h "$f" | cut -f1))"; done
