package user

import (
	"momo/pkg/entity"
	"momo/repository/sqllite"
	"momo/repository/sqllite/user/dto"
)

type IUserRepository interface {
	Create(dto.Create) (*entity.User, error)
}

type User struct {
	db *sqllite.SqlliteDB
}

func New(sqlite *sqllite.SqlliteDB) IUserRepository {
	return &User{
		db: sqlite,
	}
}
