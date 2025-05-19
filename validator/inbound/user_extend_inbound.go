package inbound

import (
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateExtendingInboundByUser(inbound *entity.Inbound, user *entity.User) error {
	scope := "validation.inbound.ValidateExtendingInboundByUser"
	now := time.Now()
	if now.Before(inbound.End) && inbound.TrafficLimit > inbound.TrafficUsage {
		return momoError.Scope(scope).Input(inbound, user).BadRequest().DeactiveWrite().ErrorWrite()
	}
	if inbound.UserID != user.ID {
		return momoError.Scope(scope).Input(inbound, user).BadRequest().ErrorWrite()
	}
	return nil
}
