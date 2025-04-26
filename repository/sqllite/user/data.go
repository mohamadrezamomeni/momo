package user

import userDto "momo/dto/repository/user"

var (
	user1 = &userDto.Create{
		Username:  "andy",
		FirstName: "andy",
		LastName:  "arodoa",
	}

	user2 = &userDto.Create{
		Username:  "micheal",
		FirstName: "micheal",
		LastName:  "casta",
	}

	user3 = &userDto.Create{
		Username:  "arodoa",
		FirstName: "micka",
		LastName:  "castarica",
	}
)
