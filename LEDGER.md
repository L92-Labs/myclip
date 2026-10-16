# Finding ledger — plan phase

Per-finding dispositions. Rounds: R1 (v1), R2 (v2), R3 (v3), R4 (v4), R5 (v5), R6 (v7), R7 (v8), R8 (v9), R9 (v10), R10 (v11), R11 (v12).

## Round 1 → v2 (45 findings: Cursor 14 + Opus 17 + GLM 14)

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
2. D1 automation CLI — default-off; **residual: once enabled, any same-user process can read history (and DB files) and trigger clear/delete requests — requests never complete without in-app confirmation, even with --yes (SecCode peer check = tamper tripwire, not confidentiality)** — approve with that understanding.
3. D5 no `defaults write` CLI (MyClip CLI replaces) — approve.
4. D6 size caps / DB budget — approve or drop.

**Before M3/M5:**
5. D7 poll floor 100ms — approve.
6. D9 Vision OCR skip in v1 — approve or add spike.
7. D2 third language: ES ok?
8. D10 copy sounds via MyGo audio — approve.
9. M5 signing: Developer ID owner? (fallback: self-signed identity via the production signing path — emits hardened runtime; plain ad-hoc only as recorded deviation)

## Round 3 → v4 (19 findings: Cursor 9 + Opus 10; + GLM run incomplete/BLOCKED)

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
| GLM | r3-finding-count-mismatch | low | adopted — **superseded by R5 recount**: R3 = 25 (Cursor 9 + Opus 16); GLM run incomplete; ledger-check rows are not findings | LEDGER R3 header |
| Opus | hardened-runtime-and-library-validation-not-produced-by-mygo | high | adopted — post-build re-sign step + sqlite-under-LV spike + Team ID escalation path | PLAN §3.6, M1, PRE-M1 |
| Opus | cli-seccode-peer-check-satisfied-by-launching-the-bundled-cli | high | adopted — honest residual + per-verb opt-ins + in-app confirm even with --yes + D1/UI surfacing | PLAN §3.7, PARITY D1 |
| Opus | artifact-version-labels-stale | low | adopted (headers) + D6 "blocks M1" | headers, PARITY D6 |
| Opus | parity-§4-header-contradicts-scope-rows | low | adopted (same as §4 header) | PARITY §4 |
| Opus | updater-n1-gate-partial | low-med | adopted — SetAutomaticChecks(false) every launch, plugin-window writes allowed, harness named | PLAN N1 |
| Opus | n4-pass-criterion-not-binary | low-med | adopted — numeric gates vs same-host Maccy baseline | PLAN acceptance |

## Round 5 → v6 (16 findings + 1 evidence-rejection)

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

## Round 6 → v7 (11 findings + GLM run incomplete ×2)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | pathuserdata-vs-pinned-bundle-dir | high | adopted — SetPath(PathUserData, pinned dir) before Run + no-default-dir test | PLAN §3.4 |
| Cursor | m1-hard-gate-omits-cli-socket | med | adopted — socket/CLI bullets in M1 gate checklist | PLAN §3.7, M1 |
| Cursor | m1-update-archive-codesign-gate-infeasible | med | adopted — DMG-only at M1; update-archive check at M5 | PLAN §3.6 |
| Cursor | d1-residual-overclaims-silent-clear | med | adopted — residual reworded (request ≠ silent clear) | PLAN §3.7, PARITY D1 |
| Opus | updater-n1-gate-vacuous-without-updates-config | high | adopted — updates config + Enabled()==true precondition, else explicit N/A→M5 | PLAN §3.6 |
| Opus | mygo-userdata-dir-diverges-from-pinned-data-dir | med | adopted (same as SetPath) | PLAN §3.4 |
| Opus | hold-to-cycle-unimplementable-on-mygo-globalshortcut | med | adopted — M1 Carbon/flagsChanged spike + D11 escalation | PLAN §3.6, PARITY §4 |
| Opus | ignore-next-coupling-as-written-is-a-noop | med | adopted — verbatim source quote pre-task + CLI verb semantics + unpaused test | PLAN §3.3, PARITY §1 |
| Opus | copy-notification-content-escapes-forensics-and-clear | med | adopted — ClearNotifications on clear; usernoted documented as system-level parity | PARITY §1, PLAN acceptance |
| Opus | cli-skeleton-missing-from-m1-gate-list | low | adopted (same as M1 bullets) | PLAN M1 |
| Opus | seccode-peer-identity-by-pid-race | low | adopted — LOCAL_PEERTOKEN audit token + exec-after-connect test | PLAN §3.7 |
| GLM | (no verdict — run incomplete ×2, infra SIGILL) | — | **BLOCKED** — retry R7 | — |

Ledger-count audit (settled R10): counting rule = ledger-check/BLOCKED rows are not findings. R1 = 45 (14/17/14); R2 = 34 (9/16/9); R3 = 19 (Cursor 9 + Opus 10; earlier claims of 18 then 25 were both wrong); R4 = 20 (5/9/6); R5 = 16 + 1 evidence-rejection; R6 = 11 (4/7); R7 = 18 (6/7/5); R8 = 18 labeled (Cursor 8 + Opus 6 + GLM 4; verdict counts 8/6/6 incl. 2 GLM meta rows; itemized below); R9 = 8 themes (14 labeled = Cursor 4 + Opus 6 + GLM 4); R10 = 6 themes (17 labeled = Cursor 5 + Opus 7 + GLM 5). Historical footnote claiming "16 Opus rows in R3 / 55 total" was wrong — superseded by this recount.

## Round 7 → v8 (18 findings: Cursor 6 + Opus 7 + GLM 5)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | m1-update-archive-codesign-self-contradiction | high | adopted — single stance: updates config at M1 ⇒ archive exists ⇒ archive check at M1 | PLAN §3.6 |
| Cursor | setpath-disposition-incomplete-vs-default-userdata | high | adopted — before-Run ordering + no-default-dir test (M1 gate list) | PLAN §3.4/§3.6 |
| Cursor | d1-residual-still-overclaims-clear | med | adopted — "read + trigger requests; never completes without in-app confirm" | PLAN §3.7, PARITY D1, LEDGER #2 |
| Cursor | d11-missing-from-parity-§6-and-n2 | med | adopted — D11 row (contingent) + N2 D1–D11 | PARITY §6, PLAN N2 |
| Cursor | hold-cycle-carbon-vs-globalshortcut-exclusive-path | med | adopted — single exclusive owner, signature check, grep gate | PLAN §3.6 |
| Cursor | cli-destructive-confirm-vs-shortcuts-automation | med | adopted — confirm UX contract (activation + modal, decline default) | PARITY D1 |
| Opus | updater-pref-source-contradiction | med | adopted — pref only in MyClip store; differentiated tests | PLAN N1 |
| Opus | sqlite-env-override-honored-in-packaged-app | med | adopted — Unsetenv on IsPackaged + resolved-path test | PLAN N1 |
| Opus | n6-type-set-first-item-only | med-high | adopted — union of item types + 2-item test | PLAN §3.3 |
| Opus | m1-update-archive-gate-contradiction | med | adopted (same as Cursor #1) | PLAN §3.6 |
| Opus | m5-adhoc-hardened-runtime-fallback-infeasible | med | adopted — self-signed via production path as fallback | PLAN M5 |
| Opus | hold-to-cycle-spike-criterion-wrong-event | med | adopted — modifier-release criterion, not key-up | PLAN §3.6 |
| Opus | carbon-handler-collision-with-mygo-globalshortcut | med | adopted (same as exclusive owner) | PLAN §3.6 |
| GLM | plan-§3.6-update-archive-gate-self-contradiction | med | adopted (same stance fix) | PLAN §3.6 |
| GLM | r6-setpath-disposition-unimplemented | med | adopted — ordering + test + anchors fixed | PLAN §3.4/§3.6 |
| GLM | r6-d1-residual-reword-not-applied | med | adopted — reword applied in all three spots | PLAN/PARITY/LEDGER |
| GLM | ledger-r3-count-false-audit | med | adopted — awk recount; headers/footnote/PLAN corrected (R1=45, R3=20) | LEDGER headers |
| GLM | m1-gate-list-incomplete-propagation | low-med | adopted — canonical M1 gate checklist enumerates everything | PLAN M1 |

## Round 8 → v9 (18 labeled findings: Cursor 8 + Opus 6 + GLM 4 + 2 meta rows; verdicts reported 8/6/6)

Root cause recorded: several R7 disposition edits silently failed to match in the edit script while the ledger claimed "adopted" — reviewers caught every one. v9 uses ASSERTED, sentence-bounded replacements (script exits non-zero if any target is absent; no greedy line-eating regexes). All 20 adopted: update-archive single stance (M1 archive check; M5 only in infeasible branch); SetPath before-Run ordering + fresh-HOME test; D1 residual reword applied (PLAN + LEDGER); N6 union + named `concealed-on-second-item` test; canonical 15-item M1 gate checklist materialized in Milestones; MyGo `GlobalShortcut` banned entirely (Signature-blind handler); N1 nettop split criterion; LEDGER decision #9 aligned; PARITY hold-cycle criterion aligned. (D11/N2, pref-source, Unsetenv, M5 fallback, confirm-UX were already correctly present in v8.)

## Round 8 → v9 — itemized (18 labeled findings: Cursor 8 + Opus 6 + GLM 4 + 2 meta rows; verdicts reported 8/6/6)

| Reviewer | Label | Sev | Disposition | Anchor |
|---|---|---|---|---|
| Cursor | update-archive-gate-still-self-contradictory | high | adopted — single stance; M5 sentence removed | PLAN §3.6 |
| Cursor | setpath-tilde-literal-skips-mkdir | high | adopted — UserHomeDir/Join + MkdirAll/fchmod; gate #11 absolute+exists asserts | PLAN §3.6, M1 #11 |
| Cursor | d1-residual-overclaim-persists-in-plan | med | adopted — PLAN reworded to PARITY wording | PLAN §3.7 |
| Cursor | n6-multi-item-ignore-test-not-specified | med-high | adopted — union + `concealed-on-second-item` | PLAN N6/§3.3 |
| Cursor | canonical-m1-gate-checklist-still-incomplete | med | adopted — 15-item checklist materialized | PLAN M1 |
| Cursor | ledger-open-m5-signing-fallback-stale | med | adopted — decision #9 = self-signed production path | LEDGER #9 |
| Cursor | parity-hold-cycle-spike-criterion-drifts-from-plan | med | adopted — PARITY criterion = PLAN | PARITY §4 |
| Cursor | globalshortcut-any-register-starves-exclusive-carbon-path | med | adopted — total GlobalShortcut ban (R9 refined to literal grep) | PLAN §3.7 |
| Opus | F1 setpath-tilde-literal | high | adopted (same as Cursor#2) | PLAN §3.6 |
| Opus | F2 n6-union-not-propagated-to-shim | med | adopted — §3.1 union wording (completed R10) | PLAN §3.1 |
| Opus | F3 updater-constructor-wrong-frontend-and-grep-hole | med | adopted — native.New + extended grep (completed R10) | PLAN N1 |
| Opus | F4 ledger-r8-not-itemized | med | adopted — this table; R6 row superseded note; count rule settled (completed R10) | LEDGER |
| Opus | F5 hold-cycle-text-spliced-and-misanchored | low-med | adopted — block rewritten clean (completed R10) | PLAN §3.7 |
| Opus | F6 doc-version-identity-regressed | low | adopted — PARITY stamp bumped (completed R10) | PARITY header |
| GLM | parity-version-stamp-v7 | high | adopted — stamp (v9) | PARITY header |
| GLM | ledger-r8-unitemized | high | adopted — this table + index extended | LEDGER |
| GLM | m1-crossref-claim-unimplemented | low | adopted — checklist = gate SSOT; crossref claim dropped (completed R10) | PLAN M1 |
| GLM | ledger-rounds-index-stale | low | adopted — index through R8 | LEDGER header |
| GLM | (meta row — not a finding: GLM R8 emitted 4 labeled findings, 2 overlapping themes above) | ledger-check | — | — |
| (meta row — not a finding) | R8 verdict counts: Cursor 8 + Opus 6 + GLM 6 reported; itemized GLM labeled findings = 4 (F1–F4) + 2 meta ⇒ **R8 = 18 labeled findings** | ledger-check | PLAN/headers use 18 | LEDGER |

## Round 9 → v10 (8 themes from 14 labeled findings: Cursor 4 + Opus 6 + GLM 4)

All adopted: SetPath `~` expansion + MkdirAll ordering (Cursor/Opus); §3.6 stale nettop-10s hedge removed (Cursor); PARITY stamp (all); R8 ledger itemized + index + count rule settled (Cursor/Opus/GLM); N6 union propagated to §3.1 (Opus); updater `native.New` constructor + extended grep (Opus); §3.7 hold-cycle splice rewritten + literal grep (Opus); M1 SSOT wording (GLM).

## Round 9 → v10 — itemized (15 labeled findings = 9 themes: Cursor 5 + Opus 6 + GLM 4)

| Reviewer | Label | Sev | Theme | Disposition | Anchor |
|---|---|---|---|---|---|
| Cursor | setpath-tilde-literal-skips-mkdir | high | T1 SetPath recipe | adopted — UserHomeDir/Join + MkdirAll/fchmod + gate #11 asserts | PLAN §3.6 |
| Opus | F1 setpath-tilde-literal | high | T1 | adopted (dup) | PLAN §3.6 |
| Cursor | n1-nettop-split-vs-§3.6-vacuous-10s | med | T3 stale hedge | adopted — hedge removed | PLAN §3.6 |
| Cursor | parity-version-identity-still-v7 | med | T4 stamps | adopted — PARITY (v9) | PARITY header |
| GLM | parity-version-stamp-v7 | high | T4 | adopted (dup) | PARITY header |
| Opus | F6 doc-version-identity-regressed | low | T4 | adopted (dup) | PARITY header |
| Cursor | r8-ledger-unitemized | med | T5 ledger itemization | adopted — R8/R9 tables + index + count rule | LEDGER |
| Opus | F4 ledger-r8-not-itemized | med | T5 | adopted (dup) | LEDGER |
| GLM | ledger-r8-unitemized | high | T5 | adopted (dup) | LEDGER |
| GLM | ledger-rounds-index-stale | low | T5 | adopted — index through R10 (completed v11) | LEDGER header |
| Opus | F2 n6-union-not-propagated-to-shim | med | T6 shim spec | adopted — §3.1 union | PLAN §3.1 |
| Opus | F3 updater-constructor-wrong-frontend-and-grep-hole | med | T7 native updater | adopted — native.New + extended grep | PLAN N1 |
| Opus | F5 hold-cycle-text-spliced-and-misanchored | low-med | T8 splice repair | adopted — block rewritten, literal grep | PLAN §3.7 |
| GLM | m1-crossref-claim-unimplemented | low | T9 SSOT wording | adopted — checklist = gate SSOT | PLAN M1 |

Gate-9 rows (m1-gate-9-na-overbroad / m1-gate-9-waiver-overbroad / m1-gate-na-branch-overbroad-for-item-9) were initially itemized here in error — raw-output grep of all three R10 reviewer files (cursor-plan-v10 / claude-plan-v10 / opencode-plan-v10) attributes them to R10 T1; rows moved out, R9 = 14 labeled (4/6/4) = 8 themes. Opus F4 relabeled to its raw-output label (`ledger-r8-not-itemized`).

## Round 10 → v11 — itemized (17 labeled findings = 6 themes: Cursor 5 + Opus 7 + GLM 5)

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 gate #9 split 9a/9b | Cursor m1-gate-9-na-overbroad (high); Opus m1-gate-9-waiver-overbroad (med); GLM m1-gate-na-branch-overbroad-for-item-9 (medium) | adopted — 9a unconditional, 9b+12 conditional | PLAN M1 |
| T2 acceptance N1 split mirror | Cursor acceptance-n1-ignores-opt-in-split (med) | adopted — mirrored | PLAN acceptance |
| T3 acceptance CI branches on runner decision | Cursor acceptance-ci-contradicts-no-ci-mode (med) | adopted | PLAN acceptance |
| T4 socket path recipe (tilde + run/ mkdir) | Cursor socket-path-tilde-and-run-mkdir-undercorrected (med) | adopted — Join(dataRoot,"run") + MkdirAll | PLAN §3.7 |
| T5 ledger bookkeeping | Cursor ledger-count-and-index-still-diverged (low-med); Opus ledger-r9-not-itemized (med); Opus ledger-r8-itemization-incomplete + round-mislabels (low-med); Opus ledger-rounds-index-stale regression (low); Opus r3-count-self-contradiction (low); GLM r3-revision-count-stale-in-plan-header (low); GLM r9-count-arithmetic-does-not-sum (low); GLM ledger-rounds-index-still-stops-at-r5 (low); GLM r9-ledger-unitemized (low) | adopted — R3=19 everywhere; R9 itemized (15→9); R8 meta rows relabeled (18 labeled); index through R10 | LEDGER, PLAN headers |
| T6 updater init sequence + nettop host | Opus updater-init-sequence-underspecified (low-med); Opus nettop-pref-on-host-set-wrong (low-med) | adopted — pinned 8-step main() order + M1 test; feed-host-pinned egress | PLAN N1, §3.6 |

## Round 11 → v12 (13 themes + 1 meta from 20 labeled findings: Cursor 4 + Opus 10 + GLM 6)

Context: first R11 launch lost Cursor (0B) + GLM (killed mid-stream) to a session interruption — both BLOCKED, relaunched fresh on the same frozen input (8b5e7c8). Opus completed in the first launch. All three verdicts: FINDINGS.

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 N1 pref=on leg vacuous + host set wrong + AutomaticDownloads egress | Opus F1 n1-nettop-pref-on-leg-vacuous-and-host-set-wrong (med); Cursor acceptance-n1-m1-vs-gate12-vacuous-na (high); Cursor n1-pref-on-egress-omits-update-asset-hosts (med) | adopted — enumerated allowed-host set (test feed w/ signed manifest, asserted ≥1 connection + LastCheck; GitHub multi-host verified in internal/update/github.go); AutomaticDownloads pinned false on nettop leg; acceptance mirrors waiver branch (never vacuous) | PLAN N1, acceptance |
| T2 nettop pref=off proof lost when #12 waived | Opus F2 m1-n1-proof-lost-when-item-12-waived (med) | adopted — new unconditional gate item 16 (pref=off zero-outbound + no-download); only pref=on leg waivable | PLAN M1 |
| T3 LV-loss branch contradicts 9a/#11 + probe ordering | Opus F3 lv-loss-branch-contradicts-9a-and-#11 (med) | adopted — LV-loss no longer free (voids 9a injection check + #11; requires SHA-256 dylib pin + explicit deviation); LV probe moved PRE-M1 | PLAN §3.6, PRE-M1(c) |
| T4 N6 TOCTOU between union and read | Opus F4 n6-toctou-type-enum-vs-data-read (med) | adopted — changeCount re-check protocol + `concealed-mid-read` test | PLAN N6, M1 #5 |
| T5 SSOT checklist omits promised tests | Opus F5 m1-checklist-ssot-omits-tests-promised-in-§3 (low-med) | adopted — own-write matrix tests in #2; main()-order test in #12; grep gates in CI line | PLAN M1 |
| T6 updater assert channel + opt-in mechanism | Opus F6 packaged-updater-assert-has-no-readout (low); GLM m1-updater-optin-mechanism-undefined (low-med) | adopted — assert on persisted `<dataRoot>/updater.json`; opt-in via settings-store API test seam | PLAN M1 #12 |
| T7 clearSystemClipboard misbound to quit only | Opus F7 clearSystemClipboard-misbound (low-med) | adopted — verified: guards pasteboard-clear inside History.clear()/clearAll() (Clipboard.swift:141-144; History.swift:242/264); quit hook runs full History.clear() when clearOnQuit; PARITY row split (M2/M4) | PARITY §1, PLAN M2 |
| T8 forgeable own-write marker | Opus F8 own-write-marker-forgeable-bypass-unspecified (low) | adopted — marker gates nothing (full ignore chain + N6 before merge); `forged-marker-concealed` test | PLAN §3.2 |
| T9 ledger counts still diverge (R8=20 residues; R9 arithmetic) | Opus F9 ledger-counts-still-diverge (low); Cursor r10-t5-r8-count-disposition-incomplete (med); GLM ledger-r9-itemized-arithmetic-broken (low-med); GLM plan-r8-count-still-diverges (low) | adopted — R8=18 everywhere; R9 recounted from raw outputs: gate-9 trio → R10 (grep of all three R10 files), Opus F4 relabeled `ledger-r8-not-itemized`; R9 = 14 labeled (4/6/4) = 8 themes; GLM's "reclassify 2 overflow rows" alternative rejected — rows proven R10's by raw-output labels | PLAN L7/L63, LEDGER |
| T10 review packet template stale | Opus F10 packet-template-stale (low, meta) | adopted — packet regenerated from LEDGER for R12 | /tmp packet |
| T11 N6 "first item only" rationale false; union framing | Cursor n6-first-item-only-rationale-contradicts-maccy-241 (med) | adopted — verified Clipboard.swift:173-176 + #241: `pasteboard.types` = superset incl. non-item types; rationale rewritten; union kept as ≥Maccy hardening (recorded, not parity) | PLAN N6/§3.1, PARITY §1 |
| T12 ignored-types default unverified | GLM ignored-types-default-unverified (med) | adopted — default verified = exactly 5 entries, no `com.apple.is-remote-clipboard` (Defaults.Keys+Names.swift:39-46); assert added to gate #15 + PARITY row; §3.3 clause evidence-backed | PLAN §3.3, M1 #15, PARITY §1 |
| T13 step① fchmod drift | GLM n1-step1-fchmod-divergence (low) | adopted — ① = MkdirAll(0700) + fchmod 0700 if pre-existing | PLAN N1 |
| T14 grep third clause non-greppable | GLM n1-grep-third-clause-nongreppable (low) | adopted — `updater.New(` banned outright; `native.New(` clause via statement-level slurped CI check | PLAN N1 |

All 20 labeled findings adopted; no rejections this round. Packet template (T10) was an implementer artifact, corrected for R12.

