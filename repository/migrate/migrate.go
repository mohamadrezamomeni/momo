package migrate

import (
	"database/sql"

	"momo/repository/sqllite"

	momoError "momo/pkg/error"

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
}

func New(cfg *sqllite.DBConfig) IMigrator {
	migrations := &migrate.FileMigrationSource{
		Dir: cfg.Migrations,
	}

	return &Migrator{path: cfg.Path, dialect: cfg.Dialect, migrations: migrations}
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
