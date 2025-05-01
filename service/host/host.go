package host

import (
	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
)

type WorkerProxy interface {
	Close()
	GetAvailablePorts(uint32, []string) ([]string, error)
	GetMetric() (uint32, entity.HostStatus, error)
}

type WorkerFactor func(string, string) (WorkerProxy, error)

type HostRepo interface {
	Create(*hostRepoDto.AddHost) (*entity.Host, error)
	Filter(*hostRepoDto.FilterHosts) ([]*entity.Host, error)
	Update(int, *hostRepoDto.UpdateHost) error
	FindByID(int) (*entity.Host, error)
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
