#!/bin/bash
# Verify the source code is correct before building

echo "========================================"
echo "Source Code Verification"
echo "========================================"
echo ""

ERRORS=0

# Check 1: go.mod version
echo "Check 1: go.mod version..."
if grep -q "terraform-plugin-framework v1.12.0" go.mod; then
    echo "✅ go.mod has correct framework version (v1.12.0)"
else
    echo "❌ go.mod has WRONG version!"
    echo "Expected: terraform-plugin-framework v1.12.0"
    grep "terraform-plugin-framework" go.mod
    ERRORS=$((ERRORS+1))
fi
echo ""

# Check 2: WireGuard file has fix
echo "Check 2: WireGuard server resource fix..."
if grep -q "readServerKeys" internal/provider/resource_wireguard_server.go; then
    echo "✅ WireGuard resource has been fixed (readServerKeys)"
else
    echo "❌ WireGuard resource NOT fixed!"
    echo "Line 183 should call: r.readServerKeys(ctx, &data)"
    echo "Currently shows:"
    sed -n '183p' internal/provider/resource_wireguard_server.go
    ERRORS=$((ERRORS+1))
fi
echo ""

# Check 3: Method signature
echo "Check 3: readServerKeys method signature..."
if grep -q "func (r \*WireguardServerResource) readServerKeys(ctx context.Context, data \*WireguardServerResourceModel)" internal/provider/resource_wireguard_server.go; then
    echo "✅ readServerKeys method signature correct"
else
    echo "❌ Method signature wrong or missing!"
    grep "func.*readServer" internal/provider/resource_wireguard_server.go || echo "Method not found!"
    ERRORS=$((ERRORS+1))
fi
echo ""

echo "======================================"
if [ $ERRORS -eq 0 ]; then
    echo "✅ ALL CHECKS PASSED"
    echo "======================================"
    echo ""
    echo "Source code is correct. Run:"
    echo "  ./clean-build.sh"
    exit 0
else
    echo "❌ $ERRORS CHECK(S) FAILED"
    echo "======================================"
    echo ""
    echo "You may have extracted an OLD version!"
    echo ""
    echo "Solution:"
    echo "  1. Delete this directory"
    echo "  2. Re-extract the LATEST .tar.gz file"
    echo "  3. Run ./verify.sh again"
    exit 1
fi
