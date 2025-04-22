package vpnmanager

import (
	"fmt"
	"reflect"
	"strings"

	"momo/entity"
	momoError "momo/pkg/error"
	"momo/repository/sqllite/vpn_manager/dto"
)

func (v *VPN) Create(inpt *dto.Add_VPN) (*entity.VPN, error) {
	var vpn *entity.VPN = &entity.VPN{
		Domain:         inpt.Domain,
		IsActive:       inpt.IsActive,
		ApiPort:        inpt.ApiPort,
		StartRangePort: inpt.StartRangePort,
		EndRangePort:   inpt.EndRangePort,
		VPNType:        inpt.VPNType,
		UserCount:      inpt.UserCount,
	}
	err := v.db.Conn().QueryRow(`
	INSERT INTO vpns (domain, is_active, api_port, start_range_port, end_range_port, vpn_type, user_count)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING id
	`, inpt.Domain, inpt.IsActive, inpt.ApiPort, inpt.StartRangePort, inpt.EndRangePort, vpn.VPNTypeString(), inpt.UserCount).Scan(
		&vpn.ID,
	)
	if err != nil {
		return vpn, momoError.DebuggingErrorf("something wrong has happend the problem was %v", err)
	}
	return vpn, nil
}

func (i *VPN) Delete(id int) error {
	sql := fmt.Sprintf("DELETE FROM vpns WHERE id = %v", id)
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Errorf("something went wrong to delete record follow error, the error was %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Errorf("something went wrong to delete record follow error, the error was %v", err)
	}

	if rowsAffected == 0 {
		return momoError.Error("None of the records have been affected.")
	}
	return nil
}

func (v *VPN) activeVPN(id int) error {
	return v.updateActivationVPN(id, true)
}

func (v *VPN) deactiveVPN(id int) error {
	return v.updateActivationVPN(id, false)
}

func (v *VPN) updateActivationVPN(id int, status bool) error {
	sql := fmt.Sprintf("UPDATE vpns SET is_active = %v WHERE id = %v", status, id)

	_, err := v.db.Conn().Exec(sql)
	if err != nil {
		return momoError.DebuggingErrorf("something bad has happend the error was %v", err)
	}
	return nil
}

func (v *VPN) Filter(inpt *dto.FilterVPNs) ([]*entity.VPN, error) {
	query := v.makeSQlFilter(inpt)

	rows, err := v.db.Conn().Query(query)
	if err != nil {
		return []*entity.VPN{}, momoError.DebuggingErrorf("error has occured err: %v", err)
	}

	vpns := make([]*entity.VPN, 0)

	for rows.Next() {
		vpn := &entity.VPN{}
		var createdAt, updatedAt interface{}
		var vpnType string
		err := rows.Scan(
			&vpn.ID,
			&vpn.Domain,
			&vpn.IsActive,
			&vpn.ApiPort,
			&vpn.StartRangePort,
			&vpn.EndRangePort,
			&vpnType,
			&vpn.UserCount,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return []*entity.VPN{}, momoError.DebuggingErrorf("error has occured err: %v", err)
		}
		vpnTypeEnum, err := entity.ConvertStringVPNTypeToEnum(vpnType)
		if err != nil {
			return []*entity.VPN{}, momoError.DebuggingErrorf("fail to convert vpnType %s to enum", vpnType)
		}
		vpn.VPNType = vpnTypeEnum
		vpns = append(vpns, vpn)
	}

	return vpns, nil
}

func (v *VPN) makeSQlFilter(inpt *dto.FilterVPNs) string {
	sql := "SELECT * FROM vpns"

	val := reflect.ValueOf(*inpt)
	t := reflect.TypeOf(*inpt)
	subQueries := make([]string, 0)

	convertKeysToColumns := func(k string) string {
		switch k {
		case "Domain":
			return "domain"
		case "IsActive":
			return "is_active"
		case "VPNType":
			return "vpn_type"
		}
		return ""
	}

	for i := 0; i < t.NumField(); i++ {
		v := val.Field(i)
		field := t.Field(i)
		if v.Kind() == reflect.Pointer && !v.IsNil() && v.Elem().Kind() == reflect.Bool {
			subQueries = append(subQueries, fmt.Sprintf("%s = %v", convertKeysToColumns(field.Name), v.Elem().Bool()))
		} else if v.Kind() == reflect.String && v.String() != "" {
			subQueries = append(subQueries, fmt.Sprintf("%s = '%s'", convertKeysToColumns(field.Name), v.String()))
		}
	}
	joinSQL := strings.Join(subQueries, " AND ")

	if len(joinSQL) > 0 {
		sql += fmt.Sprintf(" WHERE %s", joinSQL)
	}

	return sql
}
