package entity

import "momo/proxy/vpn"

type Inbound struct {
	IsActive bool
	ID       int
	Protocol string
	VPNType  vpn.VPNType
	Domain   string
	Port     string
	UserID   string
	Tag      string
	IsBlock  bool
}
