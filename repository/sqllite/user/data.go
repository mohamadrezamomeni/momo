package user

import (
	userDto "github.com/mohamadrezamomeni/momo/dto/repository/user"
	"github.com/mohamadrezamomeni/momo/entity"
)

var (
	user1 = &userDto.Create{
		Username:         "andy",
		FirstName:        "andy",
		LastName:         "arodoa",
		IsAdmin:          true,
		IsApproved:       true,
		Password:         "12342",
		IsSuperAdmin:     true,
		TelegramUsername: "momo1",
		Language:         entity.EN,
	}

	user2 = &userDto.Create{
		Username:         "micheal",
		FirstName:        "micheal",
		LastName:         "casta",
		IsAdmin:          true,
		Password:         "12334",
		IsSuperAdmin:     false,
		TelegramUsername: "momo2",
		Language:         entity.EN,
	}

	user3 = &userDto.Create{
		Username:         "arodoa",
		FirstName:        "micka",
		LastName:         "castarica",
		IsAdmin:          true,
		Password:         "1244",
		IsSuperAdmin:     false,
		TelegramUsername: "momo3",
		Language:         entity.EN,
	}

	user4 = &userDto.Create{
		Username:         "madona",
		FirstName:        "micka",
		LastName:         "castarica",
		IsAdmin:          true,
		Password:         "1244",
		IsSuperAdmin:     false,
		IsApproved:       false,
		TelegramUsername: "momo4",
		Language:         entity.EN,
	}
)
