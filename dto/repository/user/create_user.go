package dto

type Create struct {
	Username     string
	FirstName    string
	LastName     string
	IsAdmin      bool
	Password     string
	IsSuperAdmin bool
	TelegramID   string
	IsApproved   bool
}
