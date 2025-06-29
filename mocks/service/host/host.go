package host

import (
	"strconv"

	hostService "github.com/mohamadrezamomeni/momo/dto/service/host"
	"github.com/mohamadrezamomeni/momo/entity"
)

type MockHost struct{}

func New() *MockHost {
	return &MockHost{}
}

func (h *MockHost) ResolveHostPortPair(
	hostPortUsed map[string][]string,
	hostPortsRequired map[string]uint32,
) (map[string][]*hostService.HostAddress, error) {
	hostPairs := make([]*hostService.HostAddress, 0)
	for domain, countPortsRequired := range hostPortsRequired {
		for i := 0; i < int(countPortsRequired); i++ {
			hostPairs = append(hostPairs, &hostService.HostAddress{
				Domain: domain,
				Port:   strconv.Itoa(1000 + i),
			},
			)
		}
	}

	ret := map[string][]*hostService.HostAddress{}
	for _, hostPair := range hostPairs {
		ret[hostPair.Domain] = append(ret[hostPair.Domain], hostPair)
	}

	return ret, nil
}

func (h *MockHost) FindRightHosts(status entity.HostStatus) ([]*entity.Host, error) {
	return []*entity.Host{
		{
			ID:     1,
			Rank:   10,
			Status: status,
			Port:   "666",
			Domain: "google.com",
		},
		{
			ID:     2,
			Rank:   10,
			Status: status,
			Port:   "666",
			Domain: "twitter.com",
		},
	}, nil
}
