# Homebrew Tap for OKF Agent Memory

The official Homebrew Tap for **OKF Agent Memory** (`okf`) is maintained in the dedicated repository:
👉 **[uw-ssec/homebrew-tap](https://github.com/uw-ssec/homebrew-tap)**

The file in `Formula/okf.rb.tmpl` serves as the reference template for distribution packaging.

---

## 📦 How Users Install It

Users on macOS (Apple Silicon & Intel) and Linux can install `okf` with a single command:

```bash
brew install uw-ssec/tap/okf
```

Or by tapping the repository first:

```bash
brew tap uw-ssec/tap
brew install okf
```

---

## 🔄 Automatic Formula Synchronization

When a new version tag (e.g. `v0.2.0`) is published in this repository:
1. `.github/workflows/release.yml` compiles cross-platform binaries and generates `Formula/okf.rb` with accurate SHA256 checksums.
2. If `HOMEBREW_TAP_TOKEN` is configured in repository secrets, the workflow automatically commits and pushes the updated formula directly to `uw-ssec/homebrew-tap`. No manual editing is required!
