package host

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/entity"
)

type MockHost struct{}

func New() *MockHost {
	return &MockHost{}
}

func (h *MockHost) ResolveHostPortPair(hostPortUsed map[string][]string, requiredPorts int) ([][2]string, error) {
	hostPairs := make([][2]string, 0)
	for i := 0; i < requiredPorts; i++ {
		hostPairs = append(hostPairs, [2]string{"google.com", strconv.Itoa(1000 + i)})
	}

	return hostPairs, nil
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
