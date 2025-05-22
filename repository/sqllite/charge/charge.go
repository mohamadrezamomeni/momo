package charge

import (
	"database/sql"
	"fmt"
	"reflect"
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

	charge := &entity.Charge{}
	status := ""
	err := c.db.Conn().QueryRow(`
	INSERT INTO charges (status, detail, admin_comment, inbound_id, user_id)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id, status, detail, admin_comment, inbound_id, user_id
	`, entity.TranslateChargeStatus(inpt.Status),
		inpt.Detail,
		"",
		inpt.InboundID,
		inpt.UserID,
	).Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
		&charge.UserID,
	)

	if err == nil {
		charge.Status = entity.ConvertStringToChargeStatus(status)
		return charge, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
}

func (c *Charge) FindChargeByID(id string) (*entity.Charge, error) {
	scope := "chargeRepository.FindChargeByID"

	var createdAt interface{}
	charge := &entity.Charge{}
	status := ""
	s := fmt.Sprintf("SELECT * FROM charges WHERE id=%s LIMIT 1", id)
	err := c.db.Conn().QueryRow(s).Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
		&charge.UserID,
		&createdAt,
	)

	if err == nil {
		charge.Status = entity.ConvertStringToChargeStatus(status)
		return charge, nil
	}
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).NotFound().Input(id).DebuggingError()
	}
	return nil, momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
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
	v := reflect.ValueOf(*filterChargesDto)
	t := reflect.TypeOf(*filterChargesDto)

	subQueries := []string{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Name == "UserID" && len(value.String()) > 0 {
			subQueries = append(subQueries, fmt.Sprintf("user_id = '%s'", value.String()))
		} else if field.Name == "InboundID" && len(value.String()) > 0 {
			subQueries = append(subQueries, fmt.Sprintf("inbound_id = %s", value.String()))
		} else if field.Name == "Status" && value.Int() > 0 {
			subQueries = append(subQueries, fmt.Sprintf("status = '%s'", entity.TranslateChargeStatus(int(value.Int()))))
		}
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
	var createdAt interface{}
	err := rows.Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
		&charge.UserID,
		&createdAt,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}
	charge.Status = entity.ConvertStringToChargeStatus(status)
	return charge, nil
}
