package dto

import "momo/proxy/vpn"

type FilterInbound struct {
	Protocol    string
	IsAvailable *bool
	Domain      string
	VPNType     vpn.VPNType
	Port        string
	UserID      string
}
