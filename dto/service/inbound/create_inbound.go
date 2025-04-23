package inbound

import (
	"time"

	"momo/entity"
)

type CreateInbound struct {
	UserID     string
	Start      time.Time
	End        time.Time
	VPNType    entity.VPNType
	ServerType entity.HostStatus
}
