package dto

import "momo/proxy/vpn"

type Add_VPN struct {
	Domain         string
	IsActive       bool
	ApiPort        string
	StartRangePort int
	EndRangePort   int
	VPNType        vpn.VPNType
	UserCount      int
}
