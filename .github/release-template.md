# CCM v{{ .Version }} Release

## 📥 Download

| Platform | Architecture | Download |
|----------|-------------|----------|
| macOS | Intel | [ccm-darwin-amd64]({{ .DownloadURL }}/ccm-darwin-amd64) |
| macOS | Apple Silicon | [ccm-darwin-arm64]({{ .DownloadURL }}/ccm-darwin-arm64) |
| Linux | x86_64 | [ccm-linux-amd64]({{ .DownloadURL }}/ccm-linux-amd64) |
| Linux | ARM64 | [ccm-linux-arm64]({{ .DownloadURL }}/ccm-linux-arm64) |
| Windows | x86_64 | [ccm-windows-amd64.exe]({{ .DownloadURL }}/ccm-windows-amd64.exe) |
| Windows | ARM64 | [ccm-windows-arm64.exe]({{ .DownloadURL }}/ccm-windows-arm64.exe) |

## 🔒 Verify

Check the integrity of your download:

```bash
# Download checksums
curl -L {{ .DownloadURL }}/ccm_{{ .Version }}_checksums.txt -o checksums.txt

# Verify (Linux/macOS)
sha256sum -c checksums.txt

# Verify (macOS with shasum)
shasum -a 256 -c checksums.txt
```

## 📦 Installation

### macOS / Linux

```bash
# Download and install
curl -L {{ .DownloadURL }}/ccm-{{ .Os }}-{{ .Arch }} -o ccm
chmod +x ccm
sudo mv ccm /usr/local/bin/ccm

# Or use Homebrew (if available)
brew install ccm
```

### Windows

```powershell
# PowerShell
Invoke-WebRequest -Uri {{ .DownloadURL }}/ccm-windows-amd64.exe -OutFile ccm.exe
.\ccm.exe --help
```

## 🔧 From Source

```bash
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model
git checkout v{{ .Version }}
make install
```

## 📝 Changelog

{{ .Changes }}

## 🔗 Links

- [Project Home](https://github.com/taliove/go-claude-model)
- [Documentation](https://github.com/taliove/go-claude-model#readme)
- [Report Issues](https://github.com/taliove/go-claude-model/issues)
