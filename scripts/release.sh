#!/bin/bash
#
# release.sh - CCM Release Script
#
# Usage: ./release.sh <version>
# Example: ./release.sh 0.2.0
#
# This script:
# 1. Updates version in Makefile
# 2. Creates git commit
# 3. Creates git tag
# 4. Then you push to trigger GitHub Actions
#
# GitHub Actions will automatically:
# - Run tests and lint
# - Build all platforms (linux/darwin/windows x amd64/arm64)
# - Calculate checksums
# - Create GitHub Release
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}ℹ${NC} $1"; }
log_success() { echo -e "${GREEN}✓${NC} $1"; }
log_warn() { echo -e "${YELLOW}⚠${NC} $1"; }
log_error() { echo -e "${RED}✗${NC} $1"; }

# Check arguments
if [ -z "$1" ]; then
    log_error "Usage: $0 <version>"
    echo ""
    echo "Example: $0 0.2.0"
    echo ""
    echo "Version format: MAJOR.MINOR.PATCH"
    exit 1
fi

VERSION=$1

# Validate semver format
SEMVER_REGEX='^[0-9]+\.[0-9]+\.[0-9]+$'
if ! echo "$VERSION" | grep -qE "$SEMVER_REGEX"; then
    log_error "Invalid version format: $VERSION"
    echo "Expected: MAJOR.MINOR.PATCH (e.g., 0.2.0)"
    exit 1
fi

log_info "🚀 准备发布 CCM v$VERSION"
echo ""

# Step 1: Check git status
log_info "1. 检查 git 状态..."
if [ -n "$(git status --porcelain)" ]; then
    log_warn "有未提交的更改:"
    git status --short
    echo ""
    read -p "继续? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_error "已取消"
        exit 1
    fi
fi

# Step 2: Update version
log_info "2. 更新版本号..."
sed -i "s/^VERSION := .*/VERSION := $VERSION/" Makefile
log_success "版本已更新"

# Step 3: Create commit
log_info "3. 创建 commit..."
git add -A
git commit -m "Release v$VERSION" --quiet
log_success "Commit created: $(git rev-parse --short HEAD)"

# Step 4: Create tag
log_info "4. 创建 tag..."
git tag -a "v$VERSION" -m "Release v$VERSION" --quiet
log_success "Tag created: v$VERSION"

echo ""
echo "========================================"
log_success "Release v$VERSION 已准备就绪!"
echo "========================================"
echo ""
echo "📝 下一步:"
echo "   git push origin main && git push origin v$VERSION"
echo ""
echo "🤖 GitHub Actions 将自动:"
echo "   ✓ 运行测试和 lint"
echo "   ✓ 编译 6 个平台"
echo "   ✓ 计算 checksums"
echo "   ✓ 创建 GitHub Release"
echo ""
read -p "立即推送到远程? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git push origin main
    git push origin "v$VERSION"
    log_success "已推送! GitHub Actions 正在运行..."
    echo ""
    echo "🎉 Release 页面: https://github.com/taliove/go-claude-model/releases/tag/v$VERSION"
fi
