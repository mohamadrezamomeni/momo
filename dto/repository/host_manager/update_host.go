package hostmanager

import "github.com/mohamadrezamomeni/momo/entity"

type UpdateHost struct {
	Status entity.HostStatus
	Rank   uint32
}
