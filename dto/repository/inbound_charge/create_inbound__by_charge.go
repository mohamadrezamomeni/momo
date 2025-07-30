package inboundcharge

import (
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
)

type CreateInboundByCharge struct {
	TrafficLimit uint64
	TrafficUsage uint64
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
	Country      string
}
