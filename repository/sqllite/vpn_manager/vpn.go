package vpnmanager

import (
	"fmt"
	"reflect"
	"strings"

	vpnManagerDto "momo/dto/repository/vpn_manager"
	"momo/entity"
	momoError "momo/pkg/error"
)

func (v *VPN) Create(inpt *vpnManagerDto.Add_VPN) (*entity.VPN, error) {
	var vpn *entity.VPN = &entity.VPN{
		Domain:    inpt.Domain,
		IsActive:  inpt.IsActive,
		ApiPort:   inpt.ApiPort,
		VPNType:   inpt.VPNType,
		UserCount: inpt.UserCount,
	}
	err := v.db.Conn().QueryRow(`
	INSERT INTO vpns (domain, is_active, api_port, vpn_type, user_count)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id
	`, inpt.Domain, inpt.IsActive, inpt.ApiPort, entity.VPNTypeString(vpn.VPNType), inpt.UserCount).Scan(
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

func (i *VPN) DeleteAll() error {
	sql := "DELETE FROM vpns"
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Errorf("something went wrong to delete record follow error, the error was %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Errorf("something went wrong to delete record follow error, the error was %v", err)
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

func (v *VPN) Filter(inpt *vpnManagerDto.FilterVPNs) ([]*entity.VPN, error) {
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
			&vpnType,
			&vpn.UserCount,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return []*entity.VPN{}, momoError.DebuggingErrorf("error has occured err: %v", err)
		}
		vpnTypeEnum := entity.ConvertStringVPNTypeToEnum(vpnType)
		if vpnTypeEnum == entity.UknownVPNType {
			return []*entity.VPN{}, momoError.DebuggingErrorf("fail to convert vpnType %s to enum", vpnType)
		}
		vpn.VPNType = vpnTypeEnum
		vpns = append(vpns, vpn)
	}

	return vpns, nil
}

func (v *VPN) makeSQlFilter(inpt *vpnManagerDto.FilterVPNs) string {
	sql := "SELECT * FROM vpns"

	val := reflect.ValueOf(*inpt)
	t := reflect.TypeOf(*inpt)
	subQueries := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		v := val.Field(i)
		field := t.Field(i)
		if v.Kind() == reflect.Pointer && !v.IsNil() && v.Elem().Kind() == reflect.Bool && field.Name == "IsActive" {
			subQueries = append(subQueries, fmt.Sprintf("is_active = %v", v.Elem().Bool()))
		} else if v.Kind() == reflect.String && v.String() != "" && field.Name == "Domain" {
			subQueries = append(subQueries, fmt.Sprintf("domain = '%s'", v.String()))
		} else if field.Name == "VPNType" && v.Int() != 0 {
			subQueries = append(subQueries, fmt.Sprintf("vpn_type = '%s'", entity.VPNTypeString(int(v.Int()))))
		}
	}
	joinSQL := strings.Join(subQueries, " AND ")

	if len(joinSQL) > 0 {
		sql += fmt.Sprintf(" WHERE %s", joinSQL)
	}
	return sql
}
