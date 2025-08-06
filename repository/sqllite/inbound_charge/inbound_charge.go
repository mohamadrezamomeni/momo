package inboundcharge

import (
	"database/sql"
	"time"

	inboundChargeRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound_charge"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
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

func (ic *InboundCharge) CreateInbound(
	chargeID string,
	createInboundByCharge *inboundChargeRepoDto.CreateInboundByCharge,
) error {
	scope := "inboundcharge.repository.CreateInbound"

	tx, err := ic.db.Conn().Begin()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	inboundID, err := ic.createInbound(tx, createInboundByCharge)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ic.updateChargeInboundID(tx, chargeID, inboundID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("error to commit")
	}
	return nil
}

func (ic *InboundCharge) updateChargeInboundID(tx *sql.Tx, chargeID string, inboundID string) error {
	scope := "inboundcharge.repository.UpdateChargeInboundID"

	sql := "UPDATE charges SET inbound_id = $1, status = $2 WHERE id = $3"
	_, err := tx.Exec(sql, inboundID, entity.TranslateChargeStatus(entity.AssignedCharged), chargeID)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (ic *InboundCharge) createInbound(
	tx *sql.Tx,
	createInboundByCharge *inboundChargeRepoDto.CreateInboundByCharge,
) (string, error) {
	scope := "inboundcharge.repository.createInbound"

	var id string
	err := tx.QueryRow(`
	INSERT INTO inbounds (protocol, domain, vpn_type, port, user_id, tag, is_active, start, end, is_block, is_assigned, is_notified, traffic_usage, traffic_limit, country)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id
	`, createInboundByCharge.Protocol,
		createInboundByCharge.Domain,
		entity.VPNTypeString(createInboundByCharge.VPNType),
		createInboundByCharge.Port,
		createInboundByCharge.UserID,
		createInboundByCharge.Tag,
		createInboundByCharge.IsActive,
		createInboundByCharge.Start,
		createInboundByCharge.End,
		createInboundByCharge.IsBlock,
		createInboundByCharge.IsAssigned,
		createInboundByCharge.IsNotified,
		createInboundByCharge.TrafficUsage,
		createInboundByCharge.TrafficLimit,
		createInboundByCharge.Country,
	).Scan(&id)

	if errorRepository.IsDuplicateError(err) {
		return "", momoError.Wrap(err).Input(createInboundByCharge).Duplicate().Scope(scope).DebuggingError()
	}
	if err != nil {
		return "", momoError.Wrap(err).Input(createInboundByCharge).DebuggingError()
	}
	return id, nil
}
