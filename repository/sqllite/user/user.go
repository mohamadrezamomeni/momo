package user

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"momo/entity"
	momoError "momo/pkg/error"
	"momo/repository/sqllite/user/dto"
)

func (u *User) Create(inpt *dto.Create) (*entity.User, error) {
	var id, username, lastName, firstName string
	err := u.db.Conn().QueryRow(`
	INSERT INTO users (username, lastName, firstName)
	VALUES (?, ?, ?)
	RETURNING id, username, lastName, firstName
`, inpt.Username, inpt.LastName, inpt.FirstName).Scan(&id, &username, &lastName, &firstName)
	if err != nil {
		return &entity.User{}, momoError.Errorf("somoething went wrong to save user error: %v", err)
	}

	return &entity.User{
		ID:        id,
		Username:  username,
		LastName:  lastName,
		FirstName: firstName,
	}, nil
}

func (u *User) Delete(id string) error {
	sql := fmt.Sprintf("DELETE FROM users WHERE id='%s'", id)
	res, err := u.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Errorf("something went wrong to delete record follow error, the error was %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Errorf("something went wrong to delete record follow error, the error was %v", err)
	}

	if rowsAffected == 0 {
		return momoError.Error("None of the records have been affected.")
	}
	return nil
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
		var username string
		var createdAt interface{}
		err = rows.Scan(&id, &username, &createdAt, &lastName, &firstName)
		if err != nil {
			momoError.DebuggingErrorf("error has occured err: %v", err)
		}
		users = append(users, &entity.User{ID: id, Username: username, FirstName: firstName, LastName: lastName})
	}
	return users, nil
}

func (u *User) generateFilterUserQuery(q *dto.FilterUsers) (string, error) {
	query := "SELECT * FROM `users`"

	v := reflect.ValueOf(*q)
	t := reflect.TypeOf(*q)

	conditions := []string{}
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

		if v != "" {
			conditions = append(conditions, fmt.Sprintf(" %s LIKE '%%%s%%'", strings.ToLower(field.Name), v))
		}
	}
	if len(conditions) != 0 {
		joinedConditions := strings.Join(conditions, " AND ")

		query += " WHERE "
		query += joinedConditions

	}
	return query, nil
}

func (u *User) FindUserByUsername(username string) (*entity.User, error) {
	return u.findUser("username", username)
}

func (u *User) FindUserByID(ID string) (*entity.User, error) {
	return u.findUser("id", ID)
}

func (u *User) findUser(key string, value string) (*entity.User, error) {
	var id string
	var firstName string
	var lastName string
	var username string
	var createdAt interface{}
	s := fmt.Sprintf("SELECT * FROM users WHERE %s='%s' LIMIT 1", key, value)
	err := u.db.Conn().QueryRow(s).Scan(&id, &username, &createdAt, &lastName, &firstName)
	if err == nil {
		return &entity.User{
			ID:        id,
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
		}, err
	}
	if err == sql.ErrNoRows {
		return &entity.User{}, err
	}
	return &entity.User{}, momoError.Errorf("some thing went wrong please follow the problem - error: %v", err)
}
