package xray

import (
	"testing"

	"momo/proxy/vpn/internal/xray/dto"
)

var xrayInbound *Xray = New(&XrayConfig{
	Address: "192.168.116.129",
	ApiPort: "62789",
})

var (
	protocolInbound string = "vmess"
	portInbound     string = "1081"
	tagInbound      string = "inbound-1081"

	emailInbound string = "mohamadian@gmail.com"
	levelInbound string = "0"
	uuidInbound  string = "0393ed06-29bb-41c2-b3f4-6382a6729c3e"

	inboundDoesntExist string = "inbound-1083"
)

func TestAddInbound(t *testing.T) {
	_, err := xrayInbound.addInbound(&dto.AddInbound{
		Port:     portInbound,
		Tag:      tagInbound,
		Protocol: protocolInbound,
		User: &dto.InboundUser{
			Email: emailInbound,
			Level: level,
			UUID:  uuidInbound,
		},
	})
	if err != nil {
		t.Errorf("error has occured please follow error: %v", err)
	}
	xrayInbound.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
}

func TestRemoveInbound(t *testing.T) {
	xrayInbound.addInbound(&dto.AddInbound{
		Port:     portInbound,
		Tag:      tagInbound,
		Protocol: protocolInbound,
		User: &dto.InboundUser{
			Email: emailInbound,
			Level: level,
			UUID:  uuidInbound,
		},
	})
	_, err := xrayInbound.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
	if err != nil {
		t.Errorf("error has occured please follow error: %v", err)
	}
}

func TestGetTraffic(t *testing.T) {
	xrayInbound.addInbound(&dto.AddInbound{Port: portInbound, Tag: tagInbound, Protocol: protocolInbound})

	ret, err := xrayInbound.receiveInboundTraffic(&dto.ReceiveInboundTraffic{Tag: tagInbound})
	if err != nil {
		xrayInbound.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
		t.Errorf("error has happend and the error was %v", err)
	}
	if ret.UpLink != 0 || ret.DownLink != 0 {
		xrayInbound.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
		t.Error("service gave wrong answer")
	}
	xrayInbound.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
}

func TestInboundDoesntExist(t *testing.T) {
	res, err := xrayInbound.getUsers(inboundDoesntExist)

	if err == nil {
		t.Error("error could be existed. It was unexpected situation")
	}
	if len(res.Emails) != 0 {
		t.Error("the number of emails must be 0")
	}
}
