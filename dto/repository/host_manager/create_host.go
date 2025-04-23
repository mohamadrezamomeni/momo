package hostmanager

import "momo/entity"

type AddHost struct {
	Domain string
	Port   string
	Status entity.HostStatus
}
