package vpnpackage

import (
	"database/sql"
	"fmt"
	"strings"

	vpnPackageRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (vp *VPNPackage) Create(inpt *vpnPackageRepositoryDto.CreateVPNPackage) (*entity.VPNPackage, error) {
	scope := "vpnPackageRepository.Create"

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

func (vp *VPNPackage) Update(id string, update *vpnPackageRepositoryDto.UpdateVPNPackage) error {
	scope := "vpnPackageRepository.Update"
	subUpdates := []string{}

	if update.IsActive != nil {
		subUpdates = append(subUpdates, fmt.Sprintf("is_active = %v", *update.IsActive))
	}

	sql := fmt.Sprintf(
		"UPDATE vpn_package SET %s WHERE id = %s",
		strings.Join(subUpdates, ", "),
		id,
	)

	result, err := vp.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id, update).ErrorWrite()
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id, update).ErrorWrite()
	}
	return nil
}

func (vp *VPNPackage) FindVPNPackageByID(id string) (*entity.VPNPackage, error) {
	scope := "inboundRepository.FindVPNPackageByID"

	var createdAt interface{}
	vpnPackage := &entity.VPNPackage{}

	s := fmt.Sprintf("SELECT * FROM vpn_package WHERE id = %s LIMIT 1", id)
	err := vp.db.Conn().QueryRow(s).Scan(
		&vpnPackage.ID,
		&vpnPackage.PriceTitle,
		&vpnPackage.Price,
		&vpnPackage.Days,
		&vpnPackage.Months,
		&vpnPackage.TrafficLimit,
		&vpnPackage.TrafficLimitTitle,
		&vpnPackage.IsActive,
		&createdAt,
	)

	if err == nil {
		return vpnPackage, nil
	}
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).NotFound().Input(id).DebuggingError()
	}
	return nil, momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
}
