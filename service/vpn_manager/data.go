package vpnmanager

import (
	vpnManagerDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	"github.com/mohamadrezamomeni/momo/entity"
)

var (
	vpn1 = &vpnManagerDto.AddVPN{
		VPNType:   entity.XRAY_VPN,
		Domain:    "twitter.com",
		ApiPort:   "2002",
		UserCount: 20,
		IsActive:  true,
		StartPort: 2000,
		EndPort:   3000,
	}

	vpn2 = &vpnManagerDto.AddVPN{
		VPNType:   entity.XRAY_VPN,
		Domain:    "google.com",
		ApiPort:   "2002",
		UserCount: 20,
		IsActive:  true,
		StartPort: 2000,
		EndPort:   3000,
	}

	vpn3 = &vpnManagerDto.AddVPN{
		VPNType:   entity.XRAY_VPN,
		Domain:    "facebook.com",
		ApiPort:   "2002",
		UserCount: 20,
		IsActive:  false,
		StartPort: 2000,
		EndPort:   3000,
	}
)
