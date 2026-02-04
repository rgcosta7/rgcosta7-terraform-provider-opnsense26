# 🧪 Testing Your OPNsense Terraform Provider

Congratulations on building the provider! Let's test it with a simple firewall alias.

## Step 1: Get Your OPNsense API Credentials

1. Log into your OPNsense web interface
2. Navigate to **System → Access → Users**
3. Select your user (or create one for API access)
4. Scroll to **API keys** section
5. Click **+ (plus)** to generate a new API key
6. Click **Download** to save the key file
7. Open the downloaded file - it contains:
   ```
   key=your-long-api-key-here
   secret=your-long-api-secret-here
   ```

## Step 2: Create Test Configuration

Create a new directory and file:

```bash
mkdir ~/opnsense-test
cd ~/opnsense-test
```

Create `main.tf`:

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
  host       = "https://192.168.1.1"  # ← Change to your OPNsense IP
  api_key    = "your-api-key"         # ← Paste your key
  api_secret = "your-api-secret"      # ← Paste your secret
  insecure   = true                   # Use true for self-signed certs
}

# Simple test: Create a firewall alias
resource "opnsense_firewall_alias" "test" {
  name        = "terraform_test"
  type        = "host"
  content     = ["8.8.8.8", "8.8.4.4"]
  description = "Test alias created by Terraform"
  enabled     = true
}

output "alias_id" {
  value = opnsense_firewall_alias.test.id
}
```

## Step 3: Better - Use Environment Variables (Recommended)

Instead of putting credentials in the file, use environment variables:

```bash
export OPNSENSE_HOST="https://192.168.1.1"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
```

Then your `main.tf` becomes simpler:

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
  insecure = true  # Only this is needed if using env vars
}

resource "opnsense_firewall_alias" "test" {
  name        = "terraform_test"
  type        = "host"
  content     = ["8.8.8.8", "8.8.4.4"]
  description = "Test alias created by Terraform"
  enabled     = true
}
```

## Step 4: Initialize Terraform

```bash
terraform init
```

**Expected output:**
```
Initializing provider plugins...
- Finding rgcosta7/opnsense versions matching "0.1.0"...
- Installing rgcosta7/opnsense v0.1.0...
- Installed rgcosta7/opnsense v0.1.0 (unauthenticated)

Terraform has been successfully initialized!
```

## Step 5: Plan the Changes

```bash
terraform plan
```

**You should see:**
```
Terraform will perform the following actions:

  # opnsense_firewall_alias.test will be created
  + resource "opnsense_firewall_alias" "test" {
      + content     = [
          + "8.8.8.8",
          + "8.8.4.4",
        ]
      + description = "Test alias created by Terraform"
      + enabled     = true
      + id          = (known after apply)
      + name        = "terraform_test"
      + type        = "host"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

## Step 6: Apply the Changes

```bash
terraform apply
```

Type `yes` when prompted.

**Expected output:**
```
opnsense_firewall_alias.test: Creating...
opnsense_firewall_alias.test: Creation complete after 2s [id=abc123-uuid-here]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

alias_id = "abc123-uuid-here"
```

## Step 7: Verify in OPNsense

1. Log into OPNsense web interface
2. Go to **Firewall → Aliases**
3. You should see **terraform_test** alias with IPs: 8.8.8.8, 8.8.4.4

## Step 8: Clean Up (Optional)

```bash
terraform destroy
```

This removes the test alias.

---

## 🎯 More Complex Test - Firewall Rule

After the alias works, try creating a rule:

```hcl
resource "opnsense_firewall_rule" "test_rule" {
  description      = "Test rule from Terraform"
  interface        = "lan"
  protocol         = "tcp"
  source_net       = "192.168.1.0/24"
  destination_net  = opnsense_firewall_alias.test.name  # Use the alias!
  destination_port = "53"
  action           = "pass"
  enabled          = true
  log              = true
}
```

---

## 🐛 Troubleshooting

### Error: "Provider not found"
```bash
terraform init -upgrade
rm -rf .terraform .terraform.lock.hcl
terraform init
```

### Error: "401 Unauthorized"
- Check API key and secret are correct
- Verify the API user has proper permissions
- Make sure API is enabled on OPNsense

### Error: "Connection refused"
- Check OPNsense host address
- Verify HTTPS is enabled on OPNsense
- Check firewall allows access from your machine

### Error: "Certificate verify failed"
- Set `insecure = true` in provider config (for self-signed certs)

### Changes not visible in OPNsense
- Check **Firewall → Automation → Filter** (not Rules)
- Or **Firewall → Aliases**
- Refresh the page

---

## ✅ Success Checklist

- [ ] `terraform init` succeeds
- [ ] `terraform plan` shows the resource to be created
- [ ] `terraform apply` completes without errors
- [ ] Alias appears in OPNsense web interface
- [ ] `terraform destroy` removes the alias

---

## 📚 Next Steps

Once the simple alias works, explore:

1. **More aliases**: Try network type, port type
2. **Firewall rules**: Create rules using your aliases
3. **DHCP**: Set up Kea DHCP subnets and reservations
4. **WireGuard**: Configure VPN servers and peers

Check the `examples/` directory for more complex configurations!

---

## 💡 Pro Tips

1. **Use variables** for reusable values:
   ```hcl
   variable "lan_network" {
     default = "192.168.1.0/24"
   }
   ```

2. **Use modules** for common patterns
3. **Keep state remote** (S3, Terraform Cloud) for team collaboration
4. **Always plan** before applying
5. **Test in staging** before production

---

**Need help?** Check `README.md` for full documentation!
