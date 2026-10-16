# Finding ledger — plan phase

Per-finding dispositions. Rounds: R1 (plan v1, 2026-10-16), R2 (plan v2). Code-phase ledger entries appended below as PRs proceed.

## Round 1 → v2 (44 findings)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | parity-matrix-incomplete-vs-defaults | high | adopted — PARITY.md rebuilt from Defaults.Keys @2.7.1 | PARITY §1–5 |
| Cursor | own-write-recapture-marker-missing | high | adopted; **corrected in R2**: merge semantics, not skip | PLAN §3.2, PARITY §1 |
| Cursor | tray-modifier-clicks-unsupported-by-mygo | high | adopted — M1 spike + D8 fallback | PLAN §3.1, PARITY §6 |
| Cursor | updater-default-violates-local-only | high | adopted — DisableAutomaticChecks + CI gate | PLAN N1 |
| Cursor | placeholder-bundle-id-breaks-tcc-login-updates | high | adopted — com.l92labs.myclip locked; signing identity added R2 | PLAN header, §3.1 |
| Cursor | ignore-apps-and-regex-omitted | high | adopted | PARITY §1 |
| Cursor | osascript-fallback-extra-tcc | med | adopted — CGEvent-only, fallback removed | PLAN §3.6 |
| Cursor | paste-keycode-layout-switching | med | adopted — TIS/QWERTY-⌘ rule | PLAN §3.6, PARITY §3 |
| Cursor | rich-type-fidelity-via-mygo-transfer | med | adopted — round-trip spike + setData fallback | PLAN §3.2, PARITY §1 |
| Cursor | menubar-activation-and-panel-unspecified | med | adopted — accessory + non-activating popup spike moved to M1 (R2) | PLAN §3.5, M1 |
| Cursor | preview-and-app-intents-absent | med | adopted — preview rows; intents = D1 | PARITY §2, §6 |
| Cursor | paste-back-ungated-and-linux-ci-blind | med | adopted — M1 hard gate; CI split R2 (TCC not in CI) | PLAN M1 |
| Cursor | cgo-escape-vs-pure-go-crosscompile | med | adopted — cgo hard non-goal, decisions explicit | PLAN N5 |
| Cursor | idle-cpu-acceptance-unachievable | low | adopted — N4 reworded measurable | PLAN N4 |
| Opus | F1-parity-ref-unpinned | high | adopted — pinned 2.7.1 @ eb03eba | PARITY header |
| Opus | F2-matrix-missing-prefs | high | adopted — Defaults.Keys inventory | PARITY §1–5 |
| Opus | F3-source-app-not-captured | high | adopted — bundle ID column + frontmost shim | PLAN §3.1–3.2 |
| Opus | F4-no-changecount-poll-cost | high | adopted — purego changeCount, read-on-change | PLAN §3.1 |
| Opus | F5-ignore-before-read | high | adopted — sharpened R2: any-intersection on full raw set | PLAN N6, §3.3 |
| Opus | F6-own-write-recapture | med | adopted; corrected R2 (merge) | PLAN §3.2 |
| Opus | F7-writeback-fidelity | med | adopted | PARITY §1 |
| Opus | F8-osascript-fallback-tcc | med | adopted — removed | PLAN §3.6 |
| Opus | F9-keyboard-layout-keycode | med | adopted | PARITY §3 |
| Opus | F10-bundle-id-tcc-stability | med | adopted + signing identity M1 (R2) | PLAN M1 |
| Opus | F11-updater-default-on | med | adopted | PLAN N1 |
| Opus | F12-data-at-rest | med | adopted + journal spec R2; caps → D6 | PLAN §3.4, PARITY §6 |
| Opus | F13-darwin-untested-until-M5 | med | adopted — macOS CI M1; tiered R2 (TCC not in CI) | PLAN M1 |
| Opus | F14-notice-cleanroom-contradiction | med | adopted — NOTICE verbatim notice, no clean-room claim | NOTICE.md |
| Opus | F15-sqlite-runtime-download | low | adopted — prod embed too (R2) + M1 nettop | PLAN N1, §3.4 |
| Opus | F16-i18n-retrofit-order | low | adopted — keys from first string | PLAN M2 |
| Opus | F17-n4-unmeasured | low | adopted — per-milestone measurement | PLAN N4 |
| GLM | paste-back-spec | high | adopted — bounded activate/wait + TCC enumeration | PLAN §3.6 |
| GLM | item-representation-model | high | adopted | PLAN §3.2 |
| GLM | default-hotkey ⌥⌘C | med | **rejected** — evidence: KeyboardShortcuts.Name+Shortcuts.swift popup default ⇧⌘C (verified independently by Cursor R1 + Opus R2) | PARITY §4 |
| GLM | n1-vs-autoupdate | med | adopted | PLAN N1, PARITY D3 |
| GLM | self-echo-suppression | med | adopted; corrected R2 (merge) | PLAN §3.2 |
| GLM | bundle-identifier-placeholder | med | adopted | mygo.json |
| GLM | no-lsuielement-plan | med | adopted | PLAN §3.5 |
| GLM | n4-not-gated | med | adopted | PLAN N4 |
| GLM | macos-runtime-tests | med | adopted | PLAN M1 |
| GLM | localization-parity-scope | low | adopted — D2 named languages | PARITY §6 |
| GLM | security-claim-undefined | low | adopted — defined in storage spec | PLAN §3.4 |
| GLM | mygo-capability-unverified | low | adopted — M1/M2 spike checklist | PLAN M1–M2 |
| GLM | module-path | low | adopted | go.mod |
| GLM | review-routing-independence | low | adopted — N5 gate minimums | PLAN N5 |

## Round 2 → v3 (34 findings)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | own-write-skip-vs-maccy-dedup | high | adopted — merge semantics v3 | PLAN §3.2 |
| Cursor | ignore-filter-order-security | high | adopted — any-intersection explicit | PLAN N6, §3.3 |
| Cursor | modifier-action-matrix-underspecified | high | adopted — PARITY §2b + M3 tests | PARITY §2b |
| Cursor | image-ocr-title-missing | med | adopted — D9 pending user | PARITY §1, §6 |
| Cursor | searchVisibility-enum-wrong | med | adopted — always/duringSearch + bool | PARITY §2 |
| Cursor | d5-kills-ignoreevents-cli | med | adopted — CLI verbs | PARITY D1/D5 |
| Cursor | capture-pipeline-edge-cases-absent | med | adopted — rules + tests | PLAN §3.3 |
| Cursor | copy-notifications-absent | med | adopted — rows + D10 | PARITY §1, §6 |
| Cursor | d1-get-intent-not-equivalent | med | adopted — CLI `get` returns content | PARITY D1 |
| GLM | UNAPPROVED-DEVIATION-STORAGE-CAPS | med | adopted — D6 | PARITY §6 |
| GLM | SECURE-DELETE-INCOMPLETE | med | adopted — journal spec + forensics test | PLAN §3.4 |
| GLM | M1-GATE-CI-TESTABILITY-TCC | med | adopted — CI/spike tier split | PLAN M1 |
| GLM | NOTICE-COPYRIGHT-DRIFT | low | **rejected** — 2.7.1 LICENSE fetched: "Copyright (c) 2025 Alex Rodionov" (curl raw @2.7.1, this session); prose name aligned to "Alex" | NOTICE.md |
| GLM | HOTKEY-SCOPE-UNVERIFIED | low | adopted — scope column (all 4 = KeyboardShortcuts.Name ⇒ global) | PARITY §4 |
| GLM | PARITY-MISSING-LAUNCHATLOGIN-ROW | low | adopted | PARITY §2 |
| GLM | ENABLED-TYPES-DEFAULT-UNVERIFIED | low | adopted — M1 pre-task assert | PARITY §1 |
| GLM | N1-GREP-GATE-UNIMPLEMENTABLE | low | adopted — concrete grep + nettop | PLAN N1 |
| GLM | OWN-WRITE-COPYBACK-METADATA | low | adopted (same as merge correction) | PLAN §3.2 |
| Opus | F1 pasteboard-privacy-os-gap | high | adopted — macOS 26/27 matrix + M1 privacy spike | PLAN §3.6, M1 |
| Opus | F2 tcc-grant-vs-code-identity | high | adopted — stable signing identity M1 + tccutil doc | PLAN M1 |
| Opus | F3 m1-gate-unreachable-in-ci | high | adopted — tiered gates + runner decision item | PLAN M1, M5 |
| Opus | F4 sqlite-native-lib-n1-prod | high | adopted — prod embed + M1 nettop | PLAN N1, §3.4 |
| Opus | F5 secure-delete-incomplete | med | adopted | PLAN §3.4 |
| Opus | F6 url-scheme-trust-boundary | med | adopted — 0600 unix socket + peer-UID CLI | PLAN §3.7 |
| Opus | F7 own-write-skip-breaks-recency | med | adopted | PLAN §3.2 |
| Opus | F8 dedup-exact-repr-set | med | adopted — supersede rule ported | PLAN §3.2 |
| Opus | F9 undeclared-deviations | med | adopted — D6–D10 | PARITY §6 |
| Opus | F10 popup-spike-misordered | med | adopted — spike in M1 | PLAN §3.5, M1 |
| Opus | F11 keyboard-map-incomplete | med | adopted — PARITY §2b | PARITY §2b |
| Opus | F12 settings-rows-missing | low | adopted — row-per-control rule | PARITY §5 |
| Opus | F13 parity-inventory-unverifiable | med | adopted — SHA pinned; Source anchors by section | PARITY header |
| Opus | F14 D2-ambiguous | low | adopted — ES named, changeable | PARITY D2 |
| Opus | F15 n6-test-seam | low | adopted — fake + named pasteboard | PLAN §3.1 |
| Opus | F16 signing-dependency | low | adopted — M5 decision item | PLAN M5 |

## Open user decisions (blocking nothing before M1; needed before M3/M5)

1. D6 size caps / DB budget — approve or drop.
2. D7 poll floor 100ms — approve.
3. D9 Vision OCR skip in v1 — approve or add spike.
4. D2 third language: ES ok?
5. M5 signing: Developer ID owner? Actions billing or local release builds?
