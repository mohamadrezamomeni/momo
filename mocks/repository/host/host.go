package host

import (
	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
)

type Host struct {
	hosts []*entity.Host
	id    int
}

func New() *Host {
	return &Host{
		hosts: make([]*entity.Host, 0),
		id:    0,
	}
}

func (h *Host) Create(inpt *hostRepoDto.AddHost) (*entity.Host, error) {
	host := &entity.Host{
		Domain: inpt.Domain,
		Port:   inpt.Port,
		Status: inpt.Status,
		Rank:   inpt.Rank,
	}
	h.hosts = append(h.hosts, host)
	h.id += 1
	return host, nil
}

func (h *Host) FindRightHosts(status entity.HostStatus) ([]*entity.Host, error) {
	hosts := make([]*entity.Host, 0)

	for _, host := range hosts {
		if host.Status == status {
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}

func (h *Host) Update(id int, inpt *hostRepoDto.UpdateHost) error {
	for _, host := range h.hosts {
		if host.ID == id {
			host.Rank = uint32(inpt.Rank)
			host.Status = inpt.Status
		}
	}
	return nil
}
