package vpnpackage

import (
	vpnPackageRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
)

var vpnPackage1 = &vpnPackageRepositoryDto.CreateVPNPackage{
	Price:             10000,
	PriceTitle:        "10$",
	TrafficLimitTitle: "10G",
	TrafficLimit:      100000,
	Days:              10,
	Months:            0,
	IsActive:          true,
}
