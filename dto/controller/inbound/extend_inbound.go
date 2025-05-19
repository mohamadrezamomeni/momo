package inbound

type ExtendInboundDto struct {
	IdentifyInbounbdDto
	End                  string `json:"end_time"`
	ExtendedTrafficLimit uint64 `json:"extended_traffic_limit"`
}
