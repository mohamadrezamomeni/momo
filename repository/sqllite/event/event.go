package event

import (
	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	errorRepository "github.com/mohamadrezamomeni/momo/repository/sqllite"
)

func (e *Event) Create(inpt *eventRepositoryDto.CreateEvent) (*entity.Event, error) {
	scope := "inboundRepository.create"

	event := &entity.Event{}
	err := e.db.Conn().QueryRow(`
	INSERT INTO events (name, data, is_processed)
	VALUES (?, ?, ?)
	RETURNING id, name, data, is_processed
	`, inpt.Name,
		inpt.Data,
		false,
	).Scan(
		&event.ID,
		&event.Name,
		&event.Data,
		&event.IsProccessed,
	)

	if err == nil {
		return event, nil
	}

	if errorRepository.IsDuplicateError(err) {
		return nil, momoError.Wrap(err).Input(inpt).Duplicate().Scope(scope).DebuggingError()
	}
	return nil, momoError.Wrap(err).Input(inpt).UnExpected().Scope(scope).DebuggingError()
}
