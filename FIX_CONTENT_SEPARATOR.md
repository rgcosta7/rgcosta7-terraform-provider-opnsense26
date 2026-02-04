# Fix Applied: Content Separator

## Issue
OPNsense alias API expects **newline-separated** content, not comma-separated.

## What Was Wrong
```go
contentStr := strings.Join(contentItems, ",")  // ❌ Wrong - produces "8.8.8.8,8.8.4.4"
```

OPNsense rejected this with error:
```
"Entry \"8.8.8.8,8.8.4.4\" is not a valid hostname, IP address or range."
```

## What's Fixed
```go
contentStr := strings.Join(contentItems, "\n")  // ✅ Correct - produces "8.8.8.8\n8.8.4.4"
```

## Testing
To test with curl:
```bash
# Correct format (newline in content)
curl -k -u "$API_KEY:$API_SECRET" \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"alias":{"name":"test","type":"host","content":"8.8.8.8\n8.8.4.4","enabled":"1"}}' \
  "$HOST/api/firewall/alias/addItem"
```

## Next Steps
1. Rebuild the provider: `./clean-build.sh`
2. Reinstall: `./install.sh`  
3. Test again: `terraform apply`

This should now work correctly!
