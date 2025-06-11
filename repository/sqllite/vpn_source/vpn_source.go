package vpnsource

import (
	"database/sql"
	"fmt"

	vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (vs *VPNSource) Create(createdVPNSource *vpnSourceRepositoryDto.CreateVPNSourceDto) (*entity.VPNSource, error) {
	scope := "vpnSourceRepository.Create"

	vpnSource := &entity.VPNSource{}
	err := vs.db.Conn().QueryRow(`
	INSERT INTO vpn_source (title, english)
	VALUES (?, ?)
	RETURNING id, title, english
	`, createdVPNSource.Title,
		createdVPNSource.English,
	).Scan(
		&vpnSource.ID,
		&vpnSource.Title,
		&vpnSource.English,
	)
	if err == nil {
		return vpnSource, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(createdVPNSource).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(createdVPNSource).UnExpected().Scope(scope).DebuggingError()
}

func (vs *VPNSource) Find(id string) (*entity.VPNSource, error) {
	scope := "vpnSourceRepository.find"

	vpnSource := &entity.VPNSource{}

	s := fmt.Sprintf("SELECT * FROM vpn_source WHERE id=%s LIMIT 1", id)
	err := vs.db.Conn().QueryRow(s).Scan(
		&vpnSource.ID,
		&vpnSource.Title,
		&vpnSource.English,
	)
	if err == nil {
		return vpnSource, nil
	}
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).Input(id).NotFound().DebuggingError()
	}
	return nil, momoError.Wrap(err).Scope(scope).Input(id).UnExpected().DebuggingError()
}

func (i *VPNSource) DeleteAll() error {
	scope := "vpnSourceRepository.DeleteAll"

	sql := "DELETE FROM vpn_source"
	res, err := i.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}
