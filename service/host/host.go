package host

import (
	"sync"

	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
)

type WorkerProxy interface {
	Close()
	GetAvailablePorts(uint32, []string) ([]string, error)
	GetMetric() (uint32, string, error)
}

type WorkerFactor func(string, string) (WorkerProxy, error)

type HostRepo interface {
	Create(*hostRepoDto.AddHost) (*entity.Host, error)
	Filter(*hostRepoDto.FilterHosts) ([]*entity.Host, error)
	Update(int, *hostRepoDto.UpdateHost) error
}

type Host struct {
	hostRepo        HostRepo
	workerFactorNew WorkerFactor
}

func New(hostRepo HostRepo, workerFactorNew WorkerFactor) *Host {
	return &Host{
		hostRepo:        hostRepo,
		workerFactorNew: workerFactorNew,
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

	wp, err := h.workerFactorNew(host.Domain, host.Port)
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
