package user

import (
	userDto "github.com/mohamadrezamomeni/momo/dto/repository/user"
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
	}

	user2 = &userDto.Create{
		Username:         "micheal",
		FirstName:        "micheal",
		LastName:         "casta",
		IsAdmin:          true,
		Password:         "12334",
		IsSuperAdmin:     false,
		TelegramUsername: "momo2",
	}

	user3 = &userDto.Create{
		Username:         "arodoa",
		FirstName:        "micka",
		LastName:         "castarica",
		IsAdmin:          true,
		Password:         "1244",
		IsSuperAdmin:     false,
		TelegramUsername: "momo3",
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
	}
)
