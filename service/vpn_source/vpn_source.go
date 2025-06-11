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
	Upsert(string, *vpnSourceRepositoryDto.UpsertVPNSourceDto) (*entity.VPNSource, error)
}

func New(VPNSourceRepo VPNSourceRepository) *VPNSource {
	return &VPNSource{
		VPNSourceRepository: VPNSourceRepo,
	}
}

func (vs *VPNSource) Create(VPNSourceDto *VPNSourceServiceDto.CreateVPNSourceDto) error {
	_, err := vs.VPNSourceRepository.Upsert(VPNSourceDto.Country, &vpnSourceRepositoryDto.UpsertVPNSourceDto{
		English: VPNSourceDto.English,
	})
	if err != nil {
		return err
	}
	return err
}
