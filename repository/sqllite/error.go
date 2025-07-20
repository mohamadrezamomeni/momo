package sqllite

import "github.com/mattn/go-sqlite3"

func IsDuplicateError(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		return err.Code == sqlite3.ErrConstraint && err.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}
