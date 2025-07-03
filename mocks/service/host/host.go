package host

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/entity"
)

type MockHost struct{}

func New() *MockHost {
	return &MockHost{}
}

type Address struct {
	Domain string
	Port   string
}

func (h *MockHost) ResolveHostPortPair(
	hostPortUsed map[string][]string,
	hostPortsRequired map[string]uint32,
) (map[string][]string, error) {
	hostPairs := make([]*Address, 0)
	for domain, countPortsRequired := range hostPortsRequired {
		for i := 0; i < int(countPortsRequired); i++ {
			hostPairs = append(hostPairs, &Address{
				Domain: domain,
				Port:   strconv.Itoa(1000 + i),
			},
			)
		}
	}

	ret := map[string][]string{}
	for _, hostPair := range hostPairs {
		ret[hostPair.Domain] = append(ret[hostPair.Domain], hostPair.Port)
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
