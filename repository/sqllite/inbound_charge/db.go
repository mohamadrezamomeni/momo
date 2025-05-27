package inboundcharge

import "github.com/mohamadrezamomeni/momo/repository/sqllite"

type InboundCharge struct {
	db *sqllite.SqlliteDB
}

func New(sqlite *sqllite.SqlliteDB) *InboundCharge {
	return &InboundCharge{
		db: sqlite,
	}
}
