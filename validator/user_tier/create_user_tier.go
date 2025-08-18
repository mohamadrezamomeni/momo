package user

import (
	userTierControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/user_tier"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateCreatingUserTier(data userTierControllerDto.Create) error {
	scope := "validator.user.ValidateCreatingUserTier"

	user, err := v.userSvc.FindByID(data.UserID)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	if user != nil {
		return momoError.Scope(scope).BadRequest().Input(data).ErrorWrite()
	}

	tier, err := v.tierSvc.FindByName(data.Tier)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	if tier != nil {
		return momoError.Scope(scope).BadRequest().Input(data).ErrorWrite()
	}
	return nil
}
