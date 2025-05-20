package notification

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type TelegramConfig struct {
	Token string `koanf:"token"`
}

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func New(cfg *TelegramConfig) *Telegram {
	scope := "initialize.telegram"
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		momoError.Wrap(err).Scope(scope).Input(cfg).Fatal()
	}

	return &Telegram{
		bot: bot,
	}
}

func (t *Telegram) Send(id string, message string) error {
	scope := "sendnotification.telegram.send"
	idConverted, err := utils.ConvertToInt64(id)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().ErrorWrite()
	}

	msg := tgbotapi.NewMessage(idConverted, message)

	_, err = t.bot.Send(msg)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().ErrorWrite()
	}
	return nil
}
