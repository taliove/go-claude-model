# CCM (Claude Code Manager) Makefile

# Configuration
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
BINDIR_INSTALL ?= $(HOME)/.claude-model/bin
CONFIGDIR ?= $(HOME)/.claude-model/configs

# Version info (can be overridden at build time)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -s -w -X ccm/internal/version.Version=$(VERSION) -X ccm/internal/version.Commit=$(COMMIT) -X ccm/internal/version.Date=$(DATE)

.PHONY: all build clean install uninstall install-global uninstall-global install-completion

all: build

build:
	go build -ldflags="$(LDFLAGS)" -o ccm .

clean:
	rm -f ccm

# Local installation (to ~/.claude-model/bin)
install:
	@echo "Installing CCM to $(BINDIR_INSTALL)..."
	install -d $(BINDIR_INSTALL)
	install -m 755 ccm $(BINDIR_INSTALL)/ccm
	@echo ""
	@echo "✓ CCM installed successfully!"
	@echo ""
	@if [ "$(BINDIR_INSTALL)" != "/usr/local/bin" ] && [ "$(shell which ccm 2>/dev/null)" = "" ]; then \
		echo "💡 To use CCM, add to PATH:"; \
		echo "   export PATH=\"$(BINDIR_INSTALL):\$$PATH\""; \
		echo ""; \
		echo "   Add to ~/.bashrc or ~/.zshrc:"; \
		echo "   echo 'export PATH=\"$(BINDIR_INSTALL):\$$PATH\"' >> ~/.bashrc"; \
	fi

# Local uninstall
uninstall:
	@echo "Removing CCM..."
	rm -f $(BINDIR_INSTALL)/ccm
	@echo "✓ CCM uninstalled"

# Global installation (to /usr/local/bin)
install-global:
	@echo "Installing CCM globally to $(BINDIR)..."
	install -d $(BINDIR)
	install -m 755 ccm $(BINDIR)/ccm
	@echo "✓ CCM installed to $(BINDIR)/ccm"

# Global uninstall
uninstall-global:
	@echo "Removing CCM from $(BINDIR)..."
	rm -f $(BINDIR)/ccm
	@echo "✓ CCM uninstalled"

# Shell completion installation
install-completion:
	@echo "Installing shell completion..."
	install -d $(HOME)/.zsh/completion
	install -m 644 scripts/completion.zsh $(HOME)/.zsh/completion/_ccm
	@echo "✓ Shell completion installed"
	@echo ""
	@echo "💡 Add to ~/.zshrc:"
	@echo "   autoload -U compinit"
	@echo "   compinit"

# Full installation (binary + completion)
install-all: install install-completion
	@echo ""
	@echo "✓ Full installation complete!"

# ===========================================
# Release & Distribution
# ===========================================

# Cross-compile for all platforms
release-all: clean build
	@echo ""
	@echo "🔨 Building for all platforms..."
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then \
				ext=".exe"; \
			else \
				ext=""; \
			fi; \
			echo "  Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -ldflags="$(LDFLAGS)" -o "ccm-$$os-$$arch$$ext" .; \
		done; \
	done
	@echo ""
	@echo "✓ Built: ccm-linux-amd64, ccm-linux-arm64, ccm-darwin-amd64, ccm-darwin-arm64, ccm-windows-amd64.exe, ccm-windows-arm64.exe"

# Calculate checksums
checksum:
	@echo "📝 Calculating checksums..."
	sha256sum ccm-* > checksums.txt 2>/dev/null || shasum -a 256 ccm-* > checksums.txt
	@echo "✓ checksums.txt created"

# Clean release artifacts
release-clean:
	rm -f ccm-*
	rm -f checksums.txt

# Publish helper - shows next steps
publish: release-all checksum
	@echo ""
	@echo "🎉 Release artifacts ready!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Go to: https://github.com/$(GITHUB_USER)/go-claude-model/releases/new"
	@echo "2. Upload: ccm-*, checksums.txt"
	@echo "3. Update formula/ccm.rb with new SHA256 (run: grep 'sha256' checksums.txt)"
	@echo ""
	@echo "💡 To create a tag and push:"
	@echo "   git tag v$(VERSION) && git push origin v$(VERSION)"

# Default GitHub user (can be overridden)
GITHUB_USER ?= taliove
