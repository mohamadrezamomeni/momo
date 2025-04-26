package dto

import "momo/entity"

type FilterInbound struct {
	Protocol string
	IsActive *bool
	Domain   string
	VPNType  entity.VPNType
	Port     string
	UserID   string
}
