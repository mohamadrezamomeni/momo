package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (t *Telegram) setMenu(update *core.Update) error {
	scope := "telegram.setMenu"

	idStr, err := core.GetID(update)
	if err != nil {
		return err
	}

	id, err := utils.ConvertToInt64(idStr)

	commands, err := t.getCommands()
	if err != nil {
		return err
	}

	user, _ := t.userSvc.FindByTelegramID(idStr)

	if user == nil {
		commands = []tgbotapi.BotCommand{}
	}

	s := tgbotapi.NewBotCommandScopeChat(id)
	cfg := tgbotapi.SetMyCommandsConfig{
		Commands: commands,
		Scope:    &s,
	}

	_, err = t.bot.Request(cfg)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	return nil
}

func (t *Telegram) getCommands() ([]tgbotapi.BotCommand, error) {
	GenerateClientConfigTitle, err := telegrammessages.GetMessage(
		"inbound.client_config_button",
		map[string]string{},
		entity.EN,
	)
	if err != nil {
		return nil, err
	}

	listVPNsTitle, err := telegrammessages.GetMessage("inbound.list_button", map[string]string{}, entity.EN)
	if err != nil {
		return nil, err
	}

	createVPNTitle, err := telegrammessages.GetMessage("inbound.create_button", map[string]string{}, entity.EN)
	if err != nil {
		return nil, err
	}

	createChargeTitle, err := telegrammessages.GetMessage("charge.extend_button", map[string]string{}, entity.EN)
	if err != nil {
		return nil, err
	}

	listChargesTitle, err := telegrammessages.GetMessage("charge.list.button", map[string]string{}, entity.EN)
	if err != nil {
		return nil, err
	}

	return []tgbotapi.BotCommand{
		{
			Description: GenerateClientConfigTitle,
			Command:     "generate_client_config",
		},
		{
			Description: createVPNTitle,
			Command:     "create_charge",
		},
		{
			Description: listVPNsTitle,
			Command:     "list_inbounds",
		},
		{
			Description: createChargeTitle,
			Command:     "charge_inbound",
		},
		{
			Description: listChargesTitle,
			Command:     "list_charges",
		},
	}, nil
}
