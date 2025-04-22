package dto

import (
	"momo/entity"
)

type Add_VPN struct {
	Domain         string
	IsActive       bool
	ApiPort        string
	StartRangePort int
	EndRangePort   int
	VPNType        entity.VPNType
	UserCount      int
}
