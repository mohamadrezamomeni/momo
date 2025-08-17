package tier

import (
	"database/sql"
	"fmt"
	"strings"

	tierRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/tier"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (t *Tier) Create(createTierDto tierRepoDto.CreateTier) (*entity.Tier, error) {
	scope := "tierUserRepository.createTier"

	tier := &entity.Tier{}
	err := t.db.Conn().QueryRow(`
	INSERT INTO tiers (name, is_default)
	VALUES (?, ?)
	RETURNING name, is_default
`,
		createTierDto.Name,
		createTierDto.Default,
	).Scan(
		&tier.Name,
		&tier.IsDefault,
	)
	if err == nil {
		return tier, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(createTierDto).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(createTierDto).UnExpected().Scope(scope).DebuggingError()
}

func (t *Tier) FindByName(name string) (*entity.Tier, error) {
	scope := "tierRepository.FindByName"

	var tier *entity.Tier = &entity.Tier{}
	s := fmt.Sprintf("SELECT * FROM tiers WHERE name='%s' LIMIT 1", name)
	err := t.db.Conn().QueryRow(s).Scan(
		&tier.Name,
		&tier.IsDefault,
	)
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).NotFound().Input(name).DebuggingError()
	}
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(name).UnExpected().DebuggingError()
	}

	return tier, nil
}

func (t *Tier) DeleteAll() error {
	scope := "tierRepository.DeleteAll"

	sql := "DELETE FROM tiers"
	res, err := t.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return nil
}

func (t *Tier) Filter() ([]*entity.Tier, error) {
	scope := "tiersRepository.Filter"

	query := "SELECT * FROM tiers"
	rows, err := t.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	return t.getTiersFromRows(rows)
}

func (t *Tier) Update(name string, updateDto *tierRepoDto.Update) error {
	scope := "tierRepository.Update"

	subUpdates := []string{}
	if updateDto.IsDefault != nil {
		subUpdates = append(subUpdates, fmt.Sprintf("'is_default' = %v", *updateDto.IsDefault))
	}

	sql := fmt.Sprintf(
		"UPDATE tiers SET %s WHERE name = '%s'",
		strings.Join(subUpdates, ", "),
		name,
	)

	result, err := t.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(name).ErrorWrite()
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return momoError.Wrap(err).Scope(scope).Input(name).ErrorWrite()
	}
	return nil
}

func (t *Tier) getTiersFromRows(rows *sql.Rows) ([]*entity.Tier, error) {
	scope := "tiersRepository.getTiersFromRows"

	tiers := make([]*entity.Tier, 0)

	for rows.Next() {
		tier, err := t.scan(rows)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
		}
		tiers = append(tiers, tier)
	}

	return tiers, nil
}

func (t *Tier) scan(rows *sql.Rows) (*entity.Tier, error) {
	scope := "tiersRepository.scan"

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
