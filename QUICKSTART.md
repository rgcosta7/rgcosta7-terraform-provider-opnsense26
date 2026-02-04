# Quick Start Guide

This guide will help you get started with the OPNsense Terraform provider in under 10 minutes.

## Prerequisites

1. OPNsense 26.1 or later installed and accessible
2. Terraform 1.0 or later installed
3. Go 1.22 or later (if building from source)

## Step 1: Generate OPNsense API Credentials

1. Log into your OPNsense web interface
2. Go to **System → Access → Users**
3. Click on your user (or create a new one)
4. Scroll to the **API keys** section
5. Click the **+** button to generate a new key
6. Click **Download** to save the key (you can only download it once!)
7. The downloaded file contains your `key` and `secret`

## Step 2: Build and Install the Provider

```bash
# Clone the repository
git clone https://github.com/yourusername/terraform-provider-opnsense
cd terraform-provider-opnsense

# Build the provider
make build

# Install to local Terraform plugins directory
make install
```

## Step 3: Create Your First Configuration

Create a file named `main.tf`:

```hcl
terraform {
  required_providers {
    opnsense = {
      source = "yourusername/opnsense"
      version = "0.1.0"
    }
  }
}

provider "opnsense" {
  host       = "https://192.168.1.1"  # Your OPNsense URL
  api_key    = "your-api-key-here"
  api_secret = "your-api-secret-here"
  insecure   = true  # Use false with valid certificates
}

# Create a simple firewall rule
resource "opnsense_firewall_rule" "allow_ssh" {
  description      = "Allow SSH from LAN"
  interface        = "lan"
  protocol         = "tcp"
  source_net       = "192.168.1.0/24"
  destination_net  = "192.168.1.1"
  destination_port = "22"
  action           = "pass"
  enabled          = true
}
```

## Step 4: Initialize and Apply

```bash
# Initialize Terraform
terraform init

# See what will be created
terraform plan

# Apply the configuration
terraform apply
```

Type `yes` when prompted to create the resources.

## Step 5: Verify

1. Log into your OPNsense web interface
2. Go to **Firewall → Automation → Filter** (or **Firewall → Rules [new]**)
3. You should see your newly created rule!

## Using Environment Variables (Recommended)

Instead of hardcoding credentials in your `.tf` files, use environment variables:

```bash
export OPNSENSE_HOST="https://192.168.1.1"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
```

Then your provider configuration becomes:

```hcl
provider "opnsense" {
  insecure = true  # Only if using self-signed certs
}
```

## Common Use Cases

### 1. Create Firewall Alias and Rule

```hcl
# Define internal servers
resource "opnsense_firewall_alias" "web_servers" {
  name        = "web_servers"
  type        = "host"
  content     = ["192.168.1.10", "192.168.1.11"]
  description = "Web server pool"
}

# Allow HTTPS to web servers
resource "opnsense_firewall_rule" "allow_https" {
  description      = "Allow HTTPS to web servers"
  interface        = "wan"
  protocol         = "tcp"
  source_net       = "any"
  destination_net  = opnsense_firewall_alias.web_servers.name
  destination_port = "443"
  action           = "pass"
  enabled          = true
  log              = true
}
```

### 2. Setup DHCP with Static Reservations

```hcl
# Create DHCP subnet
resource "opnsense_kea_subnet" "office_network" {
  subnet      = "192.168.10.0/24"
  pools       = "192.168.10.100-192.168.10.200"
  description = "Office network DHCP"
}

# Reserve IP for printer
resource "opnsense_kea_reservation" "office_printer" {
  subnet      = opnsense_kea_subnet.office_network.id
  ip_address  = "192.168.10.50"
  hw_address  = "AA:BB:CC:DD:EE:FF"
  hostname    = "office-printer"
  description = "Office network printer"
}

# Reserve IP for NAS
resource "opnsense_kea_reservation" "nas" {
  subnet      = opnsense_kea_subnet.office_network.id
  ip_address  = "192.168.10.51"
  hw_address  = "11:22:33:44:55:66"
  hostname    = "nas"
  description = "Network storage"
}
```

### 3. Setup WireGuard VPN

```hcl
# Create WireGuard server
resource "opnsense_wireguard_server" "office_vpn" {
  name           = "wg0"
  enabled        = true
  listen_port    = 51820
  tunnel_address = "10.20.30.1/24"
}

# Add remote worker peer
resource "opnsense_wireguard_peer" "worker_laptop" {
  name        = "worker-laptop"
  enabled     = true
  public_key  = "laptop-public-key-here"
  allowed_ips = "10.20.30.10/32"
  keepalive   = 25
}

# Add mobile device peer
resource "opnsense_wireguard_peer" "worker_phone" {
  name        = "worker-phone"
  enabled     = true
  public_key  = "phone-public-key-here"
  allowed_ips = "10.20.30.11/32"
  keepalive   = 25
}
```

## Troubleshooting

### "Connection refused" or "Connection timeout"

- Verify OPNsense is accessible at the specified URL
- Check if HTTPS is enabled on OPNsense web interface
- Ensure no firewall rules are blocking access

### "401 Unauthorized"

- Verify API key and secret are correct
- Check that the API user has proper permissions
- Ensure API access is enabled in OPNsense

### "SSL certificate verify failed"

- If using self-signed certificates, set `insecure = true` in provider config
- For production, use valid certificates and set `insecure = false`

### Changes not appearing in OPNsense

- Check OPNsense system logs for errors
- Verify the API endpoints are accessible
- Try manually applying changes in the web interface to confirm permissions

## Next Steps

1. Review the [full README](README.md) for all available resources
2. Check the [examples directory](examples/) for more complex configurations
3. Read about [OPNsense API](https://docs.opnsense.org/development/api.html) to understand the underlying API

## Getting Help

- Check the [README](README.md) for detailed documentation
- Review [OPNsense API documentation](https://docs.opnsense.org/development/api.html)
- Open an issue on GitHub for bugs or feature requests

## Best Practices

1. **Use version control**: Keep your Terraform configurations in Git
2. **Use environment variables**: Never commit API credentials to version control
3. **Test in staging**: Test changes in a non-production environment first
4. **Use modules**: Create reusable modules for common patterns
5. **Document your rules**: Use descriptive names and descriptions for resources
6. **Regular backups**: Always maintain OPNsense backups independently of Terraform
7. **State management**: Use remote state storage (S3, Terraform Cloud, etc.) for team collaboration

Happy automating! 🚀
