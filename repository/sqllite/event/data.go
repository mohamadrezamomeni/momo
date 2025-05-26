package event

import (
	eventRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/event"
)

var (
	data1 = &eventRepositoryDto.CreateEvent{
		Name:                    "notification",
		Data:                    `{"type": "telegram"}`,
		IsNotificationProcessed: true,
	}
	data2 = &eventRepositoryDto.CreateEvent{
		Name:                    "test",
		Data:                    `{"type": "telegram"}`,
		IsNotificationProcessed: false,
	}

	data3 = &eventRepositoryDto.CreateEvent{
		Name:                    "test1",
		Data:                    `{"type": "telegram"}`,
		IsNotificationProcessed: false,
	}
)
