package host

import (
	"sync"
	"testing"

	"momo/adapter"
	"momo/entity"
	workerMock "momo/mocks/proxy/worker"
	hostRepository "momo/mocks/repository/host"
)

func registerHostSrv() (*Host, *hostRepository.Host) {
	hostRepo := hostRepository.New()
	return New(hostRepo, func(address string, port string) (adapter.WorkerProxy, error) {
		return &workerMock.MockWorkerProxy{}, nil
	}), hostRepo
}

func TestFindRightHosts(t *testing.T) {
	hostSvc, hostRepo := registerHostSrv()
	hostRepo.Create(host1)
	hostRepo.Create(host2)
	hostRepo.Create(host3)

	hosts, err := hostSvc.FindRightHosts(entity.High)
	if err != nil {
		t.Fatalf("error has occured that was %v", err)
	}
	if len(hosts) != 3 {
		t.Fatalf("we expected we got 3 items but we got %d", len(hosts))
	}
	for _, host := range hosts {
		if host.Status != entity.High {
			t.Fatal("we got unexpeted record")
		}
	}
}

func TestResolvePorts(t *testing.T) {
	hostSvc, hostRepo := registerHostSrv()

	hostCreated1, _ := hostRepo.Create(host1)
	ch := make(chan struct {
		Domain string
		Ports  []string
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go hostSvc.ResolvePorts(
		hostCreated1,
		3,
		[]string{"12345"},
		&wg,
		ch,
	)

	go func() {
		wg.Wait()
		close(ch)
	}()

	data := <-ch

	if data.Domain != hostCreated1.Domain || len(data.Ports) != 3 {
		t.Fatalf("error has happend the date that was sent was wrong")
	}
}

func TestMonitorHosts(t *testing.T) {
	hostSvc, hostRepo := registerHostSrv()

	hostCreated1, _ := hostRepo.Create(host4)
	hostCreated2, _ := hostRepo.Create(host5)
	hostCreated3, _ := hostRepo.Create(host6)
	hostSvc.MonitorHosts()

	hostFound1, _ := hostRepo.FindByID(hostCreated1.ID)
	hostFound2, _ := hostRepo.FindByID(hostCreated2.ID)
	hostFound3, _ := hostRepo.FindByID(hostCreated3.ID)
	hosts := []*entity.Host{hostFound1, hostFound2, hostFound3}
	for _, host := range hosts {
		if host.Rank != 10 || host.Status != entity.High {
			t.Fatalf("host hasn't updated")
		}
	}
}
