package inbound

import (
	"fmt"
	"time"

	proxyVpnDto "momo/dto/proxy/vpn"
	vpnProxyDto "momo/dto/proxy/vpn"
	inboundRepoDto "momo/dto/repository/inbound"
	dto "momo/dto/service/inbound"
	"momo/entity"
	"momo/pkg/utils"
	"momo/proxy/vpn/serializer"
	workerProxy "momo/proxy/worker"

	"github.com/google/uuid"
)

type Inbound struct {
	inboundRepo InboundRepo
	vpnService  VpnService
	userService UserService
	hostService HostService
}

type VpnProxy interface {
	AddInbound(*proxyVpnDto.Inbound, int) error
	Close()
	DisableInbound(*proxyVpnDto.Inbound, int) error
	GetTraffic(*proxyVpnDto.Inbound, int) (*serializer.Traffic, error)
}

type VpnService interface {
	MakeProxy() (VpnProxy, error)
}

type UserService interface {
	FindByID(string) (*entity.User, error)
}

type InboundRepo interface {
	UpdateDomainPort(int, string, string) error
	RetriveFaultyInbounds() ([]*entity.Inbound, error)
	Active(id int) error
	DeActive(id int) error
	Create(*inboundRepoDto.CreateInbound) (*entity.Inbound, error)
	FindInboundIsNotAssigned() ([]*entity.Inbound, error)
}

type HostService interface {
	FindRightHosts(entity.HostStatus) ([]*entity.Host, error)
}

func New(
	repo InboundRepo,
	vpnService VpnService,
	userService UserService,
	hostService HostService,
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

func (i *Inbound) AssignDomainToInbounds() error {
	inbounds, err := i.inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		return err
	}
	hosts, err := i.hostService.FindRightHosts(entity.High)

	workerErrors := make(chan struct{})
	assigns := make(chan struct {
		inboundID int
		port      string
		domain    string
	})
	for _, inbound := range inbounds {
		host := hosts[utils.GetRandom(0, len(hosts))]
		go i.assignDomainToInbound(inbound, host, assigns, workerErrors)
	}

	for j := 0; j < len(inbounds); j++ {
		select {
		case data := <-assigns:
			i.inboundRepo.UpdateDomainPort(data.inboundID, data.domain, data.port)
		case <-workerErrors:
		}
	}

	return nil
}

func (i *Inbound) assignDomainToInbound(inbound *entity.Inbound, host *entity.Host, assigns chan<- struct {
	inboundID int
	port      string
	domain    string
}, workerErrors chan<- struct{},
) {
	wp, err := workerProxy.New(&workerProxy.Config{
		Address: host.Domain,
		Port:    host.Port,
	})
	if err != nil {
		workerErrors <- struct{}{}
		return
	}
	defer wp.Close()
	port, err := wp.GetAvailablePort()
	if err != nil {
		workerErrors <- struct{}{}
	}
	assigns <- struct {
		inboundID int
		port      string
		domain    string
	}{
		inboundID: inbound.ID,
		port:      port,
		domain:    host.Domain,
	}
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
