package user

type UserSerialize struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Lastname         string `json:"lastname"`
	FirstName        string `json:"firstname"`
	IsAdmin          bool   `json:"is_admin"`
	IsSuperAdmin     bool   `json:"is_super_admin"`
	IsApproved       bool   `json:"is_approved"`
	TelegramUsername string `json:"telegram_username"`
}
