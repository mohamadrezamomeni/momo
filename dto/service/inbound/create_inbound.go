package inbound

import (
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
)

type CreateInbound struct {
	UserID       string
	Start        time.Time
	End          time.Time
	VPNType      entity.VPNType
	ServerType   entity.HostStatus
	TrafficLimit uint32
}
