package vpnsource

import (
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

type VPNSource struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *VPNSource {
	return &VPNSource{
		db: db,
	}
}
