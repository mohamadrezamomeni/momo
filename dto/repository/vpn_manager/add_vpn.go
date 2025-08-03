package vpnmanager

import (
	"github.com/mohamadrezamomeni/momo/entity"
)

type AddVPN struct {
	Domain    string
	IsActive  bool
	ApiPort   string
	VPNType   entity.VPNType
	UserCount int
	Country   string
	StartPort int
	EndPort   int
	VPNStatus entity.VPNStatus
}
