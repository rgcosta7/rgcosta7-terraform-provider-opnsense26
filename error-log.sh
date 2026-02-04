#!/bin/bash
# Show FULL build error

echo "Building with full output..."
echo ""
go build -o terraform-provider-opnsense 2>&1 | tee build-full.log
echo ""
echo "Exit code: $?"
echo ""
echo "Full output saved to: build-full.log"
echo ""
echo "Last 30 lines:"
tail -30 build-full.log
