package user

import (
	dto "momo/dto/repository/user"
	"momo/entity"
	"momo/repository/sqllite"
)

type IUserRepository interface {
	Delete(string) error
	Create(*dto.Create) (*entity.User, error)
	FilterUsers(q *dto.FilterUsers) ([]*entity.User, error)
	FindUserByUsername(string) (*entity.User, error)
	FindUserByID(string) (*entity.User, error)
}

type User struct {
	db *sqllite.SqlliteDB
}

func New(sqlite *sqllite.SqlliteDB) IUserRepository {
	return &User{
		db: sqlite,
	}
}
