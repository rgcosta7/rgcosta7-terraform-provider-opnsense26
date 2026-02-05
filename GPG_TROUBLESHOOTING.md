# GPG Key Troubleshooting

## ❌ Error: "Mandatory blank line missing between armor headers and armor data"

This happens when the GPG key is not properly formatted in GitHub Secrets.

## ✅ Quick Fix

### Step 1: Verify Your Key Format

```bash
# Check your exported key
cat terraform-gpg-private.asc | head -10
```

**Should look like this (with blank line):**
```
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQdGBGb1234ABCDEF...
mQINBGb1234...
(more base64 data)
```

**NOT like this (missing blank line):**
```
-----BEGIN PGP PRIVATE KEY BLOCK-----
lQdGBGb1234ABCDEF...   ← ERROR: No blank line!
```

### Step 2: Fix the Format

If the blank line is missing, fix it:

```bash
# Method 1: Re-export with proper formatting
gpg --armor --export-secret-keys YOUR_KEY_ID | sed '1 a\\' > private-key-fixed.asc

# Method 2: Add blank line manually
sed -i '1 a\\' terraform-gpg-private.asc
```

### Step 3: Verify the Fix

```bash
# This should show a blank line (line 2)
head -5 terraform-gpg-private.asc
```

Should output:
```
-----BEGIN PGP PRIVATE KEY BLOCK-----
                                      ← Line 2 is blank
lQdGBGb1234...                        ← Line 3 starts the data
```

### Step 4: Copy to GitHub Secret

```bash
# Copy the ENTIRE file content
cat terraform-gpg-private.asc
```

**Copy everything including:**
- `-----BEGIN PGP PRIVATE KEY BLOCK-----`
- The blank line
- All the base64 data
- `-----END PGP PRIVATE KEY BLOCK-----`

## 🔍 Alternative: Use the `pbcopy` Command (macOS)

```bash
# macOS - copies to clipboard
cat terraform-gpg-private.asc | pbcopy

# Linux with xclip
cat terraform-gpg-private.asc | xclip -selection clipboard

# Then paste directly into GitHub secret
```

## 🛠️ Generate New Key (if still having issues)

If the format is still wrong, generate a fresh key:

```bash
# 1. Delete old key (optional)
gpg --delete-secret-key YOUR_KEY_ID
gpg --delete-key YOUR_KEY_ID

# 2. Generate new key with batch mode (ensures proper format)
cat > gpg-batch.txt <<EOF
%no-protection
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: Your Name
Name-Email: your-email@example.com
Name-Comment: Terraform Provider Signing
Expire-Date: 0
EOF

gpg --batch --generate-key gpg-batch.txt

# 3. Get new key ID
gpg --list-secret-keys --keyid-format=long your-email@example.com

# 4. Export with proper formatting
KEY_ID="YOUR_NEW_KEY_ID"
gpg --armor --export-secret-keys $KEY_ID > new-private-key.asc

# 5. Verify format
head -5 new-private-key.asc
```

## 📋 Common Issues

### Issue 1: Extra Whitespace

**Problem:** Spaces or tabs on the blank line

**Fix:**
```bash
# Remove any whitespace from blank lines
sed -i 's/^[[:space:]]*$//' terraform-gpg-private.asc
```

### Issue 2: Windows Line Endings

**Problem:** File has `\r\n` (CRLF) instead of `\n` (LF)

**Fix:**
```bash
# Convert to Unix line endings
dos2unix terraform-gpg-private.asc

# Or with sed
sed -i 's/\r$//' terraform-gpg-private.asc
```

### Issue 3: Copy-Paste Corruption

**Problem:** Copy/paste adds or removes characters

**Fix:** Don't copy-paste from terminal. Use file copy instead:

```bash
# Use cat and copy directly
cat terraform-gpg-private.asc

# Or use clipboard tools
```

## ✅ Test the Key

Test if GitHub can import it:

```bash
# Save the key content to a test file
cat terraform-gpg-private.asc > test-import.asc

# Try importing it
gpg --import test-import.asc

# Should succeed without errors
```

## 🎯 Final Checklist

Before adding to GitHub Secret:

- [ ] Key starts with `-----BEGIN PGP PRIVATE KEY BLOCK-----`
- [ ] Second line is completely blank (no spaces)
- [ ] Third line starts with base64 data (e.g., `lQdGBGb1...`)
- [ ] Key ends with `-----END PGP PRIVATE KEY BLOCK-----`
- [ ] No extra whitespace or line endings
- [ ] File size is reasonable (typically 3-5 KB)

## 💡 Working Example

Here's what a **correct** key looks like when you view it:

```bash
$ cat terraform-gpg-private.asc
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQdGBGb1NXEBEADMqr7vXqK+hV...
mQINBGb1NXEBEACxVN8GXqwvP...
... (many more lines of base64)
=AbCd
-----END PGP PRIVATE KEY BLOCK-----
```

Notice:
1. Line 1: `-----BEGIN PGP PRIVATE KEY BLOCK-----`
2. Line 2: **BLANK**
3. Line 3+: Base64 data

## 🚨 If Still Not Working

Create an issue with:
- Output of: `head -5 terraform-gpg-private.asc | od -c`
- GPG version: `gpg --version`
- OS: Linux/macOS/Windows

This will show the exact bytes and help diagnose the issue.
