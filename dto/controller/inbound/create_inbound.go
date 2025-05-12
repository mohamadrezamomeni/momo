package inbound

type CreateInbound struct {
	Protocol     string `json:"protocol"`
	VPNType      string `json:"vpn_type"`
	Domain       string `json:"domain"`
	Port         string `json:"port"`
	UserID       string `json:"user_id"`
	Start        string `json:"start_time"`
	End          string `json:"end_time"`
	TrafficLimit uint32 `json:"traffic_limit"`
}
