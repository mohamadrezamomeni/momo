package vpnpackage

type VPNPackageSerializer struct {
	ID                string `json:"id"`
	TrafficLimit      uint64 `json:"traffic_limit"`
	TrafficLimitTitle string `json:"traffic_limit_title"`
	Days              uint64 `json:"days"`
	Months            uint64 `json:"months"`
	Price             uint64 `json:"price"`
	PriceTitle        string `json:"price_title"`
	IsActive          bool   `json:"is_active"`
}
