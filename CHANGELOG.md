# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-04

### Added
- Initial release
- Provider configuration with API key/secret authentication
- Support for OPNsense 26.1 API
- Firewall resources:
  - `opnsense_firewall_rule` - Manage firewall filter rules
  - `opnsense_firewall_alias` - Manage firewall aliases
- Kea DHCP resources:
  - `opnsense_kea_subnet` - Manage DHCP subnets
  - `opnsense_kea_reservation` - Manage DHCP reservations
- WireGuard VPN resources:
  - `opnsense_wireguard_server` - Manage WireGuard server instances
  - `opnsense_wireguard_peer` - Manage WireGuard peers/clients
- Data sources:
  - `opnsense_firewall_rule` - Fetch firewall rule information
- Automatic API key/secret from environment variables
- TLS certificate verification with option to skip (insecure mode)
- Comprehensive examples and documentation

### Notes
- This is the initial release targeting OPNsense 26.1 API
- Tested with OPNsense 26.1 "Witty Woodpecker"
- Provider uses HashiCorp Terraform Plugin Framework

## [Unreleased]

### Planned Features
- NAT rules support (source NAT, destination NAT/port forwarding)
- Traffic shaping rules
- VPN IPsec support
- OpenVPN support
- Interface management
- System settings management
- Additional firewall features (schedules, categories, advanced options)
- IPv6 extensive testing and improvements
- More data sources
