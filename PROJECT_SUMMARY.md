# Terraform Provider for OPNsense 26.1 - Project Summary

## What This Provider Does

This Terraform provider allows you to manage your OPNsense 26.1 firewall infrastructure as code. You can create, update, and delete firewall rules, aliases, DHCP configurations, and WireGuard VPN settings using Terraform.

## Quick Start

### 1. Build and Install
```bash
cd terraform-provider-opnsense
make build
make install
```

### 2. Configure Provider
```hcl
provider "opnsense" {
  host       = "https://192.168.1.1"
  api_key    = var.opnsense_api_key
  api_secret = var.opnsense_api_secret
  insecure   = true
}
```

### 3. Use Resources
```hcl
resource "opnsense_firewall_rule" "allow_ssh" {
  description      = "Allow SSH"
  interface        = "lan"
  protocol         = "tcp"
  source_net       = "192.168.1.0/24"
  destination_net  = "192.168.1.1"
  destination_port = "22"
  action           = "pass"
  enabled          = true
}
```

## Features Implemented

### ✅ Firewall Management
- **opnsense_firewall_rule** - Manage firewall filter rules
  - IPv4/IPv6 support
  - All protocols (TCP, UDP, ICMP, etc.)
  - Source/destination networks and ports
  - Actions: pass, block, reject
  - Logging and categories
  
- **opnsense_firewall_alias** - Manage firewall aliases
  - Host, network, port aliases
  - Multiple content entries
  - Easy reference in rules

### ✅ Kea DHCP Server
- **opnsense_kea_subnet** - Manage DHCP subnets
  - CIDR notation support
  - IP pool ranges
  - DHCP options
  
- **opnsense_kea_reservation** - Manage static reservations
  - MAC address to IP binding
  - Hostname assignment
  - Subnet association

### ✅ WireGuard VPN
- **opnsense_wireguard_server** - Manage WireGuard servers
  - Auto key generation
  - Configurable listen port
  - Tunnel addressing
  - Peer management
  
- **opnsense_wireguard_peer** - Manage WireGuard peers
  - Public key configuration
  - Allowed IPs
  - Endpoint configuration
  - Persistent keepalive

### ✅ Data Sources
- **opnsense_firewall_rule** - Query existing rules

## Files Created

### Core Provider Files
- `main.go` - Entry point
- `internal/provider/provider.go` - Provider configuration
- `internal/provider/resource_firewall_rule.go` - Firewall rule resource
- `internal/provider/resource_firewall_alias.go` - Firewall alias resource
- `internal/provider/resource_kea_subnet.go` - Kea subnet resource
- `internal/provider/resource_kea_reservation.go` - Kea reservation resource
- `internal/provider/resource_wireguard_server.go` - WireGuard server resource
- `internal/provider/resource_wireguard_peer.go` - WireGuard peer resource
- `internal/provider/data_source_firewall_rule.go` - Firewall rule data source

### Documentation
- `README.md` - Comprehensive documentation
- `QUICKSTART.md` - Quick start guide
- `IMPLEMENTATION_GUIDE.md` - Detailed implementation guide
- `CHANGELOG.md` - Version history

### Examples
- `examples/provider/provider.tf` - Provider configuration
- `examples/resources/opnsense_firewall_rule/resource.tf` - Rule examples
- `examples/resources/opnsense_firewall_alias/resource.tf` - Alias examples
- `examples/resources/opnsense_kea_dhcp/resource.tf` - DHCP examples
- `examples/resources/opnsense_wireguard/resource.tf` - WireGuard examples
- `examples/complete-infrastructure.tf` - Complete infrastructure example

### Build Files
- `go.mod` - Go module dependencies
- `Makefile` - Build automation
- `.gitignore` - Git ignore rules

## Key Design Decisions

### 1. Terraform Plugin Framework
Used the modern HashiCorp Terraform Plugin Framework (not the legacy SDK) for:
- Better type safety
- Improved performance
- Modern Go patterns
- Better testing support

### 2. Automatic Apply
All resources automatically call the appropriate apply/reconfigure endpoint after operations, ensuring changes are immediately active.

### 3. Environment Variables
Support for environment variables makes it easy to keep credentials out of code:
- `OPNSENSE_HOST`
- `OPNSENSE_API_KEY`
- `OPNSENSE_API_SECRET`

### 4. Comprehensive Examples
Included real-world examples covering:
- Basic configurations
- Complex multi-resource setups
- Complete infrastructure example with DMZ, management zones, VPN

## Next Steps to Get Started

1. **Get API Credentials**
   - Log into OPNsense web interface
   - Navigate to System → Access → Users
   - Generate API key/secret

2. **Build the Provider**
   ```bash
   cd terraform-provider-opnsense
   make build
   make install
   ```

3. **Test with Simple Example**
   ```bash
   cd examples/resources/opnsense_firewall_rule
   terraform init
   terraform plan
   terraform apply
   ```

4. **Verify in OPNsense**
   - Check Firewall → Automation → Filter
   - See your newly created rule

## Customization Options

### Adding More Resources
The provider is designed to be easily extensible. To add support for additional OPNsense features:

1. Create new resource file in `internal/provider/`
2. Implement the resource interface
3. Register in `provider.go`
4. Add examples and documentation

### Targeting Different Versions
While built for OPNsense 26.1, you can adapt for other versions by:
- Updating API endpoint URLs
- Adjusting field names to match API changes
- Testing with target version

### Custom Validation
Add validation logic in resource schemas to ensure:
- IP addresses are valid
- Ports are in range
- Required dependencies exist

## Important Notes

### API Coverage
This provider implements the most commonly used OPNsense features. OPNsense 26.1 has many more API endpoints that could be added in future versions.

### Testing
The provider includes:
- Schema validation
- Resource lifecycle management
- Error handling
- Import state support

For production use, you should:
- Test in a non-production environment first
- Maintain backups of OPNsense configuration
- Use version control for Terraform files
- Implement proper secret management

### Security
- API credentials are marked sensitive
- TLS verification is configurable
- Supports both self-signed and valid certificates
- All API communication uses HTTPS

## Common Use Cases

### 1. Dynamic Firewall Rules
Create rules based on infrastructure changes:
```hcl
resource "opnsense_firewall_rule" "allow_to_servers" {
  for_each = var.web_servers
  # ... rule configuration using each.value
}
```

### 2. DHCP with Infrastructure
Keep DHCP in sync with infrastructure:
```hcl
resource "opnsense_kea_reservation" "servers" {
  for_each   = var.servers
  ip_address = each.value.ip
  hw_address = each.value.mac
  hostname   = each.key
}
```

### 3. VPN for Remote Workers
Manage remote access at scale:
```hcl
resource "opnsense_wireguard_peer" "workers" {
  for_each    = var.remote_workers
  name        = each.key
  public_key  = each.value.pubkey
  allowed_ips = each.value.tunnel_ip
}
```

## Troubleshooting Resources

1. **Check API Access**
   ```bash
   curl -k -u "$KEY:$SECRET" https://opnsense/api/core/firmware/status
   ```

2. **Enable Debug Logging**
   ```bash
   export TF_LOG=DEBUG
   terraform apply
   ```

3. **Verify OPNsense Logs**
   - Check System → Log Files → General
   - Look for API-related errors

4. **Test Manually**
   - Try creating resources in web UI
   - Use browser developer tools to inspect API calls
   - Compare with provider implementation

## Support and Contribution

### Getting Help
- Review the comprehensive README.md
- Check QUICKSTART.md for setup issues
- Refer to IMPLEMENTATION_GUIDE.md for detailed information
- Review examples in the examples/ directory

### Contributing
To extend this provider:
1. Review IMPLEMENTATION_GUIDE.md
2. Follow existing code patterns
3. Add comprehensive examples
4. Update documentation
5. Test thoroughly

## Success Metrics

After setup, you should be able to:
- ✅ Create firewall rules via Terraform
- ✅ Manage DHCP configurations declaratively
- ✅ Set up WireGuard VPN through code
- ✅ Import existing resources
- ✅ Update configurations without manual intervention
- ✅ Version control your firewall configuration

## Final Notes

This provider gives you the foundation to manage OPNsense 26.1 as code. It implements the most critical features (firewall, DHCP, VPN) and is designed to be extended with additional resources as needed.

The modular design means you can add support for additional OPNsense features by following the patterns established in the existing resources.

**Happy automating your firewall! 🔥🛡️**
