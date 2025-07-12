package inbound

import (
	"reflect"
	"strconv"
	"testing"
	"unsafe"

	"github.com/mohamadrezamomeni/momo/dto/service/vpn"
	inboundRepository "github.com/mohamadrezamomeni/momo/mocks/repository/inbound"
	chargeService "github.com/mohamadrezamomeni/momo/mocks/service/charge"
	hostService "github.com/mohamadrezamomeni/momo/mocks/service/host"
	inboundChargeService "github.com/mohamadrezamomeni/momo/mocks/service/inbound_charge"
	userService "github.com/mohamadrezamomeni/momo/mocks/service/user"
	vpnService "github.com/mohamadrezamomeni/momo/mocks/service/vpn"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func registerInboundSvc() (*Inbound, *HostInbound, *HealingUpInbound, *inboundRepository.MockInbound, *vpnService.MockVPN) {
	inboundRepo := inboundRepository.New()
	hostSvc := hostService.New()
	vpnSvc := vpnService.New()

	inboundSvc := New(inboundRepo)
	inboundChargeSvc := inboundChargeService.New()
	inboundHostSvc := NewHostInbound(inboundRepo, hostSvc, vpnSvc)
	userSvc := userService.New()
	chargeSvc := chargeService.New()
	healingUpInboundSvc := NewHealingUpInbound(inboundRepo, vpnSvc, inboundChargeSvc, chargeSvc, userSvc)
	return inboundSvc, inboundHostSvc, healingUpInboundSvc, inboundRepo, vpnSvc
}

func TestApplyDomainAndPortToInbounds(t *testing.T) {
	_, inboundHostSvc, _, inboundRepo, vpnSvc := registerInboundSvc()

	defer vpnSvc.DeleteAll()

	vpnSvc.Create(&vpn.CreateVPN{
		VpnType:   inbound1.VPNType,
		Country:   inbound1.Country,
		Domain:    "instagram.com",
		UserCount: 2,
		Port:      "203",
	})

	vpnSvc.Create(&vpn.CreateVPN{
		VpnType:   inbound3.VPNType,
		Country:   inbound3.Country,
		Domain:    "facebook.com",
		UserCount: 2,
		Port:      "203",
	})
	inboundCreated1, _ := inboundRepo.Create(inbound1)
	inboundCreated2, _ := inboundRepo.Create(inbound2)
	inboundCreated3, _ := inboundRepo.Create(inbound3)
	inboundCreated4, _ := inboundRepo.Create(inbound7)

	inboundHostSvc.AssignDomainToInbounds()

	ret1, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated1.ID))
	ret2, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated2.ID))
	ret3, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated3.ID))
	ret4, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated4.ID))

	if ret1.IsAssigned != true ||
		ret2.IsAssigned != true ||
		ret1.Domain != "instagram.com" ||
		ret2.Domain != "instagram.com" {
		t.Fatal("inbounds aren't updated currectly")
	}

	if ret4.IsAssigned != false {
		t.Fatal("the inbound whose country isn't exist is updated")
	}

	if ret3.Domain != "instagram.com" {
		t.Fatal("wrong inbound is updated")
	}
}

func ReadPrivateField(obj interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName(fieldName)
	f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	return f.Interface()
}

func TestHealingUpInbounds(t *testing.T) {
	_, _, healingUpInboundSvc, inboundRepo, vpnSvc := registerInboundSvc()

	proxy, _ := vpnSvc.MakeProxy()

	ret1, _ := inboundRepo.Create(inbound4)
	ret2, _ := inboundRepo.Create(inbound5)
	ret3, _ := inboundRepo.Create(inbound6)
	healingUpInboundSvc.activeInbound(ret1, proxy)

	inboundEnableInpt := utils.ReadPrivateField(proxy, "addInboundData")
	disableInpt := utils.ReadPrivateField(proxy, "disableInboundData")

	if reflect.ValueOf(inboundEnableInpt).IsNil() ||
		!reflect.ValueOf(disableInpt).IsNil() {
		t.Error("the proxy is wrong")
	}
	proxy.Close()

	healingUpInboundSvc.deActiveInbound(ret2, proxy)

	inboundEnableInpt = utils.ReadPrivateField(proxy, "addInboundData")
	disableInpt = utils.ReadPrivateField(proxy, "disableInboundData")

	if !reflect.ValueOf(inboundEnableInpt).IsNil() ||
		reflect.ValueOf(disableInpt).IsNil() {
		t.Error("the proxy is wrong")
	}

	proxy.Close()

	healingUpInboundSvc.deActiveInbound(ret3, proxy)

	inboundEnableInpt = utils.ReadPrivateField(proxy, "addInboundData")
	disableInpt = utils.ReadPrivateField(proxy, "disableInboundData")

	if !reflect.ValueOf(inboundEnableInpt).IsNil() ||
		reflect.ValueOf(disableInpt).IsNil() {
		t.Error("the proxy is wrong")
	}

	proxy.Close()
}

func TestOpenInboundsPortMustBeOpen(t *testing.T) {
	_, hostInboundSvc, _, inboundRepo, _ := registerInboundSvc()
	inboundCreated1, _ := inboundRepo.Create(inbound8)
	inboundCreated2, _ := inboundRepo.Create(inbound9)
	inboundCreated3, _ := inboundRepo.Create(inbound10)
	hostInboundSvc.OpenInboundsPortMustBeOpen()

	inbound1, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated1.ID))
	inbound2, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated2.ID))
	inbound3, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated3.ID))
	if inbound1.IsPortOpen != false {
		t.Fatalf("error to compare data we expected inbound1 be false")
	}

	if inbound2.IsPortOpen != true {
		t.Fatalf("error to compare data we expected inbound1 be true")
	}

	if inbound3.IsPortOpen != true {
		t.Fatalf("error to compare data we expected inbound1 be true")
	}
}
