package dto

import "time"

type ExtendInboundDto struct {
	End             time.Time
	TrafficExtended uint32
}
