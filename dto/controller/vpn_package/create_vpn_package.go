package vpnpackage

type CreateVPNPackage struct {
	TrafficLimitTitle string `json:"traffic_limit_title"`
	TrafficLimit      uint64 `json:"traffic_limit"`
	Price             uint64 `json:"price"`
	PriceTitle        string `json:"price_title"`
	IsActive          bool   `json:"is_active"`
	Days              uint64 `json:"days"`
	Months            uint64 `json:"months"`
}
