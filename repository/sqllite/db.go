package sqllite

import (
	"database/sql"

	momoError "momo/pkg/error"

	_ "github.com/mattn/go-sqlite3"
)

type SqlliteDB struct {
	db *sql.DB
}

func (s *SqlliteDB) Conn() *sql.DB {
	return s.db
}

func New(cfg *DBConfig) *SqlliteDB {
	db, err := sql.Open(cfg.Dialect, cfg.Path)
	if err != nil {
		panic(momoError.Errorf("ERROR: something went wrong with connectiong db: %s", err))
	}
	return &SqlliteDB{
		db: db,
	}
}
