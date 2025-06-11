package vpnsource

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (vs *VPNSource) Upsert(country string, upsertdVPNSourceDto *vpnSourceRepositoryDto.UpsertVPNSourceDto) (*entity.VPNSource, error) {
	scope := "vpnSourceRepository.Upsert"

	vpnSource := &entity.VPNSource{}
	err := vs.db.Conn().QueryRow(`
	INSERT INTO vpn_source (country, english)
	VALUES (?, ?)
	ON CONFLICT(country) DO UPDATE SET
		english = excluded.english
	RETURNING country, english
	`, country,
		upsertdVPNSourceDto.English,
	).Scan(
		&vpnSource.Country,
		&vpnSource.English,
	)
	if err == nil {
		return vpnSource, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(upsertdVPNSourceDto).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(upsertdVPNSourceDto).UnExpected().Scope(scope).DebuggingError()
}

func (vs *VPNSource) Find(country string) (*entity.VPNSource, error) {
	scope := "vpnSourceRepository.find"

	vpnSource := &entity.VPNSource{}

	s := fmt.Sprintf("SELECT * FROM vpn_source WHERE country='%s' LIMIT 1", country)
	err := vs.db.Conn().QueryRow(s).Scan(
		&vpnSource.Country,
		&vpnSource.English,
	)
	if err == nil {
		return vpnSource, nil
	}
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).Input(country).NotFound().DebuggingError()
	}
	return nil, momoError.Wrap(err).Scope(scope).Input(country).UnExpected().DebuggingError()
}

func (vs *VPNSource) Filter(filterDto *vpnSourceRepositoryDto.FilterVPNSources) ([]*entity.VPNSource, error) {
	scope := "vpnsourceRepository.filter"
	query := vs.makeFilterQuery(filterDto)

	rows, err := vs.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	vpnSources := make([]*entity.VPNSource, 0)
	for rows.Next() {
		vpnsource, err := vs.scan(rows)
		if err != nil {
			return nil, err
		}
		vpnSources = append(vpnSources, vpnsource)
	}
	return vpnSources, nil
}

func (vs *VPNSource) makeFilterQuery(filterDto *vpnSourceRepositoryDto.FilterVPNSources) string {
	t := reflect.TypeOf(*filterDto)
	v := reflect.ValueOf(*filterDto)
	subQueries := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.Name == "Countries" && !value.IsNil() && value.Len() > 0 {
			subQueries = append(subQueries, fmt.Sprintf("country IN ('%s')", strings.Join(filterDto.Countries, "','")))
		}
	}
	sql := "SELECT * FROM vpn_source"
	if len(subQueries) > 0 {
		sql += fmt.Sprintf(" WHERE %s", strings.Join(subQueries, " AND "))
	}
	return sql
}

func (vp *VPNSource) scan(rows *sql.Rows) (*entity.VPNSource, error) {
	scope := "vpnPackage.repositroy.scan"

	vpnSource := &entity.VPNSource{}
	err := rows.Scan(
		&vpnSource.Country,
		&vpnSource.English,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}
	return vpnSource, nil
}

func (vs *VPNSource) DeleteAll() error {
	scope := "vpnSourceRepository.DeleteAll"

	sql := "DELETE FROM vpn_source"
	res, err := vs.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}
