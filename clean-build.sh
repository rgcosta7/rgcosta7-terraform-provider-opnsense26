#!/bin/bash
# Complete clean build - forces fresh dependencies

set -e

echo "========================================"
echo "COMPLETE CLEAN BUILD"
echo "========================================"
echo ""
echo "This will:"
echo "  1. Remove ALL Go caches"
echo "  2. Remove go.sum"
echo "  3. Download FRESH dependencies"
echo "  4. Build provider"
echo ""

# Check Go version
GO_VERSION=$(go version 2>/dev/null || echo "not found")
if [[ "$GO_VERSION" == "not found" ]]; then
    echo "❌ Go is not installed!"
    exit 1
fi
echo "Go: $GO_VERSION"
echo ""

# Step 1: Clean everything
echo "Step 1/6: Cleaning Go build cache..."
go clean -cache 2>/dev/null || true
echo "✅ Build cache cleaned"
echo ""

echo "Step 2/6: Cleaning Go module cache..."
go clean -modcache 2>/dev/null || true
echo "✅ Module cache cleaned"
echo ""

echo "Step 3/6: Removing go.sum..."
rm -f go.sum
echo "✅ go.sum removed"
echo ""

# Verify go.mod has correct versions
echo "Step 4/6: Verifying go.mod versions..."
if grep -q "terraform-plugin-framework v1.12.0" go.mod; then
    echo "✅ Correct framework version (v1.12.0)"
else
    echo "⚠️  WARNING: go.mod may have wrong version!"
    echo "Expected: v1.12.0"
    echo "Found:"
    grep "terraform-plugin-framework" go.mod || echo "  (not found)"
fi
echo ""

echo "Step 5/6: Downloading fresh dependencies..."
# Force download with verification
go get github.com/hashicorp/terraform-plugin-framework@v1.12.0
go get github.com/hashicorp/terraform-plugin-go@v0.24.0
go mod download
go mod tidy
echo "✅ Dependencies downloaded"
echo ""

# Verify downloaded version
echo "Verifying downloaded versions..."
go list -m github.com/hashicorp/terraform-plugin-framework
go list -m github.com/hashicorp/terraform-plugin-go
echo ""

echo "Step 6/6: Building..."
go build -v -o terraform-provider-opnsense
BUILD_EXIT=$?

echo ""
if [ $BUILD_EXIT -eq 0 ] && [ -f "terraform-provider-opnsense" ]; then
    BINARY_SIZE=$(ls -lh terraform-provider-opnsense | awk '{print $5}')
    echo "======================================"
    echo "✅ BUILD SUCCESSFUL!"
    echo "======================================"
    echo "Binary: terraform-provider-opnsense ($BINARY_SIZE)"
    echo ""
    echo "Next step: ./install.sh"
    exit 0
else
    echo "======================================"
    echo "❌ BUILD FAILED"
    echo "======================================"
    echo ""
    echo "Debug info:"
    echo "  Go version: $GO_VERSION"
    echo "  Module path: $(pwd)"
    echo ""
    echo "Checking for version conflicts..."
    go mod graph | grep terraform-plugin-framework || true
    echo ""
    exit 1
fi
