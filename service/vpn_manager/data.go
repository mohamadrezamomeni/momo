package vpnmanager

import (
	vpnManagerDto "momo/dto/repository/vpn_manager"
	"momo/entity"
)

var (
	vpn1 = &vpnManagerDto.Add_VPN{
		VPNType:   entity.XRAY_VPN,
		Domain:    "twitter.com",
		ApiPort:   "2002",
		UserCount: 20,
		IsActive:  true,
	}

	vpn2 = &vpnManagerDto.Add_VPN{
		VPNType:   entity.XRAY_VPN,
		Domain:    "google.com",
		ApiPort:   "2002",
		UserCount: 20,
		IsActive:  true,
	}

	vpn3 = &vpnManagerDto.Add_VPN{
		VPNType:   entity.XRAY_VPN,
		Domain:    "facebook.com",
		ApiPort:   "2002",
		UserCount: 20,
		IsActive:  false,
	}
)
