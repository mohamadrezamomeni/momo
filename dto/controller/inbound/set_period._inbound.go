package inbound

type SetPeriodDto struct {
	IdentifyInbounbdDto
	Start string `json:"start_time"`
	End   string `json:"end_time"`
}
