package vpnmanager

import (
	vpnManagerDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	"github.com/mohamadrezamomeni/momo/entity"
)

var (
	vpnRepo *VPN

	vpn1 = &vpnManagerDto.AddVPN{
		Country:  "uk",
		Domain:   "joi.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: false,
	}

	vpn2 = &vpnManagerDto.AddVPN{
		Country:  "uk",
		Domain:   "joi.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
	}

	vpn3 = &vpnManagerDto.AddVPN{
		Country:  "china",
		Domain:   "jordan.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
	}

	vpn4 = &vpnManagerDto.AddVPN{
		Country:  "colombia",
		Domain:   "jordan.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: false,
	}

	vpn5 = &vpnManagerDto.AddVPN{
		Country:  "china",
		Domain:   "joi.com",
		ApiPort:  "62733",
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
	}
)
