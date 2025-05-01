package host

import (
	"testing"

	"momo/entity"
	hostRepository "momo/mocks/repository/host"
)

func registerHostSrv() (*Host, *hostRepository.Host) {
	hostRepo := hostRepository.New()
	return New(hostRepo), hostRepo
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
