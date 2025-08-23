package inbound

import (
	"encoding/base64"
	"encoding/json"
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
	var protocol string
	if inpt.VPNType == entity.XRAY_VPN {
		protocol = "vmess"
	}
	tag := i.GenerateInboundTag(inpt.Country, inpt.UserID)
	inboundCreated, err := i.inboundRepo.Create(&inboundRepoDto.CreateInbound{
		Tag:          tag,
		Protocol:     protocol,
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

func (i *Inbound) GenerateInboundTag(country string, userID string) string {
	return fmt.Sprintf("%s-%s-%s", country, userID, uuid.New().String())
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

func (i *Inbound) LoadInboundURI(inboundID string) (string, error) {
	scope := "inboundService.LoadInboundURI"

	inbound, err := i.inboundRepo.FindInboundByID(inboundID)
	if err != nil {
		return "", err
	}

	cfg := struct {
		V    string `json:"v"`
		Ps   string `json:"ps"`
		Add  string `json:"add"`
		Port string `json:"port"`
		ID   string `json:"id"`
		Aid  string `json:"aid"`
		Net  string `json:"net"`
		Type string `json:"type"`
		Host string `json:"host"`
		Path string `json:"path"`
		TLS  string `json:"tls"`
	}{
		V:    "2",
		Ps:   "V2Ray-Node",
		Add:  inbound.Domain,
		Port: inbound.Port,
		ID:   inbound.UserID,
		Aid:  "0",
		Net:  "tcp",
		Type: "none",
		Host: "",
		Path: "",
		TLS:  "",
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Input(inboundID).ErrorWrite()
	}

	encoded := base64.StdEncoding.EncodeToString(jsonBytes)

	uri := "vmess://" + encoded
	return uri, nil
}
