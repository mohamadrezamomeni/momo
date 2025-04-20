package dto

import "momo/proxy/vpn"

type FilterInbound struct {
	Protocol string
	IsActice *bool
	Domain   string
	VPNType  vpn.VPNType
	Port     string
	UserID   string
}
