package host

import (
	"sync"

	"github.com/mohamadrezamomeni/momo/entity"
)

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
