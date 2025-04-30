package dto

import "momo/entity"

type Inbound struct {
	User     *User
	Port     string
	Protocol string
	Address  string
	Tag      string
	VPNType  entity.VPNType
}
