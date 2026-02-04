#!/bin/bash
# Debug Terraform apply with full logging

echo "Running Terraform with debug logging..."
echo ""

# Set debug environment variable
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform-debug.log

echo "Running: terraform apply"
terraform apply

echo ""
echo "Debug log saved to: terraform-debug.log"
echo ""
echo "To see the API response, search for:"
echo "  grep 'API Response' terraform-debug.log"
echo ""
echo "Or view the last 50 lines:"
echo "  tail -50 terraform-debug.log"
