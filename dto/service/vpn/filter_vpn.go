package vpn

import "github.com/mohamadrezamomeni/momo/entity"

type FilterVPNs struct {
	Domain  string
	VPNType entity.VPNType
}
