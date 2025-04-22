package dto

import "momo/proxy/vpn"

type FilterVPNs struct {
	IsActive *bool
	Domain   string
	VPNType  vpn.VPNType
}
