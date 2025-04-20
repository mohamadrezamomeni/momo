package dto

import "momo/proxy/vpn"

type CreateInbound struct {
	Protocol    string
	Tag         string
	Port        string
	UserID      string
	Domain      string
	VPNType     vpn.VPNType
	IsAvailable bool
}
