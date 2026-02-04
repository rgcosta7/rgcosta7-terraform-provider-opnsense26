#!/bin/bash
# Debug build - shows full error output

echo "========================================"
echo "DEBUG BUILD"
echo "========================================"
echo ""

echo "Go version:"
go version
echo ""

echo "Module info:"
go list -m
echo ""

echo "Framework version:"
go list -m github.com/hashicorp/terraform-plugin-framework
echo ""

echo "Attempting build..."
echo "========================================"
go build -v -o terraform-provider-opnsense
BUILD_STATUS=$?
echo "========================================"
echo ""

if [ $BUILD_STATUS -eq 0 ]; then
    echo "✅ Build succeeded!"
    ls -lh terraform-provider-opnsense
else
    echo "❌ Build failed with exit code: $BUILD_STATUS"
    echo ""
    echo "Last 50 lines of build output above ^"
fi

exit $BUILD_STATUS
