package inbound

import (
	"fmt"
	"strconv"
	"time"

	inboundDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
)

type MockInbound struct {
	inbounds []*entity.Inbound
	id       int
}

func New() *MockInbound {
	inbounds := make([]*entity.Inbound, 0)
	return &MockInbound{
		inbounds: inbounds,
		id:       0,
	}
}

func (i *MockInbound) Create(inpt *inboundDto.CreateInbound) (*entity.Inbound, error) {
	i.id += 1
	inbound := &entity.Inbound{
		ID:         strconv.Itoa(i.id),
		VPNType:    entity.XRAY_VPN,
		Protocol:   inpt.Protocol,
		Tag:        inpt.Tag,
		IsActive:   inpt.IsActive,
		Domain:     inpt.Domain,
		Port:       inpt.Port,
		Country:    inpt.Country,
		IsAssigned: inpt.IsAssigned,
		IsNotified: inpt.IsNotified,
		IsBlock:    inpt.IsBlock,
		Start:      inpt.Start,
		End:        inpt.End,
	}
	i.inbounds = append(i.inbounds, inbound)
	return inbound, nil
}

func (i *MockInbound) FindInboundByID(id string) (*entity.Inbound, error) {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			return inbound, nil
		}
	}
	return nil, fmt.Errorf("the record wasn't find")
}

func (i *MockInbound) Delete(id string) error {
	idx := -1
	for i, inbound := range i.inbounds {
		if inbound.ID == id {
			idx = i
		}
	}
	i.inbounds = append(i.inbounds[:idx], i.inbounds[idx+1:]...)
	return fmt.Errorf("the record wasn't find")
}

func (i *MockInbound) DeleteAll() error {
	i.inbounds = make([]*entity.Inbound, 0)
	return nil
}

func (i *MockInbound) GetListOfPortsByDomain() ([]struct {
	Domain string
	Ports  []string
}, error,
) {
	hashmap := map[string][]string{}
	for _, inbound := range i.inbounds {
		if inbound.IsAssigned {
			hashmap[inbound.Domain] = append(hashmap[inbound.Domain], inbound.Port)
		}
	}
	ret := make([]struct {
		Domain string
		Ports  []string
	}, 0)
	for k, v := range hashmap {
		ret = append(ret, struct {
			Domain string
			Ports  []string
		}{Domain: k, Ports: v})
	}
	return ret, nil
}

func (i *MockInbound) Active(id string) error {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			inbound.IsActive = true
			return nil
		}
	}
	return fmt.Errorf("the record hasn't found")
}

func (i *MockInbound) DeActive(id string) error {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			inbound.IsActive = false
			return nil
		}
	}
	return fmt.Errorf("the record hasn't found")
}

func (i *MockInbound) Filter(inpt *inboundDto.FilterInbound) ([]*entity.Inbound, error) {
	return nil, nil
}

func (i *MockInbound) RetriveActiveInboundBlocked() ([]*entity.Inbound, error) {
	inbounds := make([]*entity.Inbound, 0)
	for _, inbound := range i.inbounds {
		if inbound.IsActive && inbound.IsBlock {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds, nil
}

func (i *MockInbound) RetriveActiveInboundExpired() ([]*entity.Inbound, error) {
	inbounds := make([]*entity.Inbound, 0)
	now := time.Now()
	for _, inbound := range i.inbounds {
		if now.After(inbound.End) {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds, nil
}

func (i *MockInbound) RetriveActiveInboundsOverQuota() ([]*entity.Inbound, error) {
	inbounds := make([]*entity.Inbound, 0)
	for _, inbound := range i.inbounds {
		if inbound.TrafficLimit <= inbound.TrafficUsage {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds, nil
}

func (i *MockInbound) RetriveDeactiveInboundsCharged() ([]*entity.Inbound, error) {
	inbounds := make([]*entity.Inbound, 0)
	now := time.Now()
	for _, inbound := range i.inbounds {
		if !(now.Before(inbound.Start) || now.After(inbound.End)) &&
			inbound.TrafficLimit > inbound.TrafficUsage {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds, nil
}

func (i *MockInbound) RetriveChargedInbounds() ([]*entity.Inbound, error) {
	inbounds := make([]*entity.Inbound, 0)
	now := time.Now()
	for _, inbound := range inbounds {
		if inbound.TrafficLimit > inbound.TrafficUsage &&
			inbound.IsActive == false &&
			inbound.IsBlock == false &&
			!now.Before(inbound.Start) &&
			now.Before(inbound.End) {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds, nil
}

func (i *MockInbound) FindInboundIsNotAssigned() ([]*entity.Inbound, error) {
	result := make([]*entity.Inbound, 0)
	for _, inbound := range i.inbounds {
		if inbound.IsAssigned == false {
			result = append(result, inbound)
		}
	}
	return result, nil
}

func (i *MockInbound) UpdateDomainPort(id string, domain string, port string, VPNID string) error {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			inbound.Domain = domain
			inbound.Port = port
			inbound.IsAssigned = true
			inbound.VPNID = VPNID
		}
	}
	return nil
}

func (i *MockInbound) ChangeBlockState(_ string, _ bool) error {
	return nil
}

func (i *MockInbound) Update(_ string, _ *inboundDto.UpdateInboundDto) error {
	return nil
}

func (i *MockInbound) ExtendInbound(_ string, _ *inboundDto.ExtendInboundDto) error {
	return nil
}

func (i *MockInbound) IncreaseTrafficUsage(_ string, _ uint64) error {
	return nil
}

func (i *MockInbound) RetriveActiveInbounds() ([]*entity.Inbound, error) {
	inbounds := make([]*entity.Inbound, 0)
	now := time.Now()
	for _, inbound := range i.inbounds {
		if inbound.IsActive && inbound.End.After(now) && !inbound.IsBlock {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds, nil
}
