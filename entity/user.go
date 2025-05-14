package entity

type User struct {
	ID           string
	Username     string
	LastName     string
	FirstName    string
	IsAdmin      bool
	Password     string
	IsSuperAdmin bool
	TelegramID   string
}
