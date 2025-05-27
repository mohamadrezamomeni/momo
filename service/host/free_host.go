package host

import (
	"math/rand"
	"sync"
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (h *Host) ResolveHostPortPair(
	domainPortUsed map[string][]string,
	requiredHosts int,
) ([][2]string, error) {
	scope := "host.service.resolveHostPortPair"
	hosts, err := h.FindRightHosts(entity.High)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	ch := make(chan struct {
		Domain string
		Ports  []string
	})

	var wg sync.WaitGroup
	seen := map[string]struct{}{}

	for _, host := range hosts {
		wg.Add(1)
		if _, ok := seen[host.Domain]; ok {
			continue
		}

		seen[host.Domain] = struct{}{}

		go h.resolvePorts(
			host,
			requiredHosts,
			domainPortUsed[host.Domain],
			&wg,
			ch,
		)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	hostPortPairs := [][2]string{}

	for item := range ch {
		hostPortPairs = append(hostPortPairs, h.makeHostPairWiPort(item.Domain, item.Ports)...)
	}

	return h.shuffleHostPortPairs(hostPortPairs), nil
}

func (h *Host) resolvePorts(
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

func (i *Host) makeHostPairWiPort(host string, ports []string) [][2]string {
	hostPortPairs := [][2]string{}
	for _, port := range ports {
		hostPortPairs = append(hostPortPairs, [2]string{host, port})
	}
	return hostPortPairs
}

func (i *Host) shuffleHostPortPairs(hostPortPairs [][2]string) [][2]string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Shuffle(len(hostPortPairs), func(i, j int) {
		hostPortPairs[i], hostPortPairs[j] = hostPortPairs[j], hostPortPairs[i]
	})
	return hostPortPairs
}
