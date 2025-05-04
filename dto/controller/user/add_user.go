package user

type AddUser struct {
	IsAdmin   bool   `json:"is_admin"`
	Password  string `json:"password"`
	LastName  string `json:"lastname"`
	FirstName string `json:"firstname"`
	Username  string `json:"username"`
}
