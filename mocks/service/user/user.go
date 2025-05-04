package user

import (
	"github.com/google/uuid"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
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

func (u *MockUser) FindByUsername(username string) (*entity.User, error) {
	return &entity.User{
		ID:        uuid.New().String(),
		Username:  username,
		FirstName: utils.RandomString(5),
		LastName:  utils.RandomString(5),
	}, nil
}
