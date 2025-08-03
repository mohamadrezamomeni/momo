package charge

import (
	"net/http"

	"github.com/labstack/echo/v4"
	chargeControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/charge"
	"github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	httperror "github.com/mohamadrezamomeni/momo/pkg/http_error"
	chargeSerializer "github.com/mohamadrezamomeni/momo/serializer/charge"
)

func (h *Handler) FilterCharges(c echo.Context) error {
	var req chargeControllerDto.FilterCharges
	if err := c.Bind(&req); err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validation.ValidateFilterCharges(req)
	if err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	charges, err := h.chargeSvc.FilterCharges(&charge.FilterCharges{
		InboundID: req.InboundID,
		UserID:    req.UserID,
	})
	if err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	res := &chargeSerializer.FilterCharges{
		Charges: make([]*chargeSerializer.Charge, 0),
	}
	for _, charge := range charges {
		res.Charges = append(res.Charges, &chargeSerializer.Charge{
			Status:       entity.TranslateChargeStatus(charge.Status),
			VPNType:      entity.VPNTypeString(charge.VPNType),
			Country:      charge.Country,
			AdminComment: charge.AdminComment,
			Detail:       charge.Detail,
			InboundID:    charge.Detail,
			UserID:       charge.UserID,
			PackageID:    charge.PackageID,
		})
	}
	return c.JSON(http.StatusOK, res)
}
