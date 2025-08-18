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
	INSERT INTO vpn_package (price_tilte, price, days, months, traffic_limit, traffic_limit_title, is_active, tier)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id, price_tilte, price, days, months, traffic_limit, traffic_limit_title, is_active, tier
`,
		inpt.PriceTitle,
		inpt.Price,
		inpt.Days,
		inpt.Months,
		inpt.TrafficLimit,
		inpt.TrafficLimitTitle,
		inpt.IsActive,
		inpt.Tier,
	).Scan(
		&vpnPackage.ID,
		&vpnPackage.PriceTitle,
		&vpnPackage.Price,
		&vpnPackage.Days,
		&vpnPackage.Months,
		&vpnPackage.TrafficLimit,
		&vpnPackage.TrafficLimitTitle,
		&vpnPackage.IsActive,
		&vpnPackage.Tier,
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
		&vpnPackage.Tier,
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

func (vp *VPNPackage) DeleteAll() error {
	scope := "vpnPackageRepository.DeleteAll"

	sql := "DELETE FROM vpn_package"
	res, err := vp.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (vp *VPNPackage) Delete(id string) error {
	scope := "vpnPackageRepository.Delete"

	sql := fmt.Sprintf("DELETE FROM vpn_package WHERE id = '%s'", id)
	res, err := vp.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (vp *VPNPackage) Filter(inpt *vpnPackageRepositoryDto.FilterVPNPackage) ([]*entity.VPNPackage, error) {
	scope := "vpnpackage.repository.Filter"

	query := vp.makeQueryFilter(inpt)

	rows, err := vp.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Input(inpt).Scope(scope).DebuggingError()
	}

	vpnPackages := make([]*entity.VPNPackage, 0)

	for rows.Next() {
		pkg, err := vp.scan(rows)
		if err != nil {
			return nil, err
		}
		vpnPackages = append(vpnPackages, pkg)
	}
	return vpnPackages, nil
}

func (vp *VPNPackage) makeQueryFilter(inpt *vpnPackageRepositoryDto.FilterVPNPackage) string {
	subQueries := []string{}

	if inpt.IsActive != nil {
		subQueries = append(subQueries, fmt.Sprintf("is_active = %v", *inpt.IsActive))
	}

	if inpt.Tiers != nil && len(inpt.Tiers) > 0 {
		subQueries = append(
			subQueries,
			fmt.Sprintf("tier IN ('%s')", strings.Join(inpt.Tiers, "', '")),
		)
	}
	sql := "SELECT * FROM vpn_package"
	if len(subQueries) > 0 {
		sql += fmt.Sprintf(" WHERE %s", strings.Join(subQueries, " AND "))
	}
	return sql
}

func (vp *VPNPackage) scan(rows *sql.Rows) (*entity.VPNPackage, error) {
	scope := "vpnPackage.repositroy.scan"

	vpnPackage := &entity.VPNPackage{}
	var createdAt interface{}
	err := rows.Scan(
		&vpnPackage.ID,
		&vpnPackage.PriceTitle,
		&vpnPackage.Price,
		&vpnPackage.Days,
		&vpnPackage.Months,
		&vpnPackage.TrafficLimit,
		&vpnPackage.TrafficLimitTitle,
		&vpnPackage.IsActive,
		&vpnPackage.Tier,
		&createdAt,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}
	return vpnPackage, nil
}
