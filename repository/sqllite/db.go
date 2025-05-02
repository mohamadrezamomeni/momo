package sqllite

import (
	"database/sql"
	"path/filepath"

	momoError "momo/pkg/error"
	"momo/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

type SqlliteDB struct {
	db *sql.DB
}

func (s *SqlliteDB) Conn() *sql.DB {
	return s.db
}

func New(cfg *DBConfig) *SqlliteDB {
	scope := "initializeDB"
	root, _ := utils.GetRootOfProject()
	path := filepath.Join(root, cfg.Path)
	db, err := sql.Open(cfg.Dialect, path)
	if err != nil {
		panic(momoError.Wrap(err).Scope(scope).Errorf("something went wrong to connect db and the address of db is  %s", path))
	}
	return &SqlliteDB{
		db: db,
	}
}
