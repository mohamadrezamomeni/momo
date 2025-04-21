package xray

import (
	"testing"

	"momo/proxy/vpn/internal/xray/dto"
)

var xrayU *Xray

func TestMain(m *testing.M) {
	xrayU, _ = New(&XrayConfig{
		Address: "192.168.116.129",
		ApiPort: "62789",
	})
}

var (
	protocolUser string = "vmess"
	portUser     string = "1082"
	tagUser      string = "inbound-1082"

	username string = "mohamadian"
	level    string = "0"
	uuid     string = "0393ed06-29bb-41c2-b3f4-6382a6729c3e"
)

func TestAddUser(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{Tag: tagUser, Port: portUser, Protocol: protocolUser})
	err := xrayU.addUser(&dto.AddUser{Tag: tagUser, Level: level, Username: username, UUID: uuid})
	if err != nil {
		t.Errorf("error has happend the error was %v", err)
	}
	xrayU.removeInbound(&dto.RemoveInbound{Tag: tagUser})
}

func TestRemoveUser(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{Tag: tagUser, Port: portUser, Protocol: protocolUser})
	xrayU.addUser(&dto.AddUser{Tag: tagUser, Level: level, Username: username, UUID: uuid})
	err := xrayU.removeUser(&dto.RemoveUser{Tag: tagUser, Username: username})
	if err != nil {
		t.Errorf("error has happend the error was %v", err)
	}
	xrayU.removeInbound(&dto.RemoveInbound{Tag: tagUser})
}

func TestGetUsers(t *testing.T) {
	xrayU.addInbound(&dto.AddInbound{
		Tag:      tagUser,
		Port:     portUser,
		Protocol: protocolUser,
		User: &dto.InboundUser{
			Username: username,
			Level:    level,
			UUID:     uuid,
		},
	})

	res, err := xrayU.getUsers(tagUser)
	if err != nil {
		t.Errorf("error has happend you can follow the problem the problem was %v", err)
	}

	if len(res.Usernames) != 1 {
		t.Error("the inbound was empty It had been expected this inbound had one user")
	}
	xrayU.removeInbound(&dto.RemoveInbound{
		Tag: tagUser,
	})
}
