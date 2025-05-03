package migrate

import (
	"database/sql"
	"path/filepath"

	"github.com/mohamadrezamomeni/momo/repository/sqllite"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"

	momoLogger "github.com/mohamadrezamomeni/momo/pkg/log"

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
	scope := "migration.up"
	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(momoError.Wrap(err).Scope(scope).Errorf("error to connect db"))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(momoError.Wrap(err).Scope(scope).Errorf("unable to apply migrations: %v", err))
	}
	momoLogger.Infof("Applied %d migrations!", n)
}

func (m *Migrator) DOWN() {
	scope := "migration.down"

	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(momoError.Wrap(err).Scope(scope).Errorf("error to connect db"))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(momoError.Wrap(err).Scope(scope).Errorf("unable to undo migrations: %v", err))
	}
	momoLogger.Infof("undo %d migrations!", n)
}
