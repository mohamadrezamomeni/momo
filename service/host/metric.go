package host

import (
	hostRepoDto "momo/dto/repository/host_manager"
	"momo/entity"
	momoError "momo/pkg/error"
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
	wp, err := h.workerFactorNew(host.Domain, host.Port)
	if err != nil {
		return 0, entity.Uknown, momoError.Errorf("error to connect \"%s:%s\" and error was %v", host.Domain, host.Port, err)
	}
	rank, status, err := wp.GetMetric()
	if err != nil {
		return 0, entity.Uknown, momoError.Errorf("error to connect \"%s:%s\" and error was %v", host.Domain, host.Port, err)
	}

	return rank, status, nil
}
