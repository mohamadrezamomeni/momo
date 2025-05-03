package vpn

import (
	"fmt"

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
	if inpt.Domain == "" && inpt.IsActive == nil && inpt.VPNType == 0 {
		return mv.vpns, nil
	}
	for _, vpn := range mv.vpns {
		if inpt.Domain != "" && vpn.Domain == inpt.Domain {
			ret = append(ret, vpn)
		} else if inpt.IsActive != nil && vpn.IsActive == *inpt.IsActive {
			ret = append(ret, vpn)
		} else if inpt.VPNType != 0 && inpt.VPNType == vpn.VPNType {
			ret = append(ret, vpn)
		}
	}

	return ret, nil
}
