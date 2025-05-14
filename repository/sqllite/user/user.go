package user

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	dto "github.com/mohamadrezamomeni/momo/dto/repository/user"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (u *User) Create(inpt *dto.Create) (*entity.User, error) {
	scope := "userRepository.Create"

	user := &entity.User{}
	err := u.db.Conn().QueryRow(`
	INSERT INTO users (username, lastName, firstName, password, is_admin, is_super_admin, telegram_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING id, username, lastName, firstName, is_admin, password, is_super_admin, telegram_id
`,
		inpt.Username,
		inpt.LastName,
		inpt.FirstName,
		inpt.Password,
		inpt.IsAdmin,
		inpt.IsSuperAdmin,
		inpt.TelegramID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.LastName,
		&user.FirstName,
		&user.IsAdmin,
		&user.Password,
		&user.IsSuperAdmin,
		&user.TelegramID,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(inpt).DebuggingError()
	}

	return user, nil
}

func (u *User) Upsert(inpt *dto.Create) (*entity.User, error) {
	scope := "userRepository.Upsert"

	user := &entity.User{}
	err := u.db.Conn().QueryRow(`
	INSERT INTO users (username, lastName, firstName, password, is_admin, is_super_admin, telegram_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(username) DO UPDATE SET
		password = excluded.password,
		firstname = excluded.firstname,
		lastname = excluded.lastname,
		is_admin = excluded.is_admin
	RETURNING id, username, lastName, firstName, is_admin, password, is_super_admin, telegram_id
`, inpt.Username,
		inpt.LastName,
		inpt.FirstName,
		inpt.Password,
		inpt.IsAdmin,
		inpt.IsSuperAdmin,
		inpt.TelegramID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.LastName,
		&user.FirstName,
		&user.IsAdmin,
		&user.Password,
		&user.IsSuperAdmin,
		&user.TelegramID,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(inpt).DebuggingError()
	}

	return user, nil
}

func (u *User) Delete(id string) error {
	scope := "userRepository.Delete"

	sql := fmt.Sprintf("DELETE FROM users WHERE id='%s'", id)
	res, err := u.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}

	if rowsAffected == 0 {
		return momoError.Wrap(err).Scope(scope).Input(id).DebuggingError()
	}
	return nil
}

func (u *User) DeleteByUsername(username string) error {
	scope := "userRepository.DeleteByUsername"
	sql := fmt.Sprintf("DELETE FROM users WHERE username='%s'", username)
	res, err := u.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(username).UnExpected().DebuggingError()
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(username).UnExpected().DebuggingError()
	}

	if rowsAffected == 0 {
		return momoError.Wrap(err).Scope(scope).Input(username).UnExpected().DebuggingError()
	}
	return nil
}

func (u *User) DeleteAll() error {
	scope := "userRepository.DeleteAll"

	sql := "DELETE FROM users"
	res, err := u.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}

func (u *User) FilterUsers(q *dto.FilterUsers) ([]*entity.User, error) {
	scope := "userRepository.FilterUsers"

	query, err := u.generateFilterUserQuery(q)
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(q).UnExpected().DebuggingError()
	}
	users := []*entity.User{}
	for rows.Next() {
		user := &entity.User{}

		var createdAt interface{}
		err = rows.Scan(&user.ID, &user.Username, &createdAt, &user.LastName, &user.FirstName, &user.Password, &user.IsAdmin, &user.IsSuperAdmin, &user.TelegramID)
		if err != nil {
			return nil, momoError.Wrap(err).Scope(scope).Input(q).UnExpected().DebuggingError()
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) generateFilterUserQuery(q *dto.FilterUsers) (string, error) {
	scope := "userRepository.generateFilterUserQuery"

	query := "SELECT * FROM `users`"

	v := reflect.ValueOf(*q)
	t := reflect.TypeOf(*q)

	conditions := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if k := value.Kind(); k != reflect.String {
			return "", momoError.Scope(scope).Input(q).UnExpected().DebuggingError()
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
	scope := "userRepository.findUser"

	var user *entity.User = &entity.User{}

	var createdAt interface{}
	s := fmt.Sprintf("SELECT * FROM users WHERE %s='%s' LIMIT 1", key, value)
	err := u.db.Conn().QueryRow(s).Scan(
		&user.ID,
		&user.Username,
		&createdAt,
		&user.LastName,
		&user.FirstName,
		&user.Password,
		&user.IsAdmin,
		&user.IsSuperAdmin,
		&user.TelegramID,
	)
	if err == nil {
		return user, nil
	}
	if err == sql.ErrNoRows {
		return nil, momoError.Wrap(err).Scope(scope).Input(key, value).NotFound().DebuggingError()
	}
	return nil, momoError.Wrap(err).Scope(scope).Input(key, value).UnExpected().DebuggingError()
}

func (u *User) DeletePreviousSuperAdmins() error {
	scope := "userRepository.DeletePriviousSuperAdmins"

	res, err := u.db.Conn().Exec("DELETE FROM users WHERE is_super_admin=true")
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return nil
}
