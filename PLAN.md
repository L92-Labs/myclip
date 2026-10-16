# MyClip — Implementation Plan (v2)

Lightweight, open-source, 100%-local clipboard manager for macOS. **Parity contract: Maccy 2.7.1 — see `PARITY.md` (authoritative inventory; the matrix in v1 was superseded).** Stack: Go + [MyGo](https://mygo.egoist.dev) v0.5.5 (pure Go, no cgo, native GPU UI, tray, global shortcuts, clipboard APIs, SQLite plugin, updater).

Repo: `L92-Labs/myclip` (private now; public later). License: MIT + `NOTICE.md` crediting Maccy. Bundle ID: `com.l92labs.myclip` (locked; never change after first TCC grant). Go module: `github.com/L92-Labs/myclip`.

Revision history: v1 → v2 after adversarial plan review (Cursor Auto, Claude Opus 5.5, OpenCode GLM-5.3-Flash; all FINDINGS, ledger consolidated below).

## Non-negotiable constraints

- N1: No cloud, no telemetry, no unsolicited network. Updater compiled with `DisableAutomaticChecks: true` + explicit opt-in "Check for Updates" (grep gate in CI: `DisableAutomaticChecks` present, default-on absent).
- N2: Behavioral parity with Maccy 2.7.1 per `PARITY.md` — every row maps to a milestone or an approved deviation (D1–D5). Gaps found later = findings; never silent drops.
- N3: Open source (MIT); Maccy MIT notice included (we consult Maccy source, so no "clean-room" claim); THIRD_PARTY_NOTICES generated at M5.
- N4: Lightweight, measurable per milestone evidence gate: binary size, idle RSS, idle CPU% over ≥5min — "no busy-wait; resource envelope comparable to Maccy side-by-side in Activity Monitor" (literal "0% CPU" dropped as unmeasurable for a poller).
- N5: Review gates per `feature-review-loop` skill. Every gate reviewer gets a full fresh-defect pass on the final revision — gate minimums are not reduced by model routing. Routing preference (user): Cursor Agent Auto does the deep/heavy passes; Claude/ChatGPT used at skill-required minimum. cgo is a hard non-goal; if a capability is impossible via purego, the fallback decision is explicit in the plan, never silent.
- N6: Security: raw-UTI ignore filter runs **before** any pasteboard data is read (concealed/transient types may be renamed/dropped by MIME-mapped APIs — read `NSPasteboard.types` raw via purego). Secrets never enter the DB (test: concealed item never read into memory).

## Architecture (review-driven corrections)

1. **Purego darwin shim layer** (`internal/darwin`): `changeCount` (no public MyGo API; internal-only today), raw pasteboard UTI enumeration, `frontmostApplication` bundle ID, keycode-from-layout (TIS/UCKeyTranslate; QWERTY force when layout name ends `⌘`), CGEvent posting, tray `NSEvent.modifierFlags` probe inside `OnClick` (MyGo exposes only Clicked/RightClicked — spike; if impossible: menu-item fallback for pause/ignore-next and PARITY row demoted with finding tag).
2. **Item model**: one history item, N representations (raw UTI + bytes), source bundle ID, firstCopiedAt/lastCopiedAt, numberOfCopies, pin, title override. Dedup/equality over representation set; copy-back = write full representation set; paste-plain = write text representation only.
3. **Poller**: `changeCount` watermark; own writes carry marker UTI `com.l92labs.myclip.own-write` and are skipped (prevents recapture churn); ignore chain order: raw UTIs (enabled types → ignored types → always-ignore trio) → source app (ignoredApps/allowlist) → regex on text → capture.
4. **Storage**: SQLite (MyGo plugin). 0600 file in 0700 app-support dir; secure_delete on delete/clear/clearOnQuit; per-item cap + total budget; oversized items skipped gracefully. Dev/CI: pin SQLite native lib (`MYGO_SQLITE_LIBRARY` pre-fetch; no runtime downloads → keeps N1).
5. **App shell**: accessory app (no Dock icon / `ActivationPolicyAccessory` per MyGo menubar examples; non-activating popup per M2 spike gate). Focus model: capture frontmost app before popup shows; paste targets captured app; do not become "previous app" ourselves.
6. **Paste-back**: CGEvent-only (no osascript fallback — it adds an Automation TCC prompt for no added capability). Sequence: hide popup → bounded activate/wait on captured app → resolve keycode for current layout → post ⌘V. Both failure paths (no Accessibility grant; activation timeout) have explicit UX.

## Milestones (hard gates between them)

- **M1 engine + spikes (macOS-gated)**: shim layer items 1,2,3 + store + capture round-trip spike (RTF/TIFF/private UTI/multi-file; concealed-type-never-read test) + **paste-back spike** (CGEvent + TCC + focus + 3 layouts). Gate: macOS CI job (build + darwin-tagged tests) green AND both spikes demonstrated (Agentdrop evidence) **before any M2 merge**.
- **M2 UI**: accessory shell, tray (+modifier spike), popup (position/screen), list (all §2 toggles), search modes + highlight, keyboard map, pins, clear/alert, footer, preview slideout, i18n keys wired from first string (no retrofit).
- **M3 paste integration**: pasteByDefault, removeFormattingByDefault, plain-paste, focus edge cases (inactive/capture-failed), TCC grant/refresh UX.
- **M4 settings + i18n**: all six panes per PARITY §5, hotkey recorder, launch-at-login, 3 languages.
- **M5 ship**: packaging (DMG), signing + notarization (macOS runner), opt-in updater, THIRD_PARTY_NOTICES, README, release CI. Public flip = user decision, after parity audit.

## Acceptance criteria (global)

- Every PARITY.md row demonstrably matches Maccy 2.7.1 side-by-side on macOS 14/15 (evidence matrix: OS × behavior × screenshot; per milestone).
- CI: Linux `go vet`/`go test` (logic) + macOS job (darwin build + darwin tests) from M1 onward.
- N1 proof: zero outbound connections during scripted idle + normal session (nettop capture) in M5 evidence.
- N4 proof per milestone: binary size + RSS + idle CPU% recorded in PR evidence.

## Consolidated review ledger (plan phase, v1→v2)

All 44 findings (Cursor 13, Opus 17, GLM 14) dispositioned: adopted into v2 above (matrix rebuilt from Defaults.Keys @ 2.7.1; own-write marker; changeCount/raw-UTI/frontmost/keycode shims; bundle ID + module path locked; updater opt-in; N1/N2/N4/N5/N6 reworded; macOS CI from M1; paste-back CGEvent-only + layout spec; item model; storage security; i18n from M2; NOTICE corrected). Rejected with evidence: GLM "default hotkey ⌥⌘C" — Maccy 2.7.1 `KeyboardShortcuts.Name+Shortcuts.swift` sets popup default ⇧⌘C.

## Review & routing (per N5)

- Plan gate round 2: this v2 + PARITY.md → Cursor Auto (deep), Claude Opus 5.5, OpenCode GLM-5.3-Flash. Repeat until all CLEAR on same version.
- Code gate per PR: same set; full fresh pass each on final revision; Cursor Auto additionally does interim deep passes. ChatGPT/Codex enters via PR-feedback gate per skill.
- Evidence: Agentdrop @ drop.l92-labs.com, external links only in PR descriptions; review artifacts in /tmp, never in-repo.
