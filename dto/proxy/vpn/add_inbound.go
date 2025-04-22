package dto

type Inbound struct {
	User     *User
	Port     string
	Protocol string
	Address  string
	Tag      string
}
