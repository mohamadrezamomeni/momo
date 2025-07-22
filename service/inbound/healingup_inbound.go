package inbound

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/entity"
	vpnProxy "github.com/mohamadrezamomeni/momo/proxy/vpn"
)

type HealingUpInbound struct {
	inboundRepo      HealingUpInboundRepo
	vpnService       VpnService
	inboundChargeSvc InboundChargeService
	chargeSvc        ChargeService
	userService      UserService
}

type HealingUpInboundRepo interface {
	RetriveActiveInboundBlocked() ([]*entity.Inbound, error)
	RetriveActiveInboundExpired() ([]*entity.Inbound, error)
	RetriveActiveInboundsOverQuota() ([]*entity.Inbound, error)
	RetriveDeactiveInboundsCharged() ([]*entity.Inbound, error)
	RetriveActiveInbounds() ([]*entity.Inbound, error)
	Active(id int) error
	DeActive(id int) error
}

type UserService interface {
	FindByID(string) (*entity.User, error)
}

type VpnService interface {
	MakeProxy() (adapter.ProxyVPN, error)
}

type InboundChargeService interface {
	ChargeInbound(*entity.Charge) error
}

type ChargeService interface {
	FindAvailbleCharge(string) (*entity.Charge, error)
}

type UserHealingUpService interface {
	FindByID(string) (*entity.User, error)
}

func NewHealingUpInbound(
	inboundRepo HealingUpInboundRepo,
	vpnService VpnService,
	inboundChargeSvc InboundChargeService,
	chargeSvc ChargeService,
	userService UserHealingUpService,
) *HealingUpInbound {
	return &HealingUpInbound{
		inboundRepo:      inboundRepo,
		vpnService:       vpnService,
		inboundChargeSvc: inboundChargeSvc,
		chargeSvc:        chargeSvc,
		userService:      userService,
	}
}

func (i *HealingUpInbound) CheckInboundsActivation() {
	inbounds, err := i.inboundRepo.RetriveActiveInbounds()
	if err != nil {
		return
	}
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.checkInboundActivation(inbound, proxy)
	}
}

func (i *HealingUpInbound) HealingUpExpiredInbounds() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundExpired()
	if err != nil {
		return
	}
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.healingUpExpiredInbound(inbound, proxy)
	}
}

func (i *HealingUpInbound) HealingUpOverQuotedInbounds() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundsOverQuota()
	if err != nil {
		return
	}
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.healingUpExpiredInbound(inbound, proxy)
	}
}

func (i *HealingUpInbound) HealingUpBlockedInbounds() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundBlocked()
	if err != nil {
		return
	}
	i.deactiveInbounds(inbounds)
}

func (i *HealingUpInbound) HealingUpChargedInbounds() {
	inbounds, err := i.inboundRepo.RetriveDeactiveInboundsCharged()
	if err != nil {
		return
	}
	i.activeInbounds(inbounds)
}

func (i *HealingUpInbound) activeInbounds(inbounds []*entity.Inbound) {
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.activeInbound(inbound, proxy)
	}
}

func (i *HealingUpInbound) deactiveInbounds(inbounds []*entity.Inbound) {
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.deActiveInbound(inbound, proxy)
	}
}

func (i *HealingUpInbound) healingUpExpiredInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	charge, err := i.chargeSvc.FindAvailbleCharge(strconv.Itoa(inbound.ID))
	if err != nil {
		return i.deActiveInbound(inbound, vpnProxy)
	}
	return i.inboundChargeSvc.ChargeInbound(charge)
}

func (i *HealingUpInbound) deActiveInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return err
	}

	vpnProxyInput := adapter.GenerateVPNProxyInput(inbound, user)
	err = vpnProxy.DisableInbound(vpnProxyInput)
	if err != nil {
		return err
	}

	return i.inboundRepo.DeActive(inbound.ID)
}

func (i *HealingUpInbound) activeInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return err
	}

	vpnProxyInput := adapter.GenerateVPNProxyInput(inbound, user)

	if err != nil {
		return err
	}
	err = vpnProxy.AddInbound(vpnProxyInput)
	if err != nil {
		return err
	}

	return i.inboundRepo.Active(inbound.ID)
}

func (i *HealingUpInbound) checkInboundActivation(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) {
	isActive, err := i.isInboundActive(inbound, vpnProxy)
	if err != nil {
		return
	}
	if isActive {
		return
	}

	err = i.inboundRepo.DeActive(inbound.ID)
	if err != nil {
		return
	}
}

func (i *HealingUpInbound) isInboundActive(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) (bool, error) {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return false, err
	}

	vpnProxyInput := adapter.GenerateVPNProxyInput(inbound, user)
	isActive, err := vpnProxy.IsInboundActive(vpnProxyInput)
	if err != nil {
		return false, err
	}

	return isActive, nil
}
