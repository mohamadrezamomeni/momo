package hostmanager

import "github.com/mohamadrezamomeni/momo/repository/sqllite"

type Host struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *Host {
	return &Host{
		db: db,
	}
}
