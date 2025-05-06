package user

import (
	"os"
	"testing"

	userDto "github.com/mohamadrezamomeni/momo/dto/repository/user"

	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
)

var userRepo *User

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-user.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	userRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	defer userRepo.DeleteAll()

	userCreated, err := userRepo.Create(user1)
	if err != nil {
		t.Error(err)
		return
	}

	if userCreated.Username != user1.Username ||
		userCreated.FirstName != user1.FirstName ||
		userCreated.Password != user1.Password ||
		userCreated.IsAdmin != user1.IsAdmin ||
		userCreated.LastName != user1.LastName {
		t.Error("user creation requires some troubleshooting")
		return
	}
}

func TestFindByUsername(t *testing.T) {
	userCreated, err := userRepo.Create(user1)
	defer userRepo.DeleteAll()

	user, err := userRepo.FindUserByUsername(user1.Username)
	if err != nil {
		t.Fatalf("findByUsername needs troubleshooting error: %v", err)
		userRepo.Delete(userCreated.ID)
		return
	}
	if user.Username != user1.Username ||
		user.FirstName != user1.FirstName ||
		user.LastName != user1.LastName ||
		user.ID != userCreated.ID {
		t.Error("something went wrong to compare results")
		return
	}
}

func TestFindByID(t *testing.T) {
	userCreated, err := userRepo.Create(user1)
	defer userRepo.DeleteAll()

	user, err := userRepo.FindUserByID(userCreated.ID)
	if err != nil {
		t.Errorf("findByUsername needs troubleshooting error: %v", err)
		userRepo.Delete(userCreated.ID)
		return
	}

	if user.Username != user1.Username ||
		user.FirstName != user1.FirstName ||
		user.LastName != user1.LastName ||
		user.ID != userCreated.ID ||
		user.IsSuperAdmin != true {
		t.Error("something went wrong to compare results")
		return
	}
}

func TestFilter(t *testing.T) {
	userRepo.Create(user1)
	userRepo.Create(user2)
	userRepo.Create(user3)
	defer userRepo.DeleteAll()
	users, err := userRepo.FilterUsers(&userDto.FilterUsers{
		FirstName: "mic",
	})
	if err != nil {
		t.Errorf("1. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 2 {
		t.Errorf("1. something went wrong.")
	}

	users, err = userRepo.FilterUsers(&userDto.FilterUsers{
		FirstName: "mic",
		LastName:  "casta",
	})
	if err != nil {
		t.Errorf("2. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 2 {
		t.Errorf("2. something went wrong.")
	}

	users, err = userRepo.FilterUsers(&userDto.FilterUsers{
		LastName: "castarica",
	})
	if err != nil {
		t.Errorf("3. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 1 {
		t.Errorf("3. something went wrong.")
	}

	users, err = userRepo.FilterUsers(&userDto.FilterUsers{
		FirstName: "andy",
	})
	if err != nil {
		t.Errorf("4. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 1 {
		t.Errorf("4. something went wrong.")
	}
}

func TestDeleteByUsername(t *testing.T) {
	userRepo.Create(user1)

	err := userRepo.DeleteByUsername(user1.Username)
	if err != nil {
		t.Error("the user must be delted")
	}
}

func TestUpsert(t *testing.T) {
	userRepo.Create(user1)
	defer userRepo.DeleteAll()

	userCopy := user1
	userCopy.Password = "765432"
	userCopy.LastName = "jordanhenderson"
	userCopy.FirstName = "arthor"
	userCopy.IsAdmin = false

	userCreated, err := userRepo.Upsert(userCopy)
	if err != nil {
		t.Fatalf("upserting went wrong the error was %v", err)
	}
	if userCreated.LastName != userCopy.LastName ||
		userCreated.FirstName != userCopy.FirstName ||
		userCreated.Password != userCopy.Password ||
		userCreated.IsAdmin != userCopy.IsAdmin {
		t.Fatal("the data is returned is wrong")
	}

	userFound, _ := userRepo.FindUserByUsername(user1.Username)

	if userFound.LastName != userCopy.LastName ||
		userFound.FirstName != userCopy.FirstName ||
		userFound.Password != userFound.Password ||
		userFound.IsAdmin != userFound.IsAdmin {
		t.Fatal("the data is returned is wrong")
	}

	_, err = userRepo.Upsert(user1)
	if err != nil {
		t.Fatalf("upserting went wrong the error was %v", err)
	}
}

func TestDeletePreviousSuperAdmins(t *testing.T) {
	defer userRepo.DeleteAll()
	userRepo.Create(user1)
	userRepo.Create(user2)
	userRepo.Create(user3)

	err := userRepo.DeletePreviousSuperAdmins()
	if err != nil {
		t.Fatalf("something went wrong to delete previousSuperadmins")
	}
	superAdmin, _ := userRepo.FindUserByUsername(user1.Username)
	if superAdmin != nil {
		t.Fatal("super admine must be deleted")
	}
}
