package core

import "github.com/mohamadrezamomeni/momo/entity"

type HandlerFunc = func(*Context) (*ResponseHandler, error)

type ResponseHandler struct {
	Messages []*Message
}

type Message struct {
	MenuTab bool
	User    *entity.User
	Message string
	ID      string
}

func (c *Core) Register(path string, handler HandlerFunc) {
	c.router[path] = handler
}
