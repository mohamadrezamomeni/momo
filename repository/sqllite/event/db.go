package event

import "github.com/mohamadrezamomeni/momo/repository/sqllite"

type Event struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *Event {
	return &Event{
		db: db,
	}
}
