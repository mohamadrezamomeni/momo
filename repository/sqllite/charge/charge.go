package charge

import (
	"database/sql"
	"fmt"

	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (e *Charge) DeleteAll() error {
	scope := "chargeRepository.DeleteAll"

	sql := "DELETE FROM events"
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
	INSERT INTO charges (status, detail, admin_comment, inbound_id)
	VALUES (?, ?, ?, ?)
	RETURNING id, status, detail, admin_comment, inbound_id
	`, entity.TranslateChargeStatus(inpt.Status),
		inpt.Detail,
		"",
		inpt.InboundID,
	).Scan(
		&charge.ID,
		&status,
		&charge.Detail,
		&charge.AdminComment,
		&charge.InboundID,
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
