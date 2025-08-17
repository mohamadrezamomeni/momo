package tier

import (
	tierRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/tier"
	tierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/tier"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Tier struct {
	tierRepo TierRepo
}

type TierRepo interface {
	Create(*tierRepoDto.CreateTier) (*entity.Tier, error)
	Filter() ([]*entity.Tier, error)
	FindByName(string) (*entity.Tier, error)
	Update(string, *tierRepoDto.Update) error
}

func New(tierRepo TierRepo) *Tier {
	return &Tier{
		tierRepo: tierRepo,
	}
}

func (t *Tier) Create(createDto *tierServiceDto.CreateTier) (*entity.Tier, error) {
	tier, err := t.tierRepo.Create(&tierRepoDto.CreateTier{
		Name:      createDto.Name,
		IsDefault: createDto.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (t *Tier) Filter() ([]*entity.Tier, error) {
	return t.tierRepo.Filter()
}

func (t *Tier) FindByName(name string) (*entity.Tier, error) {
	return t.tierRepo.FindByName(name)
}

func (t *Tier) Update(name string, updateDto *tierServiceDto.Update) error {
	return t.tierRepo.Update(name, &tierRepoDto.Update{
		IsDefault: updateDto.IsDefault,
	})
}
