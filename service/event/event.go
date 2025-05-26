package event

import (
	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
	eventServiceDto "github.com/mohamadrezamomeni/momo/dto/service/event"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Event struct {
	eventRepo EventRepository
}

type EventRepository interface {
	Create(*eventRepositoryDto.CreateEvent) (*entity.Event, error)
	Filter(*eventRepositoryDto.FilterEvents) ([]*entity.Event, error)
	Update(string, *eventRepositoryDto.UpdateEvent) error
}

func New(eventRepo EventRepository) *Event {
	return &Event{
		eventRepo: eventRepo,
	}
}

func (e *Event) Create(event *eventServiceDto.CreateEventDto) {
	e.eventRepo.Create(&eventRepositoryDto.CreateEvent{
		Data: event.Data,
		Name: event.Name,
	})
}
