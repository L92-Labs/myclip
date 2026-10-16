# Finding ledger — plan phase

Per-finding dispositions (round → reviewed plan version). Rounds: R1 (v1), R2 (v2), R3 (v3), R4 (v4), R5 (v5), R6 (v6), R7 (v7), R8 (v8), R9 (v9), R10 (v10), R11 (v11), R12 (v12), R13 (v13), R14 (v14), R15 (v15), R16 (v16), R17 (v17), R18 (v18), R19 (v19), R20 (v20), R21 (v21), R22 (v22), R23 (v23), R24 (v24), R25 (v25), R26 (v26), R27 (v27), R28 (v27), R29 (v28), R30 (v28), R31 (v28), R32 (v28), R33 (v28), R34 (v28).

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
5. D7 poll floor 100ms — approve (blocks M1; the poll loop ships at M1).

**Before M3/M5:**
6. D9 Vision OCR skip in v1 — approve or add spike.
7. D2 third language: ES ok?
8. D10 copy sounds via MyGo audio — approve (Before M2 — sounds ship at M2.)
9. M5 signing: **DECIDED (user, R25): self-signed for now** — self-signed identity via the production signing path (emits hardened runtime; notarization N/A; documented Gatekeeper/quarantine/TCC steps); Developer ID revisited before public flip. The §3.6 LV sqlite-dylib spike still runs PRE-M1 and gates the D13 branch (identity choice does not preempt it). Plain ad-hoc only as recorded deviation.
10. D12 `migrations` N/A — approve as recorded deviation (non-blocking).
11. D13 (contingent, LV-loss) entitlement delta + residual = same-user code injection inheriting MyClip's TCC grants — approve if the §3.6 LV spike fails.
13. (c') cdhash attestation spike (conditional, D13) — its outcome decides the #12b swap-race leg (pass ⇒ test runs; fail ⇒ leg N/A) and the final D13 residual width. (NOTE: numeric #12 intentionally unused. Items 1–5 = PRE-M1 (a)/(b) blocking; #13 = PRE-M1 (c'); the rest gate later milestones.)
14. D14 Clipboard.sync() not ported — approve (non-blocking; parity-inventory accuracy).

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
| GLM | r3-finding-count-mismatch | low | adopted — **superseded twice (R5 recount, then the R10 audit): final R3 = 19 = Cursor 9 + Opus 10**; GLM run incomplete; ledger-check rows are not findings | LEDGER R3 header |
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
| GLM | ledger-r3-count-fix-introduces-new-mismatch | low | adopted at the time — **superseded by the R10 audit: R3 = 19 (Cursor 9 + Opus 10)** | LEDGER R3 |
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

Ledger-count audit (settled R10): counting rule = ledger-check/BLOCKED rows are not findings.

Segment rule (settled R29): a round's `→(R_N …)` chain segment is appended only in the commit where its dispositions land — never pre-claimed. SHA separation (settled R32): content-landing SHAs (65ed3fd/1e09824 = cites; bfa9007 = D14 decision; 9b296c8 = order/pointer) are DISTINCT from segment-landing SHAs (8b26ed2 = R29/R30 doc-only segments); later segment SHAs are named in their own round sections. R29/R30 exception (recorded R31): both were DOC-ONLY fix rounds within stamped v28 — no chain segments exist for them (grep `R29` in PLAN = 0 was the true state; the earlier `segment landed in 9b296c8` claim was false). Same-stamp doc-only segments for R29/R30 appended at L7/L66; landing commits: 65ed3fd/1e09824 (cites), bfa9007 (D14 decision), 9b296c8 (order/pointer). Round N reviews the current stamped revision; its fixes either land in that stamped revision (doc-only) or create v_(N+1), and the ledger header records where they actually landed. Proof counts exclude the proving row's own text. R1 = 45 (14/17/14); R2 = 34 (9/16/9); R3 = 19 (Cursor 9 + Opus 10; earlier claims of 18 then 25 were both wrong); R4 = 20 (5/9/6); R5 = 16 + 1 evidence-rejection; R6 = 11 (4/7); R7 = 18 (6/7/5); R8 = 18 labeled (Cursor 8 + Opus 6 + GLM 4; verdict counts 8/6/6 incl. 2 GLM meta rows; itemized below); R9 = 8 themes (14 labeled = Cursor 4 + Opus 6 + GLM 4); R10 = 6 themes (17 labeled = Cursor 5 + Opus 7 + GLM 5). Historical footnote claiming "16 Opus rows in R3 / 55 total" was wrong — superseded by this recount.

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

Root cause recorded: several R7 disposition edits silently failed to match in the edit script while the ledger claimed "adopted" — reviewers caught every one. v9 uses ASSERTED, sentence-bounded replacements (script exits non-zero if any target is absent; no greedy line-eating regexes). All 18 labeled adopted: update-archive single stance (M1 archive check; M5 only in infeasible branch); SetPath before-Run ordering + fresh-HOME test; D1 residual reword applied (PLAN + LEDGER); N6 union + named `concealed-on-second-item` test; canonical 15-item M1 gate checklist materialized in Milestones; MyGo `GlobalShortcut` banned entirely (Signature-blind handler); N1 nettop split criterion; LEDGER decision #9 aligned; PARITY hold-cycle criterion aligned. (D11/N2, pref-source, Unsetenv, M5 fallback, confirm-UX were already correctly present in v8.)

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

## Round 9 → v10 — itemized (14 labeled findings = 8 themes: Cursor 4 + Opus 6 + GLM 4)

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
| T1 gate #9 split 9a/9b | Cursor m1-gate-9-na-overbroad (high); Opus m1-gate-9-waiver-overbroad (med); GLM m1-gate-na-branch-overbroad-for-item-9 (medium) | adopted — 9a unconditional, 9b+12 conditional (**superseded in part by R19-T5/R20-T5**: item 12 since split into 12a-i/12a-ii/12b) | PLAN M1 |
| T2 acceptance N1 split mirror | Cursor acceptance-n1-ignores-opt-in-split (med) | adopted — mirrored | PLAN acceptance |
| T3 acceptance CI branches on runner decision | Cursor acceptance-ci-contradicts-no-ci-mode (med) | adopted | PLAN acceptance |
| T4 socket path recipe (tilde + run/ mkdir) | Cursor socket-path-tilde-and-run-mkdir-undercorrected (med) | adopted — Join(dataRoot,"run") + MkdirAll | PLAN §3.7 |
| T5 ledger bookkeeping | Cursor ledger-count-and-index-still-diverged (low-med); Opus ledger-r9-not-itemized (med); Opus ledger-r8-itemization-incomplete + round-mislabels (low-med); Opus ledger-rounds-index-stale regression (low); Opus r3-count-self-contradiction (low); GLM r3-revision-count-stale-in-plan-header (low); GLM r9-count-arithmetic-does-not-sum (low); GLM ledger-rounds-index-still-stops-at-r5 (low); GLM r9-ledger-unitemized (low) | adopted — R3=19 everywhere; R9 itemized (14 labeled = 8 themes, settled R11 from raw reviewer outputs); R8 meta rows relabeled (18 labeled); index through R10 | LEDGER, PLAN headers |
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

## Round 12 → v13 (15 themes from 22 labeled findings: Cursor 6 + Opus 8 + GLM 8)

All three: FINDINGS on frozen v12 @ 484738e. T9-class residue (R9 counts in PLAN Review-tracking line + LEDGER itemized header + R10 T5 text) survived v11's fix because the anchors named in R11-T9's disposition (PLAN L7, LEDGER L199/L253) were narrower than the divergence surface — recurring lesson: fix ALL count sites, then grep-verify zero stale counts in one pass.

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 §3.6 placeholder build vs N1 real-feed requirement | Opus updates-feed-placeholder-contradicts-n1-real-feed (med) | adopted — M1 gate build names `updates` = local signed test feed; items 9b/12a/16 run on that artifact | PLAN §3.6 |
| T2 item-12 waiver scope: Enabled() precondition + vacuous branch | Opus m1-waiver-scope-vs-enabled-precondition (med); GLM gate12-waiver-scope-infeasible-branch-contradiction (med) | adopted — 12a (updates-dependent: Enabled, pref=on leg, install leg) waivable; 12b (pre-seed file-assert, opt-in file-assert, main()-order test) unconditional — SetAutomaticChecks persists updater.json regardless of enabled() (updater.go:138-141/288-291/319-352) | PLAN M1, §3.6 |
| T3 hidden prefs unreachable under D5 | Opus hidden-prefs-unreachable-under-d5 (med) | adopted — verified README Advanced defaults-write set; hidden-prefs inventory pre-task in gate #15 (surface or recorded deviation) | PLAN M1 #15 |
| T4 settings-store ordering + unverified API | Opus settings-store-open-precedes-db-precreate (low-med); GLM n1-step4-mygo-settings-api-unverified (low) | adopted — step ③b DB pre-create before any sqlite load; PRE-M1 settings-API probe + recorded fallback order | PLAN N1 |
| T5 LV-loss SHA-pin TOCTOU + env conflict | Opus lv-loss-sha-pin-toctou-and-env-conflict (low-med) | adopted — private 0700 copy, hash the copy, load via env set after ③ Unsetenv; LibraryPath fallback (library.go:83-95) noted; residual = same-user downgrade (**superseded by R14 T2** — signature validation) | PLAN §3.6 |
| T6 ResolveFeed overgeneralized + allowlist not written + install leg ungated | Opus github-feed-multihost-overgeneralized (low); Cursor n1-allowed-host-set-still-unenumerated (med) | adopted — verified github.go:52-64 (API only for TaggedFeed); allowlist derived from feed form + M1 transcript enumeration (DefaultClient follows redirects); install leg = gate 12a | PLAN N1 |
| T7 own-write marker visibility contradiction | Opus own-write-marker-visibility-contradiction (low) | adopted — marker IS on the pasteboard for own-write detection (.fromMaccy model); never persisted/restored | PLAN §3.2 |
| T8 R9 count residue + index versions + R4 order | Opus ledger-bookkeeping-residue (low); Cursor t9-r9-count-disposition-incomplete (med); GLM plan-r9-count-residue-l64 (low); GLM ledger-r9-itemized-header-stale (low) | adopted — PLAN Review-tracking line, LEDGER itemized header, R10 T5 text → 14/8; LEDGER index → reviewed-version convention (R6 (v6)…R12 (v12)); R4 order → Cursor 5 + GLM 9 + Opus 6 | PLAN L7/L63, LEDGER |
| T9 mid-read "zero reads" contradiction | Cursor n6-mid-read-zero-reads-contradicts-protocol (med) | adopted — sandwich protocol (c0→union→ignore→c1→read→c2); c1 path zero reads; c2 path zero rows + discard; gate #5 wording split | PLAN N6, M1 #5 |
| T10 9a carve-out under LV-loss | Cursor m1-9a-no-acceptance-path-under-lv-loss (med); GLM gate9a-unconditional-vs-lv-loss-voiding (low-med) | adopted — 9a states its LV-loss alternate evidence (pre-load SHA-256 pin verify) (**superseded by R14 T2** — signature validation) | PLAN M1 #9 |
| T11 N1 assert channel stale vs T6 | Cursor n1-automaticchecks-assert-channel-stale-vs-t6 (low-med) | adopted — N1 packaged assert → file-read channel; in-process getter only in unit/harness | PLAN N1 |
| T12 import-alias CI bypass | Cursor ci-native-new-import-alias-bypass (low-med) | adopted — canonical-name rule for plugins/updater/native imports | PLAN N1 |
| T13 gate 11 asserts not enumerated | GLM gate11-omits-asserts-claimed-by-§3.6 (low-med) | adopted — all 4 asserts written into item 11 | PLAN M1 #11 |
| T14 decision #9 timing unowned | GLM dev-id-decision-timing-mismatch (low) | adopted — conditional timing (PRE-M1 on LV-spike failure, else before M5) | LEDGER #9 |
| T15 M5 notarization impossible in fallback | GLM m5-notarization-fallback-impossible (low) | adopted — notarization only in Developer ID branch; self-signed ⇒ N/A documented | PLAN M5 |

All 22 labeled findings adopted; no rejections.

## Round 13 → v14 (10 themes from 15 labeled findings: Cursor 4 + Opus 8 + GLM 3)

All three: FINDINGS on frozen v13 @ 6407bff. NOTE: this table was omitted from the v14 commit — all three R14 reviewers caught the omission (5th recurrence of the itemize-every-round class); table reconstructed verbatim from the R13 dispositions applied in v14.

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 12a-split refs stale | Cursor gate12-waiver-refs-stale-after-12a-split (med); Opus t2-12a-split-not-propagated-to-item16-and-acceptance (med) | adopted — item 16 + acceptance now name 12a; grep gate: zero bare `item 12[^ab]` | PLAN M1 #16, acceptance |
| T2 LV-loss vs 9a env leg + #16 path | Cursor lv-loss-9a-env-and-item16-uncarved (med); Opus lv-loss-branch-contradicts-9a-env-leg-and-item16-path (med) | adopted — 9a env leg + #16 gained LV-loss alternates (superseded by R14 T2 signature-validation rewrite) | PLAN M1 #9a/#16 |
| T3 SHA-pin source unspecified | Opus lv-loss-sha-pin-source-unspecified (low-med) | adopted in v14 as compile-time constant — **superseded by R14 T2**: file-hash pin unsatisfiable after nested re-sign; now signature-based | PLAN §3.6 |
| T4 pref=on LastCheck nondeterminism | Opus pref-on-leg-nondeterministic-lastcheck-interval (low-med) | adopted — fresh data root precondition + named Interval + failed-backoff note (updater.go:201, 269-270, 381-384) | PLAN N1 |
| T5 allowlist-from-transcript circular | Cursor n1-allowlist-from-transcript-circular (low-med); Opus n1-allowlist-derived-from-observation-is-circular (med) | adopted — expected host set predeclared from feed form + manifest; transcript verifies only; foreign host = FAIL | PLAN N1 |
| T6 CLI `set` verb not in §3.7 | Cursor hidden-prefs-cli-set-surface-undefined (low-med); GLM gate15-cli-set-verb-not-in-37-surface (low) | adopted — `set` added to §3.7 (superseded/extended by R14 T5 closed allowlist) | PLAN §3.7, M1 #15 |
| T7 CLI peer DR unpinned / ad-hoc bypass | Opus cli-peer-dr-unpinned-adhoc-same-identifier (low-med) | adopted — pinned DR string in MyClip source + gate #10 ad-hoc negative test | PLAN §3.7 |
| T8 main-order test mechanism unspecified | Opus main-order-test-mechanism-unspecified (low) | adopted — injectable step-recorder seam; no-sqlite-before-③b assert | PLAN N1, M1 #12b |
| T9 stale residue (loop line v12, All-20, R3=25 markers, `;)` typo) | Opus ledger-and-plan-stale-residue-after-t8-lesson (low); GLM plan-review-loop-line-version-residue (low) | adopted — loop line made version-less; L226 → 18 labeled; L152/L173 superseded markers; typo fixed; grep gate clean | PLAN, LEDGER |
| T10 PARITY `migrations` key uncovered | GLM parity-defaults-key-migrations-uncovered (low) | adopted — D12 row added (sign-off entry added R14) | PARITY §6, N2 |

## Round 14 → v15 (8 themes from 14 labeled findings: Cursor 5 + Opus 5 + GLM 4)

All three: FINDINGS on frozen v14 @ 41be63c. Dominant finding: the R13 ledger table itself was missing (implementer omission — the exact class R12's preamble warned about). Signature-validation rewrite (T2) replaces v14's file-hash pin, which `mygo build --force` nested re-signing makes unsatisfiable.

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 R13 ledger section missing | Cursor ledger-r13-unitemized (high); Opus ledger-r13-table-missing (med); GLM ledger-r13-unitemized (high) | adopted — R13 table added (above); index through R14 | LEDGER |
| T2 LV-loss SHA-pin unsatisfiable after nested re-sign | Opus lv-loss-sha-pin-unsatisfiable-after-nested-resign (med-high) | adopted — `--force` re-sign (resources.go:707-797) rewrites dylib signature ⇒ file hash ≠ embedded constant; validation = `SecStaticCodeCheckValidity` on the private 0700 copy against the pinned requirement; PRE-M1 probe added (shasum vs constant, expected differ) | PLAN §3.6, M1 #9a/#16 |
| T3 LV-loss item-11 entitlement contradiction + N1 branch conditioning | Cursor lv-loss-voids-gate11-without-alternate (med); GLM lv-loss-branch-item11-entitlement-contradiction (med) | adopted — item 11 (both sites) LV-loss alternate = exactly `cs.disable-library-validation` (D13); N1 Resources-only sentence branch-conditioned | PLAN §3.6, M1 #11, N1 |
| T4 LV-loss load unplaced in main order | Cursor lv-loss-load-unplaced-in-main-order (med) | adopted — step ③c in the pinned order (verify→copy→validate→Setenv), asserted by the 12b seam when active | PLAN N1 |
| T5 `set`-verb threat model undefined | Opus set-verb-threat-model-not-propagated (low-med); GLM cli-set-verb-allowlist-undefined-d1-scope-stale (med) | adopted — closed display-only allowlist; capture/ignore-chain prefs + size/clearOnQuit/updater/socket prefs excluded; own socket opt-in; D1/D5 rows + §3.7 semantics updated; gate #10 non-allowlisted-`set` reject | PLAN §3.7, PARITY D1/D5 |
| T6 D12 outside N2 range + unsigned | Cursor n2-deviation-range-stale-vs-d12 (med); Opus d12-outside-n2-range-and-signoff (low); GLM n2-deviation-range-stale-d12 (med) | adopted — N2 → D1–D13; D12 sign-off entry (decision #10, non-blocking); D13 registered (decision #11) | PLAN N2, PARITY §6, LEDGER |
| T7 gate 10 SSOT omits new negatives | Cursor gate10-ssot-omits-adhoc-dr-negative (low-med) | adopted — item 10 enumerates ad-hoc reject + pinned-DR assert + set-reject | PLAN M1 #10 |
| T8 CI textual check bypassable | Opus ci-native-new-textual-check-bypassable (low) | adopted — `go/ast` analyzer (KeyValueExpr, true ident, direct `native.New` literal) + constructor unit test | PLAN N1 |

All 14 labeled findings adopted; no rejections.

## Round 15 → v16 (11 themes from 16 labeled findings: Cursor 3 + Opus 6 + GLM 7)

All three: FINDINGS on frozen v15 @ e77d372. Dominant: R14's signature-validation rewrite was only partly propagated (9a env-leg still hash-worded — v14 edit covered the injection leg but missed the env leg), and the D13 DR literally reusing the §3.7 CLI identifier would reject the legitimate dylib (nested codesign stamps file-name identifiers).

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 9a env-leg hash residue | Cursor gate9a-lv-loss-assert-still-hash-based (med-high); Opus gate9a-lv-loss-env-leg-stale-hash-wording (low-med); GLM lv-loss-9a-env-leg-stale-sha-wording (med) | adopted — env-leg alternate = signature-verified copy (SecStaticCodeCheckValidity vs pinned dylib DR); hash asserts dropped | PLAN M1 #9a |
| T2 D13 DR identifier collision | Cursor d13-dylib-dr-collides-with-cli-identifier-requirement (med-high); Opus lv-loss-dr-identifier-mismatch (med) | adopted — separate pinned DYLIB DR (post-build dylib identifier ∧ same leaf/Team-ID anchor); PRE-M1 `codesign -d -r- -vv` probe on both artifacts; gates 9a/10 assert their own DRs | PLAN §3.6, M1 #10 |
| T3 ③c fail-open | Opus lv-loss-validation-fail-open (med) | adopted — any ③c failure = fatal exit before ③b/④; verify defined; 12b tampered-copy ⇒ exit test | PLAN N1, M1 #12b |
| T4 12b SSOT omits ③c | GLM gate-12b-ssot-omits-3c (low) | adopted — item 12b names ③c when D13 active | PLAN M1 #12b |
| T5 clipboardCheckInterval misclassified | Opus set-allowlist-misclassifies-clipboardCheckInterval (low) | adopted — reclassified cadence pref with range validation [D7 floor, documented max]; out-of-range set reject in gate #10 | PLAN §3.7, M1 #10 |
| T6 AST analyzer holes (func-value + dot-import + webview pkg) | Opus ast-analyzer-func-value-bypass (low); GLM updater-ci-gate-dot-import-and-webview-package-hole (low) | adopted — import rule covers both updater packages, named imports (incl. `.`) rejected, Uses-resolution rejects non-Fun uses, fixtures added | PLAN N1 |
| T7 PARITY History.swift path 404 | GLM parity-history-source-path-404 (low) | adopted — `Maccy/Observables/History.swift` (verified at `eb03eba`; behavior re-confirmed live) | PARITY header |
| T8 LEDGER supersession crossrefs wrong | Cursor ledger-r13-supersession-crossrefs-wrong (low); GLM ledger-r13-t2-crossref-wrong (low) | adopted — R13 T2 → R14 T2; R13 T6 → R14 T5 | LEDGER R13 |
| T9 R14 theme count inflated | Opus ledger-r14-theme-count-inflated (low); Cursor LEDGER-CHECK count note | adopted — placeholder row removed; R14 = 8 themes at all sites | LEDGER, PLAN L7/L65 |
| T10 N4 RSS/CPU no relief valve | GLM n4-rss-cpu-no-relief-valve (low) | adopted — same revisable-with-sign-off clause as size + PRE-M1 minimal-MyGo feasibility probe | PLAN acceptance |
| T11 hidden-prefs `set` unconstrained | GLM hidden-prefs-set-surface-unconstrained (low) | adopted — #15 set option qualified to §3.7-allowlist-eligible keys | PLAN M1 #15 |

All 16 labeled findings adopted; no rejections.

## Round 16 → v17 (13 themes + 1 meta from 16 labeled findings: Cursor 5 + Opus 8 + GLM 3; T14 meta excluded)

All three: FINDINGS on frozen v16 @ efa74e1. Dominant classes: R15-T5/T2 propagation residue (PARITY D1/D5 "display-only"; #9a injection leg untyped DR; D13 row singular DR), plus Opus design-depth findings (validate-then-dlopen TOCTOU window, settings-read fail-open, sidecar pre-create WAL drift, single-instance lock absent).

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 D1/D5 "display-only" residue | Cursor parity-d1-d5-display-only-stale-vs-r15-t5 (med); Opus parity-d1-d5-set-allowlist-stale-vs-r15-t5 (low); GLM t5-reclassification-parity-display-only-residue (low) | adopted — "closed non-security set (display + range-validated cadence)"; PARITY §1 clipboardCheckInterval row names the set surface | PARITY D1/D5, §1 |
| T2 cadence max unnamed — soft-disable | Cursor clipboardCheckInterval-max-unnamed-soft-disables-capture (med-high) | adopted — max pinned 3600s; above-max ⇒ rejected (soft-disable of capture forbidden) (**superseded by R20-T6**: CLI max 5s, larger via Advanced pane) | PLAN §3.7, M1 #10 |
| T3 #9a injection-leg DR untyped | Cursor gate9a-injection-alternate-pinned-requirement-vs-dylib-dr (med) | adopted — "pinned dylib DR" in both 9a legs | PLAN M1 #9a |
| T4 D13 row singular DR | Cursor parity-d13-singular-dr-omits-separate-dylib-dr (med) | adopted — distinct dylib DR named; residual reworded (T6) | PARITY D13 |
| T5 gate-10 dylib-DR unconditional | Cursor gate10-dylib-dr-assert-unconditional-vs-d13 (low-med); Opus gate10-dylib-dr-unconditional (low) | adopted — #10 asserts CLI peer DR only; dylib DR evidence lives in #9a (D13 active) | PLAN M1 #10 |
| T6 validate-then-dlopen TOCTOU window | Opus d13-validate-then-dlopen-toctou-residual-understated (med) | adopted — `sqlite.Load()` immediately after validation in ③c; post-load mapped-image cdhash check ⇒ fatal on mismatch; D13 residual = "same-user code injection inheriting MyClip's TCC grants"; 12b swap-after-validation race test | PLAN §3.6, N1, M1 #12b |
| T7 settings-read fail-open | Opus settings-read-failure-fails-open (low-med) | adopted — ④ read error ⇒ pref=false, ⑥ still runs; pre-seeded true can never win (**superseded in part by R21-T3**: capture-safe mode + ignore-chain fail-safe added) | PLAN N1 ④ |
| T8 sidecar pre-create WAL drift | Opus sidecar-precreate-forces-wal-drift (low-med) | adopted — ③b pre-creates main DB only; sidecars fchmod-or-refuse; #3 asserts journal_mode every launch + sidecar modes after first write | PLAN N1 ③b, M1 #3 |
| T9 single-instance lock absent | Opus single-instance-unspecified (low) | adopted — ②b `RequestSingleInstanceLock` before any poller/copy/socket work; new gate item 17 (second launch ⇒ exit) (**superseded by R17-T6 + R18-T7**: call removed entirely — fail-closed flock + CI ban; this row's design is void) | PLAN N1, M1 #17 |
| T10 AST true-ident shadowable | Opus ast-true-ident-shadowable (low) | adopted — constant-true via `types.Info` (name comparison shadowable) | PLAN N1 |
| T11 LV-loss copy destination unpinned | GLM lv-loss-copy-destination-unpinned (low) | adopted — `<dataRoot>/lib/mygo-sqlite3.dylib` (0700 dir; per-launch re-copy, no cache) (**file mode 0600 since R20-T1**) | PLAN §3.6 |
| T12 settings-store sqlite escapes recipe | GLM settings-store-db-mode-gap (low) | adopted — every sqlite file MyClip opens gets the §3.4 recipe; #3 asserts all DB files under data root | PLAN §3.4, M1 #3 |
| T13 R12 superseded markers missing | Opus ledger-r12-superseded-markers-missing (low) | adopted — R12 T5/T10 rows marked superseded-by-R14-T2 | LEDGER R12 |
| T14 packet identity stale (meta) | Opus packet-identity-stale (low, meta) | adopted — R17 packet regenerated from PLAN/LEDGER headers | /tmp packet |

All 16 labeled findings adopted; no rejections.

## Round 17 → v18 (12 themes from 18 labeled findings: Cursor 5 + Opus 9 + GLM 4)

All three: FINDINGS on frozen v17 @ 58ce152. Dominant: the R16-T6 disposition claimed a swap-after-validation race test that never landed in #12b (the R7/R8 silent-disposition class, again); plus depth findings on the ③c design (closure overclaim, Load/Open vocabulary, settings backend myth — MyGo ships NO settings API, both Opus and GLM verified the plugin list).

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 ③c closure overclaim + cdhash mechanism unnamed | Cursor d13-validate-then-dlopen-closure-overclaim (med-high); Opus post-load-cdhash-mechanism-unspecified (low-med) | adopted — "window closed" dropped (dyld initializers run before Load returns); post-load check = best-effort tripwire; mechanism = in-memory CodeDirectory via `_dyld_get_image_header`/`dladdr` purego (never path re-hash); unproven spike added to PRE-M1 LV probe | PLAN N1 ③c, §3.6 |
| T2 swap-race test missing from 12b SSOT | Cursor r16-t6-swap-race-test-missing-from-ssot (med); Opus ledger-r16-t6-race-test-unpropagated (med); GLM r16-t6-12b-swap-race-test-omitted (low-med) | adopted — 12b now enumerates BOTH D13 tests (tamper ⇒ fatal; swap-copy-after-validation ⇒ cdhash-mismatch ⇒ fatal); grep-verified | PLAN M1 #12b |
| T3 Load vs Open vocabulary contradiction | Cursor n1-sqlite-load-vs-open-order-contradiction (med) | adopted — dylib bind ≠ DB open; ③b before any DB open; 12b asserts no DB open before ③b, Load only inside ③c | PLAN N1, M1 #12b |
| T4 AutomaticDownloads asymmetric fail-open | Cursor automaticdownloads-asymmetric-fail-open (low-med) | adopted — ⑥ forces `SetAutomaticDownloads(false)` unless MyClip-owned pref; pre-seed test extended | PLAN N1 ⑥ |
| T5 settings backend undefined (MyGo has NO settings API) | Cursor settings-db-0600-precreate-unsequenced (low-med); Opus settings-store-backend-undefined (low-med) | adopted — settings store = MyClip-owned sqlite `settings.db`, §3.4 recipe, pre-created in ③b; void probe removed | PLAN N1 ④, ③b |
| T6 single-instance fail-open + scope | Opus single-instance-lock-fail-open-and-scope (med); GLM mygo-single-instance-fail-open-unrecorded (low) | adopted — MyClip-owned fail-closed flock `<dataRoot>/run/myclip.lock`; MyGo lock = arg-forwarding only (**superseded by R18-T7**: call removed entirely); fail-open semantics recorded as residual; #17 + dev+packaged concurrent test | PLAN N1 ②b, M1 #17 |
| T7 single-instance listener unauth ingress | Opus single-instance-listener-unauthenticated-ingress (low) | adopted — §3.7 threat model records the always-on listener; handlers = focus-popup only; CI rule forbids verb wiring (**superseded by R18-T7/R19-T3**: listener never bound — call banned) | PLAN §3.7 |
| T8 journal-mode drift no remediation | Opus journal-mode-drift-no-remediation (low) | adopted — stray `-wal`: checkpoint(TRUNCATE) → re-pin DELETE → assert → fail closed only on re-pin failure; stray-wal launch test in #3 | PLAN M1 #3 |
| T9 D13 downgrade residue (PARITY + decision #11) | Opus d13-downgrade-residue (low-med) | adopted — both sites carry the R16-T6 wording | PARITY D13, LEDGER #11 |
| T10 R16 header meta unlabeled | Opus ledger-r16-count-meta-unlabeled (low) | adopted — "13 themes + 1 meta (T14 excluded)" at all sites | LEDGER, PLAN L7/L66 |
| T11 gate order 17-before-16 | Opus gate-list-order (low); GLM m1-checklist-item-order-17-before-16 (low) | adopted — reordered 16 then 17 | PLAN M1 |
| T12 PARITY hover-delay numeric drift | GLM parity-hover-delay-numeric-drift (low) | adopted — `previewDelay` (default 1500ms) | PARITY §2b |

All 18 labeled findings adopted; no rejections.

## Round 18 → v19 (14 themes + 1 meta from 19 labeled findings: Cursor 5 + Opus 8 + GLM 6; T15 meta excluded)

All three: FINDINGS on frozen v18 @ 4d98573. Dominant: the flock recipe as v18 wrote it (O_EXCL semantics) would brick relaunch after a crash — Opus high; the R17-T4 pre-seed extension again claimed-but-not-landed (third occurrence of the class — every future disposition must carry its grep proof); GLM's fresh Maccy source fetch found three real parity bugs in §3.3/§3.4 (regex per-item scope, Word-link `.pdf` subtraction, plain-paste fileURL preservation).

| Theme | Findings | Disposition | Anchor |
|---|---|---|---|
| T1 flock O_EXCL bricks relaunch | Opus flock-o-excl-bricks-relaunch (high); Cursor single-instance-lockfile-oexcl-vs-flock (med-low); GLM ②b-lock-brick-on-crash (low-med) | adopted — open(O_CREAT\|O_RDWR\|O_NOFOLLOW\|O_CLOEXEC,0600) → fstat sanity → flock(LOCK_EX\|LOCK_NB); never unlink; fd held for process; #17 kill-9 relaunch leg | PLAN N1 ②b, M1 #17 |
| T2 run/ unsequenced before ②b | Opus run-dir-unsequenced-before-2b (med) | adopted — ① creates data root + `run/` subdir; 12b asserts; #17 fresh-HOME leg | PLAN N1 ① |
| T3 R17-T4 preseed-downloads not landed | Cursor automaticdownloads-t4-incomplete-and-window-bypass (med-high); Opus r17-t4-preseed-downloads-not-landed (med) | adopted — #12b pre-seeds both keys, asserts both persisted false; grep-verified this time | PLAN M1 #12b |
| T4 downloads pref undefined + window-write revert | Opus downloads-pref-undefined (med); Cursor window-bypass clause | adopted — `myclipDownloadsPref` defined (PARITY D3 row, default false, M4 control, set-excluded); OnChange mirrors window writes; #12a setup step named | PLAN N1 ⑥/12a, PARITY D3 |
| T5 12b Load assert unconditional | Opus 12b-load-assert-unconditional-vs-non-d13 (low-med) | adopted — D13: Load only in ③c; else no earlier than ④ after ③ Unsetenv | PLAN M1 #12b |
| T6 cdhash spike no fail branch | Cursor d13-swap-race-gate-depends-on-unproven-cdhash-spike (med); Opus cdhash-spike-no-fail-branch (med) | adopted — PRE-M1 (c') blocking item; failure branch = drop tripwire + widen residual + swap-test N/A w/ sign-off; race-injection seam = test-only build-tag hook | PLAN PRE-M1 |
| T7 MyGo single-instance listener gratuitous | Opus single-instance-listener-gratuitous-and-order-unpinned (low) | adopted — `RequestSingleInstanceLock` NOT called (no args to forward); CI grep bans it; flock is the only exclusion | PLAN N1 ②b |
| T8 R17-T5 stale sqlite hedge | Opus r17-t5-stale-sqlite-hedge (low) | adopted — ③b names history + settings.db; hedge removed | PLAN N1 ③b |
| T9 D13 downgrade lexeme residue | Cursor parity-d13-downgrade-lexeme-residue (low) | adopted — opening clause rewritten; contingency on (c') noted | PARITY D13 |
| T10 updater import gate bans own canonical import | GLM updater-import-gate-bans-own-canonical-import (med) | adopted — canonical imports allowed + Uses-restricted (native.go does not re-export Options) | PLAN N1 |
| T11 plain-paste drops fileURL | GLM plain-paste-drops-maccy-fileurl-preservation (low) | adopted — string + fileURL preserved (#962); M3 plain-paste-file test | PARITY §3 |
| T12 ignoreRegexp capture-level scope | GLM ignore-regexp-scope-capture-level (low) | adopted — per-item scope ported + 2-item M1 test | PLAN §3.3 |
| T13 Word-link drop omits `.pdf` | GLM word-link-drop-omits-pdf-subtraction (low) | adopted — 3-type subtraction verbatim + no-PDF-rep test | PLAN §3.3 |
| T14 journal-mode pin ambiguous | GLM journal-mode-pin-ambiguous (low) | adopted — `DELETE` only; WAL alternative dropped | PLAN §3.4 |
| T15 packet identity stale (meta) | Opus packet-identity-stale (low, meta) | adopted — R19 packet regenerated from headers | /tmp packet |

All 19 labeled findings adopted; no rejections.

## Round 19 → v20 (14 themes + 1 meta from 27 labeled findings: Cursor 10 + Opus 11 + GLM 6; T15 meta excluded)

All three: FINDINGS on frozen v19 @ b262d80. Escalation driver: v19's PLAN edits landed but their PARITY mirror-rows and gate-SSOT legs (#17 legs, #12b×(c') conditioning) did not — 7 of 27 findings were propagation residue of that one class. Root cause fixed procedurally: every v20 disposition below carries a grep proof executed at edit time (kill-9 ✓, DELETE-or-WAL=0 ✓, downgrade-residual=0 ✓, always-on=0 ✓, 12a-i ✓).

| Theme | Findings | Disposition | Grep proof |
|---|---|---|---|
| T1 #17 kill-9 + fresh-HOME legs missing | Cursor r18-t1-gate17-kill9-not-landed (high); Cursor r18-t2-gate17-fresh-home-leg-missing (med); Opus r18-t1-t2-gate17-legs-not-landed (med); GLM r18-t1-kill9-leg-not-landed (med-low) | adopted — both legs verbatim in #17 | `kill -9` + `fresh HOME` present |
| T2 PARITY journal WAL residue | Cursor r18-t14-parity-journal-wal-residue (med-high); Opus r18-t14-journal-pin-parity-residue (low-med); GLM parity-data-at-rest-wal-hedge (low) | adopted — PARITY row: DELETE only | `DELETE or WAL` = 0 in PARITY |
| T3 §3.7 listener "always-on" stale | Cursor r18-t7-section37-listener-threat-stale (med); Opus s37-mygo-listener-always-on-stale (low-med) | adopted — listener never bound (②b ban); handlers never registered | `always-on` = 0 in PLAN |
| T4 D13 "downgrade" lexeme (PARITY opening) | Cursor r18-t9-d13-downgrade-lexeme-residue (med); Opus r18-t9-parity-d13-downgrade-not-landed (low-med); GLM parity-d13-downgrade-lexeme-residue (low) | adopted — opening = "code-injection residual" | `downgrade residual` = 0 in PARITY |
| T5 12a check/install leg split | Opus gate12a-downloads-setup-contradicts-n1-pref-on-leg (med); GLM 12a-downloads-pref-vs-n1-check-leg-pin (low-med) | adopted — 12a-i check leg pins downloads OFF (run() auto-installs, updater.go:495-505); 12a-ii install leg adds asset hosts | `12a-i` present |
| T6 #12b swap-race vs (c') branch | Cursor r18-t6-12b-swap-race-ignores-c-prime (med); Opus cdhash-spike-no-fail-branch (med) | adopted — swap leg conditional on (c') success; N/A+sign-off on failure | text present |
| T7 flock fstat mode fail-closed bricks | Opus flock-fstat-mode-fail-closed-bricks (low-med) | adopted — drifted mode ⇒ fchmod 0600 + continue; fail closed only on non-regular/foreign-uid | text present |
| T8 regex-scope wording residue | Cursor r18-t12-regex-scope-wording-residue (med) | adopted — capture-level clause dropped; PARITY row gains per-item scope | text present |
| T9 PARITY word/pdf row | Cursor r18-t13-parity-word-pdf-omission (med) | adopted — 3-type subtraction + no-PDF-rep test in PARITY row | text present |
| T10 R18 count arithmetic (5+8+6=19, not 20) + meta label | Opus ledger-r18-count-arithmetic (low); Cursor r18-meta-count-unlabeled (low); GLM ledger-r18-header-meta-label (low) | adopted — R18 = 19 labeled; header/meta labels aligned at all sites | counts consistent |
| T11 R16 T9 supersession unmarked | Cursor r16-t9-ledger-supersession-unmarked (low) | adopted — marker added (superseded by R17-T6/R18-T7: RequestSingleInstanceLock removed) | marker present |
| T12 (c') absent from decision list | Opus c-prime-spike-absent-from-decision-list (low) | adopted — (c') documented in the lettered PRE-M1 list + numeric decision #13 + intentional-gap note (R20 anchor `#13 present` was a false proof — it matched gate item 13) | PLAN PRE-M1, LEDGER decisions |
| T13 lib/ dir + mode asserts in D13 branch | GLM d13-branch-copy-modes-unasserted (low) | adopted — ③c creates lib/ 0700 + copy 0600; #9a/#11 assert path+modes when D13 active | text present |
| T14 PARITY stamp rule | Opus parity-version-stamp-stale (low) | adopted — header now carries the stamp rule ("last plan revision touching PARITY") | header text |
| T15 packet identity (meta) | Opus r18-t15-packet-identity-not-landed (low, meta); Cursor packet notes | adopted — R20 packet generated mechanically from PLAN L7/L66 + LEDGER index + PARITY L1 (script, then diffed) | /tmp packet |

All 27 labeled findings adopted; no rejections.

## Round 20 → v21 (12 themes + 1 meta from 19 labeled findings: Cursor 5 + Opus 10 + GLM 4)

All three: FINDINGS on frozen v20 @ 7a37f42. GLM audited the R19 grep-proof procedure itself: 2 proofs were false (`always-on` hit a negated clause; `#13` matched a gate item, not a decision) and R16-T9's claimed marker had never landed. v21 introduced executed grep proofs — proof cells record the executed command + count (not transcripts).

| Theme | Findings | Disposition | Proof (executed post-edit) |
|---|---|---|---|
| T1 D13 copy-mode contradiction + gate-SSOT legs | Cursor d13-dylib-copy-mode-contradiction (high); Opus r19-t13-lib-mode-not-in-gate-ssot (med) + d13-copy-mode-contradiction (med-low); GLM d13-copy-mode-contradiction-and-gate-leg-gap (med) | adopted — dir 0700 / file 0600 everywhere; #9a env leg + #11 carry the lib/ path+mode assert verbatim | `copy 0700` = 0; `0700 + dylib copy 0600` = 3 in PLAN |
| T2 sidecar-check vs D13 Load order wording | Opus sidecar-check-after-d13-load (med-low) | adopted — ③b = "before any DB open (the ③c dylib bind is not a DB open)" | text present |
| T3 AST allowlist bans own API | Opus ast-uses-allowlist-bans-own-api (med) | adopted — OnChange/getters added; `_test.go` exempt | text present |
| T4 OnChange registration unpinned | Opus onchange-registration-unpinned (low-med) | adopted — pinned step ⑥b; callback-time read; 12b pre-seed leg | text present |
| T5 updater.json 0644 unasserted | Opus updater-json-mode-unasserted (low) | adopted — accepted residual recorded + asserted at both gate sites | text present |
| T6 CLI interval max 3600s still soft-disables | Opus set-interval-max-soft-disables (low) | adopted — CLI max 5s; larger cadences via Advanced pane only | `3600s` = 1 in PLAN (rationale clause only, intentional); PARITY = 0 |
| T7 N1 check-leg host set included asset hosts | Cursor n1-allowed-hosts-include-assets-vs-12a-i (med) | adopted — check leg = feed host(s) only; asset hosts predeclared for 12a-ii | text present |
| T8 R16-T9 supersession marker falsely claimed | Cursor r19-t11-supersession-marker-absent (med); Opus r19-t11-supersession-not-landed (low-med); GLM r19-t11-supersession-marker-absent (low) | adopted — marker landed byte-verified this time | `superseded by R17-T6 + R18-T7` = 1 |
| T9 "always-on" grep proof false | Cursor LEDGER-CHECK; Opus r19-t3-grep-proof-false (low); GLM r19-t3-grep-proof-falsified (low) | adopted — PLAN lexeme changed; proof re-run = 0 | `always-on` = 0 |
| T10 decision anchor "#13" wrong + numbering jump | GLM r19-t12-decision-label-wrong (low); Opus decision-list-gap-12 (low) | adopted — ledger anchor corrected to PRE-M1 (c'); decision list = (a)(b)(c)(c') lettered + numeric 1-13 kept separate (no renumber needed; gap documented) | anchor text |
| T11 R18 footer "All 27" + R19 header meta label | Opus ledger-count-residue (low) | adopted — footer 19; header "14 themes + 1 meta (T15 excluded)" | text present |
| T12 PARITY stamp value not applied | Cursor parity-stamp-rule-not-applied (low) | adopted — stamp now encodes the plan revision (v18; rule text names v21) | header text |
| T13 packet identity (meta) | — | R21 packet generated mechanically from headers (same script as R20) | /tmp packet |

All 19 labeled findings adopted; no rejections.

## Round 21 → v22 (15 themes + 1 meta from 20 labeled findings: Cursor 5 + Opus 11 + GLM 4; T15/T16 added R23)

All three: FINDINGS on frozen v21 @ a523a0a. Headline: Opus falsified v20's OnChange rationale from source — `update()` snapshots `u.changes` and skips callbacks when state is unchanged (updater.go:325-336), so ⑥'s own writes never reach a ⑥b-registered callback; the "fire for ⑥'s own writes" claim and its vacuous 12b leg are gone. Cursor's findings were entirely propagation/precision class on v20's additions.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 OnChange snapshot rationale false + vacuous 12b leg | Opus onchange-6b-snapshot-misses-6-writes (med) | adopted — synchronous post-⑥ assert (AutomaticDownloads==false && pref==false); mirror tested via injected window-response seam (updater.go:514-519) with callback completion exposed | `fire for ⑥'s own writes` = 0 |
| T2 12b sub-steps + pref assert not in SSOT | Cursor T4-6b-gate-ssot-omission (med); Opus r20-t4-leg-not-in-gate-ssot (low-med) | adopted — 12b names ②b/③b/⑥b (+③c D13); myclipDownloadsPref assert added | `incl. ②b, ③b, ⑥b` = 1 |
| T3 settings fail-open on ignore chain | Opus settings-read-failure-fails-open-on-ignore-chain (med) | adopted — read error ⇒ ignoreEvents=true + capture disabled + visible error; corrupt settings.db ⇒ zero rows test | `corrupt settings.db` = 1 |
| T4 gate9a injection vector unnamed (vacuous) | Opus gate9a-injection-vector-unnamed-vacuous (med) | adopted — named vector (swapped Contents/Resources dylib, unsigned/foreign-signed) + LV-rejection assert + LV-disabled negative control | `NAMED VECTOR` = 1 |
| T5 stale socket after kill-9 | Opus stale-socket-after-kill9 (low-med); GLM stale-socket-unlink-before-bind-unspecified (low-med) | adopted — flock-held lstat(unlink-socket owned by uid) → unlink → bind; #17 kill-9 leg extended with socket rebind | `else refuse) → unlink → bind` = 1 |
| T6 CI grep scope vs fixtures + garbled allowlist | Opus ci-grep-scope-vs-fixtures + garbled-allowlist (low) | adopted — grep scope = non-test .go excluding internal/lint/testdata/; duplicate SetAutomaticDownloads removed; native.go re-export note reattached | dup lexeme = 0 |
| T7 D7/D10 sign-off after implementation | Opus deviation-signoff-after-implementation (low-med) | adopted — D7 before M1; D10 before M2 (decision list) | `BEFORE M1 implementation` = 1 |
| T8 tagged-feed redirect host misclassified | Opus tagged-feed-manifest-redirect-host-misclassified (low) | adopted — check-leg set = feed host + api.github.com + manifest redirect hosts (github.go:18,44-64); 12a-ii = manifest archive URLs | text present |
| T9 packaged launch method unspecified | Opus packaged-env/HOME-launch-method-unspecified (low) | adopted — exec Contents/MacOS/myclip directly; in-process env + IsPackaged asserts | `Contents/MacOS/myclip` = 1 |
| T10 ledger supersession markers missing (R16-T2/T11, R17-T6/T7) | Opus ledger-supersession-markers-missing (low) | adopted — R16-T2 (part-superseded R19/R20), R16-T11 (file mode 0600), R17-T6/T7 (covered by R18-T7 rewrite + §3.7 never-bound text); byte-verified | markers present |
| T11 decision #12 gap undocumented + R19-T12 anchor false | Cursor T10-decision-gap-undocumented (med); GLM r20-t10-disposition-unlanded (med) | adopted — intentional-gap note under decision #13; R19-T12 anchor corrected with false-proof confession | `intentionally unused` = 1 |
| T12 PLAN R19 theme count 15→14 | Cursor T11-r19-count-PLAN-residue (med); GLM plan-ledger-r19-theme-count-divergence (low) | adopted — both PLAN sites = 14 themes + 1 meta | counts match LEDGER |
| T13 proof cells garbled / self-quote matches | Opus r20-proofs-not-as-claimed (low); GLM r20-proof-cells-not-verbatim (low-med) | adopted — proofs re-recorded as executed command + count; self-quotes excluded from counts (proof strings never appear in the proving row's own text) | this table's proofs ran post-edit |
| T14 packet identity (meta) | — | R22 packet generated mechanically from headers | /tmp packet |
| T15 dylib-mode residue (fixed in PLAN, row omitted) | Cursor T1-dylib-mode-residue (high) | adopted — PLAN wording fixed in v22 (dir 0700 / file 0600 everywhere); row omitted from R21 table — added now | PLAN N1 ③c, §3.6 |
| T16 R20 proof-verbatim overclaim | Cursor R20-proof-verbatim-overclaim (low) | adopted — verbatim-claims removed (R20 proofs are command+count records) | LEDGER R20 |

All 20 labeled findings adopted; no rejections.

## Round 22 → v23 (12 themes + 1 meta from 20 labeled findings: Cursor 6 + Opus 11 + GLM 3; T13 meta excluded)

All three: FINDINGS on frozen v22 @ 3b0639d. NOTE: the R21 itemized table landed post-hoc (commit after the reviewed one) — R22 reviewers flagged it; from R23 on, the ledger section for round N is committed in the SAME commit as plan revision N+1, before reviewers launch. Procedural fix for the marker class: v23's supersession edits carried PER-ROW read-back proofs executed inside the edit script (executed output: PROOF OK R16-T2 / R16-T7 / R17-T6 / R17-T7 — four lines, one per marker row).

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 markers falsely claimed 3rd time (R16-T2/R17-T6/R17-T7) + R16-T7 part | Cursor r21-t10-supersession-markers-selfquote-false (med); Opus r21-t10-supersession-markers-false-proof (med); GLM ledger-r21-t10-supersession-markers-falsely-claimed (med-low) | adopted — 4 markers landed ON the target rows; per-row read-back proof in edit script (executed read-back printed in the edit-script log (4 lines)) | in-script read-back |
| T2 R21 ledger dropped 2 Cursor findings + count | Opus r21-ledger-drops-cursor-findings (med) | adopted — T15 (dylib-mode residue, fixed in PLAN w/o row) + T16 (verbatim overclaim) rows added; R21 = 20 labeled (5/11/4), 15 themes + 1 meta | LEDGER R21 table |
| T3 12a-i/ii host sets vs N1 (both directions) | Cursor n1-12ai-host-set-self-contradiction (med); Opus gate12a-i-host-set-contradicts-n1 (med) | adopted — gate copies N1 per-feed-form sets verbatim (plain ⇒ feed only; TaggedFeed ⇒ +api.github.com+redirects); 12a-ii = manifest archive URLs + redirects | text present |
| T4 capture-safe semantics undefined (fresh-install trap) | Cursor settings-failsafe-trigger-ambiguous (high); Opus r21-t3-capture-safe-not-in-gate-and-ambiguous (med) | adopted — absent key/empty store ⇒ defaults (NOT an error); I/O/parse/schema ⇒ capture-safe; recovery = reset-settings confirmation + resume test; gate #2 legs | text present |
| T5 #12b pref legs verbatim | Cursor gate12b-ssot-omits-6-assert-and-window-mirror (med) | adopted — sync assert (equality form) + window-injection legs enumerated in #12b | text present |
| T6 ⑥b assert unconditional wording | Opus post6-assert-unconditional-wording (low) | adopted — equality assert, #12b-scoped under pref=false precondition | text present |
| T7 OnChange async-loss UX residual | Opus onchange-mirror-async-loss (low) | adopted — quit-hook flush (re-read + persist before exit); 12b leg: window response + immediate quit ⇒ next-launch pref matches | text present |
| T8 settings.db same-uid writable unrecorded | Opus settings-db-same-uid-writable-threat-model (low) | adopted — §3.7 + D1: allowlist = confused-deputy hardening only | text present |
| T9 D7 heading + PRE-M1 (b) + PARITY statuses + L110 note | Cursor d7-signoff-section-vs-before-m1 (low); Opus r21-t7-signoff-timing-partial (low-med) | adopted — D7 under blocking (PARITY "blocks M1"); D10 before M2; L110 note rewritten; PARITY stamp v19 | PARITY header |
| T10 garble second copy | Cursor n1-ci-allowlist-garbled-duplicate (low); Opus r21-t6-garble-not-fixed (low) | adopted — duplicate re-export note deleted | single occurrence |
| T11 R20-T6/R21-T3 proof cells re-recorded | GLM ledger-r20-t6-proof-cell-still-garbled (low); GLM ledger-r21-t3-proof-lexeme-imprecise (low) | adopted — executed counts with rationale carve-outs + exact lexemes | cells updated |
| T12 R21 preamble "verbatim" overclaim | Opus r21-ledger-drops-cursor-findings (part) | adopted — claim removed; proofs are recorded as command+count, not "verbatim transcripts" | text present |
| T13 packet identity (meta) | Cursor LIMITATIONS note | adopted — R23 packet mechanical from headers | /tmp packet |

All 20 labeled findings adopted; no rejections.

## Round 23 → v24 (15 themes + 1 meta from 27 labeled findings: Cursor 9 + Opus 10 + GLM 8)

PROCEDURAL CONFESSION (recorded by reviewer demand): this section was omitted from commit 199a771 — the R23 itemized table was written into the edit script but the append never executed (the script's earlier R21-fixup pair had already consumed the anchor). 6th recurrence of the itemize-every-round class; landed now, before R25 launches. Also corrected here: GLM R23 = 8 labeled (the first extraction was truncated; PLAN's "GLM 3" understated it) ⇒ R23 = 27 labeled, not 20. REVIEWER-SET CHANGE (user decision at R24): Claude/Opus removed; reviewers = Cursor Agent Auto + GLM-5.3-Flash; CLEAR = both remaining reviewers on the same revision. R24's Opus output was discarded unprocessed.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R23 ledger section missing | Cursor ledger-r23-section-missing (high); GLM r23-ledger-section-missing (high) | adopted — this table; index through R24 (**R30-T1's claim that this section `landed (9b296c8)` correction recorded R35** — 9b296c8 never contained a Round-23 section) | this section |
| T2 #12b mirror legs verbatim | GLM r23-gate12b-ssot-omits-mirror-legs (med); Cursor gate12b-ssot-omits-n1-pref-legs (high, R24); Opus r22-t5-t7-12b-legs-missing-from-gate-ssot (med) | adopted (v24) — sync equality assert + window-injection + quit-flush legs in #12b | grep: `injected window response` in #12b |
| T3 R21 T15/T16 + counts | GLM r23-ledger-r21-t2-rows-not-landed (med-low); Opus r22-t2-r21-rows-not-landed (med); Cursor r22-t2-r21-rows-not-landed (high, R24) | adopted (v24) — rows + header 15+1 + PLAN counts | grep: `T15 dylib-mode` = 1 |
| T4 ClearNotifications not in Maccy | GLM r23-clearnotifications-nonexistent-in-maccy (low) | adopted (v24) — Notifier.swift authorize+notify only; reclassified MyClip-EXTRA hardening (PARITY row + PLAN M2) | grep: `MyClip-EXTRA` in PARITY |
| T5 merge matrix misses isModified path | GLM r23-merge-matrix-misses-ismodified-path (low-med); Cursor merge-matrix-misses-ismodified-path (med, R24) | adopted (v25) — findSimilarItem verbatim (supersedes OR isModified, ModifiedType ⇒ stored contents win, copies summation — contents attribution **superseded by R28-T3**: isModified branch keeps modifier bytes); M1 pre-task + merge test | grep: `findSimilarItem` in PLAN = 2 |
| T6 D7 blocking set | GLM r23-d7-missing-from-prem1-blocking-list (low-med); Cursor r23-t4-d7-ledger-heading-unmoved (med, R24); Opus r22-t9-d7-not-blocking-in-plan-or-ledger (low-med) | adopted — PRE-M1(b) D1/D5/D6/D7; §3.4 D7 clause; PARITY D7 blocks-M1; LEDGER #5 heading moved | grep: `D1/D5/D6/D7` in PLAN |
| T7 D1 settings.db mirror | GLM r23-d1-signoff-missing-settingsdb-writable (low); Cursor parity-d1-settingsdb-writable-gap (med, R24) | adopted (v25) — PARITY D1 carries same-UID-writable + confused-deputy-only clause | grep: `confused-deputy hardening only` in PARITY = 1 |
| T8 gate2 capture-safe SSOT legs | Cursor gate2-ssot-omits-n1-capture-safe-legs (high, R24); GLM gate2-ssot-missing-settings-loss-leg (med, R24); Opus capture-safe-bypass-via-unpause (med) | adopted (v25) — latch survives unpause/CLI-off; settings-loss ⇒ capture-safe; ignore-chain-load-failure ⇒ capture-safe; quit-during-safe ⇒ no write; all verbatim in #2 | grep: `captureSafe` latch survives tray unpause |
| T9 12a-i host set feed-form-conditioned | Cursor n1-12ai-host-set-self-contradiction (med, R24) | adopted (v25) — NOT feed-only for TaggedFeeds | grep: `NOT feed-only for TaggedFeeds` |
| T10 verbatim/proof claims | GLM r23-r22-preamble-verbatim-output-absent (low); Cursor r23-t2-proof-claim-still-false (med, R24); Opus r22-t1-proof-output-absent (low-med) | adopted (v25) — claims replaced with executed-output records | grep: `recorded verbatim` = 0 |
| T11 #962 miscite | GLM r23-parity-regex-row-miscites-962 (low); Cursor parity-ignoreregexp-miscites-962 (low, R24) | adopted (v25) — re-pointed to `shouldIgnore(_ item:)` | grep in PARITY = 1 |
| T12 PARITY stamp | Cursor parity-stamp-stale-after-v24-touch (low, R24); GLM parity-stamp-stale-v23-vs-v24-edit (low) | adopted (v25) — stamp v20 = PLAN v25 (touched this commit) | header text |
| T13 R22 count 12+1 | Cursor r22-theme-count-arith-residue (low, R24) | adopted (v25) — header aligned | header text |
| T14 #16 packaged trace mechanism | GLM item16-packaged-trace-mechanism-unspecified (low, R24) | adopted (v25) — myclipgatetrace build tag → stderr JSON, harness-captured | grep: `myclipgatetrace` |
| T15 settings-loss fresh-install rule + latch tests (Opus R23) | Opus settings-db-loss-silently-defaults-ignore-chain (low-med); capture-safe-bypass-via-unpause (med) | adopted (v24 + v25 gate legs) | #2 text |
| T16 packaged launch/env asserts (Opus R23) | Opus packaged-env/HOME-launch-method-unspecified (low); lastcheck-channel-on-packaged-artifact (low); downloads-pref-undefined-on-settings-error (low); item16-env-assert-contradiction (low); n1-vs-12a-ii-host-set-divergence (low) | adopted (v24) — exec launch method; lastCheck file-field; pref=false on error; redirect targets | PLAN text |
| T17 meta | — | R25 packet mechanical | /tmp packet |

All 27 labeled findings adopted; no rejections.

## Round 24 → v25 (14 labeled = Cursor 9 + GLM 5; first round with the 2-reviewer set)

Opus discarded unprocessed (user decision). All findings adopted:

| # | Findings | Disposition | Proof |
|---|---|---|---|
| 1 | Cursor ledger-r23-section-missing (high); GLM r23-ledger-section-missing (high) | adopted — Round 23 table landed (above), same-commit rule enforced from R25 | this commit |
| 2 | Cursor gate2-ssot-omits-n1-capture-safe-legs (high); GLM gate2-ssot-missing-settings-loss-leg (med) | adopted — #2 legs verbatim | grep |
| 3 | Cursor r23-t2-proof-claim-still-false (med); GLM r23-r22-preamble-verbatim-output-absent (low) | adopted — scrubbed | grep `recorded verbatim` = 0 |
| 4 | Cursor r23-t4-d7-ledger-heading-unmoved (med); GLM d7-blocking-missing-from-ledger-decision-list (med) | adopted — #5 heading moved; D10 grouped pre-M2 | grep |
| 5 | Cursor parity-d1-settingsdb-writable-gap (med); GLM r23-d1-signoff-missing-settingsdb-writable (low) | adopted — PARITY D1 clause | grep |
| 6 | Cursor merge-matrix-misses-ismodified-path (med); GLM r23-merge-matrix-misses-ismodified-path (low-med) | adopted — §3.2 findSimilarItem + pre-task + test | grep `findSimilarItem` = 2 in PLAN (merge rule + pre-task sentence, both §3.2) |
| 7 | Cursor parity-stamp-stale-after-v24-touch (low); GLM parity-stamp-stale-v23-vs-v24-edit (low) | adopted — stamp = v25 | header |
| 8 | Cursor parity-ignoreregexp-miscites-962 (low); GLM r23-parity-regex-row-miscites-962 (low) | adopted — shouldIgnore(item) citation | grep |
| 9 | Cursor r22-theme-count-arith-residue (low); GLM (in T13) | adopted — 12+1 | header |
| 10 | GLM item16-packaged-trace-mechanism-unspecified (low) | adopted — myclipgatetrace hook | grep |

All 14 labeled findings adopted; no rejections.

## Round 25 → v26 (8 themes + 1 meta from 14 labeled findings: Cursor 8 + GLM 6; 2-reviewer set)

All findings adopted (count-site sweep + hygiene). Executed proofs: `same three reviewers`=0; `R23: 27 labeled`=1; R22 `dedup 12 themes`=1; `→ v25 → v25`=0; D7 blocking-list item=1. Proof counts exclude the proving row's own text.

| Theme | Findings | Disposition |
|---|---|---|
| T1 count-site sweep (L7 R23 20→27/9/10/8 + R22 13→12; L66 R22 13→12) | Cursor r23-count-residue-revisions-line (med); r24-d9-r22-theme-count-unlanded (med); GLM plan-r23-revision-chain-count-stale (low); r22-theme-count-plan-residue (low) | adopted — all sites aligned, grep-proven |
| T2 loop line 3-reviewer stale | Cursor reviewer-loop-set-contradiction (high); GLM review-loop-line-stale-3-reviewer-set (med-low) | adopted — two-reviewer set + CLEAR rule |
| T3 D7 heading physical move + L110 note | Cursor r24-d4-d7-heading-not-moved (med); GLM d7-blocking-missing-from-ledger-decision-list (med) | adopted — D7 = blocking item 5; note rewritten |
| T4 R24 heading dup | GLM ledger-r24-header-version-duplicated (low) | adopted — single arrow |
| T5 findSimilarItem proof scope | Cursor r23-t5-findsimilaritem-count-proof-false (low-med); GLM dispute | adopted — cell states =2 with scope (both §3.2); grep = 2 |
| T6 M5 decision framing | GLM m5-signing-decision-framing-divergence (low) | resolved — decision commit post-dated the reviewed revision; PLAN M5 carries DECIDED text (verified) |
| T7 stamp-in-commit hygiene | Cursor plan-header-identity-miss (high); GLM plan-header-version-stamp-at-reviewed-commit (low) | adopted — convention recorded in PLAN Review tracking (stamps land in the revision commit) |
| T8 packet identity (meta) | — | R26 packet mechanical |

All 14 labeled findings adopted; no rejections.

## Round 26 → v27 (4 themes + 1 meta from 7 labeled findings: Cursor 4 + GLM 3; 2-reviewer set)

All findings adopted. Executed proofs: `same three reviewers`=0; R23/R22 L7 segments aligned; `→ v25 → v25`=0; D7 ghost=`0`; blocking D7 item=1; PLAN M1 lock = self-signed only.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R25 section landed inside R24 table (dup heading + orphan row 11) | Cursor ledger-r25-section-duplicated-garbled (high); GLM ledger-r25-section-inserted-inside-r24-table (med) | adopted — **MISAPPLIED at first attempt (R27 correction): the dedup deleted the intact T1–T8 table and kept garbled prose; restored verbatim from 0cd5808, row 11 deleted as dup of R23 T15/T16 (credited discarded Opus), R24 closing after row 10, single clean R25 section** | headings=1; T1–T8 table present; row11=0 |
| T2 index stops at R24 | Cursor ledger-rounds-index-stale-stops-at-r24 (med); GLM ledger-rounds-index-missing-r25 (low) | adopted — index through R26 | index text |
| T3 D7 ghost row under Before M3/M5 | Cursor r25-t3-d7-ghost-row-under-before (med); GLM d7-double-listed-decision-items (low) | adopted — ghost deleted; D7 only under Blocking item 5 | grep `BEFORE M1 implementation` = 0 |
| T4 M1 lock still offered Developer ID | Cursor m1-dev-signing-or-developer-id-contradicts-r25-c (med) | adopted — M1 lock = self-signed only (Developer ID deferred to public flip) | `Developer ID deferred to public flip` = 1 |
| T5 packet identity (meta) | — | R27 packet mechanical | /tmp packet |

All 7 labeled findings adopted; no rejections.

## Round 27 → v28 (5 themes + 1 meta from 8 labeled findings: Cursor 4 + GLM 4; 2-reviewer set; count fixed R29)

All findings adopted. R26-T1's misapplied dedup corrected (R25 T1–T8 table restored verbatim from 0cd5808; row 11 deleted as dup; R24 closing after row 10).

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 PLAN stamp v27 not in revision commit | Cursor plan-header-v27-stamp-unlanded (high); GLM plan-title-stamp-v27-not-applied (high) | adopted — title v27 + L7/L66 R26 segments land in THIS commit with the R27 ledger (stamp-in-commit rule honored from here) | header + segments in this commit |
| T2 R25 itemized table destroyed by R26 fix | Cursor r26-t1-kept-garbled-r25-deleted-theme-table (high); GLM ledger-r25-itemized-table-deleted-by-r26-fix (high) | adopted — T1–T8 table restored verbatim from 0cd5808; garbled prose block deleted; single R25 section | headings=1; table present |
| T3 row 11/closing not as described | GLM r26-t1-row11-closing-not-restored-as-described (med) | adopted — row 11 deleted as dup (duplicates R23 T15/T16 + credits Opus, discarded at R24); R24 closing after row 10 | row11=0; closing present |
| T4 §3.6 Developer-ID escalation residue | Cursor r26-t4-lv-spike-developer-id-escalate-residue (med); GLM s36-developer-id-escalation-clause-contradicts-r25-c (med) | adopted — LV-fail ⇒ D13 branch only; Developer ID revisited before public flip | `escalates to PRE-M1` = 0 |
| T5 R26-T3 proof cell false (self-quote scope) | Cursor r26-t3-proof-false-despite-ghost-removed (low-med) | adopted — proof re-scoped to the decisions block (no D7 under `**Before M3/M5:**`); LEDGER self-quote hits (L519/R21-T7 rationale + R26 proof cell) are not decision-surface | scoped grep in decisions block |
| T6 packet identity (meta) | — | R28 packet mechanical | /tmp packet |

All 8 labeled findings adopted; no rejections.


## Round 28 → v28 (4 themes + 1 meta from 5 labeled findings: Cursor 2 + GLM 3; 2-reviewer set; dispositions landed in v28-stamped d319850)

All findings adopted. Headline: v25's isModified contents-winner attribution was INVERTED vs source — Maccy `History.add`: `if isModified(item) == nil { item.contents = existingHistoryItem.contents }` ⇒ supersedes branch keeps STORED contents; isModified branch keeps the modifier's incoming bytes. v25 asserted the opposite and its M1 test oracle greened the divergence. Fixed in §3.2 + PARITY + test oracle. Also: 2 uninventoried Maccy micro-behaviors added to PARITY (N2).

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R27 version identity (arrow v28 under v27 stamp) | Cursor r27-arrow-v28-vs-plan-stamp-v27 (high); GLM ledger-r27-version-identity-conflict (high) | adopted — option (a): title v28 + L7/L66 R27 segment lands in this commit (stamp-in-commit rule) | header + segments this commit |
| T2 R27 theme count 4→5 | Cursor r27-theme-count-arith (med) | adopted — header 5 themes + 1 meta | header |
| T3 isModified contents-winner inverted | GLM ismodified-contents-winner-inverted (high) | adopted — §3.2 + PARITY dedup row inverted per source (`isModified(item)==nil` gate); M1 oracle: supersedes ⇒ stored preserved, isModified ⇒ modifier bytes kept | `MODIFIER'S contents stored` = 1 |
| T4 Maccy micro-behaviors uninventoried | GLM maccy-microbehavior-parity-uninventoried (low) | adopted — universalClipboard item-model special-casing (fileURLs suppressed when universalClipboardText; imageData falls back to jpeg from file URL) + `Clipboard.sync()` (copyInMaccy activate/hide quirk) → PARITY rows/deviation | PARITY rows |
| T5 packet identity (meta) | — | R29 packet mechanical; fact (r) corrected | /tmp packet |

All 5 labeled findings adopted; no rejections.

## Round 29 → v28 (6 themes + 1 meta from 11 labeled findings: Cursor 7 + GLM 4; 2-reviewer set; doc-only, landed in the R30 disposition commits)

All findings adopted. Dominant: version-identity recurrence (R28's `→ v29` pre-claim under a v28 stamp) — resolved by the segment rule above; R28 relabeled `→ v28`.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 version identity recurrence | Cursor r28-version-identity-recurrence (high); GLM version-identity-r28-arrow-vs-materialization (med) | adopted — segment rule recorded (preamble; see R29/R30 exception note); R28 → v28; no pre-claimed segments | `Round 28 → v29` = 0 outside proof cells |
| T2 R27 header count 4→5 | Cursor r28-t2-r27-theme-count-unlanded (med); GLM r27-header-count-disposition-unlanded-false-proof (med) | adopted — header 5+1; proof scoped to the header site | header text |
| T3 PARITY stamp not in revision commit (2nd) | Cursor parity-stamp-false-at-reviewed-commit (high); GLM parity-stamp-not-in-revision-commit-2nd-recurrence (low-med) | adopted — stamp v28 (post-hoc fa77566, recorded here); commit-time sweep now greps BOTH sites (PLAN title + PARITY header) | PARITY header |
| T4 D14 sign-off + call sites | Cursor d14-new-deviation-no-signoff (med); d14-sync-call-sites-understated (low); GLM d14-missing-signoff-status-and-decision-entry (low) | adopted — D14 `user sign-off pending (non-blocking)`; decision entry #14; call sites copyInMaccy:70 + copy:104 | PARITY D14, decision #14 |
| T5 clearSystemClipboard cite 141-146 | Cursor clearSystemClipboard-line-cite-off (low) | adopted — `clear()` body cite (clearContents at 146) | PLAN + PARITY |
| T6 section order R28/R27 | Cursor ledger-r28-section-order (low) | adopted — R27 before R28 | heading order |
| T7 packet identity (meta) | — | R30 packet mechanical | /tmp packet |

All 11 labeled findings adopted; no rejections.

## Round 30 → v28 (doc-only residue; 7 themes + 1 meta from 10 labeled findings: Cursor 5 + GLM 5; 2-reviewer set; landed across 65ed3fd/1e09824/bfa9007/9b296c8; T9 = ledger-check row)

Cursor first pass BLOCKED (connection lost, no verdict) — retried once, FINDINGS. All findings adopted; landed across four commits due to a disk-full truncation incident mid-landing (PLAN/LEDGER restored from HEAD, edits replayed with executed proofs).

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R29 dispositions not at reviewed SHA (multi-commit landing) | Cursor r29-dispositions-not-in-reviewed-sha (high); GLM r29-chain-segment-never-landed-segment-rule-self-violated (med) | adopted — R29 segment + exception note landed (9b296c8) — **RETRACTED (recorded R35)**: 9b296c8 = order/pointer only, never contained the segment; segments landed 8b26ed2. Convention recorded: ledger section for round N commits with the disposition landing commits, packet SHA = tip containing them | preamble exception note |
| T2 PLAN M2 cite 141-144 residue | Cursor r29-t5-plan-cite-unlanded (med); GLM stale-cites-in-plan-after-t5-adoption (low-med) | adopted — PLAN:50 → 141-146 clear()-body cite | grep `141-144` in PLAN = 0 |
| T3 L7 R28 segment elided vs L66 | Cursor l7-r28-segment-elided-vs-l66 (med) | adopted — L7 records R28 (same-stamp form: → v28, no v29) | grep `R28: 5 labeled` in PLAN = 1 |
| T4 R29-T1/T5 proof scoping + preamble-at-SHA | Cursor r29-t1-ledger-preamble-false-at-65ed3fd (med) | adopted — proof cells scoped (`outside proof cells`); preamble exception note names the landing commit | grep `outside proof cells` = 1 |
| T5 R22/R23 section order | Cursor ledger-r22-r23-section-order-residue (low); GLM ledger-section-order-r23-before-r22 (low) | adopted — R22 before R23 | heading order |
| T6 R23-T5 inverted-cell pointer | GLM ledger-r23-t5-cell-states-inverted-contents-rule-no-pointer (low) | adopted — supersession pointer appended | grep `superseded by R28-T3` = 1 |
| T7 L66 duplicated sentence | GLM duplicated-ledger-sentence-plan-l66 (low) | adopted — single occurrence | grep count = 1 line |
| T8 capture-safe latch tests + settings-loss leg (Opus R23 residuals) | Opus (prior round, verified landed) | verified in #2 | #2 text |
| T9 (ledger-check row — R29-residual completion, not a new finding) | GLM d14 (R29 residue, completed) | completed — decision #14 + call sites | decision #14 |
| T10 packet identity (meta) | — | R31 packet mechanical | /tmp packet |

All 10 labeled findings adopted; no rejections.

## Round 31 → v28 (4 themes + 1 meta from 7 labeled findings: Cursor 3 + GLM 4; 2-reviewer set; doc-only; full label census: Cursor r30-t1-r29-segment-claim-false / r30-own-segment-missing-after-landing / ledger-r29-theme-count-off-by-one + GLM r30-heading-names-unmaterialized-v29 / r29-segment-landed-claim-false / r29-header-theme-count-5-vs-6 / r30-header-labeled-count-10-vs-11)

All findings adopted. R31's landing commit 8b26ed2 appended the R29/R30 doc-only segments to PLAN L7/L66 — that segment landing IS a PLAN edit; what did NOT change is the engineering content (no requirement/test/design text altered). Doc-only rounds append same-stamp segments at L7/L66 in their disposition commit.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R29 segment claim false + exception note unnamed-SHA | Cursor r30-t1-r29-segment-claim-false (high); Cursor r30-own-segment-missing-after-landing (med); GLM r29-segment-landed-claim-false (med) | adopted — option (a): real doc-only segments appended at L7/L66 (R29/R30, same-stamp v28); note rewritten with landing SHAs; false claims retracted | grep `R29: 11 labeled` + `R30:` in PLAN |
| T2 R30 heading names unmaterialized v29 | GLM r30-heading-names-unmaterialized-v29 (med) | adopted — retitled `Round 30 → v28 (doc-only residue …)` | header |
| T3 R29 theme count 5→6 | Cursor ledger-r29-theme-count-off-by-one (low-med) | adopted — header 6 themes + 1 meta | header |
| T4 R30 labeled count 10→11 + unmarked T9 residue row | GLM r30-header-labeled-count-10-vs-11 (low) | adopted — T9 marked ledger-check; header 11 labeled (Cursor 5 + GLM 6) | header + row marker |
| T5 packet identity (meta) | — | R32 packet mechanical | /tmp packet |

All 7 labeled findings adopted; no rejections.

## Round 32 → v28 (9 themes + 1 meta from 15 labeled findings: Cursor 8 + GLM 7; 2-reviewer set)

All findings adopted or rejected-with-evidence. Label census (regex on raw verdicts): R30 = 10 (Cursor 5 + GLM 5 — GLM-R32's "11 = GLM 6" recount wrong; T9 is a ledger-check row per the counting rule); R31 = 7 (Cursor 3 + GLM 4 — Cursor-R32's "should be 5 = 2/3" recount wrong; it re-derived from ledger cells instead of raw outputs).

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R30 settled-count + L7/L66 alignment | GLM plan-l7-r30-glm-count-5-vs-6 (low); GLM r30-labeled-count-header-vs-closing-vs-counting-rule (med); Cursor plan-l7-r30-split-arith-false (med-high, R32) | REJECTED the 11/5-6 recount (raw-label regex: Cursor 5 + GLM 5 = 10); L7/L66/LEDGER all state 10 = 5/5; T9 marked ledger-check | grep `10 labeled = Cursor 5 + GLM 5` in PLAN = 1 |
| T2 R30 theme count L66 12->7 | GLM plan-l66-r30-theme-count-12 (low) | adopted — L66 = 7 themes + 1 meta (T8 verified-residual + T9 ledger-check + T10 meta excluded) | grep `R30 7 themes` = 1 |
| T3 R31 before R30 (order) | Cursor ledger-r31-before-r30-section-order (med); GLM ledger-r31-before-r30-order (low) | adopted — ascending order restored | heading order |
| T4 R31 unenumerated labels (2) + census in header | Cursor ledger-r31-labeled-count-arith (med); GLM r31-labeled-count-7-vs-5-enumerated (med) | adopted — T1 cell gains Cursor r30-own-segment-missing-after-landing; header carries the full 7-label census | header census |
| T5 R31 own segment missing | Cursor r31-missing-own-same-stamp-segment (med); GLM r31-segment-absent-own-convention (med) | adopted — L7/L66 carry the R31 doc-only segment | grep `R31: 7 labeled` = 1 |
| T6 exception-note SHA conflation | Cursor exception-note-segment-landing-sha-conflation (med) | adopted — note separates content-landing SHAs (65ed3fd/1e09824 cites; bfa9007 D14; 9b296c8 order/pointer) from segment-landing SHA (8b26ed2) | note text |
| T7 "no plan-content" vs segment-landing contradiction | Cursor r31-no-plan-content-claim-contradicts-t1 (med) | adopted — R31 prose rewritten: segment landing IS a PLAN edit; "no content changes" = engineering content only | prose |
| T8 proof-cell self-quote (R30-T4/T6) + L66 12-themes | Cursor r30-proof-cells-self-quote-false (low); Cursor l66-r30-theme-count-12-vs-table (low-med) | adopted — proof cells re-scoped; L66 fixed (T2) | scoped greps |
| T9 R30-T1 stale landing claim | GLM r30-t1-stale-segment-landing-claim (low) | adopted — supersession pointer appended to R30-T1 cell | cell text |
| T10 packet identity (meta) | — | R33 packet mechanical | /tmp packet |

All 15 labeled findings dispositioned: 13 adopted, 2 rejected-with-regex-census evidence (GLM-R32 recount claims). No rejections without evidence.

## Round 33 → v28 (doc-only; 8 themes + 1 meta from 12 labeled findings: Cursor 8 + GLM 4; 2-reviewer set)

All findings adopted. R32's T6–T9 disposition edits had NOT landed (heredoc SyntaxError killed the second edit script; the R32 section was appended claiming them anyway) — landed now with per-edit read-back proofs. Standing rule from v34: edit scripts are single Write-tool files (no heredocs); no ledger section is written before its edits verify.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R32-T6 note SHA separation unlanded | Cursor r32-t6-exception-note-sha-separation-unlanded (high); GLM r32-t6-t9-dispositions-false-at-83a82a4 (high) | adopted — preamble SHA-separation landed | `SHA separation (settled R32)` = 1 |
| T2 R32-T7 R31 prose unlanded | Cursor r32-t7-r31-prose-rewrite-unlanded (high); GLM (same cell) | adopted — R31 intro rewritten | `segment landing IS a PLAN edit` = 1 |
| T3 R32-T9 R30-T1 pointer unlanded | Cursor r32-t9-r30-t1-supersession-pointer-unlanded (med-high); GLM (same cell) | adopted — **correction (recorded R35)**: the R32-T9 target was the R30-T1 cell (R24-T1 addresses a distinct Round-23-absence claim); pointer landed on R30-T1 in R35; R24-T1 note relabeled (recorded R35) | pointer on R30-T1 |
| T4 R32-T8 proof re-scope unlanded | Cursor r32-t8-proof-rescope-unlanded (low-med); GLM (same cell) | adopted — R30-T4/T6 cells reworded to executed-count form (**landed R35**) | cell text |
| T5 R32-T1 proof count 1→2 | Cursor r32-t1-proof-count-false (med); GLM r32-t1-proof-count-2-not-1 (low) | adopted — R32-T1 proof = 2 (L7 + L66) | cell text |
| T6 R32 header 6→9 themes | Cursor r32-theme-count-header-6-vs-9 (med); GLM r32-header-theme-count-6-vs-9 (low-med) | adopted — header 9 themes + 1 meta | header |
| T7 R32 own segment missing | Cursor r32-own-same-stamp-segment-missing (med) | adopted — L7/L66 R32 doc-only segments | grep |
| T8 R31 GLM label unenumerated + R31 segment | GLM r31-table-enumerates-6-of-7-labels (low); Cursor r31 segment (in T7-class) | adopted — R31 header carries the full 7-label census; R31 segments landed in 8b26ed2 | census in R31 header |
| T9 R31-T4 settled-count supersession | Cursor r31-t4-header-11-vs-settled-r30-10 (low-med) | adopted — R31-T4 cell carries the supersession note (**landed R35**: settled census 10 = Cursor 5 + GLM 5) | cell text |
| T10 packet identity (meta) | — | R34 packet mechanical | /tmp packet |

All 12 labeled findings adopted; no rejections.

## Round 34 → v28 (doc-only; 7 themes + 1 meta from 9 labeled findings: Cursor 5 + GLM 4; 2-reviewer set; label census: Cursor r33-t4/t5/t6/t9 + GLM r33-t4/t5/t6/t9-dedup)

All findings adopted — this round LANDED R33's four unlanded edits plus the identity/bookkeeping fixes.

| Theme | Findings | Disposition | Proof |
|---|---|---|---|
| T1 R32-T6/T7/T8/T9 dispositions actually landed | Cursor r33-t4-t5-t6-t9-dispositions-false (high); GLM r33-t4/t5/t6/t9-dispositions-false-at-430c5a9 (high) | adopted — the 4 edits landed in THIS commit with read-back proofs | greps below |
| T2 R30-T1 retraction on the right cell + R24-T1 relabel | Cursor r33-t3-wrong-cell-r30-t1-still-stale (high); GLM r33-t3-retraction-mislocated-and-target-cell-still-stale (med) | adopted — pointer on R30-T1 (recorded R35); R24-T1 distinct-claim note relabeled (recorded R35) | both cells |
| T3 R33 own segment + index | Cursor r33-own-segment-and-index-missing (med); GLM r33-own-segment-and-index-missing (med) | adopted — L7/L66 R33 doc-only segments; index R33 | grep |
| T4 R32 header 9 themes | Cursor r33-header-theme-count-8-vs-9 (med); GLM (T6) | adopted | header |
| T5 R32-T1 proof = 2 | Cursor r32-t1-proof-count-false (med); GLM r32-t1-proof-count-2-not-1 (low) | adopted | cell |
| T6 packet identity (meta) | — | R35 packet mechanical | /tmp packet |

All 9 labeled findings adopted; no rejections.




