package tier

import (
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

type Tier struct {
	db *sqllite.SqlliteDB
}

func New(sqlite *sqllite.SqlliteDB) *Tier {
	return &Tier{
		db: sqlite,
	}
}
