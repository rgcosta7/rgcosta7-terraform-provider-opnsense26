# Complete OPNsense Infrastructure Example

terraform {
  required_version = ">= 1.0"
  required_providers {
    opnsense = {
      source  = "yourusername/opnsense"
      version = "0.1.0"
    }
  }
}

provider "opnsense" {
  host       = var.opnsense_host
  api_key    = var.opnsense_api_key
  api_secret = var.opnsense_api_secret
  insecure   = var.opnsense_insecure
}

# Variables
variable "opnsense_host" {
  description = "OPNsense host URL"
  type        = string
  default     = "https://192.168.1.1"
}

variable "opnsense_api_key" {
  description = "OPNsense API key"
  type        = string
  sensitive   = true
}

variable "opnsense_api_secret" {
  description = "OPNsense API secret"
  type        = string
  sensitive   = true
}

variable "opnsense_insecure" {
  description = "Skip TLS verification"
  type        = bool
  default     = true
}

variable "lan_network" {
  description = "LAN network CIDR"
  type        = string
  default     = "192.168.1.0/24"
}

variable "dmz_network" {
  description = "DMZ network CIDR"
  type        = string
  default     = "192.168.100.0/24"
}

variable "vpn_network" {
  description = "VPN tunnel network CIDR"
  type        = string
  default     = "10.20.30.0/24"
}

# ========================================
# Firewall Aliases
# ========================================

# Internal networks
resource "opnsense_firewall_alias" "internal_networks" {
  name        = "internal_networks"
  type        = "network"
  content     = [var.lan_network, var.dmz_network]
  description = "All internal networks"
  enabled     = true
}

# Web servers in DMZ
resource "opnsense_firewall_alias" "dmz_web_servers" {
  name        = "dmz_web_servers"
  type        = "host"
  content     = ["192.168.100.10", "192.168.100.11", "192.168.100.12"]
  description = "Web servers in DMZ"
  enabled     = true
}

# Database servers
resource "opnsense_firewall_alias" "database_servers" {
  name        = "database_servers"
  type        = "host"
  content     = ["192.168.1.20", "192.168.1.21"]
  description = "Database servers"
  enabled     = true
}

# Management hosts
resource "opnsense_firewall_alias" "management_hosts" {
  name        = "management_hosts"
  type        = "host"
  content     = ["192.168.1.5", "192.168.1.6"]
  description = "Management workstations"
  enabled     = true
}

# Common service ports
resource "opnsense_firewall_alias" "web_ports" {
  name        = "web_ports"
  type        = "port"
  content     = ["80", "443", "8080", "8443"]
  description = "Common web service ports"
  enabled     = true
}

# ========================================
# Firewall Rules - WAN Interface
# ========================================

# Allow HTTPS to DMZ web servers
resource "opnsense_firewall_rule" "wan_to_dmz_https" {
  description      = "Allow HTTPS to DMZ web servers"
  interface        = "wan"
  direction        = "in"
  protocol         = "tcp"
  source_net       = "any"
  destination_net  = opnsense_firewall_alias.dmz_web_servers.name
  destination_port = "443"
  action           = "pass"
  enabled          = true
  log              = true
  category         = "public_services"
}

# Allow WireGuard VPN
resource "opnsense_firewall_rule" "wan_wireguard" {
  description      = "Allow WireGuard VPN"
  interface        = "wan"
  direction        = "in"
  protocol         = "udp"
  source_net       = "any"
  destination_net  = "wan_address"
  destination_port = "51820"
  action           = "pass"
  enabled          = true
  log              = true
  category         = "vpn"
}

# ========================================
# Firewall Rules - LAN Interface
# ========================================

# Allow LAN to internet
resource "opnsense_firewall_rule" "lan_to_internet" {
  description     = "Allow LAN to internet"
  interface       = "lan"
  direction       = "in"
  protocol        = "any"
  source_net      = var.lan_network
  destination_net = "any"
  action          = "pass"
  enabled         = true
  log             = false
  category        = "internet_access"
}

# Allow management to SSH OPNsense
resource "opnsense_firewall_rule" "mgmt_to_ssh" {
  description      = "Allow management SSH to OPNsense"
  interface        = "lan"
  direction        = "in"
  protocol         = "tcp"
  source_net       = opnsense_firewall_alias.management_hosts.name
  destination_net  = "firewall"
  destination_port = "22"
  action           = "pass"
  enabled          = true
  log              = true
  category         = "management"
}

# Allow management to web interface
resource "opnsense_firewall_rule" "mgmt_to_webgui" {
  description      = "Allow management to web interface"
  interface        = "lan"
  direction        = "in"
  protocol         = "tcp"
  source_net       = opnsense_firewall_alias.management_hosts.name
  destination_net  = "firewall"
  destination_port = "443"
  action           = "pass"
  enabled          = true
  log              = true
  category         = "management"
}

# ========================================
# Firewall Rules - DMZ Interface
# ========================================

# Allow DMZ web servers to database
resource "opnsense_firewall_rule" "dmz_to_database" {
  description      = "Allow DMZ web to database"
  interface        = "dmz"
  direction        = "in"
  protocol         = "tcp"
  source_net       = opnsense_firewall_alias.dmz_web_servers.name
  destination_net  = opnsense_firewall_alias.database_servers.name
  destination_port = "3306"
  action           = "pass"
  enabled          = true
  log              = true
  category         = "application"
}

# Block DMZ to LAN (except database)
resource "opnsense_firewall_rule" "block_dmz_to_lan" {
  description     = "Block DMZ to LAN"
  interface       = "dmz"
  direction       = "in"
  protocol        = "any"
  source_net      = var.dmz_network
  destination_net = var.lan_network
  action          = "block"
  enabled         = true
  log             = true
  category        = "security"
}

# ========================================
# Kea DHCP Configuration
# ========================================

# LAN DHCP subnet
resource "opnsense_kea_subnet" "lan_dhcp" {
  subnet      = var.lan_network
  pools       = "192.168.1.100-192.168.1.200"
  description = "LAN DHCP subnet"
}

# DMZ DHCP subnet
resource "opnsense_kea_subnet" "dmz_dhcp" {
  subnet      = var.dmz_network
  pools       = "192.168.100.50-192.168.100.100"
  description = "DMZ DHCP subnet"
}

# Static reservations for management hosts
resource "opnsense_kea_reservation" "mgmt_host1" {
  subnet      = opnsense_kea_subnet.lan_dhcp.id
  ip_address  = "192.168.1.5"
  hw_address  = "00:11:22:33:44:55"
  hostname    = "mgmt-workstation-1"
  description = "Management workstation 1"
}

resource "opnsense_kea_reservation" "mgmt_host2" {
  subnet      = opnsense_kea_subnet.lan_dhcp.id
  ip_address  = "192.168.1.6"
  hw_address  = "00:11:22:33:44:66"
  hostname    = "mgmt-workstation-2"
  description = "Management workstation 2"
}

# Database server reservations
resource "opnsense_kea_reservation" "db_primary" {
  subnet      = opnsense_kea_subnet.lan_dhcp.id
  ip_address  = "192.168.1.20"
  hw_address  = "AA:BB:CC:DD:EE:01"
  hostname    = "db-primary"
  description = "Primary database server"
}

resource "opnsense_kea_reservation" "db_secondary" {
  subnet      = opnsense_kea_subnet.lan_dhcp.id
  ip_address  = "192.168.1.21"
  hw_address  = "AA:BB:CC:DD:EE:02"
  hostname    = "db-secondary"
  description = "Secondary database server"
}

# DMZ web server reservations
resource "opnsense_kea_reservation" "web1" {
  subnet      = opnsense_kea_subnet.dmz_dhcp.id
  ip_address  = "192.168.100.10"
  hw_address  = "BB:CC:DD:EE:FF:01"
  hostname    = "web-server-1"
  description = "Web server 1"
}

resource "opnsense_kea_reservation" "web2" {
  subnet      = opnsense_kea_subnet.dmz_dhcp.id
  ip_address  = "192.168.100.11"
  hw_address  = "BB:CC:DD:EE:FF:02"
  hostname    = "web-server-2"
  description = "Web server 2"
}

# ========================================
# WireGuard VPN Configuration
# ========================================

# WireGuard server
resource "opnsense_wireguard_server" "main_vpn" {
  name           = "wg0"
  enabled        = true
  listen_port    = 51820
  tunnel_address = "10.20.30.1/24"
  disable_routes = false
}

# Remote worker peers
resource "opnsense_wireguard_peer" "remote_worker_1" {
  name        = "remote-worker-1-laptop"
  enabled     = true
  public_key  = "worker1-public-key-replace-me"
  allowed_ips = "10.20.30.10/32"
  keepalive   = 25
}

resource "opnsense_wireguard_peer" "remote_worker_2" {
  name        = "remote-worker-2-laptop"
  enabled     = true
  public_key  = "worker2-public-key-replace-me"
  allowed_ips = "10.20.30.11/32"
  keepalive   = 25
}

# Mobile device peers
resource "opnsense_wireguard_peer" "mobile_device_1" {
  name        = "mobile-phone-1"
  enabled     = true
  public_key  = "mobile1-public-key-replace-me"
  allowed_ips = "10.20.30.20/32"
  keepalive   = 25
}

# Site-to-site VPN to branch office
resource "opnsense_wireguard_peer" "branch_office" {
  name          = "branch-office-vpn"
  enabled       = true
  public_key    = "branch-office-public-key-replace-me"
  allowed_ips   = "10.20.30.50/32,172.16.0.0/24"
  endpoint      = "branch.example.com"
  endpoint_port = 51820
  keepalive     = 25
}

# ========================================
# Outputs
# ========================================

output "firewall_rules_created" {
  description = "Number of firewall rules created"
  value = {
    wan_rules = 2
    lan_rules = 3
    dmz_rules = 2
  }
}

output "dhcp_subnets" {
  description = "DHCP subnet IDs"
  value = {
    lan = opnsense_kea_subnet.lan_dhcp.id
    dmz = opnsense_kea_subnet.dmz_dhcp.id
  }
}

output "wireguard_server_id" {
  description = "WireGuard server UUID"
  value       = opnsense_wireguard_server.main_vpn.id
}

output "wireguard_public_key" {
  description = "WireGuard server public key"
  value       = opnsense_wireguard_server.main_vpn.public_key
  sensitive   = false
}
