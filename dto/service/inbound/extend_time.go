package inbound

import "time"

type ExtendInboundDto struct {
	Start                time.Time
	End                  time.Time
	ExtendedTrafficLimit uint64
}
