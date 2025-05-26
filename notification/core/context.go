package core

import (
	"encoding/json"
	"reflect"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type Context struct {
	data string
}

func NewContext(data string) *Context {
	return &Context{
		data: data,
	}
}

func (c *Context) Bind(i interface{}) error {
	scope := "notification.context.bind"
	t := reflect.ValueOf(i)

	if t.Kind() != reflect.Ptr {
		return momoError.Scope(scope).DebuggingError()
	}

	if t.Elem().Kind() != reflect.Struct {
		return momoError.Scope(scope).DebuggingError()
	}

	return json.Unmarshal([]byte(c.data), i)
}
