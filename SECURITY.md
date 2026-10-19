# Security policy

Report vulnerabilities with [GitHub private vulnerability reporting](https://github.com/L92-Labs/myclip/security/advisories/new).
Include the affected version/commit, platform, reproduction steps using
synthetic data, and expected versus actual behavior. Do not open a public
issue containing sensitive clipboard content or an unpatched exploit.

MyClip is an early preview. Security fixes target the latest version on main
and the latest published preview; older previews are not maintained.

## Security model

- Clipboard data is stored locally in a plaintext SQLite database, protected
  with user-only filesystem permissions. Other processes running as the same
  user can read it; this is not a password vault.
- Concealed/transient pasteboard types are excluded before reading their
  contents. Protection relies on the source application marking secrets.
- Unknown application identities are rejected in allowlist mode. Native
  Wayland source attribution is limited; see the README.
- SQLite secure_delete overwrites deleted database content. It does not
  securely erase SSD blocks, snapshots, backups or data held by other apps.
- Packaged builds include SQLite and do not fetch runtime libraries.
  Development/test tooling may download dependencies. There is no automatic
  updater; install a newer release manually.
- macOS preview builds are ad-hoc signed, without Developer ID notarization.
  GitHub artifact attestations identify the build workflow; they do not
  replace Apple code signing or prove that application behavior is safe.

Real Mac permissions/paste/OCR/audio and native Wayland desktop verification
remain separate from CI build and unit-test coverage.
