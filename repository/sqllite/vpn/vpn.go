package vpn

import (
	"fmt"

	"momo/entity"
	momoError "momo/pkg/error"
	"momo/repository/sqllite/vpn/dto"
)

func (v *VPN) Create(inpt *dto.Add_VPN) (*entity.VPN, error) {
	var vpn *entity.VPN = &entity.VPN{}
	err := v.db.Conn().QueryRow(`
	INSERT INTO vpns (domain, is_active, api_port, start_range_port, end_range_port, vpn_type, user_count)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING id, domain, is_active, api_port, start_range_port, end_range_port, vpn_type, user_count
	`, inpt.Domain, inpt.IsActive, inpt.ApiPort, inpt.StartRangePort, inpt.EndRangePort, inpt.VPNType, inpt.UserCount).Scan(
		&vpn.ID,
		&vpn.Domain,
		&vpn.IsActive,
		&vpn.ApiPort,
		&vpn.StartRangePort,
		&vpn.EndRangePort,
		&vpn.VPNType,
		&vpn.UserCount,
	)
	if err != nil {
		return vpn, momoError.DebuggingErrorf("something wrong has happend the problem was %v", err)
	}
	return vpn, nil
}

func (i *VPN) Delete(id int) error {
	sql := fmt.Sprintf("DELETE FROM vpns WHERE id=%v", id)
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

func (i *VPN) activeVPN(id int) error {
	return i.updateActivationVPN(id, true)
}

func (i *VPN) deactiveVPN(id int) error {
	return i.updateActivationVPN(id, false)
}

func (i *VPN) updateActivationVPN(id int, status bool) error {
	sql := fmt.Sprintf("UPDATE vpns SET is_active = %v WHERE id = %v", status, id)

	_, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.DebuggingErrorf("something bad has happend the error was %v", err)
	}
	return nil
}
