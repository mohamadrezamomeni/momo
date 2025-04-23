package hostmanager

import "momo/entity"

type AddHost struct {
	Domain         string
	Port           string
	StartRangePort int
	EndRangePort   int
	Status         entity.HostStatus
}
