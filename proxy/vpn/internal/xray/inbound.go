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
	scope := "xrayProxy.addInbound"
	port, err := utils.ConvertToUint16(inpt.Port)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	users := make([]*protocol.User, 0)
	if inpt.User != nil {
		level, err := utils.ConvertToUint32(inpt.User.Level)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", inpt)
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
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", inpt)
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
	scope := "xrayProxy.removeInbound"
	client := x.hsClient
	_, err := client.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: inpt.Tag,
	})
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", inpt)
	}
	return &serializer.RemoveInbound{}, err
}

func (x *Xray) resetTraffic(tag string) error {
	scope := "xrayProxy.resetTraffic"
	_, err := x.getInboundTrafficWithoutBeigReseted(tag)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the input is %s", tag)
	}
	return nil
}

func (x *Xray) getInboundTrafficWithoutBeigReseted(tag string) (*serializer.ReceiveInboundTraffic, error) {
	return x.receiveInboundTraffic(tag, false)
}

func (x *Xray) getInboundTrafficWithBeigReseted(tag string) (*serializer.ReceiveInboundTraffic, error) {
	return x.receiveInboundTraffic(tag, true)
}

func (x *Xray) receiveInboundTraffic(tag string, reset bool) (*serializer.ReceiveInboundTraffic, error) {
	scope := "xrayProxy.receiveInboundTraffic"

	ptn := fmt.Sprintf("inbound>>>%s>>>traffic", tag)
	resp, err := x.ssClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: ptn,
		Reset_:  reset,
	})
	stats := resp.GetStat()
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the tag is %s and reset is %v", tag, reset)
	}
	if len(stats) == 0 {
		return nil, momoError.Scope(scope).Errorf("the tag is %s and reset is %v", tag, reset)
	}

	data, err := x.convertStatsToMap(stats)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the tag is %s and reset is %v", tag, reset)
	}

	return &serializer.ReceiveInboundTraffic{
		UpLink:   data["uplink"],
		DownLink: data["downlink"],
	}, nil
}

func (x *Xray) convertStatsToMap(stats []*statsService.Stat) (map[string]int64, error) {
	scope := "xrayProxy.convertStatsToMap"

	res := map[string]int64{}
	for _, stat := range stats {
		if strings.Contains(stat.Name, "uplink") {
			res["uplink"] = x.getValStat(stat)
		} else if strings.Contains(stat.Name, "downlink") {
			res["downlink"] = x.getValStat(stat)
		} else {
			return map[string]int64{}, momoError.Scope(scope).Errorf("the input is %+v", stats)
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

func (x *Xray) fakeReceiveInboundTraffic() error {
	scope := "xrayProxy.fakeReceiveInboundTraffic"

	_, err := x.ssClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{})
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}
