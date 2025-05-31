package charge

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) AskDetail(update *core.Update) (*core.ResponseHandlerFunc, error) {
	message, err := telegrammessages.GetMessage("charge.detail_button", map[string]string{})
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}
	msgConfig := tgbotapi.NewMessage(id, message)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
	}, nil
}

func (h *Handler) AnswerDetail(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.AnswerDetail"

	message := update.Message.Text
	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)
	if !isExist {
		return nil, momoError.Scope(scope).DebuggingError()
	}

	state.SetData("detail", message)
	return nil, nil
}
