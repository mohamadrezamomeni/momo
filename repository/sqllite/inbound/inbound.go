package inbound

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	inboundDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (i *Inbound) Create(inpt *inboundDto.CreateInbound) (*entity.Inbound, error) {
	scope := "inboundRepository.create"

	inbound := &entity.Inbound{}
	err := i.db.Conn().QueryRow(`
	INSERT INTO inbounds (protocol, domain, vpn_type, port, user_id, tag, is_active, start, end, is_block, is_assigned, is_notified, traffic_usage, traffic_limit, country)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id, protocol, is_active, domain, port, user_id, tag, is_block, start, end, is_notified, is_assigned, traffic_usage, traffic_limit, country
	`, inpt.Protocol,
		inpt.Domain,
		entity.VPNTypeString(inpt.VPNType),
		inpt.Port,
		inpt.UserID,
		inpt.Tag,
		inpt.IsActive,
		inpt.Start,
		inpt.End,
		inpt.IsBlock,
		inpt.IsAssigned,
		inpt.IsNotified,
		inpt.TrafficUsage,
		inpt.TrafficLimit,
		inpt.Country,
	).Scan(
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
		&inbound.TrafficUsage,
		&inbound.TrafficLimit,
		&inbound.Country,
	)

	if err == nil {
		inbound.VPNType = inpt.VPNType
		return inbound, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
}

func (i *Inbound) FindInboundByID(id string) (*entity.Inbound, error) {
	scope := "inboundRepository.FindInboundByID"

	var createdAt, updatedAt interface{}
	var inbound *entity.Inbound = &entity.Inbound{}
	var vpnType string
	var VPNID sql.NullInt64
	s := fmt.Sprintf("SELECT * FROM inbounds WHERE id=%s LIMIT 1", id)
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
		&inbound.TrafficUsage,
		&inbound.TrafficLimit,
		&inbound.Country,
		&createdAt,
		&updatedAt,
		&VPNID,
	)

	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).NotFound().Input(id).DebuggingError()
	}
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
	}
	inbound.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
	if VPNID.Valid {
		inbound.VPNID = strconv.Itoa(int(VPNID.Int64))
	}
	return inbound, nil
}

func (i *Inbound) Delete(id int) error {
	scope := "inboundRepository.Delete"

	sql := fmt.Sprintf("DELETE FROM inbounds WHERE id=%v", id)
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}

	if rowsAffected == 0 {
		return momoError.Scope(scope).Input(id).NotFound().DebuggingError()
	}
	return nil
}

func (i *Inbound) DeleteAll() error {
	scope := "inboundRepository.DeleteAll"

	sql := "DELETE FROM inbounds"
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return nil
}

func (i *Inbound) GetListOfPortsByDomain() ([]struct {
	Domain string
	Ports  []string
}, error,
) {
	scope := "inboundRepository.GetListOfPortsByDomain"

	sql := "SELECT domain, GROUP_CONCAT(port) AS ports FROM inbounds WHERE is_assigned = true GROUP BY domain"
	rows, err := i.db.Conn().Query(sql)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	defer rows.Close()

	var res []struct {
		Domain string
		Ports  []string
	}

	for rows.Next() {
		var domain, portsStr string
		if err := rows.Scan(&domain, &portsStr); err != nil {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
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
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return res, nil
}

func (i *Inbound) changeStatus(id string, state bool) error {
	scope := "inboundRepository.changeStatus"

	sql := fmt.Sprintf("UPDATE inbounds SET is_active = %v WHERE id = %v", state, id)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id, state).DebuggingError()
	}
	return nil
}

func (i *Inbound) Active(id string) error {
	return i.changeStatus(id, true)
}

func (i *Inbound) DeActive(id string) error {
	return i.changeStatus(id, false)
}

func (i *Inbound) Filter(inpt *inboundDto.FilterInbound) ([]*entity.Inbound, error) {
	scope := "inboundRepository.Filter"

	query := i.makeQueryFilter(inpt)
	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Input(inpt).Scope(scope).DebuggingError()
	}

	return i.getInboundsFromRows(rows)
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
	if len(subSQLs) > 0 {
		subQuery := strings.Join(subSQLs, " AND ")
		sql += fmt.Sprintf(" WHERE %s", subQuery)
	}
	return sql
}

func (i *Inbound) RetriveActiveInboundBlocked() ([]*entity.Inbound, error) {
	scope := "inboundRepository.RetriveActiveInboundBlocked"

	query := "SELECT * FROM inbounds WHERE  is_block = true AND is_active = true"
	now := time.Now()
	rows, err := i.db.Conn().Query(query, now, now, now)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return i.getInboundsFromRows(rows)
}

func (i *Inbound) RetriveActiveInboundExpired() ([]*entity.Inbound, error) {
	scope := "inboundRepository.RetriveActiveInboundExpired"

	query := "SELECT * FROM inbounds WHERE is_active = true AND end < ?"
	now := time.Now()
	rows, err := i.db.Conn().Query(query, now)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return i.getInboundsFromRows(rows)
}

func (i *Inbound) RetriveFinishedInbounds() ([]*entity.Inbound, error) {
	scope := "inboundRepository.RetriveDeactivedInbounds"

	query := "SELECT * FROM inbounds WHERE (is_active = false AND is_block = false) AND (traffic_limit < traffic_usage   OR  $1 >= end)"

	now := time.Now()
	rows, err := i.db.Conn().Query(query, now)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return i.getInboundsFromRows(rows)
}

func (i *Inbound) RetriveActiveInbounds() ([]*entity.Inbound, error) {
	scope := "inboundRepository.RetriveDeactivedInbounds"

	query := "SELECT * FROM inbounds WHERE is_active = true AND is_block = false AND traffic_limit > traffic_usage   AND  $1 < end"

	now := time.Now()
	rows, err := i.db.Conn().Query(query, now)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return i.getInboundsFromRows(rows)
}

func (i *Inbound) RetriveActiveInboundsOverQuota() ([]*entity.Inbound, error) {
	scope := "inboundRepository.RetriveFaultyInbounds"

	query := "SELECT * FROM inbounds WHERE (is_active = true AND traffic_limit <= traffic_usage)"

	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return i.getInboundsFromRows(rows)
}

func (i *Inbound) RetriveDeactiveInboundsCharged() ([]*entity.Inbound, error) {
	scope := "inboundRepository.RetriveFaultyInbounds"

	query := "SELECT * FROM inbounds WHERE (is_block = false AND start <= ? AND ? <= end AND is_active = false AND traffic_limit > traffic_usage)"
	now := time.Now()
	rows, err := i.db.Conn().Query(query, now, now)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return i.getInboundsFromRows(rows)
}

func (i *Inbound) FindInboundIsNotAssigned() ([]*entity.Inbound, error) {
	scope := "inboundRepository.FindInboundIsNotAssigned"

	query := "SELECT * FROM inbounds WHERE is_assigned = false AND is_block = false"
	rows, err := i.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
	}
	inbounds := make([]*entity.Inbound, 0)
	for rows.Next() {
		inbound, err := i.scan(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
		}
		inbounds = append(inbounds, inbound)
	}
	return inbounds, nil
}

func (i *Inbound) UpdateDomainPort(id string, domain string, port string, VPNID string) error {
	scope := "inboundRepository.UpdateDomainPort"

	sql := fmt.Sprintf(
		"UPDATE inbounds SET domain = '%s', port = '%s', vpn_id=%s, is_assigned = true WHERE id = %v",
		domain,
		port,
		VPNID,
		id,
	)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id, domain, port).DebuggingError()
	}
	return nil
}

func (i *Inbound) ChangeBlockState(id string, state bool) error {
	scope := "inboundRepository.Block"

	sql := fmt.Sprintf(
		"UPDATE inbounds SET is_block = %t WHERE id = %s",
		state,
		id,
	)
	result, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	return nil
}

func (i *Inbound) Update(id string, inpt *inboundDto.UpdateInboundDto) error {
	scope := "inboundRepository.Update"

	subUpdates := []string{}
	if !inpt.Start.IsZero() {
		subUpdates = append(subUpdates, fmt.Sprintf("start = '%s'", inpt.Start.Format(time.DateTime)))
	}

	if !inpt.End.IsZero() {
		subUpdates = append(subUpdates, fmt.Sprintf("end = '%s'", inpt.End.Format(time.DateTime)))
	}

	if len(subUpdates) == 0 {
		return momoError.Scope(scope).Input(inpt, id).UnExpected().DebuggingError()
	}

	if inpt.TrafficLimit != 0 {
		subUpdates = append(subUpdates, fmt.Sprintf("traffic_limit = %d", inpt.TrafficLimit))
	}

	sql := fmt.Sprintf(
		"UPDATE inbounds SET %s WHERE id = %s",
		strings.Join(subUpdates, ", "),
		id,
	)

	result, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	return nil
}

func (i *Inbound) IncreaseTrafficUsage(id string, trafficUsage uint32) error {
	scope := "inboundRepository.IncreaseTrafficUsage"
	sql := fmt.Sprintf(
		"UPDATE inbounds SET traffic_usage = traffic_usage + %v WHERE id = %s",
		trafficUsage,
		id,
	)

	result, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	return nil
}

func (i *Inbound) ExtendInbound(id string, inpt *inboundDto.ExtendInboundDto) error {
	scope := "inboundRepository.extendinbound"

	sql := fmt.Sprintf(
		"UPDATE inbounds SET end = '%s', start = '%s', traffic_limit = traffic_limit + %d WHERE id = %s",
		inpt.End.Format(time.DateTime),
		inpt.Start.Format(time.DateTime),
		inpt.TrafficExtended,
		id,
	)
	result, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	return nil
}

func (i *Inbound) getInboundsFromRows(rows *sql.Rows) ([]*entity.Inbound, error) {
	scope := "inboundRepository.getInboundsFromRows"

	inbounds := make([]*entity.Inbound, 0)

	for rows.Next() {
		inbound, err := i.scan(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
		}
		inbounds = append(inbounds, inbound)
	}

	return inbounds, nil
}

func (i *Inbound) scan(rows *sql.Rows) (*entity.Inbound, error) {
	scope := "inboundRepository.scan"

	inbound := &entity.Inbound{}
	var createdAt, updatedAt interface{}
	var vpnType string
	var VPNID sql.NullInt64
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
		&inbound.TrafficUsage,
		&inbound.TrafficLimit,
		&inbound.Country,
		&createdAt,
		&updatedAt,
		&VPNID,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}
	inbound.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
	inbound.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
	if VPNID.Valid {
		inbound.VPNID = strconv.Itoa(int(VPNID.Int64))
	}
	return inbound, nil
}
