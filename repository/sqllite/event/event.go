package event

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (e *Event) Create(inpt *eventRepositoryDto.CreateEvent) (*entity.Event, error) {
	scope := "eventRepository.create"

	event := &entity.Event{}
	err := e.db.Conn().QueryRow(`
	INSERT INTO events (name, data, is_notification_processed)
	VALUES (?, ?, ?)
	RETURNING id, name, data, is_notification_processed
	`, inpt.Name,
		inpt.Data,
		inpt.IsNotificationProcessed,
	).Scan(
		&event.ID,
		&event.Name,
		&event.Data,
		&event.IsNotificationProcessed,
	)

	if err == nil {
		return event, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
}

func (e *Event) Filter(inpt *eventRepositoryDto.FilterEvents) ([]*entity.Event, error) {
	scope := "eventRepository.filter"

	query := e.makeQueryFilter(inpt)
	rows, err := e.db.Conn().Query(query)
	if err != nil {
		return nil, momoError.Wrap(err).Input(inpt).Scope(scope).DebuggingError()
	}

	events := make([]*entity.Event, 0)

	for rows.Next() {
		event, err := e.scan(rows)
		if err != nil {
			return []*entity.Event{}, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (e *Event) makeQueryFilter(inpt *eventRepositoryDto.FilterEvents) string {
	v := reflect.ValueOf(*inpt)
	t := reflect.TypeOf(*inpt)

	subQueries := []string{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Name == "Name" && len(value.String()) > 0 {
			subQueries = append(subQueries, fmt.Sprintf("name = '%s'", value.String()))
		} else if field.Name == "IsNotificationProcessed" && !value.IsNil() {
			subQueries = append(subQueries, fmt.Sprintf("is_notification_processed = %v", value.Elem().Bool()))
		} else if field.Name == "Names" && (value.Kind() == reflect.Slice || value.Kind() == reflect.Array) && value.Len() > 0 {
			subQueries = append(subQueries, fmt.Sprintf("name IN (\"%s\")", strings.Join(inpt.Names, "\", \"")))
		}
	}

	sql := "SELECT * FROM events"
	if len(subQueries) > 0 {
		sql += fmt.Sprintf(" WHERE %s", strings.Join(subQueries, " AND "))
	}

	return sql
}

func (e *Event) scan(rows *sql.Rows) (*entity.Event, error) {
	scope := "eventRepository.scan"

	event := &entity.Event{}
	var createdAt interface{}
	err := rows.Scan(
		&event.ID,
		&event.Name,
		&event.Data,
		&event.IsNotificationProcessed,
		&createdAt,
	)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(rows).DebuggingError()
	}
	return event, nil
}

func (e *Event) DeleteAll() error {
	scope := "eventRepository.DeleteAll"

	sql := "DELETE FROM events"
	res, err := e.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return momoError.Wrap(err).Scope(scope).UnExpected().DebuggingError()
	}

	return nil
}

func (e *Event) Update(id string, inpt *eventRepositoryDto.UpdateEvent) error {
	scope := "eventRepository.Update"
	subModifies := []string{}

	if inpt.IsNotificationProcessed != nil {
		subModifies = append(subModifies, fmt.Sprintf("is_notification_processed = %v", *inpt.IsNotificationProcessed))
	}
	if len(subModifies) == 0 {
		return momoError.Scope(scope).DebuggingErrorf("input was empty")
	}
	sql := fmt.Sprintf(
		"UPDATE events SET %s WHERE id = %v",
		strings.Join(subModifies, ", "),
		id,
	)

	_, err := e.db.Conn().Exec(sql)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(id, inpt).DebuggingError()
	}
	return nil
}
