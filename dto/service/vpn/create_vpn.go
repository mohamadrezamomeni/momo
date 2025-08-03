package vpn

import "github.com/mohamadrezamomeni/momo/entity"

type CreateVPN struct {
	VpnType   entity.VPNType
	Port      string
	Domain    string
	UserCount int
	Country   string
	StartPort int
	EndPort   int
	VPNStatus entity.VPNStatus
}
