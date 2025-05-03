package host

import (
	hostRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/host_manager"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (h *Host) MonitorHosts() {
	hosts, _ := h.hostRepo.Filter(&hostRepoDto.FilterHosts{})
	for _, host := range hosts {
		rank, status, _ := h.monitorHost(host)
		h.hostRepo.Update(host.ID, &hostRepoDto.UpdateHost{
			Rank:   rank,
			Status: status,
		})
	}
}

func (h *Host) monitorHost(host *entity.Host) (uint32, entity.HostStatus, error) {
	scope := "hostService.metric.monitorHost"

	wp, err := h.workerFactorNew(host.Domain, host.Port)
	if err != nil {
		return 0, entity.Uknown, momoError.Wrap(err).Scope(scope).Errorf("error to connect \"%s:%s\"", host.Domain, host.Port)
	}
	rank, status, err := wp.GetMetric()
	if err != nil {
		return 0, entity.Uknown, momoError.Wrap(err).Scope(scope).Errorf("error to connect \"%s:%s\"", host.Domain, host.Port)
	}

	return rank, status, nil
}
