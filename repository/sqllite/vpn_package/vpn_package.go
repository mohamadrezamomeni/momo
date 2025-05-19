package vpnpackage

import (
	vpnPackageRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (vp *VPNPackage) Create(inpt *vpnPackageRepositoryDto.CreateVPNPackage) (*entity.VPNPackage, error) {
	scope := "userRepository.Create"

	vpnPackage := &entity.VPNPackage{}
	err := vp.db.Conn().QueryRow(`
	INSERT INTO vpn_package (price_tilte, price, days, months, traffic_limit, traffic_limit_title, is_active)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING id, price_tilte, price, days, months, traffic_limit, traffic_limit_title, is_active
`,
		inpt.PriceTitle,
		inpt.Price,
		inpt.Days,
		inpt.Months,
		inpt.TrafficLimit,
		inpt.TrafficLimitTitle,
		inpt.IsActive,
	).Scan(
		&vpnPackage.ID,
		&vpnPackage.PriceTitle,
		&vpnPackage.Price,
		&vpnPackage.Days,
		&vpnPackage.Months,
		&vpnPackage.TrafficLimit,
		&vpnPackage.TrafficLimitTitle,
		&vpnPackage.IsActive,
	)
	if err == nil {
		return vpnPackage, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
}
