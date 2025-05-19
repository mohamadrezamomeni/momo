package inbound

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/pkg/cache"
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
	pkg     *entity.Package
	state   ExtendingInboundStatus
}

func getExtendingInboundKey(id string) string {
	return id + "-extend_Inbound"
}

func getExtendingInboundState(id string) (*ExtendingInboundState, bool, error) {
	scope := "telegram.extendinginbound.getExtendingInboundState"

	value, isExist := cache.Get(id)
	if !isExist {
		return nil, false, nil
	}

	state, ok := value.(*ExtendingInboundState)
	if !ok {
		return nil, false, momoError.Scope(scope).UnExpected().ErrorWrite()
	}

	return state, true, nil
}

func (h *Handler) SetExtendingInboundState(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		user := update.UserSystem

		key := getExtendingInboundKey(user.ID)

		_, isExist, err := getExtendingInboundState(key)
		if err != nil {
			return nil, err
		}

		if !isExist {
			extendingInboundState := &ExtendingInboundState{
				state: askID,
			}
			cache.Set(key, extendingInboundState)
		}
		return next(update)
	}
}

func (h *Handler) SelectInboundIDInExtendingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.SelectInboundIDInExtendingInbound"
		user := update.UserSystem

		key := getExtendingInboundKey(user.ID)

		state, isExist, err := getExtendingInboundState(key)
		if err != nil || !isExist {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
		}

		if state.state != askID {
			return next(update)
		}

		inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
			UserID: user.ID,
		})

		if len(inbounds) == 0 {
			return h.sendNotFoundInbounds(user)
		}

		state.state = answerID
		defer cache.Set(key, state)

		id, err := utils.ConvertToInt64(user.TelegramID)

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
		user := update.UserSystem

		key := getExtendingInboundKey(user.ID)

		state, isExist, err := getExtendingInboundState(key)
		if err != nil || !isExist {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
		}

		if state.state != answerID {
			return next(update)
		}

		inboundID := update.CallbackQuery.Data

		inbound, err := h.inboundSvc.FindInboundByID(inboundID)
		if err != nil {
			return nil, err
		}

		err = h.inboundValidator.ValidateExtendingInboundByUser(inbound, user)
		if err != nil {
			return nil, err
		}
		state.state = askPackage
		state.inbound = inbound
		cache.Set(key, state)

		return next(update)
	}
}

func (h *Handler) AskPackagesExtending(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AskServices"
		user := update.UserSystem

		key := getExtendingInboundKey(user.ID)

		state, isExist, err := getExtendingInboundState(key)
		if err != nil || !isExist {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
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

func (h *Handler) AnswerPackageExtendingInbound(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.AnswerPackage"
		user := update.UserSystem

		key := getExtendingInboundKey(user.ID)

		state, isExist, err := getExtendingInboundState(key)
		if err != nil || !isExist {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
		}

		packageID := update.CallbackQuery.Data

		pkg, err := h.answerPackage(packageID)
		if err != nil {
			return nil, err
		}
		state.pkg = pkg
		state.state = extendingInboundDone
		cache.Set(key, state)
		return next(update)
	}
}

func (h *Handler) ExtendInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.ExtendInbound"

	user := update.UserSystem
	key := getExtendingInboundKey(user.ID)

	state, isExist, err := getExtendingInboundState(key)
	if err != nil || !isExist || state.state != extendingInboundDone {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
	}

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
