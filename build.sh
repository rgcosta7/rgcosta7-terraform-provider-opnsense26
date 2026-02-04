#!/bin/bash
# Build script for terraform-provider-opnsense

set -e

echo "Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "Building terraform-provider-opnsense..."
go build -o terraform-provider-opnsense26

echo ""
echo "✅ Build complete!"
echo ""
echo "To install locally, run:"
echo "  ./install.sh"
echo ""
echo "Or install manually:"
echo "  mkdir -p ~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense26/0.1.0/linux_amd64/"
echo "  cp terraform-provider-opnsense26 ~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense26/0.1.0/linux_amd64/"
