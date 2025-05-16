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
	SettingState CreatingTypeState = iota
	VPNType
	SetTrafficLimit
)

type CreateVPNState struct {
	state CreatingTypeState
	data  *inboundServiceDto.CreateInbound
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
		id, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		key := generateCreatingState(id)
		_, isExist, err := getCreatingVPNState(key)
		if err != nil {
			return nil, err
		}

		if !isExist {
			now := time.Now()
			creatingVPN := &inboundServiceDto.CreateInbound{
				VPNType:    0,
				ServerType: entity.High,
				UserID:     update.UserSystem.ID,
				Start:      time.Now(),
				End:        now.AddDate(0, 1, 0),
			}
			state := &CreateVPNState{
				data:  creatingVPN,
				state: SettingState,
			}
			cache.Set(key, state)
		}

		return next(update)
	}
}

func (h *Handler) FillVPNType(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.controller.fillVPNType"
		idStr, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		key := generateCreatingState(idStr)

		state, isExist, err := getCreatingVPNState(key)
		if err != nil || !isExist {
			return nil, momoError.Scope(scope).ErrorWrite()
		}

		if state.data.VPNType == 0 && SettingState == state.state {
			xray := tgbotapi.NewInlineKeyboardButtonData(
				entity.VPNTypeString(entity.XRAY_VPN), entity.VPNTypeString(entity.XRAY_VPN),
			)

			markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(xray),
			)

			id, _ := utils.ConvertToInt64(idStr)

			text, err := telegrammessages.GetMessage(
				"inbound.create.select_vpn",
				map[string]string{},
			)

			msg := tgbotapi.NewMessage(id, text)

			msg.ReplyMarkup = markup
			state.state = VPNType

			cache.Set(key, state)

			return &core.ResponseHandlerFunc{
				Result: msg,
			}, err
		}

		if state.data.VPNType == 0 && VPNType == state.state {
			message := update.CallbackQuery.Data
			state.data.VPNType = entity.ConvertStringVPNTypeToEnum(message)
			cache.Set(idStr, state)
		}

		return next(update)
	}
}

func (h *Handler) FillTrafficLimit(next core.HandlerFunc) core.HandlerFunc {
	return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
		scope := "telegram.controller.fillVPNType"
		idStr, err := core.GetID(update)
		if err != nil {
			return nil, err
		}
		key := generateCreatingState(idStr)

		state, isExist, err := getCreatingVPNState(key)
		if err != nil || !isExist {
			return nil, momoError.Scope(scope).ErrorWrite()
		}

		if state.data.TrafficLimit == 0 && VPNType == state.state {
			fiftyG := tgbotapi.NewInlineKeyboardButtonData(
				"50GB", "50",
			)

			markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(fiftyG),
			)

			id, _ := utils.ConvertToInt64(idStr)

			text, err := telegrammessages.GetMessage(
				"inbound.create.select_traffic_limit",
				map[string]string{},
			)

			msg := tgbotapi.NewMessage(id, text)

			msg.ReplyMarkup = markup
			state.state = SetTrafficLimit
			cache.Set(key, state)

			return &core.ResponseHandlerFunc{
				Result: msg,
			}, err
		}

		if state.data.TrafficLimit == 0 && SetTrafficLimit == state.state {
			trafficLimitStr := update.CallbackQuery.Data
			traffLimit, _ := utils.ConvertToUint32(trafficLimitStr)
			state.data.TrafficLimit = traffLimit
			cache.Set(idStr, state)
		}

		return next(update)
	}
}

func (h *Handler) CreateVPN(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}
	key := generateCreatingState(idStr)

	state, isExist, err := getCreatingVPNState(key)
	if err != nil || !isExist {
		return nil, err
	}

	defer func() {
		cache.Delete(key)
	}()

	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, err
	}

	text, err := telegrammessages.GetMessage("inbound.create.successfully_creation", map[string]string{})
	if err != nil {
		return nil, err
	}

	_, err = h.inboundSvc.Create(state.data)
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
