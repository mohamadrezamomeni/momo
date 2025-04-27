package host

import (
	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
	momoError "momo/pkg/error"
	utils "momo/pkg/utils"
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

func (h *Host) FindRightHost(status entity.HostStatus) (string, string, error) {
	host, err := h.findAppropriateHostByStatus([]entity.HostStatus{status})
	if err != nil {
		return "", "", err
	}

	return host.Domain, host.Port, nil
}

func (h *Host) findAppropriateHostByStatus(statuses []entity.HostStatus) (*entity.Host, error) {
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
		return highHosts[utils.GetRandom(0, len(highHosts))], nil
	}

	if len(mediumHosts) != 0 {
		return highHosts[utils.GetRandom(0, len(mediumHosts))], nil
	}

	if len(lowHosts) != 0 {
		return highHosts[utils.GetRandom(0, len(lowHosts))], nil
	}
	return nil, momoError.Error("appropriate server isn't selected")
}
