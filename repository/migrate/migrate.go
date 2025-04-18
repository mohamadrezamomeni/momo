package migrate

import (
	"database/sql"
	"path/filepath"

	"momo/repository/sqllite"

	momoError "momo/pkg/error"
	"momo/pkg/utils"

	momoLogger "momo/pkg/log"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	path       string
	migrations *migrate.FileMigrationSource
}

type IMigrator interface {
	UP()
	DOWN()
}

func New(cfg *sqllite.DBConfig) IMigrator {
	root, _ := utils.GetRootOfProject()
	migrations := &migrate.FileMigrationSource{
		Dir: filepath.Join(root, cfg.Migrations),
	}

	return &Migrator{path: filepath.Join(root, cfg.Path), dialect: cfg.Dialect, migrations: migrations}
}

func (m *Migrator) UP() {
	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(momoError.Errorf("unable to open sqllite db: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(momoError.Errorf("unable to apply migrations: %v", err))
	}
	momoLogger.Infof("Applied %d migrations!", n)
}

func (m *Migrator) DOWN() {
	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(momoError.Errorf("unable to open sqllite db: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(momoError.Errorf("unable to undo migrations: %v", err))
	}
	momoLogger.Infof("undo %d migrations!", n)
}
