package vpnmanager

import (
	"fmt"
	"strings"

	vpnManagerDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (v *VPN) Create(inpt *vpnManagerDto.AddVPN) (*entity.VPN, error) {
	scope := "vpnRepository.Create"

	var vpn *entity.VPN = &entity.VPN{
		Domain:    inpt.Domain,
		IsActive:  inpt.IsActive,
		ApiPort:   inpt.ApiPort,
		VPNType:   inpt.VPNType,
		UserCount: inpt.UserCount,
		Country:   inpt.Country,
		VPNStatus: inpt.VPNStatus,
	}
	err := v.db.Conn().QueryRow(`
	INSERT INTO vpns (domain, is_active, api_port, vpn_type, user_count, country, start_port, end_port, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id
	`, inpt.Domain,
		inpt.IsActive,
		inpt.ApiPort,
		entity.VPNTypeString(vpn.VPNType),
		inpt.UserCount,
		inpt.Country,
		inpt.StartPort,
		inpt.EndPort,
		entity.VPNStatusString(inpt.VPNStatus),
	).Scan(
		&vpn.ID,
	)
	if err == nil {
		return vpn, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
}

func (i *VPN) Delete(id int) error {
	scope := "vpnRepository.Delete"

	sql := fmt.Sprintf("DELETE FROM vpns WHERE id = %v", id)
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}

	if rowsAffected == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}
	return nil
}

func (i *VPN) DeleteAll() error {
	scope := "vpnRepository.DeleteAll"

	sql := "DELETE FROM vpns"
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (v *VPN) ActiveVPN(id int) error {
	return v.updateActivationVPN(id, true)
}

func (v *VPN) DeactiveVPN(id int) error {
	return v.updateActivationVPN(id, false)
}

func (v *VPN) updateActivationVPN(id int, status bool) error {
	scope := "vpnRepository.updateActivationVPN"

	sql := fmt.Sprintf("UPDATE vpns SET is_active = %v WHERE id = %v", status, id)

	_, err := v.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id, status).DebuggingError()
	}
	return nil
}

func (v *VPN) Filter(inpt *vpnManagerDto.FilterVPNs) ([]*entity.VPN, error) {
	scope := "vpnRepository.Filter"

	query := v.makeSQlFilter(inpt)
	rows, err := v.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(inpt).DebuggingError()
	}

	vpns := make([]*entity.VPN, 0)

	for rows.Next() {
		vpn := &entity.VPN{}
		var createdAt, updatedAt interface{}
		var vpnType string
		var vpnStatusLabel string
		err := rows.Scan(
			&vpn.ID,
			&vpn.Domain,
			&vpn.IsActive,
			&vpn.ApiPort,
			&vpnType,
			&vpn.UserCount,
			&vpn.Country,
			&vpn.StartPort,
			&vpn.EndPort,
			&vpnStatusLabel,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).Input(inpt).DebuggingError()
		}

		vpn.VPNStatus = entity.ConvertVPNStatusLabelToVPNStatus(vpnStatusLabel)
		vpn.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
		vpns = append(vpns, vpn)
	}

	return vpns, nil
}

func (v *VPN) makeSQlFilter(inpt *vpnManagerDto.FilterVPNs) string {
	sql := "SELECT * FROM vpns"
	subQueries := make([]string, 0)

	if inpt.IsActive != nil {
		subQueries = append(subQueries, fmt.Sprintf("is_active = %v", *inpt.IsActive))
	}

	if inpt.Domain != "" {
		subQueries = append(subQueries, fmt.Sprintf("domain LIKE  '%%%s%%'", inpt.Domain))
	}

	if inpt.VPNTypes != nil && len(inpt.VPNTypes) > 0 {
		subQueries = append(subQueries, fmt.Sprintf("vpn_type IN ('%s')", strings.Join(v.convertVPNTypesToVPNstrings(inpt.VPNTypes), ",")))
	}

	if inpt.Coountries != nil && len(inpt.Coountries) > 0 {
		subQueries = append(subQueries, fmt.Sprintf("country IN ('%s')", strings.Join(inpt.Coountries, "', '")))
	}

	if inpt.VPNStatuses != nil && len(inpt.VPNStatuses) > 0 {
		subQueries = append(
			subQueries,
			fmt.Sprintf("status IN ('%s')", strings.Join(
				entity.ConvertVPNStatusesToVPNStatusLabels(inpt.VPNStatuses),
				"', '",
			)),
		)
	}

	joinSQL := strings.Join(subQueries, " AND ")

	if len(joinSQL) > 0 {
		sql += fmt.Sprintf(" WHERE %s", joinSQL)
	}
	return sql
}

func (v *VPN) GroupAvailbleVPNsByCountry() ([]string, error) {
	scope := "vpnRepository.GroupAvailbleVPNsByCountry"

	query := "SELECT country FROM vpns WHERE is_active = true GROUP BY country"

	rows, err := v.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	countries := make([]string, 0)

	for rows.Next() {
		country := ""
		err := rows.Scan(
			&country,
		)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
		}

		countries = append(countries, country)
	}
	return countries, nil
}

func (v *VPN) convertVPNTypesToVPNstrings(vpnTypes []entity.VPNType) []string {
	res := make([]string, 0)
	for _, vpnType := range vpnTypes {
		res = append(res, entity.VPNTypeString(vpnType))
	}
	return res
}
