package inbound

import (
	"fmt"
	"time"

	inboundDto "momo/dto/repository/inbound"
	"momo/entity"
	"momo/pkg/utils"

	"github.com/google/uuid"
)

var (
	port1 = "1081"
	port2 = "1082"
	port3 = "1083"
	port4 = "1084"
	port5 = "1085"
	port6 = "1086"
	port7 = "1087"

	userID1 = uuid.New().String()
	userID2 = uuid.New().String()
	userID3 = uuid.New().String()
	userID4 = uuid.New().String()
	userID5 = uuid.New().String()
	userID6 = uuid.New().String()
	userID7 = uuid.New().String()

	inbound1 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port1),
		Protocol: "vmess",
		IsBlock:  false,
		Port:     port1,
		Domain:   "google.com",
		UserID:   userID1,
		VPNType:  entity.XRAY_VPN,
		Start:    utils.GetDateTime("2024-04-21 14:30:00"),
		End:      utils.GetDateTime("2024-04-22 14:30:00"),
	}

	inbound2 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port2),
		Protocol: "vmess",
		Port:     port2,
		IsBlock:  false,
		Domain:   "twitter.com",
		UserID:   userID2,
		VPNType:  entity.XRAY_VPN,
		Start:    utils.GetDateTime("2024-04-21 14:30:00"),
		End:      utils.GetDateTime("2024-04-22 14:30:00"),
	}

	inbound3 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port2),
		Protocol: "http",
		Port:     port2,
		IsActive: false,
		Domain:   "twitter.com",
		UserID:   userID2,
		IsBlock:  false,
		VPNType:  entity.XRAY_VPN,
		Start:    utils.GetDateTime("2024-04-21 14:30:00"),
		End:      utils.GetDateTime("2024-04-22 14:30:00"),
	}

	inbound4 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port3),
		Protocol: "http",
		Port:     port3,
		Domain:   "googoo.com",
		UserID:   userID3,
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
		IsBlock:  false,
		Start:    time.Now().AddDate(0, 0, -15),
		End:      time.Now().AddDate(0, 0, 15),
	}

	inbound5 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port3),
		Protocol: "http",
		Port:     port4,
		Domain:   "googoo.com",
		UserID:   userID4,
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
		IsBlock:  true,
		Start:    time.Now().AddDate(0, 0, -15),
		End:      time.Now().AddDate(0, 0, 15),
	}

	inbound6 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port3),
		Protocol: "http",
		Port:     port5,
		Domain:   "googoo.com",
		UserID:   userID5,
		VPNType:  entity.XRAY_VPN,
		IsActive: false,
		IsBlock:  true,
		Start:    time.Now().AddDate(0, 0, -15),
		End:      time.Now().AddDate(0, 0, 15),
	}

	inbound7 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port3),
		Protocol: "http",
		Port:     port5,
		Domain:   "googoo.com",
		UserID:   userID6,
		VPNType:  entity.XRAY_VPN,
		IsActive: true,
		IsBlock:  false,
		Start:    time.Now().AddDate(0, -2, 0),
		End:      time.Now().AddDate(0, -1, 0),
	}

	inbound8 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port3),
		Protocol: "http",
		Port:     port6,
		Domain:   "googoo.com",
		UserID:   userID6,
		VPNType:  entity.XRAY_VPN,
		IsActive: false,
		IsBlock:  false,
		Start:    time.Now().AddDate(0, -2, 0),
		End:      time.Now().AddDate(0, -1, 0),
	}
	inbound9 = &inboundDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", port2),
		Protocol: "http",
		Port:     port7,
		Domain:   "twitter.com",
		UserID:   userID2,
		IsBlock:  true,
		IsActive: true,
		VPNType:  entity.XRAY_VPN,
		Start:    utils.GetDateTime("2024-04-21 14:30:00"),
		End:      utils.GetDateTime("2024-04-22 14:30:00"),
	}
)
