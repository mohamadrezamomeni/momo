package inbound

import (
	"testing"

	inboundRepository "momo/mocks/repository/inbound"
	hostService "momo/mocks/service/host"
	userService "momo/mocks/service/user"
	vpnService "momo/mocks/service/vpn"
)

func registerInboundSvc() (*Inbound, *inboundRepository.MockInbound) {
	inboundRepo := inboundRepository.New()
	userSvc := userService.New()
	hostSvc := hostService.New()
	vpnSvc := vpnService.New()

	inboundSvc := New(inboundRepo, vpnSvc, userSvc, hostSvc)
	return inboundSvc, inboundRepo
}

func TestApplyDomainAndPortToInbounds(t *testing.T) {
	inboundSvc, inboundRepo := registerInboundSvc()

	inboundCreated1, _ := inboundRepo.Create(inbound1)
	inboundCreated2, _ := inboundRepo.Create(inbound2)
	inboundCreated3, _ := inboundRepo.Create(inbound3)

	inboundSvc.AssignDomainToInbounds()

	inbound1, _ := inboundRepo.FindInboundByID(inboundCreated1.ID)
	inbound2, _ := inboundRepo.FindInboundByID(inboundCreated2.ID)
	inbound3, _ := inboundRepo.FindInboundByID(inboundCreated3.ID)

	if inbound1.IsAssigned != true || inbound2.IsAssigned != true {
		t.Fatal("inbounds aren't updated currectly")
	}

	if inbound3.Domain != "instagram.com" {
		t.Fatal("wrong inbound is updated")
	}
}
