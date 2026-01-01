# CCM (Claude Code Manager) Makefile

# ===========================================
# Configuration
# ===========================================

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
BINDIR_INSTALL ?= $(HOME)/.local/bin

# ===========================================
# Build Targets
# ===========================================

.PHONY: all build clean install uninstall install-global uninstall-global \
        test lint fmt tidy release-local check help

all: build

# Local build (for development)
build:
	go build -o ccm .

clean:
	rm -f ccm
	rm -rf dist

# ===========================================
# Testing & Quality
# ===========================================

test:
	go test -v -race ./...

lint:
	golangci-lint run --timeout 5m

fmt:
	go fmt ./...

tidy:
	go mod tidy

# ===========================================
# Installation
# ===========================================

# Local installation (to ~/.local/bin)
install: build
	@echo "Installing CCM to $(BINDIR_INSTALL)..."
	install -d $(BINDIR_INSTALL)
	install -m 755 ccm $(BINDIR_INSTALL)/ccm
	@echo ""
	@echo "CCM installed successfully!"

uninstall:
	rm -f $(BINDIR_INSTALL)/ccm

# Global installation (to /usr/local/bin)
install-global: build
	@echo "Installing CCM globally to $(BINDIR)..."
	install -d $(BINDIR)
	install -m 755 ccm $(BINDIR)/ccm
	@echo "CCM installed to $(BINDIR)/ccm"

uninstall-global:
	rm -f $(BINDIR)/ccm

# ===========================================
# Release (GoReleaser)
# ===========================================

# Test GoReleaser locally (snapshot build)
release-local:
	goreleaser release --snapshot --clean

# Validate GoReleaser config
check:
	goreleaser check

# ===========================================
# Help
# ===========================================

help:
	@echo "CCM Makefile targets:"
	@echo ""
	@echo "  Build:"
	@echo "    build          - Build binary for local development"
	@echo "    clean          - Remove build artifacts"
	@echo ""
	@echo "  Quality:"
	@echo "    test           - Run tests"
	@echo "    lint           - Run linter"
	@echo "    fmt            - Format code"
	@echo "    tidy           - Tidy go.mod"
	@echo ""
	@echo "  Install:"
	@echo "    install        - Install to ~/.local/bin"
	@echo "    uninstall      - Uninstall from ~/.local/bin"
	@echo "    install-global - Install to /usr/local/bin"
	@echo ""
	@echo "  Release:"
	@echo "    release-local  - Test GoReleaser locally"
	@echo "    check          - Validate GoReleaser config"
	@echo ""
	@echo "  For actual releases, use: ./scripts/release.sh <version>"
