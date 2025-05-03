package user

import "github.com/mohamadrezamomeni/momo/entity"

type UserRepo interface {
	FindUserByID(string) (*entity.User, error)
}

type User struct {
	userRepo UserRepo
}

func New(userRepo UserRepo) *User {
	return &User{
		userRepo: userRepo,
	}
}

func (u *User) FindByID(id string) (*entity.User, error) {
	return u.FindByID(id)
}
