# MyClip

[![CI](https://github.com/L92-Labs/myclip/actions/workflows/ci.yml/badge.svg)](https://github.com/L92-Labs/myclip/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A local-only clipboard history for the menu bar on macOS and Linux, with a
Windows preview. Copy something, then find it again in MyClip.
No account, no cloud, no network, no telemetry. Open source (MIT).
<!-- source: docs/RELEASE_NOTES.md:1, docs/RELEASING.md:3, main.go:88, CONTRIBUTING.md:38, LICENSE:1 -->

## Quick start

On an x86-64 or ARM64 Linux desktop with X11 or XWayland, run this from the
root of a fresh clone. You need Bash, Go matching [go.mod](go.mod), Python 3,
GTK 3, OpenGL/EGL, Fontconfig and FreeType.
<!-- source: CONTRIBUTING.md:8, go.mod:3, scripts/build-linux.sh:1, scripts/build-linux.sh:30, main.go:50 -->

```sh
arch=$(go env GOHOSTARCH) &&
  GOARCH="$arch" scripts/build-linux.sh \
    --arch "$arch" &&
  (cd "build/linux-$arch" &&
    ./myclip --toggle)
```
<!-- source: scripts/build-linux.sh:2, scripts/build-linux.sh:24, scripts/build-linux.sh:37, mygo.json:7, .github/workflows/ci.yml:138, main.go:91, main.go:227 -->

For macOS and Linux downloads or Windows source builds, see
[Download and install](#download-and-install).
<!-- source: GitHub release v0.1.0 assets, .github/workflows/ci.yml:120 -->

## Download and install

[Release v0.1.0](https://github.com/L92-Labs/myclip/releases/tag/v0.1.0)
publishes the macOS and Linux packages below with `SHA256SUMS`.
<!-- source: GitHub release v0.1.0 assets -->

| System | Download | Install |
|---|---|---|
| macOS, Apple Silicon or Intel | `myclip-0.1.0-macos-universal.dmg` (or `.zip`) | Open the DMG or extract the ZIP, then drag MyClip to Applications |
| Linux x86-64 | `myclip-0.1.0-linux-amd64.deb` or `.tar.gz` | Debian/Ubuntu: `sudo apt install "./myclip-0.1.0-linux-amd64.deb"` |
| Linux ARM64 | `myclip-0.1.0-linux-arm64.deb` or `.tar.gz` | Debian/Ubuntu: `sudo apt install "./myclip-0.1.0-linux-arm64.deb"` |

On Linux, install the runtime libraries listed under [Build](#build). The
Debian package declares its dependencies, including MyGo's WebKitGTK package
requirement. For the portable archive, unpack it into a separate directory
and run the bundled executable from there. This example starts in the Linux
x86-64 download directory and uses a POSIX shell with `tar`.
<!-- source: v0.1.0 Linux tar.gz contents, scripts/release.py:24, MyGo v0.6.2 cmd/mygo/deb.go:36 -->

```sh
portable_dir=$(mktemp -d) &&
  tar -xzf "myclip-0.1.0-linux-amd64.tar.gz" -C "$portable_dir" &&
  (cd "$portable_dir" && ./myclip --toggle)
```

For a per-user installation, the Linux source build generates `install.sh`
separately from the archive. After building for Linux x86-64 with the
prerequisites under [Build](#build), run the installer from the clone root as
your regular user. The example below builds first, then passes the archive
path explicitly. The installer adds a desktop entry and a `myclip` command
and checks for WebKitGTK.
<!-- source: scripts/build-linux.sh:2, mygo.json:7, MyGo v0.6.2 cmd/mygo/installscript.go:124,166,203,215,249 -->

```sh
scripts/build-linux.sh --arch amd64 &&
  sh "build/linux-amd64/install.sh" "build/linux-amd64/myclip-0.1.0-linux-amd64.tar.gz"
```

Release v0.1.0 has no Windows downloads. On Windows x86-64, install Go
matching [go.mod](go.mod) and build from the clone root in PowerShell.
[CI](https://github.com/L92-Labs/myclip/actions/workflows/ci.yml) also builds a
Windows installer and portable ZIP. Ordinary CI runs retain downloadable
build artifacts for 14 days.
<!-- source: GitHub release v0.1.0 assets, .github/workflows/ci.yml:96,120,147, docs/RELEASING.md:3 -->

```powershell
go build -trimpath -ldflags "-s -w" -o myclip.exe .
```

macOS preview builds are **ad-hoc signed, not Developer ID signed or notarized**.
Gatekeeper may block first launch; use System Settings › Privacy & Security ›
Open Anyway after checking the download. Do not disable Gatekeeper globally.
Pasting needs Accessibility permission; copying and search work without it.
Checksums detect changed downloads. To verify GitHub build provenance,
install and authenticate the GitHub CLI, then run this example from the
download directory.
<!-- source: v0.1.0 .github/workflows/ci.yml:161, gh attestation verify --help -->

```sh
gh attestation verify "myclip-0.1.0-linux-amd64.tar.gz" --repo L92-Labs/myclip
```

Check downloads from the release directory using the checksum tool for your OS.

```sh
case "$(uname -s)" in
  Darwin) shasum -a 256 -c SHA256SUMS ;;
  Linux)  sha256sum -c SHA256SUMS ;;
esac
```
<!-- source: v0.1.0 SHA256SUMS release asset -->

`SHA256SUMS` lists every installer; download all files to check the entire
manifest, or compare the hash of just the installer you downloaded.

![MyClip on Linux: searching a screenshot by recognized text](docs/screenshots/linux-ocr-search.png)

The screenshot uses synthetic clipboard content on an X11 virtual desktop.

## What it does

MyClip follows the design of [Maccy](https://maccy.app) 2.7.1;
[PARITY.md](PARITY.md) lists what matches, the deviations, and MyClip's
additions. It is written in Go with
[MyGo](https://mygo.egoist.dev), using native UI drawn on the GPU. It needs
no cgo and no Xcode project. The universal macOS executable supports
both Apple Silicon and Intel.
<!-- source: scripts/build-mac.sh:6, scripts/release.py:49 -->

> **Status:** early preview. macOS and Linux (X11; Wayland through XWayland)
> and Windows. Linux desktop behavior is tested under Xvfb/Openbox. Real Mac
> permission, paste, OCR and sound checks, plus native Wayland desktop
> checks, still need hands-on validation. Windows CI builds, tests, OCR and
> the installer are exercised automatically, but clipboard, paste and tray
> behavior still need hands-on validation on real hardware.

See the [Awesome MyGo listing](https://awesomemygo.dev/apps/myclip).
<!-- source: README.md:5 at 64d8314 -->

- **History.** Everything you copy (text, rich text, images, files),
  newest first. Pin the clips you reuse; pins get a ⌘-letter of their own.
- **Search.** Exact, fuzzy (`gcm` finds `git commit`), regular expressions,
  or mixed, with the matches highlighted.
- **Paste for you.** Pick a clip and MyClip types ⌘V into the app you were
  in, with or without formatting. The keyboard layout is respected
  (QWERTY, AZERTY, Dvorak).
- **Find screenshots by their words.** Text in copied images is read on
  your computer (Vision on macOS, tesseract on Linux) and searched.
- **Paste stack.** ⇧↑/⇧↓ or ⇧-click several clips, then Enter: they paste
  as one, in order.
- **Hold to cycle (macOS).** Keep ⇧⌘ held after ⇧⌘C, press C again to move
  down, let go to paste or copy.
- **A sound on copy,** if you want one.
- **Preview.** The text (its first 10,000 characters), the image and its
  size, the source app, and how many times and when you copied it.
- **Your rules.** Ignore apps, pasteboard types or text patterns, pause
  capture, or skip just the next copy.
- **Languages.** English, French, Spanish.

## Privacy and security

- **Nothing goes online.** Packaged builds (`scripts/build-mac.sh` and
  `scripts/build-linux.sh`) open
  no network connection: SQLite is bundled in it. Development runs
  (`scripts/dev-mac.sh`, `go test`) may download MyGo's SQLite library once,
  checked against its SHA-256.
- **Password managers are skipped.** Copies marked concealed or transient
  (1Password, KeePassXC and others following the nspasteboard.org convention) are recognised from their types
  *before* any content is read, so the secret never enters memory or the
  database.
- **Files only you can read.** History lives in
  `~/Library/Application Support/com.l92labs.myclip/` (folder 0700, files
  0600). On Linux it lives in
  `${XDG_CONFIG_HOME:-~/.config}/com.l92labs.myclip/`. SQLite overwrites deleted
  content with `secure_delete`; this does not erase filesystem snapshots,
  backups, or copies already held by other applications. History is local
  plaintext, not encrypted; another process running as your user can read it.
- **Size caps.** Clips over 32 MB are skipped and counted. The history keeps
  200 items by default; pins don't count toward that.
- **Fail closed.** If MyClip cannot read its settings, it stops capturing
  rather than fall back to defaults that could miss your ignore rules. The
  menu icon offers Reset Settings.
- **Few permissions.** Copying needs none. Pasting for you needs
  Accessibility, used only to type ⌘V. *Notify on each new copy*, off by
  default, asks for the Notifications permission when you turn it on.

## Keyboard

| In the popup (default keys) | |
|---|---|
| ⇧⌘C | open MyClip from any app (configurable) |
| ↑ ↓ · ⌃N ⌃P · ⌘↑ ⌘↓ | move · first · last |
| ↵ | copy |
| ⌥↵ · ⌥⇧↵ | paste · paste without formatting |
| ⌘1 … ⌘9 | the nth clip; ⌘-letter for a pin |
| ⌥P · ⌥⌫ · ⌃Space | pin · delete · preview (configurable) |
| ⌥⌘⌫ · ⇧⌥⌘⌫ | clear history · clear including pins |
| Esc · ⌘, · ⌘Q | close · settings · quit |

With *Paste automatically* or *Paste without formatting* on, the ↵ ⌥↵ ⌥⇧↵
actions swap as Maccy's do; Settings › General shows the current map.

Click the menu icon with ⌥ to pause or resume capture. While capture is on,
click with ⌥⇧ to ignore only the next copy (while paused, ⌥⇧-click resumes,
as in Maccy).

## Configuration

| Setting | Default | Effect |
|---|---|---|
| Popup shortcut | ⇧⌘C on macOS, ⇧+Super+C on Linux | Opens clipboard history. <!-- source: internal/settings/settings.go:76, scripts/linux-smoke.sh:67 --> |
| Search mode | Exact | Also supports fuzzy, regex and mixed search. <!-- source: internal/settings/settings.go:64, internal/settings/settings.go:174 --> |
| History size | 200 items | Sets how many unpinned clips are kept. <!-- source: internal/settings/settings.go:177, internal/history/history.go:168 --> |
| Notify on each new copy | Off | Shows a notification for each new copy on supported systems. <!-- source: internal/settings/settings.go:118, internal/settings/settings.go:168, main.go:322 --> |
| `MYCLIP_SHOW` | Unset | Set before launch to `popup`, `popup-preview` or `prefs:<tab>` to open that window for screenshots. <!-- source: main.go:221, main.go:233 --> |

## Build

On a Mac with the Xcode Command Line Tools (`xcode-select --install`),
Python 3 and Go matching `go.mod`:

```sh
scripts/build-mac.sh --dmg --open       # test, build the app and DMG, launch
scripts/build-mac.sh --universal --dmg  # arm64 + x86_64
```
<!-- source: scripts/build-mac.sh:2 -->

The first launch of an ad-hoc signed build needs right-click → Open.

For a live-reload development session, run the following in the foreground.
Quit the session with Ctrl+C before running the checks below.
<!-- source: scripts/dev-mac.sh:2,11 -->

```sh
scripts/dev-mac.sh
```

```sh
scripts/test.sh                     # go vet + go test (macOS or Linux)
python3 scripts/licenses.py --check  # verify bundled third-party notices
```
<!-- source: CONTRIBUTING.md:15 -->

On Linux, use Bash, Go matching `go.mod`, Python 3, GTK 3, OpenGL/EGL,
Fontconfig and FreeType. For the full desktop smoke test, install Xvfb,
Openbox, libX11, libXtst (`libXtst.so.6`) and xdotool. Install scrot if
capturing screenshots. See [Contributing](CONTRIBUTING.md#setup).
<!-- source: CONTRIBUTING.md:8, scripts/build-linux.sh:30, scripts/linux-smoke.sh:10,25,35,67, _tools/xkeys.py:3, .github/workflows/ci.yml:132 -->

```sh
scripts/build-linux.sh    # test and build the binary, packages and install.sh
scripts/linux-smoke.sh    # end-to-end on a virtual display
```

| Linux needs | for |
|---|---|
| GTK 3, OpenGL/EGL, Fontconfig/FreeType | native window, rendering and fonts |
| `libayatana-appindicator3` | the menu icon (without it, open MyClip with its shortcut) |
| `xdotool` (X11) or `wtype` / `ydotool` (Wayland) | pasting for you; copying works without |
| `tesseract-ocr` and its language data | local image OCR; text history works without it |
| `xprop` | telling apps apart for Ignore › Applications (by WM_CLASS, typed in Settings); in *ignore all except listed* mode, an app it cannot name is ignored |

Under Wayland, MyClip runs through XWayland, whose clipboard the compositor
keeps in sync, because native Wayland apps cannot read the clipboard while
unfocused. The compositor keeps global keys for itself there, so bind
⇧+Super+C to `myclip --toggle` in your desktop's keyboard settings (running
`myclip` again does the same). The menu icon's menu has Open MyClip too. Password managers following KDE's convention (KeePassXC…) mark
secrets with `x-kde-passwordManagerHint`; MyClip skips them unread. ⌘ in
the key map is Ctrl on Linux, and ⇧⌘C is ⇧+Super+C.

## Layout

| Path | |
|---|---|
| `internal/capture` | the capture rules: ignore chain, change-count sandwich, size cap |
| `internal/history` · `internal/store` | items, merging, sorting · SQLite with the security recipe |
| `internal/engine` | capture → merge → store, copy-back, pins, clearing |
| `internal/native` | the OS: AppKit, CoreGraphics and Carbon on macOS; GTK and X11 on Linux (purego, no cgo) |
| `internal/popup` · `internal/prefs` | the popup and the settings window |
| `internal/search` · `internal/shortcut` · `internal/settings` | search modes · accelerators · preferences |
| `locales/` | EN / FR / ES catalogs (`go tool mygo generate`) |
| `PLAN.md` · `PARITY.md` | milestones and gates · the Maccy feature inventory and deviations |

## Credits

MyClip's design follows [Maccy](https://github.com/p0deje/Maccy) by Alex
Rodionov and contributors (MIT); see [NOTICE.md](NOTICE.md). Icons from
[Lucide](https://lucide.dev) (ISC).

[MIT](LICENSE) © L92 Labs

## Contributing and releases

See [CONTRIBUTING.md](CONTRIBUTING.md) for local setup and checks,
[SECURITY.md](SECURITY.md) for private vulnerability reports, and
[docs/RELEASING.md](docs/RELEASING.md) for the release process.
[THIRD_PARTY_NOTICES.txt](THIRD_PARTY_NOTICES.txt) contains dependency licenses.
Historical design notes in PLAN.md, PARITY.md and LEDGER.md include planned
features; this README describes the current app.
