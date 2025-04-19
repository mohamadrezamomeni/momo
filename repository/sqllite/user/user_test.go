package user

import (
	"fmt"
	"os"
	"testing"

	"momo/entity"
	"momo/pkg/config"
	"momo/pkg/utils"

	"momo/repository/migrate"
	"momo/repository/sqllite"
	"momo/repository/sqllite/user/dto"
)

var userRepo IUserRepository

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

func getInfo() (string, string, string) {
	username := fmt.Sprintf("%s", utils.RandomString(5))
	name := utils.RandomString(5)
	family := utils.RandomString(5)
	return username, name, family
}

func TestCreate(t *testing.T) {
	username, firstName, lastName := getInfo()

	userCreated, err := userRepo.Create(&dto.Create{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if userCreated.Username != username || userCreated.FirstName != firstName || userCreated.LastName != lastName {
		t.Error("user creation requires some troubleshooting")
		return
	}
	userRepo.Delete(userCreated.ID)
}

func TestFindByUsername(t *testing.T) {
	username, firstName, lastName := getInfo()
	userCreated, err := userRepo.Create(&dto.Create{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	})

	user, err := userRepo.FindUserByUsername(username)
	if err != nil {
		t.Errorf("findByUsername needs troubleshooting error: %v", err)
		userRepo.Delete(userCreated.ID)
		return
	}
	if user.Username != username || user.FirstName != firstName || user.LastName != lastName || user.ID != userCreated.ID {
		t.Error("something went wrong to compare results")
		userRepo.Delete(userCreated.ID)
		return
	}
	userRepo.Delete(userCreated.ID)
}

func TestFindByID(t *testing.T) {
	username, firstName, lastName := getInfo()
	userCreated, err := userRepo.Create(&dto.Create{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	})
	user, err := userRepo.FindUserByID(userCreated.ID)
	if err != nil {
		t.Errorf("findByUsername needs troubleshooting error: %v", err)
		userRepo.Delete(userCreated.ID)
		return
	}

	if user.Username != username || user.FirstName != firstName || user.LastName != lastName || user.ID != userCreated.ID {
		t.Error("something went wrong to compare results")
		return
	}
	userRepo.Delete(userCreated.ID)
}

func TestFilter(t *testing.T) {
	user1, user2, user3 := staticUsers()
	ids := []string{user1.ID, user2.ID, user3.ID}

	users, err := userRepo.FilterUsers(&dto.FilterUsers{
		FirstName: "mic",
	})
	if err != nil {
		t.Errorf("1. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 2 {
		t.Errorf("1. something went wrong.")
	}

	users, err = userRepo.FilterUsers(&dto.FilterUsers{
		FirstName: "mic",
		LastName:  "casta",
	})
	if err != nil {
		t.Errorf("2. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 2 {
		t.Errorf("2. something went wrong.")
	}

	users, err = userRepo.FilterUsers(&dto.FilterUsers{
		LastName: "castarica",
	})
	if err != nil {
		t.Errorf("3. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 1 {
		t.Errorf("3. something went wrong.")
	}

	users, err = userRepo.FilterUsers(&dto.FilterUsers{
		FirstName: "andy",
	})
	if err != nil {
		t.Errorf("4. something went wrong. please follow problem the error was %v", err)
	}

	if len(users) != 1 {
		t.Errorf("4. something went wrong.")
	}
	deleteUsers(ids)
}

func staticUsers() (*entity.User, *entity.User, *entity.User) {
	username1 := "andy"
	username2 := "micheal"
	username3 := "micka"

	firstName1 := "andy"
	firstName2 := "micheal"
	firstName3 := "micka"

	lastName1 := "arodoa"
	lastName2 := "casta"
	lastName3 := "castarica"

	user1, _ := userRepo.Create(&dto.Create{
		Username:  username1,
		FirstName: firstName1,
		LastName:  lastName1,
	})
	user2, _ := userRepo.Create(&dto.Create{
		Username:  username2,
		FirstName: firstName2,
		LastName:  lastName2,
	})
	user3, _ := userRepo.Create(&dto.Create{
		Username:  username3,
		FirstName: firstName3,
		LastName:  lastName3,
	})
	return user1, user2, user3
}

func deleteUsers(ids []string) {
	for _, id := range ids {
		userRepo.Delete(id)
	}
}
