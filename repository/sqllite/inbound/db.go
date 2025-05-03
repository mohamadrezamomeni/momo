package inbound

import "github.com/mohamadrezamomeni/momo/repository/sqllite"

type Inbound struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *Inbound {
	return &Inbound{
		db: db,
	}
}
