package service

import (
	"time"

	vpnProxyDto "momo/dto/proxy/vpn"
	"momo/entity"
	vpnSerializer "momo/proxy/vpn/serializer"
)

type Inbound struct {
	inboundRepo inboundRepo
	vpnProxy    vpnProxy
	userService userService
}

type vpnProxy interface {
	AddInbound(*vpnProxyDto.Inbound, entity.VPNType) error
	DisableInbound(*vpnProxyDto.Inbound, entity.VPNType) error
	GetTraffic(*vpnProxyDto.Inbound, entity.VPNType) (*vpnSerializer.Traffic, error)
}

type userService interface {
	FindByID(string) (*entity.User, error)
}

type inboundRepo interface {
	RetriveFaultyInbounds() ([]*entity.Inbound, error)
	Active(id int) error
	DeActive(id int) error
}

func New(repo inboundRepo, vpnService vpnProxy, userService userService) *Inbound {
	return &Inbound{
		inboundRepo: repo,
		vpnProxy:    vpnService,
		userService: userService,
	}
}

func (i Inbound) ApplyChangesToInbounds() error {
	inbounds, err := i.inboundRepo.RetriveFaultyInbounds()
	if err != nil {
		return err
	}

	for _, inbound := range inbounds {
		i.applyChangeToInbound(inbound)
	}
	return nil
}

func (i Inbound) applyChangeToInbound(inbound *entity.Inbound) {
	if i.mustItBeActive(inbound) {
		i.activeInbound(inbound, inbound.VPNType)
	} else {
		i.deActiveInbound(inbound, inbound.VPNType)
	}
}

func (i Inbound) mustItBeActive(inbound *entity.Inbound) bool {
	if now := time.Now(); inbound.IsBlock == false &&
		((now.Before(inbound.End) || now.Equal(inbound.End)) &&
			(now.After(inbound.Start) || now.Equal(inbound.Start))) &&
		inbound.IsActive {
		return true
	}
	return false
}

func (i Inbound) deActiveInbound(inbound *entity.Inbound, vpnType entity.VPNType) error {
	info, err := i.getInfo(inbound)
	if err != nil {
		return err
	}
	err = i.vpnProxy.DisableInbound(info, vpnType)
	if err != nil {
		return err
	}

	return i.inboundRepo.DeActive(inbound.ID)
}

func (i Inbound) activeInbound(inbound *entity.Inbound, vpnType entity.VPNType) error {
	info, err := i.getInfo(inbound)
	if err != nil {
		return err
	}
	err = i.vpnProxy.AddInbound(info, vpnType)
	if err != nil {
		return err
	}

	return i.inboundRepo.Active(inbound.ID)
}

func (i Inbound) getInfo(inbound *entity.Inbound) (*vpnProxyDto.Inbound, error) {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return &vpnProxyDto.Inbound{}, err
	}
	info := &vpnProxyDto.Inbound{
		User: &vpnProxyDto.User{
			Username: user.Username,
			ID:       user.ID,
			Level:    "0",
		},
		Address: inbound.Domain,
		Port:    inbound.Port,
		Tag:     inbound.Tag,
	}
	return info, nil
}
