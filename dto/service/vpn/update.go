package vpn

import "github.com/mohamadrezamomeni/momo/entity"

type Update struct {
	VPNStatus entity.VPNStatus
	Domain    string
	ApiPort   string
}
