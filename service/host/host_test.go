package host

import (
	"sync"
	"testing"

	"momo/entity"
	workerMock "momo/mocks/proxy/worker"
	hostRepository "momo/mocks/repository/host"
)

func registerHostSrv() (*Host, *hostRepository.Host) {
	hostRepo := hostRepository.New()
	return New(hostRepo, func(address string, port string) (WorkerProxy, error) {
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

	inboundCreated1, _ := hostRepo.Create(host1)
	ch := make(chan struct {
		Domain string
		Ports  []string
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go hostSvc.ResolvePorts(
		inboundCreated1,
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

	if data.Domain != inboundCreated1.Domain || len(data.Ports) != 3 {
		t.Fatalf("error has happend the date that was sent was wrong")
	}
}
