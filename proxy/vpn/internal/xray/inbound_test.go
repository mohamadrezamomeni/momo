package xray

import (
	"testing"

	"momo/proxy/vpn/internal/xray/dto"
)

var xray *Xray = New(&XrayConfig{
	Address: "192.168.116.129",
	ApiPort: "62789",
})

var (
	protocoll string = "vmess"
	port      string = "1081"
	tag       string = "inbound-1081"
)

func TestAddInbound(t *testing.T) {
	_, err := xray.addInbound(&dto.AddInbound{Port: port, Tag: tag, Protocol: protocoll})
	if err != nil {
		t.Errorf("error has occured please follow error: %v", err)
	}
	xray.removeInbound(&dto.RemoveInbound{Tag: tag})
}

func TestRemoveInbound(t *testing.T) {
	xray.addInbound(&dto.AddInbound{Port: port, Tag: tag, Protocol: protocoll})
	_, err := xray.removeInbound(&dto.RemoveInbound{Tag: tag})
	if err != nil {
		t.Errorf("error has occured please follow error: %v", err)
	}
}

func TestGetTraffic(t *testing.T) {
	xray.addInbound(&dto.AddInbound{Port: port, Tag: tag, Protocol: protocoll})

	ret, err := xray.receiveInboundTraffic(&dto.ReceiveInboundTraffic{Tag: tag})
	if err != nil {
		t.Errorf("error has happend and the error was %v", err)
		xray.removeInbound(&dto.RemoveInbound{Tag: tag})
	}
	if ret.UpLink != 0 || ret.DownLink != 0 {
		t.Error("service gave wrong answer")
		xray.removeInbound(&dto.RemoveInbound{Tag: tag})
	}
}
