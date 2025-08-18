package usertier

import (
	"database/sql"
	"fmt"

	userTierRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/user_tier"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (ut *UserTier) Create(createDto *userTierRepoDto.Create) error {
	scope := "userTierRepository.create"

	_, err := ut.db.Conn().Exec(`
	INSERT INTO user_tiers (user_id, tier)
	VALUES (?, ?)`,
		createDto.UserID,
		createDto.Tier,
	)
	if err == nil {
		return nil
	}

	if errorRepository.IsDuplicateError(err) {
		return momoError.Wrap(err).Input(createDto).Duplicate().Scope(scope).DebuggingError()
	}
	return momoError.Wrap(err).Input(createDto).UnExpected().Scope(scope).DebuggingError()
}

func (ut *UserTier) Delete(identifyUserTierDto *userTierRepoDto.IdentifyUserTier) error {
	scope := "userTierRepository.Delete"

	sql := fmt.Sprintf(
		"DELETE FROM user_tiers WHERE user_id = '%s' AND tier = '%s'",
		identifyUserTierDto.UserID,
		identifyUserTierDto.Tier,
	)
	res, err := ut.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(identifyUserTierDto).DebuggingError()
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(identifyUserTierDto).DebuggingError()
	}

	if rowsAffected == 0 {
		return momoError.Wrap(err).Scope(scope).Input(identifyUserTierDto).DebuggingError()
	}
	return nil
}

func (ut *UserTier) DeleteAll() error {
	scope := "userTierRepository.DeleteAll"

	sql := "DELETE FROM user_tiers"
	res, err := ut.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}
	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (ut *UserTier) FilterTiersBelongToUser(userID string) ([]*entity.Tier, error) {
	scope := "userTierRepository.FilterTiersBelongToUser"

	rows, err := ut.db.Conn().Query(`
    SELECT t.name, t.is_default
    FROM tiers t
    LEFT JOIN user_tiers ut ON t.name = ut.tier AND ut.user_id = ?
    WHERE ut.user_id IS NOT NULL OR t.is_default = true
`, userID)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(userID).ErrorWrite()
	}
	return ut.getTiersFromRows(rows)
}

func (ut *UserTier) FilterTiersByUser(userID string) ([]*entity.Tier, error) {
	scope := "userTierRepository.FilterTiersBelongToUser"

	rows, err := ut.db.Conn().Query(`
    SELECT t.name, t.is_default
    FROM tiers t
    LEFT JOIN user_tiers ut ON t.name = ut.tier AND ut.user_id = ?
    WHERE ut.user_id IS NOT NULL
`, userID)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(userID).ErrorWrite()
	}
	return ut.getTiersFromRows(rows)
}

func (ut *UserTier) getTiersFromRows(rows *sql.Rows) ([]*entity.Tier, error) {
	scope := "userTierRepository.getTiersFromRows"

	tiers := make([]*entity.Tier, 0)

	for rows.Next() {
		userTier, err := ut.scanTier(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
		}
		tiers = append(tiers, userTier)
	}

	return tiers, nil
}

func (i *UserTier) scanTier(rows *sql.Rows) (*entity.Tier, error) {
	scope := "userTierRepository.scan"

	tier := &entity.Tier{}

	err := rows.Scan(
		&tier.Name,
		&tier.IsDefault,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}

	return tier, nil
}
