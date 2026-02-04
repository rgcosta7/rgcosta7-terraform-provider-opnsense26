#!/bin/bash
# Fix script for Go module version issues

set -e

echo "========================================"
echo "Module Version Fix Script"
echo "========================================"
echo ""

echo "This script will:"
echo "  1. Clean Go module cache"
echo "  2. Remove go.sum"
echo "  3. Download fresh dependencies"
echo "  4. Rebuild the provider"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo ""
echo "Step 1/5: Cleaning Go module cache..."
go clean -modcache
echo "✅ Cache cleaned"
echo ""

echo "Step 2/5: Removing go.sum..."
rm -f go.sum
echo "✅ go.sum removed"
echo ""

echo "Step 3/5: Downloading dependencies..."
go mod download
echo "✅ Dependencies downloaded"
echo ""

echo "Step 4/5: Tidying modules..."
go mod tidy
echo "✅ Modules tidied"
echo ""

echo "Step 5/5: Building..."
go build -o terraform-provider-opnsense
echo "✅ Build complete!"
echo ""

if [ -f "terraform-provider-opnsense" ]; then
    BINARY_SIZE=$(ls -lh terraform-provider-opnsense | awk '{print $5}')
    echo "✅ Binary created successfully: $BINARY_SIZE"
    echo ""
    echo "Run ./install.sh to install the provider"
else
    echo "❌ Build failed - binary not created"
    exit 1
fi
