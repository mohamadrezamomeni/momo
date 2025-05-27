package inbound

import (
	"fmt"

	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/pkg/utils"

	vpnProxyDto "github.com/mohamadrezamomeni/momo/dto/proxy/vpn"
	inboundRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"

	"github.com/google/uuid"
	vpnProxy "github.com/mohamadrezamomeni/momo/proxy/vpn"
)

type Inbound struct {
	inboundRepo InboundRepo
	vpnService  VpnService
	userService UserService
	hostService HostService
}

type VpnService interface {
	MakeProxy() (adapter.ProxyVPN, error)
}

type UserService interface {
	FindByID(string) (*entity.User, error)
}

type InboundRepo interface {
	IncreaseTrafficUsage(string, uint32) error
	Update(string, *inboundRepoDto.UpdateInboundDto) error
	FindInboundByID(id string) (*entity.Inbound, error)
	UpdateDomainPort(int, string, string) error
	ChangeBlockState(string, bool) error
	RetriveActiveInboundBlocked() ([]*entity.Inbound, error)
	RetriveActiveInboundExpired() ([]*entity.Inbound, error)
	RetriveActiveInboundsOverQuota() ([]*entity.Inbound, error)
	RetriveDeactiveInboundsCharged() ([]*entity.Inbound, error)
	Active(id int) error
	Filter(*inboundRepoDto.FilterInbound) ([]*entity.Inbound, error)
	DeActive(id int) error
	Create(*inboundRepoDto.CreateInbound) (*entity.Inbound, error)
	FindInboundIsNotAssigned() ([]*entity.Inbound, error)
	GetListOfPortsByDomain() ([]struct {
		Domain string
		Ports  []string
	}, error)
	ExtendInbound(string, *inboundRepoDto.ExtendInboundDto) error
}

type HostService interface {
	ResolveHostPortPair(map[string][]string, int) ([][2]string, error)
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

func (i *Inbound) Create(inpt *inboundServiceDto.CreateInbound) (*entity.Inbound, error) {
	inboundCreated, err := i.inboundRepo.Create(&inboundRepoDto.CreateInbound{
		Tag:          fmt.Sprintf("inbound-%s", uuid.New().String()),
		Port:         "",
		Domain:       "",
		IsActive:     false,
		IsBlock:      false,
		UserID:       inpt.UserID,
		Start:        inpt.Start,
		End:          inpt.End,
		VPNType:      inpt.VPNType,
		TrafficLimit: inpt.TrafficLimit,
	})
	if err != nil {
		return nil, err
	}
	return inboundCreated, nil
}

func (i *Inbound) AssignDomainToInbounds() {
	inbounds, err := i.inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		return
	}
	portUserSummery, err := i.summeryDomainPorts()
	if err != nil {
		return
	}

	hostPortPairs, err := i.hostService.ResolveHostPortPair(portUserSummery, len(inbounds))
	if err != nil {
		return
	}

	for j := 0; j < utils.Min(len(inbounds), len(hostPortPairs)); j++ {
		hostPort := hostPortPairs[j]
		inbound := inbounds[j]
		host, port := hostPort[0], hostPort[1]
		i.inboundRepo.UpdateDomainPort(inbound.ID, host, port)
	}
}

func (i *Inbound) summeryDomainPorts() (map[string][]string, error) {
	summery, err := i.inboundRepo.GetListOfPortsByDomain()
	if err != nil {
		return nil, err
	}
	res := map[string][]string{}
	for _, item := range summery {
		res[item.Domain] = item.Ports
	}
	return res, nil
}

func (i *Inbound) HealingUpInboundExpired() {
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
		i.deActiveInbound(inbound, proxy)
	}
}

func (i *Inbound) HealingUpInboundOverQuoted() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundsOverQuota()
	if err != nil {
		return
	}
	i.deactiveInbounds(inbounds)
}

func (i *Inbound) HealingUpInboundBlocked() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundBlocked()
	if err != nil {
		return
	}
	i.deactiveInbounds(inbounds)
}

func (i *Inbound) HealingUpInboundCharged() {
	inbounds, err := i.inboundRepo.RetriveDeactiveInboundsCharged()
	if err != nil {
		return
	}
	i.activeInbounds(inbounds)
}

func (i *Inbound) activeInbounds(inbounds []*entity.Inbound) {
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.activeInbound(inbound, proxy)
	}
}

func (i *Inbound) deactiveInbounds(inbounds []*entity.Inbound) {
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.deActiveInbound(inbound, proxy)
	}
}

func (i *Inbound) deActiveInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	info, err := i.getInfo(inbound)
	if err != nil {
		return err
	}
	err = vpnProxy.DisableInbound(info)
	if err != nil {
		return err
	}

	return i.inboundRepo.DeActive(inbound.ID)
}

func (i *Inbound) activeInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	info, err := i.getInfo(inbound)
	if err != nil {
		return err
	}
	err = vpnProxy.AddInbound(info)
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
		VPNType: inbound.VPNType,
		Address: inbound.Domain,
		Port:    inbound.Port,
		Tag:     inbound.Tag,
	}
	return info, nil
}

func (i *Inbound) Filter(inpt *inboundServiceDto.FilterInbounds) ([]*entity.Inbound, error) {
	return i.inboundRepo.Filter(&inboundRepoDto.FilterInbound{
		Domain:  inpt.Domain,
		Port:    inpt.Port,
		VPNType: inpt.VPNType,
		UserID:  inpt.UserID,
	})
}

func (i *Inbound) Block(id string) error {
	err := i.inboundRepo.ChangeBlockState(id, true)
	if err != nil {
		return err
	}
	return nil
}

func (i *Inbound) UnBlock(id string) error {
	err := i.inboundRepo.ChangeBlockState(id, false)
	if err != nil {
		return err
	}
	return nil
}

func (i *Inbound) ExtendInbound(id string, inpt *inboundServiceDto.ExtendInboundDto) error {
	err := i.inboundRepo.ExtendInbound(id, &inboundRepoDto.ExtendInboundDto{
		End:             inpt.End,
		Start:           inpt.Start,
		TrafficExtended: inpt.ExtendedTrafficLimit,
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *Inbound) FindInboundByID(id string) (*entity.Inbound, error) {
	return i.inboundRepo.FindInboundByID(id)
}

func (i *Inbound) UpdateInbound(id string, inpt *inboundServiceDto.UpdateDto) error {
	err := i.inboundRepo.Update(id, &inboundRepoDto.UpdateInboundDto{
		Start:        inpt.Start,
		End:          inpt.End,
		TrafficLimit: inpt.TrafficLimit,
	})
	if err != nil {
		return err
	}
	return nil
}
