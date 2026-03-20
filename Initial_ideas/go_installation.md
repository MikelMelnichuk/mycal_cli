## Distribution Strategy for Linux: **Single Binary + Auto-Install**

Since **Linux-only** + **Go**, here's the optimal distribution flow:

## 🎯 **Method 1: One-Liner Install (Recommended)**

### Install Script (users run this):
```bash
# Save as install.sh or curl directly
curl -sSL https://github.com/yourusername/mycal/releases/latest/download/install.sh | bash
```

**`install.sh` content:**
```bash
#!/bin/bash
LATEST_URL="https://github.com/yourusername/mycal/releases/latest/download/mycal-linux-amd64"
INSTALL_DIR="/usr/local/bin"

echo "📥 Installing mycal..."
curl -L -o /tmp/mycal "$LATEST_URL"
sudo install /tmp/mycal "$INSTALL_DIR/mycal"
rm /tmp/mycal
echo "✅ mycal installed to $INSTALL_DIR/mycal"
echo "   Run: mycal --version"
```

### User experience:
```bash
curl -sSL https://get.mycal.sh | bash  # Done!
mycal today                           # Works immediately
```

## 🏗️ **Automated Releases with GoReleaser**

**`.goreleaser.yml`:**
```yaml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
    goarch:
      - amd64
      - arm64
    binary: mycal
    main: ./cmd/mycal

archives:
  - format: binary

release:
  github:
    owner: yourusername
    repo: mycal
```

**GitHub Actions (`.github/workflows/release.yml`):**
```yaml
name: Release
on:
  push:
    tags: ['v*']
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with: { go-version: 1.22 }
      - uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: '~> v1.21'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Tag to release:**
```bash
git tag v1.0.0
git push --tags  # → Auto-builds mycal-linux-amd64 binary
```

## 📦 **Method 2: Package Managers (Advanced)**

### **Homebrew** (Linux support good):
```bash
# User adds your tap
brew tap yourusername/mycal
brew install mycal

# Your tap setup (one-time):
mkdir mycal-homebrew && cd mycal-homebrew
git init
hub create yourusername/homebrew-mycal
# Add Formula/mycal.rb with download URL
```

### **Nix** (for Nix users):
```
nix-env -i mycal
```

## 🛠️ **Self-Updating Binary** (Pro feature)

Add to your Go code:
```go
func updateCmd() *cobra.Command {
    return &cobra.Command{
        Use: "update",
        RunE: func(cmd *cobra.Command, args []string) error {
            latest, current := checkVersion()
            if latest != current {
                downloadLatest()
                fmt.Println("✅ Updated to", latest)
            }
            return nil
        },
    }
}
```

## 📋 **Complete Distribution Flow**

```
1. Developer: git tag v1.2.3 && git push --tags
2. GitHub Actions: Builds mycal-linux-amd64 → GitHub Release
3. User: curl -sSL https://get.mycal.sh | bash
4. mycal: /usr/local/bin/mycal ✓ PATH ready
5. Future: mycal update  # Optional self-update
```

## 🎨 **User-Friendly Extras**

**Generate `get.mycal.sh`:**
```bash
#!/bin/bash
# Auto-detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) BINARY="mycal-linux-amd64" ;;
    aarch64|arm64) BINARY="mycal-linux-arm64" ;;
    *) echo "Unsupported arch"; exit 1 ;;
esac

LATEST_URL="https://github.com/yourusername/mycal/releases/latest/download/$BINARY"
curl -sSL "$LATEST_URL" | sudo tee /usr/local/bin/mycal > /dev/null
sudo chmod +x /usr/local/bin/mycal
mycal --version
```

## ✅ **Result**
- **Size:** ~8MB static binary
- **Install:** 1 command, 3 seconds  
- **PATH:** `/usr/local/bin/mycal` (standard)
- **No deps:** Runs everywhere (glibc only)
- **Updates:** Script + `mycal update`

**Most popular Linux CLI tools use this exact pattern** (e.g., `hugo`, `cheat`, `fzf`).

**Ready to scaffold the full GoReleaser + install script setup?** Just say "generate distribution files".