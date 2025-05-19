package inbound

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/pkg/cache"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type CreatingTypeState = int

const (
	askVPNType CreatingTypeState = iota
	answerVPNType
	creatingInboundDone
)

type CreateVPNState struct {
	state   CreatingTypeState
	VPNType entity.VPNType
	pkg     *entity.Package
}

func generateCreatingState(id string) string {
	return id + "creating_inbound"
}

func getCreatingVPNState(key string) (*CreateVPNState, bool, error) {
	scope := "telegram.inbound.getCreatingVPNState"
	value, isExist := cache.Get(key)
	if !isExist {
		return nil, false, nil
	}

	creatingVPNState, ok := value.(*CreateVPNState)
	if !ok {
		return nil, true, momoError.Scope(scope).UnExpected().DebuggingError()
	}
	return creatingVPNState, true, nil
}

func (h *Handler) SetState(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		key := generateCreatingState(update.UserSystem.ID)
		_, isExist, err := getCreatingVPNState(key)
		if err != nil {
			return nil, err
		}

		if !isExist {

			state := &CreateVPNState{
				state: askVPNType,
			}
			cache.Set(key, state)
		}

		return next(update)
	}
}

func (h *Handler) AskVPNType(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.controller.fillVPNType"

		key := generateCreatingState(update.UserSystem.ID)

		state, isExist, err := getCreatingVPNState(key)
		if err != nil || !isExist {
			return nil, momoError.Scope(scope).ErrorWrite()
		}
		if state.state != askVPNType {
			return next(update)
		}
		xray := tgbotapi.NewInlineKeyboardButtonData(
			entity.VPNTypeString(entity.XRAY_VPN), entity.VPNTypeString(entity.XRAY_VPN),
		)

		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(xray),
		)

		id, _ := utils.ConvertToInt64(update.UserSystem.TelegramID)

		text, err := telegrammessages.GetMessage(
			"inbound.create.select_vpn",
			map[string]string{},
		)

		msg := tgbotapi.NewMessage(id, text)

		msg.ReplyMarkup = markup
		state.state = answerVPNType

		cache.Set(key, state)

		return &core.ResponseHandlerFunc{
			Result: msg,
		}, err
	}
}

func (h *Handler) FillVPNType(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.controller.fillVPNType"

		key := generateCreatingState(update.UserSystem.ID)

		state, isExist, err := getCreatingVPNState(key)
		if err != nil || !isExist {
			return nil, momoError.Scope(scope).ErrorWrite()
		}
		if state.state != answerVPNType {
			return next(update)
		}

		message := update.CallbackQuery.Data
		vpnType := entity.ConvertStringVPNTypeToEnum(message)
		if vpnType == entity.UknownVPNType {
			return nil, momoError.Scope(scope).ErrorWrite()
		}

		state.VPNType = vpnType
		state.state = askPackage
		cache.Set(key, state)

		return next(update)
	}
}

func (h *Handler) AskPackagesCreatingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AskPackagesCreatingInbound"
		user := update.UserSystem

		key := generateCreatingState(user.ID)

		state, isExist, err := getCreatingVPNState(key)
		if err != nil || !isExist {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().Errorf("error to get state")
		}

		if state.state != askPackage {
			return next(update)
		}

		res, err := h.getResponseAskPackage(update.UserSystem)
		if err != nil {
			return nil, err
		}

		state.state = answerPackage
		cache.Set(key, state)
		return res, nil
	}
}

func (h *Handler) AnswerPackageCreatingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AnswerPackageCreatingInbound"
		user := update.UserSystem

		key := generateCreatingState(user.ID)

		state, isExist, err := getCreatingVPNState(key)
		if err != nil || !isExist {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
		}

		packageID := update.CallbackQuery.Data

		pkg, err := h.answerPackage(packageID)
		if err != nil {
			return nil, err
		}
		state.pkg = pkg
		state.state = creatingInboundDone
		cache.Set(key, state)
		return next(update)
	}
}

func (h *Handler) CreateInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.CreateVPN"

	key := generateCreatingState(update.UserSystem.ID)

	state, isExist, err := getCreatingVPNState(key)
	if err != nil || !isExist || state.state != creatingInboundDone {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
	}

	defer func() {
		cache.Delete(key)
	}()

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}

	text, err := telegrammessages.GetMessage("inbound.create.successfully_creation", map[string]string{})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = h.inboundSvc.Create(&inboundServiceDto.CreateInbound{
		ServerType:   entity.High,
		UserID:       update.UserSystem.ID,
		VPNType:      state.VPNType,
		TrafficLimit: state.pkg.TrafficLimit,
		Start:        now,
		End:          now.AddDate(0, int(state.pkg.Months), int(state.pkg.Days)),
	})
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(id, text)

	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
		RedirectRoot: true,
	}, nil
}
