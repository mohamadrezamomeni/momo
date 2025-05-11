package inbound

type ExtendInboundDto struct {
	IdentifyInbounbdDto
	End string `json:"end_time"`
}
