package vpnsource

import (
	vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
)

type VPNSource struct {
	VPNSourceRepository VPNSourceRepository
	VPNManagerSvc       VPNManagerService
}

type VPNSourceRepository interface {
	Upsert(string, *vpnSourceRepositoryDto.UpsertVPNSourceDto) (*entity.VPNSource, error)
	Filter(*vpnSourceRepositoryDto.FilterVPNSourcesDto) ([]*entity.VPNSource, error)
	Find(string) (*entity.VPNSource, error)
}

type VPNManagerService interface {
	GetAvailableCountries() ([]string, error)
}

func New(
	VPNSourceRepo VPNSourceRepository,
	VPNManagerSvc VPNManagerService,
) *VPNSource {
	return &VPNSource{
		VPNSourceRepository: VPNSourceRepo,
		VPNManagerSvc:       VPNManagerSvc,
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

func (vs *VPNSource) Find(country string) (*entity.VPNSource, error) {
	return vs.VPNSourceRepository.Find(country)
}

func (vs *VPNSource) FilterVPNSources(filterVPNSourcesDto *VPNSourceServiceDto.FilterVPNSourcesDto) ([]*entity.VPNSource, error) {
	var err error
	var countries []string = []string{}
	if filterVPNSourcesDto.Available {
		countries, err = vs.VPNManagerSvc.GetAvailableCountries()
	}
	if err != nil {
		return nil, err
	}

	if filterVPNSourcesDto.Countries != nil {
		countries = append(countries, filterVPNSourcesDto.Countries...)
	}

	return vs.VPNSourceRepository.Filter(&vpnSourceRepositoryDto.FilterVPNSourcesDto{
		Countries: countries,
	})
}
