package vpnpackage

import (
	vpnPackageRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
)

var (
	vpnPackage1 = &vpnPackageRepositoryDto.CreateVPNPackage{
		Price:             10000,
		PriceTitle:        "10$",
		TrafficLimitTitle: "10G",
		TrafficLimit:      100000,
		Days:              10,
		Months:            0,
		IsActive:          true,
	}
	vpnPackage2 = &vpnPackageRepositoryDto.CreateVPNPackage{
		Price:             20000,
		PriceTitle:        "20$",
		TrafficLimitTitle: "20G",
		TrafficLimit:      100000,
		Days:              10,
		Months:            0,
		IsActive:          false,
	}

	vpnPackage3 = &vpnPackageRepositoryDto.CreateVPNPackage{
		Price:             30000,
		PriceTitle:        "30$",
		TrafficLimitTitle: "30G",
		TrafficLimit:      300000,
		Days:              0,
		Months:            3,
		IsActive:          false,
	}
)
