package hostmanager

import (
	"database/sql"
	"fmt"
	"strings"

	hostmanagerDto "github.com/mohamadrezamomeni/momo/dto/repository/host_manager"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (h *Host) Create(inpt *hostmanagerDto.AddHost) (*entity.Host, error) {
	scope := "hostRepository.create"
	var host *entity.Host = &entity.Host{
		Rank:   inpt.Rank,
		Domain: inpt.Domain,
		Port:   inpt.Port,
		Status: inpt.Status,
	}
	err := h.db.Conn().QueryRow(`
	INSERT INTO hosts (domain, port, status, rank)
	VALUES (?, ?, ?, ?)
	RETURNING id
`, host.Domain, host.Port, entity.HostStatusString(host.Status), inpt.Rank).Scan(&host.ID)
	if err != nil {
		return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
	}

	return host, nil
}

func (h *Host) FindByID(id int) (*entity.Host, error) {
	scope := "hostRepository.FindByID"

	query := fmt.Sprintf("SELECT * FROM hosts WHERE id = %d LIMIT 1", id)
	host := &entity.Host{}
	var hostStatusString string
	var createdAt, updatedAt interface{}

	row := h.db.Conn().QueryRow(query)
	err := row.Scan(
		&host.ID,
		&host.Domain,
		&host.Port,
		&hostStatusString,
		&host.Rank,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Input(id).UnExpected().Scope(scope).DebuggingError()
	}
	status, er := entity.MapHostStatusToEnum(hostStatusString)
	if err != nil {
		return nil, er
	}
	host.Status = status
	return host, nil
}

func (h *Host) Update(id int, inpt *hostmanagerDto.UpdateHost) error {
	scope := "hostRepository.Update"

	sql := fmt.Sprintf(
		"UPDATE hosts SET status = '%s', rank = %v WHERE id = %v",
		entity.HostStatusString(inpt.Status),
		inpt.Rank,
		id,
	)
	_, err := h.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Input(inpt, id).Scope(scope).DebuggingError()
	}
	return nil
}

func (h *Host) Delete(id int) error {
	scope := "hostRepository.Delete"

	sql := fmt.Sprintf("DELETE FROM hosts WHERE id=%v", id)
	res, err := h.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
	}

	if rowsAffected == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
	}
	return nil
}

func (h *Host) DeleteAll() error {
	scope := "hostRepository.DeleteAll"

	sql := "DELETE FROM hosts"
	res, err := h.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().ErrorWrite()
	}

	return nil
}

func (h *Host) Filter(inpt *hostmanagerDto.FilterHosts) ([]*entity.Host, error) {
	scope := "hostRepository.Filter"

	query := h.makeQuery(inpt)
	rows, err := h.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(inpt).UnExpected().DebuggingError()
	}

	hosts := make([]*entity.Host, 0)

	for rows.Next() {
		host, err := h.scan(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).Input(inpt).UnExpected().DebuggingError()
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func (h *Host) makeQuery(inpt *hostmanagerDto.FilterHosts) string {
	subQueries := []string{}
	sql := "SELECT * FROM hosts"
	if len(inpt.Statuses) > 0 {
		subQueries = append(
			subQueries,
			fmt.Sprintf("status IN (%s)", strings.Join(h.makeQueryByListOfSatus(inpt.Statuses), ", ")),
		)
	}

	if len(subQueries) > 0 {
		sql = sql + " WHERE " + strings.Join(subQueries, " AND ")
	}
	return sql
}

func (h *Host) makeQueryByListOfSatus(statuses []entity.HostStatus) []string {
	ret := make([]string, 0)

	for _, status := range statuses {
		ret = append(ret, fmt.Sprintf("'%s'", entity.HostStatusString(status)))
	}
	return ret
}

func (i *Host) scan(rows *sql.Rows) (*entity.Host, error) {
	scope := "hostRepository.scan"

	host := &entity.Host{}
	var hostStatusString string
	var createdAt, updatedAt interface{}

	err := rows.Scan(
		&host.ID,
		&host.Domain,
		&host.Port,
		&hostStatusString,
		&host.Rank,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return host, momoError.Wrap(err).Scope(scope).Input(rows).UnExpected().DebuggingError()
	}
	status, er := entity.MapHostStatusToEnum(hostStatusString)
	if err != nil {
		return nil, er
	}
	host.Status = status
	return host, nil
}
