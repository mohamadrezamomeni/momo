package entity

import (
	"time"
)

type Inbound struct {
	IsActive bool
	ID       int
	Protocol string
	VPNType  VPNType
	Domain   string
	Port     string
	UserID   string
	Tag      string
	IsBlock  bool
	Start    time.Time
	End      time.Time
}
