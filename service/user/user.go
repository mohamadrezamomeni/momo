package user

import (
	userRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/user"
	userServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user"
	"github.com/mohamadrezamomeni/momo/entity"
	crypt "github.com/mohamadrezamomeni/momo/service/crypt"
)

type UserRepo interface {
	FindUserByID(string) (*entity.User, error)
	FindUserByUsername(string) (*entity.User, error)
	Create(*userRepoDto.Create) (*entity.User, error)
	Upsert(*userRepoDto.Create) (*entity.User, error)
	DeleteByUsername(string) error
	DeletePreviousSuperAdmins() error
	FilterUsers(*userRepoDto.FilterUsers) ([]*entity.User, error)
	FindByTelegramID(string) (*entity.User, error)
}

type User struct {
	userRepo UserRepo
	crypt    *crypt.Crypt
}

func New(userRepo UserRepo, crypt *crypt.Crypt) *User {
	return &User{
		userRepo: userRepo,
		crypt:    crypt,
	}
}

func (u *User) Create(userDto *userServiceDto.AddUser) (*entity.User, error) {
	passwordHashed, err := u.crypt.Hash(userDto.Password)
	if err != nil {
		return nil, err
	}

	return u.userRepo.Create(&userRepoDto.Create{
		IsAdmin:    userDto.IsAdmin,
		TelegramID: userDto.TelegramID,
		FirstName:  userDto.FirstName,
		LastName:   userDto.LastName,
		Username:   userDto.Username,
		Password:   passwordHashed,
	})
}

func (u *User) CreateUserAdmin(userDto *userServiceDto.AddUser) (*entity.User, error) {
	err := u.userRepo.DeletePreviousSuperAdmins()
	if err != nil {
		return nil, err
	}

	passwordHashed, err := u.crypt.Hash(userDto.Password)
	if err != nil {
		return nil, err
	}

	return u.userRepo.Upsert(&userRepoDto.Create{
		IsAdmin:   userDto.IsAdmin,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Username:  userDto.Username,
		Password:  passwordHashed,
	})
}

func (u *User) FindByID(id string) (*entity.User, error) {
	return u.userRepo.FindUserByID(id)
}

func (u *User) FindByUsername(username string) (*entity.User, error) {
	return u.userRepo.FindUserByUsername(username)
}

func (u *User) DeleteByUsername(username string) error {
	return u.userRepo.DeleteByUsername(username)
}

func (u *User) Filter() ([]*entity.User, error) {
	return u.userRepo.FilterUsers(&userRepoDto.FilterUsers{})
}

func (u *User) FindByTelegramID(tid string) (*entity.User, error) {
	return u.userRepo.FindByTelegramID(tid)
}
