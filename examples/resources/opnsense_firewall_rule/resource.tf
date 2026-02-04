# Example: Allow HTTP traffic from LAN to WAN
resource "opnsense_firewall_rule" "allow_http" {
  description      = "Allow HTTP traffic from LAN"
  interface        = "lan"
  direction        = "in"
  ip_protocol      = "inet"
  protocol         = "tcp"
  source_net       = "192.168.1.0/24"
  destination_net  = "any"
  destination_port = "80"
  action           = "pass"
  enabled          = true
  log              = true
  category         = "web_access"
}

# Example: Block specific IP
resource "opnsense_firewall_rule" "block_malicious_ip" {
  description     = "Block malicious IP"
  interface       = "wan"
  direction       = "in"
  protocol        = "any"
  source_net      = "203.0.113.0/24"
  destination_net = "any"
  action          = "block"
  enabled         = true
  log             = true
}

# Example: Allow SSH from management network
resource "opnsense_firewall_rule" "allow_ssh_mgmt" {
  description      = "Allow SSH from management network"
  interface        = "lan"
  protocol         = "tcp"
  source_net       = "10.0.0.0/24"
  destination_net  = "192.168.1.1"
  destination_port = "22"
  action           = "pass"
  enabled          = true
  log              = true
}
