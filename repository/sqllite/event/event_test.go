package event

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var eventRepo *Event

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-event.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	eventRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreatingEvent(t *testing.T) {
	event, err := eventRepo.Create(data1)
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}
	if event.Data != data1.Data ||
		event.Name != data1.Name {
		t.Error("error to compare data")
	}
}
