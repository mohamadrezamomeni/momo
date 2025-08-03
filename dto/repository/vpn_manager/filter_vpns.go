package dto

import "github.com/mohamadrezamomeni/momo/entity"

type FilterVPNs struct {
	IsActive    *bool
	Domain      string
	VPNTypes    []entity.VPNType
	Coountries  []string
	VPNStatuses []entity.VPNStatus
}
