# Finding ledger — plan phase

Per-finding dispositions. Rounds: R1 (v1), R2 (v2), R3 (v3), R4 (v4), R5 (v5). Code-phase ledger entries appended below as PRs proceed. NOTE: R2 GLM `HOTKEY-SCOPE-UNVERIFIED` disposition ("all 4 ⇒ global") is **superseded by R3/R4**: only `popup` is global.

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

## Open user decisions

**Blocking M1 start:**
1. CI runner: fix L92-Labs Actions billing / self-hosted Mac runner / no-CI mode (local test transcripts as PR evidence).
2. D1 automation CLI — default-off; **residual: once enabled, any same-user process can read/clear history (SecCode peer check = tamper tripwire, not confidentiality); destructive verbs keep in-app confirm even with --yes** — approve with that understanding.
3. D5 no `defaults write` CLI (MyClip CLI replaces) — approve.
4. D6 size caps / DB budget — approve or drop.

**Before M3/M5:**
5. D7 poll floor 100ms — approve.
6. D9 Vision OCR skip in v1 — approve or add spike.
7. D2 third language: ES ok?
8. D10 copy sounds via MyGo audio — approve.
9. M5 signing: Developer ID owner? (fallback: ad-hoc signed + hardened runtime)

## Round 3 → v4 (18 findings + 1 BLOCKED)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | n1-sqlite-nettop-vs-unpackaged-download | high | adopted — N1 evidence packaged-only; IsPackaged runtime gate | PLAN N1 |
| Cursor | n1-updater-state-defeats-disable-grep | high | adopted — updater.New pinned, Plugin forbidden, fresh-state test, SetAutomaticChecks only from toggle | PLAN N1 |
| Cursor | capture-pipeline-omits-ignoreEvents | high | adopted — early-exit step 1 | PLAN §3.3, PARITY §1 |
| Cursor | enabled-types-disjoint-missing | high | adopted — chain step 2 | PLAN §3.3 |
| Cursor | own-write-merge-incomplete-vs-History.add | high | adopted — full field matrix + pinned test | PLAN §3.2, PARITY §1 |
| Cursor | compared-types-list-claimed-but-absent | med | adopted — M1 pre-task port exact list | PARITY §1 |
| Cursor | n2-deviation-range-stale | med | adopted — D1–D10 | PLAN N2 |
| Cursor | d1-cli-destructive-confirm-nontty | med | adopted — in-app confirm / gated --yes | PLAN §3.7 |
| Cursor | d1-d5-signoff-orphaned | med | adopted — open decisions updated, D1/D5/D6 block M1 | LEDGER open decisions |
| Cursor | notice-intro-still-Aleksey | (ledger-check) | adopted — name split across lines fixed | NOTICE.md |
| Opus | cli-confused-deputy-tcc | high | adopted — default-off socket, SecCode peer check, select-never-pastes, threat model + rejection test | PLAN §3.7 |
| Opus | in-popup-hotkeys-registered-globally | high | adopted — only popup global; others local; π test | PLAN §3.6, PARITY §4 |
| Opus | m1-gate-depends-on-blocked-actions | high | adopted — runner decision = PRE-M1 blocking | PLAN M1 |
| Opus | sqlite-system-lib-fallback-impossible + download-not-compile-time | med | adopted — fallback dropped, packaged-only evidence, MYGO_SQLITE_LIBRARY dev-only asserted | PLAN N1 |
| Opus | tcc-grant-inheritance-no-hardened-runtime | med | adopted — hardened runtime + library validation from M1 + injection-refusal spike | PLAN §3.6, M5 |
| Opus | d8-fallback-kills-tray-click-popup | med | adopted — coexistence spike + D8 loss recording | PARITY §1, §6 |
| Opus | os-matrix-omits-current-macos | med | adopted — 27 in acceptance matrix; hardware named pre-M1 | PLAN acceptance |
| Opus | sign-off-list-inconsistent | low | adopted — N2 D1–D10; D1/D5/D10 in decisions; D1/D6 blocking | LEDGER |
| Opus | n1-grep-gate-misses-real-knob | low | adopted — SetAutomaticChecks gate + fresh-state test | PLAN N1 |
| Opus | db-file-mode-unspecified | low | adopted — O_EXCL 0600 pre-create + sidecar asserts | PLAN §3.4 |
| GLM | (no verdict emitted — distill SIGILL mid-run) | — | **BLOCKED** — retry round 4 | — |

## Round 4 → v5 (20 findings)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | doc-version-identity-mismatch | high | adopted — headers/revisions bumped (this revision) | PLAN/PARITY/LEDGER headers |
| Cursor | parity-popup-hold-cycle-absent | high | adopted — PARITY §2b row (M2) | PARITY §2b |
| Cursor | parity-§4-header-still-asserts-global-scope | med | adopted — header reworded | PARITY §4 |
| Cursor | cli-socket-fs-boundary-regressed | med | adopted — container 0700 / socket 0600 / disabled = not bound | PLAN §3.7 |
| Cursor | cli-seccode-m1-test-outside-m1-gate | med | adopted — socket+SecCode evidence in M1 bullets | PLAN §3.7, M1 |
| GLM | plan-parity-version-labels-stale | med | adopted (same as headers) | PLAN/PARITY headers |
| GLM | cli-seccode-peer-boundary-overclaim-and-dr-pinning | med | adopted — honest residual (tripwire not confidentiality), spawn-the-CLI threat in text | PLAN §3.7 |
| GLM | r2-f6-socket-mode-disposition-regressed | med | adopted (same as socket fs boundary) | PLAN §3.7 |
| GLM | n1-updater-gate-mygo-api-unverified | med | adopted — M1 API probe + MyClip-owns-state fallback | PLAN N1 |
| GLM | parity-inventory-defaults-keys-blindspot | med | adopted — behavior-diff pre-task | PARITY §2b note |
| GLM | own-write-merge-tests-skip-2-preserve-fields | low | adopted — firstCopiedAt/title asserts | PLAN §3.2 |
| GLM | db-precreate-o-excl-existing-file-unspecified | low | adopted — EEXIST fstat+fchmod or refuse | PLAN §3.4 |
| GLM | parity-s4-header-contradicts-scope-table | low | adopted (same as §4 header) | PARITY §4 |
| GLM | r3-finding-count-mismatch | low | adopted — counts fixed (R3 = 18 findings: Cursor 9 + Opus 9; GLM run incomplete; ledger-check rows are not findings) | LEDGER R3 header |
| Opus | hardened-runtime-and-library-validation-not-produced-by-mygo | high | adopted — post-build re-sign step + sqlite-under-LV spike + Team ID escalation path | PLAN §3.6, M1, PRE-M1 |
| Opus | cli-seccode-peer-check-satisfied-by-launching-the-bundled-cli | high | adopted — honest residual + per-verb opt-ins + in-app confirm even with --yes + D1/UI surfacing | PLAN §3.7, PARITY D1 |
| Opus | artifact-version-labels-stale | low | adopted (headers) + D6 "blocks M1" | headers, PARITY D6 |
| Opus | parity-§4-header-contradicts-scope-rows | low | adopted (same as §4 header) | PARITY §4 |
| Opus | updater-n1-gate-partial | low-med | adopted — SetAutomaticChecks(false) every launch, plugin-window writes allowed, harness named | PLAN N1 |
| Opus | n4-pass-criterion-not-binary | low-med | adopted — numeric gates vs same-host Maccy baseline | PLAN acceptance |

## Round 5 → v6 (17 findings + 1 rejection)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | paste-stack-multiselect-omitted | high | **rejected with evidence** — `AppState.swift:11` `let multiSelectionEnabled = false` @ `eb03eba` ⇒ `startPasteStack` guard unreachable; PasteStack = disabled dead code in 2.7.1 (fetch-verified this session); noted as future feature if upstream enables | PLAN §3.3 |
| Cursor | ignore-only-next-semantics-wrong | high | adopted — exact coupling ported + 3 tests | PLAN §3.3 |
| Cursor | d1-parity-contradicts-plan-destructive-and-residual | high | adopted — D1 rewritten to match PLAN §3.7 | PARITY D1, LEDGER decisions |
| Cursor | n1-updater-force-false-vs-owned-opt-in | med | adopted — SetAutomaticChecks(myclipPref) sequenced apply | PLAN N1 |
| Cursor | cli-socket-m1-tests-without-m1-deliverable | med | adopted — minimal socket+CLI skeleton in M1 | PLAN §3.7, M1 |
| GLM | parity-2b-missing-readme-documented-behaviors | med | adopted — 4 rows added | PARITY §2b |
| GLM | capture-chain-is-remote-clipboard-union-ambiguous | med | adopted — user-setting only | PLAN §3.3 |
| GLM | d1-signoff-summary-missing-residual | med | adopted | PARITY D1, LEDGER #2 |
| GLM | sandbox-posture-and-container-path-undeclared | med | adopted — non-sandboxed declared, real paths, entitlements none | PLAN §3.6–3.7 |
| GLM | ledger-r3-count-fix-introduces-new-mismatch | low | adopted — R3 = 25 (Cursor 9 + Opus 16) | LEDGER R3 |
| GLM | n4-binary-size-gate-feasibility-unvalidated | low | adopted — measure-then-lock | PLAN acceptance |
| Opus | resign-after-build-stales-dmg-notarize-update-archive | high | adopted — no post-sign; signingIdentity in mygo.json; extracted-artifact codesign gate | PLAN §3.6, M1/M5 |
| Opus | updater-optin-reset-every-launch | med | adopted (same as pref-apply) | PLAN N1 |
| Opus | n1-updater-spike-hedge-obsolete | low | adopted — hedge dropped, concrete pre-seed test | PLAN N1 |
| Opus | parity-d1-contradicts-plan-on-yes-and-residual | med | adopted (same as D1 rewrite) | PARITY D1 |
| Opus | own-write-marker-not-excluded-from-reps-and-equality | med | adopted — marker = metadata; 2 tests | PLAN §3.2 |
| Opus | socket-path-undefined-for-nonsandboxed-app | low | adopted — pinned path + fchmod + sun_path check | PLAN §3.7 |
