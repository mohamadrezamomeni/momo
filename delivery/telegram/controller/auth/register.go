package auth

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	"github.com/mohamadrezamomeni/momo/pkg/cache"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegramMessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
)

type UserRegisteration struct {
	Username  string
	FirstName string
	LastName  string
	ID        int64
}

func generateRegistrationKey(id string) string {
	return id + "-registeration"
}

func generateStateKey(id string) string {
	return id + "-state"
}

func (h *Handler) CheckDuplicateRegistration(next core.HandlerFunc) core.HandlerFunc {
	scope := "telegram.controller.CheckDuplicateRegistration"
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		id, err := core.GetID(update)
		if err != nil {
			return nil, err
		}

		user, err := h.userSvc.FindByTelegramID(id)

		if user != nil {
			msg, err := telegramMessages.GetMessage("auth.registeration.duplicate_register", map[string]string{})

			return &core.ResponseHandlerFunc{
				Result:       tgbotapi.NewMessage(update.FromChat().ID, msg),
				ReleaseState: false,
			}, err
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

func (h *Handler) SetUserRegisteration(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)
		stateKey := generateStateKey(idStr)

		_, isExist := cache.Get(key)
		if !isExist {
			userRegisteration := &UserRegisteration{
				ID:        update.FromChat().ID,
				FirstName: update.FromChat().FirstName,
				LastName:  update.FromChat().LastName,
			}
			cache.Set(key, userRegisteration)

			msg, err := telegramMessages.GetMessage("auth.registeration.input_username", map[string]string{})

			cache.Set(stateKey, "username")
			return &core.ResponseHandlerFunc{
				Result:       tgbotapi.NewMessage(update.FromChat().ID, msg),
				ReleaseState: false,
			}, err
		}

		return next(update)
	}
}

func (h *Handler) SetUsername(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)
		stateKey := generateStateKey(idStr)

		value, isExistState := cache.Get(stateKey)
		state := ""
		if isExistState {
			state, _ = value.(string)
		}

		value, _ = cache.Get(key)
		userRegistertion, _ := value.(*UserRegisteration)
		if userRegistertion.Username == "" && state == "username" {
			userRegistertion.Username = update.Message.Text
			cache.Set(key, userRegistertion)
		}

		if userRegistertion.FirstName == "" && state != "firstname" {
			msg, err := telegramMessages.GetMessage("auth.registeration.input_firstname", map[string]string{})
			cache.Set(stateKey, "firstname")
			return &core.ResponseHandlerFunc{
				Result:       tgbotapi.NewMessage(update.FromChat().ID, msg),
				ReleaseState: false,
			}, err
		}
		return next(update)
	}
}

func (h *Handler) SetFirstname(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)
		stateKey := generateStateKey(idStr)

		value, isExistState := cache.Get(stateKey)
		state := ""
		if isExistState {
			state, _ = value.(string)
		}

		value, _ = cache.Get(key)
		userRegistertion, _ := value.(*UserRegisteration)

		if userRegistertion.FirstName == "" && state == "firstname" {
			userRegistertion.FirstName = update.Message.Text
			cache.Set(key, userRegistertion)
		}

		if userRegistertion.LastName == "" && state != "lastname" {
			cache.Set(stateKey, "lastname")
			msg, err := telegramMessages.GetMessage("auth.registeration.input_lastname", map[string]string{})
			return &core.ResponseHandlerFunc{
				Result:       tgbotapi.NewMessage(update.FromChat().ID, msg),
				ReleaseState: false,
			}, err
		}
		return next(update)
	}
}

func (h *Handler) SetLastname(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)
		stateKey := generateStateKey(idStr)

		value, isExistState := cache.Get(stateKey)
		state := ""
		if isExistState {
			state, _ = value.(string)
		}

		value, _ = cache.Get(key)
		userRegistertion, _ := value.(*UserRegisteration)

		if userRegistertion.LastName == "" && state == "lastname" {
			userRegistertion.LastName = update.Message.Text
			cache.Set(key, userRegistertion)
		}

		return next(update)
	}
}

func (h *Handler) Register(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr := strconv.Itoa(int(update.FromChat().ID))
	key := generateRegistrationKey(idStr)
	stateKey := generateStateKey(idStr)

	defer func() {
		cache.Delete(key)
		cache.Delete(stateKey)
	}()

	value, _ := cache.Get(key)
	userRegistertion, _ := value.(*UserRegisteration)
	_, err := h.authSvc.Register(&authServiceDto.RegisterDto{
		Username:   userRegistertion.Username,
		Firstname:  userRegistertion.FirstName,
		Lastname:   userRegistertion.LastName,
		TelegramID: idStr,
	})
	if err != nil {
		return nil, err
	}

	msg, err := telegramMessages.GetMessage("auth.registeration.successfully_registeration", map[string]string{
		"username":  userRegistertion.Username,
		"firstname": userRegistertion.FirstName,
		"lastname":  userRegistertion.LastName,
	})

	return &core.ResponseHandlerFunc{
		Result:       tgbotapi.NewMessage(update.Message.Chat.ID, msg),
		ReleaseState: true,
		RedirectRoot: true,
	}, err
}
