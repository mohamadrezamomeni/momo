package vpnpackage

type CreateVPNPackage struct {
	TrafficLimitTitle string
	TrafficLimit      uint64
	Price             uint64
	PriceTitle        string
	IsActive          bool
	Days              uint64
	Months            uint64
}
