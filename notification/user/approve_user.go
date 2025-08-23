package user

import (
	eventUser "github.com/mohamadrezamomeni/momo/event/user"
	"github.com/mohamadrezamomeni/momo/notification/core"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

func (h *Handler) ApproveUserByAddmin(c *core.Context) (*core.ResponseHandler, error) {
	var approveUserByAdmin eventUser.UserApproved
	if err := c.Bind(&approveUserByAdmin); err != nil {
		return nil, err
	}
	user, err := h.userSvc.FindByID(approveUserByAdmin.UserID)
	if err != nil {
		return nil, err
	}

	title, err := telegrammessages.GetMessage("auth.approve_user", map[string]string{}, user.Language)
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
