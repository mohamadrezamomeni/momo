package time

import (
	"time"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func ConvertStrToTime(timeStr string) (time.Time, error) {
	scope := "validation.ConvertStrToTime"
	t, err := time.Parse(time.DateTime, timeStr)
	if err != nil {
		return time.Time{}, momoError.Wrap(err).Scope(scope).Input(timeStr).ErrorWrite()
	}
	return t, nil
}
