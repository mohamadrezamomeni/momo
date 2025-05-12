package inbound

import "time"

type ExtendInboundDto struct {
	End                  time.Time
	ExtendedTrafficLimit uint32
}
