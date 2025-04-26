package vpnmanager

import (
	vpnManagerDto "momo/dto/repository/vpn_manager"
	"momo/entity"
)

var (
	vpnRepo *VPN

	vpn1 = &vpnManagerDto.Add_VPN{
		Domain:   "joi.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: false,
	}

	vpn2 = &vpnManagerDto.Add_VPN{
		Domain:   "joi.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
	}

	vpn3 = &vpnManagerDto.Add_VPN{
		Domain:   "jordan.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
	}
)
