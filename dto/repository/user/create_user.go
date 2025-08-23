package dto

import "github.com/mohamadrezamomeni/momo/entity"

type Create struct {
	Username         string
	FirstName        string
	LastName         string
	IsAdmin          bool
	Password         string
	IsSuperAdmin     bool
	TelegramID       string
	IsApproved       bool
	TelegramUsername string
	Language         entity.Language
}
