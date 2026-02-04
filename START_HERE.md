# 🚀 START HERE - OPNsense Terraform Provider

## ⚠️ IMPORTANT: Fresh Start Required

If you've tried building before and got errors, you MUST start fresh:

```bash
# 1. Delete any old extracted directory
cd ..
rm -rf terraform-provider-opnsense-26  # or whatever you named it

# 2. Extract fresh copy
tar -xzf terraform-provider-opnsense.tar.gz
cd terraform-provider-opnsense

# 3. Verify source code is correct
./verify.sh
```

The `verify.sh` script will check if you have the correct source files.

---

## ✅ Build Instructions (In Order)

### Step 1: Verify Source Code
```bash
./verify.sh
```

**Expected output:**
```
✅ ALL CHECKS PASSED
```

If you see ❌ errors, re-extract the archive!

### Step 2: Clean Build
```bash
./clean-build.sh
```

This script will:
- Clean all Go caches
- Remove old go.sum
- Download correct versions (v1.12.0)
- Build the provider

**Expected output:**
```
✅ BUILD SUCCESSFUL!
```

### Step 3: Install
```bash
./install.sh
```

This installs the provider to your Terraform plugins directory.

---

## 🔍 Troubleshooting

### Error: "CloseEphemeralResource" or version v1.10.0 mentioned

**Problem:** Old cached version still in use

**Solution:**
```bash
./clean-build.sh
```

This aggressively cleans ALL caches.

### Error: "readServer" or "ReadResponse" 

**Problem:** You have old source code

**Solution:**
```bash
cd ..
rm -rf terraform-provider-opnsense  # Delete old version
tar -xzf terraform-provider-opnsense.tar.gz  # Extract fresh
cd terraform-provider-opnsense
./verify.sh  # Should show all checks passed
./clean-build.sh
```

### Error: Still failing after clean build

**Debug info to share:**
```bash
# Run these and share output:
./verify.sh
go version
go list -m github.com/hashicorp/terraform-plugin-framework
cat go.mod | grep terraform-plugin
```

---

## 📋 Quick Commands Summary

```bash
# Complete fresh build process:
tar -xzf terraform-provider-opnsense.tar.gz
cd terraform-provider-opnsense
./verify.sh        # Check source is correct
./clean-build.sh   # Build with clean cache
./install.sh       # Install provider
```

---

## 🎯 What Each Script Does

| Script | Purpose |
|--------|---------|
| `verify.sh` | Checks if source code is correct |
| `clean-build.sh` | Aggressive clean + build (USE THIS) |
| `setup.sh` | Simple build (try this first) |
| `build.sh` | Just build, no cache clean |
| `fix-versions.sh` | Version-specific fixes |
| `install.sh` | Install to Terraform plugins |

**Recommendation:** Use `clean-build.sh` - it's the most thorough.

---

## ✅ Success Checklist

After building, verify:

- [ ] `./verify.sh` shows all checks passed
- [ ] `terraform-provider-opnsense` binary exists
- [ ] Binary is ~40-50 MB in size
- [ ] `./install.sh` completes without errors
- [ ] `terraform init` finds the provider

---

## 🆘 Still Having Issues?

1. **Make sure you extracted the LATEST archive**
   ```bash
   ls -lh ../terraform-provider-opnsense.tar.gz
   # Should show recent date
   ```

2. **Check Go version** (need 1.22+)
   ```bash
   go version
   ```

3. **Run verification**
   ```bash
   ./verify.sh
   ```

4. **Share output from:**
   ```bash
   ./verify.sh
   ./clean-build.sh 2>&1 | tee build.log
   ```

---

## 📚 Next Steps After Install

See these files:
- `QUICKSTART.md` - Create your first resource
- `README.md` - Full documentation  
- `examples/` - Usage examples

---

**TL;DR:** Run `./verify.sh` then `./clean-build.sh` then `./install.sh`
