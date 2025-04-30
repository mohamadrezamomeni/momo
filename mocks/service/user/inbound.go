package user

import (
	"momo/entity"
	"momo/pkg/utils"
)

type MockUser struct{}

func New() *MockUser {
	return &MockUser{}
}

func (u *MockUser) FindByID(id string) (*entity.User, error) {
	return &entity.User{
		ID:        id,
		Username:  utils.RandomString(5),
		FirstName: utils.RandomString(5),
		LastName:  utils.RandomString(5),
	}, nil
}
