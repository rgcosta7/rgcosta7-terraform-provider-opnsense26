# Example: Create a WireGuard server
resource "opnsense_wireguard_server" "vpn_server" {
  name            = "wg0"
  enabled         = true
  listen_port     = 51820
  tunnel_address  = "10.20.30.1/24"
  disable_routes  = false
}

# Example: Create WireGuard peers
resource "opnsense_wireguard_peer" "laptop" {
  name          = "laptop"
  enabled       = true
  public_key    = "your-laptop-public-key-here"
  allowed_ips   = "10.20.30.10/32"
  keepalive     = 25
}

resource "opnsense_wireguard_peer" "mobile" {
  name          = "mobile-phone"
  enabled       = true
  public_key    = "your-mobile-public-key-here"
  allowed_ips   = "10.20.30.11/32"
  keepalive     = 25
}

# Example: Site-to-site VPN peer
resource "opnsense_wireguard_peer" "remote_office" {
  name          = "remote-office"
  enabled       = true
  public_key    = "remote-office-public-key"
  allowed_ips   = "10.20.30.20/32,192.168.100.0/24"
  endpoint      = "remote.example.com"
  endpoint_port = 51820
  keepalive     = 25
}

# Link peers to server
resource "opnsense_wireguard_server" "vpn_with_peers" {
  name            = "wg1"
  enabled         = true
  listen_port     = 51821
  tunnel_address  = "10.30.40.1/24"
  peers           = [
    opnsense_wireguard_peer.laptop.id,
    opnsense_wireguard_peer.mobile.id
  ]
}
