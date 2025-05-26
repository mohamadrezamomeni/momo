package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/entity"
	telegramPKG "github.com/mohamadrezamomeni/momo/pkg/telegram"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type Core struct {
	router map[string]HandlerFunc
	bot    *tgbotapi.BotAPI
}

func New(
	bot *tgbotapi.BotAPI,
) *Core {
	return &Core{
		router: map[string]HandlerFunc{},
		bot:    bot,
	}
}

func (c *Core) Notify(path string, data string) {
	handler := c.router[path]
	res, err := handler(NewContext(data))
	if err != nil {
		return
	}
	for _, msg := range res.Messages {
		msgConfig, _ := c.makeMessageConfig(msg.ID, msg.Message)
		c.bot.Send(msgConfig)

		if msg.MenuTab {
			c.SendMenuTab(msg.ID, msg.User)
		}
	}
}

func (c *Core) SendMenuTab(telegramID string, user *entity.User) {
	msgConfig, err := telegramPKG.MenuConfigMessage(telegramID, user)
	if err != nil {
		return
	}
	c.bot.Send(msgConfig)
}

func (c *Core) makeMessageConfig(idStr string, message string) (*tgbotapi.MessageConfig, error) {
	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, err
	}
	msgConfig := tgbotapi.NewMessage(id, message)
	return &msgConfig, nil
}

func (c *Core) GetRoutes() []string {
	routes := make([]string, 0, len(c.router))
	for k := range c.router {
		routes = append(routes, k)
	}
	return routes
}
