package core

type HandlerFunc = func(*Context) ([]*ResponseHandler, error)

type ResponseHandler struct {
	Message string
	ID      string
}

func (c *Core) Register(path string, handler HandlerFunc) {
	c.router[path] = handler
}
