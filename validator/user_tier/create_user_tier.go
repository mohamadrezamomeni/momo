package user

import (
	"fmt"

	userTierControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/user_tier"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateCreatingUserTier(data userTierControllerDto.Create) error {
	scope := "validator.user.ValidateCreatingUserTier"
	user, err := v.userSvc.FindByID(data.UserID)
	fmt.Println(user, err)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	if user == nil {
		return momoError.Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	tier, err := v.tierSvc.FindByName(data.Tier)
	fmt.Println(tier, err)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	if tier == nil {
		return momoError.Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	return nil
}
