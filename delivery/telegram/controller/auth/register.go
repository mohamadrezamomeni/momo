package auth

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	"github.com/mohamadrezamomeni/momo/dto/service/auth"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegramMessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

type UserRegisteration struct {
	Username  string
	FirstName string
	LastName  string
	ID        int64
	State     RegisterationState
}

var registerationKey = "registeration"

type RegisterationState = int

const (
	AskUsername = iota
	AnswerUsername
	RegisterationDone
)

func (h *Handler) CheckDuplicateRegistration(next core.HandlerFunc) core.HandlerFunc {
	scope := "telegram.controller.CheckDuplicateRegistration"
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		id, err := core.GetID(update)
		if err != nil {
			return nil, err
		}

		user, err := h.userSvc.FindByTelegramID(id)

		if user != nil {
			return h.getDuplicateUserResponse(update)
		}

		momoErr, ok := err.(*momoError.MomoError)
		if !ok {
			return nil, momoError.Wrap(err).Scope(scope).Input(update).ErrorWrite()
		}
		if momoErr.GetErrorType() != momoError.NotFound {
			return nil, err
		}
		return next(update)
	}
}

func (h *Handler) getDuplicateUserResponse(update *core.Update) (*core.ResponseHandlerFunc, error) {
	msg, err := telegramMessages.GetMessage("auth.registeration.duplicate_register", map[string]string{})
	if err != nil {
		return nil, err
	}

	return &core.ResponseHandlerFunc{
		Result:       tgbotapi.NewMessage(update.FromChat().ID, msg),
		ReleaseState: false,
		RedirectRoot: true,
	}, nil
}

func (h *Handler) SetUserRegisteration(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))

		val := telegramState.GetControllerState(idStr, registerationKey)

		_, ok := val.(*UserRegisteration)

		if !ok {
			userRegisteration := &UserRegisteration{
				ID:        update.FromChat().ID,
				FirstName: update.FromChat().FirstName,
				LastName:  update.FromChat().LastName,
				State:     AskUsername,
			}
			telegramState.SetControllerState(idStr, registerationKey, userRegisteration)
		}

		return next(update)
	}
}

func (h *Handler) AskUsername(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "register.askUser.AskUsername"
		idStr := strconv.Itoa(int(update.FromChat().ID))
		val := telegramState.GetControllerState(idStr, registerationKey)
		registerationState, ok := val.(*UserRegisteration)
		if !ok {
			return nil, momoError.Scope(scope)
		}

		if registerationState.State != AskUsername {
			return next(update)
		}

		title, err := telegramMessages.GetMessage("auth.registeration.input_username", map[string]string{})
		if err != nil {
			return nil, err
		}

		registerationState.State = AnswerUsername
		telegramState.SetControllerState(idStr, registerationKey, registerationState)

		return &core.ResponseHandlerFunc{
			Result:       tgbotapi.NewMessage(update.FromChat().ID, title),
			ReleaseState: false,
			RedirectRoot: false,
		}, nil
	}
}

func (h *Handler) AnswerUsername(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "register.askUser.AnswerUsername"
		idStr := strconv.Itoa(int(update.FromChat().ID))
		val := telegramState.GetControllerState(idStr, registerationKey)
		registerationState, ok := val.(*UserRegisteration)
		if !ok {
			return nil, momoError.Scope(scope)
		}
		if registerationState.State != AnswerUsername {
			return next(update)
		}
		username := update.Message.Text
		registerationState.Username = username
		registerationState.State = RegisterationDone
		telegramState.SetControllerState(idStr, registerationKey, registerationState)
		return next(update)
	}
}

func (h *Handler) Register(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "register.askUser.Register"
	idStr := strconv.Itoa(int(update.FromChat().ID))
	val := telegramState.GetControllerState(idStr, registerationKey)
	registerationState, ok := val.(*UserRegisteration)
	if !ok {
		return nil, momoError.Scope(scope).DebuggingError()
	}
	if registerationState.State != RegisterationDone {
		return nil, momoError.Scope(scope).DebuggingErrorf("error to compare state")
	}

	_, err := h.authSvc.Register(&auth.RegisterDto{
		Username:   registerationState.Username,
		Firstname:  registerationState.FirstName,
		Lastname:   registerationState.LastName,
		TelegramID: strconv.Itoa(int(registerationState.ID)),
	})
	if err != nil {
		return nil, err
	}

	title, err := telegramMessages.GetMessage("auth.registeration.successfully_registeration", map[string]string{
		"username": registerationState.Username,
	})
	if err != nil {
		return nil, err
	}
	return &core.ResponseHandlerFunc{
		Result:       tgbotapi.NewMessage(update.Message.Chat.ID, title),
		ReleaseState: true,
		RedirectRoot: true,
	}, nil
}
