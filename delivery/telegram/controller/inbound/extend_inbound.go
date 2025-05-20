package inbound

import (
	"strconv"
	"strings"
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

type ExtendingInboundStatus = int

const (
	askID ExtendingInboundStatus = iota
	answerID
	extendingInboundDone
)

type ExtendingInboundState struct {
	inbound *entity.Inbound
	pkg     *entity.VPNPackage
	state   ExtendingInboundStatus
}

var exteningInboundKey = "extend_Inbound"

func (h *Handler) SetExtendingInboundState(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		val := telegramState.GetControllerState(update.UserSystem.TelegramID, exteningInboundKey)
		_, ok := val.(*ExtendingInboundState)

		if !ok {
			extendingInboundState := &ExtendingInboundState{
				state: askID,
			}
			telegramState.SetControllerState(update.UserSystem.TelegramID, exteningInboundKey, extendingInboundState)
		}

		return next(update)
	}
}

func (h *Handler) SelectInboundIDInExtendingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.SelectInboundIDInExtendingInbound"
		val := telegramState.GetControllerState(update.UserSystem.TelegramID, exteningInboundKey)

		state, ok := val.(*ExtendingInboundState)
		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
		}

		if state.state != askID {
			return next(update)
		}

		inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
			UserID: update.UserSystem.ID,
		})

		if len(inbounds) == 0 {
			return h.sendNotFoundInbounds(update.UserSystem)
		}

		id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
		if err != nil {
			return nil, err
		}
		askInboundID, err := telegrammessages.GetMessage("inbound.extend.ask_id", map[string]string{})
		if err != nil {
			return nil, err
		}
		msg := tgbotapi.NewMessage(id, askInboundID)

		title, err := telegrammessages.GetMessage("inbound.extend.ask_id", map[string]string{})
		if err != nil {
			return nil, err
		}

		var sb strings.Builder

		sb.WriteString(title)

		var rows [][]tgbotapi.InlineKeyboardButton

		for _, inbound := range inbounds {
			button, err := h.makeExtendingInboundButtom(inbound)
			if err != nil {
				return nil, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(*button))
		}

		msg.ParseMode = "HTML"
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		state.state = answerID
		telegramState.SetControllerState(update.UserSystem.TelegramID, exteningInboundKey, state)

		return &core.ResponseHandlerFunc{
			Result: msg,
		}, nil
	}
}

func (h *Handler) makeExtendingInboundButtom(inbound *entity.Inbound) (*tgbotapi.InlineKeyboardButton, error) {
	itemTtitle, err := telegrammessages.GetMessage("inbound.extend.item", map[string]string{
		"VPNType": entity.VPNTypeString(inbound.VPNType),
		"ID":      strconv.Itoa(inbound.ID),
	})
	if err != nil {
		return nil, err
	}
	button := tgbotapi.NewInlineKeyboardButtonData(itemTtitle, strconv.Itoa(inbound.ID))
	return &button, nil
}

func (h *Handler) ChooseInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.ChooseInbound"
		val := telegramState.GetControllerState(update.UserSystem.TelegramID, exteningInboundKey)

		state, ok := val.(*ExtendingInboundState)

		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
		}

		if state.state != answerID {
			return next(update)
		}
		inboundID := update.CallbackQuery.Data

		inbound, err := h.inboundSvc.FindInboundByID(inboundID)
		if err != nil {
			return nil, err
		}

		err = h.inboundValidator.ValidateExtendingInboundByUser(inbound, update.UserSystem)
		if err != nil {
			return nil, err
		}
		state.state = askPackage
		state.inbound = inbound
		telegramState.SetControllerState(update.UserSystem.TelegramID, exteningInboundKey, state)

		return next(update)
	}
}

func (h *Handler) AskPackagesExtending(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AskServices"
		val := telegramState.GetControllerState(update.UserSystem.TelegramID, exteningInboundKey)

		state, ok := val.(*ExtendingInboundState)

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
		telegramState.SetControllerState(update.UserSystem.TelegramID, exteningInboundKey, state)
		return res, nil
	}
}

func (h *Handler) AnswerPackageExtendingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AnswerPackage"

		val := telegramState.GetControllerState(update.UserSystem.TelegramID, exteningInboundKey)

		state, ok := val.(*ExtendingInboundState)

		if !ok {
			return nil, momoError.Scope(scope).Errorf("the type is difference")
		}

		packageID := update.CallbackQuery.Data

		pkg, err := h.answerPackage(packageID)
		if err != nil {
			return nil, err
		}
		state.pkg = pkg
		state.state = extendingInboundDone
		telegramState.SetControllerState(update.UserSystem.TelegramID, exteningInboundKey, state)
		return next(update)
	}
}

func (h *Handler) ExtendInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.ExtendInbound"
	val := telegramState.GetControllerState(update.UserSystem.TelegramID, exteningInboundKey)

	state, ok := val.(*ExtendingInboundState)

	if !ok {
		return nil, momoError.Scope(scope).Errorf("the type is difference")
	}
	if state.state != extendingInboundDone {
		return nil, momoError.Scope(scope).Errorf("state must be done")
	}
	now := time.Now()

	err := h.inboundSvc.ExtendInbound(strconv.Itoa(state.inbound.ID), &inboundServiceDto.ExtendInboundDto{
		Start:                now,
		End:                  now.AddDate(0, int(state.pkg.Months), int(state.pkg.Days)),
		ExtendedTrafficLimit: state.pkg.TrafficLimit,
	})

	extendingInboundTitle, err := telegrammessages.GetMessage("inbound.extend.successfully_extending", map[string]string{})
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(id, extendingInboundTitle)
	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
		RedirectRoot: true,
	}, nil
}
