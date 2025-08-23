package charge

import (
	chargeEvent "github.com/mohamadrezamomeni/momo/event/charge"
	"github.com/mohamadrezamomeni/momo/notification/core"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

func (h *Handler) ApproveCharge(c *core.Context) (*core.ResponseHandler, error) {
	var approveChargeByAdmin chargeEvent.ApproveChargeEvent
	if err := c.Bind(&approveChargeByAdmin); err != nil {
		return nil, err
	}
	charge, err := h.chargeSvc.FindByID(approveChargeByAdmin.ID)
	if err != nil {
		return nil, err
	}

	user, err := h.userSvc.FindByID(charge.UserID)
	if err != nil {
		return nil, err
	}

	title, err := telegrammessages.GetMessage("charge.approve_charge_admin", map[string]string{
		"id": charge.ID,
	}, user.Language)
	if err != nil {
		return nil, err
	}

	return &core.ResponseHandler{
		Messages: []*core.Message{
			{
				User:    user,
				MenuTab: true,
				ID:      user.TelegramID,
				Message: title,
			},
		},
	}, nil
}
