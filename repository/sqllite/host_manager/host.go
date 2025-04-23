package hostmanager

import (
	"database/sql"
	"fmt"
	"strings"

	hostmanagerDto "momo/dto/repository/host_manager"
	"momo/entity"
	momoError "momo/pkg/error"
)

func (h *Host) Create(inpt *hostmanagerDto.AddHost) (*entity.Host, error) {
	var host *entity.Host = &entity.Host{
		Domain:         inpt.Domain,
		Port:           inpt.Port,
		Status:         inpt.Status,
		StartRangePort: inpt.StartRangePort,
		EndRangePort:   inpt.EndRangePort,
	}
	err := h.db.Conn().QueryRow(`
	INSERT INTO hosts (domain, port, status, start_range_port, end_range_port)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id
`, host.Domain, host.Port, entity.HostStatusString(host.Status), inpt.StartRangePort, inpt.EndRangePort).Scan(&host.ID)
	if err != nil {
		return &entity.Host{}, momoError.Errorf("somoething went wrong to save host error: %v", err)
	}

	return host, nil
}

func (h *Host) Delete(id int) error {
	sql := fmt.Sprintf("DELETE FROM hosts WHERE id=%v", id)
	res, err := h.db.Conn().Exec(sql)
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

func (h *Host) Filter(inpt *hostmanagerDto.FilterHosts) ([]*entity.Host, error) {
	query := h.makeQuery(inpt)
	rows, err := h.db.Conn().Query(query)
	if err != nil {
		return []*entity.Host{}, momoError.DebuggingErrorf("error has occured err: %v", err)
	}

	hosts := make([]*entity.Host, 0)

	for rows.Next() {
		host, err := h.scan(rows)
		if err != nil {
			return []*entity.Host{}, err
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
	host := &entity.Host{}
	var hostStatusString string
	var createdAt, updatedAt interface{}

	err := rows.Scan(
		&host.ID,
		&host.Domain,
		&host.Port,
		&hostStatusString,
		&host.StartRangePort,
		&host.EndRangePort,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return host, momoError.DebuggingErrorf("error has occured err: %v", err)
	}
	status, er := entity.MapTuStatus(hostStatusString)
	if err != nil {
		return nil, er
	}
	host.Status = status
	return host, nil
}
