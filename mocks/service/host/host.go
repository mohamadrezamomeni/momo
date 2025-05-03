package host

import (
	"sync"

	"github.com/mohamadrezamomeni/momo/entity"
)

type MockHost struct{}

func New() *MockHost {
	return &MockHost{}
}

func (h *MockHost) ResolvePorts(
	host *entity.Host,
	numberPortNedded int,
	ports []string,
	wg *sync.WaitGroup,
	ch chan<- struct {
		Domain string
		Ports  []string
	},
) {
	defer wg.Done()

	ch <- struct {
		Domain string
		Ports  []string
	}{
		Domain: host.Domain,
		Ports:  []string{"1234", "3456"},
	}
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
