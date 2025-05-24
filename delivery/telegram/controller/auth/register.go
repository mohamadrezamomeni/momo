package auth

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	"github.com/mohamadrezamomeni/momo/dto/service/auth"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegramMessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

func (h *Handler) SetState(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegramauth.register.setstate"
		id, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		state, isExist := telegramState.FindState(id)

		if !isExist {
			telegramState.NewState(id, "askUsername", "answerUsername", "register")
			return nil, nil
		}

		path := state.GetPath()
		if path != "register" {
			return nil, momoError.Scope(scope).DebuggingError()
		}
		return next(update)
	}
}

func (h *Handler) Register(update *core.Update) (*core.ResponseHandlerFunc, error) {
	username, firstname, lastname, id, err := getData(update)
	if err != nil {
		return nil, err
	}
	_, err = h.authSvc.Register(&auth.RegisterDto{
		Username:   username,
		Firstname:  firstname,
		Lastname:   lastname,
		TelegramID: id,
	})
	if err != nil {
		return nil, err
	}

	title, err := telegramMessages.GetMessage("auth.registeration.successfully_registeration", map[string]string{
		"username": username,
	})
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, title)
	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  false,
	}, nil
}

func getData(update *core.Update) (string, string, string, string, error) {
	scope := "auth.getdata.register"

	id, err := core.GetID(update)
	if err != nil {
		return "", "", "", "", err
	}

	state, isExist := telegramState.FindState(id)
	if !isExist {
		return "", "", "", "", momoError.Scope(scope).UnExpected().ErrorWrite()
	}
	val, isExist := state.GetData("username")
	if !isExist {
		return "", "", "", "", momoError.Scope(scope).UnExpected().ErrorWrite()
	}
	username, isExist := val.(string)
	if !isExist {
		return "", "", "", "", momoError.Scope(scope).UnExpected().ErrorWrite()
	}

	return username, update.Message.From.FirstName, update.Message.From.LastName, id, nil
}
