package entity

type VPNPackage struct {
	ID                string
	TrafficLimit      uint64
	TrafficLimitTitle string
	Days              uint64
	Months            uint64
	Price             uint64
	PriceTitle        string
	IsActive          bool
}
