package inbound

import "github.com/mohamadrezamomeni/momo/entity"

type FilterInbounds struct {
	Domain  string
	Port    string
	VPNType entity.VPNType
	UserID  string
}
