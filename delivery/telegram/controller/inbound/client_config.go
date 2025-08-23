package inbound

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) SetStateGettingClientConfig(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		_, isExist := telegramState.FindState(idStr)
		if !isExist {
			telegramState.NewState(idStr,
				"render_client_config_buttons",
				"generate_client_config",
			)
			return nil, nil
		}

		return next(update)
	}
}

func (h *Handler) GetClientConfig(update *core.Update) (*core.ResponseHandlerFunc, error) {
	inboundID := update.CallbackQuery.Data

	telegramID, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}
	errorTitle, err := telegrammessages.GetMessage(
		"inbound.client_config.worng_button_selected",
		map[string]string{},
		update.UserSystem.Language,
	)
	if err != nil {
		return nil, err
	}

	template, err := h.inboundSvc.GetClientConfig(inboundID)

	if momoErr, ok := momoError.GetMomoError(err); ok && momoErr.GetErrorType() == momoError.Forbidden {
		msgConfig := tgbotapi.NewMessage(telegramID, errorTitle)

		return &core.ResponseHandlerFunc{
			MessageConfig: &msgConfig,
			ReleaseState:  true,
			RedirectRoot:  true,
			MenuTab:       true,
		}, nil
	} else if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(telegramID, template)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  true,
		MenuTab:       true,
	}, nil
}
