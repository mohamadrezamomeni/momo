package telegrammessages

import (
	"fmt"
	"strings"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func GetMessage(path string, data map[string]string) (string, error) {
	scope := "telegrammessages.getmessages"

	msg := k.String(path)
	if msg == "" {
		return "", momoError.Scope(scope).Input(path, data).ErrorWrite()
	}

	for k, v := range data {
		msg = strings.ReplaceAll(msg, fmt.Sprintf(":%s:", k), v)
	}
	return msg, nil
}
