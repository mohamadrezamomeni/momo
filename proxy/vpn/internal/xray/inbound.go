package xray

import (
	"context"
	"fmt"
	"strings"

	"momo/proxy/vpn/internal/xray/dto"

	"momo/pkg/utils"
	"momo/proxy/vpn/internal/xray/serializer"

	momoError "momo/pkg/error"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vmess"
	vmessInbound "github.com/xtls/xray-core/proxy/vmess/inbound"
)

func (x *Xray) addInbound(inpt *dto.AddInbound) (*serializer.AddInboundSerializer, error) {
	port, err := utils.ConvertToUint16(inpt.Port)
	if err != nil {
		return &serializer.AddInboundSerializer{}, momoError.Error("the port that is given is wrong")
	}

	users := make([]*protocol.User, 0)
	if inpt.User != nil {
		level, err := utils.ConvertToUint32(inpt.User.Level)
		if err != nil {
			return &serializer.AddInboundSerializer{}, momoError.Error("user's level is wrong.")
		}
		user := &protocol.User{
			Level: level,
			Email: inpt.User.Username,
			Account: serial.ToTypedMessage(&vmess.Account{
				Id: inpt.User.UUID,
			}),
		}
		users = append(users, user)
	}

	addInboundRequest := &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: inpt.Tag,
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortList: &net.PortList{
					Range: []*net.PortRange{net.SinglePortRange(net.Port(port))},
				},
				Listen: net.NewIPOrDomain(net.AnyIP),
				SniffingSettings: &proxyman.SniffingConfig{
					Enabled:             true,
					DestinationOverride: []string{"http", "tls"},
				},
			}),

			ProxySettings: serial.ToTypedMessage(&vmessInbound.Config{
				User: users,
			}),
		},
	}

	_, err = x.hsClient.AddInbound(context.Background(), addInboundRequest)
	if err != nil {
		return &serializer.AddInboundSerializer{}, momoError.Errorf("the error has happend the problem was %v", err)
	}
	return &serializer.AddInboundSerializer{}, nil
}

func (x *Xray) isUserFilled(u *dto.InboundUser) bool {
	if u.Username != "" && u.Level != "" && u.UUID != "" {
		return true
	}
	return false
}

func (x *Xray) removeInbound(inpt *dto.RemoveInbound) (*serializer.RemoveInbound, error) {
	client := x.hsClient
	_, err := client.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: inpt.Tag,
	})
	return &serializer.RemoveInbound{}, err
}

func (x *Xray) receiveInboundTraffic(inpt *dto.ReceiveInboundTraffic) (*serializer.ReceiveInboundTraffic, error) {
	ptn := fmt.Sprintf("inbound>>>%s>>>traffic", inpt.Tag)
	resp, err := x.ssClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: ptn,
		Reset_:  false,
	})
	stats := resp.GetStat()
	if err != nil {
		return &serializer.ReceiveInboundTraffic{}, err
	}
	if len(stats) == 0 {
		return &serializer.ReceiveInboundTraffic{}, momoError.Error("result wasn't found")
	}

	data, err := x.convertStatsToMap(stats)

	return &serializer.ReceiveInboundTraffic{
		UpLink:   data["uplink"],
		DownLink: data["downlink"],
	}, nil
}

func (x *Xray) convertStatsToMap(stats []*statsService.Stat) (map[string]int64, error) {
	res := map[string]int64{}
	for _, stat := range stats {
		if strings.Contains(stat.Name, "uplink") {
			res["uplink"] = x.getValStat(stat)
		} else if strings.Contains(stat.Name, "downlink") {
			res["downlink"] = x.getValStat(stat)
		} else {
			return map[string]int64{}, momoError.Error("something went wrong. we faced unexpected situation")
		}
	}
	return res, nil
}

func (x *Xray) getValStat(stat *statsService.Stat) int64 {
	if stat.Value == 0 {
		return 0
	}
	return stat.Value
}
