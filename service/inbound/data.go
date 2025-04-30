package inbound

import (
	"time"

	inboundRepositoryDto "momo/dto/repository/inbound"
	"momo/entity"

	"github.com/google/uuid"
)

var (
	now = time.Now()

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
		Start:      now.AddDate(0, -1, 0),
		End:        now,
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
		Start:      now.AddDate(0, -1, 0),
		End:        now,
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
		Start:      now.AddDate(0, -1, 0),
		End:        now,
		IsBlock:    false,
	}

	inbound4 = &inboundRepositoryDto.CreateInbound{
		Protocol:   "vmess",
		Tag:        "example-tag",
		Port:       "",
		UserID:     uuid.New().String(),
		Domain:     "instagram.com",
		VPNType:    entity.XRAY_VPN,
		IsAssigned: true,
		IsNotified: false,
		IsActive:   false,
		Start:      now.AddDate(0, 0, -3),
		End:        now.AddDate(0, 0, 27),
		IsBlock:    false,
	}
	inbound5 = &inboundRepositoryDto.CreateInbound{
		Protocol:   "vmess",
		Tag:        "example-tag",
		Port:       "",
		UserID:     uuid.New().String(),
		Domain:     "instagram.com",
		VPNType:    entity.XRAY_VPN,
		IsAssigned: true,
		IsNotified: false,
		IsActive:   true,
		Start:      now.AddDate(0, 0, -2),
		End:        now.AddDate(0, 0, 27),
		IsBlock:    true,
	}
	inbound6 = &inboundRepositoryDto.CreateInbound{
		Protocol:   "vmess",
		Tag:        "example-tag",
		Port:       "",
		UserID:     uuid.New().String(),
		Domain:     "instagram.com",
		VPNType:    entity.XRAY_VPN,
		IsAssigned: true,
		IsNotified: false,
		IsActive:   false,
		Start:      now.AddDate(0, 0, -30),
		End:        now.AddDate(0, 0, -1),
		IsBlock:    true,
	}
)
