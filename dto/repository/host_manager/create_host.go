package hostmanager

import "github.com/mohamadrezamomeni/momo/entity"

type AddHost struct {
	Domain string
	Port   string
	Rank   uint32
	Status entity.HostStatus
}
