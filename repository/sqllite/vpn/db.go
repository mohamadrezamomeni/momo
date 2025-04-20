package vpn

import (
	"momo/repository/sqllite"
)

type VPN struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *VPN {
	return &VPN{
		db: db,
	}
}
