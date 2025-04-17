package sqllite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqlliteDB struct {
	db *sql.DB
}

func (s *SqlliteDB) Conn() *sql.DB {
	return s.db
}

func New(cfg *DBConfig) (*SqlliteDB, error) {
	db, err := sql.Open(cfg.Dialect, cfg.Path)

	return &SqlliteDB{
		db: db,
	}, err
}
