package auth

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) AskUsername(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}

	usernameButtonText, err := telegrammessages.GetMessage("auth.username_button", map[string]string{})
	if err != nil {
		return nil, err
	}
	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, usernameButtonText)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
	}, nil
}

func (h *Handler) AnswerUsername(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}

	state, isExist := telegramState.FindState(idStr)
	if !isExist {
		return nil, err
	}

	state.SetData("username", update.Message.Text)
	return nil, nil
}
