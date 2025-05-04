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
		IsAdmin:   userDto.IsAdmin,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Username:  userDto.Username,
		Password:  passwordHashed,
	})
}

func (u *User) FindByID(id string) (*entity.User, error) {
	return u.FindByID(id)
}

func (u *User) FindByUsername(username string) (*entity.User, error) {
	return u.FindByUsername(username)
}
