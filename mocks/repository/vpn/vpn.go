package vpn

import (
	"fmt"
	"strconv"

	vpnManagerDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	"github.com/mohamadrezamomeni/momo/entity"
)

type MockVPN struct {
	vpns []*entity.VPN
	id   int
}

func New() *MockVPN {
	return &MockVPN{
		vpns: make([]*entity.VPN, 0),
		id:   0,
	}
}

func (mv *MockVPN) Create(inpt *vpnManagerDto.AddVPN) (*entity.VPN, error) {
	vpn := &entity.VPN{
		ID:        mv.id,
		IsActive:  inpt.IsActive,
		VPNType:   inpt.VPNType,
		ApiPort:   inpt.ApiPort,
		UserCount: inpt.UserCount,
	}
	mv.id += 1
	mv.vpns = append(mv.vpns, vpn)
	return vpn, nil
}

func (mv *MockVPN) FindByID(id int) *entity.VPN {
	for _, vpn := range mv.vpns {
		if vpn.ID == id {
			return vpn
		}
	}
	return nil
}

func (mv *MockVPN) ActiveVPN(id int) error {
	for _, vpn := range mv.vpns {
		if vpn.ID == id {
			vpn.IsActive = true
			return nil
		}
	}
	return fmt.Errorf("record hasn't found")
}

func (mv *MockVPN) DeactiveVPN(id int) error {
	for _, vpn := range mv.vpns {
		if vpn.ID == id {
			vpn.IsActive = false
			return nil
		}
	}
	return fmt.Errorf("record hasn't found")
}

func (mv *MockVPN) Filter(inpt *vpnManagerDto.FilterVPNs) ([]*entity.VPN, error) {
	ret := []*entity.VPN{}
	VPNTypeMap := make(map[entity.VPNType]struct{})
	countryMap := make(map[string]struct{})
	for _, vpnType := range inpt.VPNTypes {
		VPNTypeMap[vpnType] = struct{}{}
	}
	for _, country := range inpt.Coountries {
		countryMap[country] = struct{}{}
	}
	for _, vpn := range mv.vpns {
		_, isExistVPNType := VPNTypeMap[vpn.VPNType]
		_, isExistCountry := countryMap[vpn.Country]
		if (inpt.Domain == "" || vpn.Domain == inpt.Domain) &&
			(inpt.IsActive == nil || vpn.IsActive == *inpt.IsActive) &&
			(inpt.VPNTypes == nil || isExistVPNType) &&
			(inpt.Coountries == nil || isExistCountry) {
			ret = append(ret, vpn)
		}
	}

	return ret, nil
}

func (mv *MockVPN) GroupAvailbleVPNsByCountry() ([]string, error) {
	countriesRefrence := map[string]struct{}{}
	res := make([]string, 0)
	for _, vpn := range mv.vpns {
		if _, isExist := countriesRefrence[vpn.Country]; !isExist {
			res = append(res, vpn.Country)
		}
	}
	return res, nil
}

func (mv *MockVPN) Update(id string, updateVPNDto *vpnManagerDto.UpdateVPN) error {
	for _, mv := range mv.vpns {
		if id == strconv.Itoa(mv.ID) {
			mv.VPNStatus = updateVPNDto.VPNStatus
		}
	}
	return nil
}
