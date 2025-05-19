package vpnpackage

import "github.com/mohamadrezamomeni/momo/repository/sqllite"

type VPNPackage struct {
	db *sqllite.SqlliteDB
}

func New(db *sqllite.SqlliteDB) *VPNPackage {
	return &VPNPackage{
		db: db,
	}
}
