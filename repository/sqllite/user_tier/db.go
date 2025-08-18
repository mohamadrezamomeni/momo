package usertier

import (
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

type UserTier struct {
	db *sqllite.SqlliteDB
}

func New(sqlite *sqllite.SqlliteDB) *UserTier {
	return &UserTier{
		db: sqlite,
	}
}
