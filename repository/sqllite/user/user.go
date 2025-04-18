package user

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"momo/pkg/entity"
	momoError "momo/pkg/error"
	"momo/repository/sqllite/user/dto"
)

func (u *User) Create(inpt *dto.Create) (*entity.User, error) {
	var id, email, lastName, firstName string
	err := u.db.Conn().QueryRow(`
	INSERT INTO users (email, lastName, firstName)
	VALUES (?, ?, ?)
	RETURNING id, email, lastName, firstName
`, inpt.Email, inpt.LastName, inpt.FirstName).Scan(&id, &email, &lastName, &firstName)
	if err != nil {
		return &entity.User{}, momoError.Errorf("somoething went wrong to save user error: %v", err)
	}

	return &entity.User{
		Email:     email,
		LastName:  lastName,
		FirstName: firstName,
	}, nil
}

func (u *User) FilterUsers(q *dto.FilterUsers) ([]*entity.User, error) {
	query, err := u.generateFilterUserQuery(q)
	if err != nil {
		return []*entity.User{}, err
	}

	rows, err := u.db.Conn().Query(query)
	if err != nil {
		return []*entity.User{}, momoError.Errorf("error has occured err: %v", err)
	}
	users := []*entity.User{}
	for rows.Next() {
		var id string
		var firstName string
		var lastName string
		var email string
		var createdAt interface{}
		err = rows.Scan(&id, &email, &createdAt, &lastName, &firstName)
		if err != nil {
			momoError.DebuggingErrorf("error has occured err: %v", err)
		}
		users = append(users, &entity.User{ID: id, Email: email, FirstName: firstName, LastName: lastName})
	}
	return users, nil
}

func (u *User) generateFilterUserQuery(q *dto.FilterUsers) (string, error) {
	query := "SELECT * FROM `users`"

	v := reflect.ValueOf(*q)
	t := reflect.TypeOf(*q)
	isWherePut := false

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if k := value.Kind(); k != reflect.String {
			return "", momoError.DebuggingErrorf(
				"error has occured in filtering Field %-10s | Type: %-8s | Value: %v",
				field.Name,
				value.Kind(),
				value.Interface(),
			)
		}
		v := value.String()
		if !isWherePut && v != "" {
			query += " WHERE "
			isWherePut = true
		}
		if v != "" {
			query += fmt.Sprintf(" %s='%s'", strings.ToLower(field.Name), v)
		}
	}

	return query, nil
}

func (u *User) FindUserByEmail(email string) (*entity.User, error) {
	return u.findUser("email", email)
}

func (u *User) FindUserByID(ID string) (*entity.User, error) {
	return u.findUser("id", ID)
}

func (u *User) findUser(key string, value string) (*entity.User, error) {
	var id string
	var firstName string
	var lastName string
	var email string
	var createdAt interface{}
	s := fmt.Sprintf("SELECT * FROM users WHERE %s='%s' LIMIT 1", key, value)
	err := u.db.Conn().QueryRow(s).Scan(&id, &email, &createdAt, &lastName, &firstName)
	if err == nil {
		return &entity.User{
			ID:        id,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}, err
	}
	if err == sql.ErrNoRows {
		return &entity.User{}, err
	}
	return &entity.User{}, momoError.Errorf("some thing went wrong please follow the problem - error: %v", err)
}
