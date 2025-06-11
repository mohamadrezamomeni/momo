package vpnsource

import (
	vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
)

type VPNSource struct {
	VPNSourceRepository VPNSourceRepository
}

type VPNSourceRepository interface {
	Create(*vpnSourceRepositoryDto.CreateVPNSourceDto) (*entity.VPNSource, error)
}

func New(VPNSourceRepo VPNSourceRepository) *VPNSource {
	return &VPNSource{
		VPNSourceRepository: VPNSourceRepo,
	}
}

func (vs *VPNSource) Create(VPNSourceDto *VPNSourceServiceDto.CreateVPNSourceDto) error {
	_, err := vs.VPNSourceRepository.Create(&vpnSourceRepositoryDto.CreateVPNSourceDto{
		Country: VPNSourceDto.Country,
		English: VPNSourceDto.English,
	})
	if err != nil {
		return err
	}
	return err
}
