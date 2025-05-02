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
	scope := "inboundRepository.create"

	inbound := &entity.Inbound{}
	err := i.db.Conn().QueryRow(`
	INSERT INTO inbounds (protocol, domain, vpn_type, port, user_id, tag, is_active, start, end, is_block, is_assigned, is_notified)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id, protocol, is_active, domain, port, user_id, tag, is_block, start, end, is_notified, is_assigned
	`, inpt.Protocol,
		inpt.Domain, entity.VPNTypeString(inpt.VPNType), inpt.Port, inpt.UserID, inpt.Tag, inpt.IsActive, inpt.Start, inpt.End, inpt.IsBlock, inpt.IsAssigned, inpt.IsNotified).Scan(
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
		&inbound.IsNotified,
		&inbound.IsAssigned,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", *inpt)
	}
	inbound.VPNType = inpt.VPNType
	return inbound, nil
}

func (i *Inbound) FindInboundByID(id int) (*entity.Inbound, error) {
	scope := "inboundRepository.FindInboundByID"

	var createdAt, updatedAt interface{}
	var inbound *entity.Inbound = &entity.Inbound{}
	var vpnType string
	s := fmt.Sprintf("SELECT * FROM inbounds WHERE id=%v LIMIT 1", id)
	err := i.db.Conn().QueryRow(s).Scan(
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
		&inbound.IsNotified,
		&inbound.IsAssigned,
		&createdAt,
		&updatedAt,
	)

	if err == nil {
		inbound.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
		return inbound, nil
	}
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).Errorf("the id is %d", id)
	}
	return nil, momoError.Wrap(err).Scope(scope).Errorf("the id is %d", id)
}

func (i *Inbound) Delete(id int) error {
	scope := "inboundRepository.FindInboundByID"

	sql := fmt.Sprintf("DELETE FROM inbounds WHERE id=%v", id)
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the id is %d", id)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the id is %d", id)
	}

	if rowsAffected == 0 {
		return momoError.Scope(scope).Errorf("no row is affected")
	}
	return nil
}

func (i *Inbound) DeleteAll() error {
	scope := "inboundRepository.DeleteAll"

	sql := "DELETE FROM inbounds"
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

func (i *Inbound) GetListOfPortsByDomain() ([]struct {
	Domain string
	Ports  []string
}, error,
) {
	scope := "inboundRepository.GetListOfPortsByDomain"

	sql := "SELECT domain, GROUP_CONCAT(port) AS ports FROM inbounds GROUP BY domain"
	rows, err := i.db.Conn().Query(sql)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	defer rows.Close()

	var res []struct {
		Domain string
		Ports  []string
	}

	for rows.Next() {
		var domain, portsStr string
		if err := rows.Scan(&domain, &portsStr); err != nil {
			return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
		}

		res = append(res, struct {
			Domain string
			Ports  []string
		}{
			Domain: domain,
			Ports:  strings.Split(portsStr, ","),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	return res, nil
}

func (i *Inbound) changeStatus(id int, state bool) error {
	scope := "inboundRepository.changeStatus"

	sql := fmt.Sprintf("UPDATE inbounds SET is_active = %v WHERE id = %v", state, id)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).DebuggingErrorf("the id is %d and the state is %v", id, state)
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
	scope := "inboundRepository.Filter"

	query := i.makeQueryFilter(inpt)

	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingErrorf("the input is %+v", *inpt)
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
	scope := "inboundRepository.RetriveFaultyInbounds"

	query := "SELECT * FROM inbounds WHERE (is_active = true AND end < ?) OR (is_block = true AND is_active = true) OR (is_block = false AND start >= ? AND ? <= end AND is_active = false)"
	now := time.Now()
	rows, err := i.db.Conn().Query(query, now, now, now)

	inbounds := make([]*entity.Inbound, 0)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	for rows.Next() {
		inbound, err := i.scan(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
		}
		inbounds = append(inbounds, inbound)
	}

	return inbounds, nil
}

func (i *Inbound) FindInboundIsNotAssigned() ([]*entity.Inbound, error) {
	scope := "inboundRepository.FindInboundIsNotAssigned"

	query := "SELECT * FROM inbounds WHERE is_assigned = false"
	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	inbounds := make([]*entity.Inbound, 0)
	for rows.Next() {
		inbound, err := i.scan(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
		}
		inbounds = append(inbounds, inbound)
	}
	return inbounds, nil
}

func (i *Inbound) UpdateDomainPort(id int, domain string, port string) error {
	scope := "inboundRepository.UpdateDomainPort"

	sql := fmt.Sprintf(
		"UPDATE inbounds SET domain = '%s', port = '%s' WHERE id = %v",
		domain,
		port,
		id,
	)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the id ids %d and the domain is %s and the port is %s", id, domain, port)
	}
	return nil
}

func (i *Inbound) scan(rows *sql.Rows) (*entity.Inbound, error) {
	scope := "inboundRepository.scan"

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
		&inbound.IsNotified,
		&inbound.IsAssigned,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	inbound.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
	return inbound, nil
}
