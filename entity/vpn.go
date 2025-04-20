package entity

import "momo/proxy/vpn"

type VPN struct {
	ID             int
	Domain         string
	IsActive       bool
	ApiPort        string
	StartRangePort int
	EndRangePort   int
	VPNType        vpn.VPNType
	UserCount      int
}
