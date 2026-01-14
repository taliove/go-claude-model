# CCM (Claude Code Manager) Makefile

# ===========================================
# Configuration
# ===========================================

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
BINDIR_INSTALL ?= $(HOME)/.local/bin
BINARY ?= bin/ccm

# ===========================================
# Build Targets
# ===========================================

.PHONY: all build dev run clean install uninstall install-global uninstall-global \
        test lint fmt tidy verify ci coverage deps release-local check help

all: build

# Local build (for development)
build:
	@mkdir -p bin
	go build -o $(BINARY) .

# Development build (with debug symbols, no optimization)
dev:
	@mkdir -p bin
	go build -gcflags="all=-N -l" -o $(BINARY) .

# Quick run (build and execute)
run: build
	./$(BINARY) $(ARGS)

clean:
	rm -rf bin
	rm -rf dist
	rm -f coverage.out coverage.html

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

# One-click verification (for local development)
verify: fmt lint test

# Simulate CI environment
ci: tidy verify

# Test coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser/v2@latest
	@echo "Done!"

# ===========================================
# Installation
# ===========================================

# Local installation (to ~/.local/bin)
install: build
	@echo "Installing CCM to $(BINDIR_INSTALL)..."
	install -d $(BINDIR_INSTALL)
	install -m 755 $(BINARY) $(BINDIR_INSTALL)/ccm
	@echo ""
	@echo "CCM installed successfully!"

uninstall:
	rm -f $(BINDIR_INSTALL)/ccm

# Global installation (to /usr/local/bin)
install-global: build
	@echo "Installing CCM globally to $(BINDIR)..."
	install -d $(BINDIR)
	install -m 755 $(BINARY) $(BINDIR)/ccm
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
	@echo "    dev            - Build with debug symbols (for debugging)"
	@echo "    run            - Build and run (use ARGS= for arguments)"
	@echo "    clean          - Remove build artifacts"
	@echo ""
	@echo "  Quality:"
	@echo "    test           - Run tests"
	@echo "    lint           - Run linter"
	@echo "    fmt            - Format code"
	@echo "    tidy           - Tidy go.mod"
	@echo "    verify         - Run all checks (fmt + lint + test)"
	@echo "    ci             - Simulate CI (tidy + verify)"
	@echo "    coverage       - Generate test coverage report"
	@echo ""
	@echo "  Setup:"
	@echo "    deps           - Install development dependencies"
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
