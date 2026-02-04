# Example: Create a network alias for internal servers
resource "opnsense_firewall_alias" "internal_servers" {
  name        = "internal_servers"
  type        = "network"
  content     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  description = "Internal server networks"
  enabled     = true
}

# Example: Create a host alias for specific IPs
resource "opnsense_firewall_alias" "dns_servers" {
  name        = "dns_servers"
  type        = "host"
  content     = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
  description = "Public DNS servers"
  enabled     = true
}

# Example: Create a port alias
resource "opnsense_firewall_alias" "web_ports" {
  name        = "web_ports"
  type        = "port"
  content     = ["80", "443", "8080", "8443"]
  description = "Common web service ports"
  enabled     = true
}

# Example: Use alias in a firewall rule
resource "opnsense_firewall_rule" "allow_dns" {
  description     = "Allow DNS to public servers"
  interface       = "lan"
  protocol        = "udp"
  source_net      = "any"
  destination_net = opnsense_firewall_alias.dns_servers.name
  destination_port = "53"
  action          = "pass"
  enabled         = true
}
