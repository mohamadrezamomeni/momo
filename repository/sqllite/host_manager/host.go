package hostmanager

import (
	"fmt"

	hostmanagerDto "momo/dto/repository/host_manager"
	"momo/entity"
	momoError "momo/pkg/error"
)

func (h *Host) Create(inpt *hostmanagerDto.AddHost) (*entity.Host, error) {
	var host *entity.Host = &entity.Host{
		Domain: inpt.Domain,
		Port:   inpt.Port,
		Status: entity.Deactive,
	}
	err := h.db.Conn().QueryRow(`
	INSERT INTO hosts (domain, port, status)
	VALUES (?, ?, ?)
	RETURNING id
`, 0, host.Domain, host.Port, host.HostStatusString()).Scan(&host.ID)
	if err != nil {
		return &entity.Host{}, momoError.Errorf("somoething went wrong to save host error: %v", err)
	}

	return host, nil
}

func (h *Host) Delete(id int) error {
	sql := fmt.Sprintf("DELETE FROM inbounds WHERE id=%v", id)
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
