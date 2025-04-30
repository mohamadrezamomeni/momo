package inbound

import (
	"time"

	inboundRepositoryDto "momo/dto/repository/inbound"
	"momo/entity"

	"github.com/google/uuid"
)

var (
	inbound1 = &inboundRepositoryDto.CreateInbound{
		Protocol:   "vmess",
		Tag:        "example-tag",
		Port:       "",
		UserID:     uuid.New().String(),
		Domain:     "",
		VPNType:    entity.XRAY_VPN,
		IsAssigned: false,
		IsNotified: false,
		IsActive:   false,
		Start:      time.Now().AddDate(0, -1, 0),
		End:        time.Now(),
		IsBlock:    false,
	}
	inbound2 = &inboundRepositoryDto.CreateInbound{
		Protocol:   "vmess",
		Tag:        "example-tag",
		Port:       "",
		UserID:     uuid.New().String(),
		Domain:     "",
		VPNType:    entity.XRAY_VPN,
		IsAssigned: false,
		IsNotified: false,
		IsActive:   false,
		Start:      time.Now().AddDate(0, -1, 0),
		End:        time.Now(),
		IsBlock:    false,
	}
	inbound3 = &inboundRepositoryDto.CreateInbound{
		Protocol:   "vmess",
		Tag:        "example-tag",
		Port:       "",
		UserID:     uuid.New().String(),
		Domain:     "instagram.com",
		VPNType:    entity.XRAY_VPN,
		IsAssigned: true,
		IsNotified: false,
		IsActive:   false,
		Start:      time.Now().AddDate(0, -1, 0),
		End:        time.Now(),
		IsBlock:    false,
	}
)
