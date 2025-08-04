package entity

import (
	"time"
)

type Inbound struct {
	IsActive     bool
	ID           string
	TrafficUsage uint32
	TrafficLimit uint32
	Country      string
	Protocol     string
	VPNType      VPNType
	Domain       string
	Port         string
	VPNID        string
	IsNotified   bool
	IsAssigned   bool
	UserID       string
	Tag          string
	IsBlock      bool
	Start        time.Time
	End          time.Time
}
