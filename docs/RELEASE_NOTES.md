MyClip is an early-preview, fully local clipboard manager for macOS and Linux.

- Search text and screenshots, pin clips, combine a paste stack, and paste
  into another application.
- Optional copy sound; macOS hold-to-cycle; English, French and Spanish.
- Local SQLite history, user-only file permissions, password-manager type
  exclusions, and fail-closed capture when settings cannot be read.

Choose the macOS universal DMG/ZIP or your Linux architecture's DEB/tar.gz.
Linux archives include an install.sh for a user-local installation. Install
GTK 3 and the optional paste/OCR/tray tools described in the README.

macOS builds are ad-hoc signed, not Developer ID signed or notarized.
Gatekeeper may require Privacy & Security › Open Anyway. Do not disable
Gatekeeper globally. Real Mac permissions, paste focus, Vision OCR and sound,
plus native Wayland desktop behavior, still need hands-on validation.
Linux ARM64 packages are cross-built, without a real ARM64 desktop test.

Verify SHA256SUMS and build provenance with
`gh attestation verify <download> --repo L92-Labs/myclip`.
The installers contain LICENSE, NOTICE.md and THIRD_PARTY_NOTICES.txt.
