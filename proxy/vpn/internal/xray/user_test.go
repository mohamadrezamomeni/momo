package xray

import (
	"testing"

	"momo/proxy/vpn/internal/xray/dto"
)

var xrayU *Xray = New(&XrayConfig{
	Address: "192.168.116.129",
	ApiPort: "62789",
})

var (
	protocolUser string = "vmess"
	portUser     string = "1082"
	tagUser      string = "inbound-1082"

	email string = "mohamadian@gmail.com"
	level string = "0"
	uuid  string = "0393ed06-29bb-41c2-b3f4-6382a6729c3e"
)

func TestAddUser(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{Tag: tagUser, Port: portUser, Protocol: protocolUser})
	err := xrayU.addUser(&dto.AddUser{Tag: tagUser, Level: level, Email: email, UUID: uuid})
	if err != nil {
		t.Errorf("error has happend the error was %v", err)
	}
	xrayU.removeInbound(&dto.RemoveInbound{Tag: tagUser})
}

func TestRemoveUser(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{Tag: tagUser, Port: portUser, Protocol: protocolUser})
	xrayU.addUser(&dto.AddUser{Tag: tagUser, Level: level, Email: email, UUID: uuid})
	err := xrayU.removeUser(&dto.RemoveUser{Tag: tagUser, Email: email})
	if err != nil {
		t.Errorf("error has happend the error was %v", err)
	}
	xrayU.removeInbound(&dto.RemoveInbound{Tag: tagUser})
}
