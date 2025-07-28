package inbound

import (
	"time"

	inboundRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/entity"

	"github.com/google/uuid"
)

var (
	now = time.Now()

	inbound1 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   false,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		Country:      "china",
		TrafficLimit: 50000,
	}
	inbound2 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "china",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   false,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}
	inbound3 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Port:         "332",
		Country:      "us",
		UserID:       uuid.New().String(),
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}

	inbound4 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, 0, -3),
		End:          now.AddDate(0, 0, 27),
		IsBlock:      false,
		TrafficLimit: 50000,
	}
	inbound5 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   false,
		IsActive:     true,
		Start:        now.AddDate(0, 0, -2),
		End:          now.AddDate(0, 0, 27),
		IsBlock:      true,
		TrafficLimit: 50000,
	}
	inbound6 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, 0, -30),
		End:          now.AddDate(0, 0, -1),
		IsBlock:      true,
		TrafficLimit: 50000,
	}

	inbound7 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "uk",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   false,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}

	inbound8 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "uk",
		Port:         "3333",
		UserID:       uuid.New().String(),
		Domain:       "facebook.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   true,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}

	inbound9 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "uk",
		Port:         "2222",
		UserID:       uuid.New().String(),
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   false,
		IsActive:     true,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}

	inbound10 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "uk",
		Port:         "2222",
		UserID:       uuid.New().String(),
		Domain:       "instagram.com",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   true,
		IsNotified:   false,
		IsActive:     true,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}
	inbound11 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "brazil",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   false,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}
	inbound12 = &inboundRepositoryDto.CreateInbound{
		Protocol:     "vmess",
		Tag:          "example-tag",
		Country:      "brazil",
		Port:         "",
		UserID:       uuid.New().String(),
		Domain:       "",
		VPNType:      entity.XRAY_VPN,
		IsAssigned:   false,
		IsNotified:   false,
		IsActive:     false,
		Start:        now.AddDate(0, -1, 0),
		End:          now,
		IsBlock:      false,
		TrafficLimit: 50000,
	}
)
