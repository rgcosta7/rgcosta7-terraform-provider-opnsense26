# Example: Create a DHCP subnet
resource "opnsense_kea_subnet" "lan_subnet" {
  subnet      = "192.168.1.0/24"
  pools       = "192.168.1.100-192.168.1.200"
  description = "LAN DHCP Subnet"
}

# Example: Create DHCP reservations
resource "opnsense_kea_reservation" "server1" {
  subnet      = opnsense_kea_subnet.lan_subnet.id
  ip_address  = "192.168.1.10"
  hw_address  = "00:11:22:33:44:55"
  hostname    = "server1"
  description = "Main application server"
}

resource "opnsense_kea_reservation" "printer" {
  subnet      = opnsense_kea_subnet.lan_subnet.id
  ip_address  = "192.168.1.50"
  hw_address  = "AA:BB:CC:DD:EE:FF"
  hostname    = "office-printer"
  description = "Office network printer"
}

# Example: Guest network subnet
resource "opnsense_kea_subnet" "guest_subnet" {
  subnet      = "10.10.10.0/24"
  pools       = "10.10.10.50-10.10.10.250"
  description = "Guest WiFi network"
}
