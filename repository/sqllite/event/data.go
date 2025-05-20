package event

import (
	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
)

var data1 = &eventRepositoryDto.CreateEvent{
	Name: "notification",
	Data: `{"type": "telegram"}`,
}
