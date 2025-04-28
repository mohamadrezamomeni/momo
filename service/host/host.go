package host

import (
	"sync"

	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
	momoError "momo/pkg/error"
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
	hosts, err := h.findAppropriateHostByStatus([]entity.HostStatus{status})
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
		domain string
		ports  []string
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
		domain string
		ports  []string
	}{
		domain: host.Domain,
		ports:  ports,
	}
}

func (h *Host) findAppropriateHostByStatus(statuses []entity.HostStatus) ([]*entity.Host, error) {
	hosts, err := h.hostRepo.Filter(
		&hostRepoDto.FilterHosts{
			Statuses: statuses,
		},
	)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	highHosts := []*entity.Host{}
	mediumHosts := []*entity.Host{}
	lowHosts := []*entity.Host{}

	for _, host := range hosts {

		if host.Status == entity.High {
			highHosts = append(highHosts, host)
		}
		if host.Status == entity.Medium {
			mediumHosts = append(mediumHosts, host)
		}

		if host.Status == entity.Low {
			lowHosts = append(lowHosts, host)
		}
	}

	if len(highHosts) != 0 {
		return highHosts, nil
	}

	if len(mediumHosts) != 0 {
		return mediumHosts, nil
	}

	if len(lowHosts) != 0 {
		return lowHosts, nil
	}
	return nil, momoError.Error("appropriate server isn't selected")
}
