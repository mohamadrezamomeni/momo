package xray

import (
	"context"
	"fmt"

	"momo/proxy/xray/dto"

	"momo/pkg/utils"
	"momo/proxy/xray/serializer"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/protocol/tls/cert"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/proxy/vmess/inbound"
	vmessInbound "github.com/xtls/xray-core/proxy/vmess/inbound"
	"github.com/xtls/xray-core/transport/internet"
	_ "github.com/xtls/xray-core/transport/internet"
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	"github.com/xtls/xray-core/transport/internet/tls"
	_ "github.com/xtls/xray-core/transport/internet/tls"
	_ "github.com/xtls/xray-core/transport/internet/websocket"
)

func (x *Xray) AddInbound(inpt *dto.AddInbound) (*serializer.AddInboundSerializer, error) {
	port, err := utils.ConvertToUint16(inpt.Port)
	if err != nil {
		return &serializer.AddInboundSerializer{}, fmt.Errorf("the port that is given is wrong")
	}
	client := x.HsClient
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
				StreamSettings: &internet.StreamConfig{
					ProtocolName: "websocket",

					SecurityType: serial.GetMessageType(&tls.Config{}),
					SecuritySettings: []*serial.TypedMessage{
						serial.ToTypedMessage(&tls.Config{
							Certificate: []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil))},
						}),
					},
				},
			}),

			ProxySettings: serial.ToTypedMessage(&vmessInbound.Config{
				User: []*protocol.User{},
			}),
		},
	}

	_, err = client.AddInbound(context.Background(), addInboundRequest)
	return &serializer.AddInboundSerializer{}, err
}

func (x *Xray) RemoveInbound() {}

func (x *Xray) QueryInbound() {}
