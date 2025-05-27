package inboundcharge

import (
	"database/sql"
	"time"

	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (ic *InboundCharge) AssignChargeToInbound(
	inbound *entity.Inbound,
	charge *entity.Charge,
	vpnPackage *entity.VPNPackage,
) error {
	scope := "inboundcharge.repository.assignChargeToInbound"

	tx, err := ic.db.Conn().Begin()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	err = ic.updateInbound(tx, inbound, charge, vpnPackage)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ic.updateChargeStatusToAssigned(tx, charge)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("error to commit")
	}
	return nil
}

func (ic *InboundCharge) updateChargeStatusToAssigned(tx *sql.Tx, charge *entity.Charge) error {
	scope := "inboundcharge.repository.UpdateChargeStatusToAssigned"

	sql := "UPDATE charges SET status = $1 WHERE id = $2"
	_, err := tx.Exec(sql, entity.TranslateChargeStatus(entity.AssignedCharged), charge.ID)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (ic *InboundCharge) updateInbound(
	tx *sql.Tx,
	inbound *entity.Inbound,
	charge *entity.Charge,
	vpnPackage *entity.VPNPackage,
) error {
	scope := "inboundcharge.repository.UpdateInbound"
	now := time.Now()
	end := now.AddDate(0, int(vpnPackage.Months), int(vpnPackage.Days))
	sql := "UPDATE inbounds SET start = $1, end = $2, traffic_usage = $3, traffic_limit = $4 WHERE id = $5"
	_, err := tx.Exec(sql, now, end, 0, vpnPackage.TrafficLimit, inbound.ID)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}
