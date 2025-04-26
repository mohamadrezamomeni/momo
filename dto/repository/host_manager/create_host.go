package hostmanager

import "momo/entity"

type AddHost struct {
	Domain string
	Port   string
	Rank   uint32
	Status entity.HostStatus
}
