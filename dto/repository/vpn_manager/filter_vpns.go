package dto

import "github.com/mohamadrezamomeni/momo/entity"

type FilterVPNs struct {
	IsActive *bool
	Domain   string
	VPNType  entity.VPNType
}
