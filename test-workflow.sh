#!/bin/bash
# Test Release Workflow Locally
# This script simulates what GitHub Actions does

set -e

VERSION="0.1.0"
PLATFORMS=("linux_amd64" "linux_arm64" "darwin_amd64" "darwin_arm64" "windows_amd64")

echo "🧪 Testing Release Workflow"
echo "Version: $VERSION"
echo ""

# Clean up
rm -rf test-artifacts test-release
mkdir -p test-artifacts test-release

# Simulate building for each platform
echo "📦 Simulating multi-platform builds..."
for platform in "${PLATFORMS[@]}"; do
    echo "  Building for $platform..."
    mkdir -p "test-artifacts/terraform-provider-opnsense-$platform"
    
    if [[ $platform == windows* ]]; then
        echo "fake binary" > "test-artifacts/terraform-provider-opnsense-$platform/terraform-provider-opnsense_${VERSION}_${platform}.zip"
    else
        echo "fake binary" > "test-artifacts/terraform-provider-opnsense-$platform/terraform-provider-opnsense_${VERSION}_${platform}.tar.gz"
    fi
done

# Simulate downloading artifacts
echo ""
echo "📥 Collecting all artifacts..."
find test-artifacts -type f \( -name "*.tar.gz" -o -name "*.zip" \) -exec cp {} test-release/ \;

# Generate checksums
echo ""
echo "🔐 Generating checksums..."
cd test-release
sha256sum terraform-provider-opnsense_* > "terraform-provider-opnsense_${VERSION}_SHA256SUMS"

echo ""
echo "📋 Files before signing:"
ls -lah

# Check if GPG is available
if command -v gpg &> /dev/null; then
    echo ""
    echo "🔑 Testing GPG signing..."
    
    # Check if we have a key
    if gpg --list-secret-keys | grep -q "sec"; then
        SHASUMS_FILE="terraform-provider-opnsense_${VERSION}_SHA256SUMS"
        
        echo "Signing ${SHASUMS_FILE}..."
        gpg --detach-sign --armor "${SHASUMS_FILE}"
        
        if [ -f "${SHASUMS_FILE}.sig" ]; then
            echo "✅ Signature created successfully"
            
            # Verify signature
            if gpg --verify "${SHASUMS_FILE}.sig" "${SHASUMS_FILE}"; then
                echo "✅ Signature verified successfully"
            else
                echo "❌ Signature verification failed"
            fi
        else
            echo "❌ Signature file not created"
        fi
    else
        echo "⚠️  No GPG key found. Run ./setup-gpg.sh first"
    fi
else
    echo "⚠️  GPG not installed, skipping signature test"
fi

echo ""
echo "📋 Final release files:"
ls -lah

echo ""
echo "✅ Workflow test complete!"
echo ""
echo "Expected files:"
echo "  - terraform-provider-opnsense_${VERSION}_linux_amd64.tar.gz"
echo "  - terraform-provider-opnsense_${VERSION}_linux_arm64.tar.gz"
echo "  - terraform-provider-opnsense_${VERSION}_darwin_amd64.tar.gz"
echo "  - terraform-provider-opnsense_${VERSION}_darwin_arm64.tar.gz"
echo "  - terraform-provider-opnsense_${VERSION}_windows_amd64.zip"
echo "  - terraform-provider-opnsense_${VERSION}_SHA256SUMS"
echo "  - terraform-provider-opnsense_${VERSION}_SHA256SUMS.sig"
