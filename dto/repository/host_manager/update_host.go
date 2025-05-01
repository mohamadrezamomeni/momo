package hostmanager

import "momo/entity"

type UpdateHost struct {
	Status entity.HostStatus
	Rank   uint32
}
