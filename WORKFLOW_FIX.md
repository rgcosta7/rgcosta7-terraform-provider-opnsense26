# Workflow Fix Summary

## 🐛 Problem

The workflow failed with:
```
Error: Signature file not created
```

## 🔍 Root Cause

Each build job was creating its own `SHA256SUMS` file, but they were getting lost when artifacts were downloaded. The signing step couldn't find the checksums file.

## ✅ Solution

**Changed the workflow to:**

1. **Build job** - Each platform creates checksums locally (not uploaded)
2. **Release job** - Downloads all binaries, creates ONE combined checksum file
3. **Sign** - Signs the combined checksum file
4. **Publish** - Uploads everything

## 📋 Key Changes

### Before (❌ Broken)
```yaml
# Each job created its own SHA256SUMS
- name: Generate SHA256 checksums
  run: |
    cd dist
    sha256sum * > terraform-provider-opnsense_VERSION_SHA256SUMS

# This file got lost in artifact collection
```

### After (✅ Fixed)
```yaml
# Release job creates ONE combined SHA256SUMS
- name: Prepare release files
  run: |
    mkdir -p release
    # Copy all binaries
    find artifacts -type f \( -name "*.tar.gz" -o -name "*.zip" \) -exec cp {} release/ \;
    # Generate combined checksums
    cd release
    sha256sum terraform-provider-opnsense_* > terraform-provider-opnsense_VERSION_SHA256SUMS
```

## 🧪 Test Locally

Before pushing, test the logic:

```bash
./test-workflow.sh
```

This simulates what GitHub Actions does and verifies GPG signing works.

## 🚀 Next Steps

1. **Commit the fixed workflow:**
   ```bash
   git add .github/workflows/release.yml
   git commit -m "fix: Correct checksum generation and GPG signing"
   git push origin main
   ```

2. **Create a test release:**
   ```bash
   git commit --allow-empty -m "BUILD: v0.1.0 - Initial release"
   git push origin main
   ```

3. **Monitor GitHub Actions:**
   - Go to Actions tab
   - Watch the workflow run
   - Verify all steps complete successfully

4. **Check the release:**
   - Go to Releases page
   - Should have 7 files:
     - 5 platform binaries (3 .tar.gz, 1 .tar.gz, 1 .zip)
     - 1 SHA256SUMS file
     - 1 SHA256SUMS.sig file ✅

## 📝 Expected Release Files

```
terraform-provider-opnsense_0.1.0_linux_amd64.tar.gz
terraform-provider-opnsense_0.1.0_linux_arm64.tar.gz
terraform-provider-opnsense_0.1.0_darwin_amd64.tar.gz
terraform-provider-opnsense_0.1.0_darwin_arm64.tar.gz
terraform-provider-opnsense_0.1.0_windows_amd64.zip
terraform-provider-opnsense_0.1.0_SHA256SUMS         ← Combined checksums
terraform-provider-opnsense_0.1.0_SHA256SUMS.sig     ← GPG signature ✅
```

## ✅ Success Indicators

When the workflow succeeds, you'll see:

```
✅ Checksums signed successfully
-rw-r--r-- 1 runner docker  XXX terraform-provider-opnsense_0.1.0_SHA256SUMS
-rw-r--r-- 1 runner docker  XXX terraform-provider-opnsense_0.1.0_SHA256SUMS.sig
gpg: Signature made...
gpg: Good signature from "Your Name <email@example.com>"
```

## 🎯 Terraform Registry Requirements

With this fix, the release will have everything Terraform Registry needs:

- ✅ Binaries for multiple platforms
- ✅ SHA256 checksums file
- ✅ GPG signature file (.sig)
- ✅ Proper file naming convention

## 🆘 If It Still Fails

Check these in order:

1. **GPG secrets set correctly?**
   - Settings → Secrets → Actions
   - GPG_PRIVATE_KEY (with blank line after header!)
   - GPG_PASSPHRASE

2. **Key imported successfully?**
   - Check "Import GPG key" step in Actions log
   - Should show: `fingerprint=...`, `keyid=...`

3. **Files exist before signing?**
   - Check "Prepare release files" step
   - Should show all 5 platform files + SHA256SUMS

4. **Still having issues?**
   - Run `./test-workflow.sh` locally
   - Check the output for errors
