package xray

import (
	"testing"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/proxy/vpn/internal/xray/dto"
)

var (
	protocolInbound string = "vmess"
	portInbound     string = "1081"
	tagInbound      string = "inbound-1081"

	usernameInbound string = "mohamadian"
	levelInbound    string = "0"
	uuidInbound     string = "0393ed06-29bb-41c2-b3f4-6382a6729c3e"

	inboundDoesntExist string = "inbound-1083"
)

func TestAddInbound(t *testing.T) {
	_, err := xrayU.addInbound(&dto.AddInbound{
		Port:     portInbound,
		Tag:      tagInbound,
		Protocol: protocolInbound,
		User: &dto.InboundUser{
			Username: usernameInbound,
			Level:    level,
			UUID:     uuidInbound,
		},
	})
	if err != nil {
		t.Errorf("error has occured please follow error: %v", err)
	}
	xrayU.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
}

func TestRemoveInbound(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{
		Port:     portInbound,
		Tag:      tagInbound,
		Protocol: protocolInbound,
		User: &dto.InboundUser{
			Username: usernameInbound,
			Level:    level,
			UUID:     uuidInbound,
		},
	})
	_, err := xrayU.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
	if err != nil {
		t.Errorf("error has occured please follow error: %v", err)
	}
}

func TestGetTraffic(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{Port: portInbound, Tag: tagInbound, Protocol: protocolInbound})

	ret, err := xrayU.getInboundTrafficWithoutBeigReseted(tagInbound)
	if err != nil {
		xrayU.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
		t.Errorf("error has happend and the error was %v", err)
	}
	if ret.UpLink != 0 || ret.DownLink != 0 {
		xrayU.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
		t.Error("service gave wrong answer")
	}
	xrayU.removeInbound(&dto.RemoveInbound{Tag: tagInbound})
}

func TestInboundDoesntExist(t *testing.T) {
	_, err := xrayU.getUsers(inboundDoesntExist)
	if momoErr, ok := err.(*momoError.MomoError); !ok || momoErr.GetErrorType() != momoError.NotFound {
		t.Fatal("error could be existed. It was unexpected situation")
	}

	xrayU.addInbound(&dto.AddInbound{
		Port:     portInbound,
		Tag:      inboundDoesntExist,
		Protocol: protocolInbound,
		User: &dto.InboundUser{
			Username: usernameInbound,
			Level:    level,
			UUID:     uuidInbound,
		},
	})

	data, err := xrayU.getUsers(inboundDoesntExist)

	if err != nil || len(data.Usernames) == 0 {
		t.Fatal("error could be existed. It was unexpected situation")
	}

	xrayU.removeInbound(&dto.RemoveInbound{Tag: inboundDoesntExist})
	_, err = xrayU.getUsers(inboundDoesntExist)
	if momoErr, ok := err.(*momoError.MomoError); !ok || momoErr.GetErrorType() != momoError.NotFound {
		t.Fatal("error could be existed. It was unexpected situation")
	}
}

func TestFakeRequest(t *testing.T) {
	err := xrayU.fakeReceiveInboundTraffic()
	if err != nil {
		t.Errorf("error has happend that was %v", err)
	}
}
