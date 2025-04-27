package dto

import (
	"time"

	"momo/entity"
)

type CreateInbound struct {
	Protocol   string
	Tag        string
	Port       string
	UserID     string
	Domain     string
	VPNType    entity.VPNType
	IsAssigned bool
	IsNotified bool
	IsActive   bool
	Start      time.Time
	End        time.Time
	IsBlock    bool
}
