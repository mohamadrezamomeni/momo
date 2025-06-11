package inbound

import (
	"fmt"

	"github.com/google/uuid"
	inboundRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	templategenerator "github.com/mohamadrezamomeni/momo/templates"
)

type Inbound struct {
	inboundRepo InboundRepo
}

type InboundRepo interface {
	Update(string, *inboundRepoDto.UpdateInboundDto) error
	FindInboundByID(id string) (*entity.Inbound, error)
	ChangeBlockState(string, bool) error
	Filter(*inboundRepoDto.FilterInbound) ([]*entity.Inbound, error)
	Create(*inboundRepoDto.CreateInbound) (*entity.Inbound, error)
	ExtendInbound(string, *inboundRepoDto.ExtendInboundDto) error
}

func New(
	repo InboundRepo,
) *Inbound {
	return &Inbound{
		inboundRepo: repo,
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
		Country:      inpt.Country,
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

func (i *Inbound) GetClientConfig(id string) (string, error) {
	scope := "inboundService.GetClientConfig"

	inbound, err := i.inboundRepo.FindInboundByID(id)
	if err != nil {
		return "", err
	}

	if inbound.IsBlock || !inbound.IsAssigned {
		return "", momoError.Scope(scope).Forbidden().ErrorWrite()
	}
	template, err := templategenerator.LoadClientConfig(inbound.VPNType, inbound.Domain, inbound.Port, inbound.UserID)
	if err != nil {
		return "", err
	}
	return template, nil
}
