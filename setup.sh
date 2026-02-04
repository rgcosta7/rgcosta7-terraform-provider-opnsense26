#!/bin/bash
# Complete setup script for terraform-provider-opnsense

set -e

echo "========================================"
echo "Terraform Provider for OPNsense 26.1"
echo "Setup Script"
echo "========================================"
echo ""

# Check Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo ""
    echo "Please install Go 1.22 or later:"
    echo "  Ubuntu/Debian: sudo apt-get install golang-go"
    echo "  CentOS/RHEL:   sudo yum install golang"
    echo "  macOS:         brew install go"
    echo ""
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go $GO_VERSION detected"
echo ""

# Step 1: Download dependencies
echo "Step 1/3: Downloading Go dependencies..."
go mod download
echo "✅ Dependencies downloaded"
echo ""

# Step 2: Tidy modules
echo "Step 2/3: Tidying Go modules..."
go mod tidy
echo "✅ Modules tidied"
echo ""

# Step 3: Build
echo "Step 3/3: Building terraform-provider-opnsense..."
go build -o terraform-provider-opnsense
echo "✅ Build complete!"
echo ""

# Check binary was created
if [ ! -f "terraform-provider-opnsense" ]; then
    echo "❌ Error: Binary was not created"
    exit 1
fi

BINARY_SIZE=$(ls -lh terraform-provider-opnsense | awk '{print $5}')
echo "Binary size: $BINARY_SIZE"
echo ""

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
        echo "⚠️  Warning: Unsupported architecture: $ARCH"
        echo "Defaulting to amd64"
        ARCH="amd64"
        ;;
esac

PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense/0.1.0/${OS}_${ARCH}"

echo "========================================"
echo "Installation"
echo "========================================"
echo ""
echo "Plugin will be installed to:"
echo "  $PLUGIN_DIR"
echo ""
read -p "Install now? (y/n) " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Installing..."
    mkdir -p "$PLUGIN_DIR"
    cp terraform-provider-opnsense "$PLUGIN_DIR/"
    echo "✅ Installation complete!"
    echo ""
    echo "========================================"
    echo "Next Steps"
    echo "========================================"
    echo ""
    echo "1. Create a Terraform configuration file:"
    echo ""
    cat << 'EOF'
terraform {
  required_providers {
    opnsense = {
      source  = "rgcosta7/opnsense"
      version = "0.1.0"
    }
  }
}

provider "opnsense" {
  host       = "https://192.168.1.1"
  api_key    = var.opnsense_api_key
  api_secret = var.opnsense_api_secret
  insecure   = true
}
EOF
    echo ""
    echo "2. Initialize Terraform:"
    echo "   terraform init"
    echo ""
    echo "3. Check the examples/ directory for usage examples"
    echo ""
    echo "📚 Documentation:"
    echo "   - QUICKSTART.md    - Quick start guide"
    echo "   - README.md        - Full documentation"
    echo "   - examples/        - Usage examples"
    echo ""
else
    echo ""
    echo "Binary built but not installed."
    echo "To install manually, run:"
    echo "  mkdir -p $PLUGIN_DIR"
    echo "  cp terraform-provider-opnsense $PLUGIN_DIR/"
    echo ""
fi

echo "✅ Setup complete!"
