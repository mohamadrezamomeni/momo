package inbound

import "time"

type UpdateDto struct {
	Start        time.Time
	End          time.Time
	TrafficLimit uint32
}
