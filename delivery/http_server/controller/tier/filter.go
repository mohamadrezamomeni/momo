package tier

import (
	"net/http"

	"github.com/labstack/echo/v4"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	tierSerializer "github.com/mohamadrezamomeni/momo/serializer/tier"
)

func (h *Handler) Filter(c echo.Context) error {
	tiers, err := h.tierSvc.Filter()
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	tiersSerializer := make([]*tierSerializer.Tier, 0)
	for _, tier := range tiers {
		tiersSerializer = append(tiersSerializer, &tierSerializer.Tier{
			IsDefault: tier.IsDefault,
			Name:      tier.Name,
		})
	}
	return c.JSON(http.StatusAccepted, &tierSerializer.Tiers{
		Tiers: tiersSerializer,
	})
}
