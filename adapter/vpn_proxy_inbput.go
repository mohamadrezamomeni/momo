package adapter

import (
	vpnProxyDto "github.com/mohamadrezamomeni/momo/dto/proxy/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
)

func GenerateVPNProxyInput(inbound *entity.Inbound, user *entity.User) *vpnProxyDto.Inbound {
	return &vpnProxyDto.Inbound{
		User: &vpnProxyDto.User{
			Username: user.Username,
			ID:       user.ID,
			Level:    "0",
		},
		VPNType: inbound.VPNType,
		Address: inbound.Domain,
		Port:    inbound.Port,
		Tag:     inbound.Tag,
	}
}
