package event

import (
	"os"
	"testing"

	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
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
	defer eventRepo.DeleteAll()
	event, err := eventRepo.Create(data1)
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}
	if event.Data != data1.Data ||
		event.Name != data1.Name {
		t.Error("error to compare data")
	}
}

func TestFilter(t *testing.T) {
	defer eventRepo.DeleteAll()
	eventRepo.Create(data1)
	eventRepo.Create(data2)
	eventRepo.Create(data3)

	events, err := eventRepo.Filter(&eventRepositoryDto.FilterEvents{})
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	if len(events) > 3 {
		t.Errorf("we expected lengh of result be 3 but we got %d", len(events))
	}

	active := true

	events, err = eventRepo.Filter(&eventRepositoryDto.FilterEvents{
		IsProcessed: &active,
	})
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	if len(events) != 1 {
		t.Errorf("we expected lengh of result be 1 but we got %d", len(events))
	}

	events, err = eventRepo.Filter(&eventRepositoryDto.FilterEvents{
		IsProcessed: &active,
		Name:        "notification",
	})
	if err != nil {
		t.Fatalf("something went wrong the problem was %v", err)
	}

	if len(events) != 0 {
		t.Errorf("we expected lengh of result be 0 but we got %d", len(events))
	}
}
