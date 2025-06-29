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

	vpnSvc.Create(&vpn.CreateVPN{
		VpnType:   inbound1.VPNType,
		Country:   inbound1.Country,
		Domain:    "instagram.com",
		UserCount: 2,
		Port:      "203",
	})
	vpnSvc.Create(&vpn.CreateVPN{
		VpnType:   inbound2.VPNType,
		Country:   inbound2.Country,
		Domain:    "instagram.com",
		UserCount: 2,
		Port:      "203",
	})
	vpnSvc.Create(&vpn.CreateVPN{
		VpnType:   inbound3.VPNType,
		Country:   inbound3.Country,
		Domain:    "instagram.com",
		UserCount: 2,
		Port:      "203",
	})
	inboundCreated1, _ := inboundRepo.Create(inbound1)
	inboundCreated2, _ := inboundRepo.Create(inbound2)
	inboundCreated3, _ := inboundRepo.Create(inbound3)

	inboundHostSvc.AssignDomainToInbounds()

	ret1, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated1.ID))
	ret2, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated2.ID))
	ret3, _ := inboundRepo.FindInboundByID(strconv.Itoa(inboundCreated3.ID))

	if ret1.IsAssigned != true || ret2.IsAssigned != true {
		t.Fatal("inbounds aren't updated currectly")
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
