package charge

import (
	"database/sql"
	"fmt"
	"strings"

	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (e *Charge) DeleteAll() error {
	scope := "chargeRepository.DeleteAll"

	sql := "DELETE FROM charges"
	res, err := e.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return nil
}

func (c *Charge) Create(inpt *chargeRepositoryDto.CreateDto) (*entity.Charge, error) {
	scope := "chargeRepository.create"

	charge := &entity.Charge{
		VPNType: inpt.VPNType,
	}

	status := ""
	err := c.db.Conn().QueryRow(`
	INSERT INTO charges (status, detail, admin_comment, inbound_id, user_id, package_id, country, vpn_type)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id, status, detail, admin_comment, inbound_id, user_id, package_id, country
	`, entity.TranslateChargeStatus(inpt.Status),
		inpt.Detail,
		"",
		inpt.InboundID,
		inpt.UserID,
		inpt.PackageID,
		inpt.Country,
		entity.VPNTypeString(inpt.VPNType),
	).Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
		&charge.UserID,
		&charge.PackageID,
		&charge.Country,
	)

	if err != nil && errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	if err != nil {
		return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
	}

	charge.Status = entity.ConvertStringToChargeStatus(status)
	return charge, nil
}

func (c *Charge) FindChargeByID(id string) (*entity.Charge, error) {
	scope := "chargeRepository.FindChargeByID"

	var createdAt interface{}
	charge := &entity.Charge{}
	status := ""
	var VPNType string
	s := fmt.Sprintf("SELECT * FROM charges WHERE id=%s LIMIT 1", id)
	err := c.db.Conn().QueryRow(s).Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
		&charge.UserID,
		&charge.PackageID,
		&charge.Country,
		&VPNType,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).NotFound().Input(id).DebuggingError()
	}
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
	}
	charge.VPNType = entity.ConvertStringVPNTypeToEnum(VPNType)
	charge.Status = entity.ConvertStringToChargeStatus(status)

	return charge, nil
}

func (c *Charge) UpdateCharge(id string, inpt *chargeRepositoryDto.UpdateChargeDto) error {
	scope := "chargeRepository.Update"
	subModifies := []string{}

	if inpt.Status != 0 {
		subModifies = append(subModifies, fmt.Sprintf("status = '%s'", entity.TranslateChargeStatus(inpt.Status)))
	}

	if inpt.AdminComment != "" {
		subModifies = append(subModifies, fmt.Sprintf("admin_comment = '%s'", inpt.AdminComment))
	}

	if inpt.Detail != "" {
		subModifies = append(subModifies, fmt.Sprintf("detail = '%s'", inpt.Detail))
	}

	if len(subModifies) == 0 {
		return momoError.Scope(scope).DebuggingErrorf("input was empty")
	}
	sql := fmt.Sprintf(
		"UPDATE charges SET %s WHERE id = %v",
		strings.Join(subModifies, ", "),
		id,
	)
	_, err := c.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id, inpt).DebuggingError()
	}
	return nil
}

func (c *Charge) RetriveAvailbleChargesForInbounds(inboundIDs []string) ([]*entity.Charge, error) {
	scope := "chargeRepository.RetriveAvailbleChargesForInbounds"

	queryFormat := `SELECT c.*
		FROM charges c
		JOIN (
			SELECT id, MIN(created_at) AS min_created_at
			FROM charges
			WHERE inbound_id IN (%s) AND status = '%s'
			GROUP BY inbound_id
		) available_charge
		ON c.id = available_charge.id`

	query := fmt.Sprintf(
		queryFormat,
		strings.Join(inboundIDs, ", "),
		entity.TranslateChargeStatus(entity.ApprovedStatusCharge),
	)

	rows, err := c.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Input(inboundIDs).Scope(scope).DebuggingError()
	}

	charges := make([]*entity.Charge, 0)

	for rows.Next() {
		charge, err := c.scan(rows)
		if err != nil {
			return nil, err
		}
		charges = append(charges, charge)
	}
	return charges, nil
}

func (c *Charge) FilterCharges(filterChargesDto *chargeRepositoryDto.FilterChargesDto) ([]*entity.Charge, error) {
	scope := "chargeRepository.Filter"

	query := c.makeQuery(filterChargesDto)
	rows, err := c.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Input(filterChargesDto).Scope(scope).DebuggingError()
	}

	charges := make([]*entity.Charge, 0)

	for rows.Next() {
		charge, err := c.scan(rows)
		if err != nil {
			return nil, err
		}
		charges = append(charges, charge)
	}
	return charges, nil
}

func (c *Charge) makeQuery(filterChargesDto *chargeRepositoryDto.FilterChargesDto) string {
	subQueries := make([]string, 0)
	if filterChargesDto.InboundID != "" {
		subQueries = append(subQueries, fmt.Sprintf("inbound_id = %s", filterChargesDto.InboundID))
	}
	if filterChargesDto.UserID != "" {
		subQueries = append(subQueries, fmt.Sprintf("user_id = '%s'", filterChargesDto.UserID))
	}
	if len(filterChargesDto.Statuses) > 0 {
		subQueries = append(
			subQueries,
			fmt.Sprintf("status IN ('%s')",
				strings.Join(entity.ConvertStatusesToStatusLabels(filterChargesDto.Statuses), "', '"),
			),
		)
	}

	sql := "SELECT * FROM charges"
	if len(subQueries) > 0 {
		sql += fmt.Sprintf(" WHERE %s", strings.Join(subQueries, " AND "))
	}

	return sql
}

func (e *Charge) scan(rows *sql.Rows) (*entity.Charge, error) {
	scope := "chargeRepository.scan"

	charge := &entity.Charge{}
	status := ""
	var vpnType string
	var adminComment sql.NullString
	var inboundID sql.NullString
	var createdAt interface{}
	err := rows.Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&adminComment,
		&inboundID,
		&charge.UserID,
		&charge.PackageID,
		&charge.Country,
		&vpnType,
		&createdAt,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}

	if adminComment.Valid {
		charge.AdminComment = adminComment.String
	}
	if inboundID.Valid {
		charge.InboundID = inboundID.String
	}
	charge.VPNType = entity.ConvertStringVPNTypeToEnum(vpnType)
	charge.Status = entity.ConvertStringToChargeStatus(status)

	return charge, nil
}

func (e *Charge) GetFirstApprovedInboundCharge(inboundID string) (*entity.Charge, error) {
	scope := "chargeRepository.GetFirstAvailbleInboundCharge"

	query := fmt.Sprintf("SELECT * FROM charges WHERE inbound_id = %s AND status = 'approved' ORDER BY created_at DESC LIMIT 1", inboundID)

	charge := &entity.Charge{}
	status := ""
	var createdAt interface{}
	var vpnType string
	err := e.db.Conn().QueryRow(query).Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
		&charge.UserID,
		&charge.PackageID,
		&charge.Country,
		&vpnType,
		&createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).NotFound().Input(inboundID).DebuggingError()
	}
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(inboundID).UnExpected().DebuggingError()
	}

	charge.Status = entity.ConvertStringToChargeStatus(status)
	return charge, nil
}

func (c *Charge) RetriveChargesApprovedWithoutInbound() ([]*entity.Charge, error) {
	scope := "chargeRepository.RetriveChargesApprovedWithoutInbound"

	query := "SELECT * FROM charges WHERE (inbound_id IS NULL OR inbound_id = \"\") AND status = 'approved'"
	rows, err := c.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	charges := make([]*entity.Charge, 0)

	for rows.Next() {
		charge, err := c.scan(rows)
		if err != nil {
			return nil, err
		}
		charges = append(charges, charge)
	}
	return charges, nil
}
