package inbound

type ExtendInboundDto struct {
	IdentifyInbounbdDto
	End                  string `json:"end_time"`
	ExtendedTrafficLimit uint32 `json:"extended_traffic_limit"`
}
