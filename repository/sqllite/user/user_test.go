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
	email := fmt.Sprintf("%s@gmail.com", utils.RandomString(5))
	name := utils.RandomString(5)
	family := utils.RandomString(5)
	return email, name, family
}

func TestCreate(t *testing.T) {
	email, firstName, lastName := getInfo()

	userCreated, err := userRepo.Create(&dto.Create{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if userCreated.Email != email || userCreated.FirstName != firstName || userCreated.LastName != lastName {
		t.Error("user creation requires some troubleshooting")
		return
	}
	userRepo.Delete(userCreated.ID)
}

func TestFindByEmail(t *testing.T) {
	email, firstName, lastName := getInfo()
	userCreated, err := userRepo.Create(&dto.Create{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	})

	user, err := userRepo.FindUserByEmail(email)
	if err != nil {
		t.Errorf("findByEmail needs troubleshooting error: %v", err)
		userRepo.Delete(userCreated.ID)
		return
	}
	if user.Email != email || user.FirstName != firstName || user.LastName != lastName || user.ID != userCreated.ID {
		t.Error("something went wrong to compare results")
		userRepo.Delete(userCreated.ID)
		return
	}
	userRepo.Delete(userCreated.ID)
}

func TestFindByID(t *testing.T) {
	email, firstName, lastName := getInfo()
	userCreated, err := userRepo.Create(&dto.Create{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	})
	user, err := userRepo.FindUserByID(userCreated.ID)
	if err != nil {
		t.Errorf("findByEmail needs troubleshooting error: %v", err)
		userRepo.Delete(userCreated.ID)
		return
	}

	if user.Email != email || user.FirstName != firstName || user.LastName != lastName || user.ID != userCreated.ID {
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
		deleteUsers(ids)
	}

	if len(users) != 2 {
		t.Errorf("1. something went wrong.")
		deleteUsers(ids)
	}

	users, err = userRepo.FilterUsers(&dto.FilterUsers{
		FirstName: "mic",
		LastName:  "casta",
	})
	if err != nil {
		t.Errorf("2. something went wrong. please follow problem the error was %v", err)
		deleteUsers(ids)
	}

	if len(users) != 2 {
		t.Errorf("2. something went wrong.")
		deleteUsers(ids)
	}

	users, err = userRepo.FilterUsers(&dto.FilterUsers{
		LastName: "castarica",
	})
	if err != nil {
		t.Errorf("3. something went wrong. please follow problem the error was %v", err)
		deleteUsers(ids)
	}

	if len(users) != 1 {
		t.Errorf("3. something went wrong.")
		deleteUsers(ids)
	}

	users, err = userRepo.FilterUsers(&dto.FilterUsers{
		FirstName: "andy",
	})
	if err != nil {
		t.Errorf("4. something went wrong. please follow problem the error was %v", err)
		deleteUsers(ids)
	}

	if len(users) != 1 {
		t.Errorf("4. something went wrong.")
		deleteUsers(ids)
	}
}

func staticUsers() (*entity.User, *entity.User, *entity.User) {
	email1 := "andy@gmail.com"
	email2 := "micheal@gmail.com"
	email3 := "micka@gmail.com"

	firstName1 := "andy"
	firstName2 := "micheal"
	firstName3 := "micka"

	lastName1 := "arodoa"
	lastName2 := "casta"
	lastName3 := "castarica"

	user1, _ := userRepo.Create(&dto.Create{
		Email:     email1,
		FirstName: firstName1,
		LastName:  lastName1,
	})
	user2, _ := userRepo.Create(&dto.Create{
		Email:     email2,
		FirstName: firstName2,
		LastName:  lastName2,
	})
	user3, _ := userRepo.Create(&dto.Create{
		Email:     email3,
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
