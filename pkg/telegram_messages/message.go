package telegrammessages

import (
	"fmt"
	"strings"

	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func GetMessage(
	path string,
	data map[string]string,
	language entity.Language,
) (string, error) {
	scope := "telegrammessages.getmessages"

	msg := k.String(
		fmt.Sprintf("%s.%s", entity.LanguageString(language), path),
	)
	if msg == "" {
		return "", momoError.Scope(scope).Input(path, data).ErrorWrite()
	}

	for k, v := range data {
		msg = strings.ReplaceAll(msg, fmt.Sprintf(":%s:", k), v)
	}
	return msg, nil
}
