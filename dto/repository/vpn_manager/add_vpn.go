package dto

import (
	"momo/entity"
)

type AddVPN struct {
	Domain    string
	IsActive  bool
	ApiPort   string
	VPNType   entity.VPNType
	UserCount int
}
