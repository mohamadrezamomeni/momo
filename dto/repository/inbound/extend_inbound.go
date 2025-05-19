package dto

import "time"

type ExtendInboundDto struct {
	End             time.Time
	Start           time.Time
	TrafficExtended uint64
}
