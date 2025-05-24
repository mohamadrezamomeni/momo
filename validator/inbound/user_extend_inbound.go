package inbound

import (
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateExtendingInboundByUser(inbound *entity.Inbound, user *entity.User) error {
	scope := "validation.inbound.ValidateExtendingInboundByUser"

	if inbound.UserID != user.ID {
		return momoError.Scope(scope).Input(inbound, user).BadRequest().ErrorWrite()
	}
	return nil
}
