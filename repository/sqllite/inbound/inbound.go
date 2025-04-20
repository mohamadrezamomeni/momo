package inbound

import (
	"fmt"
	"reflect"
	"strings"

	"momo/entity"
	momoError "momo/pkg/error"
	"momo/repository/sqllite/inbound/dto"
)

func (i *Inbound) Create(inpt *dto.CreateInbound) (*entity.Inbound, error) {
	var (
		id          int
		protocol    string
		domain      string
		vpnType     string
		port        string
		userID      string
		tag         string
		isAvailable bool
	)
	err := i.db.Conn().QueryRow(`
	INSERT INTO inbounds (protocol, domain, vpn_type, port, user_id, tag, is_available)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING id, protocol, is_available, domain, vpn_type, port, user_id, tag
	`, inpt.Protocol, inpt.Domain, inpt.VPNType, inpt.Port, inpt.UserID, inpt.Tag, inpt.IsAvailable).Scan(
		&id,
		&protocol,
		&isAvailable,
		&domain,
		&vpnType,
		&port,
		&userID,
		&tag,
	)
	if err != nil {
		return &entity.Inbound{}, momoError.Errorf("somoething went wrong to save inbound error: %v", err)
	}

	return &entity.Inbound{
		ID:          id,
		Protocol:    protocol,
		Domain:      domain,
		VPNType:     vpnType,
		Port:        port,
		UserID:      userID,
		Tag:         tag,
		IsAvailable: isAvailable,
	}, nil
}

func (i *Inbound) Delete(id int) error {
	sql := fmt.Sprintf("DELETE FROM inbounds WHERE id=%v", id)
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

func (i *Inbound) changeStatus(id int, state bool) error {
	sql := fmt.Sprintf("UPDATE inbounds SET is_available = %v WHERE id = %v", state, id)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.DebuggingErrorf("something bad has happend the error was %v", err)
	}
	return nil
}

func (i *Inbound) MakeAvailable(id int) error {
	return i.changeStatus(id, true)
}

func (i *Inbound) MakeNotAvailable(id int) error {
	return i.changeStatus(id, false)
}

func (i *Inbound) Filter(inpt *dto.FilterInbound) ([]*entity.Inbound, error) {
	query := i.makeQueryFilter(inpt)

	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return []*entity.Inbound{}, momoError.DebuggingErrorf("error has occured err: %v", err)
	}

	inbounds := make([]*entity.Inbound, 0)

	for rows.Next() {
		var (
			protocol, port, domain, vpnTpe, userID, tag string
			isAvailable                                 bool
			createdAt, updatedAt                        interface{}
			id                                          int
		)

		err := rows.Scan(&id, &protocol, &isAvailable, &domain, &vpnTpe, &port, &userID, &tag, &createdAt, &updatedAt)
		if err != nil {
			return []*entity.Inbound{}, momoError.DebuggingErrorf("error has occured err: %v", err)
		}

		inbounds = append(inbounds, &entity.Inbound{
			ID:          id,
			Domain:      domain,
			Protocol:    protocol,
			VPNType:     vpnTpe,
			Tag:         tag,
			IsAvailable: isAvailable,
		})
	}
	return inbounds, nil
}

func (i *Inbound) makeQueryFilter(inpt *dto.FilterInbound) string {
	sql := "SELECT * FROM inbounds"

	v := reflect.ValueOf(*inpt)
	t := reflect.TypeOf(*inpt)

	subSQLs := make([]string, 0)

	convertKeysToColumns := func(k string) string {
		switch k {
		case "Protocol":
			return "protocol"
		case "IsAvailable":
			return "is_available"
		case "Domain":
			return "domain"
		case "VPNType":
			return "vpn_type"
		case "Port":
			return "port"
		case "UserID":
			return "user_id"
		default:
			return ""
		}
	}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.String && value.String() != "" {
			subSQLs = append(subSQLs, fmt.Sprintf("%s = '%s'", convertKeysToColumns(field.Name), value.String()))
		} else if value.Kind() == reflect.Pointer && !value.IsNil() && value.Elem().Kind() == reflect.Bool {
			subSQLs = append(subSQLs, fmt.Sprintf("%s = %v", convertKeysToColumns(field.Name), true))
		}
	}
	subQuery := strings.Join(subSQLs, " AND ")
	sql += fmt.Sprintf(" WHERE %s", subQuery)
	return sql
}
