package usertier

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userTierControllereDto "github.com/mohamadrezamomeni/momo/dto/controller/user_tier"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	tierSerializer "github.com/mohamadrezamomeni/momo/serializer/tier"
)

func (h *Handler) FilterTiersByUser(c echo.Context) error {
	var req userTierControllereDto.FilterTiersByUserID

	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	tiers, err := h.userTierSvc.FilterTiersByUser(req.UserID)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	userTiersSerializer := make([]*tierSerializer.Tier, 0)
	for _, tier := range tiers {
		userTiersSerializer = append(userTiersSerializer, &tierSerializer.Tier{
			IsDefault: tier.IsDefault,
			Name:      tier.Name,
		})
	}
	return c.JSON(http.StatusAccepted, &tierSerializer.Tiers{
		Tiers: userTiersSerializer,
	})
}
