package user

import (
	"momo/pkg/entity"
	momoError "momo/pkg/error"
	"momo/repository/sqllite/user/dto"
)

func (u *User) Create(inpt dto.Create) (*entity.User, error) {
	var id, email, lastName, firstName string
	err := u.db.Conn().QueryRow(`
	INSERT INTO users (email, lastName, firstName)
	VALUES (?, ?, ?)
	RETURNING id, email, lastName, firstName
`, inpt.Email, inpt.LastName, inpt.FirstName).Scan(&id, &email, &lastName, &firstName)
	if err != nil {
		return &entity.User{}, momoError.Errorf("somoething went wrong to save user")
	}

	return &entity.User{
		Email:     email,
		LastName:  lastName,
		FirstName: firstName,
	}, nil
}
