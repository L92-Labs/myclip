# scripts

| Script | Runs on | Does |
|---|---|---|
| `test.sh` | macOS, Linux | `go vet` + `go test` (cgo off; `-race` turns cgo on) |
| `dev-mac.sh` | macOS | `mygo dev`: live reload, ad-hoc signed |
| `build-mac.sh` | macOS (`--cross` elsewhere) | preflight → test → `mygo build` → codesign report |
| `build-linux.sh` | any (run on Linux) | test → `mygo build linux/<arch>`: binary, `.deb`, `.tar.gz` |
| `linux-smoke.sh` | Linux | end-to-end on Xvfb: capture, password hint skipped, hotkey + popup copy-back, capture-safe |

First time on a Mac: `xcode-select --install`, plus Go matching `go.mod`.

```sh
scripts/build-mac.sh                       # .app for this Mac, no DMG
scripts/build-mac.sh --dmg --open          # .app + .dmg, then launch
scripts/build-mac.sh --universal --dmg     # arm64 + x86_64
scripts/build-mac.sh --sign "MyClip Dev"   # self-signed dev identity
```

Output goes to `build/darwin-<arch>/` (`out` in `mygo.json`). `--cross` builds on Linux, unsigned and with no DMG: it only checks that the code compiles.
