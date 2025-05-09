package inbound

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/pkg/utils"

	vpnProxyDto "github.com/mohamadrezamomeni/momo/dto/proxy/vpn"
	inboundRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	dto "github.com/mohamadrezamomeni/momo/dto/service/inbound"

	vpnProxy "github.com/mohamadrezamomeni/momo/proxy/vpn"

	"github.com/google/uuid"
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
	UpdateDomainPort(int, string, string) error
	RetriveFaultyInbounds() ([]*entity.Inbound, error)
	Active(id int) error
	Filter(*inboundRepoDto.FilterInbound) ([]*entity.Inbound, error)
	DeActive(id int) error
	Create(*inboundRepoDto.CreateInbound) (*entity.Inbound, error)
	FindInboundIsNotAssigned() ([]*entity.Inbound, error)
	GetListOfPortsByDomain() ([]struct {
		Domain string
		Ports  []string
	}, error)
}

type HostService interface {
	FindRightHosts(entity.HostStatus) ([]*entity.Host, error)
	ResolvePorts(
		*entity.Host,
		int,
		[]string,
		*sync.WaitGroup,
		chan<- struct {
			Domain string
			Ports  []string
		},
	)
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

func (i *Inbound) Create(inpt *dto.CreateInbound) (*entity.Inbound, error) {
	inboundCreated, err := i.inboundRepo.Create(&inboundRepoDto.CreateInbound{
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
		return nil, err
	}
	return inboundCreated, nil
}

func (i *Inbound) AssignDomainToInbounds() {
	inbounds, err := i.inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		return
	}

	hosts, err := i.hostService.FindRightHosts(entity.High)
	if err != nil {
		return
	}

	portUserSummery, err := i.summeryDomainPorts()
	if err != nil {
		return
	}

	ch := make(chan struct {
		Domain string
		Ports  []string
	})

	var wg sync.WaitGroup
	seen := map[string]struct{}{}

	for _, host := range hosts {
		wg.Add(1)
		if _, ok := seen[host.Domain]; ok {
			continue
		}

		seen[host.Domain] = struct{}{}

		go i.hostService.ResolvePorts(
			host,
			len(inbounds),
			portUserSummery[host.Domain],
			&wg,
			ch,
		)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	hostPortPairs := [][2]string{}

	for item := range ch {
		hostPortPairs = append(hostPortPairs, i.makeHostPairWiPort(item.Domain, item.Ports)...)
	}

	hostPortPairs = i.shuffleHostPortPairs(hostPortPairs)
	for j := 0; j < utils.Min(len(inbounds), len(hostPortPairs)); j++ {
		hostPort := hostPortPairs[j]
		inbound := inbounds[j]
		host, port := hostPort[0], hostPort[1]
		i.inboundRepo.UpdateDomainPort(inbound.ID, host, port)
	}
	return
}

func (i *Inbound) shuffleHostPortPairs(hostPortPairs [][2]string) [][2]string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Shuffle(len(hostPortPairs), func(i, j int) {
		hostPortPairs[i], hostPortPairs[j] = hostPortPairs[j], hostPortPairs[i]
	})
	return hostPortPairs
}

func (i *Inbound) makeHostPairWiPort(host string, ports []string) [][2]string {
	hostPortPairs := [][2]string{}
	for _, port := range ports {
		hostPortPairs = append(hostPortPairs, [2]string{host, port})
	}
	return hostPortPairs
}

func (i *Inbound) HealingUpInbounds() {
	inbounds, err := i.inboundRepo.RetriveFaultyInbounds()
	if err != nil {
		return
	}
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.HealingUpInbound(inbound, proxy)
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

func (i *Inbound) HealingUpInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) {
	if i.mustItBeActive(inbound) {
		i.activeInbound(inbound, vpnProxy)
	} else {
		i.deActiveInbound(inbound, vpnProxy)
	}
}

func (i *Inbound) mustItBeActive(inbound *entity.Inbound) bool {
	if now := time.Now(); inbound.IsBlock == false &&
		((now.Before(inbound.End) || now.Equal(inbound.End)) &&
			(now.After(inbound.Start) || now.Equal(inbound.Start))) &&
		!inbound.IsActive {
		return true
	}
	return false
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

func (i *Inbound) Filter(inpt *dto.FilterInbounds) ([]*entity.Inbound, error) {
	return i.inboundRepo.Filter(&inboundRepoDto.FilterInbound{
		Domain:  inpt.Domain,
		Port:    inpt.Port,
		VPNType: inpt.VPNType,
		UserID:  inpt.UserID,
	})
}
