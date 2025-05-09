package inbound

import (
	"fmt"
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
		ID:         i.id,
		VPNType:    entity.XRAY_VPN,
		Protocol:   inpt.Protocol,
		Tag:        inpt.Tag,
		IsActive:   inpt.IsActive,
		Domain:     inpt.Domain,
		Port:       inpt.Port,
		IsAssigned: inpt.IsAssigned,
		IsNotified: inpt.IsNotified,
		IsBlock:    inpt.IsBlock,
		Start:      inpt.Start,
		End:        inpt.End,
	}
	i.inbounds = append(i.inbounds, inbound)
	return inbound, nil
}

func (i *MockInbound) FindInboundByID(id int) (*entity.Inbound, error) {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			return inbound, nil
		}
	}
	return nil, fmt.Errorf("the record wasn't find")
}

func (i *MockInbound) Delete(id int) error {
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
		if val, ok := hashmap[inbound.Domain]; ok {
			hashmap[inbound.Domain] = append(val, inbound.Port)
		} else {
			hashmap[inbound.Domain] = []string{inbound.Port}
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

func (i *MockInbound) Active(id int) error {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			inbound.IsActive = true
			return nil
		}
	}
	return fmt.Errorf("the record hasn't found")
}

func (i *MockInbound) DeActive(id int) error {
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

func (i *MockInbound) RetriveFaultyInbounds() ([]*entity.Inbound, error) {
	result := make([]*entity.Inbound, 0)
	now := time.Now()
	for _, inbound := range i.inbounds {
		if now.After(inbound.End) ||
			(inbound.IsBlock == true && inbound.IsActive == true) ||
			(inbound.IsBlock == false && !now.After(inbound.Start) && !inbound.End.After(now) && inbound.IsBlock == false) {
			result = append(result, inbound)
		}
	}
	return result, nil
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

func (i *MockInbound) UpdateDomainPort(id int, domain string, port string) error {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			inbound.Domain = domain
			inbound.Port = port
			inbound.IsAssigned = true
		}
	}
	return nil
}

func (i *MockInbound) ChangeBlockState(_ string, _ bool) error {
	return nil
}
