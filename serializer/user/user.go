package user

type UserSerialize struct {
	Username     string `json:"username"`
	Lastname     string `json:"lastname"`
	FirstName    string `json:"firstname"`
	IsAdmin      bool   `json:"is_admin"`
	IsSuperAdmin bool   `json:"is_super_admin"`
}
