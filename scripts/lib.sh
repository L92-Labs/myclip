# Shared helpers for scripts/*.sh. Source it; do not run it.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [ -t 1 ]; then
	_b=$'\033[1m' _g=$'\033[32m' _r=$'\033[31m' _y=$'\033[33m' _x=$'\033[0m'
else
	_b='' _g='' _r='' _y='' _x=''
fi

step() { printf '%s==>%s %s\n' "$_b" "$_x" "$*"; }
ok()   { printf '%s ok%s %s\n' "$_g" "$_x" "$*"; }
warn() { printf '%swarn%s %s\n' "$_y" "$_x" "$*" >&2; }
die()  { printf '%sfail%s %s\n' "$_r" "$_x" "$*" >&2; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "missing '$1'${2:+ — $2}"; }

# Output dir from mygo.json "out" (default dist, as mygo build does).
out_dir() {
	local out
	out="$(sed -nE 's/^[[:space:]]*"out"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/p' mygo.json | head -1)"
	printf '%s\n' "${out:-dist}"
}
