# Contributing

MyClip is an MIT-licensed, local-only clipboard manager. Contributions stay
under the project's MIT license; keep third-party attribution intact.

## Setup

Install the Go version in go.mod and Python 3. macOS builds need Xcode Command
Line Tools (`xcode-select --install`). Linux native tests need GTK 3,
OpenGL/EGL, Fontconfig and FreeType. Desktop smoke tests additionally need
Xvfb, Openbox, libXtst, xdotool and scrot (for screenshots).

```sh
git clone https://github.com/L92-Labs/myclip.git
cd myclip
scripts/test.sh -count=1
go test -race -count=1 ./...
python3 scripts/licenses.py --check
scripts/build-linux.sh  # Linux; use scripts/build-mac.sh on macOS
```

After a dependency update, run `python3 scripts/licenses.py` and commit the
updated THIRD_PARTY_NOTICES.txt. Do not hand-edit generated locale accessors;
edit locales/*.json, then run `go tool mygo generate`.

## Pull requests

Describe the observable problem, resulting behavior, and validation. Add a
regression test for a bug where practical. CI runs on Linux and macOS and
builds all supported release targets. Linux ARM64 is cross-built; the CI
test runners do not exercise an ARM64 Linux desktop. macOS unit tests do not
validate Accessibility permissions, Vision OCR, audio devices or paste focus.

Use synthetic clipboard data in tests and screenshots. Never attach real
clipboard databases, credentials, signing certificates or personal history.
Report vulnerabilities privately through [SECURITY.md](SECURITY.md).

The app must keep clipboard content on the device. Changes introducing
networking, telemetry, cloud sync or automatic downloads require an explicit
product decision and corresponding privacy documentation.

Current release instructions live in [docs/RELEASING.md](docs/RELEASING.md).
PLAN.md, PARITY.md and LEDGER.md are historical design records, not a checklist
of already implemented features.
