package inbound

import (
	"reflect"
	"testing"
	"unsafe"

	inboundRepository "momo/mocks/repository/inbound"
	hostService "momo/mocks/service/host"
	userService "momo/mocks/service/user"
	vpnService "momo/mocks/service/vpn"
	"momo/pkg/utils"
)

func registerInboundSvc() (*Inbound, *inboundRepository.MockInbound, *vpnService.MockVPN) {
	inboundRepo := inboundRepository.New()
	userSvc := userService.New()
	hostSvc := hostService.New()
	vpnSvc := vpnService.New()

	inboundSvc := New(inboundRepo, vpnSvc, userSvc, hostSvc)
	return inboundSvc, inboundRepo, vpnSvc
}

func TestApplyDomainAndPortToInbounds(t *testing.T) {
	inboundSvc, inboundRepo, _ := registerInboundSvc()

	inboundCreated1, _ := inboundRepo.Create(inbound1)
	inboundCreated2, _ := inboundRepo.Create(inbound2)
	inboundCreated3, _ := inboundRepo.Create(inbound3)

	inboundSvc.AssignDomainToInbounds()

	ret1, _ := inboundRepo.FindInboundByID(inboundCreated1.ID)
	ret2, _ := inboundRepo.FindInboundByID(inboundCreated2.ID)
	ret3, _ := inboundRepo.FindInboundByID(inboundCreated3.ID)

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
	inboundSvc, inboundRepo, vpnSvc := registerInboundSvc()

	proxy, _ := vpnSvc.MakeProxy()

	ret1, _ := inboundRepo.Create(inbound4)
	ret2, _ := inboundRepo.Create(inbound5)
	ret3, _ := inboundRepo.Create(inbound6)
	inboundSvc.HealingUpInbound(ret1, proxy)

	inboundEnableInpt := utils.ReadPrivateField(proxy, "addInboundData")
	disableInpt := utils.ReadPrivateField(proxy, "disableInboundData")

	if reflect.ValueOf(inboundEnableInpt).IsNil() ||
		!reflect.ValueOf(disableInpt).IsNil() {
		t.Error("the proxy is wrong")
	}
	proxy.Close()

	inboundSvc.HealingUpInbound(ret2, proxy)

	inboundEnableInpt = utils.ReadPrivateField(proxy, "addInboundData")
	disableInpt = utils.ReadPrivateField(proxy, "disableInboundData")

	if !reflect.ValueOf(inboundEnableInpt).IsNil() ||
		reflect.ValueOf(disableInpt).IsNil() {
		t.Error("the proxy is wrong")
	}

	proxy.Close()

	inboundSvc.HealingUpInbound(ret3, proxy)

	inboundEnableInpt = utils.ReadPrivateField(proxy, "addInboundData")
	disableInpt = utils.ReadPrivateField(proxy, "disableInboundData")

	if !reflect.ValueOf(inboundEnableInpt).IsNil() ||
		reflect.ValueOf(disableInpt).IsNil() {
		t.Error("the proxy is wrong")
	}

	proxy.Close()
}
