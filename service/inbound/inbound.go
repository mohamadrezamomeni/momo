package inbound

import (
	"fmt"
	"time"

	vpnProxyDto "momo/dto/proxy/vpn"
	inboundRepoDto "momo/dto/repository/inbound"
	dto "momo/dto/service/inbound"
	"momo/entity"
	vpnSerializer "momo/proxy/vpn/serializer"

	"github.com/google/uuid"
)

type Inbound struct {
	inboundRepo inboundRepo
	vpnService  VpnService
	userService userService
	hostService hostService
}

type VpnProxy interface {
	AddInbound(*vpnProxyDto.Inbound, entity.VPNType) error
	DisableInbound(*vpnProxyDto.Inbound, entity.VPNType) error
	GetTraffic(*vpnProxyDto.Inbound, entity.VPNType) (*vpnSerializer.Traffic, error)
	Close()
}

type VpnService interface {
	MakeProxy() (VpnProxy, error)
}

type userService interface {
	FindByID(string) (*entity.User, error)
}

type inboundRepo interface {
	RetriveFaultyInbounds() ([]*entity.Inbound, error)
	Active(id int) error
	DeActive(id int) error
	Create(*inboundRepoDto.CreateInbound) (entity.Inbound, error)
	FindInboundIsNotAssigned()
}

type hostService interface {
	FindRightHost(entity.HostStatus) (string, string, error)
}

func New(
	repo inboundRepo,
	vpnService VpnService,
	userService userService,
	hostService hostService,
) *Inbound {
	return &Inbound{
		inboundRepo: repo,
		vpnService:  vpnService,
		userService: userService,
		hostService: hostService,
	}
}

func (i *Inbound) Create(inpt *dto.CreateInbound) error {
	_, err := i.inboundRepo.Create(&inboundRepoDto.CreateInbound{
		Tag:      fmt.Sprintf("inbound-%s", uuid.New().String()),
		Port:     "",
		Domain:   "",
		IsActive: false,
		IsBlock:  false,
		UserID:   inpt.UserID,
		Start:    inpt.Start,
		End:      inpt.End,
		VPNType:  inpt.VPNType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *Inbound) ApplyChangesToInbounds() error {
	inbounds, err := i.inboundRepo.RetriveFaultyInbounds()
	if err != nil {
		return err
	}
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return err
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.applyChangeToInbound(inbound, proxy)
	}
	return nil
}

func (i *Inbound) applyChangeToInbound(inbound *entity.Inbound, vpnProxy VpnProxy) {
	if i.mustItBeActive(inbound) {
		i.activeInbound(inbound, inbound.VPNType, vpnProxy)
	} else {
		i.deActiveInbound(inbound, inbound.VPNType, vpnProxy)
	}
}

func (i *Inbound) mustItBeActive(inbound *entity.Inbound) bool {
	if now := time.Now(); inbound.IsBlock == false &&
		((now.Before(inbound.End) || now.Equal(inbound.End)) &&
			(now.After(inbound.Start) || now.Equal(inbound.Start))) &&
		inbound.IsActive {
		return true
	}
	return false
}

func (i *Inbound) deActiveInbound(inbound *entity.Inbound, vpnType entity.VPNType, vpnProxy VpnProxy) error {
	info, err := i.getInfo(inbound)
	if err != nil {
		return err
	}
	err = vpnProxy.DisableInbound(info, vpnType)
	if err != nil {
		return err
	}

	return i.inboundRepo.DeActive(inbound.ID)
}

func (i *Inbound) activeInbound(inbound *entity.Inbound, vpnType entity.VPNType, vpnProxy VpnProxy) error {
	info, err := i.getInfo(inbound)
	if err != nil {
		return err
	}
	err = vpnProxy.AddInbound(info, vpnType)
	if err != nil {
		return err
	}

	return i.inboundRepo.Active(inbound.ID)
}

func (i *Inbound) getInfo(inbound *entity.Inbound) (*vpnProxyDto.Inbound, error) {
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
