# CCM (Claude Code Manager) Makefile

# ===========================================
# Configuration
# ===========================================

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
BINDIR_INSTALL ?= $(HOME)/.claude-model/bin
CONFIGDIR ?= $(HOME)/.claude-model/configs

# Output directory for releases
DIST_DIR ?= dist

# Version info (can be overridden at build time)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -s -w -X ccm/internal/version.Version=$(VERSION) -X ccm/internal/version.Commit=$(COMMIT) -X ccm/internal/version.Date=$(DATE)

# ===========================================
# Build Targets
# ===========================================

.PHONY: all build clean install uninstall install-global uninstall-global install-completion release

all: build

# Local build (outputs to project root for development)
build:
	go build -ldflags="$(LDFLAGS)" -o ccm .

clean:
	rm -f ccm
	rm -rf $(DIST_DIR)

# ===========================================
# Installation
# ===========================================

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

# Release builds for all platforms (outputs to dist/)
release:
	@echo "🔨 Building releases to $(DIST_DIR)/..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/ccm-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/ccm-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/ccm-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/ccm-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/ccm-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/ccm-windows-arm64.exe .
	@echo ""
	@echo "✓ Built releases to $(DIST_DIR)/"

# Calculate checksums
checksum:
	@echo "📝 Calculating checksums..."
	cd $(DIST_DIR) && \
	sha256sum ccm-* > checksums.txt 2>/dev/null || shasum -a 256 ccm-* > checksums.txt
	@echo "✓ checksums.txt created in $(DIST_DIR)/"

# Clean release artifacts
release-clean:
	rm -rf $(DIST_DIR)

# ===========================================
# Development Helper
# ===========================================

# Default GitHub user (can be overridden)
GITHUB_USER ?= taliove

# Show version info
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"
