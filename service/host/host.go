package host

import (
	"sync"

	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
	"momo/proxy/worker"
)

type HostRepo interface {
	Create(*hostRepoDto.AddHost) (*entity.Host, error)
	Filter(*hostRepoDto.FilterHosts) ([]*entity.Host, error)
	Update(int, *hostRepoDto.UpdateHost) error
}

type Host struct {
	hostRepo HostRepo
}

func New(hostRepo HostRepo) *Host {
	return &Host{
		hostRepo: hostRepo,
	}
}

func (h *Host) FindRightHosts(status entity.HostStatus) ([]*entity.Host, error) {
	hosts, err := h.hostRepo.Filter(
		&hostRepoDto.FilterHosts{
			Statuses: []entity.HostStatus{status},
		},
	)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (h *Host) ResolvePorts(
	host *entity.Host,
	requiredPorts int,
	portsUsed []string,
	wg *sync.WaitGroup,
	ch chan<- struct {
		Domain string
		Ports  []string
	},
) {
	defer wg.Done()

	wp, err := worker.New(&worker.Config{
		Address: host.Domain,
		Port:    host.Port,
	})
	if err != nil {
		return
	}

	ports, err := wp.GetAvailablePorts(uint32(requiredPorts), portsUsed)
	if err != nil {
		return
	}

	ch <- struct {
		Domain string
		Ports  []string
	}{
		Domain: host.Domain,
		Ports:  ports,
	}
}
