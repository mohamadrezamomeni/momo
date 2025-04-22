package dto

import "momo/entity"

type FilterInbound struct {
	Protocol string
	IsActice *bool
	Domain   string
	VPNType  entity.VPNType
	Port     string
	UserID   string
}
