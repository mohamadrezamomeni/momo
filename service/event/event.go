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

func (e *Event) MarkNotificationProcessed(id string) error {
	active := true
	return e.eventRepo.Update(id, &eventRepositoryDto.UpdateEvent{
		IsNotificationProcessed: &active,
	})
}

func (e *Event) FilterNotifications(filterDto *eventServiceDto.FilterEvents) ([]*entity.Event, error) {
	return e.eventRepo.Filter(&eventRepositoryDto.FilterEvents{
		Name:                    filterDto.Name,
		IsNotificationProcessed: filterDto.IsNotificationProcessed,
	})
}
