#!/bin/bash
# Build script for terraform-provider-opnsense

set -e

echo "Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "Building terraform-provider-opnsense..."
go build -o terraform-provider-opnsense

echo ""
echo "✅ Build complete!"
echo ""
echo "To install locally, run:"
echo "  ./install.sh"
echo ""
echo "Or install manually:"
echo "  mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yourusername/opnsense/0.1.0/linux_amd64/"
echo "  cp terraform-provider-opnsense ~/.terraform.d/plugins/registry.terraform.io/yourusername/opnsense/0.1.0/linux_amd64/"
