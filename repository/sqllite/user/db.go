package user

import (
	"momo/repository/sqllite"
)

type User struct {
	db *sqllite.SqlliteDB
}

func New(sqlite *sqllite.SqlliteDB) *User {
	return &User{
		db: sqlite,
	}
}
