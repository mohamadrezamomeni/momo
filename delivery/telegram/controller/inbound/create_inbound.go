package inbound

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
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
	pkg     *entity.VPNPackage
}

var creatingInboundKey = "creating_inbound"

func (h *Handler) SetState(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		val := telegramState.GetControllerState(update.UserSystem.TelegramID, creatingInboundKey)
		_, ok := val.(*CreateVPNState)
		if !ok {
			state := &CreateVPNState{
				state: askVPNType,
			}
			telegramState.SetControllerState(update.UserSystem.TelegramID, creatingInboundKey, state)
		}

		return next(update)
	}
}

func (h *Handler) AskVPNType(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.controller.AskVPNType"

		val := telegramState.GetControllerState(update.UserSystem.TelegramID, creatingInboundKey)

		state, ok := val.(*CreateVPNState)

		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
		}

		if state.state != askVPNType {
			return next(update)
		}

		markup := h.getVPNTypeButtoms()

		id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
		if err != nil {
			return nil, err
		}

		text, err := telegrammessages.GetMessage(
			"inbound.create.select_vpn",
			map[string]string{},
		)
		if err != nil {
			return nil, err
		}

		msg := tgbotapi.NewMessage(id, text)

		msg.ReplyMarkup = markup
		state.state = answerVPNType

		telegramState.SetControllerState(update.UserSystem.TelegramID, creatingInboundKey, state)

		return &core.ResponseHandlerFunc{
			Result: msg,
		}, err
	}
}

func (h *Handler) getVPNTypeButtoms() tgbotapi.InlineKeyboardMarkup {
	xray := tgbotapi.NewInlineKeyboardButtonData(
		entity.VPNTypeString(entity.XRAY_VPN), entity.VPNTypeString(entity.XRAY_VPN),
	)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(xray),
	)
}

func (h *Handler) FillVPNType(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.controller.fillVPNType"

		val := telegramState.GetControllerState(update.UserSystem.TelegramID, creatingInboundKey)
		state, ok := val.(*CreateVPNState)

		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
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

		telegramState.SetControllerState(update.UserSystem.TelegramID, creatingInboundKey, state)

		return next(update)
	}
}

func (h *Handler) AskPackagesCreatingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AskPackagesCreatingInbound"

		val := telegramState.GetControllerState(update.UserSystem.TelegramID, creatingInboundKey)

		state, ok := val.(*CreateVPNState)

		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
		}

		if state.state != askPackage {
			return next(update)
		}

		res, err := h.getResponseAskPackage(update.UserSystem)
		if err != nil {
			return nil, err
		}

		state.state = answerPackage
		telegramState.SetControllerState(update.UserSystem.TelegramID, creatingInboundKey, state)
		return res, nil
	}
}

func (h *Handler) AnswerPackageCreatingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AnswerPackageCreatingInbound"
		val := telegramState.GetControllerState(update.UserSystem.TelegramID, creatingInboundKey)

		state, ok := val.(*CreateVPNState)
		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
		}

		packageID := update.CallbackQuery.Data

		pkg, err := h.answerPackage(packageID)
		if err != nil {
			return nil, err
		}
		state.pkg = pkg
		state.state = creatingInboundDone
		telegramState.SetControllerState(update.UserSystem.TelegramID, creatingInboundKey, state)
		return next(update)
	}
}

func (h *Handler) CreateInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.CreateVPN"

	val := telegramState.GetControllerState(update.UserSystem.TelegramID, creatingInboundKey)

	state, ok := val.(*CreateVPNState)

	if !ok {
		return nil, momoError.Scope(scope).Errorf("the type is difference")
	}

	if state.state != creatingInboundDone {
		return nil, momoError.Scope(scope).Errorf("state must be done")
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
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
	text, err := telegrammessages.GetMessage("inbound.create.successfully_creation", map[string]string{})
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
