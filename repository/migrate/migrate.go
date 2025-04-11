package migrate

import (
	"database/sql"
	"fmt"

	"momo/repository/sqllite"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	path       string
	migrations *migrate.FileMigrationSource
}

type IMigrator interface {
	Up()
}

func New(cfg *sqllite.DBConfig) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: cfg.Migrations,
	}

	return Migrator{path: cfg.Path, dialect: cfg.Dialect, migrations: migrations}
}

func (m *Migrator) UP() {
	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(fmt.Errorf("can't open sqllite db: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %v", err))
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
