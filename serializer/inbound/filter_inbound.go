package inbound

type FilterInboundsSerializer struct {
	Inbounds []*InboundSerializer `json:"inbounds"`
}
