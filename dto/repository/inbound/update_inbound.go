package dto

import "time"

type UpdateInboundDto struct {
	Start        time.Time
	End          time.Time
	TrafficLimit uint32
}
