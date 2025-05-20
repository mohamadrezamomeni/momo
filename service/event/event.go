package event

import (
	"sync"

	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
	eventServiceDto "github.com/mohamadrezamomeni/momo/dto/service/event"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Event struct {
	eventRepo EventRepository
	ch        chan *eventServiceDto.CreateEventDto
	wg        sync.WaitGroup
	once      sync.Once
}

type EventRepository interface {
	Create(*eventRepositoryDto.CreateEvent) (*entity.Event, error)
	Filter(*eventRepositoryDto.FilterEvents) ([]*entity.Event, error)
	Update(string, *eventRepositoryDto.UpdateEvent) error
}

func New(eventRepo EventRepository) *Event {
	return &Event{
		eventRepo: eventRepo,
		ch:        make(chan *eventServiceDto.CreateEventDto),
	}
}

func (e *Event) Create(event *eventServiceDto.CreateEventDto) {
	e.ch <- event
}

func (e *Event) ListEvents() {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		for event := range e.ch {
			e.eventRepo.Create(&eventRepositoryDto.CreateEvent{
				Data: event.Data,
				Name: event.Name,
			})
		}
	}()
}

func (e *Event) Shutdown() {
	e.once.Do(func() {
		close(e.ch)
		e.wg.Wait()
	})
}
