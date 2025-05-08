package host

import (
	"strings"

	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func TransformStringsToHostStatus(hostStatusStrings string) ([]entity.HostStatus, error) {
	scope := "tansformer.TransformStringsToHostStatus"
	if len(hostStatusStrings) == 0 {
		return []entity.HostStatus{}, nil
	}

	hostStatues := make([]entity.HostStatus, 0)
	for _, s := range strings.Split(hostStatusStrings, ",") {
		hostStatus, err := entity.MapHostStatusToEnum(s)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
		}
		hostStatues = append(hostStatues, hostStatus)
	}
	return hostStatues, nil
}
