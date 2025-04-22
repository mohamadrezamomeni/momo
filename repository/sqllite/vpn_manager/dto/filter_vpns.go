package dto

import "momo/entity"

type FilterVPNs struct {
	IsActive *bool
	Domain   string
	VPNType  entity.VPNType
}
