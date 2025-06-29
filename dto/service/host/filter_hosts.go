package host

import "github.com/mohamadrezamomeni/momo/entity"

type FilterHosts struct {
	Status  []entity.HostStatus
	Domains []string
}
