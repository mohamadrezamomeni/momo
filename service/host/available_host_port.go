package host

import (
	"math/rand"
	"sync"
	"time"

	hostServiceDto "github.com/mohamadrezamomeni/momo/dto/service/host"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type Address struct {
	Domain string
	Port   string
}

func (h *Host) ResolveHostPortPair(
	domainPortUsed map[string][]string,
	requiredHostPorts map[string]uint32,
) (map[string][]string, error) {
	scope := "host.service.resolveHostPortPair"

	domains := h.extractRequiredHostPorts(requiredHostPorts)
	hosts, err := h.Filter(&hostServiceDto.FilterHosts{
		Status:  []entity.HostStatus{entity.High},
		Domains: domains,
	})
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	ch := make(chan struct {
		Domain string
		Ports  []string
	})

	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)

		go h.resolvePorts(
			host,
			requiredHostPorts[host.Domain],
			domainPortUsed[host.Domain],
			&wg,
			ch,
		)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	hostPortPairs := []*Address{}

	for item := range ch {
		hostPortPairs = append(hostPortPairs, h.makeHostPairWiPort(item.Domain, item.Ports)...)
	}

	return h.shuffleHostPortPairs(hostPortPairs), nil
}

func (h *Host) resolvePorts(
	host *entity.Host,
	requiredPorts uint32,
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

func (h *Host) makeHostPairWiPort(host string, ports []string) []*Address {
	hostPortPairs := []*Address{}
	for _, port := range ports {
		hostPortPairs = append(hostPortPairs, &Address{
			Domain: host,
			Port:   port,
		})
	}
	return hostPortPairs
}

func (h *Host) shuffleHostPortPairs(hostPortPairs []*Address) map[string][]string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Shuffle(len(hostPortPairs), func(i, j int) {
		hostPortPairs[i], hostPortPairs[j] = hostPortPairs[j], hostPortPairs[i]
	})

	ret := map[string][]string{}
	for _, hostPair := range hostPortPairs {
		ret[hostPair.Domain] = append(ret[hostPair.Domain], hostPair.Port)
	}
	return ret
}

func (i *Host) extractRequiredHostPorts(requiredHostPorts map[string]uint32) []string {
	domains := []string{}
	for domain := range requiredHostPorts {
		domains = append(domains, domain)
	}
	return domains
}
