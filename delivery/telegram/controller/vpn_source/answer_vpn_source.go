package vpnsource

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (h *Handler) AnswerVPNSource(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.vpnSource.AnswerVPNSource"
	vpnSourceCountry := update.CallbackQuery.Data
	vpnSource, err := h.vpnSourceService.Find(vpnSourceCountry)
	if err != nil {
		return nil, err
	}

	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)

	if !isExist {
		return nil, momoError.Scope(scope).DebuggingError()
	}
	state.SetData("vpn_source", vpnSource)
	return nil, nil
}
