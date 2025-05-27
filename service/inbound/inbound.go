package inbound

import (
	"fmt"

	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/entity"

	inboundRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"

	"github.com/google/uuid"
)

type Inbound struct {
	inboundRepo      InboundRepo
	vpnService       VpnService
	userService      UserService
	hostService      HostService
	chargeSvc        ChargeService
	inboundChargeSvc InboundChargeService
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

type ChargeService interface {
	FindAvailbleCharge(string) (*entity.Charge, error)
}

type InboundChargeService interface {
	ChargeInbound(*entity.Inbound, *entity.Charge) error
}

func New(
	repo InboundRepo,
	vpnService VpnService,
	userService UserService,
	hostService HostService,
	chargeService ChargeService,
	inboundChargeService InboundChargeService,
) *Inbound {
	return &Inbound{
		inboundRepo:      repo,
		vpnService:       vpnService,
		userService:      userService,
		hostService:      hostService,
		chargeSvc:        chargeService,
		inboundChargeSvc: inboundChargeService,
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
