# Releases

The CI workflow tests Linux/macOS and builds macOS universal (Apple Silicon
and Intel) plus Linux amd64/arm64. Linux ARM64 is cross-built on an amd64
runner. Non-tag runs upload installers and Linux smoke evidence as Actions
artifacts retained for 14 days.

## Publish a preview

1. Merge reviewed changes into main after CI is green. Update the version
   in mygo.json and the release notes in docs/RELEASE_NOTES.md in that PR.
2. Regenerate notices (`python3 scripts/licenses.py`) after dependency
   changes. CI verifies notices and installer contents.
3. Tag the main commit with exactly `v` plus the mygo.json version:

   ```sh
   git switch main
   git pull --ff-only
   git tag v0.1.0
   git push origin v0.1.0
   ```

4. The tag workflow verifies the version and main ancestry, repeats the
   full gates, gathers all six installers, creates SHA256SUMS and GitHub
   build-provenance attestations, and uploads them to a draft Release. It
   publishes the prerelease only after all artifacts have uploaded.
5. Download and check a release on a real desktop. Never label the preview
   stable based only on build or unit-test results.

If a tag build fails, fix the failure before retrying it. Do not overwrite a
published tag or installer. If publishing stopped with an incomplete draft,
remove that draft before rerunning the release job, or prepare a new version.

## Artifacts

- `myclip-VERSION-macos-universal.dmg` and `.zip`
- `myclip-VERSION-linux-amd64.deb` and `.tar.gz`
- `myclip-VERSION-linux-arm64.deb` and `.tar.gz`
- `SHA256SUMS`

All installers include the project license, Maccy/Lucide attribution and
dependency notices. macOS verification checks both architectures of the app
and SQLite library and verifies the app's signature. Linux verification
checks notices and SQLite in both the Debian and portable packages.

## Signing and permissions

Preview builds use MyGo's ad-hoc signature. No Developer ID certificate or
notarization credentials are configured. The release notes and README must
retain that limitation until signing/notarization is implemented and tested.

Normal jobs have contents:read only. The tag-only publisher gets contents,
OIDC and attestation write permissions. It never runs pull-request code in
a privileged pull_request_target workflow. Actions are pinned to commit
SHAs; Dependabot proposes updates. GitHub-hosted runners build releases;
no shared workstation credentials or private signing files are bundled.
