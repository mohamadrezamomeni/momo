package inbound

import "time"

type SetPeriodDto struct {
	Start time.Time
	End   time.Time
}
