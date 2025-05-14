package auth

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	"github.com/mohamadrezamomeni/momo/pkg/cache"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
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

func (h *Handler) CheckDuplicateRegistration(next core.HandlerFunc) core.HandlerFunc {
	scope := "telegram.controller.CheckDuplicateRegistration"
	return func(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
		id, err := core.GetID(update)
		if err != nil {
			return nil, err
		}

		user, err := h.userSvc.FindByTelegramID(id)

		if user != nil {
			return nil, momoError.Scope(scope).Input(update).Errorf("you were registered before")
		}

		momoErr, ok := err.(*momoError.MomoError)
		if !ok {
			return nil, err
		}

		if momoErr.GetErrorType() != momoError.NotFound {
			return nil, err
		}
		return next(update)
	}
}

func (h *Handler) SetUserRegisteration(next core.HandlerFunc) core.HandlerFunc {
	return func(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)

		_, isExist := cache.Get(key)
		if !isExist {
			userRegisteration := &UserRegisteration{
				ID:        update.FromChat().ID,
				FirstName: update.FromChat().FirstName,
				LastName:  update.FromChat().LastName,
			}
			cache.Set(key, userRegisteration)

			msg := tgbotapi.NewMessage(update.FromChat().ID, "please input your username:")
			return &core.ResponseHandlerFunc{
				Result:       msg,
				ReleaseState: false,
			}, nil
		}

		return next(update)
	}
}

func (h *Handler) SetUsername(next core.HandlerFunc) core.HandlerFunc {
	return func(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)

		value, _ := cache.Get(key)

		userRegistertion, _ := value.(*UserRegisteration)
		if userRegistertion.Username == "" {
			userRegistertion.Username = update.Message.Text
			cache.Set(key, userRegistertion)
		}

		if userRegistertion.FirstName == "" {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "please put your firstname:")
			return &core.ResponseHandlerFunc{
				Result:       msg,
				ReleaseState: false,
			}, nil
		}
		return next(update)
	}
}

func (h *Handler) SetFirstname(next core.HandlerFunc) core.HandlerFunc {
	return func(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)

		value, _ := cache.Get(key)
		userRegistertion, _ := value.(*UserRegisteration)

		if userRegistertion.FirstName == "" {
			userRegistertion.FirstName = update.Message.Text
			cache.Set(key, userRegistertion)
		}

		if userRegistertion.LastName == "" {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "please put your lastname:")
			return &core.ResponseHandlerFunc{
				Result:       msg,
				ReleaseState: false,
			}, nil
		}
		return next(update)
	}
}

func (h *Handler) SetLastname(next core.HandlerFunc) core.HandlerFunc {
	return func(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
		idStr := strconv.Itoa(int(update.FromChat().ID))
		key := generateRegistrationKey(idStr)

		value, _ := cache.Get(key)
		userRegistertion, _ := value.(*UserRegisteration)

		if userRegistertion.LastName == "" {
			userRegistertion.LastName = update.Message.Text
			cache.Set(key, userRegistertion)
		}

		return next(update)
	}
}

func (h *Handler) Register(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
	idStr := strconv.Itoa(int(update.FromChat().ID))
	key := generateRegistrationKey(idStr)

	value, _ := cache.Get(key)
	userRegistertion, _ := value.(*UserRegisteration)
	_, err := h.authSvc.Register(&authServiceDto.RegisterDto{
		Username:   userRegistertion.Username,
		Firstname:  userRegistertion.FirstName,
		Lastname:   userRegistertion.LastName,
		TelegramID: idStr,
	})
	if err != nil {
		cache.Delete(key)
		return nil, err
	}

	message := fmt.Sprintf(
		"your data is :\nusername: %s\nfirstname: %s\nlastname: %s",
		userRegistertion.Username,
		userRegistertion.FirstName,
		userRegistertion.LastName,
	)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)

	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
		RedirectRoot: true,
	}, nil
}
