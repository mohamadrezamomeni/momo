package entity

import (
	"time"
)

type Inbound struct {
	IsActive     bool
	ID           int
	ChargeCount  uint32
	TrafficUsage uint32
	TrafficLimit uint32
	Protocol     string
	VPNType      VPNType
	Domain       string
	Port         string
	IsNotified   bool
	IsAssigned   bool
	UserID       string
	Tag          string
	IsBlock      bool
	Start        time.Time
	End          time.Time
}
