package user

import (
	"os"
	"testing"

	"momo/pkg/config"

	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/repository/sqllite/user/dto"
)

var (
	email     string = "momo@world.com"
	firstName string = "momo"
	lastName  string = "momoian"
	userRepo  IUserRepository
)

func TestMain(m *testing.M) {
	cfg, err := config.Load("config_test.yaml")
	if err != nil {
		os.Exit(1)
	}
	db := sqllite.New(&cfg.DB)

	migrate := migrate.New(&cfg.DB)

	migrate.UP()

	userRepo = New(db)

	code := m.Run()

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	user, err := userRepo.Create(&dto.Create{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if user.Email != email || user.FirstName != firstName || user.LastName != lastName || user.ID != "" {
		t.Error("user creation requires some troubleshooting")
		return
	}
}
