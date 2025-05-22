package charge

import "github.com/mohamadrezamomeni/momo/repository/sqllite"

type Charge struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *Charge {
	return &Charge{
		db: db,
	}
}
