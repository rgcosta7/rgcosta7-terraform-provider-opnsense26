#!/bin/bash
# Install script for terraform-provider-opnsense

set -e

if [ ! -f "terraform-provider-opnsense" ]; then
    echo "Error: terraform-provider-opnsense binary not found."
    echo "Please run ./build.sh first"
    exit 1
fi

echo "Installing terraform-provider-opnsense..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense/0.1.0/${OS}_${ARCH}"

echo "Creating plugin directory: $PLUGIN_DIR"
mkdir -p "$PLUGIN_DIR"

echo "Copying binary..."
cp terraform-provider-opnsense "$PLUGIN_DIR/"

echo ""
echo "✅ Installation complete!"
echo ""
echo "Plugin installed to: $PLUGIN_DIR"
echo ""
echo "To use in your Terraform configuration:"
echo ""
echo "terraform {"
echo "  required_providers {"
echo "    opnsense = {"
echo "      source  = \"rgcosta7/opnsense\""
echo "      version = \"0.1.0\""
echo "    }"
echo "  }"
echo "}"
