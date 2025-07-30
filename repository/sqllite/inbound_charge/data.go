package inboundcharge

import (
	"time"

	"github.com/google/uuid"
	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	inboundRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	inboundChargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound_charge"
	vpnPackageRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
)

var (
	userID1  = uuid.New().String()
	inbound1 = &inboundRepositoryDto.CreateInbound{
		TrafficLimit: 500000,
		TrafficUsage: 0,
		Protocol:     "vmess",
		Tag:          "inbound-33333",
		Port:         "323",
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsBlock:      false,
		IsAssigned:   true,
		IsActive:     true,
		Start:        time.Now().AddDate(0, -2, 0),
		End:          time.Now().AddDate(0, -1, 0),
		UserID:       userID1,
	}
	vpnPackage1 = &vpnPackageRepositoryDto.CreateVPNPackage{
		TrafficLimitTitle: "50G",
		TrafficLimit:      500000,
		Price:             2000000,
		PriceTitle:        "1$",
		IsActive:          true,
		Days:              0,
		Months:            1,
	}
	charge1 = &chargeRepositoryDto.CreateDto{
		Detail: "test",
		Status: entity.ApprovedStatusCharge,
		UserID: userID1,
	}

	inbound2 = &inboundChargeRepositoryDto.CreateInboundByCharge{
		TrafficLimit: 500000,
		TrafficUsage: 0,
		Protocol:     "vmess",
		Tag:          "inbound-23456",
		Port:         "323",
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsBlock:      false,
		IsAssigned:   true,
		IsActive:     true,
		Start:        time.Now().AddDate(0, -2, 0),
		End:          time.Now().AddDate(0, -1, 0),
		UserID:       userID1,
	}
	charge2 = &chargeRepositoryDto.CreateDto{
		Detail: "test",
		Status: entity.ApprovedStatusCharge,
		UserID: userID1,
	}
)
