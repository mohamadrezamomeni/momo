package dto

import (
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
)

type CreateInbound struct {
	TrafficLimit uint32
	TrafficUsage uint32
	Protocol     string
	Tag          string
	Port         string
	UserID       string
	Domain       string
	VPNType      entity.VPNType
	IsAssigned   bool
	IsNotified   bool
	IsActive     bool
	Start        time.Time
	End          time.Time
	IsBlock      bool
}
