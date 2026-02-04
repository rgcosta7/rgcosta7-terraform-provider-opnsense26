# Provider Naming - opnsense vs opnsense26

## Current Setup

The provider is currently registered as:
- **Module name**: `github.com/rgcosta7/terraform-provider-opnsense-26`
- **Registry name**: `rgcosta7/opnsense`
- **Version**: `0.1.0`

## To Use a Different Name (opnsense26)

If you want to use `opnsense26` instead of `opnsense` to differentiate from the v25 provider:

### Option 1: Change in Your Terraform Config Only

You can use `required_providers` with a local name:

```hcl
terraform {
  required_providers {
    opnsense26 = {
      source  = "rgcosta7/opnsense"
      version = "0.1.0"
    }
  }
}

provider "opnsense26" {
  # configuration
}

resource "opnsense26_firewall_alias" "test" {
  # ...
}
```

### Option 2: Change the Provider Registry Name

To truly register it as `opnsense26`, you need to:

1. **Change install path**:
```bash
# Instead of:
~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense/0.1.0/

# Use:
~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense26/0.1.0/
```

2. **Update install.sh**:
Edit `install.sh` and change:
```bash
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense26/0.1.0/${OS}_${ARCH}"
```

3. **Update main.go**:
Edit `main.go` and change:
```go
opts := providerserver.ServeOpts{
    Address: "registry.terraform.io/rgcosta7/opnsense26",
    Debug:   debug,
}
```

4. **Reinstall**:
```bash
./clean-build.sh
./install.sh
```

5. **Use in Terraform**:
```hcl
terraform {
  required_providers {
    opnsense26 = {
      source  = "rgcosta7/opnsense26"
      version = "0.1.0"
    }
  }
}

provider "opnsense26" {
  # ...
}
```

## Recommendation

**Option 1** is simpler - just use a different local name in your Terraform configs while keeping the provider registered as `rgcosta7/opnsense`. This way:
- Easier to maintain
- No confusion with file paths
- Can coexist with other providers

**Option 2** is better if you plan to:
- Publish to Terraform Registry
- Share with others who might have the v25 provider
- Want complete separation

## Current Status

Right now, your provider is:
- ✅ Built and installed as `rgcosta7/opnsense`
- ✅ Can be used alongside other providers
- ✅ Version 0.1.0 for OPNsense 26.1

If you just want to use it without conflicts, **Option 1** works perfectly!
