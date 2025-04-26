package inbound

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	inboundDto "momo/dto/repository/inbound"
	"momo/entity"
	momoError "momo/pkg/error"
)

func (i *Inbound) Create(inpt *inboundDto.CreateInbound) (*entity.Inbound, error) {
	inbound := &entity.Inbound{}
	err := i.db.Conn().QueryRow(`
	INSERT INTO inbounds (protocol, domain, vpn_type, port, user_id, tag, is_active, start, end, is_block)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id, protocol, is_active, domain, port, user_id, tag, is_block, start, end
	`, inpt.Protocol, inpt.Domain, entity.VPNTypeString(inpt.VPNType), inpt.Port, inpt.UserID, inpt.Tag, inpt.IsActive, inpt.Start, inpt.End, inpt.IsBlock).Scan(
		&inbound.ID,
		&inbound.Protocol,
		&inbound.IsActive,
		&inbound.Domain,
		&inbound.Port,
		&inbound.UserID,
		&inbound.Tag,
		&inbound.IsBlock,
		&inbound.Start,
		&inbound.End,
	)
	if err != nil {
		return nil, momoError.Errorf("somoething went wrong to save inbound error: %v", err)
	}
	inbound.VPNType = inpt.VPNType
	return inbound, nil
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

func (i *Inbound) DeleteAll() error {
	sql := "DELETE FROM inbounds"
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

func (i *Inbound) changeStatus(id int, state bool) error {
	sql := fmt.Sprintf("UPDATE inbounds SET is_active = %v WHERE id = %v", state, id)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.DebuggingErrorf("something bad has happend the error was %v", err)
	}
	return nil
}

func (i *Inbound) Active(id int) error {
	return i.changeStatus(id, true)
}

func (i *Inbound) DeActive(id int) error {
	return i.changeStatus(id, false)
}

func (i *Inbound) Filter(inpt *inboundDto.FilterInbound) ([]*entity.Inbound, error) {
	query := i.makeQueryFilter(inpt)

	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return []*entity.Inbound{}, momoError.DebuggingErrorf("error has occured err: %v", err)
	}

	inbounds := make([]*entity.Inbound, 0)

	for rows.Next() {
		inbound, err := i.scan(rows)
		if err != nil {
			return []*entity.Inbound{}, err
		}
		inbounds = append(inbounds, inbound)
	}
	return inbounds, nil
}

func (i *Inbound) makeQueryFilter(inpt *inboundDto.FilterInbound) string {
	sql := "SELECT * FROM inbounds"

	v := reflect.ValueOf(*inpt)
	t := reflect.TypeOf(*inpt)

	subSQLs := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.String && value.String() != "" && field.Name == "Protocol" {
			subSQLs = append(subSQLs, fmt.Sprintf("protocol = '%s'", value.String()))
		} else if field.Name == "IsActive" && !value.IsNil() && value.Elem().Kind() == reflect.Bool {
			subSQLs = append(subSQLs, fmt.Sprintf("is_active = %v", value.Elem().Bool()))
		} else if value.Kind() == reflect.String && value.String() != "" && field.Name == "Domain" {
			subSQLs = append(subSQLs, fmt.Sprintf("domain = '%s'", value.String()))
		} else if value.Kind() == reflect.String && value.String() != "" && field.Name == "Port" {
			subSQLs = append(subSQLs, fmt.Sprintf("port = '%s'", value.String()))
		} else if value.Kind() == reflect.String && value.String() != "" && field.Name == "VPNType" {
			subSQLs = append(subSQLs, fmt.Sprintf("vpn_type = '%s'", entity.VPNTypeString(int(value.Int()))))
		} else if value.Kind() == reflect.String && value.String() != "" && field.Name == "UserID" {
			subSQLs = append(subSQLs, fmt.Sprintf("user_id = '%s'", value.String()))
		}

	}
	subQuery := strings.Join(subSQLs, " AND ")
	sql += fmt.Sprintf(" WHERE %s", subQuery)
	return sql
}

func (i *Inbound) RetriveFaultyInbounds() ([]*entity.Inbound, error) {
	query := "SELECT * FROM inbounds WHERE (end < ?) OR (is_block = true AND is_active = true) OR (is_block = false AND start >= ? AND ? <= end AND is_active = false)"
	now := time.Now()
	rows, err := i.db.Conn().Query(query, now, now, now)

	inbounds := make([]*entity.Inbound, 0)
	if err != nil {
		return inbounds, momoError.DebuggingErrorf("something wrong has happend the problem was %v", err)
	}

	for rows.Next() {
		inbound, err := i.scan(rows)
		if err != nil {
			return inbounds, err
		}
		inbounds = append(inbounds, inbound)
	}

	return inbounds, nil
}

func (i *Inbound) CountingUsedPortEachHost() ([]*struct {
	domain string
	count  uint16
}, error,
) {
	query := "SELECT domain, COUNT(*) FROM inbounds WHERE is_active = true GROUP BY domain"
	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.DebuggingErrorf("the problem has happend that was %v", err)
	}
	var summeries []*struct {
		domain string
		count  uint16
	}
	for rows.Next() {
		var data struct {
			domain string
			count  uint16
		}

		if err := rows.Scan(&data.domain, &data.count); err != nil {
			return nil, momoError.DebuggingErrorf("the problem has happend that was in %v", err)
		}
		summeries = append(summeries, &data)
	}
	return summeries, nil
}

func (i *Inbound) scan(rows *sql.Rows) (*entity.Inbound, error) {
	inbound := &entity.Inbound{}
	var createdAt, updatedAt interface{}
	var vpnType string
	err := rows.Scan(
		&inbound.ID,
		&inbound.Protocol,
		&inbound.IsActive,
		&inbound.Domain,
		&vpnType,
		&inbound.Port,
		&inbound.UserID,
		&inbound.Tag,
		&inbound.IsBlock,
		&inbound.Start,
		&inbound.End,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return inbound, momoError.DebuggingErrorf("error has occured err: %v", err)
	}
	inbound.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
	return inbound, nil
}
