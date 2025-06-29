package host

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	hostRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/host_manager"
	hostServiceDto "github.com/mohamadrezamomeni/momo/dto/service/host"
	"github.com/mohamadrezamomeni/momo/entity"
)

type WorkerFactor func(string, string) (adapter.WorkerProxy, error)

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

func (h *Host) Filter(filterDto *hostServiceDto.FilterHosts) ([]*entity.Host, error) {
	return h.hostRepo.Filter(
		&hostRepoDto.FilterHosts{
			Statuses: filterDto.Status,
			Domains:  filterDto.Domains,
		},
	)
}

func (h *Host) Create(createHostDto *hostServiceDto.CreateHostDto) (*entity.Host, error) {
	return h.hostRepo.Create(&hostRepoDto.AddHost{
		Domain: createHostDto.Domain,
		Port:   createHostDto.Port,
		Status: entity.Deactive,
		Rank:   0,
	})
}
