package vpnmanager

import (
	"fmt"
	"reflect"
	"strings"

	vpnManagerDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_manager"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *VPN) Create(inpt *vpnManagerDto.AddVPN) (*entity.VPN, error) {
	scope := "vpnRepository.Create"

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
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", *inpt)
	}
	return vpn, nil
}

func (i *VPN) Delete(id int) error {
	scope := "vpnRepository.Delete"

	sql := fmt.Sprintf("DELETE FROM vpns WHERE id = %v", id)
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the id is %d", id)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the id is %d", id)
	}

	if rowsAffected == 0 {
		return momoError.Wrap(err).Scope(scope).Errorf("no row is affected", id)
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
		return momoError.Wrap(err).Scope(scope).Errorf("the id is %d and the status is %v", id, status)
	}
	return nil
}

func (v *VPN) Filter(inpt *vpnManagerDto.FilterVPNs) ([]*entity.VPN, error) {
	scope := "vpnRepository.Filter"

	query := v.makeSQlFilter(inpt)
	rows, err := v.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", *inpt)
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
			return nil, momoError.Wrap(err).Scope(scope).Errorf("error to scand data, the input is %+v", *inpt)
		}

		vpnTypeEnum := entity.ConvertStringVPNTypeToEnum(vpnType)
		if vpnTypeEnum == entity.UknownVPNType {
			return nil, momoError.Wrap(err).Scope(scope).Errorf("error to convert vpnType, the input is %+v", *inpt)
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
