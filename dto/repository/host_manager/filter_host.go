package hostmanager

import "github.com/mohamadrezamomeni/momo/entity"

type FilterHosts struct {
	Statuses []entity.HostStatus
	Domains  []string
}
