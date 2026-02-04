# OPNsense API Endpoint Naming Guide

## 🔍 Discovery

Thanks to your research, we found the **actual** API documentation showing endpoint naming!

## 📋 Naming Conventions by Module

### ✅ Firewall Aliases - **camelCase**
- `/api/firewall/alias/addItem`
- `/api/firewall/alias/setItem`
- `/api/firewall/alias/delItem`
- `/api/firewall/alias/reconfigure`

**Status**: ✅ WORKING (tested) + Fixed content separator (newline not comma)

### ✅ Destination NAT - **snake_case** 
- `/api/firewall/d_nat/add_rule`
- `/api/firewall/d_nat/set_rule`
- `/api/firewall/d_nat/del_rule`
- `/api/firewall/d_nat/get_rule`

**Status**: ⏳ To be implemented

### ✅ Kea DHCP - **snake_case**
- `/api/kea/dhcpv4/add_subnet`
- `/api/kea/dhcpv4/set_subnet`
- `/api/kea/dhcpv4/del_subnet`
- `/api/kea/dhcpv4/get_subnet`
- `/api/kea/dhcpv4/add_reservation`
- `/api/kea/dhcpv4/set_reservation`
- `/api/kea/dhcpv4/del_reservation`
- `/api/kea/dhcpv4/get_reservation`

**Status**: ✅ FIXED (was camelCase, now snake_case)

### ✅ WireGuard - **snake_case**
- `/api/wireguard/server/add_server`
- `/api/wireguard/server/set_server`
- `/api/wireguard/server/del_server`
- `/api/wireguard/server/get_server`
- `/api/wireguard/client/add_client`
- `/api/wireguard/client/set_client`
- `/api/wireguard/client/del_client`
- `/api/wireguard/client/get_client`

**Status**: ✅ FIXED (was camelCase, now snake_case)

### ✅ Firewall Filter Rules - **camelCase**
- `/api/firewall/filter/addRule`
- `/api/firewall/filter/setRule`
- `/api/firewall/filter/delRule`
- `/api/firewall/filter/getRule`
- `/api/firewall/filter/apply`

**Status**: ✅ Code uses correct naming

## 🎯 Pattern Recognition

**OLD MVC APIs (pre-26.1 plugins)**: Use **camelCase**
- Firewall aliases
- Firewall filter rules

**NEW MVC APIs (26.1+)**: Use **snake_case**
- Destination NAT
- Source NAT  
- Kea DHCP
- WireGuard

## 📝 What We Fixed

### Firewall Aliases
- Content separator: Comma → **Newline** ✅

### Kea DHCP (ALL 8 endpoints)
Changed from camelCase → **snake_case**:
- `addSubnet` → `add_subnet` ✅
- `setSubnet` → `set_subnet` ✅
- `delSubnet` → `del_subnet` ✅
- `getSubnet` → `get_subnet` ✅
- `addReservation` → `add_reservation` ✅
- `setReservation` → `set_reservation` ✅  
- `delReservation` → `del_reservation` ✅
- `getReservation` → `get_reservation` ✅

### WireGuard (ALL 8 endpoints)
Changed from camelCase → **snake_case**:
- `addServer` → `add_server` ✅
- `setServer` → `set_server` ✅
- `delServer` → `del_server` ✅
- `getServer` → `get_server` ✅
- `addClient` → `add_client` ✅
- `setClient` → `set_client` ✅
- `delClient` → `del_client` ✅
- `getClient` → `get_client` ✅

## 🚀 Ready to Test

All endpoints are now correctly named! You can test:

1. **Firewall Aliases** ✅ (already tested, working!)
2. **Firewall Rules** ✅ (ready to test)
3. **Kea DHCP** ✅ (fixed, ready to test)
4. **WireGuard VPN** ✅ (fixed, ready to test)

## 🎯 For Your Traefik Setup

Based on the **d_nat** API you showed, the Destination NAT endpoints are:
- `POST /api/firewall/d_nat/add_rule`
- `POST /api/firewall/d_nat/set_rule/{uuid}`
- `POST /api/firewall/d_nat/del_rule/{uuid}`
- `GET /api/firewall/d_nat/get_rule/{uuid}`

Perfect for port 443 → Traefik!

## 📚 Resources for API Discovery

1. **Browser DevTools** - Watch network tab when using OPNsense GUI
2. **OPNsense Source Code** - Check controller PHP files:
   - https://github.com/opnsense/core/tree/master/src/opnsense/mvc/app/controllers
3. **Model XML Files** - Define data structures:
   - https://github.com/opnsense/core/tree/master/src/opnsense/mvc/app/models
4. **API Docs** - https://docs.opnsense.org/development/api/
5. **Forum** - Community discussions

## ✅ Verification Checklist

Before implementing a new resource:

1. Find the controller PHP file
2. Check model XML for data structure
3. Test with curl to see actual request/response
4. Note if camelCase or snake_case (NEW APIs use snake_case!)
5. Check separators (comma vs newline vs other)
6. Verify UUID field name in response
7. Find apply/reconfigure endpoint

## 🎯 Next: Destination NAT Implementation

Once we implement Destination NAT:

```hcl
resource "opnsense_nat_destination" "traefik_https" {
  interface        = "wan"
  protocol         = "tcp"
  destination_port = "443"
  target_ip        = "192.168.1.100"
  target_port      = "443"
  description      = "HTTPS to Traefik"
  enabled          = true
}

resource "opnsense_nat_destination" "traefik_http" {
  interface        = "wan"
  protocol         = "tcp"
  destination_port = "80"
  target_ip        = "192.168.1.100"
  target_port      = "80"
  description      = "HTTP to Traefik"
  enabled          = true
}
```

Perfect for your setup! 🎉
