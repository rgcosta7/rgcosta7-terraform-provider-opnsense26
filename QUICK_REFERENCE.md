# Quick Reference - Build Commands

## ⚡ Fastest Way (Recommended)

```bash
./setup.sh
```

This single script will:
- Check Go installation
- Download dependencies
- Build the provider
- Offer to install it
- Show next steps

---

## 📋 Step-by-Step Manual Build

### 1. Download Dependencies
```bash
go mod download
go mod tidy
```

### 2. Build
```bash
go build -o terraform-provider-opnsense
```

### 3. Install
```bash
# Detect your OS/arch automatically
./install.sh

# Or install manually for Linux AMD64:
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense/0.1.0/linux_amd64/
cp terraform-provider-opnsense ~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense/0.1.0/linux_amd64/
```

---

## 🔨 Using Make

```bash
make build    # Downloads deps, builds binary
make install  # Builds and installs
make clean    # Removes binary
```

---

## 🧪 First Test

Create `test.tf`:
```hcl
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
  api_key    = "YOUR_API_KEY"
  api_secret = "YOUR_API_SECRET"
  insecure   = true
}

resource "opnsense_firewall_rule" "test" {
  description      = "Test rule"
  interface        = "lan"
  protocol         = "tcp"
  source_net       = "192.168.1.0/24"
  destination_net  = "any"
  destination_port = "443"
  action           = "pass"
}
```

Then:
```bash
terraform init
terraform plan
```

---

## ❓ Common Issues

### Issue: `missing go.sum entry` or `CloseEphemeralResource` error
**Solution:** Version incompatibility - run the fix script:
```bash
./fix-versions.sh
```

This will:
- Clean the Go module cache
- Remove old go.sum
- Download fresh dependencies
- Build with correct versions

### Issue: `go: command not found`
**Solution:** Install Go
```bash
# Ubuntu/Debian
sudo apt-get install golang-go

# CentOS/RHEL
sudo yum install golang
```

### Issue: Provider not found after install
**Solution:**
```bash
# Re-run install
./install.sh

# Clear Terraform cache
rm -rf .terraform .terraform.lock.hcl

# Reinitialize
terraform init -upgrade
```

### Issue: Wrong architecture
**Solution:** Check your architecture:
```bash
uname -m

# If x86_64: use linux_amd64
# If aarch64: use linux_arm64
# If arm64: use darwin_arm64 (Mac M1/M2)
```

---

## 📂 Important Files

- `setup.sh` - All-in-one setup script ⭐
- `build.sh` - Just build
- `install.sh` - Just install
- `Makefile` - Make commands
- `go.mod` - Dependencies list

---

## 🎯 Provider Configuration

Use environment variables (recommended):
```bash
export OPNSENSE_HOST="https://192.168.1.1"
export OPNSENSE_API_KEY="your-key"
export OPNSENSE_API_SECRET="your-secret"
```

Then in Terraform:
```hcl
provider "opnsense" {
  insecure = true  # Only for self-signed certs
}
```

---

## 📚 Documentation

- `README.md` - Complete documentation
- `QUICKSTART.md` - 10-minute tutorial
- `BUILD.md` - Build troubleshooting
- `examples/` - Usage examples
- `IMPLEMENTATION_GUIDE.md` - Technical details

---

## ✅ Verification Checklist

- [ ] Go 1.22+ installed
- [ ] Dependencies downloaded (`go mod download`)
- [ ] Binary built (`terraform-provider-opnsense` exists)
- [ ] Provider installed (check `~/.terraform.d/plugins/`)
- [ ] `terraform init` succeeds
- [ ] OPNsense API accessible

---

**Need help?** Check `BUILD.md` for detailed troubleshooting!
