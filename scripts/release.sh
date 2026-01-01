#!/bin/bash
#
# release.sh - CCM Release Script with Pre-checks
#
# Usage: ./release.sh <version>
# Example: ./release.sh 0.3.0
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check arguments
if [ -z "$1" ]; then
    log_error "Usage: $0 <version>"
    echo ""
    echo "Example: $0 0.3.0"
    echo "         $0 0.3.0-beta.1"
    echo ""
    echo "Version format: MAJOR.MINOR.PATCH[-prerelease]"
    exit 1
fi

VERSION=$1

# Validate semver format (supports prerelease)
SEMVER_REGEX='^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$'
if ! echo "$VERSION" | grep -qE "$SEMVER_REGEX"; then
    log_error "Invalid version format: $VERSION"
    echo "Expected: MAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH-prerelease"
    exit 1
fi

log_info "Preparing release CCM v$VERSION"
echo ""

# Step 1: Check dependencies
log_info "Step 1/6: Checking dependencies..."
MISSING_DEPS=""

if ! command -v go &> /dev/null; then
    MISSING_DEPS="$MISSING_DEPS go"
fi

if ! command -v golangci-lint &> /dev/null; then
    log_warn "golangci-lint not found, skipping lint check"
    SKIP_LINT=true
fi

if ! command -v goreleaser &> /dev/null; then
    log_warn "goreleaser not found, skipping local config validation"
    SKIP_GORELEASER=true
fi

if [ -n "$MISSING_DEPS" ]; then
    log_error "Missing required dependencies:$MISSING_DEPS"
    exit 1
fi
log_success "Dependencies OK"

# Step 2: Check git status
log_info "Step 2/6: Checking git status..."
if [ -n "$(git status --porcelain)" ]; then
    log_error "Working directory is not clean:"
    git status --short
    echo ""
    read -p "Continue anyway? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_error "Aborted"
        exit 1
    fi
else
    log_success "Working directory clean"
fi

# Step 3: Check if tag exists
log_info "Step 3/6: Checking for existing tag..."
if git tag -l "v$VERSION" | grep -q "v$VERSION"; then
    log_error "Tag v$VERSION already exists"
    exit 1
fi
log_success "Tag v$VERSION is available"

# Step 4: Run tests
log_info "Step 4/6: Running tests..."
if ! go test ./...; then
    log_error "Tests failed"
    exit 1
fi
log_success "Tests passed"

# Step 5: Run linter
log_info "Step 5/6: Running linter..."
if [ "$SKIP_LINT" != "true" ]; then
    if ! golangci-lint run --timeout 5m; then
        log_error "Linter found issues"
        exit 1
    fi
    log_success "Linter passed"
else
    log_warn "Linter skipped"
fi

# Step 6: Validate GoReleaser config
log_info "Step 6/6: Validating GoReleaser config..."
if [ "$SKIP_GORELEASER" != "true" ]; then
    if ! goreleaser check; then
        log_error "GoReleaser config invalid"
        exit 1
    fi
    log_success "GoReleaser config valid"
else
    log_warn "GoReleaser validation skipped"
fi

echo ""
echo "========================================"
log_success "All pre-checks passed!"
echo "========================================"
echo ""

# Create tag
log_info "Creating tag v$VERSION..."
git tag -a "v$VERSION" -m "Release v$VERSION"
log_success "Tag created: v$VERSION"

echo ""
echo "Next steps:"
echo "   git push origin main && git push origin v$VERSION"
echo ""
echo "GitHub Actions will automatically:"
echo "   - Run tests and lint"
echo "   - Build for all 6 platforms"
echo "   - Create GitHub Release with changelog"
echo "   - Update Homebrew formula"
echo ""

read -p "Push to remote now? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git push origin main
    git push origin "v$VERSION"
    log_success "Pushed! GitHub Actions is running..."
    echo ""
    echo "Release page: https://github.com/taliove/go-claude-model/releases/tag/v$VERSION"
    echo "Actions: https://github.com/taliove/go-claude-model/actions"
fi
