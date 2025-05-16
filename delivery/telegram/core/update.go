package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Update struct {
	*tgbotapi.Update
	UserSystem *entity.User
}
