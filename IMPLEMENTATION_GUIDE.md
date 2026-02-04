# OPNsense 26.1 Terraform Provider - Implementation Guide

## Overview

This Terraform provider enables Infrastructure as Code management of OPNsense 26.1 firewall configurations. It provides resources for managing firewall rules, aliases, Kea DHCP, and WireGuard VPN.

## Project Structure

```
terraform-provider-opnsense/
├── main.go                          # Provider entry point
├── go.mod                           # Go module dependencies
├── internal/provider/
│   ├── provider.go                  # Provider configuration
│   ├── resource_firewall_rule.go    # Firewall rule resource
│   ├── resource_firewall_alias.go   # Firewall alias resource
│   ├── resource_kea_subnet.go       # Kea DHCP subnet resource
│   ├── resource_kea_reservation.go  # Kea DHCP reservation resource
│   ├── resource_wireguard_server.go # WireGuard server resource
│   ├── resource_wireguard_peer.go   # WireGuard peer resource
│   └── data_source_firewall_rule.go # Firewall rule data source
├── examples/
│   ├── provider/                    # Provider configuration examples
│   ├── resources/                   # Resource usage examples
│   └── complete-infrastructure.tf   # Complete infrastructure example
├── README.md                        # Comprehensive documentation
├── QUICKSTART.md                    # Quick start guide
├── CHANGELOG.md                     # Version history
├── Makefile                         # Build automation
└── .gitignore                       # Git ignore rules
```

## Key Components

### 1. Provider Configuration (provider.go)

The provider handles:
- Authentication with API key/secret
- TLS certificate verification control
- HTTP client configuration
- Resource and data source registration

**Key Features:**
- Environment variable support for credentials
- Configurable timeout
- Insecure mode for self-signed certificates
- Basic authentication for API requests

### 2. Resources

#### Firewall Rule Resource
- **Endpoint**: `/api/firewall/filter/*`
- **Features**: Create, read, update, delete filter rules
- **Auto-apply**: Yes (calls `/api/firewall/filter/apply`)
- **Supports**: IPv4/IPv6, all protocols, source/destination networks, ports, actions

#### Firewall Alias Resource
- **Endpoint**: `/api/firewall/alias/*`
- **Features**: Manage host, network, port, and other alias types
- **Auto-apply**: Yes (calls `/api/firewall/alias/reconfigure`)
- **Supports**: Multiple content entries, various alias types

#### Kea DHCP Resources
- **Endpoints**: `/api/kea/dhcpv4/*` and `/api/kea/service/*`
- **Features**: Subnet management, static reservations
- **Auto-apply**: Yes (calls `/api/kea/service/reconfigure`)
- **Supports**: IPv4 DHCP, pools, options

#### WireGuard Resources
- **Endpoints**: `/api/wireguard/server/*` and `/api/wireguard/client/*`
- **Features**: Server instances, peer management
- **Auto-apply**: Yes (calls `/api/wireguard/service/reconfigure`)
- **Supports**: Key generation, persistent keepalive, endpoint configuration

### 3. Data Sources

#### Firewall Rule Data Source
Allows querying existing firewall rules by UUID.

## API Integration

### Authentication
The provider uses HTTP Basic Authentication with API key as username and API secret as password:

```go
req.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
```

### Request Flow
1. Create/Update/Delete operation builds JSON payload
2. HTTP request sent to appropriate endpoint
3. Response parsed for UUID or error
4. Configuration applied via reconfigure/apply endpoint

### Error Handling
- HTTP status codes checked
- JSON response parsing
- Diagnostics added for user feedback

## Building and Testing

### Build the Provider

```bash
# Install dependencies
go mod download

# Build
make build

# Install locally
make install
```

### Run Tests

```bash
# Unit tests
make test

# Acceptance tests (requires running OPNsense)
make testacc
```

### Format Code

```bash
# Format Go code
make fmt

# Run linter
make lint
```

## Customization Guide

### Adding a New Resource

1. **Create resource file**: `internal/provider/resource_<name>.go`

2. **Implement resource interface**:
```go
type MyResource struct {
    client *Client
}

type MyResourceModel struct {
    ID   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    // ... other fields
}

func (r *MyResource) Metadata(...)
func (r *MyResource) Schema(...)
func (r *MyResource) Create(...)
func (r *MyResource) Read(...)
func (r *MyResource) Update(...)
func (r *MyResource) Delete(...)
func (r *MyResource) ImportState(...)
```

3. **Register in provider.go**:
```go
func (p *opnsenseProvider) Resources(...) []func() resource.Resource {
    return []func() resource.Resource{
        // ... existing resources
        NewMyResource,
    }
}
```

4. **Create examples**: `examples/resources/opnsense_<name>/resource.tf`

### Adding API Endpoints

To add support for additional OPNsense API endpoints:

1. **Identify the API endpoint** from [OPNsense API docs](https://docs.opnsense.org/development/api.html)

2. **Understand the request/response format** by inspecting browser network traffic

3. **Create the resource** following the pattern in existing resources

4. **Test thoroughly** with your OPNsense instance

### Modifying for Different OPNsense Versions

If API endpoints change in different OPNsense versions:

1. Update endpoint URLs in resource files
2. Update field names to match API changes
3. Update examples and documentation
4. Test with target version

## Common API Patterns

### Standard CRUD Pattern
Most OPNsense API endpoints follow this pattern:
- `add<Resource>` - Create (returns UUID)
- `get<Resource>/{uuid}` - Read
- `set<Resource>/{uuid}` - Update
- `del<Resource>/{uuid}` - Delete
- `reconfigure` or `apply` - Apply changes

### Search Pattern
Some resources support search:
- `search<Resource>` - Search with filters

### Utility Pattern
Some resources have utility endpoints:
- `toggle<Resource>/{uuid}` - Enable/disable
- `copy<Resource>/{uuid}` - Duplicate

## Security Considerations

### API Credentials
- Never commit credentials to version control
- Use environment variables or secret management
- Rotate API keys regularly

### TLS Verification
- Use valid certificates in production
- Only use `insecure = true` for testing

### Firewall Rule Management
- Test rules in staging environment first
- Always have out-of-band access
- Implement rule categorization
- Enable logging for auditing

### DHCP Management
- Avoid IP conflicts
- Document all static reservations
- Keep DHCP ranges separate from static IPs

## Troubleshooting

### Provider Issues

**Problem**: Provider not found
```bash
# Solution: Reinstall provider
make install
terraform init -upgrade
```

**Problem**: Authentication failure
```bash
# Solution: Verify credentials
curl -k -u "$OPNSENSE_API_KEY:$OPNSENSE_API_SECRET" \
  https://your-opnsense/api/core/firmware/status
```

### API Issues

**Problem**: Changes not applied
- Check OPNsense system logs
- Manually apply configuration in web UI
- Verify API user permissions

**Problem**: UUID not returned
- Check API response format
- Verify OPNsense version compatibility
- Check for API errors in logs

### Resource Issues

**Problem**: Resource already exists
```bash
# Solution: Import existing resource
terraform import opnsense_firewall_rule.my_rule <uuid>
```

**Problem**: Resource not deleted
- Check if resource is in use
- Verify dependencies
- Check OPNsense logs for errors

## Best Practices

### Infrastructure as Code
1. **Version control**: Use Git for all Terraform configurations
2. **State management**: Use remote state (S3, Terraform Cloud)
3. **Modules**: Create reusable modules for common patterns
4. **Variables**: Use variables for environment-specific values
5. **Documentation**: Document all resources and their purpose

### Firewall Management
1. **Rule order**: Remember OPNsense processes rules in order
2. **Default deny**: Implement default deny policy
3. **Logging**: Enable logging for security-relevant rules
4. **Categories**: Use categories to organize rules
5. **Descriptions**: Write clear, descriptive rule descriptions

### DHCP Management
1. **Network planning**: Plan IP address allocation carefully
2. **Documentation**: Document all static reservations
3. **Consistency**: Use consistent naming conventions
4. **Backup**: Keep backups of DHCP configurations

### VPN Management
1. **Key management**: Securely store WireGuard keys
2. **Rotation**: Regularly rotate keys and credentials
3. **Monitoring**: Monitor VPN connections and logs
4. **Segmentation**: Use separate tunnels for different purposes

## Next Steps

### Phase 1 (Current)
- ✅ Firewall rules and aliases
- ✅ Kea DHCP
- ✅ WireGuard VPN

### Phase 2 (Planned)
- NAT rules (source NAT, destination NAT)
- Traffic shaping
- Interface management
- Additional firewall features

### Phase 3 (Future)
- IPsec VPN
- OpenVPN
- System settings
- High availability configuration

## Contributing

To contribute to this provider:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Submit a pull request

## Resources

- [OPNsense Documentation](https://docs.opnsense.org/)
- [OPNsense API Reference](https://docs.opnsense.org/development/api.html)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Go Documentation](https://go.dev/doc/)

## Support

For issues, questions, or contributions:
- GitHub Issues: Report bugs and request features
- GitHub Discussions: Ask questions and share ideas
- OPNsense Forums: OPNsense-specific questions

## License

Mozilla Public License 2.0 (MPL-2.0)

---

**Created**: February 2026
**Target**: OPNsense 26.1 "Witty Woodpecker"
**Status**: Initial Release (v0.1.0)
