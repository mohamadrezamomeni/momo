package usertiers

import (
	userTierRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/user_tier"
	userTierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user_tier"
	"github.com/mohamadrezamomeni/momo/entity"
)

type UserTiers struct {
	userTierRepo UserTierRepo
}

type UserTierRepo interface {
	Create(*userTierRepoDto.Create) error
	Delete(*userTierRepoDto.IdentifyUserTier) error
	FilterTiersBelongToUser(string) ([]*entity.Tier, error)
	FilterTiersByUser(string) ([]*entity.Tier, error)
}

func New(userTierRepo UserTierRepo) *UserTiers {
	return &UserTiers{
		userTierRepo: userTierRepo,
	}
}

func (ut *UserTiers) Create(createDto *userTierServiceDto.Create) error {
	return ut.userTierRepo.Create(&userTierRepoDto.Create{
		UserID: createDto.UserID,
		Tier:   createDto.Tier,
	})
}

func (ut *UserTiers) Delete(identifyUserTierDto *userTierServiceDto.IdentifyUserTier) error {
	return ut.userTierRepo.Create(&userTierRepoDto.Create{
		UserID: identifyUserTierDto.UserID,
		Tier:   identifyUserTierDto.Tier,
	})
}

func (ut *UserTiers) FilterTiersBelongToUser(userID string) ([]string, error) {
	tiers, err := ut.userTierRepo.FilterTiersBelongToUser(userID)
	if err != nil {
		return nil, err
	}
	tierStrs := make([]string, 0)
	for _, tier := range tiers {
		tierStrs = append(tierStrs, tier.Name)
	}

	return tierStrs, nil
}

func (ut *UserTiers) FilterTiersByUser(userID string) ([]*entity.Tier, error) {
	return ut.userTierRepo.FilterTiersByUser(userID)
}
