package dto

import (
	"time"

	"momo/proxy/vpn"
)

type CreateInbound struct {
	Protocol string
	Tag      string
	Port     string
	UserID   string
	Domain   string
	VPNType  vpn.VPNType
	IsActive bool
	Start    time.Time
	End      time.Time
	IsBlock  bool
}
